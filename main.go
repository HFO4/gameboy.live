package main

import (
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	//t:=byte(0)
	//t = util.SetBit(t,1)
	//t = util.SetBit(t,4)
	//fmt.Printf("%X",t)
	//fmt.Printf("%t",util.TestBit(t,0))
	//fmt.Printf("%t",util.TestBit(t,1))
	//fmt.Printf("%t",util.TestBit(t,2))
	//fmt.Printf("%t",util.TestBit(t,3))
	//fmt.Printf("%t\n",util.TestBit(t,4))
	//t = util.ClearBit(t,0)
	//t = util.ClearBit(t,1)
	//fmt.Printf("%t",util.TestBit(t,0))
	//fmt.Printf("%t",util.TestBit(t,1))
	//fmt.Printf("%t",util.TestBit(t,2))
	//fmt.Printf("%t",util.TestBit(t,3))
	//fmt.Printf("%t\n",util.TestBit(t,4))
	//t = util.ClearBit(t,4)
	//fmt.Printf("%t",util.TestBit(t,0))
	//fmt.Printf("%t",util.TestBit(t,1))
	//fmt.Printf("%t",util.TestBit(t,2))
	//fmt.Printf("%t",util.TestBit(t,3))
	//fmt.Printf("%t\n",util.TestBit(t,4))
	Driver := &driver.LCD{}
	core := gb.Core{
		FPS:           60,
		Clock:         4194304,
		Debug:         true,
		DebugControl:  255,
		DisplayDriver: Driver,
		Controller:    Driver,
		DrawSignal:    make(chan bool),
		SpeedMultiple: 0,
		ToggleSound:   true,
	}
	core.Init("G:\\LearnGo\\gb\\test.gb")
	go core.DisplayDriver.Run(core.DrawSignal)
	core.Run()
}
