package genius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/songbook"
	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/utils"
)

func GetTopHitSong(query string) (*songbook.Song, error) {
	accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")
	utils.LogDebug(accessToken)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.genius.com/search?q=%s", query), nil)
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
		return nil, fmt.Errorf("no songs were found: %v\n", res.Body)
	}

	hit := mappedResult.Response.Hits[0].Result
	lyricsPageRes, err := http.Get(hit.Url)
	if err != nil {
		return nil, err
	}
	defer lyricsPageRes.Body.Close()

	doc, err := goquery.NewDocumentFromReader(lyricsPageRes.Body)
	if err != nil {
		return nil, err
	}

	lyrics := GetLyrics(doc)

	song := &songbook.Song{
		Title:  hit.Title,
		Artist: hit.PrimaryArtist.Name,
		Lyrics: lyrics,
	}

	return song, nil
}

func GetLyrics(doc *goquery.Document) []songbook.Stanza {
	var stanzaList []songbook.Stanza
	var currentStanza songbook.Stanza
	var stanzaType string

	doc.Find(`div[class^="Lyrics__Container"]`).Each(func(_ int, wrapper *goquery.Selection) {
		wrapper.Children().Children().Children().Each(func(i int, s *goquery.Selection) {
			if s.Is("br") {
				if stanzaType == "" {
					return
				}
				stanzaList = append(stanzaList, currentStanza)
				currentStanza = songbook.Stanza{}
				stanzaType = ""
				return
			}

			verse := handleTextNodes(s)
			if verse != nil {
				currentStanza.Verses = append(currentStanza.Verses, verse)
			}

			linked := handleLinkedTextNodes(s)
			if linked != nil {
				currentStanza.Verses = append(currentStanza.Verses, linked...)
			}

			if stanzaType == "" {
				stanzaType = getType(s.Text())
			}
		})
	})

	if currentStanza.Verses != nil {
		currentStanza.Type = stanzaType
		stanzaList = append(stanzaList, currentStanza)
	}

	return stanzaList
}

func handleTextNodes(node *goquery.Selection) songbook.Verse {
	var verse songbook.Verse
	node.Contents().Each(func(_ int, child *goquery.Selection) {
		if child.Is(":not(a)") {
			for _, word := range child.Text() {
				verse = append(verse, songbook.Word{Text: string(word)})
			}
		}
	})
	return verse
}

func handleLinkedTextNodes(node *goquery.Selection) []songbook.Verse {
	var verses []songbook.Verse
	node.Find("a").Each(func(_ int, link *goquery.Selection) {
		link.Contents().Each(func(_ int, child *goquery.Selection) {
			if child.Is(":not(a)") {
				verse := handleTextNodes(child)
				verses = append(verses, verse)
			}
		})
	})
	return verses
}

func getType(value string) string {
	switch {
	case strings.Contains(value, "Verse"):
		return "VERSE"
	case value == "Pre-Chorus":
		return "PRECHORUS"
	case value == "Chorus":
		return "CHORUS"
	case value == "Bridge":
		return "BRIDGE"
	default:
		return ""
	}
}
