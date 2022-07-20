package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	logging "cloud.google.com/go/logging/apiv2"
	rm "cloud.google.com/go/resourcemanager/apiv3"
	"golang.org/x/oauth2/google"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

var scopes = []string{"https://www.googleapis.com/auth/logging.read",
	"https://www.googleapis.com/auth/cloud-platform.read-only"}

func main() {
	ctx := context.Background()
	rm.NewProjectsClient(ctx)
	_, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		cmd := exec.Command("gcloud", "auth", "application-default", "login")
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	}

	c, err := logging.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	lastHour := time.Now().Add(-1 * time.Minute).Format(time.RFC3339)
	it := c.ListLogEntries(ctx, &loggingpb.ListLogEntriesRequest{
		ResourceNames: []string{"projects/anz-x-apps-np-e1bb39"},
		Filter: `resource.labels.namespace_name="gearbox-sit" AND ` +
			`resource.labels.container_name="credit" AND ` +
			`NOT textPayload:"ignoring custom sampler for span" AND ` +
			fmt.Sprintf(`timestamp > "%s"`, lastHour),
		PageSize: 1,
	})
	var resp *loggingpb.LogEntry
	for {
		resp, err = it.Next()

		b, _ := json.Marshal(resp)
		fmt.Println(string(b))
		//t := resp.GetTextPayload()
		//p := resp.GetJsonPayload()
		//if p != nil {
		//	b, _ := p.MarshalJSON()
		//	fmt.Println(string(b))
		//} else if len(t) > 0 {
		//	fmt.Println(t)
		//}
		fmt.Println("------------------------")
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	_ = resp
	//stream, err := c.TailLogEntries(ctx)
	//if err != nil {
	//	// TODO: Handle error.
	//}
	//go func() {
	//	reqs := []*loggingpb.TailLogEntriesRequest{
	//		{
	//			ResourceNames: []string{"projects/anz-x-apps-np-e1bb39/logs"},
	//			//Filter:        `resource.labels.namespace_name="gearbox-sit"`,
	//			BufferWindow: nil,
	//		},
	//	}
	//	for _, req := range reqs {
	//		if err := stream.Send(req); err != nil {
	//			// TODO: Handle error.
	//			fmt.Println(err)
	//		}
	//	}
	//	stream.CloseSend()
	//}()
	//for {
	//	resp, err := stream.Recv()
	//	//if err == io.EOF {
	//	//	fmt.Println(err)
	//	//	break
	//	//}
	//	if err != nil {
	//		// TODO: handle error.
	//		//fmt.Println(err)
	//		time.Sleep(time.Millisecond * 1500)
	//	}
	//	// TODO: Use resp.
	//	//_ = resp
	//	if err == nil {
	//		fmt.Println(resp.String())
	//	} else {
	//		fmt.Print(".")
	//	}
	//}
}
