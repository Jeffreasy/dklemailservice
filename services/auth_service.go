package services

import (
	"context"
	"crypto/rand"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"encoding/base64"
	"errors" // Import toegevoegd
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm" // Import toegevoegd
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
	// DEPRECATED: Legacy field - will be removed in future version
	// Frontend should only use Roles array
	Role       string   `json:"role,omitempty"` // DEPRECATED - use Roles instead
	Roles      []string `json:"roles"`          // RBAC roles from user_roles table
	RBACActive bool     `json:"rbac_active"`    // Indicates if RBAC system is active
	jwt.RegisteredClaims
}

// AuthServiceImpl implementeert de AuthService interface
// V30+RBAC: Uitgebreid met participant repository voor app access checks
type AuthServiceImpl struct {
	gebruikerRepo    repository.GebruikerRepository
	refreshTokenRepo repository.RefreshTokenRepository
	userRoleRepo     repository.UserRoleRepository
	participantRepo  repository.ParticipantRepository // V30+RBAC: Voor participant app access checks
	jwtSecret        []byte
	tokenExpiry      time.Duration
}

// NewAuthService maakt een nieuwe AuthService
func NewAuthService(gebruikerRepo repository.GebruikerRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return NewAuthServiceWithRBAC(gebruikerRepo, refreshTokenRepo, nil)
}

// NewAuthServiceWithRBAC maakt een nieuwe AuthService met RBAC support
func NewAuthServiceWithRBAC(gebruikerRepo repository.GebruikerRepository, refreshTokenRepo repository.RefreshTokenRepository, userRoleRepo repository.UserRoleRepository) AuthService {
	return NewAuthServiceWithParticipantSupport(gebruikerRepo, refreshTokenRepo, userRoleRepo, nil)
}

// NewAuthServiceWithParticipantSupport maakt een nieuwe AuthService met volledige participant integratie
// V30+RBAC: Voegt participant repository toe voor app access validatie
func NewAuthServiceWithParticipantSupport(
	gebruikerRepo repository.GebruikerRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	userRoleRepo repository.UserRoleRepository,
	participantRepo repository.ParticipantRepository,
) AuthService {
	// Haal JWT secret uit omgevingsvariabele - VERPLICHT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("JWT_SECRET omgevingsvariabele is niet ingesteld. Dit is verplicht voor security.")
	}

	// Valideer minimale lengte voor security
	if len(jwtSecret) < 32 {
		logger.Fatal("JWT_SECRET moet minimaal 32 karakters bevatten voor adequate security", "length", len(jwtSecret))
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
		userRoleRepo:     userRoleRepo,
		participantRepo:  participantRepo, // V30+RBAC
		jwtSecret:        []byte(jwtSecret),
		tokenExpiry:      tokenExpiry,
	}
}

