package mixpanel

import (
	"github.com/beatlabs/patron/log"
	"github.com/dukex/mixpanel"
	"nanomsg.org/go/mangos/v2/protocol/sub"

	"nanomsg.org/go/mangos/v2"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
)

type Handler struct {
	p Processor
	mangos.Socket
}

func (hdl *Handler) Run() {
	go func() {

		var msg []byte
		var err error

		for {
			if msg, err = hdl.Recv(); err != nil {
				log.Errorf("cannot recv: %s", err.Error())
				continue
			}

			err = hdl.p.Process(msg)
			if err != nil {
				log.Error(err)
			}
		}
	}()
}

type Processor interface {
	mixpanel.Mixpanel
	Process(msg []byte) error

}

// NewHandler returns a new Mixpanel handler
func NewHandler(p Processor, pubs []string) Handler{
	var sock mangos.Socket
	var err error

	if sock, err = sub.NewSocket(); err != nil {
		log.Fatal("can't get new sub socket: %s", err.Error())
	}

	for _,v := range pubs {
		if err = sock.Dial(v); err != nil {
			log.Fatal("can't dial on sub socket: %s", err.Error())
		}
		log.Debugf("listening to %s", v)
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		log.Fatal("cannot subscribe: %s", err.Error())
	}


	return Handler{Socket: sock, p: p}
}