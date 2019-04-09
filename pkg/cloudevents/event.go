package cloudevents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/datacodec"
)

// Event represents the canonical representation of a CloudEvent.
type Event struct {
	Context EventContext
	Data    interface{}
}

func New(specVersion string) Event {
	e := &Event{}
	_ = e.SetSpecVersion(specVersion)
	return *e
}

// DataAs attempts to populate the provided data object with the event payload.
// data should be a pointer type.
func (e Event) DataAs(data interface{}) error {
	return datacodec.Decode(e.Context.GetDataMediaType(), e.Data, data)
}

// ExtensionAs returns Context.ExtensionAs(name, obj)
func (e Event) ExtensionAs(name string, obj interface{}) error {
	return e.Context.ExtensionAs(name, obj)
}

// Validate performs a spec based validation on this event.
// Validation is dependent on the spec version specified in the event context.
func (e Event) Validate() error {
	if e.Context == nil {
		return fmt.Errorf("every event conforming to the CloudEvents specification MUST include a context")
	}

	if err := e.Context.Validate(); err != nil {
		return err
	}

	// TODO: validate data.

	return nil
}

// String returns a pretty-printed representation of the Event.
func (e Event) String() string {
	b := strings.Builder{}

	b.WriteString("Validation: ")

	valid := e.Validate()
	if valid == nil {
		b.WriteString("valid\n")
	} else {
		b.WriteString("invalid\n")
	}
	if valid != nil {
		b.WriteString(fmt.Sprintf("Validation Error: \n%s\n", valid.Error()))
	}

	b.WriteString(e.Context.String())

	if e.Data != nil {
		b.WriteString("Data,\n  ")
		if strings.HasPrefix(e.DataContentType(), "application/json") {
			var prettyJSON bytes.Buffer

			data, ok := e.Data.([]byte)
			if !ok {
				var err error
				data, err = json.Marshal(e.Data)
				if err != nil {
					data = []byte(err.Error())
				}
			}
			err := json.Indent(&prettyJSON, data, "  ", "  ")
			if err != nil {
				b.Write(e.Data.([]byte))
			} else {
				b.Write(prettyJSON.Bytes())
			}
		} else {
			b.Write(e.Data.([]byte))
		}
		b.WriteString("\n")
	}
	return b.String()
}
