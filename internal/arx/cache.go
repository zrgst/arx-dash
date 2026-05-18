package arx

import (
	"context"
	"sync"
	"time"
)

// Cache holder siste export fra ARX i minnet.
//
// Dette gjør at dashboardet ikke trenger å spørre ARX-serveren for hver request.
// Foreløpig er dette kun memory-cache. Restert av appen tømmer cache.
type Cache struct {
	client *Client

	mu        sync.RWMutex
	export    PersonsExport
	loaded    bool
	updatedAt time.Time
}

// NewCache lager en ny ARX-Cache.
func NewCache(client *Client) *Cache {
	return &Cache{
		client: client,
	}
}

// GetPersonsExport returnerer cached export.
// Hvis cache er tom, hentes data fra ARX først.
func (c *Cache) GetPersonsExport(ctx context.Context) (PersonsExport, error) {
	c.mu.RLock()

	if c.loaded {
		export := c.export
		c.mu.RUnlock()
		return export, nil
	}

	c.mu.RUnlock()

	return c.RefreshPersonsExport(ctx)
}

// RefreshPersonsExport tvinger ny henting fra ARX.
func (c *Cache) RefreshPersonsExport(ctx context.Context) (PersonsExport, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	export, err := c.client.ExportPersons(ctx)
	if err != nil {
		return PersonsExport{}, err
	}

	c.export = export
	c.loaded = true
	c.updatedAt = time.Now()

	return export, nil
}

// UpdatedAt returnerer når cache sist ble oppdatert.
func (c *Cache) UpdatedAt() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.updatedAt
}

// Loaded sier om cache har data.
func (c *Cache) Loaded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.loaded
}
