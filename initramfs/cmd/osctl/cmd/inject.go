package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/autonomy/dianemo/initramfs/pkg/userdata"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// injectCmd represents the inject command
var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Inject data into fields in the user data.",
	Long:  ``,
}

// injectOSCmd represents the gen inject os command
var injectOSCmd = &cobra.Command{
	Use:   "os",
	Short: "Populates fields in the user data that are generated for the OS",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if len(args) != 1 {
			os.Exit(1)
		}
		filename := args[0]
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			os.Exit(1)
		}
		data := &userdata.UserData{}
		if err = yaml.Unmarshal(fileBytes, data); err != nil {
			os.Exit(1)
		}
		if data.OS.Security == nil {
			data.OS.Security = &userdata.Security{}
			data.OS.Security.Identity = &userdata.PEMEncodedCertificateAndKey{}
			data.OS.Security.CA = &userdata.PEMEncodedCertificateAndKey{}
		}

		encoded := &bytes.Buffer{}
		encoder := base64.NewEncoder(base64.StdEncoding, encoded)
		// nolint: errcheck
		defer encoder.Close()
		if identity != "" {
			fileBytes, err = ioutil.ReadFile(identity + ".crt")
			if err != nil {
				os.Exit(1)
			}
			if _, err = encoder.Write(fileBytes); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			data.OS.Security.Identity.Crt = encoded.Bytes()
			encoded.Reset()

			fileBytes, err = ioutil.ReadFile(identity + ".key")
			if err != nil {
				os.Exit(1)
			}
			if _, err = encoder.Write(fileBytes); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			data.OS.Security.Identity.Key = encoded.Bytes()
			encoded.Reset()
		}
		if ca != "" {
			fileBytes, err = ioutil.ReadFile(ca + ".crt")
			if err != nil {
				os.Exit(1)
			}
			if _, err = encoder.Write(fileBytes); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			data.OS.Security.CA.Crt = encoded.Bytes()
			encoded.Reset()
		}

		dataBytes, err := yaml.Marshal(data)
		if err != nil {
			os.Exit(1)
		}
		if err := ioutil.WriteFile(filename, dataBytes, 0700); err != nil {
			os.Exit(1)
		}
	},
}

// injectKubernetesCmd represents the gen inject kubernetes command
var injectKubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Populates fields in the user data that are generated for Kubernetes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			os.Exit(1)
		}
		filename := args[0]
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data := &userdata.UserData{}
		if err = yaml.Unmarshal(fileBytes, data); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if data.Kubernetes.CA == nil {
			data.Kubernetes.CA = &userdata.PEMEncodedCertificateAndKey{}
		}
		if ca != "" {
			encoded := &bytes.Buffer{}
			encoder := base64.NewEncoder(base64.StdEncoding, encoded)
			fileBytes, err = ioutil.ReadFile(ca + ".crt")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if _, err = encoder.Write(fileBytes); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			data.Kubernetes.CA.Crt = encoded.Bytes()
			encoded.Reset()

			fileBytes, err = ioutil.ReadFile(ca + ".key")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if _, err = encoder.Write(fileBytes); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			data.Kubernetes.CA.Key = encoded.Bytes()
			encoded.Reset()
		}

		dataBytes, err := yaml.Marshal(data)
		if err != nil {
			os.Exit(1)
		}
		if err := ioutil.WriteFile(filename, dataBytes, 0700); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	// Inject OS
	injectOSCmd.Flags().StringVar(&ca, "ca", "", "the basename of the key pair to use as the CA")
	injectOSCmd.Flags().StringVar(&identity, "identity", "", "the basename of the key pair to use as the identity")
	// Inject Kubernetes
	injectKubernetesCmd.Flags().StringVar(&ca, "ca", "", "the basename of the key pair to use as the CA")
	injectKubernetesCmd.Flags().StringVar(&hash, "hash", "", "the basename of the CA to use as the hash")

	injectCmd.AddCommand(injectOSCmd, injectKubernetesCmd)
	rootCmd.AddCommand(injectCmd)
}
