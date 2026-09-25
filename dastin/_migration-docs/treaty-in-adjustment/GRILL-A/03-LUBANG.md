> Modul  : Treaty In Adjustment · Ronde A · 2026-09-24
> Peran  : interogator
> Masukan: seluruh temuan ronde A · `METODE` §4.6, §6.2, §7.2
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 03 · LUBANG

## 1. Yang ditutup ronde ini dengan MEMBACA — bukan dengan bertanya

| Lubang | Ditutup oleh |
|---|---|
| Untuk revisi kedua, `OLDID` menunjuk apa? (`METODE` §8.6 Q2) | `TreatyInRevisi_post` langkah 5 |
| Apakah selisih historis dapat direproduksi? | `TreatyInSetAddendumToHistory` + ruas kanan seluruh `TreatyEDMDifference*` |
| Apakah jenis addendum dapat diubah pengguna? | dua radio `pyReadOnly=false` di picker Revisi |
| Apakah membuka kontrak mengubahnya? | precondition `param.revisionstate == 1` + parameter kedelapan kontrolnya |
| Berapa akar `ValueDifference`, dan apa saja? | sapuan + daftar di §5.6 |

## 2. Pertanyaan `QA-x` ronde ini

| # | Pertanyaan | Putusan |
|---|---|---|
| **QA-1** | Hubungan dengan migrasi modul induk | GRL-01 |
| **QA-2** | Rekonsiliasi ADR-0036 | GRL-02 |
| **QA-3** | Rekonsiliasi ADR-0048 | GRL-03 |
| **QA-4** | Apakah penyesuaian entitas tersendiri (`METODE` §8.6 Q1) | GRL-04 |
| **QA-5** | Rekonsiliasi ADR-0040 | GRL-05 |
| **QA-6** | Rekonsiliasi ADR-0049 | GRL-06 — **ditutup dengan rujukan**; anggaran berkurang satu |
| **QA-7** | Rekonsiliasi ADR-0052 | GRL-07 |
| **QA-8** | Rekonsiliasi ADR-0055 | GRL-08 — **cabang A ditutup** |

## 3. Pengambilan bukti yang belum dijalankan

### 3.1 Ekspor tambahan

| # | Ekspor | Memblokir |
|---|---|---|
| **EXP-1** | Aturan **Field Value** `EDMState` dan `EDMMaterialType` | menutup TDA-13; menuntaskan cabang C |

### 3.2 Kueri data — tidak satu pun memblokir rancangan (`METODE` §7.2)

`UA-1` bentuk pengenal yang benar-benar terjadi (tiga kueri) · `UA-2` sebaran jenis x materialitas ·
`UA-3` non-material yang nilainya berubah · `UA-4` `RevisionState` tertinggal · `UA-5` daftar yang
berubah panjang · `UA-6` mata uang yang berubah · `UA-7` detail yatim · `UA-8` ukuran `JSONDATA` ·
**`UA-9`** addendum ber-`RevisionState = 1` warisan · **`UA-10`** addendum yang melanggar lapisan
beku, **kelima** field · **`UA-11`** kontrak yang kehilangan status akseptasi · **`UA-12`** jejak peran kelompok, **termasuk sejarahnya** · **`UA-13`** berapa revisi-di-tempat per bulan, untuk memperkirakan tambahan beban penyetuju.

**Tenggat** ditetapkan pada bagian serah terima, bukan di sini.

**Dua butir yang DICORET dari cabang D, karena sudah dijawab ADR-0055:**

| Butir | Dijawab oleh |
|---|---|
| **TDA-02** — nasib penolakan yang menghapus baris | ADR-0055 perubahan 24 Sep, **syarat (a)**: baris tidak dihapus; ditambah `DITOLAK` terminal dan keadaan `DIBATALKAN` |
| **TDA-03** — nasib sampah yang ditinggalkannya | mekanismenya hilang bersama (a); **sisa warisannya** pindah ke cabang I, diukur `UA-7` |

Anggaran karena itu turun **42 -> 40**.

### 3.3 Daftar untuk dibantah ke bisnis

`DB-1` selisih sebagai angka persetujuan · `DB-2` rujukan lama menunjuk versi · `DB-3` s.d. `DB-5`
dari GRL-04 · `DB-6`, `DB-7` dari GRL-05 · `DB-8`, `DB-9` dari NA-08 · **`DB-10`** dari GRL-07,
tentang tingkat persetujuan revisi. Bunyinya di `07-AUDIT`.

