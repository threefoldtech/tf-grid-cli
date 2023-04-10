package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/grid3-go/graphql"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/router/mocks"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

func TestK8s(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cl := mocks.NewMockTFGridClient(ctrl)

	r := Router{
		client: cl,
	}

	t.Run("k8s_deploy_success", func(t *testing.T) {
		projectName := "project1"
		model := K8sCluster{
			Master: &K8sNode{
				Name:      "master",
				NodeID:    1,
				DiskSize:  10,
				PublicIP:  false,
				PublicIP6: true,
				Planetary: true,
				Flist:     "hamada",
				CPU:       1,
				Memory:    2,
			},
			Workers: []K8sNode{
				{
					Name:      "w1",
					NodeID:    2,
					DiskSize:  10,
					PublicIP:  true,
					PublicIP6: false,
					Planetary: false,
					Flist:     "hamada2",
					CPU:       3,
					Memory:    5,
				},
			},
			Name:   "cluster1",
			Token:  "token1",
			SSHKey: "key1",
		}

		want := K8sCluster{
			Master: &K8sNode{
				Name:        "master",
				NodeID:      1,
				DiskSize:    10,
				PublicIP:    false,
				PublicIP6:   true,
				Planetary:   true,
				Flist:       "hamada",
				CPU:         1,
				Memory:      2,
				ComputedIP4: "ip4",
				ComputedIP6: "ip6",
				WGIP:        "wgip",
				YggIP:       "yggip",
			},
			Workers: []K8sNode{
				{
					Name:        "w1",
					NodeID:      2,
					DiskSize:    10,
					PublicIP:    true,
					PublicIP6:   false,
					Planetary:   false,
					Flist:       "hamada2",
					CPU:         3,
					Memory:      5,
					ComputedIP4: "ip4",
					ComputedIP6: "ip6",
					WGIP:        "wgip",
					YggIP:       "yggip",
				},
			},
			Name:        "cluster1",
			Token:       "token1",
			NetworkName: fmt.Sprintf("%s_network", model.Name),
			SSHKey:      "key1",
		}

		cl.
			EXPECT().
			GetProjectContracts(gomock.Any(), projectName).
			Return(graphql.Contracts{}, nil)

		ipRange, err := gridtypes.ParseIPNet("10.1.0.0/16")
		assert.NoError(t, err)

		znet := workloads.ZNet{
			Name:         fmt.Sprintf("%s_network", model.Name),
			Nodes:        []uint32{2, 1},
			IPRange:      ipRange,
			SolutionType: projectName,
		}

		cl.EXPECT().DeployNetwork(gomock.Any(), &znet).Return(nil, nil)

		model.NetworkName = fmt.Sprintf("%s_network", model.Name)
		k8s := newK8sClusterFromModel(model, projectName)

		retK8s := workloads.K8sCluster{
			Master: &workloads.K8sNode{
				Name:        "master",
				Node:        1,
				DiskSize:    10,
				PublicIP:    false,
				PublicIP6:   true,
				Planetary:   true,
				Flist:       "hamada",
				CPU:         1,
				Memory:      2,
				ComputedIP:  "ip4",
				ComputedIP6: "ip6",
				IP:          "wgip",
				YggIP:       "yggip",
			},
			Workers: []workloads.K8sNode{
				{
					Name:        "w1",
					Node:        2,
					DiskSize:    10,
					PublicIP:    true,
					PublicIP6:   false,
					Planetary:   false,
					Flist:       "hamada2",
					CPU:         3,
					Memory:      5,
					ComputedIP:  "ip4",
					ComputedIP6: "ip6",
					IP:          "wgip",
					YggIP:       "yggip",
				},
			},
			SolutionType: projectName,
			Token:        "token1",
			NetworkName:  fmt.Sprintf("%s_network", model.Name),
			SSHKey:       "key1",
		}

		cl.EXPECT().DeployK8sCluster(gomock.Any(), &k8s).Return(&retK8s, nil)

		got, err := r.k8sDeploy(context.Background(), model, projectName)
		assert.NoError(t, err)

		assert.Equal(t, want, got)
	})
}
