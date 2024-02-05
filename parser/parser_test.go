package parser_test

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/handler"
	"github.com/codecrafters-io/redis-starter-go/parser"
)

type readCommandTestcase struct {
	commandData []byte
	res         handler.Command
}

func (tc *readCommandTestcase) check(index int, t *testing.T) {
	res, err := parser.ReadCommand(context.Background(), bytes.NewBuffer(tc.commandData))
	if err != nil {
		t.Logf("Test case %d failed with error, case %v, err %v", index, tc, err)
		t.Fail()
	}

	if !reflect.DeepEqual(res, tc.res) {
		t.Logf("Test case %d failed, wrong output, expected %v, output %v", index, tc.res, res)
		t.Fail()
	}
}

func TestReadCommain(t *testing.T) {
	tcs := []readCommandTestcase{
		{
			commandData: []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"),
			res: handler.Command{
				Name:      "ECHO",
				ArgLength: 1,
				Args:      []string{"hey"},
			},
		},
		{
			commandData: []byte("*1\r\n$4\r\nping\r\n"),
			res: handler.Command{
				Name:      "ping",
				ArgLength: 0,
				Args:      []string{},
			},
		},
	}

	for i, tc := range tcs {
		tc.check(i, t)
	}
}
