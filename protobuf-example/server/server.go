package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	pb "github.com/shufanhao/go-example/protobuf-example/service"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		myClient := pb.Client{}
		data, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println(err)
		}

		if err := proto.Unmarshal(data, &myClient); err != nil {
			fmt.Println(err)
		}

		println(myClient.Id, ":", myClient.Name, ":", myClient.Email, ":", myClient.Country)

		for _, mail := range myClient.Inbox {
			fmt.Println(mail.RemoteEmail, ":", mail.Body)
		}
	})
	fmt.Println("Listen on port :3000")
	http.ListenAndServe(":3000", nil)
}
