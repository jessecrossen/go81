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
	// toggle cards
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
			// check for a set
			if card.selected {
				g.checkForSet()
			}
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

// number of frames it takes to deal a card
const dealSteps = MaxShrink

// number of frames it takes to collect a card
const collectSteps = MaxShrink

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

// animate a card being collected to the given location
func collectAnimation(card *Card, col int, row int) *Animation {
	card.layer = LayerCollecting
	startCol := card.col
	startRow := card.row
	return &Animation{
		action: func(step int) bool {
			p := float32(step) / collectSteps
			card.shrink = int(float32(MaxShrink) * p)
			card.col = startCol + int(float32(col-startCol)*p)
			card.row = startRow + int(float32(row-startRow)*p)
			if step >= collectSteps {
				card.layer = LayerNotDealt
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
	for _, card := range g.table {
		if (card != nil) && (card.turn != FrontTurn) {
			g.animator.Animate(*revealAnimation(card))
		}
	}
}

// check to see whether the user has selected a set
func (g *Game) checkForSet() {
	// check if three cards are selected
	selected := g.selectedCards()
	if len(selected) == 3 {
		if areSet(selected[0], selected[1], selected[2]) {
			// the cards are a set, add to the score
			col, row := scoreCoords()
			for _, card := range selected {
				g.removeCardFromTable(card)
				g.animator.Animate(*collectAnimation(card, col, row))
			}
			g.score++
			g.tidyTable()
		} else {
			// the cards are not a set, subtract from the score
			g.score--
			for _, card := range selected {
				card.selected = false
			}
		}
	}
}

// deal and consolidate cards
func (g *Game) tidyTable() {
	dealt := g.countCardsDealt()
	if dealt < 12 {
		g.dealRandom(12 - dealt)
	}
}

// count cards on the table
func (g *Game) countCardsDealt() int {
	count := 0
	for _, card := range g.table {
		if (card != nil) && (card.layer == LayerDealt) {
			count++
		}
	}
	return count
}

// remove a card from the table
func (g *Game) removeCardFromTable(removeCard *Card) {
	for i, card := range g.table {
		if card == removeCard {
			g.table[i] = nil
		}
	}
}

// get all selected cards
func (g *Game) selectedCards() []*Card {
	selected := make([]*Card, 0, len(g.table))
	for _, card := range g.table {
		if (card != nil) && (card.selected) {
			selected = append(selected, card)
		}
	}
	return selected
}

// return whether three cards form make a set
func areSet(a, b, c *Card) bool {
	a1, a2, a3, a4 := a.Attributes()
	b1, b2, b3, b4 := b.Attributes()
	c1, c2, c3, c4 := c.Attributes()
	return (true &&
		areSameOrDifferent(a1, b1, c1) &&
		areSameOrDifferent(a2, b2, c2) &&
		areSameOrDifferent(a3, b3, c3) &&
		areSameOrDifferent(a4, b4, c4))
}
func areSameOrDifferent(a, b, c int) bool {
	return areSame(a, b, c) || areDifferent(a, b, c)
}
func areSame(a, b, c int) bool {
	return (a == b) && (b == c)
}
func areDifferent(a, b, c int) bool {
	return (a != b) && (b != c) && (a != c)
}
