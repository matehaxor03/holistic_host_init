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
	
	number_of_users_value, number_of_users_errors := host_client.GetEnviornmentVariableValue(common.ENV_HOLISTIC_HOST_NUMBER_OF_USERS())
	if number_of_users_errors != nil {
		errors = append(errors, number_of_users_errors...)
	}

	users_offset_value, users_offset_value_errors := host_client.GetEnviornmentVariableValue(common.ENV_HOLISTIC_HOST_USERS_USERID_OFFSET())
	if users_offset_value_errors != nil {
		errors = append(errors, users_offset_value_errors...)
	}

	if len(errors) > 0 {
		fmt.Println(fmt.Errorf("%s", errors))
		os.Exit(1)
	}

	number_of_users, number_of_users_uint64_errors := number_of_users_value.GetUInt64Value()
	if number_of_users_uint64_errors != nil {
		fmt.Println(fmt.Errorf("%s", number_of_users_uint64_errors))
		os.Exit(1)
	}

	userid_offset, userid_offset_uint64_errors := users_offset_value.GetUInt64Value()
	if userid_offset_uint64_errors != nil {
		fmt.Println(fmt.Errorf("%s", userid_offset_uint64_errors))
		os.Exit(1)
	}

	host_installer,  host_installer_errors := host_installer.NewHostInstaller(number_of_users, userid_offset)
	if host_installer_errors != nil {
		fmt.Println(fmt.Errorf("%s", host_installer_errors))
		os.Exit(1)
	}

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("no args provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many args provided")
		os.Exit(1)
	} else {
		if args[0] == "install" {
			fmt.Println("installing...")
			install_errors := host_installer.Install()
			if install_errors != nil {
				fmt.Println(fmt.Errorf("%s", install_errors))
				os.Exit(1)
			}
		} else if args[0] == "uninstall" {
			fmt.Println("uninstalling...")
			install_errors := host_installer.Uninstall()
			if install_errors != nil {
				fmt.Println(fmt.Errorf("%s", install_errors))
				os.Exit(1)
			}
		} else {
			fmt.Println("provide either install or uninstall")
			os.Exit(1)
		}
	}
		

	

	os.Exit(0)
}