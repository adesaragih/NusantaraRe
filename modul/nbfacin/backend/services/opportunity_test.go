package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// caseNBTiruan - PenulisCaseNB tiruan: mencatat langkah dan argumennya.
type caseNBTiruan struct {
	langkah   *[]string
	pembuat   *string
	o         *models.Opportunity
	galatPada string
}

func (c caseNBTiruan) catat(l string) error {
	*c.langkah = append(*c.langkah, l)
	if l == c.galatPada {
		return errors.New("ORA-UJI " + l)
	}
	return nil
}

func (c caseNBTiruan) PengenalBerikut(context.Context, *db.Tx) (string, error) {
	return "NB-900001", c.catat("nomor")
}

func (c caseNBTiruan) SisipCase(_ context.Context, _ *db.Tx, id, pembuat string) error {
	*c.pembuat = pembuat
	return c.catat("case " + id)
}

func (c caseNBTiruan) SisipOpportunity(_ context.Context, _ *db.Tx, id string, o models.Opportunity) error {
	*c.o = o
	return c.catat("opportunity " + id)
}

// transaksiTiruan - menjalankan fn tanpa Oracle dan mencatat apakah ia dibatalkan.
func transaksiTiruan(batal *bool) Transaksi {
	return func(_ context.Context, fn func(*db.Tx) error) error {
		err := fn(nil)
		*batal = err != nil
		return err
	}
}

func isianSah() IsianOpportunity {
	return IsianOpportunity{EstimatedClosingDate: "31-10-2026", BusinessProspectName: " UJI Prospek ", AccountID: "UJI-AKUN",
		InsuredID: "UJI-INS", GroupBusinessID: "UJI-GB", GroupBusiness: "UJI GRUP", ClassOfBusiness: "UJI KELAS",
		TypeOfInward: "Facultative", TypeOfFacultative: "Facultative In", Phase: "Proposal", Stage: "Opportunity",
		BusinessStatus: "New Business", Description: "UJI"}
}

// TestBuatOpportunity - tiket 29: satu transaksi berurutan nomor -> case -> opportunity;
// pembuat = pengenal akun; isian disimpan APA ADANYA (tidak dipangkas); tanggal tanpa jam.
func TestBuatOpportunity(t *testing.T) {
	var langkah []string
	var pembuat string
	var o models.Opportunity
	var batal bool
	svc := Baru(nil).DenganCaseNB(caseNBTiruan{langkah: &langkah, pembuat: &pembuat, o: &o}).DenganTransaksi(transaksiTiruan(&batal))
	id, err := svc.BuatOpportunity(context.Background(), inti.Pelaku{AkunID: "UJI-USER"}, isianSah())
	if err != nil || id != "NB-900001" || batal {
		t.Fatalf("id %q, galat %v, batal %v", id, err, batal)
	}
	if strings.Join(langkah, " | ") != "nomor | case NB-900001 | opportunity NB-900001" || pembuat != "UJI-USER" {
		t.Errorf("langkah %v, pembuat %q", langkah, pembuat)
	}
	if !o.EstimatedClosingDate.Equal(time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)) || o.BusinessProspectName != " UJI Prospek " ||
		o.TypeOfFacultative != "Facultative In" || o.Description != "UJI" || o.ClassOfBusiness != "UJI KELAS" {
		t.Errorf("opportunity %+v", o)
	}
}

