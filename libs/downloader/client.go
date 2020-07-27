package downloader

type Worker func()

type Client struct {
	message chan Worker
	maxSize int
}

func New(maxParallelWorkerSize int) *Client {
	return &Client{
		message: make(chan Worker, maxParallelWorkerSize),
		maxSize: maxParallelWorkerSize,
	}
}

func (c *Client) Size() int {
	return len(c.message)
}

func (c *Client) Start() {
	for i := 0; i < c.maxSize; i++ {
		go func() {
			for task := range c.message {
				task()
			}
		}()
	}
}

func (c *Client) PutTask(t Worker) {
	c.message <- t
}

func (c *Client) ShutDown() {
	close(c.message)
}
