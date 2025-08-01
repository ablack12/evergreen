package model

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/model/build"
	"github.com/evergreen-ci/evergreen/model/distro"
	"github.com/evergreen-ci/evergreen/model/event"
	"github.com/evergreen-ci/evergreen/model/patch"
	"github.com/evergreen-ci/evergreen/model/task"
	"github.com/evergreen-ci/evergreen/model/user"
	"github.com/evergreen-ci/evergreen/util"
	"github.com/evergreen-ci/utility"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v2"
)

type TaskVariantPairs struct {
	ExecTasks    TVPairSet
	DisplayTasks TVPairSet
}

// VariantTasksToTVPairs takes a set of variants and tasks (from both the old and new
// request formats) and builds a universal set of pairs
// that can be used to expand the dependency tree.
func VariantTasksToTVPairs(in []patch.VariantTasks) TaskVariantPairs {
	out := TaskVariantPairs{}
	for _, vt := range in {
		for _, t := range vt.Tasks {
			out.ExecTasks = append(out.ExecTasks, TVPair{vt.Variant, t})
		}
		for _, dt := range vt.DisplayTasks {
			out.DisplayTasks = append(out.DisplayTasks, TVPair{vt.Variant, dt.Name})
			for _, et := range dt.ExecTasks {
				out.ExecTasks = append(out.ExecTasks, TVPair{vt.Variant, et})
			}
		}
	}
	return out
}

// TVPairsToVariantTasks takes a list of TVPairs (task/variant pairs), groups the tasks
// for the same variant together under a single list, and returns all the variant groups
// as a set of patch.VariantTasks.
func (tvp *TaskVariantPairs) TVPairsToVariantTasks() []patch.VariantTasks {
	vtMap := map[string]patch.VariantTasks{}
	for _, pair := range tvp.ExecTasks {
		vt := vtMap[pair.Variant]
		vt.Variant = pair.Variant
		vt.Tasks = append(vt.Tasks, pair.TaskName)
		vtMap[pair.Variant] = vt
	}
	for _, dt := range tvp.DisplayTasks {
		variant := vtMap[dt.Variant]
		variant.Variant = dt.Variant
		variant.DisplayTasks = append(variant.DisplayTasks, patch.DisplayTask{Name: dt.TaskName})
		vtMap[dt.Variant] = variant
	}
	vts := []patch.VariantTasks{}
	for _, vt := range vtMap {
		vts = append(vts, vt)
	}
	return vts
}

// ValidateTVPairs checks that all of a set of variant/task pairs exist in a given project.
func ValidateTVPairs(p *Project, in []TVPair) error {
	for _, pair := range in {
		v := p.FindTaskForVariant(pair.TaskName, pair.Variant)
		if v == nil {
			return errors.Errorf("task name '%s' in variant '%s' does not exist", pair.TaskName, pair.Variant)
		}
	}
	return nil
}

// Given a patch version and a list of variant/task pairs, creates the set of new builds that
// do not exist yet out of the set of pairs, and adds tasks for builds which already exist.
func addNewTasksAndBuildsForPatch(ctx context.Context, p *patch.Patch, creationInfo TaskCreationInfo, caller string) error {
	existingBuilds, err := build.Find(ctx, build.ByIds(creationInfo.Version.BuildIds).WithFields(build.IdKey, build.BuildVariantKey, build.CreateTimeKey, build.RequesterKey))
	if err != nil {
		return err
	}
	activatedTasks, activatedDeps, err := addNewBuilds(ctx, creationInfo, existingBuilds)
	if err != nil {
		return errors.Wrap(err, "adding new builds")
	}
	activatedTasksInExistingBuilds, activatedDepsInExistingBuilds, err := addNewTasksToExistingBuilds(ctx, creationInfo, existingBuilds, caller)
	if err != nil {
		return errors.Wrap(err, "adding new tasks")
	}
	err = activateExistingInactiveTasks(ctx, creationInfo, existingBuilds, caller)
	totalActivatedTasks := len(activatedTasks) + len(activatedTasksInExistingBuilds) + len(activatedDeps) + len(activatedDepsInExistingBuilds)
	if totalActivatedTasks > 0 {
		// It's necessary to check if new tasks were actually created because
		// it's valid for a user to reconfigure a patch without changing the
		// patch's tasks at all (e.g. if they just updated the patch
		// description). A patch is only considered reconfigured if new tasks
		// were added.
		if err := p.SetIsReconfigured(ctx, true); err != nil {
			return errors.Wrap(err, "marking patch as reconfigured")
		}
	}
	return errors.Wrap(err, "activating existing inactive tasks")
}

