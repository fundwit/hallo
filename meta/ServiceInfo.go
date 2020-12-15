package meta

type ServiceInfo struct {
	ServiceInstance string `json:"serviceInstance"`
	ServiceName     string `json:"serviceName"`
	Description     string `json:"description"`

	BuildInfo *BuildInfo `json:"buildInfo"`
}

var serviceInfo *ServiceInfo

func GetServiceInfo() *ServiceInfo {
	if serviceInfo == nil {
		serviceInfo = &ServiceInfo{
			ServiceName:     "hallo",
			ServiceInstance: "hallo-xxx", // 优先从环境变量读取，其次自动生成Id
			Description:     "An authentication service",

			BuildInfo: GetBuildInfo(),
		}
	}
	return serviceInfo
}
