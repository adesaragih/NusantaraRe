package services

// Penyimpanan berkas NYATA - keputusan work owner 03-10-2026 ("untuk document masih belum berfungsi, ikuti dari XML
// nya aja"; OQ-MPNL-10 dibalik), SELALU dipasang (keputusan work owner 03-10-2026 "selalu nyata, ikut XML" - tanpa
// saklar `PELAKSANA_STORAGE`). Stub lokal (`mpnl_penyimpanan.go`) tinggal untuk uji dan objek yang dicatatnya dulu.
//
//	unggah  `ProductNameSaveAttachment` 2.4 b904 → `InsertGoogleStorage_Act`: token 7 b1155 (`GetTokenStorage_SQL`),
//	        Set Data 8 b1339 (Durasi, Folder, Namafile, Image), ext 3 b566, MimeType 4 b714, SET JSON 9 b1572 (Durasi
//	        tanpa kutip, b1641), alamat 10 b1730 (`GetLinkService` "Google"/"upload"), Connect-REST 11 b1846
//	        (`ServiceGoogle`: POST JSON b148, 300 detik b29), jawaban → `InsertDoc` b2366–b2515 → `Insert_T_Storage_SQL`;
//	        `URLImage` kosong = gagal (b2902).
//	unduh   `DownloadAttProdName_Act` 6 b953 → `GetUrlGoogleStorage_Act` (ImageID, Durasi "1800" b999): `GetLinkStorage_SQL`
//	        b671; URL tersimpan dipakai selama EXPDATE belum lewat (b855–b1029 "JIKA EXPDATE SUDAH EXPIRED"); selain itu
//	        Folder = APPFOLDER tanpa Namafile (b1260) tanpa awalan gs+App (b1325), "Google"/"geturl", lalu
//	        `Update_T_Storage_SQL` b2375 (UpdateDoc b2125–b2295).
//	office  `View Office Online` b69291 → `DownloadAttProdName_Act` (ViewOffice): URL bertanda tangan langkah 6 b953;
//	        pembungkusan penampil langkah 7 b1103 di frontend (`penampilOffice.ts`).
//	hapus   `DeleteAttacProdName_act` 2 b411 → `DeleteGoogleStorage_Act`: `GetLinkStorage_SQL` b639, Namafile = APPFOLDER
//	        tanpa awalan gs+App (b1091), tanpa Folder, "Google"/"delete" (b1349); gagal = rekam tetap
//	        (`StepStatusFail` b444). `DeleteStorage_SQL` dijalankan layanan bersama hapus rekam.
//
// ⛔ Token: `GetTokenStorage_SQL` memanggil procedure `GET_TOKEN_STORAGE`; di sini ditiru `inti/backend/layanan`
// ([keputusan work owner] "jangan ada lagi pemanggilan procedure"; garam `STORAGE_TOKEN_SALT`).
// ⛔ Alamat dari `M_LINK_SERVICE` saat jalan (ADR-0013) - tidak ada di kode. Alamat, token, dan garam tidak pernah
// masuk pesan galat maupun log.
// ⚠️ PENYIMPANGAN SADAR: (1) pengiriman tetap efek keluar outbox (P5) - berkas ditahan di antrean lokal
// `UNGGAHAN_DIR` sampai objeknya tercatat; (2) isi diunduh backend dari URL bertanda tangan (hanya https, tanpa
// pengalihan) lalu diteruskan ke peramban - Pega membuka URL itu di jendela peramban; (3) jawaban geturl tanpa
// `appfolder` tidak mengosongkan APPFOLDER (hapus membutuhkannya); (4) objek yang dicatat stub lokal (URLPUBLIC
// kosong) dibaca dan dihapus di folder stub - riwayat sebelum penyambungan tetap terbaca.

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
	"os"
	"regexp"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// BatasWaktuStorage - `ServiceGoogle.xml` b29 `pyResponseTimeout 300000`.
const BatasWaktuStorage = 300 * time.Second

