package repository

// Penyimpan dan pembaca pohon klaim - tiket 14.
//
// Untuk apa berkas ini: menulis satu klaim utuh - akar work object, header,
// peserta, baris adjustment, spreading, sampai spreading retro - dan membacanya
// kembali sampai tingkat terdalam.
//
// Dibaca sesudah: klaimlife.go (pembacaan tiket 01) dan migrasidata.go.
//
// Dua penulisan, satu transaksi:
// sistem baru menulis DUA tempat - tabel relasional baru DAN baris datar ke
// OS_AKSEPTASI_KLAIM_LIFE, sebab hilir masih membaca dari sana. Yang dibuang
// hanya blob JSON. Keduanya berada di dalam SATU transaksi: kegagalan pada
// salah satunya membatalkan keduanya.
//
// Istilah:
//   - transaksi : sekumpulan perubahan yang berlaku semua atau tidak sama
//                 sekali. Dibuka dan ditutup oleh services (ADR-U-0029).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/db"
	"nusantarare/inti/kontrak"
	"nusantarare/inti/uang"
	"nusantarare/inti/utils"
	"nusantarare/modul/claimlife/models"
)

// namaTabelLama adalah tabel datar warisan yang tetap ditulis.
const namaTabelLama = "OS_AKSEPTASI_KLAIM_LIFE"

// namaTabelRetro dinyatakan sekali. Namanya pernah diperpendek dari 36 menjadi
// 29 byte (batas pengenal Oracle di bawah 12.2), dan literal yang tersebar
// itulah yang membuat perubahan semacam itu mahal.
const namaTabelRetro = "T_CLAIMLF_ADJ_SPREADING_RETRO"

// PohonKlaim menulis dan membaca satu klaim utuh.
type PohonKlaim struct {
	db *db.DB
}

// NewPohonKlaim membuat penyimpan pohon klaim.
func NewPohonKlaim(db *db.DB) *PohonKlaim { return &PohonKlaim{db: db} }

// exec menjalankan satu pernyataan di dalam transaksi, sesudah memeriksanya.
func (r *PohonKlaim) exec(ctx context.Context, tx *db.Tx, q string, args ...any) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

// teksDesimal mengubah nilai uang menjadi TEKS untuk dikirim ke Oracle.
//
// ⛔ Uang tidak pernah melewati float (ADR-U-0003, ADR-U-0016). Nilai kosong
// dikirim sebagai NULL, bukan sebagai nol (ADR-U-0027).
func teksDesimal(m uang.Money) any {
	if m.Kosong() {
		return nil
	}
	return utils.FormatDecimal(m.Amount)
}

func teksRasio(r uang.Ratio) any {
	if r.Kosong() {
		return nil
	}
	return utils.FormatDecimal(r.Value)
}

// uraiUang membungkus uraiDesimal menjadi nilai uang beserta mata uangnya.
func uraiUang(idBaris, kolom string, v sql.NullString, mataUang string) (uang.Money, error) {
	m := uang.Money{Currency: mataUang}
	d, err := db.UraiDesimal(idBaris, kolom, v)
	if err != nil {
		return m, err
	}
	m.Amount = d
	return m, nil
}

// uraiRasio sama dengan uraiUang, untuk nilai yang berupa perbandingan.
//
// Rasio dan uang sengaja bertipe berbeda supaya keduanya tidak pernah
// terjumlahkan (ADR-F-0004): 30 persen ditambah Rp30 tidak berarti apa-apa.
func uraiRasio(idBaris, kolom string, v sql.NullString) (uang.Ratio, error) {
	var r uang.Ratio
	d, err := db.UraiDesimal(idBaris, kolom, v)
	if err != nil {
		return r, err
	}
	r.Value = d
	return r, nil
}

func waktuJadiNil(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}

// ErrIdentitasBelumAda menolak menyimpan pohon yang akar atau headernya belum
// punya nomor.
//
// ⛔ Nomor CLM-xxxxxx TIDAK dibangkitkan di sini. Pembangkitnya masih
// `[terbuka]` di tiket 14 - siapa yang membuatnya, apakah ada sequence di
// belakang prefiks, apakah di-reset per tahun - dan pemiliknya DBA atau work
// owner. Mengarang nomornya akan menetapkan yang belum diputuskan, jadi yang
// dilakukan di sini adalah menolak dengan terang.
var ErrIdentitasBelumAda = errors.New(
	"repository: T_WORK_CLAIM.ID belum terisi; pembangkit nomor CLM-/KMT- masih terbuka")

// sisipBarisAdjustment menulis SATU baris adjustment.
//
// Dipakai `Simpan` (pohon utuh) dan `SisipkanBaris` (satu baris lanjutan).
// Kolom dan urutannya hidup di sini saja.
func (r *PohonKlaim) sisipBarisAdjustment(ctx context.Context, tx *db.Tx,
	tabel, pesertaID string, adj models.BarisAdjustment) error {

	return r.exec(ctx, tx, fmt.Sprintf(`INSERT INTO %s
		(ID, PREMIUM_LIST_DETAIL_ID, CLAIM_AMOUNT, CURRENCY, STS_REJECT,
		 ACCEPTED_NO, ACCEPTATION_DATE, KOMITE_ID,
		 NAME_OF_BANK, ID_BANK, ACCOUNT_NO,
		 SHARE_NUSANTARA_RE, CEDING_RETENTION, SUM_REASURED, SUM_INSURED,
		 SHARE_RETRO, RETROCEDED_SHARE, CURRENCY_ID)
		VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18)`, tabel),
		adj.ID, pesertaID, teksDesimal(adj.JumlahKlaim),
		db.KosongJadiNil(adj.JumlahKlaim.Currency), db.KosongJadiNil(adj.KodeStatus),
		db.KosongJadiNil(adj.NomorAkseptasi), waktuJadiNil(adj.TanggalAkseptasi),
		db.KosongJadiNil(adj.KomiteID),
		db.KosongJadiNil(adj.NamaBank), db.KosongJadiNil(adj.IDBank),
		db.KosongJadiNil(adj.NomorRekening),
		teksDesimal(adj.ShareNusantaraRe), teksDesimal(adj.CedingRetention),
		teksDesimal(adj.SumReasured), teksDesimal(adj.SumInsured),
		teksDesimal(adj.ShareRetro), teksDesimal(adj.RetrocededShare),
		db.KosongJadiNil(adj.CurrencyID))
}

