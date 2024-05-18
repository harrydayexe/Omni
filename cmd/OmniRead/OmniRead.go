package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/omniread/api"
)

func main() {
	ctx := context.Background()
	logger := slog.Default()

	if err := cmd.Run(ctx, api.NewHandler(logger), os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
