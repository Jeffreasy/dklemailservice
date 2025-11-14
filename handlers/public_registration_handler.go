package handlers

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// PublicRegistrationHandler bevat handlers voor publieke registratie (zonder authenticatie)
// V30: Implementeert duaal registratiesysteem (full vs temporary accounts)
// V30+RBAC: Integreert met RBAC systeem voor automatische rol toewijzing
type PublicRegistrationHandler struct {
	participantRepo     repository.ParticipantRepository
	eventRegRepo        repository.EventRegistrationRepository
	eventRepo           repository.EventRepository
	gebruikerRepo       repository.GebruikerRepository
	emailService        *services.EmailService
	notificationService services.NotificationService
	permissionService   services.PermissionService // V30+RBAC: Voor rol assignment
	roleRepo            repository.RBACRoleRepository
	userRoleRepo        repository.UserRoleRepository
}

// NewPublicRegistrationHandler maakt een nieuwe public registration handler
func NewPublicRegistrationHandler(
	participantRepo repository.ParticipantRepository,
	eventRegRepo repository.EventRegistrationRepository,
	eventRepo repository.EventRepository,
	gebruikerRepo repository.GebruikerRepository,
	emailService *services.EmailService,
	notificationService services.NotificationService,
	permissionService services.PermissionService,
	roleRepo repository.RBACRoleRepository,
	userRoleRepo repository.UserRoleRepository,
) *PublicRegistrationHandler {
	return &PublicRegistrationHandler{
		participantRepo:     participantRepo,
		eventRegRepo:        eventRegRepo,
		eventRepo:           eventRepo,
		gebruikerRepo:       gebruikerRepo,
		emailService:        emailService,
		notificationService: notificationService,
		permissionService:   permissionService,
		roleRepo:            roleRepo,
		userRoleRepo:        userRoleRepo,
	}
}

// RegisterRoutes registreert de publieke registratie routes (GEEN authenticatie vereist)
func (h *PublicRegistrationHandler) RegisterRoutes(app *fiber.App) {
	// Publieke routes - GEEN AuthMiddleware
	publicApi := app.Group("/api/public")

	// V30: Nieuwe duale registratie endpoint
	publicApi.Post("/aanmelden", h.RegisterParticipant)

	// V30: Upgrade endpoint (vereist alleen email verificatie, geen volledige auth)
	publicApi.Post("/upgrade-to-full-account", h.UpgradeToFullAccount)

	// Info endpoints
	publicApi.Get("/events/active", h.GetActiveEvent)
}

const DKL_2026_EVENT_ID = "f1a75cc7-303e-4207-b501-8eea557bff33" // Default event ID - De Koninklijke Loop 2025

// RegisterParticipant - Publieke registratie endpoint (V30: Duaal systeem)
// @Summary Registreer deelnemer voor evenement
// @Description Registreert een new deelnemer met keuze tussen full account (met app toegang) of temporary registratie (alleen voor dit jaar)
// @Tags Public Registration
// @Accept json
// @Produce json
// @Param registration body models.PublicRegistrationRequest true "Registratie gegevens"
// @Success 200 {object} models.PublicRegistrationResponse
// @Router /api/public/aanmelden [post]
func (h *PublicRegistrationHandler) RegisterParticipant(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse request
	var req models.PublicRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
			"code":  "INVALID_INPUT",
		})
	}

	// Valideer verplichte velden
	if err := h.validateRegistrationRequest(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
	}

	// Bepaal event ID
	eventID := DKL_2026_EVENT_ID
	if req.EventID != nil && *req.EventID != "" {
		eventID = *req.EventID
	}

	// Haal event op
	event, err := h.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		logger.Error("Event niet gevonden", "event_id", eventID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Evenement niet gevonden",
			"code":  "EVENT_NOT_FOUND",
		})
	}

	// Check of email al bestaat voor dit jaar (alleen voor temporary accounts)
	if !req.WantAccount {
		existing, _ := h.checkExistingTemporaryRegistration(ctx, req.Email, 2026)
		if existing != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Je bent al ingeschreven voor dit jaar met dit e-mailadres",
				"code":  "ALREADY_REGISTERED",
			})
		}
	} else {
		// Voor full accounts: check of gebruiker al bestaat
		existingGebruiker, _ := h.gebruikerRepo.GetByEmail(ctx, req.Email)
		if existingGebruiker != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Er bestaat al een account met dit e-mailadres",
				"code":  "EMAIL_EXISTS",
			})
		}
	}

	// Roep juiste registratie flow aan
	var response *models.PublicRegistrationResponse
	if req.WantAccount {
		response, err = h.registerFullAccount(ctx, &req, event)
	} else {
		response, err = h.registerTemporaryAccount(ctx, &req, event)
	}

	if err != nil {
		logger.Error("Registratie mislukt", "error", err, "email", req.Email)
		// Stuur de specifieke foutmelding terug naar de client
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Registratie mislukt: %s", err.Error()),
			"code":  "REGISTRATION_FAILED",
		})
	}

	// ==============================================================================
	// LINTER FIX AANROEP 1
	// ==============================================================================
	// Verstuur bevestigingsemail (async)
	go h.sendConfirmationEmail(&req, response)

	// ==============================================================================
	// LINTER FIX AANROEP 2
	// ==============================================================================
	// Verstuur notificatie naar admins
	go h.sendAdminNotification(context.Background(), &req)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// registerFullAccount - Registreert met volledig account + app toegang
