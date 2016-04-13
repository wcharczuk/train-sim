package simulation

import "github.com/blendlabs/go-util/collections"

func NewQueueOfTrain() *QueueOfTrain {
	return &QueueOfTrain{
		innerQueue: collections.NewQueue(),
	}
}

type QueueOfTrain struct {
	innerQueue *collections.Queue
}

func (q *QueueOfTrain) Enqueue(t *Train) {
	q.innerQueue.Push(t)
}

func (q *QueueOfTrain) Dequeue() *Train {
	value := q.innerQueue.Dequeue()
	if value != nil {
		if typed, isTyped := value.(*Train); isTyped {
			return typed
		}
	}
	return nil
}

func (q *QueueOfTrain) Peek() *Train {
	value := q.innerQueue.Peek()
	if value != nil {
		if typed, isTyped := value.(*Train); isTyped {
			return typed
		}
	}
	return nil
}

func (q *QueueOfTrain) PeekBack() *Train {
	value := q.innerQueue.PeekBack()
	if value != nil {
		if typed, isTyped := value.(*Train); isTyped {
			return typed
		}
	}
	return nil
}

func (q *QueueOfTrain) Len() int {
	return q.innerQueue.Length()
}

func NewQueueOfPassenger() *QueueOfPassenger {
	return &QueueOfPassenger{
		innerQueue: collections.NewQueue(),
	}
}

type QueueOfPassenger struct {
	innerQueue *collections.Queue
}

func (q *QueueOfPassenger) Enqueue(p *Passenger) {
	q.innerQueue.Push(p)
}

func (q *QueueOfPassenger) Peek() *Passenger {
	value := q.innerQueue.Peek()
	if value != nil {
		if typed, isTyped := value.(*Passenger); isTyped {
			return typed
		}
	}
	return nil
}

func (q *QueueOfPassenger) PeekBack() *Passenger {
	value := q.innerQueue.PeekBack()
	if value != nil {
		if typed, isTyped := value.(*Passenger); isTyped {
			return typed
		}
	}
	return nil
}

func (q *QueueOfPassenger) Dequeue() *Passenger {
	value := q.innerQueue.Dequeue()
	if value != nil {
		if typed, isTyped := value.(*Passenger); isTyped {
			return typed
		}
	}
	return nil
}

func (q *QueueOfPassenger) Len() int {
	return q.innerQueue.Length()
}
