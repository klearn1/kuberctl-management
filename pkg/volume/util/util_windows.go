//go:build windows
// +build windows

/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://docs.microsoft.com/en-us/windows/win32/fileio/file-access-rights-constants
const (
	// read = read data | read attributes
	READ_PERMISSIONS = 0x0001 | 0x0080

	// write = write data | append data | write attributes | write EA
	WRITE_PERMISSIONS = 0x0002 | 0x0004 | 0x0100 | 0x0010

	// execute = read data | file execute
	EXECUTE_PERMISSIONS = 0x0001 | 0x0020
)

const (
	EVERYONE = "Everyone"
)

var (
	advapi32                       = windows.MustLoadDLL("advapi32.dll")
	procSetEntriesInAclW           = advapi32.MustFindProc("SetEntriesInAclW")
	procGetExplicitEntriesFromAclW = advapi32.MustFindProc("GetExplicitEntriesFromAclW")
	procGetNamedSecurityInfoW      = advapi32.MustFindProc("GetNamedSecurityInfoW")
)

// Change the permissions of the specified file. Only the nine
// least-significant bytes are used, allowing access by the file's owner, the
// file's group, and everyone else to be explicitly controlled.
func Chmod(name string, fileMode os.FileMode) error {
	// https://support.microsoft.com/en-us/help/243330/well-known-security-identifiers-in-windows-operating-systems
	creatorOwnerSID, err := windows.StringToSid("S-1-3-0")
	if err != nil {
		return err
	}
	creatorGroupSID, err := windows.StringToSid("S-1-3-1")
	if err != nil {
		return err
	}
	everyoneSID, err := windows.StringToSid("S-1-1-0")
	if err != nil {
		return err
	}

	mode := windows.ACCESS_MASK(fileMode)
	return apply(
		name,
		true,
		false,
		grantSid(((mode&0700)<<23)|((mode&0200)<<9), creatorOwnerSID),
		grantSid(((mode&0070)<<26)|((mode&0020)<<12), creatorGroupSID),
		grantSid(((mode&0007)<<29)|((mode&0002)<<15), everyoneSID),
	)
}

// apply the provided access control entries to a file. If the replace
// parameter is true, existing entries will be overwritten. If the inherit
// parameter is true, the file will inherit ACEs from its parent.
func apply(name string, replace, inherit bool, entries ...windows.EXPLICIT_ACCESS) error {
	var oldAcl windows.Handle
	if !replace {
		var secDesc windows.Handle
		getNamedSecurityInfo(
			name,
			windows.SE_FILE_OBJECT,
			windows.DACL_SECURITY_INFORMATION,
			nil,
			nil,
			&oldAcl,
			nil,
			&secDesc,
		)
		defer windows.LocalFree(secDesc)
	}
	var acl *windows.ACL
	if err := setEntriesInAcl(
		entries,
		oldAcl,
		&acl,
	); err != nil {
		return err
	}
	defer windows.LocalFree((windows.Handle)(unsafe.Pointer(acl)))
	var secInfo windows.SECURITY_INFORMATION
	if !inherit {
		secInfo = windows.PROTECTED_DACL_SECURITY_INFORMATION
	} else {
		secInfo = windows.UNPROTECTED_DACL_SECURITY_INFORMATION
	}
	return windows.SetNamedSecurityInfo(
		name,
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|secInfo,
		nil,
		nil,
		acl,
		nil,
	)
}

