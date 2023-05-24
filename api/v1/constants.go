package v1

const (
	ClusterProgressing = "progressing"

	ServiceMaintainSuffix  = "-mt"
	ServiceReadWriteSuffix = "-rw"

	labelPrefix  = "opengemini-operator.opengemini.org/"
	LabelCluster = labelPrefix + "cluster"
)
