package ethrlp

// workerPool is the pool of worker threads
// that parse hashing jobs
type workerPool struct {
	resultsCh chan *workerResult // The channel to relay results to
}

// workerJob is a single encoding job performed
// by the worker thread
type workerJob struct {
	storeIndex int    // the final store index after encoding
	sourceData []byte // the reference to the item that need to be encoded
}

// workerResult is the result of the worker thread's encoding job
type workerResult struct {
	storeIndex  int    // the final store index after encoding
	encodedData []byte // the actual RLP encoded result data
}

// newWorkerPool spawns a new worker pool
func newWorkerPool(expectedNumResults int) *workerPool {
	return &workerPool{
		resultsCh: make(chan *workerResult, expectedNumResults),
	}
}

// addJob adds a new job asynchronously to be processed by the worker pool
func (wp *workerPool) addJob(job *workerJob) {
	go wp.runJob(job)
}

// getResult takes out a result from the worker pool [Blocking]
func (wp *workerPool) getResult() *workerResult {
	return <-wp.resultsCh
}

// close closes the worker pool and their corresponding
// channels
func (wp *workerPool) close() {
	close(wp.resultsCh)
}

// runJob is the main activity method for the
// worker threads where new jobs are parsed and results sent out
func (wp *workerPool) runJob(job *workerJob) {
	// Encode the data to RLP
	encodedBytes := EncodeBytes(job.sourceData)

	// Construct a hash result from the fast hasher
	result := &workerResult{
		storeIndex:  job.storeIndex,
		encodedData: encodedBytes,
	}

	// Report the result back
	wp.resultsCh <- result
}
