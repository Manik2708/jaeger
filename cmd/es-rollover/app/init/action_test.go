// Copyright (c) 2021 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jaegertracing/jaeger/cmd/es-rollover/app"
	"github.com/jaegertracing/jaeger/pkg/es/client"
	"github.com/jaegertracing/jaeger/pkg/es/client/mocks"
)

func TestIndexCreateIfNotExist(t *testing.T) {
	t.Run("index exists", func(t *testing.T) {
		testCreateIndexIfNotExist(t, "index")
	})
	t.Run("alias exists", func(t *testing.T) {
		testCreateIndexIfNotExist(t, "alias")
	})
	tests := []struct {
		name                     string
		createIndexReturnErr     error
		createIndexExpectedError error
	}{
		{
			name:                     "generic error from create index",
			createIndexReturnErr:     errors.New("may be an http error from create index"),
			createIndexExpectedError: errors.New("may be an http error from create index"),
		},
		{
			name: "success",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexClient := &mocks.IndexAPI{}
			indexClient.On("IndexExists", "jaeger-span").Return(notFountError())
			indexClient.On("AliasExists", "jaeger-span").Return(notFountError())
			indexClient.On("CreateIndex", "jaeger-span").Return(test.createIndexReturnErr)
			err := createIndexIfNotExist(indexClient, "jaeger-span")
			assert.Equal(t, test.createIndexExpectedError, err)
		})
	}
}

func notFountError() error {
	return client.ResponseError{
		StatusCode: http.StatusNotFound,
		Err:        errors.New("not found"),
	}
}

func testCreateIndexIfNotExist(t *testing.T, existence string) {
	existsTests := []struct {
		name                string
		existsReturnError   error
		existsExpectedError error
	}{
		{
			name: "success",
		},
		{
			name:                "generic error from existence method",
			existsReturnError:   errors.New("may be an http error from index exists"),
			existsExpectedError: errors.New("may be an http error from index exists"),
		},
	}

	for _, test := range existsTests {
		t.Run(test.name, func(t *testing.T) {
			indexClient := &mocks.IndexAPI{}
			if existence == "index" {
				indexClient.On("IndexExists", "jaeger-span").Return(test.existsReturnError)
				indexClient.On("AliasExists", "jaeger-span").Return(nil)
			} else if existence == "alias" {
				indexClient.On("IndexExists", "jaeger-span").Return(nil)
				indexClient.On("AliasExists", "jaeger-span").Return(test.existsReturnError)
			}
			indexClient.On("IndexExists", "jaeger-span").Return(test.existsReturnError)
			indexClient.On("CreateIndex", "jaeger-span").Return(nil)
			err := createIndexIfNotExist(indexClient, "jaeger-span")
			assert.Equal(t, test.existsExpectedError, err)
		})
	}
}

