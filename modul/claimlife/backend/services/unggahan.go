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
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

var (
	// ErrKategoriDokumenTerlaluPanjang - nama Document Claim melebihi lebar KATEGORI_2 (11-10-2026).
	ErrKategoriDokumenTerlaluPanjang = errors.New("services: this document category name is longer than the stored " +
		"category column allows; the database column must be widened first")
	// ErrKategoriDokumenTidakDikenal - kategori di luar daftar wajib.
	ErrKategoriDokumenTidakDikenal = errors.New(
		"services: kategori dokumen tidak ada di daftar kategori")
	// ErrDokumenBelumTerunggah - efek unggahnya belum selesai.
	ErrDokumenBelumTerunggah = errors.New(
		"services: berkas belum selesai diunggah")
	// ErrDokumenTidakAda dirujuk ulang supaya handler tidak mengimpor
	// repository: pengenal dokumen yang tidak ada dijawab 404, bukan 500.
	ErrDokumenTidakAda = repository.ErrDokumenTidakAda
)

// NamaAplikasiBerkas mengisi kolom `APPNAME` kartu penyimpanan.
//
// `[data DBA]` `GCP_IMAGE.APPNAME VARCHAR2(20)`; korpus memakai kolom senama
// di `T_STORAGE_IMAGE` (`Insert_T_Storage_SQL.xml` b91). Nilainya milik kita.
const NamaAplikasiBerkas = "RNM-CLAIM-LIFE"

// PenyimpananStandar adalah nilai kolom `STORAGE`.
//
// VERBATIM `Insert_T_Storage_SQL.xml` b100: literal `'standard'`.
const PenyimpananStandar = "standard"

// LebarKategoriDokumen - lebar T_CLAIMLF_DOCUMENT.KATEGORI_2 (VARCHAR2(64), migrasi 010).
const LebarKategoriDokumen = 64

// Unggahan melayani ketiga tombol `Section/DocumentLife.xml`.
type Unggahan struct {
	svc    *Service
	sumber SumberKategoriWajib
	// penyimpanan - Google Storage seperti Master Product Name Life (11-10-2026); uji memakai tiruan.
	penyimpanan PenyimpananDokumen
	folder      string
	tautan      func(id int64) string
	jakarta     *time.Location
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
		svc:         s,
		penyimpanan: s.penyimpananDokumen(),
		sumber:      KategoriWajibBelumDiketahui{},
		folder:      s.UnggahanDir(),
		tautan:      TautanUnduhBawaan,
		jakarta:     jakarta,
		batas:       unggah.BatasUkuranUnggahan,
	}
}

// DenganKategori mengganti sumber daftar kategori - dipakai test dan cmd.
func (u *Unggahan) DenganKategori(s SumberKategoriWajib) *Unggahan {
	salin := *u
	salin.sumber = s
	return &salin
}

