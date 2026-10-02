# Telusur Flow — EDM Treaty In

STEP D2, Tahap 1 konteks #2. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `EDM Treaty In/Flow/InputAddendumTreatyIn.xml`
→ `RULE-OBJ-FLOW` / `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` / `INPUTADDENDUMTREATYIN`,
`<pyStartActivity>Start1`. Ukuran 202.269 byte.

Seluruh pernyataan menyebut file yang dibaca (OQ-018).

---

## 1. Diagram alur

**4 Assignment, 8 Decision, 2 Utility, 22 connector** `[terverifikasi]`.

```
Start1 ─(Always)─> Assignment2 "Input Realitation" [WorkBasket]
                        │ FlowAction: InboxPolicyTreatyInAddendum
                        v
                    Decision3 "Is Correct?" (isApproved)
                        ├─ Status=No  ─> End3
                        └─ Status=YES ─> Decision5 "pxCreateOperator SPV?"
                             ├─ When IsSPVCreate ─> Decision10 "IS SPV TREATY 1"
                             │                        ├─ When IsSPVTreaty1 ─> Assignment7
                             │                        └─ Else ─────────────> Assignment4
                             └─ Else ──────────────> Decision6 "is Treaty 1"
                                                      ├─ When IsTreaty1 ───> Assignment4
                                                      └─ Else ─────────────> Assignment7

  Assignment4 "Acceptance by Sec Treaty" [WorkBasket] ─DeptHeadTreatyIn_UWAddendum─> Decision2
  Assignment7 "Acceptance by Sec Treaty" [WorkBasket] ─DeptHeadTreatyIn_UWAddendum─> Decision11

  Decision2  ├─ YES ─> Decision9      Decision11 ├─ YES ─> Decision9
             └─ No  ─> Assignment2               └─ No  ─> Assignment2

  Decision9 "TO TREATY DEPT HEAD?"
        ├─ When ToTREATYDEPTHEAD ─> Assignment1 "Acceptance by Dept. Head" [WorkBasket]
        │                              └─DeptHeadTreatyIn_UWAddendum─> Decision1 "Is Correct?"
        │                                     ├─ Status=YES/Yes ─> Utility1
        │                                     └─ Status=No ──────> Assignment2
        └─ Else ──────────────────> Utility1

  Utility1 "Save json policy addendum"  Activity: SaveJsonPolisTreatyInEDM_Act
     └─(ALWAYS)─> Utility2 "HIT SERVICE ARASAPAS"  Activity: serviceInsertArasapas_act
                     └─(ALWAYS)─> End3
```

`[terverifikasi]` **Graf ini tertutup** — setiap shape punya connector masuk. Berbeda dari
NB Treaty In yang punya sub-graf tanpa connector masuk (OQ-023).

### 1.1 Perbandingan dengan NB Treaty In `[terverifikasi]`

