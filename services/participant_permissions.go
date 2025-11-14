package services

// Hier definiëren we de ENIGE 'source of truth' voor deelnemer-permissies.
//
// OPLOSSING (FINALE): Deze lijst is nu gecorrigeerd om de permissies te gebruiken
// die de applicatie daadwerkelijk verwacht (bv. 'view_own'),
// in plaats van de generieke permissies uit de migraties (bv. 'read').

var participantPermissionsMap = map[string][]string{
	// Permissies die de app verwacht (gebaseerd op oude logic/logs)
	"steps": {
		"view_own", // De app vraagt hierom (zie logs)
		"create",   // De app verwacht dit (vervanging van 'write')
	},
	"participant": {
		"view_own",   // (Vervanging van 'read')
		"update_own", // (Vervanging van 'write')
	},
	"app": {
		"access",
	},
	"leaderboard": {
		"view", // (Hernoemd van 'read' naar 'view' voor consistentie)
	},
	"events": {
		"view",
		"register",
	},

	// Aliassen voor de frontend
	"profile": {
		"read",   // Alias voor participant:view_own
		"update", // Alias voor participant:update_own
	},
}

// GetParticipantPermissionsList genereert de lijst voor de frontend response (in AuthHandler)
// Dit zorgt ervoor dat de frontend 100% consistent is met de backend.
func GetParticipantPermissionsList() []map[string]string {
	list := make([]map[string]string, 0)

	for resource, actions := range participantPermissionsMap {
		for _, action := range actions {
			// Voeg de daadwerkelijke permissie toe
			list = append(list, map[string]string{
				"resource": resource,
				"action":   action,
			})
		}
	}
	return list
}
