//go:build applicationtest

package applicationtest

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-tstr/golden"
	"github.com/stretchr/testify/require"
)

func TestTasks(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, appURL+"/api/v1/tasks", strings.NewReader(`{"name":"test task 1"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	golden.Request(t, http.DefaultClient, req, 200)
}
