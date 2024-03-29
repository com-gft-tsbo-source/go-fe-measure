package femeasure

import (
	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// FeMeasure Response - Measure
// ###########################################################################
// ###########################################################################

// MeasureResponse Encapsulates the reploy of fe-measure
type MeasureResponse struct {
	microservice.Response
	RnrSvcVersion string `json:"rnrSvcVersion"`
}

// ###########################################################################

// InitMeasureResponse Constructor of a response of fe-measure
func InitMeasureResponse(r *MeasureResponse, status string, ms *FeMeasure) {
	microservice.InitResponseFromMicroService(&r.Response, ms, 200, status)
	r.RnrSvcVersion = "???"
}

// NewMeasureResponse ...
func NewMeasureResponse(status string, ms *FeMeasure) *MeasureResponse {
	var r MeasureResponse
	InitMeasureResponse(&r, status, ms)
	return &r
}
