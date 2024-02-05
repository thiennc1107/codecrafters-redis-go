package presenter_test

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/presenter"
)

type writeResponseTestcase struct {
	input []string
	res   []byte
}

func (tc *writeResponseTestcase) check(index int, t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	for _, val := range tc.input {
		presenter.WriteResponse(context.Background(), buffer, val)
	}
	res := buffer.Bytes()
	if !reflect.DeepEqual(res, tc.res) {
		t.Logf("wrong output, testcase %d, expected %v, got %v", index, tc.res, res)
		t.Fail()
	}
}

func TestWriteResponse(t *testing.T) {
	tcs := []writeResponseTestcase{
		{
			input: []string{
				"foo",
			},
			res: []byte("+foo\r\n"),
		},
		{
			input: []string{
				"foo",
				"bar",
			},
			res: []byte("+foo\r\n+bar\r\n"),
		},
	}
	for i, tc := range tcs {
		tc.check(i, t)
	}
}
