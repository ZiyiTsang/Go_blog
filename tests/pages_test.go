package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHomepage(t *testing.T) {
	baseURL := "http://localhost:3000"
	resp, err := http.Get(baseURL + "/")

	assert.NoError(t, err, "Error happen!")
	assert.Equal(t, 200, resp.StatusCode, "should return status code->200")
}
