package services

// Penyimpanan dokumen klaim di Google Storage - keputusan work owner 11-10-2026 ("untuk penyimpanan document samakan dengan
// penyimpanan document yang digunakan oleh product name"; "iya" Google Storage dicatat di T_STORAGE_IMAGE). POLA DISALIN dari
// masterproductnamelife services/mpnl_storage.go (tanpa impor modul itu), yang meniru XML yang SAMA dengan Claim Life
// (kelompok `ASM-FW-GISFW-Int-T_STORAGE_IMAGE`):
//
//	unggah  InsertDocument_Act → InsertGoogleStorage_Act (Folder "Claim", Durasi 1800): token GetTokenStorage_SQL, Set Data
//	        (Folder + "/Doc/YYYY/MM/", Namafile "yyyyMMdd-hhmmss-S - nama"), alamat GetLinkService "Google"/"upload",
//	        Connect-REST ServiceGoogle (POST JSON, 300 detik), jawaban → Insert_T_Storage_SQL; URLImage kosong = gagal.
//	unduh   DownloadDocumentClaim → GetUrlGoogleStorage_Act: URL tersimpan dipakai selama EXPDATE belum lewat; selain itu
//	        "Google"/"geturl" lalu Update_T_Storage_SQL. Isi diunduh BACKEND dari URL bertanda tangan (https, tanpa pengalihan).
//	hapus   DeleteDocument_Act → DeleteGoogleStorage_Act: "Google"/"delete" dengan jalur objek APPFOLDER tanpa awalan gs+App;
//	        gagal = rekam tetap; lalu DeleteStorage_SQL dan baris dokumen.
//
// ⛔ Alamat dari M_LINK_SERVICE saat jalan; token lewat inti/backend/layanan (TokenStorageOracle, garam STORAGE_TOKEN_SALT dari
// inti/config). Alamat, token, dan garam tidak pernah masuk pesan galat maupun log. Inti TIDAK diubah.
// ⛔ Gagal tertutup: tanpa alamat / APPNAME / garam = ErrStorageBelumSiap ("Storage is not configured"); berkas lokal
// (UNGGAHAN_DIR) tidak disentuh, dokumen tetap "Processing" dan tetap dapat di-View dari salinan lokal.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/claimlife/backend/repository"
)

// Parameter InsertDocument_Act (VERBATIM).
const (
	FolderDokumenKlaim = "Claim"
	DurasiDokumenKlaim = 1800
	// BatasWaktuStorage - ServiceGoogle.xml pyResponseTimeout 300000.
	BatasWaktuStorage   = 300 * time.Second
	batasJawabanStorage = 1 << 20
	// formatExp - bentuk EXPDATE yang ditulis dan dibaca (TO_DATE/TO_CHAR 'DD/MM/YYYY HH24:MI:SS').
	formatExp = "02/01/2006 15:04:05"
)

var (
	// ErrStorageGagal - layanan penyimpanan tak terjangkau / menjawab galat / jawabannya rusak (502).
	ErrStorageGagal = errors.New("Storage service failed")
	// ErrStorageBelumSiap - alamat M_LINK_SERVICE, APPNAME, atau token penyimpanan tidak tersedia (503).
	ErrStorageBelumSiap = errors.New("Storage is not configured")
	// ErrBerkasTidakDiStorage - URL bertanda tangan menjawab 404 (409).
	ErrBerkasTidakDiStorage = errors.New("The document file is not in storage")
	// ErrBerkasDokumenTidakAda - berkas lokal dokumen (belum tersimpan / dokumen lama) tidak ada (404).
	ErrBerkasDokumenTidakAda = errors.New("services: berkas dokumen tidak ditemukan")
)

var zonaJakartaStorage = time.FixedZone("WIB", 7*3600)

// galatStorage - kalimat layar tanpa sebab mentah; jenisnya terbaca errors.Is.
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

// PesanLayar - kalimat yang boleh tampil (tanpa alamat/token).
func (g galatStorage) PesanLayar() string { return g.layar }
func (g galatStorage) Unwrap() []error {
	if g.sebab == nil {
		return []error{g.jenis}
	}
	return []error{g.jenis, g.sebab}
}

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

func gagalStorage(format string, a ...any) error {
	return galatStorage{jenis: ErrStorageGagal, layar: ErrStorageGagal.Error() + ": " + fmt.Sprintf(format, a...)}
}

// PesanStorage - kalimat layar galat penyimpanan (tanpa rincian teknis); "" bila bukan galat penyimpanan.
func PesanStorage(err error) string {
	var g galatStorage
	if errors.As(err, &g) {
		return g.PesanLayar()
	}
	return ""
}

// PenyimpananDokumen - penyimpanan objek dokumen klaim di balik antarmuka (uji memakai tiruan HTTP / tiruan antarmuka).
type PenyimpananDokumen interface {
	// Kirim mengirim berkas lokal `jalur` sebagai objek `o` (Folder, Namafile, App, Durasi) dan menjawab objek tercatat.
	Kirim(ctx context.Context, o repository.ObjekPenyimpanan, jalur, ext, mime string) (repository.ObjekPenyimpanan, error)
	// Buka membaca isi objek; objek baru tidak nil bila URL bertanda tangannya diperbarui.
	Buka(ctx context.Context, o repository.ObjekPenyimpanan) (io.ReadCloser, *repository.ObjekPenyimpanan, error)
	// Hapus membuang objek di penyimpanan.
	Hapus(ctx context.Context, o repository.ObjekPenyimpanan) error
}

