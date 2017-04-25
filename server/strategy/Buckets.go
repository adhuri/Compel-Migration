package strategy

import (
	"sort"

	"github.com/Sirupsen/logrus"
)

type Bucket struct {
	AgentIP          string
	ContainerDetails []*Container
	PseudoFreeCPU    float32 // Psuedo - This changes every time you make a local decision of moving containers across
	PseudoFreeMemory float32 // Psuedo - This changes every time you make a local decision of moving containers across
	CpuThreshold     int64
	MemThreshold     int64
}

type Container struct {
	ContainerID      string
	CPUValue         float32
	MemValue         float32
	movableContainer bool
}

func NewBucket(AgentIP string, cpuThreshold int64, memThreshold int64) *Bucket {
	return &Bucket{
		AgentIP:          AgentIP,
		PseudoFreeCPU:    999.0, // 999 Free CPU not possible
		PseudoFreeMemory: 999.0, // 999 Free Memory not possible
		ContainerDetails: make([]*Container, 0),
		CpuThreshold:     cpuThreshold,
		MemThreshold:     memThreshold,
	}
}

func NewContainer(ContainerID string, CPUValue float32, MemValue float32, movableStatus bool) *Container {
	return &Container{
		ContainerID:      ContainerID,
		CPUValue:         CPUValue,
		MemValue:         MemValue,
		movableContainer: movableStatus,
	}
}

func (m *Bucket) GetFreeCPU() float32 {

	if m.PseudoFreeCPU == 999.0 {
		//freeCPU := float32(100.0)
		freeCPU := float32(m.CpuThreshold)

		for _, container := range m.ContainerDetails {
			freeCPU -= container.CPUValue
		}
		m.PseudoFreeCPU = freeCPU
	}
	return m.PseudoFreeCPU
}

func (m *Bucket) GetFreeMemory() float32 {

	if m.PseudoFreeMemory == 999.0 {

		//freeMem := float32(100.0)

		freeMem := float32(m.MemThreshold)
		for _, container := range m.ContainerDetails {
			freeMem -= container.MemValue
		}
		m.PseudoFreeMemory = freeMem
	}
	return m.PseudoFreeMemory
}

func (m *Bucket) GetValue(metric string) float32 {
	if metric == "cpu" {
		return m.GetFreeCPU()
	} else if metric == "memory" {
		return m.GetFreeMemory()

	}
	return 999.0
}

func (c *Container) GetValue(metric string) float32 {
	if metric == "cpu" {
		return c.CPUValue
	} else if metric == "memory" {
		return c.MemValue

	}
	return 999.0
}

func (m *Bucket) PrintBucket(log *logrus.Logger) {
	log.Infoln("Agent IP ", m.AgentIP)
	log.Infoln("Total Containers ", len(m.ContainerDetails))

	log.Infoln("Free CPU ", m.GetFreeCPU())
	log.Infoln("Free Memory ", m.GetFreeMemory())

	for containerIndex, container := range m.ContainerDetails {
		log.Infoln("Container ", containerIndex, " : ContainerID ", container.ContainerID, " : CPU ", container.GetValue("cpu"), " : Memory ", container.GetValue("memory"), ": Movable Container ", container.movableContainer)
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

// Sorting
type sortBucketByCPU []*Bucket
type sortBucketByMemory []*Bucket

type sortContainerByCPU []*Container
type sortContainerByMemory []*Container

func sortBucketsAsc(buckets []*Bucket, metric string) {
	// Sorted Ascending - to change change Less Method
	if metric == "cpu" {
		sort.Sort(sortBucketByCPU(buckets))
	} else if metric == "memory" {
		sort.Sort(sortBucketByMemory(buckets))

	}
	//return buckets
}

func sortContainersDesc(containers []*Container, metric string) {
	// Sort Descending
	if metric == "cpu" {
		sort.Sort(sortContainerByCPU(containers))
	} else if metric == "memory" {
		sort.Sort(sortContainerByMemory(containers))

	}

	//return containers
}

//Sort by CPU

func (sCPU sortBucketByCPU) Len() int {
	return len(sCPU)
}

func (sCPU sortBucketByCPU) Less(i, j int) bool {
	return sCPU[i].GetFreeCPU() < sCPU[j].GetFreeCPU()
}

func (sCPU sortBucketByCPU) Swap(i, j int) {
	sCPU[i], sCPU[j] = sCPU[j], sCPU[i]
}

// Sort By Memory

func (sMem sortBucketByMemory) Len() int {
	return len(sMem)
}

func (sMem sortBucketByMemory) Less(i, j int) bool {
	return sMem[i].GetFreeMemory() < sMem[j].GetFreeMemory()
}

func (sMem sortBucketByMemory) Swap(i, j int) {
	sMem[i], sMem[j] = sMem[j], sMem[i]
}

// Sort Container by CPU

func (cCPU sortContainerByCPU) Len() int {
	return len(cCPU)
}

func (cCPU sortContainerByCPU) Less(i, j int) bool {

	// Descending Sort >
	return cCPU[i].CPUValue > cCPU[j].CPUValue
}

func (cCPU sortContainerByCPU) Swap(i, j int) {
	cCPU[i], cCPU[j] = cCPU[j], cCPU[i]
}

// Sort Container by Memory

// Sort Container by CPU

func (cCPU sortContainerByMemory) Len() int {
	return len(cCPU)
}

func (cCPU sortContainerByMemory) Less(i, j int) bool {

	// Descending Sort >
	return cCPU[i].MemValue > cCPU[j].MemValue
}

func (cCPU sortContainerByMemory) Swap(i, j int) {
	cCPU[i], cCPU[j] = cCPU[j], cCPU[i]
}
