package handlers

import (
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler bevat handlers voor authenticatie
type AuthHandler struct {
	authService       services.AuthService
	permissionService services.PermissionService
	rateLimiter       services.RateLimiterService
}

// NewAuthHandler maakt een nieuwe AuthHandler
func NewAuthHandler(authService services.AuthService, permissionService services.PermissionService, rateLimiter services.RateLimiterService) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		permissionService: permissionService,
		rateLimiter:       rateLimiter,
	}
}

// HandleLogin godoc
// @Summary User login
// @Description Authenticates a user and returns access/refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.GebruikerLogin true "Login credentials"
// @Success 200 {object} models.AuthLoginResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 403 {object} object
// @Failure 429 {object} object
// @Failure 500 {object} object
// @Router /api/auth/login [post]
func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	// Parse request body
	var loginData models.GebruikerLogin
	if err := c.BodyParser(&loginData); err != nil {
		logger.Error("Fout bij parsen login data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige login data",
		})
	}

	// Valideer input
	if loginData.Email == "" || loginData.Wachtwoord == "" {
		logger.Warn("Ontbrekende login gegevens")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email en wachtwoord zijn verplicht",
		})
	}

	// Rate limiting voor login pogingen
	rateLimitKey := "login:" + loginData.Email
	if !h.rateLimiter.Allow(rateLimitKey) {
		logger.Warn("Rate limit overschreden voor login", "email", loginData.Email)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Te veel login pogingen, probeer het later opnieuw",
		})
	}

	// Authenticeer gebruiker
	token, refreshToken, err := h.authService.Login(c.Context(), loginData.Email, loginData.Wachtwoord)
	if err != nil {
		// Audit: Failed login
		logger.Audit(c.Context(), logger.AuditEvent{
			EventType:  logger.AuditLoginFailed,
			ActorEmail: loginData.Email,
			IPAddress:  c.IP(),
			UserAgent:  c.Get("User-Agent"),
			Result:     logger.ResultFailed,
			Reason:     err.Error(),
		})

		// Specifieke foutafhandeling
		switch err {
		case services.ErrInvalidCredentials:
			logger.Warn("Ongeldige inloggegevens", "email", loginData.Email, "ip", c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Ongeldige inloggegevens",
				"code":  "INVALID_CREDENTIALS",
			})
		case services.ErrUserInactive:
			logger.Warn("Inactieve gebruiker", "email", loginData.Email, "ip", c.IP())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Gebruiker is inactief",
				"code":  "USER_INACTIVE",
			})
		default:
			logger.Error("Fout bij login", "email", loginData.Email, "ip", c.IP(), "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Er is een fout opgetreden bij het inloggen",
				"code":  "LOGIN_ERROR",
			})
		}
	}

	// V34: Valideer token en haal claims op om user type te bepalen
	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		logger.Error("Fout bij valideren token na login", "email", loginData.Email, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Login succesvol maar token validatie gefaald",
		})
	}

	// Probeer eerst gebruiker op te halen (admin/staff)
	gebruiker, err := h.authService.GetUser(c.Context(), userID)

	// Als gebruiker niet gevonden, check of het een participant is
	if err != nil || gebruiker == nil {
		// V34: Return participant login response (met session creation)
		return h.handleParticipantLoginResponse(c, token, refreshToken, userID, loginData.Email)
	}

	// Standard gebruiker login response (admin/staff) met session creation
	return h.handleGebruikerLoginResponse(c, token, refreshToken, gebruiker)
}

