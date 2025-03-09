package models

// ContactEmailData bevat de data nodig voor het versturen van contact formulier emails
type ContactEmailData struct {
	ToAdmin    bool
	Contact    *ContactFormulier
	AdminEmail string
}

// AanmeldingEmailData bevat de data nodig voor het versturen van aanmelding emails
type AanmeldingEmailData struct {
	ToAdmin    bool                 `json:"to_admin"`
	Aanmelding *AanmeldingFormulier `json:"aanmelding"`
	AdminEmail string               `json:"admin_email,omitempty"`
}
