package wputil

func Trim(l int, s string) string {
	if len(s) <= l {
		return s
	}
	cut := len(s) - l
	return "..." + s[cut+3:]
}
