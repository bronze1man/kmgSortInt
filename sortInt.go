package kmgSortInt

// SortInt sorts a slice of ints in increasing order.
// it guarantee no alloc will happen in this function.
func SortInt(data []int) {
	quickSort(data, 0, len(data), maxDepth(len(data)))
}

func quickSort(data []int, a, b, maxDepth int) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			//data.Less(i, i-6)
			if _less(data, i, i-6) {
				//data.Swap(i, i-6)
				_swap(data, i, i-6)
			}
		}
		insertionSort(data, a, b)
	}
}

func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func insertionSort(data []int, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && _less(data, j, j-1); j-- {
			_swap(data, j, j-1)
		}
	}
}

// medianOfThree moves the median of the three values data[m0], data[m1], data[m2] into data[m1].
func medianOfThree(data []int, m1, m0, m2 int) {
	// sort 3 elements
	if _less(data, m1, m0) {
		_swap(data, m1, m0)
	}
	// data[m0] <= data[m1]
	if _less(data, m2, m1) {
		_swap(data, m2, m1)
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if _less(data, m1, m0) {
			_swap(data, m1, m0)
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func doPivot(data []int, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(data, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && _less(data, a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !_less(data, pivot, b); b++ { // data[b] <= pivot
		}
		for ; b < c && _less(data, pivot, c-1); c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		_swap(data, b, c-1)
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !_less(data, pivot, hi-1) { // data[hi-1] = pivot
			_swap(data, c, hi-1)
			c++
			dups++
		}
		if !_less(data, b-1, pivot) { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !_less(data, m, pivot) { // data[m] = pivot
			_swap(data, m, b-1)
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && !_less(data, b-1, pivot); b-- { // data[b] == pivot
			}
			for ; a < b && _less(data, a, pivot); a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			_swap(data, a, b-1)
			a++
			b--
		}
	}
	// Swap pivot into middle
	_swap(data, pivot, b-1)
	return b - 1, c
}

func heapSort(data []int, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		//data.Swap(first, first+i)
		_swap(data, first, first+i)
		siftDown(data, lo, i, first)
	}
}

// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown(data []int, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		//data.Less(first+child, first+child+1)
		if child+1 < hi && _less(data, first+child, first+child+1) {
			child++
		}
		//data.Less(first+root, first+child)
		if !_less(data, first+root, first+child) {
			return
		}
		//data.Swap(first+root, first+child)
		_swap(data, first+root, first+child)
	}
}

func _less(data []int, a int, b int) bool {
	return data[a] < data[b]
}

func _swap(data []int, a int, b int) {
	data[a], data[b] = data[b], data[a]
}
