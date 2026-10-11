package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 64d (keputusan work owner 08-10-2026): K1 Begin / End date terkunci di tahap Binding ("Ikut Pega"); K2 Share Cedant
// Type wajib di tahap Binding ("Share Cedant Type wajib diisi.").

var (
	tahapBindingUji  = models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2"}
	tahapMarketingUj = models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"}
	pelakuPeriode    = inti.Pelaku{AkunID: "UJI-MKT", Peran: []string{"ReasFacInMarketing"}}
)

// TestKunciPeriode - tahap Binding sebelum DAN sesudah tombol = true; marketing / TL / tanpa pembaca tahap / kasus asing = false.
func TestKunciPeriode(t *testing.T) {
	ctx := context.Background()
	sesudah := tahapBindingUji
	sesudah.ConfirmBinding = "1"
	svc := Baru(nil).DenganTahap(tahapTiruan{"B": tahapBindingUji, "B1": sesudah, "M": tahapMarketingUj,
		"T": {Position: "ReasFacInTeamLeader", FlagOnGoingPolicy: "2"}})
	for id, mau := range map[string]bool{"B": true, "B1": true, "M": false, "T": false, "X": false, "UJI-GALAT": false} {
		if got := svc.KunciPeriode(ctx, id); got != mau {
			t.Errorf("%s: %v", id, got)
		}
	}
	if Baru(nil).KunciPeriode(ctx, "B") {
		t.Error("tanpa pembaca tahap")
	}
}

// TestSimpanGeneralPeriodeBinding - tahap Binding: tanggal sama (bentuk kabel) -> medan lain tersimpan; Begin atau End
// berbeda -> 409 (ErrPeriodeTerkunci membungkus ErrKasusBacaSaja), tersimpan tidak berubah; tahap marketing bebas.
func TestSimpanGeneralPeriodeBinding(t *testing.T) {
	ctx := context.Background()
	jam := func() time.Time { return time.Date(2026, 10, 3, 1, 0, 0, 0, time.UTC) }
	kasus := kasusTiruan{ada: kasusUji()}
	tahap := tahapTiruan{"UJI-NB-1": tahapMarketingUj}
	svc := Baru(nil).DenganKasus(kasus).DenganTransaksi(tanpaTx).DenganJam(jam).DenganTahap(tahap)
	isian := IsianGeneral{QQName: "UJI QQ", BeginDate: "01-11-2026", EndDate: "01-11-2027", Day: "365"}
	if _, err := svc.SimpanGeneral(ctx, pelakuPeriode, "UJI-NB-1", isian); err != nil {
		t.Fatalf("marketing: %v", err)
	}
	tahap["UJI-NB-1"] = tahapBindingUji
	isian.QQName = "UJI QQ BARU"
	k, err := svc.SimpanGeneral(ctx, pelakuPeriode, "UJI-NB-1", isian)
	if err != nil || k.General.QQName != "UJI QQ BARU" || k.General.BeginDate != "01-11-2026" {
		t.Fatalf("binding tanggal sama: %v %+v", err, k.General)
	}
	for nama, ubah := range map[string]func(*IsianGeneral){
		"begin":  func(i *IsianGeneral) { i.BeginDate = "02-11-2026" },
		"end":    func(i *IsianGeneral) { i.EndDate = "02-11-2027" },
		"kosong": func(i *IsianGeneral) { i.BeginDate, i.EndDate = "", "" },
	} {
		x := isian
		x.QQName = "UJI TIDAK TERSIMPAN"
		ubah(&x)
		_, err := svc.SimpanGeneral(ctx, pelakuPeriode, "UJI-NB-1", x)
		if !errors.Is(err, ErrPeriodeTerkunci) || !errors.Is(err, ErrKasusBacaSaja) {
			t.Errorf("%s: %v", nama, err)
		}
		if g := kasus.ada["UJI-NB-1"].General; g.QQName != "UJI QQ BARU" {
			t.Errorf("%s: tersimpan berubah %+v", nama, g)
		}
	}
	// Tahap marketing: tanggal bebas diubah.
	tahap["UJI-NB-1"] = tahapMarketingUj
	isian.BeginDate = "05-11-2026"
	if k, err := svc.SimpanGeneral(ctx, pelakuPeriode, "UJI-NB-1", isian); err != nil || k.General.BeginDate != "05-11-2026" {
		t.Errorf("marketing ubah: %v %+v", err, k.General)
	}
}

// TestWajibTypeCedant - K2: tahap Binding menggerbang Share Cedant Type SAJA (wajib / proteksi C-3 lain tidak menyala bila
// gerbang ProtectShareCedant tidak berlaku); pesan persis; di luar Binding tanpa gerbang tipe kosong boleh.
func TestWajibTypeCedant(t *testing.T) {
	ctx := context.Background()
	c := &cedantTiruan{ada: models.KasusCedant{SobName: "UJI SOB", FlagOnGoingPolicy: "0"}}
	tahap := tahapTiruan{"NB-1": tahapBindingUji}
	svc := layananCedant(c).DenganTahap(tahap)
	got, err := svc.TampilanCedant(ctx, "NB-1")
	if err != nil || got.Wajib || !got.WajibType {
		t.Fatalf("binding: %v wajib %v wajibType %v", err, got.Wajib, got.WajibType)
	}
	_, err = svc.SimpanCedant(ctx, pelakuPeriode, "NB-1", "", nil)
	if !errors.Is(err, ErrMasukanCedant) || !strings.HasSuffix(err.Error(), ": Share Cedant Type wajib diisi.") {
		t.Errorf("binding tipe kosong: %v", err)
	}
	if _, err := svc.SimpanCedant(ctx, pelakuPeriode, "NB-1", "0", nil); err != nil || len(c.tulis) != 1 {
		t.Errorf("binding tipe terisi (tanpa proteksi C-3): %v", err)
	}
	tahap["NB-1"] = tahapMarketingUj
	if got, _ := svc.TampilanCedant(ctx, "NB-1"); got.WajibType {
		t.Error("marketing tanpa gerbang: wajibType")
	}
	if _, err := svc.SimpanCedant(ctx, pelakuPeriode, "NB-1", "", nil); err != nil {
		t.Errorf("marketing tipe kosong: %v", err)
	}
	// Gerbang ProtectShareCedant (FlagOnGoingPolicy 2 di data cedant) -> wajib DAN wajibType.
	c.ada.FlagOnGoingPolicy = "2"
	if got, _ := svc.TampilanCedant(ctx, "NB-1"); !got.Wajib || !got.WajibType {
		t.Errorf("gerbang: %+v", got)
	}
}
