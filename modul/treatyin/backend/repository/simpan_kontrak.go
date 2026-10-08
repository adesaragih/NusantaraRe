package repository

// Jalur SIMPAN tombol `Save`, `Submit`, `Actions`, dan `Decline offer`.
//
// ---------------------------------------------------------------------
// ⭐ KEPUTUSAN PEMILIK PROSES — 6 dan 7 Oktober 2026
// ---------------------------------------------------------------------
//
//	"submit dan save di sistem sekarang itu akan menyimpan data kedalam
//	 masing masing table yang ditentukan, jadi tidak mengikuti pega lagi"
//	"sasarannya tetap TREATYEXCHANGEYEARLY" — tambah & ubah saja
//	"saat save boleh menyentuh table treatyin"
//
// Tiga sasaran, SATU transaksi:
//
//	T_TREATY_*            seluruh dokumen `TreatyIn` lewat `MuatKontrak` —
//	                      peta yang sama dengan pemuat, jadi tiap tab mendarat
//	                      di tabelnya sendiri
//	TREATY_IN             kepala kontrak — kesembilan belas parameter
//	                      prosedur `PEGA_TREATY_IN` (RDB `SaveTreatyIn`)
//	TREATYEXCHANGEYEARLY  grid Rate of Exchange — kolom prosedur
//	                      `PEGA_TREATYEXCHANGE`; ⛔ nol `DELETE`: tabelnya
//	                      berkunci TAHUN dan dibaca modul lain
//
// ⛔ `M_TREATY_IN` dan JSON TIDAK disentuh — prosedur Pega menulis keduanya,
// jalur ini tidak. ⛔ `M_TREATY_IN_DETAIL` (`SaveTreatyInDetail_Act` saat
// Resolve Complete) tidak termasuk tabel yang pemilik proses tentukan.

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// Pengenal BARU — bentuk prosedur Pega: `M_SITE_DATABASE.ID (CURRENT_SITE=
// '1') || LPAD(n, …, '0')`.
//
// ⛔ Sequence kontrak milik prosedur Pega TIDAK dipakai: namanya berawalan
// tabel dokumen yang pemilik proses tutup total, dan penjaganya
// (`TestAplikasiTidakMenyebutMTreatyIn`) menolak nama itu. Kontrak baru
// mendapat angka tertinggi berawalan situs + 1 (lihat `idKontrakBaru`).
// Kurs tetap memakai sequence prosedurnya sendiri.
const (
	seqKursTahunan = "TREATYEXCHANGE_SEQ"
	tabelSitus     = "M_SITE_DATABASE"
	// Panjang bagian nomor di belakang kode situs — `LPAD(…, 6, '0')`.
	panjangNomorKontrak = 6
)

// KolomKepalaTreatyIn - kolom `TREATY_IN` ↔ properti `TreatyIn`, URUT
// parameter `PEGA_TREATY_IN` (RDB `SaveTreatyIn`). `ID` tidak termasuk.
//
// ⚠️ `ClassofBusiness`, `NusareSharePct`, `BrokeragePct` — ejaan PERSIS RDB
// itu (huruf kecil di tengahnya), bukan ejaan yang mungkin diharapkan.
var KolomKepalaTreatyIn = [][2]string{
	{"PROPORTIONTYPE", "ProportionType"},
	{"TREATYCONTRACTNAME", "TreatyContractName"},
	{"TERITORIALSCOPE", "TeritorialScope"},
	{"COMMENCEMENT", "Commencement"},
	{"TERMINATION", "Termination"},
	{"CLASSOFBUSINESS", "ClassofBusiness"},
	{"LEADINGREINSSOURCE", "LeadingReinsSource"},
	{"LEADINGREINSSOURCEID", "LeadingReinsSourceID"},
	{"CEDING", "Ceding"},
	{"CEDINGID", "CedingID"},
	{"LEADINGREINSID", "LeadingReinsID"},
	{"NUSARESHAREPCT", "NusareSharePct"},
	{"BROKERAGEPCT", "BrokeragePct"},
	{"INFORMATION", "Information"},
	{"POSITIONUSERNAME", "PositionUsername"},
	{"POSITION", "Position"},
	{"STATUSAKSEPTASI", "StatusAkseptasi"},
	{"CHOOSESTATUSAKSEPTASI", "ChooseStatusAkseptasi"},
	{"TREATYYEAR", "TreatyYear"},
}

