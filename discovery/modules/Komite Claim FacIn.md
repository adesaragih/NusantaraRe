# Modul — Komite Claim FacIn

STEP D3. Sintesis dari `../inventory/Komite Claim FacIn.md` (D1) +
`../flows/Komite Claim FacIn.md` (D2 Tahap 2) + `../flows/_SUMMARY-komite.md`. **114 file** —
modul Komite terbesar.
Audit: `find "Komite Claim FacIn" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Tangga persetujuan (komite) untuk klaim facultative inward**. Class kerja
`ASM-FW-GCNMFW-WORK-KOMITE` (15 rule). Titik masuk `Komite Claim FacIn/Flow/Komite_Flow.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GCNMFW-WORK-KOMITE` / `KOMITE_FLOW`, `<pyStartActivity>Start1`,
hash **`f9366db22d`**.

## 2. Proses / fitur utama

Rincian: **`../flows/Komite Claim FacIn.md`** (15 rule ditelusur).

`[terverifikasi]` Graf sama bentuknya dengan tiga modul Komite lain: 1 Assignment "KomiteRouter"
(`WorkList`/`Custom`) + 1 Decision "KomiteLoop", 4 connector. Kondisi `IsKomiteLoop`:
`.AcceptStatus = "1"` **DAN** `.KomiteCount <= .KomiteLoop`.

`[terverifikasi]` **Pembeda modul ini:** properti batas tangga adalah **`Local.TotalKomite`**, bukan
`.KomiteLoop` seperti tiga modul lain.

### 2.1 Roster dari database + ambang nominal ter-hardcode

`[terverifikasi]` `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` membangun roster lewat
`Obj-Browse` terhadap class **`ASM-FW-GCNMFW-Int-EMAILKOMITE`** ke halaman `GetKomite`, lalu
membuang duplikat (`Property-Remove` bila `KomiteList(1).KomiteID == .OPERATOR_ID`).

`[terverifikasi]` Di activity yang sama terdapat **pita nilai ter-hardcode**:

```
Local.TotalAdj > 30000000.00 && Local.TotalAdj <= 57750000.00
```

yang menentukan `Obj-Browse` mana yang dijalankan — artinya **komposisi roster bergantung besaran
nilai**. `Local.TotalAdj` diakumulasi dari `+ .ValueAdjustment`. **Mata uang tidak disebut.**
Activity ini 146 KB dan **belum habis dibaca**, jadi jumlah pita seluruhnya **belum terukur**
→ **OQ-037**.

`[terverifikasi]` Batas bawah **30.000.000,00 sama persis** dengan `Local.LimitMax` di
`Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml`, tetapi batas atasnya berbeda
(57,75 jt di sini vs 50 jt di sana) → OQ-040.

**Jadi daftar anggota komite berasal dari tabel database, sementara sasaran routing justru
ter-hardcode** — kedua mekanisme hidup berdampingan.

### 2.2 Guard identitas orang

`[terverifikasi]` **4 file** memuat guard identitas orang ter-hardcode di `<pyExpression>` —
**di dalam activity keputusan komite, bukan di lapisan UI**: `KomitePost_Adjustment`,
`KomitePost_CloseClaim`, `KomitePost_Reject`, `SetProteksiSubmiteKomite`.
**Nilai nama orang tidak disalin** → OQ-021.

Audit: `grep -rlE "OperatorID\.pyUser(Identifier|Name)[ ]*[=!]+[ ]*[\"']" "Komite Claim FacIn" --include="*.xml"`

`[terverifikasi]` Empat sasaran routing ter-hardcode: `.KomiteCount == 1..4` →
`"komitepnc"`, `"komitepnc2"`, `"komitepnc3"`, `"komitepnc4"` → OQ-036.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: **49 When** (terbanyak relatif), 28 Activity,
17 Connect-SQL, 5 FlowAction, 5 Section, 3 ReportDefinition, 3 ConnectREST, 1 DataTransform,
1 Flow, 1 DecisionTable, 1 SystemSettings. Class teratas: `@BASECLASS` (26) → OQ-009.

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 |
| `REINSURANCE.TRLOSS_DETAIL_T`, `POOLDATA.T_FOLDER_IMAGE`, `POOLDATA.SUBPROGRESSCLAIM`, `POOLDATA.MONITORING_KLAIM_LOG`, `POOLDATA.KODE_PRODUKSI`, `POOLDATA.DIRECTTOKASIR_LOG`, `POOLDATA.CLAIMREJECTED`, `HISTORYAKSEPTASIPEGA` | 1 masing-masing |

`[terverifikasi]` Menulis akseptasi lewat `SaveOSClaim_SQL` ke `OS_AKSEPTASI_KLAIM`.

## 4. Integrasi eksternal

`[terverifikasi]` Tiga `RULE-CONNECT-REST`: `KonversiKlaimNonLife`, `SendAcceptationToKasir`,
`ServiceGoogle` — seluruhnya `SETTING` → `LinkService!LinkService`, tanpa URL literal.
`POOLDATA.DIRECTTOKASIR_LOG` — jejak Kasir di sisi database.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Claim Fac In** | `OS_AKSEPTASI_KLAIM` disentuh 6 rule di Claim Fac In dan lewat `SaveOSClaim_SQL` di sini | `[terverifikasi]` |
| Roster komite | class `ASM-FW-GCNMFW-Int-EMAILKOMITE` (tabel database) | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` Lima procedure, isinya tidak ada di korpus (OQ-002): `GL.F_GET_EMAIL`,
`POOLDATA.GET_TOKEN_STORAGE`, `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`, `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` Rule penomoran khas: `GenerateNoAccept` / `GenerateNoAcceptNonFire`, bercabang
menurut `When IsFire` — **kriteria pencabangan berbeda dari Komite Claim Life** yang memakai kode
produk (`Type`).

`[terverifikasi]` `Activity/KomitePost_Adjustment.xml` **626 KB — belum habis dibaca** (batas
cakupan telusur).

`[terverifikasi]` **Kode tanpa arti terverifikasi**: `AcceptStatus` `"1"`/`"2"`,
`TransferType` `'1'`/`'2'`, kode lini bisnis `L1`–`L11` → OQ-020, OQ-038.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-009**, **OQ-011**, **OQ-020**, **OQ-021**, **OQ-028**, **OQ-029**, **OQ-032**,
**OQ-036**, **OQ-037**, **OQ-038**, **OQ-040**.
