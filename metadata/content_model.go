package metadata

import (
	"fmt"
)

type Content struct {
	UUID        string       `json:"uuid"`
	Identifiers []Identifier `json:"identifiers"`
}

type Identifier struct {
	Authority string `json:"authority"`
}

var sourceMap = map[string]string{
	"http://api.ft.com/system/FTCOM-METHODE":    "METHODE",
	"http://api.ft.com/system/FT-LABS-WP-1-335": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-24":  "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-101": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-201": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-3":   "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-272": "BLOGS",
}

func (c Content) getSource() (string, error) {
	source := sourceMap[c.Identifiers[0].Authority]
	if len(c.Identifiers) == 1 {
		return source, nil
	}

	//handle multiple sources
	isConsistent := true
	for i := 1; i < len(c.Identifiers); i++ {
		nextSource := sourceMap[c.Identifiers[i].Authority]
		if nextSource != source {
			isConsistent = false
			break
		}
		source = nextSource
	}

	if isConsistent {
		return source, nil
	}
	return source, fmt.Errorf("Cannot find source of content with UUID=[%s]", c.UUID)
}
