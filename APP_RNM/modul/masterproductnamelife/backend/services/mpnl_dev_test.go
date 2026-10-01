package services_test

// Lanjutan 1 - OQ yang ditutup data DEV (01-10-2026, agregat `JSONDATA`, baca-saja).

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// L2 (OQ-MPNL-02): 196 dari 196 produk DEV punya baris inward ber-`ID` sama; 197 dari 199
// baris inward ber-`PRODUCTID` = `ID`. Produk baru: inward ber-`ID` = ID produk = `PRODUCTID`.
func TestInwardBerIDSamaDenganProduk(t *testing.T) {
	l, g := layananMaster()
	p, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true)
	if err != nil {
		t.Fatal(err)
	}
	if p.Inward.ID != p.ID || p.Inward.ProductID != p.ID {
		t.Errorf("inward ID = PRODUCTID = ID produk: %q %q %q", p.ID, p.Inward.ID, p.Inward.ProductID)
	}
	isi, ada := g.Inward[p.ID]
	if !ada || !strings.Contains(isi, `"ID":"`+p.ID+`"`) || !strings.Contains(isi, `"PRODUCTID":"`+p.ID+`"`) {
		t.Errorf("baris M_PRODUCTINWARD_LIFE ber-ID produk: %v %s", ada, isi)
	}
	if len(g.Inward) != 1 {
		t.Errorf("satu baris inward, nol sequence inward: %d", len(g.Inward))
	}
}

// L5 (OQ-MPNL-09): 153 produk DEV ber-`OutwardList[0]`, NOL ber-`OVR_COMM` terisi - ditulis
// kosong seperti Pega. L7 (OQ-MPNL-15): `TREATYCONTRACTID` tidak ada di `TREATYCONTRACT_LIFE`
// maupun `TREATYYEAR_LIFE` - tetap kosong seperti Pega.
func TestOutwardOvrCommDanTreatyContractIDKosong(t *testing.T) {
	l, g := layananMaster()
	kontrakUji(g)
	m := produkMasuk()
	m.Umum.IsORS, m.HitungOutward, m.Inward.Mature = true, true, "2027-02-28"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.OutwardList) == 0 {
		t.Fatal("OutwardList terisi dari kontrak OR")
	}
	for i, b := range p.OutwardList {
		if b.OvrComm != "" || b.TreatyContractID != "" {
			t.Errorf("baris %d: OVR_COMM %q TREATYCONTRACTID %q", i, b.OvrComm, b.TreatyContractID)
		}
	}
	if n := strings.Count(g.Umum[p.ID], `"OVR_COMM":""`); n != len(p.OutwardList) {
		t.Errorf("setiap baris membawa kunci OVR_COMM kosong (dibaca view PRODUCT_LIFE): %d dari %d", n, len(p.OutwardList))
	}
	if n := strings.Count(g.Umum[p.ID], `"TREATYCONTRACTID":""`); n != len(p.OutwardList) {
		t.Errorf("setiap baris membawa kunci TREATYCONTRACTID kosong: %d dari %d", n, len(p.OutwardList))
	}
	_ = repository.ReinsTypeOR
}
