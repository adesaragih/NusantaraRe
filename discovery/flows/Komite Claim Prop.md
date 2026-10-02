# Telusur Flow — Komite Claim Prop

STEP D2, Tahap 2 konteks #3. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Komite Claim Prop/Flow/KomiteTreaty_Flow.xml`
→ `RULE-OBJ-FLOW` / **`ASM-FW-GCNMFW-WORK-KOMITETREATY`** / `KOMITETREATY_FLOW`,
`<pyStartActivity>Start1`.

> **PENTING — nama file yang sama, rule yang berbeda.** `Komite Claim Non Prop` punya file
> `Flow/KomiteTreaty_Flow.xml` dengan **nama yang sama persis**, tetapi class-nya
> `ASM-FW-GCNMFW-WORK-KOMITETREATY**NONPROP**`. Keduanya **dua rule berbeda** dan ditelusur
> terpisah. Dokumen ini **hanya** tentang varian Prop.
>
> Bukti hash ternormalisasi 18 tag: Prop = `5eacdb3783`, Non Prop = `63b913a161`.

---

## 1. Diagram alur

**1 Assignment, 1 Decision, 4 connector** `[terverifikasi]` — bentuk graf sama dengan ketiga modul
Komite lainnya.

```
Start1 ─(Always)─> ASSIGNMENT63 "KomiteRouter"  [WorkList, pyRouteTo=Custom]
                        │ FlowAction: ViewTransferDtl
                        v
                    Decision1 "KomiteLoop"
                        ├─ When IsKomiteLoop ─> ASSIGNMENT63   ← LOOP
                        └─ Else ─────────────> END52
```

Guard loop `Komite Claim Prop/When/IsKomiteLoop.xml`
(`ASM-FW-GCNMFW-WORK-KOMITETREATY` / `ISKOMITELOOP`):
`.AcceptStatus = "1"` DAN `.KomiteCount <= .KomiteLoop` — **teks kondisi sama** dengan ketiga modul
Komite lain, meski rule-nya berbeda.

### 1.1 Routing

`Komite Claim Prop/Activity/KomiteRouter.xml` (`ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITEROUTER`),
7 step:

| Precondition | `param.AssignTo` |
| --- | --- |
| `.KomiteCount == 1` | `"komitepnc"` |
| `.KomiteCount == 2` | `"komitepnc2"` |
| `.KomiteCount == 3` | `"komitepnc3"` |
| `.KomiteCount == 4` | `"komitepnc4"` |
| `Primary.TransferType == '2'` | `.KomiteID` / `.Komite.KomiteID` |
| `Primary.TransferType != '2'` | (cabang terpisah) |

`[terverifikasi]` Nilai lain yang di-set router ini: `.KomiteCount+1` (penambahan pencacah) dan `""`.

Pola **sama dengan Komite Claim FacIn** — empat sasaran ter-hardcode `komitepnc`…`komitepnc4`
(OQ-036), dengan fallback ke data `.KomiteID`. Berbeda dari Komite Claim Life yang **hanya** memakai
`.KomiteID`.

---

## 2. Status / state yang berubah

Dari `Komite Claim Prop/Activity/KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITEPOSTADJUSTMENT`) dan
`Activity/KomitePost_Reject.xml`.

| Properti | Nilai literal | Peran | Arti |
| --- | --- | --- | --- |
| `pyWorkPage.AcceptStatus` | `"1"`, `"2"` | dua `Property-Set` berbeda; juga gerbang `Call SaveRejectTreatyIn_Act_KMT` | **belum terverifikasi** |
| `pyWorkPage.TransferType` | `"1"` (di `KomitePost_Reject`), `'2'` (di router) | gerbang `Obj-Open-By-Handle` | **belum terverifikasi** (OQ-020) |
| `pyWorkPage.Komite.TypeComentAnalysis` | `"5"` | gerbang `Obj-Open-By-Handle` bersama `TransferType=="1"` | **belum terverifikasi** |
| `.KomiteCount`, `.KomiteLoop` | pencacah & batas tangga | §1 | — |

`[terverifikasi]` `KomitePost_Reject` step 1:
`pyWorkPage.TransferType=="1" && pyWorkPage.Komite.TypeComentAnalysis=="5"` — **dua kode digabung**
sebagai satu gerbang. Arti keduanya **belum terverifikasi**.

`[terverifikasi]` Yang menarik: di `KomitePost_Reject`, panggilan
`Call ASM-FW-GCNMFW-Work-ClaimTreaty.SaveRejectTreatyIn_Act_KMT` digerbangi
**`AcceptStatus=="1"`** — bukan `"2"`. Ini **tidak** ditafsirkan; dicatat apa adanya karena
menunjukkan bahwa pemetaan `AcceptStatus` → "setuju/tolak" **tidak dapat ditebak dari nama rule**.

---

## 3. Objek Oracle yang disentuh

Dari `KomitePostAdjustment`:

| Class | RequestType |
| --- | --- |
| `ASM-FW-GISFW-Int-policyjson` | `GETTanggalClosing_SQL` |
| **`ASM-FW-GCNMFW-Int-V_POLIS`** | **`GenerateNoAcceptTreaty`** |
| `ASM-FW-GISFW-Int-policyjson` | `GetKodeProdNonLife_SQL` |
| `ASM-FW-GISFW-Int-policyjson` | `GetSequenceNumber_SQL` |
| `ASM-FW-GISFW-int-policyjson` | `InsertHistoryAkseptasiPega_Sql` |

`[terverifikasi]` Rule penomoran akseptasi modul ini adalah **`GenerateNoAcceptTreaty`** — berbeda
dari Komite Claim Non Prop (`GenerateNOAccCTNP`) dan dari Komite Claim FacIn
(`GenerateNoAccept` / `GenerateNoAcceptNonFire`).

`InsertHistoryAkseptasiPega_Sql` terdaftar di register **OQ-011** sebagai identitas berkonflik
(`ASM-FW-GISFW-INT-POLICYJSON` / `ASM!INSERTHISTORYAKSEPTASIPEGA_SQL`) — varian modul ini yang
dirujuk; varian modul lain berbeda isi dan **tidak dibaca**.

---

## 4. Integrasi eksternal

Modul ini punya `Activity/HitServiceToKasirKMT_Act.xml`, `Activity/SendEmailKlaim_KMT.xml`,
`Activity/SendErrorDirectKasir.xml`, `Activity/KonversiKlaim_Act.xml`, dan 3 rule `ConnectREST`
(`../inventory/Komite Claim Prop.md` §7) — **belum ditelusur**, di luar jalur inti flow.

`Komite Claim Prop/When/IsPEGAPROD.xml` ada (varian grup yang sama dengan Komite Claim FacIn).
Kondisinya **tidak terbaca** → OQ-029.

---

## 5. Batas pengetahuan

### 5.1 Guard identitas orang — tiga activity keputusan `[terverifikasi]`

| File | Peran |
| --- | --- |
| `Activity/KomitePostAdjustment.xml` | pasca-keputusan penyesuaian |
| `Activity/KomitePost_Close.xml` | pasca-keputusan tutup |
| `Activity/KomitePost_Reject.xml` | pasca-keputusan tolak |

Guard berada di `<pyExpression>` (5×) dan `<pyPropertiesValue>` (1×).
**Nilai nama orang tidak disalin** (OQ-021). `pyPosition` dipakai di 1 file.

Sama seperti Komite Claim FacIn: guard identitas berada **tepat di activity keputusan komite**,
bukan di lapisan UI.

### 5.2 Rule dipanggil lintas class

`Call ASM-FW-GCNMFW-Work-ClaimTreaty.SaveRejectTreatyIn_Act_KMT` — memanggil activity di class
**`...Work-ClaimTreaty`**, bukan class Komite. Salinan yang ada di modul ini:
`Komite Claim Prop/Activity/SaveRejectTreatyIn_Act_KMT.xml` — **belum ditelusur**.

`Call UpdateWorkObject` — sama seperti Komite Claim Life, activity ini **tidak ada** di modul ini
(satu-satunya salinan di korpus: `Endorsment Fac In/Activity/UpdateWorkObject.xml`) → OQ-035.

### 5.3 Rule berkonflik yang menyentuh modul ini

Termasuk `ASM-FW-GCNMFW-WORK-CLAIMTREATY` / `COUNTLISTCLAIMAMOUNTIDR`,
`ASM-FW-GCNMFW-WORK` / `KONVERSIKLAIM_ACT`, keluarga `T_STORAGE_IMAGE`, `@BASECLASS/ISPEGAPROD`.
**Tidak ditelusur** di sini; bila tersentuh di telusur lanjutan, **wajib per-varian**.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Komite Claim Prop/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITETREATY_FLOW` / `RULE-OBJ-FLOW` | `Flow/KomiteTreaty_Flow.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY` | `Activity/KomiteRouter.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `Activity/KomitePostAdjustment.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITEPOST_REJECT` / `RULE-OBJ-ACTIVITY` | `Activity/KomitePost_Reject.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` / `KOMITEPOST_CLOSE` / `RULE-OBJ-ACTIVITY` | `Activity/KomitePost_Close.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` / `ISKOMITELOOP` / `RULE-OBJ-WHEN` | `When/IsKomiteLoop.xml` | tidak |
| (FlowAction) `VIEWTRANSFERDTL` | `FlowAction/ViewTransferDtl.xml` | tidak |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — #305** |
| 5 rule `RULE-CONNECT-SQL` (§3) | `RDBList/` | sebagian |

**13 rule ditelusur.**

---

## 7. Pertanyaan terbuka

Tidak ada OQ baru yang eksklusif konteks ini. Dikuatkan: OQ-011 (§3, §5.3), OQ-020 (§2 — kode
`TypeComentAnalysis="5"` baru), OQ-021 (§5.1), OQ-029 (§4), OQ-035 (§5.2), OQ-036 (§1.1).

`[pertanyaan terbuka]` baru yang dicatat ke OQ-020: **`TypeComentAnalysis`** dengan nilai `"5"`
(`Komite Claim Prop/Activity/KomitePost_Reject.xml`).
