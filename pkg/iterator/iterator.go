package iterator

import (
	"iter"
	"sort"
)

type Iterator[V any] struct {
	iter iter.Seq[V]
}

func From[V any](slice []V) *Iterator[V] {
	// If slice is nil, return an iterator with empty function to avoid nil checks later
	if slice == nil {
		return &Iterator[V]{
			iter: func(yield func(V) bool) {},
		}
	}

	return &Iterator[V]{
		iter: func(yield func(V) bool) {
			// Use direct indexing for better performance with large slices
			for i := 0; i < len(slice); i++ {
				if !yield(slice[i]) {
					return
				}
			}
		},
	}
}

func (i *Iterator[V]) Collect() []V {
	// Pre-allocate with initial capacity to reduce reallocations
	collect := make([]V, 0, 16)

	// Use a closure to manage batch collection
	i.iter(func(v V) bool {
		collect = append(collect, v)
		return true
	})

	// Trim excess capacity if significantly over-allocated
	if cap(collect) > 2*len(collect) && len(collect) > 100 {
		trimmed := make([]V, len(collect))
		copy(trimmed, collect)
		return trimmed
	}

	return collect
}

func (i *Iterator[V]) Each(f func(V)) *Iterator[V] {
	// Use direct function call with closure for better performance
	i.iter(func(v V) bool {
		f(v)
		return true
	})

	// Return the iterator for method chaining
	return i
}

func (i *Iterator[V]) Reverse() *Iterator[V] {
	// First collect all elements
	collect := i.Collect()
	length := len(collect)

	// For small slices, in-place reversal is more efficient
	if length <= 1 {
		return From(collect)
	}

	// In-place reversal to avoid extra allocation
	for left, right := 0, length-1; left < right; left, right = left+1, right-1 {
		collect[left], collect[right] = collect[right], collect[left]
	}

	return From(collect)
}

func (i *Iterator[V]) Map(f func(V) V) *Iterator[V] {
	original := i.iter
	i.iter = func(yield func(V) bool) {
		original(func(v V) bool {
			return yield(f(v))
		})
	}

	// Return the iterator for method chaining
	return i
}

func (i *Iterator[V]) Filter(f func(V) bool) *Iterator[V] {
	original := i.iter
	i.iter = func(yield func(V) bool) {
		original(func(v V) bool {
			if f(v) {
				return yield(v)
			}
			return true
		})
	}

	// Return the iterator for method chaining
	return i
}

func (i *Iterator[V]) Take(n int) *Iterator[V] {
	// Handle invalid n
	if n <= 0 {
		i.iter = func(yield func(V) bool) {}
		return i
	}

	original := i.iter
	i.iter = func(yield func(V) bool) {
		count := 0
		original(func(v V) bool {
			if count >= n {
				return false
			}
			if !yield(v) {
				return false
			}
			count++
			return true
		})
	}

	// Return the iterator for method chaining
	return i
}

func (i *Iterator[V]) Skip(n int) *Iterator[V] {
	// Handle invalid n
	if n <= 0 {
		return i
	}

	// Use direct function call with closure for better performance
	original := i.iter
	i.iter = func(yield func(V) bool) {
		count := 0
		original(func(v V) bool {
			if count < n {
				count++
				return true
			}
			return yield(v)
		})
	}

	// Return the iterator for method chaining
	return i
}

func (i *Iterator[V]) Count() int {
	count := 0
	i.iter(func(v V) bool {
		count++
		return true
	})
	return count
}

func (i *Iterator[V]) First(f func(V) bool) (V, bool) {
	var result V
	found := false

	i.iter(func(v V) bool {
		if f(v) {
			result = v
			found = true
			return false // Stop iteration early
		}
		return true
	})

	return result, found
}

func (i *Iterator[V]) Last(f func(V) bool) (V, bool) {
	var (
		lastVal V
		found   bool
	)

	i.iter(func(v V) bool {
		if f(v) {
			lastVal = v
			found = true
		}
		return true // Continue iteration to find last match
	})

	return lastVal, found
}

func (i *Iterator[V]) Reduce(initial V, f func(V, V) V) V {
	result := initial
	i.iter(func(v V) bool {
		result = f(result, v)
		return true
	})
	return result
}

func (i *Iterator[V]) Find(f func(V) bool) (V, bool) {
	var result V
	found := false
	i.iter(func(v V) bool {
		if f(v) {
			result = v
			found = true
			return false // Stop iteration early
		}
		return true
	})

	return result, found
}

func (i *Iterator[V]) FindIndex(f func(V) bool) int {
	index := 0
	found := false
	i.iter(func(v V) bool {
		if f(v) {
			found = true
			return false // Stop iteration early
		}
		index++
		return true
	})

	if !found {
		return -1
	}

	return index
}

func (i *Iterator[V]) Some(f func(V) bool) bool {
	found := false
	i.iter(func(v V) bool {
		if f(v) {
			found = true
			return false // Stop iteration early
		}
		return true
	})
	return found
}

func (i *Iterator[V]) Every(f func(V) bool) bool {
	found := true
	i.iter(func(v V) bool {
		if !f(v) {
			found = false
			return false // Stop iteration early
		}
		return true
	})
	return found
}

func (i *Iterator[V]) Includes(element V, equals func(V, V) bool) bool {
	found := false
	i.iter(func(v V) bool {
		if equals(v, element) {
			found = true
			return false // Stop iteration early
		}
		return true
	})
	return found
}

