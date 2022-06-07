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
	createKindClusterCmd.Flags().StringVarP(&k8sDashboardVersion, "dashboard-version", "v", "v2.3.1", "The version of the Kubernetes dashboard to install")

	kindClusterCmd.AddCommand(deleteKindClusterCmd)

	kindCmd.AddCommand(kindTokenCmd)
}

var kindClusterName string
var k8sDashboardVersion string

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
		return installK8sDashboard()
	},
}

var deleteKindClusterCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a kind cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := du.ExecCommand("kind", "delete", "cluster", "--name", kindClusterName)
		if err != nil {
			fmt.Println(output)
		}
		return err
	},
}

var kindTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get an authentication token to be used in the K8s dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := du.ExecCommand("kubectl", "-n", "kubernetes-dashboard", "get", "sa/admin-user", "-o", "jsonpath=\"{.secrets[0].name}\"")
		if err != nil {
			fmt.Print(output)
			return err
		}
		secret := strings.ReplaceAll(output, "\"", "")
		output, err = du.ExecCommand("kubectl", "-n", "kubernetes-dashboard", "get", "secret", secret, "-o", "go-template=\"{{.data.token | base64decode}}\"")
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
	output, err := du.ExecCommand("kind", "create", "cluster", "--name", kindClusterName)
	fmt.Print(output)
	return err
}

func installK8sDashboard() error {
	_, err := du.ExecCommand("kubectl", "apply", "-f", fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/dashboard/%s/aio/deploy/recommended.yaml", k8sDashboardVersion))
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

	_, err = du.ExecCommand("kubectl", "apply", "-f", configFile)
	os.Remove(configFile)

	return err
}
