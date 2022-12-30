package trips

import (
	"fmt"
	"io"
	"log"
)

type Server struct {
	UnimplementedTripsServer
}

func (srv *Server) Command(tsrv Trips_CommandServer) error {

	ctx := tsrv.Context()

	for {

		// exit if context is done, or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := tsrv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			// log receive error
			continue
		}

		rType := req.GetType()
		rData := req.GetData()
		fmt.Println(rType, rData)

		resp := Response{Type: "", Result: ""}
		if err := tsrv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
	}
}
