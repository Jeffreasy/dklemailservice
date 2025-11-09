package services

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// NotulenService handles business logic for notulen
type NotulenService struct {
	repo              *repository.PostgresNotulenRepository
	userService       *AuthService
	permissionService PermissionService
	hub               *NotulenHub
}

// NewNotulenService creates a new notulen service
func NewNotulenService(repo *repository.PostgresNotulenRepository, userService *AuthService, permissionService PermissionService, hub *NotulenHub) *NotulenService {
	return &NotulenService{
		repo:              repo,
		userService:       userService,
		permissionService: permissionService,
		hub:               hub,
	}
}

// CreateNotulen creates a new notulen document
func (s *NotulenService) CreateNotulen(ctx context.Context, userID uuid.UUID, req *models.NotulenCreateRequest) (*models.Notulen, error) {
	// Parse vergadering datum
	vergaderingDatum, err := time.Parse("2006-01-02", req.VergaderingDatum)
	if err != nil {
		return nil, fmt.Errorf("ongeldige vergadering datum: %w", err)
	}

	// Convert string UUIDs to UUIDArray for registered users
	var aanwezigenGebruikerIDs models.UUIDArray
	var afwezigenGebruikerIDs models.UUIDArray

	// GEFIXED: Gebruik AanwezigenGebruikerIDs (uit req)
	if req.AanwezigenGebruikerIDs != nil {
		uuids := make([]uuid.UUID, 0, len(req.AanwezigenGebruikerIDs))
		for _, userUUIDStr := range req.AanwezigenGebruikerIDs {
			if userUUID, err := uuid.Parse(userUUIDStr); err == nil {
				uuids = append(uuids, userUUID)
			}
		}
		aanwezigenGebruikerIDs = uuids
	} else {
		aanwezigenGebruikerIDs = models.UUIDArray{}
	}

	// GEFIXED: Gebruik AfwezigenGebruikerIDs (uit req)
	if req.AfwezigenGebruikerIDs != nil {
		uuids := make([]uuid.UUID, 0, len(req.AfwezigenGebruikerIDs))
		for _, userUUIDStr := range req.AfwezigenGebruikerIDs {
			if userUUID, err := uuid.Parse(userUUIDStr); err == nil {
				uuids = append(uuids, userUUID)
			}
		}
		afwezigenGebruikerIDs = uuids
	} else {
		afwezigenGebruikerIDs = models.UUIDArray{}
	}

	// Create notulen object
	notulen := &models.Notulen{
		ID:                     uuid.New(),
		Titel:                  req.Titel,
		VergaderingDatum:       vergaderingDatum,
		Locatie:                req.Locatie,
		Voorzitter:             req.Voorzitter,
		Notulist:               req.Notulist,
		Aanwezigen:             req.Aanwezigen,         // Legacy field
		Afwezigen:              req.Afwezigen,          // Legacy field
		AanwezigenGebruikerIDs: aanwezigenGebruikerIDs, // GEFIXED
		AfwezigenGebruikerIDs:  afwezigenGebruikerIDs,  // GEFIXED
		AanwezigenGasten:       req.AanwezigenGasten,   // GEFIXED
		AfwezigenGasten:        req.AfwezigenGasten,    // GEFIXED
		AgendaItems:            req.AgendaItems,        // Direct assignment
		Besluiten:              req.Besluiten,          // Direct assignment
		Actiepunten:            req.Actiepunten,        // Direct assignment
		Notities:               req.Notities,
		Status:                 "draft",
		Versie:                 1,
		CreatedBy:              userID,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
		// UpdatedByID wordt bij aanmaak leeg gelaten
	}

	// Save to database
	if err := s.repo.Create(ctx, notulen); err != nil {
		return nil, fmt.Errorf("failed to create notulen: %w", err)
	}

	return notulen, nil
}

