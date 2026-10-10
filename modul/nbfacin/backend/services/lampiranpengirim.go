package services

// Transport penyimpanan lampiran NB FacIn (tiket 67) - SATU-SATUNYA berkas modul ini yang memegang klien HTTP keluar
// (penjaga `TestNolAlamatLayananDiKode`, `inti/backend/penjaga/lintasaplikasi_test.go` `berkasKlienHTTPDisetujui`).
// Salinan `tco_pengirim_storage.go` Treaty Contract Out (keputusan work owner 08-10-2026 K-B "Salin ke NB FacIn") dengan
// bentuk Pega NB (lihat lampiranstorage.go). Aktif HANYA bila PELAKSANA_STORAGE=nyata; uji memakai server httptest.
//
// ⛔ Alamat dari M_LINK_SERVICE saat jalan (ADR-0013); alamat, token, garam tidak pernah masuk pesan galat / log.

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/nbfacin/backend/models"
)

// ---- rangkaian jarak jauh ------------------------------------------------------------------------------------------

// PencatatObjekStorage - Update_T_Storage_SQL sesudah geturl, dalam transaksi pendeknya sendiri (Pega: RDB-List lepas).
type PencatatObjekStorage interface {
	PerbaruiObjek(ctx context.Context, o models.ObjekStorage) error
}

type permintaanStorage struct {
	App        string `json:"App"`
	Kodestring string `json:"Kodestring"`
	Durasi     *int   `json:"Durasi,omitempty"`
	Folder     string `json:"Folder,omitempty"`
	Namafile   string `json:"Namafile"`
	Image      string `json:"Image,omitempty"`
	Ext        string `json:"ext,omitempty"`
	MimeType   string `json:"MimeType,omitempty"`
}

type jawabanStorage struct {
	URLImage  string `json:"URLImage"`
	Exp       string `json:"exp"`
	AppFolder string `json:"appfolder"`
	DateTime  string `json:"DateTime"`
}

// storageJarakJauh - resolver M_LINK_SERVICE + cache token + HTTP.
type storageJarakJauh struct {
	resolver layanan.ResolverEndpoint
	token    *CacheTokenStorage
	app      func(ctx context.Context) (string, error)
	klien    *http.Client
	pencatat PencatatObjekStorage // nil = tanpa penyegaran
	jam      func() time.Time
}

// NewStorageLampiranJarakJauh menyusun pelaksana nyata; klien nil = klien bawaan berbatas waktu 300 s.
func NewStorageLampiranJarakJauh(resolver layanan.ResolverEndpoint, token *CacheTokenStorage,
	app func(ctx context.Context) (string, error), klien *http.Client, pencatat PencatatObjekStorage,
	jam func() time.Time) KlienStorageLampiran {
	if klien == nil {
		klien = &http.Client{Timeout: BatasWaktuStorageLampiran}
	}
	if jam == nil {
		jam = time.Now
	}
	return &storageJarakJauh{resolver: resolver, token: token, app: app, klien: klien, pencatat: pencatat, jam: jam}
}

// siapkan meresolve alamat SAAT JALAN, lalu App dan token. Alamat tidak pernah disimpan di struct.
func (p *storageJarakJauh) siapkan(ctx context.Context, kunci layanan.KunciLayanan, lokasi string) (string, permintaanStorage, error) {
	if p.resolver == nil {
		return "", permintaanStorage{}, layanan.ErrResolverBelumDiputuskan
	}
	alamat, err := p.resolver.Resolve(ctx, kunci)
	if err != nil {
		return "", permintaanStorage{}, err
	}
	if p.app == nil {
		return "", permintaanStorage{}, fmt.Errorf("%w: pembaca App belum dipasang", outbox.ErrPenyimpananBelumDisetujui)
	}
	app, err := p.app(ctx)
	if err != nil {
		return "", permintaanStorage{}, err
	}
	if p.token == nil {
		return "", permintaanStorage{}, fmt.Errorf("%w: cache token belum dipasang", ErrTokenPenyimpananGagal)
	}
	tok, err := p.token.Token(ctx)
	if err != nil {
		return "", permintaanStorage{}, err
	}
	return alamat, permintaanStorage{App: app, Kodestring: tok, Namafile: lokasi}, nil
}

// galatJarakJauh - menyebut KUNCI layanan, tidak pernah alamatnya; token yang ditolak dilupakan.
func (p *storageJarakJauh) galatJarakJauh(kunci layanan.KunciLayanan, err error) error {
	if errors.Is(err, ErrStorageTokenDitolak) && p.token != nil {
		p.token.Lupakan()
	}
	return fmt.Errorf("services: penyimpanan jarak jauh (%s/%s): %w", kunci.Kategori1, kunci.Kategori2, err)
}

