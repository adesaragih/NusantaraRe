# Sintesis Domain Claim — D2 Tahap 3

Perbandingan berbukti atas **empat** siklus kerugian yang ditelusur 2026-09-13:
Claim Life, Claim Prop, Claim Non Prop, Claim Fac In.

**Ini bukan penetapan bounded context** — itu STEP D4. Dokumen ini menyajikan bukti.

---

## 1. Empat modul, empat rule — bukan salinan

`[terverifikasi]` Dua pasang berbagi nama file Flow tetapi berbeda class dan berbeda isi:

| Modul | Class flow | Nama file | Start | Hash |
| --- | --- | --- | --- | --- |
| Claim Life | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `Register_Flow.xml` | `Start2` | — |
| Claim Fac In | `ASM-FW-GCNMFW-WORK-PNC` | `Register_Flow.xml` | `Start1` | — |
| Claim Prop | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `Flow_TreatyIn.xml` | `Start1` | `9871faee8d` |
| Claim Non Prop | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `Flow_TreatyIn.xml` | `Start1` | `b899c224ab` |

Keempatnya ditelusur **terpisah**.

---

## 2. Tahapan siklus kerugian — **berbeda-beda**

`[terverifikasi]` Tidak ada dua modul yang punya rangkaian tahap yang sama:

| Modul | Shape | Tahapan di graf | Jalur mundur |
| --- | ---: | --- | --- |
| **Claim Life** | 4A + 3D | Register → Outstanding → Medical Check → Claim Analis (Akseptasi) | `IsSendtoAdmin` (3 titik), `IsSendtoMedical` |
| **Claim Fac In** | 3A + 3D | (B2B?) → Register → Estimasi → Choose Surveyor | `IsBackStage` (2 titik) |
| **Claim Prop** | 2A + 1D | Outstanding → Input Acceptation | `IsBackStage` |
| **Claim Non Prop** | 2A + 1D | Outstanding → Input Acceptation | `IsBackStage` |

**Hanya Claim Life yang punya tahap Medical Check.** **Hanya Claim Fac In yang punya Estimasi,
Surveyor, dan percabangan B2B di titik masuk.** Kedua modul treaty (Prop/Non Prop) paling ringkas —
dua tahap saja.

`[terverifikasi]` **Tahap yang disebut siklus kerugian lengkap tidak seluruhnya muncul di graf.**
Adjustment, Close, dan Reject **tidak ada sebagai shape** di modul mana pun; keduanya ada sebagai
Activity (`SaveAdjustment_Act`, `ProtectCloseClaim_act`, `RejectOSClaimLife_Act`, `AdjClaimCNP_Act`)
di luar graf → **OQ-039**.

---

## 3. Routing: campuran di keempat modul

`[terverifikasi]` Keempat modul Claim memakai **dua model penugasan sekaligus** — pola yang belum
pernah muncul di Tahap 1 atau 2:

| Modul | `WorkList` | `WorkBasket` |
| --- | --- | --- |
| Claim Life | Input Register (`Current operator`), Outstanding Claim | Medical Check, Claim Analis |
| Claim Fac In | Input Register (`Current operator`), Input Estimasi | Choose Surveyor |
| Claim Prop | Outstanding Claim (`Current operator`) | Input Acceptation |
| Claim Non Prop | Outstanding Claim (`Current operator`) | Input Acceptation |

**Polanya konsisten:** tahap input/penyusunan → `WorkList` (sering `Current operator`); tahap
penilaian/persetujuan → `WorkBasket`.

Perbandingan lintas tahap D2:

| Domain | Model |
| --- | --- |
| Treaty inward (Tahap 1) | seluruhnya `WorkBasket` |
| Life offer & PremiumList (Tahap 1) | seluruhnya `WorkList` |
| Komite (Tahap 2) | seluruhnya `WorkList` |
| **Claim (Tahap 3)** | **campuran** |

→ memperkaya **OQ-028**.

---

## 4. `PaymentType` per lini — berbeda, dan Life tidak memakainya sama sekali

`[terverifikasi]`

| Modul | Nilai `PaymentType` yang diuji | `TransferType` |
| --- | --- | --- |
| **Claim Life** | **tidak ada sama sekali** | tidak ada |
| Claim Prop | `1, 2, 3, 4, 5, 6` | `2` |
| Claim Non Prop | `1, 2, 3, 4, 5, 6, 7` | `2` |
| Claim Fac In | `1, 2, 3, 4, 5, 6, 7` | `2` |

Perintah audit:
```
grep -rhoE "(PaymentType|TransferType)[ ]*[=!]+[ ]*[\"']?[0-9]{1,2}" "<modul>" --include="*.xml" \
  | sed 's/ //g' | sort | uniq -c | sort -rn
```

Kepadatan pemakaian tertinggi di **Claim Fac In** (`=4` 23×, `=6` 21×, `==3` 19×).

