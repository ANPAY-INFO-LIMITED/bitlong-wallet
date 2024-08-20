package api

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	terminal "github.com/lightninglabs/lightning-terminal"
	"github.com/lightningnetwork/lnd"
	"github.com/lightningnetwork/lnd/signal"
	"os"
)

func StartLitd() {
	err := terminal.New().Run()
	var flagErr *flags.Error
	isFlagErr := errors.As(err, &flagErr)
	if err != nil && (!isFlagErr || flagErr.Type != flags.ErrHelp) {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startLnd() {
	// Hook interceptor for os signals.
	shutdownInterceptor, err := signal.Intercept()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Load the configuration, and parse any command line options. This
	// function will also set up logging properly.
	loadedConfig, err := lnd.LoadConfig(shutdownInterceptor)
	if err != nil {
		var e *flags.Error
		if !errors.As(err, &e) || e.Type != flags.ErrHelp {
			// Print error if not due to help request.
			err = fmt.Errorf("failed to load config: %w", err)
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Help was requested, exit normally.
		os.Exit(0)
	}
	implCfg := loadedConfig.ImplementationConfig(shutdownInterceptor)

	// Call the "real" main in a nested manner so the defers will properly
	// be executed in the case of a graceful shutdown.
	if err = lnd.Main(
		loadedConfig, lnd.ListenerCfg{}, implCfg, shutdownInterceptor,
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
