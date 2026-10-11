// Package penyimpanan adalah padanan kelas Pega `ASM-FW-GISFW-Int-T_STORAGE_IMAGE`: berkas di Google Storage lewat
// Connect-REST `ServiceGoogle`, dipakai modul MANA PUN - keputusan work owner 08-10-2026 ("kalo di pega, semua
// activity ini bisa dipake disemua modul"; "lakukan dengan XML yang ada, 3 itu biarkan dulu").
//
//	InsertGoogleStorage_Act  Unggah, lalu pemanggil Catat di transaksinya (`Insert_T_Storage_SQL`)
//	GetUrlGoogleStorage_Act  Tautan; Buka = Tautan + isi dari URL bertanda tangan
//	DeleteGoogleStorage_Act  HapusObjek, lalu pemanggil HapusCatatan di transaksinya (`DeleteStorage_SQL`)
//	GeminiAIGoogle_Act       TanyaAI
//
// `API_GetUrl_Act`, `ArchiveGoogleStorage_Act`, dan `SetArchiveAttachment` BELUM ada: XML-nya belum diekspor.
//
// Sumber: salinan rule yang SAMA di banyak folder korpus - langkahnya identik (pembacaan 08-10-2026, urutan parameter
// ekspor saja yang berbeda): `InsertGoogleStorage_Act` 14 folder, `GetUrlGoogleStorage_Act` 14, `DeleteGoogleStorage_Act`
// 8, `GeminiAIGoogle_Act` NB FacIn dan RNW Fac In. Nomor langkah di bawah = nomor langkah activity-nya. Prasyarat
// Pega: kode 2 = jalankan langkah, 3 = lewati.
//
// ⛔ Alamat dari `M_LINK_SERVICE` saat jalan (`GetLinkService`, ADR-U-0013); token `GetTokenStorage_SQL` ditiru
// `inti/backend/layanan` (GCP_IMAGE, garam `STORAGE_TOKEN_SALT`). Alamat, token, dan garam tidak pernah masuk pesan
// galat. Klien HTTP-nya hanya di `kirim.go` (izin penjaga lintas aplikasi).
//
// ⚠️ PENYIMPANGAN SADAR (bug Pega diperbaiki di Go):
//  1. Berkas yang langkah 5 `Exit-Activity` lewati DIAM-DIAM (ekstensi kosong, isi kosong, MIME octet-stream) dan
//     upload yang menjawab `URLImage` kosong (langkah 12 dilewati) GAGAL TERANG di sini - di Pega pemanggilnya tetap
//     mencatat lampiran tanpa objek.
//  2. `APPNAME` selalu dibaca dari `T_FOLDER_IMAGE` (langkah 6 Pega berprasyarat `IsPEGAPROD` lewati) - pola Master
//     Product Name Life yang terbukti di DEV.
//  3. `Buka` mengunduh isi di backend: frontend dilarang membuka jendela atau URL (penjaga lintas modul
//     `modul/claimlife/frontend/unduhdokumen.test.ts`).
//  4. Hapus jarak jauh yang gagal MENAHAN catatannya (pemanggil tidak menghapus) - Pega tetap menghapus.
package penyimpanan

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/unggah"
)

var (
	// ErrBerkasDitolak - langkah 5 Insert: ekstensi kosong, isi kosong, atau jenis berkas tidak dikenal tabel
	// `GetMimeType` (422).
	ErrBerkasDitolak = errors.New("penyimpanan: the file cannot be stored")
	// ErrStorageBelumSiap - APPNAME, token, atau alamat `M_LINK_SERVICE` tidak tersedia (503).
	ErrStorageBelumSiap = errors.New("penyimpanan: the storage service is not ready")
	// ErrStorageGagal - layanan tak terjangkau, menjawab galat, atau jawabannya kosong/rusak (502).
	ErrStorageGagal = errors.New("penyimpanan: the storage service failed")
	// ErrObjekTidakAda - IMAGEID tidak tercatat di `T_STORAGE_IMAGE` (409).
	ErrObjekTidakAda = errors.New("penyimpanan: the file is not recorded in storage")
	// ErrBerkasTidakDiStorage - URL bertanda tangan menjawab 404 (409).
	ErrBerkasTidakDiStorage = errors.New("penyimpanan: the file is no longer in storage")
)

