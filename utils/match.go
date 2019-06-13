package utils

func Match(pattern, s string) bool {
	if pattern == "" {
		return s == pattern
	}
	if pattern == "*" {
		return true
	}
	rs := make([]rune, 0, len(s))
	rpattern := make([]rune, 0, len(pattern))
	for _, r := range s {
		rs = append(rs, r)
	}
	for _, r := range pattern {
		rpattern = append(rpattern, r)
	}
	return matchRune(rs, rpattern)
}

func matchRune(rs, pattern []rune) bool {
	for len(pattern) > 0 {
		n := len(rs)
		switch pattern[0] {
		default:
			if n == 0 || rs[0] != pattern[0] {
				return false
			}
		case '?':
			if n == 0 {
				return false
			}
		case '*':
			if matchRune(rs, pattern[1:]) {
				return true
			}
			if n == 0 {
				return false
			}
			return matchRune(rs[1:], pattern)
		}
		rs = rs[1:]
		pattern = pattern[1:]
	}
	return len(rs) == 0 && len(pattern) == 0
}
