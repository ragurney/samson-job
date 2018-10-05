package main

import (
	l "github.com/ragurney/samson-job/internal/lib"
	s "github.com/ragurney/samson-job/pkg/samson"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
)

func main() {
	zerolog.TimeFieldFormat = ""

	deployTimeout, err := strconv.Atoi(l.Env("DEPLOY_TIMEOUT", "120"))
	if err != nil {
		log.Fatal().Msg("Failed to parse DEPLOY_TIMEOUT")
	}
	pollInterval, err := strconv.Atoi(l.Env("POLL_INTERVAL", "30"))
	if err != nil {
		log.Fatal().Msg("Failed to parse POLL_INTERVAL")
	}
	project := l.Env("SAMSON_PROJECT", "")
	reference := l.Env("REFERENCE", "")
	stage := l.Env("SAMSON_STAGE", "")
	token := l.Env("SAMSON_TOKEN", "")
	url := l.Env("SAMSON_URL", "")

	log.Debug().Msg("Starting Samson deploy...")

	s.NewJob(
		s.WithDeployTimeout(deployTimeout),
		s.WithPollInterval(pollInterval),
		s.WithProject(project),
		s.WithReference(reference),
		s.WithStage(stage),
		s.WithToken(token),
		s.WithURL(url),
	).Execute()
}