type PatchUpdate struct {
	Description         string               `json:"description"`
	Caller              string               `json:"caller"`
	Parameters          []patch.Parameter    `json:"parameters,omitempty"`
	PatchTriggerAliases []string             `json:"patch_trigger_aliases,omitempty"`
	VariantsTasks       []patch.VariantTasks `json:"variants_tasks,omitempty"`
}

// ConfigurePatch validates and creates the updated tasks/variants if given, and updates description if needed.
// Returns an http status code and error.
func ConfigurePatch(ctx context.Context, settings *evergreen.Settings, p *patch.Patch, version *Version, proj *ProjectRef, patchUpdateReq PatchUpdate) (int, error) {
	var err error
	project, _, err := FindAndTranslateProjectForPatch(ctx, settings, p)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "unmarshalling project config")
	}

	addDisplayTasksToPatchReq(&patchUpdateReq, *project)
	tasks := VariantTasksToTVPairs(patchUpdateReq.VariantsTasks)
	tasks.ExecTasks, err = IncludeDependencies(project, tasks.ExecTasks, p.GetRequester(), nil)
	grip.Warning(message.WrapError(err, message.Fields{
		"message": "error including dependencies for patch",
		"patch":   p.Id.Hex(),
	}))
	if err = ValidateTVPairs(project, tasks.ExecTasks); err != nil {
		return http.StatusBadRequest, err
	}

	// only modify parameters if the patch hasn't been finalized
	if len(patchUpdateReq.Parameters) > 0 && p.Version == "" {
		if err = p.SetParameters(ctx, patchUpdateReq.Parameters); err != nil {
			return http.StatusInternalServerError, errors.Wrap(err, "setting patch parameters")
		}
	}
	// update the description for both reconfigured and new patches
	if err = p.SetDescription(ctx, patchUpdateReq.Description); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "setting description")
	}

	patchVariantTasks := tasks.TVPairsToVariantTasks()
	if len(patchVariantTasks) > 0 {
		if err = p.SetVariantsTasks(ctx, patchVariantTasks); err != nil {
			return http.StatusInternalServerError, errors.Wrap(err, "setting description")
		}
	}

	if p.Version != "" {
		// This patch has already been finalized, just add the new builds and tasks
		if version == nil {
			return http.StatusInternalServerError, errors.Errorf("finding patch for version '%s'", p.Version)
		}

		if version.Message != patchUpdateReq.Description {
			if err = UpdateVersionMessage(ctx, p.Version, patchUpdateReq.Description); err != nil {
				return http.StatusInternalServerError, errors.Wrap(err, "setting version message")
			}
		}

		if len(patchVariantTasks) > 0 {
			// Add new tasks to existing builds, if necessary
			creationInfo := TaskCreationInfo{
				Project:        project,
				ProjectRef:     proj,
				Version:        version,
				Pairs:          tasks,
				ActivationInfo: specificActivationInfo{},
				GeneratedBy:    "",
			}
			err = addNewTasksAndBuildsForPatch(ctx, p, creationInfo, patchUpdateReq.Caller)
			if err != nil {
				return http.StatusInternalServerError, errors.Wrapf(err, "creating new tasks/builds for version '%s'", version.Id)
			}
		}
	}

	if p.IsGithubPRPatch() {
		numCheckRuns := project.GetNumCheckRunsFromVariantTasks(p.VariantsTasks)
		checkRunLimit := settings.GitHubCheckRun.CheckRunLimit
		if numCheckRuns > checkRunLimit {
			return http.StatusInternalServerError, errors.Errorf("total number of checkRuns (%d) exceeds maximum limit (%d)", numCheckRuns, checkRunLimit)
		}
	}

	return http.StatusOK, nil
}

func addDisplayTasksToPatchReq(req *PatchUpdate, p Project) {
	for i, vt := range req.VariantsTasks {
		bv := p.FindBuildVariant(vt.Variant)
		if bv == nil {
			continue
		}
		for i := len(vt.Tasks) - 1; i >= 0; i-- {
			task := vt.Tasks[i]
			displayTask := bv.GetDisplayTask(task)
			if displayTask == nil {
				continue
			}
			vt.Tasks = append(vt.Tasks[:i], vt.Tasks[i+1:]...)
			vt.DisplayTasks = append(vt.DisplayTasks, *displayTask)
		}
		req.VariantsTasks[i] = vt
	}
}

