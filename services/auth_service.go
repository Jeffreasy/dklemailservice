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
	"gorm.io/gorm"
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
	Email      string   `json:"email"`
	Roles      []string `json:"roles"`       // RBAC roles from user_roles table
	RBACActive bool     `json:"rbac_active"` // Indicates if RBAC system is active
	jwt.RegisteredClaims
}

// AuthServiceImpl implementeert de AuthService interface
// V30+RBAC: Uitgebreid met participant repository voor app access checks
// RBAC: Uitgebreid met rbacRoleRepo voor role management
// Email Verification: Uitgebreid met email verification token repository en email service
// Access Token Rotation: Uitgebreid met access token repository voor server-side opslag
// Session Management: Uitgebreid met session repository voor multi-device session tracking
type AuthServiceImpl struct {
	gebruikerRepo              repository.GebruikerRepository
	refreshTokenRepo           repository.RefreshTokenRepository
	accessTokenRepo            repository.AccessTokenRepository // Access Token Rotation: Voor server-side token opslag
	passwordResetTokenRepo     repository.PasswordResetTokenRepository
	emailVerificationTokenRepo repository.EmailVerificationTokenRepository
	userRoleRepo               repository.UserRoleRepository
	rbacRoleRepo               repository.RBACRoleRepository    // RBAC: Voor role management
	participantRepo            repository.ParticipantRepository // V30+RBAC: Voor participant app access checks
	sessionRepo                repository.SessionRepository     // Session Management: Voor session tracking
	emailService               EmailSender                      // Email Verification: Voor verzenden verificatie emails
	jwtSecret                  []byte
	tokenExpiry                time.Duration
	frontendURL                string // Voor dynamische email links
}

// NewAuthService maakt een nieuwe AuthService
func NewAuthService(gebruikerRepo repository.GebruikerRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return NewAuthServiceWithRBAC(gebruikerRepo, refreshTokenRepo, nil, nil)
}

// NewAuthServiceWithRBAC maakt een nieuwe AuthService met RBAC support
func NewAuthServiceWithRBAC(gebruikerRepo repository.GebruikerRepository, refreshTokenRepo repository.RefreshTokenRepository, userRoleRepo repository.UserRoleRepository, rbacRoleRepo repository.RBACRoleRepository) AuthService {
	return NewAuthServiceWithParticipantSupport(gebruikerRepo, refreshTokenRepo, nil, nil, nil, userRoleRepo, rbacRoleRepo, nil, nil, nil)
}

// NewAuthServiceWithParticipantSupport maakt een nieuwe AuthService met volledige participant integratie
// V30+RBAC: Voegt participant repository toe voor app access validatie
// RBAC: Voegt rbacRoleRepo toe voor role management
// Email Verification: Voegt email verification token repository en email service toe
// Access Token Rotation: Voegt access token repository toe voor server-side opslag
// Session Management: Voegt session repository toe voor multi-device session tracking
func NewAuthServiceWithParticipantSupport(
	gebruikerRepo repository.GebruikerRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	accessTokenRepo repository.AccessTokenRepository,
	passwordResetTokenRepo repository.PasswordResetTokenRepository,
	emailVerificationTokenRepo repository.EmailVerificationTokenRepository,
	userRoleRepo repository.UserRoleRepository,
	rbacRoleRepo repository.RBACRoleRepository,
	participantRepo repository.ParticipantRepository,
	sessionRepo repository.SessionRepository,
	emailService EmailSender,
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

	// Haal frontend URL uit omgevingsvariabele voor email links
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://dekoninklijkeloop.nl" // Default fallback
		logger.Warn("FRONTEND_URL omgevingsvariabele niet ingesteld, gebruik default", "default", frontendURL)
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
		gebruikerRepo:              gebruikerRepo,
		refreshTokenRepo:           refreshTokenRepo,
		accessTokenRepo:            accessTokenRepo, // Access Token Rotation
		passwordResetTokenRepo:     passwordResetTokenRepo,
		emailVerificationTokenRepo: emailVerificationTokenRepo,
		userRoleRepo:               userRoleRepo,
		rbacRoleRepo:               rbacRoleRepo,    // RBAC Optimization
		participantRepo:            participantRepo, // V30+RBAC
		sessionRepo:                sessionRepo,     // Session Management
		emailService:               emailService,    // Email Verification
		jwtSecret:                  []byte(jwtSecret),
		tokenExpiry:                tokenExpiry,
		frontendURL:                frontendURL,
	}
}

