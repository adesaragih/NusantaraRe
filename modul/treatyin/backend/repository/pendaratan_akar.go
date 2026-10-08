package repository

// Medan KEPALA kontrak dan grid Rate of Exchange dibaca dari tabel
// pendaratan — bukan lagi dari `M_TREATY_IN.JSONDATA`.
//
// ⛔ Ini langkah terakhir larangan pemilik proses 6 Oktober 2026. Sesudah
// berkas ini, NOL nilai layar Treaty In datang dari dokumen.
//
// ---------------------------------------------------------------------
// ⭐ `NULL` MENGGANTIKAN "kuncinya tidak ada di dokumen", dan itu lebih kuat
// ---------------------------------------------------------------------
// Jalur lama memakai penunjuk `*string` untuk membedakan DUA keadaan yang
// berarti berbeda di layar:
//
//	kunci TIDAK ADA di dokumen   -> "tidak ada di sistem lama"
//	kunci ADA bernilai kosong    -> "belum diisi"
//
// Pembedaan itu TIDAK hilang: pemuat menulis `NULL` untuk kunci yang tidak
// ada dan teks kosong untuk kunci yang ada bernilai kosong, dan
// `sql.NullString.Valid` membacanya kembali persis. Yang berubah hanya dari
// mana jawabannya datang.
//
// ---------------------------------------------------------------------
// ⛔ RATE OF EXCHANGE — `TREATYEXCHANGEYEARLY`, dan BUKAN yang lain
// ---------------------------------------------------------------------
// Pemilik proses menegaskannya dua kali — 4 dan 6 Oktober 2026: pengganti
// `MATA_UANG_KONTRAK` adalah `TREATYEXCHANGEYEARLY`, tabel WARISAN yang sudah
// hidup. `MATA_UANG_KONTRAK` sendiri dicabut migrasi `434` dan tidak disebut
// satu kali pun di berkas ini.
//
// ⭐ Dan spesifikasi membenarkannya dengan bukti, bukan dengan wewenang:
// `KOREKSI-ERD-VERSUS-POOLDATA.md` §1.2 menyatakan `CurrencyList` di
// `M_TREATY_IN.JSONDATA` adalah **SALINAN** tabel ini, bukan sumbernya — dan
// mendaratkan salinan sebagai kalau ia sumber adalah cara tercepat membuat
// dua angka kurs yang berselisih tanpa ada yang tahu mana yang benar.
//
// Keempat kolom grid memetakan langsung, dari `RDBList/GetCurrencyToIDR_SQL`:
//
//	CURRENCY  -> Currency        TOIDR   -> Value to IDR
//	STARTDATE -> Valid From      ENDDATE -> Valid Until

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// Kolom skalar akar yang layar baca, berpasangan dengan kunci dokumennya.
//
// ⛔ Kunci dokumennya dibawa serta sebab `AdaDiJSON` memakai EJAAN DOKUMEN,
// bukan nama kolom: layar sudah menanyakan `adaDiJson["Bordeaux"]`, dan
// menukar sumber data tidak boleh menukar nama yang layar tanyakan.
var kolomRevisi = [][2]string{
	{"CONTRACTREFNO", "ContractRefNo"},
	{"BORDEAUX", "Bordeaux"},
	{"BORDEREAUXNOTE", "BordereauxNote"},
	{"ACCOUNTINGMODE", "AccountingMode"},
	{"ACCOUNTINGMODENONPROP", "AccountingModeNonProp"},
	{"TREATYLEADER", "TreatyLeader"},
	{"ISMULTIPLERETRO", "IsMultipleRetro"},
	{"EDMSTATE", "EDMState"},
	{"EDMMATERIALTYPE", "EDMMaterialType"},
	// ⭐ Dibaca sejak 6 Oktober 2026: tab `Information & Submit` memakainya.
	// `Section/TreatyInfoSubmit.xml` menyembunyikan tombol `Submit` ketika
	// `TreatyIn.StatusAkseptasi == 'Resolve Complete'` — kontrak yang sudah
	// tuntas tidak dapat dikirim ulang.
	{"STATUSAKSEPTASI", "StatusAkseptasi"},

	// ⭐ ENAM BELAS kolom migrasi `444`, 6 Oktober 2026.
	//
	// ⛔ Daftar ini TERPISAH dari `PetaPendaratan`, dan keterpisahan itu
	// sudah hampir memakan korban: migrasi `444` menambah kolomnya, pemuat
	// mengisinya, dan tanpa baris-baris di bawah nilainya akan mendarat
	// dengan sempurna lalu TIDAK PERNAH sampai ke layar. Nol galat, nol
	// uji merah — hanya medan yang tetap kosong.
	//
	// Yang menutup lubang itu `TestKolomRevisiAdaDiPeta`: setiap kolom di
	// sini wajib ada di peta. Ia tidak menuntut sebaliknya — peta boleh
	// memuat kolom yang layar belum baca, dan itu keadaan yang sah.
	{"CEDINGID", "CedingID"},
	{"LEADINGREINSSOURCEID", "LeadingReinsSourceID"},
	{"LEADINGREINSID", "LeadingReinsID"},
	{"INFORMATION", "Information"},
	{"POSITION", "Position"},
	{"POSITIONUSERNAME", "PositionUsername"},
	{"CHOOSESTATUSAKSEPTASI", "ChooseStatusAkseptasi"},
	{"ACCUMULATIONPERIOD", "AccumulationPeriod"},
	// ⚠️ Kolomnya `COMMENTTEKS` — `COMMENT` kata tercadang Oracle. Ejaan
	// dokumennya tetap `Comment`, dan ITULAH kunci yang layar pakai.
	{"COMMENTTEKS", "Comment"},
	// ⭐ Ketujuh medan kepala tab Reporting Period — yang selama ini kosong
	// di layar pada 1.219 dari 1.855 kontrak yang sungguh punya nilainya
	// (`PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §17, kini ditutup).
	{"REPORTINGSTART", "ReportingStart"},
	{"REPORTINGEND", "ReportingEnd"},
	{"REPORTINGPERIOD", "ReportingPeriod"},
	{"REPORTINGINTERVAL", "ReportingInterval"},
	{"REPORTINGSUBMISSION", "ReportingSubmission"},
	{"REPORTINGCONFIRMATION", "ReportingConfirmation"},
	{"REPORTINGSETTLEMENT", "ReportingSettlement"},

	// ⭐ ENAM kolom migrasi `448`, 7 Oktober 2026 — properti yang selama ini
	// hanya hidup di penampung halaman (`halaman.tsx`): tab Share Prop
	// (`RNMShareP`, `BrokeragePercentP`, `OptionLimit`), tab Installment
	// (`InstallmentNo`), dan jalur revisi (`RevisionState`, `ViewState`).
	// ⚠️ Kolomnya boleh belum terpasang — pembaca menyaring lewat
	// `kolomTerpasang`, jadi kontrak tetap terbuka sebelum migrasinya jalan.
	{"RNMSHAREP", "RNMShareP"},
	{"BROKERAGEPERCENTP", "BrokeragePercentP"},
	{"OPTIONLIMIT", "OptionLimit"},
	{"INSTALLMENTNO", "InstallmentNo"},
	{"REVISIONSTATE", "RevisionState"},
	{"VIEWSTATE", "ViewState"},
}

