package prompts

import (
	"io"
	"os"
	"strings"

	"github.com/orochaa/go-clack/prompts/banner"
)

// BannerOptions configures a printable FIGlet banner.
type BannerOptions struct {
	Config banner.Config
	Output io.Writer
}

// Banner writes an ASCII-art heading to the configured output (stdout by default).
// To keep the banner visible above subsequent forms, use BannerComponent with SetPersistent.
func Banner(text string, options BannerOptions) error {
	cfg := options.Config
	cfg.Text = text
	art, err := banner.Render(cfg)
	if err != nil {
		return err
	}
	out := options.Output
	if out == nil {
		out = os.Stdout
	}
	_, err = io.WriteString(out, strings.TrimRight(art, "\r\n")+"\r\n")
	return err
}

// BannerComponent renders a banner once as a reusable persistent component.
// Rendering it once also ensures randomly selected fonts and gradients remain stable.
func BannerComponent(cfg banner.Config) (PersistentComponent, error) {
	art, err := banner.Render(cfg)
	if err != nil {
		return nil, err
	}
	return StaticComponent(art), nil
}
