package repository

// Baca tab Treaty In dari tabel PENDARATAN - bukan lagi dari CLOB.
//
// ⛔ Inilah seluruh guna migrasi 430. Sebelum berkas ini, setiap pembukaan
// form menarik dokumen `JSONDATA` yang rata-rata 28 KB dan terbesar 158 KB,
// lalu mengurainya untuk mengambil tiga larik. Sesudah berkas ini, ketiganya
// datang dari tiga kueri berindex.
//
// ⚠️ SATU MEKANISME untuk prop dan non-prop. Tidak ada cabang menurut
// `PROPORTIONTYPE` di berkas ini, dan tidak boleh ada: yang membedakan kedua
// cabang adalah tab mana yang DITAMPILKAN, bukan dari mana isinya dibaca.
// Sapuan 3 Oktober 2026 menunjukkan bedanya nyata tetapi bukan di sini -
// kontrak proporsional terbesar membawa `Portfolio` dan `ReportingPeriod`,
// non-proporsional membawa `Installment`; keduanya dibaca jalur yang sama
// dan yang kosong mengembalikan nol baris.
//
// ⛔ Nol tafsir. Nilai yang mendarat sebagai teks dibaca kembali sebagai
// teks; yang menerjemahkannya untuk layar adalah services, tempat yang sama
// dengan `TanggalTampil` dan `SifatProporsiTampil`.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// BacaPeriodePelaporan membaca tab Reporting Period, URUT sesuai dokumen.
//
// ⛔ `ORDER BY URUTAN`. Tanpa itu Oracle bebas mengembalikan baris dalam
// urutan apa pun, dan larik Pega BERURUT - `Q 1`, `Q 2`, `Q 3` yang
// tertukar terbaca benar.
func (g *Gudang) BacaPeriodePelaporan(ctx context.Context, masterID string) ([]models.BarisPeriodeWarisan, error) {
	baris, err := g.bacaTab(ctx, "T_TREATY_REPORTING_PERIOD",
		"PERIOD, AUTOCALCULATE, INITIALDATE, SUBMISSIONDUE, CONFIRMATIONDUE, SETTLEMENTDUE", masterID, 6)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisPeriodeWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisPeriodeWarisan{
			Periode:          b[0],
			HitungOtomatis:   b[1],
			TanggalAwal:      b[2],
			JatuhTempoKirim:  b[3],
			JatuhTempoKonfir: b[4],
			JatuhTempoBayar:  b[5],
		})
	}
	return out, nil
}

// BacaPortofolio membaca tab Portfolio.
func (g *Gudang) BacaPortofolio(ctx context.Context, masterID string) ([]models.BarisPortofolioWarisan, error) {
	baris, err := g.bacaTab(ctx, "T_TREATY_PORTFOLIO", "TYPE, TYPEPORTFOLIO, DESCRIPTION", masterID, 3)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisPortofolioWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisPortofolioWarisan{
			Jenis:          b[0],
			JenisPortfolio: b[1],
			Keterangan:     b[2],
		})
	}
	return out, nil
}

// BacaAkumulasi membaca tab Accumulation.
//
// ⚠️ Terisi pada 18 kontrak dari 1.854 - yang terjarang dari kedelapan.
// Kosong di sini adalah keadaan yang WAJAR, bukan kegagalan pembacaan.
func (g *Gudang) BacaAkumulasi(ctx context.Context, masterID string) ([]models.BarisAkumulasiWarisan, error) {
	baris, err := g.bacaTab(ctx, "T_TREATY_ACCUMULATION", "PERIOD, REPORTDATE, SUBDAYS, SUBDUEDATE", masterID, 4)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisAkumulasiWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisAkumulasiWarisan{
			Periode:         b[0],
			TanggalLapor:    b[1],
			HariKirim:       b[2],
			JatuhTempoKirim: b[3],
		})
	}
	return out, nil
}

// bacaTab menjalankan satu kueri tab dan mengembalikan teksnya apa adanya.
//
// NULL menjadi string kosong di sini, dan itu kehilangan yang DISENGAJA:
// pembeda "kuncinya tidak ada" dari "kuncinya kosong" berlaku bagi medan
// kepala kontrak - yang layar tandai sebagai medan mati - bukan bagi sel di
// dalam grid. Sel grid kosong adalah sel kosong.
func (g *Gudang) bacaTab(ctx context.Context, tabel, kolom, masterID string, n int) ([][]string, error) {
	nama, err := g.db.Qualify(tabel)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1 ORDER BY URUTAN", kolom, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", tabel, masterID, err)
	}
	defer func() { _ = rows.Close() }()

	var out [][]string
	for rows.Next() {
		sel := make([]sql.NullString, n)
		tuju := make([]any, n)
		for i := range sel {
			tuju[i] = &sel[i]
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s kontrak %s: %w", tabel, masterID, err)
		}
		teks := make([]string, n)
		for i, s := range sel {
			teks[i] = s.String
		}
		out = append(out, teks)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", tabel, masterID, err)
	}
	return out, nil
}

