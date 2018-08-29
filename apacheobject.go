package goolog2

import (
	"time"
)

// This is a logging interface used for the Apache combined logging
// format.
type ApacheObject interface {
	// Get values to be logged
	//
	// Returns:
	//     remoteHost: IP or host of the client
	//     identity: identity string (RFC 1413)
	//     user: user name determined by the http authentication
	//     requestTime: time when the request comes (local time)
	//     method: http method
	//     resource: resource from the request line
	//     protocol: used protocol (ie HTTP/1.0)
	//     status: status code returned to the user
	//     length: length of the response without the header
	//     referer: referer string
	//     agent: user agent (browser)
	GetValues() (
		remoteHost string,
		identity string,
		user string,
		requestTime time.Time,
		method string,
		resource string,
		protocol string,
		status int,
		length uint64,
		referer string,
		agent string)
}
