package sduphue

import (
	"errors"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/amimof/huego"
	"github.com/mitchellh/mapstructure"
)

func init() {
	capRegistry[CapabilitySetColorTemp] = TriggerSetCTColor
	gCapRegistry[CapabilitySetColorTemp] = GTriggerSetCTColor
}

type CTColorArgs struct {
	Ct *float32 `mapstructure:"ct"`
}

func TriggerSetCTColor(id int, args ingestmodels.DeviceCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs CTColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.Ct == nil {
		return errors.New("ct must be set")
	}

	_, err := bridge.SetLightState(id, huego.State{On: true, Ct: uint16(*pArgs.Ct)})
	return err
}

func GTriggerSetCTColor(id int, args ingestmodels.GroupCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs CTColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return err
	}
	if pArgs.Ct == nil {
		return errors.New("ct must be set")
	}

	_, err := bridge.SetGroupState(id, huego.State{On: true, Ct: uint16(*pArgs.Ct)})
	return err
}
