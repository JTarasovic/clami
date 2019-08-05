package client

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var channels = []Channel{Alpha, Beta, Stable}

func TestParse(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "stable.json"))
	assert.Nil(t, err)
	i, err := parse(f)
	assert.Nil(t, err)
	assert.Equal(t, ReleaseInfo{Platform: "aws", Version: "2135.5.0", ReleaseDate: "2019-07-02 20:52:42 +0000"}, i.ReleaseInfo)
	assert.Equal(t, Region{Hvm: "ami-6f5d190e", Pv: "ami-e35b1f82"}, i.Regions["us-gov-west-1"])
}

var channelTests = []struct {
	input    string
	expected Channel
	err      error
}{
	{
		input:    "alpha",
		expected: Alpha,
	},
	{
		input:    "beta",
		expected: Beta,
	},
	{
		input:    "stable",
		expected: Stable,
	},
	{
		input:    "nonsense",
		expected: "",
		err:      fmt.Errorf("nonsense"),
	},
}

func TestToChannel(t *testing.T) {
	for _, tt := range channelTests {
		for _, perm := range []bool{true, false} {
			val := tt.input
			if perm {
				val = strings.ToUpper(tt.input)
			}

			t.Run(val, func(t *testing.T) {
				s, err := ToChannel(val)
				assert.Equal(t, tt.expected, s)
				if tt.err == nil {
					assert.Nil(t, err)
				} else {
					assert.Contains(t, err.Error(), tt.err.Error())
				}
			})

		}
	}
}

func TestGetAMIs(t *testing.T) {
	for _, channel := range channels {
		t.Run(string(channel), func(t *testing.T) {
			c := New()
			i, err := c.GetAMIs(channel)
			assert.Nil(t, err)
			assert.Equal(t, "aws", i.ReleaseInfo.Platform)
			assert.NotEmpty(t, i.ReleaseInfo.ReleaseDate)
			assert.NotEmpty(t, i.ReleaseInfo.Version)
			assert.NotEmpty(t, i.Regions["us-gov-west-1"])
		})
	}
}

func TestGetAMIsErrorMock(t *testing.T) {
	m := errorMock{"error from mock"}
	c := New()
	c.WithClient(m)
	r, err := c.GetAMIs(Alpha)
	assert.Nil(t, r)
	assert.Equal(t, m.message, err.Error())
}

func ExampleClient_GetAMIs() {
	c := New()
	info, err := c.GetAMIs(Stable)
	if err != nil {
		panic(err)
	}

	// info now has all of the AMI information
	fmt.Println(info.ReleaseInfo.Platform)

	// Output: aws
}

type errorMock struct {
	message string
}

func (m errorMock) Get(url string) (*http.Response, error) {
	return nil, fmt.Errorf(m.message)
}
