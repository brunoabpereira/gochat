package utils

import (
    "errors"
)

/*
	Stack
*/

type Stack[E any] []*E

func (stack *Stack[E]) Push(elem *E) {
	*stack = append(*stack, elem)
}

func (stack *Stack[E]) Pop() (*E, error) {
	l := len(*stack)
	if l == 0 {
		return new(E), errors.New("Empty Stack")
	}
	val := (*stack)[l-1]
	*stack = (*stack)[:l-1]
	return val, nil
}

/*
	Queue
*/

type Queue[E any] struct {
	items []*E
}

func (queue *Queue[E]) IsEmpty() bool{
	return len((*queue).items) == 0
}

func (queue *Queue[E]) Push(elem *E) {
	(*queue).items = append((*queue).items, elem)
}

func (queue *Queue[E]) Pop() (*E, error) {
	l := len((*queue).items)
	if l == 0 {
		return new(E), errors.New("Empty Stack")
	}
	val := (*queue).items[0]
	(*queue).items = (*queue).items[1:l]
	return val, nil
}