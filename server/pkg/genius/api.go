package genius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"

	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/songbook"
	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/utils"
)

func GetTopHitSong(query string) (*songbook.Song, error) {
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

	lyrics := GetLyrics(doc) // This is the function you should create

	song := &songbook.Song{
		Title:  hit.Title,
		Artist: hit.PrimaryArtist.Name,
		Lyrics: lyrics,
	}

	return song, nil
}

func GetLyrics(doc *goquery.Document) []songbook.Stanza {
	var stanzas []songbook.Stanza
	var currentStanza songbook.Stanza
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
						currentStanza = songbook.Stanza{}
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

func addWordsToVerse(text string, stanza *songbook.Stanza) {
	words := splitWords(text)

	if len(stanza.Verses) == 0 {
		stanza.Verses = append(stanza.Verses, songbook.Verse{})
	}
	currentVerse := &stanza.Verses[len(stanza.Verses)-1]

	for _, word := range words {
		*currentVerse = append(*currentVerse, songbook.Word{Text: word, Chord: []string{}})
	}
}

func ensureVerseEnd(stanza *songbook.Stanza) {
	if len(stanza.Verses) > 0 && len(stanza.Verses[len(stanza.Verses)-1]) > 0 {
		stanza.Verses = append(stanza.Verses, songbook.Verse{})
	}
}

func splitWords(text string) []string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	var result []string
	for _, word := range words {
		subWords := strings.FieldsFunc(word, unicode.IsPunct)
		for i, subWord := range subWords {
			if subWord != "" {
				result = append(result, subWord)
				if i < len(subWords)-1 {
					result = append(result, string(word[len(subWord)]))
				}
			}
		}
	}

	return words
}

func cleanStanzaType(text string) string {
	cleanText := strings.Trim(text, "[]")
	return strings.TrimSpace(cleanText)
}
