package shared

func CreateLMS(x string) []rune {
	len := len(x)
	lms := make([]rune, len)

	// sentinel is always S
	lms[len-1] = 'S'

	for i := len - 2; i >= 0; i-- {
		if x[i] < x[i+1] {
			lms[i] = 'S'
		}
		if x[i] > x[i+1] {
			lms[i] = 'L'
		}
		if x[i] == x[i+1] {
			lms[i] = lms[i+1]
		}
	}
	return lms
}

func CreateLMSSuffix(x string, lmsarr []rune) []int {
	lmssuffixes := make([]int, len(x))
	for i := 1; i < len(x); i++ {
		if lmsarr[i] == 'S' && lmsarr[i]-1 == 'L' {
			lmssuffixes = append(lmssuffixes, i)
		}
	}
	return lmssuffixes
}

func RecursiveSais(x string, sa []int, alphabetsize int) {
	if len(x) == alphabetsize {
		return
	}

	lms := CreateLMS(x)
	lmssucc := CreateLMSSuffix(x, lms)

	buckets := CountSort(x)

}

func Sais(x string) []int {
	sa := make([]int, len(x))

	RecursiveSais(x, sa, 0)

	return sa
}
