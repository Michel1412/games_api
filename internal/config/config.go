package config

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Port                             string
	Env                              string
	GCPProjectID                     string
	GoogleApplicationCredentials     string
	GoogleApplicationCredentialsJSON string
}

func Load() Config {
	return Config{
		Port:                             getEnv("PORT", "8080"),
		Env:                              getEnv("ENV", "dev"),
		GCPProjectID:                     os.Getenv("GCP_PROJECT_ID"),
		GoogleApplicationCredentials:     os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		GoogleApplicationCredentialsJSON: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"),
	}
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func (c Config) ValidateProd() error {
	if !c.IsProd() {
		return nil
	}

	if c.GCPProjectID == "" {
		return errors.New("GCP_PROJECT_ID e obrigatorio em producao")
	}

	if c.GoogleApplicationCredentials == "" && c.GoogleApplicationCredentialsJSON == "" {
		return errors.New("GOOGLE_APPLICATION_CREDENTIALS ou GOOGLE_APPLICATION_CREDENTIALS_JSON e obrigatorio em producao")
	}

	return nil
}

func (c Config) PrepareGoogleCredentialsFile() error {
	if c.GoogleApplicationCredentialsJSON == "" {
		return nil
	}

	credentialsPath := filepath.Join(os.TempDir(), "google-application-credentials.json")
	if err := os.WriteFile(credentialsPath, []byte(c.GoogleApplicationCredentialsJSON), 0o600); err != nil {
		return err
	}

	return os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialsPath)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