### 3.4 Pertanyaan untuk pemilik proses

`PP-1` pemilik pengerjaan to-spec induk — tidak memblokir; hanya menentukan besar ongkos jadwal
pada arah dampak GRL-01.

## 4. Lubang yang dititipkan ke cabang lain

| Cabang | Titipan |
|---|---|
| **C** | definisi materialitas (**C1**, gerbang bagi E); daftar jenis sesudah EXP-1; arti bisnis `ActualValue` |
| **D** | rencana penyalaan "versi menggantikan" sebagai kemampuan baru (`METODE` §3.8) |
| **D** | **pemberitahuan kepada pembuat addendum** yang barisnya terkunci warisan lalu terbuka sejak cut-over (GRL-08 butir iii); besarannya `UA-9(a)` |
| **I** | **nasib addendum yang masih `DRAFT` saat peralihan** — termasuk yang terkunci warisan (GRL-08 butir iii) |
| **I** | **sampah warisan TDA-03**: detail addendum yatim, lampiran salinan, dan `RevisionState = 1` yang tertinggal di baris kontrak sesudah addendumnya dihapus. Mekanismenya hilang di model baru (ADR-0055 syarat a), **sisanya tidak**. Diukur `UA-7` |
| **D** | **akibat GRL-07 ke orang**: revisi yang hari ini berhenti di kepala seksi akan naik sampai direktur sejak cut-over. Siapa yang harus diberi tahu, dan berapa tambahan beban penyetuju — diukur `UA-13` |
| **B** | **daftar pilihan revisi tidak menyaring keadaan sama sekali** (NA-13, TDA-11): kontrak yang sedang direvisi-di-tempat tetap dapat dipilih, dan addendumnya lahir terkunci dengan rantai pendek. Bolehkah addendum dibuat dari kontrak yang sedang punya versi `DRAFT`? |
| **D** atau eskalasi | **tiga tombol "(dev)"** (NA-16): dijaga `pyOrgDivision = 'IT'` dan nol privilese; di dalam divisi itu `Force Edit (dev)` terlihat semua orang karena `pyVisible = ALWAYS`. Keberadaannya **sudah** ada di butir 1 daftar eskalasi induk dan ADR-0055 §4, jadi **tidak ada butir eskalasi baru** — yang perlu hanya koreksi `REV-3`. Sisanya pekerjaan pengamanan sistem lama selama masa hidup berdampingan |
| **E** | mekanisme pembekuan — tinggal **dua** pilihan; `EDMEffective` -> `TreatyCalculateProratePct` dan nasib pro rata |
| **I** | selisih mana yang benar untuk addendum historis; versi mana yang berlaku untuk kontrak warisan; **nasib addendum warisan yang melanggar lapisan beku** — besarnya diukur `UA-10` |
| **F** | model baru menyimpan **hanya rujukan** ke cedant; namanya diturunkan dari rujukan itu. **PERUBAHAN**, berubah dari: nama dan pengenal disimpan terpisah tanpa pemeriksaan (NA-08) |
| **I** | bila nama dan pengenal cedant tidak sinkron pada baris warisan, mana yang **kebenaran**? Dan pertanyaan yang lebih besar: **aturan baru mana yang boleh dijalankan atas baris warisan**, mana yang tidak, dan hasilnya dicatat dengan asal-usul bagaimana (GRL-06) |
| **C** | bingkai pemasangan jenis: **dua sumbu × hingga enam kombinasi** versus **satu sumbu × tiga jenis**; sumbernya `EXP-1` (yang ditawarkan) dan `UA-2` (yang pernah terbentuk) |
| **B** | sistem lama tampaknya punya **dua** mekanisme perubahan: revisi-di-tempat (Treaty In, lewat kontrol 7-8) dan addendum (Adjustment). Apakah revisi-di-tempat juga menjadi `VERSI_KONTRAK`? **ADR-0055 sudah memutuskannya** — perubahan atas versi yang disetujui melahirkan versi baru dari `DRAFT` — jadi yang tersisa hanya mengkonfirmasi, bukan memutuskan |
| **D** | tahap peringatan sebelum memblokir lapisan beku, beserta **angka pemicu** peralihan ke blokir (`METODE` §3.8) |
