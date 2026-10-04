package services_test

// Enam gerbang penolakan dan empat klep pembatal - tiket E04. Data SINTETIS.
//
// Dibaca sesudah: services/gerbangtolak.go.

import (
	"errors"
	"testing"

	"nusantarare/modul/endorsmentfacin/backend/services"
	"nusantarare/modul/endorsmentfacin/backend/services/predikat"
)

func gerbang(t *testing.T, ubah func(*services.MasukanGerbang)) services.HasilGerbang {
	t.Helper()
	m := services.MasukanGerbang{
		PolicyNo: "UJI-POLIS-1", EdmType: "4", EdmTypeNew: "1", EndorsementInternalRetro: "0",
		Query:    services.HasilQueryGerbang{JenisBisnis: "UJI-FIRE"},
		Predikat: services.Predikat{IsFire: true},
	}
	ubah(&m)
	h, err := services.PeriksaGerbangPenolakan(m)
	if err != nil {
		t.Fatalf("galat sistem: %v", err)
	}
	return h
}

func pesanSatu(t *testing.T, h services.HasilGerbang, mau string) {
	t.Helper()
	if !h.Ditolak || len(h.Pesan) != 1 || h.Pesan[0].Teks != mau || h.Pesan[0].Medan != "PolicyNo" {
		t.Errorf("pesan %+v, mau satu %q", h.Pesan, mau)
	}
}

func TestTanpaPenghalangLolos(t *testing.T) {
	if h := gerbang(t, func(*services.MasukanGerbang) {}); h.Ditolak || len(h.Pesan) != 0 || h.SalahGerbang {
		t.Errorf("%+v", h)
	}
}

// TestEnamGerbangMenolakPadaKondisinya - tiap gerbang menolak pada kondisinya
// dan hanya pada kondisi itu.
func TestEnamGerbangMenolakPadaKondisinya(t *testing.T) {
	for _, u := range []struct {
		nama string
		ubah func(*services.MasukanGerbang)
		mau  string
	}{
		{"1 sudah batal (1)", func(m *services.MasukanGerbang) { m.Query.StatusEDM = "1" }, "Sudah Di endorsement Batal"},
		{"1 sudah batal (2)", func(m *services.MasukanGerbang) { m.Query.StatusEDM = "2" }, "Sudah Di endorsement Batal"},
		{"2 EDM belum selesai", func(m *services.MasukanGerbang) { m.Query.EDMBelumSelesai = []string{"UJI-EDM-9"} },
			"There's EDM with this policy no that haven't finish yet! UJI-EDM-9"},
		{"3 klaim", func(m *services.MasukanGerbang) { m.Query.Klaim = []string{"1"} }, "There's already a claim with this policy no"},
		{"4 pembayaran + batal", func(m *services.MasukanGerbang) { m.Query.CacahPembayaran = 1; m.EdmType = "1" },
			"There's already payment with this policy no"},
		{"5 renewal", func(m *services.MasukanGerbang) { m.Query.CacahRenewal = 1 }, "There's RNW with this policy no!"},
		{"6 retro tanpa fac out", func(m *services.MasukanGerbang) { m.EndorsementInternalRetro = "1" }, "This policy is not spreading FACOUT"},
	} {
		t.Run(u.nama, func(t *testing.T) {
			h := gerbang(t, u.ubah)
			pesanSatu(t, h, u.mau)
			if !h.SalahGerbang {
				t.Error("CARI3 tidak SALAH")
			}
		})
	}
}

// TestGerbangHanyaPadaKondisinya - kondisi tetangga tidak memicu gerbang.
func TestGerbangHanyaPadaKondisinya(t *testing.T) {
	for nama, ubah := range map[string]func(*services.MasukanGerbang){
		"status 3":                func(m *services.MasukanGerbang) { m.Query.StatusEDM = "3" },
		"klaim CARI1=2":           func(m *services.MasukanGerbang) { m.Query.Klaim = []string{"2"} },
		"pembayaran, bukan batal": func(m *services.MasukanGerbang) { m.Query.CacahPembayaran = 1; m.EdmType = "4" },
		"retro dengan fac out":    func(m *services.MasukanGerbang) { m.EndorsementInternalRetro = "1"; m.Query.CacahFacout = 1 },
		"fac out tanpa retro":     func(m *services.MasukanGerbang) { m.Query.CacahFacout = 0 },
	} {
		if h := gerbang(t, ubah); h.Ditolak {
			t.Errorf("%s: ditolak %+v", nama, h.Pesan)
		}
	}
}

