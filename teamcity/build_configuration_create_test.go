package teamcity

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umweltdk/teamcity/types"
	"testing"
	"time"
)

func TestClientCreateBuildConfigurationMock(t *testing.T) {
	client := NewTestClient(newResponse(`{"id": "Empty_Hello", "projectId":"Empty","templateFlag":false,"name":"Hello"}`), nil)

	config := &types.BuildConfiguration{
		ProjectID: "Empty",
		Name:      "Hello",
		Template:  nil,
	}

	err := client.CreateBuildConfiguration(config)
	require.NoError(t, err, "Expected no error")

	assert.Equal(t, "Empty_Hello", config.ID, "Expected create to return ID")
}

func TestClientCreateBuildConfigurationMinimal(t *testing.T) {
	client, err := NewRealTestClient(t)
	require.NoError(t, err, "Expected no error")
	err = client.DeleteBuildConfiguration("Empty_Hello")
	require.NoError(t, err, "Expected no error")

	config := &types.BuildConfiguration{
		ProjectID: "Empty",
		Name:      "Hello",
	}
	err = client.CreateBuildConfiguration(config)
	require.NoError(t, err, "Expected no error")
	require.NotNil(t, config, "Create to return config")

	assert.Equal(t, "Empty_Hello", config.ID, "Expected create to return ID")
	assert.Equal(t, make(types.BuildSteps, 0), config.Steps, "no steps")
}

func TestClientCreateTemplate(t *testing.T) {
	client, err := NewRealTestClient(t)
	require.NoError(t, err, "Expected no error")
	err = client.DeleteBuildConfiguration("Single_Templer")
	require.NoError(t, err, "Expected no error")

	config := &types.BuildConfiguration{
		ProjectID: "Single",
		Name:      "Templer",
		TemplateFlag: true,
    VcsRootEntries: types.VcsRootEntries{
      types.VcsRootEntry{
        VcsRootID: "Single_HttpsGithubComUmweltdkDockerNodeGit",
        CheckoutRules: "+:refs/heads/master\n+:refs/heads/trigger*",
      },
    },
	}
	err = client.CreateBuildConfiguration(config)
	require.NoError(t, err, "Expected no error")
	require.NotNil(t, config, "Create to return config")

	assert.Equal(t, "Single_Templer", config.ID, "Expected create to return ID")
	assert.Equal(t, make(types.BuildSteps, 0), config.Steps, "no steps")
	assert.Equal(t, types.VcsRootEntries{
      types.VcsRootEntry{
      	ID: "Single_HttpsGithubComUmweltdkDockerNodeGit",
        VcsRootID: "Single_HttpsGithubComUmweltdkDockerNodeGit",
        CheckoutRules: "+:refs/heads/master\n+:refs/heads/trigger*",
      },
    }, config.VcsRootEntries, "vcs root entries")

	loaded, err := client.GetBuildConfiguration("Single_Templer")
	assert.Equal(t, true, config.TemplateFlag, "Expected template")
	assert.Equal(t, loaded.VcsRootEntries, config.VcsRootEntries, "vcs root entries")
}

func TestClientCreateBuildConfigurationTemplateFull(t *testing.T) {
	client, err := NewRealTestClient(t)
	require.NoError(t, err, "Expected no error")
	err = client.DeleteBuildConfiguration("Empty_TemplateFull")
	require.NoError(t, err, "Expected no error")
	time.Sleep(5 * time.Second)

	config := &types.BuildConfiguration{
		ProjectID: "Empty",
		Name:      "Template Full",
		Template: &types.TemplateID{
			ID: "Tempy",
		},
		Steps: types.BuildSteps{
			types.BuildStep{
				Name: "Muh",
				Type: "simpleRunner",
				Properties: types.Properties{
					"script.content":     types.Property{"env", nil},
					"teamcity.step.mode": types.Property{"default", nil},
					"use.custom.script":  types.Property{"true", nil},
				},
			},
			types.BuildStep{
				Name: "Env",
				Type: "simpleRunner",
				Properties: types.Properties{
					"script.content":     types.Property{"env", nil},
					"teamcity.step.mode": types.Property{"default", nil},
					"use.custom.script":  types.Property{"true", nil},
				},
			},
		},
	}
	err = client.CreateBuildConfiguration(config)
	require.NoError(t, err, "Expected no error")
	require.NotNil(t, config, "Create to return config")

	var id1 string
	var id2 string
	if len(config.Steps) >= 3 {
		id1 = config.Steps[1].ID
		id2 = config.Steps[2].ID
	}

	assert.Equal(t, "Empty_TemplateFull", config.ID, "Supplied ID")
	assert.Equal(t, types.BuildSteps{
		types.BuildStep{
			ID:   "RUNNER_3",
			Name: "Env",
			Type: "simpleRunner",
			Properties: types.Properties{
				"script.content":     types.Property{"env", nil},
				"teamcity.step.mode": types.Property{"default", nil},
				"use.custom.script":  types.Property{"true", nil},
			},
		},
		types.BuildStep{
			ID:   id1,
			Name: "Muh",
			Type: "simpleRunner",
			Properties: types.Properties{
				"script.content":     types.Property{"env", nil},
				"teamcity.step.mode": types.Property{"default", nil},
				"use.custom.script":  types.Property{"true", nil},
			},
		},
		types.BuildStep{
			ID:   id2,
			Name: "Env (1)",
			Type: "simpleRunner",
			Properties: types.Properties{
				"script.content":     types.Property{"env", nil},
				"teamcity.step.mode": types.Property{"default", nil},
				"use.custom.script":  types.Property{"true", nil},
			},
		},
	}, config.Steps, "Build steps")
}

