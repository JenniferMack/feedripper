package wputil

import "fmt"

func TrimRight(l int, s string) string {
	if len(s) <= l {
		return s
	}
	return s[:l-3] + "..."
}

func TrimLeft(l int, s string) string {
	if len(s) <= l {
		return s
	}
	cut := len(s) - l
	return "..." + s[cut+3:]
}

func FileSize(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMG"[exp])
}