// galatJaringan meringkas galat klien HTTP TANPA rinciannya (alamat / inang / IP ada di sana).
func galatJaringan(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) && ue.Timeout() {
		return fmt.Errorf("%w: batas waktu", ErrStorageTakTerjangkau)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: permintaan dibatalkan", ErrStorageTakTerjangkau)
	}
	return ErrStorageTakTerjangkau
}

// galatStatus - jawaban TITIK LAYANAN. ⛔ 404 di sini BUKAN "berkas tidak ada" (juga jawaban M_LINK_SERVICE yang salah).
func galatStatus(kode int) error {
	switch kode {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return fmt.Errorf("%w: status %d", ErrStorageMenolakPermintaan, kode)
	case http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("%w: status %d", ErrStorageTokenDitolak, kode)
	}
	return fmt.Errorf("%w: status %d", ErrStorageGagal, kode)
}

// galatStatusObjek - jawaban URL bertanda tangan: 404 = objek tidak ada; 401/403 = tanda tangan ditolak (bukan token kita).
func galatStatusObjek(kode int) error {
	switch kode {
	case http.StatusNotFound:
		return ErrBerkasTidakAdaDiPenyimpanan
	case http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("%w: status %d", ErrStorageGagal, kode)
	}
	return galatStatus(kode)
}

func (p *storageJarakJauh) kirimJSON(ctx context.Context, alamat string, badan permintaanStorage) (jawabanStorage, error) {
	isi, err := json.Marshal(badan)
	if err != nil {
		return jawabanStorage{}, fmt.Errorf("services: menyusun permintaan penyimpanan: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(isi))
	if err != nil {
		return jawabanStorage{}, fmt.Errorf("%w: alamat layanan rusak", ErrStorageMenolakPermintaan)
	}
	req.Header.Set("Content-Type", "application/json")
	jwb, err := p.klien.Do(req)
	if err != nil {
		return jawabanStorage{}, galatJaringan(err)
	}
	defer func() { _ = jwb.Body.Close() }()
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, batasJawabanStorage))
		return jawabanStorage{}, galatStatus(jwb.StatusCode)
	}
	var j jawabanStorage
	if err := json.NewDecoder(io.LimitReader(jwb.Body, batasJawabanStorage)).Decode(&j); err != nil && !errors.Is(err, io.EOF) {
		return jawabanStorage{}, fmt.Errorf("%w: JSON", ErrStorageJawabanRusak)
	}
	return j, nil
}

// Simpan - InsertGoogleStorage_Act langkah 8-12.
func (p *storageJarakJauh) Simpan(ctx context.Context, lok models.LokasiObjek, isi io.Reader, ekstensi, mime string) (
	models.ObjekStorage, error) {
	alamat, b, err := p.siapkan(ctx, layanan.KunciUnggahBerkas, lok.Namafile)
	if err != nil {
		return models.ObjekStorage{}, err
	}
	data, err := io.ReadAll(io.LimitReader(isi, unggah.BatasUkuranUnggahan+1))
	if err != nil {
		return models.ObjekStorage{}, fmt.Errorf("services: membaca berkas antrean: %w", err)
	}
	if int64(len(data)) > unggah.BatasUkuranUnggahan {
		return models.ObjekStorage{}, fmt.Errorf("%w: berkas melebihi batas", ErrStorageMenolakPermintaan)
	}
	durasi := DurasiURLStorageLampiran
	b.Durasi, b.Folder, b.Ext, b.MimeType = &durasi, lok.Folder, ekstensi, mime
	b.Image = base64.StdEncoding.EncodeToString(data)
	j, err := p.kirimJSON(ctx, alamat, b)
	if err != nil {
		return models.ObjekStorage{}, p.galatJarakJauh(layanan.KunciUnggahBerkas, err)
	}
	if strings.TrimSpace(j.URLImage) == "" {
		// Langkah 12: `UploadDoc.Response.URLImage == ""` = tidak dicatat.
		return models.ObjekStorage{}, p.galatJarakJauh(layanan.KunciUnggahBerkas,
			fmt.Errorf("%w: URLImage kosong", ErrStorageJawabanRusak))
	}
	return models.ObjekStorage{ImageID: lok.ImageID, URLPublic: j.URLImage, AppFolder: j.AppFolder,
		Exp: models.ExpStorage(j.Exp), Namafile: lok.Namafile, App: b.App}, nil
}

