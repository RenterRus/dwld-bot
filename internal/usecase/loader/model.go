package loader

import "context"

type Loader interface {
	Processor(ctx context.Context)
	Stop()
	ForceDelete(string, string)
}
