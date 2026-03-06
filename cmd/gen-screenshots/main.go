// Generates standalone HTML files for each theme showing all GitHub-style
// alert types. Run from the repo root:
//
//	go run cmd/gen-screenshots/main.go
//
// Then screenshot them with a headless browser:
//
//	mkdir -p /tmp/mpls-screenshots/png
//	for f in /tmp/mpls-screenshots/html/*.html; do
//	  name=$(basename "$f" .html)
//	  google-chrome-stable --headless=new --window-size=900,800 \
//	    --screenshot="/tmp/mpls-screenshots/png/${name}.png" "file://${f}"
//	done

//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mhersson/mpls/pkg/parser"
)

const markdown = `# Alert Showcase

Some plain text before the alerts. This demonstrates how GitHub-style
alerts render in the **mpls** markdown preview with this theme.

> [!NOTE]
> Useful information that users should know, even when skimming content.

> [!TIP]
> Helpful advice for doing things better or more easily.

> [!IMPORTANT]
> Key information users need to know to achieve their goal.

> [!WARNING]
> Urgent info that needs immediate user attention to avoid problems.

> [!CAUTION]
> Advises about risks or negative outcomes of certain actions.
`

func main() {
	themesDir := "internal/previewserver/web/themes"
	stylesPath := "internal/previewserver/web/styles.css"
	outDir := "/tmp/mpls-screenshots/html"

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Read styles.css
	stylesCSS, err := os.ReadFile(stylesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading styles.css: %v\n", err)
		os.Exit(1)
	}

	// Render markdown once
	renderedHTML, _ := parser.HTML(markdown, "file:///demo.md")

	// List themes
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading themes dir: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".css") {
			continue
		}

		themeName := strings.TrimSuffix(entry.Name(), ".css")
		themeCSS, err := os.ReadFile(filepath.Join(themesDir, entry.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading theme %s: %v\n", entry.Name(), err)
			continue
		}

		html := fmt.Sprintf(`<!doctype html>
<html>
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<style>%s</style>
<style>%s</style>
</head>
<body>
<div class="preview-content" id="content">%s</div>
</body>
</html>`, string(themeCSS), string(stylesCSS), renderedHTML)

		outPath := filepath.Join(outDir, themeName+".html")
		if err := os.WriteFile(outPath, []byte(html), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", outPath, err)
			continue
		}

		fmt.Println(outPath)
	}
}
