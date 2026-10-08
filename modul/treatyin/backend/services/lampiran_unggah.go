package services

// Tombol `Upload file` panel Attachment — keputusan pemakai 8 Oktober 2026:
// *"untuk upload file seharusnya kesini SELECT * FROM M_ATTACHMENTTREATY_2"*.
//
// ---------------------------------------------------------------------
// Rantai ekspor (korpus `Treaty In`)
// ---------------------------------------------------------------------
//
//	Section/WorkAttachments   tombol per kategori, syarat `TreatyIn.ViewState
//	                          !='1' || TreatyIn.RevisionState='1'`:
//	                          SetkategoriDoc(IDDoc=.CARI1, Note=.CARIDESC) →
//	                          localAction `TreatyAttachContent` → Refresh
//	                          `GetMasterTreatyCategory_Act`
//	FlowAction TreatyAttachContent  "ASM Attach Content" — `pyAttachmentScreen`
//	                          (pilih berkas), tombol Attach / Cancel, pasca
//	                          `TreatySaveAttachment`
//	TreatySaveAttachment      per berkas: [1.1] pyCategory = CARI41 ·
//	                          [1.4] InsertGoogleStorage_Act(Folder "Contract",
//	                          Namafile, Image, Durasi "1800", Ext) ·
//	                          [1.6] InsertAttachment2_Sql bila objek terkirim
//	InsertGoogleStorage_Act   [3] ext huruf kecil · [4] GetMimeType ·
//	                          [5] KELUAR bila ext kosong / Image < 10 /
//	                          octet-stream · [6] GetAppName_SQL · [7] token ·
//	                          [8] Folder = Folder+"/Doc/"+YYYY+"/"+MM+"/",
//	                          Namafile = "yyyyMMdd-hhmmss-S - "+nama ·
//	                          [10] GetLinkService Google/upload · [11]
//	                          Connect-REST `ServiceGoogle` · [12] URLImage
//	                          kosong = gagal; IMAGEID → Insert_T_Storage_SQL
//
// ⚠️ PENYIMPANGAN yang dinyatakan:
//   - Pega MELEWATI berkas yang ditolak [5] tanpa kata; di sini hasil tiap
//     berkas DILAPORKAN (diunggah / ditolak dengan alasannya).
//   - Spanduk "Recommended safe substitute should be . or _" TIDAK
//     ditegakkan saat unggah — Pega pun tidak (data memuat nama berspasi dan
//     berkurung); hanya garis miring yang ditolak (jalur berkas).
//   - Klien layanan disalin dari pola `masterproductnamelife`
//     (`mpnl_storage.go`) — modul tidak boleh saling impor; potongan
//     bersamanya (`inti/backend/layanan`, `inti/backend/unggah`) dipakai.

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
	"regexp"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/treatyin/backend/models"
)

var (
	// ErrSimpananBelumSiap - alamat `M_LINK_SERVICE`, App, atau token
	// penyimpanan tidak tersedia (503).
	ErrSimpananBelumSiap = errors.New("services: the storage service is not ready")
	// ErrSimpananGagal - layanan penyimpanan tak terjangkau, menjawab galat,
	// atau jawabannya rusak (502).
	ErrSimpananGagal = errors.New("services: the storage service failed")
)

const (
	// folderLampiran / durasiLampiran - parameter `TreatySaveAttachment` [1.4].
	folderLampiran = "Contract"
	durasiLampiran = 1800
	// batasWaktuSimpanan - `ServiceGoogle` `pyResponseTimeout 300000`.
	batasWaktuSimpanan   = 300 * time.Second
	batasJawabanSimpanan = 1 << 20
)

// zonaLampiran - `@CurrentDate(…, "Asia/Jakarta")`.
var zonaLampiran = time.FixedZone("WIB", 7*3600)

// PermintaanSimpanan - halaman `UploadDoc` sebagai JSON (nama medan VERBATIM;
// `Durasi` angka — `@replaceAll` langkah [9] membuang kutipnya).
//
// ⚠️ `omitempty` pada medan yang tidak semua efek kirim: hapus
// (`DeleteGoogleStorage_Act` [6]) hanya App/Kodestring/Namafile, geturl
// (`GetUrlGoogleStorage_Act` [6.2]) tanpa Image.
type PermintaanSimpanan struct {
	App        string `json:"App"`
	Kodestring string `json:"Kodestring"`
	Durasi     int    `json:"Durasi,omitempty"`
	Folder     string `json:"Folder,omitempty"`
	Namafile   string `json:"Namafile"`
	Image      string `json:"Image,omitempty"`
	Ext        string `json:"ext,omitempty"`
	MimeType   string `json:"MimeType,omitempty"`
}

