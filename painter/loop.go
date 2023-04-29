package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циелі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправленя останнього разу у Receiver

	mq messageQueue
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	// TODO: ініціалізувати чергу подій.
	// TODO: запустити рутину обробки повідомлень у черзі подій.
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	// TODO: реалізувати додавання операції в чергу. Поточна імплементація
	update := op.Do(l.next)
	if update {
		l.Receiver.Update(l.next)
		l.next, l.prev = l.prev, l.next
	}
}

// StopAndWait сигналізує
func (l *Loop) StopAndWait() {

}

// TODO: реалізувати власну чергу повідомлень.
type messageQueue struct {
	operations []Operation
	mu         sync.Mutex
	blocked    chan struct{}
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.operations = append(mq.operations, op)

	if mq.blocked != nil {
		close(mq.blocked)
		mq.blocked = nil
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for len(mq.operations) == 0 {
		mq.blocked = make(chan struct{})
		mq.mu.Unlock()
		<-mq.blocked
		mq.mu.Lock()
	}

	op := mq.operations[0]
	mq.operations[0] = nil
	mq.operations = mq.operations[1:]
	return op
}

func (mq *messageQueue) isEmpty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	return len(mq.operations) == 0
}
