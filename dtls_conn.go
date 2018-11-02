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

type dtlsPacketConnection struct {
	*dtls.Listener
}

var listener *dtlsPacketConnection = nil
var keystore *dtls.KeystoreInMemory

func initListener() error {
	if listener != nil {
		return nil
	}
	var err error
	listener, err = NewDtlsListener(":6000")
	if err != nil {
		log.Fatalln("shit", err)
	}
	keystore = dtls.NewKeystoreInMemory()
	dtls.SetKeyStores([]dtls.Keystore{keystore})
	return nil
}

func NewDtlsListener(addr string) (*dtlsPacketConnection, error) {
	listener, err := dtls.NewUdpListener(addr, 5*time.Second)
	if err != nil {
		log.Warnln("shit", err)
		return nil, err
	}
	listener.AddCipherSuite(dtls.CipherSuite_TLS_PSK_WITH_AES_128_CCM_8)
	return &dtlsPacketConnection{listener}, nil
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

func (dc *dtlsConnection) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	panic("implement me")
}

func (dc *dtlsConnection) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	panic("implement me")
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

func (dpc *dtlsPacketConnection) Close() error {
	//TODO
	return nil
}

func (dpc *dtlsPacketConnection) LocalAddr() net.Addr {
	return nil
}

func (dpc *dtlsPacketConnection) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	bytes, peer := dpc.Listener.Read()
	if err != nil {
		return 0, nil, err
	}
	copy(p, bytes)
	return len(bytes), &net.UDPAddr{Zone: peer.RemoteAddr()}, nil
}

func (dpc *dtlsPacketConnection) SetDeadline(t time.Time) error {
	return nil
}

func (dpc *dtlsPacketConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (dpc *dtlsPacketConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func (dpc *dtlsPacketConnection) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	panic("implement me")
}
