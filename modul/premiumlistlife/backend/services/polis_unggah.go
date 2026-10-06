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
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
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

// ukuranIntipJudul - jendela bita awal berkas yang diintip untuk menebak
// pemisah kolom. Cukup untuk baris judul ratusan kolom.
const ukuranIntipJudul = 64 * 1024

// PemisahCSV menebak pemisah kolom dari BARIS JUDUL: titik koma (`;`) bila
// judulnya memuat lebih banyak `;` daripada `,` di luar tanda kutip, selain
// itu koma.
//
// ⛔ [keputusan work owner 02-10-2026] Excel berlokal Indonesia menyimpan
// CSV dengan `;`. Tanpa tebakan ini seluruh baris judul terbaca SATU kolom,
// dan berkas yang benar ditolak dengan pesan "kehilangan kolom wajib" untuk
// SETIAP kolom - pesan yang menuduh kolom yang jelas terlihat di berkasnya.
//
// ⚠️ Ditebak dari judul SAJA: nama kolom tidak pernah memuat koma maupun
// titik koma, sedangkan isi baris data (nama, angka berkoma) bisa memuat
// keduanya.
func PemisahCSV(awal string) rune {
	if i := strings.IndexAny(awal, "\r\n"); i >= 0 {
		awal = awal[:i]
	}
	var koma, titikKoma int
	kutip := false
	for _, r := range awal {
		switch {
		case r == '"':
			kutip = !kutip
		case kutip:
		case r == ',':
			koma++
		case r == ';':
			titikKoma++
		}
	}
	if titikKoma > koma {
		return ';'
	}
	return ','
}

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
func BacaCSVUnggah(r io.Reader, tipe string) ([]models.BarisUnggah, error) {
	// Kolom wajib bergantung pada Type polis (keputusan work owner 03-10-2026);
	// Type yang tidak dikenal ditolak SEBELUM berkasnya dibaca.
	wajib, err := models.KolomWajibUnggah(tipe)
	if err != nil {
		return nil, err
	}
	br := bufio.NewReaderSize(io.LimitReader(r, BatasUkuranUnggahCSV), ukuranIntipJudul)
	// Galat Peek (berkas lebih pendek dari jendela) bukan galat: yang dipakai
	// hanya bita yang sempat terbaca.
	awal, _ := br.Peek(ukuranIntipJudul)
	pemisah := PemisahCSV(string(awal))
	c := csv.NewReader(br)
	c.Comma = pemisah
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
	// `SHARE_NUSANTARA_RE_GROSS` memenuhi kolom wajib `SHARE_NUSANTARA_RE`
	// (keputusan work owner 02-10-2026) - lihat `models.TerapkanAliasShare`.
	if lihat[models.KolomShareNusantaraReGross] {
		lihat[models.KolomShareNusantaraRe] = true
	}

	// ⛔ JUDUL DIPERIKSA LEBIH DAHULU. Kolom yang HILANG SAMA SEKALI dari
	// berkas akan menghasilkan satu penolakan "kolom kosong" untuk setiap
	// barisnya - seribu kalimat yang mengatakan satu hal, dan orang harus
	// membaca seribu untuk menemukan satu.
	var hilang []string
	for _, k := range wajib {
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
		models.TerapkanAliasShare(nilai)
		if pemisah == ';' {
			models.NormalisasiDesimalKoma(nilai)
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
func (u *UnggahPremiumList) Tinjau(ctx context.Context, pelaku inti.Pelaku,
	polisID string, berkas io.Reader) (HasilTinjauUnggah, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilTinjauUnggah{}, err
	}
	if polisID == "" {
		return HasilTinjauUnggah{}, fmt.Errorf("%w: id polis kosong", galat.ErrPermintaanTidakSah)
	}
	tipe, err := u.tipeUnggah(ctx, polisID)
	if err != nil {
		return HasilTinjauUnggah{}, err
	}
	// Jalan yang SAMA dengan Simpan, TANPA hitung Type QR: Validate CSV hanya
	// memeriksa bentuk/kolom wajib dan batas USIA + SUM INSURED produk
	// (keputusan work owner 05-10-2026). Rate/risk dan syarat produk QR
	// diperiksa Calculate CSV.
	_, hasil, err := u.periksaDanHitung(ctx, polisID, berkas, tipe, false)
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
func periksaBerkas(berkas io.Reader, tipe string) ([]models.BarisUnggah, models.HasilUnggah, error) {
	if berkas == nil {
		// ⛔ Galat permintaan, bukan penolakan korpus: "Please Upload CSV File
		// into attachment" milik langkah 11/33 yang ter-remark (sensus
		// 28-09-2026). Handler sudah menjawab 400 lebih dulu.
		return nil, models.HasilUnggah{}, fmt.Errorf("%w: berkas CSV tidak disertakan",
			galat.ErrPermintaanTidakSah)
	}
	baris, err := BacaCSVUnggah(berkas, tipe)
	if err != nil {
		return nil, models.HasilUnggah{}, err
	}
	// OQ-PL-12 (GILIRAN-17): langkah 2 - uang kosong = 0 SEBELUM validasi
	// (kecuali kolom uang wajib - lihat `models.IsiNolUangKosong`).
	models.IsiNolUangKosong(baris)
	hasil, err := models.ValidasiUnggah(tipe, baris)
	if err != nil {
		return nil, models.HasilUnggah{}, err
	}
	return baris, hasil, nil
}

// periksaDanHitung - periksaBerkas, lalu untuk SEMUA Type memeriksa batas
// umur / sum insured produk (`models.PeriksaBatasProduk`, SavePremiumList_Act
// 6-8.2), lalu - hanya bila `hitung` (Calculate CSV) dan Type QR - menghitung
// kolom peserta (`models.HitungPesertaQR`), dan menggabungkan seluruh
// penolakannya. Dipakai Tinjau (hitung=false) DAN Simpan (hitung=true) - satu
// jalan (keputusan work owner 05-10-2026): peserta di luar batas membuat
// Validate tidak lolos, sehingga Calculate CSV tidak dapat diklik; Validate
// CSV tidak memeriksa rate/risk.
//
// ⛔ MEMBACA saja: nol transaksi, nol tulisan. Hasil hitungan ditulis ke
// `baris[i].Nilai`; Simpan yang menyimpannya.
//
// Syarat polis P1-P4 yang gagal → `models.ErrSyaratHitungQR` (409 berpesan),
// seluruh unggahan ditolak.
func (u *UnggahPremiumList) periksaDanHitung(ctx context.Context, polisID string, berkas io.Reader,
	tipe string, hitung bool) ([]models.BarisUnggah, models.HasilUnggah, error) {

	baris, hasil, err := periksaBerkas(berkas, tipe)
	if err != nil {
		return nil, models.HasilUnggah{}, err
	}
	rujukan := repository.NewRujukan(u.svc.DB())
	// Kepala polis dibaca SEKALI: Type, Product Name ID, Class of Business,
	// Premium Payment Method, R/I SLIP RNM.
	kepala, err := rujukan.KepalaUnggah(ctx, polisID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return nil, models.HasilUnggah{}, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return nil, models.HasilUnggah{}, err
	}
	// Baris yang sudah ditolak tidak dicek ulang: validasi bentuk → batas → hitung.
	lewati := map[int]bool{}
	tandai := func() {
		for _, p := range hasil.Ditolak {
			lewati[p.Baris] = true
		}
	}
	tandai()
	// Batas produk - semua Type. Produk kosong / tidak ada di
	// M_PRODUCTNAME_LIFE tidak diperiksa (Type QR ditolak syarat hitung saat
	// Calculate CSV).
	if kepala.ProductNameID != "" {
		batas, ada, err := rujukan.BatasProduk(ctx, kepala.ProductNameID)
		if err != nil {
			return nil, models.HasilUnggah{}, err
		}
		if ada {
			hasil.Ditolak = append(hasil.Ditolak,
				models.PeriksaBatasProduk(tipe, kepala.RISlipRNM, batas, baris, lewati)...)
		}
	}
	// Hitung Type QR - Calculate CSV saja. Baris yang melewati batas tidak
	// dihitung (penolakan rate yang hanya akibat umur tidak ditambahkan).
	if hitung && models.HitungPesertaPerTipe(tipe) {
		tandai()
		m, err := u.masterQR(ctx, rujukan, kepala)
		if err != nil {
			return nil, models.HasilUnggah{}, err
		}
		hasil.Ditolak = append(hasil.Ditolak, models.HitungPesertaQR(m, baris, lewati)...)
	}
	// Urut nomor baris - penolakan satu peserta tampil berdekatan.
	sort.SliceStable(hasil.Ditolak, func(i, j int) bool { return hasil.Ditolak[i].Baris < hasil.Ditolak[j].Baris })
	return baris, hasil, nil
}

// masterQR membaca bahan hitung Type QR SEKALI per polis dan menegakkan
// syarat P1-P4.
func (u *UnggahPremiumList) masterQR(ctx context.Context, rujukan *repository.Rujukan,
	kepala repository.KepalaUnggah) (models.MasterQR, error) {

	if kepala.ProductNameID == "" {
		return models.MasterQR{}, fmt.Errorf("%w: choose the Product Name and press Save Data before Calculate CSV",
			models.ErrSyaratHitungQR)
	}
	param, plan, err := rujukan.BahanHitungQR(ctx, kepala.ProductNameID)
	if errors.Is(err, repository.ErrRincianProdukTidakAda) {
		return models.MasterQR{}, fmt.Errorf("%w: product %s is not found in Master Product Name Life",
			models.ErrSyaratHitungQR, kepala.ProductNameID)
	}
	if err != nil {
		return models.MasterQR{}, err
	}
	param.BusinessName = kepala.BusinessName
	param.ProRateType = kepala.ProRateType
	pl, err := models.PlanCocokQR(param, plan)
	if err != nil {
		return models.MasterQR{}, err
	}
	rate, err := rujukan.RateHitungQR(ctx, strings.TrimSpace(pl.RIRateID))
	if err != nil {
		return models.MasterQR{}, err
	}
	risk, err := rujukan.RiskHitungQR(ctx, param.RIRiskID)
	if err != nil {
		return models.MasterQR{}, err
	}
	return models.SiapkanMasterQR(param, plan, rate, risk)
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
func (u *UnggahPremiumList) Simpan(ctx context.Context, pelaku inti.Pelaku,
	polisID string, berkas io.Reader) (HasilSimpanUnggah, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilSimpanUnggah{}, err
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return HasilSimpanUnggah{}, db.ErrTanpaOracle
	}
	if polisID == "" {
		return HasilSimpanUnggah{}, fmt.Errorf("%w: id polis kosong", galat.ErrPermintaanTidakSah)
	}

	// ⛔ GERBANG KASUS TERTUTUP - butir bb, lewat `T_WORK_POLIS`.
	keadaanKerja, err := repository.NewWorkPolis(u.svc.DB()).Keadaan(ctx, polisID)
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
	tipe, err := u.tipeUnggah(ctx, polisID)
	if err != nil {
		return HasilSimpanUnggah{}, err
	}
	// Type QR: kolom peserta DIHITUNG di sini, sebelum transaksi (keputusan
	// work owner 05-10-2026); Type lain nilainya tetap dari CSV.
	baris, hasil, err := u.periksaDanHitung(ctx, polisID, berkas, tipe, true)
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
	peserta := repository.NewPesertaUnggah(u.svc.DB())
	err = u.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		dihapus, err := peserta.HapusPesertaPolis(ctx, tx, polisID)
		if err != nil {
			return err
		}
		disimpan, err := peserta.SisipPeserta(ctx, tx, polisID, baris)
		if err != nil {
			return err
		}
		// Rekap summary dihitung ulang dari peserta yang baru - bukan saat
		// Confirm (keputusan work owner 03-10-2026).
		if err := u.svc.SummaryPremiumList().perbaruiRekapDalam(ctx, tx, polisID); err != nil {
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
	// Batas produk (SavePremiumList_Act 6-8.2) TIDAK diperiksa lagi di sini:
	// ia penolakan periksaDanHitung di atas (keputusan work owner 05-10-2026).
	return keluar, nil
}

// tipeUnggah membaca `Type` polis - penentu kolom wajib unggahan (keputusan
// work owner 03-10-2026). MEMBACA saja: tanpa transaksi, tanpa tulis.
func (u *UnggahPremiumList) tipeUnggah(ctx context.Context, polisID string) (string, error) {
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return "", db.ErrTanpaOracle
	}
	d, err := repository.NewPenawaran(u.svc.DB()).BacaDataPolis(ctx, polisID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return "", fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return "", err
	}
	return d.Type, nil
}

// batasRingkasPenolakan - penolakan paling banyak yang disebut di ringkasan.
const batasRingkasPenolakan = 10

// RingkasPenolakan menyusun satu kalimat dari penolakan Calculate CSV - untuk
// kunci `galat` jawaban 409 (keputusan work owner 05-10-2026: rate/risk QR
// tidak diperiksa Validate CSV, jadi penolakannya baru muncul di sini).
func RingkasPenolakan(d []models.Penolakan) string {
	if len(d) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Calculate CSV rejected %d row(s); nothing was saved.", len(models.HasilUnggah{Ditolak: d}.BarisDitolak()))
	for i, p := range d {
		if i == batasRingkasPenolakan {
			fmt.Fprintf(&b, " ... and %d more.", len(d)-i)
			break
		}
		fmt.Fprintf(&b, " Row %d %s: %s.", p.Baris, p.Kolom, strings.TrimSuffix(p.Pesan, "."))
	}
	return b.String()
}
