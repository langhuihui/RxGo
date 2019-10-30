package unit_tests

import (
	. "github.com/langhuihui/RxGo/rx"
	. "testing"
)

func Test_MergeMap(t *T) {
	Range(0, 10).MergeMap(func(i interface{}) Observable {
		return Range(i.(int), 2)
	}, nil).Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
}
func Test_MergeMapTo(t *T) {
	Range(0, 10).MergeMapTo(Range(1, 2), nil).Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
}
func Test_SwitchMap(t *T) {
	Range(0, 10).SwitchMap(func(i interface{}) Observable {
		return Range(i.(int), 2)
	}, nil).Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
}
func Test_SwitchMapTo(t *T) {
	Range(0, 10).SwitchMapTo(Range(1, 2), nil).Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
}
func Test_Scan(t *T) {
	Range(0, 10).Scan(func(aac interface{}, c interface{}) interface{} { return aac.(int) + c.(int) }).Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
}
func Test_Repeat(t *T) {
	Range(0, 5).Repeat(3).Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
}
func Test_PairWise(t *T) {
	Range(0, 5).PairWise().Subscribe(NextFunc(func(event *Event) { t.Log(event.Data) }))
}
