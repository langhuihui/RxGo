package rx

import "time"

//Timeout 在一段时间后发送数据
func Timeout(duration time.Duration) Observable {
	return func(sink *Observer) error {
		timeout := time.After(duration)
		for {
			select {
			case <-sink.Done():
				return nil
			case data := <-timeout:
				sink.Next(data)
				//sink.Complete()
				return nil
			}
		}
	}
}

//Interval 按照一定频率持续发送数据
func Interval(duration time.Duration) Observable {
	return func(sink *Observer) error {
		interval := time.NewTicker(duration)
		i := 0
		defer interval.Stop()
		for {
			select {
			case <-sink.Done():
				return nil
			case <-interval.C:
				sink.Next(i)
				i++
			}
		}
	}
}

//Timer 延迟+间隔
func Timer(delay time.Duration, interval time.Duration) Observable {
	return func(sink *Observer) error {
		<-time.After(delay)
		return Interval(interval)(sink)
	}
}
