# flows/ — STEP D2

Status: **SELESAI** — Tahap 1–5 tuntas, **20 konteks / 428 rule** ditelusur.
Penutup: **`../D2-CLOSING-REPORT.md`**.

Konvensi telusur yang mengikat: **`_METHOD.md`** (baca ini dulu), ditambah
**`_METHOD-noflow.md`** untuk modul tanpa rule `Flow` (Tahap 5).
Urutan konteks yang dipakai: `../D1-CLOSING-REPORT.md` §6.2.

## Tahap 1 — konteks ber-`Flow` paling bersih (selesai)

| # | Konteks | Titik masuk | Rule | Catatan |
| ---: | --- | --- | ---: | --- |
| 1 | `NB Treaty In.md` | `Flow/InputRealizationTreatyIn.xml` → `Start1` | 15 | **blocker OQ-025** di langkah terakhir |
| 2 | `EDM Treaty In.md` | `Flow/InputAddendumTreatyIn.xml` → `Start1` | 20 | bukti konteks untuk OQ-014 |
| 3 | `Endorsement Life.md` | `Flow/InputEDMLife.xml` → `Start1` | 14 | bukti `EdmType`; OQ-012 maju (49 path JSON) |
| 4 | `PremiumList Life.md` | `InputPolicyHolder.xml` → `Start1` | 13 | titik masuk di root modul (OQ-004) |

## Tahap 2 — domain Komite / tangga persetujuan (selesai)

Sintesis lintas-modul: **`_SUMMARY-komite.md`**.

| # | Konteks | Titik masuk | Rule | Catatan |
| ---: | --- | --- | ---: | --- |
| 5 | `Komite Claim Life.md` | `Flow/KomiteLife_Flow.xml` → `Start2` | 12 | routing **berbasis data** (`.KomiteID`); 55 kolom tabel akseptasi terbaca |
| 6 | `Komite Claim FacIn.md` | `Flow/Komite_Flow.xml` → `Start1` | 15 | roster dari DB + **ambang nominal** ter-hardcode (OQ-037) |
| 7 | `Komite Claim Prop.md` | `Flow/KomiteTreaty_Flow.xml` → `Start1` (class `...KOMITETREATY`) | 13 | 3 activity ber-guard nama orang |
| 8 | `Komite Claim Non Prop.md` | `Flow/KomiteTreaty_Flow.xml` → `Start1` (class `...KOMITETREATYNONPROP`) | 13 | **otorisasi berbasis roster**, tanpa hardcode nama |

`[terverifikasi]` Konteks 7 dan 8 memakai **nama file Flow yang sama persis** tetapi class berbeda →
**dua rule berbeda**, ditelusur terpisah (hash ternormalisasi `5eacdb3783` vs `63b913a161`).

## Tahap 3 — domain Claim / siklus kerugian (selesai)

Sintesis lintas-modul: **`_SUMMARY-claim.md`**.

| # | Konteks | Titik masuk | Rule | Catatan |
| ---: | --- | --- | ---: | --- |
| 9 | `Claim Life.md` | `Flow/Register_Flow.xml` → `Start2` (class `...CLAIMLIFE`) | 16 | 4 tahap; **tanpa `PaymentType`**; `ISCLM*` tidak menyentuh modul ini |
| 10 | `Claim Prop.md` | `Flow/Flow_TreatyIn.xml` → `Start1` (class `...CLAIMTREATY`) | 14 | limit dari **database** |
| 11 | `Claim Non Prop.md` | `Flow/Flow_TreatyIn.xml` → `Start1` (class `...CLAIMTREATYNONPROP`) | 15 | limit **ter-hardcode**; baca master Treaty Out |
| 12 | `Claim Fac In.md` | `Flow/Register_Flow.xml` → `Start1` (class `...WORK-PNC`) | 17 | percabangan **B2B**; 8 ConnectREST |

`[terverifikasi]` Dua pasang berbagi nama file Flow tetapi berbeda class → **empat rule berbeda**.

## Tahap 4 — domain Facultative Inward (selesai)

Sintesis lintas-modul: **`_SUMMARY-facultative.md`**.
Domain terbesar korpus: **6.071 file** (64,8 % dari 9.369).

| # | Konteks | Titik masuk | Rule | Catatan |
| ---: | --- | --- | ---: | --- |
| 13 | `NB FacIn.md` | 6 rule `Flow`; utama `Flow/InputInwardFacultativeOffer.xml` → `Start1` | 45 | **graf terbesar korpus** (21 Assignment, 42 Decision, 170 connector) |
| 14 | `RNW Fac In.md` | `Flow/InputRenewalFacultativeIn.xml` → `Start2` | 28 | 3 dari 4 Flow **identik** dgn NB FacIn; **tanpa cabang Life** |
| 15 | `Endorsment Fac In.md` | `Flow/InputAddendumFacultativeIn.xml` → `Start2` | 52 | taksonomi 12 tipe endorsement terbaca; class `ASM-SFAGIS-*` |

