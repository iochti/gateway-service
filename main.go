package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/namsral/flag"
	"github.com/urfave/negroni"

	authpb "github.com/iochti/auth-service/proto"
	"github.com/iochti/gateway-service/handlers"
	"github.com/iochti/gateway-service/helpers"
	thingpb "github.com/iochti/thing-service/proto"
	userpb "github.com/iochti/user-service/proto"
)

type IochtiGateway struct {
	authSvc  authpb.AuthSvcClient
	userSvc  userpb.UserSvcClient
	thingSvc thingpb.ThingSvcClient
}

var store sessions.Store

const CONTENT_TYPE = "application/json"

func main() {
	caCertFile := flag.String("cacert", "", "path to PEM-encoded CA certificate")
	addr := flag.String("srv", ":3000", "TCP address to listen to (in host:port form)")
	authAddr := flag.String("auth-addr", "localhost:5000", "Address of the auth service")
	authName := flag.String("auth-name", "", "Common name of auth service")
	userAddr := flag.String("user-addr", "localhost:5001", "Address of the user service")
	userName := flag.String("user-name", "", "Common name of user service")
	thingAddr := flag.String("thing-addr", "localhost:5002", "Address of the thing service")
	thingName := flag.String("thing-name", "", "Common name of thing service")
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

	// -----------------------Thing service declaration------------------------------
	thingConn := helpers.DeclareService(thingAddr, caCertFile, thingName)
	defer thingConn.Close()
	thingClient := thingpb.NewThingSvcClient(thingConn)

	// Handlers creation
	router := mux.NewRouter()
	var store *sessions.CookieStore
	if *cst == "" {
		store = sessions.NewCookieStore([]byte(RandToken(64)))
	} else {
		store = sessions.NewCookieStore([]byte(*cst))
	}

	authHandlers := handlers.AuthHandler{
		AuthSvc:     authClient,
		UserSvc:     userClient,
		Store:       store,
		ContentType: CONTENT_TYPE,
	}

	userHandler := handlers.UserHandler{
		UserSvc:     userClient,
		ContentType: CONTENT_TYPE,
	}

	thingHandler := handlers.ThingHandler{
		ThingSvc:    thingClient,
		ContentType: CONTENT_TYPE,
	}

	router.HandleFunc("/login", authHandlers.HandleLoginURLRequest).Methods("GET")
	router.HandleFunc("/auth", authHandlers.HandleAuth).Methods("GET")
	router.HandleFunc("/user", userHandler.HandleGetUser).Methods("GET")
	router.HandleFunc("/user", userHandler.HandleCreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", userHandler.HandleDeleteUser).Methods("DELETE")
	router.HandleFunc("/thing/{id}", thingHandler.HandleGetThing).Methods("GET")
	router.HandleFunc("/thing", thingHandler.HandleCreateThing).Methods("POST")
	router.HandleFunc("/thing", thingHandler.HandleUpdateThing).Methods("PUT")
	router.HandleFunc("/thing/one/{id}", thingHandler.HandleDeleteOneThing).Methods("DELETE")
	router.HandleFunc("/thing/many", thingHandler.HandleDeleteManyThings).Methods("DELETE")
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
