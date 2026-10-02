# Telusur Flow — Komite Claim Non Prop

STEP D2, Tahap 2 konteks #4. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Komite Claim Non Prop/Flow/KomiteTreaty_Flow.xml`
→ `RULE-OBJ-FLOW` / **`ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP`** / `KOMITETREATY_FLOW`,
`<pyStartActivity>Start1`.

> **PENTING — nama file yang sama, rule yang berbeda.** `Komite Claim Prop` punya file
> `Flow/KomiteTreaty_Flow.xml` dengan nama sama persis, class `ASM-FW-GCNMFW-WORK-KOMITETREATY`
> (tanpa `NONPROP`). **Dua rule berbeda**, ditelusur terpisah. Dokumen ini **hanya** varian Non Prop.
>
> Hash ternormalisasi 18 tag: Non Prop = `63b913a161`, Prop = `5eacdb3783`.

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

Guard `Komite Claim Non Prop/When/IsKomiteLoop.xml`
(`ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` / `ISKOMITELOOP`):
`.AcceptStatus = "1"` DAN `.KomiteCount <= .KomiteLoop` — teks sama dengan tiga modul lain.

### 1.1 Routing

`Komite Claim Non Prop/Activity/KomiteRouter.xml`, 7 step — **pola sama dengan Prop dan FacIn**:
`.KomiteCount == 1..4` → `"komitepnc"`, `"komitepnc2"`, `"komitepnc3"`, `"komitepnc4"`; cabang
`TransferType == '2'` / `!= '2'` dengan fallback `.KomiteID` / `.Komite.KomiteID` (OQ-036).

`[terverifikasi]` Berbeda dari Prop, router Non Prop **tidak** memuat nilai `.KomiteCount+1`.

---

## 2. Status / state yang berubah

Dari `Komite Claim Non Prop/Activity/KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` / `KOMITEPOSTADJUSTMENT`).

| Properti / kondisi | Nilai literal | Peran | Arti |
| --- | --- | --- | --- |
| `pyWorkPage.AcceptStatus` | `"1"`, `"2"` (dan `==2` tanpa kutip, 3×) | menentukan cabang pasca-keputusan | **belum terverifikasi** |
| `...AdjustmentList(IdxAdjustment).IsSubjectivity` | `true` | gerbang 2 `Property-Set` | **belum terverifikasi** |
| `pyWorkPage.Adjustment.AcceptedNo` | `""` | gerbang `RDB-List` penomoran | — |
| terminal tangga | `KomiteCount == KomiteLoop && AcceptStatus == "1"` | — | memakai **`KomiteLoop`** (seperti Komite Claim Life), bukan `Local.TotalKomite` (seperti FacIn) |
| `Local.currentDate` | `@toDecimal(...) > 25` **digabung** dengan terminal tangga | gerbang `Property-Set` | ambang tanggal 25 (OQ-030) |
| `PaymentType` | `1`, `2`, `3`, `4`, `5`, `6` | — | **belum terverifikasi** (OQ-020) |
| `TransferType` | `'2'` | gerbang router | **belum terverifikasi** |

`[terverifikasi]` Satu precondition **menggabungkan tiga kondisi**:
`@toDecimal(Local.currentDate) > 25 && KomiteCount == KomiteLoop && AcceptStatus == "1"` — ambang
tanggal 25 ikut menentukan perilaku **di tingkat akhir tangga**. Ini pemakaian angka 25 paling
konsekuensial yang ditemukan sejauh ini (bandingkan OQ-030).

---

## 3. Otorisasi: pola yang **berbeda** dari Prop dan FacIn

`[terverifikasi]` **Tidak ada satu pun guard identitas orang ter-hardcode** di modul ini:

```
grep -rlE "OperatorID\.pyUser(Identifier|Name)[ ]*[=!]+[ ]*[\"']" "Komite Claim Non Prop" --include="*.xml"
→ (kosong)
```

Sebagai gantinya, `KomitePostAdjustment` step 2 dan 3 memakai:

```
@contains( @toUpperCase(pyWorkPage.KomiteList(Local.IdxKomite).KomiteID),
           @toUpperCase(OperatorID.pyUserIdentifier) )
```

yaitu **mencocokkan pengguna yang sedang login dengan entri roster komite**, bukan dengan daftar
nama ter-hardcode. Hasilnya dipakai untuk `Property-Set` dan `Page-Set-Messages` (pesan ke layar).

`[dugaan]` ini pemeriksaan "apakah Anda anggota komite pada tingkat ini". **Belum terverifikasi**,
tetapi **secara pola jelas berbeda** dari Komite Claim Prop (3 activity ber-guard nama) dan
Komite Claim FacIn (4 activity ber-guard nama).

Guard jabatan yang ada hanya satu: `pyPosition != "IT Developer"`.

