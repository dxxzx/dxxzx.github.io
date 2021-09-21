---
title: "Kubernetes Authenticating Proxy"
date: 2021-09-20T22:27:01+08:00
draft: true
tags: []
topics: []
description: ""
---

# 1 Overview

Refer to `net/http/httputil/reverseproxy.go`. Modify it as a Kubernetes Authenticating Proxy.

# 2 Official Docuement
> ### Authenticating Proxy
>
> The API server can be configured to identify users from request header values, such as `X-Remote-User`.
> It is designed for use in combination with an authenticating proxy, which sets the request header value.
> 
> * `--requestheader-username-headers` Required, case-insensitive. Header names to check, in order, for the user identity. The first header containing a value is used as the username.
> * `--requestheader-group-headers` 1.6+. Optional, case-insensitive. "X-Remote-Group" is suggested. Header names to check, in order, for the user's groups. All values in all specified headers are used as group names.
> * `--requestheader-extra-headers-prefix` 1.6+. Optional, case-insensitive. "X-Remote-Extra-" is suggested. Header prefixes to look for to determine extra information about the user (typically used by the configured authorization plugin). Any headers beginning with any of the specified prefixes have the prefix removed. The remainder of the header name is lowercased and [percent-decoded](https://tools.ietf.org/html/rfc3986#section-2.1) and becomes the extra key, and the header value is the extra value.
>
> Prior to 1.11.3 (and 1.10.7, 1.9.11), the extra key could only contain characters which were [legal in HTTP header labels](https://tools.ietf.org/html/rfc7230#section-3.2.6).
>
> For example, with this configuration:
>
> ```
> --requestheader-username-headers=X-Remote-User
> --requestheader-group-headers=X-Remote-Group
> --requestheader-extra-headers-prefix=X-Remote-Extra-
> ```
>
> this request:
>
> ```http
> GET / HTTP/1.1
> X-Remote-User: fido
> X-Remote-Group: dogs
> X-Remote-Group: dachshunds
> X-Remote-Extra-Acme.com%2Fproject: some-project
> X-Remote-Extra-Scopes: openid
> X-Remote-Extra-Scopes: profile
> ```
>
> would result in this user info:
>
> ```yaml
> name: fido
> groups:
> - dogs
> - dachshunds
> extra:
>   acme.com/project:
>   - some-project
>   scopes:
>   - openid
>   - profile
> ```
>
>
> In order to prevent header spoofing, the authenticating proxy is required to present a valid client
> certificate to the API server for validation against the specified CA before the request headers are
> checked. WARNING: do **not** reuse a CA that is used in a different context unless you understand
> the risks and the mechanisms to protect the CA's usage.
>
> * `--requestheader-client-ca-file` Required. PEM-encoded certificate bundle. A valid client certificate must be presented and validated against the certificate authorities in the specified file before the request headers are checked for user names.
> * `--requestheader-allowed-names` Optional. List of Common Name values (CNs). If set, a valid client certificate with a CN in the specified list must be presented before the request headers are checked for user names. If empty, any CN is allowed.

# 3 Analyse ReverseProxy

The key method of reverseproxy is `func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)`. Look at this method, first, it prepared a http.Transport, search it, we could see it used as `res, err := transport.RoundTrip(outreq)` later. Guess it is used to obtain response fron kube-apiserver. So, it need to setup TLSClientConfig.
```go
transport := &http.Transport{TLSClientConfig: TLSClientConfig}
```

Then, refer to Official document, we need to setup authentication headers in out request. Setup these after clone incomming request as out request and before send the out request. It may look like this:

```go
outReq.Header.Add("X-Remote-User", "admin")
outReq.Header.Add("X-Remote-Group", "system:master")
```

Finally, we need change the destination of this request, set it host to kube-apiserver host. Set it after clone incomming request as out request and before send the out request.
```go
outReq.URL.Host = p.KubeHost
outReq.URL.Scheme = "https"
```

# 4 Prepare Certificate

```sh
mkdir -pv certs
cp -rv /etc/kubernetes/pki/{ca.crt,front-proxy-client.{crt,key}} certs
chown ${username}:${username} -R certs  # set username as your normal username
```

# 5 All code

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"strings"
)

const (
	caCertPath               = "certs/ca.crt"
	frontProxyClientCertPath = "certs/front-proxy-client.crt"
	frontProxyClientKeyPath  = "certs/front-proxy-client.key"
)

var (
	// Hop-by-hop headers. These are removed when sent to the backend.
	// As of RFC 7230, hop-by-hop headers are required to appear in the
	// Connection header field. These are the headers defined by the
	// obsoleted RFC 2616 (section 13.5.1) and are used for backward
	// compatibility.
	hopHeaders = []string{
		"Connection",
		"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",      // canonicalized version of "TE"
		"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
		"Transfer-Encoding",
		"Upgrade",
	}
)

type KubeAccessProxy struct {
	KubeHost        string
	TLSClientConfig *tls.Config
}

func (p *KubeAccessProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	transport := &http.Transport{TLSClientConfig: p.TLSClientConfig}

	ctx := req.Context()
	outReq := req.Clone(ctx)

	if req.ContentLength == 0 {
		outReq.Body = nil
	}
	if outReq.Header == nil {
		outReq.Header = make(http.Header)
	}

	outReq.Close = false

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		prior, ok := outReq.Header["X-Forwarded-For"]
		omit := ok && prior == nil // Issue 38079: nil now means don't populate the header
		if len(prior) > 0 {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		if !omit {
			outReq.Header.Set("X-Forwarded-For", clientIP)
		}
	}
	outReq.URL.Host = p.KubeHost
	outReq.URL.Scheme = "https"
	log.Printf("request url: %s", outReq.URL)
	outReq.Header.Add("X-Remote-User", "admin")
	outReq.Header.Add("X-Remote-Group", "system:master")

	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		log.Printf("err: %s", err.Error())
		return
	}

	removeConnectionHeaders(res.Header)

	for _, h := range hopHeaders {
		res.Header.Del(h)
	}

	copyHeader(rw.Header(), res.Header)

	// The "Trailer" header isn't included in the Transport's response,
	// at least for *http.Transport. Build it up from Trailer.
	announcedTrailers := len(res.Trailer)
	if announcedTrailers > 0 {
		trailerKeys := make([]string, 0, len(res.Trailer))
		for k := range res.Trailer {
			trailerKeys = append(trailerKeys, k)
		}
		rw.Header().Add("Trailer", strings.Join(trailerKeys, ", "))
	}

	rw.WriteHeader(res.StatusCode)

	err = p.copyResponse(rw, res.Body)
	if err != nil {
		defer res.Body.Close()
		panic(http.ErrAbortHandler)
	}
	res.Body.Close() // close now, instead of defer, to populate res.Trailer

	if len(res.Trailer) > 0 {
		// Force chunking if we saw a response trailer.
		// This prevents net/http from calculating the length for short
		// bodies and adding a Content-Length.
		if fl, ok := rw.(http.Flusher); ok {
			fl.Flush()
		}
	}

	if len(res.Trailer) == announcedTrailers {
		copyHeader(rw.Header(), res.Trailer)
		return
	}

	for k, vv := range res.Trailer {
		k = http.TrailerPrefix + k
		for _, v := range vv {
			rw.Header().Add(k, v)
		}
	}
}

