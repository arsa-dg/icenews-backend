package helper

func IsEmptyStrings(s ...string) bool {
	for _, str := range s {
		if str == "" {
			return true
		}
	}

	return false
}

func IsEqualString(s1 string, s2 string) bool {
	return s1 == s2
}

func IsInRangeString(s string, first int, last int) bool {
	return len(s) >= first && len(s) <= last
}
