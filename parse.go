package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

func ParseXML(r io.ReadCloser) (Feedback, error) {
	var report Feedback

	if r == nil {
		return Feedback{}, fmt.Errorf("r is nil")
	}

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return Feedback{}, errors.Wrap(err, "readall")
	}

	if err := xml.Unmarshal(body, &report); err != nil {
		return Feedback{}, errors.Wrap(err, "unmarshall")
	}
	return report, nil
}
