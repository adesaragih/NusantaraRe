package services

// Unggah / unduh / hapus dokumen lewat outbox - butir be.
//
// Untuk apa berkas ini: memanggil aturan yang sejak kelompok Dokumen sudah
// ditiru tetapi belum punya satu pun pemanggil - `models.KunciKelompokDokumen`,
// `models.IDDokumenBaru`, `models.PerluHapusDiPenyimpanan`,
// `models.MimeDokumen`. Utang itu dinyatakan di kepala `models/dokumenbaru.go`;
// berkas inilah pelunasannya.
//
// Pohon yang ditiru, dibaca 27-09-2026:
//
//	`Activity/SaveAttachLife.xml`     1.2 b483 kunci `DL-`, 1.7 b1467 panggil
//	`Activity/InsertDocument_Act.xml` 3 b627 ID/TANGGAL/MIME, 4 b1023 storage,
//	                                  5 b1185 simpan (prasyarat b1283)
//	`Activity/DeleteDocument_Act.xml` prasyarat b472, `Obj-Delete` b513
//	`Activity/DownloadDocumentClaim.xml` -> `GetUrlGoogleStorage_Act`
//	`RDBList/Insert_T_Storage_SQL.xml`   b86-100 kolom kartu berkas
//
// ⛔ PELAKSANA STUB, BUKAN PENYAMBUNGAN NYATA. Efek `storage-unggah` dan
// `storage-hapus` masuk outbox `T_LOG_SERVICE_RNM`; yang menjalankannya di
// giliran ini adalah pelaksana lokal yang menaruh berkas di folder kita
// sendiri. Penyambungan Google Storage menuntut persetujuan manusia, dan
// TIDAK dilakukan di sini.
//
// Dibaca sesudah: dokumen.go (gerbang kategorinya) dan efekkeluar.go.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

var (
	// ErrUnggahanDirBelumDisetel - folder unggahan belum dipilih.
	//
	// ⛔ Gagal TERANG, bukan bawaan diam-diam ke `./unggahan`. Berkas
	// pelanggan tidak boleh mendarat di folder kerja siapa pun yang
	// kebetulan menjalankan server.
	ErrUnggahanDirBelumDisetel = errors.New(
		"services: UNGGAHAN_DIR belum disetel; unggahan dokumen ditolak")
	// ErrBerkasTerlaluBesar - melebihi BatasUkuranUnggahan.
	ErrBerkasTerlaluBesar = errors.New("services: berkas melebihi batas ukuran")
	// ErrBerkasKosong - nol byte.
	ErrBerkasKosong = errors.New("services: berkas kosong")
	// ErrKategoriDokumenTidakDikenal - kategori di luar daftar wajib.
	ErrKategoriDokumenTidakDikenal = errors.New(
		"services: kategori dokumen tidak ada di daftar kategori")
	// ErrDokumenBelumTerunggah - efek unggahnya belum selesai.
	ErrDokumenBelumTerunggah = errors.New(
		"services: berkas belum selesai diunggah")
)

// Jenis efek outbox untuk berkas - butir be.
//
// ⛔ TEKS, dan tersimpan di kolom `JENIS_EFEK`. Ia dibaca kembali oleh
// proses lain, mungkin berhari-hari kemudian: nilainya tidak boleh berubah
// hanya karena sebuah konstanta Go diganti namanya.
const (
	// JenisEfekStorageUnggah - padanan `InsertGoogleStorage_Act` b1023.
	JenisEfekStorageUnggah = "storage-unggah"
	// JenisEfekStorageHapus - padanan `DeleteGoogleStorage_Act`.
	JenisEfekStorageHapus = "storage-hapus"
)

