package services

// Lampiran master IKUT ke penyesuaian baru — `TreatyInEDMSetValue` [8]
// `TreatyRevisionCopyAttachment` (korpus `Treaty In Adjustment`, memo rule
// "COPY DOCUMENT DARI GOOGLE STORAGE", versi 01-01-88).
//
// ---------------------------------------------------------------------
// Rantai ekspor
// ---------------------------------------------------------------------
//
//	TreatyInEDMSetValue [3]  OLDID = ID (master yang dipilih `Choose`)
//	                    [6]  TreatyInRevisi_post — ID revisi baru
//	                    [8]  TreatyRevisionCopyAttachment (TANPA prasyarat)
//	                    [9]  SaveTreatyIn_EDM_Act
//	TreatyRevisionCopyAttachment
//	  [1]    RDB-List `CopyAllAttachment2_Sql`: `select ID, filename,
//	         CATEGORY, FILEMIMETYPE, T_STORAGE_ID, CATEGORY_ID from
//	         M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.OLDID}`
//	  [2]    per baris (`CopyDoc.pxResults`):
//	  [2.1]  StatusDoc.CARI40 = .category (CATEGORY_ID); CARI51–53 = ""
//	  [2.3]  GetUrlGoogleStorage_Act(ImageID = .type, Durasi 1800) — URL
//	         bertanda tangan objek SUMBER (`tautanBertanda`)
//	  [2.5]  Java: buka URL → byte → Base64; galat DITELAN (base64 "")
//	  [2.7]  InsertGoogleStorage_Act(Ext = .pyFileMimeType, Image = base64,
//	         Folder "Contract", Namafile = .pyFileName, Durasi 1800) —
//	         DILEWATI bila base64 kosong
//	  [2.8]  InsertAttachment2_Sql — DILEWATI bila ImageID kosong (unggah
//	         ditolak [5] atau URLImage kosong [12]): TREATYID = TreatyIn.ID
//	         (pengenal BARU), CATEGORY = .pyCategory, FILENAME = .pyFileName,
//	         FILEMIMETYPE = .pyFileMimeType, DATA_JSON = Datain.CARI50 (tidak
//	         pernah diisi → kosong), USERNAME = operator, CATEGORY_ID =
//	         StatusDoc.CARI40, T_STORAGE_ID = ImageID BARU
//	  [2.9]  Page-Remove Datain, StatusDoc, DropFile
//
// ---------------------------------------------------------------------
// ⭐ Keputusan rancangan
// ---------------------------------------------------------------------
//   - OBJEK DIUNGGAH ULANG, persis Pega — bukan baris yang menunjuk objek
//     yang sama. Salinan semacam itu TIDAK aman di aplikasi ini: `Delete`
//     panel Attachment (`HapusLampiran`) menghapus objek Google Storage DAN
//     baris `T_STORAGE_IMAGE` ber-IMAGEID itu, jadi menghapus lampiran di
//     revisi akan merusak lampiran master (dan sebaliknya).
//   - Tabel yang ditulis SAMA dengan tombol `Upload file`
//     (`CatatLampiran`: `T_STORAGE_IMAGE` + `M_ATTACHMENTTREATY_2`, satu
//     transaksi per berkas) — nol tabel/kolom baru.
//   - SAAT: Pega menyalin ketika `Choose`; di sini `Choose` hanya menyusun
//     draf (keputusan pemilik proses 7 Oktober 2026), jadi salinan dibuat
//     pada tulisan PERTAMA draf (Save, Submit, atau Actions dengan
//     `Draf == true`), SESUDAH kepala dan pendaratan tersimpan.
//   - URUTAN & TRANSAKSI: penyimpanan tidak transaksional, jadi salinan
//     berjalan sesudah transaksi kepala MENGIKAT. Gagalnya salinan TIDAK
//     menggagalkan Save — kepala sudah tersimpan, dan mengembalikan galat
//     akan menyuruh pemakai menekan Save lagi atas draf yang kini ditolak
//     (`ErrPenyesuaianSudahAda`). Per berkas: unduh sumber → unggah objek
//     baru → SATU transaksi dua INSERT. Bila INSERT gagal sesudah unggah
//     berhasil, objek baru tertinggal tanpa baris (yatim) — sama dengan
//     tombol `Upload file`.
//   - DUPLIKAT: Pega tidak berpagar (tiap `Choose` menyalin lagi). Di sini
//     salinan hanya berjalan bila INSERT kepala draf berhasil, dan
//     `SimpanPenyesuaian` menolak draf yang pengenalnya sudah ada — jadi
//     salinan berjalan TEPAT SEKALI per pengenal penyesuaian.
//   - Konteks permintaan dilepas dari pembatalan (`context.WithoutCancel`):
//     peramban yang ditutup di tengah salinan tidak memotongnya di separuh
//     jalan, sebab tidak ada Save kedua yang akan mengulanginya.
//
// ⚠️ PENYIMPANGAN yang dinyatakan:
//   - Pega MENELAN berkas yang gagal ([2.5] catch, prasyarat [2.7]/[2.8]);
//     di sini nasib tiap berkas DILAPORKAN di hasil Save
//     (`HasilSimpan.SalinanLampiran`).
//   - Pagar `OLDID`: draf datang dari layar, jadi `OLDID` yang tidak
//     sejalan dengan pengenal revisi (`IDRevisiBaru`: tujuh aksara pertama
//     sama, lalu `/R`) TIDAK disalin — mencegah lampiran kontrak lain ikut
//     ke penyesuaian ini. Pega membaca `OLDID` dari clipboard server.
//   - Berkas sumber di atas `unggah.BatasUkuranUnggahan` tidak disalin
//     (dilaporkan) — batas yang sama dengan tombol `Upload file`.
//   - Setelah layanan penyimpanan dinyatakan BELUM SIAP (alamat/App/token
//     tidak ada), berkas sisanya tidak dicoba lagi — nasibnya sama.

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/treatyin/backend/models"
)

