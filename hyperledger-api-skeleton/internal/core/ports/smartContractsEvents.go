package ports

import (
	"os"
)

type InputEvents interface {
	StartListenEvents(signal chan os.Signal)
}

type OutputEvents interface {
}
