package unit_tests

import (
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
	"time"
)

func LogData(t *T) NextFunc {
	return func(event *Event) {
		t.Log(event.Data)
	}
}
func Test_Merge(t *T) {
	Merge(Timeout(time.Second), Timeout(time.Second), Timeout(time.Second)).Subscribe(LogData(t))
}

func Test_Concat(t *T) {
	Concat(Timeout(time.Second), Of(2, 4, 5), Timeout(time.Second)).Subscribe(LogData(t))
}

func Test_Share(t *T) {
	share := Timeout(time.Second).Share()
	var a time.Time
	var b time.Time
	go share.Subscribe(NextFunc(func(event *Event) {
		a = event.Data.(time.Time)
	}))
	share.Subscribe(NextFunc(func(event *Event) {
		b = event.Data.(time.Time)
	}))
	if a != b {
		t.FailNow()
	}
}

func Test_CombineLatest(t *T) {
	CombineLatest(Of(1, 2), Timeout(time.Second)).Subscribe(NextFunc(func(event *Event) {
		data := event.Data
		x, ok := data.([]interface{})
		if ok && len(x) == 2 {
			var a int
			a, ok = x[0].(int)
			if ok && a == 2 {
				if _, ok = x[1].(time.Time); !ok {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	}))

}
func Test_Zip(t *T) {
	Zip(Of(1, 2), Interval(time.Second)).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_Race(t *T) {
	Race(Timeout(time.Second), Timeout(time.Second), Interval(time.Second)).Take(3).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
