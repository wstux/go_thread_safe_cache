package list

// Index constants to represent a null pointer in our array-based list
const NullIdx int32 = -1

// Link represents the raw pointer connections inside a flat slice.
type Link struct {
	Prev int32
	Next int32
}

// List manages a flat array-based intrusive doubly linked list topology.
type List struct {
	Head int32
	Tail int32
}

// New initializes an empty list structure.
func New() List {
	return List{
		Head: NullIdx,
		Tail: NullIdx,
	}
}

// AddToHead inserts a node at the front of the doubly linked list.
// The caller passes a slice of Link structures to modify the pointers in-place.
func (l *List) AddToHead(idx int32, links []Link) {
	links[idx].Next = l.Head
	links[idx].Prev = NullIdx

	if l.Head != NullIdx {
		links[l.Head].Prev = idx
	}
	l.Head = idx
	if l.Tail == NullIdx {
		l.Tail = idx
	}
}

// RemoveNode detaches an index node from its current position in the list topology.
func (l *List) RemoveNode(idx int32, links []Link) {
	prev := links[idx].Prev
	next := links[idx].Next

	if prev != NullIdx {
		links[prev].Next = next
	} else {
		l.Head = next
	}

	if next != NullIdx {
		links[next].Prev = prev
	} else {
		l.Tail = prev
	}
}

// MoveToHead promotes an active node directly to the front position.
func (l *List) MoveToHead(idx int32, links []Link) {
	if l.Head == idx {
		return
	}
	l.RemoveNode(idx, links)
	l.AddToHead(idx, links)
}
