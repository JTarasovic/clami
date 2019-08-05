package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/jtarasovic/clami/client"
	"github.com/spf13/cobra"
)

var (
	channelFlag string
	hvmFlag     bool
	regionsFlag []string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "clami",
	Short: "Retrieves the lastest Container Linux AMI information from the CoreOS API",
	Run:   GetAMIInfo,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		glog.Exit(err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&channelFlag, "channel", "stable", "the channel to retrieve ami info from (stable, beta, alpha)")
	RootCmd.PersistentFlags().BoolVar(&hvmFlag, "hvm", true, "the ami type to retieve")
	RootCmd.PersistentFlags().StringSliceVar(&regionsFlag, "regions", []string{}, "the aws region(s) to retieve ami info for")
}

//GetAMIInfo is the base command to retrieve ami info from CoreOS API
func GetAMIInfo(cmd *cobra.Command, args []string) {
	// validate that flags are set
	if len(regionsFlag) == 0 {
		glog.Exitln("at least one region must be specified via the regions flag")
	}

	channel, err := client.ToChannel(channelFlag)
	if err != nil {
		glog.Exitln("channel flag must be one of stable, beta or alpha")
	}

	c := client.New()
	info, err := c.GetAMIs(channel)
	if err != nil {
		glog.Fatalf("error fetching data: %v", err)
	}

	amis := filter(info, regionsFlag, hvmFlag)

	if err = writeInfo(os.Stdout, amis); err != nil {
		glog.Errorf("Error writing data to stdout: %v", err)
	}
}

func filter(c *client.ContainerLinuxAMIResponse, regions []string, hvm bool) []string {
	var out []string
	for _, region := range regions {
		data := c.Regions[region]
		ami := data.Hvm
		if !hvm {
			ami = data.Pv
		}
		out = append(out, ami)
	}
	return out
}

func writeInfo(w io.Writer, amis []string) error {
	for _, ami := range amis {
		if _, err := io.Copy(w, strings.NewReader(ami+"\n")); err != nil {
			return err
		}
	}
	return nil
}
