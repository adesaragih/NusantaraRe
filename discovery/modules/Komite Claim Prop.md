# Modul — Komite Claim Prop

STEP D3. Sintesis dari `../inventory/Komite Claim Prop.md` (D1) + `../flows/Komite Claim Prop.md`
(D2 Tahap 2) + `../flows/_SUMMARY-komite.md`. **80 file**.
Audit: `find "Komite Claim Prop" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Tangga persetujuan (komite) untuk klaim treaty inward proporsional**.
Class kerja `ASM-FW-GCNMFW-WORK-KOMITETREATY` (14 rule). Titik masuk
`Komite Claim Prop/Flow/KomiteTreaty_Flow.xml` → `RULE-OBJ-FLOW` /
`ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITETREATY_FLOW`, `<pyStartActivity>Start1`,
hash **`5eacdb3783`**.

`[terverifikasi]` **Bukan salinan** `Komite Claim Non Prop`: nama file Flow sama persis, class
berbeda (`…KOMITETREATYNONPROP`), hash berbeda (`63b913a161`) → dua rule berbeda, ditelusur
terpisah.

## 2. Proses / fitur utama

Rincian: **`../flows/Komite Claim Prop.md`** (13 rule ditelusur).

`[terverifikasi]` Graf sama bentuknya dengan tiga modul Komite lain (1 Assignment "KomiteRouter"
`WorkList`/`Custom` + 1 Decision "KomiteLoop", 4 connector). Batas tangga: `.KomiteLoop`.
Kondisi `IsKomiteLoop`: `.AcceptStatus = "1"` **DAN** `.KomiteCount <= .KomiteLoop`.

`[terverifikasi]` Routing: empat sasaran ter-hardcode (`.KomiteCount == 1..4` → `"komitepnc"` …
`"komitepnc4"`), lalu cabang `TransferType == '2'` / `!= '2'` → `.KomiteID` / `.Komite.KomiteID`
→ OQ-036. **Pembeda kecil:** router modul ini memuat nilai `.KomiteCount+1`; Non Prop tidak.

### 2.1 Guard identitas orang

`[terverifikasi]` **3 file** memuat guard identitas orang ter-hardcode di `<pyExpression>` —
di dalam activity keputusan komite: `KomitePostAdjustment`, `KomitePost_Close`,
`KomitePost_Reject`. **Nilai nama orang tidak disalin** → OQ-021.

`[terverifikasi]` Kontras terukur: `Komite Claim Non Prop` menjalankan proses setara **tanpa**
guard nama, dengan mencocokkan pengguna terhadap roster (`@contains(…KomiteID, …pyUserIdentifier)`).
**Korpus memuat dua pendekatan otorisasi untuk proses yang setara.**

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 35 Activity, 23 Connect-SQL, 6 When, 4 ReportDefinition,
3 ConnectREST, 2 DataTransform, 2 FlowAction, 2 Section, 1 Flow, 1 DecisionTable, 1 SystemSettings.

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 |
| `TREATYYEAR`, `TREATYREINSURER`, `TREATYBUSINESS`, `REINSURANCE.TRLOSS_DETAIL_T`, `PROPORTIONALARRG`, `POOLDATA.T_FOLDER_IMAGE`, `POOLDATA.MONITORING_KLAIM_LOG`, `POOLDATA.KODE_PRODUKSI`, `POOLDATA.DIRECTTOKASIR_LOG` | 1 masing-masing |

`[terverifikasi]` Class kedua terbanyak adalah **`ASM-FW-GCNMFW-WORK-CLAIMTREATY`** (12 rule) —
class kerja milik `Claim Prop`. Modul Komite ini **membawa rule berclass Claim**.

## 4. Integrasi eksternal

`[terverifikasi]` Tiga `RULE-CONNECT-REST`: `KonversiKlaimNonLife`, `SendAcceptationToKasir`,
`ServiceGoogle` — seluruhnya `SETTING` → `LinkService!LinkService`, tanpa URL literal.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Claim Prop** | 12 rule berclass `ASM-FW-GCNMFW-WORK-CLAIMTREATY`; jembatan dari sisi Claim lewat `AddKomiteTreatyChild_ACT` | `[terverifikasi]` |
| **Master treaty** | membaca `TREATYYEAR`, `TREATYREINSURER`, `TREATYBUSINESS`, `PROPORTIONALARRG` | `[terverifikasi]` |
| Roster komite | class `ASM-FW-GCNMFW-Int-EMAILKOMITE` | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` Tujuh procedure, isinya tidak ada di korpus (OQ-002): `GL.F_GET_EMAIL`,
`POOLDATA.GETCURRENCYSTANDARD`, `POOLDATA.GET_TOKEN_STORAGE`, `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`, `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`,
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` Rule penomoran khas: `GenerateNoAcceptTreaty` (satu rule, tanpa pencabangan) —
berbeda dari Life (2 rule per kode `Type`) dan FacIn (2 rule per `IsFire`).

`[terverifikasi]` **Kode tanpa arti terverifikasi**: `AcceptStatus` `"1"`/`"2"`,
`TransferType` `'1'`/`'2'` → OQ-020.

`[pertanyaan terbuka]` Sumber nilai `.KomiteLoop` **belum terverifikasi** → OQ-032.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **3 file**,
`OperatorID.pyTelephone` **0**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-020**, **OQ-021**, **OQ-028**, **OQ-029**, **OQ-032**, **OQ-035**,
**OQ-036**.
