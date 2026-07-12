package tui

import (
	"charm.land/glamour/v2/ansi"
)

const (
	logo = "\033[38;2;180;190;254m  ██████╗   ██████╗ ████████╗  ██████╗\033[0m\n" +
		"\033[38;2;137;180;250m ██╔═══██╗ ██╔════╝ ╚══██╔══╝ ██╔═══██╗\033[0m\n" +
		"\033[38;2;116;199;236m ██║   ██║ ██║         ██║    ██║   ██║\033[0m\n" +
		"\033[38;2;148;226;213m ██║   ██║ ██║         ██║    ██║   ██║\033[0m\n" +
		"\033[38;2;203;166;247m ╚██████╔╝ ╚██████╗    ██║    ╚██████╔╝\033[0m\n" +
		"\033[38;2;245;194;231m  ╚═════╝   ╚═════╝    ╚═╝     ╚═════╝\033[0m\n"

	colorUser      = "\033[1;111m" // bold blue (Catppuccin Blue)
	colorAssistant = "\033[1;150m" // bold green (Catppuccin Green)
	colorThinking  = "\033[2;245m" // dim overlay1 (Catppuccin Overlay1)
	colorReset     = "\033[0m"

	headerHeight      = 8
	footerHeight      = 2
	maxViewportHeight = 30
	minViewportHeight = 3
)

func boolPtr(b bool) *bool { return &b }
func uintPtr(u uint) *uint { return &u }
func strPtr(s string) *string { return &s }

// Catppuccin Mocha 256-color mappings:
//
//	Text:252  Subtext0:248  Overlay1:243  Overlay0:242
//	Surface0:236  Base:234
//	Lavender:147  Blue:111  Sapphire:117  Teal:122
//	Green:150  Yellow:228  Peach:216  Red:203
//	Mauve:141  Pink:218  Flamingo:224
func customStyle() ansi.StyleConfig {
	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: strPtr("252"),
			},
			Margin: uintPtr(0),
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Bold:        boolPtr(true),
				Color:       strPtr("147"),
				BlockSuffix: "\n",
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Bold:        boolPtr(true),
				Color:       strPtr("228"),
				Prefix:      "# ",
				BlockSuffix: "\n",
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Bold:        boolPtr(true),
				Color:       strPtr("147"),
				Prefix:      "## ",
				BlockSuffix: "\n",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Bold:        boolPtr(true),
				Color:       strPtr("111"),
				Prefix:      "### ",
				BlockSuffix: "\n",
			},
		},
		Strong: ansi.StylePrimitive{
			Bold: boolPtr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: boolPtr(true),
		},
		Item: ansi.StylePrimitive{
			Prefix: "  • ",
		},
		Enumeration: ansi.StylePrimitive{
			Prefix: "  ",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:           strPtr("203"),
				BackgroundColor: strPtr("236"),
				Prefix:          " ",
				Suffix:          " ",
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: strPtr("245"),
				},
				Margin: uintPtr(1),
			},
		},
		Link: ansi.StylePrimitive{
			Color:     strPtr("117"),
			Underline: boolPtr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: strPtr("141"),
			Bold:  boolPtr(true),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:  strPtr("245"),
				Prefix: "  │ ",
			},
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  strPtr("240"),
			Prefix: "\n",
			Suffix: "\n",
		},
	}
}
