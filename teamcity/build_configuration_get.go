package teamcity

import (
	//"errors"
	"fmt"
	"github.com/umweltdk/teamcity/types"
)

func (c *Client) GetBuildConfiguration(buildConfID string) (*types.BuildConfiguration, error) {
	path := fmt.Sprintf("/httpAuth/app/rest/buildTypes/id:%s", buildConfID)
	var buildConfig *types.BuildConfiguration

	err := c.doRetryRequest("GET", path, nil, &buildConfig)
	if err != nil {
		return nil, err
	}

	/*if buildConfig == nil {
		return nil, errors.New("build configuration not found")
	}*/

	return buildConfig, nil
}
