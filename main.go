package main

import "github.com/HFO4/gbc-in-cloud/gb"

func main() {
	core := gb.Core{
		FPS:   60,
		Clock: 4194304,
		Debug: true,
	}
	core.Init("G:\\LearnGo\\gb\\test.gb")
	core.Run()
}
