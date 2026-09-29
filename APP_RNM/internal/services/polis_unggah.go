package services

// Unggahan CSV premium list - tiket 04 PremiumList Life.
//
// Untuk apa berkas ini: membaca berkas CSV menjadi baris bernama kolom,
// menyerahkannya ke mesin validasi murni, dan - hanya bila nol penolakan -
// menyimpannya ke tabel peserta dalam SATU transaksi.
//
// ⛔ NOL TABEL STAGING BARU, dan sebabnya dinyatakan. Tiket menyebut
// `POOLDATA.M_TEMPUPLOADLIFE` sebagai "staging nyata"; pembacaan rule-nya
// menunjukkan ia hanya TUJUH kolom - `INSURED, DOB, BEGINDATE, ENDDATE,
// POLIVYHOLDER, NO, CEDING_RETENTION` (`InsertDataUploadLife.xml`) - dan
// dipakai SEMATA untuk pemeriksaan duplikat `CekDoubleInsured`. Baris
// lengkapnya di Pega tinggal di halaman kerja, bukan di tabel itu.
//
// Karena itu tinjauan di sini tidak menulis apa pun: berkasnya diurai,
// divalidasi, dan hasilnya dikembalikan. Yang permanen baru tersentuh saat
// pemakai menekan simpan, dan saat itu berkasnya diurai dan divalidasi ULANG.
// Tabel staging penuh adalah keputusan SKEMA, dan migrasi baru hanya dari
// keputusan yang tercatat.
//
// ⛔ BERKASNYA DIVALIDASI ULANG SAAT SIMPAN, bukan dipercaya dari tinjauan.
// Klien yang dapat melewatkan tinjauan adalah klien yang dapat menyimpan
// apa saja.
//
// Dibaca sesudah: models/polis_unggah.go (aturannya).

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// BatasBarisUnggah membatasi cacah baris satu berkas.
//
// ⛔ Berkas tanpa batas atas adalah permintaan yang memuat seluruh isinya ke
// memori satu proses. Satu polis grup memang dapat memuat ribuan peserta -
// karena itu batasnya longgar, tetapi ADA.
const BatasBarisUnggah = 20000

// BatasUkuranUnggahCSV membatasi ukuran berkasnya.
const BatasUkuranUnggahCSV = 25 << 20 // 25 MiB

var (
	// ErrCSVKosong - berkasnya tidak punya satu pun baris data.
	ErrCSVKosong = errors.New(
		"services: berkas CSV tidak punya baris data di bawah baris judul")
	// ErrCSVTerlaluBanyakBaris - melewati BatasBarisUnggah.
	ErrCSVTerlaluBanyakBaris = errors.New(
		"services: berkas CSV melebihi batas cacah baris")
	// ErrCSVJudulGanda - dua kolom bernama sama.
	//
	// ⛔ Kolom judul ganda membuat nilai yang terbaca bergantung urutan, dan
	// yang belakangan menimpa yang duluan TANPA jejak.
	ErrCSVJudulGanda = errors.New("services: berkas CSV memuat judul kolom ganda")
	// ErrCSVKolomKurang - berkasnya tidak memuat kolom yang divalidasi.
	ErrCSVKolomKurang = errors.New("services: berkas CSV kehilangan kolom wajib")
	// ErrUnggahanDitolak - ada penolakan; tidak ada yang disimpan.
	ErrUnggahanDitolak = errors.New(
		"services: unggahan masih punya penolakan; nol baris disimpan")
)

// tandaUrutanBita adalah BOM UTF-8 yang kerap ditempelkan Excel di depan
// judul kolom PERTAMA.
//
// ⛔ DIBUANG, bukan diabaikan. Tanpa ini, judul pertama terbaca sebagai
// "CERTIFICATE_NO" dan kolomnya dianggap HILANG - sehingga SETIAP
// baris berkas yang benar ditolak, dengan pesan yang menuduh kolom yang
// jelas-jelas terlihat ada di berkasnya. Itu jenis penolakan yang membuat
// orang berhenti memercayai validasinya.
//
// ⚠️ Ditulis sebagai rune, bukan sebagai harfiah di dalam teks: BOM harfiah
// di tengah berkas Go ditolak pengurainya sendiri.
var tandaUrutanBita = string(rune(0xFEFF))