// GetNotulen retrieves a notulen by ID
// GEFIXED: Retourneert nu NotulenResponse
func (s *NotulenService) GetNotulen(ctx context.Context, id uuid.UUID) (*models.NotulenResponse, error) {
	notulen, err := s.repo.GetByID(ctx, id)
	if err != nil || notulen == nil {
		return nil, err
	}

	return s.ConvertToNotulenResponse(notulen), nil
}

// UpdateNotulen updates an existing notulen
func (s *NotulenService) UpdateNotulen(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *models.NotulenUpdateRequest) (*models.Notulen, error) {
	// Get existing notulen
	notulen, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notulen: %w", err)
	}
	if notulen == nil {
		return nil, fmt.Errorf("notulen not found")
	}

	// Check if user can edit (only creator or admin)
	if notulen.CreatedBy != userID {
		// Check admin permission
		if !s.permissionService.HasPermission(ctx, userID.String(), "admin", "access") {
			return nil, fmt.Errorf("geen toestemming om notulen te bewerken")
		}
	}

	// Check if notulen is finalized
	if notulen.Status == "finalized" {
		return nil, fmt.Errorf("gefinaliseerde notulen kunnen niet worden bewerkt")
	}

	// Update fields
	if req.Titel != "" {
		notulen.Titel = req.Titel
	}
	if req.Locatie != "" {
		notulen.Locatie = req.Locatie
	}
	if req.Voorzitter != "" {
		notulen.Voorzitter = req.Voorzitter
	}
	if req.Notulist != "" {
		notulen.Notulist = req.Notulist
	}
	if req.Aanwezigen != nil {
		notulen.Aanwezigen = req.Aanwezigen
	}
	if req.Afwezigen != nil {
		notulen.Afwezigen = req.Afwezigen
	}

	// Handle new UUID participant fields
	// GEFIXED: Gebruik ...GebruikerIDs
	if req.AanwezigenGebruikerIDs != nil {
		uuids := make([]uuid.UUID, 0, len(req.AanwezigenGebruikerIDs))
		for _, userUUIDStr := range req.AanwezigenGebruikerIDs {
			if userUUID, err := uuid.Parse(userUUIDStr); err == nil {
				uuids = append(uuids, userUUID)
			}
		}
		notulen.AanwezigenGebruikerIDs = uuids
	}
	if req.AfwezigenGebruikerIDs != nil {
		uuids := make([]uuid.UUID, 0, len(req.AfwezigenGebruikerIDs))
		for _, userUUIDStr := range req.AfwezigenGebruikerIDs {
			if userUUID, err := uuid.Parse(userUUIDStr); err == nil {
				uuids = append(uuids, userUUID)
			}
		}
		notulen.AfwezigenGebruikerIDs = uuids
	}
	if req.AanwezigenGasten != nil {
		notulen.AanwezigenGasten = req.AanwezigenGasten
	}
	if req.AfwezigenGasten != nil {
		notulen.AfwezigenGasten = req.AfwezigenGasten
	}

	// Direct assignment for JSONB fields
	if req.AgendaItems != nil {
		notulen.AgendaItems = req.AgendaItems
	}
	if req.Besluiten != nil {
		notulen.Besluiten = req.Besluiten
	}
	if req.Actiepunten != nil {
		notulen.Actiepunten = req.Actiepunten
	}
	if req.Notities != "" {
		notulen.Notities = req.Notities
	}

	notulen.UpdatedAt = time.Now()
	notulen.UpdatedByID = userID // GEFIXED: Gebruik UpdatedByID

	// Save changes
	if err := s.repo.Update(ctx, notulen); err != nil {
		return nil, fmt.Errorf("failed to update notulen: %w", err)
	}

	// Broadcast update to WebSocket clients
	if s.hub != nil {
		response := s.ConvertToNotulenResponse(notulen)
		s.hub.BroadcastNotulenUpdate(notulen.ID, userID, response)
	}

	return notulen, nil
}

