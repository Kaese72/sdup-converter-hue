package sduphue

import (
	"errors"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/amimof/huego"
	"github.com/mitchellh/mapstructure"
)

func init() {
	capRegistry[CapabilityDim] = CTriggerDim
	gCapRegistry[CapabilityDim] = GTriggerDim
}

type CTDimArgs struct {
	Inc *float32 `mapstructure:"inc"`
}

func CTriggerDim(id int, args ingestmodels.DeviceCapabilityArgs) error {
	var pArgs CTDimArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.Inc == nil {
		return errors.New("inc must be set")
	}
	// FIXME Limitations on inc variable

	_, err := bridge.SetLightState(id, huego.State{On: true, BriInc: int(*pArgs.Inc)})
	return err
}

func GTriggerDim(id int, args ingestmodels.GroupCapabilityArgs) error {
	var pArgs CTDimArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.Inc == nil {
		return errors.New("inc must be set")
	}
	// FIXME Limitations on inc variable

	_, err := bridge.SetGroupState(id, huego.State{On: true, BriInc: int(*pArgs.Inc)})
	return err
}
