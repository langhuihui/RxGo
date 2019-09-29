package unit_tests

import (
	. "../rx"
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
