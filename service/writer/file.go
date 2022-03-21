package writer

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/globalshield/drone/service/search"
    "github.com/jszwec/csvutil"
    "github.com/rs/zerolog/log"
    "io/ioutil"
    "os"
)

type Writer interface {
    Write([]byte) error
}

type FileConfig struct {
    Name           string
    ForceOverwrite bool
    Append         bool
}

type Func func (filename string, data []byte, cfg FileConfig, json bool) error

type FileWriter struct {
    Config     FileConfig
    WriterFunc Func
}

func (fw FileWriter) CSV(filePath string, results search.ResultSet) error {
    for engine, entries := range results {
        filename := fmt.Sprintf(filePath, engine)

        csv, err := csvutil.Marshal(entries)
        if err != nil {
            err = fmt.Errorf("failed to marshal to CSV: %s: %w", filePath, err)
            log.Err(err).Msg("")
            return err
        }

        // if file already exists
        // append CSV without the header
        if _, err := os.Stat(filename); err == nil {
            sep := []byte("\n")
            index := bytes.Index(csv, sep)
            length := len(csv)
            if fw.Config.Append && length >= index+1 {
                csv = csv[index+1:] // remove header
            }
            if bytes.HasSuffix(csv, sep) {
                csv = bytes.TrimRight(csv, "\n")
            }
        }

        err = fw.WriterFunc(filename, csv, fw.Config, false)
        if err != nil {
            err = fmt.Errorf("failed to write CSV: %s: %w", filePath, err)
            log.Err(err).Msg("")
            return err
        }
    }

    return nil
}

func (fw FileWriter) JSON(filePath string, results search.ResultSet) error {
    for engine, entries := range results {
        filename := fmt.Sprintf(filePath, engine)

        // on append, read existing JSON and add array items
        if _, err := os.Stat(filename); err == nil && fw.Config.Append {
            file, err := ioutil.ReadFile(filename)
            if err != nil {
                return fmt.Errorf("failed to read JSON for append: %w", err)
            }
            var existingResults []*search.Result
            err = json.Unmarshal(file, &existingResults)
            if err != nil {
                return fmt.Errorf("failed to unmarshal JSON for append: %w", err)
            }

            entries = append(entries, existingResults...)
        }

        jsonObj, err := json.MarshalIndent(entries, "", "  ")
        if err != nil {
            log.Err(err).Msgf("failed to marshal JSON")
            continue
        }

        err = fw.WriterFunc(filename, jsonObj, fw.Config, true)
        if err != nil {
            log.Err(err).Msgf("failed to write file %s", filePath)
            continue
        }
    }
    return nil
}

func (fw FileWriter) Stdout(filename string, data []byte, cfg FileConfig, json bool) error {
    fmt.Print(string(data), "\n")

    return nil
}

func (fw FileWriter) File(filename string, data []byte, cfg FileConfig, json bool) error {
    var flag int
    flag = os.O_WRONLY | os.O_CREATE
    if _, err := os.Stat(filename); err == nil {
        // TODO: refactor. JSON check here is a hack
        if json && cfg.Append {
            flag = os.O_WRONLY | os.O_CREATE
        } else if cfg.Append {
            flag = os.O_APPEND | os.O_WRONLY | os.O_CREATE
        }

        if !cfg.ForceOverwrite && !cfg.Append {
            return fmt.Errorf("file already exists, use -f to overwrite")
        }
    }

    f, err := os.OpenFile(filename, flag, 0644)
    if err != nil {
        return fmt.Errorf("failed to open %s: %w", filename, err)
    }
    defer f.Close()

    if _, err = f.Write(data); err != nil {
        return fmt.Errorf("failed to write %s: %w", filename, err)
    }

    log.Trace().Msgf("finished writing %s", filename)
    return nil
}