### Aturan telusur ulang yang dipakai di Tahap 4 `[terverifikasi]`

Dua belas file `Flow` ketiga modul hanya memuat **delapan isi berbeda**:

| Hash | File | Perlakuan |
| --- | --- | --- |
| `1306f68d56` | `InputRealizationTreatyIn` (NB FacIn = NB Treaty In) | **dipakai ulang** dari `NB Treaty In.md` |
| `c099bf4ebb` | `InputInwardFacultativeRISlip` (NB FacIn = RNW Fac In) | ditelusur **sekali** di `NB FacIn.md` §1.3 |
| `832fb12b9d` | `OfferFacOut` (NB FacIn = RNW Fac In) | ditelusur **sekali** di `NB FacIn.md` §1.4 |
| `bfd6252070` | `OfferFacRetro` (NB FacIn = RNW Fac In) | varian A — `NB FacIn.md` §1.5 |
| **`b370146c63`** | `OfferFacRetro` (**Endorsment Fac In**) | varian B — **ditelusur terpisah**, `Endorsment Fac In.md` §2 |

`OfferFacRetro` ditelusur **dua kali** (OQ-011 #338). Selisihnya nyata: tangga persetujuan fac out
**2 anak tangga** (varian A) vs **4** (varian B).

## Tahap 5 — modul TANPA rule `Flow` (selesai) — **PENUTUP D2**

Konvensi tambahan: **`_METHOD-noflow.md`** — titik masuk Harness/FlowAction, dan cara
merekonstruksi mesin status dari tombol + guard ketika tidak ada graf `Flow`.
Sintesis lintas-modul: **`_SUMMARY-treaty-master.md`**.

| # | Konteks | Titik masuk (Harness) | Rule | Catatan |
| ---: | --- | --- | ---: | --- |
| 16 | `Treaty In.md` | `Harness/InputTreatyInOffer.xml` (`DATA-PORTAL`, **8,2 MB**) | 24 | mesin status di **`Akseptasi_DT`**; 10 tombol `(dev)` |
| 17 | `Treaty In Adjustment.md` | `Harness/InputTreatyInAdjustment.xml` | 31 | **277/323 file identik** dgn Treaty In; 23 layar `OldData` |
| 18 | `Treaty Contract Out.md` | `Harness/InboxTreatyContract.xml` | 26 | **OQ-022 terjawab: BUKAN outward** |
| 19 | `Master Product Name Life.md` | `Harness/InwardProductName.xml` (class entitas) | 22 | membawa jalur tulis `M_TREATY_IN` |
| 20 | `Master Contract Retro Life.md` | `Harness/InboxRetroLifeReinsurersList.xml` | 23 | modul terkecil; tanpa `ConnectREST`/`When`/JSON |

`[terverifikasi]` Kelimanya tanpa rule `Flow`, tetapi **terbelah dua bentuk**: `Treaty In` +
`Treaty In Adjustment` tetap punya proses berjenjang (`StatusAkseptasi` 4 nilai, tangga 4 tingkat)
yang ditulis sebagai **Data Transform bersarang**, bukan graf; tiga modul lain memang editor master
murni (nol rujukan `StatusAkseptasi`).

## Format

Setiap `<konteks>.md` memakai **7 bagian** sesuai `_METHOD.md` §3: diagram alur, status/state,
objek Oracle, integrasi eksternal, batas pengetahuan, daftar rule, pertanyaan terbuka.

**Empat** berkas memakai **8 bagian** karena menampung satu pemeriksaan tambahan
(`grep -c '^## [0-9]' *.md`):

| Berkas | Bagian tambahan |
| --- | --- |
| `Endorsment Fac In.md` | §2 — varian `OfferFacRetro` yang berbeda isi |
| `Treaty In.md` | §2 — mesin status `Akseptasi_DT` |
| `Treaty In Adjustment.md` | §1 & §2 — pengukuran selisih terhadap `Treaty In`, lalu mesin revisi |
| `Master Product Name Life.md` | §2 — uji "apakah membawa salinan subtree `Treaty In`" |

`Treaty Contract Out.md` tetap 7 bagian; uji nama modulnya ditampung di **§2** menggantikan
bagian status (modul itu tidak punya properti status proses).

## Langkah berikutnya

| Step | Keluaran |
| --- | --- |
| **STEP D3** | `../modules/<modul>.md` ×20 + finalisasi `../glossary-draft.md` |
| STEP D4 | `../context-map.md` + `../understanding-report.md` → **GATE manusia** |

**FASE B (ADR/BRD/spec/tiket) tidak boleh dimulai sebelum gate setelah D4.**
