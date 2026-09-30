package services

// Transport penyimpanan lampiran - OQ-TCO-08 Treaty Contract Out.
//
// Untuk apa berkas ini: implementasi `PengirimBerkasTCO` yang benar-benar
// berbicara HTTP dengan layanan penyimpanan - [keputusan work owner
// 29-09-2026, OQ-TCO-08 "sekarang"]. Satu-satunya berkas modul yang memegang
// klien HTTP keluar; penjaga `TestNolAlamatLayananDiKode` dan
// `TestTCOLampiranTanpaAlamatLiteral` mengecualikannya HANYA dari pemeriksaan
// klien HTTP (tetap diperiksa untuk skema-alamat literal, host, dan env).
//
// Bentuk panggilan ditiru dari korpus (tanpa menyalin alamat apa pun):
//
//	`ConnectREST/ServiceGoogle.xml`: POST, `Content-Type: application/json`
//	(b148), badan = `UploadDoc.JSON` (b221), jawaban -> `UploadDoc.Response`
//	(b320-b322), batas waktu 300000 ms (b29), alamat dasar dari setelan -
//	di sini dari `M_LINK_SERVICE` saat jalan (ADR-0013).
//	`Claim Life/Activity/InsertGoogleStorage_Act.xml`: badan unggah = halaman
//	`UploadDoc` (`App`, `Kodestring`, `Durasi`, `Folder`, `Namafile`, `Image`,
//	`ext`, `MimeType`; b1339-b1641); gagal bila `URLImage` kosong (b2902).
//	`Treaty Contract Out/Activity/GetUrlGoogleStorage_Act.xml` (b1346-b1431:
//	`App`, `Kodestring`, `Folder`, `Namafile`, `Durasi`; `Folder` = jalur folder
//	objek berakhiran `/`, b1260-b1326) dan `DeleteGoogleStorage_Act.xml`
//	(b1019-b1091: `App`, `Kodestring`, `Namafile` = jalur objek PENUH
//	`<folder>/<nama>`, tanpa `Folder`).
//
// ⛔ Alamat, token, dan garam TIDAK PERNAH masuk pesan galat maupun log:
// galat jaringan diringkas tanpa rinciannya (alamat IP/inang ada di sana).
// ⛔ Tidak dipanggil dari uji atau dari sesi pengembangan: uji memakai server
// tiruan lokal; penyambungan sungguhan terjadi saat `PELAKSANA_STORAGE=nyata`.
//
// Dibaca sesudah: tco_penyimpanan.go.

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nusantarare/inti/outbox"
	"nusantarare/inti/unggah"
	"nusantarare/modul/treatycontractout/models"
)

// BatasWaktuStorageTCO - `ServiceGoogle.xml` b29 `pyResponseTimeout 300000`.
const BatasWaktuStorageTCO = 300 * time.Second

// FolderStorageTCO - folder objek lampiran modul ini di layanan penyimpanan,
// BERAKHIRAN `/` seperti Pega (Claim Life `InsertGoogleStorage_Act` b1407-b1408
// `Param.Folder+"/Doc/"+YYYY+"/"+MM+"/"`). Jalur objek = folder + `IMAGEID`.
// `[terbuka — OQ-TCO-22]` untuk namanya: korpus Treaty Contract Out
// tidak memuat unggahnya.
const FolderStorageTCO = "TreatyContractOut/"

// DurasiURLStorageTCO - `Durasi` URL bertanda tangan (satuan menurut layanan).
// `[terbuka — OQ-TCO-22]`: nilainya dari pemanggil di Pega.
const DurasiURLStorageTCO = 60

// batasJawabanStorageTCO - jawaban JSON layanan tidak dibaca tanpa batas.
const batasJawabanStorageTCO = 1 << 20

var (
	// ErrStorageTakTerjangkauTCO - jaringan/batas waktu; layak dicoba ulang.
	ErrStorageTakTerjangkauTCO = errors.New("services: storage service unreachable")
	// ErrStorageGagalTCO - layanan menjawab galat (5xx, 404 titik layanan, 3xx,
	// 401, 403, 429); layak dicoba ulang.
	ErrStorageGagalTCO = errors.New("services: storage service returned an error")
	// ErrStorageTokenDitolakTCO - 401/403: cache token dikosongkan sebelum
	// percobaan berikutnya. Membungkus ErrStorageGagalTCO.
	ErrStorageTokenDitolakTCO = fmt.Errorf("%w: token rejected", ErrStorageGagalTCO)
	// ErrStorageMenolakPermintaanTCO - layanan menolak bentuk permintaan (400, 422);
	// PERMANEN sampai manusia bertindak.
	ErrStorageMenolakPermintaanTCO = errors.New("services: storage service rejected the request")
	// ErrStorageJawabanRusakTCO - jawaban tidak berbentuk (`URLImage` kosong, JSON rusak).
	ErrStorageJawabanRusakTCO = errors.New("services: storage service response is malformed")
)