**Ini temuan penting untuk OQ-021:** korpus memuat **dua pendekatan otorisasi yang berbeda untuk
proses yang setara**. Salah satunya (Non Prop) sudah berbasis data.

---

## 4. Objek Oracle yang disentuh

Dari `KomitePostAdjustment`:

| Class | RequestType |
| --- | --- |
| `ASM-FW-GISFW-Int-policyjson` | `GETTanggalClosing_SQL` |
| **`ASM-FW-GCNMFW-Int-V_POLIS`** | **`GenerateNOAccCTNP`** |
| `ASM-FW-GISFW-Int-policyjson` | `GetKodeProdNonLife_SQL` |
| `ASM-FW-GISFW-Int-policyjson` | `GetSequenceNumber_SQL` |
| `ASM-FW-GISFW-int-policyjson` | `InsertHistoryAkseptasiPega_Sql` |

`[terverifikasi]` Penomoran akseptasi modul ini memakai **`GenerateNOAccCTNP`** — berbeda dari
Prop (`GenerateNoAcceptTreaty`) dan FacIn (`GenerateNoAccept` / `GenerateNoAcceptNonFire`).
**Empat modul Komite, empat rule penomoran berbeda.**

---

## 5. Activity khas Non Prop yang tidak ada di Prop

`[terverifikasi]` Modul ini punya activity yang **tidak ada** padanannya di Komite Claim Prop:

| Activity | Ukuran |
| --- | ---: |
| `GenerateAccCNP_act.xml` | 325.389 byte |
| `KomitePostAdjustmentCWP.xml` | 434.253 byte |
| `SaveRejectOSKomiteCNP.xml` | 188.811 byte |
| `InsertOSKlaimCNP.xml` | 104.871 byte |
| `InsertXOLKlaimCNP.xml` | 92.391 byte |
| `InsertOSSubjectivityCNP.xml` | 86.641 byte |
| `ResetSubjectivityNote.xml` | — |

Sebaliknya, Komite Claim Prop punya `KomitePost_Close.xml` dan `KomitePost_Reject.xml` yang
**tidak ada** di Non Prop.

`[terverifikasi]` Akhiran `CNP` dan `CWP` **kepanjangan belum terverifikasi**. `XOL` juga
(sudah tercatat di glossary). Adanya `InsertXOLKlaimCNP` dan `InsertOSSubjectivityCNP` menunjukkan
Non Prop menangani **XOL** dan **subjectivity** sebagai entitas tersendiri — `[dugaan]`, belum
ditelusur isinya.

**Seluruh activity ini belum ditelusur** — di luar jalur inti flow.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Komite Claim Non Prop/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` / `KOMITETREATY_FLOW` / `RULE-OBJ-FLOW` | `Flow/KomiteTreaty_Flow.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY` | `Activity/KomiteRouter.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `Activity/KomitePostAdjustment.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` / `ISKOMITELOOP` / `RULE-OBJ-WHEN` | `When/IsKomiteLoop.xml` | tidak |
| (FlowAction) `VIEWTRANSFERDTL` | `FlowAction/ViewTransferDtl.xml` | tidak |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — #305** |
| `ASM-FW-GCNMFW-WORK` / `GETBASE64ATTACHMENT` / `RULE-OBJ-ACTIVITY` | `excludeXML/GetBase64Attachment.xml` | **YA** (dan di folder `excludeXML` — OQ-003) |
| 5 rule `RULE-CONNECT-SQL` (§4) | `RDBList/` | sebagian |

**13 rule ditelusur.**

---

## 7. Batas pengetahuan & pertanyaan terbuka

### 7.1 Batas pengetahuan

- **`IsPEGAPROD`** ada di modul ini; kondisinya **tidak terbaca** dari tag → OQ-029.
- **`excludeXML/GetBase64Attachment.xml`** — satu-satunya file di folder `excludeXML` di seluruh
  korpus, terverifikasi `RULE-OBJ-ACTIVITY`, **dan terdaftar berkonflik di OQ-011**. Status aktif/
  tidaknya tetap terbuka → OQ-003.
- Enam activity besar khas CNP (§5) **belum ditelusur**.
- `InsertHistoryAkseptasiPega_Sql` adalah identitas **berkonflik**; varian modul ini yang dirujuk,
  varian modul lain **tidak dibaca**.

### 7.2 Pertanyaan terbuka

Tidak ada OQ baru eksklusif. Dikuatkan: OQ-003 (§6), OQ-011 (§4, §7.1), OQ-020 (§2),
**OQ-021 (§3 — ditemukan pola otorisasi alternatif berbasis roster)**, OQ-029, OQ-030 (§2),
OQ-032 (`KomiteLoop` vs `Local.TotalKomite`), OQ-036 (§1.1).
