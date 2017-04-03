package metadata

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
	"http://api.ft.com/system/FT-LABS-WP-1-3":   "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-333": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-91":  "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-332": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-101": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-201": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-2":   "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-272": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-51":  "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-242": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-171": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-12":  "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-10":  "BLOGS",
	"http://api.ft.com/system/FT-CLAMO":         "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-252": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-106": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-302": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-9":   "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-312": "BLOGS",
	"http://api.ft.com/system/FT-LABS-WP-1-292": "BLOGS",
}

func (c Content) getSource() (string, bool) {
	source := sourceMap[c.Identifiers[0].Authority]
	if len(c.Identifiers) == 1 {
		return source, true
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
		return source, true
	}
	return source, false
}
