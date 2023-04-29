package hack

import (
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
)

var b = cloudbuildpb.Build{
	Name:         "",
	Id:           "",
	ProjectId:    "",
	Status:       0,
	StatusDetail: "",
	Source:       nil,
	Steps: []*cloudbuildpb.BuildStep{
		{
			Name:       "",
			Env:        nil,
			Args:       nil,
			Dir:        "",
			Id:         "",
			WaitFor:    nil,
			Entrypoint: "",
			SecretEnv:  nil,
			Volumes:    nil,
			Timing:     nil,
			PullTiming: nil,
			Timeout:    nil,
			Status:     0,
			Script:     "",
		},
	},
	Options: nil,
}