// handleParticipantLoginResponse creates response for participant logins (V34)
func (h *AuthHandler) handleParticipantLoginResponse(c *fiber.Ctx, token, refreshToken, participantID, email string) error {
	// Stel cookie in met token - verbeterde security check voor productie
	// Controleer X-Forwarded-Proto header (voor proxy setups) of gebruik env var
	isSecure := false
	if forwardedProto := c.Get("X-Forwarded-Proto"); forwardedProto == "https" {
		isSecure = true
	} else if os.Getenv("COOKIE_SECURE") == "true" {
		isSecure = true
	} else {
		// Fallback naar protocol check (voor directe verbindingen)
		isSecure = c.Protocol() == "https"
	}

	cookie := fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(20 * time.Minute),
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: "Strict",
	}
	c.Cookie(&cookie)

	// OPLOSSING 2: Gebruik de centrale helper voor een consistente permissielijst
	permissionList := services.GetParticipantPermissionsList()

	// Participants krijgen participant_user rol
	roleList := []map[string]interface{}{
		{
			"id":          "participant-role",
			"name":        "participant_user",
			"description": "Participant with app access",
		},
	}

	// Audit: Successful participant login
	logger.Audit(c.Context(), logger.AuditEvent{
		EventType:  logger.AuditLoginSuccess,
		ActorID:    participantID,
		ActorEmail: email,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		Result:     logger.ResultSuccess,
		Metadata: map[string]interface{}{
			"user_type":         "participant",
			"roles_count":       len(roleList),
			"permissions_count": len(permissionList),
		},
	})

	// Return participant login response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"token":         token,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":          participantID,
			"email":       email,
			"naam":        "",             // Frontend can call /api/participants/me to get full data
			"permissions": permissionList, // <-- Nu 100% consistent
			"roles":       roleList,
			"is_actief":   true,
		},
	})
}

// handleGebruikerLoginResponse creates response for gebruiker logins (admin/staff)
func (h *AuthHandler) handleGebruikerLoginResponse(c *fiber.Ctx, token, refreshToken string, gebruiker *models.Gebruiker) error {
	// Maak sessie aan voor deze login
	session, err := h.authService.CreateSessionForLogin(c.Context(), gebruiker.ID, "gebruiker", token, c.IP(), c.Get("User-Agent"))
	if err != nil {
		logger.Error("Fout bij aanmaken sessie voor gebruiker login", "user_id", gebruiker.ID, "error", err)
		// Continue met login, sessie is niet kritisch
	}

	// Update access token met session ID als sessie is aangemaakt
	if session != nil {
		if err := h.authService.UpdateAccessTokenWithSession(c.Context(), token, session.ID); err != nil {
			logger.Error("Fout bij updaten access token met session ID", "token", token[:8]+"...", "session_id", session.ID, "error", err)
			// Continue, dit is niet kritisch
		}
	}

	// Haal permissies op
	permissions, err := h.permissionService.GetUserPermissions(c.Context(), gebruiker.ID)
	if err != nil {
		logger.Error("Fout bij ophalen permissies na login", "user_id", gebruiker.ID, "error", err)
		permissions = []*models.UserPermission{} // Fallback naar lege array
	}

	// Converteer permissies naar frontend format
	permissionList := make([]map[string]string, len(permissions))
	for i, perm := range permissions {
		permissionList[i] = map[string]string{
			"resource": perm.Resource,
			"action":   perm.Action,
		}
	}

	// Stel cookie in met token (20 minuten expiry) - verbeterde security check voor productie
	// Controleer X-Forwarded-Proto header (voor proxy setups) of gebruik env var
	isSecure := false
	if forwardedProto := c.Get("X-Forwarded-Proto"); forwardedProto == "https" {
		isSecure = true
	} else if os.Getenv("COOKIE_SECURE") == "true" {
		isSecure = true
	} else {
		// Fallback naar protocol check (voor directe verbindingen)
		isSecure = c.Protocol() == "https"
	}

	cookie := fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(20 * time.Minute),
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: "Strict",
	}
	c.Cookie(&cookie)

	// Haal RBAC roles op voor response
	userRoles, err := h.permissionService.GetUserRoles(c.Context(), gebruiker.ID)
	if err != nil {
		logger.Error("Fout bij ophalen rollen na login", "user_id", gebruiker.ID, "error", err)
		userRoles = []*models.UserRole{} // Fallback naar lege array
	}

	// Converteer rollen naar frontend format
	roleList := make([]map[string]interface{}, len(userRoles))
	for i, userRole := range userRoles {
		roleList[i] = map[string]interface{}{
			"id":          userRole.Role.ID,
			"name":        userRole.Role.Name,
			"description": userRole.Role.Description,
		}
	}

	// Audit: Successful login
	logger.Audit(c.Context(), logger.AuditEvent{
		EventType:  logger.AuditLoginSuccess,
		ActorID:    gebruiker.ID,
		ActorEmail: gebruiker.Email,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		Result:     logger.ResultSuccess,
		Metadata: map[string]interface{}{
			"user_type":         "gebruiker",
			"roles_count":       len(roleList),
			"permissions_count": len(permissionList),
		},
	})

	// Stuur complete user data terug met refresh token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"token":         token,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":          gebruiker.ID,
			"email":       gebruiker.Email,
			"naam":        gebruiker.Naam,
			"permissions": permissionList,
			"roles":       roleList,
			"is_actief":   gebruiker.IsActief,
			// DEPRECATED: rol field removed - use roles array instead
		},
	})
}

