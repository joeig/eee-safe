package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joeig/eee-safe/debug"
	"github.com/joeig/eee-safe/storage"
	"github.com/joeig/eee-safe/threema"
	"strconv"
	"strings"
	"time"
)

// PutBackup stores a backup to DynamoDB
func (d *DynamoDB) PutBackup(backupInput *threema.BackupInput) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	input := d.generatePutItemInput(backupInput)
	result, err := svc.PutItem(input)
	debug.Printf("Input: %+v", input)
	debug.Printf("Result: %+v", result)
	debug.Printf("Error: %+v", err)
	if err != nil {
		return &storage.StorageBackendError{APIError: err}
	}
	return nil
}

// GetBackup returns a backup from DynamoDB
func (d *DynamoDB) GetBackup(backupID threema.BackupID) (*threema.BackupOutput, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	input := d.generateGetItemInput(backupID)
	result, err := svc.GetItem(input)
	debug.Printf("Input: %+v", input)
	debug.Printf("Result: %+v", result)
	debug.Printf("Error: %+v", err)
	if err != nil {
		return &threema.BackupOutput{}, &storage.StorageBackendError{APIError: err}
	}
	var resultItem item
	if err := dynamodbattribute.UnmarshalMap(result.Item, &resultItem); err != nil {
		return &threema.BackupOutput{}, &storage.StorageBackendError{APIError: err}
	}
	if len(resultItem.EncryptedBackup) == 0 {
		return &threema.BackupOutput{}, &storage.BackupIDNotFoundError{BackupID: backupID}
	}
	return &threema.BackupOutput{
		BackupID:        backupID,
		EncryptedBackup: threema.EncryptedBackup(resultItem.EncryptedBackup),
		CreationTime:    time.Unix(int64(resultItem.CreationTime), 0),
		ExpirationTime:  time.Unix(int64(resultItem.ExpirationTime), 0),
	}, nil
}

// DeleteBackup deletes a backup from DynamoDB
func (d *DynamoDB) DeleteBackup(backupID threema.BackupID) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	input := d.generateDeleteItemInput(backupID)
	result, err := svc.DeleteItem(input)
	debug.Printf("Input: %+v", input)
	debug.Printf("Result: %+v", result)
	debug.Printf("Error: %+v", err)
	if err != nil {
		return &storage.StorageBackendError{APIError: err}
	}
	var resultItem item
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &resultItem); err != nil {
		return &storage.StorageBackendError{APIError: err}
	}
	if len(resultItem.EncryptedBackup) == 0 {
		return &storage.BackupIDNotFoundError{BackupID: backupID}
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
	var returnValues string
	returnValues = "ALL_OLD"
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
		ReturnValues: &returnValues,
		TableName:    aws.String(d.Table),
	}
}

func (d *DynamoDB) generateGetItemInput(backupID threema.BackupID) *dynamodb.GetItemInput {
	var projectionExpression string
	projectionExpression = "encryptedBackup, creationTime, expirationTime"
	return &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"backupID": {
				S: aws.String(strings.ToLower(backupID.String())),
			},
		},
		ProjectionExpression: &projectionExpression,
		TableName:            aws.String(d.Table),
	}
}

func (d *DynamoDB) generateDeleteItemInput(backupID threema.BackupID) *dynamodb.DeleteItemInput {
	var returnValues string
	returnValues = "ALL_OLD"
	return &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"backupID": {
				S: aws.String(strings.ToLower(backupID.String())),
			},
		},
		ReturnValues: &returnValues,
		TableName:    aws.String(d.Table),
	}
}
