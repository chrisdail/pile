package gitver

import "sync"

type cachedStringResponse struct {
	sync.Once
	response string
	err      error
}
