package state

import (
	"testing"

	"github.com/braduf/fabric-chaincode-go-helpers/mocking"
)

type oneFieldMockState struct {
	Field string `json:"field"`
}

func TestPutState(t *testing.T) {
	mockStub := mocking.NewMockChaincodeStub("TestPutState", nil)
	mockTransactionContext := mocking.NewMockTransactionContext(mockStub, nil)
	key := "key"
	value := oneFieldMockState{"value"}

	mockStub.MockTransactionStart("1")
	err := PutState(mockTransactionContext, key, value)
	mockStub.MockTransactionEnd("1")

	assertError(t, err, nil)
}

func TestGetState(t *testing.T) {
	t.Run("get non-existing state", func(t *testing.T) {
		mockStub := mocking.NewMockChaincodeStub("TestGetState", nil)
		mockTransactionContext := mocking.NewMockTransactionContext(mockStub, nil)

		mockStub.MockTransactionStart("1")
		var got oneFieldMockState
		err := GetState(mockTransactionContext, "nonExistingKey", &got)
		mockStub.MockTransactionEnd("1")

		assertError(t, err, ErrStateNotFound)
	})

	t.Run("get put state", func(t *testing.T) {
		mockStub := mocking.NewMockChaincodeStub("TestGetState", nil)
		mockTransactionContext := mocking.NewMockTransactionContext(mockStub, nil)
		key := "worldStateKey"
		value := oneFieldMockState{"value"}

		mockStub.MockTransactionStart("1")
		_ = PutState(mockTransactionContext, key, value)
		mockStub.MockTransactionEnd("1")

		mockStub.MockTransactionStart("2")
		var got oneFieldMockState
		err := GetState(mockTransactionContext, key, &got)
		mockStub.MockTransactionEnd("2")

		assertError(t, err, nil)
		if got != value {
			t.Errorf("got %q want %q", got, value)
		}
	})
}

func assertError(t *testing.T, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
	if got == nil {
		if want == nil {
			return
		}
		t.Fatalf("expected to get error %q", want)
	}
}
