package repository

// Uji pengurai POHON tab Limits proporsional — tanpa Oracle.

import "testing"

const dokPohon = `{
 "Limits":[
  {"TreatyType":"QUOTA SHARE","TreatyTypeID":10042,"Cover":"riskcat","MDPList":[{"Value":"9"}],
   "Detail":[
    {"TreatyGroup":"PROPERTY","TreatyGroupID":"77","TreatyType":"QUOTA SHARE","QSPct":"25.00","RIONR":null,
     "COBList":[{"ClassOfBusiness":"FIRE","ClassOfBusinessID":"1"}],
     "IOOLimitList":[{"Currency":"IDR","Value":"1500000000","Note":"QUOTA SHARE","Layer":"1"}],
     "SpreadingList":[{"Pct":"60"}]}
   ]},
  {"TreatyType":"SURPLUS"}
 ]}`

func TestPohonLimitsTigaTingkat(t *testing.T) {
	p := PohonLimitsDariDokumen([]byte(dokPohon))
	if len(p) != 2 {
		t.Fatalf("%d Kind of Treaty, mau 2", len(p))
	}
	if p[0]["TreatyType"] != "QUOTA SHARE" || p[0]["TreatyTypeID"] != "10042" {
		t.Errorf("tingkat 1: %v", p[0])
	}
	d := p[0]["Detail"].([]map[string]any)
	if len(d) != 1 || d[0]["TreatyGroup"] != "PROPERTY" || d[0]["QSPct"] != "25.00" {
		t.Fatalf("tingkat 2: %v", d)
	}
	io := d[0]["IOOLimitList"].([]map[string]any)
	if len(io) != 1 || io[0]["Value"] != "1500000000" || io[0]["Note"] != "QUOTA SHARE" {
		t.Errorf("IOOLimitList: %v", io)
	}
}

// ⛔ Daftar-izin: kunci yang tidak diikat Section tidak dibawa; `null` ADA
// dan kosong; kunci yang tidak ada TIDAK ada.
func TestPohonLimitsDaftarIzinDanKunciTakAda(t *testing.T) {
	p := PohonLimitsDariDokumen([]byte(dokPohon))
	if _, ada := p[0]["Cover"]; ada {
		t.Error("Cover (medan layer non-prop) ikut terbawa")
	}
	if _, ada := p[0]["MDPList"]; ada {
		t.Error("MDPList ikut terbawa")
	}
	d := p[0]["Detail"].([]map[string]any)
	if v, ada := d[0]["RIONR"]; !ada || v != "" {
		t.Errorf("RIONR null harus ada dan kosong: %v %v", ada, v)
	}
	if _, ada := d[0]["Surplus"]; ada {
		t.Error("Surplus tidak ada di dokumen, tetapi terbaca ada")
	}
	if _, ada := d[0]["SpreadingList"]; ada {
		t.Error("SpreadingList tidak diikat DetailLimits, tetapi terbawa")
	}
	if cob := d[0]["COBList"].([]map[string]any); len(cob) != 1 || len(cob[0]) != 1 {
		t.Errorf("COBList harus membawa ClassOfBusiness saja: %v", cob)
	}
	// Kind of Treaty tanpa Detail tetap punya larik Detail kosong.
	if d2 := p[1]["Detail"].([]map[string]any); len(d2) != 0 {
		t.Errorf("Detail kosong: %v", d2)
	}
}

func TestPohonLimitsDokumenRusakAtauTanpaLimits(t *testing.T) {
	for _, d := range []string{`{`, `{"ID":"1"}`, `{"Limits":"bukan larik"}`} {
		if p := PohonLimitsDariDokumen([]byte(d)); len(p) != 0 {
			t.Errorf("%q: mau nol simpul, dapat %v", d, p)
		}
	}
}