// V34: Gebruikt participants.wachtwoord_hash voor authenticatie (GEEN Gebruiker record)
func (h *PublicRegistrationHandler) registerFullAccount(ctx context.Context, req *models.PublicRegistrationRequest, event *models.Event) (*models.PublicRegistrationResponse, error) {
	// Hash wachtwoord
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Wachtwoord), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("wachtwoord hash failed: %w", err)
	}
	hashedPasswordStr := string(hashedPassword)

	// ✅ Stap 1: Maak Participant aan (met wachtwoord_hash, GEEN event-data!)
	participant := &models.Participant{
		Naam:           req.Naam,
		Email:          req.Email,
		Telefoon:       req.Telefoon, // ✅ Telefoon WEL opslaan in participants
		Terms:          req.Terms,
		TestMode:       req.TestMode,
		AccountType:    models.AccountTypeFull,
		WachtwoordHash: &hashedPasswordStr, // ✅ Voor authenticatie
		HasAppAccess:   true,
		// ❌ GEEN: Rol, Afstand, Ondersteuning - die kolommen bestaan niet meer!
		// ❌ GEEN: GebruikerID - geen Gebruiker record meer voor publieke registraties
	}

	if err := h.participantRepo.Create(ctx, participant); err != nil {
		return nil, fmt.Errorf("participant create failed: %w", err)
	}

	// ❌ VERWIJDERD: Geen Gebruiker record meer aanmaken
	// ❌ VERWIJDERD: Geen RBAC rol toewijzingen meer
	// ✅ Gebruikers tabel is ALLEEN voor admin/staff
	// ✅ Full accounts authenticeren via participants.wachtwoord_hash

	// ✅ Stap 2: Maak EventRegistration aan (hier komt event-data!)
	registration := &models.EventRegistration{
		EventID:             event.ID,
		ParticipantID:       participant.ID,
		ParticipantRoleName: &req.Rol,          // ✅ Rol in event_registrations
		DistanceRoute:       &req.Afstand,      // ✅ Afstand in event_registrations
		Ondersteuning:       req.Ondersteuning, // ✅ Ondersteuning in event_registrations
		Bijzonderheden:      req.Bijzonderheden,
		Status:              "nieuw", // ✅ V27: Gebruik 'nieuw' status (niet 'registered')
		TestMode:            req.TestMode,
	}

	if err := h.eventRegRepo.Create(ctx, registration); err != nil {
		logger.Error("EventRegistration create failed", "error", err)
		// Rollback: Verwijder participant
		_ = h.participantRepo.Delete(ctx, participant.ID)
		return nil, fmt.Errorf("kon event registratie niet aanmaken: %w", err)
	}

	logger.Info("Full account registered successfully",
		"participant_id", participant.ID,
		"email", req.Email,
		"rol", req.Rol,
		"afstand", req.Afstand)

	return &models.PublicRegistrationResponse{
		Success:        true,
		Message:        "Je bent succesvol ingeschreven met een volledig account! Je hebt nu toegang tot de DKL Step App.",
		ParticipantID:  participant.ID,
		RegistrationID: registration.ID,
		AccountType:    models.AccountTypeFull,
		HasAppAccess:   true,
		GebruikerID:    nil, // ✅ Geen gebruiker_id meer
		EventName:      event.Name,
		EventDate:      event.StartTime.Format("02-01-2006"),
	}, nil
}

