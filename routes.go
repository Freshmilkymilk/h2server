// Copyright (c) 2018 Henry Slawniak <https://datacenterscumbags.com/>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"net/http"
	"os"
	"strings"
)

func addRoutesToRouter() {

	server.r.PathPrefix("/").HandlerFunc(server.indexHandler).Name("catch-all")
}

// GetIP returns the remote ip of the request, by stripping off the port from the RemoteAddr
func GetIP(r *http.Request) string {
	split := strings.Split(r.RemoteAddr, ":")
	ip := strings.Join(split[:len(split)-1], ":")
	// This is bad, and I feel bad
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	return ip
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	host := strings.Split(r.Host, ":")[0]

	// if !domainIsRegistered(host) {
	// 	log.Debugf("Host is %s", host)
	// 	// Make sure we do this syncronousley
	// 	addToDomainList(host)
	// }

	staticFolder := "./sites/" + host
	if _, err := os.Stat(staticFolder); err != nil {
		staticFolder = "./client"
	}

	var n int64
	var code int

	if inf, err := os.Stat(staticFolder + path); err == nil && !inf.IsDir() {
		n, code = serveFile(w, r, staticFolder+path)
	} else if inf, err := os.Stat(staticFolder + path + "/index.html"); err == nil && !inf.IsDir() {
		n, code = serveFile(w, r, staticFolder+path+"/index.html")
	} else {
		n, code = serveFile(w, r, staticFolder+"/index.html")
	}

	go logRequest(w, r, n, code)
}
