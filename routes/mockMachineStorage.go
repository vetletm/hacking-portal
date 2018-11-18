package routes

import (
	"errors"
	"hacking-portal/models"
)

type mockMachineStorage struct {
	data map[string]models.Machine
}

func (s mockMachineStorage) FindAll() ([]models.Machine, error) {
	machines := make([]models.Machine, 0, len(s.data))
	for _, machine := range s.data {
		machines = append(machines, machine)
	}
	return machines, nil
}

func (s *mockMachineStorage) FindByID(uuid string) (models.Machine, error) {
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

	return machines, nil
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