// Objek adalah satu baris `T_STORAGE_IMAGE`.
type Objek struct {
	ImageID   string
	URLPublic string
	// AppFolder - jalur objek penuh, `gs://<App>/<Folder><Namafile>`.
	AppFolder string
	// Exp - `EXPDATE` dalam bentuk tulisnya, `DD/MM/YYYY HH24:MI:SS`.
	Exp      string
	FileName string
	AppName  string
	// TanggalUpload - `DateTime` jawaban geturl (`MM/DD/YYYY HH24:MI:SS`), hanya untuk Perbarui.
	TanggalUpload string
}

// Catatan adalah `T_STORAGE_IMAGE` dan `T_FOLDER_IMAGE` (Oracle di `catatan.go`, tiruan di uji).
type Catatan interface {
	// NamaAplikasi - `GetAppName_SQL`.
	NamaAplikasi(ctx context.Context) (string, error)
	// Ambil - `GetLinkStorage_SQL`; ada false = IMAGEID tidak tercatat.
	Ambil(ctx context.Context, imageID string) (Objek, bool, error)
	// Catat - `Insert_T_Storage_SQL`, di transaksi pemanggil.
	Catat(ctx context.Context, tx *db.Tx, o Objek) error
	// Perbarui - `Update_T_Storage_SQL`.
	Perbarui(ctx context.Context, o Objek) error
	// Hapus - `DeleteStorage_SQL`, di transaksi pemanggil.
	Hapus(ctx context.Context, tx *db.Tx, imageID string) error
}

// SumberToken - `GetTokenStorage_SQL` (`GET_TOKEN_STORAGE(App, OperatorID.pyUserIdentifier)`).
type SumberToken func(ctx context.Context, app, pengguna string) (string, error)

// SumberAlamat - `GetLinkService` (`M_LINK_SERVICE`).
type SumberAlamat func(ctx context.Context, kunci layanan.KunciLayanan) (string, error)

// Penyimpanan menjalankan keempat activity.
type Penyimpanan struct {
	catatan Catatan
	token   SumberToken
	alamat  SumberAlamat
	klien   *http.Client
	jam     func() time.Time
}

// Baru menyusun penyimpanan di atas sumbernya; klien nil = bawaan (`kirim.go`), jam nil = `time.Now`.
func Baru(c Catatan, token SumberToken, alamat SumberAlamat, klien *http.Client, jam func() time.Time) *Penyimpanan {
	if klien == nil {
		klien = klienBawaan()
	}
	if jam == nil {
		jam = time.Now
	}
	return &Penyimpanan{catatan: c, token: token, alamat: alamat, klien: klien, jam: jam}
}

// Oracle menyusun penyimpanan di atas Oracle akar `akar`; garam dari `config.StorageTokenSalt`.
func Oracle(akar inti.Akar, garam string) *Penyimpanan {
	penerbit := layanan.TokenStorageOracle(akar, garam)
	resolver := layanan.ResolverLinkServiceOracle(akar)
	token := func(ctx context.Context, app, pengguna string) (string, error) {
		var kode string
		err := akar.DalamTransaksi(ctx, func(tx *db.Tx) error {
			var e error
			kode, e = penerbit.Token(ctx, tx, app, pengguna, time.Now())
			return e
		})
		return kode, err
	}
	alamat := func(ctx context.Context, k layanan.KunciLayanan) (string, error) {
		return layanan.AlamatLayanan(ctx, resolver, k)
	}
	return Baru(catatanOracle{db: akar.DB()}, token, alamat, nil, nil)
}

// zonaJakarta - `@CurrentDate(…, "Asia/Jakarta")`; WIB tanpa musim panas.
var zonaJakarta = time.FixedZone("WIB", 7*3600)

// formatExp - bentuk `EXPDATE` yang ditulis dan dibaca kembali.
const formatExp = "02/01/2006 15:04:05"

