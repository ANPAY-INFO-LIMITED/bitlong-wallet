package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	terminal "github.com/lightninglabs/lightning-terminal"
	"github.com/lightningnetwork/lnd"
	"github.com/lightningnetwork/lnd/signal"
	"os"
)

func StartLitd() {
	err := terminal.New().Run(context.Background())
	var flagErr *flags.Error
	isFlagErr := errors.As(err, &flagErr)
	if err != nil && (!isFlagErr || flagErr.Type != flags.ErrHelp) {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}

func startLnd() {
	shutdownInterceptor, err := signal.Intercept()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	loadedConfig, err := lnd.LoadConfig(shutdownInterceptor)
	if err != nil {
		var e *flags.Error
		if !errors.As(err, &e) || e.Type != flags.ErrHelp {
			err = fmt.Errorf("failed to load config: %w", err)
			_, _ = fmt.Fprintln(os.Stderr, err)
			return
		}

		return
	}
	implCfg := loadedConfig.ImplementationConfig(shutdownInterceptor)

	if err = lnd.Main(
		loadedConfig, lnd.ListenerCfg{}, implCfg, shutdownInterceptor,
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}
