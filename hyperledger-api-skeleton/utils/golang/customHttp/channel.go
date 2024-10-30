package customHttp

import (
	"bytes"
	"fmt"
	"repo.plexus.services/1329-004_incibe_reto06/utils/golang/customHttp/concurrent"
	"sync"
	"time"
)

// ChannelRequests main channel structure
type ChannelRequests struct {
	LimitReq            concurrent.LimitRequests `json:"limitRequests"`
	CacheReq            concurrent.CacheRequests `json:"cacheRequests"`
	BackPressureOffline *concurrent.BPOffline    `json:"bPOffline"`
	metrics             *concurrent.MetricsChannel
	CurrentCache        *concurrent.SafeMapCacheRes `json:"currentCache"`
	DontWaitRes         bool                        `json:"dontWaitRes,omitempty"`
	LimitThreads        concurrent.LimitThreads     `json:"limitThreads,omitempty"`
	resultsChan         chan concurrent.ResponseChannel
	totalThreads        int
	chanForEachReq      bool
}

// Default idempotent function
func isIdempotentRec(r concurrent.RequestsChannel) bool {
	if r.TypeR == "GET" {
		return true
	}
	return false
}

// ChannelsRequests Use of channels to process an array of requests
func (c *ChannelRequests) ChannelsRequests(p HttpInterface, jobs chan concurrent.RequestsChannel, results chan concurrent.ResponseChannel, threads int) {
	c.resultsChan = results
	c.totalThreads = threads
	c.prepareChannelReqs(p, jobs)
}

// ChannelsRequestNoResult Use of channels to process an array of requests. Each request should provide each own result channel
func (c *ChannelRequests) ChannelsRequestNoResult(p HttpInterface, jobs chan concurrent.RequestsChannel, threads int) {
	c.totalThreads = threads
	c.chanForEachReq = true
	c.prepareChannelReqs(p, jobs)
}

// processChannelReqs Use of channels to process an array of requests. Main functionality
func (c *ChannelRequests) prepareChannelReqs(p HttpInterface, jobs chan concurrent.RequestsChannel) {
	burstLimiter := make(chan time.Time, c.LimitReq.TotalTx)
	wg := sync.WaitGroup{}

	//Initialize channel to limit total request active
	c.limiterThreads(&wg)

	var cacheRequests concurrent.SafeMapCacheRes
	cacheRequests.Init()
	jobs1 := jobs
	c.metrics = &concurrent.MetricsChannel{}

	//Initalize back preasure offline
	if c.BackPressureOffline != nil && c.BackPressureOffline.Activated && c.BackPressureOffline.Capacity > 0 {
		c.BackPressureOffline.BPChannel.SignalCh = make(chan bool, c.BackPressureOffline.Capacity)
	}

	//Cache duplicate(idempotent) requests. Optional
	if c.CacheReq.Activated {
		jobs1 = make(chan concurrent.RequestsChannel, 2)
		if c.CacheReq.CacheF == nil {
			c.CacheReq.CacheF = isIdempotentRec
		}
		c.CurrentCache = &cacheRequests
		go c.checkForCachedRequests(jobs, jobs1, c.resultsChan, &cacheRequests, &wg)
	}
	//Main functionality
	for i := 0; i < c.totalThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for element := range jobs1 {
				if c.LimitReq.Activated == true {
					<-burstLimiter
				}
				c.channelsRequests(p, element, c.resultsChan, &cacheRequests)
			}
		}()
	}

	//Limit requests Transactions per time. Optional
	c.limiterFunction(burstLimiter, &wg)
}

func (c *ChannelRequests) channelsRequests(p HttpInterface, r concurrent.RequestsChannel, results chan concurrent.ResponseChannel, cacheRequests *concurrent.SafeMapCacheRes) {
	//When offline back pressure activated, then wait for requests to complete.
	if c.BackPressureOffline != nil && c.BackPressureOffline.Activated {
		c.BackPressureOffline.BPChannel.SignalCh <- true
	}
	if c.LimitThreads.Activated {
		<-c.LimitThreads.ThreadReqLimiter
	}
	if c.DontWaitRes {
		go c.doHttpRequests(p, r, results, cacheRequests)
	} else {
		c.doHttpRequests(p, r, results, cacheRequests)
	}
}

