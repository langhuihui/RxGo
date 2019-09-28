package unit_tests

import (
	. "../rx"
	. "testing"
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
