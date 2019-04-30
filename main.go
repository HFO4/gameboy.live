package main

import "github.com/HFO4/gbc-in-cloud/server"

func run() {
	streamServer := new(server.StreamServer)
	streamServer.Port = 2333
	gameList := []server.GameInfo{
		0: {
			Title: "Tetris",
			Path:  "test.gb",
		},
		1: {
			Title: "Dr. Mario",
			Path:  "Dr. Mario (JU) (V1.1).gb",
		},
		2: {
			Title: "Legend of Zelda - Link's Awakening",
			Path:  "Legend of Zelda, The - Link's Awakening (U) (V1.2) [!].gb",
		},
		3: {
			Title: "Pokemon - Blue",
			Path:  "Pokemon - Blue Version (UE) [S][!].gb",
		},
		4: {
			Title: "Super Mario Land",
			Path:  "Super Mario Land (JUE) (V1.1) [!].gb",
		},
		5: {
			Title: "Super Mario Land 2",
			Path:  "Super Mario Land 2 - 6 Golden Coins (UE) (V1.2) [!].gb",
		},
		6: {
			Title: "Kirby's Dream Land 2",
			Path:  "Kirby's Dream Land 2 (USA, Europe) (SGB Enhanced).gb",
		},
		7: {
			Title: "F-1 Race",
			Path:  "F-1 Race (JUE) (V1.1) [!].gb",
		},
	}
	streamServer.GameList = gameList
	streamServer.Run()
}

//func main() {
//	pixelgl.Run(run)
//
//}
func main() {
	run()

}

//func run() {
//	chan1 := make(chan bool)
//	Driver := new(driver.LCD)
//	core := new(gb.Core)
//	core.FPS = 60
//	core.Clock = 4194304
//	core.Debug = true
//	core.DebugControl = 255
//	core.DisplayDriver = Driver
//	core.Controller = Driver
//	core.DrawSignal = make(chan bool)
//	core.SpeedMultiple = 0
//	core.ToggleSound = true
//	core.Init("G:\\LearnGo\\gb\\Wario Land - Super Mario Land 3 (World).gb")
//	go core.DisplayDriver.Run(core.DrawSignal)
//	go core.Run()
//	//var t byte
//
//	<-chan1
//
//}
