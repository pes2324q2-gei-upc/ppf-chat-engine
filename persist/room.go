package persist

import (
	"gorm.io/gorm"
)

type RoomRecord struct {
	Id     string `gorm:"primarykey"`
	Name   string
	Driver string
	Users  []*UserRecord `gorm:"many2many:users_rooms;"` // many to many relationship
}

func (r RoomRecord) Pk() string              { return r.Id }
func (r RoomRecord) GetName() string         { return r.Name }
func (r RoomRecord) GetDriver() string       { return r.Driver }
func (r RoomRecord) GetUsers() []*UserRecord { return r.Users }

type RoomRepository struct {
	*gorm.DB
}

func (repo RoomRepository) Exists(id int) (bool, error) {
	r := repo.First(&RoomRecord{})
	return r.RowsAffected >= 1, r.Error
}

func (repo RoomRepository) Add(room RoomRecord) error {
	return repo.Create(&room).Error
}

func (repo RoomRepository) Remove(id string) error {
	return repo.Delete(&RoomRecord{Id: id}).Error
}

func (repo RoomRepository) Get(id string) (*RoomRecord, error) {
	var result *RoomRecord = &RoomRecord{Id: id}
	stm := repo.First(result)
	return result, stm.Error
}

func (repo RoomRepository) GetAll() ([]*RoomRecord, error) {
	var results []*RoomRecord = make([]*RoomRecord, 0)
	stm := repo.Find(results)
	return results, stm.Error
}

func (repo RoomRepository) AddUser(room RoomRecord, user UserRecord) error {
	return repo.Model(&room).Association("Users").Append(&user)
}

func (repo RoomRepository) RemoveUser(room RoomRecord, user UserRecord) error {
	return repo.Model(&room).Association("Users").Delete(&user)
}