func getPatchedProjectYAML(ctx context.Context, projectRef *ProjectRef, opts *GetProjectOpts, p *patch.Patch) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, fetchProjectFilesTimeout)
	defer cancel()
	env := evergreen.GetEnvironment()

	// try to get the remote project file data at the requested revision
	var projectFileBytes []byte
	hash := p.Githash
	if p.IsGithubPRPatch() {
		hash = p.GithubPatchData.HeadHash
	}
	if p.IsMergeQueuePatch() {
		hash = p.GithubMergeData.HeadSHA
	}
	opts.Revision = hash

	path := projectRef.RemotePath
	if p.Path != "" && !p.IsGithubPRPatch() && !p.IsMergeQueuePatch() {
		path = p.Path
	}
	opts.RemotePath = path
	opts.PatchOpts.env = env
	projectFileBytes, err := getFileForPatchDiff(ctx, *opts)
	if err != nil {
		return nil, errors.Wrap(err, "fetching remote configuration file")
	}

	// apply remote configuration patch if needed
	if p.ShouldPatchFileWithDiff(path) {
		opts.ReadFileFrom = ReadFromPatchDiff
		projectFileBytes, err = MakePatchedConfig(ctx, *opts, string(projectFileBytes))
		if err != nil {
			return nil, errors.Wrap(err, "patching remote configuration file")
		}
	}

	return projectFileBytes, nil
}

// GetPatchedProject creates and validates a project by fetching latest commit
// information from GitHub and applying the patch to the latest remote
// configuration. Also returns the condensed yaml string for storage. The error
// returned can be a validation error.
func GetPatchedProject(ctx context.Context, settings *evergreen.Settings, p *patch.Patch) (*Project, *PatchConfig, error) {
	if p.Version != "" {
		return nil, nil, errors.Errorf("patch '%s' already finalized", p.Version)
	}

	if p.ProjectStorageMethod != "" {
		// If the patch has already been created but has not been finalized, it
		// has already stored the parser project document
		project, pp, err := FindAndTranslateProjectForPatch(ctx, settings, p)
		if err != nil {
			return nil, nil, errors.Wrap(err, "finding and translating project")
		}

		ppOut, err := yaml.Marshal(pp)
		if err != nil {
			return nil, nil, errors.Wrap(err, "marshalling parser project into YAML")
		}
		patchedParserProjectYAML := string(ppOut)

		patchConfig := &PatchConfig{
			PatchedParserProjectYAML: patchedParserProjectYAML,
			PatchedParserProject:     pp,
			PatchedProjectConfig:     p.PatchedProjectConfig,
		}
		return project, patchConfig, nil
	}

	projectRef, opts, err := getLoadProjectOptsForPatch(ctx, p)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fetching project options for patch")
	}

	projectFileBytes, err := getPatchedProjectYAML(ctx, projectRef, opts, p)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting patched project file as YAML")
	}

	project := &Project{}
	pp, err := LoadProjectInto(ctx, projectFileBytes, opts, p.Project, project)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	// LoadProjectInto does not set the parser project's identifier, so set it
	// here.
	pp.Identifier = utility.ToStringPtr(p.Project)
	ppOut, err := yaml.Marshal(pp)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshalling parser project into YAML")
	}

	var pc string
	if projectRef.IsVersionControlEnabled() {
		pc, err = getProjectConfigYAML(p, projectFileBytes)
		if err != nil {
			return nil, nil, errors.Wrap(err, "getting patched project config")
		}
	}

	patchConfig := &PatchConfig{
		PatchedParserProjectYAML: string(ppOut),
		PatchedParserProject:     pp,
		PatchedProjectConfig:     pc,
	}
	return project, patchConfig, nil
}

// GetPatchedProjectConfig returns the project configuration by fetching the
// latest commit information from GitHub and applying the patch to the latest
// remote configuration. The error returned can be a validation error.
func GetPatchedProjectConfig(ctx context.Context, p *patch.Patch) (string, error) {
	if p.Version != "" {
		return "", errors.Errorf("patch '%s' already finalized", p.Version)
	}

	if p.ProjectStorageMethod != "" {
		// If the patch has been created but not finalized, it has either
		// already saved the project config as a string (if any) or stored
		// the parser project document (i.e. ProjectStorageMethod). Since
		// the project config is optional, the patch may be created and have
		// already evaluated the patched project config, but there simply
		// was none. Therefore, this is a valid way to check and get the
		// PatchedProjectConfig after the patch is already created.
		return p.PatchedProjectConfig, nil
	}

	// The patch has not been created yet, so do the first-time initialization
	// to get the parser project and project config.
	projectRef, opts, err := getLoadProjectOptsForPatch(ctx, p)
	if err != nil {
		return "", errors.Wrap(err, "fetching project options for patch")
	}

	if !projectRef.IsVersionControlEnabled() {
		return "", nil
	}

	projectFileBytes, err := getPatchedProjectYAML(ctx, projectRef, opts, p)
	if err != nil {
		return "", errors.Wrap(err, "getting patched project file as YAML")
	}

	return getProjectConfigYAML(p, projectFileBytes)
}

