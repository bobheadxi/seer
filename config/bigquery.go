package config

// BigQuery denotes Google BigQuery options
type BigQuery struct {
	DatasetID        string `env:"BIGQUERY_DATASET_ID"`
	MatchesTableID   string `env:"BIGQUERY_TABLE_ID_MATCHES"`
	TimelinesTableID string `env:"BIGQUERY_TABLE_ID_TIMELINES"`
}
