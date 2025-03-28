package marionette_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theopenlane/utils/marionette"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		conf marionette.Config
		err  error
	}{
		{marionette.Config{}, marionette.ErrNoWorkers},
		{marionette.Config{Workers: 4}, marionette.ErrNoServerName},
		{marionette.Config{Workers: 4, ServerName: "marionette"}, nil},
	}

	for i, tc := range testCases {
		err := tc.conf.Validate()
		require.ErrorIs(t, err, tc.err, "test case %d failed", i)
	}
}
