package emails

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// WelcomeData is used to complete the welcome email template
type WelcomeData struct {
	EmailData
	FirstName    string
	LastName     string
	Email        string
	Organization string
	Domain       string
}

// VerifyEmailData is used to complete the verify email template
type VerifyEmailData struct {
	EmailData
	FullName  string
	VerifyURL string
}

// SubscriberEmailData is used to complete the subscriber email template
type SubscriberEmailData struct {
	EmailData
	VerifySubscriberURL string
	OrgName             string
}

// InviteData is used to complete the invite email template
type InviteData struct {
	EmailData
	Email       string
	InviterName string
	OrgName     string
	Role        string
	InviteURL   string
}

// Invite data is used to hold temporary constructor information to compose the invite email
type Invite struct {
	Token     string
	OrgName   string
	Recipient string
	Requestor string
	Role      string
}

// ResetRequestData is used to complete the password reset request email template
type ResetRequestData struct {
	EmailData
	ResetURL string
}

// ResetSuccessData is used to complete the password success request email template
type ResetSuccessData struct {
	EmailData
}

// Email subject lines
const (
	WelcomeRE              = "Welcome to TheOpenLane!"
	VerifyEmailRE          = "Please verify your email address to login to TheOpenLane"
	InviteRE               = "Join Your Teammate %s on TheOpenLane!"
	PasswordResetRequestRE = "TheOpenLane Password Reset - Action Required"
	PasswordResetSuccessRE = "TheOpenLane Password Reset Confirmation"
	InviteBeenAccepted     = "You've been added to an Organization on TheOpenLane"
	Subscribed             = "You've been subscribed to %s"
)

// WelcomeEmail creates a welcome email for a new user
func WelcomeEmail(data WelcomeData) (message *mail.SGMailV3, err error) {
	var text, html string

	data.FirstName = cases.Title(language.AmericanEnglish, cases.NoLower).String(data.FirstName)

	if text, html, err = Render("welcome", data); err != nil {
		return nil, err
	}

	data.Subject = WelcomeRE

	return data.Build(text, html)
}

// SubscribeEmail creates a subscribe email for a new subscriber
func SubscribeEmail(data SubscriberEmailData) (message *mail.SGMailV3, err error) {
	var text, html string

	if text, html, err = Render("subscribe", data); err != nil {
		return nil, err
	}

	data.Subject = fmt.Sprintf(Subscribed, data.OrgName)

	return data.Build(text, html)
}

// VerifyEmail creates an email to verify a user's email address
func VerifyEmail(data VerifyEmailData) (message *mail.SGMailV3, err error) {
	var text, html string

	// we might consider using shortname or alias instead of this but today the email sends whatever is stored exactly in the db
	data.FullName = cases.Title(language.AmericanEnglish, cases.NoLower).String(data.FullName)

	if text, html, err = Render("verify_email", data); err != nil {
		return nil, err
	}

	data.Subject = VerifyEmailRE

	return data.Build(text, html)
}

// InviteEmail creates an email to invite a user to join an organization
func InviteEmail(data InviteData) (message *mail.SGMailV3, err error) {
	var text, html string

	if text, html, err = Render("invite", data); err != nil {
		return nil, err
	}

	data.Subject = fmt.Sprintf(InviteRE, data.InviterName)

	return data.Build(text, html)
}

// InviteAccepted creates an email to invite a user to join an organization
func InviteAccepted(data InviteData) (message *mail.SGMailV3, err error) {
	var text, html string

	if text, html, err = Render("invite_joined", data); err != nil {
		return nil, err
	}

	data.Subject = InviteBeenAccepted

	return data.Build(text, html)
}

// PasswordResetRequestEmail creates a password reset email
func PasswordResetRequestEmail(data ResetRequestData) (message *mail.SGMailV3, err error) {
	var text, html string

	if text, html, err = Render("password_reset_request", data); err != nil {
		return nil, err
	}

	data.Subject = PasswordResetRequestRE

	return data.Build(text, html)
}

// PasswordResetSuccessEmail creates an email to send to users once their password has been reset
func PasswordResetSuccessEmail(data ResetSuccessData) (message *mail.SGMailV3, err error) {
	var text, html string

	if text, html, err = Render("password_reset_success", data); err != nil {
		return nil, err
	}

	data.Subject = PasswordResetSuccessRE

	return data.Build(text, html)
}