// MasukUnggah adalah parameter `InsertGoogleStorage_Act`.
type MasukUnggah struct {
	// Folder - `Param.Folder` (mis. "Contract").
	Folder string
	// NamaFile - `Param.Namafile`, nama asli berkas.
	NamaFile string
	// Isi - `Param.Image` sebelum base64.
	Isi []byte
	// Durasi - `Param.Durasi`, detik.
	Durasi int
	// Ext - `Param.Ext`; kosong = dari NamaFile (langkah 2).
	Ext string
	// Pengguna - akun pengunggah (`USERINPUT` token).
	Pengguna string
}

// Ekstensi - langkah 2: sesudah titik terakhir nama berkas; kosong bila tidak ada.
func Ekstensi(nama string) string {
	i := strings.LastIndex(nama, ".")
	if i < 0 || i == len(nama)-1 {
		return ""
	}
	return strings.ToLower(nama[i+1:])
}

// namaObjek - langkah 8: `Param.Folder + "/Doc/" + YYYY + "/" + MM + "/"` dan
// `@CurrentDate("yyyyMMdd-hhmmss-S") + " - " + Param.Namafile` (pola Java: jam 12-an, milidetik tanpa nol depan).
func namaObjek(saat time.Time, folder, nama string) (string, string) {
	w := saat.In(zonaJakarta)
	f := folder + "/Doc/" + w.Format("2006") + "/" + w.Format("01") + "/"
	n := fmt.Sprintf("%s-%d - %s", w.Format("20060102-030405"), w.Nanosecond()/int(time.Millisecond), nama)
	return f, n
}

// Unggah - `InsertGoogleStorage_Act` langkah 1-12: berkas dikirim ke penyimpanan dan objeknya dijawab untuk dicatat
// pemanggil (`Catat`) di transaksi yang sama dengan rekam lampirannya.
func (p *Penyimpanan) Unggah(ctx context.Context, m MasukUnggah) (Objek, error) {
	ext := strings.TrimSpace(m.Ext)
	if ext == "" {
		ext = Ekstensi(m.NamaFile) // 2
	}
	ext = strings.ToLower(ext)                 // 3
	mime := unggah.MimeDariNamaFile("." + ext) // 4 `GetMimeType`
	image := base64.StdEncoding.EncodeToString(m.Isi)
	// 5 `Exit-Activity` - di sini gagal terang (penyimpangan 1).
	if ext == "" || len(image) < 10 || mime == unggah.MimeBawaan {
		return Objek{}, fmt.Errorf("%w: %q is empty or of a type that storage does not accept", ErrBerkasDitolak, m.NamaFile)
	}
	app, err := p.namaAplikasi(ctx) // 6
	if err != nil {
		return Objek{}, err
	}
	saat := p.jam()
	folder, nama := namaObjek(saat, m.Folder, m.NamaFile) // 8
	durasi := m.Durasi
	// 7 token, 9 SET JSON (Durasi tanpa kutip), 10 GetLinkService "Google"/"upload", 11 Connect-REST.
	j, err := p.panggil(ctx, layanan.KunciUnggahBerkas, app, m.Pengguna, func(kode string) any {
		return permintaanBerkas{App: app, Kodestring: kode, Durasi: &durasi, Folder: folder, Namafile: nama, Image: image,
			Ext: ext, MimeType: mime}
	})
	if err != nil {
		return Objek{}, err
	}
	// 12 "Insert ke table" berprasyarat `URLImage == ""` lewati - di sini gagal terang (penyimpangan 1).
	if strings.TrimSpace(j.URLImage) == "" {
		return Objek{}, gagal("upload answered an empty URLImage")
	}
	id, err := unggah.ImageIDBaru(saat) // 12.2 `GenerateImageID_SQL`
	if err != nil {
		return Objek{}, err
	}
	// 12.3 InsertDoc: exp, URLImage, appfolder, App, Namafile.
	o := Objek{ImageID: id, URLPublic: j.URLImage, AppFolder: j.AppFolder, Exp: expPega(j.Exp), FileName: nama, AppName: app}
	if strings.TrimSpace(o.AppFolder) == "" {
		// Jawaban tanpa `appfolder`: dirakit seperti bentuknya di DEV - geturl dan hapus menuntut jalur objek.
		o.AppFolder = awalanGS(app) + folder + nama
	}
	return o, nil
}

