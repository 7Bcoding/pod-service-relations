package utils

import "sync"

type Set interface {
	Add(interface{})
	Remove(interface{})
	Contains(interface{}) bool
	Elements() []interface{}
	Empty() bool
	Clear()
	Intersect(Set) Set
	Union(Set) Set
	Subtract(Set) Set
	IsSubset(Set) bool
	Size() int
}

type safeSet struct {
	data  map[interface{}]struct{}
	mutex *sync.RWMutex
}

func NewSet() Set {
	return newSafeSet()
}

func newSafeSet() *safeSet {
	return &safeSet{
		data:  make(map[interface{}]struct{}),
		mutex: new(sync.RWMutex),
	}
}

func (s *safeSet) Add(elem interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[elem] = struct{}{}
}

func (s *safeSet) Remove(elem interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, elem)
}

func (s *safeSet) Contains(elem interface{}) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.data[elem]
	return ok
}

func (s *safeSet) Empty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.data) == 0
}

func (s *safeSet) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = make(map[interface{}]struct{})
}

func (s *safeSet) Elements() []interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	res := make([]interface{}, 0, len(s.data))
	for elem := range s.data {
		res = append(res, elem)
	}
	return res
}

func (s *safeSet) Intersect(s1 Set) Set {
	res := newSafeSet()
	elems := s1.Elements()
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, elem := range elems {
		if _, ok := s.data[elem]; ok {
			res.data[elem] = struct{}{}
		}
	}
	return res
}

func (s *safeSet) Union(s1 Set) Set {
	res := newSafeSet()
	for _, elem := range s1.Elements() {
		res.data[elem] = struct{}{}
	}
	for _, elem := range s.Elements() {
		res.data[elem] = struct{}{}
	}
	return res
}

func (s *safeSet) Subtract(s1 Set) Set {
	res := newSafeSet()
	intersection := s1.Intersect(s)
	for _, elem := range s.Elements() {
		if !intersection.Contains(elem) {
			res.data[elem] = struct{}{}
		}
	}
	return res
}

func (s *safeSet) IsSubset(s1 Set) bool {
	return s.Subtract(s1).Empty()
}

func (s *safeSet) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.data)
}