// GetFileMode returns the mode of the given file.
func GetFileMode(file string) (os.FileMode, error) {
	var acl, secDesc windows.Handle
	var owner, group *windows.SID
	err := getNamedSecurityInfo(
		file,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION|windows.GROUP_SECURITY_INFORMATION|windows.DACL_SECURITY_INFORMATION,
		&owner,
		&group,
		&acl,
		nil,
		&secDesc,
	)
	if err != nil {
		return 0, err
	}

	defer windows.LocalFree(secDesc)

	entries, entriesHandle, err := getEntriesFromAcl(acl)
	if err != nil {
		return 0, err
	}

	defer windows.LocalFree(entriesHandle)

	ownerAccountName, ownerDomainName, _, err := owner.LookupAccount("")
	if err != nil {
		return 0, err
	}
	groupAccountName, groupDomainName, _, err := group.LookupAccount("")
	if err != nil {
		return 0, err
	}

	mode := 0
	for _, entry := range entries {
		accountName, domainName, err := lookupAccountSid((*windows.SID)(unsafe.Pointer(entry.Trustee.TrusteeValue)))
		if err != nil {
			return 0, err
		}

		// lookupAccountSid may return an empty string, which wouldn't match the account name from LookupAccount.
		if accountName == "" {
			accountName = "None"
		}

		perms := 0
		if (entry.AccessPermissions & READ_PERMISSIONS) == READ_PERMISSIONS {
			perms = 0x4
		}
		if (entry.AccessPermissions & WRITE_PERMISSIONS) == WRITE_PERMISSIONS {
			perms |= 0x2
		}
		if (entry.AccessPermissions & EXECUTE_PERMISSIONS) == EXECUTE_PERMISSIONS {
			perms |= 0x1
		}

		if accountName == ownerAccountName && domainName == ownerDomainName {
			mode |= perms << 6
		} else if accountName == groupAccountName && domainName == groupDomainName {
			mode |= perms << 3
		} else if accountName == EVERYONE {
			mode |= perms
		}
	}

	return os.FileMode(mode), nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-setentriesinaclw
func setEntriesInAcl(entries []windows.EXPLICIT_ACCESS, oldAcl windows.Handle, newAcl **windows.ACL) error {
	ret, _, _ := procSetEntriesInAclW.Call(
		uintptr(len(entries)),
		uintptr(unsafe.Pointer(&entries[0])),
		uintptr(oldAcl),
		uintptr(unsafe.Pointer(newAcl)),
	)
	if ret != 0 {
		return windows.Errno(ret)
	}
	return nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-getexplicitentriesfromaclw
func getEntriesFromAcl(acl windows.Handle) ([]windows.EXPLICIT_ACCESS, windows.Handle, error) {
	var entriesCount uint32
	var entriesHandle windows.Handle
	ret, _, err := procGetExplicitEntriesFromAclW.Call(
		uintptr(acl),
		uintptr(unsafe.Pointer(&entriesCount)),
		uintptr(unsafe.Pointer(&entriesHandle)),
	)
	if ret != 0 {
		return nil, entriesHandle, err
	}

	if entriesCount == 0 {
		return []windows.EXPLICIT_ACCESS{}, entriesHandle, nil
	}

	entries := make([]windows.EXPLICIT_ACCESS, entriesCount)
	entriesptr := unsafe.Pointer(uintptr(entriesHandle))
	copy(entries, (*[(1 << 31) - 1]windows.EXPLICIT_ACCESS)(entriesptr)[:entriesCount])

	return entries, entriesHandle, nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-getnamedsecurityinfow
func getNamedSecurityInfo(objectName string, objectType windows.SE_OBJECT_TYPE, secInfo windows.SECURITY_INFORMATION, owner, group **windows.SID, dacl, sacl, secDesc *windows.Handle) error {
	ret, _, _ := procGetNamedSecurityInfoW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(objectName))),
		uintptr(objectType),
		uintptr(secInfo),
		uintptr(unsafe.Pointer(owner)),
		uintptr(unsafe.Pointer(group)),
		uintptr(unsafe.Pointer(dacl)),
		uintptr(unsafe.Pointer(sacl)),
		uintptr(unsafe.Pointer(secDesc)),
	)
	if ret != 0 {
		return windows.Errno(ret)
	}
	return nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-lookupaccountsidw
func lookupAccountSid(sid *windows.SID) (string, string, error) {
	var nameSize, domainNameSize, peUse uint32

	// The first call will fail, and the necessary sizes will be written in nameSize and domainNameSize.
	windows.LookupAccountSid(
		nil,
		sid,
		nil,
		&nameSize,
		nil,
		&domainNameSize,
		&peUse,
	)

	bufferName := make([]uint16, nameSize)
	bufferDomainName := make([]uint16, domainNameSize)
	err := windows.LookupAccountSid(
		nil,
		sid,
		&bufferName[0],
		&nameSize,
		&bufferDomainName[0],
		&domainNameSize,
		&peUse,
	)

	if err != nil {
		return "", "", err
	}

	return windows.UTF16ToString(bufferName), windows.UTF16ToString(bufferDomainName), nil
}

// Create an EXPLICIT_ACCESS instance granting permissions to the provided SID.
func grantSid(accessPermissions windows.ACCESS_MASK, sid *windows.SID) windows.EXPLICIT_ACCESS {
	return windows.EXPLICIT_ACCESS{
		AccessPermissions: accessPermissions,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeValue: windows.TrusteeValueFromSID(sid),
		},
	}
}

// Create an EXPLICIT_ACCESS instance denying permissions to the provided SID.
func denySid(accessPermissions windows.ACCESS_MASK, sid *windows.SID) windows.EXPLICIT_ACCESS {
	return windows.EXPLICIT_ACCESS{
		AccessPermissions: accessPermissions,
		AccessMode:        windows.DENY_ACCESS,
		Inheritance:       windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeValue: windows.TrusteeValueFromSID(sid),
		},
	}
}
