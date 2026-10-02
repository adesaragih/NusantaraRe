# HANDOFF — Rangkuman sesi kerja, untuk dilanjutkan di sesi/asisten lain

**Dibuat:** 16 September 2026, 15:41 WIB · **Mencakup:** sesi 15–16 September 2026 (Prompt 1–11)
**Proyek:** migrasi **Facultative Inward** Nusantara Re dari **Pega** → **Go + React + Oracle**

> ## Cara memakai dokumen ini
>
> Dokumen ini **satu-satunya** yang perlu dibaca untuk memahami keadaan pekerjaan. Ia menggantikan
> riwayat percakapan, bukan meringkasnya sambil lalu: setiap keputusan, batas kerja, angka
> terverifikasi, kesalahan yang pernah terjadi, dan pekerjaan yang tersisa ada di sini.
>
> **Ia bukan sumber kebenaran.** Sumber kebenaran tetap berkas di `OUTPUT\` yang ditunjuk tiap bagian.
> Bila dokumen ini bertentangan dengan berkas yang ditunjuknya, **berkasnya yang benar** — dan
> dokumen ini yang basi.
>
> **Urutan baca yang disarankan:** §1 (batas kerja — mengikat) → §2 (peta berkas) → §4 (keputusan) →
> §7 (yang tersisa). Sisanya rujukan saat dibutuhkan.

---

## 1. Batas kerja — MENGIKAT, ditetapkan work owner

Keempat batas ini ditetapkan eksplisit di Prompt 1 dan **tidak pernah dicabut**. Melanggarnya
membatalkan hasil kerja, bukan sekadar menyalahi gaya.

| # | Batas | Bunyi |
| ---: | --- | --- |
| 1 | **Sumber bahan** | Bahan kerja/korpus **HANYA** dibaca dari `D:\migrasi\RNM`. Dilarang membaca bahan kerja dari luar folder ini. |
| 2 | **Korpus READ-ONLY** | `NB FacIn\`, `RNW Fac In\`, `Endorsment Fac In\` — jangan membuat, mengubah, atau menghapus apa pun di dalamnya. |
| 3 | **Tujuan output** | **SELURUH** output wajib ditulis ke `D:\migrasi\RNM\OUTPUT\`. Tidak ada pengecualian tanpa permintaan eksplisit work owner. |
| 4 | **Dilarang halusinasi** | Setiap klaim perilaku wajib menyebut **path berkas + nama rule + label** (`[terverifikasi]` / `[dugaan]` / `[pertanyaan terbuka]`). Setiap angka disertai perintah auditnya. |
| 5 | **Nama orang** | **Tidak pernah** disalin ke output mana pun. Guard berbasis identitas dicatat sebagai jumlah dan mekanismenya saja. |

### 1.1 GUARD SKILL — aturan yang paling sering terlupa

> Bila sebuah skill yang diminta **tidak dapat dijalankan otomatis**: **BERHENTI**, laporkan skill
> mana yang gagal, minta work owner menjalankannya manual. **JANGAN mengimprovisasi penggantinya.**

### 1.2 Satu-satunya penulisan di luar `OUTPUT\` yang pernah disahkan

`D:\migrasi\CLAUDE.md` **§4.5 saja**, atas permintaan eksplisit work owner (Prompt 5 dan satu
susulan). Work owner menegaskan: *"Jangan ubah bagian lain CLAUDE.md."* Bagian lain tidak disentuh.

### 1.3 Metodologi

Skill Matt Pocock yang tersedia di `.claude\skills\`: `research`, `grilling`, `domain-modeling`,
`to-questionnaire`. Ketiganya sudah dipakai; `domain-modeling` menghasilkan `steering/GLOSARIUM.md`
dalam format CONTEXT-FORMAT dan ADR dalam format ADR-FORMAT.

---

## 2. Peta berkas `D:\migrasi\RNM\OUTPUT\` — keadaan 16 September 2026

40 berkas, ~1,9 MB. Yang **wajib dibaca** ditandai ★.

```
OUTPUT\
├── _HANDOFF-SESI.md                     ← dokumen ini
├── _TOKEN-LOG.md                        4,3 KB   ledger token (lihat §9)
├── _EKSTRAKSI-PEGA-SELAGI-HIDUP.md ★    9,5 KB   5 permintaan ke IT/DBA — ADA TENGGAT KERAS
├── 00-KEPUTUSAN-WORK-OWNER.md ★        32,8 KB   K-001…K-018 + 3 konflik terbuka
├── 00-RINGKASAN-NB.md ★                13,7 KB   gambaran siklus NB + 10 blocker
├── 01-activity\01-inventaris-activity-nb.md      61,2 KB   609 Activity, graf panggilan
├── 02-layar\01-skema-field-nb.md                823,5 KB   432 Section, 1.805 properti
├── 03-celah\01-spreading-capacity-scoring-nb.md  99,8 KB
├── 03-celah\02-reportdefinition-nb.md            64,2 KB   123 ReportDefinition
├── 03-celah\03-rule-eksklusif-nb.md              43,4 KB   139 rule eksklusif NB
├── 03-keputusan\RINGKASAN-GRILLING.md             7,9 KB   3 ronde, 12 pertanyaan, K-007…K-013
├── 04-spec\01-modul-go.md                        13,5 KB
├── 04-spec\02-model-data.md                       7,5 KB
├── 04-kuesioner\
│   ├── 01-underwriting.md                         5,8 KB   4 pertanyaan — TERJAWAB (lihat §6)
│   ├── 02-product-aktuaria.md                     4,3 KB   2 pertanyaan — TERJAWAB
│   ├── 03-keamanan-it.md                          3,4 KB   1 pertanyaan — TERJAWAB
│   ├── 04-dba-pengukuran.md                      11,1 KB   D1/D2 — BELUM dijawab ⚠ struktur dikunci
│   └── _DITUNDA-fase-endorsement.md               3,7 KB   D-1/D-2/D-3 — ditunda, JANGAN dihapus
├── adr\ ★
│   ├── 0001-rekonsiliasi-eksak-bertahap.md        1,7 KB
│   ├── 0002-flag-fase-panic-vs-decline.md         2,1 KB
│   ├── 0003-mesin-akseptasi-fixture-sebagai-kontrak.md 1,6 KB
│   ├── 0004-money-dan-ratio-tipe-terpisah.md ★    4,1 KB   ← paling padat, memuat aturan pembagi
│   ├── 0005-presisi-pembulatan-literal-per-langkah.md  1,5 KB
│   └── 0006-mata-uang-unknown-eksplisit.md        1,8 KB
├── steering\
│   ├── GLOSARIUM.md ★                             5,9 KB   kosakata + istilah belum terverifikasi
│   └── PANDUAN-KERJA.md ★                         5,7 KB   4 jebakan korpus + kapan gagal keras
└── _ARSIP-lintas-siklus\                         15 berkas, ~660 KB
    └── _BACA-INI.md ★                             2,9 KB   apa yang tetap berlaku + koreksi R0
