//go:build !windows && !linux && !darwin

package audio

import (
	"github.com/diakovliev/oak/v4/audio/pcm"
	"github.com/diakovliev/oak/v4/oakerr"
)

func initOS(driver Driver) error {
	return oakerr.UnsupportedPlatform{
		Operation: "pcm.Init",
	}
}

func newWriter(f pcm.Format) (pcm.Writer, error) {
	return nil, oakerr.UnsupportedPlatform{
		Operation: "pcm.NewWriter",
	}
}
