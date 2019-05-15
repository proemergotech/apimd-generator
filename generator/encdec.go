package generator

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var jsons = []jsoniter.API{newJSON("param"), newJSON("query"), newJSON("json"), newJSON("geb"), newJSON("centrifuge")}

func encDec(v interface{}) (interface{}, error) {
	outM := make(map[string]interface{})
	var outI interface{}
	for _, json := range jsons {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		err = json.Unmarshal(b, &outM)
		if err != nil {
			err = json.Unmarshal(b, &outI)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if outI != nil {
				return outI, nil
			}
		}
	}

	return outM, nil
}

func newJSON(tag string) jsoniter.API {
	return jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 tag,
		OnlyTaggedField:        true,
	}.Froze()
}