func TestRolloverAction(t *testing.T) {
	tests := []struct {
		name                  string
		setupCallExpectations func(indexClient *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, ilmClient *mocks.IndexManagementLifecycleAPI)
		config                Config
		expectedErr           error
	}{
		{
			name: "Unsupported version",
			setupCallExpectations: func(_ *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, _ *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(5), nil)
			},
			config: Config{
				Config: app.Config{
					Archive: true,
					UseILM:  true,
				},
			},
			expectedErr: errors.New("ILM is supported only for ES version 7+"),
		},
		{
			name: "error getting version",
			setupCallExpectations: func(_ *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, _ *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(0), errors.New("version error"))
			},
			expectedErr: errors.New("version error"),
			config: Config{
				Config: app.Config{
					Archive: true,
					UseILM:  true,
				},
			},
		},
		{
			name: "ilm doesnt exist",
			setupCallExpectations: func(_ *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, ilmClient *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(7), nil)
				ilmClient.On("Exists", "myilmpolicy").Return(false, nil)
			},
			expectedErr: errors.New("ILM policy myilmpolicy doesn't exist in Elasticsearch. Please create it and re-run init"),
			config: Config{
				Config: app.Config{
					Archive:       true,
					UseILM:        true,
					ILMPolicyName: "myilmpolicy",
				},
			},
		},
		{
			name: "fail get ilm policy",
			setupCallExpectations: func(_ *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, ilmClient *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(7), nil)
				ilmClient.On("Exists", "myilmpolicy").Return(false, errors.New("error getting ilm policy"))
			},
			expectedErr: errors.New("error getting ilm policy"),
			config: Config{
				Config: app.Config{
					Archive:       true,
					UseILM:        true,
					ILMPolicyName: "myilmpolicy",
				},
			},
		},
		{
			name: "fail to create template",
			setupCallExpectations: func(indexClient *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, _ *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(7), nil)
				indexClient.On("CreateTemplate", mock.Anything, "jaeger-span").Return(errors.New("error creating template"))
			},
			expectedErr: errors.New("error creating template"),
			config: Config{
				Config: app.Config{
					Archive: true,
					UseILM:  false,
				},
			},
		},
		{
			name: "fail to get jaeger indices",
			setupCallExpectations: func(indexClient *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, _ *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(7), nil)
				indexClient.On("IndexExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("AliasExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("CreateTemplate", mock.Anything, "jaeger-span").Return(nil)
				indexClient.On("CreateIndex", "jaeger-span-archive-000001").Return(nil)
				indexClient.On("GetJaegerIndices", "").Return([]client.Index{}, errors.New("error getting jaeger indices"))
			},
			expectedErr: errors.New("error getting jaeger indices"),
			config: Config{
				Config: app.Config{
					Archive: true,
					UseILM:  false,
				},
			},
		},
		{
			name: "fail to create alias",
			setupCallExpectations: func(indexClient *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, _ *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(7), nil)
				indexClient.On("IndexExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("AliasExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("CreateTemplate", mock.Anything, "jaeger-span").Return(nil)
				indexClient.On("CreateIndex", "jaeger-span-archive-000001").Return(nil)
				indexClient.On("GetJaegerIndices", "").Return([]client.Index{}, nil)
				indexClient.On("CreateAlias", []client.Alias{
					{Index: "jaeger-span-archive-000001", Name: "jaeger-span-archive-read", IsWriteIndex: false},
					{Index: "jaeger-span-archive-000001", Name: "jaeger-span-archive-write", IsWriteIndex: false},
				}).Return(errors.New("error creating aliases"))
			},
			expectedErr: errors.New("error creating aliases"),
			config: Config{
				Config: app.Config{
					Archive: true,
					UseILM:  false,
				},
			},
		},
		{
			name: "create rollover index",
			setupCallExpectations: func(indexClient *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, _ *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(7), nil)
				indexClient.On("IndexExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("AliasExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("CreateTemplate", mock.Anything, "jaeger-span").Return(nil)
				indexClient.On("CreateIndex", "jaeger-span-archive-000001").Return(nil)
				indexClient.On("GetJaegerIndices", "").Return([]client.Index{}, nil)
				indexClient.On("CreateAlias", []client.Alias{
					{Index: "jaeger-span-archive-000001", Name: "jaeger-span-archive-read", IsWriteIndex: false},
					{Index: "jaeger-span-archive-000001", Name: "jaeger-span-archive-write", IsWriteIndex: false},
				}).Return(nil)
			},
			expectedErr: nil,
			config: Config{
				Config: app.Config{
					Archive: true,
					UseILM:  false,
				},
			},
		},
		{
			name: "create rollover index with ilm",
			setupCallExpectations: func(indexClient *mocks.IndexAPI, clusterClient *mocks.ClusterAPI, ilmClient *mocks.IndexManagementLifecycleAPI) {
				clusterClient.On("Version").Return(uint(7), nil)
				indexClient.On("IndexExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("AliasExists", "jaeger-span-archive-000001").Return(notFountError())
				indexClient.On("CreateTemplate", mock.Anything, "jaeger-span").Return(nil)
				indexClient.On("CreateIndex", "jaeger-span-archive-000001").Return(nil)
				indexClient.On("GetJaegerIndices", "").Return([]client.Index{}, nil)
				ilmClient.On("Exists", "jaeger-ilm").Return(true, nil)
				indexClient.On("CreateAlias", []client.Alias{
					{Index: "jaeger-span-archive-000001", Name: "jaeger-span-archive-read", IsWriteIndex: false},
					{Index: "jaeger-span-archive-000001", Name: "jaeger-span-archive-write", IsWriteIndex: true},
				}).Return(nil)
			},
			expectedErr: nil,
			config: Config{
				Config: app.Config{
					Archive:       true,
					UseILM:        true,
					ILMPolicyName: "jaeger-ilm",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexClient := &mocks.IndexAPI{}
			clusterClient := &mocks.ClusterAPI{}
			ilmClient := &mocks.IndexManagementLifecycleAPI{}
			initAction := Action{
				Config:        test.config,
				IndicesClient: indexClient,
				ClusterClient: clusterClient,
				ILMClient:     ilmClient,
			}

			test.setupCallExpectations(indexClient, clusterClient, ilmClient)

			err := initAction.Do()
			if test.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, test.expectedErr, err)
			}

			indexClient.AssertExpectations(t)
			clusterClient.AssertExpectations(t)
			ilmClient.AssertExpectations(t)
		})
	}
}
