package main

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	EggRacks []EggRack
	Users    []User
	mu       sync.Mutex
}

func NewDatabase() *Database {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"),
		bcrypt.DefaultCost)

	return &Database{
		EggRacks: make([]EggRack, 0),
		Users: []User{
			{Username: "richard", Password: string(hashedPassword)},
		},
	}
}

func (db *Database) CreateEggRack(rack EggRack) (*EggRack, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := rack.Validate(); err != nil {
		return nil, err
	}

	rack.ID = uuid.New()
	rack.DateCreated = time.Now().UTC()

	db.EggRacks = append(db.EggRacks, rack)
	return &rack, nil
}

func (db *Database) CreateBulkEggRacks(username string, jumboCount, bigCount, mediumCount, smallCount int) ([]EggRack, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if jumboCount < 0 || bigCount < 0 || mediumCount < 0 || smallCount < 0 {
		return nil, errors.New("egg counts cannot be negative")
	}

	var created []EggRack

	for i := 0; i < jumboCount; i++ {
		rack := EggRack{
			ID:          uuid.New(),
			EggType:     EggTypeJumbo,
			DateCreated: time.Now().UTC(),
			User:        username,
		}
		db.EggRacks = append(db.EggRacks, rack)
		created = append(created, rack)
	}
	for i := 0; i < bigCount; i++ {
		rack := EggRack{
			ID:          uuid.New(),
			EggType:     EggTypeBig,
			DateCreated: time.Now().UTC(),
			User:        username,
		}
		db.EggRacks = append(db.EggRacks, rack)
		created = append(created, rack)
	}
	for i := 0; i < mediumCount; i++ {
		rack := EggRack{
			ID:          uuid.New(),
			EggType:     EggTypeMedium,
			DateCreated: time.Now().UTC(),
			User:        username,
		}
		db.EggRacks = append(db.EggRacks, rack)
		created = append(created, rack)
	}
	for i := 0; i < smallCount; i++ {
		rack := EggRack{
			ID:          uuid.New(),
			EggType:     EggTypeSmall,
			DateCreated: time.Now().UTC(),
			User:        username,
		}
		db.EggRacks = append(db.EggRacks, rack)
		created = append(created, rack)
	}
	return created, nil
}

func (db *Database) GetAllEggRacks(username string) []EggRack {
	db.mu.Lock()
	defer db.mu.Unlock()

	var result []EggRack
	for _, rack := range db.EggRacks {
		if rack.User == username {
			result = append(result, rack)
		}
	}
	return result
}

func (db *Database) FindUser(username, password string) (*User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, u := range db.Users {
		if u.Username == username {
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err == nil {
				return &u, nil
			}
			return nil, errors.New("invalid password")
		}
	}
	return nil, errors.New("user not found")
}
