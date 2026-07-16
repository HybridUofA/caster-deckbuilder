package main

import (
	"fmt"
	"log"

	"github.com/HybridUofA/caster-deckbuilder/internal/cards"

)

func main() {

	repository, err := cards.LoadFile("data/cards.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded %d cards\n", len(repository.All()))

	results := repository.Filter(cards.Filter{
		Name:			"arth",
		Elements:		[]string{"Void"},
		Types:			[]string{"Caster"},
		IncludeTesting: false,
	})

	fmt.Printf("Found %d matching card(s):\n", len(results))

	for _, card := range results {
		fmt.Printf(
			"- %s | %s | %s | %s\n",
			card.Name,
			card.Type,
			card.Element,
			card.Expansion,
		)
	}

	fmt.Printf("Elements: %v\n", repository.Elements())
	fmt.Printf("Types: %v\n", repository.Types())
}