// Login authenticeert een gebruiker en geeft een access token en refresh token terug
// OPLOSSING 1: Volledig herschreven voor veiligheid en correcte volgorde
func (s *AuthServiceImpl) Login(ctx context.Context, email, wachtwoord string) (string, string, error) {
	logger.Info("Login poging", "email", email)

	// =========================================================================
	// OPLOSSING 1: STAP 1 (VOORHEEN STAP 2)
	// Probeer EERST gebruiker login (admin/staff accounts)
	// Dit voorkomt dat een admin die ook een participant-record heeft,
	// wordt ingelogd als participant.
	// =========================================================================
	gebruiker, err := s.gebruikerRepo.GetByEmail(ctx, email)

	if err == nil && gebruiker != nil {
		// Gebruiker (admin/staff) GEVONDEN. Valideer alleen hiertegen.

		// Controleer of gebruiker actief is
		if !gebruiker.IsActief {
			logger.Warn("Inactieve gebruiker probeert in te loggen", "email", email)
			return "", "", ErrUserInactive
		}

		// Verifieer wachtwoord
		if s.VerifyPassword(gebruiker.WachtwoordHash, wachtwoord) {
			// Wachtwoord klopt. Log in als admin/staff.
			logger.Info("Login succesvol (via gebruiker)", "email", email, "user_id", gebruiker.ID)

			// Update laatste login
			if err := s.gebruikerRepo.UpdateLastLogin(ctx, gebruiker.ID); err != nil {
				logger.Error("Fout bij updaten laatste login", "email", email, "error", err)
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

			// BELANGRIJK: Return hier, ga NIET door naar participant check
			return accessToken, refreshToken, nil
		} else {
			// Gebruiker (admin/staff) gevonden, maar wachtwoord is FOUT.
			// Geef direct foutmelding. Val NIET door naar participant check.
			logger.Warn("Ongeldig wachtwoord voor gebruiker (admin/staff)", "email", email)
			return "", "", ErrInvalidCredentials
		}
	}

	// Als we hier zijn, is het OF 'gebruiker niet gevonden' OF een DB-fout.
	// Als het een DB-fout is (anders dan 'niet gevonden'), stop hier.
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Databasefout bij ophalen gebruiker", "email", email, "error", err)
		return "", "", err
	}

	// =========================================================================
	// OPLOSSING 1: STAP 2 (VOORHEEN STAP 1)
	// Gebruiker (admin/staff) niet gevonden, probeer nu participant login.
	// =========================================================================
	if s.participantRepo != nil {
		participant, err := s.loginViaParticipant(ctx, email, wachtwoord)
		if err == nil && participant != nil {
			// Participant login succesvol!
			return s.generateParticipantTokens(ctx, participant)
		}
		// Als fout != ErrInvalidCredentials, log maar ga door naar 'failure'
		if err != nil && err != ErrInvalidCredentials {
			logger.Debug("Participant login failed", "error", err)
		}
	}

	// Als beide checks falen, geef standaard foutmelding
	logger.Warn("Ongeldige inloggegevens (geen gebruiker of participant gevonden)", "email", email)
	return "", "", ErrInvalidCredentials
}

// loginViaParticipant probeert in te loggen via participants.wachtwoord_hash
// V34: Nieuwe functie voor full account participant authenticatie
func (s *AuthServiceImpl) loginViaParticipant(ctx context.Context, email, wachtwoord string) (*models.Participant, error) {
	// Zoek participant met dit email
	participants, err := s.participantRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Debug("Fout bij ophalen participant", "email", email, "error", err)
		return nil, err
	}

	// Zoek full account participant met wachtwoord
	for _, participant := range participants {
		if participant.AccountType == "full" && participant.WachtwoordHash != nil {
			// Verifieer wachtwoord
			if s.VerifyPassword(*participant.WachtwoordHash, wachtwoord) {
				// Check app access
				if !participant.HasAppAccess {
					logger.Warn("Participant heeft geen app toegang", "email", email)
					return nil, errors.New("geen app toegang")
				}
				logger.Info("Participant login succesvol", "participant_id", participant.ID, "email", email)
				return participant, nil
			}
		}
	}

	return nil, ErrInvalidCredentials
}

// generateParticipantTokens genereert tokens voor een participant (V34)
func (s *AuthServiceImpl) generateParticipantTokens(ctx context.Context, participant *models.Participant) (string, string, error) {
	// Maak een fake gebruiker object voor token generatie
	// Dit zorgt ervoor dat bestaande token logica blijft werken
	fakeGebruiker := &models.Gebruiker{
		ID:             participant.ID, // Gebruik participant ID als user ID
		Email:          participant.Email,
		Naam:           participant.Naam,
		WachtwoordHash: *participant.WachtwoordHash,
		IsActief:       true,
	}

	// Genereer JWT access token (zonder RBAC roles voor participants)
	accessToken, err := s.generateParticipantToken(fakeGebruiker)
	if err != nil {
		logger.Error("Fout bij genereren participant access token", "email", participant.Email, "error", err)
		return "", "", err
	}

	// Genereer refresh token
	refreshToken, err := s.GenerateRefreshToken(ctx, participant.ID)
	if err != nil {
		logger.Error("Fout bij genereren participant refresh token", "email", participant.Email, "error", err)
		return "", "", err
	}

	logger.Info("Participant tokens gegenereerd", "participant_id", participant.ID, "email", participant.Email)
	return accessToken, refreshToken, nil
}

