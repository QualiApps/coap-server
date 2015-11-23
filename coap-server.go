// this is a coap server
// that listening coap requests and allows us
// to interact with DB

package main

import (
	"fmt"
	coap "github.com/dustin/go-coap"
	"github.com/qualiapps/coap-server/handlers"
	"log"
	"net"
	"strings"
)

const (
	protocol = "udp"
	port     = "5683"
)

var endpoints = []string{"api-db"}

func contains(array *[]string, value string) bool {
	exist := false
	value = strings.ToLower(value)
	for _, ep := range *array {
		if ep == value {
			exist = true
			break
		}
	}

	return exist
}

func check_endpoint(end_point []string) bool {
	path := strings.Join(end_point, "/")
	return contains(&endpoints, path)
}

func response(m *coap.Message, response_code coap.COAPCode, payload *[]byte) *coap.Message {
	res := &coap.Message{
		Type:      coap.Acknowledgement,
		Code:      response_code,
		MessageID: m.MessageID,
		Token:     m.Token,
		Payload:   *payload,
	}
	fmt.Println("Send %#v", res)
	res.SetOption(coap.ContentFormat, coap.TextPlain)

	log.Printf("Transmitting %#v", res)
	return res

}

func handle(conn *net.UDPConn, addr *net.UDPAddr, mess *coap.Message) *coap.Message {
	log.Printf("Got message path=%q: %#v from %v", mess.Path(), mess, addr)

	var payload []byte
	var err error

	// checks available endpoints
	check := check_endpoint(mess.Path())
	if !check {
		return response(mess, coap.BadRequest, &payload)
	}

	// invoke handler
	payload, err = handlers.SendRequest(mess)
	if err != nil {
		log.Printf("Error on handler, stopping: %v", err)
		return nil
	}

	if mess.IsConfirmable() {
		return response(mess, coap.Content, &payload)
	}

	return nil
}

func main() {
	log.Fatal(
		coap.ListenAndServe(protocol, ":"+port, coap.FuncHandler(handle)),
	)
}
