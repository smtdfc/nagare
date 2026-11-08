package system

import (
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
)

func GetActiveDisplay() (int, error) {
	n := screenshot.NumActiveDisplays()

	return n, nil
}

func Screenshot(index int, fileName string) error {
	bounds := screenshot.GetDisplayBounds(index)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return err
	}

	file, _ := os.Create(fileName)
	defer file.Close()
	return png.Encode(file, img)
}
