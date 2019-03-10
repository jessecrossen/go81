package main

import (
	"fmt"
	"math/rand"
)

// maximum cards dealt onto the table at one time
const tableSize = 21

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
	card := &g.deck[0]
	switch c {
	case 'A':
		card.turn++
	case 'B':
		card.turn = max(0, card.turn-1)
	case 'D':
		card.shrink = min(card.shrink+1, 5)
	case 'C':
		card.shrink = max(0, card.shrink-1)
	case ' ':
		card.id++
	case 's':
		card.selected = !card.selected
	default:
		return
	}
	g.needsRender = true
}

// Render the current game state to a frame buffer.
func (g *Game) Render() Frame {
	f := NewFrame()
	// draw cards in layers from back to front
	cardsTotal := len(g.deck)
	cardsFound := 0
	layer := 0
	for {
		for i := 0; i < cardsTotal; i++ {
			card := &g.deck[i]
			if card.layer == layer {
				cardsFound++
				if layer > 0 {
					card.Render(&f)
				}
			}
		}
		layer++
		// stop iterating if we've found all cards or the layer index gets too high
		if (cardsFound >= cardsTotal) || (layer >= 100) {
			break
		}
	}
	// draw the current score
	f.Draw(fmt.Sprintf("Score: %d", g.score), 1, (CardHeight*3)+1,
		ColorDefault, ColorDefault)
	return f
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

// TABLE OPERATIONS ***********************************************************

// get the card coordinates for the given index in the table
func tableCoords(i int) (col coord, row coord) {
	row = (i % 3) * CardHeight
	col = (i / 3) * CardWidth
	return
}

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
	card.layer = -1
	return &Animation{
		action: func(step int) bool {
			if step == 0 {
				card.col = 0
				card.row = 0
				card.shrink = MaxShrink
				card.turn = BackTurn
				card.layer = 2
			}
			p := float32(step) / dealSteps
			card.shrink = int(float32(MaxShrink) * (1.0 - p))
			card.col = int(float32(col) * p)
			card.row = int(float32(row) * p)
			if step >= dealSteps {
				card.layer = 1
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
		if card.layer == 0 {
			return card
		}
	}
	return nil
}

// deal a number of random cards onto the table
func (g *Game) dealRandom(count int) {
	if !(count > 0) {
		return
	}
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
