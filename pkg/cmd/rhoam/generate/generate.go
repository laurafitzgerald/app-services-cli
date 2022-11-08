package generate

import (
	"fmt"
	v1beta1 "github.com/redhat-developer/app-services-cli/apis/3scale/v1beta1"

	kafkaFlagutil "github.com/redhat-developer/app-services-cli/pkg/cmd/kafka/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

type options struct {
	name       string
	url        string
	systemName string
	// apicast environment {staging|production}
	env string
	f   *factory.Factory
}

func NewGenerateCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		f: f,
	}

	cmd := &cobra.Command{
		Use:     "generate",
		Short:   f.Localizer.MustLocalize("rhoam.generate.cmd.shortDescription"),
		Long:    f.Localizer.MustLocalize("rhoam.generate.cmd.longDescription"),
		Example: f.Localizer.MustLocalize("rhoam.generate.cmd.example"),
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stdout, "Check the test directory for the generated files")
			err := generate(opts)
			return err
		},
	}
	flags := kafkaFlagutil.NewFlagSet(cmd, f.Localizer)

	flags.StringVar(&opts.name, "name", "", f.Localizer.MustLocalize("rhoam.generate.flag.name.description"))
	flags.StringVar(&opts.url, "url", "", f.Localizer.MustLocalize("rhoam.generate.flag.url.description"))
	flags.StringVar(&opts.systemName, "systemName", "", f.Localizer.MustLocalize("rhoam.generate.flag.systemName.description"))
	flags.StringVar(&opts.systemName, "env", "", f.Localizer.MustLocalize("rhoam.generate.flag.env.description"))

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("url")

	return cmd
}

func generate(opts *options) error {

	backend := &v1beta1.Backend{
		TypeMeta: v1.TypeMeta{
			Kind:       v1beta1.BackendKind,
			APIVersion: fmt.Sprintf("%s/%s", v1beta1.GroupVersion.Group, v1beta1.GroupVersion.Version),
		},
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s-backend", opts.name),
			Annotations: map[string]string{
				"git.ops/managed": "true",
			},
			Finalizers: []string{
				fmt.Sprintf("backend.%s/finalizer", v1beta1.GroupVersion.Group),
			},
		},
		Spec: v1beta1.BackendSpec{
			Name:           opts.name,
			PrivateBaseURL: opts.url,
			SystemName:     fmt.Sprintf("%s-system", opts.name),
		},
	}

	outputWriter, err := Open("./test", "backend.yaml")
	if err != nil {
		return err
	}
	err = WriteYAML(outputWriter, backend)
	if err != nil {
		return err
	}

	product := v1beta1.Product{
		TypeMeta: v1.TypeMeta{
			Kind:       v1beta1.ProductKind,
			APIVersion: fmt.Sprintf("%s/%s", v1beta1.GroupVersion.Group, v1beta1.GroupVersion.Version),
		},
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s-product", opts.name),
			Annotations: map[string]string{
				"git.ops/managed": "true",
			},
			Finalizers: []string{
				fmt.Sprintf("product.%s/finalizer", v1beta1.GroupVersion.Group),
			},
		},
		Spec: v1beta1.ProductSpec{
			Name:       opts.name,
			SystemName: fmt.Sprintf("%s-system-default-apicast-%s", opts.name, opts.env),
			BackendUsages: map[string]v1beta1.BackendUsageSpec{
				// TODO verify this is the correct way to associate product to the backend
				fmt.Sprintf("%s-system-default-apicast-%s", opts.name, opts.env): {Path: "/"},
			},
			MappingRules: []v1beta1.MappingRuleSpec{},
		},
		Status: v1beta1.ProductStatus{},
	}

	outputWriter, err = Open("./test", "product.yaml")
	if err != nil {
		return err
	}
	err = WriteYAML(outputWriter, product)
	if err != nil {
		return err
	}
	return nil
}
