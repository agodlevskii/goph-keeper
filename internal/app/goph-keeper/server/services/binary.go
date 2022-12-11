package services

import (
	"context"

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

func GetAllBinaries(ctx context.Context, db storage.IDataRepo, uid string) ([]BinaryRes, error) {
	sd, err := db.GetAllDataByType(ctx, uid, storage.SBinary)
	if err != nil {
		return nil, err
	}

	binaries := make([]BinaryRes, 0, len(sd))
	for _, d := range sd {
		b, dErr := getBinaryFromSecureData(d)
		if dErr != nil {
			return nil, err
		}
		b.Data = make([]byte, 0)
		binaries = append(binaries, b)
	}

	return binaries, nil
}

func GetBinaryByID(ctx context.Context, db storage.IDataRepo, uid, id string) (BinaryRes, error) {
	d, err := db.GetDataByID(ctx, uid, id)
	if err != nil {
		return BinaryRes{}, err
	}
	return getBinaryFromSecureData(d)
}

func StoreBinary(ctx context.Context, db storage.IDataRepo, uid string, req BinaryReq) (string, error) {
	bin := getBinaryFromRequest(uid, req)
	return StoreSecureDataFromPayload(ctx, db, uid, bin, storage.SBinary)
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
