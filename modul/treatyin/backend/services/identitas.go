package services

// Identitas kontrak - tiket 16, 17, 18, 19.
//
// Keempatnya perilaku atas `KONTRAK`/`VERSI_KONTRAK` yang sudah berdiri; nol
// tabel baru. Satu seam, empat aturan:
//
//	16  kunci alami kembar MEMPERINGATKAN, tidak menghalangi simpan
//	17  kontrak ditemukan lewat nomor warisan sistem lama
//	18  kelima ruas beku tidak dapat diubah sesudah kontrak lahir
//	19  versi tidak dapat membawa lapisan beku yang menyimpang

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

var (
	// ErrRuasBekuBerubah - INV-19, tiket 18.
	ErrRuasBekuBerubah = errors.New("ruas beku kontrak tidak dapat diubah")
	// ErrVersiMenyimpang - tiket 19.
	ErrVersiMenyimpang = errors.New("versi membawa lapisan beku yang menyimpang dari kontraknya")
)

// KodePeringatanKunciAlami - tiket 16.
const KodePeringatanKunciAlami = "KUNCI_ALAMI_GANDA"

// Peringatan menyertai penyimpanan yang BERHASIL.
//
// ⛔ Bukan galat, dan bedanya pokok tiket 16: `ADR-0040` §2 memutuskan kunci
// alami kontrak MEMPERINGATKAN, tidak melarang. Kontrak kembar memang terjadi -
// pembaruan tahunan yang dicatat ulang, atau dua perjanjian terpisah dengan
// cedant yang sama pada tanggal yang sama. Melarangnya berarti menolak data
// yang sah, dan itu kekeliruan lingkup yang sudah tiga kali terjadi di modul ini.
type Peringatan struct {
	Kode  string `json:"kode"`
	Pesan string `json:"pesan"`
	// Pembanding adalah pengenal kontrak yang kunci alaminya sama.
	//
	// ⛔ WAJIB terisi. Peringatan tanpa menyebut pembandingnya tidak dapat
	// ditindaklanjuti: yang membacanya harus dapat MEMBUKA kontrak yang
	// dimaksud. Daftar periksa tiket 16 menolak yang tanpa ini.
	Pembanding []int64 `json:"pembanding"`
}

// peringatkanKunciAlamiGanda menyusun peringatan tiket 16, atau nil.
func peringatkanKunciAlamiGanda(serupa []int64) *Peringatan {
	if len(serupa) == 0 {
		return nil
	}
	n := make([]string, 0, len(serupa))
	for _, id := range serupa {
		n = append(n, fmt.Sprint(id))
	}
	return &Peringatan{
		Kode: KodePeringatanKunciAlami,
		Pesan: fmt.Sprintf(
			"kunci alami kontrak ini (cedant, asal bisnis, tanggal mulai, sifat proporsi) "+
				"sama dengan kontrak %s; penyimpanan tetap dilanjutkan", strings.Join(n, ", ")),
		Pembanding: serupa,
	}
}

// CariLewatNomorWarisan menemukan kontrak dengan nomor sistem lama - tiket 17.
//
// ⛔ Hasil KOSONG adalah jawaban, bukan galat: orang yang memegang selembar
// kertas harus tahu nomornya tidak ada, bukan menerima daftar penuh atau pesan
// kegagalan. Dan baris yang nomor warisannya kosong TIDAK PERNAH ikut terambil -
// kontrak baru bukan kontrak warisan.
func (l *Layanan) CariLewatNomorWarisan(ctx context.Context, p inti.Pelaku, nomor string) ([]models.Kontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	nomor = strings.TrimSpace(nomor)
	if nomor == "" {
		// Pencarian tanpa kata kunci bukan "cari semua": ia permintaan yang
		// belum lengkap, dan menjawabnya dengan daftar penuh membuat pemanggil
		// mengira ia menemukan sesuatu.
		return nil, fmt.Errorf("%w: nomor warisan wajib diisi", ErrMasukanTidakSah)
	}
	return l.gudang.CariKontrakLewatNomorWarisan(ctx, nomor)
}