// BacaKepalaTreatyIn - kepala satu kontrak di `TREATY_IN`, berkunci
// properti Pega. Kolom `NULL` tidak dimasukkan. `ada` = barisnya ada.
func (g *Gudang) BacaKepalaTreatyIn(ctx context.Context, id string) (map[string]any, bool, error) {
	nama, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return nil, false, err
	}
	kolom := make([]string, len(KolomKepalaTreatyIn))
	for i, k := range KolomKepalaTreatyIn {
		kolom[i] = k[0]
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE ID = :1", strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, false, err
	}
	sel := make([]sql.NullString, len(kolom))
	tuju := make([]any, len(kolom))
	for i := range sel {
		tuju[i] = &sel[i]
	}
	if err := g.db.QueryRowContext(ctx, q, id).Scan(tuju...); err == sql.ErrNoRows {
		return map[string]any{}, false, nil
	} else if err != nil {
		return nil, false, fmt.Errorf("repository: membaca kepala %s %s: %w", TabelWarisanKontrak, id, err)
	}
	out := map[string]any{}
	for i, k := range KolomKepalaTreatyIn {
		if sel[i].Valid {
			out[k[1]] = sel[i].String
		}
	}
	return out, true, nil
}

// BacaDokumenPendaratan menyusun ulang dokumen `TreatyIn` satu kontrak dari
// tabel pendaratan — KEBALIKAN `MuatKontrak`, dikemudikan peta yang SAMA.
//
// ⭐ Inilah "clipboard" yang Save timpai: tab yang tidak dibuka pemakai tidak
// dikirim layar, dan tanpa dokumen ini `MuatKontrak` (kosongkan lalu sisip)
// akan menghapusnya.
//
// ⛔ `NULL` → kunci TIDAK dimasukkan; kunci bertitik (`ValueDifference.X`)
// dikembalikan ke halaman tertanamnya. Larik kosong tidak dimasukkan.
func (g *Gudang) BacaDokumenPendaratan(ctx context.Context, masterID string) (map[string]any, error) {
	doc := map[string]any{}
	// Elemen tiap baris per tabel, berkunci pengenal baris — anak menempel
	// pada induknya lewat `IDINDUK`.
	elemen := map[string]map[int64]map[string]any{}
	for _, p := range PetaPendaratan {
		baris, err := g.bacaBarisDokumen(ctx, p, masterID)
		if err != nil {
			return nil, err
		}
		elemen[p.Tabel] = map[int64]map[string]any{}
		for _, b := range baris {
			if p.Akar {
				for k, v := range b.nilai {
					setelJalur(doc, k, v)
				}
				continue
			}
			el := map[string]any{}
			for k, v := range b.nilai {
				setelJalur(el, k, v)
			}
			elemen[p.Tabel][b.id] = el
			switch {
			case p.Induk == "" && len(p.LarikGabung) > 0:
				doc[b.jenis] = tambahLarik(doc[b.jenis], el)
			case p.Induk == "":
				doc[p.Larik] = tambahLarik(doc[p.Larik], el)
			default:
				induk, ada := elemen[p.Induk][b.induk]
				if !ada {
					// Anak yatim — induknya tidak ada; tidak dikarang.
					continue
				}
				kunci := p.KunciAnak
				if len(p.LarikGabung) > 0 {
					kunci = b.jenis
				}
				induk[kunci] = tambahLarik(induk[kunci], el)
			}
		}
	}
	return doc, nil
}

type barisDokumen struct {
	id, induk int64
	jenis     string
	nilai     map[string]any
}