// BacaEgnpi membaca tab EGNPI. Larik terbesar kedua sesudah Comment.
func (g *Gudang) BacaEgnpi(ctx context.Context, masterID string) ([]models.BarisEgnpiWarisan, error) {
	baris, err := g.bacaTab(ctx, "T_TREATY_EGNPI",
		"AMOUNT, AMOUNTIDR, ASDATE, CLASSOFBUSINESS, CURRENCY, NOTE, PROPORTION, TREATYGROUP", masterID, 8)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisEgnpiWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisEgnpiWarisan{
			Jumlah: b[0], JumlahIDR: b[1], PerTanggal: b[2], KelasBisnis: b[3],
			MataUang: b[4], Keterangan: b[5], Proporsi: b[6], KelompokTreaty: b[7],
		})
	}
	return out, nil
}

// BacaRetensi membaca tab Maximum Retention.
func (g *Gudang) BacaRetensi(ctx context.Context, masterID string) ([]models.BarisRetensiWarisan, error) {
	baris, err := g.bacaTab(ctx, "T_TREATY_RETENTION",
		"AMOUNT, CLASSOFBUSINESS, CURRENCY, NOTE, TREATYGROUP", masterID, 5)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisRetensiWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisRetensiWarisan{
			Jumlah: b[0], KelasBisnis: b[1], MataUang: b[2], Keterangan: b[3], KelompokTreaty: b[4],
		})
	}
	return out, nil
}

// BacaAngsuran membaca JADWAL tab Installment - tabel anak, bukan induknya.
//
// ⛔ Induk (`T_TREATY_INSTALLMENT`) memuat total; yang dibaca orang di
// layar adalah jadwalnya. Urut `URUTAN` menjaga angsuran ke-1 sebelum ke-2.
func (g *Gudang) BacaAngsuran(ctx context.Context, masterID string) ([]models.BarisAngsuranWarisan, error) {
	baris, err := g.bacaTab(ctx, "T_TREATY_INSTALLMENT_ITEM",
		"INSTALLMENT, CURRENCY, AMOUNT, INSTALLMENTPCT, DUEDATE, PAYMENTDATE, WPC", masterID, 7)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisAngsuranWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisAngsuranWarisan{
			Angsuran: b[0], MataUang: b[1], Jumlah: b[2], Persen: b[3],
			JatuhTempo: b[4], TanggalBayar: b[5], WPC: b[6],
		})
	}
	return out, nil
}

// BacaCatatan membaca tab Information & Submit.
func (g *Gudang) BacaCatatan(ctx context.Context, masterID string) ([]models.BarisCatatanWarisan, error) {
	baris, err := g.bacaTab(ctx, "T_VIEW_COMMENT",
		"TANGGAL, OPERATORNAME, ISAPPROVED, SUGGEST", masterID, 4)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisCatatanWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisCatatanWarisan{
			Tanggal: b[0], Operator: b[1], Disetujui: b[2], Catatan: b[3],
		})
	}
	return out, nil
}

// BacaSkalaKoasuransi membaca tab Co-Ins Scale - tabel pendaratan
// kesembilan, migrasi 432.
//
// ⚠️ `COINSHARE` BUKAN angka: nilainya pita seperti `>=30% up to < 50%`
// (terpanjang 17 aksara). Ia dibaca dan ditampilkan apa adanya; pemformat
// angka akan merusaknya, dan `formatNumber` memang mengembalikan teks
// bukan-angka apa adanya justru untuk kasus ini.
//
// Terisi di 186 dari 1.854 kontrak, 702 baris. Kosong di 1.668 sisanya
// adalah keadaan yang WAJAR.
func (g *Gudang) BacaSkalaKoasuransi(ctx context.Context, masterID string) ([]models.BarisSkalaKoasuransiWarisan, error) {
	baris, err := g.bacaTab(ctx, "M_TREATYIN_COINSCALE",
		"COINSHARE, PCTLIMIT, PXCREATEOPNAME, PXCREATEDATETIME", masterID, 4)
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisSkalaKoasuransiWarisan, 0, len(baris))
	for _, b := range baris {
		out = append(out, models.BarisSkalaKoasuransiWarisan{
			BagianKoasuransi: b[0], PersenLimit: b[1], Penyusun: b[2], DisusunPada: b[3],
		})
	}
	return out, nil
}
