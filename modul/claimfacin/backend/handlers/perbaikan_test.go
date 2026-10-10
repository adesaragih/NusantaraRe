package handlers_test

// Uji seam HTTP perbaikan kelainan XML (prompt §5) di jalur layanan:
//
//	butir 4  ChooseDla_Act tanpa lompatan mundur JMP1: log "AKSEPATSI" dan "DLA" masing-masing SATU kali
//	butir 5  HitDLAClaimFacin hanya di produksi (non-produksi: nol efek)
//	butir 8  CloseClaim: pesan "Please Print DLA Before Close This Claim" MENGHENTIKAN penutupan

import (
	"net/http"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

// baruUjiProduksi - layanan ber-penanda produksi (efek keluar terantre di outbox tiruan).
func baruUjiProduksi(t *testing.T) *uji {
	g, a := tiruan.Baru(), acuanUji()
	a.Tangga = g.TanggaKomite
	jam := func() time.Time { return time.Date(2026, 3, 5, 10, 0, 0, 0, models.Jakarta) }
	return &uji{t: t, srv: handlers.Router(services.Baru(g, a, jam, true), true), g: g, a: a}
}

// sampaiDLA - kasus Choose Surveyor dengan adjustment final yang sudah diterima komite (fixture) dan dicetak
// Acceptation-nya; pop-up Print DLA terbuka.
func (u *uji) sampaiDLA() string {
	u.t.Helper()
	id := u.adjustmentFinal()
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["AcceptanceStatus"], b["AcceptedNo"], b["IsApproved"], b["IsKomite"] = "1", "UJI-AKS.03.2026.00003", "1", "1"
	u.g.SetelHalaman(id, h)
	for _, p := range []services.PermintaanAksi{
		{Aksi: "Acceptation", Konteks: panelAdj(1)},
		{Aksi: "BukaDLA", Indeks: 1},
	} {
		kode, out := u.aksiT(id, p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	return id
}

func TestDLATanpaLompatanMundurDanHitDLAHanyaProduksi(t *testing.T) {
	for _, produksi := range []bool{false, true} {
		u := baruUji(t)
		if produksi {
			u = baruUjiProduksi(t)
		}
		id := u.sampaiDLA()
		awal := len(u.g.Efek)
		kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "ChooseDla", Konteks: "dla:1"})
		u.wajib(kode, http.StatusOK, out, "submit DLA")
		n := map[string]int{}
		for _, l := range u.g.Log {
			n[l.JenisService]++
		}
		if n[models.JenisServiceAkseptasi] != 1 || n[models.JenisServiceDLA] != 1 {
			t.Fatalf("produksi=%v log %v (lompatan mundur JMP1 akan menggandakannya)", produksi, n)
		}
		dla := 0
		for _, e := range u.g.Efek[awal:] {
			if len(e) > len(services.JenisEfekDLA) && e[:len(services.JenisEfekDLA)] == services.JenisEfekDLA {
				dla++
			}
		}
		if (produksi && dla != 1) || (!produksi && len(u.g.Efek) != 0) {
			t.Fatalf("produksi=%v efek %v", produksi, u.g.Efek)
		}
	}
}

func TestCloseClaimBerhentiPadaProteksiDLA(t *testing.T) {
	u := baruUji(t)
	id := u.adjustmentFinal()
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["AcceptanceStatus"], b["AcceptedNo"], b["IsApproved"], b["IsKomite"] = "1", "UJI-AKS.03.2026.00004", "1", "1"
	b["IsPrintAccept"], b["DirectToKasir"] = "1", "false"
	b[models.PropKomiteID] = "KMT-UJI9"
	ob := h.AmbilDaftar(models.DaftarObjek)[0]
	ob["IsFacretro"], ob["RemarksDLA"] = "1", ""
	u.g.SetelHalaman(id, h)
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "CloseClaim", Konteks: services.ModalTutup,
		Masukan: map[string]string{models.JalurTKRemarks: "UJI"}})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "close claim tanpa DLA")
	if !pesanMemuat(out, models.PesanDLASebelumTutup) {
		t.Fatalf("pesan %v", out["pesan"])
	}
	if k := u.g.Kasus[id]; k.Tertutup() {
		t.Fatal("klaim tertutup walau proteksi DLA berbunyi (perilaku Pega lama)")
	}
}
