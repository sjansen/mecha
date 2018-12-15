package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadProcfile(t *testing.T) {
	require := require.New(t)

	expected := map[string]string{
		"web":   "./scripts/run-in-app-venv ./manage.py runserver",
		"tasks": "./scripts/run-in-app-venv celery -A app worker -l info",
	}

	r, err := os.Open("testdata/Procfile")
	require.NoError(err)

	actual, err := ReadProcfile(r)
	require.NoError(err)
	require.Equal(expected, actual)
}
