package sduphue

import (
	"errors"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/amimof/huego"
	"github.com/mitchellh/mapstructure"
)

func init() {
	capRegistry[CapabilitySetColorXY] = TriggerSetXYColor
	gCapRegistry[CapabilitySetColorXY] = GTriggerSetXYColor
}

type XYColorArgs struct {
	X *float32 `mapstructure:"x"`
	Y *float32 `mapstructure:"y"`
}

func TriggerSetXYColor(id int, args ingestmodels.DeviceCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs XYColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.X == nil {
		return errors.New("x must be set")
	}

	if pArgs.Y == nil {
		return errors.New("y must be set")
	}

	_, err := bridge.SetLightState(id, huego.State{On: true, Xy: []float32{*pArgs.X, *pArgs.Y}})
	return err
}

func GTriggerSetXYColor(id int, args ingestmodels.GroupCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs XYColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.X == nil {
		return errors.New("x must be set")
	}

	if pArgs.Y == nil {
		return errors.New("y must be set")
	}

	_, err := bridge.SetGroupState(id, huego.State{On: true, Xy: []float32{*pArgs.X, *pArgs.Y}})
	return err
}