// FinalizeNotulen finalizes a notulen document
func (s *NotulenService) FinalizeNotulen(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *models.NotulenFinalizeRequest) error {
	// Get existing notulen
	notulen, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get notulen: %w", err)
	}
	if notulen == nil {
		return fmt.Errorf("notulen not found")
	}

	// Check permissions
	if notulen.CreatedBy != userID {
		// Check admin permission
		if !s.permissionService.HasPermission(ctx, userID.String(), "admin", "access") {
			return fmt.Errorf("geen toestemming om notulen te finaliseren")
		}
	}

	// Finalize
	if err := s.repo.Finalize(ctx, id, userID); err != nil {
		return err
	}

	// Broadcast finalization to WebSocket clients
	if s.hub != nil {
		s.hub.BroadcastNotulenFinalized(id, userID)
	}

	return nil
}

// ArchiveNotulen archives a notulen document
func (s *NotulenService) ArchiveNotulen(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	// Get existing notulen
	notulen, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get notulen: %w", err)
	}
	if notulen == nil {
		return fmt.Errorf("notulen not found")
	}

	// Check permissions
	if notulen.CreatedBy != userID {
		// Check admin permission
		if !s.permissionService.HasPermission(ctx, userID.String(), "admin", "access") {
			return fmt.Errorf("geen toestemming om notulen te archiveren")
		}
	}

	// Archive
	if err := s.repo.Archive(ctx, id); err != nil {
		return err
	}

	// Broadcast archiving to WebSocket clients
	if s.hub != nil {
		s.hub.BroadcastNotulenArchived(id, userID)
	}

	return nil
}

// DeleteNotulen deletes a notulen (soft delete)
func (s *NotulenService) DeleteNotulen(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	// Get existing notulen
	notulen, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get notulen: %w", err)
	}
	if notulen == nil {
		return fmt.Errorf("notulen not found")
	}

	// Check permissions
	if notulen.CreatedBy != userID {
		// Check admin permission
		if !s.permissionService.HasPermission(ctx, userID.String(), "admin", "access") {
			return fmt.Errorf("geen toestemming om notulen te verwijderen")
		}
	}

	// Delete
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Broadcast deletion to WebSocket clients
	if s.hub != nil {
		s.hub.BroadcastNotulenDeleted(id, userID)
	}

	return nil
}

// ListNotulen lists notulen with filtering and pagination
// GEFIXED: Retourneert []models.NotulenResponse
func (s *NotulenService) ListNotulen(ctx context.Context, filters *models.NotulenSearchFilters) ([]models.NotulenResponse, int, error) {
	notulenList, total, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response models with resolved user names
	// GEFIXED: Retourneert []models.NotulenResponse
	responses := make([]models.NotulenResponse, len(notulenList))
	for i, notulen := range notulenList {
		responses[i] = *s.ConvertToNotulenResponse(&notulen)
	}

	return responses, total, nil
}

// SearchNotulen performs full-text search on notulen
// GEFIXED: Retourneert []models.NotulenResponse
func (s *NotulenService) SearchNotulen(ctx context.Context, query string, filters *models.NotulenSearchFilters) ([]models.NotulenResponse, int, error) {
	notulenList, total, err := s.repo.Search(ctx, query, filters)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response models with resolved user names
	// GEFIXED: Retourneert []models.NotulenResponse
	responses := make([]models.NotulenResponse, len(notulenList))
	for i, notulen := range notulenList {
		responses[i] = *s.ConvertToNotulenResponse(&notulen)
	}

	return responses, total, nil
}

// GetNotulenVersions retrieves all versions of a notulen
func (s *NotulenService) GetNotulenVersions(ctx context.Context, notulenID uuid.UUID) ([]models.NotulenVersie, error) {
	versions, err := s.repo.GetVersions(ctx, notulenID)
	if err != nil {
		return nil, err
	}

	// Resolve user names for all versions
	for i := range versions {
		version := &versions[i]

		// GEFIXED: Gebruik ...GebruikerIDs
		if len(version.AanwezigenGebruikerIDs) > 0 {
			userNames := s.resolveUserNames(version.AanwezigenGebruikerIDs)
			version.Aanwezigen = append(userNames, version.AanwezigenGasten...)
		}

		// GEFIXED: Gebruik ...GebruikerIDs
		if len(version.AfwezigenGebruikerIDs) > 0 {
			userNames := s.resolveUserNames(version.AfwezigenGebruikerIDs)
			version.Afwezigen = append(userNames, version.AfwezigenGasten...)
		}
	}

	return versions, nil
}

