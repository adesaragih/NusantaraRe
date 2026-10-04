package repository

// Baca SATU kontrak warisan - `TREATY_IN` berpasangan `M_TREATY_IN`.
//
// ⛔ BACA SAJA, kedua tabel. Nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`, nol DDL,
// nol penyebutan di migrasi. Dijaga `TestWarisanHanyaDibaca`.
//
// ⛔ Jalur ini TERPISAH dari `kontrak.go` yang membaca `KONTRAK` +
// `VERSI_KONTRAK`. Keduanya menjawab pertanyaan yang berbeda, dan
// `GET /kontrak/{id}` tidak disentuh.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// TabelWarisanJSON - pasangan `TREATY_IN` yang memegang dokumen aslinya.
const TabelWarisanJSON = "M_TREATY_IN"

// ErrWarisanTidakAda - pengenalnya tidak menunjuk baris mana pun.
var ErrWarisanTidakAda = errors.New("kontrak warisan tidak ada")

// ErrJSONWarisanRusak - dokumennya ada tetapi tidak dapat diurai.
var ErrJSONWarisanRusak = errors.New("dokumen JSON warisan tidak dapat diurai")

// jsonWarisan adalah LIMA kunci yang form pakai - bukan 175.
//
// ⛔ Sengaja TIDAK memetakan seluruh dokumen. Satu dokumen memuat 175 kunci,
// dan struct 175 medan adalah 170 medan yang tidak ada yang membacanya,
// masing-masing menunggu salah ketik yang tidak akan ketahuan. Kunci
// berikutnya ditambahkan ketika ada yang memakainya, dengan buktinya sendiri.
//
// ⚠️ Penunjuk, BUKAN string. `nil` berarti KUNCINYA TIDAK ADA; string kosong
// berarti kuncinya ada dan nilainya kosong. Keduanya berarti hal yang berbeda
// di layar - "tidak ada di sistem lama" bukan "belum diisi" - dan string
// biasa menghapus bedanya tanpa suara.
type jsonWarisan struct {
	Bordeaux       *string `json:"Bordeaux"`
	BordereauxNote *string `json:"BordereauxNote"`
	AccountingMode *string `json:"AccountingMode"`
	ContractRefNo  *string `json:"ContractRefNo"`
	TreatyLeader   *string `json:"TreatyLeader"`

	// ⭐ TAB TEKS - Jalan B, keputusan pemilik proses 4 Oktober 2026
	// (`docs/KEPUTUSAN-PENYELARASAN-REPO.md` §15): isinya TETAP di
	// `JSONDATA`, nol tabel pendaratan untuknya. Layarnya baca-saja, dan
	// dokumennya sudah ditarik untuk medan di atas - satu medan lagi dari
	// dokumen yang sama nol tambahan perjalanan ke basis data.
	//
	// ⛔ LIMA kunci untuk DUA tab, dan ejaannya BUKAN sinonim. Terukur atas
	// 1.854 dokumen: dari 303 dokumen yang punya lebih dari satu ejaan
	// `SpecialConditions*`, NOL yang isinya identik. Memilih ejaan lain
	// ketika yang sesuai cabang tidak ada akan menampilkan TEKS YANG SALAH,
	// bukan salinan yang basi.
	//
	// ⚠️ `SpecialConditionsp` berhuruf kecil di akhir, dan itu BUKAN salah
	// ketik: 292 dokumen memakainya - 170 proporsional, 122 non. Ia tidak
	// terikat cabang mana pun, jadi ia SELALU "ejaan lain".
	//
	// Nol dari kelima kunci ini pernah bernilai string kosong: bila
	// kuncinya ada, isinya ada.
	ExclusionsP        *string `json:"ExclusionsP"`
	Exclusions         *string `json:"Exclusions"`
	SpecialConditionsP *string `json:"SpecialConditionsP"`
	SpecialConditions  *string `json:"SpecialConditions"`
	SpecialConditionsp *string `json:"SpecialConditionsp"`

	// ⭐ SATU LARIK yang masih dibaca dari dokumen, dan sebabnya diukur.
	//
	// Tiga saudaranya - `ReportingPeriodList`, `Portfolio`,
	// `AccumulationList` - PINDAH ke tabel pendaratan migrasi 430 pada
	// 3 Oktober 2026, dan dibaca `pendaratan_baca.go`.
	// `TestLarikYangSudahPunyaTabelTidakDiuraiLagi` menolak kembalinya.
	//
	// ⛔ `CurrencyList` TINGGAL, dan itu bukan kelalaian. Ronde pemindahan
	// melarang membuat tabel untuk Rate of Exchange, atas dasar
	// `TREATYEXCHANGEYEARLY` sudah melayaninya. Sapuan menemukan dasar itu
	// KELIRU: tabel tersebut 140 baris per TAHUN dan nol kolomnya menunjuk
	// kontrak, sementara `CurrencyList` 3.465 elemen di 1.844 kontrak.
	// Jadi tab ini tetap membaca dokumen - sampai tabel kesembilan
	// diputuskan. Uraiannya di
	// `4-erd-dan-tabel-datar/KOREKSI-ERD-VERSUS-POOLDATA.md` §12.
	//
	// ⚠️ Seluruh nilainya `string` di JSON - terukur atas seluruh 1.854
	// dokumen, nol angka dan nol boolean.
	CurrencyList []jsonKurs `json:"CurrencyList"`
}

// jsonKurs - satu baris grid Rate of Exchange. Keempat kolom layar lama
// (`Currency`, `Value to IDR`, `Valid From`, `Valid Until`) ada di sini.
type jsonKurs struct {
	Currency    *string `json:"Currency"`
	Conversion  *string `json:"Conversion"`
	PeriodStart *string `json:"PeriodStart"`
	PeriodEnd   *string `json:"PeriodEnd"`
}

