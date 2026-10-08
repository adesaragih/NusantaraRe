package services_test

// Uji tombol tulis layar ADJUSTMENT (EDM) — Save, Submit, Actions, Decline
// offer — di atas gudang tiruan.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
	"nusantarare/modul/treatyin/backend/services"
)

func (g *gudangTiruan) BacaKepalaPenyesuaian(_ context.Context, id string) (map[string]any, bool, error) {
	k, ada := g.kepalaEDM[id]
	salin := map[string]any{}
	for a, b := range k {
		salin[a] = b
	}
	return salin, ada, nil
}

func (g *gudangTiruan) SimpanPenyesuaian(_ context.Context, r models.RencanaPenyesuaian) error {
	if _, ada := g.kepalaEDM[r.ID]; ada && r.Draf {
		return repository.ErrPenyesuaianSudahAda
	}
	g.disimpanEDM = append(g.disimpanEDM, r)
	return nil
}

func (g *gudangTiruan) HapusPenyesuaian(_ context.Context, id string) error {
	g.dihapusEDM = append(g.dihapusEDM, id)
	return nil
}

func gudangEDM() *gudangTiruan {
	g := gudangSimpan()
	g.kepalaEDM = map[string]map[string]any{
		"1001001/R01": {"OLDID": "1001001", "EDMState": "1", "Position": "", "StatusAkseptasi": "", "TreatyContractName": "EDM"},
	}
	g.dokumenTersimpan["1001001/R01"] = map[string]any{"RNMShare": "10", "Portfolio": []any{map[string]any{"Description": "x"}}}
	g.dokumenTersimpan["1001001/R01#LAMA"] = map[string]any{"RNMShare": "4"}
	return g
}

