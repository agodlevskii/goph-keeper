package services

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
)

type CardReq struct {
	Name    string `json:"name"`
	Number  string `json:"number"`
	Holder  string `json:"holder"`
	ExpDate string `json:"exp_date"`
	CVV     string `json:"cvv"`
	Note    string `json:"note"`
}

type CardRes struct {
	UID     string `json:"-"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Number  string `json:"number"`
	Holder  string `json:"holder"`
	ExpDate string `json:"exp_date"`
	CVV     string `json:"cvv"`
	Note    string `json:"note"`
}

func GetAllCards(db storage.IDataStorage, uid string) ([]CardRes, error) {
	sd, err := db.GetAllDataByType(uid, storage.SCard)
	if err != nil {
		return nil, err
	}

	cards := make([]CardRes, 0, len(sd))
	for _, d := range sd {
		c, eErr := getCardFromSecureData(d)
		if eErr != nil {
			return nil, eErr
		}

		c.CVV = "***"
		cards = append(cards, c)
	}
	return cards, nil
}

func GetCardByID(db storage.IDataStorage, uid, id string) (CardRes, error) {
	d, err := db.GetDataByID(uid, id)
	if err != nil {
		return CardRes{}, nil
	}
	return getCardFromSecureData(d)
}

func StoreCard(db storage.IDataStorage, uid string, req CardReq) (string, error) {
	card := getCardFromRequest(uid, req)
	return StoreSecureDataFromPayload(db, uid, card, storage.SCard)
}

func getCardFromSecureData(d storage.SecureData) (CardRes, error) {
	c, err := GetDataFromBytes(d.Data, storage.SCard)
	if err != nil {
		return CardRes{}, err
	}

	ct := c.(CardRes)
	ct.ID = d.ID
	return ct, nil
}

func getCardFromRequest(uid string, req CardReq) CardRes {
	return CardRes{
		UID:     uid,
		Name:    req.Name,
		Number:  req.Number,
		Holder:  req.Holder,
		ExpDate: req.ExpDate,
		CVV:     req.CVV,
		Note:    req.Note,
	}
}
