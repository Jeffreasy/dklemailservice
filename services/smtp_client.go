package services

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"gopkg.in/gomail.v2"
)

// SMTPDialer is een interface voor het testen van SMTP verbindingen
type SMTPDialer interface {
	Dial() error
}

// SMTPConfig bevat de configuratie voor een SMTP verbinding
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseSSL   bool // Toegevoegd voor directe SSL verbindingen
}

// RealSMTPClient implementeert de SMTP client interface
type RealSMTPClient struct {
	defaultConf        *SMTPConfig
	regConf            *SMTPConfig
	wfcConf            *SMTPConfig // Nieuwe configuratie voor Whisky for Charity
	dialer             SMTPDialer
	defaultDialer      *gomail.Dialer
	registrationDialer *gomail.Dialer
	wfcDialer          *gomail.Dialer // Nieuwe dialer voor Whisky for Charity
	connMutex          sync.Mutex
}

// NewRealSMTPClient creates a new SMTP client with the given configuration
func NewRealSMTPClient(host, port, user, password, from, regHost, regPort, regUser, regPassword, regFrom string) *RealSMTPClient {
	defaultPortNum, err := strconv.Atoi(port)
	if err != nil {
		defaultPortNum = 587 // Default SMTP port
	}

	regPortNum, err := strconv.Atoi(regPort)
	if err != nil {
		regPortNum = defaultPortNum
	}

	defaultConf := &SMTPConfig{
		Host:     host,
		Port:     defaultPortNum,
		Username: user,
		Password: password,
		From:     from,
		UseSSL:   false,
	}

	regConf := &SMTPConfig{
		Host:     regHost,
		Port:     regPortNum,
		Username: regUser,
		Password: regPassword,
		From:     regFrom,
		UseSSL:   false,
	}

	// Maak persistente dialers voor betere performance
	defaultDialer := gomail.NewDialer(host, defaultPortNum, user, password)
	regDialer := gomail.NewDialer(regHost, regPortNum, regUser, regPassword)

	// Configureer keepalive en timeouts
	defaultDialer.SSL = false // Use STARTTLS instead of SSL
	regDialer.SSL = false     // Use STARTTLS instead of SSL

	return &RealSMTPClient{
		defaultConf:        defaultConf,
		regConf:            regConf,
		defaultDialer:      defaultDialer,
		registrationDialer: regDialer,
		connMutex:          sync.Mutex{},
	}
}

// NewRealSMTPClientWithWFC creates a new SMTP client including Whisky for Charity configuration
func NewRealSMTPClientWithWFC(host, port, user, password, from, regHost, regPort, regUser, regPassword, regFrom, wfcHost, wfcPort, wfcUser, wfcPassword, wfcFrom string, wfcUseSSL bool) *RealSMTPClient {
	// Maak eerst de standaard client
	client := NewRealSMTPClient(host, port, user, password, from, regHost, regPort, regUser, regPassword, regFrom)

	// Voeg Whisky for Charity configuratie toe
	wfcPortNum, err := strconv.Atoi(wfcPort)
	if err != nil {
		wfcPortNum = 465 // Default SSL SMTP port
	}

	wfcConf := &SMTPConfig{
		Host:     wfcHost,
		Port:     wfcPortNum,
		Username: wfcUser,
		Password: wfcPassword,
		From:     wfcFrom,
		UseSSL:   wfcUseSSL,
	}

	// Maak een nieuwe dialer
	wfcDialer := gomail.NewDialer(wfcHost, wfcPortNum, wfcUser, wfcPassword)
	wfcDialer.SSL = wfcUseSSL // Direct SSL voor port 465

	// Stel de nieuwe dialer in
	client.wfcConf = wfcConf
	client.wfcDialer = wfcDialer

	return client
}

// SetDialer stelt een custom dialer in voor tests
func (c *RealSMTPClient) SetDialer(d SMTPDialer) {
	c.dialer = d
}

// Send verzendt een email met de standaard configuratie
func (c *RealSMTPClient) Send(msg *EmailMessage) error {
	return c.sendWithDialer(msg, c.defaultConf, c.defaultDialer)
}

// SendRegistration verzendt een email met de registratie configuratie
func (c *RealSMTPClient) SendRegistration(msg *EmailMessage) error {
	return c.sendWithDialer(msg, c.regConf, c.registrationDialer)
}

// SendWFC verzendt een email met de Whisky for Charity configuratie
func (c *RealSMTPClient) SendWFC(msg *EmailMessage) error {
	if c.wfcDialer == nil || c.wfcConf == nil {
		return fmt.Errorf("whisky for charity email configuration is not set")
	}
	return c.sendWithDialer(msg, c.wfcConf, c.wfcDialer)
}

// sendWithDialer verzendt een email met de opgegeven configuratie en dialer
func (c *RealSMTPClient) sendWithDialer(msg *EmailMessage, conf *SMTPConfig, dialer *gomail.Dialer) error {
	if conf == nil {
		return fmt.Errorf("smtp configuration is nil")
	}

	if msg.To == "" {
		return fmt.Errorf("invalid recipient")
	}

	// Als we een test dialer hebben, gebruik die
	if c.dialer != nil {
		if err := c.dialer.Dial(); err != nil {
			return err
		}
		if mockSMTP, ok := c.dialer.(SMTPClient); ok {
			return mockSMTP.Send(msg)
		}
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", conf.From)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	// Gebruik connection pooling voor betere performance
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	// Probeer e-mail te verzenden met bestaande verbinding of maak nieuwe verbinding
	err := dialer.DialAndSend(m)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	return nil
}

// SendEmail is een helper functie voor backwards compatibility
func (c *RealSMTPClient) SendEmail(to, subject, body string) error {
	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return c.Send(msg)
}

// SendWFCEmail is een helper functie voor het verzenden van Whisky for Charity emails
func (c *RealSMTPClient) SendWFCEmail(to, subject, body string) error {
	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return c.SendWFC(msg)
}
