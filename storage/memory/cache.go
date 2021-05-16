package memory

func NewLogStorage() *LogCache {
	return &LogCache{map[string][]string{}}
}

type LogCache struct {
	logs map[string][]string
}

func (c *LogCache) TxIDAlreadyExists(id string) (bool, error) {
	_, ok := c.logs[id]
	return ok, nil
}

func (c *LogCache) AppendLog(id string, logData string) error {
	if logs, ok := c.logs[id]; ok {
		c.logs[id] = append(logs, logData)
	} else {
		c.logs[id] = []string{logData}
	}

	return nil
}
