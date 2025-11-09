package models

// ContactEmailData is de struct voor de contact-template
type ContactEmailData struct {
	Contact    *ContactFormulier
	AdminEmail string
	ToAdmin    bool
}

// RegistrationEmailData bevat alle data voor de registratie-templates (V28 refactor).
// Dit vervangt de oude AanmeldingEmailData.
type RegistrationEmailData struct {
	ToAdmin      bool
	AdminEmail   string
	Participant  *Participant
	Registration *EventRegistration
}
