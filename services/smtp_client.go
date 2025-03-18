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
}

// RealSMTPClient implementeert de SMTP client interface
type RealSMTPClient struct {
	defaultConf        *SMTPConfig
	regConf            *SMTPConfig
	dialer             SMTPDialer
	defaultDialer      *gomail.Dialer
	registrationDialer *gomail.Dialer
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
	}

	regConf := &SMTPConfig{
		Host:     regHost,
		Port:     regPortNum,
		Username: regUser,
		Password: regPassword,
		From:     regFrom,
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
