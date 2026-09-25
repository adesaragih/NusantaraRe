# Temuan Adjustment yang TIDAK diadili

**Tanggal:** 23 September 2026
**Sifat berkas ini:** **daftar, bukan putusan.**

Embargo modul Treaty In Adjustment dibuka **sebagian** — hanya untuk mengetahui entitasnya,
sambungannya ke Treaty In, dan bentuk tabel selisihnya. Temuan tentang **perilakunya** dijumpai
sambil jalan, dan **tidak satu pun diadili di sini**: tidak dinilai cacatnya, tidak ditetapkan
invariannya, tidak dirancang perbaikannya.

Menuliskannya bukan mengadilinya. Setiap butir membawa **bukti berkasnya**, supaya sesi Adjustment
tidak perlu menemukannya lagi dari nol.

---

## Pemilahan menurut garis 323 / 56 — ditambahkan 23 September 2026

Aturan penundaan semula ditulis dengan anggapan Adjustment adalah permukaan tersendiri. Temuan
323 berkas byte-identik membuktikan anggapan itu salah untuk 85 % permukaannya, jadi aturannya
terbelah:

| Berkasnya ada di | Perlakuan |
|---|---|
| **irisan 323** | **bukan** temuan Adjustment — ia temuan **Treaty In**, dijumpai dari pintu masuk lain. Ia berjalan pada jalur Treaty In hari ini, pada data Treaty In hari ini. Memarkirnya berarti menunda cacat yang sedang aktif di modul yang sedang dispesifikasikan. |
| **56 khas Adjustment** | tetap diparkir; aturan asli berlaku |

Hasil pemilahan:

| Butir | Berkasnya | Sisi | Nasibnya |
|---|---|---|---|
| A-1 nomor revisi dari potongan teks | `GetTreatyRevisionID.xml` | **56 khas** | tetap parkir |
| A-2 `ROWNUM = 1` mendahului `ORDER BY` | `GetTreatyRevisionID.xml` | **56 khas** | tetap parkir |
| B-1 *picker* meng-UNION tanpa pembeda | `TreatyLoadMasterJoinEdm.xml` | **56 khas** | tetap parkir |
| B-2 lapisan beku disalin ke tiap baris addendum | `SaveTreatyInEDM.xml` | **IRISAN 323** | **pindah** → temuan Treaty In |
| C-1 penyimpan detail kontrak dipanggil | `SaveTreatyIn_EDM_Act.xml` | **IRISAN 323** | sudah dipindah (§G1b) |
| C-2 tidak ada tabel selisih | sapuan, bukan satu berkas | — | sudah dipakai di G3 |
| D-1 `PositionUsername` di kelas addendum | kelas khas Adjustment | **khas** | tetap parkir |
| D-2 tombol paksa yang sama | `TreatyInForceEdit.xml` dll. | **IRISAN 323** | **sudah** jadi butir 1 daftar eskalasi Treaty In |
| G4 pengenal dibentuk dari teks, offset tetap | `TreatyInRevisi_post.xml` | **56 khas** | tetap parkir; **tetapi** akibatnya untuk migrasi diukur **Uji AA** |

**Empat tetap parkir, tiga pindah, dua bukan berkas tunggal.**

> A-2 adalah kandidat terkuat untuk pindah — `ROWNUM = 1` yang mendahului `ORDER BY` menghasilkan
> baris sembarang, sehingga nomor revisi dapat berulang atau melompat tanpa galat. Ia **tidak**
> pindah karena berkasnya benar-benar khas Adjustment: `GetTreatyRevisionID` tidak ada di ekspor
> Treaty In sama sekali. Akibatnya pada penelusuran rantai tetap diukur dari sisi Treaty In lewat
> **Uji AA-3**, yang menghitung bentuk akhiran pengenal addendum — termasuk yang berulang.

---

## Cara membaca daftar ini

| Kolom | Isinya |
|---|---|
| **Yang terlihat** | apa yang terbaca dari ekspor, apa adanya |
| **Bukti** | berkas dan potongannya |
| **Kenapa ditunda** | alasan ia tidak diputuskan sekarang |

