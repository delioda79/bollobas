package mixpanel

import (
	"bollobas"
	"bollobas/pkg/parseid"
	"encoding/json"
	"fmt"

	"github.com/beatlabs/patron/log"
	"github.com/dukex/mixpanel"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
)

// Handler subscribes to messages sent by any registered publisher in the internal registry
type Handler struct {
	mangos.Socket
	mixpanel.Mixpanel
}

// Run starts the go routine which will receive the messages
func (hdl *Handler) Run() {
	go func() {

		var msg []byte
		var err error

		for {

			if msg, err = hdl.Recv(); err != nil {
				log.Errorf("cannot recv: %s", err.Error())
				continue
			}

			idt := &bollobas.Identity{}
			err := json.Unmarshal(msg, idt)
			if err != nil {
				log.Errorf("error while receiving message", err)
				continue
			}
			fmt.Println(string(msg))
			hdl.updateIdentity(idt)
		}
	}()
}

func (hdl *Handler) updateIdentity(idt *bollobas.Identity) {
	//id := idt.ID
	prps := &Identity{
		FirstName:        idt.FirstName,
		LastName:         idt.LastName,
		RegistrationDate: idt.RegistrationDate,
		ReferralCode:     idt.ReferralCode,
		Type:             idt.Type,
		Email:            idt.Email,
		Phone:            idt.Phone,
	}

	bts, err := json.Marshal(prps)
	if err != nil {
		log.Errorf("Impossible to unmarshal", err)
		return
	}

	mp := map[string]interface{}{}

	err = json.Unmarshal(bts, &mp)
	if err != nil {
		log.Errorf("error while unmarshaling the identity", err)
	}

	err = hdl.Update(idt.ID, &mixpanel.Update{Properties: mp, Operation:"$set"})
	if err != nil {
		log.Errorf("error while updating the identity", err)
	}
}

// NewHandler returns a new mixpanel handler
func NewHandler(token string, pubs []string) *Handler {
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

	return &Handler{
		Socket: sock,
		Mixpanel: mixpanel.New(token, ""),
	}
}