```

**Tentang arsip:** `_ARSIP-lintas-siklus\` adalah discovery Prompt 2 yang dibangun dari **ketiga**
folder korpus. Ia **tidak dibatalkan** — diarsipkan saat fokus dialihkan ke NB saja. Temuan
lintas-siklus di dalamnya tetap rujukan. Baca `_BACA-INI.md` untuk tahu bagian mana yang masih
berlaku dan mana yang sudah dikoreksi.

---

## 3. Kronologi sesi — 11 prompt

| Prompt | Isi | Keluaran utama |
| ---: | --- | --- |
| 0 | Buat ledger token | `_TOKEN-LOG.md` |
| 1 | Orientasi + penetapan batas kerja | batas §1 di atas |
| 2 | Discovery penuh (skill `research`), **bangun ulang dari korpus mentah** | 15 dokumen lintas-siklus (kini di `_ARSIP-`) |
| 3 | Pass fokus **NB FacIn** | 7 dokumen NB (~1,1 MB) |
| 4 | Keputusan lingkup | K-001…K-005 |
| 5 | Koreksi `CLAUDE.md` §4.5 | §4.5 ditulis ulang |
| 6 | `/grilling` 3 ronde + Tahap B | K-006…K-013, 6 ADR, 2 steering doc, `RINGKASAN-GRILLING.md` |
| 7 | `/to-questionnaire` — 4 kuesioner per penerima | `04-kuesioner\` |
| 8 | Filter kuesioner ke fokus NB | `_DITUNDA-fase-endorsement.md` |
| 9 | Kuesioner DBA: tabel + kolom eksplisit | lampiran kolom di `04-dba-pengukuran.md` |
| 10 | Jawaban kuesioner masuk | K-014…K-016 + Konflik K-1/K-2/K-3 |
| 11 | Kunci satuan rate | **K-018**, ADR-0004 dipertajam, GLOSARIUM diperbarui |

⚠️ **`K-017` tidak ada.** Penomoran melompat K-016 → K-018. Bukan berkas hilang — nomor itu memang
tidak pernah dipakai. Jangan mencarinya, dan **jangan memakai ulang** nomor itu untuk keputusan baru.

---

## 4. Keputusan work owner — K-001…K-018

Sumber lengkap: `00-KEPUTUSAN-WORK-OWNER.md`. **Keputusan di sana mengalahkan label
`[pertanyaan terbuka]` di dokumen mana pun.** Bila sebuah dokumen masih menyebut sesuatu terbuka
padahal sudah ada keputusan, dokumen itulah yang basi.

| # | Keputusan | Status |
| --- | --- | :-: |
| **K-001** | ~~Treaty Inward di luar lingkup~~ | ⛔ **DIBATALKAN** oleh K-004 — konsekuensinya **tidak pernah diterapkan** |
| **K-002** | `IsOfferFacIn` memakai `.Quotation.BusinessFac = "F"` (ekspresi tersimpan menang atas teks tampilan) | ✅ |
| **K-003** | `IsSpreadingDepan` **tidak** perlu `panic()` — ada di folder EDM, kondisinya terbaca. Predikatnya **tetap diimplementasikan** | ✅ |
| **K-004** | Treaty Inward **TERMASUK** lingkup | ✅ |
| **K-005** | Batas kerja **tetap** `NB FacIn\` — 1.149 berkas korpus Treaty di `RNM_BRD\` **tidak dibaca** | ✅ |
| **K-006** | 11 activity yang hilang dari ekspor tidak diminta khusus — **tetapi 6 cabang pemanggilnya DITANGGUHKAN, bukan dibuang** | ⏸ sebagian |
| **K-007** | Rekonsiliasi paralel run **eksak** (nol selisih sampai digit terakhir), dicapai bertahap | ✅ ADR-0001 |
| **K-008** | Satu flag fase: `panic` saat paralel run, `decline`+log saat produksi | ✅ ADR-0002 |
| **K-009** | Mesin akseptasi ditulis **sekarang**; fixture tabel limit = kontrak data ke DBA | ✅ ADR-0003 |
| **K-010** | `Money` dan `Ratio` tipe terpisah; skala melekat pada nilai | ✅ ADR-0004 |
| **K-011** | Presisi pembulatan **literal per-langkah**, bukan registry | ✅ ADR-0005 |
| **K-012** | Mata uang `Unknown` sebagai keadaan eksplisit | ✅ ADR-0006 |
| **K-013** | Urutan kerja ditetapkan; **siapa menulis kode Go/React ditunda** — keputusan tim | ⏸ |
| **K-014** | `Decline` = penolakan final tanpa hak banding · `Reject` = dapat dibanding. **Tidak disatukan** | ✅ |
| **K-015** | Renewal dinilai atas nilai pertanggungan **penuh** — disengaja | ✅ |
| **K-016** | Satuan rate **MBU = %** | ✅ |
| **K-018** | **Satuan rate per lini bisnis DIKUNCI** — lihat §4.1 | ✅ |
| **K-019** | `IsFacout` — fitur usang secara bisnis, **kode tetap diport apa adanya**. Status ⏸ ditangguhkan **dicabut** | ✅ |
| **K-020** | Alamat email tetap literal — **pengecualian eksplisit tercatat** terhadap `CLAUDE.md` §4.4 | ✅ |
| **K-021** | Domain `ProposalAcceptStatus` dikunci ke **{1,2,3,4,7,9}**; `5`/`6` tidak dipakai, `8` nol jejak. Validasinya **perilaku baru** → gagal keras + log, bukan tolak diam | ✅ |

### 4.1 K-018 — peta satuan rate, terkunci

| Lini bisnis | Satuan `.Rate` | Pembagi yang menempel pada `.Rate` |
| --- | :-: | ---: |
| **PA · Layering · FIRE** | **‰** | 1.000 |
| **MBU · ANEKA · BONDING · GOLF · MARINE CARGO** | **%** | 100 |

**Aturan penguraian pembagi komposit — mengikat.** Pembagi gabungan pada satu `@Math.divide` adalah
**hasil kali** sumbangan tiap faktor bersatuan, bukan satu satuan tunggal. Satuan `.Rate` dibaca dari
**faktor yang menempel padanya saja**, tidak pernah dari pembagi total.

| Faktor | Sumbangan ke pembagi |
| --- | ---: |
| `.Rate` ber-‰ | 1.000 |
| `.Rate` ber-% | 100 |
| `ProRatePercent` (persen) | 100 |

```
@Math.divide((.TSI * .Rate * pyWorkPage.OfferFacIn.ProRatePercent),100000,4)
        100000  =  1000 (.Rate ber-‰)  ×  100 (ProRatePercent ber-%)
