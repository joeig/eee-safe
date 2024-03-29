package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
)

// Client defines a DynamoDB client interface.
type Client interface {
	// PutItemWithContext puts an item to DynamoDB with context.
	PutItemWithContext(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error)
	// GetItemWithContext gets an item from DynamoDB with context.
	GetItemWithContext(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error)
	// DeleteItemWithContext deletes an item from DynamoDB with context.
	DeleteItemWithContext(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error)
}

// DynamoDB defines the configuration of the DynamoDB storage backend type.
type DynamoDB struct {
	Table string `mapstructure:"table"`
	svc   Client
}

// InitializeService initializes a DynamoDB service for a specific session.
func (d *DynamoDB) InitializeService(sess client.ConfigProvider) {
	d.svc = dynamodb.New(sess)
}

// PutBackup stores a backup to DynamoDB.
func (d *DynamoDB) PutBackup(ctx context.Context, backupInput *threema.BackupInput) error {
	input := d.generatePutItemInput(backupInput)
	if _, err := d.svc.PutItemWithContext(ctx, input); err != nil {
		return &storage.ErrUnknown{APIError: err}
	}

	return nil
}

// GetBackup returns a backup from DynamoDB.
func (d *DynamoDB) GetBackup(ctx context.Context, backupID threema.BackupID) (*threema.BackupOutput, error) {
	input := d.generateGetItemInput(backupID)

	result, err := d.svc.GetItemWithContext(ctx, input)
	if err != nil {
		return &threema.BackupOutput{}, &storage.ErrUnknown{APIError: err}
	}

	var resultItem item

	if err := dynamodbattribute.UnmarshalMap(result.Item, &resultItem); err != nil {
		return &threema.BackupOutput{}, &storage.ErrUnknown{APIError: err}
	}

	if len(resultItem.EncryptedBackup) == 0 {
		return &threema.BackupOutput{}, &storage.ErrBackupIDNotFound{BackupID: backupID}
	}

	return &threema.BackupOutput{
		BackupID:        backupID,
		EncryptedBackup: resultItem.EncryptedBackup,
		CreationTime:    time.Unix(int64(resultItem.CreationTime), 0),
		ExpirationTime:  time.Unix(int64(resultItem.ExpirationTime), 0),
	}, nil
}

// DeleteBackup deletes a backup from DynamoDB.
func (d *DynamoDB) DeleteBackup(ctx context.Context, backupID threema.BackupID) error {
	input := d.generateDeleteItemInput(backupID)

	result, err := d.svc.DeleteItemWithContext(ctx, input)
	if err != nil {
		return &storage.ErrUnknown{APIError: err}
	}

	var resultItem item

	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &resultItem); err != nil {
		return &storage.ErrUnknown{APIError: err}
	}

	if len(resultItem.EncryptedBackup) == 0 {
		return &storage.ErrBackupIDNotFound{BackupID: backupID}
	}

	return nil
}

type item struct {
	BackupID        string `json:"backupID"`
	EncryptedBackup []byte `json:"encryptedBackup"`
	CreationTime    int    `json:"creationTime"`
	ExpirationTime  int    `json:"expirationTime"`
}

func (d *DynamoDB) generatePutItemInput(backup *threema.BackupInput) *dynamodb.PutItemInput {
	creationTime := strconv.FormatInt(time.Now().Unix(), 10)
	expirationTime := strconv.FormatInt(time.Now().AddDate(0, 0, int(backup.RetentionDays)).Unix(), 10)

	return &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"backupID": {
				S: aws.String(strings.ToLower(backup.BackupID.String())),
			},
			"encryptedBackup": {
				B: backup.EncryptedBackup,
			},
			"creationTime": {
				N: aws.String(creationTime),
			},
			"expirationTime": {
				N: aws.String(expirationTime),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
		TableName:    aws.String(d.Table),
	}
}

func (d *DynamoDB) generateGetItemInput(backupID threema.BackupID) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"backupID": {
				S: aws.String(strings.ToLower(backupID.String())),
			},
		},
		ProjectionExpression: aws.String("encryptedBackup, creationTime, expirationTime"),
		TableName:            aws.String(d.Table),
	}
}

func (d *DynamoDB) generateDeleteItemInput(backupID threema.BackupID) *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"backupID": {
				S: aws.String(strings.ToLower(backupID.String())),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
		TableName:    aws.String(d.Table),
	}
}
