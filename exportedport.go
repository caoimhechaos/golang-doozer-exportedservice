/**
 * Exported named Doozer port.
 * This binds to an anonymous port, exports the host:port pair through Doozer
 * and returns the port to the caller.
 */
package exportedservice

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/4ad/doozer"
)

// We need to initialize our Doozer client beforehand and keep it somewhere.
type ServiceExporter struct {
	conn    *doozer.Conn
	path    string
	pathrev int64
	uri     string
	buri    string
}

/**
 * Try to create a new exporter by connecting to Doozer.
 */
func NewExporter(uri, buri string) (*ServiceExporter, error) {
	var self *ServiceExporter = &ServiceExporter{}
	var err error

	self.conn, err = doozer.DialUri(uri, buri)

	// We couldn't connect, let our user know.
	if err != nil {
		return nil, err
	}
	return self, nil
}

/**
 * Open a new anonymous port on "ip" and export it through Doozer as
 * "servicename". If "ip" is a host:port pair, the port will be overridden.
 */
func (self *ServiceExporter) NewExportedPort(
	network, ip, servicename string) (net.Listener, error) {
	var host, hostport string
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
	// FIXME(caoimhe): Turn this into a more efficient implementation.
	for {
		var path string = fmt.Sprintf("/ns/service/%s/%d", servicename, i)
		var ok bool
		var derr *doozer.Error
		var rev int64

		rev, err = self.conn.Set(path, 0, []byte(l.Addr().String()))
		if err == nil {
			self.path = path
			self.pathrev = rev
			return l, nil
		}

		if derr, ok = err.(*doozer.Error); !ok ||
			derr.Err != doozer.ErrOldRev {
			return nil, err
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
func (self *ServiceExporter) NewExportedTLSPort(
	network, ip, servicename string,
	config *tls.Config) (net.Listener, error) {
	var l net.Listener
	var err error

	// We can just create a new port as above...
	l, err = self.NewExportedPort(network, ip, servicename)
	if err != nil {
		return nil, err
	}

	// ... and inject a TLS context.
	return tls.NewListener(l, config), nil
}

/**
 * Remove the associated exported port. This will only delete the most
 * recently exported port.
 */
func (self *ServiceExporter) UnexportPort() error {
	var derr *doozer.Error
	var err error
	var ok bool

	if len(self.path) == 0 {
		return nil
	}

	if err = self.conn.Del(self.path, self.pathrev); err != nil {
		if derr, ok = err.(*doozer.Error); !ok ||
			derr.Err != doozer.ErrOldRev {
			return err
		}
	}

	return nil
}
