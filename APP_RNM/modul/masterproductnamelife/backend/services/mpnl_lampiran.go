package services

// Lampiran produk (paket 8, tiket 08–09, PARITAS §6).
//
//	`Add attachment` b64747 → `ProductNameAttachContent` (submit `Attach` b24) → `ProductNameSaveAttachment`:
//	   2.1 b475 `·` `.pyCategory = "File"`
//	   2.4 b904 `·` `InsertGoogleStorage_Act` (Durasi 1800, Folder "Contract"; gagal → keluar b935)
//	       3 b566 `.ext = @toLowerCase(Ext)`; 4 b714 `GetMimeType`; 5 b856 Exit bila ext kosong / octet-stream
//	   2.6 b1201 `·` PRE `CARI51==""` T=3 → `InsertAttachProdName_Sql` hanya bila objek terkirim
//	`Refresh` b65270 / `View` → `LoadAttachmentProdName`;  nama berkas b68903 → `DownloadAttProdName_Act`
//	`Delete` b69714 → `DeleteAttacProdName_act`: 2 b411 `DeleteGoogleStorage_Act` (gagal → keluar), 3 b528 hapus rekam
//
// ⛔ PENYIMPANGAN SADAR (P5, tiket 08–09 `[keputusan work owner]`): Pega
// mengirim berkas DULU lalu merekam; di sini rekam + antrean outbox lebih dulu
// (satu transaksi), lalu pengiriman sebagai EFEK KELUAR - kegagalannya tercatat,
// terlihat (`status = gagal` + galat), dan dapat diulang; produk tidak pernah
// ikut batal. Pengirimnya STUB (folder lokal `UNGGAHAN_DIR`): `ServiceGoogle`,
// `LinkService`, token storage TIDAK dipanggil, alamatnya tidak ditulis.
// ⛔ Berkas berjenis tak dikenal tabel `GetMimeType` DITOLAK berkalimat - di Pega
// ia diam-diam tidak diunggah (Exit-Activity 5 b856) dan tidak direkam.

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// GudangLampiran - rekam lampiran, objek penyimpanan, dan outbox (bagian Gudang).
type GudangLampiran interface {
	DaftarLampiran(ctx context.Context, produkID string) ([]models.Lampiran, error)
	AmbilLampiran(ctx context.Context, tx *db.Tx, produkID, id string) (models.Lampiran, error)
	SisipLampiran(ctx context.Context, tx *db.Tx, l models.Lampiran) (string, error)
	HapusLampiran(ctx context.Context, tx *db.Tx, produkID, id string) error
	CatatObjek(ctx context.Context, tx *db.Tx, o models.ObjekPenyimpanan) error
	HapusObjek(ctx context.Context, tx *db.Tx, imageID string) error
	NamaAplikasi(ctx context.Context, tx *db.Tx) (string, error)
	AntreUnggah(ctx context.Context, tx *db.Tx, lampiranID, muatan string, saat time.Time) error
	AdaUnggahAntre(ctx context.Context, tx *db.Tx, lampiranID string) (bool, error)
	PungutUnggah(ctx context.Context, tx *db.Tx, lampiranID string, saat time.Time) (string, int, bool, error)
	TuntaskanUnggah(ctx context.Context, tx *db.Tx, efekID, status string, jadwal time.Time, galat string, saat time.Time) error
}

// PenyimpananBerkas - penyimpanan berkas di balik antarmuka (tiket 09 AC). Bawaan
// = stub folder lokal (`penyimpananLokal`); uji memakai tiruan.
type PenyimpananBerkas interface {
	// SimpanAntrean menahan isi berkas sampai dikirim.
	SimpanAntrean(ctx context.Context, imageID string, isi io.Reader) error
	// BuangAntrean membuang isi yang belum dikirim (transaksi rekam gagal).
	BuangAntrean(ctx context.Context, imageID string)
	// Kirim memindahkan berkas ke penyimpanan; sudah terkirim = nil (idempoten).
	Kirim(ctx context.Context, imageID string) error
	// Buka membaca berkas terkirim.
	Buka(ctx context.Context, imageID string) (io.ReadCloser, error)
	// Hapus membuang berkas; yang sudah tidak ada bukan galat.
	Hapus(ctx context.Context, imageID string) error
}

