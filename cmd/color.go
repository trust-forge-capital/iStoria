package cmd

import "fmt"

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

// Colorize returns string with color if not noColor mode
func Colorize(s string, color string, noColor bool) string {
	if noColor {
		return s
	}
	return color + s + ColorReset
}

// FmtField formats a key:value pair with optional color
func FmtField(key, val string, noColor bool) string {
	if noColor {
		return fmt.Sprintf("%s:%s", key, val)
	}
	return fmt.Sprintf("%s:%s%s%s", Colorize(key, ColorCyan, noColor), Colorize(val, ColorReset, noColor), ColorReset, ColorReset)
}

// FmtValue formats a value with optional color
func FmtValue(val string, color string, noColor bool) string {
	if noColor {
		return val
	}
	return color + val + ColorReset
}

// FmtPercent formats percentage with color based on threshold
func FmtPercent(val float64, noColor bool) string {
	s := fmt.Sprintf("%.1f%%", val)
	if noColor {
		return s
	}
	var color string
	if val >= 90 {
		color = ColorRed
	} else if val >= 70 {
		color = ColorYellow
	} else {
		color = ColorGreen
	}
	return Colorize(s, color, noColor)
}

// FmtBar creates a simple ASCII progress bar
func FmtBar(percent float64, width int, noColor bool) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(float64(width) * percent / 100)
	empty := width - filled
	
	if noColor {
		bar := ""
		for i := 0; i < filled; i++ {
			bar += "="
		}
		for i := 0; i < empty; i++ {
			bar += "-"
		}
		return bar
	}
	
	// Colored bar
	bar := ""
	if filled > 0 {
		quarter := filled / 4
		remainder := filled % 4
		// Full blocks
		for i := 0; i < quarter; i++ {
			bar += Colorize("█", ColorGreen, false)
		}
		// Partial block
		if remainder > 0 {
			if quarter < 2 {
				bar += Colorize("▎", ColorGreen, false)
			} else {
				bar += Colorize("▎", ColorYellow, false)
			}
		}
	}
	// Empty blocks
	for i := 0; i < (empty - filled/4 - 1); i++ {
		bar += "░"
	}
	return bar
}