// BatasUkuranUnggahan adalah batas byte satu berkas.
//
// ⚠️ `[keputusan kami - tidak ada di korpus]` Nol rule di `Claim Life`
// menyebutkan batas ukuran: Pega menyerahkannya ke `dragDropFileUpload`
// bawaan platform, yang batasnya hidup di konfigurasi sistem dan tidak ikut
// diekspor. 25 MiB dipilih sebab ia lebih besar daripada seluruh jenis
// dokumen di tabel MIME 48 baris yang masuk akal dilampirkan pada klaim, dan
// cukup kecil supaya satu permintaan tidak menahan memori proses.
//
// Angkanya berdiri di SINI, bernama, supaya ia dapat dibantah - bukan
// tersebar sebagai literal di handler.
const BatasUkuranUnggahan = 25 << 20

// NamaAplikasiBerkas mengisi kolom `APPNAME` kartu penyimpanan.
//
// `[data DBA]` `GCP_IMAGE.APPNAME VARCHAR2(20)`; korpus memakai kolom senama
// di `T_STORAGE_IMAGE` (`Insert_T_Storage_SQL.xml` b91). Nilainya milik kita.
const NamaAplikasiBerkas = "RNM-CLAIM-LIFE"

// PenyimpananStandar adalah nilai kolom `STORAGE`.
//
// VERBATIM `Insert_T_Storage_SQL.xml` b100: literal `'standard'`.
const PenyimpananStandar = "standard"

// BerkasMasuk adalah satu berkas yang sedang diunggah.
//
// ⚠️ `Isi` pembaca, bukan `[]byte`. Berkas 25 MiB yang dibaca seluruhnya ke
// memori sebelum diperiksa adalah 25 MiB yang dapat diminta siapa saja,
// berkali-kali, sebelum satu pun gerbang berjalan.
type BerkasMasuk struct {
	NamaFile string
	// Mime dari pemanggil. Kosong berarti diturunkan dari nama berkas.
	Mime string
	// Kategori adalah `KATEGORI_2` - yang `DocumentLife.xml` tampilkan.
	Kategori string
	Isi      io.Reader
}

// Unggahan melayani ketiga tombol `Section/DocumentLife.xml`.
type Unggahan struct {
	svc     *Service
	sumber  SumberKategoriWajib
	folder  string
	tautan  func(id int64) string
	jakarta *time.Location
	// batas adalah batas byte satu berkas, bawaannya BatasUkuranUnggahan.
	//
	// ⚠️ Medan, bukan konstanta yang dibaca langsung - dan sebabnya satu
	// kegagalan NYATA. Uji batas menulis 25 MiB dua kali ke disk; sekali ia
	// gagal lalu lulus tiga kali berturut-turut sesudahnya. Uji yang gagal
	// secara acak akan diabaikan orang, bukan dipatuhi - pelajaran yang sama
	// dengan penjaga yang menuduh hal yang benar. Batasnya kini dapat
	// dikecilkan oleh uji; yang menjaga ANGKANYA adalah uji tersendiri atas
	// konstantanya.
	batas int64
}

// Dokumen menyusun layanan unggahan.
//
// ⚠️ `tautan` disuntikkan, bukan dirakit di dalam: bentuk URL unduhnya milik
// lapisan HTTP, dan services tidak boleh tahu jalur rutenya. Bawaannya
// dinyatakan di `TautanUnduhBawaan`.
func (s *Service) Dokumen() *Unggahan {
	jakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// ⚠️ Basis data zona waktu tidak selalu ada di Windows. UTC+7 tetap
		// zona yang benar untuk Jakarta; yang hilang hanya namanya.
		jakarta = time.FixedZone("WIB", 7*60*60)
	}
	return &Unggahan{
		svc:     s,
		sumber:  KategoriWajibBelumDiketahui{},
		folder:  s.unggahanDir,
		tautan:  TautanUnduhBawaan,
		jakarta: jakarta,
		batas:   BatasUkuranUnggahan,
	}
}

// DenganKategori mengganti sumber daftar kategori - dipakai test dan cmd.
func (u *Unggahan) DenganKategori(s SumberKategoriWajib) *Unggahan {
	salin := *u
	salin.sumber = s
	return &salin
}

// DenganFolder mengganti folder unggahan - dipakai test.
func (u *Unggahan) DenganFolder(f string) *Unggahan {
	salin := *u
	salin.folder = f
	return &salin
}