func TestClientCreateBuildConfigurationTemplateReorder(t *testing.T) {
	client, err := NewRealTestClient(t)
	client.SkipOlder(t, 10, 0)
	require.NoError(t, err, "Expected no error")
	err = client.DeleteBuildConfiguration("Empty_TemplateReorder")
	require.NoError(t, err, "Expected no error")
	time.Sleep(10 * time.Second)

	config := &types.BuildConfiguration{
		ProjectID: "Empty",
		Name:      "Template Reorder",
		Template: &types.TemplateID{
			ID: "Tempy",
		},
		Steps: types.BuildSteps{
			types.BuildStep{
				Name: "Muh",
				Type: "simpleRunner",
				Properties: types.Properties{
					"script.content":     types.Property{"env", nil},
					"teamcity.step.mode": types.Property{"default", nil},
					"use.custom.script":  types.Property{"true", nil},
				},
			},
			types.BuildStep{
				ID:   "RUNNER_3",
				Name: "Env",
				Type: "simpleRunner",
				Properties: types.Properties{
					"script.content":     types.Property{"env", nil},
					"teamcity.step.mode": types.Property{"default", nil},
					"use.custom.script":  types.Property{"true", nil},
				},
			},
		},
	}
	err = client.CreateBuildConfiguration(config)
	require.NoError(t, err, "Expected no error")
	require.NotNil(t, config, "Create to return config")

	var id1 string
	if len(config.Steps) >= 1 {
		id1 = config.Steps[0].ID
	}

	assert.Equal(t, "Empty_TemplateReorder", config.ID, "Supplied ID")
	assert.Equal(t, types.BuildSteps{
		types.BuildStep{
			ID:   id1,
			Name: "Muh",
			Type: "simpleRunner",
			Properties: types.Properties{
				"script.content":     types.Property{"env", nil},
				"teamcity.step.mode": types.Property{"default", nil},
				"use.custom.script":  types.Property{"true", nil},
			},
		},
		types.BuildStep{
			ID:   "RUNNER_3",
			Name: "Env",
			Type: "simpleRunner",
			Properties: types.Properties{
				"script.content":     types.Property{"env", nil},
				"teamcity.step.mode": types.Property{"default", nil},
				"use.custom.script":  types.Property{"true", nil},
			},
		},
	}, config.Steps, "Build steps")
}

func TestClientCreateBuildConfigurationFull(t *testing.T) {
	client, err := NewRealTestClient(t)
	require.NoError(t, err, "Expected no error")
	err = client.DeleteBuildConfiguration("Empty_Daws")
	require.NoError(t, err, "Expected no error")

	config := &types.BuildConfiguration{
		ID:        "Empty_Daws",
		ProjectID: "Empty",
		Name:      "Maws",
		Steps: types.BuildSteps{
			types.BuildStep{
				Name: "Muh",
				Type: "simpleRunner",
				Properties: types.Properties{
					"script.content":     types.Property{"env", nil},
					"teamcity.step.mode": types.Property{"default", nil},
					"use.custom.script":  types.Property{"true", nil},
				},
			},
		},
	}
	err = client.CreateBuildConfiguration(config)
	require.NoError(t, err, "Expected no error")
	require.NotNil(t, config, "Create to return config")

	var id1 string
	if len(config.Steps) >= 1 {
		id1 = config.Steps[0].ID
	}
	assert.Equal(t, "Empty_Daws", config.ID, "Supplied ID")
	assert.Equal(t, types.BuildSteps{
		types.BuildStep{
			ID:   id1,
			Name: "Muh",
			Type: "simpleRunner",
			Properties: types.Properties{
				"script.content":     types.Property{"env", nil},
				"teamcity.step.mode": types.Property{"default", nil},
				"use.custom.script":  types.Property{"true", nil},
			},
		},
	}, config.Steps, "Build steps")
}

func TestClientCreateBuildConfigurationUsedID(t *testing.T) {
	client, err := NewRealTestClient(t)
	require.NoError(t, err, "Expected no error")
	client.retries = 1

	config := &types.BuildConfiguration{
		ID:        "Single_Normal",
		ProjectID: "Single",
		Name:      "Hej Med Dig",
	}

	err = client.CreateBuildConfiguration(config)
	assert.Error(t, err, "Expected error")
}

func TestClientCreateBuildConfigurationUsedName(t *testing.T) {
	client, err := NewRealTestClient(t)
	require.NoError(t, err, "Expected no error")
	client.retries = 1

	config := &types.BuildConfiguration{
		ProjectID: "Single",
		Name:      "Normal",
	}

	err = client.CreateBuildConfiguration(config)
	assert.Error(t, err, "Expected error")
}

func TestClientCreateBuildConfigurationUsedNameExplicitID(t *testing.T) {
	client, err := NewRealTestClient(t)
	require.NoError(t, err, "Expected no error")
	client.retries = 1

	config := &types.BuildConfiguration{
		ID:        "Single_Dubie",
		ProjectID: "Single",
		Name:      "Hello",
	}

	err = client.CreateBuildConfiguration(config)
	assert.Error(t, err, "Expected error")
}