// JawabanSimpanan - `UploadDoc.Response`.
type JawabanSimpanan struct {
	URLImage  string `json:"URLImage"`
	Exp       string `json:"exp"`
	AppFolder string `json:"appfolder"`
	DateTime  string `json:"DateTime"`
}

// PengirimSimpanan - Connect-REST `ServiceGoogle`, satu metode per kunci
// `M_LINK_SERVICE`. `Kodestring` (token) diisi pengirimnya sendiri.
type PengirimSimpanan interface {
	// Unggah - Google/upload (`InsertGoogleStorage_Act`).
	Unggah(ctx context.Context, p PermintaanSimpanan) (JawabanSimpanan, error)
	// URLBaru - Google/geturl (`GetUrlGoogleStorage_Act` [6]).
	URLBaru(ctx context.Context, p PermintaanSimpanan) (JawabanSimpanan, error)
	// Hapus - Google/delete (`DeleteGoogleStorage_Act`).
	Hapus(ctx context.Context, p PermintaanSimpanan) error
}

type pengirimGoogle struct {
	alamat func(ctx context.Context, kunci layanan.KunciLayanan) (string, error)
	token  func(ctx context.Context, app string) (string, error)
	klien  *http.Client
}

// pengirimSimpanan - penyimpanan NYATA, seperti XML: alamat `M_LINK_SERVICE`
// dan token `GCP_IMAGE` dibaca saat jalan.
func (s *Service) pengirimSimpanan() PengirimSimpanan {
	return pengirimGoogle{
		alamat: func(ctx context.Context, kunci layanan.KunciLayanan) (string, error) {
			if !s.PunyaDatabase() {
				return "", db.ErrTanpaOracle
			}
			return layanan.AlamatLayanan(ctx, layanan.ResolverLinkServiceOracle(s), kunci)
		},
		token: func(ctx context.Context, app string) (string, error) {
			if !s.PunyaDatabase() {
				return "", db.ErrTanpaOracle
			}
			var tok string
			err := s.DalamTransaksi(ctx, func(tx *db.Tx) error {
				var err error
				tok, err = tokenSimpanan(ctx, tx, layanan.NewPenyimpanToken(s.DB()), s.garamToken, app, time.Now())
				return err
			})
			return tok, err
		},
		klien: &http.Client{Timeout: batasWaktuSimpanan},
	}
}

func (p pengirimGoogle) Unggah(ctx context.Context, badan PermintaanSimpanan) (JawabanSimpanan, error) {
	return p.panggil(ctx, layanan.KunciUnggahBerkas, badan)
}

func (p pengirimGoogle) URLBaru(ctx context.Context, badan PermintaanSimpanan) (JawabanSimpanan, error) {
	return p.panggil(ctx, layanan.KunciURLBerkas, badan)
}

func (p pengirimGoogle) Hapus(ctx context.Context, badan PermintaanSimpanan) error {
	_, err := p.panggil(ctx, layanan.KunciHapusBerkas, badan)
	return err
}

