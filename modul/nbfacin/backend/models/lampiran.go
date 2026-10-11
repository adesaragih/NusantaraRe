package models

// Lampiran case NB FacIn (tiket 67) - "Attachment File" Pega (`Section\AttachmentGridReas.xml`): bentuk baris
// `POOLDATA.DOCUMENT_POLIS` + `T_STORAGE_IMAGE` (DDL `D:\migrasi\RNM\DDL\DOCUMENT_POLIS.txt` / `T_STORAGE_IMAGE.txt`) dan
// aturan murni di sekitarnya: ID dokumen, jalur objek Google Storage, status, kategori.

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ZonaWIB - `"Asia/Jakarta"` rumus @CurrentDate Pega (InsertDocument_Act / InsertGoogleStorage_Act).
var ZonaWIB = time.FixedZone("WIB", 7*60*60)

// Status satu berkas di popup View File (kunci kabel `status`, keputusan 9d 08-10-2026 K-C).
const (
	StatusLampiranTersimpan = ""
	StatusLampiranMenunggu  = "menunggu"
	StatusLampiranGagal     = "gagal"
)

// KategoriLampiran - satu baris grid: `CATEGORY_ATTACH_REAS.NOTE` + Count (`InputParamUploadReas_act`).
type KategoriLampiran struct {
	Kategori string
	Jumlah   int
}

// DokumenPolis - satu baris `DOCUMENT_POLIS` (kolom GetAllDocument + PXCREATEOPERATOR) dan nasib unggahnya.
type DokumenPolis struct {
	ID         string
	IDPega     string
	NamaFile   string
	Mime       string // ekstensi huruf kecil (InsertDocument_Act `@toLowerCase(Param.MIME)`, K-E)
	Kategori1  string
	Kategori2  string
	TStorageID string
	Operator   string
	// Tanggal - `TANGGAL` sebagai teks kabel `DD-MM-YYYY HH:MM` (dibentuk Oracle).
	Tanggal string
	// Tersimpan - baris `T_STORAGE_IMAGE` untuk `T_STORAGE_ID` ada.
	Tersimpan bool
	// StatusEfek - status efek `storage-unggah` TERAKHIR dokumen ini di outbox ("" = tidak ada efek).
	StatusEfek string
}

// ObjekStorage - satu baris `T_STORAGE_IMAGE` (`Insert_T_Storage_SQL` / `Update_T_Storage_SQL`): teks APA ADANYA dari
// jawaban layanan; `Exp` berbentuk `DD/MM/YYYY HH24:MI:SS` (GMT), `TanggalUpload` `MM/DD/YYYY HH24:MI:SS`.
type ObjekStorage struct {
	ImageID       string
	URLPublic     string
	AppFolder     string
	Exp           string
	Namafile      string
	App           string
	TanggalUpload string
}

// LinkStorage - `GetLinkStorage_SQL` (URLPUBLIC, APPFOLDER, EXPDATE, APPNAME, FILENAME) satu IMAGEID; `ExpUTC` =
// EXPDATE `YYYYMMDDHH24MISS` (kosong = NULL).
type LinkStorage struct {
	ImageID   string
	URLPublic string
	AppFolder string
	ExpUTC    string
	App       string
	Namafile  string
}

// LokasiObjek - letak objek di layanan penyimpanan: `Folder` relatif berakhiran "/" (tanpa "gs://<App>/") dan `Namafile`
// (InsertGoogleStorage_Act langkah 8); `ImageID` = kunci stub lokal.
type LokasiObjek struct {
	ImageID  string
	Folder   string
	Namafile string
}

// JalurPenuh - `Namafile` panggilan delete (DeleteGoogleStorage_Act langkah 6: APPFOLDER tanpa "gs://<App>/").
func (l LokasiObjek) JalurPenuh() string { return l.Folder + l.Namafile }

// IDDokumenPolis - InsertDocument_Act langkah 3 `.ID = @CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")` VERBATIM: `hh` =
// jam 12-an (01-12), jadi 01:00 dan 13:00 berbentuk sama - bentrok ditangani penulisnya (K-D).
func IDDokumenPolis(saat time.Time) string {
	w := saat.In(ZonaWIB)
	return w.Format("20060102030405") + fmt.Sprintf("%03d", w.Nanosecond()/int(time.Millisecond))
}

