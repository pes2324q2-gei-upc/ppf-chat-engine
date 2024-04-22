package persist

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type MessageKey struct {
	room   string
	sender string
}

func (k MessageKey) Pk() (string, string) { return k.room, k.sender }
func (k MessageKey) Room() string         { return k.room }
func (k MessageKey) Sender() string       { return k.sender }

type MessageRecord struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time

	Room    RoomRecord // Belongs To relation
	Sender  UserRecord // Belongs To relation
	Content string
}

func (r MessageRecord) Pk() MessageKey {
	return MessageKey{
		room:   r.Room.Pk(),
		sender: r.Sender.Pk(),
	}
}

type MessageRepository struct {
	*gorm.DB
}

func (repo MessageRepository) Exists(id MessageKey) (bool, error) {
	r := repo.First(&MessageRecord{}).Where("room_id = ? AND sender_id = ?", id.Room, id.Sender)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, r.Error
	}
	return true, nil
}

func (repo MessageRepository) Add(msg MessageRecord) error {
	return repo.Create(&msg).Error
}

func (repo MessageRepository) Remove(pk MessageKey) error {
	r := repo.Delete(&MessageRecord{}).Where("room_id = ? AND sender_id = ?", pk.Room, pk.Sender)
	return r.Error
}

func (repo MessageRepository) Get(pk MessageKey) (*MessageRecord, error) {
	result := &MessageRecord{}
	stm := repo.First(result).Where("room_id = ? AND sender_id = ?", pk.Room, pk.Sender)
	return result, stm.Error
}

func (repo MessageRepository) GetAll() ([]*MessageRecord, error) {
	var results []*MessageRecord = make([]*MessageRecord, 0)
	stm := repo.Preload("Room").Preload("Sender").Find(results)
	return results, stm.Error
}

func (repo MessageRepository) GetByRoom(room RoomRecord) ([]*MessageRecord, error) {
	var results []*MessageRecord
	err := repo.Model(&MessageRecord{Room: room}).Preload("Sender").Find(&results).Error
	return results, err
}

func (repo MessageRepository) GetBySender(sender UserRecord) ([]*MessageRecord, error) {
	var results []*MessageRecord
	err := repo.Model(&MessageRecord{Sender: sender}).Preload("Room").Find(&results).Error
	return results, err
}
