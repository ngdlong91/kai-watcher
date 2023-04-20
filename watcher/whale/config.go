package whale

import (
	"go.uber.org/zap"
)

type Config struct {
	URL             string
	Logger          *zap.Logger
	LevelOneLimit   string
	LevelTwoLimit   string
	LevelThreeLimit string
	LevelFourLimit  string

	AlertToken string
	AlertTo    int64
}
