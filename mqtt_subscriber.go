package main

import (
	"fmt"
	"log"
	"net/url"
	"time"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	fmt.Println("Client : ", client)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
		fmt.Println("waiting ...")
	}
	if err := token.Error; err != nil {
		// fmt.Println("client connected 5 Error")
		// log.Fatal(err)
		fmt.Println(token)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("wss://%s", uri.Host))
	fmt.Println("my host: ", uri.Host)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	fmt.Println("my username: ", uri.User.Username())
	fmt.Println("my pwd: ", password)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	fmt.Println("opts : ", opts)
	return opts
}

func listen(uri *url.URL, topic string) {
	client := connect("sub", uri)
	fmt.Println("client connected 2")
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s]:[%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

func main() {
	mqtt_uri, err := url.Parse(os.Getenv("MQTT_URL"))
	if err != nil {
		log.Fatal(err)
	}
	// topic := mqtt_uri.Path[1:len(mqtt_uri.Path)]
	topic := "owntracks/#"

	fmt.Println("client connected 1")
	fmt.Println("TOPIC is : ", topic)
	go listen(mqtt_uri, topic)
	fmt.Println("client connected post 1")
	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
	}
}