// Kelima ejaan tab teks panjang — `CLOB` di DDL, teks di sini.
var kolomTeksRevisi = [][2]string{
	{"EXCLUSIONSP", "ExclusionsP"},
	{"EXCLUSIONS", "Exclusions"},
	{"SPECIALCONDITIONSP", "SpecialConditionsP"},
	{"SPECIALCONDITIONS", "SpecialConditions"},
	{"SPECIALCONDITIONSLC", "SpecialConditionsp"},
}

// RevisiPendaratan adalah satu baris `T_TREATY_REVISION` yang sudah dibaca.
type RevisiPendaratan struct {
	// Nilai per EJAAN DOKUMEN. Kunci yang `NULL` di tabel tidak masuk.
	Medan map[string]string
	// Teks panjang per ejaan, aturan yang sama.
	Teks map[string]string
	// Ada menyatakan barisnya ketemu. Salah berarti kontrak ini belum
	// didaratkan — bukan berarti medannya kosong.
	Ada bool
}

// BacaRevisiPendaratan membaca medan kepala satu kontrak.
//
// ⚠️ Baris yang TIDAK ADA bukan galat. Kontrak yang belum didaratkan tetap
// harus TERBUKA — sepuluh medan kolom `TREATY_IN` masih terbaca, dan layar
// yang menolak membuka kontrak karena pendaratannya belum jalan
// menyembunyikan justru kontrak yang paling perlu dilihat.
func (g *Gudang) BacaRevisiPendaratan(ctx context.Context, masterID string) (RevisiPendaratan, error) {
	nama, err := g.db.Qualify("T_TREATY_REVISION")
	if err != nil {
		return RevisiPendaratan{}, err
	}
	// ⭐ Hanya kolom yang SUDAH terpasang — daftar di atas boleh mendahului
	// migrasinya (`448`). Tabel yang belum ada = kontrak belum didaratkan.
	terpasang, err := g.kolomTerpasang(ctx, "T_TREATY_REVISION")
	if err != nil {
		return RevisiPendaratan{}, err
	}
	kosong := RevisiPendaratan{Medan: map[string]string{}, Teks: map[string]string{}}
	if terpasang == nil {
		return kosong, nil
	}
	var medan, teks [][2]string
	for _, p := range kolomRevisi {
		if terpasang[p[0]] {
			medan = append(medan, p)
		}
	}
	for _, p := range kolomTeksRevisi {
		if terpasang[p[0]] {
			teks = append(teks, p)
		}
	}
	kolom := make([]string, 0, len(medan)+len(teks))
	for _, p := range medan {
		kolom = append(kolom, p[0])
	}
	for _, p := range teks {
		kolom = append(kolom, p[0])
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1", strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return RevisiPendaratan{}, err
	}
	sel := make([]sql.NullString, len(kolom))
	tuju := make([]any, len(kolom))
	for i := range sel {
		tuju[i] = &sel[i]
	}
	err = g.db.QueryRowContext(ctx, q, masterID).Scan(tuju...)
	if err == sql.ErrNoRows {
		return kosong, nil
	}
	if err != nil {
		return RevisiPendaratan{}, fmt.Errorf("repository: membaca T_TREATY_REVISION kontrak %s: %w", masterID, err)
	}

	r := RevisiPendaratan{Medan: map[string]string{}, Teks: map[string]string{}, Ada: true}
	for i, p := range medan {
		// ⛔ `NULL` dilewati, bukan dimasukkan sebagai teks kosong: itulah
		// yang membedakan "tidak ada di sistem lama" dari "belum diisi".
		if sel[i].Valid {
			r.Medan[p[1]] = sel[i].String
		}
	}
	for j, p := range teks {
		if s := sel[len(medan)+j]; s.Valid {
			r.Teks[p[1]] = s.String
		}
	}
	return r, nil
}