func (g *Gudang) bacaBarisDokumen(ctx context.Context, asli Pendaratan, masterID string) ([]barisDokumen, error) {
	// Hanya kolom yang SUDAH terpasang; tabel yang belum ada = nol baris.
	p, ada, err := g.petaTerpasang(ctx, asli)
	if err != nil || !ada {
		return nil, err
	}
	nama, err := g.db.Qualify(p.Tabel)
	if err != nil {
		return nil, err
	}
	pilih := []string{"ID"}
	if p.Induk != "" {
		pilih = append(pilih, "IDINDUK")
	} else {
		pilih = append(pilih, "0")
	}
	if len(p.LarikGabung) > 0 {
		pilih = append(pilih, kolomJenis)
	} else {
		pilih = append(pilih, "NULL")
	}
	pilih = append(pilih, p.Kolom...)
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1 ORDER BY URUTAN", strings.Join(pilih, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", p.Tabel, masterID, err)
	}
	defer func() { _ = rows.Close() }()
	var out []barisDokumen
	for rows.Next() {
		var id, induk int64
		var jenis sql.NullString
		sel := make([]sql.NullString, len(p.Kolom))
		tuju := []any{&id, &induk, &jenis}
		for i := range sel {
			tuju = append(tuju, &sel[i])
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", p.Tabel, err)
		}
		b := barisDokumen{id: id, induk: induk, jenis: jenis.String, nilai: map[string]any{}}
		for i, k := range p.Kunci {
			if i < len(sel) && sel[i].Valid {
				b.nilai[k] = sel[i].String
			}
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// setelJalur menulis kunci — bertitik berarti halaman tertanam.
func setelJalur(m map[string]any, kunci string, v any) {
	if !strings.Contains(kunci, ".") {
		m[kunci] = v
		return
	}
	bagian := strings.SplitN(kunci, ".", 2)
	anak, ok := m[bagian[0]].(map[string]any)
	if !ok {
		anak = map[string]any{}
		m[bagian[0]] = anak
	}
	setelJalur(anak, bagian[1], v)
}

func tambahLarik(lama any, el map[string]any) []any {
	l, _ := lama.([]any)
	return append(l, el)
}

// SimpanKontrak menulis seluruh rencana dalam SATU transaksi dan
// mengembalikan pengenal kontraknya.
func (g *Gudang) SimpanKontrak(ctx context.Context, r models.RencanaSimpan) (id string, err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if id, err = g.simpanDalam(ctx, tx, r); err != nil {
		return "", err
	}
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("repository: mengikat simpanan kontrak %s: %w", id, err)
	}
	return id, nil
}

// SimpanLaluBatalkanUntukUji - seluruh tulisan Save di dalam transaksi yang
// SELALU dibatalkan; mengembalikan pengenal dan cacah baris pendaratan yang
// sempat tersisip. Hanya untuk uji `db` — membuktikan SQL tulisnya sah di
// Oracle tanpa mengikat satu baris pun.
func (g *Gudang) SimpanLaluBatalkanUntukUji(ctx context.Context, r models.RencanaSimpan) (string, map[string]int, error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return "", nil, err
	}
	defer func() { _ = tx.Rollback() }()
	id, err := g.simpanDalam(ctx, tx, r)
	if err != nil {
		return "", nil, err
	}
	cacah, err := g.CacahBarisKontrak(ctx, tx, id)
	return id, cacah, err
}

func (g *Gudang) simpanDalam(ctx context.Context, tx *db.Tx, r models.RencanaSimpan) (string, error) {
	id := strings.TrimSpace(r.ID)
	baru := id == ""
	if baru {
		var err error
		if id, err = g.idKontrakBaru(ctx, tx); err != nil {
			return "", err
		}
	}
	doc := r.Dokumen
	if doc == nil {
		doc = map[string]any{}
	}
	doc["ID"] = id
	// ⛔ `MuatKontrakSebagian`, BUKAN `MuatKontrak` — dan perbedaannya
	// adalah perbaikan cacat kehilangan data 7 Oktober 2026.
	//
	// `MuatKontrak` mengosongkan ke-31 tabel pendaratan kontrak ini lebih
	// dulu, lalu menyisipkan ulang dari dokumen. Itu benar bagi PEMUAT
	// MASSAL, yang dokumennya selalu kontrak utuh. Dokumen tombol Save
	// TIDAK utuh: layar hanya mengirim properti tab yang pernah dibuka
	// pemakai, sebab tab dirender hanya ketika aktif.
	//
	// Akibatnya satu penekanan Save menghapus baris SELURUH tab yang tidak
	// dibuka, dan layar lalu memuat ulang keadaan yang sudah terlanjur
	// kosong itu — persis laporan pemakai: *"saat melakukan inputan tiba
	// tiba datanya hilang dan semua di reset malah tidak tersimpan"*.
	if _, err := g.MuatKontrakSebagian(ctx, tx, id, doc); err != nil {
		return "", err
	}
	if err := g.tulisKepalaTreatyIn(ctx, tx, id, doc, baru); err != nil {
		return "", err
	}
	if r.Kurs != nil {
		if err := g.tulisKurs(ctx, tx, *r.Kurs, r.Operator, r.Stempel); err != nil {
			return "", err
		}
	}
	return id, nil
}

