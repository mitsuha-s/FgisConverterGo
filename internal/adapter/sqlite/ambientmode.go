package sqlite

import "arshin/internal/domain"

var ambientMode = map[string]domain.SourceAmbient{
	"empty":  domain.AmbientModeEmpty,
	"file":   domain.AmbientModeFile,
	"random": domain.AmbientModeRandom,
}