// sisaURLMinimum - URL tersimpan yang tinggal berlaku kurang dari ini diminta ulang: unduhan tidak putus di tengah.
const sisaURLMinimum = time.Minute

// MarginToken - token `GCP_IMAGE` (umur satu menit) yang tinggal kurang dari ini tidak dipakai ulang: unggahan base64
// besar tidak berangkat dengan token yang mati di tengah jalan (preseden Treaty Contract Out `MarginTokenTCO`).
const MarginToken = 15 * time.Second

// batasJawabanStorage - jawaban JSON layanan tidak dibaca tanpa batas.
const batasJawabanStorage = 1 << 20

// formatExp - bentuk `EXPDATE` yang ditulis (`To_date(…, 'DD/MM/YYYY HH24:MI:SS')`) dan dibaca kembali (`TO_CHAR`).
const formatExp = "02/01/2006 15:04:05"

var (
	// ErrStorageGagal - layanan penyimpanan tak terjangkau, menjawab galat, atau jawabannya rusak (502).
	ErrStorageGagal = errors.New("services: the storage service failed")
	// ErrStorageBelumSiap - alamat (`M_LINK_SERVICE`), App, atau token penyimpanan tidak tersedia (503).
	ErrStorageBelumSiap = errors.New("services: the storage service is not ready")
	// ErrBerkasTidakDiStorage - URL bertanda tangan menjawab 404: berkasnya tidak ada di penyimpanan (409).
	ErrBerkasTidakDiStorage = errors.New("services: the attachment file is not in storage")
)

// errURLDitolak - URL bertanda tangan ditolak penyimpanan (400 `ExpiredToken`, 401, 403): URL baru diminta sekali.
var errURLDitolak = errors.New("the signed URL was rejected")

// galatStorage - galat penyimpanan: kalimat layar tanpa sebab mentah (teks Oracle hanya di log), jenisnya terbaca
// `errors.Is`.
type galatStorage struct {
	jenis error
	layar string
	sebab error
}

func (g galatStorage) Error() string {
	if g.sebab == nil {
		return g.layar
	}
	return g.layar + ": " + g.sebab.Error()
}
func (g galatStorage) PesanLayar() string { return g.layar }
func (g galatStorage) Unwrap() []error {
	if g.sebab == nil {
		return []error{g.jenis}
	}
	return []error{g.jenis, g.sebab}
}

// belumSiap - kunci, App, atau token tidak tersedia. Galat bernama `inti/backend/layanan` (kalimatnya menyebut kunci,
// tidak pernah nilai) ikut tampil; sebab lain hanya di log.
func belumSiap(apa string, err error) error {
	layar := fmt.Sprintf("%s: %s is not available", ErrStorageBelumSiap, apa)
	for _, bernama := range []error{layanan.ErrGaramTokenKosong, layanan.ErrAppNameKosong, layanan.ErrEndpointTidakDitemukan,
		layanan.ErrAlamatLayananTidakAda, db.ErrTanpaOracle} {
		if errors.Is(err, bernama) {
			layar += " (" + bernama.Error() + ")"
			break
		}
	}
	return galatStorage{jenis: ErrStorageBelumSiap, layar: layar, sebab: err}
}

// gagal - galat layanan tanpa sebab mentah: galat jaringan memuat alamat, inang, dan IP.
func gagal(format string, a ...any) error {
	return galatStorage{jenis: ErrStorageGagal, layar: ErrStorageGagal.Error() + ": " + fmt.Sprintf(format, a...)}
}

// SumberAlamatStorage - alamat titik layanan dari `M_LINK_SERVICE` (`GetLinkService`).
type SumberAlamatStorage func(ctx context.Context, kunci layanan.KunciLayanan) (string, error)

// SumberTokenStorage - `Kodestring` untuk App (`GetTokenStorage_SQL`).
type SumberTokenStorage func(ctx context.Context, app string) (string, error)

// permintaanStorage - halaman `UploadDoc` / `ParamJSON` sebagai JSON (nama medan VERBATIM; Durasi angka, b1641).
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

