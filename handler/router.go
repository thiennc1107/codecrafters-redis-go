package handler

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type ReadCommandFunc func(ctx context.Context, stream io.Reader) (Command, error)

type WriteResponseFunc func(ctx context.Context, writer io.Writer, value string) error

type Context struct {
	context.Context
	Command Command
}

type Handler struct {
	timeOut       uint
	readCommand   ReadCommandFunc
	writeResponse WriteResponseFunc
}

func WithTimeOut(timeout uint) func(*Handler) {
	return func(h *Handler) {
		h.timeOut = timeout
	}
}

func NewHandler(
	readCommandFunc ReadCommandFunc,
	writeResponseFunc WriteResponseFunc,
	options ...func(*Handler),
) *Handler {
	hdl := &Handler{
		readCommand:   readCommandFunc,
		writeResponse: writeResponseFunc,
	}
	for _, o := range options {
		o(hdl)
	}
	return hdl
}

func (h *Handler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	if h.timeOut > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(h.timeOut)*time.Second)
	}
	h.HandleSession(ctx, cancel, conn)

	defer cancel()
}

func FormatErrorFailedToReadCommand(err error) string {
	return fmt.Sprintf("ERR: failed to read command, %v", err)
}

func (h *Handler) handleReadCommandError(ctx context.Context, err error, conn io.ReadWriter) {
	h.writeResponse(ctx, conn, FormatErrorFailedToReadCommand(err))
}

func (h *Handler) handleCommandError(ctx context.Context, err error, conn io.ReadWriter) {

}

func (h *Handler) HandleSession(ctx context.Context, cancelFunc context.CancelFunc, sess io.ReadWriter) {
	cmd, err := h.readCommand(ctx, sess)
	if err != nil {
		h.handleReadCommandError(ctx, err, sess)
		cancelFunc()
		return
	}

	cmdHandler, err := GetCommandHandler(cmd.Name)
	if err != nil {
		h.handleReadCommandError(ctx, err, sess)
		return
	}

	res, err := cmdHandler.Execute(Context{
		Context: ctx,
		Command: cmd,
	})
	if err != nil {
		h.handleCommandError(ctx, err, sess)
	}

	h.writeResponse(ctx, sess, res)
}
