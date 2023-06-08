// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type apiProviderUploadOptions struct {
	APIServiceKey string `json:"apiServiceKey,omitempty"`
	FilePath      string `json:"filePath,omitempty"`
}

// ApiProviderUploadCommand this steps creates an API provider artifact in the API Portal.
func ApiProviderUploadCommand() *cobra.Command {
	const STEP_NAME = "apiProviderUpload"

	metadata := apiProviderUploadMetadata()
	var stepConfig apiProviderUploadOptions
	var startTime time.Time
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createApiProviderUploadCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "this steps creates an API provider artifact in the API Portal.",
		Long: `This steps creates an API provider artifact in the API Portal using the OData API.
Learn more about API Management api for creating an API provider artifact [here](https://help.sap.com/viewer/66d066d903c2473f81ec33acfe2ccdb4/Cloud/en-US/e26b3320cd534ae4bc743af8013a8abb.html).`,
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
			log.RegisterSecret(stepConfig.APIServiceKey)

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

			apiProviderUpload(stepConfig, &stepTelemetryData)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addApiProviderUploadFlags(createApiProviderUploadCmd, &stepConfig)
	return createApiProviderUploadCmd
}

func addApiProviderUploadFlags(cmd *cobra.Command, stepConfig *apiProviderUploadOptions) {
	cmd.Flags().StringVar(&stepConfig.APIServiceKey, "apiServiceKey", os.Getenv("PIPER_apiServiceKey"), "Service key JSON string to access the API Management Runtime service instance of plan 'api'")
	cmd.Flags().StringVar(&stepConfig.FilePath, "filePath", os.Getenv("PIPER_filePath"), "Specifies api provider json file relative path")

	cmd.MarkFlagRequired("apiServiceKey")
	cmd.MarkFlagRequired("filePath")
}

// retrieve step metadata
func apiProviderUploadMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "apiProviderUpload",
			Aliases:     []config.Alias{},
			Description: "this steps creates an API provider artifact in the API Portal.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "apimApiServiceKeyCredentialsId", Description: "Jenkins secret text credential ID containing the service key to the API Management Runtime service instance of plan 'api'", Type: "jenkins"},
				},
				Parameters: []config.StepParameters{
					{
						Name: "apiServiceKey",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "apimApiServiceKeyCredentialsId",
								Param: "apiServiceKey",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_apiServiceKey"),
					},
					{
						Name:        "filePath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_filePath"),
					},
				},
			},
		},
	}
	return theMetaData
}
