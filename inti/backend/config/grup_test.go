package config

import (
	"reflect"
	"strings"
	"testing"
)

// TestGrupNBFacIn - NBFACIN_GROUP_LOGIN / NBFACIN_GROUP_MARKETINGCODE (tiket 63): dipisah koma, dipangkas, butir kosong dibuang;
// kosong = nil; String hanya jumlah (nilai tidak bocor ke log). Nilai uji sintetis.
func TestGrupNBFacIn(t *testing.T) {
	t.Setenv("NBFACIN_GROUP_LOGIN", " UJI-A , ,uji-b,")
	t.Setenv("NBFACIN_GROUP_MARKETINGCODE", "UJI-MC-1")
	t.Setenv("NBFACIN_BONDING_MARKETINGNAME", "UJI NAMA A, UJI NAMA B")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	g := c.GrupNBFacIn
	if !reflect.DeepEqual(g.Login, []string{"UJI-A", "uji-b"}) || !reflect.DeepEqual(g.MarketingCode, []string{"UJI-MC-1"}) ||
		!reflect.DeepEqual(g.BondingMarketingName, []string{"UJI NAMA A", "UJI NAMA B"}) {
		t.Errorf("%#v", g)
	}
	if s := g.String(); s != "Grup{login:2 marketingCode:1 bondingMarketingName:2}" || strings.Contains(s, "UJI") {
		t.Errorf("String: %s", s)
	}
	t.Setenv("NBFACIN_GROUP_LOGIN", "")
	t.Setenv("NBFACIN_GROUP_MARKETINGCODE", " , ")
	if c, err := Load(); err != nil || c.GrupNBFacIn.Login != nil || c.GrupNBFacIn.MarketingCode != nil {
		t.Errorf("kosong: %v %#v", err, c.GrupNBFacIn)
	}
}
