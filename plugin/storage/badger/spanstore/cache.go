// Copyright (c) 2018 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package spanstore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"go.opentelemetry.io/otel/trace"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

// CacheStore saves expensive calculations from the K/V store
type CacheStore struct {
	// Given the small amount of data these will store, we use the same structure as the memory store
	cacheLock  sync.Mutex // write heavy - Mutex is faster than RWMutex for writes
	services   map[string]uint64
	operations map[string]map[trace.SpanKind]map[string]uint64

	store *badger.DB
	ttl   time.Duration
}

// NewCacheStore returns initialized CacheStore for badger use
func NewCacheStore(db *badger.DB, ttl time.Duration, prefill bool) *CacheStore {
	cs := &CacheStore{
		services:   make(map[string]uint64),
		operations: make(map[string]map[trace.SpanKind]map[string]uint64),
		ttl:        ttl,
		store:      db,
	}
	cs.cacheLock.Lock()
	defer cs.cacheLock.Unlock()
	if prefill {
		loaded := autoMigrate(cs.store, cs.services, cs.operations)
		if !loaded {
			cs.loadData()
		}
	} else {
		services := make(map[string]uint64)
		operations := make(map[string]map[trace.SpanKind]map[string]uint64)
		autoMigrate(cs.store, services, operations)
	}
	return cs
}

func (c *CacheStore) populateCaches() {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	c.loadData()
}

func autoMigrate(store *badger.DB, services map[string]uint64, operations map[string]map[trace.SpanKind]map[string]uint64) bool {
	needToMigrate := false
	loaded := false
	err := store.View(func(txn *badger.Txn) error {
		migrationKeyBytes := []byte(migrationKey)
		_, err := txn.Get(migrationKeyBytes)
		if err == nil {
			// No need to migrate
			return nil
		} else if errors.Is(err, badger.ErrKeyNotFound) {
			needToMigrate = true
		}
		return err
	})

	if needToMigrate {
		err = store.Update(func(txn *badger.Txn) error {
			iter := txn.NewIterator(badger.DefaultIteratorOptions)
			defer iter.Close()
			createNewIndexForOldData(txn, iter, services, operations)
			createMigrationKey(txn)
			loaded = true
			return nil
		})
	}

	if err != nil {
		return false
	}

	return loaded
}

func (c *CacheStore) loadData() {
	c.store.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		serviceKey := []byte{spanKindIndexKey}

		// Seek all the services first
		for it.Seek(serviceKey); it.ValidForPrefix(serviceKey); it.Next() {
			timestampStartIndex := len(it.Item().Key()) - (sizeOfTraceID + 8) // 8 = sizeof(uint64)
			serviceNameWithOperation := string(it.Item().Key()[len(serviceKey):timestampStartIndex])
			// This string is of type len(serviceName)-serviceName+operationName+spanKind
			operationKindString := string(serviceNameWithOperation[len(serviceNameWithOperation)-1])
			serviceNameLengthIndex := strings.Index(serviceNameWithOperation, "-")
			if serviceNameLengthIndex == -1 {
				// Inconsistency we can't load this in cache, just skip it!
				continue
			}
			serviceNameLengthString := serviceNameWithOperation[:serviceNameLengthIndex]
			serviceNameLength, err := strconv.Atoi(serviceNameLengthString)
			if err != nil {
				// Inconsistency we can't load this in cache, just skip it!
				continue
			}
			serviceNameWithOperation = serviceNameWithOperation[serviceNameLengthIndex+1:]
			serviceName := serviceNameWithOperation[:serviceNameLength]
			operationName := serviceNameWithOperation[serviceNameLength : len(serviceNameWithOperation)-1]
			spanKind := model.GetSpanKindFromStringOfSpanKind(operationKindString)
			keyTTL := it.Item().ExpiresAt()
			if _, found := c.operations[serviceName]; !found {
				c.operations[serviceName] = make(map[trace.SpanKind]map[string]uint64)
			}
			if _, found := c.operations[serviceName][spanKind]; !found {
				c.operations[serviceName][spanKind] = make(map[string]uint64)
			}
			if v, found := c.services[serviceName]; found {
				if v > keyTTL {
					continue
				}
			}
			c.services[serviceName] = keyTTL
			if t, found := c.operations[serviceName][spanKind][operationName]; found {
				if t > keyTTL {
					continue
				}
			}
			c.operations[serviceName][spanKind][operationName] = keyTTL
		}
		return nil
	})
}

func createNewIndexForOldData(txn *badger.Txn, it *badger.Iterator, services map[string]uint64, operations map[string]map[trace.SpanKind]map[string]uint64) {
	serviceNameIndexKeyPrefix := []byte{serviceNameIndexKey}
	// Load all the services in cache
	for it.Seek(serviceNameIndexKeyPrefix); it.ValidForPrefix(serviceNameIndexKeyPrefix); it.Next() {
		if it.Valid() {
			timestampStartIndex := len(it.Item().Key()) - (sizeOfTraceID + 8) // 8 = sizeof(uint64)
			serviceName := string(it.Item().Key()[len(serviceNameIndexKeyPrefix):timestampStartIndex])
			keyTTL := it.Item().ExpiresAt()
			if v, found := services[serviceName]; found {
				if v > keyTTL {
					continue
				}
			}
			services[serviceName] = keyTTL
		}
	}
	for service := range services {
		loadAndCreateIndexForOperation(txn, it, service, operations)
	}
}