// idKontrakBaru - situs || nomor 6 digit: angka TERTINGGI berawalan situs di
// `TREATY_IN` dan `T_TREATY_REVISION`, ditambah satu.
//
// ⛔ `TREATY_IN` tidak berkunci unik, jadi dua Save serentak yang membaca
// angka tertinggi yang sama akan membuat dua kontrak berpengenal kembar
// tanpa galat. Baris situs dikunci `FOR UPDATE` lebih dulu — Save kedua
// menunggu yang pertama mengikat, lalu membaca angka yang sudah naik.
func (g *Gudang) idKontrakBaru(ctx context.Context, tx *db.Tx) (string, error) {
	situs, err := g.kodeSitusTerkunci(ctx, tx)
	if err != nil {
		return "", err
	}
	tabel, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return "", err
	}
	rev, err := g.db.Qualify("T_TREATY_REVISION")
	if err != nil {
		return "", err
	}
	pola := fmt.Sprintf("^%s[0-9]{%d}$", situs, panjangNomorKontrak)
	q := fmt.Sprintf(`SELECT NVL(MAX(N), 0) FROM (
		SELECT TO_NUMBER(SUBSTR(ID, :1)) N FROM %s WHERE REGEXP_LIKE(ID, :2)
		UNION ALL
		SELECT TO_NUMBER(SUBSTR(MASTERID, :3)) N FROM %s WHERE REGEXP_LIKE(MASTERID, :4))`, tabel, rev)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	mulai := len(situs) + 1
	var tertinggi int64
	if err := tx.QueryRowContext(ctx, q, mulai, pola, mulai, pola).Scan(&tertinggi); err != nil {
		return "", fmt.Errorf("repository: membaca pengenal kontrak tertinggi: %w", err)
	}
	return situs + kiriNol(fmt.Sprint(tertinggi+1), panjangNomorKontrak), nil
}

// IDKontrakBaruUntukUji - pengenal yang Save BERIKUTNYA akan pakai, dihitung
// di transaksi yang SELALU dibatalkan. Hanya untuk uji `db`.
func (g *Gudang) IDKontrakBaruUntukUji(ctx context.Context) (string, error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()
	return g.idKontrakBaru(ctx, tx)
}

// kodeSitusTerkunci - kode situs, dengan barisnya dikunci sampai transaksi
// selesai (penjaga pengenal kontrak kembar).
func (g *Gudang) kodeSitusTerkunci(ctx context.Context, tx *db.Tx) (string, error) {
	return g.bacaSitus(ctx, tx, " FOR UPDATE")
}

func (g *Gudang) kodeSitus(ctx context.Context, tx *db.Tx) (string, error) {
	return g.bacaSitus(ctx, tx, "")
}

func (g *Gudang) bacaSitus(ctx context.Context, tx *db.Tx, kunci string) (string, error) {
	nama, err := g.db.Qualify(tabelSitus)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf("SELECT ID FROM %s WHERE CURRENT_SITE = '1'%s", nama, kunci)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var s string
	if err := tx.QueryRowContext(ctx, q).Scan(&s); err != nil {
		return "", fmt.Errorf("repository: membaca %s: %w", tabelSitus, err)
	}
	return strings.TrimSpace(s), nil
}

func kiriNol(s string, n int) string {
	s = strings.TrimSpace(s)
	for len(s) < n {
		s = "0" + s
	}
	return s
}

// tulisKepalaTreatyIn - `INSERT` (kontrak baru) atau `UPDATE` SELURUH
// kesembilan belas kolom, seperti `PEGA_TREATY_IN`. Properti yang tidak ada
// di dokumen ditulis `NULL` — dokumen yang disusun services SUDAH memuat
// kepala tersimpan, jadi yang `NULL` memang kosong.
func (g *Gudang) tulisKepalaTreatyIn(ctx context.Context, tx *db.Tx, id string, doc map[string]any, baru bool) error {
	nama, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return err
	}
	args := []any{}
	nilai := func(k string) any {
		v, ada := NilaiTeks(nilaiJalur(doc, k))
		if !ada {
			return nil
		}
		return v
	}
	var q string
	if baru {
		kolom := []string{"ID"}
		tanda := []string{":1"}
		args = append(args, id)
		for _, k := range KolomKepalaTreatyIn {
			kolom = append(kolom, k[0])
			args = append(args, nilai(k[1]))
			tanda = append(tanda, fmt.Sprintf(":%d", len(args)))
		}
		q = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", nama, strings.Join(kolom, ", "), strings.Join(tanda, ", "))
	} else {
		set := make([]string, 0, len(KolomKepalaTreatyIn))
		for _, k := range KolomKepalaTreatyIn {
			args = append(args, nilai(k[1]))
			set = append(set, fmt.Sprintf("%s = :%d", k[0], len(args)))
		}
		args = append(args, id)
		q = fmt.Sprintf("UPDATE %s SET %s WHERE ID = :%d", nama, strings.Join(set, ", "), len(args))
	}
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("repository: menulis kepala %s %s: %w", TabelWarisanKontrak, id, err)
	}
	return nil
}