// ErrBarisBaruBerkeputusan - baris baru tidak pernah lahir sudah diputus.
//
// Keputusan atas baris milik Komite Claim Life
// `[keputusan work owner 2026-09-15]`; jalur ini hanya melahirkan baris
// Outstanding.
var ErrBarisBaruBerkeputusan = errors.New(
	"repository: baris adjustment baru tidak boleh lahir sudah berkeputusan")

// PeriksaBarisBaru menolak kode status yang bukan pembuka putaran.
//
// Hanya Outstanding dan kosong yang lolos: baris baru memulai putaran, ia
// tidak lahir sudah diputus.
func PeriksaBarisBaru(kodeStatus string) error {
	if kodeStatus != "" && kodeStatus != kontrak.KodeOutstanding {
		return fmt.Errorf("%w: baris baru berkode %q",
			ErrBarisBaruBerkeputusan, kodeStatus)
	}
	return nil
}

// SisipkanBaris menulis satu baris adjustment putaran berikutnya (tiket 11).
//
// Pengenalnya diambil dari `SEQ_CLAIMLF_ADJ` bila kosong - identitas seluruh
// tabel `T_CLAIMLF_*` berasal dari sequence, bukan dari cap waktu maupun teks
// yang disusun sendiri (ADR-U-0006).
//
// ⛔ Ia menulis STS_REJECT baris BARU, bukan mengubah baris lama. Keputusan
// atas baris lama milik Komite Claim Life.
func (r *PohonKlaim) SisipkanBaris(ctx context.Context, tx *db.Tx,
	pesertaID string, adj models.BarisAdjustment) (string, error) {

	if strings.TrimSpace(pesertaID) == "" {
		return "", fmt.Errorf("repository: pengenal peserta kosong")
	}
	// ⛔ PENJAGA DI TEMPAT YANG BENAR. Jalur ini melahirkan baris BARU, dan
	// baris baru tidak pernah lahir sudah berkeputusan. Sebelumnya penjagaan
	// itu hanya berupa pola atas teks `hasilkomite.go` - dan pola itu dapat
	// dielakkan oleh `models.BarisAdjustment{KodeStatus: "1"}`, yang memakai
	// titik dua, bukan tanda sama dengan. Di sini tidak ada bentuk penulisan
	// yang dapat mengelakkannya: nilainya diperiksa saat jalan.
	if err := PeriksaBarisBaru(adj.KodeStatus); err != nil {
		return "", err
	}
	tabel, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return "", err
	}
	if adj.ID == "" {
		id, err := r.db.NomorBerikut(ctx, tx, "SEQ_CLAIMLF_ADJ")
		if err != nil {
			return "", err
		}
		adj.ID = id
	}
	if err := r.sisipBarisAdjustment(ctx, tx, tabel, pesertaID, adj); err != nil {
		return "", err
	}
	return adj.ID, nil
}

// isiIdentitas melengkapi identitas yang masih kosong dari sequence-nya.
//
// Identitas yang SUDAH terisi tidak disentuh, supaya pemanggil yang membawa
// nomornya sendiri - misalnya migrasi data lama - tetap dapat menentukannya.
func (r *PohonKlaim) isiIdentitas(ctx context.Context, tx *db.Tx, p *models.PohonKlaim) error {
	if p.Work.ID == "" {
		return ErrIdentitasBelumAda
	}
	if p.Klaim.ID == "" {
		// Shared primary key: header memakai nomor akar apa adanya.
		p.Klaim.ID = p.Work.ID
	}
	for i := range p.Klaim.Peserta {
		ps := &p.Klaim.Peserta[i]
		if ps.ID == "" {
			id, err := r.db.NomorBerikut(ctx, tx, "SEQ_CLAIMLF_PLD")
			if err != nil {
				return err
			}
			ps.ID = id
		}
		for j := range ps.Baris {
			adj := &ps.Baris[j]
			if adj.ID == "" {
				id, err := r.db.NomorBerikut(ctx, tx, "SEQ_CLAIMLF_ADJ")
				if err != nil {
					return err
				}
				adj.ID = id
			}
			for k := range adj.Spreading {
				spr := &adj.Spreading[k]
				if spr.ID == "" {
					id, err := r.db.NomorBerikut(ctx, tx, "SEQ_CLAIMLF_SPR")
					if err != nil {
						return err
					}
					spr.ID = id
				}
				for l := range spr.Retro {
					if spr.Retro[l].ID == "" {
						id, err := r.db.NomorBerikut(ctx, tx, "SEQ_CLAIMLF_SPR_RETRO")
						if err != nil {
							return err
						}
						spr.Retro[l].ID = id
					}
				}
			}
		}
	}
	return nil
}