func (i *Iterator[V]) IndexOf(element V, equals func(V, V) bool) int {
	index := 0
	found := false
	i.iter(func(v V) bool {
		if equals(v, element) {
			found = true
			return false // Stop iteration early
		}
		index++
		return true
	})
	if !found {
		return -1
	}
	return index
}

func (i *Iterator[V]) LastIndexOf(element V, equals func(V, V) bool) int {
	lastIndex := -1
	index := 0
	i.iter(func(v V) bool {
		if equals(v, element) {
			lastIndex = index
		}
		index++
		return true
	})
	return lastIndex
}

func (i *Iterator[V]) Slice(start, end int) *Iterator[V] {
	if start < 0 {
		start = 0
	}
	// Return an empty iterator if start > end
	if start > end {
		return &Iterator[V]{iter: func(func(V) bool) {}}
	}

	original := i.iter
	i.iter = func(yield func(V) bool) {
		original(func(v V) bool {
			if start <= 0 && end > 0 {
				if !yield(v) {
					return false
				}
				end--
			}
			start--
			return true
		})
	}
	return i
}

func (i *Iterator[V]) Unique(equals func(V, V) bool) *Iterator[V] {
	cpy := i.iter
	i.iter = func(yield func(V) bool) {
		seen := make([]V, 0, 8)
		cpy(func(v V) bool {
			isUnique := true
			for j := len(seen) - 1; j >= 0; j-- {
				if equals(seen[j], v) {
					isUnique = false
					break
				}
			}

			if !isUnique {
				return true
			}

			seen = append(seen, v)

			return yield(v)
		})
	}
	return i
}

func (i *Iterator[V]) Concat(others ...*Iterator[V]) *Iterator[V] {
	// Create a new iterator that will concatenate all iterators
	return &Iterator[V]{
		iter: func(yield func(V) bool) {
			// Yield all values from the current iterator
			i.iter(yield)

			// Yield all values from other iterators
			for _, other := range others {
				other.iter(yield)
			}
		},
	}
}

func (i *Iterator[V]) Sort(less func(V, V) bool) *Iterator[V] {
	elements := i.Collect()

	// Use a standard sorting algorithm (like Go's sort.Slice)
	sort.Slice(elements, func(a, b int) bool {
		return less(elements[a], elements[b])
	})

	return From(elements)
}

func Flat[V any](nested [][]V) *Iterator[V] {
	return &Iterator[V]{
		iter: func(yield func(V) bool) {
			for _, slice := range nested {
				for _, v := range slice {
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

func (i *Iterator[V]) Fill(value V) *Iterator[V] {
	cpy := i.iter
	i.iter = func(yield func(V) bool) {
		for range cpy {
			if !yield(value) {
				return
			}
		}
	}
	return i
}

func (i *Iterator[V]) FillRange(value V, start, end int) *Iterator[V] {
	if start < 0 {
		start = 0
	}
	original := i.iter
	i.iter = func(yield func(V) bool) {
		count := 0
		original(func(v V) bool {
			if count >= start && count < end {
				if !yield(value) {
					return false
				}
			} else {
				if !yield(v) {
					return false
				}
			}
			count++
			return true
		})
	}
	return i
}

// Chunk divides the Iterator into chunks of the specified size
func Chunk[V any](i *Iterator[V], size int) *Iterator[[]V] {
	if size <= 0 {
		size = 1
	}

	return &Iterator[[]V]{
		iter: func(yield func([]V) bool) {
			chunk := make([]V, 0, size)
			for v := range i.iter {
				chunk = append(chunk, v)
				if len(chunk) == size {
					yield(chunk)
					chunk = chunk[:0] // Reset chunk for the next set
				}
			}
			// Yield any remaining elements in the chunk
			if len(chunk) > 0 {
				yield(chunk)
			}
		},
	}
}
func (i *Iterator[V]) Append(value V) *Iterator[V] {
	original := i.iter

	i.iter = func(yield func(V) bool) {
		continueIteration := true
		original(func(v V) bool {
			if !yield(v) {
				continueIteration = false
				return false
			}
			return true
		})

		if continueIteration {
			yield(value)
		}
	}

	return i
}

func (i *Iterator[V]) AppendIfNotExist(value V, equals func(V, V) bool) *Iterator[V] {
	if equals == nil {
		return i
	}

	cpy := i.iter
	exists := false

	i.iter = func(yield func(V) bool) {
		cpy(func(v V) bool {
			if equals(v, value) {
				exists = true
			}
			return yield(v)
		})

		if !exists {
			yield(value)
		}
	}

	return i
}

func (i *Iterator[V]) RemoveIndex(index int) *Iterator[V] {
	if index < 0 {
		return i
	}

	cpy := i.iter
	i.iter = func(yield func(V) bool) {
		cpy(func(v V) bool {
			if index != 0 {
				index--
				return yield(v)
			}
			return true // Skip matched element but continue iteration
		})
	}

	return i
}

func (i *Iterator[V]) RemoveIf(f func(V) bool) *Iterator[V] {
	cpy := i.iter
	i.iter = func(yield func(V) bool) {
		cpy(func(v V) bool {
			if !f(v) {
				return yield(v)
			}
			return true // Skip matched element but continue iteration
		})
	}

	return i
}