// permintaanStorageTCO - halaman `UploadDoc` sebagai JSON (nama medan VERBATIM).
type permintaanStorageTCO struct {
	App        string `json:"App"`
	Kodestring string `json:"Kodestring"`
	Durasi     *int   `json:"Durasi,omitempty"`
	Folder     string `json:"Folder,omitempty"`
	Namafile   string `json:"Namafile"`
	Image      string `json:"Image,omitempty"`
	Ext        string `json:"ext,omitempty"`
	MimeType   string `json:"MimeType,omitempty"`
}

// jawabanStorageTCO - `UploadDoc.Response` (`URLImage`, `exp`, `appfolder`, `DateTime`).
type jawabanStorageTCO struct {
	URLImage  string `json:"URLImage"`
	Exp       string `json:"exp"`
	AppFolder string `json:"appfolder"`
	DateTime  string `json:"DateTime"`
}

// PembacaAppStorageTCO memberi `App` (`T_FOLDER_IMAGE.APPNAME`) saat jalan.
type PembacaAppStorageTCO func(ctx context.Context) (string, error)

type pengirimBerkasHTTPTCO struct {
	klien *http.Client
	app   PembacaAppStorageTCO
}

// NewPengirimBerkasHTTPTCO menyusun transport HTTP; klien nil -> klien bawaan
// modul dengan batas waktu `BatasWaktuStorageTCO`.
func NewPengirimBerkasHTTPTCO(klien *http.Client, app PembacaAppStorageTCO) PengirimBerkasTCO {
	if klien == nil {
		klien = &http.Client{Timeout: BatasWaktuStorageTCO}
	}
	return &pengirimBerkasHTTPTCO{klien: klien, app: app}
}

// galatJaringanTCO meringkas galat klien HTTP TANPA rinciannya: `url.Error`
// dan galat dial memuat alamat, inang, dan IP.
func galatJaringanTCO(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) && ue.Timeout() {
		return fmt.Errorf("%w: timeout", ErrStorageTakTerjangkauTCO)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: request cancelled", ErrStorageTakTerjangkauTCO)
	}
	return ErrStorageTakTerjangkauTCO
}

// galatStatusTCO - jawaban TITIK LAYANAN (upload/geturl/delete).
//
// ⛔ 404 dari titik layanan BUKAN "berkas tidak ada": ia juga jawaban jalur
// `M_LINK_SERVICE` yang salah, dan menelannya membuat hapus berhasil diam-diam.
// "Tidak ada" hanya dibaca dari URL bertanda tangan (`galatStatusObjekTCO`).
func galatStatusTCO(kode int) error {
	switch {
	case kode == http.StatusBadRequest || kode == http.StatusUnprocessableEntity:
		return fmt.Errorf("%w: status %d", ErrStorageMenolakPermintaanTCO, kode)
	case kode == http.StatusUnauthorized || kode == http.StatusForbidden:
		return fmt.Errorf("%w: status %d", ErrStorageTokenDitolakTCO, kode)
	}
	return fmt.Errorf("%w: status %d", ErrStorageGagalTCO, kode)
}

// galatStatusObjekTCO - jawaban URL bertanda tangan: 404 = objek tidak ada.
func galatStatusObjekTCO(kode int) error {
	if kode == http.StatusNotFound {
		return ErrBerkasTidakAdaDiPenyimpanan
	}
	if kode == http.StatusUnauthorized || kode == http.StatusForbidden {
		// Tanda tangan URL ditolak - bukan token kita.
		return fmt.Errorf("%w: status %d", ErrStorageGagalTCO, kode)
	}
	return galatStatusTCO(kode)
}

