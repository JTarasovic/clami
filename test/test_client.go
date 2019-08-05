package test

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/jtarasovic/clami/client"
)

// ClientMock implements the Getter interface to mock responses from the CoreOS API
type ClientMock struct {
}

// Get returns a mock response
func (m *ClientMock) Get(url string) (*http.Response, error) {
	f, err := os.Open(filepath.Join("..", "client", "testdata", "stable.json"))
	if err != nil {
		return nil, err
	}
	resp := http.Response{}
	resp.StatusCode = 200
	resp.Body = f
	return &resp, nil
}

// MockData returns mock data directly
func MockData() (*client.ContainerLinuxAMIResponse, error) {
	c := client.New()
	c.WithClient(&ClientMock{})

	return c.GetAMIs(client.Stable)
}