Tidak ada kolom "akibat" dan tidak ada kolom "usulan". Keduanya milik sesi Adjustment.

---

## A — Identitas dan penomoran revisi

### A-1. Nomor revisi diambil dari potongan teks ID

**Yang terlihat.** Nomor revisi sebuah addendum dibaca dari karakter ke-10 dan ke-11 pengenalnya.

**Bukti** — `RDBList/GetTreatyRevisionID.xml`:

```sql
select TO_NUMBER(SUBSTR(ID,10,2))+1 as HASIL1
from m_treaty_in_edm
where ID like {InputData.CARI1}
and ROWNUM = 1
order by ID desc
```

**Kenapa ditunda.** Identitas yang terbentuk dari penguraian teks adalah pola yang sudah dilarang di
model baru, tetapi **berapa besar akibatnya pada data yang ada** adalah pertanyaan Adjustment.

### A-2. `ROWNUM = 1` mendahului `ORDER BY` pada kueri yang sama

**Yang terlihat.** Pada kueri di A-1, Oracle menerapkan `ROWNUM` **sebelum** `ORDER BY`. Baris yang
kembali karena itu **bukan** baris dengan `ID` tertinggi.

**Bukti.** Kueri yang sama, `RDBList/GetTreatyRevisionID.xml`.

**Kenapa ditunda.** Apa yang terjadi ketika nomor revisi berikutnya dihitung dari baris yang bukan
terakhir adalah perilaku Adjustment. Ini **kandidat kuat untuk uji data di sesinya** — bukan
sekarang.

---

## B — Sambungan dan pemilihan

### B-1. *Picker* menyatukan kontrak dan addendum tanpa pembeda

**Yang terlihat.** Daftar "yang disesuaikan" adalah UNION dua tabel, dan baris keduanya disajikan
dengan kolom yang sama tanpa penanda jenisnya.

**Bukti** — `RDBList/TreatyLoadMasterJoinEdm.xml`:

```sql
select a.ID, a.TREATYCONTRACTNAME, … from pooldata.treaty_in a
UNION
select b.ID, b.TREATYCONTRACTNAME, … from pooldata.treaty_in_edm b
```

Ada varian XOL-nya: `RDBList/TreatyLoadMasterJoinEdmXOL.xml`.

**Kenapa ditunda.** Ini menjelaskan mengapa `OLDID` dapat menunjuk dua jenis baris — dan itu sudah
dipakai di gerbang G2. Tetapi **apakah pengguna dapat membedakannya di layar**, dan apa yang terjadi
bila ia salah pilih, adalah perilaku.

### B-2. Lapisan beku disalin ke setiap baris addendum

**Yang terlihat.** Setiap baris `M_TREATY_IN_EDM` menyimpan sendiri `ProportionType`,
`Commencement`, `Termination`, `Ceding`, `CedingID`, `LeadingReinsSource`, `LeadingReinsSourceID` —
kelima hal yang ADR-0040 §3 tetapkan tidak boleh berubah antar versi.

**Bukti** — `RDBList/SaveTreatyInEDM.xml`, daftar parameter `PEGA_M_TREATY_IN_EDM(…)`.

**Kenapa ditunda.** Di model baru kelimanya hanya ada di `KONTRAK`, sehingga bentuk ini tidak
terbawa dan invariannya **tidak dapat dilanggar** (`SPEC-INVARIAN.md` INV-B1). **Apakah data yang
ada sudah menyimpang** adalah pertanyaan untuk sesi Adjustment, dan jawabannya menentukan pekerjaan
migrasinya.

---

## C — Penyimpanan

### C-1. ~~Menyimpan addendum memanggil penyimpan detail kontrak DAN penyimpan detail addendum~~ — **DIPINDAHKAN**

**Butir ini tidak lagi ditunda, dan pernyataan aslinya keliru.**

