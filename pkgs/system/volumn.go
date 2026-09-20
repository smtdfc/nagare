package system

import (
	"github.com/itchyny/volume-go"
)

type VolumeControl struct{}

func (v *VolumeControl) Set(value int) error {
	return volume.SetVolume(value)
}

func (v *VolumeControl) Increment(value int) error {
	currentVol, err := volume.GetVolume()
	if err != nil {
		return err
	}
	return volume.SetVolume(currentVol + value)
}

func (v *VolumeControl) Decrement(value int) error {
	currentVol, err := volume.GetVolume()
	if err != nil {
		return err
	}

	return volume.SetVolume(currentVol - value)
}

func (v *VolumeControl) Mute() error {
	return volume.SetVolume(0)
}

func (v *VolumeControl) Get() (int, error) {
	return volume.GetVolume()
}

func NewVolumeControl() *VolumeControl {
	return &VolumeControl{}
}
