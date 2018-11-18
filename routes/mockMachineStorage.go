package routes

import (
	"errors"
	"hacking-portal/models"
)

type mockMachineStorage struct {
	data map[string]models.Machine
}

func injectMachines(s *mockMachineStorage) {
	s.data = map[string]models.Machine{}
	s.data["1111"] = models.Machine{
		Name:    "test1",
		UUID:    "1111",
		GroupID: 1,
		Address: "1.1.1.1",
	}
	s.data["2222"] = models.Machine{
		Name:    "test2",
		UUID:    "2222",
		GroupID: 2,
		Address: "2.2.2.2",
	}
	s.data["3333"] = models.Machine{
		Name:    "test2",
		UUID:    "3333",
		GroupID: 3,
		Address: "3.3.3.3",
	}
}

func (s mockMachineStorage) FindAll() ([]models.Machine, error) {
	machines := make([]models.Machine, 0, len(s.data))
	for _, machine := range s.data {
		machines = append(machines, machine)
	}
	return machines, nil
}

func (s *mockMachineStorage) FindByID(uuid string) (models.Machine, error) {
	injectMachines(s)
	var m models.Machine
	if machine, ok := s.data[uuid]; ok {
		return machine, nil
	}
	return m, errors.New("")
}

func (s mockMachineStorage) FindByGroup(groupID int) ([]models.Machine, error) {
	var machines []models.Machine
	for _, machine := range s.data {
		if machine.GroupID == groupID {
			machines = append(machines, machine)
		}
	}

	return machines, errors.New("")
}

func (s *mockMachineStorage) Upsert(machine models.Machine) error {
	if s.data == nil {
		s.data = map[string]models.Machine{}
	}

	if machine.GroupID == -1 {
		return errors.New("")
	}

	machine.GroupID = 0
	s.data[machine.UUID] = machine

	return nil
}