// HasilSalinLampiran - nasib salinan lampiran master saat draf penyesuaian
// pertama kali tersimpan.
type HasilSalinLampiran struct {
	// Sumber - `TreatyIn.OLDID`, pengenal yang lampirannya disalin.
	Sumber string `json:"sumber"`
	// Tersalin - cacah berkas yang tercatat di pengenal baru.
	Tersalin int `json:"tersalin"`
	// Berkas - nasib tiap baris sumber, urutan `BacaLampiranKontrak`.
	Berkas []HasilBerkasUnggah `json:"berkas"`
	// Pesan - galat menyeluruh (daftar sumber tak terbaca, OLDID tak sah).
	Pesan string `json:"pesan,omitempty"`
}

// salinLampiranDraf - kait jalur tulis penyesuaian: dipanggil dengan hasil
// `tulisPenyesuaian`; hanya draf yang BERHASIL tersimpan yang disalin.
func (l *Layanan) salinLampiranDraf(ctx context.Context, p inti.Pelaku, m MasukanPenyesuaian, doc map[string]any, h HasilSimpan, err error) (HasilSimpan, error) {
	if err != nil || !m.Draf {
		return h, err
	}
	h.SalinanLampiran = l.SalinLampiranRevisi(ctx, p, h.ID, teksDok(doc, "OLDID"))
	return h, nil
}

// OLDIDSejalan - `idLama` memang asal `idBaru` menurut `IDRevisiBaru`
// (`TreatyInRevisi_post` [4]/[5]): tujuh aksara pertama sama, lalu `/R`.
func OLDIDSejalan(idBaru, idLama string) bool {
	idBaru, idLama = strings.TrimSpace(idBaru), strings.TrimSpace(idLama)
	if len(idLama) < 7 || idBaru == idLama {
		return false
	}
	return strings.HasPrefix(idBaru, idLama[:7]+"/R")
}

// SalinLampiranRevisi - `TreatyRevisionCopyAttachment`: setiap lampiran
// `idLama` diunggah ulang dan dicatat di `idBaru`. nil = tidak ada yang
// perlu disalin (OLDID kosong — SQL Pega `treatyid = ''` nol baris — atau
// master tanpa lampiran).
func (l *Layanan) SalinLampiranRevisi(ctx context.Context, p inti.Pelaku, idBaru, idLama string) *HasilSalinLampiran {
	ctx = context.WithoutCancel(ctx)
	idBaru, idLama = strings.TrimSpace(idBaru), strings.TrimSpace(idLama)
	if idLama == "" {
		return nil
	}
	if !OLDIDSejalan(idBaru, idLama) {
		return &HasilSalinLampiran{Sumber: idLama, Berkas: []HasilBerkasUnggah{},
			Pesan: fmt.Sprintf("Original ID %s is not the source of %s — its attachments were not copied.", idLama, idBaru)}
	}
	// [1] `CopyAllAttachment2_Sql`.
	sumber, err := l.gudang.BacaLampiranKontrak(ctx, idLama)
	if err != nil {
		log.Printf("treaty in: copying attachments %s -> %s: %v", idLama, idBaru, err)
		return &HasilSalinLampiran{Sumber: idLama, Berkas: []HasilBerkasUnggah{},
			Pesan: "Attachment list for " + idLama + " could not be read — its attachments were not copied (see the server log)."}
	}
	if len(sumber) == 0 {
		return nil
	}
	hasil := &HasilSalinLampiran{Sumber: idLama, Berkas: []HasilBerkasUnggah{}}
	var app string
	var belumSiap error
	for _, b := range sumber {
		h := HasilBerkasUnggah{Nama: b.NamaBerkas}
		err := belumSiap
		if err == nil {
			err = l.salinSatuLampiran(ctx, p, idBaru, idLama, b, &app)
		}
		if err == nil {
			h.Berhasil, h.Pesan = true, "Tersalin"
			hasil.Tersalin++
		} else {
			h.Pesan = pesanSalinLampiran(idLama, idBaru, b.NamaBerkas, err)
			if errors.Is(err, ErrSimpananBelumSiap) {
				belumSiap = err
			}
		}
		hasil.Berkas = append(hasil.Berkas, h)
	}
	return hasil
}

