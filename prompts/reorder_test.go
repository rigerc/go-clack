package prompts_test

import (
	"os"
	"testing"

	"github.com/orochaa/go-clack/prompts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReorderReturnsValuesInOrder(t *testing.T) {
	input, send, err := os.Pipe()
	require.NoError(t, err)
	defer input.Close()
	defer send.Close()
	_, err = send.Write([]byte("\r"))
	require.NoError(t, err)

	output, err := os.CreateTemp(t.TempDir(), "reorder")
	require.NoError(t, err)
	defer output.Close()

	result, err := prompts.Reorder(prompts.ReorderParams[string]{
		Input: input,
		Output: output,
		Message: "Choose order:",
		Options: []*prompts.ReorderOption[string]{
			{Label: "First", Value: "1"},
			{Label: "Second", Value: "2"},
			{Label: "Third", Value: "3"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, result)
}

func TestReorderRequiresOptions(t *testing.T) {
	assert.Panics(t, func() {
		prompts.Reorder(prompts.ReorderParams[string]{Message: "Empty"})
	})
}
