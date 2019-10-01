package unit_tests

import (
	"errors"
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
	"time"
)

func Test_Of(t *T) {
	count := 0
	err := Of(1, 2, 3, 4).SubscribeSync(func(data interface{}) {
		t.Log(data)
		count++
	})
	if count != 4 || err != nil {
		t.FailNow()
	}
}

func Test_StartWith(t *T) {
	Timeout(time.Second).StartWith(1).SubscribeSync(func(data interface{}) {
		t.Log(data)
	})
}

func Test_IgnoreElements(t *T) {
	Of(1, 2, 3, 4).IgnoreElements().SubscribeSync(func(data interface{}) {
		t.FailNow()
	})
}

func Test_Subject(t *T) {
	input := make(chan interface{})
	Subject(input).SubscribeAsync(func(data interface{}) {
		t.Log(data)
	}, func() {

	}, func(e error) {

	})
	input <- 1
	input <- 2
	close(input)
}

func Test_Throw(t *T) {
	err := Throw(errors.New("throw")).SubscribeSync(func(data interface{}) {
		t.FailNow()
	})
	if err.Error() != "throw" {
		t.FailNow()
	}
}

func Test_Empty(t *T) {
	Empty().SubscribeSync(func(data interface{}) {
		t.FailNow()
	})
}

func Test_Never(t *T) {
	Never().SubscribeAsync(func(data interface{}) {
		t.FailNow()
	}, t.FailNow, func(err error) {
		t.FailNow()
	})
	<-time.After(time.Second)
}