// generateParticipantToken genereert een JWT token voor een participant (zonder RBAC)
// OPLOSSING 2: Deze functie is de *juiste* voor de V34 architectuur.
func (s *AuthServiceImpl) generateParticipantToken(participant *models.Gebruiker) (string, error) {
	// Participants krijgen participant_user rol in JWT
	claims := JWTClaims{
		Email:      participant.Email,
		Role:       "participant_user",           // Voor backward compatibility
		Roles:      []string{"participant_user"}, // Participant krijgt altijd deze rol
		RBACActive: false,                        // Geen RBAC voor participants
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dklemailservice",
			Subject:   participant.ID,
		},
	}

	// Maak en onderteken token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
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

	// OPLOSSING 1: Log level verlaagd van Info naar Debug om log-spam te voorkomen
	logger.Debug("Token gevalideerd", "user_id", claims.Subject) // Gebruik claims.Subject
	return claims.Subject, nil                                   // Geef Subject (user ID) terug
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
// Deze functie is nu ALLEEN voor Gebruikers (admins/staff)
func (s *AuthServiceImpl) generateToken(gebruiker *models.Gebruiker) (string, error) {
	// Haal RBAC roles op voor de gebruiker
	rbacRoles := s.getUserRBACRoles(gebruiker.ID)

	// Fallback: als geen RBAC roles, gebruik eerste role name of lege string
	legacyRole := ""
	if len(rbacRoles) > 0 {
		legacyRole = rbacRoles[0] // Eerste role voor backward compatibility
	}

	// Maak claims - RBAC is de primary bron van truth
	claims := JWTClaims{
		Email:      gebruiker.Email,
		Role:       legacyRole,         // DEPRECATED - alleen voor backward compatibility
		Roles:      rbacRoles,          // RBAC - primary bron
		RBACActive: len(rbacRoles) > 0, // True als RBAC roles aanwezig
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

// getUserRBACRoles haalt de RBAC role namen op voor een gebruiker
func (s *AuthServiceImpl) getUserRBACRoles(userID string) []string {
	// Als userRoleRepo niet beschikbaar is, return lege array
	if s.userRoleRepo == nil {
		return []string{}
	}

	// Gebruik context met timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Haal actieve user roles op
	userRoles, err := s.userRoleRepo.ListActiveByUser(ctx, userID)
	if err != nil {
		logger.Warn("Kon RBAC roles niet ophalen voor JWT", "user_id", userID, "error", err)
		return []string{}
	}

	// Extract role names
	roleNames := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		if ur.Role.Name != "" {
			roleNames = append(roleNames, ur.Role.Name)
		}
	}

	return roleNames
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

// SearchUsers searches for users by name or email
func (s *AuthServiceImpl) SearchUsers(ctx context.Context, query string, limit int) ([]*models.Gebruiker, error) {
	return s.gebruikerRepo.Search(ctx, query, limit)
}

// GenerateRefreshToken genereert een refresh token voor een gebruiker of participant
// V34: Gebruikt OwnerID om zowel gebruiker als participant IDs te ondersteunen
func (s *AuthServiceImpl) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	// Genereer random token (32 bytes)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Error("Fout bij genereren random bytes voor refresh token", "error", err)
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Sla op in database met 7 dagen expiry
	// V34: OwnerID kan zowel gebruiker als participant ID zijn
	refreshToken := &models.RefreshToken{
		OwnerID:   userID,
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
// V34: Ondersteunt zowel gebruiker als participant tokens
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

	// V34: Probeer eerst als gebruiker, dan als participant
	gebruiker, err := s.gebruikerRepo.GetByID(ctx, token.OwnerID)
	if err != nil || gebruiker == nil {
		// Probeer als participant
		if s.participantRepo != nil {
			participant, err := s.participantRepo.GetByID(ctx, token.OwnerID)
			if err != nil || participant == nil {
				logger.Error("Owner niet gevonden voor refresh token", "owner_id", token.OwnerID, "error", err)
				return "", "", errors.New("owner niet gevonden")
			}

			// Participant gevonden - genereer participant tokens
			return s.refreshParticipantToken(ctx, participant, refreshToken)
		}

		logger.Error("Gebruiker niet gevonden voor refresh token", "owner_id", token.OwnerID, "error", err)
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

// refreshParticipantToken vernieuwt tokens voor een participant
// V34: Helper functie voor participant token refresh
func (s *AuthServiceImpl) refreshParticipantToken(ctx context.Context, participant *models.Participant, oldToken string) (string, string, error) {
	// Check app access
	if !participant.HasAppAccess {
		logger.Warn("Participant heeft geen app toegang bij token refresh", "participant_id", participant.ID)
		return "", "", errors.New("geen app toegang")
	}

	// Maak fake gebruiker voor token generatie
	fakeGebruiker := &models.Gebruiker{
		ID:             participant.ID,
		Email:          participant.Email,
		Naam:           participant.Naam,
		WachtwoordHash: *participant.WachtwoordHash,
		IsActief:       true,
	}

	// Genereer nieuwe access token
	// OPLOSSING 2: Gebruik generateParticipantToken()
	accessToken, err := s.generateParticipantToken(fakeGebruiker)
	if err != nil {
		logger.Error("Fout bij genereren nieuwe participant access token", "participant_id", participant.ID, "error", err)
		return "", "", err
	}

	// Genereer nieuwe refresh token
	newRefreshToken, err := s.GenerateRefreshToken(ctx, participant.ID)
	if err != nil {
		logger.Error("Fout bij genereren nieuwe participant refresh token", "participant_id", participant.ID, "error", err)
		return "", "", err
	}

	// Revoke oude refresh token
	if err := s.refreshTokenRepo.RevokeToken(ctx, oldToken); err != nil {
		logger.Error("Fout bij revoken oude refresh token", "error", err)
		// Continue anyway
	}

	logger.Info("Participant token refresh succesvol", "participant_id", participant.ID)
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

// ==============================================================================
// V30+RBAC: PARTICIPANT-SPECIFIEKE HELPER FUNCTIES
// ==============================================================================

// validateParticipantAppAccess controleert of een gebruiker app toegang heeft via participant record
// V30+RBAC: Gebruikt voor login validatie - alleen full account participants met app access mogen inloggen
func (s *AuthServiceImpl) validateParticipantAppAccess(ctx context.Context, gebruikerID string, email string) bool {
	// Zoek participant records gekoppeld aan deze gebruiker
	participants, err := s.participantRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Error("Fout bij ophalen participant voor app access check",
			"gebruiker_id", gebruikerID, "email", email, "error", err)
		// Bij database fout: allow login (fail-open voor backwards compatibility)
		return true
	}

	// Zoek full account participant gekoppeld aan deze gebruiker
	for _, participant := range participants {
		if participant.GebruikerID != nil && *participant.GebruikerID == gebruikerID {
			// Check account type en app access
			if participant.AccountType == "full" && participant.HasAppAccess {
				logger.Debug("Participant app access validated",
					"participant_id", participant.ID,
					"gebruiker_id", gebruikerID,
					"account_type", participant.AccountType)
				return true
			}

			// Participant gevonden maar geen app access
			logger.Warn("Participant found but no app access",
				"participant_id", participant.ID,
				"account_type", participant.AccountType,
				"has_app_access", participant.HasAppAccess)
			return false
		}
	}

	// Geen participant gekoppeld aan deze gebruiker - dit is een legacy gebruiker (allow login)
	logger.Debug("Geen participant gekoppeld aan gebruiker (legacy user)", "gebruiker_id", gebruikerID)
	return true
}

// GetParticipantByGebruikerID haalt participant op basis van gebruiker ID
// V30+RBAC: Helper voor app om participant data op te halen
func (s *AuthServiceImpl) GetParticipantByGebruikerID(ctx context.Context, gebruikerID string) (*models.Participant, error) {
	if s.participantRepo == nil {
		return nil, errors.New("participant repository niet beschikbaar")
	}

	// Zoek participant met deze gebruiker_id
	// Dit vereist een nieuwe repository methode - we gebruiken een workaround via email
	gebruiker, err := s.gebruikerRepo.GetByID(ctx, gebruikerID)
	if err != nil || gebruiker == nil {
		return nil, fmt.Errorf("gebruiker niet gevonden: %w", err)
	}

	participants, err := s.participantRepo.FindByEmail(ctx, gebruiker.Email)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen participant: %w", err)
	}

	// Zoek de participant gekoppeld aan deze gebruiker
	for _, participant := range participants {
		if participant.GebruikerID != nil && *participant.GebruikerID == gebruikerID {
			return participant, nil
		}
	}

	return nil, errors.New("geen participant gevonden voor deze gebruiker")
}