// registerTemporaryAccount - Registreert zonder account (alleen voor dit jaar)
func (h *PublicRegistrationHandler) registerTemporaryAccount(ctx context.Context, req *models.PublicRegistrationRequest, event *models.Event) (*models.PublicRegistrationResponse, error) {
	currentYear := time.Now().Year()

	// Maak Participant aan (zonder wachtwoord)
	participant := &models.Participant{
		Naam:             req.Naam,
		Email:            req.Email,
		Telefoon:         req.Telefoon,
		Terms:            req.Terms,
		TestMode:         req.TestMode,
		AccountType:      models.AccountTypeTemporary,
		RegistrationYear: &currentYear,
		HasAppAccess:     false,
		// WachtwoordHash blijft NULL
		// GebruikerID blijft NULL
	}

	if err := h.participantRepo.Create(ctx, participant); err != nil {
		return nil, fmt.Errorf("participant create failed: %w", err)
	}

	// Maak EventRegistration aan
	registration := &models.EventRegistration{
		EventID:             event.ID,
		ParticipantID:       participant.ID,
		ParticipantRoleName: &req.Rol,
		DistanceRoute:       &req.Afstand,
		Ondersteuning:       req.Ondersteuning,
		Bijzonderheden:      req.Bijzonderheden,
		Status:              "nieuw", // ✅ V27: Gebruik 'nieuw' status
		TestMode:            req.TestMode,
	}

	// DIAGNOSTICS: Log exact values being saved for temporary account
	logger.Info("🔍 DIAGNOSTICS: Creating EventRegistration for TEMPORARY account",
		"event_id", event.ID,
		"participant_id", participant.ID,
		"rol_value", req.Rol,
		"rol_pointer", registration.ParticipantRoleName,
		"afstand_value", req.Afstand,
		"afstand_pointer", registration.DistanceRoute,
		"ondersteuning", registration.Ondersteuning,
		"bijzonderheden", registration.Bijzonderheden,
		"status", registration.Status)

	if err := h.eventRegRepo.Create(ctx, registration); err != nil {
		logger.Error("FATALE FOUT: EventRegistration create failed", "error", err)
		// Rollback: Verwijder de zojuist aangemaakte participant
		_ = h.participantRepo.Delete(ctx, participant.ID)
		return nil, fmt.Errorf("kon event registratie niet aanmaken: %w", err)
	}

	// DIAGNOSTICS: Verify what was actually saved for temporary account
	logger.Info("🔍 DIAGNOSTICS: EventRegistration created successfully for TEMPORARY",
		"registration_id", registration.ID,
		"saved_rol", registration.ParticipantRoleName,
		"saved_afstand", registration.DistanceRoute,
		"saved_ondersteuning", registration.Ondersteuning,
		"saved_bijzonderheden", registration.Bijzonderheden)

	return &models.PublicRegistrationResponse{
		Success:        true,
		Message:        fmt.Sprintf("Je bent succesvol ingeschreven voor %d! Je registratie is geldig voor dit evenementjaar.", currentYear),
		ParticipantID:  participant.ID,
		RegistrationID: registration.ID,
		AccountType:    models.AccountTypeTemporary,
		HasAppAccess:   false,
		EventName:      event.Name,
		EventDate:      event.StartTime.Format("02-01-2006"),
	}, nil
}

