package logrus

import (
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	FixedFields Fields
}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data) + len(f.FixedFields) + 3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			setField(data, k, v.Error())
		default:
			setField(data, k, v)
		}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	for k, v := range f.FixedFields {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			setField(data, k, v.Error())
		default:
			setField(data, k, v)
		}
	}

	setField(data, "time", entry.Time.Format(timestampFormat))
	setField(data, "msg", entry.Message)
	setField(data, "level", entry.Level.String())

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
