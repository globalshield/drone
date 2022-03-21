package search

import (
    "fmt"
    "os"
    "strings"

    "github.com/globalshield/drone/cmd/global"
    "github.com/spf13/cobra"
)

func getConfig(cmd *cobra.Command, args []string) (*CommandConfig, error) {
    outputFormat, err := cmd.Flags().GetString(FlagOutputFormat)
    if err != nil {
        return nil, fmt.Errorf("error reading [-f] output format: %w", err)
    }
    if outputFormat != "csv" && outputFormat != "json" {
        return nil, fmt.Errorf("invalid output [-f] format (supported: json, csv): %w", err)
    }

    writer, err := cmd.Flags().GetString(FlagOutputWriter)
    if err != nil {
        return nil, fmt.Errorf("error reading [-w] writer flag: %w", err)
    }
    if writer != "stdout" && writer != "file" {
        return nil, fmt.Errorf("invalid output [-w] writer (supported: stdout, file): %w", err)
    }

    fileSuffix, err := cmd.Flags().GetString(FlagFileSuffix)
    if err != nil {
        return nil, fmt.Errorf("error parsing [-s] output suffix: %w", err)
    }

    outputDirectory, err := cmd.Flags().GetString(FlagOutputDir)
    if err != nil {
        return nil, fmt.Errorf("error reading [-d] output directory: %w", err)
    }

    keywordFile, err := cmd.Flags().GetString(FlagKeywordFile)
    if err != nil {
        return nil, fmt.Errorf("error reading [-i] import flag: %w", err)
    }
    if len(keywordFile) > 0 {
        if _, err := os.Stat(keywordFile); err != nil && os.IsNotExist(err) {
            return nil, fmt.Errorf("error reading [-i] keyword file `%s`: %w", keywordFile, err)
        }
    }

    permissionsAppend, err := cmd.Flags().GetBool(FlagPermissionAppend)
    if err != nil {
        return nil, fmt.Errorf("error reading [-a] append flag: %w", err)
    }

    permissionsOverwrite, err := cmd.Flags().GetBool(FlagPermissionOverwrite)
    if err != nil {
        return nil, fmt.Errorf("error reading [-f] force overwrite flag: %w", err)
    }

    disableLogging, err := cmd.Flags().GetBool(FlagQuiet)
    if err != nil {
        return nil, fmt.Errorf("error reading [-q] quiet flag: %w", err)
    }

    language, err := cmd.Flags().GetString(global.FlagLanguage)
    if err != nil {
        return nil, fmt.Errorf("error reading [-l] language flag: %w", err)
    }

    pageCount, err := cmd.Flags().GetInt64(FlagQueryPages)
    if err != nil {
        return nil, fmt.Errorf("error reading [-p] page count: %w", err)
    }

    err = os.MkdirAll(outputDirectory, os.ModePerm)
    if err != nil {
        return nil, fmt.Errorf("error making directory: %w", err)
    }

    var query string
    if len(args) > 0 {
        query = strings.Join(args, " ")
    }

    if len(args) == 0 || query == "" {
        // TODO: maybe a good default here?
        // query = "демилитаризация и денацизация"
    }

    cfg := &CommandConfig{
        OutputFormat:    outputFormat,
        OutputDirectory: outputDirectory,
        OutputWriter:    writer,

        FileSuffix: fileSuffix,

        PermissionOverwrite: permissionsOverwrite,
        PermissionAppend:    permissionsAppend,

        Language: language,

        Query:       query,
        QueryImport: keywordFile,
        QueryPages:  pageCount,

        Quiet: disableLogging,
    }

    return cfg, nil
}
