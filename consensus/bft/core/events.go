package core

import (
	"github.com/quickchainproject/quickchain/consensus/bft"
)

type backlogEvent struct {
	src bft.Validator
	msg *message
}

type timeoutEvent struct{}
