package main

import (
	"bufio"
	"fmt"
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
	input := getInput()
	frames := getFrames()
	// start the interactive loop
	lastFrame := NewFrame()
	needsRender := true
	for {
		select {
		case c := <-input:
			needsRender = game.Input(c) || needsRender
			if c == 'q' {
				return
			}
		case _ = <-frames:
			if needsRender {
				f := NewFrame()
				game.Render(&f)
				fmt.Print(f.Replace(lastFrame))
				lastFrame = f
				needsRender = false
			}
		}
	}
}

func getInput() <-chan rune {
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

func getFrames() <-chan int {
	frames := make(chan int)
	go func() {
		counter := 0
		for {
			frames <- counter
			counter++
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return frames
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
