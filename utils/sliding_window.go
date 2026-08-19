package utils

type SlidingWindow[T any] struct {
	Start, End, Size int
	Collection       []T
}

func NewSlidingWindow[T any](start, size int, collection []T) *SlidingWindow[T] {
	if start >= size {
		panic(start)
	}

	return &SlidingWindow[T]{
		Start:      start,
		End:        start + size - 1,
		Collection: collection,
	}
}

func (w *SlidingWindow[T]) Advance() {
	w.Start++
	w.End++
}

func (w *SlidingWindow[T]) Retreat() {
	w.Start--
	w.End--
}

func (w *SlidingWindow[T]) CanAdvance() bool {
	return w.End < len(w.Collection)-1
}

func (w *SlidingWindow[T]) CanRetreat() bool {
	return w.Start > 0
}
