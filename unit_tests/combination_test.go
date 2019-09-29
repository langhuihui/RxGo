package unit_tests

import (
	. "../rx"
	. "testing"
	"time"
)

func LogData(t *T) func(interface{}) {
	return func(data interface{}) {
		t.Log(data)
	}
}
func Test_Merge(t *T) {
	Merge(Timeout(time.Second), Timeout(time.Second), Timeout(time.Second)).SubscribeSync(LogData(t))
}

func Test_Concat(t *T) {
	Concat(Timeout(time.Second), Of(2, 4, 5), Timeout(time.Second)).SubscribeSync(LogData(t))
}

func Test_Share(t *T) {
	share := Timeout(time.Second).Share()
	var a time.Time
	var b time.Time
	share.SubscribeAsync(func(data interface{}) {
		a = data.(time.Time)
	}, func() {}, func(err error) {})
	share.SubscribeSync(func(data interface{}) {
		b = data.(time.Time)
	})
	if a != b {
		t.FailNow()
	}
}
