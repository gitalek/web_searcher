package searcher

func sliceUnique(s []string) []string {
	keys := make(map[string]struct{})
	var list []string
	for _, entry := range s {
		if _, ok := keys[entry]; !ok {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
