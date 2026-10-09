package repository

// Jalur baca isi pemilih "Choose Ceding" dan "Choose Source of Business".
//
// ⛔ BACA SAJA atas `AGENT`. Nol `INSERT`, nol `UPDATE`, nol `DELETE`, nol
// DDL.
//
// ⛔ RALAT 6 Oktober 2026: berkas ini dahulu membaca pasangan `DISTINCT`
// dari `TREATY_IN` — nilai yang PERNAH dipakai kontrak. Kini ia membaca
// daftar yang Pega sendiri tawarkan (`BrowseAgentNusaRe_RD`). Sebab dan
// pengukurannya di bawah, pada `BacaDaftarCedant`.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// ---------------------------------------------------------------------
// ⛔ SUMBERNYA DIGANTI 6 Oktober 2026 — dari "nilai yang pernah dipakai"
// ke daftar agen yang Pega sendiri tawarkan.
// ---------------------------------------------------------------------
// Bentuk sebelumnya menawarkan pasangan `DISTINCT` dari `TREATY_IN` — 131
// cedant dan 96 asal bisnis yang PERNAH tercatat di kontrak. Itu bukan
// yang Pega tawarkan, dan bedanya dua arah: agen aktif yang belum pernah
// dipakai tidak dapat dipilih, dan agen yang sudah tidak aktif tetap
// ditawarkan.
//
// ⭐ YANG PEGA TAWARKAN, dibaca dari ekspor:
//
//	Section/TreatyInSearchReinsured.xml + TreatyInSearchSoB.xml
//	  -> ReportDefinition/BrowseAgentNusaRe_RD.xml, kelas ASM-FW-GISFW-Int-AGENT
//	     A  .StatusActive = 1
//	     B  .AgentType2   != "LIFE INSURANCE"
//	     D  .ChildCount   = Param.ChildCount   (seksi mengirim 0)
//	     E  .ClientName   Contains Param.ClientName   (kata kunci, tanpa
//	                      membedakan huruf — disaring di layar)
//	     G  .ClientID     IS NOT NULL
//	     logika: A AND B AND D AND E AND G
//
// ⚠️ KEDUA pemilih memakai RD yang SAMA dengan parameter yang SAMA, jadi
// keduanya menawarkan himpunan yang sama. Itu bunyi ekspornya, bukan
// penyederhanaan di sini.
//
// ⭐ PENGENAL YANG DITULIS adalah `AGENT.ID`, bukan `CLIENTID` — ruang
// pengenal itu DIUKUR, bukan diduga. Atas seluruh 1.854 baris `TREATY_IN`:
// `LEADINGREINSSOURCEID` = `AGENT.ID` pada 1.854, `CEDINGID` = `AGENT.ID`
// pada 1.744 (110 sisanya kebetulan cocok `CLIENTID` — data lama). Dan
// SQL mentah Pega sendiri memanggil `agent where id = <pengenal ceding>`
// (`Claim Non Prop/RDBList/GetAddressCeding.xml`).
//
// ⚠️ JEBAKAN EJAAN, dan ia sudah memakan satu kueri: properti Pega
// `.AgentType2` berkolom `AGENTTPYE2` di basis data — `TPYE`, bukan `TYPE`.
// Kueri yang memakai ejaan properti gagal `ORA-00904`.
//
// Terukur 6 Oktober 2026 atas seluruh `POOLDATA.AGENT` (429 baris): 255
// lolos saringan, 132 dibuang sebagai `LIFE INSURANCE`, 4 nama dipakai
// lebih dari satu ID (8 baris), 0 ID ganda. Penandaan kembar dikerjakan
// `services.pilihanWarisan`, seperti setiap pemilih lain.

// BacaDaftarCedant mengembalikan isi pemilih "Choose Ceding".
func (g *Gudang) BacaDaftarCedant(ctx context.Context) ([]models.PilihanWarisan, error) {
	return g.bacaAgenAktif(ctx)
}

// BacaDaftarAsalBisnis mengembalikan isi pemilih "Choose Source of Business".
//
// ⚠️ Sama persis dengan `BacaDaftarCedant` — `TreatyInSearchSoB.xml` memanggil
// RD dan parameter yang sama dengan `TreatyInSearchReinsured.xml`.
func (g *Gudang) BacaDaftarAsalBisnis(ctx context.Context) ([]models.PilihanWarisan, error) {
	return g.bacaAgenAktif(ctx)
}

// TabelAgen adalah tabel di balik kelas Pega `ASM-FW-GISFW-Int-AGENT`.
const TabelAgen = "AGENT"

// SaringanAgenPega adalah klausa A·B·D·G `BrowseAgentNusaRe_RD` dalam ejaan
// kolom Oracle. E (kata kunci) dikerjakan layar.
//
// ⛔ Diekspor supaya uji db mengadu kueri yang SAMA, bukan salinannya —
// salinan yang menyimpang diam-diam membuktikan kueri yang tidak dijalankan
// siapa pun.
const SaringanAgenPega = `STATUSACTIVE = '1' AND NVL(AGENTTPYE2, '~') <> 'LIFE INSURANCE' ` +
	`AND CHILDCOUNT = '0' AND CLIENTID IS NOT NULL`

// bacaAgenAktif membaca agen yang Pega tawarkan, urut nama lalu pengenal.
//
// ⚠️ `NVL(AGENTTPYE2, '~')`: `!=` Pega atas nilai kosong MELOLOSKAN baris,
// sedangkan `<>` Oracle atas `NULL` MENOLAKNYA. Tanpa `NVL` agen yang
// jenisnya kosong hilang dari pemilih tanpa satu pun galat.
func (g *Gudang) bacaAgenAktif(ctx context.Context) ([]models.PilihanWarisan, error) {
	nama, err := g.db.Qualify(TabelAgen)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(ID), CLIENTNAME FROM %s
		WHERE %s
		ORDER BY CLIENTNAME, ID`, nama, SaringanAgenPega)
	baris, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca agen dari %s: %w", TabelAgen, err)
	}
	defer func() { _ = baris.Close() }()

	keluar := []models.PilihanWarisan{}
	for baris.Next() {
		var id, nm sql.NullString
		if err := baris.Scan(&id, &nm); err != nil {
			return nil, fmt.Errorf("repository: membaca baris agen: %w", err)
		}
		keluar = append(keluar, models.PilihanWarisan{ID: id.String, Nama: nm.String})
	}
	return keluar, baris.Err()
}
