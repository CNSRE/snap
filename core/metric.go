package core

import (
	"time"

	"github.com/intelsdi-x/pulse/core/cdata"
)

type Metric interface {
	Version() int
	Namespace() []string
	LastAdvertisedTime() time.Time
	Config() *cdata.ConfigDataNode
	Data() interface{}
}
