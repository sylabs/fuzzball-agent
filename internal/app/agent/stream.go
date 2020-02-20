// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// stream allows for IO streaming over a NATS connection.
type stream struct {
	id string
	nc *nats.Conn
}

func (s stream) Write(b []byte) (n int, err error) {
	if err = s.nc.Publish(fmt.Sprintf("job.%v.output", s.id), b); err != nil {
		return 0, err
	}
	return len(b), nil
}
