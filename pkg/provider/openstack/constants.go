package openstack

const (

	// Create params
	imageID          string = "image"
	imageIDDesc      string = "OpenStack image identifier"
	instanceType     string = "flavor"
	instanceTypeDesc string = "OpenStack flavor type for the machine running the cluster. Default is m1.xlarge."
	diskSize         string = "disk-size"
	diskSizeDesc     string = "Disk size in GB for the machine running the cluster. Default is 100."
	networkName      string = "network"
	networkNameDesc  string = "OpenStack network name for the machine running the cluster."

	// default values
	ocpInstanceType               string = "m1.xlarge"
	ocpDefaultRootBlockDeviceSize int    = 100
)
