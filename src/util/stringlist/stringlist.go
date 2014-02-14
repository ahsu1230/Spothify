// Unbounded buffer, where underlying values are arbitrary values
package stringlist

import (
	"errors"
	"sync"
)

const MissingItemError = "Error: Item not found"
const EmptyListError = "Error: Empty List"

// Linked list element
type ListEle struct {
	value string
	next *ListEle
	prev *ListEle
}

// Linked list
type List struct {
	head *ListEle
	tail *ListEle
	elements map[string](*ListEle) // Map of keys to their element in the list
	mutex *sync.RWMutex // Mutex to ensure thread-safety
}

// Creates and returns a new linked list
func NewList() *List {
	newList := new(List)
	newList.head = nil
	newList.tail = nil
	newList.elements = make(map[string](*ListEle))
	newList.mutex = new(sync.RWMutex)
	return newList
}

// Inserts the given string at the end of the list
func (li *List) Insert(value string) {
	li.mutex.Lock()
	// Create new element
	ele := new(ListEle)
	ele.value = value
	ele.prev = li.tail

	if li.head == nil {
		// Inserting into empty list
		li.head = ele
	} else {
		li.tail.next = ele
	}
	li.tail = ele
	li.elements[value] = ele
	li.mutex.Unlock()
}

// Returns the front element of the list
func (li *List) Front() (string, error) {
	li.mutex.RLock()
	if li.head == nil {
		li.mutex.RUnlock()
		return "", errors.New(EmptyListError)
	}
	value := li.head.value
	li.mutex.RUnlock()
	return value, nil
}

// Returns and removes the front element of the list
func (li *List) Pop() (string, error) {
        li.mutex.Lock()
        if li.head == nil {
                li.mutex.RUnlock()
                return "", errors.New(EmptyListError)
        }
        value := li.head.value
	li.head = li.head.next
	delete (li.elements, value)
        li.mutex.Unlock()
        return value, nil
}

// Returns true if the list contains the given key and false otherwise
func (li *List) Contains(value string) bool {
	li.mutex.RLock()
	_, exists := li.elements[value]
	li.mutex.RUnlock()
	return exists
}

// Removes the element from the list
func (li *List) Remove(value string) error {
	li.mutex.Lock()
	// Check if list is empty
	if li.head == nil {
		li.mutex.Unlock()
		return errors.New(EmptyListError)
	}

	// Check if element is present in list
	ele, exists := li.elements[value]
	if !exists {
		li.mutex.Unlock()
		return errors.New(MissingItemError)
	}

	// Update list
	if ele.next != nil {
		ele.next.prev = ele.prev
	} else {
		li.tail = ele.prev
	}
	if ele.prev != nil {
		ele.prev.next = ele.next
	} else {
		li.head = ele.next
	}

	delete(li.elements, value)
	li.mutex.Unlock()

	return nil
}

// Returns the size of the the list
func (li *List) Size() int {
	li.mutex.RLock()
	size := len(li.elements)
	li.mutex.RUnlock()

	return size
}

// Returns an array containing the elements of the list in order
func (li *List) ToArray() []string {
	li.mutex.RLock()
	arr := make([]string, li.Size())
	ele := li.head
	for c := 0; c < li.Size() ; c++ {
		arr[c] = ele.value
		ele = ele.next
	}
	li.mutex.RUnlock()

	return arr
}

// Returns true if the list is empty and false otherwise
func (li *List) Empty() bool {
	li.mutex.RLock()
	empty := li.head == nil
	li.mutex.RUnlock()

	return empty
}
