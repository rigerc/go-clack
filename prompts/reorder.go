package prompts

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/orochaa/go-clack/core"
	"github.com/orochaa/go-clack/core/validator"
	"github.com/orochaa/go-clack/prompts/symbols"
	"github.com/orochaa/go-clack/prompts/test"
	"github.com/orochaa/go-clack/prompts/theme"
	"github.com/orochaa/go-clack/third_party/picocolors"
)

type ReorderOption[T comparable] struct {
	Label string
	Value T
	Hint  string
}

type ReorderParams[T comparable] struct {
	Context context.Context
	Input   *os.File
	Output  *os.File
	Message string
	Options []*ReorderOption[T]
}

// Reorder displays an interactive ordered list. Use arrows to select an item,
// Space to pick it up or drop it, and arrows to move it while picked up.
// Enter submits all item values in their final order; Ctrl+C cancels.
func Reorder[T comparable](params ReorderParams[T]) ([]T, error) {
	v := validator.NewValidator("Reorder")
	v.ValidateOptions(len(params.Options))
	options := make([]*core.ReorderOption[T], len(params.Options))
	for i, option := range params.Options {
		if option == nil {
			v.Panic("Options must not contain nil items")
		}
		options[i] = &core.ReorderOption[T]{
			Label: option.Label,
			Value: option.Value,
			Hint:  option.Hint,
		}
	}

	p := core.NewReorderPrompt(core.ReorderPromptParams[T]{
		Context: params.Context,
		Input:   params.Input,
		Output:  params.Output,
		Options: options,
		Render: func(p *core.ReorderPrompt[T]) string {
			var value string
			switch p.State {
			case core.SubmitState, core.CancelState, core.ValidateState:
				labels := make([]string, len(p.Options))
				for i, option := range p.Options {
					labels[i] = option.Label
				}
				value = strings.Join(labels, " → ")
			default:
				lines := make([]string, len(p.Options))
				for i, option := range p.Options {
					marker := picocolors.Dim(symbols.RADIO_INACTIVE)
					label := picocolors.Dim(option.Label)
					if i == p.CursorIndex {
						marker = picocolors.Green(symbols.RADIO_ACTIVE)
						if p.Picked {
							marker = picocolors.Yellow("↕")
						}
						label = option.Label
					}
					hint := ""
					if option.Hint != "" && i == p.CursorIndex {
						hint = " " + picocolors.Dim(fmt.Sprintf("(%s)", option.Hint))
					}
					lines[i] = fmt.Sprintf("%s %s%s", marker, label, hint)
				}
				instructions := picocolors.Dim("Space: pick up/drop • ↑/↓: navigate/move • Enter: confirm")
				value = p.LimitLines(lines, 4) + "\r\n" + instructions
			}

			return theme.ApplyTheme(theme.ThemeParams[[]T]{
				Context:         p.Prompt,
				Message:         params.Message,
				Value:           value,
				ValueWithCursor: value,
			})
		},
	})
	test.ReorderTestingPrompt = p
	return p.Run()
}
