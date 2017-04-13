package metadata

import (
	"net"
	"strings"
	"time"

	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ContentService interface {
	GetContent(source string) chan Content
}

type UPPContentService struct {
	dbName  string
	session *mgo.Session
}

func tcpDialServer(addr *mgo.ServerAddr) (net.Conn, error) {
	ra, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, ra)
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(30 * time.Second)
	return conn, nil
}

func InitContentService(delivery *Cluster) (*UPPContentService, error) {
	info := mgo.DialInfo{
		Timeout:    2 * time.Minute,
		Addrs:      strings.Split(delivery.address, ","),
		DialServer: tcpDialServer,
	}

	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Strong, true)
	session.SetCursorTimeout(0)
	return &UPPContentService{session: session}, nil
}

func (c *UPPContentService) GetContent(source string) chan Content {
	result := make(chan Content)

	go func() {
		defer close(result)
		coll := c.session.DB("upp-store").C("content")
		iter := coll.Find(bson.M{"mediaType": nil}).Select(bson.M{"uuid": true, "_id": false, "identifiers.authority": true}).Iter()

		var content Content
		var count int
		for iter.Next(&content) {
			cSource, ok := content.getSource()
			if ok && source == cSource {
				count++
				result <- content
			}
		}
		fmt.Printf("Read %d content items\n", count)
	}()

	return result
}
