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
	if carID >= uint64(len(d.storage)) {
		return nil, fmt.Errorf("no car with id %d", carID)
	}
	return &d.storage[carID], nil
}

func (d DummyCarService) List(cursor uint64, limit uint64) ([]insurance.Car, error) {
	if cursor >= uint64(len(d.storage)) {
		return nil, fmt.Errorf("no car with id %d", cursor)
	}
	high := uint64(len(d.storage))
	if cursor + limit < high {
		high = cursor + limit
	}
	return d.storage[cursor:high], nil
}

func (d *DummyCarService) Create(car insurance.Car) (uint64, error) {
	d.storage = append(d.storage, car)
	return uint64(len(d.storage) - 1), nil
}

func (d *DummyCarService) Update(carID uint64, car insurance.Car) error {
	if carID >= uint64(len(d.storage)) {
		return fmt.Errorf("no car with id %d", carID)
	}
	d.storage[carID] = car
	return nil
}

func (d *DummyCarService) Remove(carID uint64) (bool, error) {
	if carID >= uint64(len(d.storage)) {
		return false, fmt.Errorf("no car with id %d", carID)
	}
	copy(d.storage[carID:], d.storage[carID+1:])
	d.storage = d.storage[:len(d.storage)-1]
	return true, nil
}

type DummyCarService struct {
	storage []insurance.Car
}

func NewDummyCarService() *DummyCarService {
	return &DummyCarService{
		storage: []insurance.Car{
			{Title: "Toyota"},
			{Title: "Nissan"},
			{Title: "Infinity"},
			{Title: "Mazda"},
			{Title: "Honda"},
			{Title: "Lexus"},
			{Title: "Acura"},
			{Title: "Suzuki"},
			{Title: "Isuzu"},
			{Title: "Mitsubishi"},
			{Title: "Subaru"},
		},
	}
}