Versi pertama butir ini menyatakan `SaveTreatyIn_EDM_Act` memanggil kedua penyimpan detail. Itu
benar sebagai **yang tertulis** dan salah sebagai **yang berjalan**: panggilan ke penyimpan detail
kontrak ber-blok `//`, begitu pula langkah yang menulis ke `M_TREATY_IN`.

Saya membaca daftar langkahnya tanpa memeriksa bloknya — melewati CONTEXT.md §2.1, aturan yang
mendahului semua pembacaan lain.

Hasil pemeriksaan yang benar, beserta akibatnya, kini ada di
`KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` §G1b, **sebagai temuan Treaty In** — berkasnya ada di irisan
323, jadi ia berjalan pada jalur Treaty In hari ini dan bukan bahan parkiran.

### C-2. Tidak ada tabel selisih sama sekali

**Yang terlihat.** Sapuan seluruh SQL di kedua ekspor tidak menemukan satu pun tabel yang menyimpan
nilai selisih. Selisih hidup di dalam `JSONDATA` sebagai sub-pohon cermin `ValueDifference.*` — 161
jalur di bawah 31 akar.

**Bukti.** Daftar tabel yang disentuh seluruh `RDBList/` dan `ReportDefinition/` kedua ekspor;
tidak ada nama yang menyerupai tabel selisih.

**Kenapa ditunda.** Ini **sudah dipakai** di gerbang G3 sebagai fakta bentuk — tabel selisih adalah
benda baru. Yang ditunda adalah pertanyaan turunannya: **apa yang selama ini dilaporkan sebagai
selisih**, dan dari mana angkanya dibaca.

---

## D — Warisan penyakit yang sama dengan Treaty In

### D-1. `PositionUsername` ada juga di kelas addendum

**Yang terlihat.** `Int-treaty_in_edm` membawa `PositionUsername` — nama orang sebagai teks yang
menempel pada baris.

**Bukti.** Inventaris properti kelas `ASM-FW-GISFW-Int-treaty_in_edm`, ekspor Adjustment.

**Kenapa ditunda.** Di sisi Treaty In ini sudah diputuskan **dibuang** (CONTEXT.md §2.9,
`SPEC-MODEL-DATA.md` §12.5). Keputusan yang sama **kemungkinan besar** berlaku di sini, tetapi
"kemungkinan besar" bukan dasar memutuskan modul yang belum dibedah.

### D-2. Tombol paksa yang sama terpasang di layar Adjustment

**Yang terlihat.** `TreatyInForceEdit`, `TreatyInForceResolveComplete`, dan `TreatyInReturntoInputor`
ada di ekspor Adjustment dengan `pzInsKey` yang **identik** dengan yang di Treaty In.

**Bukti.** Perbandingan `pzInsKey`; ketiganya ada di irisan 323 berkas.

**Kenapa ditunda.** Karena aturannya **benar-benar aturan yang sama**, temuan di sisi Treaty In
berlaku pada bendanya. Yang **belum** diperiksa adalah apakah **layar Adjustment** memasangnya
dengan kondisi tampilan yang sama — dan itu pemeriksaan layar Adjustment, milik sesinya.

---

## E — Yang sengaja TIDAK dilihat

Supaya sesi Adjustment tahu apa yang belum tersentuh, dan tidak mengira daftar ini menyeluruh:

| Tidak dilihat | Alasan |
|---|---|
| 14 seksi dan flow-action ber-akhiran `OldData` | perilaku layar — inti modul Adjustment |
| `TreatyEDMCalculateDifference` dan empat aturan `TreatyEDMDifference*` | **cara selisih dihitung** — inti yang paling tidak boleh disentuh sebelum sesinya |
| `TreatyInSetAddendumToHistory` | perilaku riwayat |
| `TreatyCreateEDM` selain langkah pembuka | mesin keadaan addendum |
| `RefreshAchievement`, `TreatyRevisionCopyAttachment` | di luar sambungan |
| Pemetaan atribut `Int-treaty_in_edm` satu per satu | dilarang eksplisit |

**Daftar ini tidak lengkap, dan tidak dimaksudkan lengkap.** Ia memuat apa yang tidak dapat
dihindari saat mencari sambungannya.
