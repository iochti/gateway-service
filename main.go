package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/namsral/flag"
	"github.com/urfave/negroni"

	authpb "github.com/iochti/auth-service/proto"
	"github.com/iochti/gateway-service/handlers"
	"github.com/iochti/gateway-service/helpers"
	userpb "github.com/iochti/user-service/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type IochtiGateway struct {
	authSvc authpb.AuthSvcClient
	userSvc userpb.UserSvcClient
}

var store sessions.Store

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
	userAddr := flag.String("user-addr", "localhost:5001", "Address of the user service")
	userName := flag.String("user-name", "", "Common name of user service")
	cst := flag.String("cookie-store-token", "", "Cookie store token (only for development environments)")
	flag.Parse()
	if flag.NArg() != 0 {
		helpers.DieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	// -----------------------Auth service declaration------------------------------
	authConn := helpers.DeclareService(authAddr, caCertFile, authName)
	defer authConn.Close()
	authClient := authpb.NewAuthSvcClient(authConn)

	// -----------------------User service declaration------------------------------
	userConn := helpers.DeclareService(userAddr, caCertFile, userName)
	defer userConn.Close()
	userClient := userpb.NewUserSvcClient(userConn)

	// Handlers creation
	router := mux.NewRouter()
	var store *sessions.CookieStore
	if *cst == "" {
		store = sessions.NewCookieStore([]byte(RandToken(64)))
	} else {
		store = sessions.NewCookieStore([]byte(*cst))
	}

	authHandlers := handlers.AuthHandler{
		AuthSvc: authClient,
		UserSvc: userClient,
		Store:   store,
	}

	userHandler := handlers.UserHandler{
		UserSvc: userClient,
	}

	router.HandleFunc("/login", authHandlers.HandleLoginURLRequest).Methods("GET")
	router.HandleFunc("/auth", authHandlers.HandleAuth).Methods("GET")
	router.HandleFunc("/user", userHandler.HandleGetUser).Methods("GET")
	router.HandleFunc("/user", userHandler.HandleCreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", userHandler.HandleDeleteUser).Methods("DELETE")
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(*addr, n)
}

// RandToken generates a random token of l length
func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
