package car

import (
	"fmt"
	"github.com/ozonmp/omp-bot/internal/model/insurance"
)

type CarService interface {
	Describe(carID uint64) (*insurance.Car, error)
	List(cursor uint64, limit uint64) ([]insurance.Car, error)
	Create(insurance.Car) (uint64, error)
	Update(carID uint64, car insurance.Car) error
	Remove(carID uint64) (bool, error)
}

func (d DummyCarService) Describe(carID uint64) (*insurance.Car, error) {
	car, ok := d.storage[carID]
	if !ok {
		return nil, fmt.Errorf("no car with id %d", carID)
	}
	return &car, nil
}

func (d DummyCarService) List(cursor uint64, limit uint64) ([]insurance.Car, error) {
	return []insurance.Car{}, nil
}

func (d DummyCarService) Create(car insurance.Car) (uint64, error) {
	return 0, nil
}

func (d DummyCarService) Update(carID uint64, car insurance.Car) error {
	return nil
}

func (d DummyCarService) Remove(carID uint64) (bool, error) {
	_, ok := d.storage[carID]
	if !ok {
		return false, nil
	}
	delete(d.storage, carID)
	return true, nil
}

type DummyCarService struct {
	storage map[uint64]insurance.Car
}

func NewDummyCarService() *DummyCarService {
	return &DummyCarService{
		storage: map[uint64]insurance.Car{
			0: insurance.Car{ID: 0,	Title: "Toyota"},
			1: insurance.Car{ID: 1,	Title: "Nissan"},
			2: insurance.Car{ID: 2,	Title: "Infinity"},
		},
	}
}