// HandleRefreshToken godoc
// @Summary Refresh access token
// @Description Refreshes access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param refreshData body models.AuthRefreshRequest true "Refresh token data"
// @Success 200 {object} models.AuthRefreshResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/auth/refresh [post]
func (h *AuthHandler) HandleRefreshToken(c *fiber.Ctx) error {
	var refreshData struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Controleer of request body leeg is (vaak door automatische frontend calls)
	body := c.Body()
	if len(body) == 0 {
		logger.Debug("Lege request body ontvangen bij refresh token endpoint", "ip", c.IP(), "user_agent", c.Get("User-Agent"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token data ontbreekt",
			"code":  "MISSING_REQUEST_BODY",
		})
	}

	if err := c.BodyParser(&refreshData); err != nil {
		logger.Warn("Fout bij parsen refresh token data", "error", err, "body_length", len(body), "ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige refresh token data",
			"code":  "INVALID_JSON",
		})
	}

	if refreshData.RefreshToken == "" {
		logger.Warn("Ontbrekende refresh token in request", "ip", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is verplicht",
			"code":  "MISSING_REFRESH_TOKEN",
		})
	}

	// Refresh tokens
	accessToken, newRefreshToken, err := h.authService.RefreshAccessToken(c.Context(), refreshData.RefreshToken)
	if err != nil {
		logger.Warn("Token refresh gefaald", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Ongeldige of verlopen refresh token",
			"code":  "REFRESH_TOKEN_INVALID",
		})
	}

	logger.Info("Token refresh succesvol")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"token":         accessToken,
		"refresh_token": newRefreshToken,
	})
}

// HandleLogout godoc
// @Summary User logout
// @Description Logs out the current user and optionally revokes all sessions
// @Tags Authentication
// @Accept json
// @Produce json
// @Param options body object false "Logout options"
// @Success 200 {object} object
// @Failure 401 {object} object
// @Security BearerAuth
// @Router /api/auth/logout [post]
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	// Haal user ID op als beschikbaar
	userID, _ := c.Locals("userID").(string)

	// Parse request body voor optionele parameters
	var logoutData struct {
		RevokeAllSessions bool   `json:"revoke_all_sessions,omitempty"` // Default: false (alleen huidige sessie)
		SessionID         string `json:"session_id,omitempty"`          // Specifieke sessie om in te trekken
	}

	// Parse request body (optioneel)
	if err := c.BodyParser(&logoutData); err != nil {
		// Als body parsing faalt, ga door met standaard logout (geen sessie intrekking)
		logger.Debug("Kon logout request body niet parsen, gebruik standaard logout", "error", err)
	}

	// Audit: Logout
	if userID != "" {
		logger.Audit(c.Context(), logger.AuditEvent{
			EventType: logger.AuditLogout,
			ActorID:   userID,
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Result:    logger.ResultSuccess,
			Metadata: map[string]interface{}{
				"revoke_all_sessions": logoutData.RevokeAllSessions,
				"specific_session":    logoutData.SessionID != "",
			},
		})
	}

	// Sessie management gebaseerd op request parameters
	if userID != "" {
		if logoutData.RevokeAllSessions {
			// Trek alle sessies van de gebruiker in
			if err := h.authService.RevokeAllUserSessions(c.Context(), userID); err != nil {
				logger.Error("Fout bij intrekken alle sessies tijdens logout", "user_id", userID, "error", err)
				// Continue met logout, dit is niet kritisch
			}
			logger.Info("Alle sessies ingetrokken tijdens logout", "user_id", userID)
		} else if logoutData.SessionID != "" {
			// Trek specifieke sessie in
			if err := h.authService.RevokeSession(c.Context(), logoutData.SessionID); err != nil {
				logger.Error("Fout bij intrekken specifieke sessie tijdens logout", "user_id", userID, "session_id", logoutData.SessionID, "error", err)
				// Continue met logout, dit is niet kritisch
			}
			logger.Info("Specifieke sessie ingetrokken tijdens logout", "user_id", userID, "session_id", logoutData.SessionID)
		} else {
			// FIX: Standaard gedrag - trek ALLEEN de huidige sessie in
			// We gebruiken de sessionID die door SessionValidationMiddleware in de context is gezet
			currentSessionID, ok := c.Locals("sessionID").(string)

			if ok && currentSessionID != "" {
				if err := h.authService.RevokeSession(c.Context(), currentSessionID); err != nil {
					logger.Error("Fout bij intrekken huidige sessie tijdens logout", "user_id", userID, "session_id", currentSessionID, "error", err)
					// Continue
				} else {
					logger.Info("Huidige sessie ingetrokken tijdens logout", "user_id", userID, "session_id", currentSessionID)
				}
			} else {
				// Fallback: Als we geen sessionID hebben (bijv. middleware issue), doen we GEEN RevokeAllUserAccessTokens meer.
				// Dat was te agressief. We loggen alleen dat server-side intrekking niet gelukt is.
				logger.Warn("Kon sessie niet intrekken: geen sessionID gevonden in context", "user_id", userID)
			}
		}
	}

	// Verwijder cookie
	c.ClearCookie("auth_token")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout succesvol",
	})
}