| Aspek | NB Treaty In | EDM Treaty In |
| --- | --- | --- |
| Class flow | `ASM-FW-GISFW-WORK` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` |
| Assignment | 6 | 4 |
| Decision | 12 | 8 |
| Label assignment approval | "Acceptance by Head. Treaty" | "Acceptance by Sec Treaty" |
| FlowAction | `InboxPolicyTreatyIn`, `DeptHeadTreatyIn_UW` | sama + akhiran **`Addendum`** |
| Guard `NopolisEmpty` | ada | **tidak ada** |
| Jalur komite klaim + Direktur | ada (tanpa connector masuk) | **tidak ada** |
| Activity simpan | `SaveJsonPolisTreatyIn_Act` (8 step) | `SaveJsonPolisTreatyInEDM_Act` (**22 step**) |

Guard yang **identik isinya** di kedua modul: `IsSPVCreate`, `IsSPVTreaty1`, `IsTreaty1`,
`ToTREATYDEPTHEAD` — kondisi literalnya sama persis.

### 1.2 Routing

Keempat Assignment memakai `<pyImplementation>WorkBasket` `[terverifikasi]`. Pemetaan shape →
nama workbasket tidak terbaca, sama seperti NB Treaty In → **OQ-024**.

---

## 2. Status / state yang berubah

| Properti | Nilai literal | Diubah/diuji oleh | Arti |
| --- | --- | --- | --- |
| `OperatorID.pyTelephone` | `"TREATY1"`, `"SPVTREATY1"` | `When/IsTreaty1.xml`, `When/IsSPVTreaty1.xml` | kode peran di field telepon → **OQ-027** |
| `pyWorkPage.pxCreateOperator` | 2 identifier orang | `When/IsSPVCreate.xml` | guard identitas → OQ-021 |
| `pyWorkPage.LetterNo` | `"TREATYINDEPTHEAD"` | `When/ToTREATYDEPTHEAD.xml` | token routing; **arti belum terverifikasi** |
| `Local.TglProd` | default **`25`** bila kosong | step 3 `SaveJsonPolisTreatyInEDM_Act`: `@if(Local.TglProd=="",25,Local.TglProd)` | `[dugaan]` hari tutup buku; **belum terverifikasi** |
| `PolicyTreatyIn.ProductionDate` | `@CurrentDateTime()` / `StatementDate` / awal bulan + 1 bulan | step 4–7 (lihat §2.1) | — |
| `param.isFOR` | `"POLICY"` | precondition step 8 | mode simpan |
| `.ClaimType`, `.ClaimPaymentType` | diuji `= ""` | precondition step 10 | **arti belum terverifikasi** |
| `OutputParam.IDPEGAOUT` | diuji `@contains(...,"ERR")` | precondition RDB-List terakhir | penanda error |
| `pyWorkPage.FlagErrorKonversi` | di-set `""` | `serviceInsertArasapas_act` step 4 | penanda error konversi |

### 2.1 Aturan tanggal produksi `[terverifikasi]`

Empat langkah `Property-Set` berturut-turut mengubah `pyWorkPage.PolicyTreatyIn.ProductionDate`:

| Precondition | Nilai baru |
| --- | --- |
| `@PropertyHasValue(pyWorkPage.PolicyTreatyIn.ProductionDate)` | `@CurrentDateTime()` |
| `@toInt(@DateTime.DateTimeDifference(...StatementDate, @CurrentDateTime(), D)) < 0` | `pyWorkPage.PolicyTreatyIn.StatementDate` |
| `pyWorkPage.OfferFacIn.PolicyData.PolicyNo == ""` | `@FormatDateTime(ProductionDate,"yyyyMM","Asia/Jakarta","in_ID") + "01T050000.000 GMT"` |
| — | `@addCalendar(ProductionDate, 0, 1, 0, 0, 0, 0, 0)` (tambah 1 bulan) |

`[terverifikasi]` Zona waktu **`Asia/Jakarta`** dan locale **`in_ID`** ditulis literal di rule.
Offset `T050000.000 GMT` = 05:00 GMT `[dugaan]` setara awal hari WIB; **belum terverifikasi**.

---

## 3. Objek Oracle yang disentuh

| Rule Connect-SQL dipanggil | Class | Dari activity | Peran |
| --- | --- | --- | --- |
| `GETTanggalClosing_SQL` (`RNM`) | `ASM-FW-GISFW-Int-policyjson` | `SaveJsonPolisTreatyInEDM_Act` step 2 | hasil → `Local.TglProd` |
| `TreatyInSearchProdKe` (`RNM`) | `ASM-FW-GISFW-Work` | `SaveJsonPolisTreatyInEDM_Act` | — |
| `SavePolisTreatyInEDM_SQL` (`RNM`) | `ASM-FW-GISFW-Int-POLISTREATYIN` | `SaveJsonPolisTreatyInEDM_Act` | simpan JSON |
| `InsertIntoJsonError_SQL` | `ASM-FW-GISFW-Int-policyjson` | `SaveJsonPolisTreatyInEDM_Act` (2×, satu ber-precondition `@contains(OutputParam.IDPEGAOUT,"ERR")`) | pencatatan error |
| `GetPolicyNoByCaseId` | `ASM-FW-GISFW-Int-policyjson` | `serviceInsertArasapas_act` step 3 | nomor polis dari case id |
| `CekSTSKonversiJson` (`RNM`) | `ASM-FW-GISFW-Int-policyjson` | `serviceInsertArasapas_act` | cek status konversi |
| `DeleteDataProduction` | — | `serviceInsertArasapas_act` | — |
| `INSERTJSON_JSONPOLISMONITORING_FACIN` | — | `serviceInsertArasapas_act` | monitoring |
| `UpdateErrorNoteJsonPolisMonitoring` | — | `serviceInsertArasapas_act` | monitoring error |

**Tabel/kolom fisik tidak dicatat di sini** karena D2 membaca dari sisi rule; nama objek yang
disentuh tiap rule Connect-SQL ada di `../inventory/EDM Treaty In.md` §6. Tipe kolom tetap tidak
diketahui (OQ-001).

---

## 4. Integrasi eksternal

`[terverifikasi]` `serviceInsertArasapas_act` (varian EDM Treaty In) melakukan **dua panggilan
`Connect-REST`**, didahului `Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`.

| Hal | Nilai |
| --- | --- |
| Rule ConnectREST | `EDM Treaty In/ConnectREST/convertJsonNusareToProduction.xml` |
| Identitas | `ASM-FW-GISFW-WORK` / `CONVERTJSONNUSARETOPRODUCTION` / `RULE-CONNECT-REST` |
| `<pyServiceName>` | `convertJsonNusareToProduction` |
| `<pyBaseURLSelectionType>` | **`SETTING`** — bukan URL literal |
| `<pyBaseURLSetting>` | `LinkService!LinkService` |

**Alamat sebenarnya tidak ada di korpus** — berasal dari Dynamic System Setting `=ResponLink.URL`
(D1 §11.1). Ini pola konfigurasi yang benar dan harus dipertahankan sebagai env var.

Rule pendukung yang dirujuk: `GetLinkService` (`ASM-FW-GISFW-Int-M_LINK_SERVICE`),
`IsSuccessHitService`, `StatusService` (`ASM-FW-GISFW-Data-StatusService`), `FlagErrorKonversi`.

---

## 5. Batas pengetahuan

### 5.1 `serviceInsertArasapas_act` — varian berkonflik (OQ-011)

`[terverifikasi]` Identitas `ASM-FW-GISFW-WORK` / `SERVICEINSERTARASAPAS_ACT` /
`RULE-OBJ-ACTIVITY` **terdaftar di register OQ-011** (entri #415: 2 varian, **2 isi berbeda**):

| Varian | Modul |
| --- | --- |
| `EDM Treaty In/Activity/serviceInsertArasapas_act.xml` | **dibaca di telusur ini** |
| `Endorsment Fac In/Activity/serviceInsertArasapas_act.xml` | isi **berbeda**, tidak dibaca |

Yang dicatat di §4 berasal **dari varian EDM Treaty In**. Perilaku varian Endorsment Fac In
**tidak dinyatakan** dan tidak boleh diasumsikan sama. Mana yang berlaku di production → OQ-011.

### 5.2 Fungsi Pega kustom

`@ASM.GetPageJSONString()` dirujuk di `SaveJsonPolisTreatyInEDM_Act` — sama seperti NB Treaty In.
**Source tidak ada di korpus** → OQ-012.

### 5.3 Isi rule Connect-SQL yang dipanggil

Sembilan rule Connect-SQL di §3 dipanggil lewat method `RDB-List` dengan parameter `ClassName` +
`RequestType`. SQL-nya ada di korpus, tetapi **stored procedure yang dipanggil sebagian SQL itu
tidak** (OQ-002) — daftar procedure per modul ada di `../inventory/EDM Treaty In.md` §6.

### 5.4 `isApproved` — di modul ini hanya `DecisionTable`

`[terverifikasi]` Berbeda dari NB Treaty In (yang punya **dua** rule `isApproved`: `When` dan
`DecisionTable`), EDM Treaty In hanya punya `EDM Treaty In/DecisionTable/isApproved.xml`
(`ASM-FW-GISFW-WORK` / `ISAPPROVED` / `RULE-DECLARE-DECISIONTABLE`).

Isi yang terbaca dari tag: satu `<pyValue>3`. Struktur tabel keputusan selengkapnya **tidak
terekstraksi** oleh pembacaan tag sederhana → perlu telusur lanjutan. Arti nilai `3` **belum
terverifikasi**. Memperkuat **OQ-026**.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `EDM Treaty In/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` / `INPUTADDENDUMTREATYIN` / `RULE-OBJ-FLOW` | `Flow/InputAddendumTreatyIn.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `SAVEJSONPOLISTREATYINEDM_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveJsonPolisTreatyInEDM_Act.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `SERVICEINSERTARASAPAS_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/serviceInsertArasapas_act.xml` | **YA — entri #415** |
| `ASM-FW-GISFW-WORK` / `INBOXPOLICYTREATYINADDENDUM` / `RULE-OBJ-FLOWACTION` | `FlowAction/InboxPolicyTreatyInAddendum.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `DEPTHEADTREATYIN_UWADDENDUM` / `RULE-OBJ-FLOWACTION` | `FlowAction/DeptHeadTreatyIn_UWAddendum.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISAPPROVED` / `RULE-DECLARE-DECISIONTABLE` | `DecisionTable/isApproved.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISSPVCREATE` / `RULE-OBJ-WHEN` | `When/IsSPVCreate.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISSPVTREATY1` / `RULE-OBJ-WHEN` | `When/IsSPVTreaty1.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISTREATY1` / `RULE-OBJ-WHEN` | `When/IsTreaty1.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `TOTREATYDEPTHEAD` / `RULE-OBJ-WHEN` | `When/ToTREATYDEPTHEAD.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `CONVERTJSONNUSARETOPRODUCTION` / `RULE-CONNECT-REST` | `ConnectREST/convertJsonNusareToProduction.xml` | tidak |
| 9 rule `RULE-CONNECT-SQL` (§3) | `RDBList/` | tidak |

