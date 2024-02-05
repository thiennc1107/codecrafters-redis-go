package presenter

import (
	"context"
	"fmt"
	"io"
)

func WriteResponse(ctx context.Context, writer io.Writer, value string) error {
	value = fmt.Sprintf("+%s\r\n", value)
	_, err := writer.Write([]byte(value))
	if err != nil {
		return err
	}
	return nil
}
