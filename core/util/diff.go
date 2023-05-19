/*
	diff.go
	Purpose: diff utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package util

type _stat int

const (
	_stat_remove _stat = iota
	_stat_add
	_stat_remain
)

//-------------------------------------------------
//- String Differ                                 -
//-------------------------------------------------
type _differ struct {
	_map map[string]_stat
}

func NewDiffer() *_differ {
	return &_differ{
		_map: make(map[string]_stat),
	}
}

// SetCurrent sets the current keys
func (d *_differ) SetCurrent(key ...string) {
	for _, k := range key {
		if tmp, ok := d._map[k]; !ok {
			d._map[k] = _stat_remove
		} else if tmp == _stat_add {
			d._map[k] = _stat_remain
		}
	}
}

// SetTarget sets what keys are in the target
func (d *_differ) SetTarget(key ...string) {
	for _, k := range key {
		if tmp, ok := d._map[k]; !ok {
			d._map[k] = _stat_add
		} else if tmp == _stat_remove {
			d._map[k] = _stat_remain
		}
	}
}

// ToBeRemove returns the keys that should be remove
func (d *_differ) ToBeRemove() []string {
	result := []string{}
	for k, v := range d._map {
		if v == _stat_remove {
			result = append(result, k)
		}
	}
	return result
}

// ToBeAdd returns the keys that should be add
func (d *_differ) ToBeAdd() []string {
	result := []string{}
	for k, v := range d._map {
		if v == _stat_add {
			result = append(result, k)
		}
	}
	return result
}

// ToBeAdd returns the keys that should be add
func (d *_differ) Unchanged() []string {
	result := []string{}
	for k, v := range d._map {
		if v == _stat_remain {
			result = append(result, k)
		}
	}
	return result
}

//-------------------------------------------------
//- Int Differ                                    -
//-------------------------------------------------

type _int_differ struct {
	_map map[int]_stat
}

func NewIntDiffer() *_int_differ {
	return &_int_differ{
		_map: make(map[int]_stat),
	}
}

// SetCurrent sets the current keys
func (d *_int_differ) SetCurrent(key ...int) {
	for _, k := range key {
		if tmp, ok := d._map[k]; !ok {
			d._map[k] = _stat_remove
		} else if tmp == _stat_add {
			d._map[k] = _stat_remain
		}
	}
}

// SetTarget sets what keys are in the target
func (d *_int_differ) SetTarget(key ...int) {
	for _, k := range key {
		if tmp, ok := d._map[k]; !ok {
			d._map[k] = _stat_add
		} else if tmp == _stat_remove {
			d._map[k] = _stat_remain
		}
	}
}

// ToBeRemove returns the keys that should be remove
func (d *_int_differ) ToBeRemove() []int {
	result := []int{}
	for k, v := range d._map {
		if v == _stat_remove {
			result = append(result, k)
		}
	}
	return result
}

// ToBeAdd returns the keys that should be add
func (d *_int_differ) ToBeAdd() []int {
	result := []int{}
	for k, v := range d._map {
		if v == _stat_add {
			result = append(result, k)
		}
	}
	return result
}

// ToBeAdd returns the keys that should be add
func (d *_int_differ) Unchanged() []int {
	result := []int{}
	for k, v := range d._map {
		if v == _stat_remain {
			result = append(result, k)
		}
	}
	return result
}
