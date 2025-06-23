package model

import "errors"

type LinkedList interface {
	Insert(value interface{})
	Remove(value interface{})
	Contains(value interface{}) bool
	Get(index int) (interface{}, error)
	GetAll() []interface{}
}

type linkedListImpl struct {
	Start  *Node
	Finish *Node
	Count  int
}

type Node struct {
	Value interface{}
	Next  *Node
	Prev  *Node
}

func NewLinkedList() LinkedList {
	return linkedListImpl{
		Start:  nil,
		Finish: nil,
		Count:  0,
	}
}

func NewNode(value interface{}, previous *Node) *Node {
	return &Node{
		Value: value,
		Next:  nil,
		Prev:  previous,
	}
}

func (les linkedListImpl) Insert(value interface{}) {
	newNode := NewNode(value, les.Finish)
	if les.Start == nil || les.Finish == nil {
		les.Finish = newNode
		les.Start = newNode
		return
	}
	les.Finish.Next = newNode
	return
}

func (les linkedListImpl) Remove(value interface{}) {
	if les.Start == nil || les.Finish == nil {
		return
	}
	current := les.Start
	for current != nil {
		if current.Value == value {
			if current.Prev != nil {
				current.Prev.Next = current.Next
			} else {
				les.Start = current.Next
			}
			if current.Next != nil {
				current.Next.Prev = current.Prev
			} else {
				les.Finish = current.Prev
			}
			les.Count--
			return
		}
	}
}

func (les linkedListImpl) Contains(value interface{}) bool {
	if les.Start == nil || les.Finish == nil {
		return false
	}
	current := les.Start
	for current != nil {
		if current.Value == value {
			return true
		}
		current = current.Next
	}
	return false
}

func (les linkedListImpl) Get(index int) (interface{}, error) {
	if index < 0 || index >= les.Count {
		return nil, errors.New("index out of bounds")
	}
	current := les.Start
	for i := 0; i < index; i++ {
		if current == nil {
			return nil, errors.New("index out of bounds")
		}
		current = current.Next
	}
	return current.Value, nil
}

func (les linkedListImpl) GetAll() []interface{} {
	if les.Start == nil || les.Finish == nil {
		return []interface{}{}
	}
	current := les.Start
	values := make([]interface{}, 0, les.Count)
	for current != nil {
		values = append(values, current.Value)
		current = current.Next
	}
	return values
}
