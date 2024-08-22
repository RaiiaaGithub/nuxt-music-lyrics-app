package songbook

import (
	"strings"
	"unicode"
)

func addWordsToVerse(text string, stanza *Stanza) {
	words := splitWords(text)

	if len(stanza.Verses) == 0 {
		stanza.Verses = append(stanza.Verses, Verse{})
	}
	currentVerse := &stanza.Verses[len(stanza.Verses)-1]

	for _, word := range words {
		*currentVerse = append(*currentVerse, Word{Text: word, Chord: []string{}})
	}
}

func ensureVerseEnd(stanza *Stanza) {
	if len(stanza.Verses) > 0 && len(stanza.Verses[len(stanza.Verses)-1]) > 0 {
		stanza.Verses = append(stanza.Verses, Verse{})
	}
}

func splitWords(text string) []string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	return words
}

func cleanStanzaType(text string) string {
	cleanText := strings.Trim(text, "[]")
	return strings.TrimSpace(cleanText)
}