// UpgradeToFullAccount - Upgrade tijdelijk account naar volledig account
// @Summary Upgrade naar volledig account
// @Description Upgrade een tijdelijke registratie naar een volledig account met app toegang
// @Tags Public Registration
// @Accept json
// @Produce json
// @Param upgrade body models.UpgradeToFullAccountRequest true "Upgrade gegevens"
// @Success 200 {object} models.UpgradeToFullAccountResponse
// @Router /api/public/upgrade-to-full-account [post]
func (h *PublicRegistrationHandler) UpgradeToFullAccount(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.UpgradeToFullAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ongeldige gegevens",
			"code":  "INVALID_INPUT",
		})
	}

	// Valideer wachtwoord sterkte
	if len(req.Wachtwoord) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Wachtwoord moet minimaal 8 karakters bevatten",
			"code":  "WEAK_PASSWORD",
		})
	}

	// Zoek participant met dit email
	participants, err := h.participantRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database fout",
			"code":  "DATABASE_ERROR",
		})
	}

	// Zoek een temporary account
	var participant *models.Participant
	for _, p := range participants {
		if p.AccountType == models.AccountTypeTemporary {
			participant = p
			break
		}
	}

	if participant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Geen tijdelijke registratie gevonden met dit e-mailadres",
			"code":  "NO_TEMPORARY_ACCOUNT",
		})
	}

	// Check nog een keer of dit echt een temporary account is
	if participant.AccountType != models.AccountTypeTemporary {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dit account is al een volledig account",
			"code":  "ALREADY_FULL_ACCOUNT",
		})
	}

	// Check of email niet al bestaat als full participant
	fullParticipants, _ := h.participantRepo.FindByEmail(ctx, req.Email)
	for _, p := range fullParticipants {
		if p.AccountType == models.AccountTypeFull && p.ID != participant.ID {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Er bestaat al een volledig account met dit e-mailadres",
				"code":  "EMAIL_EXISTS",
			})
		}
	}

	// Hash wachtwoord
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Wachtwoord), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Wachtwoord hash failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Upgrade mislukt",
			"code":  "HASH_FAILED",
		})
	}
	hashedPasswordStr := string(hashedPassword)

	// ❌ VERWIJDERD: Geen Gebruiker record meer aanmaken
	// ✅ Full accounts gebruiken participants.wachtwoord_hash voor authenticatie

	// ✅ Update Participant naar Full Account
	now := time.Now()
	participant.AccountType = models.AccountTypeFull
	participant.WachtwoordHash = &hashedPasswordStr // ✅ Voor authenticatie
	participant.HasAppAccess = true
	// ❌ VERWIJDERD: GebruikerID, UpgradedToGebruikerID - geen Gebruiker meer
	participant.UpgradedAt = &now
	// RegistrationYear blijft behouden voor historie

	if err := h.participantRepo.Update(ctx, participant); err != nil {
		logger.Error("Participant update failed during upgrade", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Upgrade mislukt",
			"code":  "PARTICIPANT_UPDATE_FAILED",
		})
	}

	logger.Info("Participant upgraded to full account",
		"participant_id", participant.ID,
		"email", req.Email)

	// Verstuur bevestigingsemail
	go h.sendUpgradeConfirmationEmail(participant)

	return c.Status(fiber.StatusOK).JSON(models.UpgradeToFullAccountResponse{
		Success:       true,
		Message:       "Je account is geüpgraded! Je hebt nu toegang tot de DKL Step App.",
		GebruikerID:   "", // ❌ Geen gebruiker_id meer
		ParticipantID: participant.ID,
		HasAppAccess:  true,
	})
}

// GetActiveEvent - Haal het actieve event op
// @Summary Haal actief evenement op
// @Description Geeft informatie over het actieve evenement
// @Tags Public Registration
// @Produce json
// @Success 200 {object} models.Event
// @Router /api/public/events/active [get]
func (h *PublicRegistrationHandler) GetActiveEvent(c *fiber.Ctx) error {
	ctx := c.Context()

	// Haal default event op
	event, err := h.eventRepo.GetByID(ctx, DKL_2026_EVENT_ID)
	if err != nil || event == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Geen actief evenement gevonden",
			"code":  "NO_ACTIVE_EVENT",
		})
	}

	return c.JSON(event)
}

// Helper functions

