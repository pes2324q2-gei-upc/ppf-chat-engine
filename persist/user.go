package persist

import (
	"errors"
)

type UserRecord struct {
	Id   string
	Name string
}

func (record *UserRecord) Pk() any {
	return record.Id
}

type UserMemoryRepository struct {
	records map[string]UserRecord
}

func (repository *UserMemoryRepository) Add(record UserRecord) error {
	repository.records[record.Id] = record
	return nil
}

func (repository *UserMemoryRepository) Remove(pk string) error {
	delete(repository.records, pk)
	return nil
}

func (repository *UserMemoryRepository) Get(pk string) (*UserRecord, error) {
	record, ok := repository.records[pk]
	if !ok {
		return nil, errors.New("record not found")
	}
	return &record, nil
}

func (repository *UserMemoryRepository) GetAll() ([]UserRecord, error) {
	var records []UserRecord
	for _, record := range repository.records {
		records = append(records, record)
	}
	return records, nil
}

func (repository *UserMemoryRepository) GetBy(field string, value string) (*UserRecord, error) {
	// Reflect the field name to the struct field.
	if field == "Id" {
		for _, record := range repository.records {
			if record.Id == value {
				return &record, nil
			}
		}
	}
	return nil, errors.New("record not found")
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		records: make(map[string]UserRecord),
	}
}