// tulisKurs - grid Rate of Exchange ke `TREATYEXCHANGEYEARLY`: baris ber-ID
// DIPERBARUI, baris tanpa ID DITAMBAHKAN (ID = situs || LPAD(
// TREATYEXCHANGE_SEQ.NEXTVAL, 4, '0'), seperti `PEGA_TREATYEXCHANGE`).
//
// ⛔ Nol `DELETE` — keputusan pemilik proses: tabelnya berkunci TAHUN, satu
// baris dipakai setiap kontrak tahun itu dan dibaca modul lain.
func (g *Gudang) tulisKurs(ctx context.Context, tx *db.Tx, k models.KursSimpan, operator, stempel string) error {
	if strings.TrimSpace(k.Tahun) == "" {
		return nil
	}
	nama, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return err
	}
	var situs string
	for _, b := range k.Baris {
		if strings.TrimSpace(b.MataUangID) == "" && strings.TrimSpace(b.MataUang) == "" {
			continue // baris kosong yang baru ditambahkan — tidak ada isinya
		}
		if strings.TrimSpace(b.ID) != "" {
			// Tanggal datang sebagai stempel tengah malam (services
			// `stempelHariPega`); yang dikosongkan pemakai memang dikosongkan.
			q := fmt.Sprintf(`UPDATE %s SET IDCURRENCY = :1, CURRENCY = :2, TOIDR = :3, STARTDATE = :4,
				ENDDATE = :5, USERID = :6, DATEIU = :7 WHERE ID = :8 AND TREATYYEAR = :9`, nama)
			if err := db.PeriksaSQL(q); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, q, b.MataUangID, b.MataUang, b.NilaiKeIDR,
				b.BerlakuDari, b.BerlakuSampai, operator, stempel, b.ID, k.Tahun); err != nil {
				return fmt.Errorf("repository: memperbarui kurs %s: %w", b.ID, err)
			}
			continue
		}
		if situs == "" {
			if situs, err = g.kodeSitus(ctx, tx); err != nil {
				return err
			}
		}
		n, err := g.db.NomorBerikut(ctx, tx, seqKursTahunan)
		if err != nil {
			return err
		}
		q := fmt.Sprintf(`INSERT INTO %s (ID, TREATYYEAR, STARTDATE, ENDDATE, TOIDR, USERID, DATEIU,
			IDCURRENCY, CURRENCY, QUARTER, DATEIN) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11)`, nama)
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, q, situs+kiriNol(n, 4), k.Tahun, b.BerlakuDari, b.BerlakuSampai,
			b.NilaiKeIDR, operator, stempel, b.MataUangID, b.MataUang, "0", stempel); err != nil {
			return fmt.Errorf("repository: menambah kurs %s %s: %w", k.Tahun, b.MataUang, err)
		}
	}
	return nil
}

// PemegangPosisi - LOGIN_ID akun AKTIF yang memegang workbasket itu (menu
// Kelola User), urut nama. Pengisi `PositionUsername` jalur naik — keputusan
// pemilik proses 7 Oktober 2026: "ambil namanya dari kelola user dan
// sesuaikan posisinya".
func (g *Gudang) PemegangPosisi(ctx context.Context, workbasket string) ([]string, error) {
	lw, err := g.db.Qualify("M_LOGIN_GO_WORKBASKET")
	if err != nil {
		return nil, err
	}
	akun, err := g.db.Qualify("M_LOGIN_GO")
	if err != nil {
		return nil, err
	}
	wb, err := g.db.Qualify("M_WORKBASKET")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT l.LOGIN_ID FROM %s l
		JOIN %s a ON a.LOGIN_ID = l.LOGIN_ID
		JOIN %s w ON w.WORKBASKET_ID = l.WORKBASKET_ID
		WHERE l.WORKBASKET_ID = :1 AND a.IS_ACTIVE = '1' AND w.IS_ACTIVE = '1'`, lw, akun, wb)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, workbasket)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca pemegang %s: %w", workbasket, err)
	}
	defer func() { _ = rows.Close() }()
	out := []string{}
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("repository: membaca pemegang %s: %w", workbasket, err)
		}
		out = append(out, s)
	}
	sort.Strings(out)
	return out, rows.Err()
}
