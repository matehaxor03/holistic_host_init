package host_installer

import (
	"fmt"
	"strconv"
	validate "github.com/matehaxor03/holistic_validator/validate"
	host_client "github.com/matehaxor03/holistic_host_client/host_client"	
)

type HostInstaller struct {
	Validate func() []error
	Install func() ([]error)
}

func NewHostInstaller(users_directory []string, number_of_users uint64, userid_offset uint64) (*HostInstaller, []error) {
	verify := validate.NewValidator()
	this_users_directory := users_directory
	this_number_of_users := number_of_users
	this_userid_offset := userid_offset

	host_client, host_client_errors := host_client.NewHostClient()
	if host_client_errors != nil {
		return nil, host_client_errors
	}
	
	getUsersDirectory := func() []string {
		return this_users_directory
	}

	getNumberOfUsers := func() uint64 {
		return this_number_of_users
	}

	getUserIdOffset := func() uint64 {
		return this_userid_offset
	}

	create_user := func(username string, absolute_home_directory_path []string, primary_user_id uint64, primary_group_id uint64) []error {
		user_directory, user_directory_errors := host_client.AbsoluteDirectory(absolute_home_directory_path)
		if user_directory_errors != nil {
			return user_directory_errors
		}

		if !user_directory.Exists() {
			user_directory_create_errors := user_directory.Create()
			if user_directory_create_errors != nil {
				return user_directory_create_errors
			}
		}

		host_user, host_user_errors := host_client.HostUser(username) 

		if host_user_errors != nil {
			return host_user_errors
		}

		exists, exists_error := host_user.Exists()
		if exists_error != nil {
			return exists_error
		}

		if !*exists {
			create_errors := host_user.Create()
			if create_errors != nil {
				return create_errors
			}
		}

		set_unique_id_errors := host_user.SetUniqueId(primary_user_id)
		if set_unique_id_errors != nil {
			return set_unique_id_errors
		}

		group, group_errors := host_client.Group(username) 

		if group_errors != nil {
			return group_errors
		}

		group_exists, group_exists_errors := group.Exists()
		if group_exists_errors != nil {
			return group_exists_errors
		}

		if !*group_exists {
			group_create_errors := group.Create()
			if group_create_errors != nil {
				return group_create_errors
			}
		}

		set_group_unique_id_errors := group.SetUniqueId(primary_group_id)
		if set_group_unique_id_errors != nil {
			return set_group_unique_id_errors
		}

		add_user_to_group_errors := group.AddUser(*host_user)
		if add_user_to_group_errors != nil {
			return add_user_to_group_errors
		}

		set_user_primary_group_id_errors := host_user.SetPrimaryGroupId(primary_group_id)
		if set_user_primary_group_id_errors != nil {
			return set_user_primary_group_id_errors
		}

		create_home_directory_errors := host_user.CreateHomeDirectoryAbsoluteDirectory(*user_directory)

		if create_home_directory_errors != nil {
			return create_home_directory_errors
		}

		set_user_directory_errors := user_directory.SetOwnerRecursive(*host_user, *group)
		
		if set_user_directory_errors != nil {
			return set_user_directory_errors
		}

		return nil
	}
	
	install := func() ([]error) {
		users_directory_parts := getUsersDirectory()
		users_directory, users_directory_errors := host_client.AbsoluteDirectory(users_directory_parts)
		if users_directory_errors != nil {
			return users_directory_errors
		}

		if !users_directory.Exists() {
			directory_create_errors := users_directory.Create()
			if directory_create_errors != nil {
				return directory_create_errors
			}
		}

		temp_this_number_of_users := getNumberOfUsers()
		temp_this_userid_offset := getUserIdOffset()

		max_user_id := temp_this_userid_offset + temp_this_number_of_users

		for i := temp_this_userid_offset ; i < max_user_id; i++ {
			current_username := "holisticxyz_b" + strconv.FormatUint(i, 10) + "_"
			var current_username_directory []string
			current_username_directory = append(current_username_directory, users_directory.GetPath()...)
			current_username_directory = append(current_username_directory, current_username)
			create_user_errors := create_user(current_username, current_username_directory, i, i)
			if create_user_errors != nil {
				return create_user_errors
			}
		}

		return nil
	}

	validate := func() []error {
		var errors []error
		temp_this_users_directory := getUsersDirectory()
		temp_this_number_of_users := getNumberOfUsers()
		temp_this_userid_offset := getUserIdOffset()

		
		if temp_this_number_of_users == 0 {
			errors = append(errors, fmt.Errorf("number_of_users cannot be 0"))
		}

		if temp_this_userid_offset < 2048 {
			errors = append(errors, fmt.Errorf("userid_offset must be >= 2048"))
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

