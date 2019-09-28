package unit_tests

import (
	. "../rx"
	. "testing"
	"time"
)

func Test_Take(t *T) {
	count := 0
	Of(1, 2, 3, 4).Take(2).SubscribeAsync(func(data interface{}) {
		t.Log(data)
		count++
		if count > 2 {
			t.FailNow()
		}
	}, func() {
		if count != 2 {
			t.FailNow()
		}
	}, func(err error) {
		t.Error(err)
	})
	<-time.After(time.Second)
}

func Test_Skip(t *T) {
	count := 0
	sub := Of(1, 2, 3, 4).Skip(2).SubscribeAsync(func(data interface{}) {
		t.Log(data)
		count++
		if count > 2 {
			t.FailNow()
		}
	}, func() {
		if count != 2 {
			t.FailNow()
		}
	}, func(err error) {
		t.Error(err)
	})
	t.Log(sub != nil)
	<-time.After(time.Second)
}

func Test_TakeUntil(t *T) {
	Interval(time.Second).TakeUntil(Timeout(time.Second * 4)).SubscribeSync(func(data interface{}) {
		t.Log(data)
	})
}