// DenganPenyimpanan mengganti penyimpanan objek - dipakai uji (tiruan).
func (u *Unggahan) DenganPenyimpanan(p PenyimpananDokumen) *Unggahan {
	salin := *u
	salin.penyimpanan = p
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
// pesertaMilikKlaim - peserta aktif (tidak tercabut) milik klaim itu.
func (u *Unggahan) pesertaMilikKlaim(ctx context.Context, klaimID, pesertaID string) error {
	peserta, err := repository.NewKlaimLife(u.svc.DB()).AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	for _, p := range peserta {
		if p.ID == pesertaID {
			return nil
		}
	}
	return fmt.Errorf("%w: peserta %q bukan peserta aktif klaim %q", galat.ErrPermintaanTidakSah, pesertaID, klaimID)
}

func (u *Unggahan) pagari(ctx context.Context, pelaku inti.Pelaku, klaimID string) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
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
func (u *Unggahan) Unggah(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string, berkas unggah.BerkasMasuk, saat time.Time) (
	models.Dokumen, error) {

	if err := u.pagari(ctx, pelaku, klaimID); err != nil {
		return models.Dokumen{}, err
	}
	if strings.TrimSpace(pesertaID) == "" {
		return models.Dokumen{}, fmt.Errorf("%w: pengenal peserta wajib diisi",
			galat.ErrPermintaanTidakSah)
	}
	// ⛔ Temuan /code-review GILIRAN-17: peserta harus milik klaim ini dan
	// TIDAK tercabut (OQ-M6) - `AmbilPeserta` menyaring keduanya. Tanpanya
	// dokumen mendarat pada peserta yang tak terlihat, tak dapat dihapus.
	if err := u.pesertaMilikKlaim(ctx, klaimID, pesertaID); err != nil {
		return models.Dokumen{}, err
	}
	if strings.TrimSpace(u.folder) == "" {
		return models.Dokumen{}, unggah.ErrUnggahanDirBelumDisetel
	}
	nama := strings.TrimSpace(berkas.NamaFile)
	if nama == "" {
		return models.Dokumen{}, fmt.Errorf("%w: nama berkas wajib diisi",
			galat.ErrPermintaanTidakSah)
	}
	// ⛔ Kategori diperiksa terhadap daftar butir ar1. (Gerbang kelengkapan
	// Save ke Outstanding yang dulu memakai daftar yang sama ter-remark di XML
	// dan dibuang - butir bl; pemeriksaan unggah ini keputusan ar1 sendiri.)
	if err := u.periksaKategori(ctx, klaimID, berkas.Kategori); err != nil {
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
		Mime: unggah.MimeDokumen(berkas.Mime, nama),
		// `KATEGORI_1` = kunci kelompok `DL-` b595-596.
		// ⚠️ CELAH TERCATAT (sensus remark 28-09-2026, OQ-J): b596 memakai
		// ULANG `Primary.DOCUMENT` peserta bila sudah ada; di sini kolom
		// `DOCUMENT` peserta tidak ada, jadi setiap unggahan membuat kunci
		// baru. Daftar layar tidak terpengaruh (disaring FK peserta), tetapi
		// pengelompokan `KATEGORI_1` berbeda dengan baris warisan.
		Kategori1:  models.KunciKelompokDokumen("", saat),
		Kategori2:  strings.TrimSpace(berkas.Kategori),
		TStorageID: "",
	}

	baca := repository.NewKlaimLife(u.svc.DB())
	err = u.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := baca.SisipDokumen(ctx, tx, dok, saat); err != nil {
			return err
		}
		// ⛔ Outbox DI DALAM transaksi yang sama dengan barisnya. Efek yang
		// diantre terpisah dapat hilang sendirian, dan berkas yang sudah
		// mendarat tidak akan pernah bertaut.
		_, err := outbox.NewPenyimpan(u.svc.DB()).AntreEfek(ctx, tx,
			inti.LiniLife, ModulClaimLife, unggah.JenisEfekStorageUnggah, idTeks,
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
	// 11-10-2026 (keputusan work owner "samakan dengan penyimpanan document product name"): seperti Master Product Name
	// Life (UnggahLampiran -> kirimLampiran), efek penyimpanan dicoba SEGERA sesudah commit, di permintaan yang sama -
	// tidak menunggu pekerja latar. Gagal kirim tidak menggagalkan unggahan: efeknya tetap antre dengan galatnya.
	u.kirimSegera(ctx, idTeks, saat)
	if baru, err := repository.NewKlaimLife(u.svc.DB()).SatuDokumen(ctx, klaimID, id); err == nil {
		return baru, nil
	}
	return dok, nil
}

// kirimSegera - satu percobaan efek "storage-unggah" dokumen itu (pola masterproductnamelife kirimLampiran), dengan
// pelaksana penyimpanan Claim Life yang ada (PelaksanaBerkasLokal: kartu T_CLAIMLF_STORAGE + T_STORAGE_ID).
func (u *Unggahan) kirimSegera(ctx context.Context, rujukan string, saat time.Time) {
	ctx, batal := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
	defer batal()
	baca := repository.NewKlaimLife(u.svc.DB())
	err := u.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		b, ada, err := baca.PungutEfekUnggahDokumen(ctx, tx, ModulClaimLife, unggah.JenisEfekStorageUnggah, rujukan, saat)
		if err != nil || !ada {
			return err
		}
		// 11-10-2026: Google Storage + T_STORAGE_IMAGE (PelaksanaGoogle); PelaksanaBerkasLokal pensiun untuk dokumen baru.
		var jalan error
		if u.penyimpanan != nil {
			jalan = NewPelaksanaGoogle(u.svc, u.penyimpanan).Laksanakan(ctx, tx, b)
		} else {
			// Transport Google belum dipasang (izin penjaga inti): pelaksana lokal seperti sebelumnya.
			jalan = NewPelaksanaBerkasLokal(u.svc).Laksanakan(ctx, tx, b)
		}
		if jalan == nil {
			return outbox.NewPenyimpan(u.svc.DB()).TuntaskanEfek(ctx, tx, b.ID, outbox.StatusEfekSelesai, time.Time{}, "", saat)
		}
		status, jadwal := outbox.StatusEfekAntre, saat.Add(outbox.Backoff(b.Percobaan+1))
		if !outbox.LayakDicobaUlang(jalan) || b.Percobaan >= outbox.PercobaanMaksimum {
			status, jadwal = outbox.StatusEfekGagalPermanen, time.Time{}
		}
		log.Printf("claimlife: menyimpan dokumen %s gagal: %v", rujukan, jalan)
		return outbox.NewPenyimpan(u.svc.DB()).TuntaskanEfek(ctx, tx, b.ID, status, jadwal, jalan.Error(), saat)
	})
	if err != nil {
		log.Printf("claimlife: mencatat hasil penyimpanan dokumen %s gagal: %v", rujukan, err)
	}
}

