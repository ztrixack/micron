package log

type Config struct {
	Service     string `mapstructure:"service"`
	Development bool   `mapstructure:"development"`
	File        bool   `mapstructure:"file"`
}