func (h *PublicRegistrationHandler) validateRegistrationRequest(req *models.PublicRegistrationRequest) error {
	// Naam verplicht
	if req.Naam == "" {
		return fmt.Errorf("naam is verplicht")
	}

	// Email verplicht en valide format
	if req.Email == "" {
		return fmt.Errorf("e-mailadres is verplicht")
	}

	// Rol verplicht
	if req.Rol == "" {
		return fmt.Errorf("rol is verplicht")
	}

	// Telefoon verplicht voor Begeleider en Vrijwilliger
	if (req.Rol == "Begeleider" || req.Rol == "Vrijwilliger") && req.Telefoon == "" {
		return fmt.Errorf("telefoonnummer is verplicht voor begeleiders en vrijwilligers")
	}

	// Afstand verplicht
	if req.Afstand == "" {
		return fmt.Errorf("afstand is verplicht")
	}

	// Ondersteuning verplicht
	if req.Ondersteuning == "" {
		return fmt.Errorf("ondersteuning keuze is verplicht")
	}

	// Bijzonderheden verplicht als ondersteuning Ja of Anders
	if (req.Ondersteuning == "Ja" || req.Ondersteuning == "Anders") && req.Bijzonderheden == "" {
		return fmt.Errorf("bijzonderheden zijn verplicht als je ondersteuning nodig hebt")
	}

	// Terms verplicht
	if !req.Terms {
		return fmt.Errorf("je moet akkoord gaan met de algemene voorwaarden")
	}

	// Als WantAccount = true, dan is wachtwoord verplicht
	if req.WantAccount {
		if req.Wachtwoord == nil || *req.Wachtwoord == "" {
			return fmt.Errorf("wachtwoord is verplicht voor een volledig account")
		}
		if len(*req.Wachtwoord) < 8 {
			return fmt.Errorf("wachtwoord moet minimaal 8 karakters bevatten")
		}
	}

	return nil
}

func (h *PublicRegistrationHandler) checkExistingTemporaryRegistration(ctx context.Context, email string, year int) (*models.Participant, error) {
	// Gebruik FindByEmail om alle participants met dit email op te halen
	participants, err := h.participantRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Zoek een temporary account voor dit jaar
	for _, participant := range participants {
		if participant.AccountType == models.AccountTypeTemporary &&
			participant.RegistrationYear != nil &&
			*participant.RegistrationYear == year {
			return participant, nil
		}
	}

	return nil, nil
}

// ==============================================================================
// LINTER FIX DEFINITIE 1
// ==============================================================================
func (h *PublicRegistrationHandler) sendConfirmationEmail(req *models.PublicRegistrationRequest, resp *models.PublicRegistrationResponse) {
	var subject, body string

	if req.WantAccount {
		subject = "Welkom bij De Koninklijke Loop! Je account is aangemaakt"
		body = fmt.Sprintf(`
            <h2>Welkom %s!</h2>
            <p>Je bent succesvol ingeschreven voor De Koninklijke Loop 2026 met een volledig account.</p>
            <p><strong>Je hebt nu toegang tot de DKL Step App!</strong></p>
            <ul>
                <li>Rol: %s</li>
                <li>Afstand: %s</li>
                <li>Event datum: %s</li>
            </ul>
            <p>Download de DKL Step App om je voortgang bij te houden!</p>
        `, req.Naam, req.Rol, req.Afstand, resp.EventDate)
	} else {
		subject = "Je bent ingeschreven voor De Koninklijke Loop 2026"
		body = fmt.Sprintf(`
            <h2>Bedankt voor je inschrijving %s!</h2>
            <p>Je bent succesvol ingeschreven voor De Koninklijke Loop 2026.</p>
            <ul>
                <li>Rol: %s</li>
                <li>Afstand: %s</li>
                <li>Event datum: %s</li>
            </ul>
            <p><em>Let op: Deze registratie is geldig voor het evenement van 2026. Voor app toegang kun je later upgraden naar een volledig account.</em></p>
        `, req.Naam, req.Rol, req.Afstand, resp.EventDate)
	}

	if err := h.emailService.SendEmail(req.Email, subject, body); err != nil {
		logger.Error("Failed to send confirmation email", "error", err, "email", req.Email)
	}
}

// ==============================================================================
// LINTER FIX DEFINITIE 2
// ==============================================================================
func (h *PublicRegistrationHandler) sendUpgradeConfirmationEmail(participant *models.Participant) {
	subject := "Je account is geüpgraded!"
	body := fmt.Sprintf(`
        <h2>Welkom terug %s!</h2>
        <p>Je tijdelijke registratie is succesvol geüpgraded naar een volledig account.</p>
        <p><strong>Je hebt nu toegang tot de DKL Step App!</strong></p>
        <p>Download de app en log in met je e-mailadres en wachtwoord.</p>
    `, participant.Naam)

	if err := h.emailService.SendEmail(participant.Email, subject, body); err != nil {
		logger.Error("Failed to send upgrade confirmation email", "error", err)
	}
}

