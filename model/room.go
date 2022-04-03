package model

import (
	"sync"
)

var (
	_roomPool *roomPool
	once      sync.Once
)

type roomPool struct {
	roomMap *sync.Map
	mutex   *sync.RWMutex
	size    int
}

func (r *roomPool) AddRoom(room *Room) {
	r.mutex.Lock()
	_, ok := r.roomMap.Load(room.RoomId)
	if !ok {
		r.size++
	}
	r.mutex.Unlock()
	r.roomMap.Store(room.RoomId, room)

}

func (r roomPool) RoomSize() int {
	r.mutex.RLock()
	defer r.mutex.Unlock()
	return r.size
}

func (r *roomPool) DelRoom(roomId string) {
	r.mutex.Lock()
	_, ok := r.roomMap.Load(roomId)
	if !ok {
		r.mutex.Unlock()
		return
	}
	r.size--
	r.mutex.Unlock()
	r.roomMap.Delete(roomId)
}

func GetRoomPool() *roomPool {
	once.Do(func() {
		_roomPool = &roomPool{
			roomMap: &sync.Map{},
			mutex:   &sync.RWMutex{},
		}
	})
	return _roomPool
}

type Room struct {
	RoomId string
	// id of user who owns the room
	OwnerId  int
	RoomName string
}