// BacaCSVUnggah mengurai berkas CSV menjadi baris bernama kolom.
//
// ⛔ Baris judul WAJIB, dan namanya dinaikkan menjadi huruf besar - kolom CSV
// dari mesin yang berbeda kerap berbeda huruf besarnya, dan menolak berkas
// karena `sum_insured` bukan `SUM_INSURED` adalah penolakan yang tidak
// mengajari apa pun.
//
// ⚠️ Nomor baris yang dibawa adalah nomor BARIS DATA - 1 untuk baris pertama
// di bawah judul. Itu yang dilihat orang saat membuka berkasnya di penyunting
// setelah menyembunyikan judulnya; nomor baris berkas mentah akan meleset
// satu, dan meleset satu di berkas seribu baris lebih buruk daripada tidak
// ada nomor sama sekali.
func BacaCSVUnggah(r io.Reader) ([]models.BarisUnggah, error) {
	c := csv.NewReader(io.LimitReader(r, BatasUkuranUnggahCSV))
	// ⛔ Cacah medan TIDAK dipatok: baris yang kependekan atau kepanjangan
	// dilaporkan sebagai penolakan kolomnya sendiri, bukan sebagai galat
	// berkas yang membatalkan seluruhnya.
	c.FieldsPerRecord = -1
	c.TrimLeadingSpace = true

	judul, err := c.Read()
	if err == io.EOF {
		return nil, ErrCSVKosong
	}
	if err != nil {
		return nil, fmt.Errorf("services: membaca judul CSV: %w", err)
	}
	nama := make([]string, len(judul))
	lihat := map[string]bool{}
	for i, j := range judul {
		k := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(j, tandaUrutanBita)))
		if k != "" && lihat[k] {
			return nil, fmt.Errorf("%w: %q", ErrCSVJudulGanda, k)
		}
		lihat[k] = true
		nama[i] = k
	}

	// ⛔ JUDUL DIPERIKSA LEBIH DAHULU. Kolom yang HILANG SAMA SEKALI dari
	// berkas akan menghasilkan satu penolakan "kolom kosong" untuk setiap
	// barisnya - seribu kalimat yang mengatakan satu hal, dan orang harus
	// membaca seribu untuk menemukan satu.
	var hilang []string
	for _, k := range models.KolomWajibUnggah() {
		if !lihat[k] {
			hilang = append(hilang, k)
		}
	}
	if len(hilang) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrCSVKolomKurang, strings.Join(hilang, ", "))
	}

	var baris []models.BarisUnggah
	for {
		rec, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("services: membaca baris %d CSV: %w",
				len(baris)+1, err)
		}
		if barisKosong(rec) {
			continue
		}
		if len(baris) >= BatasBarisUnggah {
			return nil, fmt.Errorf("%w: %d", ErrCSVTerlaluBanyakBaris, BatasBarisUnggah)
		}
		nilai := make(map[string]string, len(nama))
		for i, k := range nama {
			if k == "" || i >= len(rec) {
				continue
			}
			nilai[k] = rec[i]
		}
		baris = append(baris, models.BarisUnggah{Nomor: len(baris) + 1, Nilai: nilai})
	}
	if len(baris) == 0 {
		return nil, ErrCSVKosong
	}
	return baris, nil
}

// barisKosong menjawab apakah seluruh selnya kosong.
//
// ⚠️ Baris kosong di ujung berkas adalah hal biasa - penyunting menambahkan
// baris baru di akhir. Menolaknya berarti menolak berkas yang benar.
func barisKosong(rec []string) bool {
	for _, s := range rec {
		if strings.TrimSpace(s) != "" {
			return false
		}
	}
	return true
}

// UnggahPremiumList melayani tinjauan dan penyimpanan unggahan.
type UnggahPremiumList struct{ svc *Service }

// UnggahPremiumList menyusun layanannya.
func (s *Service) UnggahPremiumList() *UnggahPremiumList {
	return &UnggahPremiumList{svc: s}
}

// HasilTinjauUnggah adalah jawaban tinjauan.
type HasilTinjauUnggah struct {
	CacahBaris   int                `json:"cacahBaris"`
	CacahDitolak int                `json:"cacahDitolak"`
	Ditolak      []models.Penolakan `json:"ditolak"`
	Lolos        bool               `json:"lolos"`
}

// Tinjau mengurai dan memvalidasi berkas TANPA menyentuh apa pun.
func (u *UnggahPremiumList) Tinjau(ctx context.Context, pelaku Pelaku,
	polisID string, berkas io.Reader) (HasilTinjauUnggah, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return HasilTinjauUnggah{}, err
	}
	if polisID == "" {
		return HasilTinjauUnggah{}, fmt.Errorf("%w: id polis kosong", ErrPermintaanTidakSah)
	}
	_, hasil, err := periksaBerkas(berkas)
	if err != nil {
		return HasilTinjauUnggah{}, err
	}
	return HasilTinjauUnggah{
		CacahBaris:   hasil.CacahBaris,
		CacahDitolak: len(hasil.Ditolak),
		Ditolak:      hasil.Ditolak,
		Lolos:        hasil.Lolos(),
	}, nil
}