// getProjectConfigYAML creates a project config from the project YAML string
// and returns the project configuration as a string.
func getProjectConfigYAML(p *patch.Patch, projectFileBytes []byte) (string, error) {
	pc, err := CreateProjectConfig(projectFileBytes, p.Project)
	if err != nil {
		return "", errors.Wrap(err, "creating project config")
	}
	if pc == nil {
		return "", nil
	}

	yamlProjectConfig, err := yaml.Marshal(pc)
	if err != nil {
		return "", errors.Wrap(err, "marshalling project config into YAML")
	}
	return string(yamlProjectConfig), nil
}

// MakePatchedConfig takes in the project's remote file path containing the
// project YAML configuration and a stringified version of the project YAML
// configuration, and returns an unmarshalled version of the project with the
// patch applied.
func MakePatchedConfig(ctx context.Context, opts GetProjectOpts, projectConfig string) ([]byte, error) {
	p := opts.PatchOpts.patch
	remoteConfigPath := opts.RemotePath
	for _, patchPart := range p.Patches {
		// we only need to patch the main project and not any other modules
		if patchPart.ModuleName != "" {
			continue
		}

		var patchFilePath, localConfigPath, renamedFilePath, patchContents string
		var err error
		if patchPart.PatchSet.Patch == "" {
			patchContents, err = patch.FetchPatchContents(ctx, patchPart.PatchSet.PatchFileId)
			if err != nil {
				return nil, errors.Wrap(err, "fetching patch contents")
			}
			patchFilePath, err = util.WriteToTempFile(patchContents)
			if err != nil {
				return nil, errors.Wrap(err, "writing temporary patch file")
			}
		} else {
			patchContents = patchPart.PatchSet.Patch
			patchFilePath, err = util.WriteToTempFile(patchContents)
			if err != nil {
				return nil, errors.Wrap(err, "writing to temporary patch file")
			}
		}
		renamedFilePath = parseRenamedOrCopiedFile(patchContents, remoteConfigPath)

		defer os.Remove(patchFilePath) //nolint:evg-lint

		// clean the working directory
		workingDirectory := filepath.Join(filepath.Dir(patchFilePath), utility.RandomString())
		defer os.RemoveAll(workingDirectory)

		var renamedProjectConfig []byte
		// If this is a renamed include file, retrieve the bytes of the file before it
		// was renamed and use it as our local config.
		if renamedFilePath != "" {
			opts.RemotePath = renamedFilePath
			if projectConfig == "" {
				renamedProjectConfig, err = getFileForPatchDiff(ctx, opts)
				if err != nil {
					return nil, errors.Wrapf(err, "retrieving renamed file '%s'", renamedFilePath)
				}
				projectConfig = string(renamedProjectConfig)
			}
			localConfigPath = filepath.Join(
				workingDirectory,
				renamedFilePath,
			)
		} else {
			localConfigPath = filepath.Join(
				workingDirectory,
				remoteConfigPath,
			)
		}
		// write project configuration
		configFilePath, err := util.WriteToTempFile(projectConfig)
		if err != nil {
			return nil, errors.Wrap(err, "writing config file")
		}
		defer os.Remove(configFilePath) //nolint:evg-lint

		if err = os.MkdirAll(filepath.Dir(localConfigPath), 0755); err != nil {
			return nil, errors.WithStack(err)
		}
		// rename the temporary config file name to the remote config
		// file path if we are patching an existing remote config
		if renamedFilePath != "" || len(projectConfig) > 0 {
			if err = os.Rename(configFilePath, localConfigPath); err != nil {
				return nil, errors.Wrapf(err, "renaming file '%s' to '%s'", configFilePath, localConfigPath)
			}
			defer os.Remove(localConfigPath)
		}

		// selectively apply the patch to the config file
		patchCommandStrings := []string{
			"set -o errexit",
			fmt.Sprintf("git apply --whitespace=fix --include=%v < '%v'",
				remoteConfigPath, patchFilePath),
		}

		output := util.NewMBCappedWriter()
		err = opts.PatchOpts.env.JasperManager().CreateCommand(ctx).Add([]string{"bash", "-c", strings.Join(patchCommandStrings, "\n")}).
			Directory(workingDirectory).SetCombinedWriter(output).Run(ctx)
		if err != nil {
			grip.Error(message.WrapError(err, message.Fields{
				"message":       "error running patch command",
				"patch_id":      p.Id.Hex(),
				"output":        output.String(),
				"patch_command": patchCommandStrings,
			}))
			return nil, errors.Wrap(err, "running patch command (possibly due to merge conflict on evergreen configuration file)")
		}

		// Read in the patched config file. If the patched config file has been renamed as
		// part of the git apply, read from the new renamed location.
		patchedConfigPath := localConfigPath
		if renamedFilePath != "" {
			patchedConfigPath = filepath.Join(
				workingDirectory,
				remoteConfigPath,
			)
			defer os.Remove(patchedConfigPath)
		}
		data, err := os.ReadFile(patchedConfigPath)
		if err != nil {
			return nil, errors.Wrap(err, "reading patched config file")
		}
		if string(data) == "" {
			return nil, errors.New(EmptyConfigurationError)
		}
		return data, nil
	}
	return nil, errors.New("no patch on project")
}

