package main

import (
	"context"
	"io"

	logging "cloud.google.com/go/logging/apiv2"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

func main() {
	ctx := context.Background()
	c, err := logging.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()
	stream, err := c.TailLogEntries(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	go func() {
		reqs := []*loggingpb.TailLogEntriesRequest{
			{
				ResourceNames: nil,
				Filter:        "",
				BufferWindow:  nil,
			},
		}
		for _, req := range reqs {
			if err := stream.Send(req); err != nil {
				// TODO: Handle error.
			}
		}
		stream.CloseSend()
	}()
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// TODO: handle error.
		}
		// TODO: Use resp.
		_ = resp
	}
}
