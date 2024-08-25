package gravatar_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theopenlane/utils/gravatar"
)

func TestGravatar(t *testing.T) {
	email := "sfunk@theopenlane.io"
	url := gravatar.New(email, nil)
	require.Equal(t, "https://www.gravatar.com/avatar/4feb2f12c4100528b14f3bed598a9598?d=robohash&r=pg&s=80", url)
}

func TestHash(t *testing.T) {
	// Test case from: https://en.gravatar.com/site/implement/hash/
	input := "sfunk@theopenlane.io"
	expected := "4feb2f12c4100528b14f3bed598a9598"
	require.Equal(t, expected, gravatar.Hash(input))
}
