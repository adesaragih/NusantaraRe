package repository

// Pemuat tabel pendaratan - SATU-SATUNYA jalur tulis ke kedelapan tabel.
//
// ⛔ Nol tulisan ke tabel warisan. `M_TREATY_IN` DIBACA di sini dan tidak
// pernah disentuh; `TREATY_IN`, `M_TREATY_IN2`, `TREATYEXCHANGEYEARLY`, dan
// `M_TREATY_IN_DETAIL` tidak disebut sama sekali. Dijaga
// `TestWarisanHanyaDibaca`.
//
// TIGA SIFAT yang ronde ini tuntut, dan di mana masing-masing diwujudkan:
//
//	idempoten   `KosongkanKontrak` berjalan SEBELUM setiap sisip, di dalam
//	            transaksi yang sama. Memuat dua kali menghasilkan keadaan
//	            yang sama persis. UNIQUE (MASTERID, URUTAN) di migrasi 430
//	            adalah jaring keduanya: kalau pengosongan terlewat, sisipnya
//	            MELEDAK alih-alih menggandakan.
//	terbalikkan `KosongkanKontrak` sendiri, dan berkas
//	            `alat/kosongkan-tab-treatyin.sql` untuk seluruh tabel.
//	            Keduanya `DELETE`, bukan `TRUNCATE` - yang kedua ber-DDL dan
//	            MENGIKAT transaksinya, sehingga `ROLLBACK` tidak lagi
//	            mengembalikan apa pun.
//	tercocokkan `CacahBarisKontrak` mengadu cacah baris dengan `CacahLarik`
//	            atas dokumen yang sama.
//
// ⛔ NOL `Commit` di berkas ini. Transaksinya milik pemanggil, dan di dalam
// uji pemanggil itu me-`Rollback`.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// TabelDokumenPendaratan - dokumen sumbernya. Dibaca, tidak pernah ditulis.
//
// ⛔⛔ SATU-SATUNYA tempat `M_TREATY_IN` boleh disebut, dan tetapannya PINDAH
// ke sini 6 Oktober 2026 justru untuk itu.
//
// Pemilik proses melarang keras APLIKASI memakai `M_TREATY_IN`. Pemuat BUKAN
// aplikasi: ia alat sekali-jalan yang memindahkan isi dokumen ke tabel
// pendaratan, dan sesudahnya layar membaca tabel itu. Larangannya berlaku
// pada jalur baca — `TestAplikasiTidakMenyebutMTreatyIn` menegakkannya.
const TabelDokumenPendaratan = "M_TREATY_IN"

