package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func postJSON(t *testing.T, url string, body any, token string) *http.Response {
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request to %s failed: %v", url, err)
	}
	return resp
}

func TestHappyPathReceptionFlow(t *testing.T) {
	resp := postJSON(t, baseURL+"/dummyLogin", map[string]string{
		"role": "employee",
	}, "")
	defer resp.Body.Close()

	var authResp struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&authResp)
	token := authResp.Token

	if token == "" {
		t.Fatal("failed to get token")
	}

	modTokenResp := postJSON(t, baseURL+"/dummyLogin", map[string]string{
		"role": "moderator",
	}, "")
	var mod struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(modTokenResp.Body).Decode(&mod)
	modToken := mod.Token

	pvzResp := postJSON(t, baseURL+"/pvz", map[string]string{
		"city": "Москва",
	}, modToken)
	var pvz struct {
		ID string `json:"id"`
	}
	_ = json.NewDecoder(pvzResp.Body).Decode(&pvz)

	if pvz.ID == "" {
		t.Fatal("failed to create PVZ")
	}

	pvzID := pvz.ID

	recResp := postJSON(t, baseURL+"/receptions", map[string]string{
		"pvzId": pvzID,
	}, token)
	if recResp.StatusCode != http.StatusCreated {
		t.Fatal("failed to create reception")
	}
	var rec struct {
		ID string `json:"id"`
	}
	_ = json.NewDecoder(recResp.Body).Decode(&rec)

	for i := 0; i < 50; i++ {
		p := map[string]any{
			"pvzId": pvzID,
			"type":  "электроника",
		}
		resp := postJSON(t, baseURL+"/products", p, token)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("failed to add product %d", i+1)
		}
	}

	closeResp := postJSON(t, baseURL+"/pvz/"+pvzID+"/close_last_reception", nil, token)
	if closeResp.StatusCode != http.StatusOK {
		t.Fatal("failed to close reception")
	}
}
