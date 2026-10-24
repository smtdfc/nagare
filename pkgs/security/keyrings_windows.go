//go:build windows

package security

import (
	"errors"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/smtdfc/nagare/pkgs/paths"
	"golang.org/x/sys/windows"
)

func SaveKey(serviceName string, key string, value string) error {
	return saveDPAPIKey(serviceName, key, value)
}

func GetKey(serviceName string, key string) (string, error) {
	return getDPAPIKey(serviceName, key)
}

func saveDPAPIKey(serviceName, key, value string) error {
	path := filepath.Join(paths.CredentialsDir, serviceName+"_"+key)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	encrypted, err := encryptDPAPI([]byte(value))
	if err != nil {
		return err
	}

	return os.WriteFile(path, encrypted, 0600)
}

func getDPAPIKey(serviceName, key string) (string, error) {
	path := filepath.Join(paths.CredentialsDir, serviceName+"_"+key)

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	decrypted, err := decryptDPAPI(data)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func bytesToBlob(b []byte) *windows.DataBlob {
	if len(b) == 0 {
		return &windows.DataBlob{Size: 0, Data: nil}
	}
	return &windows.DataBlob{
		Size: uint32(len(b)),
		Data: &b[0],
	}
}

func encryptDPAPI(data []byte) ([]byte, error) {
	var out windows.DataBlob
	err := windows.CryptProtectData(bytesToBlob(data), nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out)
	if err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))

	ret := make([]byte, out.Size)
	copy(ret, (*[1 << 31]byte)(unsafe.Pointer(out.Data))[:out.Size:out.Size])
	return ret, nil
}

func decryptDPAPI(data []byte) ([]byte, error) {
	var out windows.DataBlob
	var outString *uint16
	err := windows.CryptUnprotectData(bytesToBlob(data), &outString, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out)
	if err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	if outString != nil {
		windows.LocalFree(windows.Handle(unsafe.Pointer(outString)))
	}

	ret := make([]byte, out.Size)
	copy(ret, (*[1 << 31]byte)(unsafe.Pointer(out.Data))[:out.Size:out.Size])
	return ret, nil
}
