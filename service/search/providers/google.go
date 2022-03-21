package providers

import (
    "context"
    "fmt"
    "github.com/globalshield/drone/service/search"
    googleSearch "github.com/rocketlaunchr/google-search"
    "time"
)

type GoogleSearchQuery struct {
    CountryCode  string
    LanguageCode string
    Limit        int
}

// GoogleSearch google search wrapper
// AFAIK Google does not provide a search API therefore we need to crawl and rely on proxies
// and pretty sure Google will occasionally return different results
// which we can benefit from by implementing different wrappers for different "ways" to query Google
// the more data we collect the better because we need good accuracy and to cover corner cases (most likely related to linguistics, see below)
type GoogleSearch struct {
    query GoogleSearchQuery
}

func NewGoogleSearch(query GoogleSearchQuery) *GoogleSearch {
    return &GoogleSearch{query: query}
}

func (g GoogleSearch) ProviderName() string {
    return "google_crawl"
}

var _ search.Search = (*GoogleSearch)(nil)

// UserAgent TODO: randomize this
var UserAgent = "Mozilla/5.0 (Windows NT 11.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4737.123 Safari/537.36"

func (g GoogleSearch) Search(cfg search.Config) ([]*search.Result, error) {
    results, err := g.searchGoogle(cfg, 0)
    if err != nil {
        return nil, err
    }

    mapped, err := g.mapResults(cfg, results)
    if err != nil {
        return nil, err
    }

    return mapped, nil
}

func (g GoogleSearch) SearchPaginated(cfg search.Config, currentPage int) ([]*search.Result, error) {
    results, err := g.searchGoogle(cfg, currentPage)
    // spew.Dump(results)
    if err != nil {
        return nil, err
    }

    mapped, err := g.mapResults(cfg, results)
    if err != nil {
        return nil, err
    }
    return mapped, nil
}

func (g GoogleSearch) mapResults(cfg search.Config, results []googleSearch.Result) ([]*search.Result, error) {
    var items = make([]*search.Result, 0)
    for _, v := range results {
        items = append(items, &search.Result{
            Query:        cfg.Query,
            Language:     g.query.LanguageCode,
            URL:          v.URL,
            Score:        search.Score{SERPGoogle: int64(v.Rank)},
            Title:        v.Title,
            Description:  v.Description,
            UpdatedAt:    time.Now().Unix(),
            CreatedAt:    time.Now().Unix(),
        })
    }

    return items, nil
}

// queryGoogleSearchFunc to simplify testing
var queryGoogleSearchFunc = func(query GoogleSearchQuery, cfg search.Config, currentPage int) ([]googleSearch.Result, error) {
    googleSearch.RateLimit.SetLimit(3) // per second
    googleSearch.RateLimit.SetBurst(3) // per second
    results, err := googleSearch.Search(context.Background(), cfg.Query, googleSearch.SearchOptions{
        CountryCode:  query.CountryCode,
        LanguageCode: query.LanguageCode,
        Limit:        query.Limit,
        Start:        currentPage * query.Limit,
        UserAgent:    UserAgent,
        OverLimit:    false,
    })

    return results, err
}

func (g GoogleSearch) searchGoogle(cfg search.Config, currentPage int) ([]googleSearch.Result, error) {
    results, err := queryGoogleSearchFunc(g.query, cfg, currentPage)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch google results: %w", err)
    }

    return results, nil
}