// TestEmpatKlepMembatalkan - klep terisi membatalkan gerbang 1, 3, 4, 5.
func TestEmpatKlepMembatalkan(t *testing.T) {
	for _, u := range []struct {
		klep services.KlepPembatal
		ubah func(*services.MasukanGerbang)
	}{
		{services.KlepBatal, func(m *services.MasukanGerbang) { m.Query.StatusEDM = "1" }},
		{services.KlepKlaim, func(m *services.MasukanGerbang) { m.Query.Klaim = []string{"1"} }},
		{services.KlepPembayaran, func(m *services.MasukanGerbang) { m.Query.CacahPembayaran = 1; m.EdmType = "2" }},
		{services.KlepRenewal, func(m *services.MasukanGerbang) { m.Query.CacahRenewal = 1 }},
	} {
		h := gerbang(t, func(m *services.MasukanGerbang) {
			u.ubah(m)
			m.Query.Klep = map[services.KlepPembatal]bool{u.klep: true}
		})
		if h.Ditolak {
			t.Errorf("klep %d tidak membatalkan: %+v", u.klep, h.Pesan)
		}
		// Klep jenis LAIN tidak membatalkan.
		lain := services.KlepPembatal(u.klep%4 + 1)
		h = gerbang(t, func(m *services.MasukanGerbang) {
			u.ubah(m)
			m.Query.Klep = map[services.KlepPembatal]bool{lain: true}
		})
		if !h.Ditolak {
			t.Errorf("klep %d ikut membatalkan gerbang klep %d", lain, u.klep)
		}
	}
}

// TestGerbangDuaDanEnamTanpaKlep - keduanya tetap menolak walau seluruh klep
// terisi.
func TestGerbangDuaDanEnamTanpaKlep(t *testing.T) {
	semua := map[services.KlepPembatal]bool{services.KlepBatal: true, services.KlepKlaim: true, services.KlepPembayaran: true, services.KlepRenewal: true}
	for nama, ubah := range map[string]func(*services.MasukanGerbang){
		"2": func(m *services.MasukanGerbang) { m.Query.EDMBelumSelesai = []string{"UJI-EDM-9"} },
		"6": func(m *services.MasukanGerbang) { m.EndorsementInternalRetro = "1" },
	} {
		h := gerbang(t, func(m *services.MasukanGerbang) { ubah(m); m.Query.Klep = semua })
		if !h.Ditolak {
			t.Errorf("gerbang %s dibatalkan klep", nama)
		}
	}
}

// TestBypassRISlip - EdmType "4" dan EdmTypeNew "4": gerbang 1 dilompati (ke
// `jmp`), lalu gerbang 3-6 dilompati (ke `jmp2`); gerbang 2 TETAP berlaku.
func TestBypassRISlip(t *testing.T) {
	riSlip := func(m *services.MasukanGerbang) { m.EdmType, m.EdmTypeNew = "4", "4" }
	if h := gerbang(t, func(m *services.MasukanGerbang) {
		riSlip(m)
		m.Query.StatusEDM, m.Query.Klaim, m.Query.CacahRenewal, m.EndorsementInternalRetro = "1", []string{"1"}, 1, "1"
	}); h.Ditolak {
		t.Errorf("RI slip tidak melewati gerbang 1, 3-6: %+v", h.Pesan)
	}
	if h := gerbang(t, func(m *services.MasukanGerbang) { riSlip(m); m.Query.EDMBelumSelesai = []string{"UJI-EDM-9"} }); !h.Ditolak {
		t.Error("RI slip melewati gerbang 2")
	}
}