// parseRenamedOrCopiedFile takes a patch contents and retrieves the old file name,
// if any, that the input filename was renamed or copied from.
func parseRenamedOrCopiedFile(patchContents, filename string) string {
	lines := strings.Split(patchContents, "\n")
	var renameFrom, renameTo, copyFrom, copyTo string
	isRenamed := false
	isCopied := false
	for _, line := range lines {
		if strings.HasPrefix(line, "rename from ") {
			renameFrom = strings.TrimPrefix(line, "rename from ")
		} else if strings.HasPrefix(line, "rename to ") {
			renameTo = strings.TrimPrefix(line, "rename to ")
			if renameTo == filename {
				isRenamed = true
				break
			}
		}

		if strings.HasPrefix(line, "copy from ") {
			copyFrom = strings.TrimPrefix(line, "copy from ")
		} else if strings.HasPrefix(line, "copy to ") {
			copyTo = strings.TrimPrefix(line, "copy to ")
			if copyTo == filename {
				isCopied = true
				break
			}
		}
	}
	if isRenamed {
		return renameFrom
	}
	if isCopied {
		return copyFrom
	}
	return ""
}

// FinalizePatch finalizes a patch:
// Patches a remote project's configuration file if needed.
// Creates a version for this patch and links it.
// Creates builds based on the Version
// Creates a manifest based on the Version
func FinalizePatch(ctx context.Context, p *patch.Patch, requester string) (*Version, error) {
	projectRef, err := FindMergedProjectRef(ctx, p.Project, p.Version, true)
	if err != nil {
		return nil, errors.Wrapf(err, "finding project '%s'", p.Project)
	}
	if projectRef == nil {
		return nil, errors.Errorf("project '%s' not found", p.Project)
	}

	settings, err := evergreen.GetConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting evergreen config")
	}
	project, _, err := FindAndTranslateProjectForPatch(ctx, settings, p)
	if err != nil {
		return nil, errors.Wrapf(err, "finding and translating project for patch '%s'", p.Id.Hex())
	}
	var config *ProjectConfig
	if projectRef.IsVersionControlEnabled() {
		config, err = CreateProjectConfig([]byte(p.PatchedProjectConfig), p.Project)
		if err != nil {
			return nil, errors.Wrapf(err,
				"marshalling patched project config from repository revision '%s'",
				p.Githash)
		}
	}
	if config != nil {
		config.Project = p.Project
		config.Id = p.Id.Hex()
	}

	distroAliases, err := distro.NewDistroAliasesLookupTable(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "resolving distro alias table for patch")
	}

	var parentPatchNumber int
	if p.IsChild() {
		parentPatch, err := p.SetParametersFromParent(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "getting parameters from parent patch")
		}
		parentPatchNumber = parentPatch.PatchNumber
	}

	params, err := getFullPatchParams(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "fetching patch parameters")
	}

	authorEmail := ""
	if p.GitInfo != nil {
		authorEmail = p.GitInfo.Email
	}
	if p.Author != "" {
		u, err := user.FindOneByIdContext(ctx, p.Author)
		if err != nil {
			return nil, errors.Wrap(err, "getting user")
		}
		if u != nil {
			authorEmail = u.Email()
		}
	}

	patchVersion := &Version{
		Id:                   p.Id.Hex(),
		CreateTime:           time.Now(),
		Identifier:           p.Project,
		Revision:             p.Githash,
		Author:               p.Author,
		Message:              p.Description,
		BuildIds:             []string{},
		BuildVariants:        []VersionBuildStatus{},
		Status:               evergreen.VersionCreated,
		Requester:            requester,
		ParentPatchID:        p.Triggers.ParentPatch,
		ParentPatchNumber:    parentPatchNumber,
		ProjectStorageMethod: p.ProjectStorageMethod,
		Branch:               projectRef.Branch,
		RevisionOrderNumber:  p.PatchNumber,
		AuthorID:             p.Author,
		Parameters:           params,
		Activated:            utility.TruePtr(),
		AuthorEmail:          authorEmail,
	}

	mfst, err := constructManifest(ctx, patchVersion, projectRef, project.Modules)
	if err != nil {
		return nil, errors.Wrap(err, "constructing manifest")
	}

	tasks := TaskVariantPairs{}
	if len(p.VariantsTasks) > 0 {
		tasks = VariantTasksToTVPairs(p.VariantsTasks)
	} else {
		// handle case where the patch is being finalized but only has the old schema tasks/variants
		// instead of the new one.
		for _, v := range p.BuildVariants {
			for _, t := range p.Tasks {
				if project.FindTaskForVariant(t, v) != nil {
					tasks.ExecTasks = append(tasks.ExecTasks, TVPair{v, t})
				}
			}
		}
		p.VariantsTasks = (&TaskVariantPairs{
			ExecTasks:    tasks.ExecTasks,
			DisplayTasks: tasks.DisplayTasks,
		}).TVPairsToVariantTasks()
	}

	// if variant tasks is still empty, then the patch is empty and we shouldn't finalize
	if len(p.VariantsTasks) == 0 {
		if p.IsMergeQueuePatch() {
			return nil, errors.Errorf("no builds or tasks for merge queue version in projects '%s', githash '%s'", p.Project, p.Githash)
		}
		return nil, errors.New("cannot finalize patch with no tasks")
	}
	taskIds, err := NewTaskIdConfig(project, patchVersion, tasks, projectRef.Identifier)
	if err != nil {
		return nil, errors.Wrap(err, "creating patch's task ID table")
	}
	variantsProcessed := map[string]bool{}

	creationInfo := TaskCreationInfo{
		Version:    patchVersion,
		Project:    project,
		ProjectRef: projectRef,
	}
	createTime, err := getTaskCreateTime(ctx, creationInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "getting create time for tasks in '%s', githash '%s'", p.Project, p.Githash)
	}

	buildsToInsert := build.Builds{}
	tasksToInsert := task.Tasks{}
	for _, vt := range p.VariantsTasks {
		if _, ok := variantsProcessed[vt.Variant]; ok {
			continue
		}
		variantsProcessed[vt.Variant] = true

		var displayNames []string
		for _, dt := range vt.DisplayTasks {
			displayNames = append(displayNames, dt.Name)
		}
		taskNames := tasks.ExecTasks.TaskNames(vt.Variant)

		buildCreationArgs := TaskCreationInfo{
			Project:          creationInfo.Project,
			ProjectRef:       creationInfo.ProjectRef,
			Version:          creationInfo.Version,
			TaskIDs:          taskIds,
			BuildVariantName: vt.Variant,
			ActivateBuild:    true,
			TaskNames:        taskNames,
			DisplayNames:     displayNames,
			DistroAliases:    distroAliases,
			TaskCreateTime:   createTime,
			// When a GitHub PR patch is finalized with the PR alias, all of the
			// tasks selected by the alias must finish in order for the
			// build/version to be finished.
			ActivatedTasksAreEssentialToSucceed: requester == evergreen.GithubPRRequester,
		}
		var build *build.Build
		var tasks task.Tasks
		build, tasks, err = CreateBuildFromVersionNoInsert(ctx, buildCreationArgs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if len(tasks) == 0 {
			grip.Info(message.Fields{
				"op":      "skipping empty build for patch version",
				"variant": vt.Variant,
				"version": patchVersion.Id,
			})
			continue
		}

		buildsToInsert = append(buildsToInsert, build)
		tasksToInsert = append(tasksToInsert, tasks...)

		patchVersion.BuildIds = append(patchVersion.BuildIds, build.Id)
		patchVersion.BuildVariants = append(patchVersion.BuildVariants,
			VersionBuildStatus{
				BuildVariant:     vt.Variant,
				BuildId:          build.Id,
				DisplayName:      build.DisplayName,
				ActivationStatus: ActivationStatus{Activated: true},
			},
		)
	}
	// We must set the NumDependents field for tasks prior to inserting them in the DB.
	SetNumDependents(tasksToInsert)
	numActivatedTasks := 0
	for _, t := range tasksToInsert {
		if t.Activated {
			numActivatedTasks++
			numActivatedTasks += utility.FromIntPtr(t.EstimatedNumActivatedGeneratedTasks)
		}
	}
	if err = task.UpdateSchedulingLimit(ctx, creationInfo.Version.AuthorID, creationInfo.Version.Requester, numActivatedTasks, true); err != nil {
		return nil, errors.Wrapf(err, "fetching user '%s' and updating their scheduling limit", creationInfo.Version.Author)
	}

	env := evergreen.GetEnvironment()
	mongoClient := env.Client()
	session, err := mongoClient.StartSession()
	if err != nil {
		return nil, errors.Wrap(err, "starting DB session")
	}
	defer session.EndSession(ctx)

	txFunc := func(ctx mongo.SessionContext) (any, error) {
		db := env.DB()
		_, err = db.Collection(VersionCollection).InsertOne(ctx, patchVersion)
		if err != nil {
			return nil, errors.Wrapf(err, "inserting version '%s'", patchVersion.Id)
		}
		if config != nil {
			_, err = db.Collection(ProjectConfigCollection).InsertOne(ctx, config)
			if err != nil {
				return nil, errors.Wrapf(err, "inserting project config for version '%s'", patchVersion.Id)
			}
		}
		if mfst != nil {
			if err = mfst.InsertWithContext(ctx); err != nil {
				return nil, errors.Wrapf(err, "inserting manifest for version '%s'", patchVersion.Id)
			}
		}
		if err = buildsToInsert.InsertMany(ctx, false); err != nil {
			return nil, errors.Wrapf(err, "inserting builds for version '%s'", patchVersion.Id)
		}
		if err = tasksToInsert.InsertUnordered(ctx); err != nil {
			return nil, errors.Wrapf(err, "inserting tasks for version '%s'", patchVersion.Id)
		}
		if err = p.SetFinalized(ctx, patchVersion.Id); err != nil {
			return nil, errors.Wrapf(err, "activating patch '%s'", patchVersion.Id)
		}
		return nil, err
	}

	_, err = session.WithTransaction(ctx, txFunc)
	if err != nil {
		return nil, errors.Wrap(err, "finalizing patch")
	}

	if p.IsParent() {
		// finalize child patches or subscribe on parent outcome based on parentStatus
		for _, childPatchId := range p.Triggers.ChildPatches {
			err = finalizeOrSubscribeChildPatch(ctx, childPatchId, p, requester)
			if err != nil {
				return nil, errors.Wrap(err, "finalizing child patch")
			}
		}
	}
	if len(tasksToInsert) > evergreen.NumTasksForLargePatch {
		numTasksActivated := 0
		for _, t := range tasksToInsert {
			if t.Activated {
				numTasksActivated++
			}
		}
		grip.Info(message.Fields{
			"message":             "version has large number of activated tasks",
			"op":                  "finalize patch",
			"num_tasks_activated": numTasksActivated,
			"total_tasks":         len(tasksToInsert),
			"version":             patchVersion.Id,
		})
	}

	return patchVersion, nil
}

