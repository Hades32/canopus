package canopus

import (
	"github.com/bocajim/dtls"
	"github.com/prometheus/common/log"
	"net"
	"time"
)

type dtlsConnection struct {
	peer *dtls.Peer
}

var listener *dtls.Listener = nil
var keystore *dtls.KeystoreInMemory

func initListener() error {
	if listener != nil {
		return nil
	}
	listener, err := dtls.NewUdpListener(":6000", 5*time.Second)
	if err != nil {
		log.Warnln("shit", err)
		return err
	}
	listener.AddCipherSuite(dtls.CipherSuite_TLS_PSK_WITH_AES_128_CCM_8)
	keystore = dtls.NewKeystoreInMemory()
	dtls.SetKeyStores([]dtls.Keystore{keystore})
	return nil
}

func NewDtlsConnection(addr, identity string, psk []byte) (net.Conn, error) {
	err := initListener()
	if err != nil {
		return nil, err
	}
	peer, err := listener.AddPeer(addr, identity)
	if err != nil {
		log.Warnln("shit", err)
		return nil, err
	}
	keystore.AddKey(identity, psk)
	return &dtlsConnection{
		peer: peer,
	}, nil
}

func (dc *dtlsConnection) Read(b []byte) (n int, err error) {
	bytes, err := dc.peer.Read(5 * time.Second)
	if err != nil {
		return 0, err
	}
	copy(b, bytes)
	return len(bytes), nil
}

func (dc *dtlsConnection) Write(b []byte) (n int, err error) {
	err = dc.peer.Write(b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (dc *dtlsConnection) Close() error {
	dc.peer.Close(dtls.AlertDesc_CloseNotify)
	return nil
}

func (dc *dtlsConnection) LocalAddr() net.Addr {
	panic("implement me")
}

func (dc *dtlsConnection) RemoteAddr() net.Addr {
	panic("implement me")
}

func (dc *dtlsConnection) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (dc *dtlsConnection) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (dc *dtlsConnection) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
