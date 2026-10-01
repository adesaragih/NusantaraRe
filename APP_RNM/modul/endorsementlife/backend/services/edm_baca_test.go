package services_test

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

var pelakuUji = inti.Pelaku{AkunID: "UJI-AKUN-1"}

func layananUji(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi, nil)
}

// gudangBerkasus - satu kasus terbuka atas polis NB sistem baru, tiga peserta.
func gudangBerkasus() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["NBLF-1"] = &tiruan.Polis{ID: "NBLF-1", ProdKe: 1, Kepala: map[string]string{"TYPE": "QR"}}
	g.Polis["EDMLF-1"] = &tiruan.Polis{
		ID: "EDMLF-1", OldPolicyNo: "UJI-PL-1", EdmType: "1", TglInput: "2026-10-01 09:00:00", Pembuat: "UJI-AKUN-1",
		Kepala: map[string]string{"TYPE": "QR", "CEDING_CO_NAME": "UJI-CEDING"},
	}
	g.Polis["EDMLF-2"] = &tiruan.Polis{ID: "EDMLF-2", OldPolicyNo: "UJI-PL-2", EdmType: "3", TglInput: "2026-10-01 10:00:00"}
	g.Polis["EDMLF-3"] = &tiruan.Polis{ID: "EDMLF-3", OldPolicyNo: "UJI-PL-3", EdmType: "1", TglInput: "2026-10-01 11:00:00",
		Status: models.StatusKasusDitolak}
	for i, s := range []string{models.StatusOld, models.StatusOld, models.StatusNew} {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{
			ID: "UJI-D" + string(rune('1'+i)), PolisID: "EDMLF-1", PLNumber: "UJI-PL-1", EdmStatus: s,
			Nilai: map[string]string{"CERTIFICATE_NO": "UJI-C" + string(rune('3'-i)), "NAME_OF_INSURED": "UJI-PESERTA"},
		})
	}
	return g
}

func TestInboxHanyaKasusTerbukaTerbaruDahulu(t *testing.T) {
	h, err := layananUji(gudangBerkasus()).Inbox(context.Background(), pelakuUji, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if h.Total != 2 || len(h.Baris) != 2 || h.Baris[0].CaseID != "EDMLF-2" || h.Baris[1].CaseID != "EDMLF-1" {
		t.Fatalf("kotak masuk = %+v", h)
	}
	if h.Baris[1].EndorsementNo != "UJI-PL-1" || h.Baris[1].Ceding != "UJI-CEDING" {
		t.Errorf("baris kotak masuk = %+v", h.Baris[1])
	}
	for _, b := range h.Baris {
		if !models.KasusEDM(b.CaseID) {
			t.Errorf("baris NB %q bocor ke kotak masuk endorsement", b.CaseID)
		}
	}
}

func TestInboxBerhalamanDanMenuntutIdentitas(t *testing.T) {
	l := layananUji(gudangBerkasus())
	h, err := l.Inbox(context.Background(), pelakuUji, 2, 1)
	if err != nil || len(h.Baris) != 1 || h.Baris[0].CaseID != "EDMLF-1" || h.Total != 2 {
		t.Fatalf("halaman 2 = %+v, %v", h, err)
	}
	h, _ = l.Inbox(context.Background(), pelakuUji, 0, 0)
	if h.Halaman != 1 || h.Ukuran != services.UkuranHalaman {
		t.Errorf("halaman tak sah tidak dijepit: %+v", h)
	}
	if _, err := l.Inbox(context.Background(), inti.Pelaku{}, 1, 20); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}

func TestBacaKasusMenyalakanTandaLayar(t *testing.T) {
	g := gudangBerkasus()
	l := layananUji(g)
	k, err := l.BacaKasus(context.Background(), pelakuUji, "EDMLF-1")
	if err != nil {
		t.Fatal(err)
	}
	if k.NomorPolis != "UJI-PL-1" || k.PLNumber != "UJI-PL-1" || !k.Terbuka() {
		t.Errorf("kasus = %+v", k)
	}
	if k.Cacah[models.StatusOld] != 2 || k.Cacah[models.StatusNew] != 1 || !k.CSVTerkunci || k.SudahSimpan {
		t.Errorf("tanda layar = cacah %v, csvTerkunci %v, sudahSimpan %v", k.Cacah, k.CSVTerkunci, k.SudahSimpan)
	}
	g.Rekap["EDMLF-1"] = 1
	if k, _ = l.BacaKasus(context.Background(), pelakuUji, "EDMLF-1"); !k.SudahSimpan {
		t.Error("IsJsonPolis tidak menyala sesudah rekap ada")
	}
}

func TestBacaKasusMenolakBukanKasusEDM(t *testing.T) {
	l := layananUji(gudangBerkasus())
	for _, id := range []string{"NBLF-1", "EDMLF-99", "", "EDMLF-"} {
		if _, err := l.BacaKasus(context.Background(), pelakuUji, id); !errors.Is(err, services.ErrKasusTidakAda) {
			t.Errorf("%q: %v", id, err)
		}
	}
}

func TestDaftarPesertaBerurutanDanBercacah(t *testing.T) {
	h, err := layananUji(gudangBerkasus()).DaftarPeserta(context.Background(), pelakuUji, "EDMLF-1", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if h.Total != 3 || len(h.Baris) != 2 || h.Baris[0].Nilai["CERTIFICATE_NO"] != "UJI-C1" {
		t.Fatalf("peserta = %+v", h)
	}
}

func TestGalatGudangDiteruskan(t *testing.T) {
	g := gudangBerkasus()
	g.Galat = errors.New("UJI-galat-oracle")
	if _, err := layananUji(g).Inbox(context.Background(), pelakuUji, 1, 20); err == nil {
		t.Error("galat gudang ditelan")
	}
}
