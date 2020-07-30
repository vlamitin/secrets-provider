package persistence

import (
	"fmt"
	"github.com/vlamitin/secrets-provider/internal/crypt"
)

type Secrets map[string][]byte

var _secretsMap = Secrets(make(map[string][]byte))
var _cryptKey = make([]byte, 32)

func SetCryptKey(cryptKey string) {
	_cryptKey = crypt.PopulateKey(cryptKey)
}

func FromDB(rows []SecretRow) {
	for _, row := range rows {
		_secretsMap[row.key] = []byte(row.value)
	}
}

func CheckCryptKey(cryptKey string) bool {
	return string(_cryptKey) == string(crypt.PopulateKey(cryptKey))
}

func GetSecret(key string) (secret string, notFoundErr error, decryptErr error) {
	v, ok := _secretsMap[key]
	if !ok {
		return "", fmt.Errorf("not found secret with key %s\n", key), nil
	}

	decrypted, err := crypt.Decrypt(v, _cryptKey)
	if err != nil {
		return "", fmt.Errorf("error when decrypt secret: %v\n", err), nil
	}

	return string(decrypted), nil, nil
}

func SetSecret(key string, value string) error {
	encrypted, err := crypt.Encrypt([]byte(value), _cryptKey)
	if err != nil {
		return fmt.Errorf("error when encrypt secret: %v\n", err)
	}
	_secretsMap[key] = encrypted
	if _db != nil {
		DeleteOne(_db, key)
		InsertOne(_db, key, string(encrypted))
	}

	return nil
}

func RemoveSecret(key string) {
	delete(_secretsMap, key)
	if _db != nil {
		DeleteOne(_db, key)
	}
}
