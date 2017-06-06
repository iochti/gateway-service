package helpers

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DeclareService is helper which creates a conn with a service
func DeclareService(addr, caCertFile, svcName *string) *grpc.ClientConn {
	var creds grpc.DialOption
	if *caCertFile == "" {
		creds = grpc.WithInsecure()
	} else {
		crds, err := credentials.NewClientTLSFromFile(*caCertFile, *svcName)
		DieIf(err)
		creds = grpc.WithTransportCredentials(crds)
	}
	conn, err := grpc.Dial(*addr, creds)
	DieIf(err)
	return conn
}