// BacaDokumenMentah membaca `JSONDATA` satu kontrak sebagai teks.
//
// ⚠️ Berbeda dari `BacaKontrakWarisan`, yang mengurai LIMA kunci ke dalam
// struct. Pemuat perlu dokumen UTUH: kunci yang struct itu tidak punya
// medannya justru yang hendak didaratkan.
func (g *Gudang) BacaDokumenMentah(ctx context.Context, masterID string) (string, error) {
	nama, err := g.db.Qualify(TabelDokumenPendaratan)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf("SELECT JSONDATA FROM %s WHERE ID = :1", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var teks string
	if err := g.db.QueryRowContext(ctx, q, masterID).Scan(&teks); err != nil {
		return "", fmt.Errorf("repository: membaca dokumen %s kontrak %s: %w",
			TabelDokumenPendaratan, masterID, err)
	}
	return teks, nil
}

// DaftarMasterID membaca pengenal kontrak yang PUNYA dokumen, urut naik.
//
// Urut naik, bukan menurun: pemuatan dijalankan bertahap, dan urutan yang
// tetap membuat "sudah sampai mana" dapat dibaca dari pengenal terakhir.
func (g *Gudang) DaftarMasterID(ctx context.Context, batas int) ([]string, error) {
	nama, err := g.db.Qualify(TabelDokumenPendaratan)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID FROM %s
		ORDER BY TO_NUMBER(ID DEFAULT 0 ON CONVERSION ERROR), ID`, nama)
	args := []any{}
	if batas > 0 {
		q += " FETCH FIRST :1 ROWS ONLY"
		args = append(args, batas)
	}
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: mendaftar pengenal %s: %w", TabelDokumenPendaratan, err)
	}
	defer func() { _ = baris.Close() }()

	var out []string
	for baris.Next() {
		var id string
		if err := baris.Scan(&id); err != nil {
			return nil, fmt.Errorf("repository: membaca pengenal %s: %w", TabelDokumenPendaratan, err)
		}
		out = append(out, id)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca pengenal %s: %w", TabelDokumenPendaratan, err)
	}
	return out, nil
}

// KosongkanKontrak membuang seluruh baris pendaratan milik SATU kontrak.
//
// Anak sebelum induknya. `FK_TT_INSTALLMENT_ITEM_1` memang berkaskade, jadi
// baris pertama secara teknis mubazir - ia tetap dijalankan supaya cacah
// yang dikembalikan jujur per tabel, dan supaya fungsi ini tetap benar bila
// kaskade itu suatu hari dicabut.
func (g *Gudang) KosongkanKontrak(ctx context.Context, tx *db.Tx, masterID string) (map[string]int64, error) {
	dibuang := map[string]int64{}
	for i := len(PetaPendaratan) - 1; i >= 0; i-- {
		p := PetaPendaratan[i]
		nama, err := g.db.Qualify(p.Tabel)
		if err != nil {
			return nil, err
		}
		q := fmt.Sprintf("DELETE FROM %s WHERE MASTERID = :1", nama)
		if err := db.PeriksaSQL(q); err != nil {
			return nil, err
		}
		hasil, err := tx.ExecContext(ctx, q, masterID)
		if err != nil {
			return nil, fmt.Errorf("repository: mengosongkan %s kontrak %s: %w", p.Tabel, masterID, err)
		}
		n, _ := hasil.RowsAffected()
		dibuang[p.Tabel] = n
	}
	return dibuang, nil
}

// MuatKontrak mendaratkan seluruh larik satu kontrak.
//
// Memanggil `KosongkanKontrak` lebih dulu - itulah yang membuatnya
// idempoten. Mengembalikan cacah baris yang DISISIPKAN per tabel.
//
// ⭐ POHONNYA DITAPAKI SECARA UMUM sejak migrasi 437. Bentuk sebelumnya
// menanam satu pasangan induk-anak di dalam gelungnya lewat dua indeks
// tetap; dengan tiga tingkat dan tiga belas tabel baru, cara itu menuntut
// satu cabang kode per pasangan. Hubungan induk-anak kini DATA (`Induk`
// dan `KunciAnak` di `PetaPendaratan`), dan gelung di bawah tidak tahu
// tabel mana yang sedang dimuatnya.
//
// ⛔ Urutan `PetaPendaratan` tetap mengikat: induk WAJIB mendahului
// anaknya, sebab anak menunjuk pengenal baris yang baru saja lahir.
func (g *Gudang) MuatKontrak(ctx context.Context, tx *db.Tx, masterID string, doc map[string]any) (map[string]int, error) {
	if strings.TrimSpace(masterID) == "" {
		return nil, fmt.Errorf("repository: memuat pendaratan tanpa pengenal kontrak")
	}
	if _, err := g.KosongkanKontrak(ctx, tx, masterID); err != nil {
		return nil, err
	}

	// Tabel yang PUNYA anak perlu pengenal barisnya diambil lebih dulu;
	// daun tidak. Memaksa seluruhnya berarti satu perjalanan pulang-pergi
	// tambahan untuk puluhan ribu baris yang pengenalnya tidak dilihat
	// siapa pun.
	punyaAnak := map[string]bool{}
	for _, p := range PetaPendaratan {
		if p.Induk != "" {
			punyaAnak[p.Induk] = true
		}
	}

	// Pengenal dan elemen tiap baris yang lahir, per tabel, SEJAJAR urutan.
	// Anak membacanya untuk menemukan induknya.
	lahir := map[string][]string{}
	elemen := map[string][]map[string]any{}

	disisip := map[string]int{}
	for _, p := range PetaPendaratan {
		// Satu kelompok = satu induk beserta anak-anaknya. Untuk tabel
		// tingkat pertama kelompoknya satu, berinduk nil.
		type kelompok struct {
			induk *string
			baris []map[string]any
			// ⭐ Sejajar dengan `baris`: nama larik asal tiap baris, yang
			// menjadi isi kolom `JENIS`. Kosong bila tabelnya tidak
			// menggabungkan beberapa larik.
			jenis []string
		}
		var kel []kelompok

		switch {
		case p.Akar:
			// ⭐ Dokumennya SENDIRI yang menjadi satu-satunya elemen.
			kel = append(kel, kelompok{nil, []map[string]any{doc}, nil})
		case p.Induk == "" && len(p.LarikGabung) > 0:
			var baris []map[string]any
			var jenis []string
			for _, nama := range p.LarikGabung {
				for _, el := range BarisLarik(doc, nama) {
					baris = append(baris, el)
					jenis = append(jenis, nama)
				}
			}
			kel = append(kel, kelompok{nil, baris, jenis})
		case p.Induk == "":
			kel = append(kel, kelompok{nil, BarisLarik(doc, p.Larik), nil})
		case len(p.LarikGabung) > 0:
			indukEl, indukID := elemen[p.Induk], lahir[p.Induk]
			for ke, el := range indukEl {
				var baris []map[string]any
				var jenis []string
				for _, nama := range p.LarikGabung {
					larik, ok := el[nama].([]any)
					if !ok {
						continue
					}
					for _, e := range larik {
						if o, ok := e.(map[string]any); ok {
							baris = append(baris, o)
							jenis = append(jenis, nama)
						}
					}
				}
				if len(baris) == 0 {
					continue
				}
				if ke >= len(indukID) {
					return nil, fmt.Errorf("repository: %s kontrak %s menunjuk induk ke-%d yang belum lahir", p.Tabel, masterID, ke)
				}
				id := indukID[ke]
				kel = append(kel, kelompok{&id, baris, jenis})
			}
		default:
			indukEl, indukID := elemen[p.Induk], lahir[p.Induk]
			for ke, el := range indukEl {
				larik, ok := el[p.KunciAnak].([]any)
				if !ok {
					continue
				}
				var baris []map[string]any
				for _, e := range larik {
					if o, ok := e.(map[string]any); ok {
						baris = append(baris, o)
					}
				}
				if len(baris) == 0 {
					continue
				}
				if ke >= len(indukID) {
					// Tidak tercapai selama induknya mendahului di peta;
					// bila tercapai, yang rusak urutan pemuatan - dan itu
					// harus berhenti, bukan menyisipkan anak tanpa induk.
					return nil, fmt.Errorf("repository: %s kontrak %s menunjuk induk ke-%d yang belum lahir", p.Tabel, masterID, ke)
				}
				id := indukID[ke]
				kel = append(kel, kelompok{&id, baris, nil})
			}
		}

		n := 0
		for _, k := range kel {
			// ⛔ `URUTAN` dihitung PER INDUK, bukan menyeluruh: kunci
			// uniknya `(IDINDUK, URUTAN)`, dan penomoran menyeluruh
			// membuat anak induk kedua bertabrakan dengan yang pertama.
			for urut, el := range k.baris {
				var paksa *string
				if punyaAnak[p.Tabel] {
					id, err := g.db.NomorBerikut(ctx, tx, p.Seq)
					if err != nil {
						return nil, err
					}
					paksa = &id
					lahir[p.Tabel] = append(lahir[p.Tabel], id)
				}
				jenis := ""
				if urut < len(k.jenis) {
					jenis = k.jenis[urut]
				}
				if err := g.sisipSatu(ctx, tx, p, masterID, urut, el, paksa, k.induk, jenis); err != nil {
					return nil, err
				}
				elemen[p.Tabel] = append(elemen[p.Tabel], el)
				n++
			}
		}
		disisip[p.Tabel] = n
	}
	return disisip, nil
}

// sisipSatu menyisipkan satu baris.
//
// ⛔ DUA cara mengambil pengenal, dan bedanya disengaja:
//
//   - `idPaksa == nil` — `SEQUENCE.NEXTVAL` ditulis LANGSUNG di dalam
//     pernyataan sisip. Satu perjalanan pulang-pergi per baris, dan untuk
//     25.740 baris yang pengenalnya tidak pernah dilihat siapa pun itu
//     seluruh yang diperlukan.
//   - `idPaksa != nil` — pengenalnya sudah diambil pemanggil lewat
//     `NomorBerikut`, sebab butir angsuran harus menunjuknya lewat
//     `IDINDUK`. Hanya 796 baris yang melewati jalur ini.
//
// INV-02 terpenuhi di kedua jalur: nomornya selalu dari sequence, tidak
// pernah dari cap waktu, nomor urut, maupun teks yang disusun sendiri.
func (g *Gudang) sisipSatu(ctx context.Context, tx *db.Tx, p Pendaratan, masterID string, urut int, el map[string]any, idPaksa, idInduk *string, jenis string) error {
	nama, err := g.db.Qualify(p.Tabel)
	if err != nil {
		return err
	}

	kolom := []string{"ID", "MASTERID", "URUTAN"}
	nilai := make([]string, 0, len(p.Kunci)+4)
	args := []any{}
	if idPaksa != nil {
		nilai = append(nilai, ":1")
		args = append(args, *idPaksa)
	} else {
		seq, err := g.db.Qualify(p.Seq)
		if err != nil {
			return err
		}
		nilai = append(nilai, seq+".NEXTVAL")
	}
	nilai = append(nilai, fmt.Sprintf(":%d", len(args)+1), fmt.Sprintf(":%d", len(args)+2))
	args = append(args, masterID, urut)

	if idInduk != nil {
		kolom = append(kolom, "IDINDUK")
		nilai = append(nilai, fmt.Sprintf(":%d", len(args)+1))
		args = append(args, *idInduk)
	}
	// ⭐ `JENIS` — nama larik ASAL baris ini. Kolomnya `NOT NULL`, jadi ia
	// ditulis hanya ketika tabelnya memang menggabungkan beberapa larik.
	if jenis != "" {
		kolom = append(kolom, "JENIS")
		nilai = append(nilai, fmt.Sprintf(":%d", len(args)+1))
		args = append(args, jenis)
	}
	for j, k := range p.Kunci {
		teks, ada := NilaiTeks(nilaiJalur(el, k))
		kolom = append(kolom, p.Kolom[j])
		nilai = append(nilai, fmt.Sprintf(":%d", len(args)+1))
		if !ada {
			args = append(args, nil)
			continue
		}
		args = append(args, teks)
	}

	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		nama, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("repository: menyisipkan %s kontrak %s urutan %d: %w",
			p.Tabel, masterID, urut, err)
	}
	return nil
}

// CacahBarisKontrak menghitung baris yang SUDAH mendarat untuk satu kontrak.
//
// ⛔ `tx` boleh nil, dan bedanya penting:
//
//   - `tx != nil` — dibaca DI DALAM transaksi pemuat, sehingga baris yang
//     belum mengikat ikut terhitung. Inilah yang rekonsiliasi pakai sebelum
//     memutuskan mengikat atau membatalkan.
//   - `tx == nil` — dibaca dari luar: apa yang dilihat orang lain.
//
// Memakai yang salah membuat rekonsiliasi membuktikan hal yang keliru -
// yang pertama selalu nol bila dibaca dari luar, yang kedua tidak pernah
// melihat kerusakan yang sudah mengikat.
func (g *Gudang) CacahBarisKontrak(ctx context.Context, tx *db.Tx, masterID string) (map[string]int, error) {
	return g.cacah(ctx, tx, "WHERE MASTERID = :1", masterID)
}

// CacahBarisSeluruhnya menghitung isi kedelapan tabel - untuk rekonsiliasi
// penuh terhadap `Pendaratan.CacahTerukur`.
func (g *Gudang) CacahBarisSeluruhnya(ctx context.Context, tx *db.Tx) (map[string]int, error) {
	return g.cacah(ctx, tx, "")
}

func (g *Gudang) cacah(ctx context.Context, tx *db.Tx, saring string, args ...any) (map[string]int, error) {
	out := map[string]int{}
	for _, p := range PetaPendaratan {
		nama, err := g.db.Qualify(p.Tabel)
		if err != nil {
			return nil, err
		}
		q := strings.TrimSpace(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", nama, saring))
		if err := db.PeriksaSQL(q); err != nil {
			return nil, err
		}
		var n int
		var baris *sql.Row
		if tx != nil {
			baris = tx.QueryRowContext(ctx, q, args...)
		} else {
			baris = g.db.QueryRowContext(ctx, q, args...)
		}
		if err := baris.Scan(&n); err != nil {
			return nil, fmt.Errorf("repository: menghitung %s: %w", p.Tabel, err)
		}
		out[p.Tabel] = n
	}
	return out, nil
}

// nilaiJalur mengambil nilai sebuah kunci, termasuk kunci BERTITIK.
//
// ⭐ `ValueDifference.RNMShare` adalah medan halaman TERTANAM satu tingkat
// — persis cara Section Pega mengikatnya
// (`TreatyIn.ValueDifference.RNMShare`). Tanpa penelusuran ini kolomnya
// selalu `NULL`, dan tidak ada yang bersuara.
//
// ⚠️ Kunci TANPA titik dibaca apa adanya, jadi kunci yang namanya memang
// memuat titik tidak dapat dibedakan dari jalur. Nol kunci semacam itu ada
// di korpus hari ini; bila kelak ada, pembedanya harus dinyatakan di peta.
func nilaiJalur(el map[string]any, kunci string) any {
	if v, ada := el[kunci]; ada {
		return v
	}
	if !strings.Contains(kunci, ".") {
		return nil
	}
	kini := el
	bagian := strings.Split(kunci, ".")
	for i, b := range bagian {
		v, ada := kini[b]
		if !ada {
			return nil
		}
		if i == len(bagian)-1 {
			return v
		}
		anak, ok := v.(map[string]any)
		if !ok {
			return nil
		}
		kini = anak
	}
	return nil
}

// ===========================================================================
// DOKUMEN ADJUSTMENT — `M_TREATY_IN_EDM`
// ===========================================================================
//
// ⭐ Tabel pendaratannya SAMA. `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`
// menandai tiap kotaknya `BERSAMA -> Prop, Non Prop, EDM Prop, EDM Non Prop`:
// dokumen Adjustment berbentuk sama dengan dokumen Treaty In.
//
// ⛔ Yang membedakan hanya `MASTERID` — pengenal dokumen untuk sisi `New`,
// pengenal yang sama berakhiran `AkhiranSisiLama` untuk halaman `OLDDATA`.

// TabelDokumenPenyesuaian - dokumen Adjustment. Dibaca, tidak pernah ditulis.
const TabelDokumenPenyesuaian = "M_TREATY_IN_EDM"

// DaftarMasterIDPenyesuaian membaca pengenal penyesuaian yang punya dokumen.
//
// ⛔ Urut TEKS, bukan angka: pengenalnya berbentuk `<asal>/R<nn>`, dan urutan
// teks menaruh revisi satu kontrak bersebelahan dan berurut.
func (g *Gudang) DaftarMasterIDPenyesuaian(ctx context.Context, batas int) ([]string, error) {
	nama, err := g.db.Qualify(TabelDokumenPenyesuaian)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("SELECT ID FROM %s ORDER BY ID", nama)
	args := []any{}
	if batas > 0 {
		q += " FETCH FIRST :1 ROWS ONLY"
		args = append(args, batas)
	}
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: mendaftar pengenal %s: %w", TabelDokumenPenyesuaian, err)
	}
	defer func() { _ = baris.Close() }()
	var out []string
	for baris.Next() {
		var id string
		if err := baris.Scan(&id); err != nil {
			return nil, fmt.Errorf("repository: membaca pengenal %s: %w", TabelDokumenPenyesuaian, err)
		}
		out = append(out, id)
	}
	return out, baris.Err()
}

// BacaDokumenPenyesuaianMentah membaca `JSONDATA` satu penyesuaian.
func (g *Gudang) BacaDokumenPenyesuaianMentah(ctx context.Context, id string) (string, error) {
	nama, err := g.db.Qualify(TabelDokumenPenyesuaian)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf("SELECT JSONDATA FROM %s WHERE ID = :1", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var teks string
	if err := g.db.QueryRowContext(ctx, q, id).Scan(&teks); err != nil {
		return "", fmt.Errorf("repository: membaca dokumen %s %s: %w", TabelDokumenPenyesuaian, id, err)
	}
	return teks, nil
}

// MuatPenyesuaian mendaratkan KEDUA sisi satu dokumen Adjustment.
//
// ⛔ Sisi `Old` didaratkan meski kosong-pun tidak: halaman `OLDDATA` yang
// tidak ada berarti penyesuaian itu memang hanya punya sisi `New`, dan
// mendaratkan baris kosong untuknya akan membuat layar menampilkan
// perbandingan yang tidak pernah ada.
//
// ⚠️ `OLDDATA` DI DALAM `OLDDATA` tidak ditelusuri — Section hanya mengikat
// satu tingkat, dan tingkat kedua tidak tampil di layar mana pun.
func (g *Gudang) MuatPenyesuaian(ctx context.Context, tx *db.Tx, id string, doc map[string]any) (map[string]int, error) {
	hasil, err := g.MuatKontrak(ctx, tx, id, doc)
	if err != nil {
		return nil, err
	}
	lama, ok := doc["OLDDATA"].(map[string]any)
	if !ok {
		// Tetap kosongkan sisi `Old`, supaya pemuatan ulang sesudah halaman
		// `OLDDATA` dihapus tidak meninggalkan baris basi.
		if _, err := g.KosongkanKontrak(ctx, tx, id+AkhiranSisiLama); err != nil {
			return nil, err
		}
		return hasil, nil
	}
	sisiLama, err := g.MuatKontrak(ctx, tx, id+AkhiranSisiLama, lama)
	if err != nil {
		return nil, err
	}
	for tabel, n := range sisiLama {
		hasil[tabel] += n
	}
	return hasil, nil
}
