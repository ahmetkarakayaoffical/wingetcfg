package wingetcfg

import (
	"errors"
	"fmt"
	"strings"
)

const (
	WinGetLocalGroupResource = "xPSDesiredStateConfiguration/xGroup"
)

// AddOrModifyLocalGroup adds or modify a local group.
// ID is an optional identifier.
// GroupName is required to identify the group.
// Description is an optional text that describes the account.
// Members define the members the group should have. This property will replace all the current group members with the specified members.
// Members should be specified as list of strings, separated by a semi-colon, including the local machine accounts usernames).
func AddOrModifyLocalGroup(ID, groupName string, description string, members string, hostname string) (*WinGetResource, error) {
	return NewLocalGroupResource(ID, groupName, description, members, EnsurePresent, hostname)
}

// IncludeMembersToGroup adds or modify a local group.
// ID is an optional identifier.
// GroupName is required to identify the group.
// MembersToInclude define the members that should be added to the group. Members should be specified as list of strings, separated by a semi-colon,
// defining the local machine accounts usernames).
func IncludeMembersToGroup(ID, groupName string, membersToInclude string, hostname string) (*WinGetResource, error) {
	r := WinGetResource{}
	r.Resource = WinGetLocalGroupResource

	// ID (optional)
	if ID != "" {
		r.ID = ID
	}

	// Directives
	r.Directives.Description = "Include members to group"
	r.Directives.AllowPreRelease = true

	// Settings
	r.Settings = map[string]any{}

	if groupName == "" {
		return nil, errors.New("groupName cannot be empty")
	}
	r.Settings["GroupName"] = groupName

	if membersToInclude == "" {
		return nil, errors.New("membersToInclude cannot be empty")
	}

	r.Settings["MembersToInclude"] = addHostnameToMembersIfRequired(membersToInclude, hostname)

	r.Settings["Ensure"] = EnsurePresent

	return &r, nil
}

// ExcludeMembersFromGroup exclude members from a local group.
// ID is an optional identifier.
// GroupName is required to identify the group.
// MembersToInclude define the members that should be added to the group. Members should be specified as list of strings, separated by a semi-colon,
// defining the local machine accounts usernames).
func ExcludeMembersFromGroup(ID, groupName string, membersToExclude string, hostname string) (*WinGetResource, error) {
	r := WinGetResource{}
	r.Resource = WinGetLocalGroupResource

	// ID (optional)
	if ID != "" {
		r.ID = ID
	}

	// Directives
	r.Directives.Description = "Exclude members from group"
	r.Directives.AllowPreRelease = true

	// Settings
	r.Settings = map[string]any{}

	if groupName == "" {
		return nil, errors.New("groupName cannot be empty")
	}
	r.Settings["GroupName"] = groupName

	if membersToExclude == "" {
		return nil, errors.New("membersToExclude cannot be empty")
	}

	r.Settings["MembersToExclude"] = addHostnameToMembersIfRequired(membersToExclude, hostname)

	r.Settings["Ensure"] = EnsurePresent

	return &r, nil
}

// RemoveLocalGroup remove a local group.
// ID is an optional identifier.
// GroupName is required to identify the group.
func RemoveLocalGroup(ID, groupName string) (*WinGetResource, error) {
	return NewLocalGroupResource(ID, groupName, "", "", EnsureAbsent, "")
}

// NewLocalGroupResource creates a new WinGetResource that contains the settings to manage a local group.
// ID is an optional identifier.
// GroupName is required to identify the group.
// Description is an optional text that describes the group.
// Reference: https://github.com/dsccommunity/xPSDesiredStateConfiguration/blob/main/source/DSC_xGroupResource/DSC_xGroupResource.psm1
func NewLocalGroupResource(ID, groupName string, description string, members string, ensure string, hostname string) (*WinGetResource, error) {
	r := WinGetResource{}
	r.Resource = WinGetLocalGroupResource

	// ID (optional)
	if ID != "" {
		r.ID = ID
	}

	// Directives
	r.Directives.Description = description
	r.Directives.AllowPreRelease = true

	// Settings
	r.Settings = map[string]any{}

	if groupName == "" {
		return nil, errors.New("groupName cannot be empty")
	}

	r.Settings["GroupName"] = groupName
	r.Settings["Description"] = description

	if members != "" {
		r.Settings["Members"] = addHostnameToMembersIfRequired(members, hostname)
	}

	r.Settings["Ensure"] = SetEnsureValue(ensure)

	return &r, nil
}

func addHostnameToMembersIfRequired(members string, hostname string) []string {
	processed := []string{}
	splittedMembers := strings.Split(members, ";")
	for _, m := range splittedMembers {
		if !strings.Contains(m, "\\") && !strings.Contains(m, "@") && !strings.Contains(strings.ToLower(m), "dc") {
			processed = append(processed, fmt.Sprintf("%s\\%s", strings.ToLower(hostname), m))
		} else {
			processed = append(processed, m)
		}
	}
	return processed
}
