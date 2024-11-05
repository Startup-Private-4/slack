package slack

type Organization struct {
	TeamID              string `json:"team_id"`
	TeamName            string `json:"team_name"`
	TeamDomain          string `json:"team_domain"`
	PublicChannelCount  int    `json:"public_channel_count"`
	PrivateChannelCount int    `json:"private_channel_count"`
	IMChannelCount      int    `json:"im_channel_count"`
	MPIMChannelCount    int    `json:"mpim_channel_count"`
	ConnectedWorkspaces struct {
		WorkspaceID   string `json:"workspace_id"`
		WorkspaceName string `json:"workspace_name"`
	} `json:"connected_workspaces"`
	// SlackConnectPrefs   struct{} `json:"slack_connect_prefs"`
	ConnectionStatus    string `json:"connection_status"`
	LastActiveTimestamp int    `json:"last_active_timestamp"`
	IsSponsored         bool   `json:"is_sponsored"`
	Canvas              struct {
		TotalCount       int `json:"total_count"`
		OwnershipDetails []struct {
			TeamID string `json:"team_id"`
		} `json:"ownership_details"`
	} `json:"canvas"`
	Lists struct {
		TotalCount       int `json:"total_count"`
		OwnershipDetails []struct {
			TeamID string `json:"team_id"`
		} `json:"ownership_details"`
	} `json:"lists"`
}