// teks membaca penunjuk yang boleh nil menjadi string kosong.
func teks(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// BacaKontrakWarisan membaca satu kontrak beserta dokumen aslinya.
//
// ⚠️ `LEFT JOIN` ke `M_TREATY_IN`: sapuan menemukan 1.854 lawan 1.854 dan
// nol dokumen NULL hari ini, tetapi kontrak yang kehilangan dokumennya harus
// tetap TERBUKA - delapan medan kolomnya masih dapat dibaca, dan layar yang
// menolak membuka kontrak karena dokumennya hilang menyembunyikan justru
// kontrak yang paling perlu dilihat.
func (g *Gudang) BacaKontrakWarisan(ctx context.Context, id string) (models.KontrakWarisan, error) {
	tKontrak, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	tJSON, err := g.db.Qualify(TabelWarisanJSON)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	q := fmt.Sprintf(`SELECT t.ID, t.TREATYCONTRACTNAME, t.TERITORIALSCOPE, t.TREATYYEAR,
		t.CEDING, t.CEDINGID, t.LEADINGREINSSOURCE, t.LEADINGREINSSOURCEID,
		t.PROPORTIONTYPE, t.COMMENCEMENT, t.TERMINATION, m.JSONDATA
		FROM %s t LEFT JOIN %s m ON m.ID = t.ID
		WHERE t.ID = :1`, tKontrak, tJSON)

	var k models.KontrakWarisan
	var sid, nk, ts, ty, cd, cdid, lrs, lrsid, pt, cm, tm, dok sql.NullString
	err = g.db.QueryRowContext(ctx, q, id).Scan(&sid, &nk, &ts, &ty, &cd, &cdid,
		&lrs, &lrsid, &pt, &cm, &tm, &dok)
	if errors.Is(err, sql.ErrNoRows) {
		return models.KontrakWarisan{}, fmt.Errorf("%w: %s", ErrWarisanTidakAda, id)
	}
	if err != nil {
		return models.KontrakWarisan{}, fmt.Errorf("repository: membaca kontrak warisan %s: %w", id, err)
	}

	k.ID = sid.String
	k.NamaKontrak = nk.String
	// ⭐ `TeritorialScope` diambil dari KOLOM, bukan dari JSON. Keduanya ada
	// dan 1.844 dari 1.854 identik begitu `\n` di-unescape; kolomnya sudah
	// berbentuk teks siap tampil, nilai JSON-nya belum.
	k.LingkupWilayah = ts.String
	k.TahunTreaty = ty.String
	k.Cedant = cd.String
	k.IDCedant = cdid.String
	k.AsalBisnis = lrs.String
	k.IDAsalBisnis = lrsid.String
	k.SifatProporsiAsli = pt.String
	k.TanggalMulaiAsli = cm.String
	k.TanggalBerakhirAsli = tm.String
	k.AdaDiJSON = map[string]bool{}

	if !dok.Valid || dok.String == "" {
		// Dokumennya tidak ada - bukan galat. Kedelapan medan kolom sudah
		// terisi, dan layar menyatakan sisanya tidak tersedia.
		return k, nil
	}
	var j jsonWarisan
	if err := json.Unmarshal([]byte(dok.String), &j); err != nil {
		// ⛔ Galat yang MENYEBUT KONTRAKNYA. "invalid character" tanpa
		// pengenal membuat yang menyelidiki memeriksa 1.854 dokumen.
		return models.KontrakWarisan{}, fmt.Errorf("%w: kontrak %s: %v", ErrJSONWarisanRusak, id, err)
	}
	ambil := func(nama string, p *string, ke *string) {
		if p == nil {
			return
		}
		k.AdaDiJSON[nama] = true
		*ke = *p
	}
	ambil("Bordeaux", j.Bordeaux, &k.Bordereaux)
	ambil("BordereauxNote", j.BordereauxNote, &k.BordereauxCatatan)
	ambil("AccountingMode", j.AccountingMode, &k.CaraPembukuan)
	ambil("ContractRefNo", j.ContractRefNo, &k.NomorRujukan)
	ambil("TreatyLeader", j.TreatyLeader, &k.PemimpinTreaty)

	// ⛔ Kelima kunci teks dibawa MENTAH, lengkap dengan ejaan mana yang
	// ada. Yang MEMILIH di antaranya adalah services, sebab pemilihannya
	// bergantung pada cabang - dan cabang bukan urusan lapisan baca.
	k.TeksMentah = map[string]string{}
	for _, pasang := range []struct {
		nama string
		p    *string
	}{
		{"ExclusionsP", j.ExclusionsP},
		{"Exclusions", j.Exclusions},
		{"SpecialConditionsP", j.SpecialConditionsP},
		{"SpecialConditions", j.SpecialConditions},
		{"SpecialConditionsp", j.SpecialConditionsp},
	} {
		if pasang.p == nil {
			continue
		}
		k.AdaDiJSON[pasang.nama] = true
		k.TeksMentah[pasang.nama] = *pasang.p
	}

	// ⭐ Larik diterjemahkan APA ADANYA - nol penyusunan ulang, nol pengurutan,
	// nol penyaringan. Urutan di dokumen adalah urutan yang layar lama
	// tampilkan, dan mengurutkannya di sini mengubah arti tanpa ada yang
	// memintanya.
	for _, b := range j.CurrencyList {
		k.Kurs = append(k.Kurs, models.BarisKursWarisan{
			MataUang:      teks(b.Currency),
			NilaiKeIDR:    teks(b.Conversion),
			BerlakuDari:   teks(b.PeriodStart),
			BerlakuSampai: teks(b.PeriodEnd),
		})
	}
	return k, nil
}
