# Paritas layar dan aksi — Treaty Contract Out

> Disusun 28-09-2026 (sesi modul, tiket 01 →). Satu baris per tombol/aksi korpus → rute/kontrol sistem
> baru. Nomor baris = `sed -e 's/></>\n</g'` atas berkas korpus `D:\XML\RNM_BRD\Treaty Contract Out\`.
> Keadaan: ✅ dibangun · ⚠️ belum · ➖ sengaja tidak dibawa (kode mati / penyimpangan sadar) · 🔜 tiket berikut.

## Menu — tiga harness portal (kelas `Data-Portal`)

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Harness/InboxTreatyContract.xml` b151 `<pyLabel>InboxTreatyContract</pyLabel>` → section `GridTreatyContract` (b931) → `InputTreatyContract` (`GridTreatyContract.xml` b2482) | sidebar **Treaty Contract Out → InboxTreatyContract** | 🔜 tiket 03 |
| `Harness/InboxTreatyContractReinsType.xml` b151 `<pyLabel>InboxTreatyContractReinsType</pyLabel>` → judul `ReinsType` b1670 → section `PanggilReinsType` (b1825) → `InputTreatyContractReinsType` (`PanggilReinsType.xml` b1300) | sidebar **Treaty Contract Out → InboxTreatyContractReinsType** | 🔜 tiket 04 |
| `Harness/InboxTreatyContractDescription.xml` b359 `<pyLabel>InboxTreatyContractDescription</pyLabel>` → `NitipKurs` b3882, `SubViewDetailDescription` b11366, `ViewDetailDescriptionProp` b12583, `ViewDetailDescriptionNonProp` b13547 | sidebar **Treaty Contract Out → InboxTreatyContractDescription** | 🔜 tiket 08 |

⚠️ Kelompok **Treaty Contract Out** BELUM ada di sidebar (`labels.ts` `MODUL` memuat 17 kelompok; folder
korpus 20 — ralat §8 `PROMPT-EKSEKUSI-HULU-HILIR.md`). Ditambahkan ADITIF bersama layar pertama (tiket 03).

