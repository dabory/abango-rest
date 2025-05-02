// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

type (
	MemoryMap map[string]interface{}

	MapStore struct {
		store MemoryMap
	}
)

func (c *MapStore) Get(key string) interface{} {
	return c.store[key]
}

func (c *MapStore) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(MemoryMap)
	}
	c.store[key] = val
}