// Pesan VERBATIM dan nilai korpus.
const (
	// PesanTanpaBerkas - `ProductNameSaveAttachment` 1 b292 `Local.Err`.
	PesanTanpaBerkas = "Tidak ada file yg diattach"
	// KategoriLampiran - 2.1 b475.
	KategoriLampiran = "File"
	// FolderLampiran / DurasiLampiran - parameter 2.4 b904.
	FolderLampiran = "Contract"
	DurasiLampiran = 1800
)

// ekstensiOffice - syarat tautan `View Office Online` b69291.
var ekstensiOffice = map[string]bool{"xls": true, "xlsx": true, "doc": true, "docx": true, "ppt": true, "pptx": true}

var (
	ErrLampiranTidakAda = repository.ErrLampiranTidakAda
	// ErrNamaLampiranSudahAda - berkas bernama sama sudah terlampir (409, tiket 08 AC).
	ErrNamaLampiranSudahAda = errors.New("services: an attachment with the same file name already exists for this product; " +
		"delete it first or rename the file")
	// ErrLampiranBelumTerkirim - berkas belum ada di penyimpanan (409).
	ErrLampiranBelumTerkirim = errors.New("services: the attachment has not been sent to storage yet; retry sending it first")
	// ErrLampiranSudahTerkirim - "ulangi" atas lampiran terkirim (409).
	ErrLampiranSudahTerkirim = errors.New("services: the attachment is already in storage")
	// ErrBerkasSumberHilang - isi antrean tidak ada dan penyimpanan tidak memilikinya.
	ErrBerkasSumberHilang = errors.New("services: the attachment source file no longer exists; delete the attachment " +
		"and upload it again")
	// ErrBerkasTidakDiStub - rekam dan objek ada, tetapi berkasnya tidak ada di folder stub: berkas
	// yang diunggah Pega tinggal di penyimpanan asalnya sampai pengirim nyata tersambung (409).
	ErrBerkasTidakDiStub = errors.New("services: the attachment file is not in the storage stub; files uploaded " +
		"before this module stay in the original storage until it is connected (OQ-MPNL-10)")
	// ErrOfficeStub - `View Office Online` membungkus URL ke penampil kantor di luar (503, OQ-MPNL-11).
	ErrOfficeStub = errors.New("services: View Office Online is a stub - the external office viewer is not called " +
		"(OQ-MPNL-11); download the file instead")
	// ErrPenyimpananBelumDisetel - folder stub belum dipilih (503).
	ErrPenyimpananBelumDisetel = unggah.ErrUnggahanDirBelumDisetel
)

// zonaJakarta - `@CurrentDate(…, "Asia/Jakarta")`; WIB tanpa musim panas.
var zonaJakarta = time.FixedZone("WIB", 7*3600)

// namaObjek - `InsertGoogleStorage_Act` 8 b1339: Folder dan Namafile.
func namaObjek(saat time.Time, nama string) (folder, file string) {
	w := saat.In(zonaJakarta)
	folder = FolderLampiran + "/Doc/" + w.Format("2006") + "/" + w.Format("01") + "/"
	// Pola Java `yyyyMMdd-hhmmss-S`: jam 12-an, milidetik tanpa nol depan.
	file = fmt.Sprintf("%s-%d - %s", w.Format("20060102-030405"), w.Nanosecond()/int(time.Millisecond), nama)
	return folder, file
}

// ekstensi - `.pyFileMimeType` layar ini: ekstensi sesudah titik terakhir, huruf kecil (3 b566).
func ekstensi(nama string) string {
	i := strings.LastIndex(nama, ".")
	if i < 0 || i == len(nama)-1 {
		return ""
	}
	return strings.ToLower(nama[i+1:])
}

// DaftarLampiran - `LoadAttachmentProdName` (`Refresh` b65270, `View` b74973).
func (l *Layanan) DaftarLampiran(ctx context.Context, p inti.Pelaku, produkID string) ([]models.Lampiran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if _, err := l.gudang.AmbilProduk(ctx, nil, produkID); err != nil {
		return nil, tidakAda(err, ErrProdukTidakAda, produkID)
	}
	return l.gudang.DaftarLampiran(ctx, produkID)
}

