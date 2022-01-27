package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ideally this should be in init of this package, but viper is not available yet
// JSON_LOGS should be part of docker compose and we can move this to init and use os.GetEnv()
func SetupLogger() {
	jsonLogs := viper.GetBool("JSON_LOGS")
	if !jsonLogs {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func Info(msg string) {
	log.Logger.Info().Str("project", "BASELEDGER").Msg(msg)
}

func Infof(msg string, v ...interface{}) {
	log.Logger.Info().Str("project", "BASELEDGER").Msgf(msg, v...)
}

func Warn(msg string) {
	log.Logger.Warn().Str("project", "BASELEDGER").Msg(msg)
}

func Warnf(msg string, v ...interface{}) {
	log.Logger.Warn().Str("project", "BASELEDGER").Msgf(msg, v...)
}

func Error(msg string) {
	log.Logger.Error().Str("project", "BASELEDGER").Msg(msg)
}

func Errorf(msg string, v ...interface{}) {
	log.Logger.Error().Str("project", "BASELEDGER").Msgf(msg, v...)
}

func Debug(msg string) {
	log.Logger.Debug().Str("project", "BASELEDGER").Msg(msg)
}

func Debugf(msg string, v ...interface{}) {
	log.Logger.Debug().Str("project", "BASELEDGER").Msgf(msg, v...)
}
