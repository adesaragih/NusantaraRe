package services

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// tahapTiruan - PembacaTahapKasus dalam memori; "UJI-GALAT" = galat baca.
type tahapTiruan map[string]models.TahapKasus

func (t tahapTiruan) TahapKasus(_ context.Context, id string) (models.TahapKasus, error) {
	if id == "UJI-GALAT" {
		return models.TahapKasus{}, errors.New("UJI galat baca")
	}
	k, ada := t[id]
	if !ada {
		return k, repository.ErrKasusTidakAda
	}
	return k, nil
}

var tahapUji = tahapTiruan{
	"UJI-MKT":  {Position: "Offer", FlagOnGoingPolicy: "0"},
	"UJI-MKT2": {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"},
	"UJI-TL":   {Position: "ReasFacInTeamLeader", FlagOnGoingPolicy: "0"},
	"UJI-DEC":  {Position: "Decline", FlagOnGoingPolicy: "0"},
	"UJI-F1":   {Position: "Offer", FlagOnGoingPolicy: "1"},
}

// metodeTulis - setiap metode layanan yang MENULIS data kasus (endpoint PUT / POST / DELETE penyimpan), argumen minimal:
// penjaga baca-saja dicek sesudah identitas (dan akses IT), sebelum isian lain.
func metodeTulis(svc *Service) map[string]func(inti.Pelaku, string) error {
	ctx := context.Background()
	return map[string]func(inti.Pelaku, string) error{
		"general": func(p inti.Pelaku, id string) error {
			_, err := svc.SimpanGeneral(ctx, p, id, IsianGeneral{})
			return err
		},
		"objek":     func(p inti.Pelaku, id string) error { _, err := svc.GantiObjek(ctx, p, id, nil); return err },
		"klausa":    func(p inti.Pelaku, id string) error { _, err := svc.GantiKlausa(ctx, p, id, nil); return err },
		"spreading": func(p inti.Pelaku, id string) error { _, err := svc.SimpanSpreading(ctx, p, id, nil, nil); return err },
		"cedant":    func(p inti.Pelaku, id string) error { _, err := svc.SimpanCedant(ctx, p, id, "", nil); return err },
		"deduction": func(p inti.Pelaku, id string) error {
			_, err := svc.SimpanDeduction(ctx, p, id, IsianDeduction{})
			return err
		},
		"payment": func(p inti.Pelaku, id string) error {
			_, err := svc.SimpanPayment(ctx, p, id, IsianPayment{})
			return err
		},
		"scoring": func(p inti.Pelaku, id string) error {
			_, err := svc.SimpanScoring(ctx, p, id, [2]models.KolomScoring{})
			return err
		},
		"lossRecord": func(p inti.Pelaku, id string) error {
			_, err := svc.SimpanLossRecord(ctx, p, id, models.LossRecord{})
			return err
		},
		"korespondensiTambah": func(p inti.Pelaku, id string) error { _, err := svc.TambahKorespondensi(ctx, p, id); return err },
		"korespondensiSimpan": func(p inti.Pelaku, id string) error {
			_, err := svc.SimpanKorespondensi(ctx, p, id, "1", models.IsiKorespondensi{})
			return err
		},
		"korespondensiKirim": func(p inti.Pelaku, id string) error {
			_, err := svc.KirimKorespondensi(ctx, p, id, "1", models.IsiKorespondensi{})
			return err
		},
		"korespondensiHapus": func(p inti.Pelaku, id string) error { return svc.HapusKorespondensi(ctx, p, id, "1") },
	}
}

// tabIT - metode tab IT-only: syarat TAHAP saja (tanpa workbasket marketing).
var tabIT = map[string]bool{"lossRecord": true, "korespondensiTambah": true, "korespondensiSimpan": true,
	"korespondensiKirim": true, "korespondensiHapus": true}

// TestPenjagaBacaSaja - keputusan work owner 07-10-2026: setiap metode tulis -> 409 bila kasus di luar tahap marketing
// (posisi / FlagOnGoingPolicy), 403 tanpa workbasket ReasFacInMarketing (kecuali tab IT), 404 kasus asing; tahap marketing
// + workbasket lolos penjaga (galat lain sesudahnya boleh).
func TestPenjagaBacaSaja(t *testing.T) {
	svc := Baru(nil).DenganTahap(tahapUji).DenganJabatan(jabatanTiruan{"UJI-IT": "IT"})
	mkt := inti.Pelaku{AkunID: "UJI-IT", Peran: []string{"ReasFacInMarketing"}}
	tanpaWB := inti.Pelaku{AkunID: "UJI-IT", Peran: []string{"ReasFacInTeamLeader"}}
	metode := metodeTulis(svc)
	if len(metode) != 13 {
		t.Fatalf("%d metode", len(metode))
	}
	for nama, m := range metode {
		for _, id := range []string{"UJI-TL", "UJI-DEC", "UJI-F1"} {
			if err := m(mkt, id); !errors.Is(err, ErrKasusBacaSaja) {
				t.Errorf("%s %s: %v, mau 409 baca-saja", nama, id, err)
			}
		}
		if err := m(mkt, "UJI-NB-X"); !errors.Is(err, ErrKasusTidakAda) {
			t.Errorf("%s 404: %v", nama, err)
		}
		if err := m(inti.Pelaku{}, "UJI-TL"); !errors.Is(err, inti.ErrTanpaIdentitas) {
			t.Errorf("%s 401 lebih dulu: %v", nama, err)
		}
		err := m(tanpaWB, "UJI-MKT")
		if tabIT[nama] {
			if errors.Is(err, inti.ErrTanpaWewenang) || errors.Is(err, ErrKasusBacaSaja) {
				t.Errorf("%s tab IT tanpa workbasket marketing harus lolos penjaga: %v", nama, err)
			}
		} else if !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%s 403 tanpa workbasket: %v", nama, err)
		}
		for _, id := range []string{"UJI-MKT", "UJI-MKT2"} {
			if err := m(mkt, id); errors.Is(err, ErrKasusBacaSaja) || errors.Is(err, inti.ErrTanpaWewenang) {
				t.Errorf("%s %s tahap marketing ditolak penjaga: %v", nama, id, err)
			}
		}
	}
}

