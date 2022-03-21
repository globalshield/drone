package search

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "sync"
    "time"
)

type Config struct {
    Query string
}

type Search interface {
    ProviderName() string
    Search(config Config) ([]*Result, error)
    SearchPaginated(config Config, page int) ([]*Result, error)
}

type Engines struct {
    engines map[string]Search
    mu      sync.RWMutex
}

func NewEngines() *Engines {
    e := &Engines{
        engines: make(map[string]Search, 0),
        mu:      sync.RWMutex{},
    }

    return e
}

func (e *Engines) Remove(name string) {
    e.mu.Lock()
    delete(e.engines, name)
    e.mu.Unlock()
}

func (e *Engines) Add(engine Search) {
    e.mu.Lock()
    e.engines[engine.ProviderName()] = engine
    e.mu.Unlock()
}

type ResultSet map[string][]*Result
type ErrorPool map[string]error

// Search TODO: simplify this
func (e *Engines) Search(query string, pages int64) (*ResultSet, ErrorPool) {
    e.mu.Lock()

    var wg sync.WaitGroup
    wg.Add(len(e.engines) * int(pages))

    var mu = sync.Mutex{}
    var results = make(ResultSet, 0)
    var errorPool = make(ErrorPool, 0)
    doSearch := func(name string, engine Search) {
        for page := 0; page < int(pages); page++ {
            log.Trace().Msgf("fetching page %d for %s", page+1, engine.ProviderName())

            searchResults, err := engine.SearchPaginated(Config{Query: query}, page)
            time.Sleep(200 * time.Millisecond) // TODO: needs rate limiting
            if err != nil {
                err = fmt.Errorf("failed to fetch page %d: %w", page+1, err)
                log.Err(err).Send()
                mu.Lock()
                errorPool[name] = err
                mu.Unlock()
                return
            }

            mu.Lock()
            if _, ok := results[engine.ProviderName()]; !ok {
                results[engine.ProviderName()] = make([]*Result, 0)
            }
            mu.Unlock()

            mu.Lock()
            results[engine.ProviderName()] = append(results[engine.ProviderName()], searchResults...)
            mu.Unlock()

            wg.Done() // finish work
        }
    }

    for name, engine := range e.engines {
        go doSearch(name, engine)
    }

    wg.Wait()
    e.mu.Unlock()

    return &results, errorPool
}
