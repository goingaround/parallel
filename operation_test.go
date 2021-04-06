package parallel

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestOperation_ZeroTimeout(t *testing.T) {
	op, err := NewOperation(func() { time.Sleep(1 * time.Millisecond) })
	require.NoError(t, err)
	require.NoError(t, op.Run(0, 0))
}

func TestOperation_NilFunc(t *testing.T) {
	_, err := NewOperation(nil)
	require.EqualError(t, err, "got nil function (index 0)")
}

func TestOperation_NoFuncsError(t *testing.T) {
	_, err := NewOperation()
	require.EqualError(t, err, "got no functions")
}

func TestOperation_TimeoutExceeded(t *testing.T) {
	op, err := NewOperation(multiplyFunc(func() {}, 3)...)
	require.NoError(t, err)

	err = op.Run(time.Millisecond, 10*time.Millisecond)
	require.True(t, IsTimeoutExceeded(err))
	require.EqualError(t, err, "timeout exceeded")
}

func TestOperation_StartOnSameTime(t *testing.T) {
	res := make([]time.Time, 0, 100)
	mu := new(sync.Mutex)
	op, err := NewOperation(multiplyFunc(func() {
		stamp := time.Now()

		mu.Lock()
		defer mu.Unlock()
		res = append(res, stamp)

	}, 100)...)
	require.NoError(t, err)


	require.NoError(t, op.Run(3*time.Second, 0))
	require.Len(t, res, 100)

	t0 := res[0] // indicator
	for i := 1; i < len(res); i++ {
		require.Equal(t, int64(0), t0.Sub(res[i]).Milliseconds())
	}
}

func TestOperation_UseSameOperationTwice(t *testing.T) {
	op, err := NewOperation(multiplyFunc(func() {}, 10)...)
	require.NoError(t, err)

	require.NoError(t, op.Run(time.Second, 0))
	require.NoError(t, op.Run(time.Second, 0))
}

func multiplyFunc(f func(), n int) []func() {
	if n <= 0 {
		return []func(){}
	}
	res := make([]func(), 0, n)
	for i := 0; i < n; i++ {
		res = append(res, f)
	}
	return res
}
