package handlers_test

// Berkas salinan dokumen Pega lama (Copy Old / pemuat - laporan work owner 07-10-2026: "Commencement,Termination tdak
// muncul"): DATA_JSON hanya halaman PolicyTreatyIn, jadi TREATY_IN_ID kosong. Master dibaca ulang lewat NoOffer =
// TREATYID view (ID view tidak pernah sama dengan TREATYID; tanggal tunggal per TREATYID - DEV 07-10-2026). Fixture UJI-.

import (
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestBukaBerkasSalinanMemuatMasterLewatNoOffer(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := u.g.Halaman[id]
	h.Setel(models.HalamanMaster+".ID", "")
	h.Setel(models.HalamanPolis+".NoOffer", "UJI-TREATY-9")
	u.g.Kontrak["UJI-D-9B"] = models.BarisKontrak{"ID": "UJI-D-9B", "TREATYID": "UJI-TREATY-9",
		"COMMENCEMENT": "2025-01-01 00:00:00", "TERMINATION": "2025-12-31 00:00:00", "LAYER": "2"}
	u.g.Kontrak["UJI-D-9A"] = models.BarisKontrak{"ID": "UJI-D-9A", "TREATYID": "UJI-TREATY-9",
		"COMMENCEMENT": "2025-01-01 00:00:00", "TERMINATION": "2025-12-31 00:00:00", "LAYER": "1"}
	k := u.g.Kasus[id]
	k.StatusWork = models.StatusSelesai
	u.g.Kasus[id] = k
	kode, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	if kode != http.StatusOK {
		t.Fatalf("buka: %d %s", kode, isi)
	}
	ly := u.layar(isi)
	for j, harap := range map[string]string{
		"TreatyIn.Commencement": "2025-01-01 00:00:00",
		"TreatyIn.Termination":  "2025-12-31 00:00:00",
		// baris view PERTAMA (ORDER BY ID), sama dengan pxResults(1)
		models.HalamanPolis + ".Layer": "1",
	} {
		if v := ly.Halaman.Ambil(j); v != harap {
			t.Errorf("%s = %q, harap %q", j, v, harap)
		}
	}

	// kontrak lama tidak ada di view: berkas tetap terbuka, tanpa tanggal
	h.Setel(models.HalamanPolis+".NoOffer", "UJI-TREATY-TIDAK-ADA")
	kode, isi = u.panggil("GET", "/kasus/"+id, admin, nil)
	if kode != http.StatusOK {
		t.Fatalf("kontrak tidak ada di view: %d %s", kode, isi)
	}
	if v := u.layar(isi).Halaman.Ambil("TreatyIn.Commencement"); v != "" {
		t.Errorf("Commencement tanpa kontrak %q", v)
	}
}