// GetNotulenVersion retrieves a specific version of a notulen
func (s *NotulenService) GetNotulenVersion(ctx context.Context, notulenID uuid.UUID, versie int) (*models.NotulenVersie, error) {
	version, err := s.repo.GetVersion(ctx, notulenID, versie)
	if err != nil || version == nil {
		return version, err
	}

	// Resolve user names for display
	// GEFIXED: Gebruik ...GebruikerIDs
	if len(version.AanwezigenGebruikerIDs) > 0 {
		userNames := s.resolveUserNames(version.AanwezigenGebruikerIDs)
		version.Aanwezigen = append(userNames, version.AanwezigenGasten...)
	}

	// GEFIXED: Gebruik ...GebruikerIDs
	if len(version.AfwezigenGebruikerIDs) > 0 {
		userNames := s.resolveUserNames(version.AfwezigenGebruikerIDs)
		version.Afwezigen = append(userNames, version.AfwezigenGasten...)
	}

	return version, nil
}

// RenderMarkdown renders a notulen as Markdown
func (s *NotulenService) RenderMarkdown(notulen *models.Notulen) (string, error) {
	// Create a custom template data structure that includes resolved names
	templateData := struct {
		*models.Notulen
		AanwezigenResolved []string
		AfwezigenResolved  []string
	}{
		Notulen: notulen,
	}

	// Resolve user names from UUIDs for display
	// GEFIXED: Gebruik ...GebruikerIDs
	if len(notulen.AanwezigenGebruikerIDs) > 0 {
		templateData.AanwezigenResolved = s.resolveUserNames(notulen.AanwezigenGebruikerIDs)
	}
	if len(notulen.AfwezigenGebruikerIDs) > 0 {
		templateData.AfwezigenResolved = s.resolveUserNames(notulen.AfwezigenGebruikerIDs)
	}

	// Combine with guest names
	templateData.AanwezigenResolved = append(templateData.AanwezigenResolved, notulen.AanwezigenGasten...)
	templateData.AfwezigenResolved = append(templateData.AfwezigenResolved, notulen.AfwezigenGasten...)

	// If no new data, fall back to legacy fields
	if len(templateData.AanwezigenResolved) == 0 && len(notulen.Aanwezigen) > 0 {
		templateData.AanwezigenResolved = notulen.Aanwezigen
	}
	if len(templateData.AfwezigenResolved) == 0 && len(notulen.Afwezigen) > 0 {
		templateData.AfwezigenResolved = notulen.Afwezigen
	}

	tmpl := template.Must(template.New("notulen").Parse(notulenMarkdownTemplate))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return buf.String(), nil
}