// UnggahLampiran - `Attach` b24: rekam + antrean satu transaksi, lalu kirim.
func (l *Layanan) UnggahLampiran(ctx context.Context, p inti.Pelaku, produkID, nama string, isi io.Reader) (models.Lampiran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Lampiran{}, err
	}
	nama = strings.TrimSpace(nama)
	if nama == "" || isi == nil {
		return models.Lampiran{}, GalatValidasi{Pesan: []string{PesanTanpaBerkas}}
	}
	ext := ekstensi(nama)
	if ext == "" || unggah.MimeDariNamaFile(nama) == unggah.MimeBawaan {
		return models.Lampiran{}, GalatValidasi{Pesan: []string{fmt.Sprintf(
			"File type %q is not supported (GetMimeType); the file was not attached", ext)}}
	}
	if _, err := l.gudang.AmbilProduk(ctx, nil, produkID); err != nil {
		return models.Lampiran{}, tidakAda(err, ErrProdukTidakAda, produkID)
	}
	ada, err := l.gudang.DaftarLampiran(ctx, produkID)
	if err != nil {
		return models.Lampiran{}, err
	}
	for _, a := range ada {
		if strings.EqualFold(a.FileName, nama) {
			return models.Lampiran{}, fmt.Errorf("%w: %s", ErrNamaLampiranSudahAda, nama)
		}
	}
	saat := l.jam()
	imageID, err := unggah.ImageIDBaru(saat)
	if err != nil {
		return models.Lampiran{}, err
	}
	if err := l.berkas.SimpanAntrean(ctx, imageID, isi); err != nil {
		if errors.Is(err, unggah.ErrBerkasKosong) {
			return models.Lampiran{}, GalatValidasi{Pesan: []string{PesanTanpaBerkas}}
		}
		if errors.Is(err, unggah.ErrBerkasTerlaluBesar) {
			return models.Lampiran{}, GalatValidasi{Pesan: []string{fmt.Sprintf("file %q exceeds %d MB", nama,
				unggah.BatasUkuranUnggahan>>20)}}
		}
		return models.Lampiran{}, err
	}
	baru := models.Lampiran{ProdukID: produkID, Category: KategoriLampiran, FileName: nama, FileMimeType: ext,
		UserName: p.AkunID, StorageID: imageID}
	var lampiranID string
	err = l.tx(ctx, func(tx *db.Tx) error {
		id, err := l.gudang.SisipLampiran(ctx, tx, baru)
		if err != nil {
			return err
		}
		muatan, _ := json.Marshal(struct {
			LampiranID string `json:"lampiran_id"`
			ImageID    string `json:"image_id"`
			Waktu      string `json:"waktu"`
		}{id, imageID, saat.UTC().Format(time.RFC3339)})
		lampiranID = id
		return l.gudang.AntreUnggah(ctx, tx, id, string(muatan), saat)
	})
	if err != nil {
		l.berkas.BuangAntrean(ctx, imageID)
		return models.Lampiran{}, err
	}
	l.kirimLampiran(ctx, produkID, lampiranID)
	return l.gudang.AmbilLampiran(ctx, nil, produkID, lampiranID)
}

// kirimLampiran - satu percobaan efek keluar; hasilnya DICATAT di outbox, tidak
// dilempar ke pengunggah (lampiran opsional, tiket 08).
func (l *Layanan) kirimLampiran(ctx context.Context, produkID, lampiranID string) {
	saat := l.jam()
	err := l.tx(ctx, func(tx *db.Tx) error {
		efek, percobaan, ada, err := l.gudang.PungutUnggah(ctx, tx, lampiranID, saat)
		if err != nil || !ada {
			return err
		}
		gagal := l.kirimSatu(ctx, tx, produkID, lampiranID, saat)
		if gagal == nil {
			return l.gudang.TuntaskanUnggah(ctx, tx, efek, outbox.StatusEfekSelesai, time.Time{}, "", saat)
		}
		status, jadwal := outbox.StatusEfekAntre, saat.Add(outbox.Backoff(percobaan))
		if percobaan >= outbox.PercobaanMaksimum || errors.Is(gagal, ErrBerkasSumberHilang) {
			status, jadwal = outbox.StatusEfekGagalPermanen, time.Time{}
		}
		return l.gudang.TuntaskanUnggah(ctx, tx, efek, status, jadwal, l.pesanKirim(lampiranID, gagal), saat)
	})
	if err != nil {
		l.catat(fmt.Sprintf("master product name life: recording attachment %s send result failed: %v", lampiranID, err))
	}
}

// PesanKirimGagal - kalimat layar galat kirim lampiran yang bukan galat bernama: teks Oracle dan
// jalur folder stub tidak pernah disimpan di `GALAT_TERAKHIR` atau tampil di layar (audit 02-10-2026).
const PesanKirimGagal = "sending the attachment to storage failed; the details are in the server log"

