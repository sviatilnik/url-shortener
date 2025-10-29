package testdata

// generate:reset
type TestConfig struct {
	BaseURL    string
	MaxRetries int
	Timeout    int
	Enabled    bool
	Headers    map[string]string
	Options    []string
	Logger     *TestLogger
}

// generate:reset
type TestLogger struct {
	Level  string
	Format string
	Output []string
}
