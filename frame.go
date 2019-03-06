package main

import (
	"fmt"
	"strings"
)

// a row or column index.
type coord = int

// a color constant
type color = int8

// the maximum number of rows in a frame
const maxRows coord = 40

// A Frame represents a block of text to render to the terminal.
type Frame struct {
	Lines  [][]rune
	Colors [][]string
}

// EmptyFrame returns a Frame struct with no content.
func EmptyFrame() Frame {
	return Frame{make([][]rune, 0, maxRows), make([][]string, 0, maxRows)}
}

// Draw a set of newline-delimited lines to the given coordinates in the frame.
func (f *Frame) Draw(text string, col coord, row coord, fg color, bg color) {
	drawLines := strings.Split(text, "\n")
	f.ensureRowCount(row + len(drawLines))
	for i, drawLine := range drawLines {
		drawRunes := []rune(drawLine)
		lineIndex := row + i
		if lineIndex < cap(f.Lines) {
			f.Lines[lineIndex] = insertStringInLine(f.Lines[lineIndex], drawRunes, col)
		}
		if lineIndex < cap(f.Colors) {
			f.Colors[lineIndex] = insertColorInLine(f.Colors[lineIndex],
				changeColor(fg, bg), col, len(drawRunes))
		}
	}
}

// Render a frame to a string that can be written to the terminal.
func (f *Frame) Render() string {
	b := strings.Builder{}
	for i := 0; i < len(f.Lines); i++ {
		b.WriteString(clearLine)
		f.renderLine(&b, i)
		b.WriteString(nextLine)
	}
	return b.String()
}

// Reset returns a string that will restore the cursor after rendering a frame.
func (f *Frame) Reset() string {
	return cursorUp(len(f.Lines))
}

// Replace returns a string that will replace the given frame with the receiver.
func (f *Frame) Replace(old Frame) string {
	b := strings.Builder{}
	b.WriteString(old.Reset())
	b.WriteString(f.Render())
	for i := len(f.Lines); i < len(old.Lines); i++ {
		b.WriteString(clearLine)
		b.WriteString(nextLine)
	}
	return b.String()
}

// console colors (offset from 30 for foreground and 40 for background)
const (
	ColorBlack        color = 0
	ColorRed          color = 1
	ColorGreen        color = 2
	ColorYellow       color = 3
	ColorBlue         color = 4
	ColorMagenta      color = 5
	ColorCyan         color = 6
	ColorLightGray    color = 7
	ColorDefault      color = 9
	ColorDarkGray     color = 61
	ColorLightRed     color = 62
	ColorLightGreen   color = 63
	ColorLightYellow  color = 64
	ColorLightBlue    color = 65
	ColorLightMagenta color = 66
	ColorLightCyan    color = 67
	ColorWhite        color = 68
)

// IMPLEMENTATION *************************************************************

// make sure a frame has at least the given number of rows
func (f *Frame) ensureRowCount(rows coord) {
	for i := len(f.Lines); i < min(rows, cap(f.Lines)); i++ {
		f.Lines = append(f.Lines, make([]rune, 0))
	}
	for i := len(f.Colors); i < min(rows, cap(f.Colors)); i++ {
		f.Colors = append(f.Colors, make([]string, 0))
	}
}

// insert a string into another string, overwriting characters it overlaps
func insertStringInLine(line []rune, s []rune, col coord) []rune {
	if len(s) == 0 {
		return line
	}
	afterInsertedIndex := col + len(s)
	if afterInsertedIndex > len(line) {
		newLine := make([]rune, afterInsertedIndex)
		copy(newLine, line)
		for i := len(line); i < len(newLine); i++ {
			newLine[i] = ' '
		}
		line = newLine
	}
	copy(line[col:col+len(s)], s)
	return line
}

// insert color changes into a line
func insertColorInLine(line []string, color string, col coord, width coord) []string {
	if width == 0 {
		return line
	}
	afterInsertedIndex := col + width
	if afterInsertedIndex > len(line) {
		newLine := make([]string, afterInsertedIndex)
		copy(newLine, line)
		fill := changeColor(ColorDefault, ColorDefault)
		for i := len(line); i < len(newLine); i++ {
			newLine[i] = fill
		}
		line = newLine
	}
	for i := col; i < col+width; i++ {
		line[i] = color
	}
	return line
}

// render a line, inserting color changes where needed
func (f *Frame) renderLine(b *strings.Builder, row coord) {
	if row >= len(f.Lines) {
		return
	}
	line := f.Lines[row]
	if row >= len(f.Colors) {
		b.WriteString(string(line))
		return
	}
	colors := f.Colors[row]
	lastColor := ""
	lastIndex := 0
	minLen := min(len(colors), len(line))
	for i := 0; i < minLen; i++ {
		thisColor := colors[i]
		if thisColor != lastColor {
			b.WriteString(string(line[lastIndex:i]))
			b.WriteString(thisColor)
			lastColor = thisColor
			lastIndex = i
		}
	}
	if lastIndex < len(line) {
		b.WriteString(string(line[lastIndex:]))
	}
	b.WriteString(changeColor(ColorDefault, ColorDefault))
}

// compose terminal escape sequences
const escape = "\x1b"
const clearLine = escape + "[2K"
const nextLine = "\n"

func cursorUp(lines int) string {
	return fmt.Sprintf("%s[%dA", escape, lines)
}

func changeColor(fg, bg color) string {
	// return fmt.Sprintf("%s[%dm%s[%dm", escape, 39+fg, escape, 49+bg)
	return fmt.Sprintf("%s[%d;%dm", escape, 30+fg, 40+bg)
}

// utility functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
