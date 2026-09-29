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
//	`App`, `Kodestring`, `Folder`, `Namafile`, `Durasi`) dan
//	`DeleteGoogleStorage_Act.xml` (b1019-b1091: `App`, `Kodestring`, `Namafile`).
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
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BatasWaktuStorageTCO - `ServiceGoogle.xml` b29 `pyResponseTimeout 300000`.
const BatasWaktuStorageTCO = 300 * time.Second

// FolderStorageTCO - folder objek lampiran modul ini di layanan penyimpanan.
// `[keputusan kami]` (OQ-TCO-22): Pega Claim Life menyusun `Param.Folder +
// "/Doc/" + tahun/bulan`; korpus Treaty Contract Out tidak memuat unggahnya.
const FolderStorageTCO = "TreatyContractOut"

// DurasiURLStorageTCO - `Durasi` URL bertanda tangan (satuan menurut layanan).
// `[keputusan kami]` (OQ-TCO-22): nilainya dari pemanggil di Pega.
const DurasiURLStorageTCO = 60

// batasJawabanStorageTCO - jawaban JSON layanan tidak dibaca tanpa batas.
const batasJawabanStorageTCO = 1 << 20

var (
	// ErrStorageTakTerjangkauTCO - jaringan/batas waktu; layak dicoba ulang.
	ErrStorageTakTerjangkauTCO = errors.New("services: layanan penyimpanan tidak terjangkau")
	// ErrStorageGagalTCO - layanan menjawab galat (5xx, 401, 403, 429); layak dicoba ulang.
	ErrStorageGagalTCO = errors.New("services: layanan penyimpanan menjawab galat")
	// ErrStorageMenolakPermintaanTCO - layanan menolak bentuk permintaan (400, 422);
	// PERMANEN sampai manusia bertindak.
	ErrStorageMenolakPermintaanTCO = errors.New("services: layanan penyimpanan menolak permintaan")
	// ErrStorageJawabanRusakTCO - jawaban tidak berbentuk (`URLImage` kosong, JSON rusak).
	ErrStorageJawabanRusakTCO = errors.New("services: jawaban layanan penyimpanan tidak berbentuk")
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
		return fmt.Errorf("%w: batas waktu", ErrStorageTakTerjangkauTCO)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: permintaan dibatalkan", ErrStorageTakTerjangkauTCO)
	}
	return ErrStorageTakTerjangkauTCO
}

func galatStatusTCO(kode int) error {
	switch {
	case kode == http.StatusBadRequest || kode == http.StatusUnprocessableEntity:
		return fmt.Errorf("%w: status %d", ErrStorageMenolakPermintaanTCO, kode)
	case kode == http.StatusNotFound:
		return ErrBerkasTidakAdaDiPenyimpanan
	}
	return fmt.Errorf("%w: status %d", ErrStorageGagalTCO, kode)
}