// Catat - `Insert_T_Storage_SQL` (12.4) di transaksi pemanggil.
func (p *Penyimpanan) Catat(ctx context.Context, tx *db.Tx, o Objek) error {
	return p.catatan.Catat(ctx, tx, o)
}

// HapusCatatan - `DeleteStorage_SQL` (Delete langkah 10) di transaksi pemanggil.
func (p *Penyimpanan) HapusCatatan(ctx context.Context, tx *db.Tx, imageID string) error {
	return p.catatan.Hapus(ctx, tx, imageID)
}

// berlaku - prasyarat langkah 6 GetUrl `@CompareDates(UploadDoc.exp, @CurrentDateTime())` (benar = lewati): URL
// tersimpan dipakai selama EXPDATE SESUDAH sekarang; kosong = minta baru. EXPDATE dibaca jam Asia/Jakarta - bukti
// DEV Master Product Name Life 03-10-2026.
func (p *Penyimpanan) berlaku(o Objek) bool {
	exp, err := time.ParseInLocation(formatExp, strings.TrimSpace(o.Exp), zonaJakarta)
	return err == nil && exp.After(p.jam())
}

// Tautan - `GetUrlGoogleStorage_Act`: URL bertanda tangan objek `imageID`, berlaku `durasi` detik.
func (p *Penyimpanan) Tautan(ctx context.Context, imageID string, durasi int, pengguna string) (string, error) {
	o, ada, err := p.catatan.Ambil(ctx, imageID) // 4 GET LINK
	if err != nil {
		return "", err
	}
	if !ada {
		return "", ErrObjekTidakAda
	}
	if p.berlaku(o) {
		return o.URLPublic, nil // 6 dilewati -> 7
	}
	// 6.2 `@replaceAll(.Folder, .Namafile, "")`, `@replaceAll(.Folder, "gs://"+.App+"/", "")`.
	folder := strings.ReplaceAll(strings.ReplaceAll(o.AppFolder, o.FileName, ""), awalanGS(o.AppName), "")
	j, err := p.panggil(ctx, layanan.KunciURLBerkas, o.AppName, pengguna, func(kode string) any {
		return permintaanBerkas{App: o.AppName, Kodestring: kode, Durasi: &durasi, Folder: folder, Namafile: o.FileName}
	})
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(j.URLImage) == "" {
		return o.URLPublic, nil // 6.6 dilewati -> 7: URL tersimpan
	}
	// 6.6.2 UpdateDoc dari jawaban; `appfolder` kosong tidak menimpa jalur objek yang tercatat.
	segar := o
	segar.URLPublic, segar.Exp, segar.TanggalUpload = j.URLImage, expPega(j.Exp), j.DateTime
	if strings.TrimSpace(j.AppFolder) != "" {
		segar.AppFolder = j.AppFolder
	}
	if err := p.catatan.Perbarui(ctx, segar); err != nil { // 6.6.3
		// URL baru tetap dijawab; catatan lama hanya membuat geturl berikutnya terulang.
		log.Printf("penyimpanan: memperbarui T_STORAGE_IMAGE %s: %v", imageID, err)
	}
	return j.URLImage, nil
}

// Buka - Tautan lalu isi dari URL bertanda tangan (penyimpangan 3).
func (p *Penyimpanan) Buka(ctx context.Context, imageID string, durasi int, pengguna string) (io.ReadCloser, error) {
	bertanda, err := p.Tautan(ctx, imageID, durasi, pengguna)
	if err != nil {
		return nil, err
	}
	return p.unduh(ctx, bertanda)
}

// HapusObjek - `DeleteGoogleStorage_Act` langkah 4-9: objek dihapus dari penyimpanan. IMAGEID yang tidak tercatat
// tidak punya objek untuk dihapus - bukan galat.
func (p *Penyimpanan) HapusObjek(ctx context.Context, imageID, pengguna string) error {
	o, ada, err := p.catatan.Ambil(ctx, imageID) // 4 GET DATA
	if err != nil || !ada {
		return err
	}
	// 6 `ParamJSON.Namafile = @replaceAll(.Folder, "gs://"+.App+"/", "")`.
	jalur := strings.ReplaceAll(o.AppFolder, awalanGS(o.AppName), "")
	if strings.TrimSpace(jalur) == "" {
		return gagal("APPFOLDER is empty; the storage object cannot be named")
	}
	_, err = p.panggil(ctx, layanan.KunciHapusBerkas, o.AppName, pengguna, func(kode string) any {
		return permintaanBerkas{App: o.AppName, Kodestring: kode, Namafile: jalur}
	})
	return err
}

