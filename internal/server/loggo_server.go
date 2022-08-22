/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package server

import (
	"fmt"
	"time"

	"github.com/aurc/loggo/protobuf/pb"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *server) LogStream(req *pb.LogEntryRequest, stream pb.Loggo_LogStreamServer) error {
	startFrom := int64(0)
	if req.FromPosition != nil {
		startFrom = *req.FromPosition
	}
	curr := int(startFrom)
	for {
		if len(s.inSlice) > curr {
			v, err := structpb.NewStruct(s.inSlice[curr])
			if err != nil {
				fmt.Println("fail to build struct ", err)
				return err
			}
			if err := stream.Send(&pb.LogEntryResponse{
				Position: int64(curr),
				Entry:    v,
			}); err != nil {
				fmt.Println("fail to dispatch message ", err)
				return err
			}
			curr++
		} else {
			time.Sleep(200 * time.Millisecond)
		}
	}
	return nil
}
