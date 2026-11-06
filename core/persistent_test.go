package core_test

import (
	"os"
	"strings"
	"testing"

	"github.com/orochaa/go-clack/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runPersistentTestPrompt(t *testing.T, output *os.File, text string) {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()
	defer w.Close()
	_, err = w.Write([]byte("\r"))
	require.NoError(t, err)

	p := core.NewPrompt(core.PromptParams[string]{
		Input: r, Output: output,
		Render: func(*core.Prompt[string]) string { return text },
	})
	_, err = p.Run()
	require.NoError(t, err)
	assert.True(t, strings.Contains(p.Frame, "HEADER"))
}

func TestPersistentComponentsAcrossForms(t *testing.T) {
	// Prompt.Run does not switch to raw mode under go test.
	core.SetPersistent(func() string { return "HEADER" })
	t.Cleanup(core.ClearPersistent)

	output, err := os.CreateTemp(t.TempDir(), "terminal")
	require.NoError(t, err)
	defer output.Close()

	runPersistentTestPrompt(t, output, "Form one")
	runPersistentTestPrompt(t, output, "Form two")
	content, err := os.ReadFile(output.Name())
	require.NoError(t, err)
	assert.Equal(t, 2, strings.Count(string(content), "HEADER"))
	assert.Contains(t, string(content), "Form two")
	// The second form erases and replaces the previous session frame.
	assert.Contains(t, string(content), "\x1b[J")

	core.ClearPersistent()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()
	defer w.Close()
	_, err = w.Write([]byte("\r"))
	require.NoError(t, err)
	p := core.NewPrompt(core.PromptParams[string]{
		Input: r, Output: output,
		Render: func(*core.Prompt[string]) string { return "No header" },
	})
	_, err = p.Run()
	require.NoError(t, err)
	assert.Equal(t, "No header", p.Frame)
}
