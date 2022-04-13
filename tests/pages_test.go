package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestHomepage(t *testing.T) {
	baseURL := "http://localhost:3000"
	resp, err := http.Get(baseURL + "/")

	assert.NoError(t, err, "Error happen!")
	assert.Equal(t, 200, resp.StatusCode, "should return status code->200")
}
func TestAllPages(t *testing.T) {

	baseURL := "http://localhost:3000"

	var tests = []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/3", 200},
		{"GET", "/articles/3/edit", 200},
		{"POST", "/articles/3", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/1/delete", 404},
	}

	for _, test := range tests {
		t.Logf("URL: %v \n", test.url)
		var (
			resp *http.Response
			err  error
		)
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}
		assert.NoError(t, err, "ERR URL "+test.url)
		assert.Equal(t, test.expected, resp.StatusCode, test.url+" should return "+strconv.Itoa(test.expected))
	}
}
