package access

import (
	"fmt"

	"github.com/yinshaoze/BSProxy/common/set"
	"github.com/yinshaoze/BSProxy/config"
)

func GetTargetList(listName string) (*set.StringSet, error) {
	set, ok := config.Lists[listName]
	if ok {
		return set, nil
	}
	return nil, fmt.Errorf("list %q not found", listName)
}
