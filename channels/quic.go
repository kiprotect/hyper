package channels

/*
- Create a QUIC server and QUIC client
- Create local TCP-based servers
- Forward incoming connections to local ports
-
*/

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/tls"
	"github.com/quic-go/quic-go"
	"io"
	"net"
)

type QUICSettings struct {
	TLS         *tls.TLSSettings     `json:"tls"`
	BindAddress string               `json:"bindAddress"`
	Channels    []*QUICChannelConfig `json:"channels"`
}

type QUICChannelConfig struct {
	Remote *QUICRemoteChannel `json:"remote"`
	Local  *QUICLocalChannel  `json:"local"`
}

type QUICLocalChannel struct {
	Port int64  `json:"port"`
	Host string `json:"string"`
}

type QUICRemoteChannel struct {
	Port int64  `json:"port"`
	Host string `json:"host"`
}

type QUICChannel struct {
	hyper.BaseChannel
	Settings *QUICSettings
}

var QUICRemoteForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "host",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "port",
			Validators: []forms.Validator{
				forms.IsInteger{
					HasMin: true,
					Min:    1,
					HasMax: true,
					Max:    65535,
				},
			},
		},
	},
}

var QUICLocalForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "port",
			Validators: []forms.Validator{
				forms.IsInteger{
					HasMin: true,
					Min:    1,
					HasMax: true,
					Max:    65535,
				},
			},
		},
		{
			Name: "host",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "0.0.0.0"},
				forms.IsString{},
			},
		},
	},
}

var QUICChannelForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "remote",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &QUICRemoteForm,
				},
			},
		},
		{
			Name: "local",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &QUICLocalForm,
				},
			},
		},
	},
}

var QUICForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "tls",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &tls.TLSSettingsForm,
				},
			},
		},
		{
			Name: "bindAddress",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "channels",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &QUICChannelForm,
						},
					},
				},
			},
		},
	},
}

func QUICSettingsValidator(settings map[string]interface{}) (interface{}, error) {
	if params, err := QUICForm.Validate(settings); err != nil {
		return nil, err
	} else {
		validatedSettings := &QUICSettings{}
		if err := QUICForm.Coerce(validatedSettings, params); err != nil {
			return nil, err
		}
		return validatedSettings, nil
	}
}

func MakeQUICChannel(settings interface{}) (hyper.Channel, error) {
	quicSettings := settings.(QUICSettings)
	return &QUICChannel{
		Settings: &quicSettings,
	}, nil
}

func (q *QUICChannel) Type() string {
	return "quic"
}

func (q *QUICChannel) CanDeliverTo(*hyper.Address) bool {
	return false
}

func (q *QUICChannel) DeliverRequest(*hyper.Request) (*hyper.Response, error) {
	return nil, nil
}

func (q *QUICChannel) Close() error {
	return nil
}

func (q *QUICChannel) server(listener *quic.Listener) {
	for {
		conn, err := listener.Accept(context.Background())

		if err != nil {
			hyper.Log.Errorf("Cannot accept QUIC connection: %v", err)
			break
		}

		go func() {
			for {
				stream, err := conn.AcceptStream(context.Background())

				if err != nil {
					hyper.Log.Errorf("Cannot accept QUIC stream: %v", err)
					break
				}

				hyper.Log.Info("Accepted a QUIC stream")

				bs2 := make([]byte, 2)

				if _, err := io.ReadFull(stream, bs2); err != nil {
					hyper.Log.Error("Cannot read port")
					return
				}

				port := binary.LittleEndian.Uint16(bs2)

				conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))

				if err != nil {
					hyper.Log.Error("Cannot connect to local port %d", port)
					return
				}

				close := func() {
					conn.Close()
					stream.Close()
				}

				go pipe(stream, conn, close)
				go pipe(conn, stream, close)

			}
		}()
	}

}

type QUICChannelSettings struct {
	Address string `json:"address"`
}

var QUICChannelSettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "address",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

func getQUICChannelSettings(settings map[string]interface{}) (*QUICChannelSettings, error) {
	if params, err := QUICChannelSettingsForm.Validate(settings); err != nil {
		return nil, err
	} else {
		validatedSettings := &QUICChannelSettings{}
		if err := QUICChannelSettingsForm.Coerce(validatedSettings, params); err != nil {
			return nil, err
		}
		return validatedSettings, nil
	}
}

func (q *QUICChannel) handle(conn net.Conn, channel *QUICChannelConfig) {

	// we check if the requested service offers a gRPC server
	if entry, err := q.DirectoryEntry(channel.Remote.Host, "quic"); entry != nil {
		if settings, err := getQUICChannelSettings(entry.Channel("quic").Settings); err != nil {
			hyper.Log.Error(err)
			return
		} else {
			if err := q.pipe(conn, settings.Address, channel.Remote.Host, channel.Remote.Port); err != nil {
				hyper.Log.Errorf("Cannot connect: %v", err)
				return
			}
		}
	} else if err != nil {
		// we log this error
		hyper.Log.Error(err)
		return
	}
}

func (q *QUICChannel) channel(channel *QUICChannelConfig) error {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", channel.Local.Host, channel.Local.Port))

	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				break
			}

			go q.handle(conn, channel)
		}
	}()

	return nil

}

func (q *QUICChannel) client() {
	for _, channel := range q.Settings.Channels {
		go q.channel(channel)
	}
}

func pipe(left, right io.ReadWriteCloser, close func()) {
	buf := make([]byte, 4096)
	for {
		n, err := left.Read(buf)
		if err != nil {
			hyper.Log.Error(err)
			close()
			return
		}
		if m, err := right.Write(buf[:n]); err != nil {
			hyper.Log.Error(err)
			close()
			return
		} else if m != n {
			hyper.Log.Errorf("cannot write all data")
			close()
			return
		}
	}
}

func (q *QUICChannel) pipe(conn net.Conn, addr, serverName string, port int64) error {

	config, err := tls.TLSClientConfig(q.Settings.TLS)

	if err != nil {
		return err
	}

	config.ServerName = serverName

	config.NextProtos = []string{"hyper-quic"}

	quicConn, err := quic.DialAddr(context.Background(), addr, config, nil)

	if err != nil {
		return err
	}

	stream, err := quicConn.OpenStreamSync(context.Background())

	if err != nil {
		return err
	}

	close := func() {
		stream.Close()
		conn.Close()
	}

	bs2 := make([]byte, 2)

	binary.LittleEndian.PutUint16(bs2, uint16(port))

	if _, err := stream.Write(bs2); err != nil {
		return err
	}

	hyper.Log.Debugf("Proxying connection...")

	go pipe(stream, conn, close)
	go pipe(conn, stream, close)

	return nil

}

func (q *QUICChannel) Open() error {
	hyper.Log.Info("Opening QUIC channel...")

	config, err := tls.TLSServerConfig(q.Settings.TLS)

	if err != nil {
		return err
	}

	config.NextProtos = []string{"hyper-quic"}

	listener, err := quic.ListenAddr(q.Settings.BindAddress, config, nil)

	if err != nil {
		return err
	}

	go q.server(listener)
	go q.client()

	return nil
}
