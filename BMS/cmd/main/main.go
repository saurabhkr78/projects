package main

import (
	"log"
	"net/http"

	"BMS/pkg/configs"
	"BMS/pkg/routes"

	"github.com/gorilla/mux"
	"time"
)

func main() {
	// Initialize database.
	configs.Connect()
	log.Println("Database connected")

	// Create router.
	r := mux.NewRouter()

	// Register routes.
	routes.RegisterBookStoreRoutes(r)

	// Create HTTP server.
	/*you're creating an HTTP server object.

	Think of it like creating a car before driving it.
	The struct stores all the settings the server will use.
	*/
	server := &http.Server{
		Addr:        ":8080", //Listen on TCP port 8080.
		ReadTimeout: 5 * time.Second,
		/*3. ReadTimeout
		ReadTimeout: 5 * time.Second,

		Imagine a malicious client.

		He connects.

		But instead of sending

		GET /book

		he sends

		G

		waits 30 seconds

		then

		E

		waits another minute...

		This is called a Slowloris attack.

		Without a timeout,

		your server waits forever.

		Eventually,

		thousands of attackers connect.

		Now your server has thousands of stuck connections.

		With

		ReadTimeout: 5 * time.Second,

		you're saying

		If the client doesn't finish sending the request within 5 seconds, disconnect them.
		*/
		WriteTimeout: 10 * time.Second,
		/*
					4. WriteTimeout
			WriteTimeout: 10 * time.Second,

			Suppose your server has already processed the request.

			Now it wants to send

			{
			    ...
			}

			But the client reads extremely slowly.

			Without timeout,

			your server keeps waiting.

			With

			WriteTimeout: 10 * time.Second,

			you're saying

			If I can't finish writing my response within 10 seconds,

			close the connection.
		*/
		IdleTimeout: 60 * time.Second,
		/*
					5. IdleTimeout

			Imagine

			GET /book

			Response sent.

			Connection stays open.

			Nothing happens.

			One minute later...

			Still nothing.

			Eventually you may have

			20,000 idle connections.

			They're doing absolutely nothing.

			This wastes memory.

			So

			IdleTimeout: 60 * time.Second,

			means

			If a connection stays idle for 60 seconds,

			close it.
		*/
		Handler: r, //"Whenever a request comes in, who should handle it?"
		//MaxHeaderBytes:    1 << 20, //1 << 20 is 1 megabyte. This is the maximum size of request headers your server will accept.If someone sends 500 MB headers the server rejects them Protects against attacks.
	}

	log.Println("Server started on http://localhost:8080")

	// Start server.
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