// DenganBatas mengganti batas ukuran - dipakai test saja.
//
// ⛔ Nol dan negatif DIABAIKAN: batas yang dapat dimatikan adalah batas
// yang suatu hari akan dimatikan di jalur nyata.
func (u *Unggahan) DenganBatas(n int64) *Unggahan {
	if n <= 0 {
		return u
	}
	salin := *u
	salin.batas = n
	return &salin
}

// TautanUnduhBawaan menyusun URL publik sebuah dokumen di DEV.
//
// ⛔ Ia jalur KITA, bukan URL penyimpanan luar. Butir **be**: selama
// penyambungan nyata belum disetujui, `URLPUBLIC` berisi rute ini - dan rute
// ini berbatas identitas.
func TautanUnduhBawaan(id int64) string {
	return "/api/dokumen/" + strconv.FormatInt(id, 10) + "/isi"
}

// pagariUnggahan menjalankan gerbang yang sama untuk ketiga jalurnya.
func (u *Unggahan) pagari(ctx context.Context, pelaku Pelaku, klaimID string) error {
	if err := WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", ErrPermintaanTidakSah)
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return repository.ErrTanpaOracle
	}
	// ⛔ BUTIR bb - kasus tertutup tidak menerima dokumen baru maupun
	// kehilangan dokumen lama.
	return u.svc.PastikanKasusTerbuka(ctx, klaimID)
}

// Unggah menerima satu berkas dan mencatatnya - `Add attachment` b1245.
//
// Urutannya, dan setiap langkahnya punya sebab:
//
//  1. gerbang identitas, bentuk, kasus terbuka, kategori, folder;
//  2. berkas ditulis ke `UNGGAHAN_DIR` LEBIH DULU - byte yang sudah mendarat
//     tidak pernah tanpa catatan bila langkah 3 gagal, sebab langkah 3 yang
//     gagal membatalkan transaksinya dan berkas yatim itu terlihat di folder;
//     urutan sebaliknya kehilangan byte-nya sama sekali;
//  3. satu transaksi: baris dokumen + baris outbox.
//
// ⛔ `T_STORAGE_ID` KOSONG di langkah 3, dan penyimpangannya dijelaskan di
// `repository.SisipDokumen`.
func (u *Unggahan) Unggah(ctx context.Context, pelaku Pelaku,
	klaimID, pesertaID string, berkas BerkasMasuk, saat time.Time) (
	models.Dokumen, error) {

	if err := u.pagari(ctx, pelaku, klaimID); err != nil {
		return models.Dokumen{}, err
	}
	if strings.TrimSpace(pesertaID) == "" {
		return models.Dokumen{}, fmt.Errorf("%w: pengenal peserta wajib diisi",
			ErrPermintaanTidakSah)
	}
	if strings.TrimSpace(u.folder) == "" {
		return models.Dokumen{}, ErrUnggahanDirBelumDisetel
	}
	nama := strings.TrimSpace(berkas.NamaFile)
	if nama == "" {
		return models.Dokumen{}, fmt.Errorf("%w: nama berkas wajib diisi",
			ErrPermintaanTidakSah)
	}
	// ⛔ Kategori diperiksa terhadap daftar yang SAMA dengan yang menggerbangi
	// Save ke Outstanding (butir ar1). Dua daftar kategori berarti dokumen
	// dapat diunggah dengan kategori yang gerbang penyimpanan tolak, dan
	// pemakai baru tahu berbulan-bulan kemudian.
	if err := u.periksaKategori(ctx, berkas.Kategori); err != nil {
		return models.Dokumen{}, err
	}

	idTeks := models.IDDokumenBaru(saat, u.jakarta)
	id, err := strconv.ParseInt(idTeks, 10, 64)
	if err != nil {
		return models.Dokumen{}, fmt.Errorf(
			"services: pengenal dokumen %q bukan angka: %w", idTeks, err)
	}
	jalur := filepath.Join(u.folder, namaBerkasLokal(idTeks, nama))
	if err := u.tulisBerkas(jalur, berkas.Isi); err != nil {
		return models.Dokumen{}, err
	}

	dok := models.Dokumen{
		ID:        id,
		PesertaID: pesertaID,
		NamaFile:  nama,
		// ⛔ Pemanggil MENANG - prasyarat b586 `Param.MIME==""` `WhenTrue=2`
		// LANJUT, artinya turunan dari nama berkas hanya dipakai bila
		// pemanggil diam. `@toLowerCase` b781 ada di dalam MimeDokumen.
		Mime: models.MimeDokumen(berkas.Mime, nama),
		// `KATEGORI_1` = kunci kelompok `DL-` b595; satu unggahan satu
		// kelompok, dan kelompok yang sudah ada tidak pernah ditimpa.
		Kategori1:  models.KunciKelompokDokumen("", saat),
		Kategori2:  strings.TrimSpace(berkas.Kategori),
		TStorageID: "",
	}

	baca := repository.NewKlaimLife(u.svc.db)
	err = u.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		if err := baca.SisipDokumen(ctx, tx, dok, saat); err != nil {
			return err
		}
		// ⛔ Outbox DI DALAM transaksi yang sama dengan barisnya. Efek yang
		// diantre terpisah dapat hilang sendirian, dan berkas yang sudah
		// mendarat tidak akan pernah bertaut.
		_, err := repository.NewPohonKlaim(u.svc.db).AntreEfek(ctx, tx,
			models.LiniLife, ModulClaimLife, JenisEfekStorageUnggah, idTeks,
			muatanBerkas(klaimID, jalur, nama), saat)
		return err
	})
	if err != nil {
		// ⚠️ Berkas yang sudah ditulis DIBUANG lagi: transaksinya batal, jadi
		// tidak ada baris yang merujuknya. Kegagalan membuang dilaporkan
		// bersama sebab aslinya - berkas yatim yang diam adalah berkas yang
		// tidak akan pernah dicari.
		if e := os.Remove(jalur); e != nil && !os.IsNotExist(e) {
			return models.Dokumen{}, fmt.Errorf(
				"%w (dan berkas %s gagal dibuang: %v)", err, jalur, e)
		}
		return models.Dokumen{}, err
	}
	return dok, nil
}

