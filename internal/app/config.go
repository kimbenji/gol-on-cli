package app

type SeedMode string

type ColorMode string

const (
	SeedModeRandom      SeedMode  = "random"
	ColorModeTrueColor  ColorMode = "truecolor"
	defaultFPS                    = 5
)

type Config struct {
	FPS       int
	SeedMode  SeedMode
	ColorMode ColorMode
}

func DefaultConfig() Config {
	return Config{
		FPS:       defaultFPS,
		SeedMode:  SeedModeRandom,
		ColorMode: ColorModeTrueColor,
	}
}