func (c *ChannelRequests) doHttpRequests(p HttpInterface, r concurrent.RequestsChannel, results chan concurrent.ResponseChannel, cacheRequests *concurrent.SafeMapCacheRes) {
	var resp concurrent.ResponseChannel
	resp.Request = r.Request

	ByteBody := bytes.NewReader(r.Body)
	switch r.TypeR {
	case "GET":
		resp.Response, resp.Error = p.GetCustomMethod(r.Ctx, r.Url, ByteBody)
	case "GET_OPTIONS":
		resp.Response, resp.Error = p.PostCustomMethodOptions(r.Ctx, r.Url, ByteBody, r.Opt)
	case "POST":
		resp.Response, resp.Error = p.PostCustomMethod(r.Ctx, r.Url, ByteBody)
	case "POST_AUTH":
		resp.Response, resp.Error = p.PostCustomMethodAuth(r.Ctx, r.Url, ByteBody, r.BasicAuthValue)
	case "POST_OPTIONS":
		resp.Response, resp.Error = p.PostCustomMethodOptions(r.Ctx, r.Url, ByteBody, r.Opt)
	case "PUT":
		resp.Response, resp.Error = p.PutCustomMethod(r.Ctx, r.Url, ByteBody)
	case "PUT_AUTH":
		resp.Response, resp.Error = p.PutCustomMethodAuth(r.Ctx, r.Url, ByteBody, r.BasicAuthValue)
	case "PUT_OPTIONS":
		resp.Response, resp.Error = p.PostCustomMethodOptions(r.Ctx, r.Url, ByteBody, r.Opt)
	case "DELETE":
		resp.Response, resp.Error = p.DeleteCustomMethod(r.Ctx, r.Url, ByteBody)
	}
	resp.Code = 200
	resp.Ctx = r.Ctx

	//Count number of txs(reset if limiter activated)
	c.metrics.CompleteTxs.Add()

	//Load again to threadReqLimiter
	if c.LimitThreads.Activated {
		c.LimitThreads.ThreadReqLimiter <- true
	}

	//Check if cache requests functionality is implemented and if request is type cacheF
	if c.CacheReq.Activated && c.CacheReq.CacheF(r) {
		//Check if response exists on map of cached requests
		key := r.TypeR + r.Url + string(r.Body)
		if _, ok := cacheRequests.Get(key); !ok {
			cResp := concurrent.CacheResponse{Resp: resp, TimeCached: r.CacheTimeMs}
			cacheRequests.Set(key, cResp)
		}
	}

	if r.ResultChan == nil && !c.chanForEachReq {
		results <- resp
	} else if r.ResultChan != nil {
		r.ResultChan <- resp
	} else {
		fmt.Println("Error: The response is not being sent through any channel due to a misuse of the library.")
	}
}

func (c *ChannelRequests) limiterThreads(wg *sync.WaitGroup) {
	if c.LimitThreads.Activated {
		c.LimitThreads.ThreadReqLimiter = make(chan bool, c.LimitThreads.TotalReqActive)
		for i := 0; i <= c.LimitThreads.TotalReqActive-1; i++ {
			c.LimitThreads.ThreadReqLimiter <- true
		}
		go func() {
			//Wait for threads to be done
			(*wg).Wait()
			close(c.LimitThreads.ThreadReqLimiter)
		}()
	}
}

func (c *ChannelRequests) limiterFunction(burstLimiter chan time.Time, wg *sync.WaitGroup) {
	if c.LimitReq.Activated == true {
		quit := make(chan bool)
		capacityChannel := cap(burstLimiter)

		//Load first transactions per second
		for i := 0; i < c.LimitReq.TotalTx; i++ {
			burstLimiter <- time.Now()
		}
		//Define time channel to control time per requests
		go func() {
			for t := range time.Tick(c.LimitReq.TimeMs * time.Millisecond) {
				select {
				case <-quit:
					close(burstLimiter)
					return
				default:
					var length int
					if c.LimitReq.AutoBackPressure {
						length = c.metrics.CompleteTxs.GetAndReset()
					} else {
						length = capacityChannel - len(burstLimiter)
					}
					for i := 0; i < length; i++ {
						burstLimiter <- t
					}
				}
			}
		}()
		go func() {
			//Wait for threads to be done
			(*wg).Wait()
			//Check burstLimiter to see if is blocked
			if len(burstLimiter) != 0 {
				for i := 0; i < len(burstLimiter); i++ {
					<-burstLimiter
				}
			}
			//Close limiter go routine
			quit <- true
			//Close quit channel
			close(quit)
		}()
	} else {
		close(burstLimiter)
	}
}

func (c *ChannelRequests) checkForCachedRequests(jobs2 chan concurrent.RequestsChannel, jobs1 chan concurrent.RequestsChannel, results chan concurrent.ResponseChannel, cacheRequests *concurrent.SafeMapCacheRes, wg *sync.WaitGroup) {
	var ticker *time.Ticker
	ticker = time.NewTicker(c.CacheReq.TimeMs * time.Millisecond)

	//Receive jobs to me treated and check for already cached requests
	go func() {
		for req := range jobs2 {
			key := req.TypeR + req.Url + string(req.Body)
			if r, ok := cacheRequests.Get(key); ok && c.CacheReq.CacheF(req) {
				r.Resp.Request = req.Request
				r.Resp.Ctx = req.Ctx

				if req.ResultChan == nil && !c.chanForEachReq {
					results <- r.Resp
				} else if req.ResultChan != nil {
					req.ResultChan <- r.Resp
				} else {
					fmt.Println("error: misuse of library: no response channel")
				}

				c.metrics.TotalReqCache++ //metric
			} else {
				jobs1 <- req
			}
		}
	}()

	//Break ticker when all process end
	go func() {
		(*wg).Wait()
		ticker.Stop()
	}()

	//Clear Cache for each unit of time
	for range ticker.C {
		//Send all requests to job_1 channel
		cacheRequests.DeleteOrUpdate(c.CacheReq.TimeMs)
	}
}

// GetMetrics Obtain metrics of channel
func (c *ChannelRequests) GetMetrics() *concurrent.MetricsChannel {
	return c.metrics
}

// FlushCache to clean cache
func (c *ChannelRequests) FlushCache() *ChannelRequests {
	c.CurrentCache.Flush()
	return c
}