// pesanKirim - kalimat galat kirim yang DISIMPAN dan ditampilkan: galat bernama dengan kalimatnya
// sendiri, selain itu kalimat tetap; galat aslinya ke log.
func (l *Layanan) pesanKirim(lampiranID string, gagal error) string {
	if errors.Is(gagal, ErrBerkasSumberHilang) || errors.Is(gagal, ErrPenyimpananBelumDisetel) {
		return Pesan(gagal)
	}
	l.catat(fmt.Sprintf("master product name life: sending attachment %s failed: %v", lampiranID, gagal))
	return PesanKirimGagal
}

func (l *Layanan) kirimSatu(ctx context.Context, tx *db.Tx, produkID, lampiranID string, saat time.Time) error {
	a, err := l.gudang.AmbilLampiran(ctx, tx, produkID, lampiranID)
	if err != nil {
		return err
	}
	if err := l.berkas.Kirim(ctx, a.StorageID); err != nil {
		return err
	}
	app, err := l.gudang.NamaAplikasi(ctx, tx)
	if err != nil {
		return err
	}
	if app == "" {
		return errors.New("services: T_FOLDER_IMAGE has no APPNAME; the storage application name is required")
	}
	folder, file := namaObjek(saat, a.FileName)
	return l.gudang.CatatObjek(ctx, tx, models.ObjekPenyimpanan{ImageID: a.StorageID, AppFolder: folder, FileName: file,
		AppName: app, DurasiDetik: DurasiLampiran})
}

// UlangiLampiran - kirim ulang lampiran yang belum terkirim (tiket 09 AC).
func (l *Layanan) UlangiLampiran(ctx context.Context, p inti.Pelaku, produkID, id string) (models.Lampiran, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Lampiran{}, err
	}
	a, err := l.gudang.AmbilLampiran(ctx, nil, produkID, id)
	if err != nil {
		return models.Lampiran{}, err
	}
	if a.Status == models.StatusTerunggah {
		return models.Lampiran{}, ErrLampiranSudahTerkirim
	}
	saat := l.jam()
	err = l.tx(ctx, func(tx *db.Tx) error {
		// Efek antre yang ada dipakai ulang oleh `kirimLampiran` (satu percobaan per klik);
		// efek baru hanya bila tidak ada lagi yang antre (gagal permanen / objek hilang).
		ada, err := l.gudang.AdaUnggahAntre(ctx, tx, id)
		if err != nil || ada {
			return err
		}
		return l.gudang.AntreUnggah(ctx, tx, id, `{"lampiran_id":"`+id+`"}`, saat)
	})
	if err != nil {
		return models.Lampiran{}, err
	}
	l.kirimLampiran(ctx, produkID, id)
	return l.gudang.AmbilLampiran(ctx, nil, produkID, id)
}

// BerkasUnduhan - isi satu lampiran untuk diunduh.
type BerkasUnduhan struct {
	Nama string
	Mime string
	Isi  io.ReadCloser
}

// UnduhLampiran - tautan nama berkas b68903 (`DownloadAttProdName_Act` 6 b953).
func (l *Layanan) UnduhLampiran(ctx context.Context, p inti.Pelaku, produkID, id string) (BerkasUnduhan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return BerkasUnduhan{}, err
	}
	a, err := l.gudang.AmbilLampiran(ctx, nil, produkID, id)
	if err != nil {
		return BerkasUnduhan{}, err
	}
	if a.Status != models.StatusTerunggah {
		return BerkasUnduhan{}, ErrLampiranBelumTerkirim
	}
	isi, err := l.berkas.Buka(ctx, a.StorageID)
	if err != nil {
		if errors.Is(err, ErrPenyimpananBelumDisetel) {
			return BerkasUnduhan{}, err // 503, bukan 409
		}
		// Jalur folder stub dan kunci objek hanya di log, tidak di layar (audit 02-10-2026).
		l.catat(fmt.Sprintf("master product name life: opening attachment %s (storage %s) failed: %v", id, a.StorageID, err))
		return BerkasUnduhan{}, fmt.Errorf("%w: %s", ErrBerkasTidakDiStub, a.FileName)
	}
	return BerkasUnduhan{Nama: a.FileName, Mime: unggah.MimeDariNamaFile(a.FileName), Isi: isi}, nil
}

