package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/namsral/flag"
	"github.com/urfave/negroni"

	authpb "github.com/iochti/auth-service/proto"
	"github.com/iochti/gateway-service/handlers"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type IochtiGateway struct {
	authSvc authpb.AuthSvcClient
}

var store sessions.Store

func dieIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s. Try --help for help.\n", err)
	os.Exit(-1)
}

// extract router-specific headers to support dynamic routing & tracing
func getContext(req *http.Request) context.Context {
	headers := make(map[string]string)
	for k, values := range req.Header {
		prefixed := func(s string) bool { return strings.HasPrefix(k, s) }
		if prefixed("L5d-") || prefixed("Dtab-") || prefixed("X-Dtab-") {
			if len(values) > 0 {
				headers[k] = values[0]
			}
		}
	}
	md := metadata.New(headers)
	ctx := metadata.NewContext(context.Background(), md)
	return ctx
}

func main() {
	caCertFile := flag.String("cacert", "", "path to PEM-encoded CA certificate")
	addr := flag.String("srv", ":3000", "TCP address to listen to (in host:port form)")
	authAddr := flag.String("auth-addr", "localhost:5000", "Address of the auth service")
	authName := flag.String("auth-name", "", "Common name of auth service")
	stateStoreSecret := flag.String("state-store-secret", "", "State store secret name")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	if *stateStoreSecret == "" {
		dieIf(fmt.Errorf("Expecting stateStoreSecret not to be empty, got %s", *stateStoreSecret))
	}

	var authCreds grpc.DialOption
	if *caCertFile == "" {
		authCreds = grpc.WithInsecure()
	} else {
		creds, err := credentials.NewClientTLSFromFile(*caCertFile, *authName)
		dieIf(err)
		authCreds = grpc.WithTransportCredentials(creds)
	}
	authConn, err := grpc.Dial(*authAddr, authCreds)
	dieIf(err)
	defer authConn.Close()
	authClient := authpb.NewAuthSvcClient(authConn)

	router := mux.NewRouter()
	store = sessions.NewCookieStore([]byte(*stateStoreSecret))
	authHandlers := handlers.AuthHandler{
		AuthSvc: authClient,
		Store:   store,
	}
	router.HandleFunc("/login", authHandlers.HandleLoginURLRequest).Methods("GET")
	router.HandleFunc("/auth", authHandlers.HandleAuth).Methods("GET")
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(*addr, n)
}