// periksaKategori menolak kategori di luar daftar butir ar1.
func (u *Unggahan) periksaKategori(ctx context.Context, klaimID, kategori string) error {
	k := strings.TrimSpace(kategori)
	if k == "" {
		return fmt.Errorf("%w: kategori dokumen wajib diisi", galat.ErrPermintaanTidakSah)
	}
	// ⚠️ T_CLAIMLF_DOCUMENT.KATEGORI_2 VARCHAR2(64) (migrasi 010), sedangkan Document Claim master VARCHAR2(500): nama yang
	// lebih panjang DITOLAK terang (422), tidak dipotong. Melebarkan kolom = migrasi (jatah habis) - menunggu keputusan.
	if len(k) > LebarKategoriDokumen {
		return fmt.Errorf("%w: %d byte (kolom KATEGORI_2 hanya %d byte)", ErrKategoriDokumenTerlaluPanjang, len(k), LebarKategoriDokumen)
	}
	wajib, err := u.sumber.KategoriWajib(ctx, klaimID)
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
	return unggah.TulisBerkas(jalur, isi, u.batas)
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
func (u *Unggahan) Unduh(ctx context.Context, pelaku inti.Pelaku,
	klaimID string, dokID int64) (models.Dokumen, string, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return models.Dokumen{}, "", err
	}
	if strings.TrimSpace(klaimID) == "" {
		return models.Dokumen{}, "", fmt.Errorf("%w: pengenal klaim wajib diisi",
			galat.ErrPermintaanTidakSah)
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return models.Dokumen{}, "", db.ErrTanpaOracle
	}
	// ⚠️ TANPA `PastikanKasusTerbuka`. Unduh MEMBACA; kasus yang sudah
	// ditutup tetap boleh dibaca - butir bb menutup jalur pengubah, bukan
	// jalur baca. Menutup unduhan berarti dokumen klaim selesai tidak dapat
	// dilihat lagi oleh siapa pun.
	dok, err := repository.NewKlaimLife(u.svc.DB()).SatuDokumen(ctx, klaimID, dokID)
	if err != nil {
		return models.Dokumen{}, "", err
	}
	// ⛔ `T_STORAGE_ID` kosong berarti efek unggahnya BELUM selesai - bukan
	// berkasnya hilang. Dibedakan supaya layar dapat berkata "sedang
	// diproses" alih-alih "tidak ditemukan".
	// 11-10-2026 (keputusan work owner "View" selama storage belum jalan): dokumen yang efek unggahnya BELUM selesai
	// (T_STORAGE_ID kosong) tetap dilayani dari berkas LOKAL di UNGGAHAN_DIR - cadangan sampai pekerja penyimpanan ada.
	// Jalurnya disusun server dari pengenal dokumen (namaBerkasLokal = ID + ekstensi), tidak pernah dari klien; berkas
	// yang tidak ada dijawab 404 oleh handler. ErrDokumenBelumTerunggah tetap ada untuk pemanggil lain.
	if strings.TrimSpace(u.folder) == "" {
		return models.Dokumen{}, "", unggah.ErrUnggahanDirBelumDisetel
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
func (u *Unggahan) Hapus(ctx context.Context, pelaku inti.Pelaku,
	klaimID string, dokID int64, saat time.Time) error {

	if err := u.pagari(ctx, pelaku, klaimID); err != nil {
		return err
	}
	baca := repository.NewKlaimLife(u.svc.DB())
	// SatuDokumen berbatas klaim: dokumen klaim lain = ErrDokumenTidakAda (404).
	dok, err := baca.SatuDokumen(ctx, klaimID, dokID)
	if err != nil {
		return err
	}
	tertunda := !models.PerluHapusDiPenyimpanan(dok.TStorageID)
	// 11-10-2026 (seperti Master Product Name Life / DeleteDocument_Act → DeleteGoogleStorage_Act): objek di Google Storage
	// dihapus LEBIH DULU - gagal = rekam tetap (pesan jelas); lalu baris T_STORAGE_IMAGE dan baris dokumen satu transaksi.
	objek, diGoogle := repository.ObjekPenyimpanan{}, false
	if !tertunda {
		if objek, diGoogle, err = baca.AmbilObjekPenyimpanan(ctx, nil, dok.TStorageID); err != nil {
			return err
		}
		if diGoogle {
			if u.penyimpanan == nil {
				return ErrTransportBelumDisetujui
			}
			if err := u.penyimpanan.Hapus(ctx, objek); err != nil {
				return err
			}
		}
	}
	err = u.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := baca.HapusDokumen(ctx, tx, dokID); err != nil {
			return err
		}
		switch {
		case tertunda:
			// Efek "storage-unggah" yang MASIH ANTRE dibatalkan di transaksi yang sama (outbox inti tanpa status
			// "dibatalkan": ditutup "selesai" dengan catatan).
			return batalkanEfekUnggah(ctx, tx, u.svc.DB(), baca, strconv.FormatInt(dok.ID, 10), saat)
		case diGoogle:
			return baca.HapusObjekPenyimpanan(ctx, tx, objek.ImageID)
		default:
			// Dokumen LAMA dengan kartu stub lokal T_CLAIMLF_STORAGE (dipensiunkan untuk dokumen baru).
			return baca.HapusKartuBerkas(ctx, tx, dok.TStorageID)
		}
	})
	if err != nil {
		return err
	}
	// Berkas lokal (antrean / salinan stub) dibuang SESUDAH commit, best-effort.
	if strings.TrimSpace(u.folder) != "" {
		jalur := filepath.Join(u.folder, namaBerkasLokal(strconv.FormatInt(dok.ID, 10), dok.NamaFile))
		if e := os.Remove(jalur); e != nil && !os.IsNotExist(e) {
			log.Printf("claimlife: dokumen %d dihapus, tetapi berkas lokalnya gagal dibuang: %v", dok.ID, e)
		}
	}
	return nil
}

