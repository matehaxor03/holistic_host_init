package db_installer

import (
	"fmt"
	validate "github.com/matehaxor03/holistic_validator/validate"	
)

type HostInstaller struct {
	Validate func() []error
	Install func() ([]error)
}

func NewHostInstaller(users_directory []string, number_of_users uint64) (*HostInstaller, []error) {
	verify := validate.NewValidator()
	this_users_directory := users_directory
	this_number_of_users := number_of_users
	
	getUsersDirectory := func() []string {
		return this_users_directory
	}

	getNumberOfUsers := func() uint64 {
		return this_number_of_users
	}
	
	install := func() ([]error) {
		directory_parts := getUsersDirectory()
		directory := "/" 
		for index, directory_part := range directory_parts {
			directory += directory_part
			if index < len(directory_parts) - 1 {
				directory += "/"
			}
		}


		return nil
	}

	validate := func() []error {
		var errors []error
		temp_this_users_directory := getUsersDirectory()
		temp_this_number_of_users := getNumberOfUsers()
		
		if temp_this_number_of_users == 0 {
			errors = append(errors, fmt.Errorf("number_of_users cannot be 0"))
		}

		if temp_this_number_of_users % 10 != 0 {
			errors = append(errors, fmt.Errorf("number_of_users must be divisabe by 10 (e.g 10, 100, 1000, ...)"))
		}

		for _, directory_name_part := range temp_this_users_directory {
			directory_name_errors := verify.ValidateDirectoryName(directory_name_part)

			if directory_name_errors != nil {
				errors = append(errors, directory_name_errors...)
			}
		}

		if errors != nil {
			return errors
		}

		return nil
	}

	x := HostInstaller{
		Validate: func() []error {
			return validate()
		},
		Install: func() ([]error) {
			return install()
		},
	}

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

