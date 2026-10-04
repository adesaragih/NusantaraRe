package services_test

// Copy Old TANPA Oracle: hanya superadmin (pemegang menu Kelola User), hasil per organisasi.

import (
	"context"
	"errors"
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/companydetail/backend/models"
	"nusantarare/modul/companydetail/backend/services"
	"nusantarare/modul/companydetail/backend/tiruan"
)

func ctxMenu(kode ...string) context.Context { return inti.DenganAksesMenu(context.Background(), kode) }

func gudangLama() *tiruan.Gudang {
	g := tiruan.Contoh()
	g.Lama = []models.OrgLama{
		{ID: "ASM-SFAGIS-WORK-ORG ORG-118", IDView: "ORG-118", Nama: "UJI LAMA SATU", Baru: true, PIC: 1, Isi: []string{}},
		{ID: "ASM-SFAGIS-WORK-ORG ORG-119", IDView: "ORG-119", Nama: "UJI LAMA DUA", Nomor: 2, Isi: []string{"NOTE"}},
	}
	g.GagalSalin = map[string]bool{"ASM-SFAGIS-WORK-ORG ORG-119": true}
	return g
}

func TestCopyOldHanyaSuperadmin(t *testing.T) {
	l := layanan(gudangLama())
	if h := l.HakAkun(ctxMenu("companydetail")); h.CopyOld {
		t.Error("akun tanpa Kelola User mendapat Copy Old")
	}
	if h := l.HakAkun(context.Background()); h.CopyOld {
		t.Error("permintaan tanpa sesi login mendapat Copy Old")
	}
	if h := l.HakAkun(ctxMenu("companydetail", "kelolauser")); !h.CopyOld {
		t.Error("pemegang Kelola User tidak mendapat Copy Old")
	}
	if _, err := l.DaftarLama(ctxMenu("companydetail"), admin); !errors.Is(err, services.ErrBukanSuperadmin) {
		t.Errorf("daftar lama bukan superadmin: %v", err)
	}
	if _, err := l.SalinLama(ctxMenu("companydetail"), admin, []string{"x"}); !errors.Is(err, services.ErrBukanSuperadmin) {
		t.Errorf("salin lama bukan superadmin: %v", err)
	}
	if _, err := l.DaftarLama(ctxMenu("kelolauser"), inti.Pelaku{}); !errors.Is(err, services.ErrTanpaPelaku) {
		t.Errorf("tanpa pelaku: %v", err)
	}
}

// Process Copy: satu per satu, hasil per ID - disalin, sudah ada, bukan organisasi lama, gagal.
func TestCopyOldHasilPerOrganisasi(t *testing.T) {
	g := gudangLama()
	l := layanan(g)
	ctxAdmin := ctxMenu("kelolauser")
	d, err := l.DaftarLama(ctxAdmin, admin)
	if err != nil || len(d) != 2 {
		t.Fatalf("daftar lama %v %v", d, err)
	}
	j, err := l.SalinLama(ctxAdmin, admin, []string{" ASM-SFAGIS-WORK-ORG ORG-118", "ASM-SFAGIS-WORK-ORG ORG-118",
		anak, "UJI-ASING", "ASM-SFAGIS-WORK-ORG ORG-119"})
	if err != nil {
		t.Fatal(err)
	}
	status := []string{}
	for _, h := range j.Hasil {
		status = append(status, h.Status)
	}
	mau := []string{models.SalinDisalin, models.SalinSudahAda, models.SalinDitolak, models.SalinGagal}
	if !reflect.DeepEqual(status, mau) || j.Disalin != 1 {
		t.Errorf("hasil %v disalin %d, mau %v disalin 1", status, j.Disalin, mau)
	}
	if len(g.Lama) != 1 {
		t.Errorf("yang tersalin masih di popup: %+v", g.Lama)
	}
	if _, err := l.SalinLama(ctxAdmin, admin, []string{" "}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("tanpa pilihan: %v", err)
	}
}
