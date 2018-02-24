package main

import (
	"io/ioutil"
	"fmt"
	"encoding/xml"
)

func parseXML(file string) (Feedback, error) {

	var body []byte
	var report Feedback

	verbose("parsing %s\n", file)
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return Feedback{}, fmt.Errorf("unable to read file: %v", err)
	}

	if err := xml.Unmarshal(body, &report); err != nil {
		return Feedback{}, fmt.Errorf("unable to parse file: %v", err)
	}
	return report, nil
}
