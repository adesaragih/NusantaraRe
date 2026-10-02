# Ringkasan — Discovery Siklus New Business

> **Sumber tunggal:** `D:\migrasi\RNM\NB FacIn\` — **2.083 berkas `.xml`**, READ-ONLY.
> Pass ini dibangun atas keputusan work owner (15 September 2026) untuk mengerjakan **satu siklus
> dahulu**. Discovery lintas-siklus sebelumnya diarsipkan di `_ARSIP-lintas-siklus\` dan **tetap
> berlaku** untuk temuan yang tidak diproduksi ulang di sini.
>
> Label mengikuti `CLAUDE.md` §3. Setiap angka disertai perintah audit yang menghasilkannya.

---

## 1. Mengapa NB aman dikerjakan lebih dulu

`[terverifikasi]` Pengukuran langsung atas `<pyRuleSetVersion>` di dalam berkas:

| Ukuran | Nilai |
| --- | ---: |
| Rule bernama sama hadir di NB **dan** EDM | 1.719 |
| — berbeda versi ruleset | **95** (62 salinan EDM lebih tua, 33 lebih baru) |
| Rule bernama sama **NB vs RNW** berbeda versi | **0** |

Contoh terverifikasi `Activity\SetToInbox_ACT.xml`: NB & RNW `01-01-95` (commit 2026-08-07), EDM
`01-01-87` (commit **2025-09-18**) — selisih **11 bulan pada rule-nya sendiri**, bukan pada tanggal
ekspor (ketiga folder diekspor dalam rentang 8 hari).

**Dua konsekuensi:**

1. Ketidaksepadanan versi adalah **masalah EDM**, bukan NB. Pekerjaan NB dapat direkonsiliasi tanpa
   menunggu pertanyaan ekspor dijawab.
2. Karena NB dan RNW nol selisih, **implementasi NB mencakup RNW hampir seluruhnya.** Yang benar-benar
   tersisa setelah ini tinggal EDM.

⚠️ Ini **mengoreksi** rumusan R0 di arsip ("ketiga folder diekspor dari titik waktu berbeda"), yang
terlalu longgar. Koreksi lengkap di `_ARSIP-lintas-siklus\_BACA-INI.md`.

---

## 2. Enam temuan yang mengubah rancangan

### 2.1 ⛔ Tidak ada `RDB-Save` — seluruh tulisan Oracle lewat metode **baca**

`[terverifikasi]`

```powershell
foreach ($m in @('RDB-Save','RDB-List','RDB-Delete','RDB-Open','Obj-Save','Commit')) {
  $n = (Select-String -Path "D:\migrasi\RNM\NB FacIn\Activity\*.xml" `
        -Pattern "<pyStepsActivityName>$m</pyStepsActivityName>" -AllMatches | Measure-Object).Count
  "{0,-12} : {1}" -f $m, $n }
```

`RDB-Save` **0** · `RDB-List` 614 · `RDB-Delete` 2 · `RDB-Open` 0 · `Obj-Save` 45 · `Commit` **5**

Dan **55 dari 57** SQL yang menulis tersimpan di tag `pyBrowseSQL` — tag untuk *browse*.

**Arah operasi tidak boleh diturunkan dari nama metode atau nama tag.** Hanya verba SQL di dalam rule
`RDBList` yang membuktikannya.

### 2.2 ⛔ `COMMIT` milik database, bukan aplikasi

`[terverifikasi]` 36 SQL memuat `COMMIT`; hanya 5 langkah `Commit` di seluruh 609 Activity.
**21 SQL penulis tidak memuat `COMMIT` sama sekali** — termasuk jalur produksi. Lapisan repository Go
tidak dapat dirancang dengan asumsi transaksi yang biasa.

### 2.3 ⛔ Tidak ada rule `Property` — seluruh tipe data belum terverifikasi

`[terverifikasi]` 1.805 ekspresi properti unik (1.096 nama daun), nol deklarasi tipe. Satu-satunya
petunjuk adalah kontrol UI — dan **150 ekspresi diikat kontrol yang bertentangan**.

Diverifikasi langsung: `.ASMDateOfBirth` di `Section\CoverageSpreadingList_IsUW.xml` diikat sebagai
`pxDateTime` **dan** `pxInteger` — properti sama, berkas sama.

### 2.4 ⛔ Tidak ada aturan validasi selain wajib-isi

`[terverifikasi]` Nol aturan panjang/format/rentang/pola di lapisan Section; tipe rule
`Edit Validate`/`Validate` tidak ada di ekspor NB. Menambahkan validasi "yang masuk akal" akan
menolak data yang hari ini diterima — itu perubahan perilaku, bukan migrasi.

Ditambah **604 dropdown tanpa satu pun daftar nilainya** di korpus.

### 2.5 ⚠️ Spreading: tiga mesin, dan akibat pelanggaran ambang bukan eskalasi

`[terverifikasi]` NB memuat **tiga mesin spreading berbeda** dengan presisi berbeda (20 · 20 · 10/8
desimal). Empat sifat yang mematahkan asumsi lazim:

- **Spreading tidak mengubah `LetterNo` maupun `PositionNote`** — nol `Property-Set` di 74 Activity
  terkait. Dependensinya satu arah; `spreading` tidak memanggil `acceptance`.
- **Pelanggaran ambang mematikan tombol kirim, bukan menaikkan approver.** Kondisinya hanya muncul di
  `pyDisabledWhen` Section email — nol rujukan di `Flow`/`FlowAction`/`When`.
- **Scoring tidak menggerakkan alur sama sekali.** Pemakaian non-tampilan satu-satunya memeriksa
  *kelengkapan* (`FinalScore == ""`), bukan nilainya.
- **Validasi sisa tanpa toleransi** — menuntut `Σ Share == 100` dan
  `Σ PremiumSpreaded == PremiNusantaraRe` **persis**, tanpa distribusi galat pembulatan.

⛔ Mesin terbesar (`GetKapasitasTreaty`, 1,3 MB — Activity terbesar di NB) bergerbang parameter
`spread == "Spreading"`. Audit terverifikasi: parameter `spread` hanya dikirim **2** berkas
(`Harness\ViewPolis.xml`, `Section\InputInwardFacultativeDtl_IsUW.xml`), sedangkan `InSpreading`
dikirim 7 berkas. Apakah 1,3 MB logika itu perlu diport sama sekali — **memblokir**.

### 2.6 ✅ "Eksklusif NB" ternyata sebagian besar bukan New Business — **sudah diputuskan**

`[terverifikasi]` 139 rule hanya ada di folder NB. Tetapi rinciannya:

| Kelompok | Jumlah |
| --- | ---: |
| **Treaty Inward** (termasuk *seluruh* 22 activity perhitungan premi Treaty) | **80** |
| Produk Pega / SFA-CRM | 15 |
| Lini langsung (kendaraan, kargo, warranty) | 11 |
| Util | 5 |
| **Fac In New Business sejati** | **28** |

✅ **Keputusan work owner K-004 (15 September 2026): Treaty Inward TERMASUK lingkup.**
(Membatalkan K-001 yang semula mengeluarkannya.)

Karena itu 80 rule Treaty Inward di atas **tidak dikeluarkan**, dan batas yang tidak bersih antara
Treaty dan Fac In tidak perlu ditarik sama sekali — termasuk untuk `serviceInsertArasapas_act`
(panggilan servis konversi produksi yang berkelas Treaty tetapi dipakai Fac In).

⚠️ **Cakupan Treaty di sini terbatas, dan itu disadari.** `[terverifikasi]` Treaty Inward punya
folder korpus tersendiri berjumlah **1.149 berkas** yang **tidak dibaca** — keputusan K-005 menetapkan
batas kerja tetap di `NB FacIn\`. Yang didiscovery hanyalah **106 berkas** bernuansa Treaty yang ada
di dalam folder Fac In. Apa pun yang dokumen ini katakan tentang Treaty berlaku untuk **jejak Treaty
di dalam Fac In**, bukan untuk modul Treaty Inward yang utuh.

---

## 3. Angka pokok

| Ukuran | Nilai | Sumber |
| --- | ---: | --- |
| Berkas `.xml` NB | 2.083 | — |
| Activity | 609 | `01-activity\` |
| — langkah seluruhnya | 11.090 (4.069 tingkat atas + **7.021 bersarang**, kedalaman s/d 9) | `01-activity\` |
| — menulis Oracle | 62 activity / 225 langkah | `01-activity\` |
| — tak terjangkau (kode mati) | 5 + 1 kelompok klon 4-berkas | `01-activity\` |
| Section | 432 (152,3 MB; `ScoringRisk.xml` sendiri 7,68 MB) | `02-layar\` |
| — field terikat nyata | **6.314** (dari 7.311 sel, 997 placeholder dibuang) | `02-layar\` |
| — ekspresi properti unik / nama daun unik | 1.805 / 1.096 | `02-layar\` |
| — read-only | 4.834 = **76,6 %** | `02-layar\` |
| — RepeatGrid / PageList unik | 629 / 181 | `02-layar\` |
| ReportDefinition | 123 — **62 menggerakkan logika**, 60 tampilan, 1 yatim | `03-celah\02` |
| Rule eksklusif NB | 139 — hanya **28** Fac In NB sejati | `03-celah\03` |
| Berkas menyentuh spreading/capacity/scoring | 184 menurut nama, **381** menurut isi | `03-celah\01` |

⚠️ **Membaca hanya `pySteps` tingkat atas kehilangan 63 % logika.** Audit ulang wajib memakai parser
rekursif, bukan grep tingkat atas.

---

## 4. Jebakan "nama bukan bukti" — empat bukti baru dari pass ini

`CLAUDE.md` §3 butir 3 terbukti berulang kali, dan setiap kali dengan bentuk yang berbeda:

1. **`When/IsOfferFacIn`** — teks tampilannya berbunyi `Kode Bisnis = "02"/"58"/"SB"/"SG"`, ekspresi
   tersimpannya berbunyi `.Quotation.BusinessFac = "F"`. **Dua kondisi berbeda dalam satu rule**, dan
   rule ini menggerbangi flow masuk NB. ✅ Ditutup K-002: **ekspresi tersimpan yang berlaku**; teks
   tampilan dicatat sebagai kandidat perbaikan, tidak diimplementasikan.
2. **`SaveOfferProduction_Act` langkah 11** — metodenya `RDB-Delete`, deskripsinya *"insert to db"*.
3. **`Section/FormulaTreatyCapacityDesc`** — bernama "Treaty", berkelas `Data-OfferFacIn`, dipanggil
   dari layar Fac In.
4. **Label konektor di jalur spreading memuat nama orang**, sementara kondisi sebenarnya adalah
   `When IsTBonding`.

---

## 5. Satu blocker yang **gugur** di pass ini

`[terverifikasi]` `IsSpreadingDepan` dirujuk `NB FacIn\Activity\SpreadingAdditionalProtection.xml`
sebagai pra-kondisi langkah, dan berkasnya **tidak ada di `NB FacIn\When\`**. Ia sempat diusulkan
masuk daftar `panic()`.

Ternyata rule itu **ada di folder Endorsment**, dan kondisinya terbaca penuh:

```
pyLogic = A OR B OR C OR D
  A: Rule IsPA     evaluates to true
  B: Rule IsTravel evaluates to true
  C: Rule IsMBU    evaluates to true
  D: Rule IsFire   evaluates to true
```

⚠️ Salinan yang terbaca berversi `01-01-52` (2025-05-13). Karena salinan NB tidak ada, **tidak dapat
dipastikan** NB memakai versi yang sama. Implementasikan kondisinya, catat asalnya dari folder lain.

**Tidak perlu `panic()`** — ini rule yang hilang dari satu folder, bukan rule yang kondisinya tidak
diketahui. Sejalan dengan koreksi `CLAUDE.md` §4.5 di arsip: jumlah rule `When` yang kondisinya
benar-benar tidak terbaca tetap **0**.

---

## 6. Yang memblokir implementasi NB

Diurutkan menurut apa yang paling mahal bila ditebak.

> **Ditutup keputusan work owner 15 September 2026** (`00-KEPUTUSAN-WORK-OWNER.md`):
> `IsOfferFacIn` memakai `.Quotation.BusinessFac = "F"` (K-002) · `IsSpreadingDepan` tanpa `panic()`
> (K-003) · Treaty Inward **termasuk** lingkup (K-004, membatalkan K-001) · batas kerja **tetap**
> `NB FacIn\`, korpus Treaty tersendiri tidak dibaca (K-005).

**10 butir** — turun dari 15. Ditutup keputusan work owner: lingkup Treaty (K-004) dan `IsOfferFacIn`
(K-002). Dua lagi **ditutup dari korpus** pada sesi grilling 16 September 2026:

- ✅ **Arti `ProposalAcceptStatus`** — ditemukan dekoder literal di
  `NB FacIn\Activity\SaveViewSuggest.xml`: `1` Accept · `2` Reject · `3` Ask · `4` **Banding** ·
  `7` Decline · `9` Revise. `5`/`6` milik ranah ceding, `8` nol jejak. Sisa yang terbuka hanya
  **baris keputusan** `IsUWAccepted`, ditangani K-008.
- ✅ **Satuan `.Rate`** — bukan satu satuan: **per COB**. FIRE dan PA ‰; ANEKA, BONDING, GOLF,
  MARINE CARGO, MBU %. Risikonya berbalik — bahayanya menyeragamkan, bukan salah memilih. Ditangani
  K-010.

Rincian: `03-keputusan/RINGKASAN-GRILLING.md`.

⏸ **Satu butir berstatus ditangguhkan, bukan tertutup** (K-006): enam cabang alur yang memanggil
11 activity yang hilang dari ekspor — termasuk **Special Acceptance** dan **tangga akseptasi putaran
kedua**. Ekspornya tidak diminta secara khusus, tetapi cabangnya **tidak dibuang**: "hilang dari
ekspor" ≠ "usang", karena ekspor berasal dari titik waktu berbeda antar folder. Diverifikasi
terhadap ekspor produksi tunggal.

| # | Butir | Pemilik |
| ---: | --- | --- |
| 1 | Isi **35 stored procedure `POOLDATA.*`** — seluruh tulisan produksi melewatinya | DBA |
| 2 | **DDL / tipe kolom** — korpus nol DDL; 150 ekspresi berkontrol bertentangan | DBA |
| 3 | Siapa meng-commit jalur produksi? 21 SQL penulis tanpa `COMMIT`, nol langkah `Commit` | DBA + IT |
| 4 | Isi **604 dropdown** — tidak ada layar isian yang dapat dibangun lengkap | Product |
| 5 | Nama tabel Oracle untuk **64 kelas `Int-*`** — tanpa ini tak satu pun endpoint lookup dapat ditulis | DBA |
| 6 | Mesin spreading B: gerbang `spread`/`InSpreading` — menentukan apakah 1,3 MB logika perlu diport | IT + Underwriting |
| 7 | Alias koneksi `ASM` (529 langkah) vs `RNM` (87) menunjuk instance/skema apa | DBA |
| 8 | Arti `pyStepsPreCondition=false` (577 langkah) — tafsir berlawanan hasilnya | IT |
| 9 | 15 tautologi pembanding limit di `GetLimitAkseptasi_ActFlow` — disengaja? | Underwriting |
| 10 | Arti kode aksi prakondisi `4`/`5`/`6` — 340 percabangan NB tak dapat diport | IT |

⚠️ **Tenggat keras** kini ada pada butir 2 (`ALL_SOURCE` prosedur `POOLDATA`) dan pada ekstraksi
riwayat akseptasi — keduanya hanya mungkin selama Pega masih hidup. Daftar lengkap permintaan ke
IT/DBA: `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`.

---

## 7. Yang dapat dimulai **sekarang**

Tanpa menunggu jawaban bisnis:

1. `pkg/money` — `decimal.Decimal` + mata uang wajib. Dasarnya terverifikasi berlapis, termasuk
   temuan baru bahwa **112 Section menampilkan uang tanpa field mata uang mana pun**.
2. `internal/rules` — registry predikat `When`.
3. Kerangka agregat: satu inti + lima bagian opsional (963 nama daun khas lini vs 133 bersama
   menolak satu struct tunggal).
4. Komponen React read-only — 76,6 % field bersifat tampilan, bukan isian.

---

## 8. Peta dokumen pass NB

| Butuh | Baca |
| --- | --- |
| **Keputusan work owner yang menutup pertanyaan terbuka** | `00-KEPUTUSAN-WORK-OWNER.md` |
| Inventaris 609 Activity, graf panggilan, penulis Oracle | `01-activity\01-inventaris-activity-nb.md` |
| Skema field 432 Section + katalog 1.805 properti | `02-layar\01-skema-field-nb.md` |
| Spreading · capacity · scoring | `03-celah\01-spreading-capacity-scoring-nb.md` |
| 123 ReportDefinition | `03-celah\02-reportdefinition-nb.md` |
| 139 rule eksklusif NB | `03-celah\03-rule-eksklusif-nb.md` |
| **Spesifikasi modul Go** | `04-spec\01-modul-go.md` |
| **Spesifikasi model data** | `04-spec\02-model-data.md` |
| Temuan lintas-siklus (tetap berlaku) | `_ARSIP-lintas-siklus\_BACA-INI.md` |

---

*`CLAUDE.md` §1: perbaikan dipisahkan dari migrasi. Pass ini menemukan banyak yang tampak keliru —
tautologi pembanding limit, gerbang parameter yang tidak pernah cocok, properti yang diikat dua tipe
bertentangan, 18 cacat asal di jalur spreading. Semuanya dicatat sebagai kandidat perbaikan dan
direproduksi apa adanya. Memutuskannya milik bisnis.*