// ruasBeku adalah kelima ruas yang ADR-0040 bekukan, beserta pembacanya.
//
// Disebut SEKALI di sini dan dipakai tiket 18 maupun 19 - dua aturan atas
// himpunan ruas yang sama, dan dua daftar yang harus sepakat adalah dua daftar
// yang akan menyimpang.
var ruasBeku = []struct {
	nama  string
	nilai func(models.Kontrak) string
}{
	{"idCedant", func(k models.Kontrak) string { return fmt.Sprint(k.IDCedant) }},
	{"idAsalBisnis", func(k models.Kontrak) string { return fmt.Sprint(k.IDAsalBisnis) }},
	{"sifatProporsi", func(k models.Kontrak) string { return k.SifatProporsi }},
	{"tanggalMulai", func(k models.Kontrak) string { return k.TanggalMulai }},
	{"tanggalBerakhir", func(k models.Kontrak) string { return k.TanggalBerakhir }},
}

// bedaRuasBeku mengembalikan SELURUH ruas beku yang berbeda, bukan yang pertama.
func bedaRuasBeku(berlaku, diminta models.Kontrak) []string {
	var beda []string
	for _, r := range ruasBeku {
		if lama, baru := r.nilai(berlaku), r.nilai(diminta); lama != baru {
			beda = append(beda, fmt.Sprintf("%s (berlaku %q, diminta %q)", r.nama, lama, baru))
		}
	}
	return beda
}

// MasukanUbahKontrak adalah permintaan menyunting kepala kontrak - tiket 18.
//
// Kelima ruas beku IKUT dikirim, dan itu disengaja: borang mengirimkan apa yang
// ada di layar. Yang diperiksa bukan "apakah dikirim" melainkan "apakah
// BERBEDA" - mengirim nilai yang sama dengan yang berlaku bukan percobaan ubah.
type MasukanUbahKontrak struct {
	NomorKontrakWarisan string `json:"nomorKontrakWarisan"`
	IDCedant            int64  `json:"idCedant"`
	IDAsalBisnis        int64  `json:"idAsalBisnis"`
	SifatProporsi       string `json:"sifatProporsi"`
	TanggalMulai        string `json:"tanggalMulai"`
	TanggalBerakhir     string `json:"tanggalBerakhir"`
}

// PerbaruiKontrak menyunting kepala kontrak, dan MENOLAK setiap perubahan atas
// kelima ruas beku - tiket 18, INV-19.
//
// ⛔ Sistem lama membiarkannya: `SaveTreatyInDetail_Act` menulis ulang seluruh
// kepala pada setiap penyimpanan, tanpa satu pun pemeriksaan bahwa yang beku
// tetap beku.
//
// Ruas yang BUKAN kunci alami - `NOMOR_KONTRAK_WARISAN` - tetap dapat diubah.
func (l *Layanan) PerbaruiKontrak(ctx context.Context, p inti.Pelaku, id int64, m MasukanUbahKontrak) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("%w: pengenal kontrak %d", ErrMasukanTidakSah, id)
	}
	ada, err := l.gudang.BacaKontrak(ctx, id)
	if err != nil {
		return err
	}
	diminta := models.Kontrak{
		IDCedant:        m.IDCedant,
		IDAsalBisnis:    m.IDAsalBisnis,
		SifatProporsi:   strings.TrimSpace(m.SifatProporsi),
		TanggalMulai:    strings.TrimSpace(m.TanggalMulai),
		TanggalBerakhir: strings.TrimSpace(m.TanggalBerakhir),
	}
	if beda := bedaRuasBeku(ada.Kontrak, diminta); len(beda) > 0 {
		// Pesannya menyebut SELURUH ruas yang menyimpang, bukan yang pertama:
		// pada borang berkolom puluhan, menebak berarti mencoba satu per satu.
		return fmt.Errorf("%w: %s", ErrRuasBekuBerubah, strings.Join(beda, "; "))
	}
	ubah := ada.Kontrak
	ubah.NomorKontrakWarisan = strings.TrimSpace(m.NomorKontrakWarisan)
	return l.gudang.PerbaruiKontrak(ctx, ubah)
}