// getFullPatchParams will retrieve a merged list of parameters defined on the patch alias (if any)
// with the parameters that were explicitly user-specified, with the latter taking precedence.
func getFullPatchParams(ctx context.Context, p *patch.Patch) ([]patch.Parameter, error) {
	paramsMap := map[string]string{}
	if p.Alias == "" || !IsPatchAlias(p.Alias) {
		return p.Parameters, nil
	}
	aliases, err := findAliasesForPatch(ctx, p.Project, p.Alias, p)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving alias '%s' for patch '%s'", p.Alias, p.Id.Hex())
	}
	for _, alias := range aliases {
		if len(alias.Parameters) > 0 {
			for _, aliasParam := range alias.Parameters {
				paramsMap[aliasParam.Key] = aliasParam.Value
			}
		}
	}
	for _, param := range p.Parameters {
		paramsMap[param.Key] = param.Value
	}
	var fullParams []patch.Parameter
	for k, v := range paramsMap {
		fullParams = append(fullParams, patch.Parameter{
			Key:   k,
			Value: v,
		})
	}
	return fullParams, nil
}

func getLoadProjectOptsForPatch(ctx context.Context, p *patch.Patch) (*ProjectRef, *GetProjectOpts, error) {
	projectRef, err := FindMergedProjectRef(ctx, p.Project, p.Version, true)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	if projectRef == nil {
		return nil, nil, errors.Errorf("project '%s' doesn't exist", p.Project)
	}
	hash := p.Githash
	if p.IsGithubPRPatch() {
		hash = p.GithubPatchData.HeadHash
	}
	if p.IsMergeQueuePatch() {
		hash = p.GithubMergeData.HeadSHA
	}

	var manifestID string
	if p.ReferenceManifestID != "" {
		manifestID = p.ReferenceManifestID
	} else {
		baseVersion, err := VersionFindOne(ctx, BaseVersionByProjectIdAndRevision(p.Project, p.Githash))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "finding base version for project '%s' and revision '%s'", p.Project, p.Githash)
		}
		if baseVersion != nil {
			manifestID = baseVersion.Id
		}
	}

	opts := GetProjectOpts{
		Ref:                 projectRef,
		ReadFileFrom:        ReadFromPatch,
		Revision:            hash,
		LocalModuleIncludes: p.LocalModuleIncludes,
		ReferencePatchID:    p.ReferenceManifestID,
		ReferenceManifestID: manifestID,
		PatchOpts: &PatchOpts{
			patch: p,
		},
	}
	return projectRef, &opts, nil
}

