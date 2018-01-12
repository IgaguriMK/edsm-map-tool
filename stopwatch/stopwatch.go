package stopwatch

import (
	"fmt"
	"os"
	"time"
)

const ChannelBuffer = 4

var ch chan<- mark
var cch chan byte

type mark struct {
	Message  string
	Time     time.Time
	TierDiff int
}

func init() {
	c := make(chan mark, ChannelBuffer)
	ch = c

	cch = make(chan byte)

	go marker(c)
}

func Mark(msg string) {
	ch <- mark{
		Message:  msg,
		Time:     time.Now(),
		TierDiff: 0,
	}
}

func StartTier(msg string) {
	ch <- mark{
		Message:  msg,
		Time:     time.Now(),
		TierDiff: 1,
	}
}

func EndTier(msg string) {
	ch <- mark{
		Message:  msg,
		Time:     time.Now(),
		TierDiff: -1,
	}
}

func Close(msg string) {
	EndTier(msg)
	close(ch)
	<-cch
}

func marker(c <-chan mark) {
	var start int64 = -1
	lasts := make([]int64, 0, 32)
	tier := -1

	for {
		m, ok := <-c
		if !ok {
			cch <- 0
			return
		}
		mt := m.Time.UnixNano()

		if start == -1 {
			start = mt
		}
		if m.TierDiff == 1 {
			lasts = append(lasts, mt)
			tier++
		}

		t := mt - start
		l := mt - lasts[tier]
		total := 1e-6 * float64(t)
		lap := 1e-6 * float64(l)

		fmt.Fprintf(os.Stderr, "%s%.3f ms, %.3f ms: %s\n", indent(tier), total, lap, m.Message)

		lasts[tier] = mt

		if m.TierDiff == -1 {
			lasts = lasts[0:tier]
			tier--
		}
	}
}

var indentCache = make([]string, 0)

func indent(n int) string {
	if n < len(indentCache) {
		return indentCache[n]
	}

	if n <= 0 {
		indentCache = append(indentCache, "")
		return ""
	}

	v := "  " + indent(n-1)
	indentCache = append(indentCache, v)
	return v
}