// panggil - token, alamat, POST JSON, jawaban — satu Connect-REST.
func (p pengirimGoogle) panggil(ctx context.Context, kunci layanan.KunciLayanan, badan PermintaanSimpanan) (JawabanSimpanan, error) {
	tok, err := p.token(ctx, badan.App)
	if err != nil {
		return JawabanSimpanan{}, fmt.Errorf("%w: the storage token is not available (%v)", ErrSimpananBelumSiap, sebabBernama(err))
	}
	alamat, err := p.alamat(ctx, kunci)
	if err != nil {
		return JawabanSimpanan{}, fmt.Errorf("%w: the M_LINK_SERVICE address (%s, %s) is not available (%v)",
			ErrSimpananBelumSiap, kunci.Kategori1, kunci.Kategori2, sebabBernama(err))
	}
	badan.Kodestring = tok
	isi, err := json.Marshal(badan)
	if err != nil {
		return JawabanSimpanan{}, fmt.Errorf("services: building the storage request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(isi))
	if err != nil {
		// ⛔ Alamat tidak disebut — pesan `url.Parse` memuatnya.
		return JawabanSimpanan{}, fmt.Errorf("%w: the M_LINK_SERVICE address is malformed", ErrSimpananBelumSiap)
	}
	req.Header.Set("Content-Type", "application/json")
	klien := *p.klien
	// Badan memuat token — pengalihan TIDAK diikuti.
	klien.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	jwb, err := klien.Do(req)
	if err != nil {
		var ue *url.Error
		if errors.As(err, &ue) && ue.Timeout() {
			return JawabanSimpanan{}, fmt.Errorf("%w: %s timed out", ErrSimpananGagal, kunci.Kategori2)
		}
		return JawabanSimpanan{}, fmt.Errorf("%w: %s is unreachable", ErrSimpananGagal, kunci.Kategori2)
	}
	defer func() { _ = jwb.Body.Close() }()
	if jwb.StatusCode < 200 || jwb.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, io.LimitReader(jwb.Body, batasJawabanSimpanan))
		return JawabanSimpanan{}, fmt.Errorf("%w: %s answered status %d", ErrSimpananGagal, kunci.Kategori2, jwb.StatusCode)
	}
	var j JawabanSimpanan
	if err := json.NewDecoder(io.LimitReader(jwb.Body, batasJawabanSimpanan)).Decode(&j); err != nil && !errors.Is(err, io.EOF) {
		return JawabanSimpanan{}, fmt.Errorf("%w: the %s answer is not JSON", ErrSimpananGagal, kunci.Kategori2)
	}
	return j, nil
}

// sebabBernama - hanya galat BERNAMA yang ikut ke pesan (kalimatnya menyebut
// kunci, tidak pernah nilai); sebab lain bisa memuat alamat.
func sebabBernama(err error) string {
	for _, b := range []error{layanan.ErrGaramTokenKosong, layanan.ErrAppNameKosong, layanan.ErrEndpointTidakDitemukan,
		layanan.ErrAlamatLayananTidakAda, db.ErrTanpaOracle} {
		if errors.Is(err, b) {
			return b.Error()
		}
	}
	return "see server log"
}

// tokenSimpanan - GET TOKEN (`GetTokenStorage_SQL` → procedure
// `GET_TOKEN_STORAGE`) seperti `masterproductnamelife`: token terbaru yang
// masih berlaku dipakai ulang tanpa garam; selain itu token baru dirakit
// dengan garam (`STORAGE_TOKEN_SALT`) dan disimpan semenit.
func tokenSimpanan(ctx context.Context, tx *db.Tx, p *layanan.PenyimpanToken, garam, app string, saat time.Time) (string, error) {
	if strings.TrimSpace(app) == "" {
		return "", layanan.ErrAppNameKosong
	}
	lama, err := p.TokenBerlaku(ctx, tx, app, saat)
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

// polaExpSimpanan - `YYYYMMDDTHHMMSS` di awal `exp` sesudah `-`/`:` dibuang.
var polaExpSimpanan = regexp.MustCompile(`^\d{8}T\d{6}`)

// expSimpanan - `exp` jawaban → bentuk `To_date(…, 'DD/MM/YYYY HH24:MI:SS')`
// (`InsertGoogleStorage_Act` [12.3]). Tak terbaca → kosong (NULL).
func expSimpanan(exp string) string {
	awal := polaExpSimpanan.FindString(strings.NewReplacer("-", "", ":", "").Replace(strings.TrimSpace(exp)))
	if awal == "" {
		return ""
	}
	t, err := time.Parse("20060102T150405", awal)
	if err != nil {
		return ""
	}
	return t.Format("02/01/2006 15:04:05")
}

// BerkasUnggah - satu berkas yang layar kirim.
type BerkasUnggah struct {
	Nama string
	Isi  []byte
}

// MasukanUnggahLampiran - tombol `Attach` modal `ASM Attach Content`.
type MasukanUnggahLampiran struct {
	IDKontrak    string
	KodeKategori string
	Berkas       []BerkasUnggah
}

// HasilBerkasUnggah - nasib SATU berkas.
type HasilBerkasUnggah struct {
	Nama     string `json:"nama"`
	Berhasil bool   `json:"berhasil"`
	Pesan    string `json:"pesan"`
}

// PanelLampiran - isi panel Attachment satu kontrak.
type PanelLampiran struct {
	Lampiran         []models.BarisLampiranWarisan  `json:"lampiran"`
	KategoriLampiran []models.BarisKategoriLampiran `json:"kategoriLampiran"`
}

// HasilUnggahLampiran - hasil per berkas, lalu panel yang dibaca ULANG
// (Refresh `GetMasterTreatyCategory_Act`).
type HasilUnggahLampiran struct {
	Berkas []HasilBerkasUnggah `json:"berkas"`
	PanelLampiran
}

// BacaPanelLampiran - panel Attachment satu kontrak, tanpa membaca ulang
// seluruh form (isian yang belum di-Save tidak hilang).
func (l *Layanan) BacaPanelLampiran(ctx context.Context, p inti.Pelaku, id string) (PanelLampiran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return PanelLampiran{}, err
	}
	kepala, ada, err := l.gudang.BacaKepalaTreatyIn(ctx, strings.TrimSpace(id))
	if err != nil {
		return PanelLampiran{}, err
	}
	if !ada {
		return PanelLampiran{}, fmt.Errorf("%w: %s", ErrKontrakTidakAda, id)
	}
	return l.susunPanelLampiran(ctx, strings.TrimSpace(id), SifatProporsional(teksDok(kepala, "ProportionType")))
}