func finalizeOrSubscribeChildPatch(ctx context.Context, childPatchId string, parentPatch *patch.Patch, requester string) error {
	intent, err := patch.FindIntent(ctx, childPatchId, patch.TriggerIntentType)
	if err != nil {
		return errors.Wrap(err, "fetching child patch intent")
	}
	if intent == nil {
		return errors.New("child patch intent not found")
	}
	triggerIntent, ok := intent.(*patch.TriggerIntent)
	if !ok {
		return errors.Errorf("intent '%s' didn't not have expected type '%T'", intent.ID(), intent)
	}
	// if the parentStatus is "", finalize without waiting for the parent patch to finish running
	if triggerIntent.ParentStatus == "" {
		childPatchDoc, err := patch.FindOneId(ctx, childPatchId)
		if err != nil {
			return errors.Wrap(err, "fetching child patch")
		}
		if childPatchDoc == nil {
			return errors.Errorf("could not find child patch '%s'", childPatchId)
		}
		if _, err := FinalizePatch(ctx, childPatchDoc, requester); err != nil {
			grip.Error(message.WrapError(err, message.Fields{
				"message":       "Failed to finalize child patch document",
				"source":        requester,
				"patch_id":      childPatchId,
				"variants":      childPatchDoc.BuildVariants,
				"tasks":         childPatchDoc.Tasks,
				"variant_tasks": childPatchDoc.VariantsTasks,
				"alias":         childPatchDoc.Alias,
				"parentPatch":   childPatchDoc.Triggers.ParentPatch,
			}))
			return err
		}
	} else {
		//subscribe on parent outcome
		if err = SubscribeOnParentOutcome(ctx, triggerIntent.ParentStatus, childPatchId, parentPatch, requester); err != nil {
			return errors.Wrap(err, "getting parameters from parent patch")
		}
	}
	return nil
}

