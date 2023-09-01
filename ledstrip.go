package talkiepi

import (
	"fmt"
	"strconv"

	"github.com/periph/conn/physic"
	"github.com/periph/conn/spi"
	"github.com/periph/conn/spi/spireg"
	"github.com/periph/devices/apa102"
	"github.com/periph/host"

	"periph.io/x/periph/conn/physic@v3.6.8"
	"periph.io/x/periph/conn/spi@v3.6.8"
	"periph.io/x/periph/conn/spi/spireg@v3.6.8"
	"periph.io/x/periph/devices/apa102@v3.6.8"
	"periph.io/x/periph/host@v3.6.8"

	// "periph.io/x/periph@v3.6.8"
)

const (
	numLEDs            int = 3
	OnlineLED          int = 0
	ParticipantsLED    int = 1
	TransmitLED        int = 2
	OnlineCol          string = "FF0000"
	ParticipantsCol    string = "00FF00"
	TransmitCol        string = "0000FF"
	OffCol             string = "000000"
)

type LedStrip struct {
	buf      []byte
	display  *apa102.Dev
	spiInterface spi.PortCloser
}

func NewLedStrip() (*LedStrip, error) {
	var spiID string = "SPI0.0" //SPI port to use
	var intensity uint8 = 16 //light intensity [1-255]
	var temperature uint16 = 5000 //light temperature in °Kelvin [3500-7500]
	var hz physic.Frequency //SPI port speed
	var globalPWM bool = false

	if _, err := host.Init(); err != nil {
		return nil, err
	}

	// Open the display device.
	s, err := spireg.Open(spiID)
	if err != nil {
		return nil, err
	}
	//Set port speed
	if hz != 0 {
		if err := s.LimitSpeed(hz); err != nil {
			return nil, err
		}
	}
	if p, ok := s.(spi.Pins); ok {
		fmt.Printf("Using pins CLK: %s  MOSI: %s  MISO: %s", p.CLK(), p.MOSI(), p.MISO())
	}
	o := apa102.DefaultOpts
	o.NumPixels = numLEDs
	o.Intensity = intensity
	o.Temperature = temperature
	o.DisableGlobalPWM = globalPWM
	display, err := apa102.New(s, &o)
	if err != nil {
		return nil, err
	}
	fmt.Printf("init display: %s\n", display)

	buf := make([]byte, numLEDs*3)

	return &LedStrip{
		buf: buf,
		display: display,
		spiInterface: s,
	}, nil
}

func (ls *LedStrip) ledCtrl(num int, color string) error {
	rgb, err := strconv.ParseUint(color, 16, 32)
	if err != nil {
		return err
	}
	r := byte(rgb >> 16)
	g := byte(rgb >> 8)
	b := byte(rgb)
	ls.buf[num*numLEDs+0] = r
	ls.buf[num*numLEDs+1] = g
	ls.buf[num*numLEDs+2] = b

	_, err = ls.display.Write(ls.buf)

	fmt.Printf("%v\n", ls.buf)
	
	return err
}

func (ls *LedStrip) closePort() {
	ls.spiInterface.Close()
}