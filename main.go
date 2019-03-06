package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func cardDemo() {
	// c = Card{}
	// for {

	// }
}

func main() {
	disableLineBuffering()
	disableEcho()
	defer enableEcho()

	input := getInput()
	frames := getFrames()

	card := Card{}

	lastFrame := NewFrame()
	for {
		select {
		case c := <-input:
			switch c {
			case 'A':
				card.turn++
			case 'B':
				card.turn--
			case 'D':
				card.shrink = min(card.shrink+1, 5)
			case 'C':
				card.shrink = max(0, card.shrink-1)
			}
		case _ = <-frames:
			f := NewFrame()
			card.Render(&f)
			fmt.Print(f.Replace(lastFrame))
			lastFrame = f
		}
	}
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
