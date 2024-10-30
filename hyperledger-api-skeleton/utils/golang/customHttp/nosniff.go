package customHttp

import (
	"net/http"
)

/*
The X-Content-Type-Options header is used to protect against MIME sniffing vulnerabilities.
A response is sent back with the header X-Content-Type-Options: nosniff.
This prevents the client from “sniffing” the asset to try and determine if the file type is something other than what is declared by the server.
*/

func NoSniff(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		//Add nos niff to the header
		w.Header().Set("X-Content-Type-Options", "nosniff")

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
