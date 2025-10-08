package services

import (
	"context"
	"crypto/rand"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"strings"

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
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// AuthServiceImpl implementeert de AuthService interface
type AuthServiceImpl struct {
	gebruikerRepo    repository.GebruikerRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtSecret        []byte
	tokenExpiry      time.Duration
}

// NewAuthService maakt een nieuwe AuthService
func NewAuthService(gebruikerRepo repository.GebruikerRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	// Haal JWT secret uit omgevingsvariabele of gebruik een standaard waarde
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Warn("JWT_SECRET omgevingsvariabele niet gevonden, gebruik standaard waarde")
		jwtSecret = "default_jwt_secret_change_in_production"
	}

	// Haal token expiry uit omgevingsvariabele of gebruik een standaard waarde (20 minuten)
	tokenExpiryStr := os.Getenv("JWT_TOKEN_EXPIRY")
	tokenExpiry := 20 * time.Minute
	if tokenExpiryStr != "" {
		var err error
		tokenExpiry, err = time.ParseDuration(tokenExpiryStr)
		if err != nil {
			logger.Warn("Ongeldige JWT_TOKEN_EXPIRY waarde, gebruik standaard waarde", "error", err)
		}
	}

	return &AuthServiceImpl{
		gebruikerRepo:    gebruikerRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        []byte(jwtSecret),
		tokenExpiry:      tokenExpiry,
	}
}

// Login authenticeert een gebruiker en geeft een access token en refresh token terug
func (s *AuthServiceImpl) Login(ctx context.Context, email, wachtwoord string) (string, string, error) {
	logger.Info("Login poging", "email", email)

	// Haal gebruiker op basis van email
	gebruiker, err := s.gebruikerRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "email", email, "error", err)
		return "", "", err
	}

	// Controleer of gebruiker bestaat
	if gebruiker == nil {
		logger.Warn("Gebruiker niet gevonden", "email", email)
		return "", "", ErrInvalidCredentials
	}

	// Controleer of gebruiker actief is
	if !gebruiker.IsActief {
		logger.Warn("Inactieve gebruiker probeert in te loggen", "email", email)
		return "", "", ErrUserInactive
	}

	// Verifieer wachtwoord
	if !s.VerifyPassword(gebruiker.WachtwoordHash, wachtwoord) {
		logger.Warn("Ongeldig wachtwoord", "email", email)
		return "", "", ErrInvalidCredentials
	}

	// Update laatste login
	if err := s.gebruikerRepo.UpdateLastLogin(ctx, gebruiker.ID); err != nil {
		logger.Error("Fout bij updaten laatste login", "email", email, "error", err)
		// We gaan door ondanks de fout, omdat de login zelf succesvol was
	}

	// Genereer JWT access token
	accessToken, err := s.generateToken(gebruiker)
	if err != nil {
		logger.Error("Fout bij genereren access token", "email", email, "error", err)
		return "", "", err
	}

	// Genereer refresh token
	refreshToken, err := s.GenerateRefreshToken(ctx, gebruiker.ID)
	if err != nil {
		logger.Error("Fout bij genereren refresh token", "email", email, "error", err)
		return "", "", err
	}

	logger.Info("Login succesvol", "email", email, "user_id", gebruiker.ID)
	return accessToken, refreshToken, nil
}

