package main

import (
	"fmt"

	"github.com/orochaa/go-clack/prompts"
	"github.com/orochaa/go-clack/prompts/banner"
)

// ReorderBannerExample demonstrates a persistent banner and an ordered form.
func ReorderBannerExample() {
	header, err := prompts.BannerComponent(banner.Config{
		Text:     "CLACK",
		Font:     "standard",
		Gradient: &banner.GradientOcean,
	})
	HandleCancel(err)
	prompts.SetPersistent(header, prompts.StaticComponent("Deployment setup"))
	defer prompts.ClearPersistent()

	order, err := prompts.Reorder(prompts.ReorderParams[string]{
		Message: "Choose a deployment order:",
		Options: []*prompts.ReorderOption[string]{
			{Label: "Install", Value: "install"},
			{Label: "Build", Value: "build"},
			{Label: "Deploy", Value: "deploy"},
		},
	})
	HandleCancel(err)

	_, err = prompts.Confirm(prompts.ConfirmParams{
		Message: fmt.Sprintf("Use %v?", order),
	})
	HandleCancel(err)
}
