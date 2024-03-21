package host_installer

import (
	"fmt"
	"strconv"
	validate "github.com/matehaxor03/holistic_validator/validate"
	host_client "github.com/matehaxor03/holistic_host_client/host_client"
	common "github.com/matehaxor03/holistic_common/common"
	"time"
)

type HostInstaller struct {
	Validate func() []error
	Install func() ([]error)
	Uninstall func() ([]error)
}

func NewHostInstaller(number_of_users uint64, userid_offset uint64) (*HostInstaller, []error) {
	verify := validate.NewValidator()
	this_number_of_users := number_of_users
	this_userid_offset := userid_offset

	host_client_instance, host_client_errors := host_client.NewHostClient()
	if host_client_errors != nil {
		return nil, host_client_errors
	}

	getNumberOfUsers := func() uint64 {
		return this_number_of_users
	}

	getUserIdOffset := func() uint64 {
		return this_userid_offset
	}

	create_user := func(username string, primary_user_id uint64, primary_group_id uint64, folder_group_username string) (*host_client.User, []error) {
		var absolute_home_directory_path []string
		var absolute_ssh_directory_path []string
		var absolute_io_directoy_path []string
		
		absolute_home_directory_path = append(absolute_home_directory_path, common.GetUsersDirectory()...)
		absolute_home_directory_path = append(absolute_home_directory_path, username)
		
		absolute_ssh_directory_path = append(absolute_ssh_directory_path, absolute_home_directory_path...)
		absolute_ssh_directory_path = append(absolute_ssh_directory_path, ".ssh")

		absolute_io_directoy_path = append(absolute_io_directoy_path, absolute_home_directory_path...)
		absolute_io_directoy_path = append(absolute_io_directoy_path, ".io")

		user_directory, user_directory_errors := host_client_instance.AbsoluteDirectory(absolute_home_directory_path)
		if user_directory_errors != nil {
			fmt.Println("user_directory_errors")
			return nil, user_directory_errors
		}

		user_directory_create_errors := user_directory.CreateIfDoesNotExist()
		if user_directory_create_errors != nil {
			return nil, user_directory_create_errors
		}
		
		ssh_directory, ssh_directory_errors := host_client_instance.AbsoluteDirectory(absolute_ssh_directory_path)
		if ssh_directory_errors != nil {
			fmt.Println("ssh_directory_errors")
			return nil, ssh_directory_errors
		}

		ssh_directory_create_errors := ssh_directory.CreateIfDoesNotExist()
		if ssh_directory_create_errors != nil {
			fmt.Println("ssh_directory_create_errors")
			return nil, ssh_directory_create_errors
		}

		io_directory, io_directory_errors := host_client_instance.AbsoluteDirectory(absolute_io_directoy_path)
		if io_directory_errors != nil {
			fmt.Println("absolute_io_directoy_path")
			return nil, io_directory_errors
		}

		io_directory_create_errors := io_directory.CreateIfDoesNotExist()
		if io_directory_create_errors != nil {
			return nil, io_directory_create_errors
		}

		user, user_errors := host_client_instance.User(username) 

		if user_errors != nil {
			fmt.Println("user_errors")
			return nil, user_errors
		}

		exists, exists_error := user.Exists()
		if exists_error != nil {
			fmt.Println("exists_error")
			return nil, exists_error
		}

		if !*exists {
			create_errors := user.Create()
			if create_errors != nil {
				fmt.Println("create_errors")
				return nil, create_errors
			}
		}

		set_unique_id_errors := user.SetUniqueId(primary_user_id)
		if set_unique_id_errors != nil {
			fmt.Println("set_unique_id_errors")
			return nil, set_unique_id_errors
		}

		group, group_errors := host_client_instance.Group(username) 

		if group_errors != nil {
			fmt.Println("group_errors")
			return nil, group_errors
		}

		group_exists, group_exists_errors := group.Exists()
		if group_exists_errors != nil {
			fmt.Println("group_exists_errors")
			return nil, group_exists_errors
		}

		if !*group_exists {
			group_create_errors := group.Create()
			if group_create_errors != nil {
				fmt.Println("group_create_errors")
				return nil, group_create_errors
			}
		}

		set_group_unique_id_errors := group.SetUniqueId(primary_group_id)
		if set_group_unique_id_errors != nil {
			fmt.Println("set_group_unique_id_errors")
			return nil, set_group_unique_id_errors
		}

		add_user_to_group_errors := group.AddUser(*user)
		if add_user_to_group_errors != nil {
			fmt.Println("add_user_to_group_errors")
			return nil, add_user_to_group_errors
		}

		set_user_primary_group_id_errors := user.SetPrimaryGroupId(primary_group_id)
		if set_user_primary_group_id_errors != nil {
			fmt.Println("set_user_primary_group_id_errors")
			return nil, set_user_primary_group_id_errors
		}

		create_home_directory_errors := user.CreateHomeDirectoryAbsoluteDirectory(*user_directory)

		if create_home_directory_errors != nil {
			fmt.Println("create_home_directory_errors")
			return nil, create_home_directory_errors
		}


		folder_group, folder_group_username_errors := host_client_instance.Group(folder_group_username) 

		if folder_group_username_errors != nil {
			fmt.Println("folder_group_username_errors")
			return nil, folder_group_username_errors
		}

		set_user_directory_errors := user_directory.SetOwnerRecursive(*user, *folder_group)
		
		if set_user_directory_errors != nil {
			fmt.Println("set_user_directory_errors")
			return nil, set_user_directory_errors
		}

		enable_bash_errors := user.EnableBinBash()
		if enable_bash_errors != nil {
			fmt.Println("enable_bash_errors")
			return nil, enable_bash_errors
		}

		set_password_errors := user.SetPassword("*")
		if set_password_errors != nil {
			fmt.Println("set_password_errors")
			return nil, set_password_errors
		}

		enable_full_disk_access_errors := user.EnableRemoteFullDiskAccess()
		if enable_full_disk_access_errors != nil {
			fmt.Println("enable_full_disk_access_errors")
			return nil, enable_full_disk_access_errors
		}

		return user, nil
	}
	
	install := func() ([]error) {
		var errors []error
		localhost := "127.0.0.1"

		ramdisk, ramdisk_errors := host_client_instance.Ramdisk(common.GetBaseDiskName(), uint64(2048*1000))
		if ramdisk_errors != nil {
			fmt.Println("ramdisk_errors")
			return ramdisk_errors
		}

		if !ramdisk.Exists() {
			ramdisk_create_errors := ramdisk.Create()
			if ramdisk_create_errors != nil {
				fmt.Println("ramdisk_create_errors")
				return ramdisk_create_errors
			}
			time.Sleep(30 * time.Second)
		}

		enable_filesystem_permissions_errors := ramdisk.EnableOwnership()
		if enable_filesystem_permissions_errors != nil {
			fmt.Println("enable_filesystem_permissions_errors")
			return enable_filesystem_permissions_errors
		}

		host, host_errors := host_client_instance.Host(localhost)
		if host_errors != nil {
			return host_errors
		}

		enable_host_ssh_errors := host.EnableSSH()
		if enable_host_ssh_errors != nil {
			return enable_host_ssh_errors
		}

		host_fingerprint, host_fingerprint_errors := host.GetSSHFingerprint()
		if host_fingerprint_errors != nil {
			return host_fingerprint_errors
		} else if len(*host_fingerprint) == 0 {
			errors = append(errors, fmt.Errorf("host fingerprint scan returned no results is ssh enabled on the host?"))
			return errors
		} else if len(*host_fingerprint) == 1 &&  (*host_fingerprint)[0] == "" {
			errors = append(errors, fmt.Errorf("host fingerprint scan returned no results is ssh enabled on the host?"))
			return errors
		}

		users_directory_parts := common.GetUsersDirectory()
		users_directory, users_directory_errors := host_client_instance.AbsoluteDirectory(users_directory_parts)
		if users_directory_errors != nil {
			fmt.Println("users_directory_errors")
			return users_directory_errors
		}

		if !users_directory.Exists() {
			directory_create_errors := users_directory.Create()
			if directory_create_errors != nil {
				fmt.Println("directory_create_errors")
				return directory_create_errors
			}

			read_only_other := int(0007)
			chmod_errors := users_directory.Chmod(read_only_other)
			if chmod_errors != nil {
				fmt.Println("chmod_errors")
				return chmod_errors
			}
		}

		temp_this_number_of_users := getNumberOfUsers()
		temp_this_userid_offset := getUserIdOffset()

		end_of_branch_user_ids := temp_this_userid_offset + temp_this_number_of_users
		end_of_users_id := end_of_branch_user_ids + temp_this_number_of_users
		current_unique_id := temp_this_userid_offset

		holistic_processor_username := "holisticxyz_holistic_processor_"
		holistic_processor_unique_id := end_of_users_id + 10
		holistic_processor_user, holistic_processor_create_user_errors := create_user(holistic_processor_username, holistic_processor_unique_id, holistic_processor_unique_id, holistic_processor_username)
		if holistic_processor_create_user_errors != nil {
			return holistic_processor_create_user_errors
		}

		holistic_processor_user_host_user := host_client_instance.HostUser(*host, *holistic_processor_user)

		holistic_webserver_username := "holisticxyz_holistic_webserver_"
		{
			holistic_webserver_unique_id := end_of_users_id + 11
			_, create_user_errors := create_user(holistic_webserver_username, holistic_webserver_unique_id, holistic_webserver_unique_id, holistic_webserver_username)
			if create_user_errors != nil {
				return create_user_errors
			}
		}

		holistic_queue_username := "holisticxyz_holistic_queue_"
		{
			holistic_queue_unique_id := end_of_users_id + 12
			_, create_user_errors := create_user(holistic_queue_username, holistic_queue_unique_id, holistic_queue_unique_id, holistic_queue_username)
			if create_user_errors != nil {
				return create_user_errors
			}
		}

		holistic_username := "holisticxyz_holistic_"
		{
			holistic_unique_id := end_of_users_id + 13
			_, create_user_errors := create_user(holistic_username, holistic_unique_id, holistic_unique_id, holistic_username)
			if create_user_errors != nil {
				return create_user_errors
			}
		}

		{
			for ; current_unique_id < end_of_branch_user_ids; current_unique_id++ {
				fmt.Println(current_unique_id)
				current_username := "holisticxyz_b" + strconv.FormatUint(current_unique_id, 10) + "_"
				branch_user, create_user_errors := create_user(current_username, current_unique_id, current_unique_id, holistic_processor_username)
				if create_user_errors != nil {
					return create_user_errors
				}

				
				host_branch_user := host_client_instance.HostUser(*host, *branch_user)
				generate_keys_errors := holistic_processor_user_host_user.GenerateSSHKey(host_branch_user)
				if generate_keys_errors != nil {
					return generate_keys_errors
				}
			}
		}

		return nil
	}

	delete_user := func(username string) []error {
		user, user_errors := host_client_instance.User(username)
		if user_errors != nil {
			return user_errors
		}

		disable_remote_full_disk_access_errors := user.DisableRemoteFullDiskAccess()
		if disable_remote_full_disk_access_errors != nil {
			return disable_remote_full_disk_access_errors
		}

		delete_if_exists_errors := user.DeleteIfExists()
		if delete_if_exists_errors != nil {
			return delete_if_exists_errors
		}

		group, group_errors := host_client_instance.Group(username)
		if group_errors != nil {
			return group_errors
		}

		delete_if_group_exists_errors := group.DeleteIfExists()
		if delete_if_group_exists_errors != nil {
			return delete_if_group_exists_errors
		}

		return nil
	}

	uninstall := func() []error {
		/*localhost := "127.0.0.1"
		host, host_errors := host_client_instance.Host(localhost)
		if host_errors != nil {
			return host_errors
		}*/

		/*
		enable_host_ssh_errors := host.DisableSSH()
		if enable_host_ssh_errors != nil {
			return enable_host_ssh_errors
		}*/

		temp_this_number_of_users := getNumberOfUsers()
		temp_this_userid_offset := getUserIdOffset()

		end_of_branch_user_ids := temp_this_userid_offset + temp_this_number_of_users
		current_unique_id := temp_this_userid_offset

		holistic_processor_username := "holisticxyz_holistic_processor_"
		holistic_processor_delete_user_errors := delete_user(holistic_processor_username)
		if holistic_processor_delete_user_errors != nil {
			return holistic_processor_delete_user_errors
		}

		holistic_webserver_username := "holisticxyz_holistic_webserver_"
		{
			delete_user_errors := delete_user(holistic_webserver_username)
			if delete_user_errors != nil {
				return delete_user_errors
			}
		}

		holistic_queue_username := "holisticxyz_holistic_queue_"
		{
			delete_user_errors := delete_user(holistic_queue_username)
			if delete_user_errors != nil {
				return delete_user_errors
			}
		}

		holistic_username := "holisticxyz_holistic_"
		{
			delete_user_errors := delete_user(holistic_username)
			if delete_user_errors != nil {
				return delete_user_errors
			}
		}

		{
			for ; current_unique_id < end_of_branch_user_ids; current_unique_id++ {
				fmt.Println(current_unique_id)
				current_username := "holisticxyz_b" + strconv.FormatUint(current_unique_id, 10) + "_"
				delete_user_errors := delete_user(current_username)
				if delete_user_errors != nil {
					return delete_user_errors
				}
			}
		}

		return nil
	}

	validate := func() []error {
		var errors []error
		temp_this_number_of_users := getNumberOfUsers()
		temp_this_userid_offset := getUserIdOffset()

		
		if temp_this_number_of_users == 0 {
			errors = append(errors, fmt.Errorf("number_of_users cannot be 0"))
		}

		if temp_this_userid_offset < 2000 {
			errors = append(errors, fmt.Errorf("userid_offset must be >= 2000"))
		}

		if temp_this_number_of_users % 10 != 0 {
			errors = append(errors, fmt.Errorf("number_of_users must be divisabe by 10 (e.g 10, 100, 1000, ...)"))
		}

		temp_this_users_directory := common.GetUsersDirectory()

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
		Uninstall: func() ([]error) {
			return uninstall()
		},
	}

	errors := validate()

	if errors != nil {
		return nil, errors
	}

	return &x, nil
}

