package converter

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	historypb "go.temporal.io/api/history/v1"
	"google.golang.org/protobuf/proto"
)

type testStruct struct {
	Name string
	Age  int
}

func TestProtoJsonPayloadConverter_Google(t *testing.T) {
	pc := NewProtoJSONPayloadConverter()

	wt := &historypb.HistoryEvent{
		EventId:   1978,
		EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_TIMED_OUT,
		Attributes: &historypb.HistoryEvent_WorkflowTaskTimedOutEventAttributes{WorkflowTaskTimedOutEventAttributes: &historypb.WorkflowTaskTimedOutEventAttributes{
			ScheduledEventId: 2,
			TimeoutType:      enumspb.TIMEOUT_TYPE_SCHEDULE_TO_START,
		}}}
	payload, err := pc.ToPayload(wt)
	require.NoError(t, err)
	wt2 := &historypb.HistoryEvent{}
	err = pc.FromPayload(payload, &wt2)
	require.NoError(t, err)
	assert.Equal(t, int64(1978), wt2.EventId)

	var wt3 *historypb.HistoryEvent
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Equal(t, int64(1978), wt3.EventId)

	var wt4 historypb.HistoryEvent
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Equal(t, int64(1978), wt3.EventId)

	s := pc.ToString(payload)
	assert.JSONEq(t, `{"eventId":"1978","eventType":"EVENT_TYPE_WORKFLOW_TASK_TIMED_OUT","workflowTaskTimedOutEventAttributes":{"scheduledEventId":"2","timeoutType":"TIMEOUT_TYPE_SCHEDULE_TO_START"}}`, s)

	// Add additional field to payload data
	payload.Data = []byte(`{"eventId":"1978","eventType":"EVENT_TYPE_WORKFLOW_TASK_TIMED_OUT","workflowTaskTimedOutEventAttributes":{"scheduledEventId":"2","timeoutType":"TIMEOUT_TYPE_SCHEDULE_TO_START"},"newField":"newValue"}`)
	// Should fail, unknown field
	wt5 := &Gogo{}
	err = pc.FromPayload(payload, &wt5)
	require.Error(t, err)

	// Shouldn't fail, unknown fields are allowed
	pc = NewProtoJSONPayloadConverterWithOptions(ProtoJSONPayloadConverterOptions{
		AllowUnknownFields: true,
	})
	wt6 := &Gogo{}
	err = pc.FromPayload(payload, &wt6)
	require.NoError(t, err)
}

func TestProtoJsonPayloadConverter_Gogo(t *testing.T) {
	pc := NewProtoJSONPayloadConverter()

	wt := &Gogo{
		Name:     "qwe",
		Birthday: 12,
		Type:     Typegogo_TYPEGOGO_R,
		Values:   &Gogo_ValueS{ValueS: "asd"},
	}
	payload, err := pc.ToPayload(wt)
	require.NoError(t, err)
	wt2 := &Gogo{}
	err = pc.FromPayload(payload, &wt2)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt2.Name)
	assert.Equal(t, int64(12), wt2.Birthday)

	var wt3 *Gogo
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt3.Name)
	assert.Equal(t, int64(12), wt3.Birthday)

	var wt4 Gogo
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt4.Name)
	assert.Equal(t, int64(12), wt4.Birthday)

	s := pc.ToString(payload)
	assert.Equal(t, `{"name":"qwe","birthday":"12","type":"TYPEGOGO_R","valueS":"asd"}`, strings.Replace(s, " ", "", -1))

	// Add additional field to payload data
	payload.Data = []byte(`{"name":"qwe","birthday":"12","type":"TYPEGOGO_R","valueS":"asd","newField":"newValue"}`)
	// Should fail, unknown field
	wt5 := &Gogo{}
	err = pc.FromPayload(payload, &wt5)
	require.Error(t, err)

	// Shouldn't fail, unknown fields are allowed
	pc = NewProtoJSONPayloadConverterWithOptions(ProtoJSONPayloadConverterOptions{
		AllowUnknownFields: true,
	})
	wt6 := &Gogo{}
	err = pc.FromPayload(payload, &wt6)
	require.NoError(t, err)
}

