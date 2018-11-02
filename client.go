package canopus

import "net"

func Dial(address string) (conn Connection, err error) {
	udpConn, err := net.Dial("udp", address)
	if err != nil {
		return
	}

	conn = &UDPConnection{
		conn: udpConn,
	}

	return
}

func DialDTLS(address, identity, psk string) (Connection, error) {
	conn, err := NewDtlsConnection(address, identity, []byte(psk))
	if err != nil {
		return nil, err
	}
	return &UDPConnection{
		conn: conn,
	}, err
}

func NewObserveMessage(r string, val interface{}, msg Message) ObserveMessage {
	return &CoapObserveMessage{
		Resource: r,
		Value:    val,
		Msg:      msg,
	}
}

type CoapObserveMessage struct {
	CoapMessage
	Resource string
	Value    interface{}
	Msg      Message
}

func (m *CoapObserveMessage) GetResource() string {
	return m.Resource
}

func (m *CoapObserveMessage) GetValue() interface{} {
	return m.Value
}

func (m *CoapObserveMessage) GetMessage() Message {
	return m.GetMessage()
}
