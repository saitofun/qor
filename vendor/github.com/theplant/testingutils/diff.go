package testingutils

import (
	"encoding/json"

	"fmt"

	"io"
	"io/ioutil"

	"github.com/pmezard/go-difflib/difflib"
)

/*
It convert the two objects into pretty json, and diff them, output the result.
*/
func PrettyJsonDiff(expected interface{}, actual interface{}) (r string) {

	actualJson := marshalIfNotStringOrReader(actual)
	expectedJson := marshalIfNotStringOrReader(expected)

	if actualJson != expectedJson {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(expectedJson),
			B:        difflib.SplitLines(actualJson),
			FromFile: "Expected",
			ToFile:   "Actual",
			Context:  3,
		}
		r, _ = difflib.GetUnifiedDiffString(diff)
	}
	return
}

func PrintlnJson(vals ...interface{}) {
	var newvals []interface{}
	for _, v := range vals {
		if s, ok := v.(string); ok {
			newvals = append(newvals, s)
			continue
		}
		j, _ := json.MarshalIndent(v, "", "\t")
		newvals = append(newvals, "\n", string(j))
	}
	fmt.Println(newvals...)
}

func marshalIfNotStringOrReader(v interface{}) (r string) {
	var ok bool
	if r, ok = v.(string); ok {
		r = formatIfJson(r)
		return
	}

	var rd io.Reader
	if rd, ok = v.(io.Reader); ok {
		bs, _ := ioutil.ReadAll(rd)
		r = string(bs)
		return
	}
	rbytes, _ := json.MarshalIndent(v, "", "\t")
	r = string(rbytes)
	return
}

func formatIfJson(input string) (r string) {
	var inputRawM json.RawMessage
	var err error
	err = json.Unmarshal([]byte(input), &inputRawM)
	if err != nil {
		r = input
		return
	}

	var rbytes []byte
	rbytes, err = json.MarshalIndent(inputRawM, "", "\t")
	if err != nil {
		r = input
		return
	}
	r = string(rbytes)
	return
}
