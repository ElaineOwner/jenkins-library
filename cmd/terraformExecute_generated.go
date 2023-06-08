// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type terraformExecuteOptions struct {
	Command          string   `json:"command,omitempty"`
	TerraformSecrets string   `json:"terraformSecrets,omitempty"`
	GlobalOptions    []string `json:"globalOptions,omitempty"`
	AdditionalArgs   []string `json:"additionalArgs,omitempty"`
	Init             bool     `json:"init,omitempty"`
	CliConfigFile    string   `json:"cliConfigFile,omitempty"`
	Workspace        string   `json:"workspace,omitempty"`
}

type terraformExecuteCommonPipelineEnvironment struct {
	custom struct {
		terraformOutputs map[string]interface{}
	}
}

func (p *terraformExecuteCommonPipelineEnvironment) persist(path, resourceName string) {
	content := []struct {
		category string
		name     string
		value    interface{}
	}{
		{category: "custom", name: "terraformOutputs", value: p.custom.terraformOutputs},
	}

	errCount := 0
	for _, param := range content {
		err := piperenv.SetResourceParameter(path, resourceName, filepath.Join(param.category, param.name), param.value)
		if err != nil {
			log.Entry().WithError(err).Error("Error persisting piper environment.")
			errCount++
		}
	}
	if errCount > 0 {
		log.Entry().Error("failed to persist Piper environment")
	}
}

// TerraformExecuteCommand Executes Terraform
func TerraformExecuteCommand() *cobra.Command {
	const STEP_NAME = "terraformExecute"

	metadata := terraformExecuteMetadata()
	var stepConfig terraformExecuteOptions
	var startTime time.Time
	var commonPipelineEnvironment terraformExecuteCommonPipelineEnvironment
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createTerraformExecuteCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Executes Terraform",
		Long:  `This step executes the terraform binary with the given command, and is able to fetch additional variables from vault.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.CliConfigFile)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if err = log.RegisterANSHookIfConfigured(GeneralConfig.CorrelationID); err != nil {
				log.Entry().WithError(err).Warn("failed to set up SAP Alert Notification Service log hook")
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				commonPipelineEnvironment.persist(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {

					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.Dsn,
						GeneralConfig.HookConfig.SplunkConfig.Token,
						GeneralConfig.HookConfig.SplunkConfig.Index,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)

					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}

				if len(GeneralConfig.HookConfig.SplunkConfig.ProdDsn) > 0 {

					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.ProdDsn,
						GeneralConfig.HookConfig.SplunkConfig.ProdToken,
						GeneralConfig.HookConfig.SplunkConfig.ProdIndex,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)

					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)

			terraformExecute(stepConfig, &stepTelemetryData, &commonPipelineEnvironment)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addTerraformExecuteFlags(createTerraformExecuteCmd, &stepConfig)
	return createTerraformExecuteCmd
}

func addTerraformExecuteFlags(cmd *cobra.Command, stepConfig *terraformExecuteOptions) {
	cmd.Flags().StringVar(&stepConfig.Command, "command", `plan`, "")
	cmd.Flags().StringVar(&stepConfig.TerraformSecrets, "terraformSecrets", os.Getenv("PIPER_terraformSecrets"), "")
	cmd.Flags().StringSliceVar(&stepConfig.GlobalOptions, "globalOptions", []string{}, "")
	cmd.Flags().StringSliceVar(&stepConfig.AdditionalArgs, "additionalArgs", []string{}, "")
	cmd.Flags().BoolVar(&stepConfig.Init, "init", false, "")
	cmd.Flags().StringVar(&stepConfig.CliConfigFile, "cliConfigFile", os.Getenv("PIPER_cliConfigFile"), "Path to the terraform CLI configuration file (https://www.terraform.io/docs/cli/config/config-file.html#credentials).")
	cmd.Flags().StringVar(&stepConfig.Workspace, "workspace", os.Getenv("PIPER_workspace"), "")

}

// retrieve step metadata
func terraformExecuteMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "terraformExecute",
			Aliases:     []config.Alias{},
			Description: "Executes Terraform",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "cliConfigFileCredentialsId", Description: "Jenkins 'Secret file' credentials ID containing terraform CLI configuration. You can find more details about it in the [Terraform documentation](https://www.terraform.io/docs/cli/config/config-file.html#credentials).", Type: "jenkins"},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "command",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `plan`,
					},
					{
						Name: "terraformSecrets",
						ResourceRef: []config.ResourceReference{
							{
								Name:    "terraformFileVaultSecretName",
								Type:    "vaultSecretFile",
								Default: "terraform",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_terraformSecrets"),
					},
					{
						Name:        "globalOptions",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name:        "additionalArgs",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
					{
						Name:        "init",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name: "cliConfigFile",
						ResourceRef: []config.ResourceReference{
							{
								Name: "cliConfigFileCredentialsId",
								Type: "secret",
							},

							{
								Name:    "cliConfigFileVaultSecretName",
								Type:    "vaultSecretFile",
								Default: "terraform",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_cliConfigFile"),
					},
					{
						Name:        "workspace",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_workspace"),
					},
				},
			},
			Containers: []config.Container{
				{Name: "terraform", Image: "hashicorp/terraform:1.0.10", EnvVars: []config.EnvVar{{Name: "TF_IN_AUTOMATION", Value: "piper"}}, Options: []config.Option{{Name: "--entrypoint", Value: ""}}},
			},
			Outputs: config.StepOutputs{
				Resources: []config.StepResources{
					{
						Name: "commonPipelineEnvironment",
						Type: "piperEnvironment",
						Parameters: []map[string]interface{}{
							{"name": "custom/terraformOutputs", "type": "map[string]interface{}"},
						},
					},
				},
			},
		},
	}
	return theMetaData
}
