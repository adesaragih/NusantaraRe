package services

// Tombol `Download All` panel Attachment — `TreatyInDownloadAll` →
// `DownloadAll_Act` (`Section/WorkAttachments.xml`).
//
// ⭐ TAMPIL di Pega: selnya `pyVisible = ALWAYS`, yang MENIMPA `pyCondition =
// never` (pola yang sama dengan grid Co-Ins Scale) — tangkapan layar Pega
// pemakai 8 Oktober 2026 memperlihatkannya di samping `Refresh`.
//
//	[2]  RDB GetAllAttachment2_Sql — semua `M_ATTACHMENTTREATY_2` kontrak
//	     (`treatyid = {TreatyIn.ID}`)
//	[4]–[8] per berkas: GetUrlGoogleStorage_Act → isi (Base64)
//	[10] Java: ZipOutputStream, entri = `pyFileName` (nama berkas apa adanya),
//	     dikirim `AllDocuments.zip` (`pyCaption AllDocuments`, application/zip)
//
// ⚠️ Pega membangun zip di memori (`ByteArrayOutputStream`) — di sini juga,
// supaya galat satu berkas tidak meninggalkan unduhan setengah jadi.
// ⚠️ Nama berkas KEMBAR: Pega melempar galat (ZipEntry kembar) dan seluruh
// unduhan gagal; di sini entri berikutnya diberi akhiran ` (2)`, ` (3)`.

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	inti "nusantarare/inti/backend"
)

// NamaUnduhSemuaLampiran - `pyCaption AllDocuments` + `.zip`.
const NamaUnduhSemuaLampiran = "AllDocuments.zip"

// UnduhSemuaLampiran - zip seluruh lampiran satu kontrak.
func (l *Layanan) UnduhSemuaLampiran(ctx context.Context, p inti.Pelaku, idKontrak string) ([]byte, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	id := strings.TrimSpace(idKontrak)
	if id == "" {
		return nil, ditolak("Pengenal kontrak kosong — lampiran tidak dapat diunduh.")
	}
	daftar, err := l.gudang.BacaLampiranKontrak(ctx, id)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	dipakai := map[string]int{}
	for _, b := range daftar {
		berkas, err := l.IsiLampiran(ctx, p, id, b.ID)
		if err != nil {
			return nil, fmt.Errorf("berkas %s: %w", b.NamaBerkas, err)
		}
		w, err := zw.Create(namaEntriUnik(berkas.Nama, dipakai))
		if err == nil {
			_, err = io.Copy(w, berkas.Isi)
		}
		_ = berkas.Isi.Close()
		if err != nil {
			return nil, fmt.Errorf("berkas %s: %w", b.NamaBerkas, err)
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// namaEntriUnik - nama berkas apa adanya; kembarannya diberi ` (n)` sebelum
// ekstensinya.
func namaEntriUnik(nama string, dipakai map[string]int) string {
	nama = strings.TrimSpace(nama)
	if nama == "" {
		nama = "berkas"
	}
	dipakai[nama]++
	n := dipakai[nama]
	if n == 1 {
		return nama
	}
	ext := path.Ext(nama)
	return fmt.Sprintf("%s (%d)%s", strings.TrimSuffix(nama, ext), n, ext)
}
