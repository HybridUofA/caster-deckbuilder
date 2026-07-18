package decks

import (
	"fmt"
	"strings"
)

// Validate checks the deck's name and supported schema version.
func (deck *Deck) Validate() error {
	if strings.TrimSpace(deck.Name) == "" {
		return fmt.Errorf("deck name cannot be empty")
	}
	if deck.SchemaVersion != 1 {
		return fmt.Errorf("unsupported schema version: %d", deck.SchemaVersion)
	}
	return nil
}

// ValidateCards verifies that every main-deck entry resolves in the catalog.
func (deck *Deck) ValidateCards(repository CardCatalog) error {
	for _, entry := range deck.MainDeck {
		if _, found := repository.FindByID(entry.CardID); !found {
			return fmt.Errorf("main deck contains unknown card ID %q", entry.CardID)
		}
	}
	return nil
}
