package shared

//some simple types to simplify things
type triplet = [3]int
type word = []int

/*main function that calls recursive
we assume alphabet is at most one byte*/
func Skew(x string) []int {

	x_b := []byte(x)

	x_int := make(word, len(x_b))
	for i, v := range x_b {
		x_int[i] = int(v)
	}
	return skewRec(x_int, 256)
}

/*recursive skew function */
func skewRec(x word, alpSize int) []int {

	//create S12
	S12 := createS12(x)
	//start 2 idx to the right
	S12 = radixS(x, alpSize, S12, 2)

	//recursion if duplicate letters in current alphabet
	alphabet := triplets(x, S12)

	if len(S12) > len(alphabet) {
		//create u string
		u := createU(x, alphabet)

		//get the temporary suffix array from recursion
		tsa := skewRec(u, len(alphabet))

		//now remap the recursive sorted array to S12
		S12 = getSortedS12(tsa)
	}

	//create suffix array for S3 - only most significant digit is unsorted
	S3 := S3SemiSorted(S12, x)
	//sort most significant digit in S3
	S3 = radixS(x, alpSize, S3, 0)

	//if we get to this point S12 should be in right order
	return merge(x, S12, S3, alphabet)

}

//method to create the initial S12 array
func createS12(x word) []int {
	//integer division gets floored which we want
	S12 := make([]int, (len(x)*2)/3)
	j := 0
	for i := range x {
		if (i % 3) != 0 {
			S12[j] = i
			j++
		}
	}
	return S12
}

//called when we recieve the recursively sorted tsa to sort S12
func getSortedS12(tsa []int) []int {
	//identify middle sentinel so we can exclude it
	mid := len(tsa) / 2

	//create sorted sa for S12
	S12 := make([]int, 0)
	for _, v := range tsa {
		if v < mid {
			S12 = append(S12, 1+v*3)
		} else if v > mid {
			S12 = append(S12, 2+3*(v-mid-1))
		}

	}
	return S12
}

//get inital sorting from the sorted S3 array
func S3SemiSorted(S12 []int, x word) []int {

	var S3 []int

	//do not forget $ (normal sentinel)
	if (len(x)-1)%3 == 0 {
		S3 = append(S3, len(x)-1)
	}

	for _, v := range S12 {

		//basic case - shift all by 1 since S12 is sorted
		if v%3 == 1 {
			S3 = append(S3, v-1)
		}
	}
	return S3

}

//merge the two sorted suffix arrays together
func merge(x word, SA12 []int, SA3 []int, alphabet map[triplet]int) []int {

	sa := make([]int, 0)

	//create 'ISA' p. 162
	ISA12 := make(map[int]int)
	//SA12 is sorted so we can get ISA in order aswell
	//need this to determine which is greater in constant time (max 2 recursive calls per comparison)
	for i := range SA12 {
		ISA12[SA12[i]] = i
	}

	i, j := 0, 0

	//merge the two already sorted arrays
	for i < len(SA12) && j < len(SA3) {

		if isSmaller(x, SA12[i], SA3[j], ISA12) {
			sa = append(sa, SA12[i])
			i++
		} else {
			sa = append(sa, SA3[j])
			j++
		}
	}

	//finally append whatever is left
	sa = append(sa, SA12[i:]...)
	sa = append(sa, SA3[j:]...)

	//attempt to get rid of the special sentinel #
	return sa
}

//check which of two indexes in x that are smallest
func isSmaller(x word, xi int, xj int, ISA12 map[int]int) bool {

	//find i value
	li := valOrSentinel(x, xi)
	//find j value
	lj := valOrSentinel(x, xj)

	//if we easily can determine smallest return, otherwise we proceed recursively
	if li < lj {
		return true
	}
	if li > lj {
		return false
	}

	//if both indices are in S12 we can use ISA
	if xi%3 != 0 && xj%3 != 0 {
		return ISA12[xi] < ISA12[xj]
	}

	//recursive case
	return isSmaller(x, xi+1, xj+1, ISA12)
}

//create all triplets/letters in new alphabet
func triplets(x word, S12 []int) map[triplet]int {

	triplets := make(map[triplet]int)

	for _, v := range S12 {
		tri := getTriplet(x, v)

		//might have to add extra to represent sentinels
		if _, ok := triplets[tri]; !ok {
			triplets[tri] = len(triplets)
		}

	}
	return triplets
}

//return a triplet
func getTriplet(x word, idx int) triplet {
	return triplet{valOrSentinel(x, idx), valOrSentinel(x, idx+1), valOrSentinel(x, idx+2)}
}

//create new word from triplets. Triplets encoded as ints (p. 146)
func createU(x word, triplets map[triplet]int) word {

	var u word

	for i := range x {

		//get int (letter) representing triplet and build new word
		if (i-1)%3 == 0 {
			tri := getTriplet(x, i)
			u = append(u, triplets[tri])
		}
	}
	//seperating sentinel. might need another value
	u = append(u, 0)

	for i := range x {
		//get int (letter) representing triplet and build new word
		if (i-2)%3 == 0 {
			tri := getTriplet(x, i)
			u = append(u, triplets[tri])
		}
	}

	return u
}
