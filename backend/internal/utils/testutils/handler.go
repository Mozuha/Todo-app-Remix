package testutils

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://github.com/hack-31/point-app-backend/blob/main/testutil/handler.go#L23
func AssertResponse(t *testing.T, actualResponse *http.Response, wantStatus int, wantBody []byte) {
	t.Helper()
	t.Cleanup(func() { _ = actualResponse.Body.Close() })
	arb, err := io.ReadAll(actualResponse.Body)
	assert.NoError(t, err)

	assert.Equalf(t, wantStatus, actualResponse.StatusCode, "want http.status %d, but got %d", wantStatus, actualResponse.StatusCode)

	if len(arb) == 0 && len(wantBody) == 0 {
		return
	}

	assert.JSONEq(t, string(wantBody), string(arb))
}

// https://khigashigashi.hatenablog.com/entry/2019/04/27/150230
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bs, err := os.ReadFile(path)
	assert.NoErrorf(t, err, "cannot read from %q", path)
	return bs
}
