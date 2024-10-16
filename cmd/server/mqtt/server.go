package mqtt_server

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/klaital/tv-controller/internal/config"
	"github.com/klaital/tv-controller/vlcclient"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MqttServer struct {
	Host string
	Port int

	PublishConfigTopic  string // used by me to publish config settings whenever they change (or when requested)
	RequestConfigTopic  string // used by others to prompt for a re-publication of the config settings
	ChangePlaylistTopic string // used by others to request a new playlist be loaded

	Config   *config.Config
	Vlc      *vlcclient.Client
	mClient  mqtt.Client
	shutdown bool
}

func (s *MqttServer) Start() {
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", s.Host, s.Port)).
		SetClientID("tv-controller").
		SetOrderMatters(false).
		SetDefaultPublishHandler(s.Handler).
		SetConnectRetry(true).
		SetConnectRetryInterval(30 * time.Second).
		SetOnConnectHandler(s.OnConnect).
		SetConnectionLostHandler(s.OnConnectionLost)
	client := mqtt.NewClient(opts)
	//if token := client.Subscribe(s.ChangePlaylistTopic, 1, nil); token.Wait() && token.Error() != nil {
	//	slog.Error("Failed to subscribe", "topic", s.ChangePlaylistTopic, "err", token.Error())
	//} else {
	//	slog.Info("Subscribed to topic", "topic", s.ChangePlaylistTopic)
	//}
	//
	//if token := client.Subscribe(s.RequestConfigTopic, 1, nil); token.Wait() && token.Error() != nil {
	//	slog.Error("Failed to subscribe", "topic", s.RequestConfigTopic, "err", token.Error())
	//} else {
	//	slog.Info("Subscribed to topic", "topic", s.RequestConfigTopic)
	//}

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("Failed to connect to mqtt server", "err", token.Error())
		panic(token.Error())
	}
	s.mClient = client

	// TODO: move signal handling up to main
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	slog.Info("Shutting down MQTT server")
	s.shutdown = true
	client.Disconnect(250)
}
func (s *MqttServer) OnConnect(client mqtt.Client) {
	opts := client.OptionsReader()
	for _, url := range opts.Servers() {
		slog.Debug("MQTT Connected to", "server", url.String())
	}
	if token := client.Subscribe(s.ChangePlaylistTopic, 1, nil); token.Wait() && token.Error() != nil {
		slog.Error("Failed to subscribe", "topic", s.ChangePlaylistTopic, "err", token.Error())
	} else {
		slog.Info("Subscribed to topic", "topic", s.ChangePlaylistTopic)
	}

	if token := client.Subscribe(s.RequestConfigTopic, 1, nil); token.Wait() && token.Error() != nil {
		slog.Error("Failed to subscribe", "topic", s.RequestConfigTopic, "err", token.Error())
	} else {
		slog.Info("Subscribed to topic", "topic", s.RequestConfigTopic)
	}
}
func (s *MqttServer) OnConnectionLost(client mqtt.Client, err error) {
	slog.Error("MQTT Connection lost", "err", err)
	// Automatically reconnect unless the disconnect was on purpose.
	if !s.shutdown {
		go func() {
			time.Sleep(5 * time.Second)
			token := client.Connect()
			if err := token.Error(); err != nil {
				slog.Error("Failed to reconnect to MQTT broker", "err", err)
			}
		}()
	}
}

func (s *MqttServer) Handler(client mqtt.Client, msg mqtt.Message) {
	switch msg.Topic() {
	case s.PublishConfigTopic:
		s.PublishConfig()
	case s.RequestConfigTopic:
		s.RequestConfig()
	case s.ChangePlaylistTopic:
		s.ChangePlaylist(string(msg.Payload()))
	}
}

func (s *MqttServer) PublishConfig() {
	slog.Debug("Publishing updated config")
	s.mClient.Publish(s.PublishConfigTopic, 0, false, s.Config.ToString())
}
func (s *MqttServer) RequestConfig() {
	slog.Debug("New config requested")
	s.PublishConfig()
}
func (s *MqttServer) ChangePlaylist(playlist string) {
	slog.Debug("Changing playlist to " + playlist)
	// Tell VLC to start playing the specified playlist
	s.Config.SelectedPlaylist = playlist
	err := s.Config.StopVlc()
	if err != nil {
		slog.Error("Failed to stop existing VLC instance", "error", err)
	}
	err = s.Config.StartVlc()
	if err != nil {
		slog.Error("Failed to start new VLC instance", "error", err)
	}

	// save the setting change
	config.SaveConfig(s.Config)
}
