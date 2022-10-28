package module

type Identifier interface {
	GetID() ID
	GetDescription() string
	GetVersion() Version
}

type ID string

func (id ID) String() string { return string(id) }

type IDs []ID

func (ids IDs) Unique() UniqueIDs {
	unique := make(UniqueIDs, len(ids))
	for _, id := range ids {
		if _, ok := unique[id]; !ok {
			unique[id] = struct{}{}
		}
	}
	return unique
}

type UniqueIDs map[ID]struct{}

func (uids UniqueIDs) Iterate() chan ID {
	iterator := make(chan ID, len(uids))
	for id := range uids {
		iterator <- id
	}
	close(iterator)
	return iterator
}