// HandleResetPassword godoc
// @Summary Change password
// @Description Changes the current user's password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.AuthResetPasswordRequest true "Current and new password"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Security BearerAuth
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) HandleResetPassword(c *fiber.Ctx) error {
	// Alleen toegankelijk voor ingelogde gebruikers
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen gebruiker ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	// Parse request body
	var resetData struct {
		HuidigWachtwoord string `json:"huidig_wachtwoord"`
		NieuwWachtwoord  string `json:"nieuw_wachtwoord"`
	}
	if err := c.BodyParser(&resetData); err != nil {
		logger.Error("Fout bij parsen wachtwoord reset data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige wachtwoord reset data",
		})
	}

	// Valideer input
	if resetData.HuidigWachtwoord == "" || resetData.NieuwWachtwoord == "" {
		logger.Warn("Ontbrekende wachtwoord reset gegevens")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Huidig wachtwoord en nieuw wachtwoord zijn verplicht",
		})
	}

	// Haal gebruiker op
	gebruiker, err := h.authService.GetUserFromToken(c.Context(), c.Locals("token").(string))
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het ophalen van de gebruiker",
		})
	}

	// Verifieer huidig wachtwoord
	if !h.authService.VerifyPassword(gebruiker.WachtwoordHash, resetData.HuidigWachtwoord) {
		logger.Warn("Ongeldig huidig wachtwoord", "user_id", userID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Ongeldig huidig wachtwoord",
		})
	}

	// Reset wachtwoord
	if err := h.authService.ResetPassword(c.Context(), gebruiker.Email, resetData.NieuwWachtwoord); err != nil {
		logger.Error("Fout bij resetten wachtwoord", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het resetten van het wachtwoord",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Wachtwoord succesvol gewijzigd",
	})
}

