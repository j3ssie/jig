package core

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
)

// ConvertToJson convert map to JSON
func ConvertToJson(data map[string]string) string {
	if len(data) == 0 {
		return ""
	}
	jsonObj := gabs.New()
	for k, v := range data {
		// {"k", "v"}
		jsonObj.Set(v, k)
	}
	return fmt.Sprintf(jsonObj.String())
}
