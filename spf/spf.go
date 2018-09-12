package main

import (
	"fmt"
	"net"
	"os"

	spf "github.com/femaref/spf-1"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "spf",
	Short: "",
	Long:  ``,
}

func init() {
	RootCmd.AddCommand(validateCmd)
	RootCmd.AddCommand(extractCmd)
}

var validateCmd = &cobra.Command{
	Use:  "validate",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(spf.CheckHost(net.ParseIP("127.0.0.1"), args[0], "postmaster@"))
	},
}

var extractCmd = &cobra.Command{
	Use:  "extract",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(spf.ExtractAllowedHosts(net.ParseIP("127.0.0.1"), args[0], "postmaster@"))
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
