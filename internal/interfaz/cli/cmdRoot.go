package cli

import (
	"fmt"
	"os"
	"sync"

	app "github.com/jvillablanca/rolecheck/internal/aplicacion"
	cfg "github.com/jvillablanca/rolecheck/internal/infraestructura/puertos/cfg"

	"github.com/spf13/cobra"
)

var (
	rootCmd         *cobra.Command
	once            sync.Once
	username1       string
	host1           string
	username2       string
	host2           string
	userAdmin       string
	passAdmin       string
	nombreAmbiente1 string
	nombreAmbiente2 string
	c               = cfg.Crea // Usar el creador de configuración

)

func validateUserHostArgs() {
	// Se valida que venga al menos un usuario y dos hosts o dos usuarios y un host
	userCount := 0
	if username1 != "" {
		userCount++
	}
	if username2 != "" {
		userCount++
	}

	hostCount := 0
	if host1 != "" {
		hostCount++
	}
	if host2 != "" {
		hostCount++
	}

	if !((userCount >= 1 && hostCount >= 2) || (userCount >= 2 && hostCount >= 1)) {
		fmt.Println("Error: debes proporcionar al menos 1 usuario y 2 hosts, o 2 usuarios y 1 host.")
		os.Exit(1)
	}
}

func initRootCmd() {
	rootCmd = &cobra.Command{
		Use:   "rolecheck",
		Short: "rolecheck CLI",
		Long:  "rolecheck es una herramienta para comparar roles y permisos en Oracle Fusion ERP.",
		Run: func(cmd *cobra.Command, args []string) {

			validateUserHostArgs()
			// Si username2 o host2 están vacíos, asignar el valor de username1 o host1
			if username2 == "" {
				username2 = username1
			}
			if host2 == "" {
				host2 = host1
			}
			c.IniCfg(username1, host1, username2, host2, userAdmin, passAdmin, nombreAmbiente1, nombreAmbiente2)
			app.ComparaCuentas()
			fmt.Println("Comparación de cuentas completada.")
		},
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initPackage() {
	once.Do(func() {
		initRootCmd()
		rootCmd.Flags().StringVar(&username1, "username1", "", "Usuario a analizar entre hosts o primer usuario a comparar")
		rootCmd.Flags().StringVar(&username2, "username2", "", "Segundo usuario a comparar")
		rootCmd.Flags().StringVar(&host1, "host1", "", "Host de análisis o primer host")
		rootCmd.Flags().StringVar(&host2, "host2", "", "Segundo host")
		rootCmd.Flags().StringVar(&userAdmin, "userAdmin", "", "Usuario administrador para operaciones de administración")
		rootCmd.Flags().StringVar(&passAdmin, "passAdmin", "", "Contraseña del usuario administrador")
		rootCmd.Flags().StringVar(&nombreAmbiente1, "nombreAmbiente1", "ambiente1", "Nombre del primer ambiente")
		rootCmd.Flags().StringVar(&nombreAmbiente2, "nombreAmbiente2", "ambiente2", "Nombre del segundo ambiente")
		rootCmd.MarkFlagRequired("username1")
		rootCmd.MarkFlagRequired("host1")
		rootCmd.MarkFlagRequired("userAdmin")
		rootCmd.MarkFlagRequired("passAdmin")
	})
}

func init() {
	initPackage()
}
