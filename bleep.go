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
	var fx float32
	var x2, velx2 int
	t.Start()
	c := make(chan rune)
	go ui(c)
	go gui(c)
	objList := []Ball{}
	objList = addBall(1, 15, 1, 0, objList)
	
	w := Wall{"vertical", 19, 1, 0}

	for {
		if objList[0].PosX >= 18 {
			playBeep(&t.Buff, beepSquare(1000, 600, 0))
			objList[0].VX  = -1
		}else if objList[0].PosX <= 0 {
			playBeep(&t.Buff, beepSquare(1000, 600, 0))
			objList[0].VX  = 1
		}
		if x2 >= 18 {
			playBeep(&t.Buff, noise(4000))
			velx2 = -2
		}else if x2 <= 0 {
			playBeep(&t.Buff, noise(4000))
			velx2 = 2
		}
		fx += objList[0].VX
		objList[0].PosX = int(fx)
		x2 += velx2
		termbox.SetCell(objList[0].PosX, 2, 9673, termbox.ColorGreen, termbox.ColorDefault)
		termbox.SetCell(x2, 3, 9673, termbox.ColorRed, termbox.ColorDefault)
		w.Draw()
		termbox.Flush()
		time.Sleep(50* time.Millisecond)
		
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		
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

func addBall(x, y int, velx, vely float32, objList []Ball) []Ball{
	return append(objList, Ball{x, y, velx, vely})
}

func (t *thing) Start() {
	portaudio.Initialize()
	t.S, _ = portaudio.OpenDefaultStream(0, 2, 44100, 40, t.myCallback)
	t.S.Start()
}

type Ball struct {
	PosX		int
	PosY		int
	VX			float32
	VY			float32
}




func newBall(x, y int, vx, vy float32) (b Ball) {
	b.PosX = x
	b.PosY = y
	b.VX = vx
	b.VY = vy
	return b
}

func newWall(x, y, length int, orientation string) (w Wall) {
	w.Orientation = orientation
	w.PosX = x
	w.PosY = y
	w.Size = length
	return w
}

type objectSpace struct{
	balls		[]Ball
	walls		[]Wall
}

func (Space *objectSpace) draw() {
	for _, w := range Space.walls {
		w.Draw()
	}
	for _, b := range Space.balls {
		b.Draw()
	}
}

type object interface{
	Draw()
}

func (Space *objectSpace) add (o interface{}) {
	switch obj := o.(type) {
		case Ball:
			Space.balls = append(Space.balls, obj)
		case Wall:
			Space.walls = append(Space.walls, obj)
	}
}

type Wall struct {
 	Orientation string
 	PosX int
 	PosY int
	Size int
}

func (w * Wall) Draw(){
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
