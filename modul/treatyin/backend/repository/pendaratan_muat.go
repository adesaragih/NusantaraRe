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
const TabelDokumenPendaratan = TabelWarisanJSON

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
// Anak sebelum induknya. `FK_MTI_INSTALLMENTITEM_1` memang berkaskade, jadi
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
func (g *Gudang) MuatKontrak(ctx context.Context, tx *db.Tx, masterID string, doc map[string]any) (map[string]int, error) {
	if strings.TrimSpace(masterID) == "" {
		return nil, fmt.Errorf("repository: memuat pendaratan tanpa pengenal kontrak")
	}
	if _, err := g.KosongkanKontrak(ctx, tx, masterID); err != nil {
		return nil, err
	}

	disisip := map[string]int{}
	// ⚠️ Pengenal induk angsuran disimpan di sini: butirnya merujuknya lewat
	// `IDINDUK`, dan `IDINDUK` harus pengenal baris yang BARU SAJA lahir -
	// bukan nomor urut, bukan pengenal dari muatan sebelumnya.
	var idInduk []string

	for i, p := range PetaPendaratan {
		if i == IndeksButirAngsuran {
			n, err := g.sisipButirAngsuran(ctx, tx, masterID, doc, idInduk)
			if err != nil {
				return nil, err
			}
			disisip[p.Tabel] = n
			continue
		}
		baris := BarisLarik(doc, p.Larik)
		for urut, el := range baris {
			var paksa *string
			if i == IndeksAngsuran {
				// Induk angsuran: pengenalnya diambil LEBIH DULU, sebab
				// butirnya harus menunjuknya.
				id, err := g.db.NomorBerikut(ctx, tx, p.Seq)
				if err != nil {
					return nil, err
				}
				idInduk = append(idInduk, id)
				paksa = &id
			}
			if err := g.sisipSatu(ctx, tx, p, masterID, urut, el, paksa, nil); err != nil {
				return nil, err
			}
		}
		disisip[p.Tabel] = len(baris)
	}
	return disisip, nil
}

// sisipButirAngsuran mendaratkan `Installment[].InstallmentList`, tiap butir
// menunjuk induk yang melahirkannya.
func (g *Gudang) sisipButirAngsuran(ctx context.Context, tx *db.Tx, masterID string, doc map[string]any, idInduk []string) (int, error) {
	p := PetaPendaratan[IndeksButirAngsuran]
	n := 0
	for ke, induk := range BarisLarik(doc, PetaPendaratan[IndeksAngsuran].Larik) {
		if ke >= len(idInduk) {
			// Tidak mungkin tercapai selama induknya disisipkan lebih dulu;
			// bila tercapai, yang rusak urutan pemuatan - dan itu harus
			// berhenti, bukan menyisipkan butir tanpa induk.
			return 0, fmt.Errorf("repository: butir angsuran kontrak %s menunjuk induk ke-%d yang belum lahir", masterID, ke)
		}
		larik, ok := induk[LarikAnakAngsuran].([]any)
		if !ok {
			continue
		}
		urut := 0
		for _, el := range larik {
			objek, ok := el.(map[string]any)
			if !ok {
				continue
			}
			ind := idInduk[ke]
			if err := g.sisipSatu(ctx, tx, p, masterID, urut, objek, nil, &ind); err != nil {
				return 0, err
			}
			urut++
			n++
		}
	}
	return n, nil
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
func (g *Gudang) sisipSatu(ctx context.Context, tx *db.Tx, p Pendaratan, masterID string, urut int, el map[string]any, idPaksa, idInduk *string) error {
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
	for j, k := range p.Kunci {
		teks, ada := NilaiTeks(el[k])
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