// ValidateToken valideert een JWT token en geeft de gebruiker ID terug
func (s *AuthServiceImpl) ValidateToken(token string) (string, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	// Parse token
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Controleer signing methode
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("onverwachte signing methode: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		// Verbeterde logging om specifieke JWT validatiefouten te tonen
		if errors.Is(err, jwt.ErrTokenMalformed) {
			logger.Error("Fout bij valideren token: Malformed token", "error", err)
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			logger.Error("Fout bij valideren token: Invalid signature", "error", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			logger.Warn("Fout bij valideren token: Token expired or not valid yet", "error", err)
		} else {
			logger.Error("Fout bij valideren token: Andere fout", "error", err)
		}
		return "", ErrInvalidToken
	}

	// Controleer of token geldig is (ParseWithClaims doet dit al, maar extra check kan geen kwaad)
	if !parsedToken.Valid {
		logger.Warn("Ongeldig token (parsedToken.Valid is false)")
		return "", ErrInvalidToken
	}

	// Haal claims op
	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		logger.Error("Kon claims niet naar *JWTClaims casten")
		return "", ErrInvalidToken
	}

	// Controleer of Subject (user ID) leeg is
	if claims.Subject == "" {
		logger.Error("Token gevalideerd, maar Subject (user ID) claim is leeg")
		return "", ErrInvalidToken // Behandel lege user ID als ongeldig token
	}

	logger.Info("Token gevalideerd", "user_id", claims.Subject) // Gebruik claims.Subject
	return claims.Subject, nil                                  // Geef Subject (user ID) terug
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
		Email: gebruiker.Email,
		Role:  gebruiker.Rol,
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

// CreateUser creates a new user with hashed password
func (s *AuthServiceImpl) CreateUser(ctx context.Context, gebruiker *models.Gebruiker, password string) error {
	if password != "" {
		hashed, err := s.HashPassword(password)
		if err != nil {
			return err
		}
		gebruiker.WachtwoordHash = hashed
	} else {
		gebruiker.WachtwoordHash = "not_set"
	}
	return s.gebruikerRepo.Create(ctx, gebruiker)
}

// ListUsers lists users with pagination
func (s *AuthServiceImpl) ListUsers(ctx context.Context, limit, offset int) ([]*models.Gebruiker, error) {
	return s.gebruikerRepo.List(ctx, limit, offset)
}

// GetUser gets a user by ID
func (s *AuthServiceImpl) GetUser(ctx context.Context, id string) (*models.Gebruiker, error) {
	return s.gebruikerRepo.GetByID(ctx, id)
}

// UpdateUser updates a user, optionally changing password
func (s *AuthServiceImpl) UpdateUser(ctx context.Context, gebruiker *models.Gebruiker, password *string) error {
	if password != nil {
		hashed, err := s.HashPassword(*password)
		if err != nil {
			return err
		}
		gebruiker.WachtwoordHash = hashed
	}
	return s.gebruikerRepo.Update(ctx, gebruiker)
}

// DeleteUser deletes a user by ID
func (s *AuthServiceImpl) DeleteUser(ctx context.Context, id string) error {
	return s.gebruikerRepo.Delete(ctx, id)
}

// GenerateRefreshToken genereert een refresh token voor een gebruiker
func (s *AuthServiceImpl) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	// Genereer random token (32 bytes)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Error("Fout bij genereren random bytes voor refresh token", "error", err)
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Sla op in database met 7 dagen expiry
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		logger.Error("Fout bij opslaan refresh token", "user_id", userID, "error", err)
		return "", err
	}

	logger.Debug("Refresh token gegenereerd", "user_id", userID)
	return token, nil
}

// RefreshAccessToken vernieuwt een access token met een refresh token
func (s *AuthServiceImpl) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Valideer refresh token
	token, err := s.refreshTokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		logger.Error("Fout bij ophalen refresh token", "error", err)
		return "", "", err
	}

	if token == nil || !token.IsValid() {
		logger.Warn("Ongeldige of verlopen refresh token")
		return "", "", errors.New("ongeldige of verlopen refresh token")
	}

	// Haal gebruiker op
	gebruiker, err := s.gebruikerRepo.GetByID(ctx, token.UserID)
	if err != nil || gebruiker == nil {
		logger.Error("Gebruiker niet gevonden voor refresh token", "user_id", token.UserID, "error", err)
		return "", "", errors.New("gebruiker niet gevonden")
	}

	if !gebruiker.IsActief {
		logger.Warn("Inactieve gebruiker probeert token te refreshen", "user_id", gebruiker.ID)
		return "", "", ErrUserInactive
	}

	// Genereer nieuwe access token
	accessToken, err := s.generateToken(gebruiker)
	if err != nil {
		logger.Error("Fout bij genereren nieuwe access token", "user_id", gebruiker.ID, "error", err)
		return "", "", err
	}

	// Genereer nieuwe refresh token (token rotation voor security)
	newRefreshToken, err := s.GenerateRefreshToken(ctx, gebruiker.ID)
	if err != nil {
		logger.Error("Fout bij genereren nieuwe refresh token", "user_id", gebruiker.ID, "error", err)
		return "", "", err
	}

	// Revoke oude refresh token
	if err := s.refreshTokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		logger.Error("Fout bij revoken oude refresh token", "error", err)
		// Continue anyway, nieuwe tokens zijn al gegenereerd
	}

	logger.Info("Token refresh succesvol", "user_id", gebruiker.ID)
	return accessToken, newRefreshToken, nil
}

// RevokeRefreshToken trekt een refresh token in
func (s *AuthServiceImpl) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if err := s.refreshTokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		logger.Error("Fout bij revoken refresh token", "error", err)
		return err
	}
	logger.Debug("Refresh token ingetrokken")
	return nil
}

// RevokeAllUserRefreshTokens trekt alle refresh tokens van een gebruiker in
func (s *AuthServiceImpl) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	if err := s.refreshTokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		logger.Error("Fout bij revoken alle refresh tokens", "user_id", userID, "error", err)
		return err
	}
	logger.Info("Alle refresh tokens ingetrokken", "user_id", userID)
	return nil
}
