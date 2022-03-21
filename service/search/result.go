package search

import (
    "crypto/sha256"
)

// Score holds position/score for different locations where this article was found
type Score struct {
    // ElasticSearch calculates this at runtime
    // we can store it and use it for filtering or our own relevancy algorithms
    ElasticSearch int64 `json:"elastic_search"`
    // SERPGoogle position in Google
    SERPGoogle int64 `json:"serp_google"`
    // SERPYandex position in Yandex
    SERPYandex int64 `json:"serp_yandex"`
    // Compound is a system generated score
    Compound int64 `json:"compound"`
}

// Result is a unified search query result
type Result struct {
    ID string `json:"id"`
    // Language used in keywords
    Language string `json:"language"`
    // Query holds keywords used to fetch this result
    Query string `json:"query"`
    // URL of the article
    URL string `json:"url"`

    Score Score `json:"score"`

    Title       string `json:"title"`
    Description string `json:"description"`
    HTML        string `json:"html"`
    Text        string `json:"text"`

    // CreatedAt unix timestamp of when this article appeared in our system
    CreatedAt int64 `json:"created_at"`
    // UpdatedAt unix timestamp of when this article was updated in our system
    UpdatedAt int64 `json:"updated_at"`
    // PublishedAt unix timestamp of when this article was published by the authors
    PublishedAt int64 `json:"published_at"`
}

// NewFullResult appends meta information to the search Result
func NewFullResult(result Result) *ResultFull {
    r := &ResultFull{
        Result:        result,
        IsTranslation: false,
    }

    // u, err := url.Parse(r.URL)
    // if err != nil {
    //
    // }

    r.URLHash = sha256.Sum256([]byte(r.URL))
    // r.Domain = u.Hostname()
    // r.DomainHash = sha256.Sum256([]byte(r.Domain))

    return r
}

type Meta struct {
    Keywords    []string `json:"keywords"`
    Description string   `json:"description"`
}

type ResultFull struct {
    Result

    // Version of the article
    Version int `json:"version"`
    // OriginalLanguage is filled when a search result is translated
    OriginalLanguage string `json:"original_language"`
    // IsTranslation indicates if this is a translated record
    IsTranslation bool `json:"is_translation"`
    // URLHash CRC32 checksum of the article URL
    URLHash [32]byte `json:"url_hash"`
    // Domain name
    Domain Domain `json:"domain"`
    // DomainHash CRC32 checksum of the domain name
    DomainHash [32]byte `json:"domain_hash"`
}

type Domain struct {
    Name       DomainName
    Hash       uint32
    Subdomains []string
    Subdomain  string
    Whois      struct{}
}
type DomainName string
type Domains map[DomainName]Domain
type ResultsPerDomain map[DomainName][]*Result

type Keyword struct {
    Keyword  string
    Language string
}
