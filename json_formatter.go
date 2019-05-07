package logrus

import (
	"encoding/json"
	"fmt"
)

type fieldKey string
type FieldMap map[fieldKey]string

const (
	FieldKeyMsg   = "msg"
	FieldKeyLevel = "level"
	FieldKeyTime  = "time"
)

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	FixedFields     Fields

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// FieldMap allows users to customize the names of keys for various fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime: "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg: "@message",
	//    },
	// }
	FieldMap FieldMap
}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data)+len(f.FixedFields)+3)
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
