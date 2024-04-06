package ethrlp

import (
	iq "github.com/madz-lab/insertion-queue"
)

// workerJob is a single routine job performed
// by the worker thread
type workerJob struct {
	sourceData []byte // the reference to the item that need to be parsed
	storeIndex int    // the final store index after finalization
}

type decodeWorkerPool struct {
	resultsCh chan *workerDecodeResult // The channel to relay results to
}

// workerDecodeResult is the result of the worker thread's decoding job
type workerDecodeResult struct {
	error        error
	decodedValue Value
	storeIndex   int // the final store index after decoding
}

func (wr *workerDecodeResult) Less(other iq.Item) bool {
	return wr.storeIndex < other.(*workerDecodeResult).storeIndex
}

// newWorkerPool spawns a new worker pool
func newDecodeWorkerPool() *decodeWorkerPool {
	return &decodeWorkerPool{
		resultsCh: make(chan *workerDecodeResult),
	}
}

// addJob adds a new job asynchronously to be processed by the worker pool
func (wp *decodeWorkerPool) addJob(job *workerJob) {
	go wp.runJob(job)
}

// getResult takes out a result from the worker pool [BLOCKING]
func (wp *decodeWorkerPool) getResult() *workerDecodeResult {
	return <-wp.resultsCh
}

// close closes the worker pool and their corresponding channels
func (wp *decodeWorkerPool) close() {
	close(wp.resultsCh)
}

// runJob is the main activity method for the
// worker threads where new jobs are parsed and results sent out
func (wp *decodeWorkerPool) runJob(job *workerJob) {
	decodedValue, decodeErr := DecodeBytes(job.sourceData)

	result := &workerDecodeResult{
		storeIndex:   job.storeIndex,
		decodedValue: decodedValue,
		error:        decodeErr,
	}

	// Report the result back
	wp.resultsCh <- result
}
