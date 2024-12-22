// Copyright (c) 2019 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package spanstore

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaegertracing/jaeger/model"
)

/*
	Additional cache store tests that need to access internal parts. As such, package must be spanstore and not spanstore_test
*/

func TestExpiredItems(t *testing.T) {
	runWithBadger(t, func(store *badger.DB, t *testing.T) {
		cache := NewCacheStore(store, time.Duration(-1*time.Hour), false)

		expireTime := uint64(time.Now().Add(cache.ttl).Unix())

		// Expired service

		cache.Update("service1", "op1", model.SpanKindUnspecified, expireTime)
		cache.Update("service1", "op2", model.SpanKindUnspecified, expireTime)

		services, err := cache.GetServices()
		require.NoError(t, err)
		assert.Empty(t, services) // Everything should be expired

		// Expired service for operations

		cache.Update("service1", "op1", model.SpanKindUnspecified, expireTime)
		cache.Update("service1", "op2", model.SpanKindUnspecified, expireTime)

		operations, err := cache.GetOperations("service1", nil)
		require.NoError(t, err)
		assert.Empty(t, operations) // Everything should be expired

		// Expired operations, stable service

		cache.Update("service1", "op1", model.SpanKindUnspecified, expireTime)
		cache.Update("service1", "op2", model.SpanKindUnspecified, expireTime)

		cache.services["service1"] = uint64(time.Now().Unix() + 1e10)

		operations, err = cache.GetOperations("service1", nil)
		require.NoError(t, err)
		assert.Empty(t, operations) // Everything should be expired
	})
}

func TestCacheStore_GetOperations(t *testing.T) {
	runWithBadger(t, func(store *badger.DB, t *testing.T) {
		cache := NewCacheStore(store, 1*time.Hour, false)
		expireTime := uint64(time.Now().Add(cache.ttl).Unix())
		cache.services = make(map[string]uint64)
		serviceName := "service1"
		operationName := "op1"
		cache.services[serviceName] = uint64(time.Now().Unix() + 1e10)
		cache.operations = make(map[string]map[model.SpanKind]map[string]uint64)
		cache.operations[serviceName] = make(map[model.SpanKind]map[string]uint64)
		for i := 0; i <= 5; i++ {
			cache.operations[serviceName][model.SpanKind(i)] = make(map[string]uint64)
			cache.operations[serviceName][model.SpanKind(i)][operationName] = expireTime
		}
		operations, err := cache.GetOperations(serviceName, nil)
		require.NoError(t, err)
		assert.Len(t, operations, 6)
		var kinds []string
		for i := 0; i <= 5; i++ {
			kinds = append(kinds, model.SpanKind(i).String())
		}
		// This is necessary as we want to check whether the result is sorted or not
		sort.Strings(kinds)
		for i := 0; i <= 5; i++ {
			assert.Equal(t, kinds[i], operations[i].SpanKind)
			assert.Equal(t, operationName, operations[i].Name)
			k := model.SpanKind(i)
			kp := &k
			singleKindOperations, err := cache.GetOperations(serviceName, kp)
			require.NoError(t, err)
			assert.Len(t, singleKindOperations, 1)
			assert.Equal(t, model.SpanKind(i).String(), singleKindOperations[0].SpanKind)
			assert.Equal(t, operationName, singleKindOperations[0].Name)
		}
	})
}

func TestCacheStore_Update(t *testing.T) {
	runWithBadger(t, func(store *badger.DB, t *testing.T) {
		cache := NewCacheStore(store, 1*time.Hour, false)
		expireTime := uint64(time.Now().Add(cache.ttl).Unix())
		serviceName := "service1"
		operationName := "op1"
		for i := 0; i <= 5; i++ {
			cache.Update(serviceName, operationName, model.SpanKind(i), expireTime)
			assert.Equal(t, expireTime, cache.operations[serviceName][model.SpanKind(i)][operationName])
		}
	})
}

func TestCacheStore_NoPrefillButMigration(t *testing.T) {
	testPrefillAndMigration(t, false)
}

func TestCacheStore_PrefillAndMigration(t *testing.T) {
	testPrefillAndMigration(t, true)
}

func testPrefillAndMigration(t *testing.T, prefill bool) {
	runWithBadger(t, func(store *badger.DB, t *testing.T) {
		var spans []*model.Span
		for i := 0; i < 6; i++ {
			service := fmt.Sprintf("service%d", i)
			operation := fmt.Sprintf("operation%d", i)
			kind := model.SpanKind(i)
			span := createDummySpanWithKind(service, operation, kind)
			spans = append(spans, span)
			// This is to create old data
			err := writeOldSpan(store, span, 1*time.Hour)
			require.NoError(t, err)
		}
		cache := NewCacheStore(store, 1*time.Hour, prefill)
		checkIfCacheContainsKeys(t, store, spans, prefill, cache)
	})
}

func TestCacheStore_PrefillButNoMigration(t *testing.T) {
	runWithBadger(t, func(store *badger.DB, t *testing.T) {
		cache := NewCacheStore(store, 1*time.Hour, false)
		writer := NewSpanWriter(store, cache, 1*time.Hour)
		var spans []*model.Span
		for i := 0; i < 6; i++ {
			service := fmt.Sprintf("service%d", i)
			operation := fmt.Sprintf("operation%d", i)
			kind := model.SpanKind(i)
			span := createDummySpanWithKind(service, operation, kind)
			spans = append(spans, span)
			// This is to create old data
			err := writer.WriteSpan(context.Background(), span)
			require.NoError(t, err)
		}
		// To check whether the WriteSpan is correctly updating or not
		checkIfCacheContainsKeys(t, store, spans, true, cache)
		// We need to create a new cache to see if data is loaded in that cache or not
		newCache := NewCacheStore(store, 1*time.Hour, true)
		checkIfCacheContainsKeys(t, store, spans, true, newCache)
	})
}

