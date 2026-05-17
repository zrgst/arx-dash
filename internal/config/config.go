package config

import "os"

// Config holder all runtime config for appen.
//
// Vi starter med miljøvariabler for å slippe ekstra dependencies.
// Senere kan vi legge til config.toml hvis vi vil gjøre dette mer produktvennlig.
type Config struct {
	ARXBaseURL         string
	ARXUsername        string
	ARXPassword        string
	ARXAllowSelfSigned bool
	Addr               string
}

// Load leser config fra environment variables.
func Load() Config {
	return Config{
		ARXBaseURL:         os.Getenv("ARX_BASE_URL"),
		ARXUsername:        os.Getenv("ARX_USERNAME"),
		ARXPassword:        os.Getenv("ARX_PASSWORD"),
		ARXAllowSelfSigned: os.Getenv("ARX_ALLOW_SELF_SIGNED") == "true",
		Addr:               envOrDefault("ARX_DASHBOARD_ADDR", ":8080"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
