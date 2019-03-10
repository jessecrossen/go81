package main

import (
	"bufio"
	"os"
	"os/exec"
	"time"
)

func main() {
	disableLineBuffering()
	disableEcho()
	defer enableEcho()
	// make a new game
	game := NewGame()
	// make channels that update the game
	input := newInput()
	timer := newTimer()
	display := NewDisplay()
	// start the interactive loop
	needsRender := true
	for {
		select {
		case c := <-input:
			needsRender = game.Input(c) || needsRender
			if c == 'q' {
				return
			}
		case _ = <-timer:
			if needsRender {
				display <- game.Render()
				needsRender = false
			}
		}
	}
}

func newInput() <-chan rune {
	input := make(chan rune)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			c, _, err := reader.ReadRune()
			if err == nil {
				input <- c
			}
		}
	}()
	return input
}

func newTimer() <-chan int {
	times := make(chan int)
	go func() {
		counter := 0
		for {
			times <- counter
			counter++
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return times
}

// from: https://stackoverflow.com/a/17278730/745831
func disableLineBuffering() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
}
func disableEcho() {
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}
func enableEcho() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}
