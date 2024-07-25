package client

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/ishidawataru/sctp"
)

func MakeConnection(laddr, raddr, transport string) (diam.Conn, error) {
	mux := sm.New(newStateMachineSettings("host", "realm"))
	mux.HandleFunc("ALL", handleALL) // Catch all.

	client := newClient(mux)
	return connectClient(client, laddr, raddr, transport)
}

func connectClient(client *sm.Client, laddr, raddr, protocol string) (diam.Conn, error) {
	log.Printf("connecting to %s, from %s\n", raddr, laddr)

	var la net.Addr
	var err error

	switch protocol {
	case "sctp":
		la, err = sctp.ResolveSCTPAddr(protocol, laddr)
		if err != nil {
			return nil, err
		}
	case "tcp":
		la, err = net.ResolveTCPAddr(protocol, laddr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported protocol " + protocol)
	}

	connection, err := client.DialExt(protocol, raddr, time.Second*5, la)
	if err != nil {
		return nil, err
	}

	log.Println("client connected, handshake ok")
	return connection, nil
}

func newClient(mux *sm.StateMachine) *sm.Client {
	return &sm.Client{
		Dict:    dict.Default,
		Handler: mux,
		AuthApplicationID: []*diam.AVP{
			// Advertise support for credit control application
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)), // RFC 4006
		},
	}
}

func newStateMachineSettings(host, realm string) *sm.Settings {
	return &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(host),
		OriginRealm:      datatype.DiameterIdentity(realm),
		VendorID:         10415,
		ProductName:      datatype.UTF8String(realm),
		FirmwareRevision: 1,
	}
}

func handleALL(c diam.Conn, m *diam.Message) {
	log.Printf("received unexpected message from %s:\n%s\n", c.RemoteAddr(), m)
}
