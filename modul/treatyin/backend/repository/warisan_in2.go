package repository

// Baca `POOLDATA.M_TREATY_IN2` — empat tab dari satu tabel warisan.
//
// ⛔ BACA SAJA. Nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`, nol DDL, nol
// penyebutan di migrasi mana pun. Dijaga `TestWarisanHanyaDibaca`.
//
// ⛔ NOL TABEL BARU untuk Limits · Share · Event Limits · RNM Share. Tabelnya
// sudah ada, berisi, dan terhitung per layer; mendaratkan ulang isinya akan
// membuat dua sumber untuk satu grid.
//
// ⚠️ Jangkauannya TIDAK penuh, dan itu harus terbaca di layar, bukan
// disembunyikan. Terukur 3 Oktober 2026:
//
//	M_TREATY_IN2                     7.281 baris · 1.340 kontrak
//	(MASTERID, LAYER) berbeda        2.540
//	dokumen dengan `Limits[]` berisi  1.850 kontrak
//	`Limits[]` berisi TAPI nol baris di M_TREATY_IN2  →  510 kontrak
//	                                                     1.210 elemen
//
// Jadi grid yang kosong di sini TIDAK selalu berarti "kontrak ini memang
// tidak punya limit". Pada 510 kontrak ia berarti "sumber yang tabel ini
// ambil tidak mencakupnya". Bedanya dinyatakan di layar lewat petunjuk
// kosong yang berbeda, dan dicatat di `docs/PEMETAAN-M-TREATY-IN2.md`.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// TabelWarisanLayer - tabel per-layer sistem lama.
const TabelWarisanLayer = "M_TREATY_IN2"

// kolomLayer adalah ke-41 kolom, DALAM URUTAN yang `BacaLayerWarisan`
// pindai. Urutannya harus sama persis dengan urutan medan di `pindaiLayer`;
// `TestKolomLayerSejajarDenganPemindai` mengadu keduanya.
var kolomLayer = []string{
	"MASTERID", "PROPORTIONTYPE", "CEDINGID", "CEDING", "SOBID", "SOB",
	"TREATYGROUP", "TREATYCONTRACTNAME", "COMMENCEMENT", "TERMINATION",
	"TREATYTYPE", "CESSIONPCT", "CEDANT_RETENTION", "BASIS_COVER",
	"LAYERTYPE", "LAYER", "SPREADINGTYPE", "CURRENCY", "LIMIT_100",
	"ADJ_RATE", "EARN_PREMIUM", "MDP_RATIO", "MDP", "ROL",
	"CURRENCYRELATION", "CESSION_TO_RI", "EPI100", "RIOGR",
	"BROKERAGEPERCENTP", "EARTHQUAKE", "RNMSHARE", "CURRENCYLIMIT",
	"LIABILITY_RNM", "MDP_RNM_100", "QSOR", "QSRI", "LIABILITYQSRI",
	"LIABILITYQSOR", "EPIRNMQS100", "RNM_RETAINED_PREMI", "RNM_QS_PREMI",
}

// KolomLayerWarisan membuka daftar kolom untuk uji — ke-41 nama apa adanya.
func KolomLayerWarisan() []string { return append([]string(nil), kolomLayer...) }

