package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	k_Version    = "0.0.1"
	k_StationId  = "Station ID of the connected Piconet device."
	k_Port       = "Serial port device name to access the Piconet device."
	k_RootFolder = "Root folder where the fileserver files are stored."
)

func init() {

	rootCmd.AddCommand(fileserver)
	fileserver.Flags().IntP("station-id", "s", 32, k_StationId)
	fileserver.Flags().StringP("port", "p", "/dev/econet", k_Port)
	fileserver.Flags().StringP("root-folder", "f", "", k_RootFolder)
	//fileserver.MarkFlagRequired("port")
	fileserver.MarkFlagRequired("root-folder")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error()+".")
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "Piconet Fileserver Lite",
	Short: "Simple Econet fileserver for Piconet devices. (c) John Newcombe 2025. Version: " + k_Version,
	Long:  `Piconet Fileserver Lite, a simple single network Econet fileserver for use with Piconet devices.`,
}