// polaIDDokumen - bentuk ID buatan IDDokumenPolis.
var polaIDDokumen = regexp.MustCompile(`^\d{17}$`)

// FolderStoragePolicy - InsertDocument_Act langkah 4 `Folder = "Policy"`.
const FolderStoragePolicy = "Policy"

// tahunMingguUS - `@CurrentDate("YYYY", …)`: pola Java `YYYY` = TAHUN-MINGGU, bukan tahun kalender. `[dugaan]` lokal JVM
// en_US (minggu mulai Minggu, minggu pertama = minggu yang memuat 1 Januari): hari Desember di minggu yang sama dengan
// 1 Januari ikut tahun berikutnya. Lokal server Pega belum terverifikasi; jalur ini hanya nama folder (Pega membaca ulang
// APPFOLDER dari T_STORAGE_IMAGE, tidak menghitungnya).
func tahunMingguUS(t time.Time) int {
	if t.Month() == time.December {
		sabtu := t.AddDate(0, 0, 6-int(t.Weekday()))
		if sabtu.Year() > t.Year() {
			return t.Year() + 1
		}
	}
	return t.Year()
}

// LokasiUnggah - InsertGoogleStorage_Act langkah 8 untuk dokumen ber-ID buatan IDDokumenPolis: `.Folder = Param.Folder +
// "/Doc/" + YYYY + "/" + MM + "/"`, `.Namafile = @CurrentDate("yyyyMMdd-hhmmss-S") + " - " + Param.Namafile`. Waktunya
// diambil dari ID dokumen (bukan jam unggah) supaya setiap percobaan ulang menulis ke objek yang SAMA; `S` = milidetik
// tanpa nol di depan. ID berbentuk lain -> false.
func LokasiUnggah(idDokumen, imageID, namaFile string) (LokasiObjek, bool) {
	if !polaIDDokumen.MatchString(idDokumen) {
		return LokasiObjek{}, false
	}
	tgl, err := time.ParseInLocation("20060102", idDokumen[:8], ZonaWIB)
	if err != nil {
		return LokasiObjek{}, false
	}
	ms, _ := strconv.Atoi(idDokumen[14:17])
	folder := fmt.Sprintf("%s/Doc/%04d/%02d/", FolderStoragePolicy, tahunMingguUS(tgl), int(tgl.Month()))
	nama := idDokumen[:8] + "-" + idDokumen[8:14] + "-" + strconv.Itoa(ms) + " - " + namaFile
	return LokasiObjek{ImageID: imageID, Folder: folder, Namafile: nama}, true
}

// awalanBucket - awalan URI objek Google Storage `gs` + App + "/" yang Pega buang dari APPFOLDER (literal Pega
// `"gs://"+.App+"/"`). Dirakit pustaka standar dari skema + host (App dari T_FOLDER_IMAGE, saat jalan) - bukan alamat
// layanan (ADR-0013: alamat layanan tetap dari M_LINK_SERVICE). App kosong = tanpa awalan.
func awalanBucket(app string) string {
	if strings.TrimSpace(app) == "" {
		return ""
	}
	return (&url.URL{Scheme: "gs", Host: app}).String() + "/"
}

// buangAwalanBucket - `@replaceAll(x, awalanBucket, "")`.
func buangAwalanBucket(s, app string) string {
	if a := awalanBucket(app); a != "" {
		return strings.ReplaceAll(s, a, "")
	}
	return s
}

// LokasiDariLink - GetUrlGoogleStorage_Act langkah 6.2: `.Folder = @replaceAll(.Folder, .Namafile, "")` lalu awalan bucket
// dibuang, atas APPFOLDER tersimpan. `@replaceAll` Pega berbasis regex; di sini penggantian teks apa adanya `[dugaan]` (sama
// kecuali nama berkas memuat karakter regex).
func LokasiDariLink(l LinkStorage) LokasiObjek {
	folder := l.AppFolder
	if l.Namafile != "" {
		folder = strings.ReplaceAll(folder, l.Namafile, "")
	}
	return LokasiObjek{ImageID: l.ImageID, Folder: buangAwalanBucket(folder, l.App), Namafile: l.Namafile}
}

// JalurHapusDariLink - DeleteGoogleStorage_Act langkah 6 `Namafile` = APPFOLDER tanpa awalan bucket.
func JalurHapusDariLink(l LinkStorage) string { return buangAwalanBucket(l.AppFolder, l.App) }

