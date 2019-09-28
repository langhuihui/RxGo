package rx

//Merge 合并多个事件流
func Merge(sources ...Observable) Observable {
	return func(sink *Control) {
		scs := make(ControlSet)
		onNext := func(event *Event) {
			if event.err == nil {
				sink.Push(event)
			} else {
				scs.remove(event.control)
				if scs.isEmpty() {
					sink.Complete()
				}
			}
		}
		for _, source := range sources {
			scs.add(source.Subscribe(onNext))
		}
		for sink.Check() {

		}
		for sc := range scs {
			sc.Stop()
		}
	}
}

//Share 共享数据源
func (ob Observable) Share() Observable {
	children := make(ControlSet)
	var source *Control
	fromSource := func(event *Event) {
		for sink := range children {
			sink.Push(event)
		}
	}
	return func(sink *Control) {
		if children.isEmpty() {
			source = ob.Subscribe(fromSource)
		}
		children.add(sink)
		for sink.Check() {
		}
		children.remove(sink)
		if children.isEmpty() {
			source.Stop()
		}
	}
}