// jawabanStorage - `UploadDoc.Response` (`ServiceGoogle.xml` b319–b325).
type jawabanStorage struct {
	URLImage  string `json:"URLImage"`
	Exp       string `json:"exp"`
	AppFolder string `json:"appfolder"`
	DateTime  string `json:"DateTime"`
}

type penyimpananGoogle struct {
	// lokal - antrean P5 dan objek yang dicatat stub; nil = `UNGGAHAN_DIR` kosong.
	lokal  *penyimpananLokal
	alamat SumberAlamatStorage
	token  SumberTokenStorage
	klien  *http.Client
	jam    func() time.Time
}

// PenyimpananGoogle menyusun penyimpanan nyata: antrean di folder lokal `dirUnggahan`, alamat dan token dari
// sumbernya, klien HTTP `klien` (nil = bawaan, batas waktu 300 detik), jam (nil = `time.Now`).
func PenyimpananGoogle(dirUnggahan string, alamat SumberAlamatStorage, token SumberTokenStorage, klien *http.Client,
	jam func() time.Time) PenyimpananBerkas {
	if klien == nil {
		klien = &http.Client{Timeout: BatasWaktuStorage}
	}
	if jam == nil {
		jam = time.Now
	}
	p := penyimpananGoogle{alamat: alamat, token: token, klien: klien, jam: jam}
	if l, ok := PenyimpananLokal(dirUnggahan).(penyimpananLokal); ok {
		p.lokal = &l
	}
	return p
}

func (p penyimpananGoogle) SimpanAntrean(ctx context.Context, imageID string, isi io.Reader) error {
	if p.lokal == nil {
		return ErrPenyimpananBelumDisetel
	}
	return p.lokal.SimpanAntrean(ctx, imageID, isi)
}

func (p penyimpananGoogle) BuangAntrean(ctx context.Context, imageID string) {
	if p.lokal != nil {
		p.lokal.BuangAntrean(ctx, imageID)
	}
}

// panggil - satu Connect-REST `ServiceGoogle`: token, alamat, POST JSON, jawaban.
func (p penyimpananGoogle) panggil(ctx context.Context, kunci layanan.KunciLayanan, badan permintaanStorage) (jawabanStorage, error) {
	token, err := p.token(ctx, badan.App)
	if err != nil {
		return jawabanStorage{}, belumSiap("the storage token", err)
	}
	alamat, err := p.alamat(ctx, kunci)
	if err != nil {
		return jawabanStorage{}, belumSiap(fmt.Sprintf("the M_LINK_SERVICE address (%s, %s)", kunci.Kategori1, kunci.Kategori2), err)
	}
	badan.Kodestring = token
	isi, err := json.Marshal(badan)
	if err != nil {
		return jawabanStorage{}, fmt.Errorf("services: building the storage request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(isi))
	if err != nil {
		// ⛔ Alamat tidak disebut - pesan `url.Parse` memuatnya.
		return jawabanStorage{}, belumSiap(fmt.Sprintf("a well-formed M_LINK_SERVICE address (%s, %s)", kunci.Kategori1, kunci.Kategori2), nil)
	}
	req.Header.Set("Content-Type", "application/json")
	jwb, err := tanpaPengalihan(p.klien).Do(req)
	if err != nil {
		return jawabanStorage{}, galatJaringan(kunci.Kategori2, err)
	}
	defer func() { _ = jwb.Body.Close() }()
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, batasJawabanStorage))
		return jawabanStorage{}, gagal("%s answered status %d", kunci.Kategori2, jwb.StatusCode)
	}
	var j jawabanStorage
	if err := json.NewDecoder(io.LimitReader(jwb.Body, batasJawabanStorage)).Decode(&j); err != nil && !errors.Is(err, io.EOF) {
		return jawabanStorage{}, gagal("the %s answer is not JSON", kunci.Kategori2)
	}
	return j, nil
}

