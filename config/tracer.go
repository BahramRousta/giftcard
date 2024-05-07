package config

type tracer struct {
	HostPort string `mapstructure:"hostPort" required:"true"`
	LogSpans bool   `mapstructure:"logSpans"`
}