// urlBertanda - geturl satu lokasi (GetUrlGoogleStorage_Act 6.2-6.5); hanya https.
func (p *storageJarakJauh) urlBertanda(ctx context.Context, lok models.LokasiObjek) (string, models.ObjekStorage, error) {
	alamat, b, err := p.siapkan(ctx, layanan.KunciURLBerkas, lok.Namafile)
	if err != nil {
		return "", models.ObjekStorage{}, err
	}
	durasi := DurasiURLStorageLampiran
	b.Durasi, b.Folder = &durasi, lok.Folder
	j, err := p.kirimJSON(ctx, alamat, b)
	if err != nil {
		return "", models.ObjekStorage{}, p.galatJarakJauh(layanan.KunciURLBerkas, err)
	}
	u, err := url.Parse(strings.TrimSpace(j.URLImage))
	if err != nil || u.Host == "" || u.Scheme != "https" {
		return "", models.ObjekStorage{}, p.galatJarakJauh(layanan.KunciURLBerkas,
			fmt.Errorf("%w: URLImage tidak sah", ErrStorageJawabanRusak))
	}
	return u.String(), models.ObjekStorage{ImageID: lok.ImageID, URLPublic: j.URLImage, AppFolder: j.AppFolder,
		Exp: models.ExpStorage(j.Exp), Namafile: lok.Namafile, App: b.App,
		TanggalUpload: models.TanggalUploadStorage(j.DateTime)}, nil
}

// unduh membuka URL bertanda tangan. ⛔ Pengalihan TIDAK diikuti; `jengkal` = header Range.
func (p *storageJarakJauh) unduh(ctx context.Context, bertanda, jengkal string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bertanda, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: URL objek tidak sah", ErrStorageJawabanRusak)
	}
	if jengkal != "" {
		req.Header.Set("Range", jengkal)
	}
	klien := *p.klien
	klien.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	jwb, err := klien.Do(req)
	if err != nil {
		return nil, galatJaringan(err)
	}
	return jwb, nil
}

// segarkan - Update_T_Storage_SQL; gagal TIDAK menggagalkan unduhan (dicatat tanpa URL).
func (p *storageJarakJauh) segarkan(ctx context.Context, o models.ObjekStorage) {
	if p.pencatat == nil || o.ImageID == "" {
		return
	}
	if err := p.pencatat.PerbaruiObjek(ctx, o); err != nil {
		log.Printf("nbfacin lampiran: menyegarkan catatan objek %s gagal: %v", o.ImageID, err)
	}
}

// Buka - GetUrlGoogleStorage_Act: EXPDATE kosong = tanpa URL (langkah 4 Exit); URL tersimpan dipakai selama berlaku;
// selain itu geturl + Update_T_Storage_SQL. Isinya diunduh backend (proksi).
func (p *storageJarakJauh) Buka(ctx context.Context, link models.LinkStorage) (io.ReadCloser, error) {
	if strings.TrimSpace(link.ExpUTC) == "" {
		return nil, ErrBerkasTidakAdaDiPenyimpanan
	}
	bertanda := ""
	if models.LinkBerlaku(link, p.jam(), sisaURLTersimpan) {
		u, err := url.Parse(strings.TrimSpace(link.URLPublic))
		if err == nil && u.Host != "" && u.Scheme == "https" {
			bertanda = u.String()
		}
	}
	if bertanda == "" {
		b, objek, err := p.urlBertanda(ctx, models.LokasiDariLink(link))
		if err != nil {
			return nil, err
		}
		p.segarkan(ctx, objek)
		bertanda = b
	}
	jwb, err := p.unduh(ctx, bertanda, "")
	if err != nil {
		return nil, err
	}
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_ = jwb.Body.Close()
		return nil, galatStatusObjek(jwb.StatusCode)
	}
	return jwb.Body, nil
}

// Hapus - DeleteGoogleStorage_Act langkah 6-9 (tanpa cek keberadaan, seperti Pega); jawaban selain 2xx = galat.
func (p *storageJarakJauh) Hapus(ctx context.Context, link models.LinkStorage) error {
	jalur := models.JalurHapusDariLink(link)
	if strings.TrimSpace(jalur) == "" {
		return ErrBerkasTidakAdaDiPenyimpanan
	}
	alamat, b, err := p.siapkan(ctx, layanan.KunciHapusBerkas, jalur)
	if err != nil {
		return err
	}
	if _, err := p.kirimJSON(ctx, alamat, b); err != nil {
		return p.galatJarakJauh(layanan.KunciHapusBerkas, err)
	}
	return nil
}

// Ada - geturl lalu membuka SATU byte objeknya: 404 = tidak ada.
func (p *storageJarakJauh) Ada(ctx context.Context, lok models.LokasiObjek) (bool, models.ObjekStorage, error) {
	bertanda, objek, err := p.urlBertanda(ctx, lok)
	if err != nil {
		return false, models.ObjekStorage{}, err
	}
	jwb, err := p.unduh(ctx, bertanda, "bytes=0-0")
	if err != nil {
		return false, models.ObjekStorage{}, err
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, 1<<10))
	_ = jwb.Body.Close()
	switch {
	case jwb.StatusCode == http.StatusNotFound:
		return false, models.ObjekStorage{}, nil
	case jwb.StatusCode >= 200 && jwb.StatusCode <= 299:
		return true, objek, nil
	}
	return false, models.ObjekStorage{}, galatStatusObjek(jwb.StatusCode)
}
