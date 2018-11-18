package routes

import (
	"errors"
	"hacking-portal/models"
)

type mockStudentStorage struct {
	data map[string]models.Student
}

func (s mockStudentStorage) FindAll() ([]models.Student, error) {
	students := make([]models.Student, 0, len(s.data))
	for _, student := range s.data {
		students = append(students, student)
	}
	return students, nil
}

func (s mockStudentStorage) FindByID(id string) (models.Student, error) {
	var st models.Student
	if student, ok := s.data[id]; ok {
		return student, nil
	}
	return st, errors.New("")
}

func (s mockStudentStorage) FindByName(name string) (models.Student, error) {
	var st models.Student
	for _, student := range s.data {
		if student.Name == name {
			return student, nil
		}
	}

	return st, errors.New("")
}

func (s mockStudentStorage) FindByGroup(groupID int) ([]models.Student, error) {
	var students []models.Student
	for _, student := range s.data {
		if student.GroupID == groupID {
			students = append(students, student)
		}
	}

	return students, errors.New("")
}

func (s mockStudentStorage) FindGroups() (map[int]int, error) {
	groupIDs := map[int]int{}
	for _, student := range s.data {
		groupID := student.GroupID
		if groupID != 0 {
			if _, isset := groupIDs[groupID]; !isset {
				groupIDs[groupID] = 1
			} else {
				groupIDs[groupID]++
			}
		}
	}

	return groupIDs, nil
}

func (s *mockStudentStorage) Upsert(student models.Student) error {
	if s.data == nil {
		s.data = map[string]models.Student{}
	}

	student.GroupID = 0
	s.data[student.ID] = student
	return nil
}
