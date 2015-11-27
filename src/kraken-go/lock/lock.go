package lock

import (
	"os"
	"time"
	"../errors"
)

/**
 * Lock interface
 */
type Lock interface {
	Lock()		errors.Error
	Unlock()	errors.Error
}

/**
 * FileLock class
 */
type FileLock struct {
	filePath	string
}

/**
 * FileLock constructor
 */
func CreateFileLock(filePath string) *FileLock {
	lock := &FileLock{}

	lock.filePath = "../../data/" + filePath + ".lock"

	return lock
}

/**
 * FileLock.Lock() errors.Error
 */
func (lock *FileLock) Lock() errors.Error {
	var fp *os.File
	var err error

	for {
		if _, err := os.Stat(lock.filePath); os.IsNotExist(err) {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	if fp, err = os.Create(lock.filePath); err != nil {
		return errors.New(9, err.Error())
	}
	fp.Close()

	return nil
}

/**
 * FileLock.Unlock() errors.Error
 */
func (lock *FileLock) Unlock() errors.Error {
	if err := os.Remove(lock.filePath); err != nil {
		return errors.New(10, err.Error())
	}

	return nil
}
