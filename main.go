package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pubsub/dht22"
)

const (
	ThingName  = "xxxxxxxxxxxxxxxxxx"
	RootCAFile = "AmazonRootCA1.pem"
	CertFile   = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-certificate.pem.crt"
	KeyFile    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-private.pem.key"
	PubTopic   = "topic/to/publish"
	endpoint   = "xxxxxxxxxxxxxxxxx-xxx.iot.ap-northeast-1.amazonaws.com"
	QoS        = 1
)

func main() {

	tlsConfig, err := newTLSConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to construct tls config: %v", err))
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", endpoint, 443))
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(ThingName)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("failed to connect broker: %v", token.Error()))
	}
	for {
		var mydht dht22.MyDHT22
		PubMsg, _ := json.Marshal(mydht.Read())

		log.Printf("publishing %s...\n", PubTopic)
		if token := client.Publish(PubTopic, QoS, false, PubMsg); token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("failed to publish %s: %v", PubTopic, token.Error()))
		}

		time.Sleep(60 * 60 * 2 * time.Second)
	}
}

func newTLSConfig() (*tls.Config, error) {
	rootCA, err := ioutil.ReadFile(RootCAFile)
	if err != nil {
		return nil, err
	}
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(rootCA)
	cert, err := tls.LoadX509KeyPair(CertFile, KeyFile)
	if err != nil {
		return nil, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:            certpool,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		NextProtos:         []string{"x-amzn-mqtt-ca"}, // Port 443 ALPN
	}, nil
}