func TestProtoPayloadConverter_Google(t *testing.T) {
	pc := NewProtoPayloadConverter()

	wt := &historypb.HistoryEvent{
		EventId:   1978,
		EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_TIMED_OUT,
		Attributes: &historypb.HistoryEvent_WorkflowTaskTimedOutEventAttributes{WorkflowTaskTimedOutEventAttributes: &historypb.WorkflowTaskTimedOutEventAttributes{
			ScheduledEventId: 2,
			TimeoutType:      enumspb.TIMEOUT_TYPE_SCHEDULE_TO_START,
		}}}
	payload, err := pc.ToPayload(wt)
	require.NoError(t, err)
	wt2 := &historypb.HistoryEvent{}
	err = pc.FromPayload(payload, &wt2)
	require.NoError(t, err)
	assert.Equal(t, int64(1978), wt2.EventId)

	var wt3 *historypb.HistoryEvent
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Equal(t, int64(1978), wt3.EventId)

	var wt4 historypb.HistoryEvent
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Equal(t, int64(1978), wt4.EventId)

	s := pc.ToString(payload)
	assert.Equal(t, "CLoPGAhqBAgCGAI", s)
	assert.Equal(t, "temporal.api.history.v1.HistoryEvent", string(payload.Metadata[MetadataMessageType]))
}

func TestProtoPayloadConverter_Gogo(t *testing.T) {
	pc := NewProtoPayloadConverter()

	wt := &Gogo{
		Name:     "qwe",
		Birthday: 12,
		Type:     Typegogo_TYPEGOGO_R,
		Values:   &Gogo_ValueS{ValueS: "asd"},
	}
	payload, err := pc.ToPayload(wt)
	require.NoError(t, err)
	wt2 := &Gogo{}
	err = pc.FromPayload(payload, &wt2)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt2.Name)

	var wt3 *Gogo
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt3.Name)

	var wt4 Gogo
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt4.Name)

	s := pc.ToString(payload)
	assert.Equal(t, "CgNxd2UQDDgBQgNhc2Q", s)
	assert.Equal(t, "temporal.sdk.converter.Gogo", string(payload.Metadata[MetadataMessageType]))
}

func TestJsonPayloadConverter(t *testing.T) {
	pc := NewJSONPayloadConverter()

	wt := testStruct{Name: "qwe"}
	payload, err := pc.ToPayload(wt)
	require.NoError(t, err)
	wt2 := testStruct{}
	err = pc.FromPayload(payload, &wt2)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt2.Name)

	var wt3 *testStruct
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt3.Name)

	var wt4 testStruct
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt4.Name)

	s := pc.ToString(payload)
	assert.Equal(t, `{"Name":"qwe","Age":0}`, s)
}

func TestProtoJsonPayloadConverter_Nil(t *testing.T) {
	pc := NewProtoJSONPayloadConverter()

	var wt1 *Gogo
	payload, err := pc.ToPayload(wt1)
	require.NoError(t, err)
	assert.Equal(t, "null", string(payload.Data))

	wt1 = &Gogo{Name: "qwe"}
	err = pc.FromPayload(payload, &wt1)
	require.NoError(t, err)
	assert.Nil(t, wt1)

	var wt2 *commonpb.WorkflowType
	payload, err = pc.ToPayload(wt2)
	require.NoError(t, err)
	assert.Equal(t, "null", string(payload.Data))

	wt2 = &commonpb.WorkflowType{Name: "qwe"}
	err = pc.FromPayload(payload, &wt2)
	require.NoError(t, err)
	assert.Nil(t, wt2)

	var wt3 interface{}
	payload, err = pc.ToPayload(wt3)
	require.NoError(t, err)
	assert.Equal(t, "null", string(payload.Data))

	wt3 = 123
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Nil(t, wt3)

	var wt4 *interface{}
	payload, err = pc.ToPayload(wt4)
	require.NoError(t, err)
	assert.Equal(t, "null", string(payload.Data))

	i := interface{}(123)
	wt4 = &i
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Nil(t, wt4)
}

func TestJsonPayloadConverter_Nil(t *testing.T) {
	pc := NewJSONPayloadConverter()

	var wt1 *testStruct
	payload, err := pc.ToPayload(wt1)
	require.NoError(t, err)
	assert.Equal(t, "null", string(payload.Data))

	wt1 = &testStruct{Name: "qwe"}
	err = pc.FromPayload(payload, &wt1)
	require.NoError(t, err)
	assert.Nil(t, wt1)

	var wt3 interface{}
	payload, err = pc.ToPayload(wt3)
	require.NoError(t, err)
	assert.Equal(t, "null", string(payload.Data))

	wt3 = 123
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Nil(t, wt3)

	var wt4 *interface{}
	payload, err = pc.ToPayload(wt4)
	require.NoError(t, err)
	assert.Equal(t, "null", string(payload.Data))

	i := interface{}(123)
	wt4 = &i
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Nil(t, wt4)
}

