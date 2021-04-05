package game

import (
	"fmt"
	"sync"
)

type licenses struct {
	maxActive int
	mu        sync.Mutex
	nextID    int
	licenses  map[int]*License
	isActive  map[int]bool
}

func newLicenses(maxActive int) *licenses {
	return &licenses{
		maxActive: maxActive,
		licenses:  make(map[int]*License, maxActive),
		isActive:  make(map[int]bool, maxActive),
	}
}

func (ls *licenses) init(active []License) {
	for i := range active {
		l := active[i]
		if ls.nextID < l.ID+1 {
			ls.nextID = l.ID + 1
		}
		ls.licenses[l.ID] = &l
		ls.isActive[l.ID] = true
	}
}

func (ls *licenses) active() []License {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	active := make([]License, 0, len(ls.licenses))
	for id := range ls.licenses {
		if ls.isActive[id] {
			active = append(active, *ls.licenses[id])
		}
	}
	return active
}

func (ls *licenses) beginIssue(digAllowed int) (l License, _ error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if digAllowed < 1 || digAllowed > maxDigAllowed {
		panic(fmt.Sprintf("digAllowed=%d must be between 1 and %d", digAllowed, maxDigAllowed))
	}
	if len(ls.licenses) >= ls.maxActive {
		return l, ErrActiveLicenseLimit
	}

	l = License{
		ID:         ls.nextID,
		DigAllowed: digAllowed,
	}
	ls.licenses[ls.nextID] = &l
	ls.nextID++
	return l, nil
}

func (ls *licenses) mustBegunIssue(id int) {
	if _, ok := ls.isActive[id]; ok {
		panic("never here")
	} else if ls.licenses[id] == nil {
		panic("never here")
	}
}

func (ls *licenses) commitIssue(id int) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.mustBegunIssue(id)
	ls.isActive[id] = true
}

func (ls *licenses) rollbackIssue(id int) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.mustBegunIssue(id)
	delete(ls.licenses, id)
}

func (ls *licenses) use(id int) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if _, ok := ls.isActive[id]; !ok {
		return ErrNoSuchLicense
	}
	ls.licenses[id].DigUsed++
	if ls.licenses[id].DigUsed >= ls.licenses[id].DigAllowed {
		delete(ls.licenses, id)
		delete(ls.isActive, id)
	}
	return nil
}
