package gcp

const (

	// Create params
	imageID          string = "gcp-image-id"
	imageIDDesc      string = "GCP image identifier"
	instanceType     string = "gcp-instance-type"
	instanceTypeDesc string = "Instance type for the machine running the cluster. Default is n1-standard-8."
	diskSize         string = "gcp-disk-size"
	diskSizeDesc     string = "Disk size in GB for the machine running the cluster. Default is 100."

	// default values
	ocpInstanceType               string = "n1-standard-8"
	ocpDefaultRootBlockDeviceSize int    = 100
)
