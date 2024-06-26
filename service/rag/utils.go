package rag

func groupText(inputs []string, size int) [][]string {
	var groups [][]string
	for i := 0; i < len(inputs); i += size {
		end := i + size
		if end > len(inputs) {
			end = len(inputs)
		}
		groups = append(groups, inputs[i:end])
	}
	return groups
}