// salinSatuLampiran - langkah [2.3]–[2.8] atas SATU baris sumber.
func (l *Layanan) salinSatuLampiran(ctx context.Context, p inti.Pelaku, idBaru, idLama string, b models.BarisLampiranWarisan, app *string) error {
	if l.simpanan == nil {
		return fmt.Errorf("%w: the storage sender is not wired", ErrSimpananBelumSiap)
	}
	// [2.3] `GetUrlGoogleStorage_Act(ImageID = DropFile.type)` — URL tersimpan
	// selama berlaku, selain itu Google/geturl + `Update_T_Storage_SQL`.
	_, tautan, err := l.tautanBertanda(ctx, idLama, b.ID)
	if err != nil {
		return err
	}
	// [2.5] Java: isi berkas dari URL → Base64.
	r, err := l.simpanan.Ambil(ctx, tautan)
	if err != nil {
		return err
	}
	isi, err := io.ReadAll(io.LimitReader(r, unggah.BatasUkuranUnggahan+1))
	_ = r.Close()
	if err != nil {
		return fmt.Errorf("%w: reading the source object was interrupted", ErrSimpananGagal)
	}
	if len(isi) > unggah.BatasUkuranUnggahan {
		return ditolak(fmt.Sprintf("File exceeds %d MB — not copied.", unggah.BatasUkuranUnggahan>>20))
	}
	gambar := base64.StdEncoding.EncodeToString(isi)
	// [2.7] prasyarat `Datain.CARI53 == ""` → unggah dilewati.
	if gambar == "" {
		return ditolak("The source file is empty — not copied.")
	}
	// `InsertGoogleStorage_Act` [2] Ext = `DropFile.pyFileMimeType`; kosong →
	// dari Namafile. [3] huruf kecil. [4] `GetMimeType`. [5] keluar.
	nama := b.NamaBerkas
	ext := strings.TrimSpace(b.JenisMime)
	if ext == "" {
		if i := strings.LastIndex(nama, "."); i >= 0 {
			ext = nama[i+1:]
		}
	}
	ext = strings.ToLower(ext)
	mime := unggah.MimeDariNamaFile("berkas." + ext)
	if ext == "" || len(gambar) < 10 || mime == unggah.MimeBawaan {
		return ditolak(fmt.Sprintf("File type .%s is not recognised (GetMimeType) — not copied.", ext))
	}
	// [6] `GetAppName_SQL` — sekali per salinan.
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
	// [8] Set Data — Folder "Contract" (`TreatyRevisionCopyAttachment` [2.7]),
	// jam Asia/Jakarta, Namafile = nama asli sumber.
	kini := time.Now().In(zonaLampiran)
	folder := folderLampiran + "/Doc/" + kini.Format("2006") + "/" + kini.Format("01") + "/"
	namaObjek := NamaObjekLampiran(kini, nama)
	j, err := l.simpanan.Unggah(ctx, PermintaanSimpanan{
		App: *app, Durasi: durasiLampiran, Folder: folder, Namafile: namaObjek, Image: gambar, Ext: ext, MimeType: mime,
	})
	if err != nil {
		return err
	}
	// [12] URLImage kosong → ImageID kosong → [2.8] dilewati.
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
	// [2.8] `InsertAttachment2_Sql` — kolom baris SUMBER apa adanya, TREATYID
	// pengenal baru, USERNAME operator yang menyimpan.
	_, err = l.gudang.CatatLampiran(ctx, models.LampiranBaru{
		IDKontrak: idBaru, KodeKategori: b.KodeKategori, NamaKategori: b.NamaKategori, NamaBerkas: nama,
		Ekstensi: b.JenisMime, Pengguna: p.AkunID, ImageID: imageID, URLPublik: j.URLImage, AppFolder: appFolder,
		Exp: expSimpanan(j.Exp), NamaObjek: namaObjek, App: *app,
	})
	return err
}

// pesanSalinLampiran - kalimat nasib satu berkas. Galat penyimpanan sudah
// dirumuskan tanpa alamat/token (`panggil`); galat lain bisa memuat SQL,
// jadi hanya dicatat di log server.
func pesanSalinLampiran(idLama, idBaru, nama string, err error) string {
	var tolak galatTombol
	if errors.As(err, &tolak) {
		return tolak.pesan
	}
	if errors.Is(err, ErrSimpananBelumSiap) || errors.Is(err, ErrSimpananGagal) {
		return err.Error()
	}
	log.Printf("treaty in: copying attachment %q %s -> %s: %v", nama, idLama, idBaru, err)
	return "Gagal disalin — lihat log server."
}
