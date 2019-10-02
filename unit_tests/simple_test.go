package unit_tests

import (
	"errors"
	"github.com/langhuihui/RxGo/pipe"
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
	"time"
)

func Test_Of(t *T) {
	count := 0
	err := Of(1, 2, 3, 4).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
		count++
	}))
	if count != 4 || err != nil {
		t.FailNow()
	}
}

func Test_StartWith(t *T) {
	Timeout(time.Second).StartWith(1).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_IgnoreElements(t *T) {
	Of(1, 2, 3, 4).IgnoreElements().Subscribe(NextFunc(func(event *Event) {
		t.FailNow()
	}))
}

func Test_Subject(t *T) {
	input := make(chan interface{})
	go Subject(input).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
	input <- 1
	input <- 2
	close(input)
}

func Test_Throw(t *T) {
	err := Throw(errors.New("throw")).Subscribe(NextFunc(func(event *Event) {
		t.FailNow()
	}))
	if err.Error() != "throw" {
		t.FailNow()
	}
}

func Test_Empty(t *T) {
	Empty().Subscribe(NextFunc(func(event *Event) {
		t.FailNow()
	}))
}

func Test_Never(t *T) {
	go Never().Subscribe(NextFunc(func(event *Event) {
		t.FailNow()
	}))
	<-time.After(time.Second)
}

func Test_Pipe(t *T) {
	nextChan := make(NextChan)
	Of(1, 3, 4, 5).Pipe(pipe.Skip(1), pipe.Take(2)).SubscribeAsync(nextChan, func() {
		close(nextChan)
	}, func(e error) {
		t.FailNow()
		close(nextChan)
	})
	event := <-nextChan
	if data := event.Data.(int); data != 3 {
		t.FailNow()
	}
	event = <-nextChan
	if data := event.Data.(int); data != 4 {
		t.FailNow()
	}
	event = <-nextChan
	if event != nil {
		t.FailNow()
	}
}

func Test_Create(t *T) {
	ob := Observable(func(sink *Observer) error {
		sink.Next("hello")
		return nil
	})
	ob.Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_Do(t *T) {
	Of(1, 2, 3).Do(func(i interface{}) {
		t.Log(i)
	}).Subscribe(EmptyNext)
}
