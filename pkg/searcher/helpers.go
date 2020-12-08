package searcher

func sliceUnique(s []string) []string {
	keys := make(map[string]struct{})
	list := make([]string, 0)
	for _, entry := range s {
		if _, ok := keys[entry]; !ok {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
