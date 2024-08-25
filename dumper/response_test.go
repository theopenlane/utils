package dumper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	// TODO: look into updating the echox dependency
	echo "github.com/theopenlane/echox"
)

func TestDumper(t *testing.T) {
	responseString := gofakeit.SentenceSimple()
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		respDumper := NewDumper(c.Response())
		c.Response().Writer = respDumper

		defer func() {
			require.Equal(t, respDumper.GetResponse(), responseString)
		}()

		return c.String(http.StatusOK, responseString)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	e.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), responseString)
}
