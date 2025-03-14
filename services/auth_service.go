package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials wordt teruggegeven wanneer de inloggegevens ongeldig zijn
	ErrInvalidCredentials = errors.New("ongeldige inloggegevens")

	// ErrUserInactive wordt teruggegeven wanneer de gebruiker inactief is
	ErrUserInactive = errors.New("gebruiker is inactief")

	// ErrInvalidToken wordt teruggegeven wanneer het token ongeldig is
	ErrInvalidToken = errors.New("ongeldig token")

	// ErrUserNotFound wordt teruggegeven wanneer de gebruiker niet gevonden kan worden
	ErrUserNotFound = errors.New("gebruiker niet gevonden")
)

// JWTClaims definieert de claims in het JWT token
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthServiceImpl implementeert de AuthService interface
type AuthServiceImpl struct {
	gebruikerRepo repository.GebruikerRepository
	jwtSecret     []byte
	tokenExpiry   time.Duration
}

// NewAuthService maakt een nieuwe AuthService
func NewAuthService(gebruikerRepo repository.GebruikerRepository) AuthService {
	// Haal JWT secret uit omgevingsvariabele of gebruik een standaard waarde
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Warn("JWT_SECRET omgevingsvariabele niet gevonden, gebruik standaard waarde")
		jwtSecret = "default_jwt_secret_change_in_production"
	}

	// Haal token expiry uit omgevingsvariabele of gebruik een standaard waarde (24 uur)
	tokenExpiryStr := os.Getenv("JWT_TOKEN_EXPIRY")
	tokenExpiry := 24 * time.Hour
	if tokenExpiryStr != "" {
		var err error
		tokenExpiry, err = time.ParseDuration(tokenExpiryStr)
		if err != nil {
			logger.Warn("Ongeldige JWT_TOKEN_EXPIRY waarde, gebruik standaard waarde", "error", err)
		}
	}

	return &AuthServiceImpl{
		gebruikerRepo: gebruikerRepo,
		jwtSecret:     []byte(jwtSecret),
		tokenExpiry:   tokenExpiry,
	}
}

// Login authenticeert een gebruiker en geeft een JWT token terug
func (s *AuthServiceImpl) Login(ctx context.Context, email, wachtwoord string) (string, error) {
	logger.Info("Login poging", "email", email)

	// Haal gebruiker op basis van email
	gebruiker, err := s.gebruikerRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "email", email, "error", err)
		return "", err
	}

	// Controleer of gebruiker bestaat
	if gebruiker == nil {
		logger.Warn("Gebruiker niet gevonden", "email", email)
		return "", ErrInvalidCredentials
	}

	// Controleer of gebruiker actief is
	if !gebruiker.IsActief {
		logger.Warn("Inactieve gebruiker probeert in te loggen", "email", email)
		return "", ErrUserInactive
	}

	// Verifieer wachtwoord
	if !s.VerifyPassword(gebruiker.WachtwoordHash, wachtwoord) {
		logger.Warn("Ongeldig wachtwoord", "email", email)
		return "", ErrInvalidCredentials
	}

	// Update laatste login
	if err := s.gebruikerRepo.UpdateLastLogin(ctx, gebruiker.ID); err != nil {
		logger.Error("Fout bij updaten laatste login", "email", email, "error", err)
		// We gaan door ondanks de fout, omdat de login zelf succesvol was
	}

	// Genereer JWT token
	token, err := s.generateToken(gebruiker)
	if err != nil {
		logger.Error("Fout bij genereren token", "email", email, "error", err)
		return "", err
	}

	logger.Info("Login succesvol", "email", email, "user_id", gebruiker.ID)
	return token, nil
}

// ValidateToken valideert een JWT token en geeft de gebruiker ID terug
func (s *AuthServiceImpl) ValidateToken(token string) (string, error) {
	// Parse token
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Controleer signing methode
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("onverwachte signing methode: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		logger.Error("Fout bij valideren token", "error", err)
		return "", ErrInvalidToken
	}

	// Controleer of token geldig is
	if !parsedToken.Valid {
		logger.Warn("Ongeldig token")
		return "", ErrInvalidToken
	}

	// Haal claims op
	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		logger.Error("Kon claims niet parsen")
		return "", ErrInvalidToken
	}

	logger.Info("Token gevalideerd", "user_id", claims.UserID)
	return claims.UserID, nil
}

// GetUserFromToken haalt de gebruiker op basis van een JWT token
func (s *AuthServiceImpl) GetUserFromToken(ctx context.Context, token string) (*models.Gebruiker, error) {
	// Valideer token en haal gebruiker ID op
	userID, err := s.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Haal gebruiker op basis van ID
	gebruiker, err := s.gebruikerRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "user_id", userID, "error", err)
		return nil, err
	}

	// Controleer of gebruiker bestaat
	if gebruiker == nil {
		logger.Warn("Gebruiker niet gevonden", "user_id", userID)
		return nil, ErrUserNotFound
	}

	// Controleer of gebruiker actief is
	if !gebruiker.IsActief {
		logger.Warn("Inactieve gebruiker", "user_id", userID)
		return nil, ErrUserInactive
	}

	return gebruiker, nil
}

// HashPassword genereert een hash voor een wachtwoord
func (s *AuthServiceImpl) HashPassword(wachtwoord string) (string, error) {
	// Genereer hash met bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(wachtwoord), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Fout bij hashen wachtwoord", "error", err)
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifieert een wachtwoord tegen een hash
func (s *AuthServiceImpl) VerifyPassword(hash, wachtwoord string) bool {
	// Verifieer wachtwoord met bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(wachtwoord))
	return err == nil
}

// ResetPassword reset het wachtwoord van een gebruiker
func (s *AuthServiceImpl) ResetPassword(ctx context.Context, email, nieuwWachtwoord string) error {
	logger.Info("Wachtwoord reset poging", "email", email)

	// Haal gebruiker op basis van email
	gebruiker, err := s.gebruikerRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "email", email, "error", err)
		return err
	}

	// Controleer of gebruiker bestaat
	if gebruiker == nil {
		logger.Warn("Gebruiker niet gevonden", "email", email)
		return ErrUserNotFound
	}

	// Hash nieuw wachtwoord
	hash, err := s.HashPassword(nieuwWachtwoord)
	if err != nil {
		return err
	}

	// Update wachtwoord
	gebruiker.WachtwoordHash = hash
	if err := s.gebruikerRepo.Update(ctx, gebruiker); err != nil {
		logger.Error("Fout bij updaten wachtwoord", "email", email, "error", err)
		return err
	}

	logger.Info("Wachtwoord reset succesvol", "email", email)
	return nil
}

// generateToken genereert een JWT token voor een gebruiker
func (s *AuthServiceImpl) generateToken(gebruiker *models.Gebruiker) (string, error) {
	// Maak claims
	claims := JWTClaims{
		UserID: gebruiker.ID,
		Email:  gebruiker.Email,
		Role:   gebruiker.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dklemailservice",
			Subject:   gebruiker.ID,
		},
	}

	// Maak token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Onderteken token
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
