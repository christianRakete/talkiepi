package talkiepi

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func (b *Talkiepi) closeGPIO() {
    rpio.Close()
}

func (b *Talkiepi) initGPIO() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		b.GPIOEnabled = false
		return
	} else {
		b.GPIOEnabled = true
	}

	if ButtonPin > 0 {
	    b.Button = rpio.Pin(ButtonPin)
	}
	
	// unfortunately the gpio watcher stuff doesnt work for me in this context, so we have to poll the button instead
	go func() {
		for {
			if ButtonPin > 0 {
			    b.Button.PullDown()
			    currentState := b.Button.Read()

			    if currentState != b.ButtonState {
				b.ButtonState = currentState

				if b.Stream != nil {
					if b.ButtonState == 0 {
						fmt.Printf("Button is released\n")
						b.TransmitStop()
					} else {
						fmt.Printf("Button is pressed\n")
						b.TransmitStart()
					}
				}
			    }
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// then we can do our gpio stuff
	if OnlineLEDPin > 0 {
	    b.OnlineLED = rpio.Pin(OnlineLEDPin)
	    b.OnlineLED.Output()
	}
	if ParticipantsLEDPin > 0 {
	    b.ParticipantsLED = rpio.Pin(ParticipantsLEDPin)
	    b.ParticipantsLED.Output()
	}
	if TransmitLEDPin > 0 {
    	    b.TransmitLED = rpio.Pin(TransmitLEDPin)
	    b.TransmitLED.Output()
	}
}

func (b *Talkiepi) LEDOn(LED rpio.Pin) {
    if b.GPIOEnabled == false && LED > 0 {
	return
    }

    LED.High()
}

func (b *Talkiepi) LEDOff(LED rpio.Pin) {
    if b.GPIOEnabled == false && LED > 0 {
	return
    }

    LED.Low()
}
func (b *Talkiepi) LEDOffAll() {
	if b.GPIOEnabled == false {
		return
	}

	if OnlineLEDPin > 0 {
	    b.OnlineLED.Low()
	}
	if ParticipantsLEDPin > 0 {
	    b.ParticipantsLED.Low()
	}
	if TransmitLEDPin > 0 {
	    b.TransmitLED.Low()
	}
}
