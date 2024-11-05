package slack

import (
	"context"
	"net/url"
	"strconv"
)

const (
	DEFAULT_LOGINS_COUNT = 100
	DEFAULT_LOGINS_PAGE  = 1
)

type TeamResponse struct {
	Team TeamInfo `json:"team"`
	SlackResponse
}

type TeamInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	EmailDomain string                 `json:"email_domain"`
	Icon        map[string]interface{} `json:"icon"`
}

type TeamProfileResponse struct {
	Profile TeamProfile `json:"profile"`
	SlackResponse
}

type TeamProfile struct {
	Fields []TeamProfileField `json:"fields"`
}

type TeamProfileField struct {
	ID             string          `json:"id"`
	Ordering       int             `json:"ordering"`
	Label          string          `json:"label"`
	Hint           string          `json:"hint"`
	Type           string          `json:"type"`
	PossibleValues []string        `json:"possible_values"`
	IsHidden       bool            `json:"is_hidden"`
	Options        map[string]bool `json:"options"`
}

type LoginResponse struct {
	Logins []Login `json:"logins"`
	Paging `json:"paging"`
	SlackResponse
}

type Login struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	DateFirst int    `json:"date_first"`
	DateLast  int    `json:"date_last"`
	Count     int    `json:"count"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	ISP       string `json:"isp"`
	Country   string `json:"country"`
	Region    string `json:"region"`
}

type BillableInfoResponse struct {
	BillableInfo map[string]BillingActive `json:"billable_info"`
	SlackResponse
}

type BillingActive struct {
	BillingActive bool `json:"billing_active"`
}

// AccessLogParameters contains all the parameters necessary (including the optional ones) for a GetAccessLogs() request
type AccessLogParameters struct {
	TeamID string
	Count  int
	Page   int
}

// NewAccessLogParameters provides an instance of AccessLogParameters with all the sane default values set
func NewAccessLogParameters() AccessLogParameters {
	return AccessLogParameters{
		Count: DEFAULT_LOGINS_COUNT,
		Page:  DEFAULT_LOGINS_PAGE,
	}
}

func (api *Client) teamRequest(ctx context.Context, path string, values url.Values) (*TeamResponse, error) {
	response := &TeamResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response, response.Err()
}

func (api *Client) billableInfoRequest(ctx context.Context, path string, values url.Values) (map[string]BillingActive, error) {
	response := &BillableInfoResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}

	return response.BillableInfo, response.Err()
}

func (api *Client) accessLogsRequest(ctx context.Context, path string, values url.Values) (*LoginResponse, error) {
	response := &LoginResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}
	return response, response.Err()
}

func (api *Client) teamProfileRequest(ctx context.Context, client httpClient, path string, values url.Values) (*TeamProfileResponse, error) {
	response := &TeamProfileResponse{}
	err := api.postMethod(ctx, path, values, response)
	if err != nil {
		return nil, err
	}
	return response, response.Err()
}

func (api *Client) GetExternalTeams(params *GetExternalTeamsParameters) (organizations []*Organization, nextCursor string, err error) {
	return api.GetExternalTeamsContext(context.Background(), params)
}

type GetExternalTeamsParameters struct {
	ConnectionStatusFilter string
	Cursor                 string
	Limit                  int
	SlackConnectPrefFilter []string
	SortDirection          string
	SortField              string
	WorkspaceFilter        []string
}

func (api *Client) GetExternalTeamsContext(ctx context.Context, params *GetExternalTeamsParameters) (organizations []*Organization, nextCursor string, err error) {
	values := url.Values{
		"token": {api.token},
	}

	if params != nil {
		if params.ConnectionStatusFilter != "" {
			values.Add("connection_status_filter", params.ConnectionStatusFilter)
		}
		if params.Cursor != "" {
			values.Add("cursor", params.Cursor)
		}
		if params.Limit != 0 {
			values.Add("limit", strconv.Itoa(params.Limit))
		}
		if len(params.SlackConnectPrefFilter) > 0 {
			values.Add("slack_connect_pref_filter", params.SlackConnectPrefFilter[0])
		}
		if params.SortDirection != "" {
			values.Add("sort_direction", params.SortDirection)
		}
		if params.SortField != "" {
			values.Add("sort_field", params.SortField)
		}
		if len(params.WorkspaceFilter) > 0 {
			values.Add("workspace_filter", params.WorkspaceFilter[0])
		}
	}

	response := struct {
		Organizations    []*Organization  `json:"organizations"`
		ResponseMetaData responseMetaData `json:"response_metadata"`
		SlackResponse
	}{}

	err = api.postMethod(ctx, "team.externalTeams.list", values, response)
	if err != nil {
		return nil, "", err
	}

	return response.Organizations, response.ResponseMetadata.Cursor, response.Err()

}

// GetTeamInfo gets the Team Information of the user
func (api *Client) GetTeamInfo() (*TeamInfo, error) {
	return api.GetTeamInfoContext(context.Background())
}

// GetOtherTeamInfoContext gets Team information for any team with a custom context
func (api *Client) GetOtherTeamInfoContext(ctx context.Context, team string) (*TeamInfo, error) {
	if team == "" {
		return api.GetTeamInfoContext(ctx)
	}
	values := url.Values{
		"token": {api.token},
	}
	values.Add("team", team)
	response, err := api.teamRequest(ctx, "team.info", values)
	if err != nil {
		return nil, err
	}
	return &response.Team, nil
}

// GetOtherTeamInfo gets Team information for any team
func (api *Client) GetOtherTeamInfo(team string) (*TeamInfo, error) {
	return api.GetOtherTeamInfoContext(context.Background(), team)
}

// GetTeamInfoContext gets the Team Information of the user with a custom context
func (api *Client) GetTeamInfoContext(ctx context.Context) (*TeamInfo, error) {
	values := url.Values{
		"token": {api.token},
	}

	response, err := api.teamRequest(ctx, "team.info", values)
	if err != nil {
		return nil, err
	}
	return &response.Team, nil
}

// GetTeamProfile gets the Team Profile settings of the user
func (api *Client) GetTeamProfile(teamID ...string) (*TeamProfile, error) {
	return api.GetTeamProfileContext(context.Background(), teamID...)
}

// GetTeamProfileContext gets the Team Profile settings of the user with a custom context
func (api *Client) GetTeamProfileContext(ctx context.Context, teamID ...string) (*TeamProfile, error) {
	values := url.Values{
		"token": {api.token},
	}
	if len(teamID) > 0 {
		values["team_id"] = teamID
	}

	response, err := api.teamProfileRequest(ctx, api.httpclient, "team.profile.get", values)
	if err != nil {
		return nil, err
	}
	return &response.Profile, nil
}

// GetAccessLogs retrieves a page of logins according to the parameters given
func (api *Client) GetAccessLogs(params AccessLogParameters) ([]Login, *Paging, error) {
	return api.GetAccessLogsContext(context.Background(), params)
}

// GetAccessLogsContext retrieves a page of logins according to the parameters given with a custom context
func (api *Client) GetAccessLogsContext(ctx context.Context, params AccessLogParameters) ([]Login, *Paging, error) {
	values := url.Values{
		"token": {api.token},
	}
	if params.TeamID != "" {
		values.Add("team_id", params.TeamID)
	}
	if params.Count != DEFAULT_LOGINS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_LOGINS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}

	response, err := api.accessLogsRequest(ctx, "team.accessLogs", values)
	if err != nil {
		return nil, nil, err
	}
	return response.Logins, &response.Paging, nil
}

type GetBillableInfoParams struct {
	User   string
	TeamID string
}

// GetBillableInfo ...
func (api *Client) GetBillableInfo(params GetBillableInfoParams) (map[string]BillingActive, error) {
	return api.GetBillableInfoContext(context.Background(), params)
}

// GetBillableInfoContext ...
func (api *Client) GetBillableInfoContext(ctx context.Context, params GetBillableInfoParams) (map[string]BillingActive, error) {
	values := url.Values{
		"token": {api.token},
	}

	if params.TeamID != "" {
		values.Add("team_id", params.TeamID)
	}

	if params.User != "" {
		values.Add("user", params.User)
	}

	return api.billableInfoRequest(ctx, "team.billableInfo", values)
}
