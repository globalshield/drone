package search

import (
    "github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
    // searchCmd represents the search command
    var searchCmd = &cobra.Command{
        Use:       "search",
        ValidArgs: []string{"keywords"},
        Short:     "Search multiple search engines",
        Long:      `Search multiple search engines`,
        Run:       Command,
    }

    // output config
    searchCmd.Flags().StringP(FlagOutputFormat, "o", "json", "output format (csv, json)")
    searchCmd.Flags().StringP(FlagOutputWriter, "w", "stdout", "output writer (stdout, file)")
    // file config: --writer file
    searchCmd.Flags().StringP(FlagFileSuffix, "s", "", "output file suffix {query}_{engine}_{suffix}.{ext}")
    searchCmd.Flags().StringP(FlagOutputDir, "d", "./output", "directory path for search results")
    // file permissions
    searchCmd.Flags().BoolP(FlagPermissionOverwrite, "f", false, "overwrite an existing file")
    searchCmd.Flags().BoolP(FlagPermissionAppend, "a", false, "append results to an existing file")
    // keyword import
    searchCmd.Flags().StringP(FlagKeywordFile, "i", "", "newline separated keyword file")
    // search config
    searchCmd.Flags().Int64P(FlagQueryPages, "p", 1, "specify how many pages to scroll through")
    // logging
    searchCmd.Flags().BoolP(FlagQuiet, "q", false, "disable logging")

    return searchCmd
}
