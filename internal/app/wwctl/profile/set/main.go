package set

import (
	"fmt"
	"os"

	apiprofile "github.com/hpcng/warewulf/internal/pkg/api/profile"
	"github.com/hpcng/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/hpcng/warewulf/internal/pkg/api/util"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func CobraRunE(cmd *cobra.Command, args []string) (err error) {
	// OptionStrMap, netWithoutName := apinode.AddNetname(OptionStrMap)
	// if netWithoutName {
	// 	return errors.New("a netname must be given for any network related configuration")
	// }
	buffer, err := yaml.Marshal(profileConf)
	if err != nil {
		wwlog.Error("Cant marshall nodeInfo", err)
		os.Exit(1)
	}
	set := wwapiv1.NodeSetParameter{
		NodeConfYaml: string(buffer[:]),
		NetdevDelete: SetNetDevDel,
		AllNodes:     SetNodeAll,
		Force:        SetForce,
		NodeNames:    args,
	}

	if !SetYes {
		var nodeCount uint
		// The checks run twice in the prompt case.
		// Avoiding putting in a blocking prompt in an API.
		_, nodeCount, err = apiprofile.ProfileSetParameterCheck(&set, false)
		if err != nil {
			return
		}
		yes := util.ConfirmationPrompt(fmt.Sprintf("Are you sure you want to modify %d nodes(s)", nodeCount))
		if !yes {
			return
		}
	}
	return apiprofile.ProfileSet(&set)
}