// kirimJSON mengirim satu permintaan JSON dan membaca jawabannya.
func (p *pengirimBerkasHTTPTCO) kirimJSON(ctx context.Context, alamat string, badan permintaanStorageTCO) (jawabanStorageTCO, error) {
	isi, err := json.Marshal(badan)
	if err != nil {
		return jawabanStorageTCO{}, fmt.Errorf("services: building storage request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(isi))
	if err != nil {
		// ⛔ Alamat tidak disebut - pesan `url.Parse` memuatnya.
		return jawabanStorageTCO{}, fmt.Errorf("%w: service address is malformed", ErrStorageMenolakPermintaanTCO)
	}
	req.Header.Set("Content-Type", "application/json")
	jwb, err := p.klien.Do(req)
	if err != nil {
		return jawabanStorageTCO{}, galatJaringanTCO(err)
	}
	defer func() { _ = jwb.Body.Close() }()
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, batasJawabanStorageTCO))
		return jawabanStorageTCO{}, galatStatusTCO(jwb.StatusCode)
	}
	var j jawabanStorageTCO
	if err := json.NewDecoder(io.LimitReader(jwb.Body, batasJawabanStorageTCO)).Decode(&j); err != nil && !errors.Is(err, io.EOF) {
		return jawabanStorageTCO{}, fmt.Errorf("%w: JSON", ErrStorageJawabanRusakTCO)
	}
	return j, nil
}

func (p *pengirimBerkasHTTPTCO) dasar(ctx context.Context, token, kunci string) (permintaanStorageTCO, error) {
	if p.app == nil {
		return permintaanStorageTCO{}, fmt.Errorf("%w: App reader is not installed", outbox.ErrPenyimpananBelumDisetujui)
	}
	app, err := p.app(ctx)
	if err != nil {
		return permintaanStorageTCO{}, err
	}
	return permintaanStorageTCO{App: app, Kodestring: token, Namafile: kunci}, nil
}

// normalEkstensiTCO - `ext` unggahan dari NAMA BERKAS ASLI
// (`InsertGoogleStorage_Act` b587-b588 `@toLowerCase(Param.Ext)`): huruf
// kecil, tanpa titik. Tidak ditebak dari MIME - tabel MIME berbeda per OS.
func normalEkstensiTCO(ext string) string {
	return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(ext), "."))
}

// Kirim - unggah: isi dikirim base64 di medan `Image`; objek bernama `kunci`
// (IMAGEID) di `FolderStorageTCO` sehingga pengulangan menulis ke objek yang
// SAMA (AC 58).
func (p *pengirimBerkasHTTPTCO) Kirim(ctx context.Context, alamat, token, kunci string, isi io.Reader,
	tipe, ekstensi string) (models.ObjekPenyimpananTCO, error) {
	b, err := p.dasar(ctx, token, kunci)
	if err != nil {
		return models.ObjekPenyimpananTCO{}, err
	}
	data, err := io.ReadAll(io.LimitReader(isi, unggah.BatasUkuranUnggahan+1))
	if err != nil {
		return models.ObjekPenyimpananTCO{}, fmt.Errorf("services: reading queued file: %w", err)
	}
	if int64(len(data)) > unggah.BatasUkuranUnggahan {
		return models.ObjekPenyimpananTCO{}, fmt.Errorf("%w: file exceeds the limit", ErrStorageMenolakPermintaanTCO)
	}
	durasi := DurasiURLStorageTCO
	b.Durasi, b.Folder, b.MimeType, b.Ext = &durasi, FolderStorageTCO, tipe, normalEkstensiTCO(ekstensi)
	b.Image = base64.StdEncoding.EncodeToString(data)
	j, err := p.kirimJSON(ctx, alamat, b)
	if err != nil {
		return models.ObjekPenyimpananTCO{}, err
	}
	if strings.TrimSpace(j.URLImage) == "" {
		// `InsertGoogleStorage_Act.xml` b2902: `URLImage == ""` = gagal.
		return models.ObjekPenyimpananTCO{}, fmt.Errorf("%w: URLImage is empty", ErrStorageJawabanRusakTCO)
	}
	// tco4: yang Pega catat ke T_STORAGE_IMAGE (`Insert_T_Storage_SQL` b85):
	// URLImage, appfolder, exp, Namafile, App - apa adanya dari jawaban.
	// `exp` diubah seperti `InsertGoogleStorage_Act` b2366/b2431 (OQ-TCO-26).
	return models.ObjekPenyimpananTCO{ImageID: kunci, URLPublic: j.URLImage, AppFolder: j.AppFolder,
		Exp: models.ExpStorageTCO(j.Exp), Namafile: kunci, App: b.App}, nil
}