// periksaKategori menolak kategori di luar daftar butir ar1.
func (u *Unggahan) periksaKategori(ctx context.Context, kategori string) error {
	k := strings.TrimSpace(kategori)
	if k == "" {
		return fmt.Errorf("%w: kategori dokumen wajib diisi", ErrPermintaanTidakSah)
	}
	wajib, err := u.sumber.KategoriWajib(ctx)
	if err != nil {
		return err
	}
	for _, w := range wajib {
		if strings.EqualFold(strings.TrimSpace(w), k) {
			return nil
		}
	}
	return fmt.Errorf("%w: %q", ErrKategoriDokumenTidakDikenal, k)
}

// tulisBerkas menyalin isi ke folder unggahan, berbatas ukuran.
//
// ⛔ Batasnya ditegakkan SAAT MENYALIN, bukan dari header `Content-Length`.
// Header itu datang dari pengirim dan dapat berbohong; `io.LimitReader`
// tidak dapat.
func (u *Unggahan) tulisBerkas(jalur string, isi io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(jalur), 0o750); err != nil {
		return fmt.Errorf("services: menyiapkan folder unggahan: %w", err)
	}
	f, err := os.OpenFile(jalur, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
	if err != nil {
		return fmt.Errorf("services: membuat berkas unggahan: %w", err)
	}
	// Satu byte LEBIH daripada batas, supaya kelebihan dapat dibedakan dari
	// berkas yang panjangnya tepat sebesar batas.
	batas := u.batas
	if batas <= 0 {
		batas = BatasUkuranUnggahan
	}
	n, salinErr := io.Copy(f, io.LimitReader(isi, batas+1))
	tutupErr := f.Close()
	switch {
	case salinErr != nil:
		_ = os.Remove(jalur)
		return fmt.Errorf("services: menulis berkas unggahan: %w", salinErr)
	case tutupErr != nil:
		_ = os.Remove(jalur)
		return fmt.Errorf("services: menutup berkas unggahan: %w", tutupErr)
	case n == 0:
		_ = os.Remove(jalur)
		return ErrBerkasKosong
	case n > batas:
		_ = os.Remove(jalur)
		return fmt.Errorf("%w: %d byte, batas %d", ErrBerkasTerlaluBesar, n, batas)
	}
	return nil
}

