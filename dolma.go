package dolma

import (
	"archive/zip"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	pb "github.com/kai5263499/dolma/generated"
	uuid "github.com/satori/go.uuid"
)

// TODO: Implement in-memory file storage with https://github.com/spf13/afero

func NewDolma() *Dolma {
	return &Dolma{
		Signature: &pb.Signature{
			Ts:               time.Now().UTC().UnixNano(),
			Sections:         make([]*pb.Section, 0),
			TotalSectionSize: 0,
		},
	}
}

type Dolma struct {
	BinaryName            string
	BinaryFileInfo        os.FileInfo
	Signature             *pb.Signature
	SectionPrefixIndexMap map[string]int
}

func (d *Dolma) LoadBinary(binaryName string) error {
	d.BinaryName = binaryName

	fileInfo, err := os.Stat(d.BinaryName)
	if os.IsNotExist(err) {
		return fmt.Errorf("binary %s doesn't exist", d.BinaryName)
	}

	d.BinaryFileInfo = fileInfo

	d.LoadSignature()

	return nil
}

func (d *Dolma) LoadSignature() error {
	fileHandle, _ := os.OpenFile(d.BinaryName, os.O_RDWR, d.BinaryFileInfo.Mode())
	defer fileHandle.Close()

	fileHandle.Seek(-4, 2)

	data := make([]byte, 4)
	fileHandle.Read(data)

	signatureSize := int64(binary.BigEndian.Uint32(data))

	if signatureSize > 0 && signatureSize < d.BinaryFileInfo.Size() {
		fileHandle.Seek(int64(4+signatureSize)*-1, 2)
		sigBytes := make([]byte, signatureSize)
		fileHandle.Read(sigBytes)

		err := proto.Unmarshal(sigBytes, d.Signature)
		if err != nil {
			d.Signature.BinSize = d.BinaryFileInfo.Size()
			return err
		}
	} else {
		d.Signature.BinSize = d.BinaryFileInfo.Size()
	}

	d.buildSectionPrefixIndexMap()

	return nil
}

func (d *Dolma) buildSectionPrefixIndexMap() {
	d.SectionPrefixIndexMap = make(map[string]int)
	for idx, section := range d.Signature.Sections {
		d.SectionPrefixIndexMap[section.Prefix] = idx
	}
}

func (d *Dolma) StripSections() error {
	if d.Signature.BinSize < 1 {
		logrus.Errorf("refusing to strip a file an invalid BinSize of %d", d.Signature.BinSize)
		return nil
	}

	logrus.Debugf("truncating %s to %d", d.BinaryName, d.Signature.BinSize)

	return os.Truncate(d.BinaryName, d.Signature.BinSize)
	return nil
}

func (d *Dolma) SaveSignature() error {
	signatureOffset := d.Signature.BinSize + d.Signature.TotalSectionSize

	logrus.Debugf("truncating %s to %d", signatureOffset)

	os.Truncate(d.BinaryName, signatureOffset)

	fileHandle, _ := os.OpenFile(d.BinaryName, os.O_RDWR, d.BinaryFileInfo.Mode())
	defer fileHandle.Close()

	data, _ := proto.Marshal(d.Signature)

	fileHandle.Seek(1, 2)
	fileHandle.Write(data)

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(len(data)))
	logrus.Debugf("now writing signature lengh %#v", bs)
	fileHandle.Write(bs)

	return nil
}

func (d *Dolma) AddContent(prefix string, contentPath string) error {
	logrus.Debugf("add content prefix=%s contentPath=%s", prefix, contentPath)

	var section *pb.Section
	tmpID, _ := uuid.NewV4()
	tmpSectionFileName := filepath.Join(os.TempDir(), tmpID.String())
	tmpSectionFile, _ := os.Create(tmpSectionFileName)

	logrus.Debugf("tmpSectionFileName=%s", tmpSectionFileName)

	if idx, ok := d.SectionPrefixIndexMap[prefix]; ok {
		section = d.Signature.Sections[idx]
		// Copy section data out to a file

		fileHandle, _ := os.OpenFile(d.BinaryName, os.O_RDWR, d.BinaryFileInfo.Mode())

		fileHandle.Seek(section.Offset, 0)

		data := make([]byte, section.Size)
		n, err := fileHandle.Read(data)
		logrus.Debugf("read %d bytes err=%s", n, err)

		fileHandle.Close()

		tmpSectionFile.Write(data)
	} else {
		section = &pb.Section{
			Prefix: prefix,
			Type:   pb.SectionType_ZIP,
		}
	}

	d.zip(prefix, contentPath, tmpSectionFile)

	tmpSectionFile.Close()

	return nil
}

func (d *Dolma) zip(prefix, source string, zipFile *os.File) error {
	var err error

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	archive.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.Join(prefix, strings.TrimPrefix(path, source))
		header.Method = zip.Deflate

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
