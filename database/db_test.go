package database

import (
	"testing"

	"github.com/minhmannh2001/sms/test"
)

func TestInitDB(t *testing.T) {
	err := test.MockedDB(test.CREATE)
	if err != nil {
		t.Errorf("TestInitDB Failed")
		return
	}
	defer func() {
		if err := test.MockedDB(test.DROP); err != nil {
			t.Errorf("TestInitDB Failed")
		}
	}()
}