// BukaDokumen membuka isi dokumen untuk "View" (11-10-2026), berbatas klaim:
//   - objek di T_STORAGE_IMAGE → Google Storage (URL bertanda tangan diunduh backend; URL diperbarui bila kedaluwarsa);
//   - selain itu (belum tersimpan / dokumen lama berkartu stub lokal) → salinan lokal UNGGAHAN_DIR.
//
// Jalur lokal disusun dari pengenal dokumen, tidak pernah dari klien.
func (u *Unggahan) BukaDokumen(ctx context.Context, pelaku inti.Pelaku, klaimID string, dokID int64) (
	models.Dokumen, io.ReadCloser, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return models.Dokumen{}, nil, err
	}
	if strings.TrimSpace(klaimID) == "" {
		return models.Dokumen{}, nil, fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return models.Dokumen{}, nil, db.ErrTanpaOracle
	}
	baca := repository.NewKlaimLife(u.svc.DB())
	dok, err := baca.SatuDokumen(ctx, klaimID, dokID)
	if err != nil {
		return models.Dokumen{}, nil, err
	}
	if strings.TrimSpace(dok.TStorageID) != "" {
		o, ada, err := baca.AmbilObjekPenyimpanan(ctx, nil, dok.TStorageID)
		if err != nil {
			return models.Dokumen{}, nil, err
		}
		if ada && strings.TrimSpace(o.URLPublic) != "" {
			if u.penyimpanan == nil {
				return models.Dokumen{}, nil, ErrTransportBelumDisetujui
			}
			isi, baru, err := u.penyimpanan.Buka(ctx, o)
			if err != nil {
				return models.Dokumen{}, nil, err
			}
			if baru != nil {
				if err := baca.PerbaruiObjekPenyimpanan(ctx, nil, *baru); err != nil {
					log.Printf("claimlife: memperbarui URL objek dokumen %d gagal: %v", dok.ID, err)
				}
			}
			return dok, isi, nil
		}
	}
	if strings.TrimSpace(u.folder) == "" {
		return models.Dokumen{}, nil, unggah.ErrUnggahanDirBelumDisetel
	}
	f, err := os.Open(filepath.Join(u.folder, namaBerkasLokal(strconv.FormatInt(dok.ID, 10), dok.NamaFile)))
	if err != nil {
		return models.Dokumen{}, nil, fmt.Errorf("%w: dokumen %d", ErrBerkasDokumenTidakAda, dok.ID)
	}
	return dok, f, nil
}

