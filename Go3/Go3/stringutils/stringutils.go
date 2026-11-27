package stringutils

func Perevorot_stroki(a string) string {
	perevertysh := []rune(a)
	for i, j := 0, len(perevertysh)-1; i < j; i, j = i+1, j-1 {
		perevertysh[i], perevertysh[j] = perevertysh[j], perevertysh[i]
	}
	return string(perevertysh)
}
