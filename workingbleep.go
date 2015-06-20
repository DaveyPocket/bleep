package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"math"
	"github.com/nsf/termbox-go"
	"time"
	"math/rand"
	)


func main() {
	var t thing
	t.Buff.Data = make([]int8, 40000)
	t.Buff.Length = 40000
	
	
	// Start portAudio stream. Appoint buffer data first.
	
	termbox.Init()
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)
	var x int
	var fx float32
	var velx float32
	var x2, velx2 int
	t.Start()
	c := make(chan rune)
	var temp rune
	go ui(c)
	for {
		
		if x >= 18 {
			playBeep(&t.Buff, beep(110, 10000))
			velx = -1
		}else if x <= 0 {
		playBeep(&t.Buff, beep(220, 10000))

			velx = 1
		}
		if x2 >= 9 {
			playBeep(&t.Buff, noise(4000))
			velx2 = -1
		}else if x2 <= 0 {
			playBeep(&t.Buff, noise(4000))
			velx2 = 1
		}
		fx += velx
		x = int(fx)
		x2 += velx2
		termbox.SetCell(x, 2, 'A', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x2, 3, 'B', termbox.ColorDefault, termbox.ColorDefault)
		termbox.Flush()
		time.Sleep(50* time.Millisecond)
			select{
			case j := <-c:
				temp = j
			default:
		}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
termbox.SetCell(1, 1, temp, termbox.ColorDefault, termbox.ColorDefault)

	}


}


func ui(c chan rune) {
	for {
	switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				termbox.Close()
			}
			if ev.Ch != 0{
				c <- ev.Ch
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

func noise(length int) (wave []int8) {
	r := rand.New(rand.NewSource(12))
	wave = make([]int8, length)
	for n := 0; n < length; n++ {
		wave[n] = int8((1 - float64(n)/float64(length))*127*r.Float64())
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
