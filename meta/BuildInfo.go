package meta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type BuildInfo struct {
	BuildTime   time.Time    `json:"buildTime"`
	Version     string       `json:"version"`
	SourceCodes []SourceCode `json:"sourceCodes"`
}

type SourceCode struct {
	Repository string    `json:"repository"`
	Ref        string    `json:"ref"`
	Reversion  Reversion `json:"reversion"`
}

type Reversion struct {
	Id        string    `json:"id"`
	Author    string    `json:"author"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func GetBuildInfo() *BuildInfo {
	bytes := readBuildInfo()
	if bytes == nil {
		return nil
	}

	return parseBuildInfo(bytes)
}

func parseBuildInfo(bytes []byte) *BuildInfo {
	bi := &BuildInfo{}
	err := json.Unmarshal(bytes, bi)
	if err != nil {
		log.Printf("unexpect build info content: \"%s\", %v", string(bytes), err)
		return nil
	}
	return bi
}

func readBuildInfo() []byte {
	buildInfoFile := "./buildInfo.json"
	bytes, err := ioutil.ReadFile(buildInfoFile)
	if err != nil {
		log.Printf("build info file not found: \"%s\", %v", buildInfoFile, err)
		return nil
	}
	return bytes
}