## Tiket 01 — skema + migrasi data (tanpa layar)

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `RDBList/SaveMasterTreatyYear_SQL.xml` → `POOLDATA.PEGA_TREATYYEAR` (10 param) | tabel `T_TREATYYEAR` (migrasi 300), sequence `SEQ_T_TREATYYEAR` | ✅ skema; penulisnya 🔜 tiket 03 |
| `RDBList/SaveMasterTreatyContract_SQL.xml` → `PEGA_TREATYCONTRACT` (8 param) | `T_TREATYCONTRACT` (301) | ✅ skema; 🔜 tiket 04 |
| `RDBList/SaveMasterTreatyReinsurer_SQL.xml` → `PEGA_TREATYREINSURER` (19 param) | `T_TREATYREINSURER` (302) | ✅ skema; 🔜 tiket 05 |
| `RDBList/InsertToMTreatySecurity.xml` INSERT posisional 7 nilai, `UpdateMTreatySecurity.xml` kunci `trim()` | `T_MTREATYSECURITY` (303) PK surrogate, kolom bernama, tanpa `trim()` | ✅ skema; 🔜 tiket 06 |
| `RDBList/SaveMasterTreatyBusiness_SQL.xml` → `PEGA_TREATYBUSINESS` (12 param) | `T_TREATYBUSINESS` (304) | ✅ skema; 🔜 tiket 07 |
| `RDBList/SaveMasterProportionalArrg.xml` (35) + `SaveMasterProportionalArrgChild.xml` (26) → satu tabel | `T_PROPORTIONALARRG` (305) — 35 kolom, anak NULL di 9 kolom induk | ✅ skema; 🔜 tiket 08 |
| — (tidak ada di korpus; ADR-0007) | `T_TREATYCO_JEJAK` (306) | ✅ skema; penulis 🔜 tiket 03+ |
| enam tabel warisan `POOLDATA.TREATYYEAR` … `PROPORTIONALARRG` | `-migrate-data-treaty-contract-out` → `repository.MigrasiTCO.Pindahkan` (baca teks → konversi → tulis → rekonsiliasi tepat → sequence) | ✅ kode + uji `db` (SKIP tanpa skema uji); ⛔ tidak dijalankan di DEV (tco2) |
| kueri hilir `Claim Prop/RDBList/GetLimitPLATreatyin.xml`, `GetListRetro_Sql.xml`, `GetTreatyGroupID.xml`; `Claim Fac In/RDBList/GetLimitPLADLA_Sql.xml`, `GetQuotaShare.xml`, `GetTreatyGroup_Sql.xml`, `GetDataTreatyLimit_Sql.xml`; `Komite Claim Prop/RDBList/GetListRetro_Sql.xml` | `repository.KontrakHilirTCO` — 6 pembaca read-only berkolom VERBATIM; uji tiga sisi (DDL, korpus, nol tulis) | ✅ tco3 |
| `M_PROPORTIONALARRG`, `M_TREATYCONTRACT`, `M_TREATYYEAR`, `M_TREATYBUSINESS` + seluruh kueri `FROM m_*` | — | ➖ MATI (penyimpangan sadar 1); penjaga `TestTCONolTabelDokumenWarisan` |
| `POOLDATA.PROSESCOPY`, `RDBList/SaveMasterCopyData_SQL.xml`, `Activity/BrowseCopyData.xml`, `SaveTreatyYearMultiple_Act.xml`, `NewInputTreatyYearMultiple_Act.xml`, bagian salin `BrowseDeleteRowTreatyInContract.xml`; form From/To `InputTreatyContract.xml` b2374/b3818, tombol `Proces` b5104 → `BrowseCopyData` b5128 | — | ➖ fitur salin DIBUANG (penyimpangan sadar 3, AC 72) |
| kolom warisan `PROPORTIONALLIST`, `OBJECT` | dicacah laporan migrasi (`KolomMatiBerisi`), tidak dibawa | ➖ AC 70 |

## Tiket 02 — jenis reasuransi: master dibaca + saringan non-life

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` — `.ID NotStartsWith` 12 nilai b573/b581, `.Flag="active"` b586, `.Type` 1/2/3 b603, sort `.Note` b667; dipakai 11 grid klausul | `GET /api/treaty-contract-out/jenis-reasuransi` → `repository.MasterJenisReasuransi.DaftarNonLife` (`FLAG=`, `TYPE IN`, `ID NOT LIKE` ×12, `ORDER BY NOTE`) | ✅ nama jujur (nol "Old") |
| master `POOLDATA.REINSURANCETYPE` (`ID`, `NOTE`, `TYPE`, `FLAG`) | dibaca saja; penjaga `TestTCOWarisanHanyaDibaca` (`masterDibacaSajaTCO`) | ✅ AC 12 |
| master kosong / tersaring habis | 503 `ErrMasterJenisReasuransiKosong` menyebut `REINSURANCETYPE` | ✅ ADR-0015 |
| pemilih `ReinsType` (`InputTreatyContractReinsType.xml` b2652) / `Reinsurance Type` (`InputTreatyContract.xml` b7925) | `components/treaty-contract-out/PilihJenisReasuransi.tsx` (`Pilih` ui/dasar; nilai `.ID` teks, label `.Note`) | ✅ komponen; dipasang di layar 🔜 tiket 03/04/08 |
| RD non-Old `BrowseReinsuranceType_RD.xml` (param `Flag="active"` b2768, pemilih form kontrak) | daftar tersaring yang sama (AC 5 tiket 02) | ⚠️ OQ-TCO-06 |
| `RDBList/GetMasterReinsTypeContract.xml` (`a.JSONDATA`) | — | ➖ MATI |
