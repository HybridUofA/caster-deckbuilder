package cards

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Repository struct {
	cards  []Card
	byID   map[string]Card
	byName map[string][]Card
}

type Filter struct {
	Name           string
	Elements       []string
	Types          []string
	Traits         []string
	CostLevels     []string
	Expansions     []string
	IncludeTesting bool
}

func LoadFile(path string) (*Repository, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read card database %q: %w", path, err)
	}

	var cards []Card
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("decode card database %q: %w", path, err)
	}

	repository, err := NewRepository(cards)
	if err != nil {
		return nil, fmt.Errorf(
			"build repository from %q: %w",
			path,
			err,
		)
	}

	return repository, nil
}

func NewRepository(cards []Card) (*Repository, error) {
	if len(cards) == 0 {
		return nil, fmt.Errorf("card list cannot be empty")
	}

	repository := &Repository{
		cards:  make([]Card, 0, len(cards)),
		byID:   make(map[string]Card, len(cards)),
		byName: make(map[string][]Card),
	}

	for index, card := range cards {
		card.ID = strings.TrimSpace(card.ID)
		card.Name = strings.TrimSpace(card.Name)

		if card.ID == "" {
			return nil, fmt.Errorf(
				"card at index %d has no ID",
				index,
			)
		}

		if card.Name == "" {
			return nil, fmt.Errorf(
				"card %q has no name",
				card.ID,
			)
		}

		if _, exists := repository.byID[card.ID]; exists {
			return nil, fmt.Errorf(
				"duplicate card ID %q",
				card.ID,
			)
		}

		repository.cards = append(repository.cards, card)
		repository.byID[card.ID] = card

		nameKey := normalizeText(card.Name)

		repository.byName[nameKey] = append(
			repository.byName[nameKey],
			card,
		)
	}

	return repository, nil
}

func normalizeText(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

func (repository *Repository) FindByID(id string) (Card, bool) {
	card, found := repository.byID[strings.TrimSpace(id)]
	return card, found
}

func (repository *Repository) FindByName(name string) []Card {
	matches := repository.byName[normalizeText(name)]

	result := make([]Card, len(matches))
	copy(result, matches)

	return result
}

func (repository *Repository) SearchByName(query string) []Card {
	normalizedQuery := normalizeText(query)

	if normalizedQuery == "" {
		return []Card{}
	}

	var matches []Card

	for _, card := range repository.cards {
		normalizedCardName := normalizeText(card.Name)

		if strings.Contains(normalizedCardName, normalizedQuery) {
			matches = append(matches, card)
		}
	}

	return matches
}

func matchesAnySelected(
	cardValue string,
	selectedValues []string,
) bool {
	normalizedCard := normalizeText(cardValue)

	for _, selected := range selectedValues {
		normalizedSelected := normalizeText(selected)

		if normalizedSelected == "" {
			continue
		}

		if strings.Contains(
			normalizedCard,
			normalizedSelected,
		) {
			return true
		}
	}

	return false
}

func (repository *Repository) Filter(options Filter) []Card {
	normalizedName := normalizeText(options.Name)

	var matches []Card

	for _, card := range repository.cards {
		if normalizedName != "" &&
			!strings.Contains(normalizeText(card.Name), normalizedName) {
			continue
		}

		if len(options.Elements) > 0 &&
			!matchesAnySelected(card.Element, options.Elements) {
			continue
		}

		if len(options.Types) > 0 &&
			!matchesAnySelected(card.Type, options.Types) {
			continue
		}

		if len(options.Traits) > 0 &&
			!matchesAnySelected(card.Traits, options.Traits) {
			continue
		}

		if len(options.CostLevels) > 0 &&
			!matchesAnySelected(card.CostLevel, options.CostLevels) {
			continue
		}

		if len(options.Expansions) > 0 &&
			!matchesAnySelected(card.Expansion, options.Expansions) {
			continue
		}

		if card.IsPlaytesting && !options.IncludeTesting {
			continue
		}

		matches = append(matches, card)
	}

	return matches
}

func (repository *Repository) All() []Card {
	result := make([]Card, len(repository.cards))
	copy(result, repository.cards)

	return result
}

func uniqueSortedValues(values []string) []string {
	unique := make(map[string]string)

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		key := normalizeText(value)
		if _, exists := unique[key]; !exists {
			unique[key] = value
		}
	}

	results := make([]string, 0, len(unique))

	for _, value := range unique {
		results = append(results, value)
	}

	sort.Strings(results)

	return results
}

func (repository *Repository) Elements() []string {
	values := make([]string, 0, len(repository.cards))

	for _, card := range repository.cards {
		values = append(values, card.Element)
	}

	return uniqueSortedValues(values)
}

func (repository *Repository) Types() []string {
	values := make([]string, 0, len(repository.cards))

	for _, card := range repository.cards {
		values = append(values, card.Type)
	}

	return uniqueSortedValues(values)
}

func (repository *Repository) Traits() []string {
	unique := make(map[string]string)

	for _, card := range repository.cards {
		for _, trait := range splitTraits(card.Traits) {
			key := normalizeText(trait)

			if key == "" {
				continue
			}

			// Preserve the first spelling encountered.
			if _, exists := unique[key]; !exists {
				unique[key] = trait
			}
		}
	}

	traits := make(
		[]string,
		0,
		len(unique),
	)

	for _, trait := range unique {
		traits = append(traits, trait)
	}

	sort.Slice(
		traits,
		func(i, j int) bool {
			return strings.ToLower(traits[i]) <
				strings.ToLower(traits[j])
		},
	)

	return traits
}

func (repository *Repository) CostLevels() []string {
	values := make([]string, 0, len(repository.cards))

	for _, card := range repository.cards {
		values = append(values, card.CostLevel)
	}

	return uniqueSortedValues(values)
}

func (repository *Repository) Expansions() []string {
	values := make([]string, 0, len(repository.cards))

	for _, card := range repository.cards {
		values = append(values, card.Expansion)
	}

	return uniqueSortedValues(values)
}

func matchesAnyExact(values []string, target string) bool {
	normalizedTarget := normalizeText(target)

	for _, value := range values {
		if normalizeText(value) == normalizedTarget {
			return true
		}
	}

	return false
}

func matchesAnyContained(values []string, target string) bool {
	normalizedTarget := normalizeText(target)

	for _, value := range values {
		if strings.Contains(
			normalizedTarget,
			normalizeText(value),
		) {
			return true
		}
	}

	return false
}

func splitTraits(value string) []string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return nil
	}

	traits := make([]string, 0)
	remaining := raw

	for {
		open := strings.Index(remaining, "[")
		if open == -1 {
			break
		}

		closeOffset := strings.Index(
			remaining[open+1:],
			"]",
		)
		if closeOffset == -1 {
			break
		}

		close := open + 1 + closeOffset

		trait := strings.TrimSpace(
			remaining[open+1 : close],
		)
		if trait != "" {
			traits = append(traits, trait)
		}
		remaining = remaining[close+1:]
	}
	if len(traits) > 0 {
		return traits
	}

	for _, trait := range strings.FieldsFunc(
		raw,
		func(r rune) bool {
			return r == ',' || r == ';'
		},
	) {
		trait = strings.TrimSpace(trait)

		if trait != "" {
			traits = append(traits, trait)
		}
	}

	return traits
}
