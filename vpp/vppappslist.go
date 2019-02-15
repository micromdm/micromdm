package vpp

type VPPAppResponse struct {
	VPPAppsList []VPPAppData `json:"vpp-apps"`
}

type VPPAppData struct {
	Asset         Asset            `json:"vpp-asset,omitempty"`
	ClientContext ClientContext    `json:"client-context,omitempty"`
	Metadata      *ContentMetadata `json:"vpp-metadata,omitempty"`
}

func (c *Client) GetVPPApps(ids ...string) (*VPPAppResponse, error) {

	assetsSrvOptions := AssetsSrvOptions{
		SToken:               c.VPPToken.SToken,
		IncludeLicenseCounts: true,
	}

	assetsSrv, err := c.GetAssetsSrv(assetsSrvOptions)
	if err != nil {
		return nil, err
	}

	var assets = assetsSrv.Assets

	// Filter out Assets
	if len(ids) == 1 && ids[0] == "" {
		ids = []string{}
	}
	if len(ids) > 0 {
		keep := []Asset{}
		for _, a := range assets {
			for _, id := range ids {
				if a.AdamIDStr == id {
					keep = append(keep, a)
				}
			}
		}
		assets = keep
	}

	response := VPPAppResponse{}

	for _, a := range assets {

		var vppApp VPPAppData

		vppApp.Asset = a

		var clientContext ClientContext
		err = DecodeToClientContext(assetsSrv.ClientContext, &clientContext)
		if err != nil {
			return nil, err
		}

		vppApp.ClientContext = clientContext

		metadata, err := c.GetContentMetadata(ContentMetadataOptions{
			ID: a.AdamIDStr,
		})
		if err != nil {
			return nil, err
		}
		vppApp.Metadata = metadata

		response.VPPAppsList = append(response.VPPAppsList, vppApp)
	}

	return &response, nil
}
