package types

import (
	"sync/atomic"
	"time"
	"fmt"
)

type Guid struct {
	Guid string
}

func (g *Guid) Equal(o Guid) bool {
	return g.Guid == o.Guid
}

func (g *Guid) IsValid() bool {
	return g.Guid != "" && g.Guid != "0" && g.Guid != "0x0"
}

func GUID(guid string)(g Guid) {
	g.Guid = guid
	return
}

func (g Guid) String() string{
	return g.Guid
}

func NewGuid() (g Guid) {
	currentTime := time.Now().UTC().Unix()
	i := atomic.AddUint32(&objectIDCounter, 1) & 0x3FFFFF

	g.Guid = fmt.Sprintf("0x%x", uint64(currentTime)<<36 | uint64(i)<<14 | uint64(nodeSN))

	return
}

var nodeSN int = 6399748

var objectIDCounter uint32 = 0