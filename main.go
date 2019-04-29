package main

import (
	"fmt"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"github.com/faiface/pixel/pixelgl"
)

//func run() {
//	server.Run("2333")
//}

func main() {
	pixelgl.Run(run)

}

func run() {
	chan1 := make(chan bool)
	Driver := new(driver.LCD)
	core := new(gb.Core)
	core.FPS = 60
	core.Clock = 4194304
	core.Debug = true
	core.DebugControl = 255
	core.DisplayDriver = Driver
	core.Controller = Driver
	core.DrawSignal = make(chan bool)
	core.SpeedMultiple = 0
	core.ToggleSound = true
	core.Memory = new(gb.Memory)
	core.Init("G:\\LearnGo\\gb\\Kirby's Dream Land 2 (USA, Europe) (SGB Enhanced).gb")
	go core.DisplayDriver.Run(core.DrawSignal)
	go core.Run()
	var t byte
	fmt.Scanf("%d", &t)

	Driver2 := new(driver.LCD)
	core2 := new(gb.Core)
	core2.FPS = 60
	core2.Clock = 4194304
	core2.Debug = false
	core2.DebugControl = 0
	core2.DisplayDriver = Driver2
	core2.Controller = Driver2
	core2.DrawSignal = make(chan bool)
	core2.SpeedMultiple = 0
	core2.ToggleSound = true
	core2.Memory = new(gb.Memory)
	core2.Init("G:\\LearnGo\\gb\\Kirby's Dream Land 2 (USA, Europe) (SGB Enhanced).gb")
	go core2.DisplayDriver.Run(core2.DrawSignal)
	go core2.Run()

	<-chan1

}
