package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"github.com/HybridUofA/caster-deckbuilder/internal/speedrobo"
	"github.com/HybridUofA/caster-deckbuilder/internal/updates"
)

func main() {

	client, err := speedrobo.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	config, err := speedrobo.FetchPageConfig(client)
	if err != nil {
		log.Fatal(err)
	}

	cards, err := updates.FetchAllCards(client, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downloaded %d cards\n", len(cards))

	cardJSON, err := json.MarshalIndent(cards, "", " ")
	if err != nil {
		log.Fatalf("encode cards as JSON: %v", err)
	}
	cardJSON = append(cardJSON, '\n')

	if err := os.WriteFile("data/cards.json", cardJSON, 0644); err != nil {
		log.Fatalf("Write card database: %v", err)
	}

	fmt.Printf("Saved %d cards to data/cards.json\n", len(cards))
}
