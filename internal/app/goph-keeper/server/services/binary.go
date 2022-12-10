package services

import (
	"github.com/agodlevskii/goph-keeper/internal/app/goph-keeper/server/storage"
)

type BinaryReq struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
	Note string `json:"note"`
}

type BinaryRes struct {
	UID  string `json:"-"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Data []byte `json:"data"`
	Note string `json:"note"`
}

func GetAllBinaries(db storage.IDataRepository, uid string) ([]BinaryRes, error) {
	sd, err := db.GetAllDataByType(uid, storage.SBinary)
	if err != nil {
		return nil, err
	}

	binaries := make([]BinaryRes, 0, len(sd))
	for _, d := range sd {
		b, dErr := getBinaryFromSecureData(d)
		if dErr != nil {
			return nil, err
		}
		binaries = append(binaries, b)
	}

	return binaries, nil
}

func GetBinaryByID(db storage.IDataRepository, uid, id string) (BinaryRes, error) {
	d, err := db.GetDataByID(uid, id)
	if err != nil {
		return BinaryRes{}, err
	}
	return getBinaryFromSecureData(d)
}

func StoreBinary(db storage.IDataRepository, uid string, req BinaryReq) (string, error) {
	bin := getBinaryFromRequest(uid, req)
	return StoreSecureDataFromPayload(db, uid, bin, storage.SBinary)
}

func getBinaryFromSecureData(d storage.SecureData) (BinaryRes, error) {
	b, err := GetDataFromBytes(d.Data, storage.SBinary)
	if err != nil {
		return BinaryRes{}, err
	}

	bt := b.(BinaryRes)
	bt.ID = d.ID
	return bt, nil
}

func getBinaryFromRequest(uid string, req BinaryReq) BinaryRes {
	return BinaryRes{
		UID:  uid,
		Name: req.Name,
		Data: req.Data,
		Note: req.Note,
	}
}
