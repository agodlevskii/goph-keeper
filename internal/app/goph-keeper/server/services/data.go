package services

import (
	"context"
	"encoding/json"

	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
	"github.com/agodlevskii/goph-keeper/internal/pkg/enc"
)

func StoreSecureDataFromPayload(ctx context.Context, db storage.IDataRepo,
	uid string, payload any, t storage.Type) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encData, err := enc.EncryptData(data)
	if err != nil {
		return "", err
	}

	sd := storage.SecureData{
		UID:  uid,
		Data: encData,
		Type: t,
	}
	return db.StoreData(ctx, sd)
}

func DeleteSecureData(ctx context.Context, db storage.IDataRepo, uid, id string) error {
	return db.DeleteData(ctx, uid, id)
}

func GetDataFromBytes(b []byte, t storage.Type) (any, error) {
	db, err := enc.DecryptData(b)
	if err != nil {
		return nil, err
	}

	return getDataOfType(db, t)
}

func getDataOfType(data []byte, t storage.Type) (any, error) {
	var (
		res any
		err error
	)

	switch t {
	case storage.SBinary:
		var d BinaryRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	case storage.SCard:
		var d CardRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	case storage.SPassword:
		var d PasswordRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	case storage.SText:
		var d TextRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	}

	if err != nil {
		return nil, enc.ErrDecryption
	}
	return res, nil
}
