package model

import (
	"reflect"
	"testing"

	"github.com/evergreen-ci/evergreen/model"
	"github.com/evergreen-ci/utility"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoBuildFromService(t *testing.T) {
	repoRef := model.RepoRef{ProjectRef: model.ProjectRef{
		Id:                  "project",
		Owner:               "my_owner",
		Repo:                "my_repo",
		GithubChecksEnabled: utility.TruePtr(),
		PRTestingEnabled:    utility.FalsePtr(),
		CommitQueue: model.CommitQueueParams{
			MergeMethod: "Squash", // being partially populated shouldn't prevent enabled from being defaulted
		}},
	}
	apiRef := &APIProjectRef{}
	assert.NoError(t, apiRef.BuildFromService(repoRef.ProjectRef))
	// not defaulted yet
	require.NotNil(t, apiRef)
	assert.NotNil(t, apiRef.TaskSync)
	assert.Nil(t, apiRef.GitTagVersionsEnabled)
	assert.Nil(t, apiRef.TaskSync.ConfigEnabled)

	apiRef.DefaultUnsetBooleans()
	assert.True(t, *apiRef.GithubChecksEnabled)
	assert.False(t, *apiRef.PRTestingEnabled)
	require.NotNil(t, apiRef.GitTagVersionsEnabled) // should default
	assert.False(t, *apiRef.GitTagVersionsEnabled)

	assert.NotNil(t, apiRef.TaskSync)
	require.NotNil(t, apiRef.TaskSync.ConfigEnabled) // should default
	assert.False(t, *apiRef.TaskSync.ConfigEnabled)

	require.NotNil(t, apiRef.CommitQueue.Enabled)
	assert.False(t, *apiRef.CommitQueue.Enabled)
}

func TestRecursivelyDefaultBooleans(t *testing.T) {
	type insideStruct struct {
		InsideBool *bool
	}
	type testStruct struct {
		EmptyBool *bool
		TrueBool  *bool
		Strct     insideStruct
		PtrStruct *insideStruct
	}

	myStruct := testStruct{TrueBool: utility.TruePtr()}
	reflected := reflect.ValueOf(&myStruct).Elem()
	recursivelyDefaultBooleans(reflected)

	require.NotNil(t, myStruct.EmptyBool)
	assert.False(t, *myStruct.EmptyBool)
	assert.True(t, *myStruct.TrueBool)
	require.NotNil(t, myStruct.Strct.InsideBool)
	assert.False(t, *myStruct.EmptyBool)
	assert.Nil(t, myStruct.PtrStruct) // shouldn't be affected
}

func TestGitHubDynamicTokenPermissionGroupToService(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	t.Run("Empty permissions should result in no permissions", func(t *testing.T) {
		req := APIGitHubDynamicTokenPermissionGroups{
			APIGitHubDynamicTokenPermissionGroup{
				Name:        utility.ToStringPtr("some-group"),
				Permissions: map[string]string{},
			},
		}
		models, err := req.ToService()
		require.NoError(err)
		require.Len(models, 1)

		assert.Equal("some-group", models[0].Name)
		assert.False(models[0].AllPermissions)
		assert.Nil(models[0].Permissions.Contents)
	})

	t.Run("Invalid permissions should result in an error", func(t *testing.T) {
		req := APIGitHubDynamicTokenPermissionGroups{
			APIGitHubDynamicTokenPermissionGroup{
				Name:        utility.ToStringPtr("some-group"),
				Permissions: map[string]string{"invalid": "invalid"},
			},
		}
		_, err := req.ToService()
		require.ErrorContains(err, "decoding GitHub permissions")
	})

	t.Run("Valid permissions should be converted", func(t *testing.T) {
		req := APIGitHubDynamicTokenPermissionGroups{
			APIGitHubDynamicTokenPermissionGroup{
				Name:        utility.ToStringPtr("some-group"),
				Permissions: map[string]string{"contents": "read", "pull_requests": "write"},
			},
		}
		models, err := req.ToService()
		require.NoError(err)
		require.Len(models, 1)

		assert.Equal("some-group", models[0].Name)
		assert.False(models[0].AllPermissions)
		assert.Equal("read", utility.FromStringPtr(models[0].Permissions.Contents))
		assert.Equal("write", utility.FromStringPtr(models[0].Permissions.PullRequests))
	})

	t.Run("All permissions set should be kept", func(t *testing.T) {
		req := APIGitHubDynamicTokenPermissionGroups{
			APIGitHubDynamicTokenPermissionGroup{
				Name:           utility.ToStringPtr("some-group"),
				AllPermissions: utility.TruePtr(),
			},
		}
		models, err := req.ToService()
		require.NoError(err)
		require.Len(models, 1)

		assert.Equal("some-group", models[0].Name)
		assert.True(models[0].AllPermissions)
	})

	t.Run("If permissions are set and all permissions is true, an error is returned", func(t *testing.T) {
		req := APIGitHubDynamicTokenPermissionGroups{
			APIGitHubDynamicTokenPermissionGroup{
				Name:           utility.ToStringPtr("some-group"),
				Permissions:    map[string]string{"contents": "read", "pull_requests": "write"},
				AllPermissions: utility.TruePtr(),
			},
		}
		_, err := req.ToService()
		require.ErrorContains(err, "a group will all permissions must have no permissions set")
	})
}