// TestJenisBisnisTakDikenal - langkah 9-18: jenis bisnis kosong, atau tidak
// satu pun predikat lini benar.
func TestJenisBisnisTakDikenal(t *testing.T) {
	h := gerbang(t, func(m *services.MasukanGerbang) { m.Query.JenisBisnis = ""; m.Predikat = services.Predikat{} })
	pesanSatu(t, h, "Invalid policy no!")
	if h.SalahGerbang {
		t.Error("CARI3 SALAH dari langkah sesudah 5.24")
	}
	h = gerbang(t, func(m *services.MasukanGerbang) { m.Predikat = services.Predikat{} })
	pesanSatu(t, h, "Invalid business, please contact IT!")
}

// TestRNMLTanpaJenisBisnisDibacaLife - langkah 8 di DALAM gerbang: polis
// berpola RNML dengan CARI2 kosong TIDAK mendapat "Invalid policy no!".
func TestRNMLTanpaJenisBisnisDibacaLife(t *testing.T) {
	h := gerbang(t, func(m *services.MasukanGerbang) {
		m.PolicyNo, m.Query.JenisBisnis, m.Predikat = "UJI-RNML-1", "", services.Predikat{IsLife: true}
	})
	if h.Ditolak {
		t.Errorf("polis RNML ditolak: %+v", h.Pesan)
	}
}

// TestJenisBisnisUntukPredikat - langkah 8: polis berpola RNML dibaca Life
// SEBELUM predikat dinilai.
func TestJenisBisnisUntukPredikat(t *testing.T) {
	if got := services.JenisBisnisUntukPredikat("UJI-RNML-1", "UJI-FIRE"); got != "Life" {
		t.Errorf("%q", got)
	}
	if got := services.JenisBisnisUntukPredikat("UJI-POLIS-1", "UJI-FIRE"); got != "UJI-FIRE" {
		t.Errorf("%q", got)
	}
}

// TestKodeTidakAmbigu - `CARI20==1`, `EdmType==1`: teks yang terbaca angka
// tetapi tidak sama sebagai teks ditolak (butir 20), bukan ditebak.
func TestKodeTidakAmbigu(t *testing.T) {
	_, err := services.PeriksaGerbangPenolakan(services.MasukanGerbang{
		PolicyNo: "UJI-POLIS-1", EdmType: "4", Query: services.HasilQueryGerbang{StatusEDM: "1.0", JenisBisnis: "X"},
		Predikat: services.Predikat{IsFire: true},
	})
	if !errors.Is(err, services.ErrTafsirKodeBerbeda) {
		t.Errorf("galat %v, mau ErrTafsirKodeBerbeda", err)
	}
}

// TestGerbangLiniLewatRegistry - langkah 10-16 dengan predikat dari registry
// EDM (E01) atas jenis bisnis SESUDAH langkah 8.
func TestGerbangLiniLewatRegistry(t *testing.T) {
	periksa := func(polis, cari2, oldID string) services.HasilGerbang {
		t.Helper()
		jenis := services.JenisBisnisUntukPredikat(polis, cari2)
		p, err := services.PredikatDari(predikat.KasusEDM{
			Properti: map[string]string{
				"pyWorkPage.Quotation.BusinessOldId":                 oldID,
				"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "3",
			},
			JenisBisnisLama: &jenis,
		})
		if err != nil {
			t.Fatal(err)
		}
		return gerbang(t, func(m *services.MasukanGerbang) {
			m.PolicyNo, m.Query.JenisBisnis, m.Predikat = polis, cari2, p
		})
	}
	if h := periksa("UJI-1", "FireStyle1", ""); h.Ditolak {
		t.Errorf("fire ditolak: %+v", h.Pesan)
	}
	if h := periksa("UJI-1", "JenisLain", ""); !h.Ditolak {
		t.Error("jenis bisnis tak dikenal lolos")
	}
	// RNML → "Life", tetapi IsLife membaca BusinessOldId.
	if h := periksa("UJI-RNML-1", "", "L3"); h.Ditolak {
		t.Errorf("life L3 ditolak: %+v", h.Pesan)
	}
	if h := periksa("UJI-RNML-1", "", "X"); !h.Ditolak {
		t.Error("RNML tanpa BusinessOldId life lolos")
	}
}
