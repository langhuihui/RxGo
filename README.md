# RxGo 非官方实现版本
[![Build Status](https://travis-ci.org/langhuihui/RxGo.svg?branch=master)](https://travis-ci.org/langhuihui/RxGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/langhuihui/RxGo)](https://goreportcard.com/report/github.com/langhuihui/RxGo)
[![codecov](https://codecov.io/gh/langhuihui/RxGo/branch/master/graph/badge.svg)](https://codecov.io/gh/langhuihui/RxGo)

- 代码精简、可读性强（尽可能用最少的代码实现）
- 设计精妙、实现优雅（尽可能利用golang的语言特点和优势）
- 可扩展性强（可自定义Observable以及Operator）
- 占用系统资源低（尽一切可能减少创建goroutine和其他对象）
- 性能强（尽一切可能减少计算量）

**每一行代码都是深思熟虑……**

## 已实现的功能

### Observable

**`FromSlice`**
**`FromChan`**
**`Of`** 
**`Range`**
**`Subject`**
**`Timeout`**
**`Interval`**
**`Merge`**
**`Concat`**
**`Race`**
**`CombineLatest`**
**`Empty`**
**`Never`**
**`Throw`**
### Operator
**`Do`**
**`Take`** 
**`TakeWhile`**
**`TakeUntil`**
**`Skip`** 
**`SkipWhile`** 
**`SkipUntil`**
**`IgnoreElements`**
**`Share`**
**`StartWith`**
**`Zip`**
**`Filter`**
**`Distinct`**
**`DistinctUntilChanged`**
**`Debounce`**
**`DebounceTime`**
**`Throttle`**
**`ThrottleTime`**
**`First`**
**`Last`**
**`Count`**
**`Max`**
**`Min`**
**`Reduce`**
**`Map`**
**`MapTo`**
**`MergeMap`**
**`MergeMapTo`**
**`SwitchMap`**
**`SwitchMapTo`**
**`Scan`**
**`Repeat`**
**`PairWise`**
**`Buffer`**
## 使用方法
### 链式调用方式
```go
import (
    . "github.com/langhuihui/RxGo/rx"
)
func main(){
    err := Of(1, 2, 3, 4).Take(2).Subscribe(NextFunc(func(event *Event) {
        
    }))
}
```
### 管道模式
```go
import (
    . "github.com/langhuihui/RxGo/rx"
    . "github.com/langhuihui/RxGo/pipe"
)
func main(){
    err := Of(1, 2, 3, 4).Pipe(Skip(1),Take(2)).Subscribe(NextFunc(func(event *Event) {
        
    }))
}
```

管道模式相比链式模式，具有操作符**可扩展性**，用户可以按照规则创建属于自己的操作符

```go
type Operator func(Observable) Observable
```

操作符只需要返回Operator这个类型即可,例如
实现一个happy为false就立即完成的操作符
```go
func MyOperator(happy bool) Operator {
	return func(source Observable) Observable {
		return func (sink *Observer) error {
            if happy{
                return source(sink)
            }
            return nil
		}
	}
}

```

### 创建自定义Observable
在任何时候，您都可以创建自定义的Observable，用来发送任何事件
```go
import (
    . "github.com/langhuihui/RxGo/rx"
)
func MyObservable (sink *Control) error {
    sink.Next("hello")
    return nil
}
func main(){
    ob := Observable(MyObservable)
    ob.Subscribe(NextFunc(func(event *Event) {
        
    }))
}
```

# 设计思想

## 基本知识
所谓`Observable`，就是一个可以被订阅，然后不断发送事件的事件源，见如下示意图
```bash
                                time -->

(*)-------------(o)--------------(o)---------------(x)----------------|>
 |               |                |                 |                 |
Start          value            value             error              Done
```
该示意图代表了，事件被订阅后（Start）开始不停发送事件的过程，直到发出error或者Done（完成）为止

有的`Observable`并不会发出完成事件，比如`Never`

参考网站：
[rxmarbles](https://rxmarbles.com/)

## 总体方案

实现Rx的关键要素，是要问几个问题
1. 如何定义`Observable`，？（一个结构体？一个函数？一个接口？一个`Channel`？）
2. 如何实现**订阅**逻辑？（调用函数？发送数据？）
3. 如何实现**接受数据**？（如何实现`Observer`？）
4. 如何实现**完成/错误**的传递？
5. 如何实现**取消订阅**？（难点：在事件响应中取消，以及在任何其他goroutine中取消订阅）
6. 如何实现**操作符**？
7. **操作符**如何处理连锁反应，比如后面描述的情况
8. 如何实现链式和管道两种编程模式
9. 如何让用户扩展（自定义）`Observable`和**操作符**
10. 如何向普通用户解释复杂的概念

- 当用户需要订阅或者终止事件流，则进行链路传递，订阅或者终止所有中间过程中的事件源
```bash

Observable---------Operator----------Operator-----------Observer
             <|                <|                <|          
           订阅/取消          订阅/取消          订阅/取消         
```
- 当事件流完成或者报错时，需要通知下游事件流的完成或者报错

```bash

Observable---------Operator----------Operator-----------Observer
             |>                |>                |>          
           完成/错误          完成/错误          完成/错误         
```
实际情况远比这个复杂，后面会进行分析

### 可观察对象（事件源）`Observable`
Observable 被定义成为一个函数，该函数含有一个类型为*Observer的参数。
```go
type Observable func(*Observer) error
```
任何事件源都是这样的一个函数，当调用该函数即意味着**订阅**了该事件源，入参为一个`Observer`，具体功能见下面

如果该函数返回nil，即意味着**事件流完成**

否则意味着**事件流异常**

### 观察者对象`Observer`
```go
type Observer struct {
	next NextHandler //缓存当前的NextHandler，后续可以被替换
	disposed    bool             //是否已经取消
	disposeList []DisposeHandler //缓存的DisposeHandler
	err         error            //缓存当前的错误
	lock        sync.Mutex       //用于Dispose的锁
}
```
该控制器为一个结构体，其中next记录了当前的NextHandler，

**取消订阅**的行为是设置标志位，然后级联调用被缓存进来的DisposeHandler。
这些DisposeHandler都是由上由Observable设置的，用来释放资源。
```go
//Dispose 取消订阅
func (c *Observer) Dispose() {
	c.lock.Lock()
	if c.disposed {
		return
	}
	c.disposed = true
	c.lock.Unlock()
	for _, handler := range c.disposeList {
		handler.Dispose()
	}
}
```

`Observer`对象为`Observable`和事件处理逻辑共同持有，是二者沟通的桥梁

之所以如此设计，是因为Observable如果返回结构体会使得设计变复杂，Golang没有继承，所以不好设计

现在将Dispose方法放到Observer对象中，可以方便在下游和上游调用

### NextHandler
```go
type Event struct {
    Data    interface{}
    Target  *Observer
}
NextHandler interface {
    OnNext(*Event)
}
```
`NextHandler`是一个接口，实现`OnNext`函数，当`Observable`数据推送到`Observer`中时，即调用了该函数

`Target`属性用于存储当前发送事件的`Observer`对象，有两大重要使命
1. 更换NextHandler，用于减少数据传递过程
2. 在NextHandler过程中终止事件流

这样做的好处是可以实现不同的观察者，比如函数或者channel
```go
type(
    NextFunc func(*Event)
    NextChan chan *Event
)
func (next NextFunc) OnNext(event *Event) {
	next(event)
}
func (next NextChan) OnNext(event *Event) {
	next <- event
}
```
## 实现案例TakeUntil
```go
//TakeUntil 一直获取事件直到unitl传来事件为止
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Observer) error {
		observer := &Observer{next: sink.next}
		utilObserver := &Observer{next: NextFunc(func(event *Event) {
			//获取到任何数据就让下游完成
			event.Target.Dispose()
			observer.Dispose()
		})}
		go until(utilObserver)
		sink.Defer(utilObserver, observer)
		return ob(observer)
	}
}
```
TakeUnitl的用途是，传入一个until事件源，当这个until事件源接受到事件时，就会导致当前的事件源"完成”。相当于某种中断信号。

看似简短的代码，确考虑各种不同的情形

几大实现细节：
1. 订阅until事件源，通过go关键字创建goroutine防止阻塞当前goroutine
2. 使用函数式NextHandler，用户接受来自until事件源的事件，一旦接受任何事件，就调用observer.Dispose()来使得当前事件流完成
3. sink.Defer(utilObserver, observer)的含义是当下游取消订阅时，将这两个级联取消————取消订阅行为的向上传播
4. 最后一步是订阅上游事件源，并将返回结果返回————上游Observable完成也意味着本Observable完成即完成信号的向下传播
5. 任何情况取消订阅，或者上游事件源完成都可以使得事件源函数返回，接着TakeUntil函数也会返回,即意味着完成
6. until事件源的完成或者错误，都将忽略，所以我们没有去获取until函数返回值