// ConvertToNotulenResponse converts a Notulen model to NotulenResponse with resolved user names
// GEFIXED: Gebruikt NotulenResponse
func (s *NotulenService) ConvertToNotulenResponse(notulen *models.Notulen) *models.NotulenResponse {
	// Collect all unique user UUIDs that need to be resolved
	userUUIDs := []uuid.UUID{notulen.CreatedBy}

	// GEFIXED: Gebruik UpdatedByID
	if notulen.UpdatedByID != uuid.Nil {
		userUUIDs = append(userUUIDs, notulen.UpdatedByID)
	}

	if notulen.FinalizedBy != nil && *notulen.FinalizedBy != uuid.Nil {
		userUUIDs = append(userUUIDs, *notulen.FinalizedBy)
	}

	// Remove duplicates
	uniqueUUIDs := make([]uuid.UUID, 0, len(userUUIDs))
	seen := make(map[uuid.UUID]bool)
	for _, uid := range userUUIDs {
		if !seen[uid] {
			seen[uid] = true
			uniqueUUIDs = append(uniqueUUIDs, uid)
		}
	}

	// Get user names
	userNames := s.resolveUserNames(uniqueUUIDs)

	// Create a map for quick lookup
	nameMap := make(map[uuid.UUID]string)
	for i, uid := range uniqueUUIDs {
		if i < len(userNames) {
			nameMap[uid] = userNames[i]
		} else {
			nameMap[uid] = fmt.Sprintf("User-%s", uid.String()[:8])
		}
	}

	// Resolve participant UUIDs en combineer met guests
	var aanwezigenCombined pq.StringArray
	var afwezigenCombined pq.StringArray

	// GEFIXED: Gebruik ...GebruikerIDs
	if len(notulen.AanwezigenGebruikerIDs) > 0 {
		resolvedNames := s.resolveUserNames(notulen.AanwezigenGebruikerIDs)
		aanwezigenCombined = append(aanwezigenCombined, resolvedNames...)
	}
	if len(notulen.AanwezigenGasten) > 0 {
		aanwezigenCombined = append(aanwezigenCombined, notulen.AanwezigenGasten...)
	}
	if len(aanwezigenCombined) == 0 && len(notulen.Aanwezigen) > 0 {
		aanwezigenCombined = notulen.Aanwezigen
	}

	// GEFIXED: Gebruik ...GebruikerIDs
	if len(notulen.AfwezigenGebruikerIDs) > 0 {
		resolvedNames := s.resolveUserNames(notulen.AfwezigenGebruikerIDs)
		afwezigenCombined = append(afwezigenCombined, resolvedNames...)
	}
	if len(notulen.AfwezigenGasten) > 0 {
		afwezigenCombined = append(afwezigenCombined, notulen.AfwezigenGasten...)
	}
	if len(afwezigenCombined) == 0 && len(notulen.Afwezigen) > 0 {
		afwezigenCombined = notulen.Afwezigen
	}

	// Create response with resolved names
	// GEFIXED: Gebruik NotulenResponse
	response := &models.NotulenResponse{
		ID:                     notulen.ID,
		Titel:                  notulen.Titel,
		VergaderingDatum:       notulen.VergaderingDatum,
		Locatie:                notulen.Locatie,
		Voorzitter:             notulen.Voorzitter,
		Notulist:               notulen.Notulist,
		Aanwezigen:             aanwezigenCombined,
		Afwezigen:              afwezigenCombined,
		AanwezigenGebruikerIDs: notulen.AanwezigenGebruikerIDs,
		AfwezigenGebruikerIDs:  notulen.AfwezigenGebruikerIDs,
		AanwezigenGasten:       notulen.AanwezigenGasten,
		AfwezigenGasten:        notulen.AfwezigenGasten,
		AgendaItems:            notulen.AgendaItems,
		Besluiten:              notulen.Besluiten,
		Actiepunten:            notulen.Actiepunten,
		Notities:               notulen.Notities,
		Status:                 notulen.Status,
		Versie:                 notulen.Versie,
		CreatedBy:              notulen.CreatedBy,
		CreatedByName:          nameMap[notulen.CreatedBy],
		CreatedAt:              notulen.CreatedAt,
		UpdatedAt:              notulen.UpdatedAt,
		UpdatedByID:            notulen.UpdatedByID, // GEFIXED
		FinalizedAt:            notulen.FinalizedAt,
		FinalizedBy:            notulen.FinalizedBy,
	}

	// Set resolved names for updated_by and finalized_by
	// GEFIXED: Gebruik UpdatedByID
	if notulen.UpdatedByID != uuid.Nil {
		response.UpdatedByName = nameMap[notulen.UpdatedByID]
	}

	if notulen.FinalizedBy != nil && *notulen.FinalizedBy != uuid.Nil {
		response.FinalizedByName = nameMap[*notulen.FinalizedBy]
	}

	return response
}

