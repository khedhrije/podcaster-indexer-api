package elasticsearchv7

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/port"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/khedhrije/podcaster-indexer-api/internal/configuration"
	"github.com/khedhrije/podcaster-indexer-api/pkg/retry"
	"github.com/rs/zerolog/log"
)

const (
	latestAlias     = "latest"
	previousAlias   = "previous"
	inProgressAlias = "in-progress"
)

// adapter represents an Elasticsearch adapter wrapper.
type adapter struct {
	client *elasticsearch.Client
}

// NewElasticSearchClient initializes a new Elasticsearch adapter with the given configuration.
func NewElasticSearchClient(config *configuration.AppConfig) (port.Indexer, error) {
	esClient, err := createClient(config)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch adapter init failed: %w", err)
	}
	return &adapter{client: esClient}, nil
}

// CreateIndex creates a new index with the given name and default settings.
func (c *adapter) CreateIndex(ctx context.Context, indexName string, docsType string) error {

	metadata, err := metadataByDocsType(docsType)
	if err != nil {
		return errors.New("unknown docs type")
	}

	parsedMetadata := fmt.Sprintf(metadata, 10, 1, 5, 200)
	body := bytes.NewBufferString(parsedMetadata)
	response, err := c.client.Indices.Create(indexName, c.client.Indices.Create.WithBody(body), c.client.Indices.Create.WithContext(ctx))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("error while creating index")
	}
	return nil
}

// DeleteIndexes deletes the indices with the given names.
func (c *adapter) DeleteIndexes(ctx context.Context, indexNames []string) error {
	response, err := c.client.Indices.Delete(indexNames, c.client.Indices.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("error while deleting indexes")
	}
	return nil
}

// CreateAlias creates an alias for the given index name.
func (c *adapter) CreateAlias(ctx context.Context, indexName, aliasName string) error {
	response, err := c.client.Indices.PutAlias([]string{indexName}, aliasName, c.client.Indices.PutAlias.WithContext(ctx))
	if err != nil {
		return err
	}
	defer closeBodyResponse(response)
	if response.StatusCode != http.StatusOK {
		return errors.New("error while adding alias to index")
	}
	return nil
}

// DeleteAlias deletes the alias for the given index name.
func (c *adapter) DeleteAlias(ctx context.Context, indexName, aliasName string) error {
	response, err := c.client.Indices.DeleteAlias([]string{indexName}, []string{aliasName}, c.client.Indices.DeleteAlias.WithContext(ctx))
	if err != nil {
		return err
	}
	defer closeBodyResponse(response)
	if response.StatusCode != http.StatusOK {
		return errors.New("error while deleting alias from index")
	}
	return nil
}

// IndexByAlias retrieves the index names associated with the given alias.
func (c *adapter) IndexByAlias(ctx context.Context, aliasName string) []string {
	response, err := c.client.Indices.Get([]string{aliasName}, c.client.Indices.Get.WithContext(ctx))
	if err != nil || response.IsError() {
		log.Warn().Err(err).Str("alias", aliasName).Msg("No index found for alias " + aliasName)
		return []string{}
	}
	defer closeBodyResponse(response)

	var resp map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Error().Err(err).Msg("Error decoding IndexByAlias response")
		return nil
	}
	indexNames := make([]string, 0, len(resp))
	for name := range resp {
		indexNames = append(indexNames, name)
	}

	return indexNames
}

// MoveIndex transitions aliases between old, previous, and latest indices.
func (c *adapter) MoveIndex(ctx context.Context, indexationName string) error {

	oldIndexes := c.IndexByAlias(ctx, previousAlias)
	previousIndexes := c.IndexByAlias(ctx, latestAlias)
	latestIndexes := c.IndexByAlias(ctx, inProgressAlias)

	if len(oldIndexes) != 1 || len(previousIndexes) != 1 || len(latestIndexes) != 1 {
		return fmt.Errorf("invalid alias state for %s", indexationName)
	}

	if err := c.switchAlias(ctx, previousIndexes[0], previousAlias, oldIndexes); err != nil {

		return err
	}
	if err := c.switchAlias(ctx, latestIndexes[0], latestAlias, previousIndexes); err != nil {
		return err
	}
	if err := c.DeleteAlias(ctx, latestIndexes[0], inProgressAlias); err != nil {
		return err
	}
	if len(oldIndexes) > 0 {
		if err := c.DeleteIndexes(ctx, oldIndexes); err != nil {
			return err
		}
	}
	return nil
}