// ekstensi mengambil akhiran nama berkas, huruf kecil.
func ekstensi(nama string) string {
	return strings.ToLower(filepath.Ext(nama))
}

// namaBerkasLokal menyusun nama berkas di `UNGGAHAN_DIR`.
//
// ⛔ SATU tempat, dipakai penulis dan pembacanya. Dua tempat berarti
// unggah dan unduh dapat bergeser sendiri-sendiri, dan yang bergeser tidak
// akan berbunyi: berkasnya ada, pencarinya menengok ke nama lain.
//
// ⚠️ Kuncinya pengenal DOKUMEN, bukan `IMAGEID`. Nama berkas di disk
// kita bukan kunci penyimpanan luar; menyatukannya membuat penggantian
// pelaksana - yang memang direncanakan - mengubah letak berkas lama pula.
func namaBerkasLokal(idDokumen, namaAsli string) string {
	return idDokumen + ekstensi(namaAsli)
}

// muatanBerkas menyusun muatan JSON outbox.
//
// ⚠️ Disusun dengan tangan, bukan `encoding/json`, supaya bentuknya TERLIHAT
// di berkas ini - ia menyeberang batas proses dan batas waktu, dan struct Go
// yang berubah minggu depan tidak boleh membuat baris terantre tak terbaca.
func muatanBerkas(klaimID, jalur, nama string) string {
	return `{"klaimId":` + strconv.Quote(klaimID) +
		`,"jalur":` + strconv.Quote(jalur) +
		`,"namaFile":` + strconv.Quote(nama) + `}`
}

// Unduh menyerahkan berkas satu dokumen - `View Office Online` b3502 dan
// tautan baris `DocumentLife.xml`.
//
// `[terverifikasi]` `DownloadDocumentClaim.xml` memanggil
// `GetUrlGoogleStorage_Act`, yang membaca `GetLinkStorage_SQL` b85
// (`URLPUBLIC ... where imageid = {UploadDoc.ImageID}`). Di DEV URL itu
// adalah rute kita sendiri, jadi jalur ini yang melayaninya.
//
// ⛔ Mengembalikan JALUR berkas, bukan isinya. Membaca 25 MiB ke memori
// hanya untuk menyalinnya ke jaringan adalah 25 MiB yang dapat diminta
// serentak oleh siapa pun yang punya identitas.
func (u *Unggahan) Unduh(ctx context.Context, pelaku Pelaku,
	klaimID string, dokID int64) (models.Dokumen, string, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return models.Dokumen{}, "", err
	}
	if strings.TrimSpace(klaimID) == "" {
		return models.Dokumen{}, "", fmt.Errorf("%w: pengenal klaim wajib diisi",
			ErrPermintaanTidakSah)
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return models.Dokumen{}, "", repository.ErrTanpaOracle
	}
	// ⚠️ TANPA `PastikanKasusTerbuka`. Unduh MEMBACA; kasus yang sudah
	// ditutup tetap boleh dibaca - butir bb menutup jalur pengubah, bukan
	// jalur baca. Menutup unduhan berarti dokumen klaim selesai tidak dapat
	// dilihat lagi oleh siapa pun.
	dok, err := repository.NewKlaimLife(u.svc.db).SatuDokumen(ctx, klaimID, dokID)
	if err != nil {
		return models.Dokumen{}, "", err
	}
	// ⛔ `T_STORAGE_ID` kosong berarti efek unggahnya BELUM selesai - bukan
	// berkasnya hilang. Dibedakan supaya layar dapat berkata "sedang
	// diproses" alih-alih "tidak ditemukan".
	if !models.BolehSimpanBarisDokumen(dok.TStorageID) {
		return models.Dokumen{}, "", fmt.Errorf("%w: dokumen %d",
			ErrDokumenBelumTerunggah, dokID)
	}
	if strings.TrimSpace(u.folder) == "" {
		return models.Dokumen{}, "", ErrUnggahanDirBelumDisetel
	}
	// ⛔ Dari pengenal DOKUMEN, bukan dari `T_STORAGE_ID`. Sejak 28-09-2026
	// keduanya BERBEDA: `IMAGEID` lahir dari `GenerateImageID_SQL` (MD5 atas
	// cap waktu + `SYS_GUID()`), sedangkan nama berkas lokal adalah urusan
	// kita sendiri. Ronde pertama memakai `T_STORAGE_ID` karena saat itu
	// keduanya kebetulan sama nilainya - dan setiap unduhan akan gagal
	// seketika rumus `IMAGEID` yang benar dipasang.
	return dok, filepath.Join(u.folder,
		namaBerkasLokal(strconv.FormatInt(dok.ID, 10), dok.NamaFile)), nil
}

