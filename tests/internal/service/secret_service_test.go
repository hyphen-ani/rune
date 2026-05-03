package service

import (
	"os"
	"rune/internal/crypto"
	"rune/internal/seal"
	"rune/internal/service"
	"rune/internal/storage"
	"testing"
)

func setupTestService(t *testing.T) *service.SecretService {
	dbPath := "test.db"
	store, err := storage.NewBoltStore(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(dbPath)
	})

	err = store.CreateNamespace("default")
	if err != nil {
		t.Fatal(err)
	}

	sealer := seal.NewManger()

	salt, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatal(err)
	}

	key := crypto.DeriveKey("test-passphrase", salt)
	sealer.Unseal(key)

	return service.NewSecretService(store, sealer)
}

func TestPutFailsWithoutNamespace(t *testing.T) {

	dbPath := "test.db"
	store, _ := storage.NewBoltStore(dbPath)
	defer os.Remove(dbPath)

	sealer := seal.NewManger()
	sealer.Unseal([]byte("key"))

	svc := service.NewSecretService(store, sealer)
	err := svc.Put("default", "db/pass", "123")
	if err == nil {
		t.Fatal("expected error when namespace not created")
	}

}

func TestPutAndGetSecret(t *testing.T) {

	svc := setupTestService(t)
	err := svc.Put("default", "db/pass", "123")
	if err != nil {
		t.Fatal(err)
	}

	val, err := svc.Get("default", "db/pass")
	if err != nil {
		t.Fatal(err)
	}

	if val != "123" {
		t.Fatalf("expected 123, got %s", val)
	}
}
