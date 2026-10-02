# Ringkasan Eksekutif — Discovery Facultative Inward

> **Sumber tunggal:** korpus ekspor rule Pega di
> `D:\migrasi\RNM\{NB FacIn, RNW Fac In, Endorsment Fac In}\` — **6.071 berkas `.xml`**, READ-ONLY.
> Discovery ini dibangun **dari nol atas korpus mentah** (keputusan work owner, 15 September 2026);
> hasil discovery terdahulu tidak dipakai sebagai sumber.
>
> Label mengikuti `CLAUDE.md` §3: `[terverifikasi]` (tag/SQL dikutip langsung) · `[dugaan]` (dari
> pola/nama) · `[pertanyaan terbuka]` (tidak terjawab dari korpus). Setiap angka disertai perintah
> audit yang menghasilkannya.

---

## 1. Apa yang dimigrasikan

Aplikasi penerimaan risiko **reasuransi fakultatif masuk** (Facultative Inward). Alur intinya sama
di ketiga siklus:

```
Marketing → gerbang akseptasi → tangga persetujuan berbasis limit → binding
          → R/I slip → fac out / retro → konversi ke produksi
```

Target: **Go (backend) + React (frontend) + Oracle (skema existing, tidak berubah)**.

---

## 2. Inventaris korpus

```powershell
foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  "{0,-20} {1}" -f $f, (Get-ChildItem "D:\migrasi\RNM\$f" -Recurse -Filter *.xml -File | Measure-Object).Count }
# -> NB FacIn 2083 / RNW Fac In 1927 / Endorsment Fac In 2061   (total 6071)
```

| Tipe rule | Total | | Tipe rule | Total |
| --- | ---: | --- | --- | ---: |
| Activity | 1.747 | | Harness | 119 |
| Section | 1.263 | | DataPage | 101 |
| FlowAction | 754 | | DecisionTable | 32 |
| RDBList | 605 | | Flow | 12 |
| When | 601 | | ConnectREST | 9 |
| DataTransform | 443 | | DecisionTree | 3 |
| ReportDefinition | 379 | | SystemSettings | 3 |
| | | | **Total** | **6.071** |

---

## 3. Tujuh temuan utama

### 3.1 Tiga siklus berjalan di atas **satu basis rule** — dan NB = RNW sepenuhnya

`[terverifikasi]` Dari **2.439 identitas rule unik** (Tipe + Nama), **1.680 hadir di ketiga folder**.
Dari 1.680 itu, **1.220 logikanya identik** dan 460 bercabang.

```powershell
# skrip lengkap + tiga koreksi metodologis wajib: 05-migrasi/03-risiko-dan-pertanyaan.md §5
# -> identitas di 3 folder 1680 ; identik 1220 ; bercabang 460 ; NB==RNW 1680 dari 1680
```

⚠️ **Seluruh percabangan adalah EDM vs NB+RNW. NB dan RNW tidak pernah berbeda satu sama lain** —
1.680 dari 1.680. Di lapisan UI angkanya mutlak juga: 575 dari 575 identitas UI bersama identik.

**Konsekuensi:** satu layanan dengan `StatusBusiness` sebagai atribut case. Membangun tiga layanan
menggandakan mayoritas kode tanpa dasar.

⛔ **Angka 460 tidak boleh dibaca sebagai "percabangan bisnis".** `[terverifikasi]` Dari 1.571 rule
bernama sama, **86 punya `pyRuleSetVersion` berbeda** antara NB dan EDM — dan **arahnya tidak
konsisten**: satu rule di EDM bernomor versi lebih lama (2025-09) dan kehilangan sebuah cabang
routing, sementara rule lain di EDM justru lebih baru dan menambah klausa persetujuan. Ketiga folder
diekspor dari **titik waktu yang berbeda-beda**.

Artinya sebagian dari 460 itu adalah **selisih versi ekspor**, bukan percabangan bisnis — dan keduanya
**tidak dapat dipisahkan dari korpus ini**. Kesimpulan "satu basis rule" tetap berdiri (ditopang 1.220
yang identik dan NB=RNW 1.680/1.680); yang tidak berdiri adalah pemakaian angka 460 sebagai ukuran
kerja percabangan. **Rekonsiliasi paralel run mustahil sampai ada satu ekspor produksi tunggal yang
konsisten untuk ketiga siklus.**

`[terverifikasi]` Diskriminator siklus, dikutip langsung:

| Nilai | Rule + baris | Ekspresi |
| ---: | --- | --- |
| `1` | `NB FacIn\When\IsNB.xml:163` | `pyWorkPage.Quotation.StatusBusiness = 1` |
| `3` | `Endorsment Fac In\When\IsEDM.xml:169` | `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` |

⚠️ Perhatikan **jalur properti berbeda**, bukan sekadar beda nilai — `Quotation.StatusBusiness` vs
`OfferFacIn.QuotationData.StatusBusiness`. Apakah keduanya selalu sinkron adalah
`[pertanyaan terbuka]` yang **memblokir**.

### 3.2 Tangga akseptasi adalah mesin keadaan **satu langkah**, bukan loop

`[terverifikasi]` (`01-flow/04-mesin-akseptasi.md` §0) Satu keputusan manusia = satu putaran:

1. Tombol ditekan → DataTransform menulis `.ProposalAcceptStatus`.
2. Post-activity memanggil `GetLimitAkseptasi_Act` → SQL ke tabel limit → menulis **satu** properti:
   `pyWorkPage.LetterNo` (kode jabatan tujuan berikutnya).
3. `DecisionTable/IsUWAccepted` memetakan status → konektor flow.
4. Gerbang "Limit Akseptasi" membaca `LetterNo` lewat rule `When` bernama `To*` → mengarahkan case.
5. Bila tidak ada `To*` yang cocok → cabang `Else` → **tangga selesai, bukan galat**.

**Konsekuensi:** di Go ini mesin keadaan yang dipanggil per transisi. Merancangnya sebagai loop yang
menghitung seluruh rantai approver sekaligus menghasilkan perilaku berbeda saat rantai terputus.

### 3.3 `LetterNo` tidak menyimpan nomor surat

`[terverifikasi]` Ia menyimpan **kode jabatan approver berikutnya**. Nama kanonik di target:
**`next_approver_position`**.

`[terverifikasi]` Alur digerakkan **tiga** field state, bukan satu "status":

| Properti | Isi |
| --- | --- |
| `PositionNote` | antrean yang **sedang** memegang case |
| `LetterNo` | antrean **tujuan berikutnya** |
| `ProposalAcceptStatus` | hasil **keputusan terakhir** |

⚠️ **Dua ruang nama yang tidak boleh disatukan:** token antrean (`ReasFacInUnderwriting`) dan kode
jabatan (`SENIORUW`). Menyatukannya merusak routing.

### 3.4 Uang tidak pernah `float` — buktinya berlapis

`[terverifikasi]` (`04-aturan/02-formula-dan-status.md` §7)

- **1.057 kemunculan di 91 berkas** konversi koma→titik. Komentar pengembang lama menyatakannya
  eksplisit: *"karena property bukan decimal"*.
- **8 bentuk perbandingan ambang uang sebagai string berkutip** — dan ambang **yang sama** (30
  miliar) dibandingkan numerik di rule lain.
- Produk antara rumus premi mencapai orde **10²³**; `float64` punya 15–17 digit signifikan.
- Mata uang adalah dimensi eksplisit (`CurrencyList`).

⚠️ **Presisi pembulatan adalah parameter per-rule.** Aritmetika yang sama dibulatkan **4 desimal** di
`HitungPremi_FacInDT` dan **20 desimal** di `CountPremi_ACT`. Menyeragamkannya mengubah angka.

### 3.5 Endorsement dinilai dari **selisih**, bukan nilai absolut

`[terverifikasi]` Nilai dasar akseptasi berbeda antar siklus: NB dan renewal memakai TSI penuh;
endorsement memakai selisih terhadap before-image. Aturan bisnis yang mengikuti: endorsement kecil
tidak naik ke direktur, endorsement besar tidak lolos.

⚠️ `[terverifikasi]` **Tiga mekanisme "nilai lama" berbeda** hidup berdampingan dan tidak boleh
disatukan jadi satu konsep. Riwayat versi polis (`PRODKE`, *append-only*) berbeda dari before-image
(salinan kerja nilai sebelumnya).

### 3.6 ⛔ `ProposalAcceptStatus` — arti nilainya **tidak ada di korpus**

`[terverifikasi]` `DecisionTable/IsUWAccepted` mendeklarasikan enam hasil (`confirm`, `reject`,
`ask`, `banding`, `revise`, `decline`); Data Transform menulis `1`/`2`/`3`/`4`/`7`/`9`.
**Baris pemetaannya tidak ikut terekspor.**

⚠️ Bukti bahwa menebak berbahaya: nilai `4` ditulis rule bernama *banding* tetapi **diuji** rule
bernama *fac out* — `RNW Fac In\When\IsFacout.xml:148` → `pyWorkPage.ProposalAcceptStatus = 4`.
Satu nilai, dua nama pemakai bertema berbeda. **Nama rule tidak dapat dipakai menebak arti.**

Ini butir yang memblokir paling luas: menebaknya menyalahkan arah **setiap** keputusan underwriting.

### 3.7 Koreksi: `CLAUDE.md` §4.5 keliru — tidak ada rule `When` yang kondisinya kosong

`[terverifikasi]` §4.5 memerintahkan `panic()` untuk lima rule yang "kondisinya kosong di ekspor".
Pemeriksaan ulang atas 601 berkas: jumlah rule yang kondisinya **tidak terbaca sama sekali = 0**.

| Rule | Kondisi terbaca |
| --- | --- |
| `IsPKSASM` | `OfferFacIn.IsB2B = "ASM"` |
| `ToUW` | `LetterNo = "UNDERWRITER"` |
| `ToJUW_A` | `LetterNo = "JUW_A"` |
| `LetterNoNull` | `LetterNo = ""` |
| `IsEdmInternalRetro` | `QuotationData.EndorsementInternalRetro = 1` |

**Penyebab kekeliruan:** discovery terdahulu hanya membaca `<pyConditionString>`. **170 berkas
(28,3 %) menyembunyikan kondisinya** dari tag itu — kondisinya ada di `<pyConditionValue1String>`.
`ToUW` dan `ToJUW_A` bahkan punya **kedua** tag terisi dan sepakat — tidak pernah kosong dalam
pengertian apa pun.

⚠️ `IsPKSASM` tetap perlu konfirmasi: baris kondisinya membawa `<pyTempText>true</pyTempText>` dengan
label belum ter-resolve. Rincian dan usulan perubahan `CLAUDE.md` ada di
`02-data/03-batas-pengetahuan.md` §9.

---

### 3.8 Lapisan data: tiga hal yang mengubah rencana migrasi

`[terverifikasi]` (`02-data/02-skema-oracle.md`) Dari **605 berkas SQL** / 619 blok / **112 tabel
unik** — dan **0 SQL inline** di luar `RDBList\`, sehingga inventarisnya lengkap, bukan sampel:

**a. `HISTORYAKSEPTASIPEGA` di-SELECT untuk mengambil keputusan, bukan hanya di-INSERT.**
`GetAksepBanding_SQL` mengambil workbasket terakhir dari riwayat lalu men-join ke tabel limit untuk
memperoleh jalur banding; `GetFlagReject_SQL` menghitung baris berstatus reject. **Tabel ini tidak
boleh diarsipkan atau dipangkas** — memangkasnya mengubah keputusan.

**b. Prompt AI adalah data tabel, bukan kode.** `M_PROMPT_AI` memasok prompt yang menggerakkan model
AI di alur underwriting. Perilaku model **dapat diubah tanpa deployment dan tanpa jejak version
control**. Ini memperluas cakupan tinjauan keamanan `CLAUDE.md` §6 dari rule ke **tabel**.

**c. Jalur produksi endorsement menulis data pribadi.** `FACINPRODUCTION` menulis 79 kolom di NB/RNW
tetapi **82 di EDM** — dan bukan sekadar tambahan: EDM menulis `NO_KONTRAK`, `NO_NPWP`, `NO_KTP`,
`STARTDATE_DEBITUR`, `ENDDATE_DEBITUR` **sebagai pengganti** `DEDUCTION2_MENJADI`/`_SELISIH`.

⚠️ Bukti ketiga bahwa **nama bukan bukti**: dua rule bernama `Insert…` sesungguhnya **menghapus lalu
menyisipkan ulang**, dan deskripsi langkah `RDB-Delete`-nya sendiri berbunyi *"insert to db"*.

---

### 3.9 Endorsement: struktur alurnya berbeda, bukan sekadar variasi

`[terverifikasi]` (`01-flow/03-alur-endorsement.md`)

- **EDM tidak punya fase Offer maupun Binding.** Seluruh 30 penetapan `IsCedingConfirm` /
  `FlagOnGoingPolicy` di flow EDM bernilai `Policy`/`1` — tidak ada satu pun `Offer`/`Binding`.
  EDM juga tidak punya flow siklus polis terpisah.
- **Before-image ada tiga lapis** (dokumen JSON polis versi terakhir · properti `*Old` per baris ·
  penanda `.IsOldData`), dan ketiganya berbeda lagi dari riwayat versi polis. Menyatukannya jadi satu
  konsep akan salah.
- **Selisih dihitung di dua tingkat** — tingkat pembayaran dan tingkat baris produksi, dengan
  pemetaan kolom `_MENJADI`/`_SELISIH` yang terverifikasi.
- **Retro EDM punya empat tingkat persetujuan**, berbanding dua tingkat di NB/RNW.
- ⛔ **Endorsement jiwa melewati seluruh tangga akseptasi**, dan rule tangga versi jiwa tidak ada di
  korpus EDM.

⚠️ **Jebakan label `<pyMOName>` terbukti persis seperti diperingatkan `CLAUDE.md` §3 butir 3:** dua
shape di flow EDM sama-sama berlabel *"Err Konversi?"* — yang satu memang menggerbangi pemeriksaan
keberhasilan servis, yang lain sesungguhnya menggerbangi **`IsFacRetro`**. Akibatnya jalur simpan JSON
tidak terjangkau untuk case ber-fac-retro.

---

## 4. Yang memblokir — 45 butir

Daftar lengkap di `02-data/03-batas-pengetahuan.md`. Yang paling mahal bila salah:

| # | Butir | Dampak bila ditebak |
| ---: | --- | --- |
| 0 | **Ketiga folder diekspor dari titik waktu berbeda** (86 rule beda versi, arah tidak konsisten) | **Rekonsiliasi paralel run mustahil** — selisih versi tidak terpisahkan dari percabangan bisnis |
| 1 | Arti `ProposalAcceptStatus` | Arah **setiap** keputusan underwriting salah |
| 2 | Ambang uang dibandingkan sebagai string | Tangga akseptasi salah cabang pada rentang tertentu |
| 3 | Satuan `.Rate` (per mille vs persen) | Premi meleset **10×** |
| 4 | Isi tabel limit `M_LIMIT_*` + ejaan kolom `JABATAN` | Tangga tidak dapat direkonsiliasi sama sekali |
| 5 | Sumber pengganti tabel internal Pega (riwayat akseptasi) | Jalur banding kehilangan sumber keputusan |
| 6 | 4 rule + 12 DecisionTable + 21 Section hilang dari ekspor | Special Acceptance & konversi produksi tidak dapat dimigrasikan |
| 7 | Tipe kolom seluruh tabel — korpus memuat **nol DDL** | Bila `FACINOFFER` ternyata `VARCHAR2`, seluruh perbandingan angka atasnya adalah perbandingan string |
| 8 | Isi **24 stored procedure `POOLDATA`** | **Seluruh jalur tulis melewatinya** — mustahil direplikasi |
| 9 | Prompt AI di tabel `M_PROMPT_AI` | Perilaku model dapat berubah tanpa deployment dan tanpa jejak |

⚠️ **Dua butir punya tenggat keras, dan keduanya tidak dapat dibatalkan setelah lewat:**
butir 6 (ekspor ulang rule) dan butir 5 (ekstraksi data historis riwayat akseptasi). Keduanya hanya
mungkin **selama Pega masih hidup**.

---

## 5. Risiko metodologis — tiga jebakan yang menghasilkan angka salah

Ketiganya **terjadi** selama discovery ini sebelum dikoreksi.

| # | Jebakan | Akibat | Koreksi |
| --- | --- | --- | --- |
| M1 | Membaca kondisi `When` hanya dari `<pyConditionString>` | 170 berkas tampak "tanpa kondisi" — sumber kekeliruan §3.7 | Baca **kedua** tag |
| M2 | `Get-FileHash` mentah antar folder | Melaporkan **0** dari 1.680 identik; penyebabnya metadata ekspor | Buang baris `px*`/`pz*` |
| M3 | Tidak mengurutkan, atau memakai `Sort-Object` bawaan | Urutan tag ekspor Pega **tidak deterministik**; angka berubah 507 → **1.220** setelah dikoreksi | Urutkan `[StringComparer]::Ordinal` |

⚠️ `[terverifikasi]` Nilai `pxHostId` berbeda (`jboss122117`, `jboss1073`) mengonfirmasi
`CLAUDE.md` §3 butir 1 — korpus dirakit dari **lebih dari satu server Pega**. Menyebut nama rule
saja tidak cukup; **path berkas wajib**.

---

## 6. Kepatuhan penulisan

`[terverifikasi]` **Nol nama orang** di seluruh dokumen keluaran. Diperiksa dengan mengekstrak 232
nama operator dan 211 identitas login dari metadata korpus, lalu mengujinya terhadap seluruh berkas
`OUTPUT\` dengan pencocokan batas kata. Dua kecocokan yang muncul adalah token non-orang: `System`
(operator bawaan Pega) dan `Manager` (jabatan). Skrip: `05-migrasi/03-risiko-dan-pertanyaan.md` §3.

⚠️ `[terverifikasi]` **Literal yang seharusnya konfigurasi** (`CLAUDE.md` §4.4): tiga alamat email
tertanam di satu activity konversi produksi EDM, dan **629 kemunculan / 176 nomor polis unik** di 8
berkas Activity endorsement. Di sisi baik, **0 dari 619 blok SQL** memuat literal `http(s)://` —
alamat servis memang sudah lewat tabel lookup. Nilai-nilai itu tidak disalin ke dokumen mana pun.

