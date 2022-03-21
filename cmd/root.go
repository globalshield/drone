package cmd

import (
    "github.com/globalshield/drone/cmd/report"
    "github.com/globalshield/drone/cmd/search"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/pkgerrors"
    "github.com/spf13/cobra"
    "os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   `drone`,
    Short: `Collects propaganda and desinformation sources`,
    Long:  `Collects propaganda and desinformation sources`,
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

var v string
func init() {
    defaultFlags()
    logging()
    addCommands()
}

func addCommands() {
    rootCmd.AddCommand(search.GetCommand())
    rootCmd.AddCommand(report.GetCommand())
}

func logging() {
    rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
        zerolog.SetGlobalLevel(zerolog.TraceLevel)
        zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
        return nil
    }
}

func defaultFlags() {
    rootCmd.PersistentFlags().
        StringVarP(
            &v,
            "verbosity",
            "v",
            "trace",
            "log level (trace, debug, info, warn, error, fatal, panic)",
        )
    rootCmd.PersistentFlags().
        StringVarP(
            &v,
            "lang",
            "l",
            "ru",
            "language for search engine configuration",
        )
}
