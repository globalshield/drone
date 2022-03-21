package search

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFullResult(t *testing.T) {
	r := Result{
		Language:     "ru",
		Query:        "Вооружённый конфликт на востоке Украины",
		URL:          "https://blah.ru/event_War_in_Donbass/",
		Score:        0,
		SERPPosition: 0,
		Title:        "Вооруженный конфликт на востоке Украины - РИА Новости",
		Description:  "Вооруженный конфликт на востоке Украины. Читайте последние новости на тему в ленте новостей на сайте РИА Новости. Один человек погиб и девять ранены в",
		UpdatedAt:    0,
		CreatedAt:    0,
	}

	fr := NewFullResult(r)
	assert.Equal(t, "blah.ru", fr.Domain)
	assert.Equal(t, "https://blah.ru/event_War_in_Donbass/", fr.URL)
	assert.Equal(t, "75b425bfe8ad215185fabf4a843b64ecd98e1737bdcf9f25437494be47af38c5", hex.EncodeToString(fr.URLHash[:]))
	assert.Equal(t, "4821b984cd7dcb87240c1be5001f1bbb0565555f9ab534e6fb7e982530f7e6d1", hex.EncodeToString(fr.DomainHash[:]))
}
