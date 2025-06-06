package operations

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/agent"
	"github.com/evergreen-ci/evergreen/agent/command"
	"github.com/evergreen-ci/evergreen/agent/globals"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/grip/recovery"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const defaultAgentStatusPort = 2285

const (
	agentAPIServerURLFlagName  = "api_server"
	agentCloudProviderFlagName = "provider"
	agentHostIDFlagName        = "host_id"
	agentHostSecretFlagName    = "host_secret"
	singleTaskDistroFlagName   = "single_task_distro"
)

func Agent() cli.Command {
	const (
		workingDirectoryFlagName           = "working_directory"
		logOutputFlagName                  = "log_output"
		logPrefixFlagName                  = "log_prefix"
		statusPortFlagName                 = "status_port"
		cleanupFlagName                    = "cleanup"
		modeFlagName                       = "mode"
		podIDFlagName                      = "pod_id"
		podSecretFlagName                  = "pod_secret"
		versionFlagName                    = "version"
		sendTaskLogsToGlobalSenderFlagName = "global_task_logs"
	)

	return cli.Command{
		Name:  "agent",
		Usage: "run an evergreen agent",
		Subcommands: []cli.Command{
			agentMonitor(),
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   agentHostIDFlagName,
				Usage:  "the ID of the host the agent is running on (applies only to host mode)",
				EnvVar: evergreen.HostIDEnvVar,
			},
			cli.StringFlag{
				Name:   agentHostSecretFlagName,
				Usage:  "secret for the current host (applies only to host mode)",
				EnvVar: evergreen.HostSecretEnvVar,
			},
			cli.StringFlag{
				Name:   podIDFlagName,
				Usage:  "the ID of the pod that the agent is running on (applies only to pod mode)",
				EnvVar: "POD_ID",
			},
			cli.StringFlag{
				Name:   podSecretFlagName,
				Usage:  "the secret for the pod that the agent is running on (applies only to pod mode)",
				EnvVar: "POD_SECRET",
			},
			cli.StringFlag{
				Name:  agentAPIServerURLFlagName,
				Usage: "URL of the API server",
			},
			cli.StringFlag{
				Name:  workingDirectoryFlagName,
				Usage: "working directory for the agent",
			},
			cli.StringFlag{
				Name:  logOutputFlagName,
				Value: string(globals.LogOutputFile),
				Usage: "location for the agent's log output (file, stdout)",
			},
			cli.StringFlag{
				Name:  logPrefixFlagName,
				Value: "evg.agent",
				Usage: "if logging to a file, the prefix of the file name for the agent's log output",
			},
			cli.IntFlag{
				Name:  statusPortFlagName,
				Value: defaultAgentStatusPort,
				Usage: "port to run the status server",
			},
			cli.BoolFlag{
				Name:  cleanupFlagName,
				Usage: "clean up working directory and processes (do not set for smoke tests)",
			},
			cli.StringFlag{
				Name:  agentCloudProviderFlagName,
				Usage: "the cloud provider that manages this host",
			},
			cli.StringFlag{
				Name:  modeFlagName,
				Usage: "the mode that the agent should run in (host, pod)",
				Value: "host",
			},
			cli.BoolFlag{
				Name:  sendTaskLogsToGlobalSenderFlagName,
				Usage: "send task logs to the global agent file log",
			},
			cli.BoolFlag{
				Name:  joinFlagNames(versionFlagName, "v"),
				Usage: "print the agent revision of the current binary and exit",
			},
			cli.BoolFlag{
				Name:  singleTaskDistroFlagName,
				Usage: "marks the agent as running in single task distro",
			},
		},
		Before: mergeBeforeFuncs(
			func(c *cli.Context) error {
				if c.Bool(versionFlagName) {
					return nil
				}

				catcher := grip.NewBasicCatcher()
				catcher.Add(requireStringFlag(agentAPIServerURLFlagName)(c))
				catcher.Add(requireStringFlag(workingDirectoryFlagName)(c))
				mode := c.String(modeFlagName)
				switch mode {
				case string(globals.HostMode):
					catcher.Add(requireStringFlag(agentHostIDFlagName)(c))
					catcher.Add(requireStringFlag(agentHostSecretFlagName)(c))
				case string(globals.PodMode):
					catcher.Add(requireStringFlag(podIDFlagName)(c))
					catcher.Add(requireStringFlag(podSecretFlagName)(c))
				default:
					return errors.Errorf("invalid mode '%s'", mode)
				}
				return catcher.Resolve()
			},
			func(c *cli.Context) error {
				grip.SetName("evergreen.agent")
				return nil
			},
		),
		Action: func(c *cli.Context) error {
			if c.Bool(versionFlagName) {
				fmt.Println(evergreen.AgentVersion)
				return nil
			}

			opts := agent.Options{
				HostID:                     c.String(agentHostIDFlagName),
				HostSecret:                 c.String(agentHostSecretFlagName),
				PodID:                      c.String(podIDFlagName),
				PodSecret:                  c.String(podSecretFlagName),
				Mode:                       globals.Mode(c.String(modeFlagName)),
				StatusPort:                 c.Int(statusPortFlagName),
				LogPrefix:                  c.String(logPrefixFlagName),
				LogOutput:                  globals.LogOutputType(c.String(logOutputFlagName)),
				WorkingDirectory:           c.String(workingDirectoryFlagName),
				Cleanup:                    c.Bool(cleanupFlagName),
				CloudProvider:              c.String(agentCloudProviderFlagName),
				SendTaskLogsToGlobalSender: c.Bool(sendTaskLogsToGlobalSenderFlagName),
				SingleTaskDistro:           c.Bool(singleTaskDistroFlagName),
			}

			// Once the agent has retrieved the host ID and secret, unset those
			// env vars to prevent them from being inherited by task
			// subprocesses (e.g. shell.exec).
			if err := os.Unsetenv(evergreen.HostIDEnvVar); err != nil {
				return errors.Wrapf(err, "unsetting %s env var", evergreen.HostIDEnvVar)
			}
			if err := os.Unsetenv(evergreen.HostSecretEnvVar); err != nil {
				return errors.Wrapf(err, "unsetting %s env var", evergreen.HostSecretEnvVar)
			}

			if err := os.MkdirAll(opts.WorkingDirectory, 0777); err != nil {
				return errors.Wrapf(err, "creating working directory '%s'", opts.WorkingDirectory)
			}

			grip.Info(message.Fields{
				"message":            "starting agent",
				"commands":           command.RegisteredCommandNames(),
				"dir":                opts.WorkingDirectory,
				"host_id":            opts.HostID,
				"single_task_distro": opts.SingleTaskDistro,
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			agt, err := agent.New(ctx, opts, c.String(agentAPIServerURLFlagName))
			if err != nil {
				return errors.Wrap(err, "constructing agent")
			}

			go hardShutdownForSignals(ctx, cancel, agt.Close)

			defer agt.Close(ctx)

			sender, err := agt.GetSender(ctx, opts.LogOutput, opts.LogPrefix, "", -1)
			if err != nil {
				return errors.Wrap(err, "configuring logger")
			}

			if err = grip.SetSender(sender); err != nil {
				return errors.Wrap(err, "setting up global logger")
			}
			agt.SetDefaultLogger(sender)
			agt.SetHomeDirectory()

			err = agt.Start(ctx)
			if err != nil {
				msg := message.Fields{
					"message": "agent is exiting due to unrecoverable error",
				}
				msg = opts.AddLoggableInfo(msg)
				// Although we still want to return an error, it's an acceptable state for an agent to be unauthorized
				// as ec2 doesn't immediately shut down hosts, so avoid logging it as an emergency.
				isUnauthorizedErr := strings.Contains(err.Error(), "401 (Unauthorized)")
				grip.EmergencyWhen(!isUnauthorizedErr, message.WrapError(err, msg))
				return err
			}

			return nil
		},
	}
}

func hardShutdownForSignals(ctx context.Context, serviceCanceler context.CancelFunc, closeAgent func(context.Context)) {
	defer recovery.LogStackTraceAndExit("agent signal handler")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-sigChan:
		grip.Info("service exiting after receiving signal")
	}

	// Close may not succeed if the context is cancelled, but this is a
	// best-effort attempt to clean up before imminent shutdown anyways.
	closeAgent(ctx)

	serviceCanceler()
	os.Exit(2)
}
