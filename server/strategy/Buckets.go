package strategy

import "github.com/Sirupsen/logrus"

type Bucket struct {
	AgentIP          string
	ContainerDetails []*Container
	PseudoFreeCPU    float64 // Psuedo - This changes every time you make a local decision of moving containers across
	PseudoFreeMemory float64 // Psuedo - This changes every time you make a local decision of moving containers across
}

type Container struct {
	ContainerID string
	CPUValue    float64
	MemValue    float64
}

func NewBucket(AgentIP string) *Bucket {
	return &Bucket{
		AgentIP:          AgentIP,
		PseudoFreeCPU:    -1.0,
		PseudoFreeMemory: -1.0,
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
	totalCPU := 100.0
	freeCPU := 0.0
	for _, container := range m.ContainerDetails {
		freeCPU += container.CPUValue
	}
	return totalCPU - freeCPU
}

func (m *Bucket) GetFreeMemory() float64 {
	freeMem := 0.0
	for _, container := range m.ContainerDetails {
		freeMem += container.MemValue
	}
	return freeMem
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
