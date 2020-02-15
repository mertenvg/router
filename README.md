# router
A simple http router

## usage
```go
package main

import (
    "fmt"  
    "net/http"

    "github.com/whitecypher/router"
)

func main() {
    // create the router
    r := router.New()

    // add some middleware
    r.Middleware(
        func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
        		ra := r.RemoteAddr
        		if ff := r.Header.Get("X-Forwarded-For"); len(ff) > 0 {
        			ra = ff
        		}
        		fmt.Println(ra, r.Method, r.URL.Path)
        		next(w, r)
        	},
    )
    
    // catch everything
    r.HandleFunc(router.Default, func(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("I'm not a teapot"))
    	w.WriteHeader(http.StatusOK)
    })
    
    // start the server
    server := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }
    server.ListenAndServe()
}
```