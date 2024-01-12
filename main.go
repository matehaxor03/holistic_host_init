package main

import (
	"fmt"
	"os"
	host_client "github.com/matehaxor03/holistic_host_client/host_client"
	host_installer "github.com/matehaxor03/holistic_host_init/host_installer"
	common "github.com/matehaxor03/holistic_common/common"
)

func main() {
	var errors []error
	host_client, host_client_errors := host_client.NewHostClient()
	if host_client_errors != nil {
		fmt.Println(fmt.Errorf("%s", host_client_errors))
		os.Exit(1)
	}

	number_of_users, number_of_users_errors := host_client.GetEnviornmentVariableUInt64(common.ENV_HOLISTIC_HOST_NUMBER_OF_USERS())
	if number_of_users_errors != nil {
		errors = append(errors, number_of_users_errors...)
	}

	if len(errors) > 0 {
		fmt.Println(fmt.Errorf("%s", errors))
		os.Exit(1)
	}

	host_installer,  host_installer_errors := host_installer.NewHostInstaller(common.GetUsersDirectory(), number_of_users)
	if host_installer_errors != nil {
		fmt.Println(fmt.Errorf("%s", host_installer_errors))
		os.Exit(1)
	}

	install_errors := host_installer.Install()
	if install_errors != nil {
		fmt.Println(fmt.Errorf("%s", install_errors))
		os.Exit(1)
	}

	os.Exit(0)
}