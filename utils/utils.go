package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

func Time(strct interface{}, new bool) error {
	t := reflect.TypeOf(strct)
	v := reflect.ValueOf(strct).Elem()

	// Validate if strct is a pointer and struct
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("strct must be a pointer to a struct")
	}

	// Validate if the datecreated and dateupdated fields are in strct and are of type time.Time
	dateCreatedField, has_created := t.Elem().FieldByName("DateCreated")
	dateUpdatedField, has_updated := t.Elem().FieldByName("DateUpdated")
	if !has_created || !has_updated {
		return errors.New("strct must have DateCreated and DateUpdated fields")
	}

	if dateCreatedField.Type.Kind() != reflect.TypeOf(time.Time{}).Kind() || dateUpdatedField.Type.Kind() != reflect.TypeOf(time.Time{}).Kind() {
		return errors.New("strct DateCreated and DateUpdated fields must be of type time.Time")
	}

	// Set the time for the fields based on new arguement value
	if new {
		// Set the "DateUpdated" field to current UTC time
		v.FieldByName("DateCreated").Set(reflect.ValueOf(time.Now().UTC()))
	}

	v.FieldByName("DateUpdated").Set(reflect.ValueOf(time.Now().UTC()))
	return nil

}

func JsonResponse(w http.ResponseWriter, success bool, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		StatusCode: statusCode,
		Success:    success,
		Message:    message,
		Data:       data,
	})
}

func Unmarshall(r io.Reader, struct_ interface{}) error {
	structType := reflect.TypeOf(struct_)
	if structType.Kind() != reflect.Ptr || structType.Elem().Kind() != reflect.Struct {
		return errors.New("struct_ must be pointer to a struct")
	}

	err := json.NewDecoder(r).Decode(struct_)
	if err != nil {
		return err
	}
	return nil
}

func AddQueryParams(baseUrl string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func ConvertStringToTime(dateStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateStr)
}

// ExtractStringField extracts a string field from a map.
func ExtractStringField(rawData map[string]string, key string) string {
	return rawData[key]
}

// HandleDateField handles the conversion of date fields with default values if not provided.
func HandleDateField(dateStr string, defaultValue time.Time) (time.Time, error) {
	if dateStr == "" {
		return defaultValue, nil
	}
	parsedDate, err := ConvertStringToTime(dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedDate, nil
}
