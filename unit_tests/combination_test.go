package unit_tests

import (
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
	"time"
)

func LogData(t *T) Observer {
	return ObserverFunc(func(event *Event) {
		t.Log(event.Data)
	})
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
	go share.Subscribe(ObserverFunc(func(event *Event) {
		a = event.Data.(time.Time)
	}))
	share.Subscribe(ObserverFunc(func(event *Event) {
		b = event.Data.(time.Time)
	}))
	if a != b {
		t.FailNow()
	}
}

func Test_CombineLatest(t *T) {
	CombineLatest(Of(1, 2), Timeout(time.Second)).Subscribe(ObserverFunc(func(event *Event) {
		data := event.Data
		x, ok := data.([]interface{})
		if ok && len(x) == 2 {
			var a int
			a, ok = x[0].(int)
			if ok && a == 2 {
				_, ok = x[1].(time.Time)
				if !ok {
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
