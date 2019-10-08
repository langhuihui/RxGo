package unit_tests

import (
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
	"time"
)

func Test_Take(t *T) {
	Of(1).Take(0).Subscribe(nil)
	count := 0
	Of(1, 2, 3, 4).Take(2).Subscribe(NextFunc(func(event *Event) {
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
	Of(1).Skip(0).Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
	count := 0
	Of(1, 2, 3, 4).Skip(2).Subscribe(NextFunc(func(event *Event) {
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
	Interval(time.Second).TakeUntil(Timeout(time.Second * 4)).Subscribe(NextFunc(func(event *Event) {
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
	}).Subscribe(NextFunc(func(event *Event) {
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
	}).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_SkipUntil(t *T) {
	Interval(time.Second).Take(5).SkipUntil(Timeout(time.Second * 3)).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_Filter(t *T) {
	Range(0, 10).Filter(func(data interface{}) bool {
		return data.(int)%2 == 0
	}).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_Distinct(t *T) {
	Merge(Range(0, 10), Range(1, 11)).Distinct().Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_DistinctUntilChanged(t *T) {
	Merge(Range(0, 10), Range(1, 11)).DistinctUntilChanged().Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}

func Test_Debounce(t *T) {
	Interval(time.Second).Debounce(func(i interface{}) Observable {
		return Timeout(time.Second * time.Duration(i.(int)))
	}).Take(3).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_DebounceTime(t *T) {
	Interval(time.Second).DebounceTime(time.Second * 2).Take(3).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_Throttle(t *T) {
	Interval(time.Second).Throttle(func(i interface{}) Observable {
		return Timeout(time.Second * time.Duration(i.(int)))
	}).Take(3).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_ThrottleTime(t *T) {
	Interval(time.Second).ThrottleTime(time.Second * 2).Take(3).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_First(t *T) {
	Timer(time.Second, time.Second).First().Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_Last(t *T) {
	Timer(time.Second, time.Second).Take(2).Last().Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_ElementAt(t *T) {
	Range(0, 10).ElementAt(5).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_Find(t *T) {
	Range(0, 10).Find(func(i interface{}) bool {
		return i.(int) == 5
	}).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
func Test_FindIndex(t *T) {
	Range(1, 10).FindIndex(func(i interface{}) bool {
		return i.(int) == 5
	}).Subscribe(NextFunc(func(event *Event) {
		t.Log(event.Data)
	}))
}
