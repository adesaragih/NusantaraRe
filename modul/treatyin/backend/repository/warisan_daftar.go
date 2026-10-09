package repository

// Jalur baca tabel WARISAN `POOLDATA.TREATY_IN` - layar daftar kontrak.
//
// ⛔ BERKAS INI MEMBACA TABEL YANG MODUL INI TIDAK MILIKI, dan itu seluruh
// sebab ia berdiri terpisah dari `daftar.go`. `daftar.go` membaca `KONTRAK` +
// `VERSI_KONTRAK` - model BARU, hari ini nol baris. Berkas ini membaca
// `TREATY_IN` - tabel WARISAN, 1.854 baris, keputusan pemilik proses
// 3 Oktober 2026 agar layar pertama menampilkan data nyata.
//
// ⛔ BACA SAJA. Nol `INSERT`, nol `UPDATE`, nol `DELETE`, nol DDL terhadap
// `TREATY_IN`. Ia memuat baris produksi-bayangan dan modul ini bukan
// pemiliknya. Dijaga `TestWarisanHanyaDibaca` di `migrasi_invarian_test.go`,
// yang menyapu seluruh berkas Go modul ini.
//
// ⛔ `TREATY_IN` TIDAK ditulis di satu migrasi pun. Ia sudah ada;
// menuliskannya di migrasi berarti mengklaim kepemilikan yang bukan milik
// modul ini - dan `migrate` berikutnya akan mencoba membuatnya.
//
// ⭐ KAPAN BERKAS INI BOLEH DICABUT: ketika tiket 59 selesai memindahkan
// kepala kontrak warisan ke `KONTRAK`/`VERSI_KONTRAK` DAN cacah keduanya
// cocok. Sampai itu terjadi, layar yang membaca model baru menampilkan
// layar kosong di atas 1.854 baris yang sebenarnya ada.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// TabelWarisanKontrak - nama tabelnya, disebut SEKALI di seluruh modul.
const TabelWarisanKontrak = "TREATY_IN"

// CacahKontrakWarisan menghitung seluruh baris tabel warisan.
//
// Dipisah dari pembacaan halamannya supaya penomoran halaman dapat
// menghitung halaman terakhir tanpa menarik 1.854 baris.
func (g *Gudang) CacahKontrakWarisan(ctx context.Context) (int, error) {
	nama, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return 0, err
	}
	var n int
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s", nama)
	if err := g.db.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: menghitung %s: %w", TabelWarisanKontrak, err)
	}
	return n, nil
}

// DaftarKontrakWarisan membaca SATU halaman, urut pengenal MENURUN.
//
// ⛔ `ORDER BY TO_NUMBER(ID) DESC`, bukan `ORDER BY ID DESC`. Kolomnya teks,
// dan urutan teks menaruh "999999" sesudah "1001856". Ke-1.854 nilainya
// terukur berupa angka tujuh digit, jadi `TO_NUMBER` aman HARI INI - dan
// `DEFAULT 0 ON CONVERSION ERROR` menjaga baris ke-1.855 yang bukan angka
// tidak meledakkan seluruh kueri; ia mendarat di ujung, terlihat, bukan
// hilang.
//
// ⚠️ `OFFSET … FETCH NEXT` menuntut Oracle 12c ke atas. Modul ini sudah
// memakai `FETCH FIRST` di `daftar.go`, jadi syaratnya tidak baru.
func (g *Gudang) DaftarKontrakWarisan(ctx context.Context, offset, batas int) ([]models.BarisDaftarWarisan, error) {
	nama, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID, TREATYCONTRACTNAME, PROPORTIONTYPE, LEADINGREINSSOURCE,
		CEDING, COMMENCEMENT, TERMINATION, POSITIONUSERNAME, STATUSAKSEPTASI, POSITION
		FROM %s
		ORDER BY TO_NUMBER(ID DEFAULT 0 ON CONVERSION ERROR) DESC, ID DESC
		OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`, nama)

	baris, err := g.db.QueryContext(ctx, q, offset, batas)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", TabelWarisanKontrak, err)
	}
	defer func() { _ = baris.Close() }()

	keluar := []models.BarisDaftarWarisan{}
	for baris.Next() {
		var b models.BarisDaftarWarisan
		// ⚠️ Kesepuluh kolomnya NULLABLE - seluruh 20 kolom tabel ini
		// nullable, dan 1.822 dari 1.854 `POSITIONUSERNAME` memang NULL.
		// `sql.NullString` di sepuluh tempat, bukan di satu.
		var id, nk, pt, lrs, cd, cm, tm, pu, sa, po sql.NullString
		if err := baris.Scan(&id, &nk, &pt, &lrs, &cd, &cm, &tm, &pu, &sa, &po); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelWarisanKontrak, err)
		}
		b.ID = id.String
		b.NamaKontrak = nk.String
		b.SifatProporsiAsli = pt.String
		b.AsalBisnis = lrs.String
		b.Cedant = cd.String
		b.TanggalMulaiAsli = cm.String
		b.TanggalBerakhirAsli = tm.String
		b.PosisiKe = pu.String
		b.StatusAkseptasi = sa.String
		b.Posisi = po.String
		// Bentuk layarnya DITERJEMAHKAN DI SERVICES, bukan di sini: lapisan
		// ini mengembalikan apa yang tersimpan.
		keluar = append(keluar, b)
	}
	return keluar, baris.Err()
}
