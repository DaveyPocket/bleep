package main

import (
	"github.com/nsf/termbox-go"
	//"math"
	//"time"
)

type Ball struct {
	PosX int
	PosY int
	VX float32
	VY float32
}

type Wall struct {
 	Orientation string
 	PosX int
 	PosY int
}

func (w * Wall) Draw(){
	x,y := termbox.Size()
	if(w.Orientation == "vertical") {
		for i := 0 ; i < y;i++ {
			termbox.SetCell(w.PosX,i,'|',termbox.ColorWhite, termbox.ColorBlack)
		}
	}else if(w.Orientation == "horizontal") {
		for i := 0; i < x; i++ {
			termbox.SetCell(i,w.PosY,'_',termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func main(){
	err := termbox.Init()
	if err != nil {
        panic(err)
	}
	defer termbox.Close()

	for{
		termbox.SetCell(5,5,'A',termbox.ColorYellow,termbox.ColorBlack)
		wall1 := Wall{"vertical", 16,0}
		wall2 := Wall{"vertical", 6,0}
		wall1.Draw()
		wall2.Draw()
		termbox.Flush()
	}
}