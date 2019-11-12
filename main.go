package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	proto "github.com/golang/protobuf/proto"
)

func main() {
	u, err := url.Parse("http://datamine.mta.info/mta_esi.php")
	if err != nil {
		exit("unable to parse MTA URL", err)
	}

	q := u.Query()
	q.Set("key", os.Getenv("MTA_API_KEY"))
	q.Set("feed_id", "1")

	u.RawQuery = q.Encode()

	cli := &http.Client{}

	resp, err := cli.Get(u.String())
	if err != nil {
		exit("unable to get GTFS realtime data", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exit("unable to read response body", err)
	}

	var feed gtfs.FeedMessage
	if err := proto.Unmarshal(body, &feed); err != nil {
		exit("unable to unmarshal GTFS Feed Message", err)
	}

	fmt.Println("Successfully requested GTFS realtime data from MTA API!")
	fmt.Printf(
		"The GTFS realtime version of the response is %s and it contains %d entities.\n",
		feed.Header.GetGtfsRealtimeVersion(), len(feed.Entity),
	)
}

func exit(msg string, err error) {
	fmt.Printf("%s: %v", msg, err)
	os.Exit(1)
}
