package routers

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func ListenAndServe(mux *http.ServeMux) {
	port := 8000
	for {
		// Try to listen on the current port
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Printf("Port %d is in use, trying next...\n", port)
			port++
			continue
		}
		// Close the listener immediately so we can reuse the port
		ln.Close()

		// Start the server
		log.Printf("Server running on port: %d\n", port)
		fmt.Printf("URL: http://localhost:%d\n", port)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
		break
	}
}