// ==============================================================================
// LINTER FIX DEFINITIE 3
// ==============================================================================
func (h *PublicRegistrationHandler) sendAdminNotification(ctx context.Context, req *models.PublicRegistrationRequest) {
	if h.notificationService == nil {
		return
	}

	accountTypeStr := "tijdelijk"
	if req.WantAccount {
		accountTypeStr = "volledig"
	}

	title := "Nieuwe registratie"
	message := fmt.Sprintf(
		"<b>%s</b> heeft zich ingeschreven.\n\n"+
			"<b>Email:</b> %s\n"+
			"<b>Rol:</b> %s\n"+
			"<b>Afstand:</b> %s\n"+
			"<b>Account type:</b> %s",
		req.Naam, req.Email, req.Rol, req.Afstand, accountTypeStr,
	)

	_, err := h.notificationService.CreateNotification(
		ctx,
		models.NotificationTypeAanmelding,
		models.NotificationPriorityMedium,
		title,
		message,
	)

	if err != nil {
		logger.Error("Failed to send admin notification", "error", err)
	}
}

// ==============================================================================
// V30+RBAC: HELPER FUNCTIES VOOR ROL TOEWIJZING
// ==============================================================================

// assignParticipantUserRole wijst de basis participant_user rol toe aan een gebruiker
func (h *PublicRegistrationHandler) assignParticipantUserRole(ctx context.Context, gebruikerID string) error {
	// Zoek participant_user rol
	role, err := h.roleRepo.GetByName(ctx, "participant_user")
	if err != nil {
		return fmt.Errorf("participant_user rol niet gevonden: %w", err)
	}

	// Wijs rol toe via PermissionService (met proper audit logging)
	if err := h.permissionService.AssignRole(ctx, gebruikerID, role.ID, nil); err != nil {
		// Als de fout is dat de rol al toegewezen is, negeer we deze
		if err.Error() == "gebruiker heeft deze rol al" {
			logger.Debug("participant_user rol al toegewezen (mogelijk via trigger)", "gebruiker_id", gebruikerID)
			return nil
		}
		return fmt.Errorf("fout bij toewijzen participant_user rol: %w", err)
	}

	logger.Info("participant_user rol toegewezen", "gebruiker_id", gebruikerID, "role_id", role.ID)
	return nil
}

// assignRoleBasedOnParticipantRole wijst een rol-specifieke RBAC rol toe op basis van event rol
func (h *PublicRegistrationHandler) assignRoleBasedOnParticipantRole(ctx context.Context, gebruikerID string, participantRole string) error {
	var rbacRoleName string

	// Map participant event rol naar RBAC rol
	switch participantRole {
	case "Begeleider":
		rbacRoleName = "participant_guide"
	case "Vrijwilliger":
		rbacRoleName = "participant_volunteer"
	case "Deelnemer":
		// Deelnemers krijgen alleen participant_user rol (geen extra rol)
		return nil
	default:
		logger.Debug("Onbekende participant rol, geen extra RBAC rol", "role", participantRole)
		return nil
	}

	// Zoek RBAC rol
	role, err := h.roleRepo.GetByName(ctx, rbacRoleName)
	if err != nil {
		return fmt.Errorf("RBAC rol %s niet gevonden: %w", rbacRoleName, err)
	}

	// Wijs rol toe
	if err := h.permissionService.AssignRole(ctx, gebruikerID, role.ID, nil); err != nil {
		if err.Error() == "gebruiker heeft deze rol al" {
			logger.Debug("Rol-specifieke RBAC rol al toegewezen", "gebruiker_id", gebruikerID, "role", rbacRoleName)
			return nil
		}
		return fmt.Errorf("fout bij toewijzen rol %s: %w", rbacRoleName, err)
	}

	logger.Info("Rol-specifieke RBAC rol toegewezen",
		"gebruiker_id", gebruikerID,
		"rbac_role", rbacRoleName,
		"participant_role", participantRole)
	return nil
}

// getParticipantEventRole haalt de event rol op voor een participant (voor upgrade flow)
func (h *PublicRegistrationHandler) getParticipantEventRole(ctx context.Context, participantID string) string {
	// Haal meest recente event registratie op (gebruik ListByParticipantID)
	registrations, err := h.eventRegRepo.ListByParticipantID(ctx, participantID)
	if err != nil || len(registrations) == 0 {
		logger.Debug("Geen event registraties gevonden voor participant", "participant_id", participantID)
		return ""
	}

	// Return de rol van de meest recente registratie
	if registrations[0].ParticipantRoleName != nil {
		return *registrations[0].ParticipantRoleName
	}

	return ""
}
