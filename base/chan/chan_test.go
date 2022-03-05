package chanx_test

import (
	"bufio"
	"fmt"
	chanx "go-demo/base/chan"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestInputChan(t *testing.T) {
	ch := make(chan string)
	defer close(ch)

	scanner := bufio.NewScanner(os.Stdin)

	printInput := func(ch chan string) {
		for value := range ch {
			if value == "eof" {
				break
			}
			t.Logf("input is %v", value)
		}
	}

	go printInput(ch)

	for scanner.Scan() {
		text := scanner.Text()
		ch <- text
		if text == "eof" {
			t.Log("End the Game!")
			break
		}
	}
}

func TestChanRDWR(t *testing.T) {
	ch := make(chan interface{}, 2)
	go func() {
		fmt.Println("go routing access")
		select {
		// 当ch关闭或者有值传入时执行，否则始终阻塞
		case <-ch:
			fmt.Println("receive.")
		}
		fmt.Println("done")
	}()
	ch <- "test"
	time.Sleep(2 * time.Second)

	ch <- 1
	ch <- 2
	// 超出缓存容量 则会阻塞，导致死锁
	//ch <- 3
	close(ch)
	// 结束后再写入会报错，但依然可以读取所有缓存数据
	//ch <- 4

	for c := range ch {
		fmt.Printf("%v\n", c)
	}
}

func TestChanOrOne(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	// 异步定时执行 函数func
	time.AfterFunc(3*time.Second, func() {
		ch1 <- 1
	})

	time.AfterFunc(4*time.Second, func() {
		ch2 <- 2
	})
	t.Log("access")
	//<-chanx.Or(ch1, ch2)
	<-chanx.OrBySelect(ch1, ch2)
	t.Log("收到了信号,开始执行业务逻辑")
}

func TestFanIn(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	// 异步写ch1、ch2, 奇数->ch1、偶数->ch2
	go func() {
		for i := 0; i < 100; i++ {
			if i%2 != 0 {
				ch1 <- i
			} else {
				ch2 <- i
			}
		}
		close(ch1)
		close(ch2)
	}()

	//out := chanx.FanIn(ch1, ch2)
	out := chanx.FanInByReflect(ch1, ch2)

	go func() {
		count := 0
		for i := range out {
			t.Log("out: ", i)
			count++
		}
		// out close之后会执行到这里
		t.Logf("out 共接收到%v个数", count)
	}()

	time.Sleep(2 * time.Second)
}

func TestFanInByReflect(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	go func() {
		for i := 0; i < 100; i++ {
			if i%2 != 0 {
				ch1 <- i
			} else {
				ch2 <- i
			}
		}
		close(ch1)
		close(ch2)
	}()

	out := chanx.FanInByReflect(ch1, ch2)
	go func() {
		count := 0
		defer func() {
			t.Logf("out 共接收到%v个数", count)
		}()
	loop:
		for {
			select {
			// out close之后不会panic， 而是读到nil
			case i := <-out:
				if i == nil {
					break loop
				}
				count++
				t.Log("out: ", i)
			}
		}
	}()

	time.Sleep(1 * time.Second)
}

func TestFanOut(t *testing.T) {
	ch1 := make(chan interface{})
	chs := []chan interface{}{
		make(chan interface{}),
		make(chan interface{}),
		make(chan interface{}),
	}
	// 将ch1收到值，扇出到chs每个chan中
	chanx.FanOut(ch1, chs)

	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-t.C:
				ch1 <- rand.Intn(100)
			}
		}
	}()

	go func() {
		for {
			for _, c := range chs {
				t.Log(<-c)
			}
		}
	}()

	time.Sleep(1 * time.Second)
}

func TestFanOutByReflect(t *testing.T) {
	ch1 := make(chan interface{})
	chs := []chan interface{}{
		make(chan interface{}),
		make(chan interface{}),
		make(chan interface{}),
	}
	// 将ch1收到值，扇出到chs每个chan中
	chanx.FanOutByReflect(ch1, chs)

	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-t.C:
				ch1 <- rand.Intn(100)
			}
		}
	}()

	go func() {
		for {
			for _, c := range chs {
				if v, ok := <-c; ok {
					t.Log(v)
				}
			}
		}
	}()

	time.Sleep(1 * time.Second)
}

func TestFanOutRandom(t *testing.T) {
	ch1 := make(chan interface{})
	chs := []chan interface{}{
		make(chan interface{}, 10),
		make(chan interface{}, 10),
		make(chan interface{}, 10),
	}
	//rand.Seed(time.Now().UnixNano())

	// 将ch1收到值，扇出到chan中任意一个
	chanx.FanOutRandom(ch1, chs)

	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-t.C:
				ch1 <- rand.Intn(100)
			}
		}
	}()

	go func() {
		for {
			for i, c := range chs {
				// v, ok 会立即返回，不会阻塞
				if v, ok := <-c; ok {
					t.Log("i:", i, "v:", v)
				}
			}
		}
	}()

	time.Sleep(3 * time.Second)
}

func TestFanOutRandomByReflect(t *testing.T) {
	ch1 := make(chan interface{})
	chs := []chan interface{}{
		make(chan interface{}, 10),
		make(chan interface{}, 10),
		make(chan interface{}, 10),
	}
	//rand.Seed(time.Now().UnixNano())
	// 将ch1收到值，扇出到chan中任意一个 //实测并不是任一，而是顺序
	chanx.FanOutRandomByReflect(ch1, chs)

	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-t.C:
				ch1 <- rand.Intn(100)
			}
		}
	}()

	go func() {
		for {
			for i, c := range chs {
				if v, ok := <-c; ok {
					t.Log("i:", i, "v:", v)
				}
			}
		}
	}()

	time.Sleep(2 * time.Second)
}
