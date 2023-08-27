package v1

const (
	ClusterProgressing = "progressing"

	ServiceMaintainSuffix  = "-mt"
	ServiceReadWriteSuffix = "-rw"

	labelPrefix      = "opengemini-operator.opengemini.org/"
	LabelCluster     = labelPrefix + "cluster"
	LabelInstance    = labelPrefix + "instance"
	LabelInstanceSet = labelPrefix + "instance-set"
	LabelConfigHash  = labelPrefix + "config-hash"

	AdminUserSecretSuffix = "-admin"
	DefaultAdminUsername  = "admin"
)
