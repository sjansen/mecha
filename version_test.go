package main

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	require := require.New(t)

	_, err := semver.NewVersion(version)
	require.NoError(err)
}
