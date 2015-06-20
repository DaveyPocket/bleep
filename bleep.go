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
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)
	var x int
	var velx float32
	var x2, velx2 int
	t.Start()
	go ui()
	for {
		termbox.SetCell(x, 2, 'A', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x2, 3, 'B', termbox.ColorDefault, termbox.ColorDefault)

		termbox.Flush()
		if x > 20 {
			playBeep(&t.Buff, beep(110, 4000))
			velx = -.5
		}else if x < 10 {
		playBeep(&t.Buff, beep(880, 4000))

			velx = .5
		}
		if x2 > 20 {
			playBeep(&t.Buff, beep(880, 4000))
			velx2 = -1
		}else if x2 < 10 {
			playBeep(&t.Buff, beep(880, 4000))
			velx2 = 1
		}
		x += velx
		x2 += velx2
		time.Sleep(50* time.Millisecond)
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	}


}


func ui() {
	for {
	switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				termbox.Close()
			}
	}
	}
}

func playBeep(in *ringBuff, wave []int8){
	tempIndex := in.Index
	for i, _ := range wave{
		in.Data[tempIndex] += wave[i]/4
		tempIndex++
		if tempIndex == in.Length{
			tempIndex = 0
		}
	} 
}

func beep(freq, length int) (wave []int8) {
	wave = make([]int8, length)
	for n := 0; n < length; n++ {
		wave[n] = int8((1 - float64(n)/float64(length))*   127*math.Sin(float64(2*math.Pi)*float64(freq*n)/44100))
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
