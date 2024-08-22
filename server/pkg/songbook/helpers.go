package songbook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/utils"
)

func GetTopHitSong(query string) (*Song, error) {
	accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.genius.com/search?q=%s", url.QueryEscape(query)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var mappedResult SearchResult
	if err := json.NewDecoder(res.Body).Decode(&mappedResult); err != nil {
		return nil, err
	}

	if len(mappedResult.Response.Hits) == 0 {
		return nil, utils.LogError("no songs were found: %v", res.Body)
	}

	hit := mappedResult.Response.Hits[0].Result
	songPageRes, err := http.Get(hit.Url)
	if err != nil {
		return nil, err
	}
	defer songPageRes.Body.Close()

	doc, err := goquery.NewDocumentFromReader(songPageRes.Body)
	if err != nil {
		return nil, err
	}

	lyrics := getLyrics(doc) // This is the function you should create

	song := &Song{
		Title:  hit.Title,
		Artist: hit.PrimaryArtist.Name,
		Lyrics: lyrics,
	}

	return song, nil
}

func getLyrics(doc *goquery.Document) []Stanza {
	var stanzas []Stanza
	var currentStanza Stanza
	consecutiveBreaks := 0

	doc.Find("div[class^='Lyrics__Container']").Each(func(i int, s *goquery.Selection) {
		s.Contents().Each(func(j int, child *goquery.Selection) {
			// Ignore inline annotations
			if child.Is("span") && strings.Contains(child.AttrOr("class", ""), "InlineAnnotation__Container") {
				return
			}

			// Handle text nodes like [Verse 1] or [Chorus]
			text := strings.TrimSpace(child.Text())
			if strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]") {
				currentStanza.Type = cleanStanzaType(text)
				return
			}

			// Handle <br> tags
			if child.Is("br") {
				consecutiveBreaks++
				if consecutiveBreaks >= 2 {
					if len(currentStanza.Verses) > 0 {
						stanzas = append(stanzas, currentStanza)
						currentStanza = Stanza{}
					}
					consecutiveBreaks = 0
					return
				}
				ensureVerseEnd(&currentStanza)
				return
			}

			// Handle text nodes
			if goquery.NodeName(child) == "#text" {
				addWordsToVerse(child.Text(), &currentStanza)
				return
			}

			// Handle lyrics inside annotations
			if child.Is("a[href^='/']") {
				addWordsToVerse(child.Text(), &currentStanza)
				return
			}

			// Handle italic and bold text
			if child.Is("i") || child.Is("b") {
				addWordsToVerse(child.Text(), &currentStanza)
				return
			}

			// Custom handler
			if child.Is("span[class^='ReferentFragmentDesktop__Highlight']") {
				addWordsToVerse(child.Text(), &currentStanza)
				return
			}

			// Remaining values
			addWordsToVerse(text, &currentStanza)
		})
	})

	if len(currentStanza.Verses) > 0 {
		stanzas = append(stanzas, currentStanza)
	}

	return stanzas
}
