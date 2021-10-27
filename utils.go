package bitbucket

import (
	"bytes"
	"encoding/json"
)

//A way to transcode (map[string]interface{}) to struct
func m2s(in, out interface{}) {
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(in)
	_ = json.NewDecoder(buf).Decode(out)
}
