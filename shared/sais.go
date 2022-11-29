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

func ResSA(x string, sa *[]int) {
	s := make([]int, len(x))
	for i := range s {
		s[i] = -1
	}
}

func CreateBuckets(x string) map[rune]int {
	xs := CountSort(x)

	buckets := make(map[rune]int) //create first bucket beggining at 0
	//create buckets with accumulated values
	for i, v := range xs {
		if i == 0 {
			buckets[v] = 0
			continue
		}
		if v != rune(xs[i-1]) {
			buckets[v] = i
		}
	}

	return buckets
}

func InductionSortL(x string, sa []int, lms []rune) {
	for i := 0; i < len(x); i++ {
		if sa[i] == 0 || sa[i] == -1 {
			continue
		}

		j := sa[i] - 1
		if lms[j] == 'L' {
			// todo
		}
	}
}

func RecursiveSais(x string, sa []int, alphabetsize int) {
	if len(x) == alphabetsize {
		return
	}

	//lms := CreateLMS(x)
	// := CreateLMSSuffix(x, lms)

	ResSA(x, &sa)
	// buckets := CountSort(x)

}

func Sais(x string) []int {
	sa := make([]int, len(x))

	RecursiveSais(x, sa, 0)

	return sa
}
