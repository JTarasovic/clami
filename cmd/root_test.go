package cmd

import (
	"bytes"
	"testing"

	"github.com/jtarasovic/clami/client"
	"github.com/jtarasovic/clami/test"
	"github.com/stretchr/testify/assert"
)

var writeTests = []struct {
	name   string
	input  []string
	output string
}{
	{
		name:   "single",
		input:  []string{"a"},
		output: "a\n",
	},
	{
		name:   "multiple",
		input:  []string{"a", "b", "c"},
		output: "a\nb\nc\n",
	},
	{
		name:   "empty",
		input:  []string{},
		output: "",
	},
}

func TestWriteInfo(t *testing.T) {
	for _, tt := range writeTests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			writeInfo(&buf, tt.input)
			assert.Equal(t, tt.output, buf.String())
		})
	}
}

var filterTests = []struct {
	name    string
	regions []string
	hvms    []string
	pvs     []string
}{
	{
		name:    "single region",
		regions: []string{"us-gov-east-1"},
		hvms:    []string{"ami-0dc23aad3fa5a13c9"},
		pvs:     []string{""},
	},
	{
		name:    "multiple region",
		regions: []string{"us-gov-east-1", "ap-southeast-1"},
		hvms:    []string{"ami-0dc23aad3fa5a13c9", "ami-0bdf64786279efbbc"},
		pvs:     []string{"", "ami-06b44b502238c091d"},
	},
}

func TestFilter(t *testing.T) {
	i, err := test.MockData()
	assert.Nil(t, err)
	for _, tt := range filterTests {
		t.Run(tt.name, func(t *testing.T) {
			hvms := filter(i, tt.regions, true)
			assert.Equal(t, tt.hvms, hvms)
			pvs := filter(i, tt.regions, false)
			assert.Equal(t, tt.pvs, pvs)
		})
	}
}

var expected = &client.ContainerLinuxAMIResponse{
	ReleaseInfo: client.ReleaseInfo{
		Version:     "2135.5.0",
		ReleaseDate: "2019-07-02 20:52:42 +0000",
		Platform:    "aws",
	},
	Regions: map[string]client.Region{
		"ap-northeast-1": client.Region{Hvm: "ami-070d50353dfb032ba", Pv: "ami-09673be92fa59a18c"},
		"ap-northeast-2": client.Region{Hvm: "ami-041a583a9761bda64", Pv: ""},
		"ap-south-1":     client.Region{Hvm: "ami-083eec4a98ca0396b", Pv: ""},
		"ap-southeast-1": client.Region{Hvm: "ami-0bdf64786279efbbc", Pv: "ami-06b44b502238c091d"},
		"ap-southeast-2": client.Region{Hvm: "ami-0bb7c56044b64aa56", Pv: "ami-0cc32476b1b941cf3"},
		"ca-central-1":   client.Region{Hvm: "ami-082a1a74cfc2d2403", Pv: ""},
		"cn-north-1":     client.Region{Hvm: "ami-0d8ca8372e3b0aff4", Pv: "ami-0fe0dc6001c982cb6"},
		"cn-northwest-1": client.Region{Hvm: "ami-049ed451bb483d4be", Pv: ""},
		"eu-central-1":   client.Region{Hvm: "ami-0cfac31dd01a5f898", Pv: "ami-051b84d3e0a89fec0"},
		"eu-north-1":     client.Region{Hvm: "ami-009c476af4072d56a", Pv: ""},
		"eu-west-1":      client.Region{Hvm: "ami-053d1b6039e1098d4", Pv: "ami-0898e2390ed497160"},
		"eu-west-2":      client.Region{Hvm: "ami-09e2e4b79ea105d0f", Pv: ""},
		"eu-west-3":      client.Region{Hvm: "ami-0a409979da233373a", Pv: ""},
		"sa-east-1":      client.Region{Hvm: "ami-0b2f9ee1da741ad19", Pv: "ami-0ef096d9aa2909669"},
		"us-east-1":      client.Region{Hvm: "ami-02b51824b39a1d52a", Pv: "ami-01d492ec136ec8359"},
		"us-east-2":      client.Region{Hvm: "ami-03aa12465ead76468", Pv: ""},
		"us-gov-east-1":  client.Region{Hvm: "ami-0dc23aad3fa5a13c9", Pv: ""},
		"us-gov-west-1":  client.Region{Hvm: "ami-6f5d190e", Pv: "ami-e35b1f82"},
		"us-west-1":      client.Region{Hvm: "ami-04a1dd7b81fe80e40", Pv: "ami-084c9acb389f1801b"},
		"us-west-2":      client.Region{Hvm: "ami-071f4352a744b29aa", Pv: "ami-0108b87fd991ef10e"}},
}

func TestGetAMIInfo(t *testing.T) {
	c := client.New()
	c.WithClient(&test.ClientMock{})
	resp, err := c.GetAMIs(client.Stable)
	assert.Nil(t, err)
	assert.Equal(t, expected, resp)
}