// Login authenticeert een gebruiker en geeft een access token en refresh token terug
// OPLOSSING 1: Volledig herschreven voor veiligheid en correcte volgorde
func (s *AuthServiceImpl) Login(ctx context.Context, email, wachtwoord string) (string, string, error) {
	logger.Info("Login poging", "email", email)

	// =========================================================================
	// OPLOSSING 1: STAP 1 (VOORHEEN STAP 2)
	// Probeer EERST gebruiker login (admin/staff accounts)
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
			accessToken, err := s.generateToken(ctx, gebruiker) // << AANGEPAST: ctx toegevoegd
			if err != nil {
				logger.Error("Fout bij genereren access token", "email", email, "error", err)
				return "", "", err
			}

			// Sla access token op in database voor server-side validatie
			if err := s.storeAccessToken(ctx, accessToken, gebruiker.ID); err != nil {
				logger.Error("Fout bij opslaan access token", "email", email, "error", err)
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

	// Sla access token op in database voor server-side validatie
	if err := s.storeAccessToken(ctx, accessToken, participant.ID); err != nil {
		logger.Error("Fout bij opslaan participant access token", "email", participant.Email, "error", err)
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
func (s *AuthServiceImpl) generateToken(ctx context.Context, gebruiker *models.Gebruiker) (string, error) { // << AANGEPAST: ctx toegevoegd
	// Haal RBAC roles op voor de gebruiker
	rbacRoles := s.getUserRBACRoles(ctx, gebruiker.ID) // << AANGEPAST: ctx doorgegeven

	// Maak claims - RBAC is de primary bron van truth
	claims := JWTClaims{
		Email:      gebruiker.Email,
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
// Legacy role_id system removed - only uses user_roles
func (s *AuthServiceImpl) getUserRBACRoles(ctx context.Context, userID string) []string {
	// Als userRoleRepo niet beschikbaar is, return empty roles
	if s.userRoleRepo == nil {
		logger.Warn("UserRole repository niet beschikbaar, geen roles voor JWT", "user_id", userID)
		return []string{}
	}

	// Gebruik de doorgegeven context (met eventuele timeout/cancellation)
	// We gebruiken hier een timeout van 2 seconden voor de DB-queries om excessieve latentie te voorkomen.
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Get active user_roles
	userRoles, err := s.userRoleRepo.ListActiveByUser(ctxTimeout, userID)
	if err != nil {
		logger.Warn("Error getting user_roles for JWT", "user_id", userID, "error", err)
		return []string{}
	}

	// Extract role names
	roleNames := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		if ur.Role.Name != "" {
			roleNames = append(roleNames, ur.Role.Name)
		}
	}

	logger.Debug("Retrieved roles from user_roles", "user_id", userID, "roles", roleNames)
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
	accessToken, err := s.generateToken(ctx, gebruiker) // << AANGEPAST: ctx doorgegeven
	if err != nil {
		logger.Error("Fout bij genereren nieuwe access token", "user_id", gebruiker.ID, "error", err)
		return "", "", err
	}

	// Sla nieuwe access token op in database voor server-side validatie
	if err := s.storeAccessToken(ctx, accessToken, gebruiker.ID); err != nil {
		logger.Error("Fout bij opslaan nieuwe access token", "user_id", gebruiker.ID, "error", err)
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

	// Sla nieuwe access token op in database voor server-side validatie
	if err := s.storeAccessToken(ctx, accessToken, participant.ID); err != nil {
		logger.Error("Fout bij opslaan nieuwe participant access token", "participant_id", participant.ID, "error", err)
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

// ListUserSessions haalt alle actieve sessies op voor een gebruiker
func (s *AuthServiceImpl) ListUserSessions(ctx context.Context, userID string) ([]*models.Session, error) {
	if s.sessionRepo == nil {
		logger.Warn("Session repository niet beschikbaar, sessies kunnen niet worden opgehaald")
		return []*models.Session{}, nil // Graceful degradation
	}

	sessions, err := s.sessionRepo.ListByOwnerID(ctx, userID)
	if err != nil {
		logger.Error("Fout bij ophalen sessies", "user_id", userID, "error", err)
		return nil, err
	}

	logger.Debug("Sessies opgehaald", "user_id", userID, "count", len(sessions))
	return sessions, nil
}

// RevokeSession trekt een specifieke sessie in
func (s *AuthServiceImpl) RevokeSession(ctx context.Context, sessionID string) error {
	if s.sessionRepo == nil {
		logger.Warn("Session repository niet beschikbaar, sessie wordt niet ingetrokken")
		return nil // Graceful degradation
	}

	// Haal eerst de sessie op om de access token te krijgen
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.Error("Fout bij ophalen sessie voor intrekking", "session_id", sessionID, "error", err)
		return err
	}

	if session == nil {
		logger.Warn("Sessie niet gevonden voor intrekking", "session_id", sessionID)
		return errors.New("sessie niet gevonden")
	}

	// Trek de sessie in
	if err := s.sessionRepo.RevokeByAccessToken(ctx, session.AccessToken); err != nil {
		logger.Error("Fout bij intrekken sessie", "session_id", sessionID, "error", err)
		return err
	}

	// Trek ook de access token in
	if s.accessTokenRepo != nil {
		if err := s.accessTokenRepo.RevokeToken(ctx, session.AccessToken); err != nil {
			logger.Error("Fout bij intrekken access token van sessie", "session_id", sessionID, "error", err)
			// Continue anyway
		}
	}

	logger.Info("Sessie ingetrokken", "session_id", sessionID, "user_id", session.OwnerID)
	return nil
}

// RevokeAllUserSessions trekt alle sessies van een gebruiker in
func (s *AuthServiceImpl) RevokeAllUserSessions(ctx context.Context, userID string) error {
	if s.sessionRepo == nil {
		logger.Warn("Session repository niet beschikbaar, sessies worden niet ingetrokken")
		return nil // Graceful degradation
	}

	if err := s.sessionRepo.RevokeAllUserSessionsComplete(ctx, userID); err != nil {
		logger.Error("Fout bij intrekken alle sessies", "user_id", userID, "error", err)
		return err
	}

	// Trek ook alle access tokens in
	if s.accessTokenRepo != nil {
		if err := s.accessTokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
			logger.Error("Fout bij intrekken alle access tokens", "user_id", userID, "error", err)
			// Continue anyway
		}
	}

	logger.Info("Alle sessies ingetrokken", "user_id", userID)
	return nil
}

// RevokeOtherUserSessions trekt alle sessies van een gebruiker in behalve de huidige sessie
func (s *AuthServiceImpl) RevokeOtherUserSessions(ctx context.Context, userID, currentSessionID string) error {
	if s.sessionRepo == nil {
		logger.Warn("Session repository niet beschikbaar, andere sessies worden niet ingetrokken")
		return nil // Graceful degradation
	}

	if err := s.sessionRepo.RevokeAllUserSessions(ctx, userID, currentSessionID); err != nil {
		logger.Error("Fout bij intrekken andere sessies", "user_id", userID, "current_session_id", currentSessionID, "error", err)
		return err
	}

	logger.Info("Andere sessies ingetrokken", "user_id", userID, "current_session_id", currentSessionID)
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

// RequestPasswordReset vraagt een wachtwoord reset aan voor een email adres
func (s *AuthServiceImpl) RequestPasswordReset(ctx context.Context, email string) error {
	logger.Info("Password reset request", "email", email)

	// Controleer of gebruiker bestaat (zowel admin/staff als participant)
	var userExists bool
	var userType string

	// Check admin/staff users
	if gebruiker, err := s.gebruikerRepo.GetByEmail(ctx, email); err == nil && gebruiker != nil && gebruiker.IsActief {
		userExists = true
		userType = "gebruiker"
		logger.Debug("Found active gebruiker for password reset", "email", email, "user_id", gebruiker.ID)
	}

	// Check participants if not found as gebruiker
	if !userExists && s.participantRepo != nil {
		if participants, err := s.participantRepo.FindByEmail(ctx, email); err == nil {
			for _, participant := range participants {
				if participant.AccountType == "full" && participant.HasAppAccess && participant.WachtwoordHash != nil {
					userExists = true
					userType = "participant"
					logger.Debug("Found active participant for password reset", "email", email, "participant_id", participant.ID)
					break
				}
			}
		}
	}

	// Als gebruiker niet bestaat, geef geen foutmelding terug voor security (geen email enumeration)
	if !userExists {
		logger.Info("Password reset requested for non-existent or inactive user", "email", email)
		return nil // Silent success voor security
	}

	// Genereer reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Error("Failed to generate random bytes for password reset token", "error", err)
		return errors.New("kon geen reset token genereren")
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Sla token op in database (1 uur expiry)
	resetToken := &models.PasswordResetToken{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		IsUsed:    false,
	}

	if err := s.passwordResetTokenRepo.Create(ctx, resetToken); err != nil {
		logger.Error("Failed to save password reset token", "email", email, "error", err)
		return errors.New("kon reset token niet opslaan")
	}

	// Genereer reset link met frontend URL uit config
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)

	// TODO: Email service integreren
	// Voor nu loggen we alleen de link
	logger.Info("Password reset token generated", "email", email, "token", token, "user_type", userType, "reset_link", resetLink)

	// Hier zou normaal gesproken een email verzonden worden naar de gebruiker
	// met een link naar de frontend waar ze hun wachtwoord kunnen resetten
	// Bijvoorbeeld:
	// if err := s.emailService.SendPasswordResetEmail(email, resetLink); err != nil {
	// 	logger.Error("Failed to send password reset email", "email", email, "error", err)
	// 	return errors.New("kon reset email niet verzenden")
	// }

	return nil
}

// ResetPasswordWithToken reset het wachtwoord met een geldige reset token
func (s *AuthServiceImpl) ResetPasswordWithToken(ctx context.Context, token, newPassword string) error {
	logger.Info("Password reset with token attempt")

	// Haal token op uit database
	resetToken, err := s.passwordResetTokenRepo.GetByToken(ctx, token)
	if err != nil {
		logger.Error("Failed to get password reset token", "error", err)
		return errors.New("ongeldige reset token")
	}

	if resetToken == nil {
		logger.Warn("Password reset token not found", "token", token[:8]+"...")
		return errors.New("ongeldige reset token")
	}

	if !resetToken.IsValid() {
		logger.Warn("Password reset token is invalid or expired", "token", token[:8]+"...", "is_used", resetToken.IsUsed, "expires_at", resetToken.ExpiresAt)
		return errors.New("reset token is verlopen of al gebruikt")
	}

	// Controleer wachtwoord sterkte (minimaal 8 karakters)
	if len(newPassword) < 8 {
		logger.Warn("Password too short", "email", resetToken.Email)
		return errors.New("wachtwoord moet minimaal 8 karakters bevatten")
	}

	// Reset wachtwoord voor gebruiker of participant
	var resetErr error
	var userType string

	// Probeer eerst als admin/staff gebruiker
	if gebruiker, err := s.gebruikerRepo.GetByEmail(ctx, resetToken.Email); err == nil && gebruiker != nil && gebruiker.IsActief {
		resetErr = s.ResetPassword(ctx, resetToken.Email, newPassword)
		userType = "gebruiker"
	} else if s.participantRepo != nil {
		// Probeer als participant
		if participants, err := s.participantRepo.FindByEmail(ctx, resetToken.Email); err == nil {
			for _, participant := range participants {
				if participant.AccountType == "full" && participant.HasAppAccess && participant.WachtwoordHash != nil {
					// Hash nieuw wachtwoord
					hashedPassword, err := s.HashPassword(newPassword)
					if err != nil {
						logger.Error("Failed to hash new password for participant", "participant_id", participant.ID, "error", err)
						resetErr = errors.New("kon wachtwoord niet hashen")
						break
					}

					// Update participant wachtwoord
					participant.WachtwoordHash = &hashedPassword
					resetErr = s.participantRepo.Update(ctx, participant)
					userType = "participant"
					break
				}
			}
		}
	}

	if resetErr != nil {
		logger.Error("Failed to reset password", "email", resetToken.Email, "user_type", userType, "error", resetErr)
		return errors.New("kon wachtwoord niet resetten")
	}

	// Markeer token als gebruikt
	if err := s.passwordResetTokenRepo.MarkAsUsed(ctx, token); err != nil {
		logger.Error("Failed to mark password reset token as used", "token", token[:8]+"...", "error", err)
		// Dit is niet kritisch, continue
	}

	logger.Info("Password reset successful", "email", resetToken.Email, "user_type", userType)
	return nil
}

// SendEmailVerification verzendt een email verificatie token naar een gebruiker
func (s *AuthServiceImpl) SendEmailVerification(ctx context.Context, email, userID, userType string) error {
	logger.Info("Email verification request", "email", email, "user_id", userID, "user_type", userType)

	// Controleer of gebruiker bestaat en actief is
	var userExists bool
	if userType == "gebruiker" {
		gebruiker, err := s.gebruikerRepo.GetByID(ctx, userID)
		if err != nil || gebruiker == nil || !gebruiker.IsActief {
			logger.Warn("Gebruiker niet gevonden of inactief voor email verificatie", "user_id", userID)
			return errors.New("gebruiker niet gevonden of inactief")
		}
		userExists = true
	} else if userType == "participant" && s.participantRepo != nil {
		participant, err := s.participantRepo.GetByID(ctx, userID)
		if err != nil || participant == nil || !participant.HasAppAccess {
			logger.Warn("Participant niet gevonden of geen app toegang voor email verificatie", "user_id", userID)
			return errors.New("participant niet gevonden of geen app toegang")
		}
		userExists = true
	}

	if !userExists {
		logger.Warn("Ongeldig user type voor email verificatie", "user_type", userType)
		return errors.New("ongeldig user type")
	}

	// Controleer of er al een actieve verificatie token bestaat
	existingTokens, err := s.emailVerificationTokenRepo.GetByUserID(ctx, userID, userType)
	if err != nil {
		logger.Error("Fout bij ophalen bestaande verificatie tokens", "user_id", userID, "error", err)
		return err
	}

	// Als er al een actieve token bestaat, geef een foutmelding terug
	for _, token := range existingTokens {
		if token.IsValid() {
			logger.Info("Actieve verificatie token bestaat al", "user_id", userID, "email", email)
			return errors.New("er is al een actieve verificatie token voor dit account")
		}
	}

	// Genereer verificatie token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Error("Fout bij genereren random bytes voor verificatie token", "error", err)
		return errors.New("kon geen verificatie token genereren")
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Sla token op in database (24 uur expiry)
	verificationToken := &models.EmailVerificationToken{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsUsed:    false,
		UserType:  userType,
		UserID:    userID,
	}

	if err := s.emailVerificationTokenRepo.Create(ctx, verificationToken); err != nil {
		logger.Error("Fout bij opslaan verificatie token", "email", email, "error", err)
		return errors.New("kon verificatie token niet opslaan")
	}

	// Genereer verificatie link met frontend URL uit config
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, token)

	logger.Info("Email verification token generated", "email", email, "token", token[:8]+"...", "user_type", userType, "verification_link", verificationLink)

	// Verstuur verificatie email
	if s.emailService != nil {
		// Haal naam op voor personalisatie
		var naam string
		if userType == "gebruiker" {
			if gebruiker, err := s.gebruikerRepo.GetByID(ctx, userID); err == nil && gebruiker != nil {
				naam = gebruiker.Naam
			}
		} else if userType == "participant" && s.participantRepo != nil {
			if participant, err := s.participantRepo.GetByID(ctx, userID); err == nil && participant != nil {
				naam = participant.Naam
			}
		}

		if naam == "" {
			naam = "Gebruiker" // Fallback naam
		}

		// Verstuur email met 24 uur expiry
		if err := s.emailService.SendEmailVerificationEmail(email, naam, verificationLink, 24); err != nil {
			logger.Error("Failed to send email verification", "email", email, "error", err)
			return errors.New("kon verificatie email niet verzenden")
		}
	} else {
		logger.Warn("Email service not available, skipping email verification send", "email", email)
	}

	return nil
}

// VerifyEmailWithToken verifieert een email adres met een token
func (s *AuthServiceImpl) VerifyEmailWithToken(ctx context.Context, token string) error {
	logger.Info("Email verification attempt with token")

	// Haal token op uit database
	verificationToken, err := s.emailVerificationTokenRepo.GetByToken(ctx, token)
	if err != nil {
		logger.Error("Fout bij ophalen verificatie token", "error", err)
		return errors.New("ongeldige verificatie token")
	}

	if verificationToken == nil {
		logger.Warn("Verificatie token niet gevonden", "token", token[:8]+"...")
		return errors.New("ongeldige verificatie token")
	}

	if !verificationToken.IsValid() {
		logger.Warn("Verificatie token is verlopen of al gebruikt", "token", token[:8]+"...", "is_used", verificationToken.IsUsed, "expires_at", verificationToken.ExpiresAt)
		return errors.New("verificatie token is verlopen of al gebruikt")
	}

	// Update user verification status
	// Opgelost: Gebruik nu een tagged switch op verificationToken.UserType (QF1003)
	switch verificationToken.UserType {
	case "gebruiker":
		// Voor gebruikers: we zouden een email_verified veld kunnen toevoegen aan de gebruiker tabel
		// Voor nu markeren we alleen de token als gebruikt
		logger.Info("Gebruiker email verificatie succesvol", "user_id", verificationToken.UserID, "email", verificationToken.Email)
	case "participant":
		// Voor participants: update email_verified status
		if s.participantRepo != nil {
			participant, err := s.participantRepo.GetByID(ctx, verificationToken.UserID)
			if err != nil || participant == nil {
				logger.Error("Participant niet gevonden voor verificatie", "user_id", verificationToken.UserID, "error", err)
				return errors.New("participant niet gevonden")
			}

			// Stel email_verified in op true (ervan uitgaande dat dit veld bestaat)
			// participant.EmailVerified = true
			// if err := s.participantRepo.Update(ctx, participant); err != nil {
			// 	logger.Error("Fout bij updaten participant verificatie status", "user_id", verificationToken.UserID, "error", err)
			// 	return errors.New("kon verificatie status niet updaten")
			// }

			logger.Info("Participant email verificatie succesvol", "user_id", verificationToken.UserID, "email", verificationToken.Email)
		}
	}

	// Markeer token als gebruikt
	if err := s.emailVerificationTokenRepo.MarkAsUsed(ctx, token); err != nil {
		logger.Error("Fout bij markeren token als gebruikt", "token", token[:8]+"...", "error", err)
		// Dit is niet kritisch, continue
	}

	logger.Info("Email verification successful", "email", verificationToken.Email, "user_type", verificationToken.UserType)
	return nil
}

// ResendEmailVerification verzendt een nieuwe verificatie email
func (s *AuthServiceImpl) ResendEmailVerification(ctx context.Context, email, userID, userType string) error {
	logger.Info("Resend email verification request", "email", email, "user_id", userID, "user_type", userType)

	// Invalideer bestaande tokens voor deze gebruiker
	if err := s.emailVerificationTokenRepo.InvalidateUserTokens(ctx, userID, userType); err != nil {
		logger.Error("Fout bij invalideren bestaande tokens", "user_id", userID, "error", err)
		// Continue anyway
	}

	// Verstuur nieuwe verificatie email
	return s.SendEmailVerification(ctx, email, userID, userType)
}

// IsEmailVerified controleert of een email adres is geverifieerd
func (s *AuthServiceImpl) IsEmailVerified(ctx context.Context, userID, userType string) (bool, error) {
	logger.Debug("Checking email verification status", "user_id", userID, "user_type", userType)

	// Controleer of er gebruikte verificatie tokens bestaan voor deze gebruiker
	tokens, err := s.emailVerificationTokenRepo.GetByUserID(ctx, userID, userType)
	if err != nil {
		logger.Error("Fout bij ophalen verificatie tokens", "user_id", userID, "error", err)
		return false, err
	}

	// Als er minstens één gebruikte token is, is de email geverifieerd
	for _, token := range tokens {
		if token.IsUsed {
			return true, nil
		}
	}

	// Voorlopig: controleer ook user status (dit zou later vervangen kunnen worden door een dedicated email_verified veld)
	// Opgelost: Gebruik nu een tagged switch op userType (QF1003)
	switch userType {
	case "gebruiker":
		gebruiker, err := s.gebruikerRepo.GetByID(ctx, userID)
		if err != nil {
			return false, err
		}
		// Admin/staff accounts worden als geverifieerd beschouwd
		return gebruiker != nil && gebruiker.IsActief, nil
	case "participant":
		participant, err := s.participantRepo.GetByID(ctx, userID)
		if err != nil {
			return false, err
		}
		// Participants met app access worden als geverifieerd beschouwd (voorlopig)
		return participant != nil && participant.HasAppAccess, nil
	default:
		return false, nil
	}
}

// GetGebruikerByEmail haalt een gebruiker op basis van email adres
func (s *AuthServiceImpl) GetGebruikerByEmail(ctx context.Context, email string) (*models.Gebruiker, error) {
	return s.gebruikerRepo.GetByEmail(ctx, email)
}

// GetParticipantByEmail haalt participants op basis van email adres
func (s *AuthServiceImpl) GetParticipantByEmail(ctx context.Context, email string) ([]*models.Participant, error) {
	if s.participantRepo == nil {
		return nil, errors.New("participant repository niet beschikbaar")
	}
	return s.participantRepo.FindByEmail(ctx, email)
}

// ValidateAccessToken controleert of een access token geldig is in de database
func (s *AuthServiceImpl) ValidateAccessToken(ctx context.Context, token string) error {
	if s.accessTokenRepo == nil {
		logger.Warn("Access token repository niet beschikbaar, token validatie overgeslagen")
		return nil // Graceful degradation
	}

	accessToken, err := s.accessTokenRepo.GetByToken(ctx, token)
	if err != nil {
		logger.Error("Fout bij ophalen access token uit database", "error", err)
		return errors.New("kon access token niet valideren")
	}

	if accessToken == nil {
		logger.Warn("Access token niet gevonden in database")
		return errors.New("access token niet gevonden")
	}

	if accessToken.IsRevoked {
		logger.Warn("Access token is ingetrokken")
		return errors.New("access token ingetrokken")
	}

	if !accessToken.IsValid() {
		logger.Warn("Access token is niet meer geldig")
		return errors.New("access token niet meer geldig")
	}

	logger.Debug("Access token validatie succesvol")
	return nil
}

// RevokeAllUserAccessTokens trekt alle access tokens van een gebruiker in
func (s *AuthServiceImpl) RevokeAllUserAccessTokens(ctx context.Context, userID string) error {
	if s.accessTokenRepo == nil {
		logger.Warn("Access token repository niet beschikbaar, tokens worden niet ingetrokken")
		return nil // Graceful degradation
	}

	if err := s.accessTokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		logger.Error("Fout bij revoken alle access tokens", "user_id", userID, "error", err)
		return err
	}

	logger.Info("Alle access tokens ingetrokken", "user_id", userID)
	return nil
}

// storeAccessToken slaat een access token op in de database voor server-side validatie
func (s *AuthServiceImpl) storeAccessToken(ctx context.Context, token, ownerID string) error {
	if s.accessTokenRepo == nil {
		logger.Warn("Access token repository niet beschikbaar, token wordt niet opgeslagen")
		return nil // Graceful degradation
	}

	accessToken := &models.AccessToken{
		OwnerID:   ownerID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenExpiry),
		IsRevoked: false,
	}

	if err := s.accessTokenRepo.Create(ctx, accessToken); err != nil {
		logger.Error("Fout bij opslaan access token in database", "owner_id", ownerID, "error", err)
		return err
	}

	logger.Debug("Access token opgeslagen in database", "owner_id", ownerID)
	return nil
}

// CreateSessionForLogin maakt een nieuwe sessie aan voor een login
func (s *AuthServiceImpl) CreateSessionForLogin(ctx context.Context, userID, userType, accessToken string, ipAddress, userAgent string) (*models.Session, error) {
	if s.sessionRepo == nil {
		logger.Warn("Session repository niet beschikbaar, sessie wordt niet aangemaakt")
		return nil, nil // Graceful degradation
	}

	session := &models.Session{
		OwnerID:      userID,
		AccessToken:  accessToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		DeviceInfo:   s.extractDeviceInfo(userAgent),
		IsActive:     true,
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(s.tokenExpiry),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.Error("Fout bij aanmaken sessie", "user_id", userID, "error", err)
		return nil, err
	}

	logger.Info("Nieuwe sessie aangemaakt", "session_id", session.ID, "user_id", userID, "user_type", userType)
	return session, nil
}

// extractDeviceInfo haalt device informatie op uit de user agent string
func (s *AuthServiceImpl) extractDeviceInfo(userAgent string) models.DeviceInfo {
	// Eenvoudige device detectie gebaseerd op user agent
	userAgentLower := strings.ToLower(userAgent)

	var deviceType string
	if strings.Contains(userAgentLower, "mobile") ||
		strings.Contains(userAgentLower, "android") ||
		strings.Contains(userAgentLower, "iphone") {
		deviceType = "mobile"
	} else if strings.Contains(userAgentLower, "tablet") ||
		strings.Contains(userAgentLower, "ipad") {
		deviceType = "tablet"
	} else {
		deviceType = "desktop"
	}

	// Parse browser info (simplified)
	var browser, browserVersion string
	if strings.Contains(userAgentLower, "chrome") {
		browser = "Chrome"
	} else if strings.Contains(userAgentLower, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(userAgentLower, "safari") && !strings.Contains(userAgentLower, "chrome") {
		browser = "Safari"
	} else if strings.Contains(userAgentLower, "edge") {
		browser = "Edge"
	} else {
		browser = "Unknown"
	}

	// Parse OS info (simplified)
	var os, platform string
	if strings.Contains(userAgentLower, "windows") {
		os = "Windows"
		platform = "Windows"
	} else if strings.Contains(userAgentLower, "mac os x") || strings.Contains(userAgentLower, "macos") {
		os = "macOS"
		platform = "macOS"
	} else if strings.Contains(userAgentLower, "linux") {
		os = "Linux"
		platform = "Linux"
	} else if strings.Contains(userAgentLower, "android") {
		os = "Android"
		platform = "Android"
	} else if strings.Contains(userAgentLower, "ios") || strings.Contains(userAgentLower, "iphone") || strings.Contains(userAgentLower, "ipad") {
		os = "iOS"
		platform = "iOS"
	} else {
		os = "Unknown"
		platform = "Unknown"
	}

	return models.DeviceInfo{
		Browser:        browser,
		BrowserVersion: browserVersion,
		OS:             os,
		OSVersion:      "", // Would need more complex parsing
		DeviceType:     deviceType,
		Platform:       platform,
	}
}

// UpdateAccessTokenWithSession werkt een access token bij met session informatie
func (s *AuthServiceImpl) UpdateAccessTokenWithSession(ctx context.Context, token string, sessionID string) error {
	if s.accessTokenRepo == nil {
		logger.Warn("Access token repository niet beschikbaar, token wordt niet bijgewerkt met session")
		return nil // Graceful degradation
	}

	// Update de access token met session ID
	if err := s.accessTokenRepo.UpdateSessionID(ctx, token, sessionID); err != nil {
		logger.Error("Fout bij updaten access token met session ID", "token", token[:8]+"...", "session_id", sessionID, "error", err)
		return err
	}

	logger.Debug("Access token bijgewerkt met session ID", "token", token[:8]+"...", "session_id", sessionID)
	return nil
}

// DeleteUserAccount verwijdert een gebruikersaccount volledig (GDPR compliance)
// Ondersteunt zowel admin/staff gebruikers als participants
func (s *AuthServiceImpl) DeleteUserAccount(ctx context.Context, userID, password, reason string) error {
	logger.Info("Account deletion request", "user_id", userID)

	// Controleer of gebruiker bestaat (admin/staff)
	gebruiker, err := s.gebruikerRepo.GetByID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Fout bij ophalen gebruiker voor deletion", "user_id", userID, "error", err)
		return errors.New("kon gebruiker niet controleren")
	}

	var userType string
	var userEmail string

	if gebruiker != nil && gebruiker.IsActief {
		// Admin/staff gebruiker gevonden
		userType = "gebruiker"
		userEmail = gebruiker.Email

		// Verifieer wachtwoord voor admin/staff gebruikers
		if !s.VerifyPassword(gebruiker.WachtwoordHash, password) {
			logger.Warn("Ongeldig wachtwoord bij account deletion", "user_id", userID)
			return errors.New("ongeldige inloggegevens")
		}
	} else {
		// Check of het een participant is
		if s.participantRepo != nil {
			participant, err := s.participantRepo.GetByID(ctx, userID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error("Fout bij ophalen participant voor deletion", "user_id", userID, "error", err)
				return errors.New("kon participant niet controleren")
			}

			if participant != nil && participant.HasAppAccess && participant.WachtwoordHash != nil {
				// Participant gevonden
				userType = "participant"
				userEmail = participant.Email

				// Verifieer wachtwoord voor participants
				if !s.VerifyPassword(*participant.WachtwoordHash, password) {
					logger.Warn("Ongeldig wachtwoord bij participant account deletion", "user_id", userID)
					return errors.New("ongeldige inloggegevens")
				}
			} else {
				logger.Warn("Geen actief account gevonden voor deletion", "user_id", userID)
				return errors.New("gebruiker niet gevonden")
			}
		} else {
			logger.Warn("Geen actief account gevonden voor deletion", "user_id", userID)
			return errors.New("gebruiker niet gevonden")
		}
	}

	// Audit log de deletion request
	logger.Audit(ctx, logger.AuditEvent{
		EventType:  logger.AuditAccountDeletion,
		ActorID:    userID,
		ActorEmail: userEmail,
		IPAddress:  "", // Wordt gevuld door caller
		UserAgent:  "", // Wordt gevuld door caller
		Result:     logger.ResultSuccess,
		Metadata: map[string]interface{}{
			"user_type": userType,
			"reason":    reason,
		},
	})

	// Verwijder alle gerelateerde data afhankelijk van user type
	if userType == "gebruiker" {
		// Admin/staff gebruiker: verwijder alle gerelateerde data
		if err := s.deleteGebruikerData(ctx, userID); err != nil {
			logger.Error("Fout bij verwijderen gebruiker data", "user_id", userID, "error", err)
			return errors.New("kon gebruiker data niet verwijderen")
		}
	} else {
		// Participant: verwijder alle gerelateerde data
		if err := s.deleteParticipantData(ctx, userID); err != nil {
			logger.Error("Fout bij verwijderen participant data", "user_id", userID, "error", err)
			return errors.New("kon participant data niet verwijderen")
		}
	}

	logger.Info("Account succesvol verwijderd", "user_id", userID, "user_type", userType, "email", userEmail)
	return nil
}

// deleteGebruikerData verwijdert alle data gerelateerd aan een admin/staff gebruiker
func (s *AuthServiceImpl) deleteGebruikerData(ctx context.Context, userID string) error {
	// 1. Verwijder refresh tokens
	if s.refreshTokenRepo != nil {
		if err := s.refreshTokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
			logger.Error("Fout bij verwijderen refresh tokens", "user_id", userID, "error", err)
			return err
		}
	}

	// 2. Verwijder access tokens
	if s.accessTokenRepo != nil {
		if err := s.accessTokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
			logger.Error("Fout bij verwijderen access tokens", "user_id", userID, "error", err)
			return err
		}
	}

	// 3. Verwijder user roles
	if s.userRoleRepo != nil {
		if err := s.userRoleRepo.DeleteByUser(ctx, userID); err != nil {
			logger.Error("Fout bij verwijderen user roles", "user_id", userID, "error", err)
			return err
		}
	}

	// 4. Verwijder de gebruiker zelf (dit is de belangrijkste stap voor GDPR)
	if err := s.gebruikerRepo.Delete(ctx, userID); err != nil {
		logger.Error("Fout bij verwijderen gebruiker record", "user_id", userID, "error", err)
		return err
	}

	return nil
}

// deleteParticipantData verwijdert alle data gerelateerd aan een participant
func (s *AuthServiceImpl) deleteParticipantData(ctx context.Context, participantID string) error {
	// 1. Verwijder refresh tokens
	if s.refreshTokenRepo != nil {
		if err := s.refreshTokenRepo.RevokeAllUserTokens(ctx, participantID); err != nil {
			logger.Error("Fout bij verwijderen participant refresh tokens", "participant_id", participantID, "error", err)
			return err
		}
	}

	// 2. Verwijder access tokens
	if s.accessTokenRepo != nil {
		if err := s.accessTokenRepo.RevokeAllUserTokens(ctx, participantID); err != nil {
			logger.Error("Fout bij verwijderen participant access tokens", "participant_id", participantID, "error", err)
			return err
		}
	}

	// 3. Verwijder de participant zelf (dit is de belangrijkste stap voor GDPR)
	if s.participantRepo != nil {
		if err := s.participantRepo.Delete(ctx, participantID); err != nil {
			logger.Error("Fout bij verwijderen participant record", "participant_id", participantID, "error", err)
			return err
		}
	}

	return nil
}
