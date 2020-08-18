package progress

import (
	"io"
	"sync/atomic"
	"time"

	"github.com/cenkalti/log"
	"github.com/paulbellamy/ratecounter"
)

type Progress struct {
	r       io.Reader
	offset  int64
	size    int64
	prefix  string
	counter *ratecounter.RateCounter
	ticker  *time.Ticker
}

func New(r io.Reader, offset, size int64, prefix string) *Progress {
	return &Progress{
		r:       r,
		offset:  offset,
		size:    size,
		prefix:  prefix,
		counter: ratecounter.NewRateCounter(time.Second),
	}
}

func (r *Progress) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.counter.Incr(int64(n))
	atomic.AddInt64(&r.offset, int64(n))
	return n, err
}

func (r *Progress) Start() {
	r.ticker = time.NewTicker(time.Second)
	go r.run()
}

func (r *Progress) run() {
	for range r.ticker.C {
		offset := atomic.LoadInt64(&r.offset)
		progress := (offset * 100) / r.size
		speed := r.counter.Rate() / (1 << 10)
		log.Infof("%s %d/%d MB (%d%%) %d KB/s", r.prefix, offset/(1<<20), r.size/(1<<20), progress, speed)
	}
}

func (r *Progress) Stop() {
	r.ticker.Stop()
}