// Hapus membuang satu dokumen - `Delete` b4288 -> `ConfirmDeleteAttachment`
// b4317 -> `DeleteDocument_Act`.
//
// ⛔ URUTANNYA DARI RULE, dan ia berlawanan dengan unggah. `DeleteDocument_Act`
// b513 `Obj-Delete` berjalan TANPA prasyarat: barisnya dihapus apa pun
// keadaan penyimpanannya. Yang berprasyarat hanya panggilan penyimpanan -
// b472 `WhenTrue=3` LEWATI bila `T_STORAGE_ID` kosong, yaitu
// `models.PerluHapusDiPenyimpanan`.
//
// Jadi: baris dulu, penyimpanan menyusul lewat outbox. Berkas lokal dibuang
// oleh pelaksana efek, bukan di sini - bila transaksinya batal, berkasnya
// masih ada dan barisnya pun masih ada.
func (u *Unggahan) Hapus(ctx context.Context, pelaku Pelaku,
	klaimID string, dokID int64, saat time.Time) error {

	if err := u.pagari(ctx, pelaku, klaimID); err != nil {
		return err
	}
	baca := repository.NewKlaimLife(u.svc.db)
	dok, err := baca.SatuDokumen(ctx, klaimID, dokID)
	if err != nil {
		return err
	}
	return u.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		if err := baca.HapusDokumen(ctx, tx, dokID); err != nil {
			return err
		}
		if !models.PerluHapusDiPenyimpanan(dok.TStorageID) {
			// Prasyarat b472: nol penyimpanan yang perlu disentuh. Berkas
			// lokalnya pun belum pernah tertaut, jadi tidak ada yang
			// diantre - mengantre pekerjaan kosong membanjiri outbox.
			return nil
		}
		_, err := repository.NewPohonKlaim(u.svc.db).AntreEfek(ctx, tx,
			LiniLife, ModulClaimLife, JenisEfekStorageHapus, dok.TStorageID,
			muatanBerkas(klaimID, filepath.Join(u.folder,
				namaBerkasLokal(strconv.FormatInt(dok.ID, 10), dok.NamaFile)),
				dok.NamaFile), saat)
		return err
	})
}

