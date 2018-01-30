package dolma

import (
	"github.com/kai5263499/dolma/signature"
	"os"
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"time"
	"fmt"
)

// TODO: Implement in-memory file storage with https://github.com/spf13/afero

func NewDolma() *Dolma {
	return &Dolma{
		HasSignature: false,
		Signature: &signature.Signature{
			TotalSectionSize: 0,
			Ts: time.Now().UTC().Unix(),
		},
	}
}

type Dolma struct {
	BinaryName string
	BinaryFileInfo os.FileInfo
	HasSignature bool
	Signature *signature.Signature
}

func (w *Dolma) LoadBinary(binaryName string) error {
	w.BinaryName = binaryName

	fileInfo, err := os.Stat(w.BinaryName)
	if os.IsNotExist(err) {
		return fmt.Errorf("binary %s doesn't exist", w.BinaryName)
	}

	w.BinaryFileInfo = fileInfo

	w.LoadSignature()

	return nil
}

func (w *Dolma) LoadSignature() error {
	fileHandle, _ := os.OpenFile(w.BinaryName, os.O_RDWR, w.BinaryFileInfo.Mode())
	defer fileHandle.Close()

	fileHandle.Seek(-4, 2)

	data := make([]byte, 4)
	fileHandle.Read(data)

	signatureSize := binary.BigEndian.Uint32(data)

	if signatureSize > 0 && int64(signatureSize) < w.BinaryFileInfo.Size() {
		fileHandle.Seek(int64(4 + signatureSize) * -1, 2)
		sigBytes := make([]byte, signatureSize)
		fileHandle.Read(sigBytes)

		err := proto.Unmarshal(sigBytes, w.Signature)
		if err != nil {
			return fmt.Errorf("error unmarshalling err=%#v", err)
		}

		w.HasSignature = true
	} else {
		w.Signature.BinSize = w.BinaryFileInfo.Size()
	}

	return nil
}

func (w *Dolma) StripSections() error {
	if !w.HasSignature {
		logrus.Errorf("refusing to strip a file with no signature")
		return nil
	}

	if w.Signature.BinSize < 1 {
		logrus.Errorf("refusing to strip a file an invalid BinSize of %d", w.Signature.BinSize)
		return nil
	}

	w.HasSignature = false

	logrus.Debugf("truncating %s to %d", w.BinaryName, w.Signature.BinSize)

	return os.Truncate(w.BinaryName, w.Signature.BinSize)
}

func (w *Dolma) SaveSignature() error {
	signatureOffset := w.Signature.BinSize + w.Signature.TotalSectionSize

	logrus.Debugf("truncating %s to %d", signatureOffset)

	os.Truncate(w.BinaryName, signatureOffset)

	fileHandle, _ := os.OpenFile(w.BinaryName, os.O_RDWR, w.BinaryFileInfo.Mode())
	defer fileHandle.Close()

	data, _ := proto.Marshal(w.Signature)

	logrus.Debugf("write signature %#v", w.Signature)
	fileHandle.Seek(1, 2)

	logrus.Debugf("now writing %d bytes", len(data))

	fileHandle.Write(data)

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(len(data)))
	logrus.Debugf("now writing signature lengh %#v", bs)
	fileHandle.Write(bs)

	w.HasSignature = true

	return nil
}
