package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	terminal "github.com/lightninglabs/lightning-terminal"
	"os"
)

func PcStartLitd() {
	err := terminal.New().Run(context.Background())
	var flagErr *flags.Error
	isFlagErr := errors.As(err, &flagErr)
	if err != nil && (!isFlagErr || flagErr.Type != flags.ErrHelp) {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}
