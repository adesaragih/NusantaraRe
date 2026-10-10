package services_test

// Uji tiket 42 — arsip muatan keluar.
//
// ⛔ Uji terpenting di berkas ini bukan yang memeriksa penyimpanan, melainkan
// `TestBuktiArsipTidakDapatMembawaMuatan`: ia memeriksa BENTUK nilai baliknya,
// dan karena itu ia gagal pada hari seseorang menambahkan ruas muatan — bukan
// pada hari seseorang mulai memakainya.

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func arsipSah() models.ArsipMuatanKeluar {
	return models.ArsipMuatanKeluar{
		IDKontrak:   7,
		Tujuan:      "PEGA_TREATY_IN",
		DikirimPada: "2026-10-02T09:15:00Z",
		Berhasil:    true,
		Muatan:      `{"TreatyIn":{"ID":"TR-1"}}`,
	}
}

func TestArsipMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	if err := l.SimpanArsipMuatanKeluar(context.Background(), inti.Pelaku{}, arsipSah()); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("simpan: mau ErrTanpaIdentitas, dapat %v", err)
	}
	if _, err := l.BuktiArsipKontrak(context.Background(), inti.Pelaku{}, 7); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("bukti: mau ErrTanpaIdentitas, dapat %v", err)
	}
	if len(g.arsip) != 0 || len(g.dibaca) != 0 {
		t.Errorf("gudang tersentuh walau identitas tidak ada")
	}
}

func TestArsipMenolakMasukanTidakSah(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(*models.ArsipMuatanKeluar)
		petik string
	}{
		{"tujuan di luar keempat tujuan hilir", func(a *models.ArsipMuatanKeluar) { a.Tujuan = "GUDANG_LAIN" }, "GUDANG_LAIN"},
		{"tujuan kosong", func(a *models.ArsipMuatanKeluar) { a.Tujuan = "" }, "tujuan"},
		{"muatan kosong", func(a *models.ArsipMuatanKeluar) { a.Muatan = "   " }, "muatan kosong"},
		{"waktu kirim kosong", func(a *models.ArsipMuatanKeluar) { a.DikirimPada = "" }, "waktu KIRIM"},
		{"waktu kirim bukan RFC 3339", func(a *models.ArsipMuatanKeluar) { a.DikirimPada = "02-10-2026" }, "RFC 3339"},
		{"pengenal kontrak nol", func(a *models.ArsipMuatanKeluar) { a.IDKontrak = 0 }, "pengenal kontrak"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			g := &gudangTiruan{}
			l := services.LayananDengan(g)
			a := arsipSah()
			k.ubah(&a)

			err := l.SimpanArsipMuatanKeluar(context.Background(), pelakuAda, a)
			if !errors.Is(err, services.ErrMasukanTidakSah) {
				t.Fatalf("mau ErrMasukanTidakSah, dapat %v", err)
			}
			if !strings.Contains(err.Error(), k.petik) {
				t.Errorf("pesan tidak menyebut %q: %v", k.petik, err)
			}
			if len(g.arsip) != 0 {
				t.Errorf("masukan yang ditolak tetap sampai ke gudang")
			}
		})
	}
}

// Uji POSITIF, dan ia yang menangkap arsip yang menyaring isinya: muatan yang
// TIDAK dapat diurai sebagai JSON **tetap tersimpan**.
//
// Arsip yang hanya menerima yang sudah benar tidak melestarikan apa pun yang
// berguna — justru muatan yang rusak itulah yang paling perlu ditunjukkan.
func TestMuatanRusakTetapDiarsipkan(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	a := arsipSah()
	a.Muatan = `{"TreatyIn": <<rusak di tengah`
	if err := l.SimpanArsipMuatanKeluar(context.Background(), pelakuAda, a); err != nil {
		t.Fatalf("muatan rusak ditolak: %v", err)
	}
	if len(g.arsip) != 1 || g.arsip[0].Muatan != a.Muatan {
		t.Errorf("muatan tidak sampai apa adanya: %+v", g.arsip)
	}
}

// Positif kedua: pengiriman yang GAGAL di hilir tetap terarsip.
func TestPengirimanGagalTetapDiarsipkan(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	a := arsipSah()
	a.Berhasil = false
	if err := l.SimpanArsipMuatanKeluar(context.Background(), pelakuAda, a); err != nil {
		t.Fatalf("pengiriman gagal ditolak: %v", err)
	}
	if len(g.arsip) != 1 || g.arsip[0].Berhasil {
		t.Errorf("penanda gagal tidak sampai: %+v", g.arsip)
	}
}

// Kontrak tanpa arsip menjawab cacah nol — JAWABAN, bukan galat.
func TestKontrakTanpaArsipMenjawabNol(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	b, err := l.BuktiArsipKontrak(context.Background(), pelakuAda, 7)
	if err != nil {
		t.Fatalf("kontrak tanpa arsip menghasilkan galat %v; ia keadaan yang sah", err)
	}
	if b.Cacah != 0 || b.TerakhirDikirim != "" {
		t.Errorf("bukti kosong tidak kosong: %+v", b)
	}
	// Daftar kosong, bukan nil — pemanggil JSON tidak perlu membedakan
	// `null` dari `[]` untuk pertanyaan yang jawabannya "belum ada".
	if b.Tujuan == nil || len(b.Tujuan) != 0 {
		t.Errorf("Tujuan %v, mau daftar kosong", b.Tujuan)
	}
}

func TestBuktiArsipMenjawabCacahDanTujuan(t *testing.T) {
	g := &gudangTiruan{bukti: models.BuktiArsip{
		Cacah: 4, TerakhirDikirim: "2026-10-02T09:15:00Z",
		Tujuan: []string{"PEGA_M_TREATY_IN_DETAIL", "PEGA_TREATY_IN"},
	}}
	l := services.LayananDengan(g)

	b, err := l.BuktiArsipKontrak(context.Background(), pelakuAda, 7)
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if b.Cacah != 4 || b.IDKontrak != 7 || len(b.Tujuan) != 2 {
		t.Errorf("bukti: %+v", b)
	}
}

// ⛔ PENJAGA BENTUK — `INV-61` diwujudkan, bukan sekadar ditulis.
//
// `models.BuktiArsip` adalah SATU-SATUNYA bentuk yang jalur baca arsip
// kembalikan, dan ia tidak boleh punya ruas yang membawa isi arsip. Uji ini
// gagal pada hari seseorang menambahkannya — bukan pada hari seseorang mulai
// memakainya, yang bisa setahun kemudian dan di modul lain.
func TestBuktiArsipTidakDapatMembawaMuatan(t *testing.T) {
	tipe := reflect.TypeOf(models.BuktiArsip{})
	for i := 0; i < tipe.NumField(); i++ {
		nama := strings.ToUpper(tipe.Field(i).Name)
		for _, terlarang := range []string{"MUATAN", "ISI", "PAYLOAD", "JSON", "DOKUMEN"} {
			if strings.Contains(nama, terlarang) {
				t.Errorf("models.BuktiArsip punya ruas %q; ADR-0034 menyatakan arsip "+
					"TIDAK punya jalur baca, dan bentuk ini yang mewujudkannya",
					tipe.Field(i).Name)
			}
		}
	}
	// Dan ia harus benar-benar punya ruas — struct kosong lolos uji di atas
	// tanpa menjaga apa pun.
	if tipe.NumField() < 3 {
		t.Fatalf("models.BuktiArsip hanya %d ruas; penjaganya yang rusak", tipe.NumField())
	}
}
