package gui

import "gopkg.in/qml.v1"

import "log"

//import "github.com/ArchieT/3manchess/movedet"
import "github.com/ArchieT/3manchess/game"
import "math/cmplx"
import "math"

const (
	Center             = 350 + 350i
	InnerRadius        = 70
	SubsequentRadiiAdd = 35
	DefaultRotation    = -math.Pi / 6
	OneFile            = math.Pi / 12
)

//adowbiowl — Angle Depening On Whether Black Is On Whites Left
func adowbiowl(p float64, biowl bool) float64 {
	if !biowl {
		return math.Remainder(p, 2*math.Pi)
	}
	return math.Remainder(p+math.Pi, 2*math.Pi)
}

type GUI struct {
	appears             chan<- appearing
	BlackIsOnWhitesLeft bool
	fromtos             <-chan game.FromTo
	Rotated             float64 //zerofile blackmost boundary angle
	errchan             chan error
	ErrorChan           <-chan error
	engine              *qml.Engine
	component           *qml.Object
	window              *qml.Window
}

type boardmap [6][24]string

func (bm *boardmap) Appear(w game.BoardDiff) {
	bm[w.Pos[0]][w.Pos[1]] = FigURIs[w.Fig.Uint8()]
}

type boardclicker chan complex128

func (bckr boardclicker) ClickedIt(x, y int) {
	bckr <- complex(float64(x), float64(y))
}

func clicking(s <-chan complex128, d chan<- game.Pos, rot *float64, biowl *bool) {
	var c complex128
	var r, p float64
	var m uint16
	var pr, pf int8
	for {
		c = <-s
		log.Println("RawClick:", c)
		c -= Center
		r, p = cmplx.Polar(c)
		p -= *rot
		r -= InnerRadius
		if r < 0 {
			continue
		}
		p = adowbiowl(p, *biowl)
		m = uint16(r) / 35
		if m < 24 {
			pr = int8(m)
		} else {
			continue
		}
		pf = int8(p / OneFile)
		d <- game.Pos{pr, pf}
	}
}

func fromtoing(s <-chan game.Pos, d chan<- game.FromTo) {
	var f game.Pos
	for {
		f = <-s
		d <- game.FromTo{f, <-s}
	}
}

func NewGUI() (*GUI, error) {
	gui := new(GUI)
	clicks := make(boardclicker)
	clickpos := make(chan game.Pos)
	appears := make(chan appearing)
	fromtos := make(chan game.FromTo)
	gui.appears = appears
	gui.Rotated = DefaultRotation
	gui.fromtos = fromtos
	go clicking(clicks, clickpos, &(gui.Rotated), &(gui.BlackIsOnWhitesLeft))
	go fromtoing(clickpos, fromtos)
	gui.engine = qml.NewEngine()
	gui.engine.Context().SetVar("clickinto", clicks)
	component, err := gui.engine.LoadFile("okno.qml")
	gui.component = &component
	if err != nil {
		return gui, err
	}
	gui.errchan = make(chan error)
	gui.ErrorChan = gui.errchan
	gui.window = component.CreateWindow(nil)
	gui.window.Show()
	go gui.run()
	return gui, nil
}

func (gui *GUI) run() {
	gui.window.Wait()
}
