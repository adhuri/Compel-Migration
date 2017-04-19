package strategy

type Bucket struct {
	AgentIP          string
	ContainerDetails []*Container
}

type Container struct {
	ContainerID string
	CPUValue    float64
	MemValue    float64
}

func NewBucket(AgentIP string) *Bucket {
	return &Bucket{
		AgentIP:          AgentIP,
		ContainerDetails: make([]*Container, 0),
	}
}

func NewContainer(ContainerID string, CPUValue float64, MemValue float64) *Container {
	return &Container{
		ContainerID: ContainerID,
		CPUValue:    CPUValue,
		MemValue:    MemValue,
	}
}

func (m *Bucket) GetFreeCPU() float64 {
	freeCPU := 0.0
	for _, container := range m.ContainerDetails {
		freeCPU += container.CPUValue
	}
	return freeCPU
}

func (m *Bucket) GetFreeMemory() float64 {
	freeMem := 0.0
	for _, container := range m.ContainerDetails {
		freeMem += container.MemValue
	}
	return freeMem
}
