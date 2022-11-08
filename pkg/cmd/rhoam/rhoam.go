package rhoam

import (
	"github.com/redhat-developer/app-services-cli/internal/doc"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/rhoam/generate"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

func NewRHOAMCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "rhoam",
		Annotations: map[string]string{doc.AnnotationName: "RHOAM commands"},
		Short:       f.Localizer.MustLocalize("rhoam.cmd.shortDescription"),
		Long:        f.Localizer.MustLocalize("rhoam.cmd.longDescription"),
		Example:     f.Localizer.MustLocalize("rhoam.cmd.example"),
		Args:        cobra.MinimumNArgs(1),
	}

	// add sub-commands
	cmd.AddCommand(
		generate.NewGenerateCommand(f),
	)

	return cmd
}
