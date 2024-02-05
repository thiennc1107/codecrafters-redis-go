package handler_test

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/handler"
	"github.com/codecrafters-io/redis-starter-go/parser"
	"github.com/codecrafters-io/redis-starter-go/presenter"
)

type handlerTestCase struct {
	input  []byte
	output []byte
}

func TestHandler(t *testing.T) {
	tcs := []handlerTestCase{
		{
			input:  []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"),
			output: []byte("+hey\r\n"),
		},
	}
	for i, tc := range tcs {
		testHandler(t, tc, i)
	}
}

type mockReadWriter struct {
	io.Reader
	io.Writer
}

func testHandler(t *testing.T, tc handlerTestCase, i int) {
	inputBuffer := bytes.NewReader(tc.input)
	outputBuffer := bytes.NewBuffer([]byte{})
	conn := mockReadWriter{
		Reader: inputBuffer,
		Writer: outputBuffer,
	}
	hdl := handler.NewHandler(parser.ReadCommand, presenter.WriteResponse)
	ctx, cancel := context.WithCancel(context.Background())
	hdl.HandleSession(ctx, cancel, conn)
	out := outputBuffer.Bytes()
	if !reflect.DeepEqual(out, tc.output) {
		t.Logf("test case %d failed, unexpected output, expect %s, got %s", i, out, tc.output)
		t.Fail()
	}
}