```

Pembagi `100000` karena itu **menegaskan** `.Rate` = ‰ — bukan membantahnya.

⚠️ **Bukti bahwa aturan ini bukan tafsir:** `FillPremiMBU_FacIn.xml` L1144 berbentuk **bersarang** —
pembagi **dalam** (`100`) menempel pada `.Rate`, pembagi **luar** (`100`) menempel pada
`ProRatePercent`. Keterikatan faktor→pembagi tertulis eksplisit dalam struktur ekspresinya.

Bukti lengkap per baris (PA L713/L858/L1003/L1146 · MBU L1144/L1626 · Layering L909) ada di K-018 dan
ADR-0004, **sudah diverifikasi langsung ke korpus**, bukan dikutip dari ingatan.

---

## 5. Fakta terverifikasi yang mengubah rancangan

Semua berlabel `[terverifikasi]` dan punya perintah audit di dokumen sumbernya.

### 5.1 Struktur aplikasi

- **`StatusBusiness` adalah pembeda siklus**: `1` NB · `2` RNW · `3` EDM. Ketiganya berjalan di atas
  **satu basis rule**.
- **NB dan RNW tidak pernah berbeda** — 1.680/1.680 identitas rule, UI 575/575, nol selisih versi.
  **Implementasi NB mencakup RNW hampir seluruhnya.** Yang tersisa hanya EDM.
- **Tangga akseptasi adalah mesin keadaan satu langkah, bukan loop.** Satu keputusan manusia = satu
  transisi. Cabang `Else` = **penyelesaian normal**, bukan galat.
- **Tiga field state terpisah**, jangan disatukan jadi satu "status":
  `PositionNote` (antrean sekarang) · `LetterNo` → **kode jabatan tujuan** (bukan nomor surat!) ·
  `ProposalAcceptStatus` (hasil keputusan terakhir).
- **`ProposalAcceptStatus`** terdekode dari `SaveViewSuggest.xml` (identik di ketiga folder):
  `1` Accept · `2` Reject · `3` Ask · `4` **Banding** · `7` Decline · `9` Revise.
  `5`/`6` milik ranah ceding (properti terpisah); `8` nol jejak di 6.071 berkas.

### 5.2 Lapisan data — empat kejutan

| Temuan | Konsekuensi |
| --- | --- |
| **`RDB-Save` = 0.** Seluruh tulisan Oracle lewat `RDB-List` — metode **baca** — dan 55 dari 57 SQL penulis tersimpan di tag `pyBrowseSQL` | Arah operasi **tidak boleh** diturunkan dari nama metode atau nama tag. Hanya verba SQL yang membuktikan |
| **`COMMIT` milik database, bukan aplikasi.** 36 SQL memuat `COMMIT`; hanya 5 langkah `Commit` di 609 Activity; **21 SQL penulis tanpa `COMMIT` sama sekali** | Lapisan repository Go tidak dapat mengasumsikan transaksi biasa |
| **Nol rule `Property`** — 1.805 ekspresi properti, nol deklarasi tipe. **150 ekspresi diikat kontrol yang bertentangan** (`.ASMDateOfBirth` = `pxDateTime` **dan** `pxInteger`, berkas sama) | Seluruh tipe data **belum terverifikasi**. DDL dari DBA memblokir |
| **Nol aturan validasi selain wajib-isi** + **604 dropdown tanpa daftar nilai** | Menambahkan validasi "yang masuk akal" **menolak data yang hari ini diterima** — itu perubahan perilaku, bukan migrasi |

### 5.3 Spreading — empat sifat yang mematahkan asumsi lazim

- **Tiga mesin spreading berbeda** dengan presisi berbeda (20 · 20 · 10/8 desimal). Jangan disatukan.
- **Spreading tidak mengubah `LetterNo` maupun `PositionNote`** — nol `Property-Set` di 74 Activity
  terkait. `spreading` **tidak** memanggil `acceptance`.
- **Pelanggaran ambang kapasitas mematikan tombol kirim, bukan menaikkan approver.**
- **Scoring tidak menggerakkan alur sama sekali** — satu-satunya pemakaian non-tampilan memeriksa
  *kelengkapan* (`FinalScore == ""`), bukan nilainya.

### 5.4 Angka pokok NB

| Ukuran | Nilai |
| --- | ---: |
| Berkas `.xml` NB | 2.083 (korpus total 6.071) |
| Activity | 609 |
| — langkah seluruhnya | 11.090 (4.069 tingkat atas + **7.021 bersarang**, s/d kedalaman 9) |
| — menulis Oracle | 62 activity / 225 langkah |
| Section | 432 (152,3 MB) |
| — field terikat nyata | 6.314 · **76,6 % read-only** |
| — ekspresi properti unik / nama daun | 1.805 / 1.096 |
| ReportDefinition | 123 — **62 menggerakkan logika** |
| Rule `When` | 601 — **170 (28,3 %) menyembunyikan kondisi**, **0 tidak terbaca** |
| Rule eksklusif NB | 139 — hanya **28** Fac In NB sejati |

⚠️ **Membaca hanya `pySteps` tingkat atas kehilangan 63 % logika.**

---

## 6. Kuesioner — keadaan sekarang

Empat kuesioner ditulis dengan skill `to-questionnaire`, satu berkas per penerima.

| Berkas | Penerima | Keadaan |
| --- | --- | --- |
| `01-underwriting.md` | Underwriting | ✅ **dijawab lengkap & jawaban sudah ditulis balik ke berkasnya** → K-014, K-015, K-019, K-021 |
| `02-product-aktuaria.md` | Product/Aktuaria | 🟡 **dijawab & ditulis balik** — butir 2 (label) ✅ tertutup → K-016, K-018; **butir 1 (satuan di slip) masih menunggu D2** |
| `03-keamanan-it.md` | Keamanan/IT | ✅ **dijawab lengkap & ditulis balik** → K-020 (lewat Konflik K-3) |
| `04-dba-pengukuran.md` | DBA | **belum dijawab** — D1 dan D2 masih terbuka |
| `_DITUNDA-fase-endorsement.md` | Underwriting | **ditunda ke fase endorsement** — D-1/D-2/D-3 |

### ⚠️ Dua hal yang wajib diketahui penerus

1. ✅ **Jawaban sudah ditulis balik ke ketiga kuesioner** (16 September 2026) — `01-underwriting.md`,
   `02-product-aktuaria.md`, dan `03-keamanan-it.md` semuanya memuat jawaban + keputusan yang
   dihasilkannya, berikut banner status di kepala berkas. Ketiganya kini **dapat dibaca berdiri
   sendiri** tanpa register keputusan. Naskah pertanyaannya tidak diubah.
   ⚠️ **Yang tersisa:** stub *"Anything else?"* di ketiganya memang kosong (tidak ada jawaban yang
   diterima), dan **butir 1 `02-product-aktuaria.md` masih menunggu D2**.
   ⚠️ **Penomoran `02-product-aktuaria.md` terbalik dari `RINGKASAN-GRILLING.md`**: berkas kuesioner
   memakai *butir 1 = slip · butir 2 = label*, ringkasan grilling memakai *P1 = label · P2 = slip*.
   Isinya sama, hanya urutannya berbeda — sudah dicatat di banner berkasnya.
2. **Struktur `04-dba-pengukuran.md` dikunci.** Work owner: *"jangan rubah struktur tabel yang sudah
   ada, jika sebelumnya mata uang dipisah, jangan dirubah."* Berkas ini **tidak boleh** dirapikan,
   digabung kolomnya, atau diformat ulang.
3. **`_DITUNDA-fase-endorsement.md` jangan dihapus.** Work owner: *"jangan dihapus, hanya ditunda."*

---

## 7. Yang masih terbuka — pekerjaan tersisa

### 7.1 Tiga konflik — **semuanya sudah tertutup** (16 September 2026)

| Konflik | Isi | Penyelesaian |
| --- | --- | --- |
| **K-1** | ~~Satuan rate PA & Layering~~ | ✅ **DITUTUP oleh K-018** — rumus korpus menang, PA & Layering = ‰ |
| **K-2** | ~~*"`IsFacout` tidak dipakai lagi"* vs 10 rule perujuk~~ | ✅ **DITUTUP oleh K-019** — fitur usang secara bisnis, **kode tetap diport apa adanya**; status ⏸ ditangguhkan dicabut, jadi porting normal |
| **K-3** | ~~*"alamat email tidak usah dipindahkan"* vs `CLAUDE.md` §4.4~~ | ✅ **DITUTUP oleh K-020** — diterima sebagai **pengecualian eksplisit tercatat**; `CLAUDE.md` **tidak** diubah, pengecualian hidup di register keputusan |

**Tidak ada lagi konflik terbuka.** Yang tersisa di bawah adalah blocker data (butuh pihak luar), bukan
pertanyaan keputusan.

### 7.2 Sepuluh blocker implementasi NB

Sumber: `00-RINGKASAN-NB.md` §6. Diurutkan menurut apa yang paling mahal bila ditebak.

| # | Butir | Pemilik |
| ---: | --- | --- |
| 1 | Isi **35 stored procedure `POOLDATA.*`** — seluruh tulisan produksi melewatinya | DBA |
| 2 | **DDL / tipe kolom** — korpus nol DDL; 150 ekspresi berkontrol bertentangan | DBA |
| 3 | Siapa meng-commit jalur produksi? 21 SQL penulis tanpa `COMMIT` | DBA + IT |
| 4 | Isi **604 dropdown** | Product |
| 5 | Nama tabel Oracle untuk **64 kelas `Int-*`** | DBA |
| 6 | Mesin spreading B: gerbang `spread`/`InSpreading` — menentukan apakah 1,3 MB logika perlu diport | IT + Underwriting |
| 7 | Alias koneksi `ASM` (529 langkah) vs `RNM` (87) menunjuk instance/skema apa | DBA |
| 8 | Arti `pyStepsPreCondition=false` (577 langkah) | IT |
| 9 | 15 tautologi pembanding limit di `GetLimitAkseptasi_ActFlow` — disengaja? | Underwriting |
| 10 | Arti kode aksi prakondisi `4`/`5`/`6` — 340 percabangan NB | IT |

### 7.3 Lima permintaan ke IT/DBA — **ada tenggat keras**

Sumber: `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`. **Hanya dapat diambil selama Pega masih hidup.**

| # | Butir | Pemilik |
| ---: | --- | --- |
| 1 | Ekspor produksi **tunggal**, satu titik waktu, tiga siklus — juga alat verifikasi 6 cabang K-006 | IT |
| 2 | `IsUWAccepted` + 11 DecisionTable lain **beserta barisnya** (baris pemetaan tidak ikut terekspor) | IT |
| 3 | `ALL_SOURCE` 35 prosedur `POOLDATA` — definisi kode saja, tanpa data | DBA |
| 4 | Isi + DDL 7 tabel `M_LIMIT_*`, terutama **ejaan pasti nilai kolom `JABATAN`** | DBA |
| 5 | Ekstraksi `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` + DDL `HISTORYAKSEPTASIPEGA` | DBA + IT |

⛔ **Butir 5 adalah satu-satunya yang datanya benar-benar musnah bersama Pega.** Bila hanya satu yang
dapat dikerjakan hari ini, kerjakan butir 5.

### 7.4 Pekerjaan yang dapat dimulai **sekarang**, tanpa menunggu siapa pun

1. `pkg/money` — `Money{Amount decimal, Currency}` + `Ratio{Value decimal, Scale}` (ADR-0004, 0006).
2. `internal/rules` — registry predikat `When` + **resolver COB → skala** memakai peta K-018 terkunci.
3. `services/acceptance` — mesin tangga dengan fixture tabel limit (K-009); fixture-nya **dikirim ke
   DBA sebagai kontrak bentuk data**.
4. Kerangka agregat: satu inti + lima bagian opsional (963 nama daun khas lini vs 133 bersama menolak
   satu struct tunggal).
5. Komponen React read-only — 76,6 % field bersifat tampilan.

### 7.5 Belum diminta work owner — jangan dikerjakan tanpa permintaan

Gelombang kedua kuesioner sempat teridentifikasi tetapi **tidak pernah diminta**: kamus enum Product
(98 `BusinessCode` + 87 `BusinessOldId` + 604 dropdown), semantik konstruksi Pega ke IT, infrastruktur
ke DBA.

---

## 8. Jebakan yang sudah terbukti — jangan diulang

### 8.1 Empat jebakan membaca korpus Pega

Keempatnya **benar-benar terjadi** di sesi ini dan menghasilkan angka salah sebelum dikoreksi.
Sumber: `steering/PANDUAN-KERJA.md` §2.

| # | Jebakan | Koreksi |
| --- | --- | --- |
| 1 | Membaca kondisi `When` hanya dari `<pyConditionString>` | Baca **kedua** tag — `<pyConditionString>` **dan** `<pyConditionValue1String>`. 170 berkas menyembunyikan kondisinya di tag kedua |
| 2 | `Get-FileHash` mentah untuk membandingkan antar folder | Buang baris `px*`/`pz*` — metadata ekspor berbeda di **setiap** berkas |
| 3 | Membandingkan tanpa mengurutkan, atau dengan `Sort-Object` bawaan | Urutkan dengan **`[StringComparer]::Ordinal`** — urutan tag ekspor tidak deterministik, dan ada nama rule yang hanya beda kapitalisasi |
| 4 | Membaca hanya `pySteps` tingkat atas | Parser **rekursif** — 63 % logika ada di langkah bersarang, sampai kedalaman 9 |

Tambahan: nama section pada FlowAction ada di `pySectionReference`, **bukan** `pyStreamName`.

### 8.2 "Nama bukan bukti" — enam bukti dari korpus ini

1. `When/IsFacout` namanya *fac out*, isinya menguji **banding**.
2. `When/IsOfferFacIn` memuat **dua kondisi berbeda** — teks tampilan vs ekspresi tersimpan.
3. Dua rule bernama `Insert…` sesungguhnya **menghapus lalu menyisipkan ulang**.
4. Dua shape flow berlabel sama *"Err Konversi?"*; yang satu menggerbangi hal lain sama sekali.
5. `Section/FormulaTreatyCapacityDesc` bernama *Treaty*, berkelas Fac In, dipanggil dari layar Fac In.
6. Seluruh tulisan Oracle memakai `RDB-List` — **metode baca**.

### 8.3 Kapan gagal keras (`panic`)

| Situasi | Sikap |
| --- | --- |
| Arti enumerasi belum dijawab bisnis | `panic` |
| Isi stored procedure / tabel lookup tidak diketahui | `panic` |
| Rule dirujuk tetapi berkasnya tidak ada **di folder mana pun** | `panic` bila tercapai — **cabangnya tidak dihapus** |
| COB tidak ada di peta resolver skala | `panic` |
| Mata uang `Unknown` dipakai lintas mata uang / ditulis ke Oracle | `panic` |
| Status dikenali, baris keputusannya belum terverifikasi | `panic` saat paralel run · `decline`+log saat produksi |
| Kondisi tidak dikenali sama sekali | `panic` di **kedua** fase |

⚠️ **`IsPKSASM` gagal keras** — kondisinya terbaca (`IsB2B = "ASM"`) tetapi label **belum ter-resolve**
dan `<pyTempText>true`. Ia gagal keras bukan karena kondisinya hilang, melainkan karena belum terbukti
sebagai yang dieksekusi.

---

## 9. Kesalahan yang terjadi di sesi ini — dan koreksinya

Bagian ini sengaja ada. Penerus yang tidak tahu kesalahan ini akan mengulanginya.

### 9.1 Koreksi dari work owner — yang paling penting

> ⛔ **"Hilang dari ekspor" ≠ "usang".**
>
> Enam cabang alur yang memanggil 11 activity yang hilang dari ekspor sempat **disimpulkan sebagai
> kode mati**. Work owner mencabut kesimpulan itu: ekspor korpus diambil dari **titik waktu berbeda
> antar folder** (95 rule berbeda versi, arah tidak konsisten), sehingga ketiadaan sebuah rule
> **bukan bukti** ia tidak dipakai. Keenamnya kini **⏸ DITANGGUHKAN**, menunggu verifikasi terhadap
> ekspor produksi tunggal. Di antaranya **jalur Special Acceptance** dan **tangga akseptasi putaran
> kedua** — keduanya inti aplikasi.
>
> Bukti R0 sudah ada di tangan saat kesimpulan itu dibuat. Ini kesalahan penalaran, bukan kekurangan
> data. **Pelajaran umum: jangan menyimpulkan ketiadaan sebagai keusangan.**

### 9.2 Koreksi yang diajukan dan diterima work owner

> Usulan mengganti default `decline` dengan `panic` di semua fase **ditarik**.
> `<pyDefaultResult>decline</pyDefaultResult>` terbaca langsung dari `IsUWAccepted` — itu **perilaku
> terekam**, bukan celah pengetahuan. Menggantinya melanggar `CLAUDE.md` §1. Diselesaikan lewat
> pemisahan fase (K-008).

### 9.3 Kesalahan pengukuran yang terdeteksi sendiri dan diperbaiki

| Klaim awal | Kenyataan | Sebab |
| --- | --- | --- |
| "0 dari 1.680 berkas identik" | 1.220 identik | `Get-FileHash` mentah → metadata ekspor; lalu urutan tag non-deterministik |
| "507 identik" | 1.220 | `Sort-Object` bawaan tidak stabil → `[StringComparer]::Ordinal` |
| "132 activity hilang" | **11** terverifikasi | heuristik pertama ikut menangkap DataTransform dan nama properti (`Country`, `CheckBox1`) |
| "19 activity hilang" | **11** | over-claim, dikoreksi |
| "24 prosedur `POOLDATA`" | **35** | terukur ulang |
| "11 butir blocker" | **10** | tabelnya memang 10 |
| R0 "ketiga folder diekspor dari titik waktu berbeda" | 95 rule berbeda **versi ruleset**, arah tidak konsisten, **NB↔RNW = 0** | rumusan lama terlalu longgar dan salah sasaran |

### 9.4 Pemindai kebocoran nama orang — positif palsu

Pemindaian substring sempat menandai empat nilai — di antaranya `System`, `Manager`, `Administrator`.
Diperbaiki dengan pencocokan **batas kata**, lalu keempatnya diverifikasi sebagai bawaan Pega / nama
jabatan — bukan nama orang. **Hasil akhir: nol nama orang bocor ke output mana pun.**

---

## 10. Ledger token

`_TOKEN-LOG.md` mencatat setiap tahap. **Seluruh angka adalah ESTIMASI** yang dihitung asisten dari
selisih penghitung anggaran token sesi — **bukan** pembacaan saldo/kuota akun.

Kolom: `Tahap | Waktu START | Waktu END | Durasi | Est. Token Input | Est. Token Output | Est. Total |
Est. Total Kumulatif`. Zona waktu **WIB**. Satu baris = satu tahap; tahap yang diulang jadi baris
**terpisah**, tidak menimpa baris lama.

Kumulatif s/d Prompt 11: **~3.329.500 [ESTIMASI]**.

---

## 11. Konvensi menulis — ikuti persis

- Seluruh dokumen berbahasa **Indonesia**, kecuali istilah teknis Pega/Go yang memang tidak
  diterjemahkan.
- Setiap klaim perilaku: **path berkas + nama rule + label** `[terverifikasi]` / `[dugaan]` /
  `[pertanyaan terbuka]`.
- Setiap angka: **perintah audit** yang menghasilkannya, ditulis apa adanya sebagai blok PowerShell.
- **Keputusan yang dibalik tidak dihapus** — ditandai dibatalkan dan menunjuk penggantinya, supaya
  alasan perubahannya tetap terbaca. Lihat K-001 dan Konflik K-1 sebagai contoh bentuknya.
- Menambah keputusan: entri `K-0nn` baru dengan **tanggal, pertanyaan asal, dan konsekuensinya**.
- Kosakata: pakai `steering/GLOSARIUM.md`, termasuk daftar _Avoid_-nya. Istilah yang artinya belum
  diketahui **tetap dicantumkan** dengan penanda ☐ — menghapusnya membuat kekurangan itu tak terlihat.

---

## 12. Lima kalimat yang paling mudah salah — ringkasan terakhir

1. **Ini migrasi, bukan greenfield.** Perilaku direkam dari korpus, bukan direka ulang. Perbaikan
   dipisahkan dari migrasi, dan memutuskannya milik bisnis.
2. **`LetterNo` tidak berisi nomor surat** — isinya kode jabatan tujuan berikutnya.
3. **`.Rate` tidak punya satu satuan** — per COB, dan petanya sudah dikunci K-018.
4. **Ketiadaan bukan keusangan.** Rule yang hilang dari ekspor ditangguhkan, bukan dibuang.
5. **Aplikasi yang berjalan dan salah diam-diam jauh lebih berbahaya daripada aplikasi yang berhenti.**

---

*Dokumen ini tidak memuat nama orang, kredensial, endpoint, maupun data pelanggan.*
