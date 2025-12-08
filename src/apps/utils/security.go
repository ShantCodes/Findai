package utils

import (
	"net/url"

	"github.com/microcosm-cc/bluemonday"
)

/*
* XSS Prevention
 */
func SanitizeURLValues(vals url.Values, p *bluemonday.Policy) {
	for key, list := range vals {
		for i, val := range list {
			vals[key][i] = p.Sanitize(val)
		}
	}
}

func SanitizeMap(m map[string]interface{}, p *bluemonday.Policy) {
	for k, v := range m {
		switch t := v.(type) {
		case string:
			m[k] = p.Sanitize(t)
		case map[string]interface{}:
			SanitizeMap(t, p)
		case []interface{}:
			for i, item := range t {
				switch val := item.(type) {
				case string:
					t[i] = p.Sanitize(val)
				case map[string]interface{}:
					SanitizeMap(val, p)
				}
			}
		}
	}
}
