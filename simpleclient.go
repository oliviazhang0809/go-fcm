package fcm

import (
	"fmt"
	"sync"
	"time"
)

// fcmSimpleClient is a container for http clients.
type fcmSimpleClient struct {
	sync.RWMutex
	cerr    chan error
	sandbox bool
	debug   bool
	// Clients.
	httpClient httpC
	// FCM auth.
	senderID string
	apiKey   string
}

// fcmSimpleClient can only be created through client creation with XMPPClientEnabled disabled
func newSimpleClient(config *Config, h MessageHandler) (Client, error) {
	return newSimpleFCMClient(
		newHTTPClient(
			config.APIKey,
			config.Debug,
			config.HTTPTimeout,
		),
		config,
		h,
	)
}

// ID returns client unique identification.
func (c *fcmSimpleClient) ID() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("%p", c.httpClient)
}

// SendHTTP sends a message using the HTTP FCM connection server (blocking).
func (c *fcmSimpleClient) SendHTTP(m HTTPMessage) (*HTTPResponse, error) {
	return c.httpClient.Send(m)
}

func (c *fcmSimpleClient) SendXMPP(m XMPPMessage) (string, int, error) {
	return "", 0, nil
}

// Close will stop and close the corresponding client, releasing all resources (blocking).
func (c *fcmSimpleClient) Close() error {
	return nil
}

// newFCMClient creates an instance of fcmClient.
func newSimpleFCMClient(httpc httpC, config *Config, h MessageHandler) (*fcmSimpleClient, error) {
	c := &fcmSimpleClient{
		httpClient: httpc,
		cerr:       make(chan error, 1),
		senderID:   config.SenderID,
		apiKey:     config.APIKey,
		debug:      config.Debug,
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
