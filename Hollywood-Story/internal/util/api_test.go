package util

import "testing"

func TestAPIInternalPrompt(t *testing.T) {

	resp, err := APIInternalPrompt("6583af4c32314355db5f7b06", nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
