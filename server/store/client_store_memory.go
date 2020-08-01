package store

import "errors"

var (
	// _memoryClientStore 内存模式 Client Store
	_memoryClientStore *memoryClientStore
)

// memoryClientStore 内存类型 Client Store
type memoryClientStore struct {
	clientMap map[string]*ClientInfo
}

// GetClient 获得客户端信息
// clientId 客户端ID
func (s *memoryClientStore) GetClient(clientId string) (*ClientInfo, error) {
	client, ok := s.clientMap[clientId]
	if ok {
		return client, nil
	} else {
		return nil, errors.New("客户端不存在！")
	}
}

// AddClient 追加客户端信息
func (s *memoryClientStore) AddClient(clientInfo *ClientInfo) error {
	if len(clientInfo.ClientId) > 0 {
		s.clientMap[clientInfo.ClientId] = clientInfo
	}
	return nil
}

// 移除客户端信息
func (s *memoryClientStore) RemoveClient(clientId string) error {
	delete(s.clientMap, clientId)
	return nil
}

// 获取内存类型 Client Store
func NewMemoryClientStoreInstance() ClientStore {
	if _memoryClientStore == nil {
		_memoryClientStore = &memoryClientStore{
			clientMap: make(map[string]*ClientInfo),
		}
	}
	return _memoryClientStore
}
