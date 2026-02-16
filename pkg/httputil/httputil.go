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

	var f *fault.Fault
	if errors.As(err, &f) {
		w.WriteHeader(f.GetHTTPCode())
		_ = json.NewEncoder(w).Encode(f)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(fault.NewInternalServerError("an unexpected error occurred"))
}

// WriteSuccess writes a JSON success response with the specified HTTP status code.
func WriteSuccess(w http.ResponseWriter, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

// WriteJSON writes a JSON response with the specified HTTP status code and the provided data.
func WriteJSON(w http.ResponseWriter, code int, dst any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(dst)
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
func ReadQueryInt(qs url.Values, key string, defaultValue int) int {
	val := qs.Get(key)
	if val == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
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
func ReadQueryBool(qs url.Values, key string, defaultValue bool) bool {
	val := qs.Get(key)
	if val == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
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
func ReadQueryString(qs url.Values, key string, defaultValue string) string {
	val := qs.Get(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// ReadQueryIntOptional reads a query string parameter and returns a pointer to the parsed integer.
// Returns nil if the parameter is missing or cannot be parsed.
//
// Example:
//
//	if page := httputil.ReadQueryIntOptional(r.URL.Query(), "page"); page != nil {
//	    // parameter was explicitly provided
//	}
func ReadQueryIntOptional(qs url.Values, key string) *int {
	val := qs.Get(key)
	if val == "" {
		return nil
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return nil
	}

	return &i
}

// ReadQueryBoolOptional reads a query string parameter and returns a pointer to the parsed boolean.
// Returns nil if the parameter is missing or cannot be parsed.
//
// Example:
//
//	if archived := httputil.ReadQueryBoolOptional(r.URL.Query(), "archived"); archived != nil {
//	    // parameter was explicitly provided
//	}
func ReadQueryBoolOptional(qs url.Values, key string) *bool {
	val := qs.Get(key)
	if val == "" {
		return nil
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return nil
	}

	return &b
}

// ReadQueryStringOptional reads a query string parameter and returns a pointer to the string value.
// Returns nil if the parameter is missing.
//
// Example:
//
//	if sort := httputil.ReadQueryStringOptional(r.URL.Query(), "sort"); sort != nil {
//	    // parameter was explicitly provided
//	}
func ReadQueryStringOptional(qs url.Values, key string) *string {
	val := qs.Get(key)
	if val == "" {
		return nil
	}
	return &val
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
//	err := httputil.ReadRequestBody(r, &body)
//	if err != nil {
//		// handle error here
//	}
func ReadRequestBody(r *http.Request, dst any) error {
	reader := io.LimitReader(r.Body, maxRequestBodyBytes+1)

	d := json.NewDecoder(reader)
	d.DisallowUnknownFields()

	err := d.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fault.NewBadRequest(
				fmt.Sprintf("body contains badly-formed JSON (at character %d)", syntaxError.Offset),
			)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fault.NewBadRequest("body contains incorrect JSON field type")
			}
			return fault.NewBadRequest(
				fmt.Sprintf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset),
			)
		case errors.Is(err, io.EOF):
			return fault.NewBadRequest("request body cannot be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fault.NewBadRequest("request body contains unknown field")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return fault.NewBadRequest(err.Error())
		}
	}

	err = d.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return fault.NewBadRequest("body must only contain a single JSON value")
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
