# Tangga akseptasi — apa yang sudah terpetakan, dan apa yang menahan NB-11, NB-12, NB-14

**Untuk:** work owner (meneruskan ke pemilik aplikasi Pega, DBA, dan Underwriting) · **Tanggal:** 1 Oktober 2026

NB-13 (domain hasil keputusan + flag fase) **sudah diport** (`backend/services/acceptance/keputusan.go`).
Tangga akseptasi sendiri (NB-11) **belum**, dan NB-12/NB-14 bergantung padanya. Sebabnya bukan kurang
waktu: tiga hal di bawah menentukan perilaku, dan korpus tidak menjelaskannya. Mengisinya dengan
tebakan akan memindahkan approver kasus ke jabatan yang salah.

Semua berkas di bawah: `D:\migrasi\RNM\NB FacIn\`.

---

## A. Yang sudah terpetakan

`[terverifikasi]` dicek langsung 01-10-2026: SQL L80, langkah 2–20 `_ActFlow`, tautologi (15 per activity,
dihitung dua cara: grep baris dan pengurai XML — sepakat), nomor polis L4688/L5109, efek Revise, dekoder.
Baris yang ditandai lain berasal dari pemetaan agen penjelajah dan belum dicek ulang.

| Bagian | Asal |
| --- | --- |
| Flow Offer memanggil `GETLIMITAKSEPTASI_ACTFLOW` (shape `Utility2`); RI Slip memanggil `GETLIMITAKSEPTASI_ACT`; tingkat JUW memakai `GETLIMITAKSEPTASI_JUW_UW` | `Flow\InputInwardFacultativeOffer.xml` (nomor baris dari pemetaan agen, **belum dicek ulang**) |
| Langkah 2 `_ActFlow`: `LetterNo = ""`, `DataSearch.CARID2 = round0(TotalTSINusaRe)` (nilai dasar, K-015); langkah 4 menimpanya dengan `TotalTSITopRisk` bila ada top risk; langkah 5 selisih hanya bila `StatusBusiness == "3"` | `Activity\GetLimitAkseptasi_ActFlow.xml` |
| Langkah 6–10: satu `RDB-List` per kelompok lini (Preferred / Non-Preferred / Preferred Commercial / Engineering / Non-Prop&Eng) | idem |
| SQL bentuk A: `SELECT JABATAN, MAX_LIMIT_IDR, MAX_LIMIT_USD, BATAS_WAKTU, NAMA … WHERE team_group = {TeamGroup} AND LIMIT_BOTTOM > (SELECT LIMIT_BOTTOM … WHERE … LOGIN = {OperatorID.pyUserIdentifier}) AND LIMIT_BOTTOM < {CARID2} ORDER BY LIMIT_BOTTOM ASC` | `RDBList\GetLimitAkseptasi_SQL.xml` L80 (+13 varian JUWA/Banding) |
| Langkah 20: `LetterNo = LimitAkseptasi.pxResults(1).CARI1` — baris pertama hasil = jabatan berikutnya | `GetLimitAkseptasi_ActFlow.xml` L7544 |
| Antrean (`PositionNote`) ditulis **flow**, bukan activity: konektor `Decision23` memetakan `LetterNo` → antrean (`ToDepHeadUW` → `ReasFacInDepHeadUnderwriting`, `ToManagerTeknik`, `ToDirMarketing`, `ToDirTeknik`, `ToKadivTeknik` → `ReasFacInGroupLeader`, `ToKadivFin`); **`ToKadivFacultative` tidak menulis antrean**; cabang Else = tangga selesai | `Flow\InputInwardFacultativeOffer.xml` konektor `Decision23` (nama `To*` mulai L1653; nomor baris tiap pemetaan dari pemetaan agen, **belum dicek ulang**) |
| Dekoder hasil keputusan, penulisnya, efek samping Reject **dan Revise** | NB-13, sudah diport |
| 15 tautologi `TotalTSI >= Limit \|\| TotalTSI <= Limit` di `_ActFlow` (langkah 19.1.2–19.1.9, 21.1.2–21.1.8) dan 15 lagi di `_Act` | lihat `issues/11-…` bab Comments |
| Empat nomor polis literal `OldPolicyNo = "RNM-F29.05.2023.02880"/"…02879"/"…02877"/"…02883"` di langkah 17–18 | `_ActFlow` L4688, L5109 |

## B. Yang menahan — mohon dijawab

### B1. Arti kode transisi `5` — pemilik aplikasi Pega

Langkah 19 (blok tangga) dan langkah 20 (`LetterNo = pxResults(1).CARI1`) bergerbang:
`IsNonPropertyandNonEngineering` → **T=5**, F=2; lalu `IsPropertyandEngineering` → T=2, F=5 / F=3.
`[terverifikasi]` audit 02 mencatat kode 2 = jalankan, 3 = lewati `[dugaan kuat]`, 6 = keluar
`[dugaan kuat]` — **arti kode 5 tidak tertambat**. Tanpa itu, tidak dapat dipastikan kapan blok 19
berjalan dan kapan langkah 20 menimpanya — padahal itulah yang menentukan approver berikutnya.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Buka langkah 19 `GetLimitAkseptasi_ActFlow` di UI Pega dan lihat tindakan pada dropdown prakondisinya.
> Daftar tindakan itu menerjemahkan **semua** kode sekaligus (2, 3, 5, 6) — satu tangkapan layar dapat
> mengunci arti ratusan langkah serupa di seluruh korpus. Dugaan: kode 5 adalah tindakan loop (mis. "keluar
> iterasi"), karena langkah 19 memang loop — `[dugaan]`, tetap harus dilihat.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` — belum dijawab; letaknya ditanyakan
> Letak kode `5` di `Activity\GetLimitAkseptasi_ActFlow.xml`: **23 kemunculan** (22 "bila benar",
> 1 "bila salah") — langkah 12, 13, 19 (×2), 19.1.2/.3/.5, 20, 21.1.2–21.1.8, 22 (×4), 23 (×4). Termudah:
> **langkah 12, prakondisi 1 `IsLimitSBondKBG`, "bila benar" = 5** (L3477). Buka langkah itu di UI Pega dan
> lihat nama tindakan pada dropdown "bila benar" — **masih menahan NB-11**.


