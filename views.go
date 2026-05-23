package main

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
	zone "github.com/lrstanley/bubblezone"
)

var (
	colAccent = lipgloss.Color("#C9A227") // warm gold
	colText   = lipgloss.Color("#E6E6E6")
	colDim    = lipgloss.Color("#6C6C6C")
	colStar   = lipgloss.Color("#C9A227")

	portraitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9A9A9A"))
	bannerStyle   = lipgloss.NewStyle().Foreground(colText).Bold(true)
	bioStyle      = lipgloss.NewStyle().Foreground(colText)
	dimStyle      = lipgloss.NewStyle().Foreground(colDim)

	footerActive   = lipgloss.NewStyle().Foreground(colAccent).Bold(true)
	footerInactive = lipgloss.NewStyle().Foreground(colDim)

	ruleStyle     = lipgloss.NewStyle().Foreground(colDim)
	categoryStyle = lipgloss.NewStyle().Foreground(colDim)
	itemStyle     = lipgloss.NewStyle().Foreground(colText)
	itemSelected  = lipgloss.NewStyle().Foreground(colAccent).Bold(true)
	starStyle     = lipgloss.NewStyle().Foreground(colStar)
)

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var body string
	if m.page == pageHome {
		body = m.viewHome()
	} else {
		body = m.viewSection()
	}

	content := lipgloss.JoinVertical(lipgloss.Left, body, "", m.viewFooter())
	placed := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	return zone.Scan(placed)
}

// portraitScale shrinks the art by keeping every Nth column and row. The
// source is 100×68 chars; scale 2 → 50×34, scale 3 → ~34×23. Bump it up to
// make the portrait smaller, down toward 1 for full resolution.
const portraitScale = 1

// portraitBlock is the art downsampled, then padded to a perfect rectangle
// (every line the same display width), computed once so the right-hand column
// always lines up.
var portraitBlock = rectangle(downsample(portrait, portraitScale))

// downsample shrinks ASCII art by keeping every nth column and row.
func downsample(s string, n int) string {
	if n <= 1 {
		return s
	}
	lines := strings.Split(strings.Trim(s, "\n"), "\n")
	var out []string
	for y := 0; y < len(lines); y += n {
		runes := []rune(lines[y])
		var b strings.Builder
		for x := 0; x < len(runes); x += n {
			b.WriteRune(runes[x])
		}
		out = append(out, b.String())
	}
	return strings.Join(out, "\n")
}

func rectangle(s string) string {
	lines := strings.Split(strings.Trim(s, "\n"), "\n")
	w := 0
	for _, l := range lines {
		if x := lipgloss.Width(l); x > w {
			w = x
		}
	}
	for i, l := range lines {
		if pad := w - lipgloss.Width(l); pad > 0 {
			lines[i] = l + strings.Repeat(" ", pad)
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) viewHome() string {
	left := portraitStyle.Render(portraitBlock)

	var right strings.Builder
	right.WriteString(shimmerBanner(m.frame))
	right.WriteString("\n\n")
	for _, line := range bio {
		right.WriteString(bioStyle.Render(line))
		right.WriteString("\n")
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, left, "    ", right.String())
}

func (m model) viewSection() string {
	var b strings.Builder
	entry := nav[m.page]
	title := entry.label
	b.WriteString(footerActive.Render(title))
	b.WriteString("\n")
	b.WriteString(ruleStyle.Render(strings.Repeat("─", lipgloss.Width(title)+4)))
	b.WriteString("\n\n")

	idx := 0
	render := func(it item) {
		cursor := "  "
		style := itemStyle
		if idx == m.cursor {
			cursor = "› "
			style = itemSelected
		}
		star := "  "
		if it.star {
			star = starStyle.Render("✦ ")
		}
		line := cursor + star + style.Render(it.title)
		b.WriteString(zone.Mark(rowZoneID(idx), line))
		b.WriteString("\n")
		idx++
	}

	if entry.contacts {
		for _, c := range contacts {
			render(item{title: fmt.Sprintf("%-4s %s", c.label, c.handle), url: c.url})
			b.WriteString("\n")
		}
	} else {
		for _, g := range entry.section.groups {
			b.WriteString(categoryStyle.Render(g.category))
			b.WriteString("\n")
			for _, it := range g.items {
				render(it)
			}
			b.WriteString("\n")
		}
	}

	if m.opened != "" {
		b.WriteString(starStyle.Render("→ ") + hyperlink(m.opened, m.opened))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("[↑ ↓ to select · enter to open · esc back]"))
	return b.String()
}

func (m model) viewFooter() string {
	parts := make([]string, len(nav))
	for i, entry := range nav {
		text := entry.label
		style := footerInactive
		if (m.page == pageHome && i == m.footerIdx) || m.page == i {
			text = "✦ " + entry.label
			style = footerActive
		}
		parts[i] = zone.Mark(footerZoneID(i), style.Render(text))
	}
	return strings.Join(parts, dimStyle.Render("    "))
}

// shimmerBanner renders the figlet name with a bright highlight band that
// sweeps left-to-right across the letters, advanced by the frame counter.
func shimmerBanner(frame int) string {
	base, _ := colorful.Hex("#8A6D1F")    // deep gold
	highlight, _ := colorful.Hex("#FFF1B8") // pale gold

	lines := strings.Split(nameBanner, "\n")
	width := 0
	for _, l := range lines {
		if w := utf8.RuneCountInString(l); w > width {
			width = w
		}
	}

	const band = 7.0 // width of the highlight falloff, in columns
	period := float64(width) + band*3
	head := math.Mod(float64(frame)*0.9, period) - band

	var b strings.Builder
	for _, line := range lines {
		for i, r := range []rune(line) {
			if r == ' ' {
				b.WriteRune(' ')
				continue
			}
			t := math.Max(0, 1-math.Abs(float64(i)-head)/band)
			c := base.BlendLab(highlight, t).Clamped()
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color(c.Hex())).
				Bold(true).
				Render(string(r)))
		}
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

// hyperlink wraps text in an OSC 8 escape so terminals that support it render
// a clickable link.
func hyperlink(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}
