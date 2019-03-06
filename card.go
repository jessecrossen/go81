package main

// A Card describes one card in a deck of cards.
type Card struct {
	id     int   // which card this is, coded from 0 to 80
	col    coord // the column to render the left edge of the card at
	row    coord // the row to render the top edge of the card at
	turn   int   // vary this to animate the card flipping over
	shrink int   // vary this to animate the card shrinking
}

// Attributes returns the categories the card is a member of.
func (c *Card) Attributes() (count, shape, fill, color int) {
	count = (c.id % 3) + 1 // return count as 1-based for clarity
	shape = (c.id / 3) % 3 // all others range from 0 to 2
	fill = (c.id / 9) % 3
	color = (c.id / 27) % 3
	return
}

// RENDERING *****************************************************************

// Render the card into the given frame buffer.
func (c *Card) Render(f *Frame) {
	f.Draw(c.renderOutline(), c.col, c.row, ColorDefault, ColorDefault)
	// TODO
	// shrink, turn := c.normalizedShrinkAndTurn()
	// if shrink == 0 && turn <= 1 {
	// 	count, shape, fill, color := c.Attributes()
	// 	f.Draw(c.renderFace(),
	// }
}

// repeat and mirror the turn animation, limit shrink range
func (c *Card) normalizedShrinkAndTurn() (shrink, turn int) {
	turn = max(0, c.turn%10)
	if turn > 5 {
		turn = 10 - turn
	}
	shrink = min(max(0, c.shrink), 5)
	return
}

func (c *Card) renderOutline() string {
	shrink, turn := c.normalizedShrinkAndTurn()
	switch shrink {
	case 0:
		switch turn {
		case 0:
			return "" +
				"╭────╮\n" +
				"│    │\n" +
				"│    │\n" +
				"│    │\n" +
				"╰────╯"
		case 1:
			return "" +
				"╭───╮\n" +
				"│   │\n" +
				"│   │\n" +
				"│   │\n" +
				"╰───╯"
		case 2:
			return "" +
				"╭──╮\n" +
				"│  │\n" +
				"│  │\n" +
				"│  │\n" +
				"╰──╯"
		case 3:
			return "" +
				"╭─╮\n" +
				"│ │\n" +
				"│ │\n" +
				"│ │\n" +
				"╰─╯"
		case 4:
			return "" +
				"╭╮\n" +
				"││\n" +
				"││\n" +
				"││\n" +
				"╰╯"
		case 5:
			return "" +
				"╷\n" +
				"│\n" +
				"│\n" +
				"│\n" +
				"╵"
		}
	case 1:
		switch turn {
		case 0:
			return "" +
				"╭───╮\n" +
				"│   │\n" +
				"│   │\n" +
				"╰───╯"
		case 1:
			return "" +
				"╭──╮\n" +
				"│  │\n" +
				"│  │\n" +
				"╰──╯"
		case 2:
			return "" +
				"╭─╮\n" +
				"│ │\n" +
				"│ │\n" +
				"╰─╯"
		case 3:
			return "" +
				"╭╮\n" +
				"││\n" +
				"││\n" +
				"╰╯"
		default:
			return "" +
				"╷\n" +
				"│\n" +
				"│\n" +
				"╵"
		}
	case 2:
		switch turn {
		case 0:
			return "" +
				"╭──╮\n" +
				"│  │\n" +
				"╰──╯"
		case 1:
			return "" +
				"╭─╮\n" +
				"│ │\n" +
				"╰─╯"
		case 2:
			return "" +
				"╭╮\n" +
				"││\n" +
				"╰╯"
		default:
			return "" +
				"╷\n" +
				"│\n" +
				"╵"
		}
	case 3:
		switch turn {
		case 0:
			return "" +
				"╭╮\n" +
				"╰╯"
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
			return "⏐"
		}
	default:
		return "·"
	}
	return ""
}
