package security

type Config struct {
	SaltLength   uint32 `yaml:"saltLength" env:"SALT_LENGTH" env-default:"16"`
	ArgonTime    uint32 `yaml:"argonTime" env:"ARGON_TIME" env-default:"1"`
	ArgonMemory  uint32 `yaml:"argonMemory" env:"ARGON_MEMORY" env-default:"65536"`
	ArgonThreads uint8  `yaml:"argonThreads" env:"ARGON_THREADS" env-default:"4"`
	ArgonKeyLen  uint32 `yaml:"argonKeyLen" env:"ARGON_KEY_LEN" env-default:"32"`
}
