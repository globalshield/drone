package providers

import (
    "github.com/globalshield/drone/service/search"
    googleSearch "github.com/rocketlaunchr/google-search"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestSearchPaginate(t *testing.T) {
    var testResults = []googleSearch.Result{
        {
            Rank:        1,
            URL:         "url",
            Title:       "title",
            Description: "desc",
        },
    }

    queryGoogleSearchFunc = func(query GoogleSearchQuery, cfg search.Config, currentPage int) ([]googleSearch.Result, error) {
        return testResults, nil
    }

    query := "демилитаризация и денацизация"

    gs := NewGoogleSearch(GoogleSearchQuery{
        Limit:        10,
        CountryCode:  "ru",
        LanguageCode: "ru",
    })

    results, err := gs.Search(search.Config{Query: query})
    assert.NoError(t, err)
    assert.Equal(t, results[0].URL, testResults[0].URL)
    assert.Equal(t, results[0].Score, int64(testResults[0].Rank))
    assert.Equal(t, results[0].SERPPosition, int64(testResults[0].Rank))
    assert.Equal(t, results[0].Description, testResults[0].Description)
    assert.Equal(t, results[0].Title, testResults[0].Title)
    assert.NotZero(t, results[0].CreatedAt)
}