func loadAndCreateIndexForOperation(txn *badger.Txn, it *badger.Iterator, service string, operations map[string]map[trace.SpanKind]map[string]uint64) {
	operationKey := make([]byte, len(service)+1)
	operationKey[0] = operationNameIndexKey
	copy(operationKey[1:], service)
	for it.Seek(operationKey); it.ValidForPrefix(operationKey); it.Next() {
		if it.Valid() {
			timestampStartIndex := len(it.Item().Key()) - (sizeOfTraceID + 8)
			operation := string(it.Item().Key()[len(operationKey):timestampStartIndex])
			timeStampAndTraceId := string(it.Item().Key()[timestampStartIndex:])
			kind := getSpanKind(txn, service, timeStampAndTraceId)
			keyTTL := it.Item().ExpiresAt()
			if _, found := operations[service]; !found {
				operations[service] = make(map[trace.SpanKind]map[string]uint64)
			}
			if _, found := operations[service][kind]; !found {
				operations[service][kind] = make(map[string]uint64)
			}
			key := fmt.Sprintf("%d-", len(service)) + service + operation + fmt.Sprintf("%d", rune(kind))
			err := createSpanKindIndex(txn, []byte(key), it.Item().Key()[timestampStartIndex:], keyTTL)
			if err != nil {
				return
			}
			if v, found := operations[service][kind][operation]; found {
				if v > keyTTL {
					continue
				}
			}
			operations[service][kind][operation] = keyTTL
		}
	}
}

func getSpanKind(txn *badger.Txn, service string, timestampAndTraceId string) trace.SpanKind {
	for i := 0; i < 6; i++ {
		value := service + model.SpanKindKey + trace.SpanKind(i).String()
		valueBytes := []byte(value)
		operationKey := make([]byte, 1+len(valueBytes)+8+sizeOfTraceID)
		operationKey[0] = tagIndexKey
		copy(operationKey[1:], valueBytes)
		copy(operationKey[1+len(valueBytes):], timestampAndTraceId)
		_, err := txn.Get(operationKey)
		if err == nil {
			return trace.SpanKind(i)
		}
	}
	return trace.SpanKindUnspecified
}

func createMigrationKey(txn *badger.Txn) {
	txn.SetEntry(&badger.Entry{
		Key:   []byte(migrationKey),
		Value: nil,
	})
}

func createSpanKindIndex(txn *badger.Txn, key []byte, timeStampAndTraceId []byte, ttL uint64) error {
	timeStamp := timeStampAndTraceId[:len(timeStampAndTraceId)-sizeOfTraceID]
	timeStamp64 := binary.BigEndian.Uint64(timeStamp)
	traceId, err := model.TraceIDFromBytes(timeStampAndTraceId[len(timeStampAndTraceId)-sizeOfTraceID:])
	if err != nil {
		return err
	}
	indexKey := createIndexKey(spanKindIndexKey, key, timeStamp64, traceId)
	err = txn.SetEntry(&badger.Entry{
		Key:       indexKey,
		Value:     nil,
		ExpiresAt: ttL,
	})
	return err
}

// Update caches the results of service and service + operation indexes and maintains their TTL
func (c *CacheStore) Update(service, operation string, kind trace.SpanKind, expireTime uint64) {
	c.cacheLock.Lock()
	c.services[service] = expireTime
	if _, ok := c.operations[service]; !ok {
		c.operations[service] = make(map[trace.SpanKind]map[string]uint64)
	}
	if _, ok := c.operations[service][kind]; !ok {
		c.operations[service][kind] = make(map[string]uint64)
	}
	c.operations[service][kind][operation] = expireTime
	c.cacheLock.Unlock()
}

// GetOperations returns all operations for a specific service & spanKind traced by Jaeger
func (c *CacheStore) GetOperations(service string, kind *trace.SpanKind) ([]spanstore.Operation, error) {
	operations := make([]spanstore.Operation, 0, len(c.services))
	//nolint: gosec // G115
	t := uint64(time.Now().Unix())
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	if v, ok := c.services[service]; ok {
		if v < t {
			// Expired, remove
			delete(c.services, service)
			delete(c.operations, service)
			return []spanstore.Operation{}, nil // empty slice rather than nil
		}
		if kind != nil {
			for o, e := range c.operations[service][*kind] {
				operations = insertOperations(c, operations, service, o, *kind, t, e)
				sort.Slice(operations, func(i, j int) bool {
					return operations[i].Name < operations[j].Name
				})
			}
		} else {
			for sKind := range c.operations[service] {
				for o, e := range c.operations[service][sKind] {
					operations = insertOperations(c, operations, service, o, sKind, t, e)
					sort.Slice(operations, func(i, j int) bool {
						if operations[i].SpanKind == operations[j].SpanKind {
							return operations[i].Name < operations[j].Name
						}
						return operations[i].SpanKind < operations[j].SpanKind
					})
				}
			}
		}
	}
	return operations, nil
}

func insertOperations(c *CacheStore, operations []spanstore.Operation, service, operation string, kind trace.SpanKind, t, e uint64) []spanstore.Operation {
	if e > t {
		op := spanstore.Operation{Name: operation, SpanKind: kind.String()}
		operations = append(operations, op)
		return operations
	}
	delete(c.operations[service][kind], operation)
	return operations
}

// GetServices returns all services traced by Jaeger
func (c *CacheStore) GetServices() ([]string, error) {
	services := make([]string, 0, len(c.services))
	//nolint: gosec // G115
	t := uint64(time.Now().Unix())
	c.cacheLock.Lock()
	// Fetch the items
	for k, v := range c.services {
		if v > t {
			services = append(services, k)
		} else {
			// Service has expired, remove it
			delete(c.services, k)
		}
	}
	c.cacheLock.Unlock()

	sort.Strings(services)

	return services, nil
}
