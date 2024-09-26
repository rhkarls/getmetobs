package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	outputPath string
	version    string
	ext        string
	rootCmd    = &cobra.Command{
		Use:   "getmetobs <parameter> <station> <period>",
		Short: "Download SMHI meteorological observation data",
		Long: `getmetobs is a CLI tool that downloads meteorological observation data provided by SMHI.
Data is downloaded and saved in the specified directory (or current directory if not specified) with a standardized filename.`,
		Args:    cobra.ExactArgs(3),
		Run:     run,
		Example: "getmetobs 1 159880 latest-day --output /path/to/directory",
	}
)

func init() {
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "Directory to save the downloaded file (default is current directory)")
	rootCmd.Flags().StringVarP(&version, "version", "v", "1.0", "API version to use")
	rootCmd.Flags().StringVarP(&ext, "ext", "e", "csv", "File extension for the data (e.g., csv, json)")

	// Add help messages for arguments
	rootCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

Arguments:
  parameter    The meteorological parameter to retrieve provided as integer ID, see https://opendata.smhi.se/apidocs/metobs/parameter.html
  station      The weather station identifier as integer ID, see https://www.smhi.se/data/meteorologi/ladda-ner-meteorologiska-observationer
  period       One of four Periods. Valid values are latest-hour, latest-day, latest-months or corrected-archive. Notice that all Stations do not have all four Periods so make sure to check which ones are available in the Period level. 
`)
}

func run(cmd *cobra.Command, args []string) {
	parameter := args[0]
	station := args[1]
	period := args[2]

	fullURL := fmt.Sprintf("https://opendata-download-metobs.smhi.se/api/version/%s/parameter/%s/station/%s/period/%s/data.%s",
		version, parameter, station, period, ext)

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status: %s\nCheck argument values\n", resp.Status)
		os.Exit(1)
	}

	// Make sure that the output directory exists, then create the file
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	outputFilePath := filepath.Join(outputPath, fmt.Sprintf("smhi_metobs_%s_%s_%s.%s", parameter, station, period, ext))
	out, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	// Write the contents to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File downloaded successfully: %s\n", outputFilePath)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