**Arti seluruh nilai belum terverifikasi.** Bukti Tahap 2 menegaskan mengapa tidak boleh ditebak:
`Komite Claim Prop/Activity/KomitePost_Reject.xml` justru digerbangi `AcceptStatus=="1"`, bukan
`"2"` — pemetaan kode ke makna bisa berlawanan dengan dugaan dari nama rule. → **OQ-020**.

Kode lain yang ditemukan: `STS_REJECT` (`'0'`, `'1'`, `'2'`) di Claim Life; `.pyNote = "Back"`
(gerbang `IsBackStage`, tiga modul); `BusinessCode` `L1`…`L11` (Claim Life, OQ-038);
`ContentNote = "DEATH"` (Claim Life).

---

## 5. Limit wewenang: **dua mekanisme yang hidup berdampingan**

Ini temuan paling material Tahap 3.

`[terverifikasi]`

| Modul | Mekanisme | Bukti |
| --- | --- | --- |
| **Claim Non Prop** | **ter-hardcode** di Activity | `CreateChildKomiteCNP_Act.xml`: `Local.LimitMax = 30000000.00`, `Local.LimitMaxDivHead = 50000000.00`, `Local.LimitPersenMax = 30.00` |
| **Claim Prop** | **dari database** | `RDBList/GetLimitDirekturUtama_SQL.xml`, `GetLimitPLATreatyin.xml`, `GetLimitsTreatyIn_SQL.xml` |
| Claim Non Prop | **juga** dari database | `GetLimitsTreatyIn_SQL.xml`, `GetLimitTONPPLA.xml` |
| Komite Claim FacIn (Tahap 2) | ter-hardcode | `ApprovalKomite_Act.xml`: `> 30000000.00 && <= 57750000.00` |

`[terverifikasi]` **Batas bawah 30.000.000,00 muncul identik** di `CreateChildKomiteCNP_Act`
(Claim Non Prop) dan `ApprovalKomite_Act` (Komite Claim FacIn), tetapi **batas atasnya berbeda**:
50.000.000,00 vs 57.750.000,00.

Mata uang **tidak disebut** di mana pun. `DivHead` dan arti `LimitPersenMax = 30.00` **belum
terverifikasi**. → **OQ-040**.

---

## 6. Jembatan Claim → Komite

`[terverifikasi]` Penyerahan ke tangga Komite dilakukan **Activity, bukan shape flow**:

| Modul Claim | Activity pembuat case anak Komite | Ukuran |
| --- | --- | ---: |
| Claim Non Prop | `CreateChildKomiteCNP_Act.xml` (`ASM-FW-GCNMFW-DATA-ADJUSTMENT`) | 756.836 |
| Claim Prop | `AddKomiteTreatyChild_ACT.xml` (`ASM-FW-GCNMFW-DATA-ADJUSTMENT`) | 420.247 |
| Claim Life | `CreateKMTLife_Act.xml`, `GetListKomiteLife.xml` | belum diukur |

Keduanya di class **`ASM-FW-GCNMFW-DATA-ADJUSTMENT`** — class yang sama, bukan class modulnya
sendiri. Digerbangi perbandingan nilai terhadap limit (§5).

Menjawab sebagian **OQ-039**: Komite dipicu dari sisi Claim oleh Activity ber-gerbang limit.

---

## 7. Tabel yang dibagi dengan Komite — **satu rule, dua pemanggil**

`[terverifikasi]` Klarifikasi penting terhadap Tahap 2.

