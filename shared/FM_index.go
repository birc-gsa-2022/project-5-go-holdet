package shared

import (
	"fmt"
)

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

/*create reverse bwt array which we use for R0.
This could probably be done in preprocessing*/
func BuildROtable(bwt []byte) []map[byte]int {
	fmt.Println(bwt)
	rbwt := make([]byte, len(bwt))
	for i, j := 0, len(rbwt)-1; j >= 0; i, j = i+1, j-1 {
		rbwt[i] = bwt[j]
	}
	fmt.Println(rbwt, bwt)
	return BuildOtable(rbwt)
}

/*Used to terminate early in Li-Durbin
p. 263-264 in book*/
func BuildDTable(p string, rec FMRecs) []int {

	D := make([]int, len(p))

	min_edits := 0
	L := 0
	R := len(rec.Bwt)

	for i, v := range []byte(p) {
		a := v

		L = rec.C[a] + rec.RO[L][a]
		R = rec.C[a] + rec.RO[R][a]

		if L >= R {
			min_edits++
			L = 0
			R = len(rec.Bwt)
		}
		D[i] = min_edits

	}
	return D
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

//simple function to initiate variables etc for the recursive search
func FM_search_approx(rec FMRecs, read Recs, edits int) {
	p := read.Rec
	d_global = BuildDTable(p, rec)

	L, R := 0, len(rec.Bwt)
	i := len(p) - 1

	x_global = rec
	p_global = read
	//initiate recursive search
	RecApproxMatching(L, R, i, edits, rec, p, []rune{})

}

//right now we use a lot of global variables in an attempt to limit amount of variables getting passed around
var x_global FMRecs
var p_global Recs
var d_global []int

func RecApproxMatching(L int, R int, idx int, edits int, rec FMRecs, p string, cigar []rune) {

	/*L, R interval contains matches
	this also prevents deletions in front of match */
	if idx == -1 {
		if edits >= 0 {
			matchFound(L, R, cigar, rec.BS)
		}
		return
	}

	//we are out of available edits
	if edits < d_global[idx] {
		return
	}

	//take I step
	RecApproxMatching(L, R, idx-1, edits-1, rec, p, append(cigar, 'I'))

	//iterate over alphabet ($ excluded)
	for a := range rec.C {
		if a == '$' {
			continue
		}
		//decide if this letter is considered an edit
		cost := 1
		if p[idx] == a {
			cost = 0
		}
		//no edits available
		if edits-cost < 0 {
			continue
		}

		//do a single FM step
		newL := rec.C[a] + rec.O[L][a]
		newR := rec.C[a] + rec.O[R][a]

		//no interval to consider
		if newL >= newR {
			continue
		}

		//take M step
		RecApproxMatching(newL, newR, idx-1, edits-cost, rec, p, append(cigar, 'M'))

		/*take D step
		recursive so we do not allow first iteration (last 'cigar letter') to be a D*/
		if len(cigar) > 0 {
			RecApproxMatching(newL, newR, idx, edits-1, rec, p, append(cigar, 'D'))
		}

	}
}

/*Should return the sam format but we need somehow to
pass p */
func matchFound(L int, R int, cigar []rune, sa []int) {
	for i := L; i < R; i++ {
		SamMID(p_global.Name, x_global.Name, sa[i], p_global.Rec, GetCompactCigar(cigar))
	}

}
