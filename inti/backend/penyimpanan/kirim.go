package penyimpanan

// Connect-REST `ServiceGoogle` - SATU-SATUNYA berkas paket ini yang memegang klien HTTP keluar (izin bernama di
// `inti/backend/penjaga/lintasaplikasi_test.go`, keputusan work owner 08-10-2026).
//
// `ServiceGoogle.xml`: POST, badan JSON `UploadDoc.JSON` (b148), jawaban JSON ke `UploadDoc.Response` (b319-b325),
// batas waktu 300 detik (b29). Tanpa URL literal - alamatnya dari `GetLinkService`.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nusantarare/inti/backend/layanan"
)

// BatasWaktu - `ServiceGoogle.xml` b29 `pyResponseTimeout 300000`.
const BatasWaktu = 300 * time.Second

// batasJawaban - jawaban JSON layanan tidak dibaca tanpa batas.
const batasJawaban = 1 << 20

func klienBawaan() *http.Client { return &http.Client{Timeout: BatasWaktu} }

// permintaanBerkas - halaman `UploadDoc` / `ParamJSON` upload, geturl, dan delete sebagai JSON (nama medan VERBATIM;
// Durasi angka tanpa kutip).
type permintaanBerkas struct {
	App        string `json:"App"`
	Kodestring string `json:"Kodestring"`
	Durasi     *int   `json:"Durasi,omitempty"`
	Folder     string `json:"Folder,omitempty"`
	Namafile   string `json:"Namafile"`
	Image      string `json:"Image,omitempty"`
	Ext        string `json:"ext,omitempty"`
	MimeType   string `json:"MimeType,omitempty"`
}

// permintaanAI - `ParamJSON` `GeminiAIGoogle_Act` langkah 8-9 (`ListNamafile` diganti nama `Namafile`, Temp angka).
type permintaanAI struct {
	App        string   `json:"App"`
	Kodestring string   `json:"Kodestring"`
	Tanya      string   `json:"Tanya"`
	Temp       float64  `json:"Temp"`
	Namafile   []string `json:"Namafile,omitempty"`
}

// jawaban - `UploadDoc.Response`.
type jawaban struct {
	URLImage  string `json:"URLImage"`
	Exp       string `json:"exp"`
	AppFolder string `json:"appfolder"`
	DateTime  string `json:"DateTime"`
	Hasil     string `json:"Hasil"`
}

// belumSiap - APPNAME, token, atau alamat tidak tersedia. Sebabnya ikut terbungkus (galat bernama `inti/` menyebut
// kunci, tidak pernah nilai).
func belumSiap(apa string, sebab error) error {
	return fmt.Errorf("%w: %s is not available: %w", ErrStorageBelumSiap, apa, sebab)
}

// PesanBelumSiap - kalimat layar untuk `ErrStorageBelumSiap`: APA yang belum siap, menyebut nama kunci, tidak pernah
// nilainya. Sebab lain = kalimat umum (rinciannya di log server).
func PesanBelumSiap(err error) string {
	switch {
	case errors.Is(err, layanan.ErrGaramTokenKosong):
		return "File storage is not ready: there is no valid storage token, and STORAGE_TOKEN_SALT is not set on the server"
	case errors.Is(err, layanan.ErrAppNameKosong):
		return "File storage is not ready: APPNAME is missing in T_FOLDER_IMAGE"
	case errors.Is(err, layanan.ErrAlamatLayananTidakAda), errors.Is(err, layanan.ErrEndpointTidakDitemukan):
		return "File storage is not ready: the Google storage address is missing in M_LINK_SERVICE"
	}
	return "File storage is not ready; the details are in the server log"
}

// gagal - galat layanan TANPA sebab mentah: galat jaringan memuat alamat, inang, dan IP.
func gagal(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrStorageGagal, fmt.Sprintf(format, a...))
}

// panggil - satu Connect-REST: token (`GetTokenStorage_SQL`), alamat (`GetLinkService`), POST JSON, jawaban.
func (p *Penyimpanan) panggil(ctx context.Context, kunci layanan.KunciLayanan, app, pengguna string,
	susun func(kodestring string) any) (jawaban, error) {
	kode, err := p.token(ctx, app, pengguna)
	if err != nil {
		return jawaban{}, belumSiap("the storage token", err)
	}
	alamat, err := p.alamat(ctx, kunci)
	if err != nil {
		return jawaban{}, belumSiap(fmt.Sprintf("the M_LINK_SERVICE address (%s, %s)", kunci.Kategori1, kunci.Kategori2), err)
	}
	isi, err := json.Marshal(susun(kode))
	if err != nil {
		return jawaban{}, fmt.Errorf("penyimpanan: building the storage request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(isi))
	if err != nil {
		// ⛔ Alamat tidak disebut - pesan `url.Parse` memuatnya.
		return jawaban{}, fmt.Errorf("%w: the M_LINK_SERVICE address (%s, %s) is malformed", ErrStorageBelumSiap,
			kunci.Kategori1, kunci.Kategori2)
	}
	req.Header.Set("Content-Type", "application/json")
	jwb, err := tanpaPengalihan(p.klien).Do(req)
	if err != nil {
		return jawaban{}, galatJaringan(kunci.Kategori2, err)
	}
	defer func() { _ = jwb.Body.Close() }()
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, batasJawaban))
		return jawaban{}, gagal("%s answered status %d", kunci.Kategori2, jwb.StatusCode)
	}
	var j jawaban
	if err := json.NewDecoder(io.LimitReader(jwb.Body, batasJawaban)).Decode(&j); err != nil && !errors.Is(err, io.EOF) {
		return jawaban{}, gagal("the %s answer is not JSON", kunci.Kategori2)
	}
	return j, nil
}

// unduh membuka URL bertanda tangan. ⛔ Hanya https dan pengalihan TIDAK diikuti: jawaban luar tidak dapat menyuruh
// backend membuka alamat lain.
func (p *Penyimpanan) unduh(ctx context.Context, bertanda string) (io.ReadCloser, error) {
	u, err := url.Parse(strings.TrimSpace(bertanda))
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return nil, gagal("the signed URL is not an https address")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, gagal("the signed URL is malformed")
	}
	jwb, err := tanpaPengalihan(p.klien).Do(req)
	if err != nil {
		return nil, galatJaringan("the signed URL", err)
	}
	if jwb.StatusCode >= 200 && jwb.StatusCode <= 299 {
		return jwb.Body, nil
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, batasJawaban))
	_ = jwb.Body.Close()
	if jwb.StatusCode == http.StatusNotFound {
		return nil, ErrBerkasTidakDiStorage
	}
	return nil, gagal("the signed URL answered status %d", jwb.StatusCode)
}

// tanpaPengalihan - salinan klien yang TIDAK mengikuti pengalihan: badan permintaan memuat token, dan URL bertanda
// tangan menunjuk objeknya langsung; 3xx dibaca sebagai jawaban galat.
func tanpaPengalihan(k *http.Client) *http.Client {
	salin := *k
	salin.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &salin
}

// galatJaringan meringkas galat klien HTTP TANPA rinciannya (`url.Error` memuat alamat, dan URL bertanda tangan
// memuat kredensial).
func galatJaringan(titik string, err error) error {
	var ue *url.Error
	if errors.As(err, &ue) && ue.Timeout() {
		return gagal("%s timed out", titik)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return gagal("%s request was cancelled", titik)
	}
	return gagal("%s is unreachable", titik)
}
