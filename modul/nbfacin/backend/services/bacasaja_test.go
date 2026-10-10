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
	"UJI-UW":   {Position: "ReasFacInJuniorUnderwritingB", FlagOnGoingPolicy: "0"},
	"UJI-TL1":  {Position: "ReasFacInTeamLeader", FlagOnGoingPolicy: "1"},
	"UJI-DEC":  {Position: "Decline", FlagOnGoingPolicy: "0"},
	"UJI-F1":   {Position: "Offer", FlagOnGoingPolicy: "1"},
	// tiket 64 lanjutan: tahap Binding sebelum / sesudah tombol Binding.
	"UJI-BND":  {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2"},
	"UJI-BND0": {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2", ConfirmBinding: "0"},
	"UJI-BND1": {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2", ConfirmBinding: "1"},
	"UJI-TLB":  {Position: "ReasFacInTeamLeader", FlagOnGoingPolicy: "2"},
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
		"objek": func(p inti.Pelaku, id string) error { _, err := svc.GantiObjek(ctx, p, id, nil); return err },
		"objekAneka": func(p inti.Pelaku, id string) error {
			_, err := svc.GantiObjekAneka(ctx, p, id, nil, nil)
			return err
		},
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
		"csvObjek":           func(p inti.Pelaku, id string) error { _, err := svc.SimpanCsvObjek(ctx, p, id, nil); return err },
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
	if len(metode) != 15 { // + objekAneka tiket 73 P1
		t.Fatalf("%d metode", len(metode))
	}
	for nama, m := range metode {
		for _, id := range []string{"UJI-UW", "UJI-DEC", "UJI-F1"} {
			if nama == "spreading" && id == "UJI-UW" {
				continue // tiket 59: tab Spreading terbuka di tahap uw (403 tanpa workbasket posisi) - TestMatriksPenjagaSpreading
			}
			if err := m(mkt, id); !errors.Is(err, ErrKasusBacaSaja) {
				t.Errorf("%s %s: %v, mau 409 baca-saja", nama, id, err)
			}
		}
		if err := m(mkt, "UJI-NB-X"); !errors.Is(err, ErrKasusTidakAda) {
			t.Errorf("%s 404: %v", nama, err)
		}
		if err := m(inti.Pelaku{}, "UJI-UW"); !errors.Is(err, inti.ErrTanpaIdentitas) {
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
		// Tiket 58: Team Leader mengubah case di posisinya; marketing tidak (403, kecuali tab IT - syarat tahap saja).
		tl := inti.Pelaku{AkunID: "UJI-IT", Peran: []string{"ReasFacInTeamLeader"}}
		if err := m(tl, "UJI-TL"); errors.Is(err, ErrKasusBacaSaja) || errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%s TL di posisi TL ditolak penjaga: %v", nama, err)
		}
		if err := m(mkt, "UJI-TL"); tabIT[nama] == errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%s marketing di posisi TL: %v", nama, err)
		}
		// K69-8 (RALAT 69a-1): TL tahap Policy (flag 1) - seluruh tab baca-saja (409), juga bagi pemegang workbasket TL.
		for _, p := range []inti.Pelaku{tl, mkt} {
			if err := m(p, "UJI-TL1"); !errors.Is(err, ErrTahapPolicyBacaSaja) {
				t.Errorf("%s %v di posisi TL Policy: %v, mau 409 tahap Policy", nama, p.Peran, err)
			}
		}
	}
}

// TestPenjagaBacaSajaBinding - tiket 64 lanjutan (work owner 08-10-2026 "inikan lagi dimarketing harusnya msh bisa di edit,
// kecuali sudah di klik binding, baru smua tidak bisa diedit lagi"): tahap Binding (ReasFacInMarketing, flag 2) sebelum
// tombol Binding (CONFIRM_BINDING kosong / "0") - setiap metode tulis lolos penjaga bagi pemegang workbasket marketing, 403
// tanpa workbasket (kecuali tab IT); sesudah tombol Binding ("1") - 409 baca-saja. Scoring tetap baca-saja di flag 2
// (CallDualScoringRisk "Write" hanya FlagOnGoingPolicy 0). Flag 2 di posisi lain tetap baca-saja.
func TestPenjagaBacaSajaBinding(t *testing.T) {
	ctx := context.Background()
	svc := Baru(nil).DenganTahap(tahapUji).DenganJabatan(jabatanTiruan{"UJI-IT": "IT"})
	mkt := inti.Pelaku{AkunID: "UJI-IT", Peran: []string{"ReasFacInMarketing"}}
	tanpaWB := inti.Pelaku{AkunID: "UJI-IT", Peran: []string{"ReasFacInTeamLeader"}}
	for nama, m := range metodeTulis(svc) {
		for _, id := range []string{"UJI-BND", "UJI-BND0"} {
			err := m(mkt, id)
			if errors.Is(err, ErrKasusBacaSaja) || errors.Is(err, inti.ErrTanpaWewenang) {
				t.Errorf("%s %s sebelum Binding ditolak penjaga: %v", nama, id, err)
			}
			if nama == "scoring" && !errors.Is(err, ErrScoringBacaSaja) && !errors.Is(err, ErrScoringTanpaDatabase) {
				t.Errorf("scoring %s: %v, mau tetap baca-saja (CallDualScoringRisk)", id, err)
			}
			if err := m(tanpaWB, id); tabIT[nama] == errors.Is(err, inti.ErrTanpaWewenang) {
				t.Errorf("%s %s tanpa workbasket: %v", nama, id, err)
			}
		}
		for _, id := range []string{"UJI-BND1", "UJI-TLB"} {
			if err := m(mkt, id); !errors.Is(err, ErrKasusBacaSaja) {
				t.Errorf("%s %s: %v, mau 409 baca-saja", nama, id, err)
			}
		}
	}
	for id, mau := range map[string]bool{"UJI-BND": true, "UJI-BND0": true, "UJI-BND1": false} {
		boleh, alasan := svc.AksesUbahKasus(ctx, mkt, id)
		if boleh != mau || (!mau && alasan != "Case baca-saja: tombol Binding sudah ditekan.") {
			t.Errorf("AksesUbahKasus %s: %v %q", id, boleh, alasan)
		}
	}
	if boleh, _ := svc.AksesSpreading(ctx, mkt, "UJI-BND"); !boleh {
		t.Error("spreading tahap Binding sebelum tombol harus dapat diubah")
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
		{mkt, "UJI-UW", false, `Case baca-saja: sudah tidak berada di tahap marketing / Team Leader (posisi "ReasFacInJuniorUnderwritingB").`},
		{mkt, "UJI-TL", false, "Akun tidak memegang workbasket ReasFacInTeamLeader."},
		{inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInTeamLeader"}}, "UJI-TL", true, ""},
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
	// Satu sumber aturan: Submit dan Scoring memakai tahapMarketing. Pengecualian tiket 64 lanjutan: tahap Binding - Submit
	// boleh sebelum DAN sesudah tombol Binding (Accept mensyaratkannya), pengubahan hanya sebelumnya.
	semua := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing", "ReasFacInTeamLeader"}}
	for id, k := range tahapUji {
		_, err := alasanTidakBoleh(models.KonteksSubmit{Position: k.Position, FlagOnGoingPolicy: k.FlagOnGoingPolicy}, semua)
		bolehSubmit := tahapDari(k) != "" || tahapSubmitBinding(k) || tahapTLPolicy(k) // K69-8: TL Policy Submit saja
		if (err == nil) != bolehSubmit || bolehUbahScoring(models.KasusScoring{Position: k.Position, FlagOnGoingPolicy: k.FlagOnGoingPolicy}) != tahapMarketing(k) {
			t.Errorf("%s: submit %v / scoring tidak sama dengan tahapDari / tahapMarketing", id, err)
		}
	}
	for id, mau := range map[string]string{"UJI-MKT": "marketing", "UJI-MKT2": "marketing", "UJI-TL": "teamLeader", "UJI-TL1": "",
		"UJI-UW": "", "UJI-DEC": "", "UJI-F1": "", "UJI-BND": "binding", "UJI-BND0": "binding", "UJI-BND1": "", "UJI-TLB": ""} {
		if got := tahapDari(tahapUji[id]); got != mau {
			t.Errorf("tahap %s = %q, mau %q", id, got, mau)
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