func SubscribeOnParentOutcome(ctx context.Context, parentStatus string, childPatchId string, parentPatch *patch.Patch, requester string) error {
	subscriber := event.NewRunChildPatchSubscriber(event.ChildPatchSubscriber{
		ParentStatus: parentStatus,
		ChildPatchId: childPatchId,
		Requester:    requester,
	})
	patchSub := event.NewParentPatchSubscription(parentPatch.Id.Hex(), subscriber)
	if err := patchSub.Upsert(ctx); err != nil {
		return errors.Wrapf(err, "inserting child patch subscription '%s'", childPatchId)
	}
	return nil
}

// CancelPatch aborts all of a patch's in-progress tasks and deactivates its undispatched tasks.
func CancelPatch(ctx context.Context, p *patch.Patch, reason task.AbortInfo) error {
	if p.Version != "" {
		if err := SetVersionActivation(ctx, p.Version, false, reason.User); err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(task.AbortVersionTasks(ctx, p.Version, reason))
	}

	return errors.WithStack(patch.Remove(ctx, patch.ById(p.Id)))
}

// AbortPatchesWithGithubPatchData aborts patches created
// before the given time, with the same PR number, and base repository. Tasks
// which are abortable will be aborted, while completed tasks will not be
// affected.
func AbortPatchesWithGithubPatchData(ctx context.Context, createdBefore time.Time, closed bool, newPatch, owner, repo string, prNumber int) error {
	patches, err := patch.Find(ctx, patch.ByGithubPRAndCreatedBefore(createdBefore, owner, repo, prNumber))
	if err != nil {
		return errors.Wrap(err, "fetching initial patch")
	}
	grip.Info(message.Fields{
		"source":         "github hook",
		"created_before": createdBefore.String(),
		"owner":          owner,
		"repo":           repo,
		"message":        "fetched patches to abort",
		"num_patches":    len(patches),
	})

	catcher := grip.NewSimpleCatcher()
	for _, p := range patches {
		if p.Version == "" { // Skip anything unfinalized
			continue
		}

		err = CancelPatch(ctx, &p, task.AbortInfo{User: evergreen.GithubPatchUser, NewVersion: newPatch, PRClosed: closed})
		msg := message.Fields{
			"source":         "github hook",
			"created_before": createdBefore.String(),
			"owner":          owner,
			"repo":           repo,
			"message":        "aborting patch's version",
			"patch_id":       p.Id.Hex(),
			"pr":             p.GithubPatchData.PRNumber,
			"project":        p.Project,
			"version":        p.Version,
		}
		grip.Error(message.WrapError(err, msg))
		catcher.Add(err)
	}

	return errors.Wrap(catcher.Resolve(), "aborting patches")
}