// HandleGetProfile godoc
// @Summary Get user profile
// @Description Gets the current user's profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} models.AuthProfileResponse
// @Failure 401 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Security BearerAuth
// @Router /api/auth/profile [get]
func (h *AuthHandler) HandleGetProfile(c *fiber.Ctx) error {
	// Haal user ID op uit context (gezet door AuthMiddleware)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen user ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
		})
	}

	// Haal gebruiker op
	gebruiker, err := h.authService.GetUser(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen gebruiker", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon gebruiker niet ophalen",
		})
	}

	if gebruiker == nil {
		logger.Warn("Gebruiker niet gevonden", "user_id", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Gebruiker niet gevonden",
		})
	}

	// Haal permissies op via RBAC systeem
	permissions, err := h.permissionService.GetUserPermissions(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen permissies", "user_id", userID, "error", err)
		// Fallback naar lege array als permissies niet opgehaald kunnen worden
		permissions = []*models.UserPermission{}
	}

	// Converteer permissies naar frontend format
	permissionList := make([]map[string]string, len(permissions))
	for i, perm := range permissions {
		permissionList[i] = map[string]string{
			"resource": perm.Resource,
			"action":   perm.Action,
		}
	}

	// Haal rollen op
	userRoles, err := h.permissionService.GetUserRoles(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen rollen", "user_id", userID, "error", err)
		userRoles = []*models.UserRole{}
	}

	// Converteer rollen naar frontend format
	roleList := make([]map[string]interface{}, len(userRoles))
	for i, userRole := range userRoles {
		roleList[i] = map[string]interface{}{
			"id":          userRole.Role.ID,
			"name":        userRole.Role.Name,
			"description": userRole.Role.Description,
			"assigned_at": userRole.AssignedAt,
			"is_active":   userRole.IsActive,
		}
	}

	// Stuur gebruikersprofiel terug met permissies en rollen
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            gebruiker.ID,
		"naam":          gebruiker.Naam,
		"email":         gebruiker.Email,
		"permissions":   permissionList,
		"roles":         roleList,
		"is_actief":     gebruiker.IsActief,
		"laatste_login": gebruiker.LaatsteLogin,
		"created_at":    gebruiker.CreatedAt,
		// DEPRECATED: rol field removed - use roles array instead
	})
}

// HandleForgotPassword handelt wachtwoord vergeten verzoeken af
func (h *AuthHandler) HandleForgotPassword(c *fiber.Ctx) error {
	// Parse request body
	var forgotData struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&forgotData); err != nil {
		logger.Error("Fout bij parsen forgot password data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Valideer input
	if forgotData.Email == "" {
		logger.Warn("Ontbrekend email adres bij forgot password")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email adres is verplicht",
		})
	}

	// Rate limiting voor forgot password requests
	rateLimitKey := "forgot_password:" + forgotData.Email
	if !h.rateLimiter.Allow(rateLimitKey) {
		logger.Warn("Rate limit overschreden voor forgot password", "email", forgotData.Email)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Te veel wachtwoord reset verzoeken, probeer het later opnieuw",
		})
	}

	// Verzoek wachtwoord reset
	if err := h.authService.RequestPasswordReset(c.Context(), forgotData.Email); err != nil {
		logger.Error("Fout bij aanvragen wachtwoord reset", "email", forgotData.Email, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het verwerken van je verzoek",
		})
	}

	// Altijd succes teruggeven voor security (geen email enumeration)
	logger.Info("Password reset request processed", "email", forgotData.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Als er een account bestaat met dit email adres, ontvang je binnenkort een email met instructies om je wachtwoord te resetten.",
	})
}

// HandleResetPasswordWithToken handelt wachtwoord reset met token verzoeken af
func (h *AuthHandler) HandleResetPasswordWithToken(c *fiber.Ctx) error {
	// Parse request body
	var resetData struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&resetData); err != nil {
		logger.Error("Fout bij parsen reset password data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Valideer input
	if resetData.Token == "" || resetData.NewPassword == "" {
		logger.Warn("Ontbrekende token of wachtwoord bij reset")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token en nieuw wachtwoord zijn verplicht",
		})
	}

	// Controleer wachtwoord sterkte
	if len(resetData.NewPassword) < 8 {
		logger.Warn("Wachtwoord te kort bij reset")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Wachtwoord moet minimaal 8 karakters bevatten",
		})
	}

	// Reset wachtwoord met token
	if err := h.authService.ResetPasswordWithToken(c.Context(), resetData.Token, resetData.NewPassword); err != nil {
		logger.Warn("Fout bij resetten wachtwoord met token", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige of verlopen reset token",
		})
	}

	logger.Info("Password reset successful with token")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Je wachtwoord is succesvol gewijzigd. Je kunt nu inloggen met je nieuwe wachtwoord.",
	})
}

