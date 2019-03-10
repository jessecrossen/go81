package main

// Game stores the complete state of a game in progress.
type Game struct {
	deck Deck
}

// NewGame returns a game with initial state.
func NewGame() Game {
	return Game{
		deck: NewDeck(),
	}
}

// Input updates the game state based on an input character and return whether anything changed.
func (g *Game) Input(c rune) bool {
	card := &g.deck.cards[0]
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
		return false
	}
	return true
}

// Render the current game state to a frame buffer.
func (g *Game) Render() Frame {
	f := NewFrame()
	g.deck.cards[0].Render(&f)
	return f
}