func TestNilPayloadConverter(t *testing.T) {
	pc := NewNilPayloadConverter()

	var wt1 *testStruct
	payload, err := pc.ToPayload(wt1)
	require.NoError(t, err)
	assert.Nil(t, payload.Data)

	wt1 = &testStruct{Name: "qwe"}
	err = pc.FromPayload(payload, &wt1)
	require.NoError(t, err)
	assert.Nil(t, wt1)

	var wt3 interface{}
	payload, err = pc.ToPayload(wt3)
	require.NoError(t, err)
	assert.Nil(t, payload.Data)

	wt3 = 123
	err = pc.FromPayload(payload, &wt3)
	require.NoError(t, err)
	assert.Nil(t, wt3)

	var wt4 *interface{}
	payload, err = pc.ToPayload(wt4)
	require.NoError(t, err)
	assert.Nil(t, payload.Data)

	i := interface{}(123)
	wt4 = &i
	err = pc.FromPayload(payload, &wt4)
	require.NoError(t, err)
	assert.Nil(t, wt4)
}

func TestProtoPayloadConverter_WithOptions(t *testing.T) {
	pc := NewProtoPayloadConverterWithOptions(ProtoPayloadConverterOptions{ExcludeProtobufMessageTypes: true})

	wt := commonpb.WorkflowType{Name: "qwe"}
	payload, err := pc.ToPayload(&wt)
	require.NoError(t, err)

	_, ok := payload.Metadata[MetadataMessageType]
	assert.False(t, ok)
}

func TestProtoJSONPayloadConverter_WithOptions(t *testing.T) {
	pc := NewProtoJSONPayloadConverterWithOptions(ProtoJSONPayloadConverterOptions{ExcludeProtobufMessageTypes: true})

	wt := commonpb.WorkflowType{Name: "qwe"}
	payload, err := pc.ToPayload(&wt)
	require.NoError(t, err)

	_, ok := payload.Metadata[MetadataMessageType]
	assert.False(t, ok)
}

func TestProtoJsonPayloadConverter_FromPayload_Errors(t *testing.T) {
	pc := NewProtoJSONPayloadConverter()

	wt := commonpb.WorkflowType{Name: "qwe"}
	payload, err := pc.ToPayload(&wt)
	require.NoError(t, err)

	var wt2 *int
	err = pc.FromPayload(payload, &wt2)
	require.Error(t, err)
	assert.Equal(t, "type: *int: type doesn't implement proto.Message", err.Error())
	assert.True(t, errors.Is(err, ErrTypeNotImplementProtoMessage))

	var wt3 *commonpb.WorkflowType
	err = pc.FromPayload(payload, wt3)
	require.Error(t, err)
	assert.Equal(t, "type: *common.WorkflowType: unable to set value", err.Error())
	assert.True(t, errors.Is(err, ErrUnableToSetValue))

	// But 31, 32, and 33 work
	var wt31 commonpb.WorkflowType
	err = pc.FromPayload(payload, &wt31)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt31.Name)

	wt32 := &commonpb.WorkflowType{}
	err = pc.FromPayload(payload, wt32)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt32.Name)

	var wt33 *commonpb.WorkflowType //lint:ignore S1021 as it indicates exactly this case
	wt33 = &commonpb.WorkflowType{}
	err = pc.FromPayload(payload, wt33)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt33.Name)

	var wt5 interface{}
	err = pc.FromPayload(payload, wt5)
	require.Error(t, err)
	assert.Equal(t, "type: <nil>: not a pointer type", err.Error())
	assert.True(t, errors.Is(err, ErrValuePtrIsNotPointer))

	var wt6 *interface{}
	err = pc.FromPayload(payload, wt6)
	require.Error(t, err)
	assert.Equal(t, "type: *interface {}: unable to set value", err.Error())
	assert.True(t, errors.Is(err, ErrUnableToSetValue))

	// supported by JSON serializer but not by ProtoJson
	var wt7 interface{}
	err = pc.FromPayload(payload, &wt7)
	require.Error(t, err)
	assert.Equal(t, "value type: interface {}: must be a concrete type, not interface", err.Error())
	assert.True(t, errors.Is(err, ErrValuePtrMustConcreteType))

	var wt8 proto.Message
	err = pc.FromPayload(payload, &wt8)
	require.Error(t, err)
	assert.Equal(t, "value type: protoreflect.ProtoMessage: must be a concrete type, not interface", err.Error())
	assert.True(t, errors.Is(err, ErrValuePtrMustConcreteType))

	var wt9 string
	err = pc.FromPayload(payload, &wt9)
	require.Error(t, err)
	assert.Equal(t, "type: *string: type doesn't implement proto.Message", err.Error())
	assert.True(t, errors.Is(err, ErrTypeNotImplementProtoMessage))
}