⚠️ `[terverifikasi]` Sistem lama **menampilkan nama orang ke pengguna**: 669 kemunculan / 23 label
status, plus 92 kemunculan pola *"IS IN … INBOX"* dengan 14 nama unik di satu flow saja. Di sistem
baru ini harus menjadi rujukan ke antrean/jabatan — **keputusan bisnis** apakah teks yang dilihat
pengguna boleh berubah.

---

## 7. Peta dokumen

| Butuh | Baca |
| --- | --- |
| **Mesin tangga persetujuan** (inti aplikasi) | `01-flow/04-mesin-akseptasi.md` |
| Alur per siklus | `01-flow/01-alur-new-business.md`, `02-alur-renewal.md`, `03-alur-endorsement.md` |
| Model domain & agregat `OfferFacIn` | `02-data/01-model-domain.md` |
| Tabel & kolom Oracle | `02-data/02-skema-oracle.md` |
| **Yang tidak diketahui dari korpus** | `02-data/03-batas-pengetahuan.md` |
| Layar & komponen React | `03-ui/01-inventaris-layar.md` |
| Katalog 226 rule `When` | `04-aturan/01-katalog-when.md` |
| Rumus premi & enumerasi status | `04-aturan/02-formula-dan-status.md` |
| Arsitektur target & kontrak API | `05-migrasi/01-arsitektur-target.md`, `02-kontrak-api.md` |
| **Risiko + pertanyaan untuk bisnis** | `05-migrasi/03-risiko-dan-pertanyaan.md` |