// MasukAI adalah parameter `GeminiAIGoogle_Act`.
type MasukAI struct {
	// ImageID - objek gambar yang ikut ditanyakan; kosong = tanpa gambar.
	ImageID string
	// Pesan - `Param.Pesan`, pertanyaannya.
	Pesan    string
	Pengguna string
}

// suhuAI - langkah 8 `ParamJSON.Temp = 0.7`, dikirim tanpa kutip (langkah 9).
const suhuAI = 0.7

// TanyaAI - `GeminiAIGoogle_Act`: pertanyaan (dan gambar dari penyimpanan) ke `ServiceGoogle` "getAI"; jawabannya
// `Param.Hasil` (langkah 12).
func (p *Penyimpanan) TanyaAI(ctx context.Context, m MasukAI) (string, error) {
	app, err := p.namaAplikasi(ctx) // 4
	if err != nil {
		return "", err
	}
	var namafile []string
	if strings.TrimSpace(m.ImageID) != "" {
		o, ada, err := p.catatan.Ambil(ctx, m.ImageID) // 5 JIKA PAKE GAMBAR
		if err != nil {
			return "", err
		}
		if !ada {
			return "", ErrObjekTidakAda
		}
		if strings.TrimSpace(o.AppName) != "" {
			app = o.AppName // `GetLinkStorage_SQL` mengisi `UploadDoc.App`
		}
		namafile = []string{o.AppFolder} // 7 `ListNamafile(<APPEND>) = UploadDoc.Folder`
	}
	// 6 token, 8 SET DATA, 9 JSON (`ListNamafile` -> `Namafile`), 10 "Google"/"getAI", 11 Connect-REST.
	j, err := p.panggil(ctx, layanan.KunciTanyaAI, app, m.Pengguna, func(kode string) any {
		return permintaanAI{App: app, Kodestring: kode, Tanya: m.Pesan, Temp: suhuAI, Namafile: namafile}
	})
	if err != nil {
		return "", err
	}
	return j.Hasil, nil
}

// namaAplikasi - `GetAppName_SQL`; kosong = belum siap (penyimpangan 2).
func (p *Penyimpanan) namaAplikasi(ctx context.Context) (string, error) {
	app, err := p.catatan.NamaAplikasi(ctx)
	if err != nil {
		return "", belumSiap("APPNAME (T_FOLDER_IMAGE)", err)
	}
	if strings.TrimSpace(app) == "" {
		return "", belumSiap("APPNAME (T_FOLDER_IMAGE)", layanan.ErrAppNameKosong)
	}
	return app, nil
}

// awalanGS - awalan objek di APPFOLDER: skema gs, inang App, garis miring.
func awalanGS(app string) string { return (&url.URL{Scheme: "gs", Host: app, Path: "/"}).String() }

// polaExp - `YYYYMMDDTHHMMSS` di awal `exp` sesudah `-`/`:` dibuang.
var polaExp = regexp.MustCompile(`^\d{8}T\d{6}`)

// expPega - `exp` jawaban -> `To_date(exp, 'DD/MM/YYYY HH24:MI:SS')`: `@substring(@pxReplaceAllViaRegex(exp, "[-:]",
// ""), 0, 19) + " GMT"` lalu `@FormatDateTime(…, "dd/MM/yyyy HH:mm:ss", "GMT")`. Bentuk tak terbaca -> kosong (NULL).
func expPega(exp string) string {
	awal := polaExp.FindString(strings.NewReplacer("-", "", ":", "").Replace(strings.TrimSpace(exp)))
	if awal == "" {
		return ""
	}
	t, err := time.Parse("20060102T150405", awal)
	if err != nil {
		return ""
	}
	return t.Format(formatExp)
}
