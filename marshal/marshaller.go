package marshal

import (
	"encoding/json"
	"github.com/juju/errors"
)

// Marshal to converts the log to json string
func Marshal(log interface{}) (string, error) {
	if b, err := json.Marshal(log); err != nil {
		return "", errors.Annotatef(err, "could not marshal the log %v", log)
	} else {
		return string(b), nil
	}
}

// Unmarshal to unmarshal the data bytes and fills in the Log object
func Unmarshal(data []byte, log interface{}) error {
	if err := json.Unmarshal(data, log); err != nil {
		return errors.Annotatef(err, "could not unmarshal the log data: %s", data)
	}

	return nil
}
