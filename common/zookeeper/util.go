package zookeeper

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/samuel/go-zookeeper/zk"
)

// ValidatePath will make sure a path is valid before sending the request
func ValidatePath(path string, isSequential bool) error {
	if path == "" {
		return zk.ErrInvalidPath
	}

	if path[0] != '/' {
		return zk.ErrInvalidPath
	}

	n := len(path)
	if n == 1 {
		// path is just the root
		return nil
	}

	if !isSequential && path[n-1] == '/' {
		return zk.ErrInvalidPath
	}

	// Start at rune 1 since we already know that the first character is
	// a '/'.
	for i, w := 1, 0; i < n; i += w {
		r, width := utf8.DecodeRuneInString(path[i:])
		switch {
		case r == '\u0000':
			return zk.ErrInvalidPath
		case r == '/':
			last, _ := utf8.DecodeLastRuneInString(path[:i])
			if last == '/' {
				return zk.ErrInvalidPath
			}
		case r == '.':
			last, lastWidth := utf8.DecodeLastRuneInString(path[:i])

			// Check for double dot
			if last == '.' {
				last, _ = utf8.DecodeLastRuneInString(path[:i-lastWidth])
			}

			if last == '/' {
				if i+1 == n {
					return zk.ErrInvalidPath
				}

				next, _ := utf8.DecodeRuneInString(path[i+w:])
				if next == '/' {
					return zk.ErrInvalidPath
				}
			}
		case r >= '\u0000' && r <= '\u001f',
			r >= '\u007f' && r <= '\u009f',
			r >= '\uf000' && r <= '\uf8ff',
			r >= '\ufff0' && r < '\uffff':
			return zk.ErrInvalidPath
		}
		w = width
	}
	return nil
}

// ParseACL parse acl string to []zk.ACL
func ParseACL(acl string) ([]zk.ACL, error) {
	acls := make([]zk.ACL, 0)

	aclstr := strings.Split(acl, ",")
	for _, a := range aclstr {
		acl, err := parseACL(a)
		if err != nil {
			return nil, err
		}

		acls = append(acls, acl)
	}

	return acls, nil
}

func parseACL(acl string) (zk.ACL, error) {
	as := strings.Split(acl, ":")
	if len(as) < 3 || len(as) > 4 {
		return zk.ACL{}, zk.ErrInvalidACL
	}

	var id string
	if len(as) == 3 {
		id = as[1]
	}
	if len(as) == 4 {
		id = as[1] + ":" + as[2]
	}

	var perms int32
	bs := []byte(as[len(as)-1])
	for _, b := range bs {
		switch b {
		case 'c':
			perms += zk.PermCreate
		case 'd':
			perms += zk.PermDelete
		case 'r':
			perms += zk.PermRead
		case 'w':
			perms += zk.PermWrite
		case 'a':
			perms += zk.PermAdmin
		}
	}

	return zk.ACL{
		Perms:  perms,
		Scheme: as[0],
		ID:     id,
	}, nil
}

// FormatACLs format ACL to string
func FormatACLs(acls []zk.ACL) string {
	if len(acls) == 0 {
		return ""
	}

	aclstrs := make([]string, len(acls))
	for i, a := range acls {
		ps := permsToString(a.Perms)
		aclstrs[i] = a.Scheme + ":" + a.ID + ":" + ps
	}

	return strings.Join(aclstrs, ",")
}

func permsToString(perms int32) string {
	if perms == 0 {
		return ""
	}

	var permstr string

	p := fmt.Sprintf("%05b", perms)
	for i, pb := range []byte(p) {
		switch i {
		case 0:
			if pb == '1' {
				permstr += "a"
			}
		case 1:
			if pb == '1' {
				permstr += "d"
			}
		case 2:
			if pb == '1' {
				permstr += "c"
			}
		case 3:
			if pb == '1' {
				permstr += "w"
			}
		case 4:
			if pb == '1' {
				permstr += "r"
			}
		}
	}

	return permstr
}
