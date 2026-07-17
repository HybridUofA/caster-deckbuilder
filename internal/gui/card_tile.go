package deckgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"

	"github.com/HybridUofA/caster-deckbuilder/internal/cardimages"
	"github.com/HybridUofA/caster-deckbuilder/internal/cards"
)

type CardTile struct {
	widget.BaseWidget

	Card cards.Card

	image *canvas.Image

	OnSelected  func(cards.Card)
	OnRightClick func(cards.Card, bool)
}

func NewCardTile(
	card cards.Card,
	onSelected func(cards.Card),
	onRightClick func(cards.Card, bool),
) *CardTile {
	tile := &CardTile{
		Card:         card,
		OnSelected:   onSelected,
		OnRightClick: onRightClick,
	}

	thumbPath, found := cardimages.FindThumbnail(card.ID)

	if found {
		tile.image = canvas.NewImageFromFile(thumbPath)
	} else {
		fullPath, found := cardimages.Find(card.ID)
		if found {
			tile.image = canvas.NewImageFromFile(fullPath)
		} else {
			tile.image = canvas.NewImageFromResource(theme.BrokenImageIcon())
		}
	}

	tile.image.FillMode = canvas.ImageFillContain
	tile.image.ScaleMode = canvas.ImageScaleSmooth

	tile.ExtendBaseWidget(tile)

	return tile
}

func (tile *CardTile) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(tile.image)
}

// Normal left click.
func (tile *CardTile) Tapped(_ *fyne.PointEvent) {
	if tile.OnSelected != nil {
		tile.OnSelected(tile.Card)
	}
}

// Required by desktop.Mouseable.
func (tile *CardTile) MouseDown(_ *desktop.MouseEvent) {
}

func (tile *CardTile) MouseUp(event *desktop.MouseEvent) {
	if event.Button != desktop.MouseButtonSecondary {
		return
	}

	shiftHeld :=
		event.Modifier&fyne.KeyModifierShift != 0

	if tile.OnRightClick != nil {
		tile.OnRightClick(tile.Card, shiftHeld)
	}
}