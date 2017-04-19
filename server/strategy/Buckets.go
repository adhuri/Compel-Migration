package strategy

import "github.com/Sirupsen/logrus"

type Bucket struct {
	AgentIP          string
	ContainerDetails []*Container
	PseudoFreeCPU    float32 // Psuedo - This changes every time you make a local decision of moving containers across
	PseudoFreeMemory float32 // Psuedo - This changes every time you make a local decision of moving containers across
}

type Container struct {
	ContainerID string
	CPUValue    float32
	MemValue    float32
}

func NewBucket(AgentIP string) *Bucket {
	return &Bucket{
		AgentIP:          AgentIP,
		PseudoFreeCPU:    -1.0,
		PseudoFreeMemory: -1.0,
		ContainerDetails: make([]*Container, 0),
	}
}

func NewContainer(ContainerID string, CPUValue float32, MemValue float32) *Container {
	return &Container{
		ContainerID: ContainerID,
		CPUValue:    CPUValue,
		MemValue:    MemValue,
	}
}

func (m *Bucket) GetFreeCPU() float32 {
	totalCPU := float32(100.0)
	freeCPU := float32(0.0)
	for _, container := range m.ContainerDetails {
		freeCPU += container.CPUValue
	}
	return totalCPU - freeCPU
}

func (m *Bucket) GetFreeMemory() float32 {
	totalMem := float32(0.0)
	freeMem := float32(0.0)
	for _, container := range m.ContainerDetails {
		freeMem += container.MemValue
	}
	return totalMem - freeMem
}

func (m *Bucket) PrintBucket(log *logrus.Logger) {
	log.Infoln("Agent IP ", m.AgentIP)
	log.Infoln("Total Containers ", len(m.ContainerDetails))
	if m.PseudoFreeCPU == -1.0 {
		log.Infoln("Free CPU ", m.GetFreeCPU())
	} else {
		log.Infoln("Free CPU ", m.PseudoFreeCPU)
	}

	if m.PseudoFreeMemory == -1.0 {
		log.Infoln("Free Memory ", m.GetFreeMemory())
	} else {
		log.Infoln("Free Memory ", m.PseudoFreeMemory)

	}
	for containerIndex, container := range m.ContainerDetails {
		log.Infoln("Container ", containerIndex, " : ContainerID ", container.ContainerID)
	}
}

func PrintAllBuckets(Buckets []*Bucket, log *logrus.Logger) {
	log.Info("===============Print Buckets=====================")
	for bucketIndex, bucket := range Buckets {
		log.Info("-----------------------------------------------")
		log.Info("Bucket ", bucketIndex)
		bucket.PrintBucket(log)
		log.Info("-----------------------------------------------")

	}
	log.Info("==========================================")

}