func TestSaveDrafPenyesuaianMenulisKeduaSisi(t *testing.T) {
	g := gudangEDM()
	h, err := services.LayananDengan(g).SimpanPenyesuaian(context.Background(), admin, services.MasukanPenyesuaian{
		ID: "1001001/R02", Draf: true,
		Baru: services.SisiKiriman{
			Medan: map[string]any{"OLDID": "1001001/R01", "EDMState": "2", "ValueDifference.RNMShare": "1", "Position": "X"},
			Larik: map[string]any{"Portfolio": []any{map[string]any{"Description": "baru"}}},
		},
		Lama: &services.SisiKiriman{Medan: map[string]any{"RNMShare": "4"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := g.disimpanEDM[0]
	// TreatyInAddNew: Position Admin, PositionUsername operator.
	if !r.Draf || r.Baru["Position"] != models.PosisiAdmin || r.Baru["PositionUsername"] != "ADESAMUEL" || r.Lama["RNMShare"] != "4" {
		t.Errorf("rencana %+v", r)
	}
	// Kunci bertitik kembali ke halaman tertanamnya.
	if vd, _ := r.Baru["ValueDifference"].(map[string]any); vd["RNMShare"] != "1" {
		t.Errorf("ValueDifference %v", r.Baru["ValueDifference"])
	}
	if h.Pesan != "Data Sudah Disimpan Dengan ID : 1001001/R02" {
		t.Errorf("pesan %q", h.Pesan)
	}
}

func TestSaveDrafBerpengenalTersimpanDitolak(t *testing.T) {
	g := gudangEDM()
	_, err := services.LayananDengan(g).SimpanPenyesuaian(context.Background(), admin, services.MasukanPenyesuaian{ID: "1001001/R01", Draf: true})
	if !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("galat %v", err)
	}
}

func TestSavePenyesuaianTersimpanMenimpaClipboardDanMelaporkanKurs(t *testing.T) {
	g := gudangEDM()
	h, err := services.LayananDengan(g).SimpanPenyesuaian(context.Background(), admin, services.MasukanPenyesuaian{
		ID: "1001001/R01",
		Baru: services.SisiKiriman{
			Medan: map[string]any{"TreatyContractName": "UBAH", "StatusAkseptasi": "Resolve Complete"},
			Larik: map[string]any{"CurrencyList": []any{map[string]any{"Currency": "USD"}}},
		},
		Lama: &services.SisiKiriman{Medan: map[string]any{"RNMShare": "99"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := g.disimpanEDM[0]
	// Kepala + pendaratan tersimpan, ditimpa kiriman; kunci milik server diabaikan.
	if r.Baru["TreatyContractName"] != "UBAH" || r.Baru["RNMShare"] != "10" || r.Baru["StatusAkseptasi"] != "" || r.Baru["OLDID"] != "1001001" {
		t.Errorf("dokumen %v", r.Baru)
	}
	// Sisi Old penyesuaian tersimpan tidak disentuh.
	if r.Lama != nil {
		t.Errorf("sisi Old ditulis: %v", r.Lama)
	}
	if _, ada := r.Baru["CurrencyList"]; ada || len(h.KunciTakTersimpan) == 0 || h.KunciTakTersimpan[0] != "CurrencyList" {
		t.Errorf("kurs %v laporan %v", r.Baru["CurrencyList"], h.KunciTakTersimpan)
	}
}

func TestSubmitPenyesuaianNaikTanpaValidasiDanMenghitungSelisih(t *testing.T) {
	g := gudangEDM()
	// Ceding kosong — Submit EDM TIDAK memeriksanya ([2]–[4] mati).
	_, err := services.LayananDengan(g).KirimPenyesuaian(context.Background(), admin, services.MasukanKirimPenyesuaian{
		MasukanPenyesuaian: services.MasukanPenyesuaian{ID: "1001001/R01", Baru: services.SisiKiriman{Medan: map[string]any{"Comment": "ok", "Ceding": ""}}},
		Aksi:               services.AksiSubmit,
	})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpanEDM[0].Baru
	if d["Position"] != models.PosisiSecHead || d["StatusAkseptasi"] != "Accept" || d["PositionUsername"] != "SEC1, SEC2" {
		t.Errorf("langkah %v %v %v", d["Position"], d["StatusAkseptasi"], d["PositionUsername"])
	}
	// [7] EDMState 1: ValueDifference = New − Old (10 − 4).
	vd, _ := d["ValueDifference"].(map[string]any)
	if vd["RNMShare"] != "6" {
		t.Errorf("ValueDifference.RNMShare %v", vd["RNMShare"])
	}
	k, _ := d["CommentList"].([]any)
	if el, _ := k[len(k)-1].(map[string]any); el["Suggest"] != "ok" || el["IsApproved"] != "Accept" {
		t.Errorf("komentar %v", k)
	}
}

func TestActionsPenyesuaianHanyaPemegangPosisi(t *testing.T) {
	g := gudangEDM()
	g.kepalaEDM["1001001/R01"]["Position"] = models.PosisiSecHead
	g.kepalaEDM["1001001/R01"]["StatusAkseptasi"] = "Accept"
	l := services.LayananDengan(g)
	if _, err := l.KirimPenyesuaian(context.Background(), admin, services.MasukanKirimPenyesuaian{
		MasukanPenyesuaian: services.MasukanPenyesuaian{ID: "1001001/R01"}, Aksi: services.AksiAkseptasi, Pilihan: "Accept",
	}); !errors.Is(err, services.ErrBukanPemegangPosisi) {
		t.Errorf("Admin: %v", err)
	}
	if _, err := l.KirimPenyesuaian(context.Background(), secHead, services.MasukanKirimPenyesuaian{
		MasukanPenyesuaian: services.MasukanPenyesuaian{ID: "1001001/R01"}, Aksi: services.AksiAkseptasi, Pilihan: "Reject",
	}); err != nil {
		t.Fatal(err)
	}
	if d := g.disimpanEDM[0].Baru; d["Position"] != models.PosisiAdmin || d["StatusAkseptasi"] != "Reject" {
		t.Errorf("Reject %v %v", d["Position"], d["StatusAkseptasi"])
	}
}

func TestDeclineOfferPenyesuaianMenghapus(t *testing.T) {
	g := gudangEDM()
	l := services.LayananDengan(g)
	if _, err := l.HapusPenyesuaian(context.Background(), admin, services.MasukanHapusPenyesuaian{ID: "1001001/R01"}); err != nil {
		t.Fatal(err)
	}
	if len(g.dihapusEDM) != 1 || g.dihapusEDM[0] != "1001001/R01" {
		t.Errorf("dihapus %v", g.dihapusEDM)
	}
	// ⚠️ Pagar di luar ekspor: penyesuaian tuntas tidak dihapus.
	g.kepalaEDM["1001001/R01"]["StatusAkseptasi"] = models.StatusTuntas
	if _, err := l.HapusPenyesuaian(context.Background(), admin, services.MasukanHapusPenyesuaian{ID: "1001001/R01"}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("tuntas: %v", err)
	}
}
