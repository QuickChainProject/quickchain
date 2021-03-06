package core

import (
	bft "github.com/quickchainproject/quickchain/consensus/dbft"
)

type backlogEvent struct {
	src bft.Validator
	msg *message
}

type timeoutEvent struct{}
