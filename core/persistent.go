package core

import (
	"os"
	"strings"
	"sync"
)

// PersistentComponent renders a piece of content above the current prompt.
// The functions are re-evaluated on each prompt redraw.
type PersistentComponent func() string

var persistent = struct {
	sync.Mutex
	components  []PersistentComponent
	output      *os.File
	lastFrame   string
	completed   bool
	activeRender func()
}{}

// SetPersistent replaces the components displayed above prompts. The components
// remain in place as the user advances through sequential prompts. Only one
// interactive prompt should be running at a time.
func SetPersistent(components ...PersistentComponent) {
	persistent.Lock()
	defer persistent.Unlock()
	persistent.components = append([]PersistentComponent(nil), components...)
	if len(components) == 0 {
		persistent.lastFrame = ""
		persistent.completed = false
	}
}

// ClearPersistent removes all persistent components for subsequent renders.
// It does not erase already written terminal history.
func ClearPersistent() {
	SetPersistent()
}

// RefreshPersistent re-renders the active prompt with the latest component
// content. This should be invoked on the same execution path as prompt input;
// callers must not mutate component data concurrently with prompt state.
func RefreshPersistent() {
	persistent.Lock()
	render := persistent.activeRender
	persistent.Unlock()
	if render != nil {
		render()
	}
}

func beginPersistent(output *os.File, render func()) {
	persistent.Lock()
	defer persistent.Unlock()
	persistent.output = output
	persistent.activeRender = render
}

func endPersistent(output *os.File) {
	persistent.Lock()
	defer persistent.Unlock()
	if persistent.output == output {
		persistent.activeRender = nil
		persistent.completed = len(persistent.components) > 0 && persistent.lastFrame != ""
	}
}

func persistentFrame(form string) string {
	parts := make([]string, 0, len(persistent.components)+1)
	for _, component := range persistent.components {
		if component == nil {
			continue
		}
		content := strings.TrimRight(strings.ReplaceAll(strings.ReplaceAll(component(), "\r\n", "\n"), "\n", "\r\n"), "\r\n")
		if content != "" {
			parts = append(parts, content)
		}
	}
	parts = append(parts, form)
	return strings.Join(parts, "\r\n")
}
