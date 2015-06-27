package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"math"
	"github.com/nsf/termbox-go"
//	"time"
	"math/rand"
//	"fmt"
//	"reflect"
	)

func main() {
	//var t thing
	//t.Buff.Data = make([]int8, 40000)
	//t.Buff.Length = 40000
	// Start portAudio stream. Appoint buffer data first.

	termbox.Init()
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)
	//t.Start()
	c := make(chan rune)
	go ui(c)
	go gui(c)

	w := makeWorld()		// Create a world
	l := Wall{"horizontal", 10, 1, 50}
	w.Add(&l, newBall(10, 15, 0, 0))
	w.Draw()
	mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				break mainloop
			}
	}
	}
}

func gui(c chan rune) {
	for {
	select{
		case j := <-c:
			termbox.SetCell(1, 1, j, termbox.ColorDefault, termbox.ColorDefault)
			termbox.Flush()
		default:
			termbox.SetCell(1, 1, 'd', termbox.ColorDefault, termbox.ColorDefault)
	}
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

func beepPhase(freq, length, deg int) (wave []int8) {
	wave = make([]int8, length)
	for n := 0; n < length; n++ {
		wave[n] = int8((1 - float64(n)/float64(length))*63*math.Sin(float64(math.Pi)*float64(deg)/float64(180) + float64(2*math.Pi)*float64(freq*n)/44100)) + int8((1 - float64(n)/float64(length))*63*math.Sin(float64(2*math.Pi)*float64(freq*n)/44100))

	}
	return wave
}

func beepSquare(freq, length, duty int) (wave []int8) {
	// Duty cycle in percent
	wave = make([]int8, length)
	for n := 0; n < length; n++ {
		// Make separate envelope function?????!?!?!?!?11?1One!!
		wave[n] = int8((1 - float64(n)/float64(length)) * 127 * float64((((8*n)/(freq)) % 2)) + float64(duty))
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
	Data		[]int8
	Length	int
	Index		int
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

const (
	Edit	=	iota
	Run
)

func makeWorld() (w *world) {
	s := new(objectSpace)
	w = new(world)
	*w = world{s, 0}
	return
}
type world struct {
						*objectSpace
	State				byte
}

type Ball struct {
	PosX		int
	PosY		int
	VX			float32
	VY			float32
}

func newBall(x, y int, vx, vy float32) (b *Ball) {
	b = new(Ball)
	b.PosX = x
	b.PosY = y
	b.VX = vx
	b.VY = vy
	return b
}

func newWall(x, y, length int, orientation string) (w *Wall) {
	w = new(Wall)
	w.Orientation = orientation
	w.PosX = x
	w.PosY = y
	w.Size = length
	return w
}

// Can this be replaced by a signal slice of object interfaces??
type objectSpace struct{
	list		[]object
}

func (Space *objectSpace) Draw() {
	for _, o := range Space.list {
		o.Draw()
	}
	termbox.Flush()
}

type object interface{
	Draw()
}

func (Space *objectSpace) Add (o ...object) {
	for _, thing := range o {
			Space.list = append(Space.list, thing)
			//fmt.Println("Added ", reflect.TypeOf(thing))
	}
}

type Wall struct {
	Orientation		string
	PosX				int
	PosY				int
	Size				int
}

func (w *Wall) Draw(){
	if(w.Orientation == "vertical") {
		for i := 0 ; i < w.Size; i++ {
			termbox.SetCell(w.PosX, i,' ', termbox.ColorWhite, termbox.ColorRed)
		}
	}else if(w.Orientation == "horizontal") {
		for i := 0; i < w.Size; i++ {
			termbox.SetCell(i, w.PosY, ' ', termbox.ColorWhite, termbox.ColorRed)
		}
	}
}

func (b *Ball) Draw() {
	termbox.SetCell(b.PosX, b.PosY, 9673, termbox.ColorGreen, termbox.ColorDefault)
}
