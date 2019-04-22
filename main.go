package main

import "github.com/HFO4/gbc-in-cloud/gb"

func main() {
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
	core := gb.Core{
		FPS:   60,
		Clock: 4194304,
		Debug: true,
	}
	core.Init("G:\\LearnGo\\gb\\test.gb")
	core.Run()
}
