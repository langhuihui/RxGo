package rx

import "context"

//Subscribe 对外同步订阅Observable，可以在收到的事件中访问Observer的Dispose函数来终止事件流
func (ob Observable) Subscribe(onNext NextHandler) error {
	ctx, cancel := context.WithCancel(context.Background())
	return ob(&Observer{ctx, cancel, onNext})
}

//SubscribeAsync 对外异步订阅模式，返回一个可用于终止事件流的控制器,可以在未收到数据时也能终止事件流
func (ob Observable) SubscribeAsync(onNext NextHandler, onError func(error), onComplete func()) *Observer {
	ctx, cancel := context.WithCancel(context.Background())
	source := &Observer{ctx, cancel, onNext}
	go func() {
		err := ob(source)
		if !source.IsDisposed() {
			if err != nil {
				onError(err)
			} else {
				onComplete()
			}
		}
	}()
	return source
}
