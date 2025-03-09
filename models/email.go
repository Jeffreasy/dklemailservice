package models

// ContactEmailData bevat de data nodig voor het versturen van contact formulier emails
type ContactEmailData struct {
	ToAdmin    bool              `json:"to_admin"`
	Contact    *ContactFormulier `json:"contact"`
	AdminEmail string            `json:"admin_email,omitempty"`
}

// AanmeldingEmailData bevat de data nodig voor het versturen van aanmelding emails
type AanmeldingEmailData struct {
	ToAdmin    bool                 `json:"to_admin"`
	Aanmelding *AanmeldingFormulier `json:"aanmelding"`
	AdminEmail string               `json:"admin_email,omitempty"`
}
