# Modul — Komite Claim Non Prop

STEP D3. Sintesis dari `../inventory/Komite Claim Non Prop.md` (D1) +
`../flows/Komite Claim Non Prop.md` (D2 Tahap 2) + `../flows/_SUMMARY-komite.md`. **59 file**.
Audit: `find "Komite Claim Non Prop" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Tangga persetujuan (komite) untuk klaim treaty inward non-proporsional**.
Class kerja `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` (12 rule). Titik masuk
`Komite Claim Non Prop/Flow/KomiteTreaty_Flow.xml` → `RULE-OBJ-FLOW` /
`ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` / `KOMITETREATY_FLOW`, `<pyStartActivity>Start1`,
hash **`63b913a161`**.

`[terverifikasi]` **Bukan salinan** `Komite Claim Prop` (hash `5eacdb3783`, class tanpa `NONPROP`)
— dua rule berbeda meski nama file identik.

## 2. Proses / fitur utama

Rincian: **`../flows/Komite Claim Non Prop.md`** (13 rule ditelusur).

`[terverifikasi]` Graf sama bentuknya dengan tiga modul Komite lain. Batas tangga `.KomiteLoop`;
kondisi `IsKomiteLoop`: `.AcceptStatus = "1"` **DAN** `.KomiteCount <= .KomiteLoop`.

### 2.1 Otorisasi berbasis roster — pola yang paling siap untuk RBAC

`[terverifikasi]` `Komite Claim Non Prop/Activity/KomitePostAdjustment.xml` step 2–3:

```
@contains( @toUpperCase(pyWorkPage.KomiteList(Local.IdxKomite).KomiteID),
           @toUpperCase(OperatorID.pyUserIdentifier) )
```

Pengguna yang login **dicocokkan dengan entri roster**, bukan dengan daftar nama ter-hardcode.

`[terverifikasi]` **Nol file** memuat guard identitas orang ter-hardcode — bersama
`Komite Claim Life`, satu dari dua modul Komite yang demikian. Bandingkan: FacIn 4 file,
Prop 3 file → OQ-021.

`[terverifikasi]` Routing tetap memuat empat sasaran ter-hardcode (`komitepnc`…`komitepnc4`)
→ OQ-036. **Jadi otorisasi sudah berbasis data, tetapi routing belum.**

`[terverifikasi]` **Activity `KomitePost_Close` dan `KomitePost_Reject` tidak ada** di modul ini —
keduanya ada di Komite Claim FacIn dan Komite Claim Prop.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 21 Activity, 16 Connect-SQL, 6 When, 5 ConnectREST,
3 ReportDefinition, 2 FlowAction, 2 Section, 1 DataTransform, 1 Flow, 1 DecisionTable,
1 SystemSettings.

| Objek | Rule perujuk |
| --- | ---: |
| `T_STORAGE_IMAGE`, `REINSURANCE.TRLOSS_DETAIL_T`, `POOLDATA.T_FOLDER_IMAGE`, `POOLDATA.KODE_PRODUKSI`, `POOLDATA.DIRECTTOKASIR_LOG`, **`OS_AKSEPTASI_KLAIM`**, `HISTORYAKSEPTASIPEGA` | 1 masing-masing |

`[terverifikasi]` Class kedua & keenam terbanyak adalah `ASM-FW-GCNMFW-INT-V_POLIS` (6) dan
**`ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP`** (5) — class kerja milik `Claim Non Prop`.

## 4. Integrasi eksternal

`[terverifikasi]` **5 rule `RULE-CONNECT-REST` — terbanyak di domain Komite**:
`KonversiKlaimNonLife`, `SendAcceptationToKasir`, `ServiceGoogle`, `insertClaimFinalOrClosed_NP`,
**`insertClaimReject_NP`** (yang terakhir hanya di modul ini). Seluruhnya `SETTING` →
`LinkService!LinkService`, tanpa URL literal.

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Claim Non Prop** | 5 rule berclass `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP`; `OS_AKSEPTASI_KLAIM` ditulis lewat `InsertOSKlaimCNP` | `[terverifikasi]` |
| Roster komite | class `ASM-FW-GCNMFW-Int-EMAILKOMITE` | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Delapan procedure — terbanyak di domain Komite**, isinya tidak ada di korpus
(OQ-002): `GL.F_GET_EMAIL`, `POOLDATA.GET_TOKEN_STORAGE`, `POOLDATA.PEGA_JSON_KLAIM_PNC`,
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`, `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`,
**`POOLDATA.PEGA_JSON_OS_AKSEP_SUBJECTIVITY`**, `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`,
**`POOLDATA.XOL2_AKSEP_KLAIM`** (dua terakhir hanya di modul ini).

`[terverifikasi]` Rule penomoran khas: `GenerateNOAccCTNP` (satu rule).

`[terverifikasi]` Activity khas **belum ditelusur**: `KomitePostAdjustmentCWP` (434 KB),
`GenerateAccCNP_act` (325 KB), `InsertXOLKlaimCNP`, `InsertOSKlaimCNP` — batas **cakupan telusur**.
Akhiran `CWP`, `CNP`, `XOL` **kepanjangan belum terverifikasi**.

`[terverifikasi]` **Kode tanpa arti terverifikasi**: `AcceptStatus` `"1"`/`"2"`,
`TransferType` `'1'`/`'2'` → OQ-020.

`[pertanyaan terbuka]` Sumber `.KomiteLoop` **belum terverifikasi** → OQ-032.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **3 file** (dipakai untuk
pencocokan roster, bukan sebagai konstanta), `OperatorID.pyTelephone` **0**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-011**, **OQ-018**, **OQ-020**, **OQ-021**, **OQ-028**, **OQ-029**, **OQ-032**,
**OQ-036**.
