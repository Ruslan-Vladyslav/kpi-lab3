package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправленя останнього разу у Receiver

	mq messageQueue

	stopped chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stopped = make(chan struct{})

	go func() {
		for !l.stopReq || !l.mq.empty() {
			op := l.mq.pull()
			update := op.Do(l.next)
			if update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stopped)
	}()

}

func (l *Loop) Post(op Operation) {
	if op != nil {
		l.mq.push(op)
	}
}

func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(screen.Texture) {
		l.stopReq = true
	}))
	<-l.stopped
}

type messageQueue struct {
	ops   []Operation
	mut   sync.Mutex
	noacs chan struct{}
}

func (mq *messageQueue) push(op Operation) {
	mq.mut.Lock()
	defer mq.mut.Unlock()

	mq.ops = append(mq.ops, op)

	if mq.noacs != nil {
		close(mq.noacs)
		mq.noacs = nil
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mut.Lock()
	defer mq.mut.Unlock()

	for len(mq.ops) == 0 {
		mq.noacs = make(chan struct{})
		mq.mut.Unlock()
		<-mq.noacs
		mq.mut.Lock()
	}

	op_res := mq.ops[0]
	mq.ops[0] = nil
	mq.ops = mq.ops[1:]

	return op_res
}

func (mq *messageQueue) empty() bool {
	mq.mut.Lock()
	defer mq.mut.Unlock()

	return len(mq.ops) == 0
}
