package rx

import "time"

//Timeout 在一段时间后发送数据
func Timeout(duration time.Duration) Observable {
	return func(sink *Control) error {
		timeout := time.After(duration)
		for {
			select {
			case <-sink.stop:
				return nil
			case data := <-timeout:
				sink.Next(data)
				//sink.Complete()
				return sink.err
			}
		}
	}
}

//Interval 按照一定频率持续发送数据
func Interval(duration time.Duration) Observable {
	return func(sink *Control) error {
		interval := time.NewTicker(duration)
		i := 0
		defer interval.Stop()
		for {
			select {
			case <-sink.stop:
				return sink.err
			case <-interval.C:
				i++
				sink.Next(i)
			}
		}
	}
}
