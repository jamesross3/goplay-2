package zipplay_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURLDEcode(t *testing.T) {
	parsed, err := url.Parse("https://maps.googleapis.com/maps/api/js/ViewportInfoService.GetViewportInfo?1m6&1m2&1d36.364288159166776&2d-120.79917641914992&2m2&1d38.29458427720474&2d-117.81033081702661&2u10&4sen-US&5e4&6sr%40687000000&7b0&8e0&12e1&13shttps%3A%2F%2Fcalscape.org%2FArctostaphylos-densiflora-(Vine-Hill-Manzanita)&14b1&callback=_xdc_._ccgjm4&key=AIzaSyCDcz9Bbd1Vj6mv9elv2nqaqo3ExvoGlZo&token=24976")
	require.NoError(t, err)
	t.Logf("Parsed query params:\n")
	queryParams := parsed.Query()
	asJSON, err := json.MarshalIndent(queryParams, "", "  ")
	require.NoError(t, err)
	t.Logf(string(asJSON))
}
