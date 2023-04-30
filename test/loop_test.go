package test

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"github.com/hrystynaa/lab3-go/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	screen := new(screenMock)
	texture := new(textureMock)
	receiver := new(receiverMock)
	tx := image.Pt(800, 800)
	l := painter.Loop{
		Receiver: receiver,
	}

	screen.On("NewTexture", tx).Return(texture, nil)
	receiver.On("Update", texture).Return()

	l.Start(screen)

	op1 := new(operationMock)
	op2 := new(operationMock)
	op3 := new(operationMock)

	texture.On("Bounds").Return(image.Rectangle{})
	op1.On("Do", texture).Return(false)
	op2.On("Do", texture).Return(true)
	op3.On("Do", texture).Return(true)

	assert.Empty(t, l.Mq.Operations)
	l.Post(op1)
	l.Post(op2)
	l.Post(op3)
	time.Sleep(1 * time.Second)
	assert.Empty(t, l.Mq.Operations)

	op1.AssertCalled(t, "Do", texture)
	op2.AssertCalled(t, "Do", texture)
	op3.AssertCalled(t, "Do", texture)
	receiver.AssertCalled(t, "Update", texture)
	screen.AssertCalled(t, "NewTexture", image.Pt(800, 800))
}

func TestMessageQueue_Push(t *testing.T) {
	mq := &painter.MessageQueue{}

	op1 := &operationQueueMock{}
	mq.Push(op1)
	if len(mq.Operations) != 1 {
		t.Errorf("Expected queue length to be 1, but got %d", len(mq.Operations))
	}
	if !reflect.DeepEqual(op1, mq.Operations[0]) {
		t.Error("Expected pushed operation to be in the queue")
	}

	op2 := &operationQueueMock{}
	mq.Push(op2)
	if len(mq.Operations) != 2 {
		t.Errorf("Expected queue length to be 2, but got %d", len(mq.Operations))
	}
	if !reflect.DeepEqual(op2, mq.Operations[0]) {
		t.Error("Expected pushed operation to be in the queue")
	}
}

func TestMessageQueue_Pull(t *testing.T) {
	mq := &painter.MessageQueue{}

	op1 := &operationQueueMock{}
	go func() {
		time.Sleep(50 * time.Millisecond)
		mq.Push(op1)
	}()
	start := time.Now()
	op := mq.Pull()
	elapsed := time.Since(start)

	if !reflect.DeepEqual(op, op1) {
		t.Errorf("Expected pulled operation to be the same as the pushed")
	}
	if elapsed < 50*time.Millisecond {
		t.Errorf("Expected Pull to block when pulling from an empty queue")
	}
	if len(mq.Operations) != 0 {
		t.Errorf("Expected queue length to be 0, but got %d", len(mq.Operations))
	}
	op2 := &operationQueueMock{}
	op3 := &operationQueueMock{}
	mq.Push(op2)
	mq.Push(op3)
	op = mq.Pull()
	if len(mq.Operations) != 1 {
		t.Errorf("Expected queue length to be 1, but got %d", len(mq.Operations))
	}
	if !reflect.DeepEqual(op, op2) {
		t.Error("Expected pulled operation to be the first pushed operation")
	}
}

type receiverMock struct {
	mock.Mock
}

func (rm *receiverMock) Update(t screen.Texture) {
	rm.Called(t)
}

type screenMock struct {
	mock.Mock
}

func (sm *screenMock) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (sm *screenMock) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

func (sm *screenMock) NewTexture(size image.Point) (screen.Texture, error) {
	args := sm.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

type textureMock struct {
	mock.Mock
}

func (tm *textureMock) Release() {
	tm.Called()
}

func (tm *textureMock) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	tm.Called(dp, src, sr)
}

func (tm *textureMock) Bounds() image.Rectangle {
	args := tm.Called()
	return args.Get(0).(image.Rectangle)
}

func (tm *textureMock) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	tm.Called(dr, src, op)
}

func (tm *textureMock) Size() image.Point {
	args := tm.Called()
	return args.Get(0).(image.Point)
}

type operationMock struct {
	mock.Mock
}

func (om *operationMock) Do(t screen.Texture) bool {
	args := om.Called(t)
	return args.Bool(0)
}

type operationQueueMock struct{}

func (m *operationQueueMock) Do(t screen.Texture) (ready bool) {
	return false
}
