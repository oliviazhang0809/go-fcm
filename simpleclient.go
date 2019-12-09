package gcm

import (
	"fmt"
	"sync"
	"time"
)

// gcmSimpleClient is a container for http clients.
type gcmSimpleClient struct {
	sync.RWMutex
	mh      MessageHandler
	cerr    chan error
	sandbox bool
	debug   bool
	// Clients.
	httpClient httpC
	// GCM auth.
	senderID string
	apiKey   string
	// XMPP config.
	pingInterval time.Duration
	pingTimeout  time.Duration
}

// gcmSimpleClient can only be created through client creation with XMPPClientEnabled disabled
func newSimpleClient(config *Config, h MessageHandler) (Client, error) {
	return newSimpleGCMClient(newHTTPClient(config.APIKey, config.Debug), config, h)
}

// ID returns client unique identification.
func (c *gcmSimpleClient) ID() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("%p", c.httpClient)
}

// SendHTTP sends a message using the HTTP GCM connection server (blocking).
func (c *gcmSimpleClient) SendHTTP(m HTTPMessage) (*HTTPResponse, error) {
	return c.httpClient.Send(m)
}

// newGCMClient creates an instance of gcmClient.
func newSimpleGCMClient(httpc httpC, config *Config, h MessageHandler) (*gcmClient, error) {
	c := &gcmClient{
		httpClient:   httpc,
		cerr:         make(chan error, 1),
		senderID:     config.SenderID,
		apiKey:       config.APIKey,
		mh:           h,
		debug:        config.Debug,
		sandbox:      config.Sandbox,
		pingInterval: time.Duration(config.PingInterval) * time.Second,
		pingTimeout:  time.Duration(config.PingTimeout) * time.Second,
	}
	if c.pingInterval <= 0 {
		c.pingInterval = DefaultPingInterval
	}
	if c.pingTimeout <= 0 {
		c.pingTimeout = DefaultPingTimeout
	}

	// Wait a bit to see if the newly created client is ok.
	// TODO: find a better way (success notification, etc).
	select {
	case err := <-c.cerr:
		return nil, err
	case <-time.After(time.Second): // TODO: configurable
		// Looks good.
	}

	return c, nil

}
