package gateway

func convertStringSliceToPointer(strings []string) []*string {
	var ptrSlice []*string
	for _, s := range strings {
		ptrSlice = append(ptrSlice, &s)
	}
	return ptrSlice
}
