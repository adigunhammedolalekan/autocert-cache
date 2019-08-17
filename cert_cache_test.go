package main

import (
	"context"
	"testing"
)

const testDbUrl = "root:passcode@tcp(localhost:3306)/test_db?charset=utf8&parseTime=True"

func TestDbCertificateCache_Get(t *testing.T) {
	ctx := context.Background()
	cache, err := NewDbCache("mysql", testDbUrl)
	if err != nil {
		t.Fatal(err)
	}
	err = cache.Put(ctx, "testKey", []byte("testData"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := cache.Get(ctx, "testKey")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "testData" {
		t.Fatalf("inconsistent value! Expecting %s, got %s", "testData", string(data))
	}
}

func TestDbCertificateCache_Put(t *testing.T) {
	ctx := context.Background()
	cache, err := NewDbCache("mysql", testDbUrl)
	if err != nil {
		t.Fatal(err)
	}
	err = cache.Put(ctx, "testKey1", []byte("testData1"))
	if err != nil {
		t.Fatal(err)
	}
	// test if the key is getting updated
	err = cache.Put(ctx, "testKey1", []byte("updatedTestData"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := cache.Get(ctx, "testKey1")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "updatedTestData" {
		t.Fatalf("inconsistent data! expected %s, got %s", "updatedTestData", string(data))
	}
}

func TestDbCertificateCache_Delete(t *testing.T) {
	ctx := context.Background()
	cache, err := NewDbCache("mysql", testDbUrl)
	if err != nil {
		t.Fatal(err)
	}
	err = cache.Put(ctx, "deleteKey", []byte("testDelete"))
	if err != nil {
		t.Fatal(err)
	}

	err = cache.Delete(ctx, "deleteKey")
	if err != nil {
		t.Fatal(err)
	}

	data, err := cache.Get(ctx, "deleteKey")
	if err == nil || data != nil {
		t.Fatal("data should be nil!")
	}
}