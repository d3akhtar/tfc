package utils

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

type HashSet[T cmp.Ordered] struct {
	items map[T]any
}

func NewHashSet[T cmp.Ordered]() HashSet[T] {
	return HashSet[T]{items: make(map[T]any)}
}

func NewHashSetWithItems[T cmp.Ordered](items []T) HashSet[T] {
	set := HashSet[T]{items: make(map[T]any, len(items))}
	for _, item := range items {
		set.Add(item)
	}

	return set
}

func (set *HashSet[T]) Add(item T) {
	set.items[item] = struct{}{}
}

func (set *HashSet[T]) Remove(item T) {
	set.items[item] = struct{}{}
}

func (set *HashSet[T]) Contains(item T) bool {
	_, res := set.items[item]
	return res
}

func (set *HashSet[T]) Length() int {
	return len(set.items)
}

func (set *HashSet[T]) Items() iter.Seq[T] {
	items := slices.Collect(maps.Keys(set.items))
	slices.Sort(items)

	return func(yield func(T) bool) {
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}
