package test

func Contains(set []string, element string) bool {
	for _, v := range set {
		if element == v {
			return true
		}
	}
	return false
}
