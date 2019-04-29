package main

import (
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"github.com/faiface/pixel/pixelgl"
)

//func run() {
//	listener, err := net.Listen("tcp", ":2333")
//	if err != nil {
//		fmt.Println("Error listening", err.Error())
//		return //终止程序
//	}
//
//	// 监听并接受来自客户端的连接
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			fmt.Println("Error accepting", err.Error())
//			return // 终止程序
//		}
//		doServerStuff(conn)
//	}
//	//pixelgl.Run(run)
//}
//
//func doServerStuff(conn net.Conn) {
//	conn.Write([]byte{255, 253, 34})
//	conn.Write([]byte{255, 250, 34, 1, 0, 255, 240})
//	conn.Write([]byte{0xFF, 0xFB, 0x01})
//	Driver := &driver.ASCII{
//		Conn: conn,
//	}
//	controller := &driver.TelnetController{}
//	core := gb.Core{
//		FPS:           60,
//		Clock:         4194304,
//		Debug:         true,
//		DebugControl:  255,
//		DisplayDriver: Driver,
//		Controller:    controller,
//		DrawSignal:    make(chan bool),
//		SpeedMultiple: 0,
//		ToggleSound:   true,
//	}
//	go func() {
//		go core.DisplayDriver.Run(core.DrawSignal)
//		core.Init("G:\\LearnGo\\gb\\Donkey Kong Land III (U) [S][!].gb")
//		core.Run()
//	}()
//	for {
//		buf := make([]byte, 512)
//		n, err := conn.Read(buf)
//		if err != nil {
//			fmt.Println("Error reading", err.Error())
//			os.Exit(1)
//			return //终止程序
//		}
//		log.Println(buf[:n])
//		core.Controller.NewInput(buf[:n])
//
//	}
//}

func main() {
	pixelgl.Run(run)
}

func run() {
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
	core.Init("G:\\LearnGo\\gb\\Pokemon - Blue Version (UE) [S][!].gb")
	go core.DisplayDriver.Run(core.DrawSignal)
	core.Run()
}
