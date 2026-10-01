//go:build db

package handlers_test

// Seam HTTP hapus berjenjang terhadap Oracle NYATA (paket 7): DEV nol FK
// (K1), jadi yatim hanya tercegah bila Go menghapus anaknya.

import (
	"net/http"
	"strings"
	"testing"
)

func TestDBHapusKontrakTanpaYatimBatalNolBaris(t *testing.T) {
	u := pasangDB(t)
	s := u.skema
	u.exec(t, `INSERT INTO `+s+`.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR) VALUES ('1000002', '1000001')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR) VALUES ('1000009', '1000001')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYREINSURER_LIFE (ID, TREATYCONTRACTID) VALUES ('1000003', '1000002')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYSECURITYREINSURER_LIFE (ID, TREATYREINSURERID) VALUES ('1000004', '1000003')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYBUSINESS_LIFE (ID, TREATYCONTRACTID) VALUES ('1000005', '1000002')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYBUSINESS_LIFE (ID, TREATYCONTRACTID) VALUES ('1000006', '1000009')`)

	kode, badan := u.get(t, "/api/master-contract-retro-life/kontrak/1000002/dampak-hapus")
	if kode != http.StatusOK || !strings.Contains(badan, `"dampak":{"security":1,"reinsurer":1,"business":1}`) {
		t.Fatalf("popup: %d %s", kode, badan)
	}
	kode, _ = u.kirim(t, "DELETE", "/api/master-contract-retro-life/kontrak/1000002", `{"dampak":{"security":0,"reinsurer":0,"business":0}}`)
	if kode != http.StatusConflict || u.cacah(t, "TREATYSECURITYREINSURER_LIFE", "") != 1 {
		t.Fatalf("batal/basi harus nol baris terhapus: %d", kode)
	}
	kode, badan = u.kirim(t, "DELETE", "/api/master-contract-retro-life/kontrak/1000002", `{"dampak":{"security":1,"reinsurer":1,"business":1}}`)
	if kode != http.StatusOK {
		t.Fatalf("hapus: %d %s", kode, badan)
	}
	if u.cacah(t, "TREATYSECURITYREINSURER_LIFE", "") != 0 || u.cacah(t, "TREATYREINSURER_LIFE", "") != 0 ||
		u.cacah(t, "TREATYBUSINESS_LIFE", "") != 1 || u.cacah(t, "TREATYCONTRACT_LIFE", "") != 1 {
		t.Error("yatim tertinggal atau pohon lain tersentuh")
	}
}
