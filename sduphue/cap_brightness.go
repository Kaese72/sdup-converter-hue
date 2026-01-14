package sduphue

import (
	"errors"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/amimof/huego"
	"github.com/mitchellh/mapstructure"
)

func init() {
	capRegistry[CapabilitySetBrightness] = CTriggerSetBrightness
	gCapRegistry[CapabilitySetBrightness] = GTriggerSetBrightness
}

type CTBrightnessArgs struct {
	Value *float32 `mapstructure:"value"`
}

func CTriggerSetBrightness(id int, args ingestmodels.DeviceCapabilityArgs) error {
	var pArgs CTBrightnessArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.Value == nil {
		return errors.New("value must be set")
	}
	// FIXME Limitations on value variable

	_, err := bridge.SetLightState(id, huego.State{On: true, Bri: uint8(*pArgs.Value)})
	return err
}

func GTriggerSetBrightness(id int, args ingestmodels.GroupCapabilityArgs) error {
	var pArgs CTBrightnessArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.Value == nil {
		return errors.New("value must be set")
	}
	// FIXME Limitations on value variable

	_, err := bridge.SetGroupState(id, huego.State{On: true, Bri: uint8(*pArgs.Value)})
	return err
}
