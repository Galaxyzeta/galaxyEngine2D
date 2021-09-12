package ds

import "container/list"

// IndexedList is a data structure consists of a list and an indexer which indicate
// the specific element in the list.
type IndexedList struct {
	list    *list.List
	indexer map[interface{}]*list.Element
}

// Pushback will insert the item to the end of a doubly-linked list.
// Will record it in the indexed map.
func (s IndexedList) PushBack(item interface{}) {
	s.indexer[item] = s.list.PushBack(item)
}

// Remove will look up the item from an indexed map, and then delete it from list.
// Cost is O(1).
func (s IndexedList) Remove(item interface{}) {
	s.list.Remove(s.indexer[item])
	delete(s.indexer, item)
}

func (s IndexedList) Len() int {
	return s.list.Len()
}

func (s IndexedList) Front() *list.Element {
	return s.list.Front()
}

func (s IndexedList) Back() *list.Element {
	return s.list.Back()
}
