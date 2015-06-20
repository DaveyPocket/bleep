package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"math"
	"github.com/nsf/termbox-go"
	"time"
	)


func main() {
	var t thing
	t.Buff.Data = make([]int8, 4000)
	t.Buff.Length = 4000
	
	
	// Start portAudio stream. Appoint buffer data first.
	
	
	
	termbox.Init()
	termbox.SetInputMode(termbox.InputMouse)
	var x, velx int
	t.Start()

	for {
		termbox.SetCell(x, 2, 'A', termbox.ColorDefault, termbox.ColorDefault)
		termbox.Flush()
		if x > 20 {
			t.Buff.Data = beep(440, 4000)
			velx = -1
		}else if x < 3 {
			t.Buff.Data = beep(880, 4000)
			velx = 1
		}
		x += velx
		time.Sleep(100* time.Millisecond)
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	}


}


func beep(freq, length int) (wave []int8) {
	wave = make([]int8, length)
	for n := 0; n < length; n++ {
		wave[n] = int8(127*math.Sin(float64(freq*n)))
	}
	return wave
}

type ringBuff struct {
	Data 	[]int8
	Length	int
	Index	int
}

type thing struct {
	S *portaudio.Stream
	Buff ringBuff
}

func (rB *ringBuff) Next() {
	rB.Index++
	if rB.Index == rB.Length {
		rB.Index = 0
	}
}

func (t *thing) myCallback(_, out []int8) {
	for i := range out {
		tempIndex := t.Buff.Index
		out[i] = t.Buff.Data[tempIndex]
		t.Buff.Next()
		t.Buff.Data[tempIndex] = 0
	}
}

func (t *thing) Start() {
	portaudio.Initialize()
	t.S, _ = portaudio.OpenDefaultStream(0, 2, 44100, 40, t.myCallback)
	t.S.Start()
}