// HandleSendEmailVerification handelt verzoeken voor het verzenden van email verificatie af
func (h *AuthHandler) HandleSendEmailVerification(c *fiber.Ctx) error {
	// Parse request body
	var verificationData struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&verificationData); err != nil {
		logger.Error("Fout bij parsen email verification data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Valideer input
	if verificationData.Email == "" {
		logger.Warn("Ontbrekend email adres bij email verification")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email adres is verplicht",
		})
	}

	// Controleer of gebruiker bestaat en haal user info op
	var userID, userType string

	// Check admin/staff users eerst
	if gebruiker, err := h.authService.GetGebruikerByEmail(c.Context(), verificationData.Email); err == nil && gebruiker != nil && gebruiker.IsActief {
		userID = gebruiker.ID
		userType = "gebruiker"
	} else {
		// Check participants
		if participants, err := h.authService.GetParticipantByEmail(c.Context(), verificationData.Email); err == nil {
			for _, participant := range participants {
				if participant.AccountType == "full" && participant.HasAppAccess && participant.WachtwoordHash != nil {
					userID = participant.ID
					userType = "participant"
					break
				}
			}
		}
	}

	if userID == "" {
		logger.Warn("Geen gebruiker gevonden voor email verification", "email", verificationData.Email)
		// Voor security: geef geen foutmelding terug (geen email enumeration)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Als er een account bestaat met dit email adres, ontvang je binnenkort een verificatie email.",
		})
	}

	// Rate limiting voor email verification requests
	rateLimitKey := "email_verification:" + verificationData.Email
	if !h.rateLimiter.Allow(rateLimitKey) {
		logger.Warn("Rate limit overschreden voor email verification", "email", verificationData.Email)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Te veel verificatie verzoeken, probeer het later opnieuw",
		})
	}

	// Verstuur verificatie email
	if err := h.authService.SendEmailVerification(c.Context(), verificationData.Email, userID, userType); err != nil {
		logger.Error("Fout bij verzenden email verificatie", "email", verificationData.Email, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het verzenden van de verificatie email",
		})
	}

	logger.Info("Email verification request processed", "email", verificationData.Email, "user_type", userType)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Als er een account bestaat met dit email adres, ontvang je binnenkort een verificatie email.",
	})
}

// HandleVerifyEmail handelt email verificatie verzoeken af
func (h *AuthHandler) HandleVerifyEmail(c *fiber.Ctx) error {
	// Parse request body
	var verificationData struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&verificationData); err != nil {
		logger.Error("Fout bij parsen email verification data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Valideer input
	if verificationData.Token == "" {
		logger.Warn("Ontbrekende verificatie token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Verificatie token is verplicht",
		})
	}

	// Verificeer email met token
	if err := h.authService.VerifyEmailWithToken(c.Context(), verificationData.Token); err != nil {
		logger.Warn("Fout bij verifiëren email met token", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige of verlopen verificatie token",
		})
	}

	logger.Info("Email verification successful")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Je email adres is succesvol geverifieerd! Je kunt nu volledig gebruik maken van je account.",
	})
}

// HandleResendEmailVerification handelt verzoeken voor het opnieuw verzenden van email verificatie af
func (h *AuthHandler) HandleResendEmailVerification(c *fiber.Ctx) error {
	// Parse request body
	var resendData struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&resendData); err != nil {
		logger.Error("Fout bij parsen resend email verification data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
		})
	}

	// Valideer input
	if resendData.Email == "" {
		logger.Warn("Ontbrekend email adres bij resend email verification")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email adres is verplicht",
		})
	}

	// Controleer of gebruiker bestaat en haal user info op
	var userID, userType string

	// Check admin/staff users eerst
	if gebruiker, err := h.authService.GetGebruikerByEmail(c.Context(), resendData.Email); err == nil && gebruiker != nil && gebruiker.IsActief {
		userID = gebruiker.ID
		userType = "gebruiker"
	} else {
		// Check participants
		if participants, err := h.authService.GetParticipantByEmail(c.Context(), resendData.Email); err == nil {
			for _, participant := range participants {
				if participant.AccountType == "full" && participant.HasAppAccess && participant.WachtwoordHash != nil {
					userID = participant.ID
					userType = "participant"
					break
				}
			}
		}
	}

	if userID == "" {
		logger.Warn("Geen gebruiker gevonden voor resend email verification", "email", resendData.Email)
		// Voor security: geef geen foutmelding terug
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Als er een account bestaat met dit email adres, ontvang je binnenkort een nieuwe verificatie email.",
		})
	}

	// Rate limiting voor resend requests
	rateLimitKey := "resend_email_verification:" + resendData.Email
	if !h.rateLimiter.Allow(rateLimitKey) {
		logger.Warn("Rate limit overschreden voor resend email verification", "email", resendData.Email)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Te veel verzoeken om verificatie email opnieuw te verzenden, probeer het later opnieuw",
		})
	}

	// Verstuur nieuwe verificatie email
	if err := h.authService.ResendEmailVerification(c.Context(), resendData.Email, userID, userType); err != nil {
		logger.Error("Fout bij opnieuw verzenden email verificatie", "email", resendData.Email, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het opnieuw verzenden van de verificatie email",
		})
	}

	logger.Info("Email verification resend request processed", "email", resendData.Email, "user_type", userType)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Als er een account bestaat met dit email adres, ontvang je binnenkort een nieuwe verificatie email.",
	})
}

