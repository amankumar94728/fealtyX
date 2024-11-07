package storage

import (
    "errors"
    "sync"

    "github.com/amankumar94728/fealtyx-student-api/internal/models"
)

type Storage struct {
    students map[int]models.Student
    mutex    sync.RWMutex
    nextID   int
}

func NewStorage() *Storage {
    return &Storage{
        students: make(map[int]models.Student),
        nextID:   1,
    }
}

func (s *Storage) Create(student models.Student) (models.Student, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    student.ID = s.nextID
    s.students[student.ID] = student
    s.nextID++
    return student, nil
}

func (s *Storage) GetAll() []models.Student {
    s.mutex.RLock()
    defer s.mutex.RUnlock()

    students := make([]models.Student, 0, len(s.students))
    for _, student := range s.students {
        students = append(students, student)
    }
    return students
}

func (s *Storage) GetByID(id int) (models.Student, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()

    student, ok := s.students[id]
    if !ok {
        return models.Student{}, errors.New("student not found")
    }
    return student, nil
}

func (s *Storage) Update(id int, student models.Student) (models.Student, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if _, ok := s.students[id]; !ok {
        return models.Student{}, errors.New("student not found")
    }

    student.ID = id
    s.students[id] = student
    return student, nil
}

func (s *Storage) Delete(id int) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if _, ok := s.students[id]; !ok {
        return errors.New("student not found")
    }

    delete(s.students, id)
    return nil
}