package utils

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// DumpIncomingContext :nodoc:
func DumpIncomingContext(c context.Context) string {
	md, _ := metadata.FromIncomingContext(c)
	return Dump(md)
}
