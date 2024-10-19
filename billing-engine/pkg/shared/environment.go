package shared

// Environment additional in this service
type Environment struct {
	// more additional environment with struct tag is environment key example:
	// ExampleHost string `env:"EXAMPLE_HOST"`
	JWTSecret string `env:"JWT_SECRET"`
}

var sharedEnv Environment

// GetEnv get global additional environment
func GetEnv() Environment {
	return sharedEnv
}

// SetEnv get global additional environment
func SetEnv(env Environment) {
	sharedEnv = env
}
