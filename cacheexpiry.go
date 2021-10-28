package apigee

import (
	"encoding/json"
	"fmt"
)

// CacheExpiry represents the expiry settings on a cache.  This struct marshals
// and unmarshals between the json format Edge uses and a reasonably clear
// golang struct.
type CacheExpiry struct {
	ExpiryType  string
	ExpiryValue string
	ValuesNull  bool
}

// serialization format:
//  {
//    -one of-
//      "expiryDate": { "value": "{mm-dd-yyyy}" },
//      "timeoutInSec" : { "value" : "300" },
//      "timeOfDay": { "value" : "hh:mm:ss" },
//    -and-
//    "valuesNull" : false
//  }
//

// MarshalJSON implements the json.Marshaler interface. It marshals from
// the form used by Apigee Edge into a CacheExpiry struct. Eg,
//
//     {
//       "expiryDate": { "value": "{mm-dd-yyyy}" },
//       "valuesNull" : false
//     }
//
func (ce CacheExpiry) MarshalJSON() ([]byte, error) {
	valueMap := map[string]string{}
	valueMap["value"] = ce.ExpiryValue
	m1 := map[string]interface{}{}
	m1["valuesNull"] = ce.ValuesNull
	m1[ce.ExpiryType] = valueMap
	j, _ := json.Marshal(m1)
	return []byte(j), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. It unmarshals from
// a string like
//
//     {
//       "expiryDate": { "value": "{mm-dd-yyyy}" },
//       "valuesNull" : false
//     }
//
// ...into a CacheExpiry struct.
//
func (ce *CacheExpiry) UnmarshalJSON(b []byte) error {
	var m1 map[string]interface{}
	e := json.Unmarshal(b, &m1)
	if e == nil {
		entry := CacheExpiry{}
		if v, ok := m1["valuesNull"]; ok {
			entry.ValuesNull = v.(bool)
		} else {
			entry.ValuesNull = false
		}
		candidates := []string{"expiryDate", "timeOfDay", "timeoutInSec"}
		for _, k := range candidates {
			if v, ok := m1[k]; ok {
				entry.ExpiryType = k

				entry.ExpiryValue = ((v.(map[string]interface{}))["value"]).(string)
			}
		}
		*ce = entry
	}
	return nil
}

func (ce CacheExpiry) String() string {
	return fmt.Sprintf("CacheExpiry[%s,%s,%t]", ce.ExpiryType, ce.ExpiryValue, ce.ValuesNull)
}
