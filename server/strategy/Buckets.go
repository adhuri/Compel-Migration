package strategy

import "github.com/Sirupsen/logrus"

type Bucket struct {
	AgentIP          string
	ContainerDetails []*Container
	PseudoFreeCPU    float32 // Psuedo - This changes every time you make a local decision of moving containers across
	PseudoFreeMemory float32 // Psuedo - This changes every time you make a local decision of moving containers across
}

type Container struct {
	ContainerID     string
	CPUValue        float32
	MemValue        float32
	staticContainer bool
}

func NewBucket(AgentIP string) *Bucket {
	return &Bucket{
		AgentIP:          AgentIP,
		PseudoFreeCPU:    999.0, // 999 Free CPU not possible
		PseudoFreeMemory: 999.0, // 999 Free Memory not possible
		ContainerDetails: make([]*Container, 0),
	}
}

func NewContainer(ContainerID string, CPUValue float32, MemValue float32) *Container {
	return &Container{
		ContainerID:     ContainerID,
		CPUValue:        CPUValue,
		MemValue:        MemValue,
		staticContainer: false,
	}
}

func (m *Bucket) GetFreeCPU() float32 {

	if m.PseudoFreeCPU == 999.0 {
		freeCPU := float32(100.0)

		for _, container := range m.ContainerDetails {
			freeCPU -= container.CPUValue
		}
		m.PseudoFreeCPU = freeCPU
	}
	return m.PseudoFreeCPU
}

func (m *Bucket) GetFreeMemory() float32 {

	if m.PseudoFreeMemory == 999.0 {

		freeMem := float32(100.0)
		for _, container := range m.ContainerDetails {
			freeMem -= container.MemValue
		}
		m.PseudoFreeMemory = freeMem
	}
	return m.PseudoFreeMemory
}

func (m *Bucket) PrintBucket(log *logrus.Logger) {
	log.Infoln("Agent IP ", m.AgentIP)
	log.Infoln("Total Containers ", len(m.ContainerDetails))

	log.Infoln("Free CPU ", m.PseudoFreeCPU)
	log.Infoln("Free Memory ", m.PseudoFreeMemory)

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
