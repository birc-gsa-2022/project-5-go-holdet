package shared

/*radix sort, the offset describes how far we want to radix sort.
a complete radix sort can be done by passing length of word-1 as offset param
in skew we use 2 for S12 and 0 for S3*/
func radixS(x word, alphaSize int, xInd []int, offsetIdx int) []int {
	//radix sort left to right and finish with offset 0
	for offsetIdx >= 0 {
		xInd = bucketSortS(x, alphaSize, xInd, offsetIdx)
		offsetIdx--
	}
	return xInd
}

//bucket sort where we take the use the offset to specify what 'column' to sort
func bucketSortS(x word, alphabetSize int, xInd []int, offsetIdx int) []int {
	xSub := make([]int, len(xInd))
	for i, v := range xInd {
		xSub[i] = valOrSentinel(x, v+offsetIdx)
	}
	//fill buckets
	buckets := countCharS(xSub, alphabetSize)
	//accumulate buckets
	cbuckets := countAccS(buckets)

	//insert all letters in 'sorted' order and increment 'bucket pointer' each iteration
	xSorted := make([]int, len(xInd))
	for _, v := range xInd {
		letter := valOrSentinel(x, v+offsetIdx)
		xSorted[cbuckets[letter]] = v
		cbuckets[letter]++
	}
	return xSorted
}

//create and fill buckets
func countCharS(x_sub []int, alphabetSize int) []int {
	buckets := make([]int, alphabetSize)
	for _, v := range x_sub {
		buckets[v] += 1
	}
	return buckets
}

//accumulate buckets
func countAccS(counts []int) []int {
	csum := 0
	cbuckets := make([]int, len(counts))
	for i, v := range counts {
		cbuckets[i] = csum
		csum += v
	}
	return cbuckets
}

//if index is in bounds we return val,  otherwise return 0
func valOrSentinel(x word, idx int) int {
	if len(x) > idx {
		return x[idx]
	} else {
		return 0
	}
}