// Simpan menulis satu pohon klaim, beserta baris datarnya, dalam SATU
// transaksi.
//
// Transaksinya dibuka dan ditutup oleh pemanggil di services (ADR-U-0029);
// fungsi ini hanya ikut di dalamnya.
//
// Identitas yang masih kosong di tingkat T_CLAIMLF_* diisi dari sequence
// (ADR-U-0006). Nomor akar tidak pernah dikarang - lihat ErrIdentitasBelumAda.
func (r *PohonKlaim) Simpan(ctx context.Context, tx *db.Tx, p models.PohonKlaim) error {
	if err := r.isiIdentitas(ctx, tx, &p); err != nil {
		return err
	}
	work, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	header, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return err
	}
	pesertaT, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	adjT, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return err
	}
	sprT, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT_SPREADING")
	if err != nil {
		return err
	}
	retroT, err := r.db.Qualify(namaTabelRetro)
	if err != nil {
		return err
	}
	datarT, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return err
	}

	// Tingkat 1 - akar work object.
	err = r.exec(ctx, tx, fmt.Sprintf(`INSERT INTO %s
		(ID, COVER_KEY, LINI, PY_POSITION, SENDTO_ADMIN,
		 SENDTO_MEDICAL, TYPE, CASE_ID, CREATE_OP, CREATE_OP_NAME, TGL_UPDATE,
		 TAHAP, TGL_CREATE)
		VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13)`, work),
		p.Work.ID, db.KosongJadiNil(p.Work.CoverKey), db.KosongJadiNil(p.Work.Lini),
		db.KosongJadiNil(p.Work.PyPosition),
		db.KosongJadiNil(p.Work.SendtoAdmin), db.KosongJadiNil(p.Work.SendtoMedical),
		db.KosongJadiNil(p.Work.Type), db.KosongJadiNil(p.Work.CaseID),
		db.KosongJadiNil(p.Work.CreateOp), db.KosongJadiNil(p.Work.CreateOpName),
		waktuJadiNil(p.Work.TglUpdate),
		db.KosongJadiNil(p.Work.Tahap), waktuJadiNil(p.Work.TglCreate))
	if err != nil {
		return err
	}

	// Tingkat 2 - header klaim. ID-nya SAMA dengan ID work object: shared PK.
	err = r.exec(ctx, tx, fmt.Sprintf(`INSERT INTO %s
		(ID, CLAIM_NO, POLICY_NO, BUSINESS_NAME, STS_REJECT, CURRENCY)
		VALUES (:1,:2,:3,:4,:5,:6)`, header),
		p.Work.ID, db.KosongJadiNil(p.Klaim.NomorKlaim), db.KosongJadiNil(p.Klaim.NomorPolis),
		db.KosongJadiNil(p.Klaim.NamaBisnis), db.KosongJadiNil(p.Klaim.KodeStatus),
		// Butir z1: mata uang header, supaya ClaimRetro yang dibaca kembali
		// punya mata uang - uang tanpa mata uang tidak bermakna.
		db.KosongJadiNil(p.Klaim.ClaimRetro.Currency))
	if err != nil {
		return err
	}

	for _, ps := range p.Klaim.Peserta {
		// Daftar kolomnya datang dari kolompeserta.go - satu daftar untuk
		// tulis dan baca, supaya urutan bind tidak mungkin berselisih.
		q, nilai := insertPeserta(pesertaT, ps.ID, p.Work.ID, ps)
		if err = r.exec(ctx, tx, q, nilai...); err != nil {
			return err
		}

		for _, adj := range ps.Baris {
			// Kedelapan kolom warisan ikut ditulis sejak 26-09-2026. Sebelumnya
			// hanya CURRENCY yang ditulis, sehingga tujuh kolom lain SELALU
			// NULL walau DDL menyediakannya - pewarisan yang tiket janjikan
			// tidak pernah sampai ke basis data.
			// ⛔ SATU pernyataan insert baris adjustment di seluruh
			// repository, dipakai juga oleh `SisipkanBaris` (tiket 11). Dua
			// salinan berarti dua kesempatan untuk berbeda - dan yang paling
			// mudah tertinggal justru kolom warisannya, yang sudah pernah
			// tertinggal sekali.
			if err := r.sisipBarisAdjustment(ctx, tx, adjT, ps.ID, adj); err != nil {
				return err
			}

			for _, spr := range adj.Spreading {
				err = r.exec(ctx, tx, fmt.Sprintf(`INSERT INTO %s
					(ID, ADJUSTMENT_ID, TREATY_TYPE_ID, TREATY_TYPE_NAME,
					 TREATY_YEAR_LIFE, RETROCADED_SHARE, RATE, IDR, USD, CURRENCY)
					VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10)`, sprT),
					spr.ID, adj.ID, db.KosongJadiNil(spr.TreatyTypeID),
					db.KosongJadiNil(spr.TreatyTypeName), db.KosongJadiNil(spr.TreatyYearLife),
					teksDesimal(spr.RetrocadedShare), teksRasio(spr.Rate),
					teksDesimal(spr.IDR), teksDesimal(spr.USD),
					db.KosongJadiNil(spr.Currency))
				if err != nil {
					return err
				}

				for _, rt := range spr.Retro {
					err = r.exec(ctx, tx, fmt.Sprintf(`INSERT INTO %s
						(ID, SPREADING_ID, REINSURER_NAME, PERCENT_SHARE, AMOUNT, RATE,
						 PREMIUM_SPREADED_GROSS, PREMIUM_SPREADED_NET, COMMISION,
						 OVR_COMM, TREATY_TYPE_ID, TREATY_TYPE_NAME)
						VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12)`, retroT),
						rt.ID, spr.ID, db.KosongJadiNil(rt.ReinsurerName),
						teksRasio(rt.PercentShare), teksDesimal(rt.Amount),
						teksRasio(rt.Rate), teksDesimal(rt.PremiumSpreadedGross),
						teksDesimal(rt.PremiumSpreadedNet), teksRasio(rt.Commision),
						teksRasio(rt.OvrComm), db.KosongJadiNil(rt.TreatyTypeID),
						db.KosongJadiNil(rt.TreatyTypeName))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// Penulisan kedua - baris datar warisan, di dalam transaksi yang SAMA.
	//
	// ⚠️ [terbuka - tiket 02/03] Delapan belas dari 55 kolom yang ditulis; 37
	// sisanya tinggal NULL. Hilir - Arasapas - membaca tabel ini, sehingga
	// kekosongan itu terlihat olehnya. Ia TIDAK diperbaiki di tiket 14: sebagian
	// besar kolom yang hilang adalah atribut polis (POLICY_HOLDER,
	// NAME_OF_INSURED, SEX, DOB, AGE, PLAN, BEGIN_DATE, LAPSE_DATE,
	// EXPIRED_DATE, SUM_INSURED, SUM_REASURED, CEDING_RETENTION,
	// SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE, EM_PERCENT, WPC,
	// CEDINGCO, CEDINGCONAME, SOB, SOBNAME, BUSINESSID, TYPECEDING) yang
	// keputusannya dibaca HIDUP dari tabel polis, dan pembacaan itu milik tiket
	// 02 dan 03. Sisanya (DISEASE, ICD_CODE, NOTES, KETERANGAN, STATUS,
	// RETROID, RETRONAME, SECURITYREINSURERID, SECURITYREINSURER,
	// CONFIRMATION_DATE, CLAIM_RECEIVED_DATE, COMPLETE_DATE, PRODUCTNAMEID,
	// PRODUCTNAME) menunggu tahap yang mengisinya - medical check, penutupan,
	// dan modul PremiumList Life.
	for _, b := range BarisLamaDari(p) {
		// ⛔ Dipagari SEBELUM bind: kolom yang di tabel warisan bertipe NUMBER
		// harus menerima bilangan, dan galat yang menyebut kolom serta nilainya
		// jauh lebih berguna daripada ORA-01722 (butir s1).
		if err = PeriksaNilaiWarisan(b); err != nil {
			return err
		}
		err = r.exec(ctx, tx, fmt.Sprintf(`INSERT INTO %s
			(ID, CASEID, NO_CLAIM, POLICY_NO, CERTIFICATE_NO, PL_NUMBER, BUSINESSNAME,
			 CLAIM_RETRO, CURRENCY, CLAIM_AMOUNT, STS_REJECT, NO_ACCEPTATION,
			 ACCEPTATION_DATE, TYPE, CREATEOPNAME,
			 NAME_OF_BANK, IDBANK, ACCOUNTNO)
			VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18)`, datarT),
			b.ID, db.KosongJadiNil(b.CASEID), db.KosongJadiNil(b.NO_CLAIM),
			db.KosongJadiNil(b.POLICY_NO), db.KosongJadiNil(b.CERTIFICATE_NO),
			db.KosongJadiNil(b.PL_NUMBER), db.KosongJadiNil(b.BUSINESSNAME),
			db.KosongJadiNil(b.CLAIM_RETRO), db.KosongJadiNil(b.CURRENCY),
			db.KosongJadiNil(b.CLAIM_AMOUNT), db.KosongJadiNil(b.STS_REJECT),
			db.KosongJadiNil(b.NO_ACCEPTATION), db.KosongJadiNil(b.ACCEPTATION_DATE),
			db.KosongJadiNil(b.TYPE), db.KosongJadiNil(b.CREATEOPNAME),
			db.KosongJadiNil(b.NAME_OF_BANK), db.KosongJadiNil(b.IDBANK),
			db.KosongJadiNil(b.ACCOUNTNO))
		if err != nil {
			return err
		}
	}
	// ⛔ OQ-N2 DITUTUP (GILIRAN-17): NAME_OF_INSURED, DOB, CEDINGCO cermin
	// diisi DI DALAM SQL dari baris sumber peserta, transaksi yang sama -
	// nama dan tanggal lahir tidak pernah melintasi Go. Peserta tanpa baris
	// sumber (`SOURCE_ID` kosong) hanya lahir dari uji/migrasi: pendaftaran
	// membaca peserta ULANG dari sumbernya, dan Save to RNM menolak peserta
	// tanpa SOURCE_ID.
	cermin := NewKlaimLife(r.db)
	for _, ps := range p.Klaim.Peserta {
		if strings.TrimSpace(ps.SumberID) == "" {
			continue
		}
		k := KunciPesertaSumber{PLNumber: ps.NomorPremiList, Sertifikat: ps.NomorSertifikat, SumberID: ps.SumberID}
		for _, adj := range ps.Baris {
			if err = cermin.IsiTertanggungCermin(ctx, tx, adj.ID, p.Work.CaseID, k, p.Klaim.NomorPolis); err != nil {
				return err
			}
		}
	}
	return nil
}

// sqlHapusCerminBelumDisimpan - OQ-M6 (temuan /code-review GILIRAN-17): baris
// cermin peserta yang DICABUT sebelum Save to RNM dibuang. Di Pega baris
// cermin baru lahir saat Save Outstanding (`InsertJsonKlaimLife_sql`), jadi
// peserta yang dilepas sebelumnya tidak pernah punya baris; aplikasi ini
// menulisnya saat pendaftaran. Dikunci `CASEID` + `STS_REJECT IS NULL` -
// baris era Pega dan baris yang sudah berstatus tidak tersentuh.
func sqlHapusCerminBelumDisimpan(lama string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND CASEID = :2 AND STS_REJECT IS NULL`, lama)
}

// HapusCerminBelumDisimpan menjalankannya di transaksi pemanggil; nol baris
// bukan galat (baris lahir sebelum cermin ditulis, atau sudah tidak ada).
func (r *KlaimLife) HapusCerminBelumDisimpan(ctx context.Context, tx *db.Tx, adjID, caseID string) error {
	lama, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return err
	}
	q := sqlHapusCerminBelumDisimpan(lama)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, adjID, caseID); err != nil {
		return fmt.Errorf("repository: membuang cermin peserta tercabut %s: %w", adjID, err)
	}
	return nil
}

// sqlSetelCerminOutstanding - OQ-N13 (GILIRAN-18) `[keputusan asisten dari
// bukti; veto work owner]`: Save to RNM menyetel status cermin seperti Pega.
//
// `[terverifikasi]` `SaveOutStandingLife_Act` langkah 22.1.3.1 (b10115, hidup)
// memanggil `InsertJsonKlaimLife_sql` per baris adjustment yang baru disimpan;
// rule itu menulis `STS_REJECT = '0'` (b176). Di Pega baris cerminnya LAHIR
// di sana; di aplikasi ini ia lahir saat pendaftaran (status NULL), jadi yang
// ditiru di sini hanya statusnya.
//
// ⛔ Larangan lama "tabel warisan baca-saja di Save to RNM" DICABUT UNTUK
// KOLOM INI SAJA (tiket 03, 30-09-2026). SET satu kolom; dikunci `ID` +
// `CASEID` + `STS_REJECT IS NULL` - baris era Pega dan baris yang sudah
// berkeputusan (Komite) tidak tersentuh.
func sqlSetelCerminOutstanding(lama string) string {
	return fmt.Sprintf(`UPDATE %s SET STS_REJECT = :1 WHERE ID = :2 AND CASEID = :3 AND STS_REJECT IS NULL`, lama)
}

// SetelCerminOutstanding menjalankannya di transaksi pemanggil dan
// mengembalikan cacah baris cermin yang disetel.
//
// ⚠️ Nol baris BUKAN galat: baris adjustment putaran Komite
// (`SisipkanBaris`) tidak punya baris cermin, dan cermin yang sudah
// berkeputusan tidak ditimpa. Menyisip cermin yang hilang berarti menulis
// kolom lain - di luar pencabutan larangan yang sempit ini.
func (r *KlaimLife) SetelCerminOutstanding(ctx context.Context, tx *db.Tx, adjID, caseID string) (int, error) {
	if strings.TrimSpace(adjID) == "" || strings.TrimSpace(caseID) == "" {
		return 0, fmt.Errorf("repository: menolak menyetel cermin tanpa ID adjustment dan CASEID (%q, %q)",
			adjID, caseID)
	}
	lama, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return 0, err
	}
	q := sqlSetelCerminOutstanding(lama)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	hasil, err := tx.ExecContext(ctx, q, kontrak.KodeOutstanding, adjID, caseID)
	if err != nil {
		return 0, fmt.Errorf("repository: menyetel status cermin %s: %w", adjID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repository: mencacah cermin tersetel: %w", err)
	}
	return int(n), nil
}

// sqlIkutkanStatusCermin - temuan /code-review GILIRAN-18 atas OQ-N13: status
// cermin MENGIKUTI setiap penulis `PerbaruiStatusBaris`.
//
// `[terverifikasi]` Di Pega penolakan Admin (`RejectOSClaimLife_Act` langkah 5
// `Insert ke OS` b1970, `pyStepsBlockName` kosong b1982) dan akseptasi
// (`SaveAdjustment_Act` 1.6.2 b2363, kosong b2375) menulis cermin lewat
// `UpdateOsAkseptasiClaimLife_sql`. Tanpa ini, cermin yang disetel '0' oleh
// Save to RNM TERTAHAN '0' sesudah baris ditolak, dan gerbang 11.4
// (`==0 || ==1` pesan) memblokir setiap klaim kematian berikutnya atas
// tertanggung itu - regresi yang dilahirkan N13 sendiri.
//
// ⛔ `STS_REJECT` SAJA - batas pencabutan larangan yang sama dengan N13.
// Rule itu meng-`INSERT` baris baru; di sini baris cerminnya sudah ada sejak
// pendaftaran, jadi DIPERBARUI (penyimpangan yang sama dengan
// `RekamAkhirWarisan` Komite). Kunci `ID` baris adjustment + `CASEID` klaim
// pemilik peserta (`NVL(CASE_ID, ID)`, bentuk yang sama dengan pengecualian
// pemeriksa ganda). Hanya cermin yang masih NULL atau masih berkode lama yang
// ikut - cermin yang sudah ditulis jalur lain (Komite) tidak ditimpa, dan
// cermin NULL milik klaim yang disimpan sebelum N13 ikut pulih.
func sqlIkutkanStatusCermin(lama, work, pes string, dariKosong bool) string {
	syarat := "o.STS_REJECT IS NULL"
	if !dariKosong {
		syarat = "(o.STS_REJECT IS NULL OR o.STS_REJECT = :4)"
	}
	return fmt.Sprintf(`UPDATE %s o SET o.STS_REJECT = :1
	 WHERE o.ID = :2
	   AND o.CASEID = (SELECT NVL(w.CASE_ID, w.ID) FROM %s w, %s p
	                    WHERE p.ID = :3 AND w.ID = p.CLAIM_ID)
	   AND %s`, lama, work, pes, syarat)
}

// sqlIsiCedingCermin - CEDINGCO cermin dari POLIS (OQ-N2, GILIRAN-17):
// `pyWorkPage.PolicyDataLife.CedingCo` b9226 -> `InsertJsonKlaimLife_sql`
// b114. `T_PREMIUM_LIST.CEDING_CO` baris PROD_KE terakhir - bentuk yang sama
// dengan `sqlPolisRingkas` - sumber yang sama dengan yang dibandingkan
// pemeriksa klaim ganda (`o.CEDINGCO = :4`). Polis yang tidak ada = NULL,
// sama seperti sebelum GILIRAN-17.
func sqlIsiCedingCermin(lama, polis string) string {
	return fmt.Sprintf(`UPDATE %s o
	   SET CEDINGCO = (SELECT pl.CEDING_CO FROM %s pl WHERE pl.NO_POLIS = :1
	                    ORDER BY NVL(pl.PROD_KE, 0) DESC FETCH FIRST 1 ROWS ONLY)
	 WHERE o.ID = :2 AND o.CASEID = :3`, lama, polis)
}

// isiCedingCermin menjalankan `sqlIsiCedingCermin` di transaksi pemanggil.
func (r *KlaimLife) isiCedingCermin(ctx context.Context, tx *db.Tx, lama, adjID, caseID, nomorPolis string) error {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return err
	}
	q := sqlIsiCedingCermin(lama, polis)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, db.KosongJadiNil(nomorPolis), adjID, caseID); err != nil {
		return fmt.Errorf("repository: mengisi CEDINGCO cermin %s: %w", adjID, err)
	}
	return nil
}

// Hapus membuang satu pohon klaim.
//
// ⛔ HANYA untuk test kaskade. Ini DELETE fisik: barisnya benar-benar hilang.
// Jalur pengguna TIDAK boleh memakainya - ADR-U-0031 menetapkan penghapusan
// klaim berupa penanda ditambah nilai balik, bukan DELETE, dan pelaksananya
// tiket 15. Fungsi ini ada semata supaya AC 38 - kaskade sampai cicit - dapat
// dibuktikan dengan test. Nol pemanggil di luar berkas _test.go, dan
// TestHapusTidakDipanggilDiLuarTest menjaganya tetap begitu.
//
// Tabel anak di relasi 3, 4, 5, dan 6 ikut terhapus oleh kaskade Oracle sampai
// tingkat terdalam. Yang TIDAK kaskade dan karena itu dihapus di sini adalah
// T_CLAIMLF_DOCUMENT (lintas-lini), baris datar warisan, dan akar work object.
//
// Urutannya penting: T_CLAIMLF_DOCUMENT dibuang LEBIH DULU. Kunci tamunya sengaja
// tanpa ON DELETE, sehingga menghapus header selagi masih ada baris dokumen
// akan ditolak Oracle dengan ORA-02292 (induk masih punya anak).
func (r *PohonKlaim) HapusFisik(ctx context.Context, tx *db.Tx, id, caseID string) error {
	header, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return err
	}
	work, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	datarT, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return err
	}
	dokT, err := r.db.Qualify("T_CLAIMLF_DOCUMENT")
	if err != nil {
		return err
	}
	pesT, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}

	// Dokumen milik seluruh peserta klaim ini - dihapus di Go, bukan kaskade.
	if err := r.exec(ctx, tx, fmt.Sprintf(
		`DELETE FROM %s WHERE PREMIUM_LIST_DETAIL_ID IN
		   (SELECT ID FROM %s WHERE CLAIM_ID = :1)`, dokT, pesT), id); err != nil {
		return err
	}

	// Menghapus header mengkaskade ke peserta, adjustment, spreading, dan
	// spreading retro - empat tingkat sekaligus.
	if err := r.exec(ctx, tx,
		fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, header), id); err != nil {
		return err
	}
	if err := r.exec(ctx, tx,
		fmt.Sprintf(`DELETE FROM %s WHERE CASEID = :1`, datarT), caseID); err != nil {
		return err
	}
	return r.exec(ctx, tx, fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, work), id)
}