// AppFolderUji - APPFOLDER berbentuk jawaban layanan (awalan bucket + jalur) - dipakai uji.
func AppFolderUji(app, jalur string) string { return awalanBucket(app) + jalur }

// StatusDokumen - tersimpan MENANG (objek tercatat); efek unggah terakhir gagal-permanen atau tidak ada sama sekali
// (dokumen tanpa objek yang tidak pernah diantre sistem ini) = gagal; selain itu menunggu.
func StatusDokumen(d DokumenPolis, statusGagalPermanen string) string {
	switch {
	case d.Tersimpan:
		return StatusLampiranTersimpan
	case d.StatusEfek == "" || d.StatusEfek == statusGagalPermanen:
		return StatusLampiranGagal
	default:
		return StatusLampiranMenunggu
	}
}

// KategoriLampiranSah - pilihan pemakai dicocokkan (tanpa beda huruf) dengan master; teks MASTER yang disimpan.
func KategoriLampiranSah(master []string, pilihan string) (string, bool) {
	p := strings.TrimSpace(pilihan)
	if p == "" {
		return "", false
	}
	for _, m := range master {
		if strings.EqualFold(strings.TrimSpace(m), p) {
			return m, true
		}
	}
	return "", false
}

// EkstensiBerkas - ekstensi huruf kecil tanpa titik dari titik TERAKHIR nama dasar; "" bila tidak ada.
func EkstensiBerkas(nama string) string {
	e := strings.ToLower(path.Ext(NamaDasarBerkas(nama)))
	return strings.TrimPrefix(e, ".")
}

// NamaDasarBerkas - buang jalur yang dikirim sebagian peramban.
func NamaDasarBerkas(nama string) string {
	d := path.Base(strings.ReplaceAll(strings.TrimSpace(nama), `\`, "/"))
	if d == "." || d == "/" || d == ".." {
		return ""
	}
	return d
}

// akhiranSah - akhiran berkas antrean di disk: huruf kecil dan angka, 1-10.
var akhiranSah = regexp.MustCompile(`^[a-z0-9]{1,10}$`)

// NamaBerkasAntreLampiran - nama berkas antrean di UNGGAHAN_DIR: hanya kunci (IMAGEID heksa) + akhiran yang lolos
// saringan; nama unggahan tidak pernah menentukan jalur.
func NamaBerkasAntreLampiran(imageID, namaAsli string) string {
	if e := EkstensiBerkas(namaAsli); akhiranSah.MatchString(e) {
		return imageID + "." + e
	}
	return imageID
}

// ExpStorage - `exp` jawaban layanan -> `dd/MM/yyyy HH:mm:ss` GMT (InsertGoogleStorage_Act 12.3 / GetUrlGoogleStorage_Act
// 6.6.2: `-`/`:` dibuang, 19 karakter + " GMT"). Bentuk tak terbaca -> "" (NULL), bukan galat To_date.
func ExpStorage(exp string) string {
	polos := strings.NewReplacer("-", "", ":", "").Replace(strings.TrimSpace(exp))
	if len(polos) < 15 {
		return ""
	}
	t, err := time.Parse("20060102T150405", polos[:15])
	if err != nil {
		return ""
	}
	return t.Format("02/01/2006 15:04:05")
}

// TanggalUploadStorage - `DateTime` geturl apa adanya bila berbentuk `MM/DD/YYYY HH24:MI:SS` (To_date
// Update_T_Storage_SQL); selain itu "" (NULL).
func TanggalUploadStorage(teks string) string {
	t := strings.TrimSpace(teks)
	if _, err := time.Parse("01/02/2006 15:04:05", t); err != nil {
		return ""
	}
	return t
}

// LinkBerlaku - GetUrlGoogleStorage_Act langkah 6 `@CompareDates(UploadDoc.exp, @CurrentDateTime())`: URL tersimpan dipakai
// selama EXPDATE (GMT `[dugaan]`) masih sesudah `kini + sisa`. EXPDATE kosong / tak terbaca = tidak berlaku.
func LinkBerlaku(l LinkStorage, kini time.Time, sisa time.Duration) bool {
	if strings.TrimSpace(l.URLPublic) == "" || len(l.ExpUTC) != 14 {
		return false
	}
	exp, err := time.ParseInLocation("20060102150405", l.ExpUTC, time.UTC)
	if err != nil {
		return false
	}
	return exp.After(kini.Add(sisa))
}