// HandleDeleteAccount godoc
// @Summary Delete account
// @Description Deletes the current user's account (GDPR compliance)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.AuthDeleteAccountRequest true "Account deletion confirmation"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @Failure 429 {object} object
// @Failure 500 {object} object
// @Security BearerAuth
// @Router /api/auth/account [delete]
func (h *AuthHandler) HandleDeleteAccount(c *fiber.Ctx) error {
	// Haal user ID op uit context (gezet door AuthMiddleware)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen user ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
			"code":  "UNAUTHORIZED",
		})
	}

	// Parse request body
	var deleteData struct {
		Password string `json:"password" binding:"required"`
		Reason   string `json:"reason,omitempty"` // Optioneel voor GDPR compliance
	}

	if err := c.BodyParser(&deleteData); err != nil {
		logger.Error("Fout bij parsen account deletion data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige request data",
			"code":  "INVALID_INPUT",
		})
	}

	// Valideer input
	if deleteData.Password == "" {
		logger.Warn("Ontbrekend wachtwoord bij account deletion")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Wachtwoord is verplicht voor account verwijdering",
			"code":  "MISSING_PASSWORD",
		})
	}

	// Rate limiting voor account deletion requests
	rateLimitKey := "delete_account:" + userID
	if !h.rateLimiter.Allow(rateLimitKey) {
		logger.Warn("Rate limit overschreden voor account deletion", "user_id", userID)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Te veel account verwijdering verzoeken, probeer het later opnieuw",
			"code":  "RATE_LIMIT_EXCEEDED",
		})
	}

	// Verwijder account
	if err := h.authService.DeleteUserAccount(c.Context(), userID, deleteData.Password, deleteData.Reason); err != nil {
		switch err.Error() {
		case "ongeldige inloggegevens":
			logger.Warn("Ongeldig wachtwoord bij account deletion", "user_id", userID)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Ongeldig wachtwoord",
				"code":  "INVALID_PASSWORD",
			})
		case "gebruiker niet gevonden":
			logger.Warn("Gebruiker niet gevonden bij account deletion", "user_id", userID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Account niet gevonden",
				"code":  "ACCOUNT_NOT_FOUND",
			})
		default:
			logger.Error("Fout bij account deletion", "user_id", userID, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Er is een fout opgetreden bij het verwijderen van je account",
				"code":  "DELETION_FAILED",
			})
		}
	}

	logger.Info("Account succesvol verwijderd", "user_id", userID)

	// Logout na succesvolle verwijdering
	c.ClearCookie("auth_token")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Je account is succesvol verwijderd. Alle persoonlijke gegevens zijn permanent verwijderd.",
	})
}