// BacaLayerWarisan membaca seluruh baris layer satu kontrak.
//
// ⛔ `ORDER BY` yang BERLAPIS, dan tiap lapisnya ada sebabnya:
//
//   - `TO_NUMBER(LAYER …)` lebih dulu, sebab `LAYER` kolom TEKS dan urutan
//     teks menaruh layer 10 di antara 1 dan 2. Nilai seperti `1A` ada
//     (terukur di `Limits[].Layer`), jadi `DEFAULT 0 ON CONVERSION ERROR`
//     menjaganya tidak meledakkan kueri — ia mendarat di depan, terlihat.
//   - `LAYER` apa adanya sesudahnya, supaya `1` dan `1A` punya urutan tetap.
//   - `LAYERTYPE` terakhir: satu layer punya baris `layer` dan `sublayer`,
//     dan tanpa lapis ini keduanya bertukar tempat antar pembacaan.
//
// Tanpa `ORDER BY` Oracle bebas mengembalikan baris dalam urutan apa pun,
// dan grid limit yang layernya tertukar TERBACA BENAR.
func (g *Gudang) BacaLayerWarisan(ctx context.Context, masterID string) ([]models.BarisLayerWarisan, error) {
	nama, err := g.db.Qualify(TabelWarisanLayer)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE MASTERID = :1
		ORDER BY TO_NUMBER(LAYER DEFAULT 0 ON CONVERSION ERROR), LAYER, LAYERTYPE`,
		strings.Join(kolomLayer, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}

	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelWarisanLayer, masterID, err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.BarisLayerWarisan{}
	for rows.Next() {
		sel := make([]sql.NullString, len(kolomLayer))
		tuju := make([]any, len(kolomLayer))
		for i := range sel {
			tuju[i] = &sel[i]
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s kontrak %s: %w", TabelWarisanLayer, masterID, err)
		}
		out = append(out, pindaiLayer(sel))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelWarisanLayer, masterID, err)
	}
	return out, nil
}

// pindaiLayer menyusun satu baris. Urutannya WAJIB sama dengan `kolomLayer`.
//
// ⚠️ Ditulis mendatar dan bernomor, bukan lewat refleksi: pergeseran satu
// medan di sini memindahkan SELURUH nilai sesudahnya ke kolom tetangganya,
// dan hasilnya tetap berupa grid yang terisi rapi. Nomor di komentar adalah
// indeks di `kolomLayer`, supaya keduanya dapat diadu dengan mata.
func pindaiLayer(s []sql.NullString) models.BarisLayerWarisan {
	return models.BarisLayerWarisan{
		MasterID:        s[0].String, // MASTERID
		SifatProporsi:   s[1].String, // PROPORTIONTYPE
		IDCedant:        s[2].String, // CEDINGID
		Cedant:          s[3].String, // CEDING
		IDAsalBisnis:    s[4].String, // SOBID
		AsalBisnis:      s[5].String, // SOB
		KelompokTreaty:  s[6].String, // TREATYGROUP
		NamaKontrak:     s[7].String, // TREATYCONTRACTNAME
		TanggalMulai:    s[8].String, // COMMENCEMENT
		TanggalBerakhir: s[9].String, // TERMINATION

		JenisTreaty:     s[10].String, // TREATYTYPE
		PersenCession:   s[11].String, // CESSIONPCT
		RetensiCedant:   s[12].String, // CEDANT_RETENTION
		DasarCover:      s[13].String, // BASIS_COVER
		JenisLayer:      s[14].String, // LAYERTYPE
		Layer:           s[15].String, // LAYER
		JenisPenyebaran: s[16].String, // SPREADINGTYPE
		MataUang:        s[17].String, // CURRENCY
		Limit100:        s[18].String, // LIMIT_100
		AdjRate:         s[19].String, // ADJ_RATE
		PremiEarned:     s[20].String, // EARN_PREMIUM
		RasioMDP:        s[21].String, // MDP_RATIO
		MDP:             s[22].String, // MDP
		ROL:             s[23].String, // ROL
		RelasiMataUang:  s[24].String, // CURRENCYRELATION
		CessionKeRI:     s[25].String, // CESSION_TO_RI
		EPI100:          s[26].String, // EPI100
		RIOGR:           s[27].String, // RIOGR
		PersenBrokerage: s[28].String, // BROKERAGEPERCENTP
		Gempa:           s[29].String, // EARTHQUAKE
		RNMShare:        s[30].String, // RNMSHARE
		MataUangLimit:   s[31].String, // CURRENCYLIMIT

		LiabilityRNM:     s[32].String, // LIABILITY_RNM
		MDPRNM100:        s[33].String, // MDP_RNM_100
		QSOR:             s[34].String, // QSOR
		QSRI:             s[35].String, // QSRI
		LiabilityQSRI:    s[36].String, // LIABILITYQSRI
		LiabilityQSOR:    s[37].String, // LIABILITYQSOR
		EPIRNMQS100:      s[38].String, // EPIRNMQS100
		RNMRetainedPremi: s[39].String, // RNM_RETAINED_PREMI
		RNMQSPremi:       s[40].String, // RNM_QS_PREMI
	}
}
