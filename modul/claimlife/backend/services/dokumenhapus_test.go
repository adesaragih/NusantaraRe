package services_test

import (
	"os"
	"strings"
	"testing"

	"nusantarare/modul/claimlife/backend/services"
)

// 11-10-2026 ("buat document bisa dihapus"): Delete dokumen yang masih "Processing" membatalkan efek unggahnya di
// transaksi yang sama dan membuang berkas lokal SESUDAH commit; dokumen berstorage tetap lewat efek storage-hapus.
func TestHapusDokumenTertundaMembatalkanEfekDanBerkas(t *testing.T) {
	isi, err := os.ReadFile("unggahan.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	fn := s[strings.Index(s, "func (u *Unggahan) Hapus("):]
	fn = fn[:strings.Index(fn, "\n}\n")]
	// 11-10-2026 (Google Storage seperti master): objek Google dihapus LEBIH DULU (gagal = rekam tetap), lalu satu transaksi:
	// baris dokumen + (tertunda: batalkan efek | Google: T_STORAGE_IMAGE | lama: kartu stub lokal), lalu berkas lokal.
	hapusGoogle := strings.Index(fn, "u.penyimpanan.Hapus(ctx, objek)")
	transaksi := strings.Index(fn, "u.svc.DalamTransaksi(")
	hapusBaris := strings.Index(fn, "baca.HapusDokumen(ctx, tx, dokID)")
	batal := strings.Index(fn, "batalkanEfekUnggah(ctx, tx,")
	objek := strings.Index(fn, "baca.HapusObjekPenyimpanan(ctx, tx, objek.ImageID)")
	kartu := strings.Index(fn, "baca.HapusKartuBerkas(ctx, tx, dok.TStorageID)")
	buangBerkas := strings.Index(fn, "os.Remove(jalur)")
	if hapusGoogle < 0 || transaksi < hapusGoogle || hapusBaris < transaksi || batal < hapusBaris || objek < batal ||
		kartu < objek || buangBerkas < kartu {
		t.Errorf("urutan Hapus: google %d, transaksi %d, baris %d, batal %d, objek %d, kartu %d, berkas %d",
			hapusGoogle, transaksi, hapusBaris, batal, objek, kartu, buangBerkas)
	}
	if !strings.Contains(fn, "log.Printf(") {
		t.Error("kegagalan membuang berkas lokal harus dicatat di log (best-effort)")
	}
	if services.CatatanBatalUnggah == "" || !strings.Contains(s, "outbox.StatusEfekSelesai, time.Time{},\n\t\t\tCatatanBatalUnggah, saat)") {
		t.Error("efek unggah dibatalkan dengan catatan pembatalan")
	}
	// Pelaksana melewati dokumen yang sudah dihapus SEBELUM menulis kartu berkas.
	pel := s[strings.Index(s, "func (p *PelaksanaBerkasLokal) unggah("):]
	if a, b := strings.Index(pel, "p.baca.DokumenAda(ctx, tx, id)"), strings.Index(pel, "p.baca.SisipKartuBerkas("); a < 0 || a > b {
		t.Errorf("pelaksana unggah tidak memeriksa dokumen ada sebelum menulis kartu berkas (%d, %d)", a, b)
	}
	// View lokal: Unduh tidak lagi menolak dokumen yang belum terunggah; jalur dari pengenal dokumen.
	und := s[strings.Index(s, "func (u *Unggahan) Unduh("):]
	und = und[:strings.Index(und, "\n}\n")]
	if strings.Contains(und, "ErrDokumenBelumTerunggah, dokID") || !strings.Contains(und, "namaBerkasLokal(strconv.FormatInt(dok.ID, 10), dok.NamaFile)") {
		t.Error("Unduh: cadangan lokal untuk dokumen tertunda, jalur dari pengenal dokumen")
	}
}
