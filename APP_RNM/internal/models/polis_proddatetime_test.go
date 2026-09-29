package models

// OQ-PL-13 DITUTUP (GILIRAN-17): `ProdDateTime` ikut XML - ambang 25 TERTANAM.

import (
	"os"
	"strings"
	"testing"
)

func TestAmbangProdDateTimeTertanamB1170(t *testing.T) {
	if AmbangProdDateTimePega != 25 {
		t.Errorf("ambang %d, mau 25", AmbangProdDateTimePega)
	}
	isi, err := os.ReadFile(`D:\XML\RNM_BRD\PremiumList Life\Activity\InsertJsonPolisLife_Act.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	baris := strings.Split(strings.ReplaceAll(string(isi), "><", ">\n<"), "\n")
	if len(baris) < 1170 || !strings.Contains(baris[1169], "@toDecimal(Local.currentdate)&gt;25") {
		t.Errorf("b1170 bukan gerbang `>25`: %q", baris[1169])
	}
}
