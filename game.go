package main

import (
	"fmt"
	"math/rand"
)

// Game stores the complete state of a game in progress.
type Game struct {
	deck        Deck             // all cards in the game
	table       [tableSize]*Card // cards currently dealt to the table
	animator    Animator         // animations that modify game state
	needsRender bool             // whether game state has changed since the last render
	random      rand.Source      // a source of randomness for the game
	score       int              // the current player's score
}

// NewGame returns a game with initial state.
func NewGame() *Game {
	g := Game{
		deck:        NewDeck(),
		animator:    NewAnimator(),
		needsRender: true,
		random:      rand.NewSource(0),
	}
	g.dealRandom(12)
	return &g
}

// Input updates the game state based on an input character and return whether anything changed.
func (g *Game) Input(c rune) {
	tableIndex := -1
	if (c >= 'a') && (c < 'a'+tableSize) {
		tableIndex = int(c - 'a')
	} else if (c >= 'A') && (c < 'A'+tableSize) {
		tableIndex = int(c - 'A')
	}
	if tableIndex >= 0 {
		card := g.table[tableIndex]
		if card != nil && card.layer == LayerDealt {
			card.selected = !card.selected
			g.needsRender = true
			return
		}
	}
}

// Render the current game state to a frame buffer.
func (g *Game) Render() Frame {
	f := NewFrame()
	g.renderLetters(&f)
	g.renderCards(&f)
	g.renderScore(&f)
	return f
}

// draw cards in layers from back to front
func (g *Game) renderCards(f *Frame) {
	cardsTotal := len(g.deck)
	cardsFound := 0
	layer := 0
	for {
		for i := 0; i < cardsTotal; i++ {
			card := &g.deck[i]
			if card.layer == layer {
				cardsFound++
				if layer > 0 {
					card.Render(f)
				}
			}
		}
		layer++
		// stop iterating if we've found all cards or the layer index gets too high
		if (cardsFound >= cardsTotal) || (layer >= 100) {
			break
		}
	}
}

// draw letters marking each card position
func (g *Game) renderLetters(f *Frame) {
	for i, card := range g.table {
		if card != nil && card.layer == LayerDealt {
			color := ColorDarkGray
			if card.selected {
				color = ColorCyan
			}
			col, row := letterCoords(i)
			f.Draw(fmt.Sprintf("%c", 'A'+i), col, row, color, ColorDefault)
		}
	}
}

// render the player's current score
func (g *Game) renderScore(f *Frame) {
	col, row := scoreCoords()
	f.Draw(fmt.Sprintf("Score: %d", g.score), col, row, ColorDefault, ColorDefault)
}

// Update the game state and render to the given display if needed.
func (g *Game) Update(display chan<- Frame) {
	// apply animations
	if g.animator.Step() {
		g.needsRender = true
	}
	// render if needed
	if g.needsRender {
		display <- g.Render()
		g.needsRender = false
	}
}

// LAYOUT *********************************************************************

// maximum cards dealt onto the table at one time
const tableSize = 21

// get the card coordinates for the given index in the table
func tableCoords(i int) (col coord, row coord) {
	row = (i % 3) * CardHeight
	col = 1 + ((i / 3) * (CardWidth + 2))
	return
}

// the coords of the letter marking the given index in the table
func letterCoords(i int) (col coord, row coord) {
	col, row = tableCoords(i)
	col += CardWidth
	row += CardHeight / 2
	return
}

// get the coordinates for the score display
func scoreCoords() (col coord, row coord) {
	col = 1
	row = (CardHeight * 3) + 1
	return
}

// TABLE OPERATIONS ***********************************************************

// number of frames it takes to deal or remove a card
const dealSteps = 5

// animate dealing a card from the top left corner
func (g *Game) dealAnimation(card *Card) *Animation {
	// find the first empty spot on the table for the card
	tableIndex := 0
	for i := range g.table {
		if g.table[i] == nil {
			tableIndex = i
			break
		}
	}
	g.table[tableIndex] = card
	col, row := tableCoords(tableIndex)
	// ensure the card is invisible but not dealt twice
	card.layer = LayerToDeal
	return &Animation{
		action: func(step int) bool {
			if step == 0 {
				card.selected = false
				card.col = 0
				card.row = 0
				card.shrink = MaxShrink
				card.turn = BackTurn
				card.layer = LayerDealing
			}
			p := float32(step) / dealSteps
			card.shrink = int(float32(MaxShrink) * (1.0 - p))
			card.col = int(float32(col) * p)
			card.row = int(float32(row) * p)
			if step >= dealSteps {
				card.layer = LayerDealt
				return false
			}
			return true
		},
	}
}

// pick a random card that is not on the table
func (g *Game) pickCard() *Card {
	// iterate a limited number of times, just in case all cards are dealt
	for tries := 0; tries <= 1000; tries++ {
		i := g.random.Int63() % int64(len(g.deck))
		card := &g.deck[i]
		if card.layer == LayerNotDealt {
			return card
		}
	}
	return nil
}

// deal a number of random cards onto the table
func (g *Game) dealRandom(count int) *Animation {
	var firstAnimation *Animation
	var lastAnimation *Animation
	for i := 0; i < count; i++ {
		if firstAnimation == nil {
			firstAnimation = g.dealAnimation(g.pickCard())
			lastAnimation = firstAnimation
		} else {
			animation := g.dealAnimation(g.pickCard())
			lastAnimation.andThen = animation
			lastAnimation = animation
		}
	}
	lastAnimation.andThen = &Animation{
		action: func(step int) bool {
			g.revealAll()
			return false
		},
	}
	g.animator.Animate(*firstAnimation)
	// return the final animation for chaining
	return lastAnimation
}

// animate flipping the given card over to reveal its front
func revealAnimation(card *Card) *Animation {
	return &Animation{
		action: func(step int) bool {
			if card.turn < FrontTurn {
				card.turn++
				return true
			} else if card.turn > FrontTurn {
				card.turn--
				return true
			} else {
				return false
			}
		},
	}
}

// reveal all cards on the table
func (g *Game) revealAll() {
	for i := range g.table {
		card := g.table[i]
		if (card != nil) && (card.turn != FrontTurn) {
			g.animator.Animate(*revealAnimation(card))
		}
	}
}
