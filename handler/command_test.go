package handler_test

import (
	"context"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/handler"
)

type echoTest struct {
	cmd       handler.Command
	errString string
	isPanic   bool
	desc      string
}

func TestEchoCommands(t *testing.T) {
	tcs := []echoTest{
		{
			cmd: handler.Command{
				Name: handler.EchoCommand,
				Args: []string{
					"hey",
				},
				ArgLength: 1,
			},
			errString: "",
			isPanic:   false,
			desc:      "Echo should response exactly the input argument",
		},
		{
			cmd: handler.Command{
				Name: handler.EchoCommand,
				Args: []string{
					"hey",
					"hello",
				},
				ArgLength: 2,
			},
			errString: handler.ErrEchoArgumentCount,
			isPanic:   false,
			desc:      "Echo should only accept 1 argument",
		},
		{
			cmd: handler.Command{
				Name:      handler.EchoCommand,
				Args:      []string{},
				ArgLength: 1,
			},
			errString: handler.ErrEchoArgumentCount,
			isPanic:   true,
			desc:      "Echo should panic if the ArgLength is incorrect",
		},
	}

	for i, tc := range tcs {
		testEchoCommand(t, tc, i)
	}

}

func testEchoCommand(t *testing.T, tc echoTest, i int) {
	defer func() {
		if err := recover(); err != nil {
			if !tc.isPanic {
				t.Logf("testcase %d failed, unexpected panic, case %v, err %v", i, tc, err)
				t.Fail()
				return
			}
		}
	}()
	cmdHandler, err := handler.GetCommandHandler(tc.cmd.Name)
	if err != nil {
		t.Logf("testcase %d failed, cannot get command, case %v, err %v", i, tc, err)
		t.Fail()
		return
	}

	res, err := cmdHandler.Execute(handler.Context{
		Context: context.Background(),
		Command: tc.cmd,
	})

	if tc.errString != "" {
		if err == nil {
			t.Logf("testcase %d failed, expect non-nil error, case %v", i, tc)
			t.Fail()
			return
		}
		if tc.errString != err.Error() {
			t.Logf("testcase %d failed, invalid error, case %v, got %v", i, tc, err)
			t.Fail()
			return
		}
		return
	}

	if res != tc.cmd.Args[0] {
		t.Logf("testcase %d failed, invalid respond, case %v, got %v", i, tc, err)
		t.Fail()
	}
}
