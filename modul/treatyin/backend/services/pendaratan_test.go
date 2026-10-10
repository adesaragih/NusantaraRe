package services_test

// Uji jalur baca tab dari tabel PENDARATAN - nol koneksi Oracle.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// ⭐ Ketiga tab datang dari GUDANG, bukan dari dokumen. Uji ini sengaja
// memberi kontrak yang `AdaDiJSON`-nya penuh dan lariknya kosong: kalau
// isinya ternyata muncul dari CLOB, ketiganya akan kosong di sini.
func TestTigaTabDibacaDariTabelBukanDariDokumen(t *testing.T) {
	g := &gudangTiruan{
		kontrakWarisan: kontrakWarisanUji(),
		periodePelaporan: []models.BarisPeriodeWarisan{
			{Periode: "Q 1", TanggalAwal: "20181231T170000.000 GMT"},
			{Periode: "Q 2"},
		},
		portofolio: []models.BarisPortofolioWarisan{
			{Jenis: "Premium", JenisPortfolio: "Withdrawal"},
		},
		akumulasi: []models.BarisAkumulasiWarisan{
			{Periode: "1", HariKirim: "60"},
		},
	}
	l := services.LayananDengan(g)

	k, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if len(k.PeriodePelaporan) != 2 {
		t.Errorf("periode pelaporan %d baris, mau 2", len(k.PeriodePelaporan))
	}
	if len(k.Portofolio) != 1 {
		t.Errorf("portofolio %d baris, mau 1", len(k.Portofolio))
	}
	if len(k.Akumulasi) != 1 {
		t.Errorf("akumulasi %d baris, mau 1", len(k.Akumulasi))
	}
	// ⛔ URUTANNYA UTUH. Larik Pega berurut, dan services tidak boleh
	// mengurutkan ulang apa pun.
	if k.PeriodePelaporan[0].Periode != "Q 1" || k.PeriodePelaporan[1].Periode != "Q 2" {
		t.Errorf("urutan periode berubah: %v", k.PeriodePelaporan)
	}
}

// Tab KOSONG bukan galat - 1.836 dari 1.854 kontrak tidak punya akumulasi,
// dan menolak membukanya akan menolak hampir seluruh tabel.
func TestTabKosongBukanGalat(t *testing.T) {
	g := &gudangTiruan{kontrakWarisan: kontrakWarisanUji()}
	l := services.LayananDengan(g)

	k, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if err != nil {
		t.Fatalf("tab kosong menghasilkan galat %v; ia keadaan yang sah", err)
	}
	if len(k.PeriodePelaporan) != 0 || len(k.Portofolio) != 0 || len(k.Akumulasi) != 0 {
		t.Errorf("tab terisi dari suatu tempat: %+v", k)
	}
}

// ⛔ Galat BACA TAB diteruskan, bukan ditelan menjadi tab kosong.
//
// Tab kosong dan tab yang gagal dibaca terlihat sama di layar, dan bedanya
// menentukan apakah seseorang menyelidiki. Yang pertama keadaan sah; yang
// kedua basis data yang tidak menjawab.
func TestGalatBacaTabDiteruskan(t *testing.T) {
	rusak := errors.New("oracle tidak menjawab")
	g := &gudangTiruan{kontrakWarisan: kontrakWarisanUji(), galatTab: rusak}
	l := services.LayananDengan(g)

	_, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if !errors.Is(err, rusak) {
		t.Fatalf("galat baca tab menjadi %v; ia harus sampai ke pemanggil", err)
	}
	// Dan ia BUKAN "kontraknya tidak ada" — kedua pembedanya tidak boleh
	// tertukar: yang satu 404, yang lain 503.
	if services.WarisanTidakAda(err) {
		t.Error("galat baca tab terbaca sebagai kontrak tidak ada")
	}
}