// PelaksanaBerkasLokal menjalankan efek berkas TANPA layanan luar.
//
// ⛔ INI STUB, dan namanya menyebutnya. Ia memenuhi bentuk
// `InsertGoogleStorage_Act` / `DeleteGoogleStorage_Act` - kartu berkas
// tertulis, `T_STORAGE_ID` terisi, tautan publik ada - tanpa menghubungi
// Google Storage sama sekali. Penyambungan nyata menuntut persetujuan
// manusia; sampai itu datang, seluruh jalur unggah-unduh-hapus tetap dapat
// dijalankan dan diuji dari ujung ke ujung.
//
// ⚠️ Pelaksana NYATA kelak adalah implementasi LAIN dari antarmuka yang
// sama, dipilih lewat env - bukan berkas ini yang diubah. Itu sebab ia
// antarmuka.
type PelaksanaBerkasLokal struct {
	baca   *repository.KlaimLife
	tautan func(id int64) string
	// jam menyuplai cap waktu `TANGGAL_UPLOAD` dan masukan `IMAGEID`.
	//
	// ⚠️ Medan, bukan `time.Now()` di dalam badan. Fungsi yang membaca
	// jamnya sendiri tidak dapat diuji tanpa menunggu waktu berlalu, dan
	// `IMAGEID` adalah hash ATAS jam itu.
	jam func() time.Time
}

// NewPelaksanaBerkasLokal menyusun pelaksana stub-nya.
func NewPelaksanaBerkasLokal(svc *Service) *PelaksanaBerkasLokal {
	return &PelaksanaBerkasLokal{
		baca:   repository.NewKlaimLife(svc.db),
		tautan: TautanUnduhBawaan,
		jam:    time.Now,
	}
}

// DenganJam mengganti sumber waktunya - dipakai uji.
func (p *PelaksanaBerkasLokal) DenganJam(j func() time.Time) *PelaksanaBerkasLokal {
	salin := *p
	salin.jam = j
	return &salin
}

// Laksanakan menjalankan satu baris outbox jenis berkas.
//
// ⛔ Jenis yang TIDAK dikenalnya dikembalikan sebagai galat, bukan diam-diam
// dianggap berhasil. Outbox lintas modul: baris milik jenis lain yang
// terlanjur dipungut harus terlihat, bukan tertandai selesai tanpa pernah
// dikerjakan.
func (p *PelaksanaBerkasLokal) Laksanakan(ctx context.Context, tx *repository.Tx,
	b repository.BarisEfekKeluar) error {

	switch b.Jenis {
	case JenisEfekStorageUnggah:
		return p.unggah(ctx, tx, b)
	case JenisEfekStorageHapus:
		return p.hapus(ctx, tx, b)
	default:
		return fmt.Errorf("%w: jenis efek %q bukan milik pelaksana berkas",
			ErrPermintaanTidakSah, b.Jenis)
	}
}

// unggah menulis kartu berkas lalu menautkannya ke barisnya.
func (p *PelaksanaBerkasLokal) unggah(ctx context.Context, tx *repository.Tx,
	b repository.BarisEfekKeluar) error {

	id, err := strconv.ParseInt(b.Rujukan, 10, 64)
	if err != nil {
		// ⛔ Permanen, bukan layak dicoba ulang: rujukan yang bukan angka
		// tidak akan menjadi angka pada percobaan kedelapan.
		return fmt.Errorf("%w: rujukan %q bukan pengenal dokumen",
			ErrPermintaanTidakSah, b.Rujukan)
	}
	saat := p.jam()
	// ⛔ `IMAGEID` DITERBITKAN DI SINI, dan bukan pengenal dokumennya.
	// `GenerateImageID_SQL.xml` b85-b88 adalah rule pembangkitnya sendiri, dan
	// `InsertGoogleStorage_Act.xml` b2226 memanggilnya tepat di titik ini -
	// saat berkas masuk penyimpanan, bukan saat barisnya lahir.
	//
	// ⚠️ RALAT 28-09-2026: ronde pertama memakai `b.Rujukan` (pengenal
	// dokumen, yaitu cap waktu) sebagai `IMAGEID`. Kunci penyimpanan yang
	// dapat ditebak dari waktu unggah bukan kunci.
	imageID, err := models.ImageIDBaru(saat)
	if err != nil {
		return err
	}
	if err := p.baca.SisipKartuBerkas(ctx, tx, repository.KartuBerkas{
		ImageID: imageID,
		// ⛔ URL-nya jalur KITA, berbatas identitas - butir be.
		URLPublic: p.tautan(id),
		FileName:  namaDariMuatan(b.Muatan),
		AppName:   NamaAplikasiBerkas,
		Storage:   PenyimpananStandar,
		// `Update_T_Storage_SQL.xml` b89. Di sistem lama nilainya datang dari
		// jawaban layanan penyimpanan; di sini layanan itu kita sendiri.
		TanggalUpload: saat,
	}); err != nil {
		return err
	}
	return p.baca.TandaiDokumenTerunggah(ctx, tx, id, imageID)
}

