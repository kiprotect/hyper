package channels

/*
- Create a QUIC server and QUIC client
- Create local TCP-based servers
- Forward incoming connections to local ports
-
*/

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/hyper"
)

type QUICSettings struct {
}

type QUICChannel struct {
	hyper.BaseChannel
	Settings *QUICSettings
}

var QUICForm = forms.Form{
	Fields: []forms.Field{},
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

func (q *QUICChannel) Open() error {
	hyper.Log.Info("Opening QUIC channel...")
	return nil
}
