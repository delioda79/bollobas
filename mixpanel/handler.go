package mixpanel

import (
	"bollobas"
	"bollobas/pkg/parseid"
	"encoding/json"
	"fmt"

	"github.com/beatlabs/patron/log"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
)

type Handler struct {
	mangos.Socket
}

func (hdl *Handler) Run() {
	go func() {

		var msg []byte
		var err error

		for {
			if msg, err = hdl.Recv(); err != nil {
				log.Fatal("Cannot recv: %s", err.Error())
			}

			idt := &bollobas.Identity{}
			err := json.Unmarshal(msg, idt)
			if err != nil {
				log.Errorf("Error while receiving message", err)
				continue
			}
			fmt.Printf("%+v\n", idt)
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
		log.Fatal(err)
	}

	mp := map[string]interface{}{}

	err = json.Unmarshal(bts, &mp)

	fmt.Println(mp)

	//Here be temp cipher code..
	fmt.Println(parseid.DecryptString(idt.ID))
}

func NewHandler(name string) *Handler {
	var sock mangos.Socket
	var err error

	if sock, err = sub.NewSocket(); err != nil {
		log.Fatal("can't get new sub socket: %s", err.Error())
	}
	if err = sock.Dial("inproc://driver-publisher"); err != nil {
		log.Fatal("can't dial on sub socket: %s", err.Error())
	}

	if err = sock.Dial("inproc://passenger-publisher"); err != nil {
		log.Fatal("can't dial on sub socket: %s", err.Error())
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		log.Fatal("cannot subscribe: %s", err.Error())
	}

	return &Handler{
		Socket: sock,
	}
}
