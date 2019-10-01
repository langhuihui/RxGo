package unit_tests

import (
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
	"time"
)

func Test_Take(t *T) {
	count := 0
	Of(1, 2, 3, 4).Take(2).Subscribe(ObserverFunc(func(event *Event) {
		t.Log(event.Data)
		count++
		if count > 2 {
			t.FailNow()
		}
	}))
	if count != 2 {
		t.FailNow()
	}
}

func Test_Skip(t *T) {
	count := 0
	Of(1, 2, 3, 4).Skip(2).Subscribe(ObserverFunc(func(event *Event) {
		t.Log(event.Data)
		count++
		if count > 2 {
			t.FailNow()
		}
	}))
	if count != 2 {
		t.FailNow()
	}
}

func Test_TakeUntil(t *T) {
	Interval(time.Second).TakeUntil(Timeout(time.Second * 4)).Subscribe(ObserverFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_TakeWhile(t *T) {
	Of(1, 2, 3, 19, 1).TakeWhile(func(data interface{}) bool {
		num, ok := data.(int)
		if ok {
			ok = num < 10
		}
		return ok
	}).Subscribe(ObserverFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_SkipWhile(t *T) {
	Of(1, 2, 3, 19, 1).SkipWhile(func(data interface{}) bool {
		num, ok := data.(int)
		if ok {
			ok = num < 10
		}
		return ok
	}).Subscribe(ObserverFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_SkipUntil(t *T) {
	Interval(time.Second).Take(5).SkipUntil(Timeout(time.Second * 3)).Subscribe(ObserverFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