// AmbilSpreading membaca spreading beserta pecahan retronya untuk satu klaim,
// dikelompokkan per baris adjustment.
//
// Pembacaan ini turun sampai CICIT - tingkat terdalam pohon.
func (r *PohonKlaim) AmbilSpreading(ctx context.Context, klaimID string) (
	map[string][]models.Spreading, error) {
	sprT, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT_SPREADING")
	if err != nil {
		return nil, err
	}
	adjT, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return nil, err
	}
	pesT, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return nil, err
	}
	retroT, err := r.db.Qualify(namaTabelRetro)
	if err != nil {
		return nil, err
	}

	// Seluruh kolom angka diminta lewat TO_CHAR ber-argumen NLS, sehingga ia
	// tiba sebagai TEKS dan tidak pernah melewati float (ADR-U-0003).
	// RETROCADED_SHARE dan RATE ikut dibaca: keduanya rasio, dan ronde 1
	// melewatkannya sehingga pecahan retro pulang tanpa dasar pembagiannya.
	q := fmt.Sprintf(`SELECT s.ADJUSTMENT_ID, s.ID, s.TREATY_TYPE_ID,
		       s.TREATY_TYPE_NAME, s.TREATY_YEAR_LIFE, s.CURRENCY,
		       `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`
		  FROM %s s
		  JOIN %s a ON a.ID = s.ADJUSTMENT_ID
		  JOIN %s p ON p.ID = a.PREMIUM_LIST_DETAIL_ID
		 WHERE p.CLAIM_ID = :1 AND p.STS_HAPUS IS NULL
		 ORDER BY s.ADJUSTMENT_ID, s.ID`,
		"s.IDR", "s.USD", "s.RETROCADED_SHARE", "s.RATE", sprT, adjT, pesT)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, klaimID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca spreading: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := map[string][]models.Spreading{}
	var idSpreading []string
	indeks := map[string]string{} // idSpreading -> adjustmentID
	for rows.Next() {
		var adjID, id string
		var tipeID, tipeNama, tahun, mataUang sql.NullString
		var idr, usd, retrocaded, rate sql.NullString
		if err := rows.Scan(&adjID, &id, &tipeID, &tipeNama, &tahun,
			&mataUang, &idr, &usd, &retrocaded, &rate); err != nil {
			return nil, err
		}
		s := models.Spreading{
			ID: id, AdjustmentID: adjID,
			TreatyTypeID: tipeID.String, TreatyTypeName: tipeNama.String,
			TreatyYearLife: tahun.String, Currency: mataUang.String,
			IDR: uang.Money{Currency: "IDR"}, USD: uang.Money{Currency: "USD"},
		}
		if s.IDR, err = uraiUang(id, "IDR", idr, "IDR"); err != nil {
			return nil, err
		}
		if s.USD, err = uraiUang(id, "USD", usd, "USD"); err != nil {
			return nil, err
		}
		// RETROCADED_SHARE adalah uang, dan mata uangnya kolom CURRENCY baris
		// ini - bukan IDR/USD, yang di sini nama kolom KAPASITAS treaty-year.
		if s.RetrocadedShare, err = uraiUang(id, "RETROCADED_SHARE", retrocaded,
			mataUang.String); err != nil {
			return nil, err
		}
		if s.Rate, err = uraiRasio(id, "RATE", rate); err != nil {
			return nil, err
		}
		out[adjID] = append(out[adjID], s)
		idSpreading = append(idSpreading, id)
		indeks[id] = adjID
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Tingkat terdalam, dibaca per baris spreading.
	for _, sid := range idSpreading {
		retro, err := r.ambilRetro(ctx, retroT, sid)
		if err != nil {
			return nil, err
		}
		adjID := indeks[sid]
		for i := range out[adjID] {
			if out[adjID][i].ID == sid {
				out[adjID][i].Retro = retro
			}
		}
	}
	return out, nil
}