// periksaBerkas mengurai lalu memvalidasi - satu jalan, dipakai dua pintu.
//
// ⛔ SATU JALAN, dan itu inti gerbangnya. Tinjauan dan penyimpanan yang
// memakai pengurai atau validator berbeda akan berselisih, dan selisihnya
// berarti berkas yang lolos tinjauan lalu ditolak saat simpan - atau, jauh
// lebih buruk, sebaliknya.
func periksaBerkas(berkas io.Reader) ([]models.BarisUnggah, models.HasilUnggah, error) {
	if berkas == nil {
		// ⛔ Galat permintaan, bukan penolakan korpus: "Please Upload CSV File
		// into attachment" milik langkah 11/33 yang ter-remark (sensus
		// 28-09-2026). Handler sudah menjawab 400 lebih dulu.
		return nil, models.HasilUnggah{}, fmt.Errorf("%w: berkas CSV tidak disertakan",
			ErrPermintaanTidakSah)
	}
	baris, err := BacaCSVUnggah(berkas)
	if err != nil {
		return nil, models.HasilUnggah{}, err
	}
	// OQ-PL-12 (GILIRAN-17): langkah 2 - uang kosong = 0 SEBELUM validasi.
	models.IsiNolUangKosong(baris)
	return baris, models.ValidasiUnggah(baris), nil
}

// HasilSimpanUnggah adalah jawaban penyimpanan.
type HasilSimpanUnggah struct {
	CacahBaris    int                `json:"cacahBaris"`
	CacahDisimpan int                `json:"cacahDisimpan"`
	CacahDihapus  int                `json:"cacahDihapus"`
	Ditolak       []models.Penolakan `json:"ditolak"`
}

// Simpan menyimpan unggahan ke tabel peserta - hanya bila NOL penolakan.
//
// ⛔ MENGGANTI, BUKAN MENUMPUK (AC tiket 04). Baris peserta polis ini dihapus
// lebih dahulu di dalam transaksi yang sama. Unggahan kedua yang menumpuk
// menghasilkan peserta ganda yang tidak seorang pun minta, dan nomor
// sertifikat yang tadinya unik menjadi tidak.
//
// ⛔ SATU TRANSAKSI. Hapus dan sisip yang terpisah meninggalkan polis TANPA
// peserta sama sekali bila yang kedua gagal - keadaan yang lebih buruk
// daripada keduanya tidak pernah berjalan.
func (u *UnggahPremiumList) Simpan(ctx context.Context, pelaku Pelaku,
	polisID string, berkas io.Reader) (HasilSimpanUnggah, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return HasilSimpanUnggah{}, err
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return HasilSimpanUnggah{}, repository.ErrTanpaOracle
	}
	if polisID == "" {
		return HasilSimpanUnggah{}, fmt.Errorf("%w: id polis kosong", ErrPermintaanTidakSah)
	}

	// ⛔ GERBANG KASUS TERTUTUP - butir bb, lewat `T_WORK_POLIS`.
	keadaanKerja, err := repository.NewWorkPolis(u.svc.db).Keadaan(ctx, polisID)
	if err != nil {
		return HasilSimpanUnggah{}, err
	}
	if models.KasusPolisTertutup(keadaanKerja.Status) {
		return HasilSimpanUnggah{}, fmt.Errorf("%w: polis %q berstatus %q",
			ErrKasusPolisTertutup, polisID, keadaanKerja.Status)
	}

	// ⛔ DIVALIDASI ULANG. Tinjauan yang lolos bukan izin menyimpan.
	//
	// ⚠️ Barisnya diambil dari pengurai yang SAMA, bukan dibaca ulang dari
	// `berkas`: aliran itu sudah habis terbaca, dan membacanya dua kali
	// menuntut menahan seluruh berkas di memori.
	baris, hasil, err := periksaBerkas(berkas)
	if err != nil {
		return HasilSimpanUnggah{}, err
	}
	if !hasil.Lolos() {
		return HasilSimpanUnggah{
			CacahBaris: hasil.CacahBaris,
			Ditolak:    hasil.Ditolak,
		}, ErrUnggahanDitolak
	}

	var keluar HasilSimpanUnggah
	peserta := repository.NewPesertaUnggah(u.svc.db)
	err = u.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		dihapus, err := peserta.HapusPesertaPolis(ctx, tx, polisID)
		if err != nil {
			return err
		}
		disimpan, err := peserta.SisipPeserta(ctx, tx, polisID, baris)
		if err != nil {
			return err
		}
		keluar = HasilSimpanUnggah{
			CacahBaris:    hasil.CacahBaris,
			CacahDisimpan: disimpan,
			CacahDihapus:  dihapus,
			Ditolak:       []models.Penolakan{},
		}
		return nil
	})
	if err != nil {
		return HasilSimpanUnggah{}, err
	}
	return keluar, nil
}
