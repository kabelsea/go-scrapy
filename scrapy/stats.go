package scrapy

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
)

type SpiderStats interface {
	String() string
	Clear()
	SetValue(k string, v interface{})
	GetValue(k string) (interface{}, bool)
	IncValue(k string) error
}

// Spider stats collection
type Stats struct {
	values map[string]interface{}
	mutex  *sync.Mutex
}

func NewStats() *Stats {
	return &Stats{
		values: make(map[string]interface{}),
		mutex:  &sync.Mutex{},
	}
}

// Return collected statistics on json format.
func (s *Stats) String() string {
	res, err := json.Marshal(s.values)
	if err != nil {
		logger.Error(err)
	}
	return string(res)
}

func (s *Stats) SetValue(k string, v interface{}) {
	s.values[k] = v
}

func (s *Stats) GetValue(k string) (interface{}, bool) {
	v, ok := s.values[k]
	return v, ok
}

func (s *Stats) IncValue(k string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch v := s.values[k].(type) {
	case int:
		s.values[k] = v + 1
	default:
		return errors.New(fmt.Sprintf("Wrong value type in key %s, its not integer", k))
	}
	return nil
}

// Clear spider stats
func (s *Stats) Clear() {
	s.values = make(map[string]interface{})
}
