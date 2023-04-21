package dynamodb

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joeig/eee-safe/pkg/threema"
	"testing"
)

type mockDynamoDBClient struct {
	PutItemWithContextOutput    *dynamodb.PutItemOutput
	PutItemWithContextErr       error
	GetItemWithContextOutput    *dynamodb.GetItemOutput
	GetItemWithContextErr       error
	DeleteItemWithContextOutput *dynamodb.DeleteItemOutput
	DeleteItemWithContextErr    error
}

func (m *mockDynamoDBClient) PutItemWithContext(_ aws.Context, _ *dynamodb.PutItemInput, _ ...request.Option) (*dynamodb.PutItemOutput, error) {
	return m.PutItemWithContextOutput, m.PutItemWithContextErr
}

func (m *mockDynamoDBClient) GetItemWithContext(_ aws.Context, _ *dynamodb.GetItemInput, _ ...request.Option) (*dynamodb.GetItemOutput, error) {
	return m.GetItemWithContextOutput, m.GetItemWithContextErr
}

func (m *mockDynamoDBClient) DeleteItemWithContext(_ aws.Context, _ *dynamodb.DeleteItemInput, _ ...request.Option) (*dynamodb.DeleteItemOutput, error) {
	return m.DeleteItemWithContextOutput, m.DeleteItemWithContextErr
}

func TestDynamoDB_PutBackup(t *testing.T) {
	dynamoDB := &DynamoDB{
		Table: "table",
		svc:   &mockDynamoDBClient{},
	}

	if err := dynamoDB.PutBackup(context.Background(), &threema.BackupInput{}); err != nil {
		t.Error("unexpected error")
	}
}

func TestDynamoDB_PutBackup_error(t *testing.T) {
	dynamoDB := &DynamoDB{
		Table: "table",
		svc:   &mockDynamoDBClient{PutItemWithContextErr: errors.New("mock")},
	}

	if err := dynamoDB.PutBackup(context.Background(), &threema.BackupInput{}); err == nil {
		t.Error("no error")
	}
}

func TestDynamoDB_GetBackup(t *testing.T) {
	marshalledItem, _ := dynamodbattribute.MarshalMap(&item{
		BackupID:        "foo",
		EncryptedBackup: []byte{'f'},
		CreationTime:    0,
		ExpirationTime:  0,
	})
	dynamoDB := &DynamoDB{
		Table: "table",
		svc:   &mockDynamoDBClient{GetItemWithContextOutput: &dynamodb.GetItemOutput{Item: marshalledItem}},
	}

	output, err := dynamoDB.GetBackup(context.Background(), threema.BackupID{})

	if output == nil {
		t.Error("no output")
	}

	if err != nil {
		t.Error("unexpected error")
	}
}

func TestDynamoDB_GetBackup_error(t *testing.T) {
	dynamoDB := &DynamoDB{
		Table: "table",
		svc:   &mockDynamoDBClient{GetItemWithContextErr: errors.New("mock")},
	}

	output, err := dynamoDB.GetBackup(context.Background(), threema.BackupID{})

	if output == nil {
		t.Error("no output")
	}

	if err == nil {
		t.Error("no error")
	}
}

func TestDynamoDB_DeleteBackup(t *testing.T) {
	marshalledItem, _ := dynamodbattribute.MarshalMap(&item{
		BackupID:        "foo",
		EncryptedBackup: []byte{'f'},
		CreationTime:    0,
		ExpirationTime:  0,
	})
	dynamoDB := &DynamoDB{
		Table: "",
		svc:   &mockDynamoDBClient{DeleteItemWithContextOutput: &dynamodb.DeleteItemOutput{Attributes: marshalledItem}},
	}

	if err := dynamoDB.DeleteBackup(context.Background(), threema.BackupID{}); err != nil {
		t.Error("unexpected error")
	}
}

func TestDynamoDB_DeleteBackup_error(t *testing.T) {
	dynamoDB := &DynamoDB{
		Table: "",
		svc:   &mockDynamoDBClient{DeleteItemWithContextErr: errors.New("mock")},
	}

	if err := dynamoDB.DeleteBackup(context.Background(), threema.BackupID{}); err == nil {
		t.Error("no error")
	}
}