func TestProtoPayloadConverter_FromPayload_Errors(t *testing.T) {
	pc := NewProtoPayloadConverter()

	wt := commonpb.WorkflowType{Name: "qwe"}
	payload, err := pc.ToPayload(&wt)
	require.NoError(t, err)

	var wt2 *int
	err = pc.FromPayload(payload, &wt2)
	require.Error(t, err)
	assert.Equal(t, "type: *int: type doesn't implement proto.Message", err.Error())
	assert.True(t, errors.Is(err, ErrTypeNotImplementProtoMessage))

	var wt3 *commonpb.WorkflowType
	err = pc.FromPayload(payload, wt3)
	require.Error(t, err)
	assert.Equal(t, "type: *common.WorkflowType: unable to set value", err.Error())
	assert.True(t, errors.Is(err, ErrUnableToSetValue))

	// But 31, 32, and 33 work
	var wt31 commonpb.WorkflowType
	err = pc.FromPayload(payload, &wt31)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt31.Name)

	wt32 := &commonpb.WorkflowType{}
	err = pc.FromPayload(payload, wt32)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt32.Name)

	var wt33 *commonpb.WorkflowType //lint:ignore S1021 as it indicates exactly this case
	wt33 = &commonpb.WorkflowType{}
	err = pc.FromPayload(payload, wt33)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt33.Name)

	var wt5 interface{}
	err = pc.FromPayload(payload, wt5)
	require.Error(t, err)
	assert.Equal(t, "type: <nil>: not a pointer type", err.Error())
	assert.True(t, errors.Is(err, ErrValuePtrIsNotPointer))

	var wt6 *interface{}
	err = pc.FromPayload(payload, wt6)
	require.Error(t, err)
	assert.Equal(t, "type: *interface {}: unable to set value", err.Error())
	assert.True(t, errors.Is(err, ErrUnableToSetValue))

	// supported by JSON serializer but not by ProtoJson
	var wt7 interface{}
	err = pc.FromPayload(payload, &wt7)
	require.Error(t, err)
	assert.Equal(t, "value type: interface {}: must be a concrete type, not interface", err.Error())
	assert.True(t, errors.Is(err, ErrValuePtrMustConcreteType))

	var wt8 proto.Message
	err = pc.FromPayload(payload, &wt8)
	require.Error(t, err)
	assert.Equal(t, "value type: protoreflect.ProtoMessage: must be a concrete type, not interface", err.Error())
	assert.True(t, errors.Is(err, ErrValuePtrMustConcreteType))

	var wt9 string
	err = pc.FromPayload(payload, &wt9)
	require.Error(t, err)
	assert.Equal(t, "type: *string: type doesn't implement proto.Message", err.Error())
	assert.True(t, errors.Is(err, ErrTypeNotImplementProtoMessage))
}

func TestJsonPayloadConverter_FromPayload_Errors(t *testing.T) {
	pc := NewJSONPayloadConverter()

	wt := testStruct{Name: "qwe"}
	payload, err := pc.ToPayload(wt)
	require.NoError(t, err)

	var wt2 *int
	err = pc.FromPayload(payload, &wt2)
	require.Error(t, err)
	assert.Equal(t, "unable to decode: json: cannot unmarshal object into Go value of type int", err.Error())

	var wt3 *testStruct
	err = pc.FromPayload(payload, wt3)
	require.Error(t, err)
	assert.Equal(t, "unable to decode: json: Unmarshal(nil *converter.testStruct)", err.Error())

	// But 31, 32, and 33 work
	var wt31 testStruct
	err = pc.FromPayload(payload, &wt31)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt31.Name)

	wt32 := &testStruct{}
	err = pc.FromPayload(payload, wt32)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt32.Name)

	var wt33 *testStruct //lint:ignore S1021 as it indicates exactly this case
	wt33 = &testStruct{}
	err = pc.FromPayload(payload, wt33)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt33.Name)

	var wt4 testStruct
	err = pc.FromPayload(payload, wt4)
	require.Error(t, err)
	assert.Equal(t, "unable to decode: json: Unmarshal(non-pointer converter.testStruct)", err.Error())

	var wt5 interface{}
	err = pc.FromPayload(payload, wt5)
	require.Error(t, err)
	assert.Equal(t, "unable to decode: json: Unmarshal(nil)", err.Error())

	var wt6 *interface{}
	err = pc.FromPayload(payload, wt6)
	require.Error(t, err)
	assert.Equal(t, "unable to decode: json: Unmarshal(nil *interface {})", err.Error())

	// supported by JSON serializer (wt7 will be map[string]interface{})
	var wt7 interface{}
	err = pc.FromPayload(payload, &wt7)
	require.NoError(t, err)
	assert.Equal(t, "qwe", wt7.(map[string]interface{})["Name"])
}