`RDBList/UpdateOsAkseptasiClaimLife_sql.xml` ada di **Claim Life** dan **Komite Claim Life** dengan
identitas sama (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL`) dan
**hash ternormalisasi identik** (`c50bfd9a12`), serta **tidak terdaftar di register OQ-011**.

**Jadi bukan dua penulis independen ke tabel akseptasi, melainkan satu rule Connect-SQL yang sama
dipanggil dari dua sisi.** 55 kolom `OS_AKSEPTASI_KLAIM_LIFE` sudah terdaftar di
`Komite Claim Life.md` §3.1 — **tidak diulang**.

Rujukan `OS_AKSEPTASI_KLAIM` per modul Claim: Fac In **6 rule**, Non Prop 3, Prop 2.

---

## 8. Kebocoran batas ke domain Treaty

`[terverifikasi]` Dua modul **klaim** membaca objek domain **treaty**:

| Modul | Objek treaty | Rule |
| --- | --- | --- |
| Claim Non Prop | `M_TREATY_OUT`, `M_TREATY_OUT_DETAIL`, `TREATY_OUT` | `GetDataMasterTOutNP`, `BrowseDtlTreatyOutNP`, `GetLimitTONPPLA` |
| Claim Fac In | `TREATYCONTRACT`, `TREATYBUSINESS`, `PROPORTIONALARRG` | 2–3 rule masing-masing; + Activity `DLAFacintoTreaty_Act` (644 KB) |

`[dugaan]` jalur recovery/retrosesi dan jembatan Fac → Treaty. **Belum terverifikasi.**

Catatan silang: D1 §21.3 menemukan modul bernama **`Treaty Contract Out` justru tidak merujuk satu
pun objek treaty outward**, sementara `Claim Non Prop` **merujuk** → memperkuat **OQ-022**.

---

## 9. Konfirmasi temuan D1: trio Claim **tidak** berbagi ruleset

`[terverifikasi]` Overlap identitas rule antar modul Claim, diukur ulang di Tahap 3:

| Pasangan | Identitas dibagi | Identitas modul terkecil | ~% |
| --- | ---: | ---: | ---: |
| Claim Fac In ↔ Claim Prop | 108 | 270 | 40% |
| Claim Non Prop ↔ Claim Prop | 102 | 270 | 38% |
| Claim Fac In ↔ Claim Non Prop | 83 | 279 | 30% |
| Claim Life ↔ Claim Prop | 19 | 136 | 14% |
| Claim Fac In ↔ Claim Life | 18 | 136 | 13% |
| Claim Life ↔ Claim Non Prop | 16 | 136 | 12% |

Identitas unik per modul: Claim Fac In 481, Claim Non Prop 279, Claim Prop 270, Claim Life 136.

**Temuan D1 dikonfirmasi:** domain klaim **tidak** berbagi ruleset seperti domain facultative
(di mana `RNW Fac In` berbagi 99,0% dengan `NB FacIn`). `Claim Life` paling terpisah (12–14%).

Bukti D2 menguatkan dari sisi perilaku: tahapan berbeda (§2), `PaymentType` tidak dipakai Life (§4),
mekanisme limit berbeda (§5).

---

## 10. Batas pengetahuan Tahap 3

### 10.1 `ISCLM*` — berkonflik **dan** tidak terbaca

`[terverifikasi]` `@BASECLASS` / `ISCLM`, `ISCLMP`, `ISCLMNP`, `ISPEGASYARIAH` — masing-masing
**6 varian, 6 isi berbeda** (register #291, #293, #292, #306). Hash varian yang relevan:

| When | Claim Prop | Claim Non Prop | Claim Fac In |
| --- | --- | --- | --- |
| `ISCLM` | `2245b132` | `541179a0` | `4ee647f1` |
| `ISCLMP` | `4919199b` | `9b10dc97` | `a44d0c98` |
| `ISCLMNP` | `50c584af` | `d3255d2c` | `91c9b037` |
| `ISPEGASYARIAH` | `0a7077fd` | `1bbf8269` | `cac64576` |

**Masalah berlapis:** kondisi keduabelas varian **tidak terbaca** dari tag (`<pyLabel>` hanya
template kosong). Kita tahu semuanya berbeda, tetapi **tidak bisa melihat bedanya** → **OQ-041**.

`[terverifikasi]` **Claim Life tidak memiliki keempat rule ini** — konflik tersebut tidak
menyentuhnya.

### 10.2 When lain yang kondisinya tidak terbaca

| Rule | Modul | Yang digerbangi |
| --- | --- | --- |
| `IsSendtoMedical` | Claim Life | pengembalian ke tahap medis |
| `IsSPK` | Claim Fac In | **percabangan di titik masuk flow** (B2B) |
| `IsPEGAPROD` | keempat modul | integrasi keluar |

→ OQ-029.

### 10.3 Activity besar yang belum habis dibaca

| Modul | Activity terbesar | Ukuran |
| --- | --- | ---: |
| Claim Fac In | `CLaimFaceSheet_Act.xml` | **804.401** |
| Claim Non Prop | `CreateChildKomiteCNP_Act.xml` | 756.836 |
| Claim Fac In | `GenerateDLAFacin_Act.xml` | 690.099 |
| Claim Fac In | `DLAFacintoTreaty_Act.xml` | 644.505 |
| Claim Life | `SaveOutStandingLife_Act.xml` | 642.787 |
| Claim Non Prop | `CountLossAllocation_act.xml` | 635.029 |

**Perhitungan alokasi kerugian, XOL, DLA, dan face sheet berada di dalam activity ini dan belum
ditelusur** — batas cakupan telusur, bukan batas pengetahuan korpus.

### 10.4 Stored procedure

| Modul | Procedure |
| --- | ---: |
| Claim Prop | 8 |
| Claim Fac In | 5 (termasuk `PEGA_PROGRESSCLAIM`, `PEGA_SUBPROGRESSCLAIM` — hanya di modul ini) |
| Claim Non Prop | 4 |
| Claim Life | 1 |

**Isi seluruhnya tidak ada di korpus** (OQ-002).
