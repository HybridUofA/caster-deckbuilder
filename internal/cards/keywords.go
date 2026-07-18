package cards

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Keywords returns the distinct rules keywords discovered in card ability text.
//
// The normalized card data does not provide a separate keyword field, so the
// extractor recognizes the label forms used by the rules text: bracketed
// labels, labels followed by a definition in parentheses, and short labels
// before an ability separator. This lets newly published keywords appear in
// the search filter without requiring an application update.
func (repository *Repository) Keywords() []string {
	keywords := make(map[string]string)

	for _, card := range repository.cards {
		for _, line := range strings.Split(card.Ability, "\n") {
			keyword, found := abilityKeyword(line)
			if !found {
				continue
			}

			key := normalizeText(keyword)
			if _, exists := keywords[key]; !exists {
				keywords[key] = keyword
			}
		}
	}

	result := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
		result = append(result, keyword)
	}

	sort.Slice(result, func(i, j int) bool {
		return normalizeText(result[i]) < normalizeText(result[j])
	})

	return result
}

// abilityKeyword extracts a rules label from the beginning of one ability line.
func abilityKeyword(line string) (string, bool) {
	line = strings.TrimSpace(line)
	line = strings.TrimSpace(strings.TrimLeft(line, "•*-"))
	if line == "" {
		return "", false
	}

	if strings.HasPrefix(line, "[") {
		if closing := strings.IndexRune(line, ']'); closing > 1 {
			return validKeyword(line[1:closing])
		}
	}

	separator := strings.IndexAny(line, "(:,→")
	if separator >= 0 {
		return validKeyword(line[:separator])
	}

	if strings.ContainsAny(line, ".!?;") {
		return "", false
	}

	return validKeyword(line)
}

// validKeyword accepts concise labels while rejecting ordinary rules prose.
func validKeyword(candidate string) (string, bool) {
	candidate = strings.Join(strings.Fields(candidate), " ")
	if candidate == "" || len([]rune(candidate)) > 40 {
		return "", false
	}

	words := strings.Fields(candidate)
	if len(words) == 0 || len(words) > 5 {
		return "", false
	}

	for _, r := range candidate {
		if unicode.IsLetter(r) || unicode.IsDigit(r) ||
			r == ' ' || r == '-' || r == '\'' {
			continue
		}
		return "", false
	}

	// Multi-word labels are title-cased in the source data. Requiring that
	// convention avoids treating costs such as "Discard a card:" as keywords.
	if len(words) > 1 {
		for _, word := range words {
			first, _ := utf8.DecodeRuneInString(word)
			if !unicode.IsUpper(first) && !unicode.IsDigit(first) {
				return "", false
			}
		}
	}

	return candidate, true
}

// matchesAnyKeyword reports whether ability declares any selected keyword as
// a rules label. Mentions inside effect prose do not make the card a keyword
// card; for example, "cards they control lose Break" does not declare Break.
func matchesAnyKeyword(ability string, selected []string) bool {
	selectedKeywords := make(map[string]bool, len(selected))
	for _, keyword := range selected {
		if normalized := normalizeText(keyword); normalized != "" {
			selectedKeywords[normalized] = true
		}
	}

	for _, line := range strings.Split(ability, "\n") {
		keyword, found := abilityKeyword(line)
		if found && selectedKeywords[normalizeText(keyword)] {
			return true
		}
	}

	return false
}
