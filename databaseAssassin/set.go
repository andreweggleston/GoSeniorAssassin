package databaseAssassin

import "fmt"

type Set struct {
	list *LinkedList
}

func NewSet() *Set {
	return &Set{
		list: nil,
	}
}

func (s *Set) Add(val interface{}) error {
	if s.Contains(val) {
		return fmt.Errorf("Value already found in set")
	}
	// ensure there is a list
	if s.list == nil {

		genericList := make([]interface{}, 1)
		genericList[0] = val
		s.list = NewLinkedList(genericList)
	} else {
		s.list.Push(val)
	}
	return nil
}

func (s *Set) Contains(val interface{}) bool {
	if s.list == nil {
		return false
	} else {
		_, err := s.list.Find(val)
		return err == nil // if no err it is in the linkedlist
	}
}