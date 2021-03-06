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

package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/CanonicalLtd/serial-vault/service/log"
)

// JSONHeader is the JSON HTTP header
const JSONHeader = "application/json; charset=UTF-8"

// StandardResponse is the JSON response from an API method, indicating success or failure.
type StandardResponse struct {
	Success      bool   `json:"success"`
	ErrorCode    string `json:"error_code"`
	ErrorSubcode string `json:"error_subcode"`
	ErrorMessage string `json:"message"`
}

// FormatStandardResponse returns a JSON response from an API method, indicating success or failure.
func FormatStandardResponse(success bool, errorCode, errorSubcode, message string, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := StandardResponse{Success: success, ErrorCode: errorCode, ErrorSubcode: errorSubcode, ErrorMessage: message}

	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error forming the boolean response (%v)\n. %v", response, err)
		return err
	}
	return nil
}

// ParseStandardResponse parses the response body and returns a standard response object
func ParseStandardResponse(w *httptest.ResponseRecorder) (StandardResponse, error) {
	// Check the JSON response
	result := StandardResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	return result, err
}
