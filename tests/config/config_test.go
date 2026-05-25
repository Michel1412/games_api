package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"games_api/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAppliesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ENV", "")
	t.Setenv("GCP_PROJECT_ID", "")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS_JSON", "")

	cfg := config.Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "dev", cfg.Env)
	assert.False(t, cfg.IsProd())
	assert.NoError(t, cfg.ValidateProd())
}

func TestLoadReadsEnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ENV", "prod")
	t.Setenv("GCP_PROJECT_ID", "my-project")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/tmp/creds.json")

	cfg := config.Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.True(t, cfg.IsProd())
	assert.Equal(t, "my-project", cfg.GCPProjectID)
	assert.NoError(t, cfg.ValidateProd())
}

func TestValidateProdRequiresProjectID(t *testing.T) {
	cfg := config.Config{Env: "prod", GoogleApplicationCredentials: "/tmp/creds.json"}

	err := cfg.ValidateProd()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "GCP_PROJECT_ID")
}

func TestValidateProdRequiresCredentials(t *testing.T) {
	cfg := config.Config{Env: "prod", GCPProjectID: "my-project"}

	err := cfg.ValidateProd()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "GOOGLE_APPLICATION_CREDENTIALS")
}

func TestPrepareGoogleCredentialsFileWritesJSON(t *testing.T) {
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	cfg := config.Config{GoogleApplicationCredentialsJSON: `{"type":"service_account"}`}

	require.NoError(t, cfg.PrepareGoogleCredentialsFile())

	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	require.NotEmpty(t, path)
	assert.Equal(t, filepath.Join(os.TempDir(), "google-application-credentials.json"), path)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, `{"type":"service_account"}`, string(content))

	_ = os.Remove(path)
}

func TestPrepareGoogleCredentialsFileSkipsWhenJSONEmpty(t *testing.T) {
	cfg := config.Config{}

	assert.NoError(t, cfg.PrepareGoogleCredentialsFile())
}
