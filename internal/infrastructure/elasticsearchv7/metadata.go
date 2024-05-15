package elasticsearchv7

import (
	"errors"
	"github.com/khedhrije/podcaster-indexer-api/internal/domain/api"
)

const catsMetadataJSON = `{
  "settings": {
    "index": {
      "refresh_interval": "%vs",
      "number_of_shards": "%v",
      "number_of_replicas": "%v",
      "mapping.nested_fields.limit": %v
    }
  },
  "mappings": {
    "dynamic": false,
      "properties": {
        "ID": {
          "type": "keyword"
        },
        "Name": {
          "type": "text"
        },
        "Description": {
          "type": "text"
        },
        "ParentUUID": {
          "type": "keyword"
        }
      }
  }
}`

const tagsMetadataJSON = `{
  "settings": {
    "index": {
      "refresh_interval": "%vs",
      "number_of_shards": "%v",
      "number_of_replicas": "%v",
      "mapping.nested_fields.limit": %v
    }
  },
  "mappings": {
    "dynamic": false,
      "properties": {
        "ID": {
          "type": "keyword"
        },
        "Name": {
          "type": "text"
        },
        "Description": {
          "type": "text"
        }
      }
  }
}`
const wallsMetadataJSON = `{
  "settings": {
    "index": {
      "refresh_interval": "%vs",
      "number_of_shards": "%v",
      "number_of_replicas": "%v",
      "mapping.nested_fields.limit": %v
    }
  },
  "mappings": {
    "dynamic": false,
      "properties": {
        "ID": {
          "type": "keyword"
        },
        "Name": {
          "type": "text"
        },
        "Description": {
          "type": "text"
        }
      }
  }
}`
const blocksMetadataJSON = `{
  "settings": {
    "index": {
      "refresh_interval": "%vs",
      "number_of_shards": "%v",
      "number_of_replicas": "%v",
      "mapping.nested_fields.limit": %v
    }
  },
  "mappings": {
    "dynamic": false,
      "properties": {
        "ID": {
          "type": "keyword"
        },
        "Name": {
          "type": "text"
        },
        "Description": {
          "type": "text"
        },
		"Kind": {
          "type": "keyword"
        }
      }
  }
}`
const programsMetadataJSON = `{
  "settings": {
    "index": {
      "refresh_interval": "%vs",
      "number_of_shards": "%v",
      "number_of_replicas": "%v",
      "mapping.nested_fields.limit": %v
    }
  },
  "mappings": {
    "dynamic": false,
      "properties": {
        "ID": {
          "type": "keyword"
        },
        "Name": {
          "type": "text"
        },
        "Description": {
          "type": "text"
        }
      }
  }
}`
const episodesMetadataJSON = `{
  "settings": {
    "index": {
      "refresh_interval": "%vs",
      "number_of_shards": "%v",
      "number_of_replicas": "%v",
      "mapping.nested_fields.limit": %v
    }
  },
  "mappings": {
    "dynamic": false,
      "properties": {
        "ID": {
          "type": "keyword"
        },
        "Name": {
          "type": "text"
        },
        "Description": {
          "type": "text"
        },
		"Position": {
          "type": "integer"
        },
		"ProgramID": {
          "type": "keyword"
		}
      }
  }
}`
const mediasMetadataJSON = `{
  "settings": {
    "index": {
      "refresh_interval": "%vs",
      "number_of_shards": "%v",
      "number_of_replicas": "%v",
      "mapping.nested_fields.limit": %v
    }
  },
  "mappings": {
    "dynamic": false,
      "properties": {
        "ID": {
          "type": "keyword"
        },
        "DirectLink": {
          "type": "text"
        },
        "Kind": {
          "type": "keyword"
        },
		"EpisodeID": {
          "type": "keyword"
        }
      }
  }
}`

func metadataByDocsType(docsType string) (string, error) {

	switch docsType {
	case api.Cats:
		return catsMetadataJSON, nil
	case api.Tags:
		return tagsMetadataJSON, nil
	case api.Walls:
		return wallsMetadataJSON, nil
	case api.Blocks:
		return blocksMetadataJSON, nil
	case api.Programs:
		return programsMetadataJSON, nil
	case api.Episodes:
		return episodesMetadataJSON, nil
	case api.Medias:
		return mediasMetadataJSON, nil
	default:
		return "", errors.New("unknown docs type")

	}

}
