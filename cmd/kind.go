package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	du "io.twasyl/devcore/pkg/utils"
)

func init() {
	rootCmd.AddCommand(kindCmd)

	kindCmd.AddCommand(kindClusterCmd)
	kindClusterCmd.PersistentFlags().StringVarP(&kindClusterName, "name", "n", "", "The name of cluster")

	kindClusterCmd.AddCommand(createKindClusterCmd)
	createKindClusterCmd.MarkFlagRequired("name")
	createKindClusterCmd.Flags().StringVarP(&k8sDashboardVersion, "dashboard-version", "v", "v2.7.0", "The version of the Kubernetes dashboard to install")
	createKindClusterCmd.Flags().StringVarP(&kindClusterConfigFile, "config", "c", "", "The config file to use to create the cluster")
	createKindClusterCmd.Flags().BoolVar(&installNginx, "nginx", false, "Specify wether to install nginx or not")
	createKindClusterCmd.Flags().StringVar(&k8sVersion, "k8s-version", "", "The Kubernetes version to install")
	deleteKindClusterCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Indicates to delete all kind clusters")
	kindClusterCmd.AddCommand(deleteKindClusterCmd)

	kindCmd.AddCommand(kindTokenCmd)
}

var kindClusterName string
var kindClusterConfigFile string
var installNginx bool
var k8sDashboardVersion string
var allFlag bool
var k8sVersion string

var kindCmd = &cobra.Command{
	Use:   "kind",
	Short: "Provide kind utilities",
	Long:  "kind command provide utilities related to kind such as installing, creating and deleting a cluster",
}

var kindClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster operations",
}

var createKindClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a kind cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := createKindCluster(); err != nil {
			return err
		}

		if installNginx {
			if err := installNginxInCluster(); err != nil {
				return err
			}
		}

		return installK8sDashboard()
	},
}

var deleteKindClusterCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a kind cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if allFlag {
			if clusters, err := du.ExecCommandAndGetSync("kind", "get", "clusters"); err != nil {
				return err
			} else if !strings.Contains(clusters, "No kind clusters found.") {
				for _, cluster := range strings.Fields(clusters) {
					if err = du.ExecCommandSync("kind", "delete", "cluster", "--name", cluster); err != nil {
						return err
					}
				}
			}
		} else {
			return du.ExecCommandSync("kind", "delete", "cluster", "--name", kindClusterName)
		}
		return nil
	},
}

var kindTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get an authentication token to be used in the K8s dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := du.ExecCommandAndGetSync("kubectl", "-n", "kubernetes-dashboard", "get", "sa/admin-user", "-o", "jsonpath=\"{.secrets[0].name}\"")
		if err != nil {
			fmt.Print(output)
			return err
		}
		secret := strings.ReplaceAll(output, "\"", "")
		output, err = du.ExecCommandAndGetSync("kubectl", "-n", "kubernetes-dashboard", "get", "secret", secret, "-o", "go-template=\"{{.data.token | base64decode}}\"")
		if err != nil {
			fmt.Print(output)
			return err
		}

		du.ToClipboard([]byte(strings.ReplaceAll(output, "\"", "")))
		fmt.Println("Token copied to your clipboard.")
		du.OpenBrowser("http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/")
		return nil
	},
}

func createKindCluster() error {
	cmdArgs := []string{"create", "cluster", "--name", kindClusterName}
	if kindClusterConfigFile != "" {
		cmdArgs = append(cmdArgs, "--config", kindClusterConfigFile)
	}

	if k8sVersion != "" {
		cmdArgs = append(cmdArgs, "--image", fmt.Sprintf("kindest/node:v%s", k8sVersion))
	}

	return du.ExecCommandSync("kind", cmdArgs...)
}

func installNginxInCluster() error {
	return du.ExecCommandsSync([][]string{
		{"kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml"},
		{"kubectl", "wait", "--namespace", "ingress-nginx", "--for=condition=ready", "pod", "--selector=app.kubernetes.io/component=controller",
			"--timeout=300s"},
		{"sleep", "15s"},
		{"kubectl", "delete", "-A", "ValidatingWebhookConfiguration", "ingress-nginx-admission"},
	})
}

func installK8sDashboard() error {
	err := du.ExecCommandSync("kubectl", "apply", "-f", fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/dashboard/%s/aio/deploy/recommended.yaml", k8sDashboardVersion))
	if err != nil {
		return err
	}

	configFile := "./tmp.yaml"
	configFileContent := []byte(`apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kubernetes-dashboard`)
	err = os.WriteFile(configFile, configFileContent, 0744)
	if err != nil {
		return err
	}

	err = du.ExecCommandSync("kubectl", "apply", "-f", configFile)
	os.Remove(configFile)

	return err
}
