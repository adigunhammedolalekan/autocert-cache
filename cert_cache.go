package autocertcache

import (
	"context"
	"encoding/base64"
	"gorm.io/gorm"
	"sync"
)

type Cert struct {
	CertKey string
	Data    string
}

// DbCertificateCache implements autocert.Cache interface to provide
// a db cache module for our https certs.
type DbCertificateCache struct {
	db            *gorm.DB
	mtx           sync.RWMutex
	inMemoryCache map[string][]byte
}

// NewDbCache creates a *DbCertificate. This function will try to connect
// to db using the provided url and dialect.
func NewDbCache(db *gorm.DB) (*DbCertificateCache, error) {
	db.AutoMigrate(&Cert{})
	c := &DbCertificateCache{db: db}
	c.inMemoryCache = make(map[string][]byte)
	return c, nil
}

// Get returns a cert record for key.
// it first consult the inmemory cache for faster response time
func (c *DbCertificateCache) Get(ctx context.Context, key string) ([]byte, error) {
	// check inMemory cache
	c.mtx.Lock()
	v, ok := c.inMemoryCache[key]
	c.mtx.Unlock()

	if ok && v != nil {
		return v, nil
	}

	cert, err := c.getCert(key)
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(cert.Data)
	if err != nil {
		return nil, err
	}
	// update inMemory cache
	c.mtx.Lock()
	c.inMemoryCache[key] = data
	c.mtx.Unlock()
	return data, nil
}

// getCert gets a cert record by key, straight from database
func (c *DbCertificateCache) getCert(key string) (*Cert, error) {
	cert := &Cert{}
	err := c.db.Table("certs").Where("cert_key = ?", key).First(cert).Error
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// Put insert or updates a new cert record.
// it also update the inMemory cache
func (c *DbCertificateCache) Put(ctx context.Context, key string, data []byte) error {
	_, err := c.getCert(key)
	if err == gorm.ErrRecordNotFound {
		newCert := &Cert{CertKey: key, Data: base64.StdEncoding.EncodeToString(data)}
		if err := c.db.Create(newCert).Error; err != nil {
			return err
		}
	} else {
		if err := c.db.Table("certs").Where("cert_key = ?", key).Update("data", base64.StdEncoding.EncodeToString(data)).Error; err != nil {
			return err
		}
	}
	c.mtx.Lock()
	c.inMemoryCache[key] = data
	c.mtx.Unlock()
	return nil
}

// Delete removes a cert from db
func (c *DbCertificateCache) Delete(ctx context.Context, key string) error {
	err := c.db.Table("certs").Where("cert_key = ?", key).Delete(&Cert{}).Error
	if err != nil {
		return err
	}
	c.mtx.Lock()
	c.inMemoryCache[key] = nil
	c.mtx.Unlock()
	return nil
}
