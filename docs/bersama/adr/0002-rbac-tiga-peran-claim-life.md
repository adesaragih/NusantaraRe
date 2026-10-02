---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 1 Q4, keputusan work owner
---

# RBAC Claim — Life memakai tiga peran yang sudah ada; rangkap peran ditolak

Model peran untuk Claim — Life **dirancang**, bukan dimigrasikan — korpus Pega tidak memuat satu pun
rule identitas atau otorisasi (**OQ-007**). Bahannya diambil dari tiga **kode peran** yang sudah
dipakai modul ini, dikonfirmasi work owner sebagai **daftar lengkap**. **Rangkap peran tidak
diperbolehkan**, kecuali akses ditambahkan eksplisit pada role akun.

## Peran dan pemetaan tahap

`[terverifikasi work owner 2026-09-14]`

| Tahap (shape pada `Register_Flow`) | Peran |
| --- | --- |
| Register (`Assignment2`) + Outstanding (`Assignment1`) | **`ReasLifeAdmin`** |
| Medical Check (`Assignment3`) | **`ReasLifeMedicalAdvisor`** |
| Claim Analis (`Assignment4`) | **`ReasLifeSPV`** |

Jalur balik yang melekat pada peran: `SendtoAdmin = 1` mengembalikan kasus dari
`ReasLifeMedicalAdvisor` **atau** `ReasLifeSPV` ke `ReasLifeAdmin`; `SendtoMedical = 1`
mengembalikan dari `ReasLifeSPV` ke `ReasLifeMedicalAdvisor`.

## Mengapa tiga kode ini, bukan peran baru

`[terverifikasi]` Claim — Life **satu-satunya konteks di korpus yang otorisasinya sudah berbasis
peran, bukan identitas orang**:

- `pyPosition` dibandingkan terhadap `'ReasLifeAdmin'`, `'ReasLifeSPV'`,
  `'ReasLifeMedicalAdvisor'` di **17 berkas** `Claim Life`, mis.
  `pyPosition != 'ReasLifeAdmin' || pyWorkPage.ClaimData.…` (8 kemunculan),
  `pyPosition == 'ReasLifeMedicalAdvisor'` (5), `pyPosition == 'ReasLifeSPV'` (6).
- `OperatorID.pyUserIdentifier` (2 berkas) dan `pyUserName` (4 berkas) dipakai sebagai **data jejak**,
  **bukan** guard terhadap literal nama orang:
  ```
  grep -rhoE "(pyUserIdentifier|pyUserName)[^<]{0,45}" "Claim Life" --include="*.xml" \
    | sed 's/&amp;#61;/=/g;s/\]\[/ /g;s/[][]//g' | grep -E "[=!]"      # -> kosong
  ```

Bandingkan konteks lain: identitas orang ter-hardcode di `NB FacIn` (43 berkas),
`RNW Fac In` (36), `Endorsment Fac In` (33), `NB Treaty In` (12), `Komite Claim FacIn` (4),
`Komite Claim Prop` (3) — dan **OQ-053**, di mana 5 identitas orang justru *ditetapkan* sebagai
pemilik tugas berikutnya.

**Koreksi artefak:** nilai `IT Developer` yang tercatat di D1 untuk `pyPosition` berasal dari sapuan
korpus-wide — **tidak ada di `Claim Life`**.

## Consequences

- Otorisasi di sistem baru memakai **tiga peran ini sebagai himpunan kanonik**; penambahan peran
  adalah keputusan bisnis, bukan detail implementasi.
- Larangan rangkap peran harus **ditegakkan sistem**, bukan sekadar konvensi — dan pengecualiannya
  (“akses ditambahkan eksplisit di role akun”) perlu jalur yang terdefinisi.
- Pemetaan tahap → peran ini **tidak dapat digeneralisasi** ke konteks lain: facultative dan treaty
  inward memakai nama workbasket (`ReasFacIn*`, `ReasTreatyIn*`), bukan `pyPosition`.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-007** | Tidak ada rule identitas/otorisasi di korpus — model peran korpus-wide tetap harus dirancang |
| **OQ-021** (terjawab untuk Claim — Life) | Identitas orang ter-hardcode di 6 modul lain **tetap terbuka** |
| **OQ-024** (maju sebagian) | Pemetaan Assignment → nama workbasket di facultative & treaty tetap tidak terbaca |
| **OQ-028** | `WorkList` vs `WorkBasket` — Claim Life memakai **campuran**: Register/Outstanding `WorkList`, Medical Check/Claim Analis `WorkBasket`. Arti perbedaannya belum terverifikasi |