func (l *Layanan) susunPanelLampiran(ctx context.Context, id string, proporsional bool) (PanelLampiran, error) {
	lampiran, err := l.gudang.BacaLampiranKontrak(ctx, id)
	if err != nil {
		return PanelLampiran{}, err
	}
	katalog, err := l.gudang.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		return PanelLampiran{}, err
	}
	nama := models.NamaKategoriLampiranProp
	if !proporsional {
		nama = models.NamaKategoriLampiranNonProp
	}
	kategori := SusunKategoriLampiran(nama, katalog, lampiran)
	for i := range lampiran {
		lampiran[i].Diunggah = TanggalTampil(lampiran[i].Diunggah)
	}
	return PanelLampiran{Lampiran: lampiran, KategoriLampiran: kategori}, nil
}

// UnggahLampiran - `TreatySaveAttachment` atas setiap berkas yang dipilih.
func (l *Layanan) UnggahLampiran(ctx context.Context, p inti.Pelaku, m MasukanUnggahLampiran) (HasilUnggahLampiran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilUnggahLampiran{}, err
	}
	id := strings.TrimSpace(m.IDKontrak)
	kepala, ada, err := l.gudang.BacaKepalaTreatyIn(ctx, id)
	if err != nil {
		return HasilUnggahLampiran{}, err
	}
	if !ada {
		return HasilUnggahLampiran{}, ditolak("Simpan kontrak lebih dulu — lampiran menempel pada ID kontrak yang sudah tersimpan.")
	}
	proporsional := SifatProporsional(teksDok(kepala, "ProportionType"))
	katalog, err := l.gudang.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		return HasilUnggahLampiran{}, err
	}
	kode := strings.TrimSpace(m.KodeKategori)
	namaKategori, adaKategori := katalog[kode]
	if !adaKategori {
		return HasilUnggahLampiran{}, ditolak(fmt.Sprintf("Kategori %q tidak ada di M_KATEGORIMASTERTREATY.", kode))
	}
	// [1.1] `.pyCategory = StatusDoc.CARI41` — nama yang panel TAMPILKAN.
	namaKategori = namaKategoriTampil(namaKategori, proporsional)
	if len(m.Berkas) == 0 {
		return HasilUnggahLampiran{}, ditolak("Tidak ada berkas yang dipilih.")
	}
	if l.simpanan == nil {
		return HasilUnggahLampiran{}, fmt.Errorf("%w: the storage sender is not wired", ErrSimpananBelumSiap)
	}

	var app string
	hasil := HasilUnggahLampiran{}
	for _, b := range m.Berkas {
		h := HasilBerkasUnggah{Nama: b.Nama}
		if err := l.unggahSatu(ctx, p, id, kode, namaKategori, b, &app); err != nil {
			var tolak galatTombol
			if !errors.As(err, &tolak) {
				// Layanan/basis data gagal — berkas berikutnya bernasib sama;
				// yang sudah terunggah tetap tercatat (satu transaksi per berkas).
				return HasilUnggahLampiran{}, err
			}
			// Berkas ini ditolak [5]; berkas berikutnya tetap dicoba.
			h.Pesan = tolak.pesan
		} else {
			h.Berhasil, h.Pesan = true, "Terunggah"
		}
		hasil.Berkas = append(hasil.Berkas, h)
	}
	panel, err := l.susunPanelLampiran(ctx, id, proporsional)
	if err != nil {
		return HasilUnggahLampiran{}, err
	}
	hasil.PanelLampiran = panel
	return hasil, nil
}

