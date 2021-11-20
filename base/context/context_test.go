package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextWithCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	select {
	case <-ctx.Done():
		fmt.Println("手动取消")
		return
	}
}
func TestContextWithCancel1(t *testing.T) {
	//contextWithCancel()
	contextWithDeadLine(time.Now().Add(1500 * time.Millisecond))
}

func TestContextWithValue1(t *testing.T) {
	kv := make(map[string]interface{})
	kv["name"] = "pibigstar"
	kv["age"] = 20
	contextWithValue(kv)
}

func TestContextWithDeadLine(t *testing.T) {
	ctx := context.Background()
	deadTime := time.Now().Add(time.Second * 3)
	ctx, _ = context.WithDeadline(ctx, deadTime)

	select {
	case <-ctx.Done():
		fmt.Println("到达deadLine,结束执行")
		return
	}
}

// Internally, context.WithTimeout calls the
// context.WithDeadline function and generates
// the deadline by adding timeout to the current time.
func TestContextWithTimeout(t *testing.T) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*3) // = WithDeadline(parent, time.Now().Add(time.Second*3))
	select {
	case <-ctx.Done():
		fmt.Println("已超时,结束执行")
	}
}

func TestContextWithValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "name", "派大星")
	fmt.Println(ctx.Value("name"))
}
