package chanx

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
)

// 当等待多个信号的时候，如果收到任意一个信号， 就执行业务逻辑，忽略其它的还未收到的信号
func Or(chs ...<-chan interface{}) <-chan interface{} {
	// 声明为无缓冲chan，可在外部阻塞读
	out := make(chan interface{})
	// 启动异步携程，不阻塞主线程
	go func() {
		var once sync.Once
		// 任意一个chan收到信号，会在异步协程中close out，外部就会接收到信号，其他等待携程就会执行到case <-out, 携程就会结束
		for _, c := range chs {
			// 开启n个携程，阻塞监听各个chan，有任一信号则close out，同时释放其他n-1个阻塞携程
			go func(c <-chan interface{}) {
				select {
				case <-c:
					once.Do(func() {
						close(out) // 关闭out,提醒外部可以继续执行了
					})
					// 用于out关闭后，退出select
				case <-out:
				}
			}(c)
		}
	}()
	return out
}

// OrBySelect 通过select方式，实现起来更简洁
func OrBySelect(channels ...<-chan interface{}) <-chan interface{} {
	// edge case process
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	// 创建对外通知用的chan，用作handler
	orDone := make(chan interface{})
	// 启动异步协程监听多源chan信号
	go func() {
		// 监听到信号后对外通知，通过关闭handler的方式
		defer close(orDone)
		// 声明用于反射select的case数组
		var cases []reflect.SelectCase
		for _, c := range channels {
			// 填充case数组：case chan为只读chan
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		// 执行select，监听任一chan
		reflect.Select(cases)
	}()
	return orDone
}

// FanIn 将多个同样类型的输入channel合并成一个同样类型的输出channel
func FanIn(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		// 使用wg阻塞该携程close，保证当n个chan都close后，再close out
		var wg sync.WaitGroup
		wg.Add(len(chs))
		// 启动n个携程，内部使用for来接受各自chan写入，转送到out；chan close之后通知wg结束
		for _, c := range chs {
			go func(c <-chan interface{}) {
				for v := range c {
					out <- v
				}
				wg.Done()
			}(c)
		}
		// wg统计n个chan都结束
		wg.Wait()
		fmt.Println("out close")
	}()
	return out
}

// 将多个同样类型的输入channel合并成一个同样类型的输出channel
func FanInByReflect(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		var cases []reflect.SelectCase
		for _, c := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		for len(cases) != 0 {
			// 执行 select，也就是从 chan 中接受值
			// 这里将会一直阻塞，直到输入chan close，
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:i], cases[i+1:]...)
			} else {
				out <- v
			}
		}
		fmt.Println("fin in chan close")
	}()

	return out
}

// 将一个输入channel扇出为多个channel
func FanOut(ch <-chan interface{}, out []chan interface{}) {
	// 异步处理输入，然后扇出
	go func() {
		// 保证关闭所有fan out chan
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		// 接收ch的所有输入
		for v := range ch {
			// 这里为啥要重新复制一下？
			//v := v
			for i := 0; i < len(out); i++ {
				out[i] <- v
			}
		}
	}()
}

// 将一个输入channel扇出为多个channel
func FanOutByReflect(ch <-chan interface{}, out []chan interface{}) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		cases := make([]reflect.SelectCase, len(out))

		for i := range cases {
			cases[i].Dir = reflect.SelectSend
			cases[i].Chan = reflect.ValueOf(out[i])
		}

		// 阻塞一直到ch被close
		for v := range ch {
			for i := range cases {
				cases[i].Send = reflect.ValueOf(v)
			}
			for range cases {
				// 执行 select，也就是将 v 发送到 所有的cases中
				reflect.Select(cases)
			}
		}
	}()
}

// 将一个输入channel扇出到out chan中任意一个
func FanOutRandom(ch <-chan interface{}, out []chan interface{}) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		var n = len(out)
		for v := range ch {
			i := rand.Intn(n)
			fmt.Println("random: ", i, " ", v)
			out[i] <- v
		}
	}()
}

// 将一个输入channel扇出到out chan中任意一个
func FanOutRandomByReflect(ch <-chan interface{}, out []chan interface{}) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		cases := make([]reflect.SelectCase, len(out))

		for i := range cases {
			cases[i].Dir = reflect.SelectSend
			cases[i].Chan = reflect.ValueOf(out[i])
		}

		// 阻塞一直到ch被close
		for v := range ch {
			for i := range cases {
				cases[i].Send = reflect.ValueOf(v)
			}
			// 执行 select，也就是将 v 发送到 任意一个cases中
			i, _, _ := reflect.Select(cases)
			fmt.Println("random: ", i, " ", v)
		}
	}()
}