// namaObjekDokumen - InsertGoogleStorage_Act Set Data: Folder "Claim/Doc/YYYY/MM/" dan Namafile "yyyyMMdd-hhmmss-S - nama"
// (pola Java: jam 12-an, milidetik tanpa nol depan), waktu Asia/Jakarta.
func namaObjekDokumen(saat time.Time, nama string) (folder, file string) {
	w := saat.In(zonaJakartaStorage)
	folder = FolderDokumenKlaim + "/Doc/" + w.Format("2006") + "/" + w.Format("01") + "/"
	file = fmt.Sprintf("%s-%d - %s", w.Format("20060102-030405"), w.Nanosecond()/int(time.Millisecond), nama)
	return folder, file
}

// PelaksanaGoogle menjalankan efek "storage-unggah" dokumen klaim ke Google Storage dan mencatat objeknya di T_STORAGE_IMAGE
// (pengganti PelaksanaBerkasLokal untuk dokumen baru).
type PelaksanaGoogle struct {
	baca        *repository.KlaimLife
	penyimpanan PenyimpananDokumen
	jam         func() time.Time
}

// NewPelaksanaGoogle menyusunnya.
func NewPelaksanaGoogle(svc *Service, p PenyimpananDokumen) *PelaksanaGoogle {
	return &PelaksanaGoogle{baca: repository.NewKlaimLife(svc.DB()), penyimpanan: p, jam: time.Now}
}

// Laksanakan - hanya jenis "storage-unggah"; jenis lain galat permanen (tidak dipungut kirimSegera).
func (p *PelaksanaGoogle) Laksanakan(ctx context.Context, tx *db.Tx, b outbox.BarisEfekKeluar) error {
	if b.Jenis != unggah.JenisEfekStorageUnggah {
		return fmt.Errorf("%w: jenis efek %q bukan milik pelaksana dokumen", galat.ErrPermintaanTidakSah, b.Jenis)
	}
	id, err := strconv.ParseInt(b.Rujukan, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: rujukan %q bukan pengenal dokumen", galat.ErrPermintaanTidakSah, b.Rujukan)
	}
	ada, err := p.baca.DokumenAda(ctx, tx, id)
	if err != nil {
		return err
	}
	if !ada {
		// Dokumen sudah dihapus sebelum dikirim - tidak ada yang diunggah.
		return nil
	}
	app, err := p.baca.NamaAplikasiPenyimpanan(ctx, tx)
	if err != nil {
		return err
	}
	if app == "" {
		return belumSiap("the T_FOLDER_IMAGE APPNAME", layanan.ErrAppNameKosong)
	}
	saat := p.jam()
	imageID, err := unggah.ImageIDBaru(saat)
	if err != nil {
		return err
	}
	nama := namaDariMuatan(b.Muatan)
	folder, file := namaObjekDokumen(saat, nama)
	o, err := p.penyimpanan.Kirim(ctx, repository.ObjekPenyimpanan{ImageID: imageID, AppFolder: folder, FileName: file,
		AppName: app, DurasiDetik: DurasiDokumenKlaim, Exp: saat.UTC().Add(DurasiDokumenKlaim * time.Second).Format(formatExp)},
		jalurDariMuatan(b.Muatan), strings.TrimPrefix(ekstensi(nama), "."), unggah.MimeDariNamaFile(nama))
	if err != nil {
		return err
	}
	if err := p.baca.SisipObjekPenyimpanan(ctx, tx, o); err != nil {
		return err
	}
	return p.baca.TandaiDokumenTerunggah(ctx, tx, id, imageID)
}

// penyimpananDokumen - transport Google Storage untuk Claim Life.
//
// ⛔ BELUM DIPASANG (11-10-2026): penjaga inti TestNolAlamatLayananDiKode (inti/backend/penjaga/lintasaplikasi_test.go)
// menolak klien HTTP keluar di modul kecuali berkasnya ada di daftar izin MILIK INTI (masterproductnamelife dan
// treatycontractout sudah ada di sana). Menambah izin = perubahan inti; menunggu tim inti / work owner. Transport yang sudah
// ditulis dan diuji (salinan pola mpnl_storage.go) siap dipasang sebagai berkas services/penyimpanangoogle_transport.go.
// Sampai itu, nil: unggahan memakai PelaksanaBerkasLokal (kartu T_CLAIMLF_STORAGE, perilaku sebelumnya); objek
// T_STORAGE_IMAGE tidak lahir; View/Delete objek T_STORAGE_IMAGE gagal tertutup ErrStorageBelumSiap.
func (s *Service) penyimpananDokumen() PenyimpananDokumen { return nil }

// ErrTransportBelumDisetujui - penyimpanan dokumen Google belum dipasang (izin klien HTTP di penjaga inti).
var ErrTransportBelumDisetujui = belumSiap("the Google Storage transport for Claim Life (pending core approval)", nil)