// kirimJSON mengirim satu permintaan JSON dan membaca jawabannya.
func (p *pengirimBerkasHTTPTCO) kirimJSON(ctx context.Context, alamat string, badan permintaanStorageTCO) (jawabanStorageTCO, error) {
	isi, err := json.Marshal(badan)
	if err != nil {
		return jawabanStorageTCO{}, fmt.Errorf("services: merakit permintaan penyimpanan: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(isi))
	if err != nil {
		// ⛔ Alamat tidak disebut - pesan `url.Parse` memuatnya.
		return jawabanStorageTCO{}, fmt.Errorf("%w: alamat layanan tidak berbentuk", ErrStorageMenolakPermintaanTCO)
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
		return permintaanStorageTCO{}, fmt.Errorf("%w: pembaca App belum dipasang", ErrPenyimpananBelumDisetujui)
	}
	app, err := p.app(ctx)
	if err != nil {
		return permintaanStorageTCO{}, err
	}
	return permintaanStorageTCO{App: app, Kodestring: token, Namafile: kunci}, nil
}

// ekstensiDariMime - `ext` unggahan (`InsertGoogleStorage_Act` b587-b588:
// huruf kecil, tanpa titik).
func ekstensiDariMime(tipe string) string {
	ext, _ := mime.ExtensionsByType(tipe)
	if len(ext) == 0 {
		return ""
	}
	return strings.ToLower(strings.TrimPrefix(ext[0], "."))
}

// Kirim - unggah: isi dikirim base64 di medan `Image`; objek bernama `kunci`
// (IMAGEID) sehingga pengulangan menulis ke objek yang SAMA (AC 58).
func (p *pengirimBerkasHTTPTCO) Kirim(ctx context.Context, alamat, token, kunci string, isi io.Reader, tipe string) error {
	b, err := p.dasar(ctx, token, kunci)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(io.LimitReader(isi, BatasUkuranUnggahan+1))
	if err != nil {
		return fmt.Errorf("services: membaca berkas antrean: %w", err)
	}
	if int64(len(data)) > BatasUkuranUnggahan {
		return fmt.Errorf("%w: berkas melebihi batas", ErrStorageMenolakPermintaanTCO)
	}
	durasi := DurasiURLStorageTCO
	b.Durasi, b.Folder, b.MimeType, b.Ext = &durasi, FolderStorageTCO, tipe, ekstensiDariMime(tipe)
	b.Image = base64.StdEncoding.EncodeToString(data)
	j, err := p.kirimJSON(ctx, alamat, b)
	if err != nil {
		return err
	}
	if strings.TrimSpace(j.URLImage) == "" {
		// `InsertGoogleStorage_Act.xml` b2902: `URLImage == ""` = gagal.
		return fmt.Errorf("%w: URLImage kosong", ErrStorageJawabanRusakTCO)
	}
	return nil
}

// urlBertanda meminta URL bertanda tangan satu objek (`geturl`).
func (p *pengirimBerkasHTTPTCO) urlBertanda(ctx context.Context, alamat, token, kunci string) (string, error) {
	b, err := p.dasar(ctx, token, kunci)
	if err != nil {
		return "", err
	}
	durasi := DurasiURLStorageTCO
	b.Durasi, b.Folder = &durasi, FolderStorageTCO
	j, err := p.kirimJSON(ctx, alamat, b)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(strings.TrimSpace(j.URLImage))
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return "", fmt.Errorf("%w: URLImage tidak berbentuk", ErrStorageJawabanRusakTCO)
	}
	return u.String(), nil
}

// unduh membuka URL bertanda tangan yang DIBERIKAN layanan saat jalan.
func (p *pengirimBerkasHTTPTCO) unduh(ctx context.Context, bertanda string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bertanda, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: URLImage tidak berbentuk", ErrStorageJawabanRusakTCO)
	}
	jwb, err := p.klien.Do(req)
	if err != nil {
		return nil, galatJaringanTCO(err)
	}
	return jwb, nil
}

// Ambil - `geturl` lalu unduh isinya.
func (p *pengirimBerkasHTTPTCO) Ambil(ctx context.Context, alamat, token, kunci string) (io.ReadCloser, error) {
	bertanda, err := p.urlBertanda(ctx, alamat, token, kunci)
	if err != nil {
		return nil, err
	}
	jwb, err := p.unduh(ctx, bertanda)
	if err != nil {
		return nil, err
	}
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_ = jwb.Body.Close()
		return nil, galatStatusTCO(jwb.StatusCode)
	}
	return jwb.Body, nil
}

// Buang - `delete` (`App`, `Kodestring`, `Namafile`).
func (p *pengirimBerkasHTTPTCO) Buang(ctx context.Context, alamat, token, kunci string) error {
	b, err := p.dasar(ctx, token, kunci)
	if err != nil {
		return err
	}
	_, err = p.kirimJSON(ctx, alamat, b)
	return err
}

// Periksa - `geturl` lalu memastikan objeknya dapat dibuka.
func (p *pengirimBerkasHTTPTCO) Periksa(ctx context.Context, alamat, token, kunci string) (bool, error) {
	bertanda, err := p.urlBertanda(ctx, alamat, token, kunci)
	if errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	jwb, err := p.unduh(ctx, bertanda)
	if err != nil {
		return false, err
	}
	_ = jwb.Body.Close()
	switch {
	case jwb.StatusCode == http.StatusNotFound:
		return false, nil
	case jwb.StatusCode >= 200 && jwb.StatusCode <= 299:
		return true, nil
	}
	return false, galatStatusTCO(jwb.StatusCode)
}
