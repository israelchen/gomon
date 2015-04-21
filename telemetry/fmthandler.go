package telemetry

import (
	"fmt"
	"github.com/israelchen/gomon/util"
)

type FmtHandler struct{}

func (handler *FmtHandler) Started(telemetry *Telemetry) {
	util.Require(telemetry != nil, "telemetry: telemetry cannot be nil.")

	fmt.Printf("Telemetry %s started.\n", telemetry.name)
}

func (handler *FmtHandler) Ended(telemetry *Telemetry) {
	util.Require(telemetry != nil, "telemetry: telemetry cannot be nil.")

	elapsed := telemetry.EndTime().Sub(*telemetry.StartTime())

	if telemetry.Error() == nil {
		fmt.Printf("Telemetry %s ended. Elapsed: %v.\n", telemetry.Name(), elapsed)
	} else {
		fmt.Printf("Telemetry %s ended. Elapsed: %v, Error: %s\n", telemetry.Name(), elapsed, telemetry.Error())
	}
}
