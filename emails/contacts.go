package emails

import (
	"github.com/theopenlane/utils/sendgrid"
)

// AddContact adds a contact to SendGrid, adding them to the sign-up marketing list if
// it is configured. This is an upsert operation so existing contacts will be updated.
// The caller can optionally specify additional lists that the contact should be added
// to. If no lists are configured or specified, then the contact is added or updated in
// SendGrid but is not added to any marketing lists - the intent of this is to track within SendGrid the
// sign-ups and track PLG-related stuff
func (m *EmailManager) AddContact(contact *sendgrid.Contact, listIDs ...string) (err error) {
	if !m.Enabled() {
		return ErrSendgridNotEnabled
	}

	// Setup the request data
	sgdata := &sendgrid.AddContactData{
		Contacts: []*sendgrid.Contact{contact},
	}

	// Add the contact to the specified marketing lists
	if m.conf.ListID != "" {
		sgdata.ListIDs = append(sgdata.ListIDs, m.conf.ListID)
	}

	for _, listID := range listIDs {
		if listID != "" {
			sgdata.ListIDs = append(sgdata.ListIDs, listID)
		}
	}

	// Invoke the SendGrid API to add the contact
	if err = sendgrid.AddContacts(m.conf.SendGridAPIKey, sgdata); err != nil {
		return err
	}

	return nil
}
