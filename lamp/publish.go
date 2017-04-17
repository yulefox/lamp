package lamp

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// NSQPublish publish message to given `topic`
func NSQPublish(addr string, topic string, msg []byte) (err error) {
	url := fmt.Sprintf("http://%s/pub?topic=%s", addr, topic)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(msg))
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	res.Body.Close()
	return
}