// MasukanVersiTambahan adalah versi berikutnya pada kontrak yang sudah ada -
// tiket 19.
//
// Kelima ruas beku dikirim sebagai PERNYATAAN, bukan sebagai kolom: `INV-59`
// melarang menyalin lapisan beku ke `VERSI_KONTRAK`, dan `KAMUS-KOLOM.md` §10.2
// memang tidak punya kolomnya. Yang dilakukan di sini membandingkan pernyataan
// itu dengan kontraknya, lalu membuangnya.
type MasukanVersiTambahan struct {
	NomorUrutVersi     int64  `json:"nomorUrutVersi"`
	KeadaanSiklusHidup string `json:"keadaanSiklusHidup"`
	NamaKontrak        string `json:"namaKontrak"`

	// Pernyataan lapisan beku - dibandingkan, tidak disimpan.
	IDCedant        int64  `json:"idCedant"`
	IDAsalBisnis    int64  `json:"idAsalBisnis"`
	SifatProporsi   string `json:"sifatProporsi"`
	TanggalMulai    string `json:"tanggalMulai"`
	TanggalBerakhir string `json:"tanggalBerakhir"`
}

// TambahVersi menambahkan versi pada kontrak yang sudah ada, dan MENOLAK versi
// yang membawa lapisan beku menyimpang - tiket 19.
//
// ⛔ Kenapa ini penting: pemisahan KONTRAK/VERSI_KONTRAK hanya berarti bila
// lapisan bekunya benar-benar beku. Bila sebuah versi dapat menyatakan cedant
// yang berbeda, pertanyaan "siapa cedant kontrak ini" punya jawaban berbeda-beda
// tergantung versi mana yang dibaca - dan di sistem lama memang begitu: tiap
// addendum menyimpan salinan penuh kepala kontraknya di JSONDATA-nya sendiri.
func (l *Layanan) TambahVersi(ctx context.Context, p inti.Pelaku, idKontrak int64, m MasukanVersiTambahan) (int64, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return 0, err
	}
	if idKontrak <= 0 {
		return 0, fmt.Errorf("%w: pengenal kontrak %d", ErrMasukanTidakSah, idKontrak)
	}
	var salah []string
	if strings.TrimSpace(m.NamaKontrak) == "" {
		salah = append(salah, "namaKontrak wajib diisi")
	}
	if strings.TrimSpace(m.KeadaanSiklusHidup) == "" {
		salah = append(salah, "keadaanSiklusHidup wajib diisi")
	}
	if m.NomorUrutVersi <= 0 {
		salah = append(salah, "nomorUrutVersi wajib bilangan positif")
	}
	if len(salah) > 0 {
		return 0, fmt.Errorf("%w: %s", ErrMasukanTidakSah, strings.Join(salah, "; "))
	}

	ada, err := l.gudang.BacaKontrak(ctx, idKontrak)
	if err != nil {
		return 0, err
	}
	diminta := models.Kontrak{
		IDCedant:        m.IDCedant,
		IDAsalBisnis:    m.IDAsalBisnis,
		SifatProporsi:   strings.TrimSpace(m.SifatProporsi),
		TanggalMulai:    strings.TrimSpace(m.TanggalMulai),
		TanggalBerakhir: strings.TrimSpace(m.TanggalBerakhir),
	}
	if beda := bedaRuasBeku(ada.Kontrak, diminta); len(beda) > 0 {
		return 0, fmt.Errorf("%w: %s", ErrVersiMenyimpang, strings.Join(beda, "; "))
	}

	nomor := m.NomorUrutVersi
	return l.gudang.TambahVersi(ctx, idKontrak, models.VersiKontrak{
		IDKontrak:          idKontrak,
		NomorUrutVersi:     &nomor,
		KeadaanSiklusHidup: strings.TrimSpace(m.KeadaanSiklusHidup),
		NamaKontrak:        strings.TrimSpace(m.NamaKontrak),
	})
}