### B2. Ejaan jabatan berspasi atau tidak — DBA

Di activity yang sama, blok tangga mencocokkan `.CARI1` dengan literal **berspasi**
(`"DIREKTUR TEKNIK"` L4635, `"SENIOR UNDERWRITER"`, …) dan juga **tanpa spasi** (`"KADIVFACULTATIVE"`
L4929). K-023 menyatakan isi `JABATAN` bentuk A tanpa spasi. Bila benar, cabang berspasi **tidak pernah
cocok**. Mohon pastikan ejaan nilai `JABATAN` di kelima tabel `M_LIMIT_*` bentuk A.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Satu query tanpa nama orang: `SELECT DISTINCT JABATAN FROM` kelima tabel `M_LIMIT_*` bentuk A. Bila
> hasilnya tanpa spasi (sesuai K-023), cabang berspasi di blok tangga adalah kode mati — diport apa adanya
> dan ditandai kandidat perbaikan.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 31
> **Isi `JABATAN` sudah sesuai: ada yang berspasi dan ada yang tanpa spasi.** Kedua cabang hidup;
> tidak ada kode mati; dicocokkan persis per literal.


### B3. Guard berbasis identitas `LOGIN` — IAM / work owner

`[terverifikasi]` 9 dari 14 rule SQL bentuk A membaca `LOGIN = {OperatorID.pyUserIdentifier}` untuk
mengambil `LIMIT_BOTTOM` milik pengguna yang sedang login. Dihitung dua cara pada 01-10-2026 atas
`RDBList\GetLimitAkseptasi*_SQL.xml` + `GetLimitAccEngineering*_SQL.xml` tanpa varian Bond/Kredit/Life
(14 berkas): `grep -l LOGIN` → 9; pengurai XML atas `<pyBrowseSQL>` → 9. `GetLimitAkseptasiLife_SQL`
(tabel `M_LIMIT_LIFE`, bukan bentuk A) tidak membaca `LOGIN`. K-025 mewajibkan kolom `LOGIN` dibuang dari
fixture. Usulan port: limit pengguna menjadi **masukan** tangga (didapat dari model peran, OQ RBAC),
bukan dicari lewat login. Mohon persetujuan, karena ini mengubah sumber datanya.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> **Setujui** limit pengguna sebagai MASUKAN tangga, diambil dari model peran — bukan lewat login di query.
> Sejalan dengan K-025 (kolom `LOGIN` dibuang dari fixture) dan ADR bersama tentang aturan peran lintas
> modul; pemetaan login → jabatan tetap dipegang IAM.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 33
> **Setuju** — limit pengguna menjadi masukan tangga dari model peran.


### B4. Data tabel limit — DBA

