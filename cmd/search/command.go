package search

import (
    "fmt"
    "github.com/globalshield/drone/service/search"
    "github.com/globalshield/drone/service/search/providers"
    "github.com/globalshield/drone/service/writer"
    "github.com/gosimple/slug"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/spf13/cobra"
    "io/ioutil"
    "strings"
    "time"
)

const (
    FlagOutputFormat        string = "format"
    FlagOutputWriter        string = "writer"
    FlagOutputDir           string = "directory"
    FlagFileSuffix          string = "suffix"
    FlagPermissionOverwrite string = "overwrite"
    FlagPermissionAppend string = "append"
    FlagKeywordFile      string = "import"
    FlagQueryPages       string = "pages"
    FlagQuiet               string = "quiet"
)

type CommandConfig struct {
    OutputFormat        string
    OutputDirectory     string
    OutputWriter        string
    FileSuffix          string
    PermissionOverwrite bool
    PermissionAppend    bool
    QueryImport         string
    QueryPages          int64
    Query               string
    Quiet               bool
    Language            string
}

var Command = func(cmd *cobra.Command, args []string) {
    cfg, err := getConfig(cmd, args)
    if err != nil {
        log.Err(fmt.Errorf("search failed: %w", err)).Send()
        return
    }

    if cfg.Quiet {
        zerolog.SetGlobalLevel(zerolog.NoLevel)
    }

    if len(cfg.QueryImport) > 0 {
        queries, _ := readKeywordFile(cfg.QueryImport)
        for _, query := range queries {
            fmt.Println("search start...")
            performSearch(cfg, query)
            fmt.Println("search end...")
        }
        return
    }

    performSearch(cfg, cfg.Query)
}

func readKeywordFile(filename string) ([]string, error) {
    contents, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    input := strings.Trim(string(contents), "\n")
    queries := strings.Split(input, "\n")

    return queries, nil
}

func performSearch(cfg *CommandConfig, query string) {
    log.Trace().Msgf("query: %s", query)

    fileNameTpl := fileNameTemplate(cfg, query)
    log.Trace().Msgf("filename template: %s", fileNameTpl)

    results := getResults(cfg, query)

    err := write(cfg, fileNameTpl, results)
    if err != nil {
        log.Err(err).Msgf("failed to write results")
    }
}

func write(cfg *CommandConfig, fileNameTpl string, results *search.ResultSet) error {
    config := writer.FileConfig{
        ForceOverwrite: cfg.PermissionOverwrite,
        Append:         cfg.PermissionAppend,
    }
    w := writer.FileWriter{Config: config}
    w.WriterFunc = w.Stdout
    if cfg.OutputWriter == "file" {
        w.WriterFunc = w.File
    }

    var err error
    if cfg.OutputFormat == "csv" {
        err = w.CSV(fileNameTpl, *results)
    }
    if cfg.OutputFormat == "json" {
        err = w.JSON(fileNameTpl, *results)
    }
    return err
}

func fileNameTemplate(cfg *CommandConfig, query string) string {
    normalized := slug.MakeLang(query, cfg.Language)

    suffix := cfg.FileSuffix
    if len(suffix) > 0 {
        suffix = "_" + suffix
    }

    //               %s engine name
    // {dir}/{timestamp}_{query}_%s_{suffix}.{format}
    // we keep %s for later injection
    fileNameTpl := fmt.Sprintf(
        "%s/%s_%s_%%s%s.%s",
        strings.TrimRight(cfg.OutputDirectory, "/"),
       time.Now().Format("2006-01-02_15-04-05"),
        normalized,
        suffix,
        cfg.OutputFormat,
    )
    return fileNameTpl
}

func getResults(cfg *CommandConfig, query string) *search.ResultSet {
    config := providers.GoogleSearchQuery{
        CountryCode:  cfg.QueryImport,
        LanguageCode: cfg.QueryImport,
        Limit:        100,
    }

    // TODO: ideally this should be configurable at runtime
    // and maybe extended with shared libs? idk
    engines := search.NewEngines()
    engines.Add(providers.NewGoogleSearch(config))

    // TODO: retry
    results, errPool := engines.Search(query, cfg.QueryPages)
    log.Trace().Msgf("results: %d, errors: %d", len(*results), len(errPool))

    if len(errPool) > 0 {
        printErrors(errPool)
        return nil
    }

    return results
}

func printErrors(pool search.ErrorPool) {
    for engine, err := range pool {
        log.Err(err).Msgf("error in %s engine query: %s", engine)
    }
}