// tanpaPengalihan - salinan klien yang TIDAK mengikuti pengalihan: badan permintaan memuat token, dan URL bertanda
// tangan menunjuk objeknya langsung; 3xx dibaca sebagai jawaban galat.
func tanpaPengalihan(k *http.Client) *http.Client {
	salin := *k
	salin.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &salin
}

// galatJaringan meringkas galat klien HTTP TANPA rinciannya (`url.Error` memuat alamat, inang, dan IP).
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

// Kirim - `InsertGoogleStorage_Act`: isi antrean dikirim base64 di `Image`; objek yang dijawab dicatat.
func (p penyimpananGoogle) Kirim(ctx context.Context, o models.ObjekPenyimpanan, ext, mime string) (models.ObjekPenyimpanan, error) {
	if p.lokal == nil {
		return models.ObjekPenyimpanan{}, ErrPenyimpananBelumDisetel
	}
	jalur, err := p.lokal.jalur("antre", o.ImageID)
	if err != nil {
		return models.ObjekPenyimpanan{}, err
	}
	f, err := os.Open(jalur)
	if errors.Is(err, os.ErrNotExist) {
		return models.ObjekPenyimpanan{}, ErrBerkasSumberHilang
	}
	if err != nil {
		return models.ObjekPenyimpanan{}, fmt.Errorf("services: reading the queued file: %w", err)
	}
	data, err := io.ReadAll(io.LimitReader(f, unggah.BatasUkuranUnggahan+1))
	_ = f.Close()
	if err != nil {
		return models.ObjekPenyimpanan{}, fmt.Errorf("services: reading the queued file: %w", err)
	}
	durasi := o.DurasiDetik
	j, err := p.panggil(ctx, layanan.KunciUnggahBerkas, permintaanStorage{App: o.AppName, Durasi: &durasi, Folder: o.AppFolder,
		Namafile: o.FileName, Image: base64.StdEncoding.EncodeToString(data), Ext: ext, MimeType: mime})
	if err != nil {
		return models.ObjekPenyimpanan{}, err
	}
	if strings.TrimSpace(j.URLImage) == "" {
		return models.ObjekPenyimpanan{}, gagal("upload answered an empty URLImage")
	}
	folder := o.AppFolder
	o.URLPublic, o.AppFolder, o.Exp = j.URLImage, j.AppFolder, expPega(j.Exp)
	if strings.TrimSpace(o.AppFolder) == "" {
		// Jawaban tanpa `appfolder`: dirakit seperti bentuknya di DEV (awalan gs+App, Folder, Namafile) - hapus dan
		// geturl membutuhkan jalur objek.
		o.AppFolder = awalanGS(o.AppName) + folder + o.FileName
	}
	return o, nil
}

// awalanGS - awalan objek di APPFOLDER: skema gs, inang App, garis miring (b1325 / b1091).
func awalanGS(app string) string { return (&url.URL{Scheme: "gs", Host: app, Path: "/"}).String() }

// jalurObjek - `@replaceAll(.Folder, <awalan gs+App>, "")`: jalur objek penuh (Namafile hapus, b1091).
func jalurObjek(o models.ObjekPenyimpanan) string {
	return strings.ReplaceAll(o.AppFolder, awalanGS(o.AppName), "")
}

// berlaku - URL tersimpan masih dapat dipakai (EXPDATE belum lewat, dengan sisa minimum).
//
// ⚠️ EXPDATE dibaca jam Asia/Jakarta, BUKAN GMT seperti rumus Pega (b2366 `+ " GMT"`): bukti DEV 03-10-2026 - URL
// ber-EXPDATE `13:28:00` ditolak Google `ExpiredToken` sesudah 13:28 WIB (layanan menjawab `exp` jam lokal). Dibaca GMT,
// URL mati dianggap berlaku tujuh jam lagi. Bila ternyata GMT, akibatnya hanya geturl lebih awal.
func (p penyimpananGoogle) berlaku(o models.ObjekPenyimpanan) bool {
	if strings.TrimSpace(o.URLPublic) == "" {
		return false
	}
	exp, err := time.ParseInLocation(formatExp, strings.TrimSpace(o.Exp), zonaJakarta)
	return err == nil && exp.After(p.jam().Add(sisaURLMinimum))
}