// TabelKursTahunan - master kurs warisan. Dibaca, tidak pernah ditulis.
const TabelKursTahunan = "TREATYEXCHANGEYEARLY"

// BacaKursTahunan membaca grid Rate of Exchange satu kontrak.
//
// ⚠️ DISARING MENURUT `TREATYYEAR`, dan itu BACAAN SAYA atas bentuk
// tabelnya — bukan salinan aturan Pega. Yang ekspor tunjukkan
// (`RDBList/GetCurrencyToIDR_SQL.xml`) menyaring menurut `IDCURRENCY` lalu
// mengambil `STARTDATE` terbaru, tanpa menyebut tahun:
//
//	SELECT TOIDR FROM (SELECT * FROM treatyexchangeyearly
//	  WHERE IDCURRENCY = :cur ORDER BY startdate DESC) WHERE rownum = 1
//
// Tetapi tabelnya TERBAGI PER TAHUN — 140 baris, 25 mata uang, 6–25 baris
// per `TREATYYEAR` — dan grid yang tidak menyaring tahun akan menampilkan
// kurs tahun lain pada kontrak lama. `TREATY_IN.TREATYYEAR` ada dan sudah
// terbaca, jadi saringannya memakai itu.
//
// ⛔ Kalau pemilik proses menghendaki aturan ekspor APA ADANYA — kurs
// terbaru tanpa memandang tahun kontrak — yang berubah hanya klausa `WHERE`
// di bawah. Bedanya nyata: kontrak 2023 akan menampilkan kurs 2026.
func (g *Gudang) BacaKursTahunan(ctx context.Context, tahunTreaty string) ([]models.BarisKursWarisan, error) {
	// Kontrak tanpa `TREATYYEAR` tidak punya tahun untuk disaring; nol baris
	// adalah jawaban yang jujur, dan bukan galat.
	if strings.TrimSpace(tahunTreaty) == "" {
		return []models.BarisKursWarisan{}, nil
	}
	nama, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID, CURRENCY, IDCURRENCY, TOIDR, STARTDATE, ENDDATE
		FROM %s WHERE TREATYYEAR = :1 ORDER BY CURRENCY`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, tahunTreaty)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s tahun %s: %w", TabelKursTahunan, tahunTreaty, err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.BarisKursWarisan{}
	for rows.Next() {
		var id, cur, idCur, toidr, mulai, akhir sql.NullString
		if err := rows.Scan(&id, &cur, &idCur, &toidr, &mulai, &akhir); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelKursTahunan, err)
		}
		out = append(out, models.BarisKursWarisan{
			ID:            id.String,
			MataUang:      cur.String,
			MataUangID:    idCur.String,
			NilaiKeIDR:    toidr.String,
			BerlakuDari:   mulai.String,
			BerlakuSampai: akhir.String,
		})
	}
	return out, rows.Err()
}

// KolomRevisiUntukUji membuka daftar kolom pembaca untuk penjaga di paket
// uji — `TestKolomRevisiAdaDiPeta`.
//
// ⛔ Hanya untuk uji, dan itu disengaja terlihat dari namanya. Daftarnya
// tetap tidak dapat diubah dari luar: yang dikembalikan salinan.
func KolomRevisiUntukUji() [][2]string {
	out := make([][2]string, 0, len(kolomRevisi)+len(kolomTeksRevisi))
	out = append(out, kolomRevisi...)
	out = append(out, kolomTeksRevisi...)
	return out
}
