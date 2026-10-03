# 02: Jenis reasuransi — master dibaca + saringan non-life

**Status:** selesai (28-09-2026)

**Blocked by:** 01 (skema harus ada)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin memilih **jenis reasuransi** dari daftar master — dan
daftar itu hanya memuat jenis yang **berlaku untuk non-life** — sehingga saya tidak pernah mengetik
bebas dan tidak pernah salah memilih jenis milik lini lain. *(User story 8–10 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Baca master jenis reasuransi (read-only) |
| `internal/services` | Saringan non-life |
| `internal/handlers` | Endpoint daftar jenis reasuransi |
| `frontend/` | Pemilih jenis reasuransi yang dipakai layar kontrak dan seluruh grid klausul |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| **RD dominan** | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` | ⚠️ dipakai **11 grid** |
| RD kedua | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseReinsuranceType_RD.xml` | dipakai 1 layar — **layar master**, di luar konteks ini |
| `GetMasterReinsTypeContract` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterReinsTypeContract.xml` | ⚠️ membaca tabel **JSON** — mati |

⚠️ **Nama berbohong — dan yang "Old" justru yang benar** `[terverifikasi]`:

| Rule | RuleSet | Commit | Dipakai |
| --- | --- | --- | ---: |
| `…!BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` | **01-01-91** | **2026-02-05** | **11** grid |
| `…!BROWSEREINSURANCETYPE_RD` | 01-01-83 | 2025-06-02 | 1 layar |

`[keputusan work owner]` Yang **dimigrasikan** adalah perilaku yang dipakai 11 grid:

```
.ID NOT IN ("10004","10011","10012","10021","10022","10025","10026","10028",
            "10248","10249","10018","10217")
AND .Flag = "active"
AND .Type IN ("1","2","3")
```

⚠️ **Penyimpangan sadar 8 — nama jujur.** `[keputusan work owner]` Komponen di sistem baru
**tidak** mengandung kata `"Old"`.

`[terverifikasi]` Class `ASM-FW-GISFW-INT-REINSURANCETYPE` adalah **class yang sama** dengan master
jenis reasuransi Life (`Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml`).
Jadi master ini **melayani life dan non-life sekaligus**; saringanlah yang memisahkan.

## ADR terkait

**ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] Jenis reasuransi dipilih dari **master**, tidak pernah diketik bebas. *(AC 10 spec)*
- [ ] ⚠️ Daftar disaring **persis** seperti existing: **dua belas ID di-blacklist** (`10004`,
      `10011`, `10012`, `10021`, `10022`, `10025`, `10026`, `10028`, `10248`, `10249`, `10018`,
      `10217`) **dan** `Flag = "active"` **dan** `Type` termasuk `1`, `2`, `3`.
      *(AC 11 spec; `[keputusan work owner]` — ditiru apa adanya, **jangan digeneralkan**)*
- [ ] Master jenis reasuransi **tidak ditulis** oleh konteks ini. Test yang menemukan tulisan ke
      master itu **gagal**. *(AC 12 spec)*
- [ ] ⚠️ Nama komponen, endpoint, dan fungsi di sistem baru **tidak mengandung kata "Old"**
      meskipun rule sumbernya bernama demikian. *(AC 13 spec; penyimpangan sadar 8)*
- [ ] Daftar dibaca dari **satu tempat**, dipakai layar kontrak maupun seluruh grid klausul —
      bukan disalin per layar.
- [ ] Master yang **kosong atau tidak terbaca** menghasilkan kegagalan yang **terlihat**, bukan
      daftar kosong yang diam. *(**ADR-0015**)*

## Blocker

**Tidak ada.** **OQ-020 ditutup** — arti `ReinsTypeID` sudah jelas: master bersama life & non-life,
dipisahkan oleh saringan di atas.

## Catatan

⚠️ **Beda ejaan antar konteks.** `.Flag` bernilai **`1`** di Master Contract Retro Life, tetapi
**`"active"`** di sini. Jangan menyalin nilai dari konteks Life.

⚠️ **Jangan ikut memigrasikan `GetMasterReinsTypeContract`** — kueri itu membaca `M_TREATYCONTRACT`
dan `M_TREATYYEAR` lewat `a.JSONDATA.…`, dua tabel JSON yang sudah **mati** (tiket 01).

⚠️ **Layar master jenis reasuransi ada di luar konteks ini.** `[keputusan work owner]`
`Harness/InboxTreatyContractReinsType.xml` (`DATA-PORTAL!INBOXTREATYCONTRACTREINSTYPE`) dan RD kedua
adalah perilaku **layar master**, konteks/menu tersendiri. Di sini: **baca saja**.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — saringan blacklist hanya terbukti benar
bila diuji terhadap master yang benar-benar memuat ID yang dikecualikan.

```
go test ./internal/...
cd frontend && npm test
make check
```


---

## Pembacaan ulang XML — 28-09-2026 (sesi modul)

Nomor baris = nomor baris mentah berkas korpus (satu tag per baris; sama dengan `sed -e 's/></>\n</g'`).

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| RD dominan `BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` (kelas `ASM-FW-GISFW-Int-REINSURANCETYPE` b44) | `<pyFilterLogic>A AND B AND C</pyFilterLogic>` b557; **A** `.ID` b567 nilai `"10004","10011","10012","10021","10022","10025","10026","10028","10248","10249","10018","10217"` b573 dengan **`<pyFilterOperation>NotStartsWith</pyFilterOperation>` b581**; **C** `.Flag` b586 `"active"` (tanpa `pyFilterOperation` = Equal); **B** `.Type` b603 `"1","2","3"` (Equal atas daftar = IN); kolom `.ID`, `.Note` (sort ASC b667), `.Type`, `.SOANote`, `.Code`; `pyMaxRecords` 500 b742 | `repository/tco_jenisreasuransi.go`: `FLAG = 'active' AND TYPE IN ('1','2','3') AND ID NOT LIKE '<awalan>%'` ×12, `ORDER BY NOTE ASC, ID ASC`; tabel kebenaran `LolosSaringanNonLifeTCO` |
| Tabel fisik master | `Claim Non Prop/RDBList/GetReinsuranceTypeBYName_SQL.xml` `select ID, NOTE … from REINSURANCETYPE where type='4' and FLAG='active'`; `EDM Treaty In/RDBList/GetReinstypeIDbyName_SQL.xml` `from REINSURANCETYPE where note=…`; `Claim Non Prop/RDBList/GetTreatyName_SQL.xml` `c.nourut from pooldata.reinsurancetype`; korpus: 18× `from pooldata.reinsurancetype` | `POOLDATA.REINSURANCETYPE` kolom `ID`, `NOTE`, `TYPE`, `FLAG` (dibaca); `NOURUT` (ada, tidak dipakai); `SOANOTE`/`CODE` `[terbuka]` nama fisiknya |
| Pemakai RD dominan | 11 grid klausul (xlsx r57, r73, r84, r115, …: `GridTreatyArrangementProfitCommision`, `GridTreatyArrTreatyExGratiaChildList`, `gridTreatyArrangementExGratiaList`, `GridTreatyArrangementCashLossLimit`, …) | satu endpoint `GET /api/treaty-contract-out/jenis-reasuransi`; satu komponen `PilihJenisReasuransi` |
| RD non-Old `BrowseReinsuranceType_RD.xml` | saringan **berparameter** `Param.ID` b607, `Param.Note` b622, `Param.Flag` b640, `Param.Type` b658; sort `.ID DESC` b698; **satu-satunya pemakai**: `Section/InputTreatyContractReinsType.xml` b2652 pemilih `ReinsType` (`InputTreatyContract.ReinsTypeName`, tampil `.Note` b2724, `"active"` b2768) | ⚠️ pemilih FORM KONTRAK (tiket 04) — lihat ralat 1 |
| `RDBList/GetMasterReinsTypeContract.xml` | membaca `M_TREATYCONTRACT`/`M_TREATYYEAR` lewat `a.JSONDATA.…` | ➖ MATI (penyimpangan sadar 1; penjaga `TestTCONolTabelDokumenWarisan`) |
| Label pemilih | `Section/InputTreatyContractReinsType.xml` b2652 `<pyLabelFieldValue>ReinsType</pyLabelFieldValue>`; `Section/InputTreatyContract.xml` b7925 `<pyLabelFieldValue>Reinsurance Type</pyLabelFieldValue>`, b7948 `<pyLabelPreview>` | `assets/labels.treaty-contract-out.ts` `JENIS_REASURANSI_TCO` — keduanya VERBATIM, dijaga `labels.treaty-contract-out.test.ts` |

### Ralat bertanggal 28-09-2026

1. **Tiket ini, bab Catatan** dan **spec §11**: *"`Harness/InboxTreatyContractReinsType.xml` dan RD kedua adalah
   perilaku layar master jenis reasuransi, konteks/menu tersendiri"* — **DIBANTAH XML**. Harness itu →
   `PanggilReinsType` b1825 → `InputTreatyContractReinsType` b1300 = editor **kontrak** per tahun (`Save` →
   `SaveTreatyContract_Act` b3642, grid business b14117, grid reinsurer b13364). RD non-Old adalah **pemilih
   `ReinsType` di form kontrak itu** (b2652), disaring `Flag="active"` (b2768) tanpa blacklist. Tidak ada layar
   master jenis reasuransi di modul ini. Keputusan tiket ini (satu daftar untuk seluruh layar, AC 5) tetap
   dipakai form kontrak di tiket 04; bahwa itu menyingkirkan 12 awalan yang pemilih kontrak Pega tidak
   singkirkan dicatat **OQ-TCO-06** untuk work owner.
2. **Spec §11 dan tiket ini** menulis saringan `.ID NOT IN (…)` — operator XML adalah **`NotStartsWith`**
   (b581): *tidak berawalan* salah satu dari 12 nilai. Untuk ID lima karakter keduanya sama; untuk ID yang
   lebih panjang tidak. SQL memakai `NOT LIKE '<awalan>%'`, uji `100041` (berawalan `10004`) tersingkir.
3. Nama RD menyebut `Old_Ljt_id_isnotnull`, tetapi ketiga saringannya **tidak** memuat `OLD_LJT_ID`
   (kolom `M_REINSURANCETYPE` menurut katalog DEV) — nama sisa (OQ-066), tidak dibawa.

### Yang dibangun

`repository/tco_jenisreasuransi.go` (+uji SQL/tabel kebenaran/korpus), `services/tco_jenisreasuransi.go`
(pembaca disuntik; master kosong = `ErrMasterJenisReasuransiKosong` → 503, ADR-0015),
`handlers/rute_treaty_contract_out.go` (`GET /api/treaty-contract-out/jenis-reasuransi`, satu baris di
`Router`), `handlers/tco_db_test.go` (seam HTTP terhadap master tiruan yang memuat blacklist, awalan, flag,
tipe 4), `skemauji/tco_tiruan.go` (+tiruan `REINSURANCETYPE`), penjaga nama jujur dua sisi
(`tco_nama_jujur_test.go`, `namaJujur.test.ts`), frontend `assets/labels.treaty-contract-out.ts`,
`components/treaty-contract-out/PilihJenisReasuransi.tsx`, `services/api.ts` `ambilJenisReasuransiTreaty`.

**OQ dibuka:** OQ-TCO-06 — pemilih form kontrak (tiket 04) memakai daftar tersaring blacklist (keputusan
tiket ini) padahal Pega menyaring `Flag="active"` saja; OQ-TCO-07 — nama fisik `.SOANote`/`.Code`.

**Status:** selesai 28-09-2026 — commit `treaty-contract-out: tiket 02 — jenis reasuransi master dibaca + saringannya`.
