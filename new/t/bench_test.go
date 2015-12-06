package gomet

import "testing"

func BenchmarkArray(b *testing.B) {
	q := Queue{}
	var result int
	for i := 0; i < b.N; i++ {
		q.Push(i)
		result = q.Pull()

	}
	b.Logf("resul: %d\n", result)
}

func BenchmarkSliceAppend(b *testing.B) {
	q := make([]int, 0, 10)
	var result int
	for i := 0; i < b.N; i++ {
		q = append(q, i)
		result = q[0]

		q = q[1:]
	}
	b.Logf("resul: %d\n", result)
}

func BenchmarkSliceWinner(b *testing.B) {
	q := make([]int, 0, 1000)
	var result int
	for i := 0; i < b.N; i++ {
		q = append(q, i)
		result = q[0]

		n := copy(q, q[1:])
		q = q[:n]
	}
	b.Logf("resul: %d\n", result)
}

func BenchmarkSliceCopy(b *testing.B) {
	var arr [100]int
	q := arr[:0]
	var result int
	for i := 0; i < b.N; i++ {
		q = append(q, i)
		result = q[0]

		n := copy(arr[:], q[1:])
		q = arr[:n]
	}
	b.Logf("resul: %d\n", result)
}

func BenchmarkChan(b *testing.B) {
	q := make(chan int, 100)
	var result int
	for i := 0; i < b.N; i++ {
		q <- i
		result = <-q

	}
	b.Logf("resul: %d\n", result)

}

type Queue struct {
	q       [100]int
	in, out int
}

func (q *Queue) Push(i int) {

	q.q[q.in] = i
	q.in++

	if q.in >= len(q.q) {
		q.in = 0
	}
}

func (q *Queue) Pull() int {
	out := q.q[q.out]
	q.out++

	if q.out >= len(q.q) {
		q.out = 0
	}
	return out
}
