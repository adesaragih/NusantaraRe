# Modul — Claim Fac In

STEP D3. Sintesis dari `../inventory/Claim Fac In.md` (D1) + `../flows/Claim Fac In.md` (D2 Tahap 3).
**482 file** — modul klaim terbesar di korpus.
Audit: `find "Claim Fac In" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **siklus kerugian (klaim) untuk bisnis facultative inward**. Class kerja
utamanya `ASM-FW-GCNMFW-WORK-PNC` (89 rule), dengan titik masuk
`Claim Fac In/Flow/Register_Flow.xml` → `RULE-OBJ-FLOW` / `ASM-FW-GCNMFW-WORK-PNC` /
`REGISTER_FLOW`, `<pyStartActivity>Start1`.

`[terverifikasi]` **Bukan salinan** `Claim Life`: modul itu punya `Flow/Register_Flow.xml` dengan
nama file sama tetapi class `ASM-FW-GCNMFW-WORK-CLAIMLIFE` dan `Start2` → dua rule berbeda.

`PNC` **kepanjangan belum terverifikasi**.

## 2. Proses / fitur utama

Rincian: **`../flows/Claim Fac In.md`** (17 rule ditelusur).

`[terverifikasi]` Tahapan: (percabangan **B2B**) → Register → Estimasi → Choose Surveyor → selesai,
dengan dua jalur mundur lewat `When IsBackStage` (`.pyNote = "Back"`).

`[terverifikasi]` **Satu-satunya modul Claim yang bercabang sebelum tahap pertama**: `Decision4`
berlabel "B2B", digerbangi `When IsSPK` (`Claim Fac In/When/IsSPK.xml`,
`ASM-FW-GCNMFW-WORK` / `ISSPK`) — bila benar, Register dan Estimasi dilewati.

`[pertanyaan terbuka]` Kondisi `IsSPK` **tidak terbaca** dari tag (`<pyLabel>` hanya template kosong)
→ OQ-029. `SPK` kepanjangan belum terverifikasi.

`[terverifikasi]` Routing campuran: `Assignment1` "Input Register" = `WorkList`/`Current operator`;
`Assignment7` "Input Estimasi" = `WorkList`/`Custom`; `Assignment3` "Choose Surveyor" =
**`WorkBasket`**/`Custom` → OQ-028.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 178 Activity, 64 Connect-SQL, 60 When, 56 Section,
33 FlowAction, 28 DataTransform, 27 ReportDefinition, 17 Harness, 8 DataPage, 8 ConnectREST,
1 Flow, 1 DecisionTable.

Objek Oracle teratas:

| Objek | Rule perujuk |
| --- | ---: |
| **`OS_AKSEPTASI_KLAIM`** | 6 — terbanyak dari 4 modul Claim |
| `JSON_POLIS` | 5 |
| `FACINPRODUCTION` | 4 |
| `T_STORAGE_IMAGE`, **`REINSURANCE.TRLOSS_DETAIL_T`**, `PROPORTIONALARRG`, `JSON_KLAIM`, `AGENT` | 3 masing-masing |
| `TREATYCONTRACT`, `TREATYBUSINESS`, `POOLDATA.DIRECTTOKASIR_LOG`, `M_CLIENT` | 2 masing-masing |

Audit:
```
awk -F'\t' '$1 ~ "^Claim Fac In/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
 | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

`[terverifikasi]` Skema **`REINSURANCE`** dirujuk — salah satu dari 10 skema aplikasi (OQ-016).

## 4. Integrasi eksternal

`[terverifikasi]` **8 rule `RULE-CONNECT-REST` — terbanyak di korpus**: `HitDLAClaimFacin`,
`KonversiKlaimNonLife`, `SendAcceptationToKasir`, `ServiceGoogle`, `getPayAttachment`,
`getPaymentClaim`, `getPremiumPaidOn`, **`getPremiumPaidOnMarine`**.

`[terverifikasi]` Tujuh memakai `pyBaseURLSelectionType = SETTING` (`LinkService!LinkService`);
**`getPremiumPaidOnMarine` memakai `URL` dengan endpoint literal** — satu-satunya di korpus.
Nilainya **tidak disalin** ke artefak D2/D3 karena memuat contoh data → OQ-018.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Komite Claim FacIn** | `OS_AKSEPTASI_KLAIM` disentuh 6 rule di sini dan lewat `SaveOSClaim_SQL` di modul Komite | `[terverifikasi]` |
| **Domain Treaty** | membaca `TREATYCONTRACT`, `TREATYBUSINESS`, `PROPORTIONALARRG`; ditambah `Claim Fac In/Activity/DLAFacintoTreaty_Act.xml` (`ASM-FW-GCNMFW-DATA-OBJECT` / `DLAFACINTOTREATY_ACT`, 644.505 byte) | `[terverifikasi]` keberadaan / `[dugaan]` arah — OQ-042 |
| **Produksi facultative** | `FACINPRODUCTION` (4 rule) | `[terverifikasi]` |

`DLA` **kepanjangan belum terverifikasi**.

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Lima stored procedure** dipanggil, isinya tidak ada di korpus (OQ-002):
`POOLDATA.PEGA_D_CAUSE_OF_LOSS`, `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`,
**`POOLDATA.PEGA_PROGRESSCLAIM`**, **`POOLDATA.PEGA_SUBPROGRESSCLAIM`** (dua terakhir hanya di modul
ini), `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`. Ditambah `GL.F_GET_EMAIL`,
`POOLDATA.GETCURRENCYSTANDARD`, `POOLDATA.GET_TOKEN_STORAGE`, `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`DBMS_LOB.CREATETEMPORARY`.

`[terverifikasi]` **Activity besar yang belum habis dibaca** — batas **cakupan telusur**, bukan
batas korpus:

| Activity | Ukuran |
| --- | ---: |
| `CLaimFaceSheet_Act.xml` | **804.401** — terbesar yang ditemui sepanjang D2 |
| `GenerateDLAFacin_Act.xml` | 690.099 |
| `DLAFacintoTreaty_Act.xml` | 644.505 |
| `DraftGenerateDLAFacin_Act.xml` | 630.950 |
| `CheckLimitSpreadingTreaty_Act.xml` | 602.168 |

`[terverifikasi]` **Kode tanpa arti terverifikasi**: `PaymentType` `1`–`7` (pemakaian terpadat dari
4 modul Claim: `=4` 23×, `=6` 21×, `==3` 19×), `TransferType` `2` → OQ-020.

`[terverifikasi]` `ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` — varian modul ini berhash
`4ee647f1`, `a44d0c98`, `91c9b037`, `cac64576`; **berkonflik antar modul DAN kondisinya tidak
terbaca** → OQ-041. Varian modul lain **tidak dibaca**.

`[terverifikasi]` `Claim Fac In/When/IsPEGAPROD.xml` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **7 file**,
`OperatorID.pyTelephone` **0**. Nilai nama orang **tidak disalin**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-016**, **OQ-018**, **OQ-020**, **OQ-028**, **OQ-029**, **OQ-039**,
**OQ-041**, **OQ-042**.
