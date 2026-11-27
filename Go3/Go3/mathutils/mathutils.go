package mathutils

func Faktorial_vychisli_mne_bistro(a int) int {
	b := 1
	for i := 1; i <= a; i++ {
		b = b * i
	}
	return b
}