// Tautan - `GetUrlGoogleStorage_Act`: URL bertanda tangan objek terkirim; objek baru dijawab bila geturl dipanggil.
// Objek yang dicatat stub lokal tidak punya URL (`ErrOfficeStub`).
func (p penyimpananGoogle) Tautan(ctx context.Context, o models.ObjekPenyimpanan) (string, *models.ObjekPenyimpanan, error) {
	if strings.TrimSpace(o.URLPublic) == "" {
		return "", nil, ErrOfficeStub
	}
	if p.berlaku(o) {
		return o.URLPublic, nil, nil
	}
	return p.mintaURL(ctx, o)
}

// mintaURL - geturl (`GetUrlGoogleStorage_Act` b1238–b1781): URL bertanda tangan BARU dan objek yang diperbarui.
func (p penyimpananGoogle) mintaURL(ctx context.Context, o models.ObjekPenyimpanan) (string, *models.ObjekPenyimpanan, error) {
	durasi := DurasiLampiran
	folder := strings.TrimSuffix(jalurObjek(o), o.FileName)
	j, err := p.panggil(ctx, layanan.KunciURLBerkas, permintaanStorage{App: o.AppName, Durasi: &durasi, Folder: folder,
		Namafile: o.FileName})
	if err != nil {
		return "", nil, err
	}
	if strings.TrimSpace(j.URLImage) == "" {
		return "", nil, gagal("geturl answered an empty URLImage")
	}
	segar := o
	segar.URLPublic, segar.Exp, segar.TanggalUpload = j.URLImage, expPega(j.Exp), tanggalUploadPega(j.DateTime)
	if strings.TrimSpace(j.AppFolder) != "" {
		segar.AppFolder = j.AppFolder
	}
	return j.URLImage, &segar, nil
}

// Buka - `Tautan` lalu isi dari URL bertanda tangan; objek baru dijawab bila geturl dipanggil.
func (p penyimpananGoogle) Buka(ctx context.Context, o models.ObjekPenyimpanan) (io.ReadCloser, *models.ObjekPenyimpanan, error) {
	if strings.TrimSpace(o.URLPublic) == "" {
		// Dicatat stub lokal - berkasnya di folder stub, tidak pernah di penyimpanan.
		if p.lokal == nil {
			return nil, nil, ErrPenyimpananBelumDisetel
		}
		return p.lokal.Buka(ctx, o)
	}
	bertanda, baru, err := p.Tautan(ctx, o)
	if err != nil {
		return nil, nil, err
	}
	isi, err := p.unduh(ctx, bertanda)
	if err != nil && baru == nil && errors.Is(err, errURLDitolak) {
		// URL TERSIMPAN ditolak walau EXPDATE belum lewat (selisih jam / zona): URL baru, diulang SEKALI.
		if bertanda, baru, err = p.mintaURL(ctx, o); err != nil {
			return nil, nil, err
		}
		isi, err = p.unduh(ctx, bertanda)
	}
	if err != nil {
		return nil, nil, err
	}
	return isi, baru, nil
}

// unduh membuka URL bertanda tangan. ⛔ Hanya https (isi lampiran tidak melintas tanpa sandi; jawaban luar tidak
// dapat menyuruh backend membuka alamat polos) dan pengalihan TIDAK diikuti.
func (p penyimpananGoogle) unduh(ctx context.Context, bertanda string) (io.ReadCloser, error) {
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
	_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, batasJawabanStorage))
	_ = jwb.Body.Close()
	switch jwb.StatusCode {
	case http.StatusNotFound:
		return nil, galatStorage{jenis: ErrBerkasTidakDiStorage, layar: ErrBerkasTidakDiStorage.Error()}
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		return nil, galatStorage{jenis: ErrStorageGagal, sebab: errURLDitolak,
			layar: fmt.Sprintf("%s: the signed URL answered status %d", ErrStorageGagal, jwb.StatusCode)}
	}
	return nil, gagal("the signed URL answered status %d", jwb.StatusCode)
}

