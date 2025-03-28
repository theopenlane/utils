package keyring

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

var (
	// ErrNoResultsFound is returned when no results are found in the keyring
	ErrNoResultsFound = errors.New("no results found")
)

// QueryKeyring queries the keyring for the first item with the given service and account
// and return the password as a string
func QueryKeyring(service, account string) (string, error) {
	result, err := keyring.Get(service, account)
	if err != nil {
		return "", err
	}

	if result == "" {
		return "", fmt.Errorf("%w for account %s", ErrNoResultsFound, account)
	}

	return result, nil
}

// SetKeying sets the keyring with the given service, account and secret
func SetKeying(service, account string, secret []byte) error {
	if err := keyring.Set(service, account, string(secret)); err != nil {
		return err
	}

	return nil
}

// DeleteKeyring deletes the keyring with the given service and account
func DeleteKeyring(service, account string) error {
	return keyring.Delete(service, account)
}
