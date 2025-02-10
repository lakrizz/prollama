package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/lakrizz/prollama/cmd"
	"github.com/lakrizz/prollama/config"
)

func main() {
	var cfg = &config.Config{
		Model:    "",
		Repo:     "",
		Endpoint: "http://localhost:11434",
		Debug:    false,
		Timeout:  30,
	}

	app := &cli.App{
		Name:  "prollama",
		Usage: "Review Github repository with Ollama AI.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "model",
				Value:       "",
				Usage:       "The Ollama model to use. If left blank, we will try to auto-detect a well fitting model.",
				Destination: &cfg.Model,
			},
			&cli.StringFlag{
				Name:        "repo",
				Value:       "",
				Usage:       "Github repository to review, e.g. 'lakrizz/prollama'. If left blank, it will determine the repository of the current directory",
				Destination: &cfg.Repo,
			},
			&cli.StringFlag{
				Name:        "endpoint",
				Value:       cfg.Endpoint,
				Usage:       "API endpoint for remote Ollama instance, including the listening port.",
				Destination: &cfg.Endpoint,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Value:       cfg.Debug,
				Usage:       "Enable debug output.",
				Destination: &cfg.Debug,
			},
			&cli.IntFlag{
				Name:        "timeout",
				Value:       cfg.Timeout,
				Usage:       "Ollama request timeout in seconds.",
				Destination: &cfg.Timeout,
			},
			&cli.StringFlag{
				Name:        "access-token",
				Value:       "",
				Usage:       "Access token for Ollama authentication. (optional)",
				Destination: &cfg.AccessToken,
			},
			&cli.BoolFlag{
				Name:        "no-color",
				Value:       cfg.NoColor,
				Usage:       "Disables color output.",
				Destination: &cfg.NoColor,
			},
		},
		Action: func(c *cli.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
			defer cancel()

			g, ctx := errgroup.WithContext(ctx)
			ctx = config.NewContext(ctx, cfg)

			if cfg.Debug {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}

			// enable colored output
			if !cfg.NoColor {
				// set global logger with custom options
				slog.SetDefault(slog.New(
					tint.NewHandler(os.Stderr, &tint.Options{
						Level:      slog.LevelDebug,
						TimeFormat: time.Kitchen,
					}),
				))
			}

			// Example task within the errgroup
			g.Go(func() error {
				res := make(chan error, 1)
				go func() {
					// Simulate command execution
					res <- cmd.Prollama(ctx) // Replace with actual command execution logic
				}()
				select {
				case <-ctx.Done():
					return ctx.Err()
				case err := <-res:
					return err
				}
			})

			if err := g.Wait(); err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error("cancelling execution", "msg", err.Error())
	}
}