// BukaDokumenLewatPengenal - "View" rute lama /api/dokumen/{id}/isi: klaim dicari dulu, lalu BukaDokumen berbatas klaim.
func (u *Unggahan) BukaDokumenLewatPengenal(ctx context.Context, pelaku inti.Pelaku, dokID int64) (
	models.Dokumen, io.ReadCloser, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return models.Dokumen{}, nil, err
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return models.Dokumen{}, nil, db.ErrTanpaOracle
	}
	klaimID, err := repository.NewKlaimLife(u.svc.DB()).KlaimDokumen(ctx, dokID)
	if err != nil {
		return models.Dokumen{}, nil, err
	}
	return u.BukaDokumen(ctx, pelaku, klaimID, dokID)
}

// CatatanBatalUnggah - GALAT_TERAKHIR efek "storage-unggah" yang dibatalkan karena dokumennya dihapus.
const CatatanBatalUnggah = "dibatalkan: dokumen dihapus sebelum diunggah ke penyimpanan"

// batalkanEfekUnggah menutup efek "storage-unggah" dokumen yang masih antre, di transaksi pemanggil.
func batalkanEfekUnggah(ctx context.Context, tx *db.Tx, d *db.DB, baca *repository.KlaimLife, rujukan string,
	saat time.Time) error {

	id, err := baca.EfekUnggahAntre(ctx, tx, ModulClaimLife, unggah.JenisEfekStorageUnggah, rujukan, outbox.StatusEfekAntre)
	if err != nil {
		return err
	}
	for _, x := range id {
		if err := outbox.NewPenyimpan(d).TuntaskanEfek(ctx, tx, x, outbox.StatusEfekSelesai, time.Time{},
			CatatanBatalUnggah, saat); err != nil {
			return err
		}
	}
	return nil
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
		baca:   repository.NewKlaimLife(svc.DB()),
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
func (p *PelaksanaBerkasLokal) Laksanakan(ctx context.Context, tx *db.Tx,
	b outbox.BarisEfekKeluar) error {

	switch b.Jenis {
	case unggah.JenisEfekStorageUnggah:
		return p.unggah(ctx, tx, b)
	case unggah.JenisEfekStorageHapus:
		return p.hapus(ctx, tx, b)
	default:
		return fmt.Errorf("%w: jenis efek %q bukan milik pelaksana berkas",
			galat.ErrPermintaanTidakSah, b.Jenis)
	}
}

// unggah menulis kartu berkas lalu menautkannya ke barisnya.
func (p *PelaksanaBerkasLokal) unggah(ctx context.Context, tx *db.Tx,
	b outbox.BarisEfekKeluar) error {

	id, err := strconv.ParseInt(b.Rujukan, 10, 64)
	if err != nil {
		// ⛔ Permanen, bukan layak dicoba ulang: rujukan yang bukan angka
		// tidak akan menjadi angka pada percobaan kedelapan.
		return fmt.Errorf("%w: rujukan %q bukan pengenal dokumen",
			galat.ErrPermintaanTidakSah, b.Rujukan)
	}
	// 11-10-2026: dokumen yang sudah DIHAPUS sebelum efeknya berjalan dilewati (dibuang berkas lokalnya, efek dianggap
	// tuntas) - berkas yang dihapus pemakai tidak pernah diunggah.
	ada, err := p.baca.DokumenAda(ctx, tx, id)
	if err != nil {
		return err
	}
	if !ada {
		if jalur := jalurDariMuatan(b.Muatan); jalur != "" {
			if e := os.Remove(jalur); e != nil && !os.IsNotExist(e) {
				log.Printf("claimlife: efek unggah dokumen %d dilewati (dokumen terhapus); berkas gagal dibuang: %v", id, e)
			}
		}
		return nil
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
	imageID, err := unggah.ImageIDBaru(saat)
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
func (p *PelaksanaBerkasLokal) hapus(ctx context.Context, tx *db.Tx,
	b outbox.BarisEfekKeluar) error {

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
func (u *Unggahan) UnduhLewatPengenal(ctx context.Context, pelaku inti.Pelaku,
	dokID int64) (models.Dokumen, string, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return models.Dokumen{}, "", err
	}
	if u == nil || u.svc == nil || !u.svc.PunyaDatabase() {
		return models.Dokumen{}, "", db.ErrTanpaOracle
	}
	klaimID, err := repository.NewKlaimLife(u.svc.DB()).KlaimDokumen(ctx, dokID)
	if err != nil {
		return models.Dokumen{}, "", err
	}
	return u.Unduh(ctx, pelaku, klaimID, dokID)
}