// TestBuatOpportunityGagal - 401 tanpa identitas, 400 per medan (sebelum basis data),
// 503 tanpa basis data, galat di dalam transaksi membatalkan seluruhnya.
func TestBuatOpportunityGagal(t *testing.T) {
	ctx := context.Background()
	akun := inti.Pelaku{AkunID: "UJI-USER"}
	var langkah []string
	var pembuat string
	var o models.Opportunity
	var batal bool
	tiruan := caseNBTiruan{langkah: &langkah, pembuat: &pembuat, o: &o}
	svc := Baru(nil).DenganCaseNB(tiruan).DenganTransaksi(transaksiTiruan(&batal))

	if _, err := svc.BuatOpportunity(ctx, inti.Pelaku{AkunID: "  "}, isianSah()); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	ubah := func(f func(*IsianOpportunity)) IsianOpportunity { i := isianSah(); f(&i); return i }
	for nama, u := range map[string]struct {
		isian IsianOpportunity
		pesan string
	}{
		"tanggal kosong":    {ubah(func(i *IsianOpportunity) { i.EstimatedClosingDate = "" }), "estimatedClosingDate wajib"},
		"tanggal 31-02":     {ubah(func(i *IsianOpportunity) { i.EstimatedClosingDate = "31-02-2026" }), "estimatedClosingDate bukan"},
		"tanggal ISO":       {ubah(func(i *IsianOpportunity) { i.EstimatedClosingDate = "2026-10-31" }), "estimatedClosingDate bukan"},
		"tanggal tanpa nol": {ubah(func(i *IsianOpportunity) { i.EstimatedClosingDate = "1-10-2026" }), "estimatedClosingDate bukan"},
		"prospek spasi":     {ubah(func(i *IsianOpportunity) { i.BusinessProspectName = "   " }), "businessProspectName wajib"},
		"kelas kosong":      {ubah(func(i *IsianOpportunity) { i.ClassOfBusiness = "" }), "classOfBusiness wajib"},
		"inward kosong":     {ubah(func(i *IsianOpportunity) { i.TypeOfInward = "" }), "typeOfInward wajib"},
		"fakultatif kosong": {ubah(func(i *IsianOpportunity) { i.TypeOfFacultative = "" }), "typeOfFacultative wajib"},
		"phase kosong":      {ubah(func(i *IsianOpportunity) { i.Phase = "" }), "phase wajib"},
		"status kosong":     {ubah(func(i *IsianOpportunity) { i.BusinessStatus = "" }), "businessStatus wajib"},
		"prospek 256 bita":  {ubah(func(i *IsianOpportunity) { i.BusinessProspectName = strings.Repeat("U", 256) }), "businessProspectName paling banyak 255 byte"},
		"grup 33 karakter":  {ubah(func(i *IsianOpportunity) { i.GroupBusinessID = strings.Repeat("é", 33) }), "groupBusinessId paling banyak 32 karakter"},
		"keterangan > 4000": {ubah(func(i *IsianOpportunity) { i.Description = strings.Repeat("U", 4001) }), "description paling banyak 4000 byte"},
	} {
		langkah = nil
		_, err := svc.BuatOpportunity(ctx, akun, u.isian)
		if !errors.Is(err, ErrMasukanOpportunity) || !strings.Contains(err.Error(), u.pesan) || len(langkah) != 0 {
			t.Errorf("%s: %v (langkah %v), mau 400 %q", nama, err, langkah, u.pesan)
		}
		if _, err := Baru(nil).BuatOpportunity(ctx, akun, u.isian); !errors.Is(err, ErrMasukanOpportunity) {
			t.Errorf("%s tanpa DB: %v, mau 400 dulu", nama, err)
		}
	}
	// Batas tepat lolos: 32 karakter dua-bita = 64 bita (semantik CHAR), 255 bita.
	boleh := ubah(func(i *IsianOpportunity) {
		i.GroupBusinessID, i.BusinessProspectName = strings.Repeat("é", 32), strings.Repeat("U", 255)
	})
	if _, err := svc.BuatOpportunity(ctx, akun, boleh); err != nil {
		t.Errorf("tepat di batas: %v", err)
	}
	// Bukan Facultative: Type Of Facultative tidak wajib.
	if _, err := svc.BuatOpportunity(ctx, akun, ubah(func(i *IsianOpportunity) { i.TypeOfInward, i.TypeOfFacultative = "UJI LAIN", "" })); err != nil {
		t.Errorf("bukan Facultative: %v", err)
	}
	if _, err := Baru(nil).BuatOpportunity(ctx, akun, isianSah()); !errors.Is(err, ErrOpportunityTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	for _, l := range []string{"nomor", "case NB-900001", "opportunity NB-900001"} {
		langkah = nil
		gagal := tiruan
		gagal.galatPada = l
		id, err := Baru(nil).DenganCaseNB(gagal).DenganTransaksi(transaksiTiruan(&batal)).BuatOpportunity(ctx, akun, isianSah())
		if err == nil || id != "" || !batal || langkah[len(langkah)-1] != l {
			t.Errorf("galat di %q: id %q, galat %v, batal %v, langkah %v", l, id, err, batal, langkah)
		}
	}
}
