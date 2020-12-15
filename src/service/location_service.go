package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/models"
	uuid "github.com/satori/go.uuid"
	"log"
)

type LocationService struct {
	db *sql.DB
}

func NewLocationService(db *sql.DB) *LocationService {
	return &LocationService{
		db: db,
	}
}

func (l LocationService) CreateLocation(location dto.LocationRequest) (int64, error) {
	isExist, err := l.isLocationExist(location)
	if err != nil {
		return 0, err
	}
	if isExist {
		return 0, fmt.Errorf("Location already exist for %v %v %v", location.Country, location.CityName, location.Village)
	}
	ctx := context.Background()
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	id := uuid.NewV4()
	sid := id.String()

	sql := fmt.Sprintf("INSERT INTO location (id, country, city_name, address, village) VALUES ('%v', '%v', '%v', '%v', '%v')",
		sid, location.Country, location.CityName, location.Address, location.Village)
	result, err := tx.ExecContext(ctx, sql)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return result.RowsAffected()
}

func (l LocationService) FindLocationById(id string) (*models.Location, error) {
	rows, err := l.db.Query("SELECT id, country, city_name, village, address  FROM location WHERE  id= $1 ", id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		location := &models.Location{}
		err = rows.Scan(&location.ID, &location.Country, &location.CityName, &location.Village, &location.Address)
		if err != nil {
			return nil, err
		} else {
			return location, nil
		}
	}
	return nil, nil
}

func (l LocationService) isLocationExist(location dto.LocationRequest) (bool, error) {
	rows, err := l.db.Query("SELECT country, city_name  FROM location WHERE  country= $1 and city_name= $2",
		location.Country, location.CityName)
	if err != nil {
		log.Print(err)
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}
