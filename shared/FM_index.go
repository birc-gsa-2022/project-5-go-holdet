package shared

import "fmt"

func getSortedKeysOfCountSlice(counts map[byte]int) map[byte]int {
	keys := make([]int, 256)
	C := make(map[byte]int)
	for i, k := range counts {
		keys[i] += k
	}
	tot := 0
	for i, v := range keys {
		if v != 0 {

			C[byte(i)] = tot
			tot += v
		}
	}
	return C
}
func BuildOtable(bwt []byte) []map[byte]int {
	o := make([]map[byte]int, len(bwt)+1)
	counts := make(map[byte]int)
	copyOfCounts := make(map[byte]int)

	o[0] = copyOfCounts

	for i, v := range bwt {
		copyOfCounts := make(map[byte]int)

		counts[v] += 1
		for key, value := range counts {
			copyOfCounts[key] = value
		}

		o[i+1] = copyOfCounts
	}
	return o
}

// Data might need to represented differently
func FM_build(sa []int, genome string) ([]byte, map[byte]int) {
	bwt := make([]byte, len(sa))
	counts := make(map[byte]int)
	activeSymbol := genome[len(genome)-1]

	for i, v := range sa {
		copyOfCounts := make(map[byte]int)
		// Copy from the original map to the target map
		for key, value := range counts {
			copyOfCounts[key] = value
		}

		//add current letter to o table
		if v == 0 {
			bwt[i] = genome[len(sa)-1]
		} else {
			bwt[i] = genome[v-1]
		}
		counts[bwt[i]] += 1

		if activeSymbol != genome[v] {
			activeSymbol = genome[v]
		}
	}

	//create buckets
	C := getSortedKeysOfCountSlice(counts)

	return bwt, C
}

//locate interval for pattern p
func FM_search(bwt []byte, C map[byte]int, O []map[byte]int, p string) (int, int) {
	L := 0
	R := len(bwt)

	for i := len(p) - 1; i >= 0; i-- {
		if L == R {
			return L, R
		}

		a := p[i]

		L = C[a] + O[L][a]
		R = C[a] + O[R][a]
	}
	return L, R
}

func FM_search_approx(bwt []byte, C map[byte]int, O []map[byte]int, p string, k int) {

	for k, v := range C {
		fmt.Println(k, v)
	}

}

func ReverseBWT(bwt []byte, C map[byte]int, O []map[byte]int) []int {
	//remember O is the same as rank

	rev := make([]int, len(bwt))

	st := -1
	//find sentinel
	for i, v := range bwt {
		if v == '$' {
			st = i
			break
		}
	}

	bwt_idx := st
	bar_idx := 0
	//reversing transformation
	for rot := len(bwt); rot >= 0; rot-- {
		letter := bwt[bwt_idx]
		bar_idx = C[letter] + O[bwt_idx][letter]
		rev[bwt_idx] = rot
		bwt_idx = bar_idx
	}
	return rev
}
