package storage

import (
	"bytes"
	"strings"
)

/**
 * DataRecord class
 */
type DataRecord map[string]string

/**
 * DataRecord constructor
 */
func CreateDataRecord() *DataRecord {
	return &DataRecord{}
}

/**
 * DataRecord.Set(string, string)
 */
func (record *DataRecord) Set(key string, val string) {
	(*record)[key] = val
}

/**
 * DataRecord.Get(string) string
 */
func (record *DataRecord) Get(key string) string {
	return (*record)[key]
}

/**
 * DataRecord.ToMap() map[string]string
 */
func (record *DataRecord) ToMap() map[string]string {
	return map[string]string(*record)
}

/**
 * DataRecord.FromMap(map[string]string) *DataRecord
 */
func (record *DataRecord) FromMap(data map[string]string) *DataRecord {
	record = &DataRecord{}

	for key, val := range data {
		record.Set(key, val)
	}

	return record
}

/**
 * DataRecord.Exists(string) bool
 */
func (record *DataRecord) Exists(key string) bool {
	if _, ok := (*record)[key]; ok {
		return true
	}

	return false
}

/**
 * DataRecord.ToString() string
 */
func (record *DataRecord) ToString() string {
	var buffer bytes.Buffer

	for key, val := range *record {
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString(val)
		buffer.WriteString(",")
	}

	s := buffer.String()
	return s[0:len(s)-1]
}

/**
 * DataRecord.FromString(string) *DataRecord
 */
func (record *DataRecord) FromString(line string) *DataRecord {
	opts := strings.Split(line, ",")
	record = &DataRecord{}

	for _, val := range opts {
		opt := strings.Split(val, "=")
		record.Set(opt[0], opt[1])
	}

	return record
}

/**
 * DataRecord.Equals(*DataRecord) bool
 */
func (record *DataRecord) Equals(needle *DataRecord) bool {
	for key, _ := range needle.ToMap() {
		if !record.Exists(key) || record.Get(key) != needle.Get(key) {
			return false
		}
	}

	return true
}
