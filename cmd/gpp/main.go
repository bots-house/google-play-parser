package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bots-house/google-play-parser/cmd/gpp/cli"
)

func main() {
	if err := cli.New().Run(os.Args); err != nil {
		log.Error().Err(err).Send()
	}
}