// NamaDaftarTakTersedia - entri zip yang menyebut lampiran terkirim yang berkasnya
// tidak ada di folder stub (lampiran Pega lama) - kekurangan dinyatakan, bukan disembunyikan.
const NamaDaftarTakTersedia = "_not-available.txt"

// UnduhSemuaLampiran - tombol `Download All` b67657 → satu arsip zip lampiran
// TERKIRIM produk ini (R15, OQ-MPNL-07). Mengembalikan jumlah berkas. Lampiran
// yang berkasnya tidak ada di folder stub dicantumkan di `_not-available.txt`;
// bila TIDAK SATU PUN tersedia, jawabannya 409 berkalimat.
func (l *Layanan) UnduhSemuaLampiran(ctx context.Context, p inti.Pelaku, produkID string, ke io.Writer) (int, error) {
	daftar, err := l.DaftarLampiran(ctx, p, produkID)
	if err != nil {
		return 0, err
	}
	type tersedia struct {
		nama string
		isi  io.ReadCloser
	}
	var ada []tersedia
	var hilang []string
	defer func() {
		for _, t := range ada {
			_ = t.isi.Close()
		}
	}()
	for _, a := range daftar {
		if a.Status != models.StatusTerunggah {
			continue
		}
		isi, err := l.berkas.Buka(ctx, a.StorageID)
		if err != nil {
			hilang = append(hilang, a.FileName)
			continue
		}
		ada = append(ada, tersedia{nama: a.FileName, isi: isi})
	}
	if len(ada) == 0 && len(hilang) > 0 {
		return 0, fmt.Errorf("%w: %s", ErrBerkasTidakDiStub, strings.Join(hilang, ", "))
	}
	z := zip.NewWriter(ke)
	for _, t := range ada {
		w, err := z.Create(t.nama)
		if err == nil {
			_, err = io.Copy(w, t.isi)
		}
		if err != nil {
			return 0, err
		}
	}
	if len(hilang) > 0 {
		w, err := z.Create(NamaDaftarTakTersedia)
		if err == nil {
			_, err = io.WriteString(w, Pesan(ErrBerkasTidakDiStub)+":\r\n"+strings.Join(hilang, "\r\n")+"\r\n")
		}
		if err != nil {
			return 0, err
		}
	}
	return len(ada), z.Close()
}

// HapusLampiran - `Delete` b69714 (`DeleteAttacProdName_act`): berkas dulu (2
// b411; gagal = rekam tetap), lalu objek dan rekam (3 b528) di satu transaksi.
// Berkas yang sudah tidak ada di penyimpanan bukan galat (tiket 09 AC).
func (l *Layanan) HapusLampiran(ctx context.Context, p inti.Pelaku, produkID, id string) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	a, err := l.gudang.AmbilLampiran(ctx, nil, produkID, id)
	if err != nil {
		return err
	}
	if err := l.berkas.Hapus(ctx, a.StorageID); err != nil {
		return err
	}
	saat := l.jam()
	return l.tx(ctx, func(tx *db.Tx) error {
		// Efek yang masih antre untuk lampiran ini ditutup - tidak ada lagi yang dikirim.
		for {
			efek, _, ada, err := l.gudang.PungutUnggah(ctx, tx, id, saat)
			if err != nil {
				return err
			}
			if !ada {
				break
			}
			if err := l.gudang.TuntaskanUnggah(ctx, tx, efek, outbox.StatusEfekGagalPermanen, time.Time{},
				"attachment deleted before it was sent", saat); err != nil {
				return err
			}
		}
		if err := l.gudang.HapusObjek(ctx, tx, a.StorageID); err != nil {
			return err
		}
		return l.gudang.HapusLampiran(ctx, tx, produkID, id)
	})
}

// LihatOffice - tautan `View Office Online` b69291 (stub, OQ-MPNL-11).
func (l *Layanan) LihatOffice(ctx context.Context, p inti.Pelaku, produkID, id string) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	a, err := l.gudang.AmbilLampiran(ctx, nil, produkID, id)
	if err != nil {
		return err
	}
	if !ekstensiOffice[a.FileMimeType] {
		return GalatValidasi{Pesan: []string{fmt.Sprintf(
			"View Office Online is only offered for xls, xlsx, doc, docx, ppt and pptx files, not %q", a.FileMimeType)}}
	}
	return ErrOfficeStub
}