**20 rule ditelusur.**

---

## 7. Pertanyaan terbuka & bukti untuk OQ-014 (arti `EDM`)

### 7.1 Bukti yang terkumpul untuk OQ-014 `[terverifikasi]`

Telusur ini mengumpulkan bukti penamaan, **bukan** kepanjangan singkatannya:

| Bukti | Path |
| --- | --- |
| Flow bernama **`InputAddendumTreatyIn`** berada di modul bernama **`EDM Treaty In`** | `Flow/InputAddendumTreatyIn.xml` |
| Class flow adalah `ASM-FW-GISFW-WORK-**ENDORSEMENTTREATY**` | `<pxInsName>` flow |
| FlowAction memakai akhiran **`Addendum`** | `FlowAction/InboxPolicyTreatyInAddendum.xml`, `FlowAction/DeptHeadTreatyIn_UWAddendum.xml` |
| Activity simpan memakai akhiran **`EDM`** untuk proses yang sama | `Activity/SaveJsonPolisTreatyInEDM_Act.xml` |
| Di NB Treaty In, `param.isFOR` bernilai `"EDM"` atau `"POLICY"` sebagai **dua mode simpan** | `NB Treaty In/Activity/SaveJsonPolisTreatyIn_Act.xml` |

`[dugaan]` `EDM`, `Addendum`, dan `Endorsement` menunjuk **konsep yang sama** dalam korpus ini:
perubahan/amandemen atas polis treaty yang sudah ada, dibedakan dari penutupan baru (`POLICY`).

**Kepanjangan `EDM` tetap `belum terverifikasi`** — korpus tidak memuat satu pun tempat yang
menjabarkannya. OQ-014 **tetap terbuka**, tetapi kini punya bukti konteks pemakaian.

### 7.2 OQ baru

Tidak ada OQ baru dari konteks ini. OQ yang dikuatkan: **OQ-011** (§5.1), **OQ-012** (§5.2),
**OQ-014** (§7.1), **OQ-024** (§1.2), **OQ-026** (§5.4), **OQ-027** (§2).
