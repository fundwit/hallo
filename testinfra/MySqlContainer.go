package testinfra

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
)

type MysqlContainer struct {
	MappedPort nat.Port
	Host       string
	Username   string
	Password   string

	Ctx          context.Context
	ContainerReq testcontainers.ContainerRequest
	Container    testcontainers.Container
}

func (containerService *MysqlContainer) Stop() {
	if containerService.Container != nil {
		err := containerService.Container.Terminate(containerService.Ctx)
		log.Fatalln(err)
	}
}

func NewMysqlContainer() (*MysqlContainer, error) {
	mysqlService := new(MysqlContainer)

	mysqlService.Username = "root"
	mysqlService.Password = "root"

	mysqlService.Ctx = context.Background()
	mysqlService.ContainerReq = testcontainers.ContainerRequest{
		Image: "mysql:5.7",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": mysqlService.Password,
		},
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor:   wait.ForListeningPort("3306/tcp"),
	}
	var err error
	mysqlService.Container, err = testcontainers.GenericContainer(mysqlService.Ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mysqlService.ContainerReq,
		Started:          true,
	})
	if err != nil {
		mysqlService.Stop()
		return nil, err
	}
	mysqlService.Host, err = mysqlService.Container.Host(mysqlService.Ctx)
	if err != nil {
		mysqlService.Stop()
		return nil, err
	}
	mysqlService.MappedPort, err = mysqlService.Container.MappedPort(mysqlService.Ctx, "3306")
	if err != nil {
		mysqlService.Stop()
		return nil, err
	}

	return mysqlService, nil
}
