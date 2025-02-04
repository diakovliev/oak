//go:build darwin

package audio

import (
	"fmt"

	"github.com/diakovliev/oak/v4/audio/pcm"
	"github.com/diakovliev/oak/v4/oakerr"
	"github.com/jfreymuth/pulse"
)

func initOS(driver Driver) error {
	switch driver {
	case DriverDefault:
		fallthrough
	case DriverPulse:
		// Sanity check that pulse is installed and a sink is defined
		client, err := pulse.NewClient()
		if err != nil {
			// osx: brew install pulseaudio
			// linux: sudo apt install pulseaudio
			return oakerr.UnsupportedPlatform{
				Operation: "pcm.Init:" + driver.String(),
			}
		}
		defer client.Close()
		_, err = client.DefaultSink()
		if err != nil {
			return err
		}
		newWriter = newPulseWriter
	default:
		return oakerr.UnsupportedPlatform{
			Operation: "pcm.Init:" + driver.String(),
		}
	}
	return nil
}

var newWriter = func(f pcm.Format) (pcm.Writer, error) {
	return nil, fmt.Errorf("this package has not been initialized")
}
