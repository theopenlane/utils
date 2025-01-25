//go:build darwin
// +build darwin

package keyring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestKeying(t *testing.T) {
	keyring.MockInit()

	service := "test-service-1"
	account := "test-account-1"
	passwd := "test-password"

	// cleanup
	defer DeleteKeyring(service, account) //nolint:errcheck

	err := SetKeying(service, account, []byte(passwd))
	require.NoError(t, err)

	result, err := QueryKeyring(service, account)
	require.NoError(t, err)

	assert.Equal(t, passwd, result)

	// update password
	updatedPasswd := "updated-password"
	err = SetKeying(service, account, []byte(updatedPasswd))
	require.NoError(t, err)

	result, err = QueryKeyring(service, account)
	require.NoError(t, err)

	assert.Equal(t, updatedPasswd, result)

	// bad query
	result, err = QueryKeyring("meows", "eddy")
	require.Error(t, err)
	assert.Empty(t, result)
}