---

## 8. Cakupan yang **belum** didiscovery

Jangan dianggap "tidak ada".

| Area | Status |
| --- | --- |
| **Spreading · capacity · scoring** | Tersentuh lewat rumus dan proteksi spreading; **belum ada dokumen alurnya sendiri** |
| Isi 460 rule bercabang EDM | Terhitung, **belum dipetakan field-per-field**. Sebelum dipetakan, tidak boleh diklaim "hanya beda kosmetik" |
| `ReportDefinition\` (379 berkas) | Belum dibaca sama sekali |
| 141 Section yatim | Termasuk rumpun `*FacOut*` dan `*TreatyIn*` — apakah termasuk lingkup Fac **In**? |
| `AnalysLocationbyAI`, `AttachDoc_AI` | Ada di NB+RNW, tidak di EDM. **Isinya belum dibaca.** `CLAUDE.md` §6 mengharuskan tinjauan keamanan — dan cakupannya kini meluas ke tabel `M_PROMPT_AI` (§3.8b) |

---

*`CLAUDE.md` §1: perbaikan dipisahkan dari migrasi. Sistem lama memuat hal yang tampak keliru —
perbandingan angka sebagai string, `ProRateType != 3` yang menihilkan premi, guard berbasis identitas
orang. Semuanya kandidat perbaikan, tetapi memutuskannya milik bisnis. Migrasi yang diam-diam
"memperbaiki" menghasilkan selisih angka yang tidak dapat dijelaskan saat rekonsiliasi paralel run.*
