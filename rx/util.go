package rx

//Do 可以在中间执行一个逻辑
func (ob Observable) Do(f func(interface{})) Observable {
	return func(sink *Observer) error {
		return ob(sink.New3(NextFunc(func(event *Event) {
			f(event.Data)
			sink.Push(event)
		})))
	}
}

func justComplete(event *Event) {
	event.Target.Complete()
}