// urlBertanda meminta URL bertanda tangan satu objek (`geturl`). Objeknya =
// yang `GetUrlGoogleStorage_Act` b2125-b2295 salin ke `UpdateDoc` untuk
// `Update_T_Storage_SQL` (OQ-TCO-26).
func (p *pengirimBerkasHTTPTCO) urlBertanda(ctx context.Context, alamat, token, kunci string) (
	string, models.ObjekPenyimpananTCO, error) {

	b, err := p.dasar(ctx, token, kunci)
	if err != nil {
		return "", models.ObjekPenyimpananTCO{}, err
	}
	durasi := DurasiURLStorageTCO
	b.Durasi, b.Folder = &durasi, FolderStorageTCO
	j, err := p.kirimJSON(ctx, alamat, b)
	if err != nil {
		return "", models.ObjekPenyimpananTCO{}, err
	}
	// ⛔ HANYA https: isi lampiran tidak melintas jaringan tanpa sandi, dan
	// jawaban layanan tidak dapat menyuruh backend membuka alamat polos.
	u, err := url.Parse(strings.TrimSpace(j.URLImage))
	if err != nil || u.Host == "" || u.Scheme != "https" {
		return "", models.ObjekPenyimpananTCO{}, fmt.Errorf("%w: URLImage is malformed", ErrStorageJawabanRusakTCO)
	}
	return u.String(), models.ObjekPenyimpananTCO{ImageID: kunci, URLPublic: j.URLImage, AppFolder: j.AppFolder,
		Exp: models.ExpStorageTCO(j.Exp), TanggalUpload: models.TanggalUploadStorageTCO(j.DateTime)}, nil
}

// unduh membuka URL bertanda tangan yang DIBERIKAN layanan saat jalan.
//
// ⛔ Pengalihan TIDAK diikuti (3xx = galat): URL bertanda tangan menunjuk
// objeknya langsung, dan pengalihan dari jawaban luar adalah jalan membuka
// alamat internal. `jengkal` = header Range (cek keberadaan tanpa mengunduh).
func (p *pengirimBerkasHTTPTCO) unduh(ctx context.Context, bertanda, jengkal string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bertanda, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: URLImage is malformed", ErrStorageJawabanRusakTCO)
	}
	if jengkal != "" {
		req.Header.Set("Range", jengkal)
	}
	klien := *p.klien
	klien.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	jwb, err := klien.Do(req)
	if err != nil {
		return nil, galatJaringanTCO(err)
	}
	return jwb, nil
}

// Ambil - `geturl` lalu unduh isinya; objek geturl ikut dijawab.
func (p *pengirimBerkasHTTPTCO) Ambil(ctx context.Context, alamat, token, kunci string) (
	io.ReadCloser, models.ObjekPenyimpananTCO, error) {

	bertanda, objek, err := p.urlBertanda(ctx, alamat, token, kunci)
	if err != nil {
		return nil, models.ObjekPenyimpananTCO{}, err
	}
	jwb, err := p.unduh(ctx, bertanda, "")
	if err != nil {
		return nil, models.ObjekPenyimpananTCO{}, err
	}
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_ = jwb.Body.Close()
		return nil, models.ObjekPenyimpananTCO{}, galatStatusObjekTCO(jwb.StatusCode)
	}
	return jwb.Body, objek, nil
}

// Buang - `delete`: `Namafile` = jalur objek PENUH folder + nama
// (`DeleteGoogleStorage_Act` b1091: `.Folder` tersimpan dikurangi awalan
// skema `gs` + `.App`),
// tanpa `Folder`. Jawaban selain 2xx = galat - termasuk 404 (lihat galatStatusTCO).
func (p *pengirimBerkasHTTPTCO) Buang(ctx context.Context, alamat, token, kunci string) error {
	b, err := p.dasar(ctx, token, FolderStorageTCO+kunci)
	if err != nil {
		return err
	}
	_, err = p.kirimJSON(ctx, alamat, b)
	return err
}

// Periksa - `geturl` lalu membuka SATU byte objeknya: 404 = tidak ada. Objek
// geturl dijawab hanya bila objeknya ada.
func (p *pengirimBerkasHTTPTCO) Periksa(ctx context.Context, alamat, token, kunci string) (
	bool, models.ObjekPenyimpananTCO, error) {

	bertanda, objek, err := p.urlBertanda(ctx, alamat, token, kunci)
	if err != nil {
		return false, models.ObjekPenyimpananTCO{}, err
	}
	jwb, err := p.unduh(ctx, bertanda, "bytes=0-0")
	if err != nil {
		return false, models.ObjekPenyimpananTCO{}, err
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, 1<<10))
	_ = jwb.Body.Close()
	switch {
	case jwb.StatusCode == http.StatusNotFound:
		return false, models.ObjekPenyimpananTCO{}, nil
	case jwb.StatusCode >= 200 && jwb.StatusCode <= 299:
		return true, objek, nil
	}
	return false, models.ObjekPenyimpananTCO{}, galatStatusObjekTCO(jwb.StatusCode)
}
