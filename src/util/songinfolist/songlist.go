// Unbounded buffer, where underlying values are arbitrary values
package songlist

import (
	"errors"
	"sync"
	"..\songinfo"
)

const MAX_LIST_SIZE 5

// Linked list element
type ListEle struct {
	value songinfo.SongInfo
	data  []byte
	next *ListEle
	prev *ListEle
}

// Linked list
type List struct {
	head *ListEle
	tail *ListEle
	elements map[songinfo.SongInfo](*ListEle) 	// Map of keys to their element in the list
	mutex *sync.RWMutex // Mutex to ensure thread-safety
}

// Creates and returns a new linked list
func NewList() *List {
	newList := new(List)
	newList.head = nil
	newList.tail = nil
	newList.elements = make(map[songinfo.SongInfo](*ListEle))
	newList.mutex = new(sync.RWMutex)
	return newList
}


// Move element to the front of list
func (li *List) MoveToFront(ele *ListEle) {
	li.mutex.Lock()
	if ele == li.head {
		return
	}
	else if ele == li.tail {
		// set ele's previous as tail
		li.tail = ele.prev
	}
	else {
		// set ele's previous as ele's next
		ele.prev.next = ele.next
	}
	// set ele's next to current head & set ele as new head
	ele.next = li.head
	li.head = ele
	li.mutex.Unlock()
	return
}

// Check if element is in list
func (li *List) Contains(song songinfo.SongInfo) bool {
	li.mutex.RLock()
	_,exists := li.elements[song]
	li.mutex.RUnlock()
	return exists
}

// Retreive SongData - byte array of music file
// Aftre retrieving, move to front
func (li *List) Get(song songinfo.SongInfo) []byte {
	// Is the song in the list?
	if !li.Contains(song) {
		return nil
	}
	li.mutex.RLock()
	ele := li.elements[song]
	li.mutex.RUnlock()

	li.MoveToFront(ele)
	return ele.data
}

// Returns the size of the the list
func (li *List) Size() int {
	li.mutex.RLock()
	size := len(li.elements)
	li.mutex.RUnlock()
	return size
}
// Returns true if the list is empty and false otherwise
func (li *List) Empty() bool {
	li.mutex.RLock()
	empty := li.head == nil
	li.mutex.RUnlock()
	return empty
}

// Add element to the front of list
// If list's size exceeds MAX, remove back element
func (li *List) AddToFront(song songinfo.SongInfo, songbytes []byte) error {
	li.mutex.Lock()
	newEle := new(ListEle)
	newEle.value = song
	newEle.data = songbytes
	
	newEle.next = li.head
	li.head.prev = newEle
	li.head = newEle
	li.elements[song] = songbytes
	li.mutex.Unlock()
	
	if li.Size() > MAX_LIST_SIZE {
		li.PopLast()
	}	
	return nil
}

// Remove element from back of list
func (li *List) PopLast() error {
	if li.Empty() {
		return nil
	}
	li.mutex.Lock()
	lastSong := li.tail.value
	delete(li.elements, lastSong)
	
	lastEle := li.tail
	li.tail = lastEle.prev
	li.tail.next = nil
	lastEle.prev = nil

	li.mutex.Unlock()
	return nil
}