package service

import (
	"context"
	"fmt"

	{{cookiecutter.app_name|lower}} "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
)

type svc struct {
	prefix string
}

func (s *svc) Echo(_ context.Context, req *{{cookiecutter.app_name|lower}}.EchoRequest) (*{{cookiecutter.app_name|lower}}.EchoResponse, error) {
	return &{{cookiecutter.app_name|lower}}.EchoResponse{
		Msg: fmt.Sprintf("%s: %s", s.prefix, req.GetMsg()),
	}, nil
}

// Creates a new Service
func New() {{cookiecutter.app_name|lower}}.{{cookiecutter.service_name}}ServiceServer {

	return &svc{
		prefix: config.Get().Prefix,
	}
}