// Hapus - `DeleteGoogleStorage_Act` untuk objek terkirim; antrean dan berkas stub lokal ikut dibuang.
func (p penyimpananGoogle) Hapus(ctx context.Context, imageID string, o *models.ObjekPenyimpanan) error {
	if o != nil && strings.TrimSpace(o.URLPublic) != "" {
		jalur := jalurObjek(*o)
		if strings.TrimSpace(jalur) == "" {
			return gagal("APPFOLDER is empty; the storage object cannot be named")
		}
		if _, err := p.panggil(ctx, layanan.KunciHapusBerkas, permintaanStorage{App: o.AppName, Namafile: jalur}); err != nil {
			return err
		}
	}
	if p.lokal == nil {
		return nil
	}
	return p.lokal.Hapus(ctx, imageID, nil)
}

// polaExp - `YYYYMMDDTHHMMSS` di awal `exp` sesudah `-`/`:` dibuang.
var polaExp = regexp.MustCompile(`^\d{8}T\d{6}`)

// expPega - `exp` jawaban → `To_date(exp, 'DD/MM/YYYY HH24:MI:SS')`: `@substring(@pxReplaceAllViaRegex(exp,
// "[-:]", ""), 0, 19) + " GMT"` lalu `@FormatDateTime(…, "dd/MM/yyyy HH:mm:ss", "GMT")` (b2366, b2431). Bentuk yang
// tidak terbaca → kosong (NULL), bukan galat To_date yang menggagalkan pencatatan.
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

// tanggalUploadPega - `DateTime` geturl apa adanya bila berbentuk `MM/DD/YYYY HH24:MI:SS` (`Update_T_Storage_SQL`);
// bentuk lain → kosong (NULL), bukan ORA-01843.
func tanggalUploadPega(teks string) string {
	t := strings.TrimSpace(teks)
	if _, err := time.Parse("01/02/2006 15:04:05", t); err != nil {
		return ""
	}
	return t
}

// penyimpanTokenStorage - bacaan dan tulisan `GCP_IMAGE` (`inti/backend/layanan.PenyimpanToken`).
type penyimpanTokenStorage interface {
	TokenBerlaku(ctx context.Context, tx *db.Tx, appName string, saat time.Time) (string, error)
	SimpanToken(ctx context.Context, tx *db.Tx, appName, token, pengguna string, sampai time.Time) error
}

// tokenStorage - `GET_TOKEN_STORAGE` ditiru (`inti/backend/layanan/token.go`): token yang masih berlaku LEBIH dari
// `MarginToken` dipakai ulang - tanpa garam; selain itu token baru dirakit dengan garam (`STORAGE_TOKEN_SALT`) dan
// disimpan dengan umur satu menit, pengguna `Job` (`NVL(masukan, 'Job')`). ⛔ Garam dan token tidak pernah masuk
// pesan galat.
func tokenStorage(ctx context.Context, tx *db.Tx, p penyimpanTokenStorage, garam, app string, saat time.Time) (string, error) {
	if strings.TrimSpace(app) == "" {
		return "", layanan.ErrAppNameKosong
	}
	lama, err := p.TokenBerlaku(ctx, tx, app, saat.Add(MarginToken))
	if err != nil || lama != "" {
		return lama, err
	}
	if strings.TrimSpace(garam) == "" {
		return "", layanan.ErrGaramTokenKosong
	}
	baru, err := layanan.RakitToken(garam, saat)
	if err != nil {
		return "", err
	}
	if err := p.SimpanToken(ctx, tx, app, baru, layanan.PenggunaTokenBawaan, saat.Add(layanan.UmurToken)); err != nil {
		return "", fmt.Errorf("services: saving the storage token: %w", err)
	}
	return baru, nil
}
