package core

import "github.com/quickchainproject/quickchain/common"

func (c *core) handleFinalCommitted() error {
	logger := c.logger.New("state", c.state)
	logger.Trace("Received a final committed proposal")
	c.startNewRound(common.Big0)
	return nil
}
