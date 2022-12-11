package services

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
)

type TextReq struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Note string `json:"note"`
}

type TextRes struct {
	UID  string `json:"-"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
	Note string `json:"note"`
}

func GetAllTexts(db storage.IDataRepo, uid string) ([]TextRes, error) {
	sd, err := db.GetAllDataByType(uid, storage.SText)
	if err != nil {
		return nil, err
	}

	texts := make([]TextRes, 0, len(sd))
	for _, d := range sd {
		t, dErr := getTextFromSecureData(d)
		if dErr != nil {
			return nil, err
		}
		t.Data = ""
		texts = append(texts, t)
	}

	return texts, nil
}

func GetTextByID(db storage.IDataRepo, uid, id string) (TextRes, error) {
	sd, err := db.GetDataByID(uid, id)
	if err != nil {
		return TextRes{}, err
	}
	return getTextFromSecureData(sd)
}

func StoreText(db storage.IDataRepo, uid string, req TextReq) (string, error) {
	text := getTextFromRequest(uid, req)
	return StoreSecureDataFromPayload(db, uid, text, storage.SPassword)
}

func getTextFromSecureData(d storage.SecureData) (TextRes, error) {
	t, err := GetDataFromBytes(d.Data, storage.SText)
	if err != nil {
		return TextRes{}, err
	}

	tt := t.(TextRes)
	tt.ID = d.ID
	return tt, nil
}

func getTextFromRequest(uid string, req TextReq) TextRes {
	return TextRes{
		UID:  uid,
		Name: req.Name,
		Data: req.Data,
		Note: req.Note,
	}
}
