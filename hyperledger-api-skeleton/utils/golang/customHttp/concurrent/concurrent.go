package concurrent

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type LimitThreads struct {
	Activated        bool `json:"activated"`
	TotalReqActive   int  `json:"totalReqActive"`
	ThreadReqLimiter chan bool
}

// RequestsChannel channel requests structure
type RequestsChannel struct {
	Ctx            context.Context
	TypeR          string        `json:"typeR"`
	Url            string        `json:"url"`
	Body           []byte        `json:"body,omitempty"`
	Request        []byte        `json:"request,omitempty"`
	CacheTimeMs    time.Duration `json:"cacheTimeMs,omitempty"`
	BasicAuthValue string
	ResultChan     chan ResponseChannel `json:"resultChan,omitempty"`
	Opt            *ProviderOptions
}

// LimitRequests Define limit requests options
type LimitRequests struct {
	Activated        bool          `json:"activated"`
	TotalTx          int           `json:"totalTx,omitempty"`
	TimeMs           time.Duration `json:"timeMs,omitempty"`
	AutoBackPressure bool          `json:"autoBackPressure,omitempty"`
}

// CacheRequests Define cache request options
type CacheRequests struct {
	Activated bool          `json:"activated"`
	TimeMs    time.Duration `json:"timeMs"`
	CacheF    func(r RequestsChannel) bool
}

// BPOffline to handle manual backpreasure(Offline requests)
type BPOffline struct {
	Activated bool          `json:"activated"`
	BPChannel *StructOfChan `json:"BPChannel,omitempty"`
	Capacity  int           `json:"capacity"`
}

// ChangeBPChannelCap Change channel bp capacity
func (c *BPOffline) ChangeBPChannelCap() {
	if c.BPChannel != nil {
		capacity := cap(c.BPChannel.SignalCh)
		if capacity != c.Capacity {
			for i := 0; i < capacity; i++ {
				c.BPChannel.SignalCh <- true
			}
			c.BPChannel.SignalCh = make(chan bool, c.Capacity)
		}
	}
}

// GetBPChannelCapAndLen get channel back pressure capacity and length
func (c *BPOffline) GetBPChannelCapAndLen() (int, int) {
	var capacity, length int
	if c.BPChannel != nil {
		capacity = cap(c.BPChannel.SignalCh)
		length = len(c.BPChannel.SignalCh)
	}
	return capacity, length
}

type StructOfChan struct {
	Name         string
	SignalCh     chan bool
	ReqWaiting   int32
	ReqInProcess int32
}

func (s *StructOfChan) ReleaseChannel(ctx context.Context, id string) {
	headerLogMessage := "ReleaseChannel - Offline request. TraceID: " + id + " - "

	<-s.SignalCh

	atomic.AddInt32(&s.ReqInProcess, -1)
	requestInProcess := atomic.LoadInt32(&s.ReqInProcess)
	s.LogCustom.InfoCtx(ctx, headerLogMessage+"End offline request - Total offline requests in process:", requestInProcess)
}

type MetricsChannel struct {
	TotalReqCache int
	CompleteTxs   Counter `json:"completeTx,omitempty"`
}

// ProviderOptions provider options: timeout, etc.
type ProviderOptions struct {
	Timeout          int
	Auth             string
	Attempts         int
	AttemptException *AttemptException
	URLToLogRemotely map[string]string
}

type AttemptException struct {
	TypeR int
	Url   string
}

// ResponseChannel channel response structure
type ResponseChannel struct {
	Code    int             `json:"code"`
	Ctx     context.Context `json:"ctx,omitempty"`
	Request []byte          `json:"request,omitempty"`
}

func (r *RequestsChannel) AddBasicAuth(basicAuthValue string) {
	r.TypeR = r.TypeR + "_AUTH"
	r.BasicAuthValue = basicAuthValue
}

func (r *RequestsChannel) AddOptionsRequest(opt *ProviderOptions) {
	r.TypeR = r.TypeR + "_OPTIONS"
	r.Opt = opt
}

// func (r *ResponseChannel) FindResponse(ctx context.Context, results chan ResponseChannel) {

// 	return response
// }

type CacheResponse struct {
	Resp       ResponseChannel
	TimeCached time.Duration
}

// SafeMapCacheRes Map type that can be safely shared between
// goroutines that require read/write access to a map
type SafeMapCacheRes struct {
	sync.RWMutex
	items map[string]CacheResponse
}

// SafeMapItem Concurrent map item
type SafeMapItem struct {
	Key   string
	Value CacheResponse
}

// Init - Initialize struct
func (cm *SafeMapCacheRes) Init() {
	cm.items = make(map[string]CacheResponse)
}

// Set - Sets a key in a concurrent map
func (cm *SafeMapCacheRes) Set(key string, value CacheResponse) {

	cm.items[key] = value
}

// Get - Gets a key from a concurrent map
func (cm *SafeMapCacheRes) Get(key string) (CacheResponse, bool) {

	value, ok := cm.items[key]

	return value, ok
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *SafeMapCacheRes) iter() <-chan SafeMapItem {
	c := make(chan SafeMapItem)

	return c
}

// Flush - delete all elements from map
func (cm *SafeMapCacheRes) Flush() {
	for k := range cm.items {
		delete(cm.items, k)
	}
}

// ///////////////////////////////////////////////////////////
// Counter can be safely shared between

type Counter struct {
	mu sync.RWMutex
	n  int
}

func (c *Counter) Add() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *Counter) Get() int {
	c.mu.Lock()
	n := c.n
	c.mu.Unlock()
	return n
}

func (c *Counter) Reset() {
	c.mu.Lock()
	if c.n > 8190 {
		c.n = 0
	}
	c.mu.Unlock()
}

func (c *Counter) GetAndReset() int {
	c.mu.Lock()
	n := c.n
	c.n = 0
	c.mu.Unlock()
	return n
}
