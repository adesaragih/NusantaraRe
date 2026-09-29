package services

// Inbox Komite - tiket 01 Komite Claim Life. TANPA Oracle.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/modul/komite/repository"
)

func kasusUji() repository.KasusKomite {
	return repository.KasusKomite{
		Baris: repository.BarisInboxKomite{KasusID: "KMTLF-UJI1", TingkatBerjalan: 2,
			KomiteLoop: 3, NilaiKlaim: "1500000.25", MataUang: "IDR", StsReject: "0"},
		Tangga: []repository.AnggotaKasus{
			{Urut: 1, OperatorID: "UJI-A", Approval: "1"},
			{Urut: 2, OperatorID: "UJI-B", Approval: "0"},
			{Urut: 3, OperatorID: "UJI-C", Approval: "0"},
		},
	}
}

// TestKasusKomiteHanyaUntukAnggotaTangga - ADR-0014.
func TestKasusKomiteHanyaUntukAnggotaTangga(t *testing.T) {
	if _, err := susunKasusTampil(kasusUji(), "UJI-LUAR"); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("bukan anggota: %v, mau ErrTanpaWewenang", err)
	}
	for akun, giliran := range map[string]bool{"UJI-A": false, "UJI-B": true, "UJI-C": false} {
		k, err := susunKasusTampil(kasusUji(), akun)
		if err != nil {
			t.Fatalf("%s: %v", akun, err)
		}
		if k.GiliranSaya != giliran {
			t.Errorf("%s: giliran %v, mau %v - hanya anggota berjalan", akun, k.GiliranSaya, giliran)
		}
	}
}

// TestKasusKomiteTertutupBukanGiliranSiapaPun.
func TestKasusKomiteTertutupBukanGiliranSiapaPun(t *testing.T) {
	k := kasusUji()
	k.Baris.StatusWork = "Resolved-Completed"
	got, err := susunKasusTampil(k, "UJI-B")
	if err != nil {
		t.Fatal(err)
	}
	if got.GiliranSaya {
		t.Error("kasus tertutup masih menjadi giliran anggota berjalan")
	}
}

// TestStatusDanUangSebagaiTeks - AC 17, AC 29 spec.
func TestStatusDanUangSebagaiTeks(t *testing.T) {
	k, _ := susunKasusTampil(kasusUji(), "UJI-B")
	if k.Kasus.StatusBaris != "Outstanding" || k.Kasus.NilaiKlaim != "1500000.25" {
		t.Errorf("status %q nilai %q", k.Kasus.StatusBaris, k.Kasus.NilaiKlaim)
	}
	if KataApprovalKomite("0") != "Menunggu" || KataApprovalKomite("9") != "Kode 9" ||
		KataApprovalKomite("1") != "Setuju" || KataApprovalKomite("") != KataTingkatDilewati {
		t.Errorf("kata approval %q / %q", KataApprovalKomite("0"), KataApprovalKomite("9"))
	}
}

// TestInboxKomiteMenuntutIdentitas.
func TestInboxKomiteMenuntutIdentitas(t *testing.T) {
	i := New(nil).InboxKomite()
	if _, err := i.Ambil(context.Background(), inti.Pelaku{}, 1, 10); err == nil {
		t.Error("inbox tanpa identitas diterima")
	}
	if _, err := i.Ambil(context.Background(), inti.Pelaku{AkunID: "UJI"}, 1, 10); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("inbox tanpa Oracle: %v", err)
	}
}
