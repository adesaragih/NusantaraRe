# Modul — Komite Claim Life

STEP D3. Sintesis dari `../inventory/Komite Claim Life.md` (D1) + `../flows/Komite Claim Life.md`
(D2 Tahap 2) + `../flows/_SUMMARY-komite.md`. **47 file** — modul terkecil kedua di korpus.
Audit: `find "Komite Claim Life" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Tangga persetujuan (komite) untuk klaim lini life**. Class kerja
`ASM-FW-GCNMFW-WORK-KOMITELIFE` (12 rule). Titik masuk
`Komite Claim Life/Flow/KomiteLife_Flow.xml` → `RULE-OBJ-FLOW` /
`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW`, `<pyStartActivity>Start2`,
hash **`b69929c407`**.

`[dugaan]` Komite adalah **tahap proses dengan model data, roster, penomoran, dan penulisan database
sendiri** — bukan wrapper approval tanpa data. Bukti: class kerja sendiri, `KomiteList` sebagai
model data, roster dari class `ASM-FW-GCNMFW-Int-EMAILKOMITE`, rule penomoran sendiri.
**Belum terverifikasi penuh** — sebagian activity belum ditelusur.

## 2. Proses / fitur utama

Rincian: **`../flows/Komite Claim Life.md`** (12 rule ditelusur).

`[terverifikasi]` Graf: **1 Assignment, 1 Decision, 4 connector** — bentuk yang sama di keempat
modul Komite:

```
Start ─(Always)─> Assignment "KomiteRouter" [WorkList, route=Custom]
                      │ FlowAction: ViewTransferDtl
                      v
                  Decision "KomiteLoop"
                      ├─ When IsKomiteLoop ─> Assignment   ← LOOP
                      └─ Else ─────────────> End
```

`[terverifikasi]` **Tangga persetujuan tidak dimodelkan sebagai shape** — ia loop satu assignment
yang dikendalikan data. Kondisi `IsKomiteLoop`: `.AcceptStatus = "1"` **DAN**
`.KomiteCount <= .KomiteLoop`.

`[terverifikasi]` Jejak tiap tingkat ditulis sebagai entri `KomiteList(KomiteCount)` berisi
`KomiteAproval` (= `AcceptStatus`), `KomiteComment`, `DateApprove` (`@CurrentDateTime()`).
**Berapa tingkat tangga dan siapa penyetujunya ditentukan data, bukan struktur proses.**

### 2.1 Routing berbasis data — khas modul ini

`[terverifikasi]` **Satu-satunya dari empat modul Komite yang routing-nya sepenuhnya berbasis data**:
`param.AssignTo` diisi **hanya** dari `.KomiteID`, digerbangi `TransferType == '2'`.
Tiga modul Komite lain memuat empat sasaran ter-hardcode (`komitepnc`…`komitepnc4`) → OQ-036.

`[terverifikasi]` **Tidak ada guard identitas orang ter-hardcode** di modul ini (0 file) — berbeda
dari Komite Claim FacIn (4 file) dan Komite Claim Prop (3 file) → OQ-021.

`[terverifikasi]` Memakai `WorkList`, bukan `WorkBasket` → OQ-028.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 18 Activity, 14 Connect-SQL, 3 ReportDefinition,
3 FlowAction, 3 Section, 2 When, 1 Flow, 1 DecisionTable, 1 ConnectREST, 1 SystemSettings.

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 |
| `POOLDATA.T_FOLDER_IMAGE`, `POOLDATA.TANGGAL_CLOSING`, **`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`**, `POOLDATA.KODE_PRODUKSI`, `POOLDATA.DIRECTTOKASIR_LOG`, `CURRENCY` | 1 masing-masing |

### 3.1 Fragmen skema nyata — meringankan sebagian OQ-001

`[terverifikasi]` `Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` adalah blok PL/SQL
`INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE (...)` + `COMMIT` — **bukan** stored procedure,
sehingga **55 nama kolom terbaca** (daftar lengkap di `../flows/Komite Claim Life.md` §3.1).
**Tipe kolom tetap tidak diketahui** (OQ-001).

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `ServiceGoogle` (`SETTING` → `LinkService!LinkService`,
tanpa URL literal). `POOLDATA.DIRECTTOKASIR_LOG` — jejak Kasir di sisi database.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Claim Life** | `UpdateOsAkseptasiClaimLife_sql` **identitas dan hash sama (`c50bfd9a12`)** di kedua modul dan **tidak di register OQ-011** → satu rule bersama, dipanggil dari dua sisi | `[terverifikasi]` |
| **Master Contract Retro Life / Product Life** | class `ASM-FW-GISFW-INT-TREATYYEAR_LIFE`, `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` dirujuk | `[terverifikasi]` |

`[terverifikasi]` D1 mengukur **59 dari 133 identitas modul Komite (44,4 %) tidak muncul di modul
Claim mana pun** (`../inventory/_summary.md` §15.4) — Komite bukan lapisan tipis di atas Claim.

## 6. Batasan & batas pengetahuan

`[terverifikasi]` Empat procedure, isinya tidak ada di korpus (OQ-002): `GL.F_GET_EMAIL`,
`POOLDATA.GET_TOKEN_STORAGE`, `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` Rule penomoran akseptasi khas modul ini: `Generate_NoAccept_KMT_Life` /
`Generate_NoAccept_KMT_LifeRetro`, bercabang menurut kode `Type` (`QP`/`QR` vs `TP`/`TR`) —
**arti kode belum terverifikasi** (OQ-020). `KMT` kepanjangan belum terverifikasi.

`[terverifikasi]` Tiga rule dipanggil keempat modul Komite: `GETTanggalClosing_SQL`,
`GetSequenceNumber_SQL`, `InsertHistoryAkseptasiPega_Sql` — yang terakhir **berkonflik** di
register OQ-011.

`[pertanyaan terbuka]` Sumber nilai `.KomiteLoop` (batas tangga) **belum terverifikasi** → OQ-032.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **1 file**,
`OperatorID.pyTelephone` **0**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-001** (sebagian ringan, §3.1), **OQ-002**, **OQ-011**, **OQ-018**, **OQ-020**, **OQ-021**,
**OQ-024**, **OQ-028**, **OQ-029**, **OQ-032**, **OQ-036**.