func (p *KubeAccessProxy) copyResponse(dst io.Writer, src io.Reader) error {
	var buf []byte
	_, err := p.copyBuffer(dst, src, buf)
	return err
}

// copyBuffer returns any write errors or non-EOF read errors, and the amount
// of bytes written.
func (p *KubeAccessProxy) copyBuffer(dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	if len(buf) == 0 {
		buf = make([]byte, 32*1024)
	}
	var written int64
	for {
		nr, rerr := src.Read(buf)
		if rerr != nil && rerr != io.EOF && rerr != context.Canceled {
			log.Printf("httputil: ReverseProxy read error during body copy: %v", rerr)
		}
		if nr > 0 {
			nw, werr := dst.Write(buf[:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if werr != nil {
				return written, werr
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}
		if rerr != nil {
			if rerr == io.EOF {
				rerr = nil
			}
			return written, rerr
		}
	}
}
func removeConnectionHeaders(h http.Header) {
	for _, f := range h["Connection"] {
		for _, sf := range strings.Split(f, ",") {
			if sf = textproto.TrimString(sf); sf != "" {
				h.Del(sf)
			}
		}
	}
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	caCertData, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}
	frontProxyClientCertData, err := ioutil.ReadFile(frontProxyClientCertPath)
	if err != nil {
		panic(err)
	}
	frontProxyClientKeyData, err := ioutil.ReadFile(frontProxyClientKeyPath)
	if err != nil {
		panic(err)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertData)
	if !ok {
		panic("failed to parse root certificate")
	}
	proxyCert, err := tls.X509KeyPair(frontProxyClientCertData, frontProxyClientKeyData)
	if err != nil {
		panic(err)
	}
	proxyCerts := []tls.Certificate{proxyCert}
	tlsConfig := &tls.Config{
		Certificates: proxyCerts,
		RootCAs:      roots,
	}
	handler := &KubeAccessProxy{
		KubeHost:        "192.168.0.2:6443",
		TLSClientConfig: tlsConfig,
	}
	log.Fatal(http.ListenAndServe("localhost:8080", handler))
}
```