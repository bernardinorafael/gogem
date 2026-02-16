package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bernardinorafael/gogem/pkg/fault"
)

const maxRequestBodyBytes = 1_048_576 // 1MB

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	var f *fault.Fault
	if errors.As(err, &f) {
		w.WriteHeader(f.GetHTTPCode())
		_ = encoder.Encode(f)
		return
	}

	switch {
	case errors.Is(err, ErrUnknownRequestBodyKey):
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(fault.NewBadRequest("request body contains unknown field"))
		return
	case errors.Is(err, ErrEmptyRequestBody):
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(fault.NewBadRequest("request body cannot be empty"))
		return
	case errors.Is(err, ErrInvalidJSONField):
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(fault.NewBadRequest("request body contains incorrect JSON field type"))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(fault.NewInternalServerError("an unexpected error occurred"))
		return
	}
}

// WriteSuccess writes a JSON success response with the specified HTTP status code.
// The JSON response will contain a single key "message" with the value "success".
func WriteSuccess(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

// WriteJSON writes a JSON response with the specified HTTP status code and the provided data.
// The data to be encoded as JSON should be passed as the 'dst' parameter.
func WriteJSON(w http.ResponseWriter, code int, dst any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(dst)
}

// ReadQueryInt reads a query string parameter from the URL values and parses it into an integer.
// If the parameter is missing or cannot be parsed, the provided default value 'defval' is returned.
//
// Example:
//
//	page := httputil.ReadQueryInt(r.URL.Query(), "page", 1)
//
// If the query string parameter "page" is not present, the default value 1 is returned.
// If the query string parameter "page" is present but cannot be parsed as an integer,
// the default value 1 is returned.
func ReadQueryInt(qs url.Values, key string) int {
	val := qs.Get(key)
	if val == "" {
		return 0
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}

	return i
}

// ReadQueryBool reads a query string param from the URL and returns it as a boolean
// If the parameter is missing or cannot be parsed, the provided default value 'defval' is returned.
//
// Example:
//
//	includeArchived := httputil.ReadQueryBool(r.URL.Query(), "include_archived", false)
//
// If the query string parameter "include_archived" is not present, the default value false is returned.
// If the query string parameter "include_archived" is present but cannot be parsed as a boolean.
func ReadQueryBool(qs url.Values, key string) bool {
	val := qs.Get(key)
	if val == "" {
		return false
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return b
}

// ReadQueryString reads a query string parameter from the URL values and returns it as a string.
// If the parameter is missing, the provided default value 'defval' is returned.
//
// Example:
//
//	sort := httputil.ReadQueryString(r.URL.Query(), "sort", "asc")
//
// If the query string parameter "sort" is not present, the default value "asc" is returned.
func ReadQueryString(qs url.Values, key string) string {
	val := qs.Get(key)
	if val == "" {
		return ""
	}
	return val
}

// ReadQueryArray reads a query string parameter from the URL values and returns it as a string slice.
// The parameter value is split by commas, and each element is trimmed of whitespace.
// If the parameter is missing or empty, an empty slice is returned.
//
// Example:
//
//	tags := httputil.ReadQueryArray(r.URL.Query(), "tags")
//
// If the query string parameter "tags" contains "user,admin,guest", the returned slice will be ["user", "admin", "guest"].
// If the query string parameter "tags" is not present or empty, an empty slice [] is returned.
func ReadQueryArray(qs url.Values, key string) []string {
	val := qs.Get(key)
	if val == "" {
		return []string{}
	}

	parts := strings.Split(val, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}

	return result
}

// ReadRequestBody reads and parses the JSON body of an HTTP request into the provided destination struct.
// It limits the size of the request body to 1MB and returns detailed error messages for various parsing issues.
//
// Example:
//
//	var body struct {
//		ID string `json:"id"`
//		Name string `json:"name"`
//	}
//
//	err := httputil.ReadRequestBody(w, r, &body)
//	if err != nil {
//		// handle error here
//	}
func ReadRequestBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxRequestBodyBytes))

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			// JSON syntax error in the request body
			// Offset is the exact byte where the error occurred
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			// JSON value and struct type do not match
			if unmarshalTypeError.Field != "" {
				return ErrInvalidJSONField
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			// io.EOF (End of File) indicates that there are no more bytes left to read
			return ErrEmptyRequestBody
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return ErrUnknownRequestBodyKey
		case errors.As(err, &invalidUnmarshalError):
			// Received a non-nil pointer into Decode()
			panic(err)
		default:
			return err
		}
	}

	// Calling decode again to check if there's more data after the JSON object
	// This will return an io.EOF error, indicating that the client sent more data
	err = d.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode: body must only contain a single JSON value: %w", err)
	}

	return nil
}

// GetClientIP extracts the client's real IP address from HTTP request headers.
// It checks multiple headers in order of preference:
// 1. X-Forwarded-For (load balancers/proxies)
// 2. X-Real-IP (reverse proxies)
// 3. RemoteAddr (direct connection)
//
// For X-Forwarded-For, it returns the first IP in the comma-separated list,
// which represents the original client IP.
//
// WARNING: Proxy headers (X-Forwarded-For, X-Real-IP) can be easily spoofed by clients.
// This function should only be used behind trusted infrastructure (e.g., load balancers, reverse proxies)
// that sanitize these headers. Do not use this function to determine client identity in untrusted environments.
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (most common for load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		// We want the first one (original client)
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	// Check X-Real-IP header (used by some reverse proxies)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to RemoteAddr (direct connection)
	// RemoteAddr format is "IP:port", so we need to strip the port
	ip := r.RemoteAddr

	// Handle IPv6 addresses wrapped in brackets before stripping port
	if strings.HasPrefix(ip, "[") {
		if closingBracket := strings.Index(ip, "]"); closingBracket != -1 {
			ip = ip[1:closingBracket]
		}
	} else {
		// For IPv4, strip the port
		if colon := strings.LastIndex(ip, ":"); colon != -1 {
			ip = ip[:colon]
		}
	}

	return ip
}
