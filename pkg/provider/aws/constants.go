package aws

const (

	// Create params
	amiID            string = "aws-ami-id"
	amiIDDesc        string = "AMI identifier"
	instanceType     string = "aws-instance-type"
	instanceTypeDesc string = "Instance type for the machine running the cluster. Default is c6a.2xlarge."
	diskSize         string = "aws-disk-size"
	diskSizeDesc     string = "Disk size in GB for the machine running the cluster. Default is 100."

	// default values
	ocpInstanceType               string = "c6a.2xlarge"
	ocpDefaultRootBlockDeviceSize int    = 100
)
