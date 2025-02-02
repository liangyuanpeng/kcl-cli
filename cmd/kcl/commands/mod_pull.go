package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
)

const (
	modPullDesc = `This command pulls kcl modules from the registry.
`
	modPullExample = `  # Pull the the module named "k8s" to the local path from the registry
  kcl mod pull k8s

  # Pull the module dependency named "k8s" with the version "1.28"
  kcl mod pull k8s:1.28

  # Pull the module from the GitHub by git url
  kcl mod pull git://github.com/kcl-lang/konfig --tag v0.4.0

  # Pull the module from the OCI Registry by oci url
  kcl mod pull oci://ghcr.io/kcl-lang/helloworld --tag 0.1.0

  # Pull the module from the Git by flag
  kcl mod pull --git https://github.com/kcl-lang/konfig --tag v0.4.0

  # Pull the module from the Git by flag with ssh url
  kcl mod pull --git ssh://github.com/kcl-lang/konfig --tag v0.4.0

  # Pull the module from the OCI Registry by flag
  kcl mod pull --oci https://ghcr.io/kcl-lang/helloworld --tag 0.1.0
  
  # Pull the module from the OCI Registry by flag and specify the module spce
  kcl mod pull subhelloworld --oci https://ghcr.io/kcl-lang/helloworld --tag 0.1.4
  
  # Pull the module from the OCI Registry by flag and specify the module spce with version
  kcl mod pull subhelloworld:0.0.1 --oci https://ghcr.io/kcl-lang/helloworld --tag 0.1.4
  
  # Pull the module from the Git Repo by flag and specify the module spce
  kcl mod pull cc --git git://github.com/kcl-lang/flask-demo-kcl-manifests.git --commit 8308200
  
  # Pull the module from the Git Repo by flag and specify the module spce with version
  kcl mod pull cc:0.0.1 --git git://github.com/kcl-lang/flask-demo-kcl-manifests.git --commit 8308200`
)

// NewModPullCmd returns the mod pull command.
func NewModPullCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull",
		Short:   "pull kcl package from the registry",
		Long:    modPullDesc,
		Example: modPullExample,
		RunE: func(_ *cobra.Command, args []string) error {
			localPath := argsGet(args, 1)
			return pull(cli, args, localPath)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&git, "git", "", "git repository url")
	cmd.Flags().StringVar(&oci, "oci", "", "oci repository url")
	cmd.Flags().StringVar(&tag, "tag", "", "git or oci repository tag")
	cmd.Flags().StringVar(&commit, "commit", "", "git repository commit")
	cmd.Flags().StringVar(&branch, "branch", "", "git repository branch")
	cmd.Flags().BoolVar(&insecureSkipTLSverify, "insecure-skip-tls-verify", false, "skip tls certificate checks for the KCL module download")

	return cmd
}

func pull(cli *client.KpmClient, args []string, localPath string) error {
	source, err := ParseSourceFromArgs(cli, args)
	if err != nil {
		return err
	}

	cli.SetInsecureSkipTLSverify(insecureSkipTLSverify)
	_, err = cli.Pull(
		client.WithPullSource(source),
		client.WithLocalPath(localPath),
	)

	if err != nil {
		return err
	}

	return nil
}
