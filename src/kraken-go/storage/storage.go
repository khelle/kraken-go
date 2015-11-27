package storage

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"../lock"
	"../errors"
)

/**
 * Storage interface
 */
type Storage interface {
	Open()								(bool, errors.Error)
	Close()								(bool, errors.Error)
	Add(*DataRecord)					(bool, errors.Error)
	AddSeveral([]*DataRecord)			(bool, errors.Error)
	Remove(*DataRecord) 				(bool, errors.Error)
	RemoveSeveral([]*DataRecord) 		(bool, errors.Error)
	RemoveAll()							(bool, errors.Error)
	Get(*DataRecord)					([]*DataRecord, errors.Error)
	GetSeveral([]*DataRecord)			([]*DataRecord, errors.Error)
	GetAll() 							([]*DataRecord, errors.Error)
	Exclude(*DataRecord)				([]*DataRecord, errors.Error)
	ExcludeSeveral([]*DataRecord)		([]*DataRecord, errors.Error)
	Erase() 							(bool, errors.Error)
}

/**
 * FileStorage class
 */
type FileStorage struct {
	filePath	string
	fileLock	*lock.FileLock
}

/**
 * FileStorage constructor
 */
func NewFileStorage(filePath string) (*FileStorage, errors.Error) {
	storage := &FileStorage{}

	storage.filePath = "../../data/" + filePath + ".data"
	storage.fileLock = lock.CreateFileLock(filePath)

	return storage, nil
}

/**
 * FileStorage.Open() (bool, errors.Error)
 */
func (fs *FileStorage) Open() (bool, errors.Error) {
	err := fs.fileLock.Lock()
	status := true

	if err != nil {
		status = false
	}

	return status, err
}

/**
 * FileStorage.Close() (bool, errors.Error)
 */
func (fs *FileStorage) Close() (bool, errors.Error) {
	err := fs.fileLock.Unlock()
	status := true

	if err != nil {
		status = false
	}

	return status, err
}

/**
 * FileStorage.Add(*DataRecord) (bool, errors.Error)
 */
func (fs *FileStorage) Add(record *DataRecord) (bool, errors.Error) {
	if _, err := fs.AddSeveral([]*DataRecord{record}); err != nil {
		return false, err
	}

	return true, nil
}

/**
 * FileStorage.AddSeveral([]*DataRecord) (bool, errors.Error)
 */
func (fs *FileStorage) AddSeveral(records []*DataRecord) (bool, errors.Error) {
	var file *os.File
	var err errors.Error
	if file, err = fs.GetStore(); err != nil {
		return false, err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	for _, record := range records {
		fmt.Fprintln(w, record.ToString())
	}

	w.Flush()

	return true, nil
}

/**
 * FileStorage.Remove(*DataRecord) (bool, errors.Error)
 */
func (fs *FileStorage) Remove(record *DataRecord) (bool, errors.Error) {
	return fs.RemoveSeveral([]*DataRecord{record})
}

/**
 * FileStorage.RemoveSeveral([]*DataRecord) (bool, errors.Error)
 */
func (fs *FileStorage) RemoveSeveral(records []*DataRecord) (bool, errors.Error) {
	records, err := fs.ExcludeSeveral(records)
	if err != nil {
		return false, err
	}

	if _, err = fs.Erase(); err != nil {
		return false, err
	}
	
	if _, err = fs.AddSeveral(records); err != nil {
		return false, err
	}

	return true, nil
}

/**
 * FileStorage.RemoveAll() (bool, errors.Error)
 */
func (fs *FileStorage) RemoveAll() (bool, errors.Error) {
	return fs.Erase()
}

/**
 * FileStorage.Get(*DataRecord) ([]*DataRecord, errors.Error)
 */
func (fs *FileStorage) Get(needle *DataRecord) ([]*DataRecord, errors.Error) {
	return fs.GetSeveral([]*DataRecord{needle})
}

/**
 * FileStorage.GetSeveral([]*DataRecord) ([]*DataRecord, errors.Error)
 */
func (fs *FileStorage) GetSeveral(needles []*DataRecord) ([]*DataRecord, errors.Error) {
	var records []*DataRecord
	var err errors.Error

	if records, err = fs.GetAll(); err != nil {
		return nil, err
	}

	find := []*DataRecord{}
	for _, record := range records {
		for _, needle := range needles {
			if record.Equals(needle) {
				find = append(find, record)
				break
			}
		}
	}

	return find, nil
}

/**
 * FileStorage.GetAll() ([]*DataRecord, errors.Error)
 */
func (fs *FileStorage) GetAll() ([]*DataRecord, errors.Error) {
	var file *os.File
	var err errors.Error

	if file, err = fs.GetStore(); err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	records := []*DataRecord{}

	for scanner.Scan() {
		opts := strings.Split(scanner.Text(), ",")
		record := &DataRecord{}

		for _, val := range opts {
			opt := strings.Split(val, "=")
			record.Set(opt[0], opt[1])
		}

		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New(4, err.Error())
	}

	return records, nil
}

/**
 * FileStorage.Exclude(*DataRecord) []*DataRecord
 */
func (fs *FileStorage) Exclude(needle *DataRecord) ([]*DataRecord, errors.Error) {
	return fs.ExcludeSeveral([]*DataRecord{needle})
}

/**
 * FileStorage.ExcludeSeveral([]*DataRecord) ([]*DataRecord, errors.Error)
 */
func (fs *FileStorage) ExcludeSeveral(needles []*DataRecord) ([]*DataRecord, errors.Error) {
	var records []*DataRecord
	var err errors.Error
	var valid bool

	if records, err = fs.GetAll(); err != nil {
		return nil, err
	}

	find := []*DataRecord{}
	for _, record := range records {
		valid = true

		for _, needle := range needles {
			if record.Equals(needle) {
				valid = false
				break
			}
		}

		if valid {
			find = append(find, record)
		}
	}

	return find, nil
}

/**
 * FileStorage.Erase() (bool, errors.Error)
 */
func (fs *FileStorage) Erase() (bool, errors.Error) {
	if err := os.Remove(fs.filePath); err != nil {
		return false, errors.New(13, err.Error())
	}

	return true, nil
}

/**
 * FileStorage.GetStore(filePath string) (*os.File, errors.Error)
 */
func (fs *FileStorage) GetStore() (*os.File, errors.Error) {
	var fp *os.File
	var err error

	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		fp, err = os.Create(fs.filePath)
	} else {
		fp, err = os.OpenFile(fs.filePath, os.O_APPEND|os.O_RDWR, 0666)
	}

	if err != nil {
		return nil, errors.New(5, err.Error())
	}

	return fp, nil
}
