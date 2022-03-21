package report

import "github.com/spf13/cobra"

var Domains = func(cmd *cobra.Command, args []string) {}
var Phrases = func(cmd *cobra.Command, args []string) {}

var Command = func(cmd *cobra.Command, args []string) {

}

func GetCommand() *cobra.Command {
    var reportCmd = &cobra.Command{
        Use: "report",
        Short: "Aggregated search result reports",
        Long:  "Aggregated search result reports",
        Run:   Command,
    }
    var domainsCmd = &cobra.Command{
        Use:   "domain",
        Short: "Domain reports",
        Long:  "Domain reports",
        Run:   Domains,
    }
    var phrasesCmd = &cobra.Command{
        Use:   "keyword",
        Short: "Phrase and keyword reports",
        Long:  "Phrase and keyword reports",
        Run:   Phrases,
    }

    reportCmd.AddCommand(domainsCmd)
    reportCmd.AddCommand(phrasesCmd)

    return reportCmd
}