// resolveUserNames converts UUID array to user names
func (s *NotulenService) resolveUserNames(userUUIDs []uuid.UUID) []string {
	if len(userUUIDs) == 0 {
		return []string{}
	}

	// Query database for user names
	names := make([]string, 0, len(userUUIDs))

	// Build query to get user names by UUIDs
	query := `SELECT naam FROM gebruikers WHERE id = ANY($1) AND is_actief = true ORDER BY naam`
	rows, err := s.repo.DB().Raw(query, pq.Array(userUUIDs)).Rows()
	if err != nil {
		// Fallback to placeholder names on error
		for _, uid := range userUUIDs {
			names = append(names, fmt.Sprintf("User-%s", uid.String()[:8]))
		}
		return names
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			names = append(names, name)
		}
	}

	// If we got fewer names than UUIDs, add placeholders for missing ones
	if len(names) < len(userUUIDs) {
		for i := len(names); i < len(userUUIDs); i++ {
			names = append(names, fmt.Sprintf("User-%s", userUUIDs[i].String()[:8]))
		}
	}

	return names
}

// Hub returns the WebSocket hub for external access
func (s *NotulenService) Hub() *NotulenHub {
	return s.hub
}

// ValidateNotulen validates notulen data
func (s *NotulenService) ValidateNotulen(req *models.NotulenCreateRequest) error {
	if strings.TrimSpace(req.Titel) == "" {
		return fmt.Errorf("titel is verplicht")
	}
	if len(req.Titel) > 255 {
		return fmt.Errorf("titel mag niet langer zijn dan 255 karakters")
	}
	if req.VergaderingDatum == "" {
		return fmt.Errorf("vergadering datum is verplicht")
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", req.VergaderingDatum); err != nil {
		return fmt.Errorf("ongeldige vergadering datum formaat (gebruik YYYY-MM-DD)")
	}

	return nil
}

// Markdown template for notulen rendering
// GEFIXED: Gebruikt AanwezigenResolved en AfwezigenResolved
// GEFIXED: Correcte range loops
const notulenMarkdownTemplate = `# {{.Titel}}

**Datum:** {{.VergaderingDatum.Format "2 January 2006"}}
**Locatie:** {{.Locatie}}
**Voorzitter:** {{.Voorzitter}}
**Notulist:** {{.Notulist}}

## Aanwezigen
{{range .AanwezigenResolved}}- {{.}}
{{end}}

## Afwezigen
{{range .AfwezigenResolved}}- {{.}}
{{end}}

## Agenda Items
{{range .AgendaItems}}
### {{.Titel}}
{{if .Beschrijving}}{{.Beschrijving}}{{end}}
{{if .Spreker}}**Spreker:** {{.Spreker}}{{end}}
{{if .Tijdslot}}**Tijdslot:** {{.Tijdslot}}{{end}}
{{end}}

## Besluiten
{{range .Besluiten}}
- {{.Beschrijving}}
  {{if .Verantwoordelijke}}**Verantwoordelijke:** {{.Verantwoordelijke}}{{end}}
  {{if .Deadline}}**Deadline:** {{.Deadline.Format "2-1-2006"}}{{end}}
{{end}}

## Actiepunten
{{range .Actiepunten}}
- [{{if eq .Status "completed"}}x{{else}} {{end}}] {{.Beschrijving}}
  **Verantwoordelijke:** {{.Verantwoordelijke}}
  {{if .Deadline}}**Deadline:** {{.Deadline.Format "2-1-2006"}}{{end}}
  **Status:** {{.Status}}
{{end}}

## Notities
{{.Notities}}

---
*Gemaakt op {{.CreatedAt.Format "2-1-2006 15:04"}} door {{.CreatedByName}}*
*Status: {{.Status}} | Versie: {{.Versie}}*
{{if .FinalizedAt}}*Gefinaliseerd op {{.FinalizedAt.Format "2-1-2006 15:04"}}*{{end}}`
