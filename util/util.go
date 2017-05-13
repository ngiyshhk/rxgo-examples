package util

func distinct(t []string) map[string]bool {
	result := map[string]bool{}
	for _, v := range t {
		if !result[v] {
			result[v] = true
		}
	}
	return result
}

func Diff(before []string, after []string) []string {
	b := distinct(before)
	a := distinct(after)
	result := make([]string, 0)
	for k := range a {
		if !b[k] {
			result = append(result, k)
		}
	}
	return result
}
