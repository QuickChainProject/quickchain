package core

import "github.com/quickchainproject/quickchain/consensus/bft"

func (c *core) handleRequest(request *bft.Request) error {
	logger := c.logger.New("state", c.state, "seq", c.current.sequence)

	if err := c.checkRequestMsg(request); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid request")
			return err
		}
		logger.Warn("unexpected request", "err", err, "number", request.Proposal.Number(), "hash", request.Proposal.Hash())
		return err
	}

	logger.Trace("handleRequest", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.current.pendingRequest = request
	if c.state == StateAcceptRequest {
		c.sendPreprepare(request)
	}
	return nil
}

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
func (c *core) checkRequestMsg(request *bft.Request) error {
	if request == nil || request.Proposal == nil {
		return errInvalidMessage
	}

	if c := c.current.sequence.Cmp(request.Proposal.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(request *bft.Request) {
	logger := c.logger.New("state", c.state)

	logger.Trace("Store future request", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(request, float32(-request.Proposal.Number().Int64()))
}

func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()
		r, ok := m.(*bft.Request)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkRequestMsg(r)
		if err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())
				c.pendingRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash(), "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())

		go c.sendEvent(bft.RequestEvent{
			Proposal: r.Proposal,
		})
	}
}
