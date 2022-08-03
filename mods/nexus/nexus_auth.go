package nexus

/*
import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/url"
)

func authorize() error {
	u := url.URL{
		Scheme: "wss://",
		Host:   "sso.nexusmods.com",
		Path:   "/",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	// todo try to get this from secret
	body := msg{
		Uuid:     uuid.New().String(),
		Token:    "",
		Protocol: 2,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}
	fyne.CurrentApp().OpenURL("https://www.nexusmods.com/sso?id=" + uuid + "&application=" + application_slug)
}

type msg struct {
	Uuid     string `json:"uuid"`
	Token    string `json:"token"`
	Protocol int    `json:"protocol"`
}
*/
