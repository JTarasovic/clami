// Package client exposes a Client for accessing the CoreOS provided API for retrieving Container Linux AMI info
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Alpha, Beta and Stable represent the alpha, beta and stable channels respectively
const (
	urlTemplate            = "https://coreos.com/dist/aws/aws-%s.json"
	releaseInfoKey         = "release_info"
	Alpha          Channel = "alpha"
	Beta           Channel = "beta"
	Stable         Channel = "stable"
)

// Getter is an interface for mocking. *Testing only.
type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

// Client exposes the CoreOS API for retrieving Container Linux AMI information
type Client struct {
	getter Getter
}

// ReleaseInfo represents the nested `release_info` object within the response
type ReleaseInfo struct {
	Version     string `mapstructure:"version"`
	ReleaseDate string `mapstructure:"release_date"`
	Platform    string `mapstructure:"platform"`
}

// Region is the generic representation of each `region` object within the response
type Region struct {
	Hvm string `mapstructure:"hvm"`
	Pv  string `mapstructure:"pv"`
}

// ContainerLinuxAMIResponse represents the response from getting CoreOS AMI endpoint
type ContainerLinuxAMIResponse struct {
	ReleaseInfo ReleaseInfo
	Regions     map[string]Region
}

// Channel is the Container Linux release channel
type Channel string

// New returns a new Client with the given configuration
func New() *Client {
	return &Client{getter: http.DefaultClient}
}

// WithClient allows for mocking the CoreOS API
func (c *Client) WithClient(getter Getter) {
	c.getter = getter
}

// GetAMIs requests the full, parsed list of AMIs for the given channel
func (c *Client) GetAMIs(channel Channel) (*ContainerLinuxAMIResponse, error) {
	i, err := c.fetchData(channel)
	if err != nil {
		return nil, err
	}
	defer i.Close()

	return parse(i)
}

func (c *Client) fetchData(channel Channel) (io.ReadCloser, error) {
	url := fmt.Sprintf(urlTemplate, channel)
	resp, err := c.getter.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// parse reads the Reader but does not close it
func parse(r io.Reader) (*ContainerLinuxAMIResponse, error) {
	m := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}

	// pull out the release_info key
	var release ReleaseInfo
	if err := mapstructure.Decode(m[releaseInfoKey], &release); err != nil {
		return nil, err
	}
	delete(m, releaseInfoKey)

	// iterate over the remainder - the regions
	regions := make(map[string]Region)
	for k, v := range m {
		var region Region
		if err := mapstructure.Decode(v.(map[string]interface{}), &region); err != nil {
			return nil, err
		}
		regions[k] = region
	}

	return &ContainerLinuxAMIResponse{ReleaseInfo: release, Regions: regions}, nil
}

// ToChannel converts a raw channel string into the typed representation of that channel or returns an error
func ToChannel(s string) (Channel, error) {
	val := strings.ToLower(s)
	for _, c := range []Channel{Alpha, Beta, Stable} {
		if c == Channel(val) {
			return c, nil
		}
	}
	return "", fmt.Errorf("unable to parse %s as a valid Container Linux release channel", val)
}