Berkas `DDL\M_LIMIT_*.xls` berformat biner BIFF dan tidak terbaca alat yang ada. Fixture tangga butuh
isi kolom yang dibaca query (`JABATAN`, `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `MAX_LIMIT_IDR`,
`MAX_LIMIT_USD`, `BATAS_WAKTU`, `TEAM_GROUP`) dalam bentuk teks/CSV **tanpa `NAMA` dan `LOGIN`**.

## C. Koreksi atas dokumen lama — untuk dicatat, bukan ditanyakan

1. Spec Modul 4: "antrean dari `WORKBASKET`" — `[terverifikasi]` tidak satu pun dari 14 rule SQL bentuk A
   membaca `WORKBASKET` atau `JABATAN_ATASAN`; antrean ditulis flow `Decision23`.
2. K-014: efek samping Reject (`ConfirmBinding = 0`, `ReceivedRiSlip = false`) **juga** dimiliki Revise
   (`DataTransform\SetReviseProposal.xml` L179, L209).
3. Spec Modul 4: dekoder hasil keputusan "identik di ketiga folder" — salinan Endorsment
   `IsUWAccepted.xml` berbeda isi baris; NB dan RNW identik byte.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Minta ekspor **CSV** kelima tabel bentuk A dan `M_LIMIT_FINANCIALINS`, berisi hanya kolom yang dibaca
> query: `JABATAN`, `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `MAX_LIMIT_IDR`, `MAX_LIMIT_USD`, `BATAS_WAKTU`,
> `TEAM_GROUP` — tanpa `NAMA`, tanpa `LOGIN`. Berkas itu sekaligus fixture dan kontrak data ke DBA (K-009).

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 35
> Sesuai rekomendasi: CSV dari DBA. **Masih menahan** fixture tangga.

## D. Ralat 1 Oktober 2026 — dicek ulang atas pertanyaan work owner

Semua berkas di bawah: `D:\migrasi\RNM\`. Kolom dihitung per berkas lewat pengurai XML atas
`<pyBrowseSQL>`, **tak peka huruf** (SQL menulis `team_group` huruf kecil — `grep` peka huruf pertama
melewatkannya, jebakan sensus no. 4).

1. **B4 — permintaan kolom untuk `M_LIMIT_FINANCIALINS` keliru.** Tabel itu tidak punya kolom bentuk A.
   `[terverifikasi]` `DDL\M_LIMIT_FINANCIALINS.txt`: `ID`, `JABATAN`, `LIMITBOND_BOTTOM`,
   `LIMITCREDITCL_BOTTOM`, `LIMITCREDITNCL_BOTTOM`, `LIMITTRADE_BOTTOM`, `NAMA`. Yang dibaca query:
   `JABATAN` + satu kolom `LIMIT*_BOTTOM` per rule (`GetLimitAkseptasiBond_SQL`, `…KreditCL_SQL`,
   `…KreditNCL_SQL`); `LIMITTRADE_BOTTOM` tidak dibaca rule SQL NB mana pun.
2. **B4 — `DDL\M_LIMIT_*.txt` ada, tetapi isinya DDL** (`CREATE TABLE "POOLDATA"."M_LIMIT_…"`), bukan
   baris data. Isi baris tetap hanya di `.xls` biner — permintaan CSV tetap berlaku.
3. **C.1 — perlu nuansa.** Ke-14 rule bentuk A memang tidak membaca `WORKBASKET` / `JABATAN_ATASAN`. Tetapi
   rule ke-15, **`RDBList\GetAksepBanding_SQL.xml`** (atas `M_LIMIT_PROPERTYY`), membaca `JABATAN`,
   `TEAM_GROUP`, `WORKBASKET`, `JABATAN_ATASAN`. Dipakai `Activity\GetAksepBanding.xml`, dirujuk
   `Section\EmailSection.xml` dan `EmailSectionCeding.xml` — rantai pemanggilnya **belum ditelusuri**.
4. **B1 — kode `5` tidak hanya di tangga.** `[terverifikasi]` `grep -cE` per berkas `NB FacIn\Activity`:
   `GetLimitAkseptasi_ActFlow` 23, `GetLimitAkseptasi_Act` 19, `GetLimitAkseptasi_JUW_UW` 2,
   `GetLimitAkseptasiLife_Act` 1 — dan di luar tangga, mis. `SumTSIPremiSpreadedRNM_ANEKA_Act` 52,
   `CountPremiAndTSIRNMFireAnekaGolf_ACT` 24. Satu jawaban dari UI Pega berlaku untuk semuanya.
5. **`DDL\TABLEOFLIMIT`** bukan tabel tangga: kolom `BIZCODE`, `TAHUN`, `CATEGORY`, `PCTLIMIT`, dibaca
   `Activity\GetLowestPctLimit_ACT.xml` — tidak diminta untuk NB-11/12.

**Permintaan CSV yang berlaku (mengganti B4):**

| Untuk | Tabel | Kolom (tanpa `NAMA`, tanpa `LOGIN`) |
| --- | --- | --- |
| NB-11 tangga bentuk A | `M_LIMIT_PROPERTYY`, `M_LIMIT_PROPERTY_NON_PREFERREDD`, `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL`, `M_LIMIT_ENGINEERINGG`, `M_LIMIT_NONPROPANDENGG` | `JABATAN`, `TEAM_GROUP`, `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `MAX_LIMIT_IDR`, `MAX_LIMIT_USD`, `BATAS_WAKTU` (+ `WORKBASKET`, `JABATAN_ATASAN` untuk `M_LIMIT_PROPERTYY`, butir 3) |
| NB-12 bentuk B financial | `M_LIMIT_FINANCIALINS` | `JABATAN`, `LIMITBOND_BOTTOM`, `LIMITCREDITCL_BOTTOM`, `LIMITCREDITNCL_BOTTOM` |

## E. B1 terjawab — 1 Oktober 2026

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 39

| Kode | Tampil di UI Pega | Bukti |
| :-: | --- | --- |
| `2` | **Continue Whens** | prakondisi 1 "bila salah" **dan** prakondisi 2 "bila benar" — dua baris sepakat |
| `3` | **Skip Step** | prakondisi 2 "bila salah" |
| `5` | **Skip Whens** | prakondisi 1 "bila benar" |

Tangkapan layar UI Pega dari work owner, 1 Oktober 2026: `GetLimitAkseptasi_ActFlow` langkah 12, dicocokkan
dengan `NB FacIn\Activity\GetLimitAkseptasi_ActFlow.xml` (`ASM-FW-GISFW-WORK!GETLIMITAKSEPTASI_ACTFLOW`,
`pyRuleSetVersion` 01-01-87) — kedua baris prakondisi identik: `IsLimitSBondKBG` benar=5 salah=2,
`IsLimitCustomBond` benar=2 salah=3. UI tidak menampilkan kodenya — hanya nama tindakan; karena itu "isian bila benar = 5" tidak
terlihat di layar. ⚠️ Dugaan "kode 5 = tindakan loop" di rekomendasi B1 **salah**.

**Masih terbuka (U-2): kode `1`, `4`, `6`.** Langkah termudah untuk dilihat dengan cara yang sama:

| Kode | Activity (`NB FacIn\Activity\`) | Langkah | Prakondisi |
| :-: | --- | --- | --- |
| `1` | `CopyFacRetro_ACT` (`ASM-FW-GISFW-WORK`) | 1 | `Param.status=="CopyAll"`: benar = **1**, salah = 3 |
| `4` | `setDiscountRetro_act` (`ASM-FW-GISFW-DATA-SPREADINGRISK`) | 2 | `.ID==Param.id`: benar = 2, salah = **4** |
| `6` | `AddStatusCoverageNew_Act` (`ASM-FW-GISFW-WORK`) | 1 | `pyWorkPage.RetrieveEDM==1`: salah = **6** |

## F. Kode 1, 4, 6 terjawab — 1 Oktober 2026

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 40

| Kode | Tindakan di UI Pega | Kemunculan* | Sumber |
| :-: | --- | ---: | --- |
| `1` | **Jump To Later Step** | 392 | butir 40 — `CopyFacRetro_ACT` langkah 1 |
| `2` | **Continue Whens** | 48.340 | butir 39 — `GetLimitAkseptasi_ActFlow` langkah 12 |
| `3` | **Skip Step** | 17.650 | butir 39 |
| `4` | **Exit Iteration** | 122 | butir 40 — `setDiscountRetro_act` langkah 2 |
| `5` | **Skip Whens** | 1.034 | butir 39 |
| `6` | **Exit Activity** | 277 | butir 40 — `AddStatusCoverageNew_Act` langkah 1 |
| *(kosong)* | `[pertanyaan terbuka]` | 5.953 | — |

\* Jendela: seluruh `*.xml` di bawah `D:\migrasi\RNM\`, elemen `pyStepsPreCondParamsWhenTrue` dan
`…WhenFalse`. Dua cara sepakat untuk `1`–`6`: `grep -rhoE '<pyStepsPreCondParamsWhen(True|False)>[^<]*</'`
dan pengurai XML (0 berkas gagal diurai). Isian kosong hanya terhitung pengurai — `grep` melewatkannya
karena ditulis sebagai tag tutup-sendiri. Tidak ada kode lain.

**Satu sisa:** arti isian **kosong**. Cukup lihat kolom "if true" di `AddStatusCoverageNew_Act` langkah 1 —
layar yang sama dengan kode 6.

## G. Isian kosong terjawab — 1 Oktober 2026

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 41
> **Kosong = Continue Whens**, sama dengan kode `2`. Peta kode transisi selesai: tidak ada lagi yang terbuka.

## H. CSV tabel limit diterima — 1 Oktober 2026 (B4)

> ✅ **Data dari work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 42. Disimpan sebagai
> fixture `backend/services/acceptance/testdata/limit/`.

- **Format** `[terverifikasi]`: UTF-8, CRLF, pemisah `;`, tanpa `NAMA`/`LOGIN`, setiap baris selebar header.
  Angka bertitik ribuan: 541 nilai bertitik semuanya ≥ 2 titik berpola ribuan, 581 bulat polos, nol bentuk
  lain (pola `\d{1,3}(\.\d{3})+` dan `\d+` atas sembilan kolom angka) — tidak ada desimal.
- **`JABATAN`** `[terverifikasi]`: lima tabel bentuk A — 10 nilai, semuanya **tanpa spasi**
  (`DEPHEADUNDERWRITER`, `DIREKTURMARKETING`, `DIREKTURTEKNIK`, `JUW_B`, `KADIVFACULTATIVE`, `KADIVTEKNIK`,
  `LEADER`, `MANAGERTEKNIK`, `SENIORUW`, `UNDERWRITER`); `M_LIMIT_FINANCIALINS` — 4 nilai, semuanya
  **berspasi** (`DIREKTUR MARKETING`, `DIREKTUR TEKNIK`, `KADIV KEUANGAN`, `SENIOR UNDERWRITER`).
- **Letak literal berspasi** di `GetLimitAkseptasi_ActFlow` (`.CARI1 == "…"`, pengurai XML per langkah):
  17.1, 19.1.2, 19.1.8, 19.1.9, 21.1.7, 21.1.8, 22.2.2–22.2.4. Hanya **langkah 22** ("Limit Bond & KBG &
  Kredit & Trade") yang membaca daftar FINANCIALINS — cocok dengan data. Langkah **19 dan 21 berlabel `//`**
  (label = `pyStepsBlockName`, terbukti dari tangkapan layar langkah 12 berlabel `NON1`). Langkah **17**
  ("hapus …") bergerbang satu nomor kasus literal `pyWorkPage.pyID` — tambalan untuk satu kasus.
- **Label `//`** — 11 langkah (pengurai dan `grep -c '<pyStepsBlockName>//</'` sepakat): tingkat atas
  **15, 16, 19, 21**, dan 19.1.2/.3/.5/.6/.7, 21.1.2/.4. `[dugaan kuat]` `//` = langkah dinonaktifkan
  (konvensi Pega; U-6 masih terbuka).

### H1. Apakah label `//` menonaktifkan langkah? — pemilik aplikasi Pega

**Yang kami tanyakan:** di `GetLimitAkseptasi_ActFlow`, langkah **19** ("NON PROPERTY AND ENGINEERING") dan
**21** ("IsFire untuk spesial acceptance") berlabel `//`. Apakah langkah itu **tidak dijalankan**?

**Mengapa:** bila nonaktif, tangga bentuk A hanya ditentukan langkah 20, dan seluruh blok 19/21 (termasuk
literal berspasi dan 15 tautologi) tidak perlu diport. Ini juga menjawab U-6 `_DAFTAR-ISSUE-TERBUKA.md`.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Buka langkah 19 di UI Pega — cukup lihat apakah barisnya tampil sebagai remark (biasanya teks redup/abu).
> Satu kali lihat itu mengonfirmasi semua 11 langkah berlabel `//`.

## I. H1 terjawab — 1 Oktober 2026

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 43
> **Label `//` = di-remark, tidak dipakai lagi.** Langkah 15, 16, 19, 21 `GetLimitAkseptasi_ActFlow` tidak
> diport; tangga bentuk A diputuskan langkah 20. **Seluruh penahan NB-11 (B1–B4, H1) kini terjawab.**