// unggahSatu - [1.4] `InsertGoogleStorage_Act` lalu [1.6] pencatatan.
func (l *Layanan) unggahSatu(ctx context.Context, p inti.Pelaku, id, kode, namaKategori string, b BerkasUnggah, app *string) error {
	nama := strings.TrimSpace(b.Nama)
	if nama == "" || strings.ContainsAny(nama, `/\`) {
		return ditolak("Nama berkas tidak sah.")
	}
	// [3]–[5] ext huruf kecil, `GetMimeType`, keluar bila tak dikenal.
	ext := ""
	if i := strings.LastIndex(nama, "."); i >= 0 && i < len(nama)-1 {
		ext = strings.ToLower(nama[i+1:])
	}
	mime := unggah.MimeDariNamaFile(nama)
	gambar := base64.StdEncoding.EncodeToString(b.Isi)
	if ext == "" || len(gambar) < 10 || mime == unggah.MimeBawaan {
		return ditolak(fmt.Sprintf("Jenis berkas .%s tidak dikenal (GetMimeType) atau berkas kosong — tidak diunggah.", ext))
	}
	if len(b.Isi) > unggah.BatasUkuranUnggahan {
		return ditolak(fmt.Sprintf("Berkas melebihi %d MB.", unggah.BatasUkuranUnggahan>>20))
	}
	// [6] `GetAppName_SQL` — sekali per penekanan.
	if *app == "" {
		a, err := l.gudang.NamaAplikasiSimpanan(ctx)
		if err != nil {
			return err
		}
		if strings.TrimSpace(a) == "" {
			return fmt.Errorf("%w: T_FOLDER_IMAGE.APPNAME is empty", ErrSimpananBelumSiap)
		}
		*app = a
	}
	// [8] Set Data — jam Asia/Jakarta.
	kini := time.Now().In(zonaLampiran)
	folder := folderLampiran + "/Doc/" + kini.Format("2006") + "/" + kini.Format("01") + "/"
	namaObjek := NamaObjekLampiran(kini, nama)
	j, err := l.simpanan.Unggah(ctx, PermintaanSimpanan{
		App: *app, Durasi: durasiLampiran, Folder: folder, Namafile: namaObjek, Image: gambar, Ext: ext, MimeType: mime,
	})
	if err != nil {
		return err
	}
	// [12] `UploadDoc.Response.URLImage == ""` → objek tidak tercatat.
	if strings.TrimSpace(j.URLImage) == "" {
		return fmt.Errorf("%w: upload answered an empty URLImage", ErrSimpananGagal)
	}
	appFolder := j.AppFolder
	if strings.TrimSpace(appFolder) == "" {
		appFolder = (&url.URL{Scheme: "gs", Host: *app, Path: "/"}).String() + folder + namaObjek
	}
	imageID, err := unggah.ImageIDBaru(time.Now())
	if err != nil {
		return err
	}
	_, err = l.gudang.CatatLampiran(ctx, models.LampiranBaru{
		IDKontrak: id, KodeKategori: kode, NamaKategori: namaKategori, NamaBerkas: nama, Ekstensi: ext,
		Pengguna: p.AkunID, ImageID: imageID, URLPublik: j.URLImage, AppFolder: appFolder,
		Exp: expSimpanan(j.Exp), NamaObjek: namaObjek, App: *app,
	})
	return err
}

// NamaObjekLampiran - `@CurrentDate("yyyyMMdd-hhmmss-S","Asia/Jakarta")+" - "+
// Param.Namafile` (`InsertGoogleStorage_Act` [8]). ⚠️ `hh` jam 12 dan `S`
// milidetik TANPA nol pengisi — format Java apa adanya.
func NamaObjekLampiran(kini time.Time, nama string) string {
	return kini.Format("20060102-030405") + "-" + strconv.Itoa(kini.Nanosecond()/int(time.Millisecond)) + " - " + nama
}