func checkIfCacheContainsKeys(t *testing.T, store *badger.DB, spans []*model.Span, prefill bool, cache *CacheStore) {
	for i := 0; i < 6; i++ {
		_, foundService := cache.services[fmt.Sprintf("service%d", i)]
		if prefill {
			assert.True(t, foundService)
		} else {
			assert.False(t, foundService)
		}
		_, foundOperation := cache.operations[fmt.Sprintf("service%d", i)][model.SpanKind(i)][fmt.Sprintf("operation%d", i)]
		if prefill {
			assert.True(t, foundOperation)
		} else {
			assert.False(t, foundOperation)
		}
		v := fmt.Sprintf("8-service%doperation%d%d", i, i, i)
		startTime := model.TimeAsEpochMicroseconds(spans[i].StartTime)
		key := createIndexKey(spanKindIndexKey, []byte(v), startTime, spans[i].TraceID)
		store.View(func(txn *badger.Txn) error {
			_, err := txn.Get(key)
			require.NoError(t, err)
			return nil
		})
	}
	store.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(migrationKey))
		require.NoError(t, err)
		return nil
	})
}

func TestOldReads(t *testing.T) {
	runWithBadger(t, func(store *badger.DB, t *testing.T) {
		timeNow := model.TimeAsEpochMicroseconds(time.Now())
		s1o1Key := createIndexKey(spanKindIndexKey, []byte("8-service1operation10"), timeNow, model.TraceID{High: 0, Low: 0})

		tid := time.Now().Add(1 * time.Minute)

		writer := func() {
			store.Update(func(txn *badger.Txn) error {
				txn.SetEntry(&badger.Entry{
					Key:       s1o1Key,
					ExpiresAt: uint64(tid.Unix()),
				})
				return nil
			})
		}

		cache := NewCacheStore(store, time.Duration(-1*time.Hour), false)
		writer()

		nuTid := tid.Add(1 * time.Hour)

		cache.Update("service1", "operation1", model.SpanKindUnspecified, uint64(tid.Unix()))
		cache.services["service1"] = uint64(nuTid.Unix())
		cache.operations["service1"][model.SpanKindUnspecified]["operation1"] = uint64(nuTid.Unix())

		cache.populateCaches()

		// Now make sure we didn't use the older timestamps from the DB
		assert.Equal(t, uint64(nuTid.Unix()), cache.services["service1"])
		assert.Equal(t, uint64(nuTid.Unix()), cache.operations["service1"][model.SpanKindUnspecified]["operation1"])
	})
}

// func runFactoryTest(tb testing.TB, test func(tb testing.TB, sw spanstore.Writer, sr spanstore.Reader)) {
func runWithBadger(t *testing.T, test func(store *badger.DB, t *testing.T)) {
	opts := badger.DefaultOptions("")

	opts.SyncWrites = false
	dir := t.TempDir()
	opts.Dir = dir
	opts.ValueDir = dir

	store, err := badger.Open(opts)
	defer func() {
		store.Close()
	}()

	require.NoError(t, err)

	test(store, t)
}

func createDummySpanWithKind(service string, operation string, kind model.SpanKind) *model.Span {
	tid := time.Now()
	testSpan := model.Span{
		TraceID: model.TraceID{
			Low:  uint64(0),
			High: 1,
		},
		SpanID:        model.SpanID(0),
		OperationName: operation,
		Process: &model.Process{
			ServiceName: service,
		},
		StartTime: tid.Add(1 * time.Millisecond),
		Duration:  1 * time.Millisecond,
		Tags: model.KeyValues{
			model.KeyValue{
				Key:   "span.kind",
				VType: model.StringType,
				VStr:  kind.String(),
			},
		},
	}
	return &testSpan
}

func writeOldSpan(store *badger.DB, span *model.Span, ttl time.Duration) error {
	expireTime := uint64(time.Now().Add(ttl).Unix())
	startTime := model.TimeAsEpochMicroseconds(span.StartTime)
	entriesToStore := make([]*badger.Entry, 0, len(span.Tags)+4+len(span.Process.Tags)+len(span.Logs)*4)
	entriesToStore = append(entriesToStore, createBadgerEntry(createIndexKey(serviceNameIndexKey, []byte(span.Process.ServiceName), startTime, span.TraceID), nil, expireTime))
	entriesToStore = append(entriesToStore, createBadgerEntry(createIndexKey(operationNameIndexKey, []byte(span.Process.ServiceName+span.OperationName), startTime, span.TraceID), nil, expireTime))
	for _, kv := range span.Tags {
		entriesToStore = append(entriesToStore, createBadgerEntry(createIndexKey(tagIndexKey, []byte(span.Process.ServiceName+kv.Key+kv.AsString()), startTime, span.TraceID), nil, expireTime))
	}
	err := store.Update(func(txn *badger.Txn) error {
		for i := range entriesToStore {
			err := txn.SetEntry(entriesToStore[i])
			if err != nil {
				// Most likely primary key conflict, but let the caller check this
				return err
			}
		}
		return nil
	})
	return err
}

func createBadgerEntry(key []byte, value []byte, expireTime uint64) *badger.Entry {
	return &badger.Entry{
		Key:       key,
		Value:     value,
		ExpiresAt: expireTime,
	}
}
