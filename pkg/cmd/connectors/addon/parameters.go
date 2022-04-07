package addon

import (
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
	"github.com/redhat-developer/app-services-cli/pkg/shared/connection"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"

	"github.com/spf13/cobra"
)

type options struct {
	outputFormat string
	id           string

	f *factory.Factory
}

func NewParametersCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:     "addon-parameters",
		Short:   f.Localizer.MustLocalize("connector.addon.params.cmd.shortDescription"),
		Long:    f.Localizer.MustLocalize("connector.addon.params.cmd.longDescription"),
		Example: f.Localizer.MustLocalize("connector.addon.params.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, flagutil.ValidOutputFormats...) {
				return flagutil.InvalidValueError("output", opts.outputFormat, flagutil.ValidOutputFormats...)
			}

			if opts.id == "" {
				// TODO validation for id?
			}

			return runList(opts)
		},
	}

	flags := flagutil.NewFlagSet(cmd, f.Localizer)

	flags.AddOutput(&opts.outputFormat)
	flags.StringVar(&opts.id, "id", "", f.Localizer.MustLocalize("connectors.common.id.flag.page"))

	return cmd
}

func runList(opts *options) error {
	f := opts.f
	conn, err := f.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	api := conn.API()

	a := api.ConnectorsMgmt().ConnectorClustersApi.GetConnectorClusterAddonParameters(f.Context, opts.id)

	response, _, err := a.Execute()
	if err != nil {
		return err
	}

	return dump.Formatted(f.IOStreams.Out, opts.outputFormat, response)

}