// switchAlias switches the alias from an old index to a new index.
func (c *adapter) switchAlias(ctx context.Context, newIndex, alias string, oldIndexes []string) error {
	if err := c.CreateAlias(ctx, newIndex, alias); err != nil {
		return err
	}
	if len(oldIndexes) != 0 {
		if err := c.DeleteAlias(ctx, oldIndexes[0], alias); err != nil {
			return err
		}
	}
	return nil
}

// RecordBulkItems indexes a batch of items using bulk indexing with retry logic.
func (c *adapter) RecordBulkItems(ctx context.Context, indexName string, items []interface{}, backoffRetryCount, backoffTimeSeconds int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(backoffRetryCount+1)*time.Duration(backoffTimeSeconds)*time.Second)
	defer cancel()

	var buf bytes.Buffer
	for _, item := range items {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" } }%s`, indexName, "\n"))
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("an error occurred while encoding: %w", err)
		}
		data = append(data, '\n')

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	response, nbTry, err := retry.ExecuteWithBackoffRetry(func() (interface{}, error) {
		return c.client.Bulk(bytes.NewReader(buf.Bytes()), c.client.Bulk.WithIndex(indexName), c.client.Bulk.WithContext(ctx))
	}, backoffRetryCount, time.Duration(backoffTimeSeconds)*time.Second)
	if err != nil {
		log.Warn().Err(err).Msgf("error while bulk indexing in %s after %d tries", indexName, nbTry)
		return fmt.Errorf("error while bulk indexing in %s", indexName)
	}

	bulkResp := response.(*esapi.Response)
	defer closeBodyResponse(bulkResp)

	if bulkResp.IsError() {
		return handleBulkResponseError(bulkResp)
	}

	return handleBulkResponse(bulkResp)
}

// handleBulkResponseError processes bulk response errors.
func handleBulkResponseError(bulkResp *esapi.Response) error {
	var raw map[string]interface{}
	if err := json.NewDecoder(bulkResp.Body).Decode(&raw); err != nil {
		return fmt.Errorf("failure to parse response body: %s", err)
	}
	errorType := raw["error"].(map[string]interface{})["type"]
	errorReason := raw["error"].(map[string]interface{})["reason"]
	return fmt.Errorf("error: [%d] %s: %s", bulkResp.StatusCode, errorType, errorReason)
}

// handleBulkResponse processes successful bulk responses.
func handleBulkResponse(bulkResp *esapi.Response) error {
	var blk BulkResponse
	if err := json.NewDecoder(bulkResp.Body).Decode(&blk); err != nil {
		return fmt.Errorf("failure to parse response body: %s", err)
	}

	var numErrors int
	for _, item := range blk.Items {
		if item.Index.Status > 201 {
			numErrors++
			log.Error().Int("status", item.Index.Status).
				Str("type", item.Index.Error.Type).
				Str("reason", item.Index.Error.Reason).
				Str("cause_type", item.Index.Error.Cause.Type).
				Str("cause_reason", item.Index.Error.Cause.Reason).
				Msg("Bulk index error")
		}
	}

	if numErrors > 0 {
		return fmt.Errorf("indexed documents with [%d] errors", numErrors)
	}
	return nil
}

// createClient creates an Elasticsearch adapter using the given configuration.
func createClient(config *configuration.AppConfig) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{config.Elasticsearch.URL},
		Username:  config.Elasticsearch.User,
		Password:  config.Elasticsearch.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: 20 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("elastic adapter creation failure: %w", err)
	}
	return client, nil
}

// closeBodyResponse handles closing the body of an esapi.Response.
func closeBodyResponse(response *esapi.Response) {
	if response == nil || response.Body == nil {
		log.Debug().Msg("Response or Response.Body is nil")
		return
	}
	if err := response.Body.Close(); err != nil {
		log.Debug().Err(err).Msg("Error while closing Response.Body")
	}
}

// BulkResponse represents the response structure for bulk operations.
type BulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

// Function to read an io.ReadCloser and return the content as a string
func readToString(rc io.ReadCloser) (string, error) {
	defer rc.Close()
	bytes, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
