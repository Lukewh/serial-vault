// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	settingsFile = "../settings.yaml"
	config := ConfigSettings{}
	err := ReadConfig(&config)
	if err != nil {
		t.Errorf("Error reading config file: %v", err)
	}
}

func TestReadConfigInvalidPath(t *testing.T) {
	settingsFile = "not a good path"
	config := ConfigSettings{}
	err := ReadConfig(&config)
	if err == nil {
		t.Error("Expected an error with an invalid config file.")
	}
}

func TestReadConfigInvalidFile(t *testing.T) {
	settingsFile = "../README.md"
	config := ConfigSettings{}
	err := ReadConfig(&config)
	if err == nil {
		t.Error("Expected an error with an invalid config file.")
	}
}

func TestFormatModelsResponse(t *testing.T) {
	var models []ModelSerialize
	models = append(models, ModelSerialize{ID: 1, BrandID: "Vendor", Name: "Alder 聖誕快樂", Revision: 1})
	models = append(models, ModelSerialize{ID: 2, BrandID: "Vendor", Name: "Ash", Revision: 7})

	w := httptest.NewRecorder()
	err := formatModelsResponse(true, "", "", "", models, w)
	if err != nil {
		t.Errorf("Error forming models response: %v", err)
	}

	var result ModelsResponse
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the models response: %v", err)
	}
	if len(result.Models) != len(models) || !result.Success || result.ErrorMessage != "" {
		t.Errorf("Models response not as expected: %v", result)
	}
	if result.Models[0].Name != models[0].Name {
		t.Errorf("Expected the first model name of '%s', got: %s", models[0].Name, result.Models[0].Name)
	}
}

func TestFormatKeypairsResponse(t *testing.T) {
	var keypairs []Keypair
	keypairs = append(keypairs, Keypair{ID: 1, AuthorityID: "Vendor", KeyID: "12345678abcde", Active: true})
	keypairs = append(keypairs, Keypair{ID: 2, AuthorityID: "Vendor", KeyID: "abcdef123456", Active: true})

	w := httptest.NewRecorder()
	err := formatKeypairsResponse(true, "", "", "", keypairs, w)
	if err != nil {
		t.Errorf("Error forming keypairs response: %v", err)
	}

	var result KeypairsResponse
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the keypairs response: %v", err)
	}
	if len(result.Keypairs) != len(keypairs) || !result.Success || result.ErrorMessage != "" {
		t.Errorf("Keypairs response not as expected: %v", result)
	}
	if result.Keypairs[0].KeyID != keypairs[0].KeyID {
		t.Errorf("Expected the first key ID '%s', got: %s", keypairs[0].KeyID, result.Keypairs[0].KeyID)
	}
}

func TestDecodePublicKeyInvalid(t *testing.T) {
	_, err := decodePublicKey([]byte(""))
	if err == nil {
		t.Error("Expected an error with an invalid public key")
	}

	_, err = decodePublicKey([]byte("ThisIsAnInvalidKey"))
	if err == nil {
		t.Error("Expected an error with an invalid public key")
	}

	_, err = decodePublicKey([]byte("openpgp ThisIsAnInvalidKey"))
	if err == nil {
		t.Error("Expected an error with an invalid public key")
	}

	base64InvalidKey := base64.StdEncoding.EncodeToString([]byte("ThisIsAnInvalidKey"))
	unsupportedKey := fmt.Sprintf("unsupported %s", base64InvalidKey)
	_, err = decodePublicKey([]byte(unsupportedKey))
	if err == nil {
		t.Error("Expected an error with an invalid public key")
	}

	key := fmt.Sprintf("openpgp %s", base64InvalidKey)
	_, err = decodePublicKey([]byte(key))
	if err == nil {
		t.Error("Expected an error with an invalid public key")
	}
}

func TestFormatSigningLogResponse(t *testing.T) {
	var signingLog []SigningLog
	signingLog = append(signingLog, SigningLog{ID: 1, Make: "System", Model: "Router 3400", SerialNumber: "A1", Fingerprint: "a1", Created: time.Now()})
	signingLog = append(signingLog, SigningLog{ID: 2, Make: "System", Model: "Router 3400", SerialNumber: "A2", Fingerprint: "a2", Created: time.Now()})

	w := httptest.NewRecorder()
	err := formatSigningLogResponse(true, "", "", "", signingLog, w)
	if err != nil {
		t.Errorf("Error forming signing log response: %v", err)
	}

	var result SigningLogResponse
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the signing log response: %v", err)
	}
	if len(result.SigningLog) != len(signingLog) || !result.Success || result.ErrorMessage != "" {
		t.Errorf("Signing log response not as expected: %v", result)
	}
	if result.SigningLog[0].Fingerprint != signingLog[0].Fingerprint {
		t.Errorf("Expected the first fingerprint '%s', got: %s", signingLog[0].Fingerprint, result.SigningLog[0].Fingerprint)
	}
}