func (r *PohonKlaim) ambilRetro(ctx context.Context, tabel, spreadingID string) (
	[]models.SpreadingRetro, error) {
	// Tujuh kolom angka, bukan satu. Ronde 1 hanya membaca AMOUNT, sehingga
	// pecahan retro pulang tanpa persentase, rate, premi, maupun komisinya.
	q := fmt.Sprintf(`SELECT ID, REINSURER_NAME, TREATY_TYPE_ID, TREATY_TYPE_NAME,
		       `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`,
		       `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`
		  FROM %s WHERE SPREADING_ID = :1 ORDER BY ID`,
		"AMOUNT", "PERCENT_SHARE", "RATE", "PREMIUM_SPREADED_GROSS",
		"PREMIUM_SPREADED_NET", "COMMISION", "OVR_COMM", tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, spreadingID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca spreading retro: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.SpreadingRetro{}
	for rows.Next() {
		var id string
		var nama, tipeID, tipeNama sql.NullString
		var jumlah, persen, rate, gross, net, komisi, ovr sql.NullString
		if err := rows.Scan(&id, &nama, &tipeID, &tipeNama, &jumlah,
			&persen, &rate, &gross, &net, &komisi, &ovr); err != nil {
			return nil, err
		}
		rt := models.SpreadingRetro{
			ID: id, SpreadingID: spreadingID, ReinsurerName: nama.String,
			TreatyTypeID: tipeID.String, TreatyTypeName: tipeNama.String,
		}
		var err error
		if rt.Amount, err = uraiUang(id, "AMOUNT", jumlah, ""); err != nil {
			return nil, err
		}
		if rt.PercentShare, err = uraiRasio(id, "PERCENT_SHARE", persen); err != nil {
			return nil, err
		}
		if rt.Rate, err = uraiRasio(id, "RATE", rate); err != nil {
			return nil, err
		}
		if rt.PremiumSpreadedGross, err = uraiUang(id, "PREMIUM_SPREADED_GROSS", gross, ""); err != nil {
			return nil, err
		}
		if rt.PremiumSpreadedNet, err = uraiUang(id, "PREMIUM_SPREADED_NET", net, ""); err != nil {
			return nil, err
		}
		if rt.Commision, err = uraiRasio(id, "COMMISION", komisi); err != nil {
			return nil, err
		}
		if rt.OvrComm, err = uraiRasio(id, "OVR_COMM", ovr); err != nil {
			return nil, err
		}
		out = append(out, rt)
	}
	return out, rows.Err()
}

// AmbilBarisLama membaca seluruh baris datar warisan milik satu kasus.
//
// Inilah pasangan baca dari BongkarBarisLama: yang ini berbicara ke Oracle,
// yang itu murni. Memisahkan keduanya membuat logika pembongkaran - bagian yang
// paling mudah salah - dapat diuji tanpa instance sama sekali.
//
// ⛔ Seluruh kolom angka dan tanggal diminta lewat TO_CHAR ber-argumen NLS,
// sehingga nilainya tiba sebagai TEKS dengan titik sebagai pemisah desimal apa
// pun setelan sesi. BarisLama memang bertipe teks seluruhnya: pada tahap ini
// nilainya belum diurai, dan yang tidak terurai kelak DILAPORKAN oleh
// BongkarBarisLama, bukan ditebak. Daftar kolomnya ada di barislamakolom.go.
//
// Urutan baris dikunci ORDER BY ID supaya hasilnya dapat diulang.
func (r *PohonKlaim) AmbilBarisLama(ctx context.Context, caseID string) ([]BarisLama, error) {
	tabel, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE CASEID = :1 ORDER BY ID`,
		ekspresiSelectLama(), tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, caseID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca baris lama: %w", err)
	}
	defer func() { _ = rows.Close() }()

	jumlahKolom := len(medanBarisLama(&BarisLama{}))
	out := []BarisLama{}
	for rows.Next() {
		nilai, tujuan := tujuanScanLama(jumlahKolom)
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris lama: %w", err)
		}
		var b BarisLama
		salinKeBarisLama(&b, nilai)
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca baris lama: %w", err)
	}
	return out, nil
}

// CacahBarisLama menghitung baris datar warisan milik satu CASEID. Dipakai
// untuk membuktikan bahwa penulisan kedua benar-benar terjadi.
func (r *PohonKlaim) CacahBarisLama(ctx context.Context, caseID string) (int, error) {
	tabel, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return 0, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE CASEID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, q, caseID).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: mencacah baris lama: %w", err)
	}
	return n, nil
}

// Dampak menghitung baris yang akan ikut terhapus bersama sebuah klaim.
//
// ⛔ Dihitung dari DATA, bukan diperkirakan: angka yang ditampilkan peringatan
// wajib sama persis dengan yang benar-benar terhapus, dan perkiraan tidak
// pernah sama persis.
//
// ⚠️ Tujuh `SELECT COUNT` terpisah, bukan satu kueri berjenjang. Satu kueri
// akan lebih rapi dan lebih sulit dibaca, sedangkan tiap angka di sini muncul
// sendiri-sendiri di layar dan wajib dapat ditelusuri sendiri-sendiri.
func (r *PohonKlaim) Dampak(ctx context.Context, klaimID, caseID string) (
	models.DampakHapus, error) {
	var d models.DampakHapus

	nama := []string{"T_CLAIMLF_PREMIUMLIST_DETAIL", "T_CLAIMLF_ADJUSTMENT",
		"T_CLAIMLF_ADJUSTMENT_SPREADING", namaTabelRetro, "T_CLAIMLF_DOCUMENT",
		"T_WORK_CLAIM", namaTabelLama, "T_GENERAL_CLAIM"}
	tabel := make([]string, len(nama))
	for i, n := range nama {
		t, err := r.db.Qualify(n)
		if err != nil {
			return d, err
		}
		tabel[i] = t
	}
	pes, adj, spr, retro, dok, work, datar, hdr := tabel[0], tabel[1], tabel[2],
		tabel[3], tabel[4], tabel[5], tabel[6], tabel[7]

	// Tiap baris: kueri, argumen, dan tujuan angkanya.
	langkah := []struct {
		q     string
		arg   string
		tuju  *int
		jenis string
	}{
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE CLAIM_ID = :1`, pes),
			klaimID, &d.Peserta, "peserta"},
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s a JOIN %s p
		      ON p.ID = a.PREMIUM_LIST_DETAIL_ID WHERE p.CLAIM_ID = :1`, adj, pes),
			klaimID, &d.Adjustment, "adjustment"},
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s s JOIN %s a ON a.ID = s.ADJUSTMENT_ID
		      JOIN %s p ON p.ID = a.PREMIUM_LIST_DETAIL_ID WHERE p.CLAIM_ID = :1`,
			spr, adj, pes), klaimID, &d.Spreading, "spreading"},
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s t JOIN %s s ON s.ID = t.SPREADING_ID
		      JOIN %s a ON a.ID = s.ADJUSTMENT_ID
		      JOIN %s p ON p.ID = a.PREMIUM_LIST_DETAIL_ID WHERE p.CLAIM_ID = :1`,
			retro, spr, adj, pes), klaimID, &d.SpreadingRetro, "spreading retro"},
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s d JOIN %s p
		      ON p.ID = d.PREMIUM_LIST_DETAIL_ID WHERE p.CLAIM_ID = :1`, dok, pes),
			klaimID, &d.Dokumen, "dokumen"},
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, work),
			klaimID, &d.WorkClaim, "baris work"},
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE CASEID = :1`, datar),
			caseID, &d.BarisDatarWarisan, "baris datar warisan"},
		// Header klaim itu sendiri - ia baris pertama yang dihapus, dan
		// melewatkannya membuat Total kurang satu dari yang sebenarnya.
		{fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, hdr),
			klaimID, &d.Header, "header klaim"},
	}
	for _, l := range langkah {
		if err := db.PeriksaSQL(l.q); err != nil {
			return d, err
		}
		if err := r.db.QueryRowContext(ctx, l.q, l.arg).Scan(l.tuju); err != nil {
			return d, fmt.Errorf("repository: mencacah %s: %w", l.jenis, err)
		}
	}
	return d, nil
}
