package marionette_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/theopenlane/utils/marionette"
)

func TestTaskErrors(t *testing.T) {
	werr := errors.New("significant badness happened") // nolint: err113
	err := marionette.Errorw(werr)

	require.ErrorIs(t, errors.Unwrap(err), werr, "expected to be able to unwrap an error")
	require.ErrorIs(t, err, werr, "expected the error to wrap an error")
	require.EqualError(t, err, "after 0 attempts: significant badness happened")

	// Append some errors
	err.Append(errors.New("could not reach database"))      // nolint: err113
	err.Append(errors.New("could not reach database"))      // nolint: err113
	err.Append(errors.New("could not reach database"))      // nolint: err113
	err.Append(errors.New("failed precondition"))           // nolint: err113
	err.Append(errors.New("maximum backoff limit reached")) // nolint: err113

	err.Since(time.Now().Add(-10 * time.Second))
	require.EqualError(t, err, "after 5 attempts: significant badness happened")
}

func TestNilTaskError(t *testing.T) {
	err := &marionette.Error{}

	require.Nil(t, errors.Unwrap(err))
	require.EqualError(t, err, "task failed after 0 attempts")

	// Append some errors
	err.Append(errors.New("could not reach database"))      // nolint: err113
	err.Append(errors.New("could not reach database"))      // nolint: err113
	err.Append(errors.New("could not reach database"))      // nolint: err113
	err.Append(errors.New("failed precondition"))           // nolint: err113
	err.Append(errors.New("maximum backoff limit reached")) // nolint: err113

	err.Since(time.Now().Add(-10 * time.Second))
	require.EqualError(t, err, "task failed after 5 attempts")
}
