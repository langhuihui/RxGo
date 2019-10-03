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

## 使用方法
### 链式调用方式
```go
import (
    . "github.com/langhuihui/RxGo/rx"
)
func main(){
    err := Of(1, 2, 3, 4).Take(2).Subscribe(ObserverFunc(func(event *Event) {
        
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
    err := Of(1, 2, 3, 4).Pipe(Skip(1),Take(2)).Subscribe(ObserverFunc(func(event *Event) {
        
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
    ob.Subscribe(ObserverFunc(func(event *Event) {
        
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
type Stop chan bool
type Observer struct {
	next NextHandler //缓存当前的Handler，后续可以被替换
	stop     Stop     //取消订阅的信号，只用来close
	err      error    //缓存当前错误
}
```
该控制器为一个结构体，其中next记录了当前的NextHandler，

在任何时候，如果关闭了stop这个channel，就意味着**取消订阅**。
```go
//Stop 取消订阅
func (c *Observer) Stop() {
	if !c.IsStopped() {
		close(c.stop)
	}
}
//IsStopped 判断是否已经取消订阅
func (c *Observer) IsStopped() bool {
	select {
	case <-c.stop:
		return true
	default:
		return false
	}
}
```

由于Channel的close可以引发所有读取该Channel的阻塞行为唤醒，所以可以在不同层级复用该channel

并且，由于已经close的channel可以反复读取以取得是否close的状态信息，所以不需要再额外记录

`Observer`对象为`Observable`和事件处理逻辑共同持有，是二者沟通的桥梁

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
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Observer) error {
		go until(NewObserver(NextFunc(func(event *Event) {
			//获取到任何数据就让下游完成
			sink.Stop()
		}), sink.stop))
		return ob(sink)
	}
}
```
TakeUnitl的用途是，传入一个until事件源，当这个until事件源接受到事件时，就会导致当前的事件源"完成”。相当于某种中断信号。

看似简短的代码，确考虑各种不同的情形

几大实现细节：
1. 订阅until事件源，通过go关键字创建goroutine防止阻塞当前goroutine
2. 使用函数式NextHandler，用户接受来自until事件源的事件，一旦接受任何事件，就调用sink.Stop()来使得所有使用该关闭channel（此处为sink.stop)的事件源全部取消订阅，并且导致事件源函数返回
3. 订阅until事件源的Observer复用了sink.stop,当用户在代码中取消了订阅，就会引发该until事件源的取消订阅行为
4. 最后一步是订阅上游事件源，我们不创建新的Observer，而直接把下游的Observer传入，避免了不必要的转发逻辑
5. 任何情况取消订阅，或者上游事件源完成都可以使得事件源函数返回，接着TakeUntil函数也会返回,即意味着完成
6. until事件源的完成或者错误，都将忽略，所以我们没有去获取until函数返回值