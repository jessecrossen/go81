package main

import (
	"fmt"
)

func main() {
	f := EmptyFrame()
	f.Draw(""+
		"╭─╮\n"+
		"│e│\n"+
		"│h│\n"+
		"│k│\n"+
		"╰─╯", 1, 2, ColorRed, ColorDefault)
	f.Draw("1\n2\n3", 2, 3, ColorYellow, ColorDefault)
	fmt.Print(f.Render())
}
