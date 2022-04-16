package main

import (
	"encoding/json"
	"github.com/antchfx/xmlquery"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

// cb rf daily quotes url
const url = "https://www.cbr.ru/scripts/XML_daily.asp"

type quotesMsg struct {
	Timestamp string `json:"timestamp"`
	Date      string `json:"date"`
	USD       string `json:"usd"` // dollar exchange rate
	EUR       string `json:"eur"` // euro exchange rate
}

const (
	streamName     = "QUOTES"
	streamSubjects = "QUOTES.*"
	subjectName    = "QUOTES.publish"
)

func createStream(js nats.JetStreamContext) error {
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
			MaxMsgs:  100, // limited just fo fun
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func publishQuotes(js nats.JetStreamContext, msg quotesMsg) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = js.Publish(subjectName, b)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("Pusher"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	err = createStream(js)

	// publish quotes to a stream in a loop
	for {
		doc, err := xmlquery.LoadURL(url)
		if err != nil {
			log.Println(err)
			continue // just skip
		}

		eurNode := xmlquery.Find(doc, "//Valute[@ID='R01239']") // search eur
		eurValue := xmlquery.Find(eurNode[0], "//Value")        // search eur value

		usdNode := xmlquery.Find(doc, "//Valute[@ID='R01235']") // search usd
		usdValue := xmlquery.Find(usdNode[0], "//Value")        // search usd value

		res := quotesMsg{
			Timestamp: time.Now().Format("20060102150405"),
			Date:      doc.FirstChild.NextSibling.Attr[0].Value,
			USD:       usdValue[0].FirstChild.Data,
			EUR:       eurValue[0].FirstChild.Data,
		}

		// send message to stream
		if err = publishQuotes(js, res); err != nil {
			log.Println(err)
		}

		time.Sleep(5 * time.Second)
	}

}
