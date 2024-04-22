package persist

import (
	"errors"

	"gorm.io/gorm"
)

type UserRecord struct {
	Id   string `gorm:"primarykey"`
	Name string
}

func (r UserRecord) Pk() string      { return r.Id }
func (r UserRecord) GetName() string { return r.Name }

type UserRepository struct {
	*gorm.DB
}

func (repo *UserRepository) Exists(id string) (bool, error) {
	r := repo.First(&UserRecord{})
	return r.RowsAffected >= 1, r.Error
}

func (repo UserRepository) Add(user UserRecord) error {
	return repo.Create(&user).Error
}

func (repo UserRepository) Remove(id string) error {
	return repo.Delete(&UserRecord{}, id).Error
}

func (repo UserRepository) Get(id string) (*UserRecord, error) {
	var result *UserRecord = &UserRecord{Id: id}
	stm := repo.First(result).Preload("Rooms")
	return result, stm.Error
}

func (repo UserRepository) GetAll() ([]*UserRecord, error) {
	var users []*UserRecord = make([]*UserRecord, 0)
	stm := repo.Find(users)
	return users, stm.Error
}

func (repo UserRepository) AddRoom(user UserRecord, room RoomRecord) error {
	stm := repo.Model(&user).Association("Rooms").Append(&room)
	return errors.New(stm.Error())
}

func (repo UserRepository) RemoveRoom(user UserRecord, room RoomRecord) error {
	stm := repo.Model(&user).Association("Rooms").Delete(&room)
	return errors.New(stm.Error())
}
