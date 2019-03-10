package main

// A Card describes one card in a deck of cards.
type Card struct {
	id       int   // which card this is, coded from 0 to 80
	col      coord // the column to render the left edge of the card at
	row      coord // the row to render the top edge of the card at
	turn     int   // vary this to animate the card flipping over (0 to 8)
	shrink   int   // vary this to animate the card shrinking (0 to 5)
	selected bool  // whether the card has been selected by the user
	layer    int   // z-index of the card, where layers 0 or lower are never drawn
}

// Possible values for the layer property of a card.
const (
	LayerToDeal   = -1
	LayerNotDealt = 0
	LayerDealt    = 1
	LayerDealing  = 2
)

// MaxShrink is the maximum value of the shrink property of a card.
const MaxShrink = 5

// FrontTurn is the value of a card's turn property that shows the front.
const FrontTurn = 0

// BackTurn is the value of a card's turn property that shows the back.
const BackTurn = 4

// Attributes returns the categories the card is a member of.
func (c *Card) Attributes() (count, shape, fill, clr int) {
	count = (c.id % 3) + 1 // return count as 1-based for clarity
	shape = (c.id / 3) % 3 // all others range from 0 to 2
	fill = (c.id / 9) % 3
	clr = (c.id / 27) % 3
	return
}

// A Deck stores a complete deck of cards.
type Deck = [81]Card

// NewDeck creates a complete deck of cards.
func NewDeck() Deck {
	var d Deck
	for id := 0; id < len(d); id++ {
		d[id].id = id
	}
	return d
}

// RENDERING *****************************************************************

// CardWidth is the width of a rendered card in characters.
const CardWidth = 5

// CardHeight is the height of a rendered card in lines.
const CardHeight = 5

// Render the card into the given frame buffer.
func (c *Card) Render(f *Frame) {
	outlineColor := ColorLightGray
	if c.selected {
		outlineColor = ColorLightCyan
	}
	f.Draw(c.renderOutline(), c.col, c.row, outlineColor, ColorDefault)
	shrink, turn := c.normalizedShrinkAndTurn()
	if shrink == 0 {
		if turn <= 1 || turn >= 7 {
			f.Draw(c.renderFace(), c.col+2, c.row+1, c.faceColor(), ColorDefault)
		} else if turn >= 3 && turn <= 5 {
			f.Draw(c.renderBack(), c.col+2, c.row+1, outlineColor, ColorDefault)
		}
	}
}

// get the color for the card's face symbols
func (c *Card) faceColor() color {
	_, _, _, clr := c.Attributes()
	switch clr {
	case 0:
		return ColorRed
	case 1:
		return ColorGreen
	case 2:
		return ColorBlue
	}
	return ColorDefault
}

// limit the range of the animation parameters
func (c *Card) normalizedShrinkAndTurn() (shrink, turn int) {
	turn = max(0, c.turn%8)
	shrink = min(max(0, c.shrink), 5)
	return
}

func (c *Card) renderOutline() string {
	shrink, turn := c.normalizedShrinkAndTurn()
	// the turn animation sequence is mirrored to get a full flip
	//	and repeated for the front and back of the card
	//	 turn:    0 1 2 3 4 5 6 7 8
	//	 outline: 0 1 2 1 0 1 2 1 0
	turn = turn % 4
	if turn > 2 {
		turn = 4 - turn
	}
	switch shrink {
	case 0:
		switch turn {
		case 0:
			return "" +
				"╭───╮\n" +
				"│   │\n" +
				"│   │\n" +
				"│   │\n" +
				"╰───╯"
		case 1:
			return "" +
				" ╭─╮\n" +
				" │ │\n" +
				" │ │\n" +
				" │ │\n" +
				" ╰─╯"
		default:
			return "" +
				"  ╷\n" +
				"  │\n" +
				"  │\n" +
				"  │\n" +
				"  ╵"
		}
	case 1:
		switch turn {
		case 0:
			return "" +
				"╭──╮\n" +
				"│  │\n" +
				"│  │\n" +
				"╰──╯"
		case 1:
			return "" +
				" ╭─╮\n" +
				" │ │\n" +
				" │ │\n" +
				" ╰─╯"
		default:
			return "" +
				" ╷\n" +
				" │\n" +
				" │\n" +
				" ╵"
		}
	case 2:
		switch turn {
		case 0:
			return "" +
				"╭─╮\n" +
				"│ │\n" +
				"╰─╯"
		default:
			return "" +
				" ╷\n" +
				" │\n" +
				" ╵"
		}
	case 3:
		switch turn {
		case 0:
			return "" +
				"┌┐\n" +
				"└┘"
		default:
			return "" +
				"╷\n" +
				"╵"
		}
	case 4:
		switch turn {
		case 0:
			return "▯"
		default:
			return "│"
		}
	default:
		return "·"
	}
}

func (c *Card) renderFace() string {
	count, shape, fill, _ := c.Attributes()
	symbol := " "
	switch shape {
	case 0:
		switch fill {
		case 0:
			symbol = "△"
		case 1:
			symbol = "◮"
		case 2:
			symbol = "▲"
		}
	case 1:
		switch fill {
		case 0:
			symbol = "□"
		case 1:
			symbol = "◨"
		case 2:
			symbol = "■"
		}
	case 2:
		switch fill {
		case 0:
			symbol = "○"
		case 1:
			symbol = "◑"
		case 2:
			symbol = "●"
		}
	}
	switch count {
	case 1:
		return "\n" + symbol
	case 2:
		return symbol + "\n\n" + symbol
	case 3:
		return symbol + "\n" + symbol + "\n" + symbol
	}
	return ""
}

func (c *Card) renderBack() string {
	return "\n?"
}
