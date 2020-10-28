package pvtdata

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	implicitCollectionPrefix = "_implicit_org_"
)

// Technical pvtdata errors
// TODO: make errors constant
var (
	ErrWrongTransientFieldName  = errors.New("Field is not present in transient data map")
	ErrEmptyTransientFieldValue = errors.New("Transient field has empty value")
	ErrPrivateDataNotFound      = errors.New("Private data was not found")
)

// GetTransientDataValue is a function that obtains the value of a specific field of the transient data
func GetTransientDataValue(ctx contractapi.TransactionContextInterface, fieldName string) (value []byte, err error) {
	TransientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return
	}

	value, ok := TransientMap[fieldName]
	if !ok {
		err = ErrWrongTransientFieldName
		return
	}

	if len(value) == 0 {
		err = ErrEmptyTransientFieldValue
		return
	}

	return
}

// GetTransientDataValueUnmarshaled is a function that obtains the value of a specific field of the transient data that should be a JSON string and unmarshals it
func GetTransientDataValueUnmarshaled(ctx contractapi.TransactionContextInterface, fieldName string, v interface{}) (err error) {
	valueBytes, err := GetTransientDataValue(ctx, fieldName)
	if err != nil {
		return
	}

	err = json.Unmarshal(valueBytes, v)
	return
}

// PutImplicitPrivateData is a function to store private data in the implicit collection of the specified organization
func PutImplicitPrivateData(ctx contractapi.TransactionContextInterface, collectionMSP string, key string, v interface{}) (err error) {
	value, err := json.Marshal(v)
	if err != nil {
		return
	}

	err = PutImplicitPrivateDataBytes(ctx, collectionMSP, key, value)
	return
}

// PutImplicitPrivateDataBytes is a function to store private data in the implicit collection of the specified organization
func PutImplicitPrivateDataBytes(ctx contractapi.TransactionContextInterface, collectionMSP string, key string, value []byte) (err error) {
	collection := implicitCollectionPrefix + collectionMSP
	err = ctx.GetStub().PutPrivateData(collection, key, value)
	return
}

// GetImplicitPrivateData is a function to retrieve data stored in the implicit private data collection of the specified organization
func GetImplicitPrivateData(ctx contractapi.TransactionContextInterface, collectionMSP string, key string, v interface{}) (err error) {
	collection := implicitCollectionPrefix + collectionMSP
	bytes, err := ctx.GetStub().GetPrivateData(collection, key)
	if err != nil {
		return
	}
	if bytes == nil {
		err = ErrPrivateDataNotFound
		return
	}

	err = json.Unmarshal(bytes, v)
	return
}
