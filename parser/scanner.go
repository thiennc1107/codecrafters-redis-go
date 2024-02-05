package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/handler"
)

func ReadCommand(ctx context.Context, stream io.Reader) (handler.Command, error) {
	scanner := bufio.NewScanner(stream)
	scanner.Split(bufio.ScanLines)
	cmd := handler.Command{}

	l, err := scanArgsLength(scanner)
	if err != nil {
		return cmd, err
	}
	cmd.ArgLength = l

	n, err := scanCommandName(scanner)
	if err != nil {
		return cmd, err
	}
	cmd.Name = n

	args, err := scanCommandArgs(scanner, cmd.ArgLength)
	if err != nil {
		return cmd, err
	}
	cmd.Args = args
	return cmd, nil
}

func scanArgsLength(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, io.ErrUnexpectedEOF
	}
	header := scanner.Text()
	if header[0] != '*' {
		return 0, fmt.Errorf("invalid command header, header must start with *")
	}
	n, err := strconv.Atoi(header[1:])
	if err != nil {
		return n, err
	}
	return n - 1, nil
}

func scanCommandName(scanner *bufio.Scanner) (string, error) {
	for i := 0; i < 2; i++ {
		if !scanner.Scan() {
			return "", io.ErrUnexpectedEOF
		}
		val := scanner.Text()
		if val[0] == '$' {
			continue
		}
		return val, nil
	}

	return "", fmt.Errorf("invalid command name stream")
}

func scanCommandArgs(scanner *bufio.Scanner, argLength int) ([]string, error) {
	args := []string{}
	for n := 0; n <= argLength; n++ {
		if !scanner.Scan() {
			return nil, io.ErrUnexpectedEOF
		}
		resp := scanner.Text()
		switch resp[0] {
		case '$':
			continue
		default:
			args = append(args, resp)
		}
	}
	return args, nil
}
