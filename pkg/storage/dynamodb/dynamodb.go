package dynamodb

// DynamoDB defines the configuration of the DynamoDB storage backend type
type DynamoDB struct {
	Table string `mapstructure:"table"`
}
