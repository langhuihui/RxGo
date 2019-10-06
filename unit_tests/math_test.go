package unit_tests

import (
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
)

func Test_Count(t *T) {
	Range(0, 10).Count().Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_Max(t *T) {
	Range(-1, 10).Max().Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_Min(t *T) {
	Range(-1, 10).Min().Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_Reduce(t *T) {
	Range(0, 10).Reduce(func(aac interface{}, c interface{}) interface{} {
		return aac.(int) + c.(int)
	}).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