// hapus membuang kartu berkas dan berkas lokalnya.
//
// ⚠️ Berkas yang SUDAH tidak ada bukan kegagalan. Efek dijamin berjalan
// SETIDAKNYA sekali; percobaan kedua atas berkas yang sudah terhapus harus
// berhasil, bukan berputar sampai jatah percobaannya habis.
func (p *PelaksanaBerkasLokal) hapus(ctx context.Context, tx *repository.Tx,
	b repository.BarisEfekKeluar) error {

	if err := p.baca.HapusKartuBerkas(ctx, tx, b.Rujukan); err != nil {
		return err
	}
	jalur := jalurDariMuatan(b.Muatan)
	if jalur == "" {
		return nil
	}
	if err := os.Remove(jalur); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("services: menghapus berkas %s: %w", jalur, err)
	}
	return nil
}

// muatanBerkasTerbaca adalah bentuk muatan outbox jenis berkas.
//
// ⚠️ Bentuknya SATU tempat, dan penulisnya (`muatanBerkas`) berdiri di atas.
// Dua tempat berarti pembaca dan penulis dapat bergeser sendiri-sendiri -
// bentuk cacat lintas-lapis yang sudah enam kali terjadi di modul ini, kali
// ini menyeberangi batas WAKTU alih-alih batas bahasa.
type muatanBerkasTerbaca struct {
	KlaimID  string `json:"klaimId"`
	Jalur    string `json:"jalur"`
	NamaFile string `json:"namaFile"`
}

// bacaMuatanBerkas mengurai muatan outbox jenis berkas.
//
// ⚠️ Muatan yang tidak terbaca mengembalikan nilai KOSONG, bukan galat.
// Baris outbox berumur panjang; bentuk JSON-nya dapat berubah di antara saat
// ia ditulis dan saat ia dijalankan. Yang penting - `RUJUKAN` - tidak ada di
// muatan, jadi efeknya tetap dapat diselesaikan tanpa muatan sama sekali.
func bacaMuatanBerkas(teks string) muatanBerkasTerbaca {
	var m muatanBerkasTerbaca
	if err := json.Unmarshal([]byte(teks), &m); err != nil {
		return muatanBerkasTerbaca{}
	}
	return m
}

func namaDariMuatan(teks string) string  { return bacaMuatanBerkas(teks).NamaFile }
func jalurDariMuatan(teks string) string { return bacaMuatanBerkas(teks).Jalur }

// UnduhLewatPengenal melayani `GET /api/dokumen/{id}/isi`.
//
// ⛔ Jalur rute itu TIDAK menyebut klaim - sebab `URLPUBLIC` di sistem lama
// pun tidak (`GetLinkStorage_SQL` b91 mencari dengan `imageid` saja). Klaimnya
// dicari lebih dulu, lalu pembacaan yang SAMA dengan rute lain dipakai: satu
// pembacaan berbatas klaim, bukan dua.
func (u *Unggahan) UnduhLewatPengenal(ctx context.Context, pelaku Pelaku,
	dokID int64) (models.Dokumen, string, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return models.Dokumen{}, "", err
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return models.Dokumen{}, "", repository.ErrTanpaOracle
	}
	klaimID, err := repository.NewKlaimLife(u.svc.db).KlaimDokumen(ctx, dokID)
	if err != nil {
		return models.Dokumen{}, "", err
	}
	return u.Unduh(ctx, pelaku, klaimID, dokID)
}
