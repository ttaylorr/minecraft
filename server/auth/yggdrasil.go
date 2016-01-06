package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Yggdrasil struct {
	pubKey []byte
}

func NewYggdrasil(pubKey []byte) *Yggdrasil {
	return &Yggdrasil{
		pubKey: pubKey,
	}
}

func (y *Yggdrasil) GetSession(username string, secret []byte) (*Session, error) {
	hash := HashSession("", secret, y.pubKey)

	res, err := http.Get(y.Url(username, hash))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var session *Session
	if err = json.NewDecoder(res.Body).Decode(&session); err != nil {
		return nil, err
	}

	return session, nil
}

func (y *Yggdrasil) Url(username, session string) string {
	return fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s", username, session)
}
