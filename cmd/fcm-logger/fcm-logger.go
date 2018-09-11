// Program fcm-logger logs and echoes as a FCM "server".
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/aliafshar/toylog"

	"github.com/lilichen/go-fcm"
)

var (
	serverKey = kingpin.Flag("server_key", "The server key to use for FCM.").Short('k').Required().String()
	senderId  = kingpin.Flag("sender_id", "The sender ID to use for FCM.").Short('s').Required().String()
	g fcm.Client
)

// onMessage receives messages, logs them, and echoes a response.
func onMessage(cm fcm.CCSMessage) error {
	toylog.Infoln("Message, from:", cm.From, "with:", cm.Data)
	// Echo the message with a tag.
	cm.Data["echoed"] = true
	m := fcm.HTTPMessage{To: cm.From, Data: cm.Data}
	r, err := g.SendHTTP(m)
	if err != nil {
		toylog.Errorln("Error sending message.", err)
		return err
	}
	toylog.Infof("Sent message. %+v -> %+v", m, r)
	return nil
}

func main() {
	toylog.Infoln("FCM Logger, starting.")
	kingpin.Parse()

	// Setup signal handler and wait for a terminating signal.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// Init fcm client.
	gconf := &fcm.Config{
		SenderID:          *senderId,
		APIKey:            *serverKey,
		Sandbox:           false,
		MonitorConnection: false,
		Debug:             true,
	}
	var err error
	if g, err = fcm.NewClient(gconf, onMessage); err != nil {
		toylog.Errorf("fcm client init error %s", err.Error())
		return
	}

	<-sig
}