// TestAksesUbahKasus - bolehUbah / alasanBacaSaja GET kasus: lunak; aturan = Submit (alasanTidakBoleh) dan Scoring (tahap).
func TestAksesUbahKasus(t *testing.T) {
	ctx := context.Background()
	svc := Baru(nil).DenganTahap(tahapUji)
	mkt := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing"}}
	for _, c := range []struct {
		pelaku inti.Pelaku
		id     string
		boleh  bool
		alasan string
	}{
		{mkt, "UJI-MKT", true, ""},
		{mkt, "UJI-TL", false, `Case baca-saja: sudah tidak berada di tahap marketing Offer (posisi "ReasFacInTeamLeader").`},
		{inti.Pelaku{AkunID: "UJI"}, "UJI-MKT", false, "Akun tidak memegang workbasket ReasFacInMarketing."},
		{mkt, "UJI-GALAT", false, "Tahap case tidak terbaca."},
		{mkt, "UJI-NB-X", false, "Tahap case tidak terbaca."},
	} {
		if b, a := svc.AksesUbahKasus(ctx, c.pelaku, c.id); b != c.boleh || a != c.alasan {
			t.Errorf("%s: %v %q", c.id, b, a)
		}
	}
	if b, a := Baru(nil).AksesUbahKasus(ctx, mkt, "UJI-TL"); !b || a != "" {
		t.Errorf("tanpa pembaca tahap: %v %q", b, a)
	}
	// Satu sumber aturan: Submit dan Scoring memakai tahapMarketing.
	for id, k := range tahapUji {
		_, err := alasanTidakBoleh(models.KonteksSubmit{Position: k.Position, FlagOnGoingPolicy: k.FlagOnGoingPolicy}, mkt)
		if (err == nil) != tahapMarketing(k) || bolehUbahScoring(models.KasusScoring{Position: k.Position, FlagOnGoingPolicy: k.FlagOnGoingPolicy}) != tahapMarketing(k) {
			t.Errorf("%s: submit %v / scoring tidak sama dengan tahapMarketing", id, err)
		}
	}
}

// TestDariDasarMemasangTahap - penjaga baca-saja aktif di produksi: DariDasar memasang pembaca tahap (uji teks sumber -
// DariDasar butuh basis data sungguhan).
func TestDariDasarMemasangTahap(t *testing.T) {
	b, err := os.ReadFile("layanan.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "DenganKasus(kasus).DenganTahap(kasus)") {
		t.Error("DariDasar tidak memasang DenganTahap")
	}
}
