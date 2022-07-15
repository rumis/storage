package skafka

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/rumis/storage/pkg/ujson"
	"github.com/segmentio/kafka-go"
)

func TestReader(t *testing.T) {
	addr := ""
	username := ""
	password := ""
	ca := ""

	InitDefaultDialer(D_WithCA(ca), D_WithUserNamePassword(username, password))
	r, c := NewReaderChannel(
		R_WithBrokers(strings.Split(addr, ",")),
		R_WithDialer(DefaultDialer()),
		R_WithGroupID("GID_fz_cservice_thumbnail_page_screenshot"),
		R_WithTopic("topic_fz_cservice_thumbnail_page_task"))

	defer c()
	for {
		msg := <-r
		fmt.Println(msg)
	}

}

type CoursewareChanged struct {
	ResID      int64  `json:"resId"`
	Tag        string `json:"tag"`
	SlideURL   string `json:"slideUrl"`
	Force      int    `json:"force"`
	CreateID   int    `json:"createId"`
	CreateName string `json:"createName"`
}

func TestWriter(t *testing.T) {

	addr := ""
	username := ""
	password := ""
	ca := ""

	InitDefaultDialer(D_WithCA(ca), D_WithUserNamePassword(username, password))
	w1, c := NewWriter1(
		W_WithBrokers(strings.Split(addr, ",")),
		W_WithDialer(DefaultDialer()),
		W_WithTopic("topic_fz_cservice_courseware_content_changed"))

	defer c()

	val, err := ujson.Marshal(CoursewareChanged{
		ResID:      202278586147668300,
		Tag:        "4ce4fc3dd7f0aa556e146d44898d0595",
		SlideURL:   "https://file1-fz.jiaoyanyun.com/course/202278586147668300/pro/2022785861476683001657844832155.json",
		Force:      1,
		CreateID:   100010,
		CreateName: "mu",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = w1(context.TODO(), kafka.Message{
		Value: val,
	})
	if err != nil {
		t.Fatal(err)
	}
}
