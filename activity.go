package awsshadow

import (
	"encoding/json"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-rest")

const (
	ivThingName = "thingName"
	ivOp        = "op"
	ivDesired     = "desired"
	ivReported    = "reported"

	//ivAwsEndpoint = "awsEndpoint"

	ovResult = "result"
)

// AwsIoT is an Activity that is used to update an Aws IoT device shadow
// inputs : {method,uri,params}
// outputs: {result}
type AwsIoT struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AwsIoT activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AwsIoT{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *AwsIoT) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a Aws Iot Shadow Update
func (a *AwsIoT) Eval(context activity.Context) (done bool, err error) {

	thingName := context.GetInput(ivThingName).(string)
	op := context.GetInput(ivOp).(string)
	op = strings.ToLower(op)

	sess, err := session.NewSession()
	if err != nil {
		return false, err
	}

	idp := iotdataplane.New(sess)

	var payload []byte

	switch op {
	case "update":

		req := &ShadowRequest{State: &ShadowState{}}

		if context.GetInput(ivDesired) != nil {
			desired := context.GetInput(ivDesired).(map[string]string)
			req.State.Desired = desired
		}

		if context.GetInput(ivReported) != nil {
			reported := context.GetInput(ivReported).(map[string]string)
			req.State.Reported = reported
		}

		reqJSON, err := json.Marshal(req)

		sInput := &iotdataplane.UpdateThingShadowInput{}
		sInput.SetThingName(thingName)
		sInput.SetPayload(reqJSON)
		out, err := idp.UpdateThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	case "get":
		sInput := &iotdataplane.GetThingShadowInput{}
		sInput.SetThingName(thingName)
		out, err := idp.GetThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	case "delete":

		sInput := &iotdataplane.DeleteThingShadowInput{}
		sInput.SetThingName(thingName)
		out, err := idp.DeleteThingShadow(sInput)
		if err != nil {
			return false, err
		}
		payload = out.Payload
	}

	var result interface{}
	err = json.Unmarshal(payload, &result)
	if err != nil {
		return false, err
	}

	context.SetOutput(ovResult, result)
	return true, nil
}

// ShadowRequest is a simple structure representing a Aws Shadow Update Request
type ShadowRequest struct {
	State *ShadowState `json:"state"`
}

// ShadowState is the state to be updated
type ShadowState struct {
	Desired  map[string]string `json:"desired,omitempty"`
	Reported map[string]string `json:"reported,omitempty"`
}
