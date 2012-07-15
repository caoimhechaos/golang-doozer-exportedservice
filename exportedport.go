/**
 * Exported named Doozer port.
 * This binds to an anonymous port, exports the host:port pair through Doozer
 * and returns the port to the caller.
 */
package exportedservice

import (
	"crypto/tls"
	"fmt"
	"github.com/4ad/doozer"
	"net"
)

/**
 * Open a new anonymous port on "ip" and export it through Doozer as
 * "servicename". If "ip" is a host:port pair, the port will be overridden.
 */
func NewExportedPort(network, ip, servicename string) (net.Listener, error) {
	var host, hostport string
	var conn *doozer.Conn
	var l net.Listener
	var err error
	var i uint

	if host, _, err = net.SplitHostPort(ip); err != nil {
		// Apparently, it's not in host:port format.
		host = ip
	}

	hostport = net.JoinHostPort(host, "0")
	if l, err = net.Listen(network, hostport); err != nil {
		return nil, err
	}

	// Now write our host:port pair to Doozer. First, determine the next
	// free number.
	conn, err = doozer.Dial("doozer.l.internetputzen.com:8046")
	if err != nil {
		l.Close()
		return nil, err
	}

	// FIXME(tonnerre): Turn this into a more efficient implementation.
	for {
		var path string = fmt.Sprintf("/ns/service/%s/%d",
			servicename, i)
		_, err = conn.Set(path, 0, []byte(l.Addr().String()))
		if err == nil {
			return l, nil
		}

		i += 1
	}
	return nil, err
}

/**
 * Open a new anonymous port on "ip" and export it through Doozer as
 * "servicename". Associate the TLS configuration "config". If "ip" is
 * a host:port pair, the port will be overridden.
 */
func NewExportedPort(network, ip, servicename string,
		     config *tls.Config) (net.Listener, error) {
	var l net.Listener
	var err error

	// We can just create a new port as above...
	l, err = NewExportedPort(network, ip, servicename)
	if err != nil {
		return nil, err
	}

	// ... and inject a TLS context.
	return tls.NewListener(l, config)
}