// HandleListSessions godoc
// @Summary List user sessions
// @Description Lists all active sessions for the current user
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Security BearerAuth
// @Router /api/auth/sessions [get]
func (h *AuthHandler) HandleListSessions(c *fiber.Ctx) error {
	// Haal user ID op uit context (gezet door AuthMiddleware)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen user ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
			"code":  "UNAUTHORIZED",
		})
	}

	// Haal sessies op
	sessions, err := h.authService.ListUserSessions(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen sessies", "user_id", userID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het ophalen van sessies",
			"code":  "SESSIONS_ERROR",
		})
	}

	// Converteer naar response format
	sessionResponses := make([]map[string]interface{}, 0, len(sessions))
	for _, session := range sessions {
		sessionResponses = append(sessionResponses, map[string]interface{}{
			"id":            session.ID,
			"device_info":   session.DeviceInfo,
			"ip_address":    session.IPAddress,
			"user_agent":    session.UserAgent,
			"login_time":    session.LoginTime,
			"last_activity": session.LastActivity,
			"is_current":    session.IsCurrent,
			"display_name":  session.GetDisplayName(),
			"location_info": session.GetLocationInfo(),
		})
	}

	logger.Info("Sessies opgehaald", "user_id", userID, "count", len(sessions))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"sessions": sessionResponses,
	})
}

// HandleRevokeSession godoc
// @Summary Revoke specific session
// @Description Revokes a specific user session
// @Tags Authentication
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Security BearerAuth
// @Router /api/auth/sessions/{sessionId} [delete]
func (h *AuthHandler) HandleRevokeSession(c *fiber.Ctx) error {
	// Haal user ID op uit context (gezet door AuthMiddleware)
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen user ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
			"code":  "UNAUTHORIZED",
		})
	}

	sessionID := c.Params("sessionId")
	if sessionID == "" {
		logger.Warn("Geen session ID opgegeven")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is verplicht",
			"code":  "MISSING_SESSION_ID",
		})
	}

	// Controleer of de sessie van de huidige gebruiker is (voor security)
	// Dit zou normaal gesproken in de service gebeuren, maar voor nu doen we een basic check
	sessions, err := h.authService.ListUserSessions(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij controleren sessie ownership", "user_id", userID, "session_id", sessionID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het controleren van de sessie",
			"code":  "SESSION_CHECK_ERROR",
		})
	}

	// Controleer of de sessie bestaat en van de gebruiker is
	sessionExists := false
	for _, session := range sessions {
		if session.ID == sessionID {
			sessionExists = true
			break
		}
	}

	if !sessionExists {
		logger.Warn("Sessie niet gevonden of niet van gebruiker", "user_id", userID, "session_id", sessionID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sessie niet gevonden",
			"code":  "SESSION_NOT_FOUND",
		})
	}

	// Trek de sessie in
	if err := h.authService.RevokeSession(c.Context(), sessionID); err != nil {
		logger.Error("Fout bij intrekken sessie", "user_id", userID, "session_id", sessionID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het intrekken van de sessie",
			"code":  "REVOKE_SESSION_ERROR",
		})
	}

	logger.Info("Sessie ingetrokken", "user_id", userID, "session_id", sessionID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Sessie succesvol ingetrokken",
	})
}

// HandleRevokeOtherSessions godoc
// @Summary Revoke other sessions
// @Description Revokes all other sessions except the current one
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Security BearerAuth
// @Router /api/auth/sessions/revoke-others [post]
func (h *AuthHandler) HandleRevokeOtherSessions(c *fiber.Ctx) error {
	// Haal user ID en huidige sessie op uit context
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		logger.Warn("Geen user ID gevonden in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Niet geautoriseerd",
			"code":  "UNAUTHORIZED",
		})
	}

	// Haal huidige sessie ID op uit access token
	// Dit vereist dat we de access token kunnen koppelen aan een sessie
	// Voor nu gebruiken we een placeholder - dit zou normaal gesproken uit de context komen
	currentSessionID := c.Get("X-Session-ID") // Placeholder - zou uit middleware moeten komen
	if currentSessionID == "" {
		logger.Warn("Geen huidige sessie ID gevonden")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Huidige sessie ID kon niet worden bepaald",
			"code":  "MISSING_CURRENT_SESSION",
		})
	}

	// Trek andere sessies in
	if err := h.authService.RevokeOtherUserSessions(c.Context(), userID, currentSessionID); err != nil {
		logger.Error("Fout bij intrekken andere sessies", "user_id", userID, "current_session_id", currentSessionID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Er is een fout opgetreden bij het intrekken van andere sessies",
			"code":  "REVOKE_OTHERS_ERROR",
		})
	}

	logger.Info("Andere sessies ingetrokken", "user_id", userID, "current_session_id", currentSessionID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Alle andere sessies zijn succesvol ingetrokken",
	})
}
