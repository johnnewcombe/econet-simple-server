package cobra

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
	k_Debug      = "Provides debug output to stdout."
)

func init() {

	rootCmd.AddCommand(fileserver)
	rootCmd.AddCommand(monitor)

	fileserver.Flags().IntP("station-id", "s", 254, k_StationId)
	fileserver.Flags().StringP("root-folder", "f", "", k_RootFolder)
	fileserver.Flags().StringP("port", "p", "/dev/econet", k_Port)
	fileserver.Flags().Bool("debug", false, k_Debug)

	if err := fileserver.MarkFlagRequired("root-folder"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error()+".")
		os.Exit(1)
	}

	monitor.Flags().StringP("port", "p", "/dev/econet", k_Port)
	monitor.Flags().Bool("debug", false, k_Debug)

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
