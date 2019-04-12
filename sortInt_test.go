package kmgSortInt

import (
	"github.com/bronze1man/kmgTest"
	"sort"
	"testing"
)

func TestInts(t *testing.T) {
	a := []int{4, 3, 1, 2, 5, 6}
	SortInt(a)
	for i, v := range a {
		kmgTest.Ok(v == i+1)
	}
	kmgTest.BenchmarkWithRepeatNum(1<<10, func() {
		SortInt(a)
	})
	kmgTest.BenchmarkWithRepeatNum(1<<10, func() {
		sort.Ints(a)
	})

	size:=1024
	a=make([]int,size)
	for i:=0;i<size;i++{
		a[i] = size-i
	}
	kmgTest.Benchmark(func(){
		kmgTest.BenchmarkSetNum(size)
		SortInt(a)
	})
	kmgTest.Benchmark(func(){
		kmgTest.BenchmarkSetNum(size)
		sort.Ints(a)
	})
}
