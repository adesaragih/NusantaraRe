# Sisa dan serah terima — sesi to-spec lapisan data Treaty In

**Tanggal:** 24 September 2026
**Keadaan:** **keenam langkah selesai.** `L-1` ditutup.
**Sepuluh keputusan P-1…P-10 dijawab pemilik proses dan SELURUHNYA sudah dikerjakan** —
lihat §7.

---

## 1. Gerbang selesai — dijawab satu per satu

| Gerbang | Keadaan |
|---|---|
| `2-to-spec/` berdiri: `SPEC-MODEL-DATA.md`, `KAMUS-KOLOM.md`, `ddl-usulan/` | **YA**, dengan satu penyimpangan yang diputuskan pemilik proses: **`SPEC-MODEL-DATA.md` tetap di akar modul.** Memindahkannya memutus ratusan rujukan; membiarkannya dapat dibatalkan kapan saja |
| 14 paket uang bertingkat, atau terdaftar sebagai pertanyaan | **YA** — 8 terbaca dari ekspor, 6 menjadi `T-1`…`T-5` |
| ADR-0056 ditulis; kelima procedure terpetakan habis | **YA** — 674 baris dibaca seluruhnya |
| lima kunci alami ber-invarian dan ber-uji negatif | **YA** — INV-64…INV-68, uji `N-6a`…`N-6e` + **empat uji positif** |
| **L-1 ditutup** | **YA** — `ddl-usulan/` **36 berkas** sesudah P-8 |
| L-8, L-4, L-6 dinyatakan lingkupnya | **YA** — L-6 ditutup, L-8 ditutup sebagian, L-4 diberi pemilik |
| angka §1, §3, §4, §7 dan kepala berkas konsisten | **YA** — §1a ditulis; kepala berkas diperbaiki |
| daftar sisa memisahkan **yang menahan** dari **yang mengukur kerusakan** | **YA** — §5 |

---

## 2. Yang dihasilkan

| Berkas | Isi |
|---|---|
| `2-to-spec/KAMUS-KOLOM.md` | **34 entitas, 240 kolom** — dibangkitkan (§7.2) |
| `2-to-spec/ddl-usulan/` | **36 berkas**: 34 tabel, `00_SKEMA_DAN_AKUN.sql`, `Z00_KUNCI_ALAMI.sql` |
| `2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md` | 8 terbaca, 6 ditanyakan |
| `2-to-spec/PERMINTAAN-TEKNIK-TREATY.md` | **sebelas butir** — sekaligus `daftar wawancara` yang selama ini dirujuk tanpa pernah ada |
| `2-to-spec/PEMETAAN-PROCEDURE.md` | lima procedure, tiga golongan |
| `2-to-spec/INVENTARIS-BERKAS-INDUK.md` | 99 berkas; 7 biner ditolak dengan namanya |
| `2-to-spec/TEMUAN-0-5-ACTUALVALUE-TERBALIK.md` | koreksi berkas induk |
| `docs/adr/0056-tanpa-stored-procedure.md` | K-4 |
| `SPEC-INVARIAN.md` §7 + INV-64…INV-68 · `UJI-NEGATIF-INVARIAN.md` N-6 | lima kunci alami |
| `SPEC-MODEL-DATA.md` kepala · §1a · §2.3 | koreksi, rekonsiliasi angka, tujuh + satu entitas |
| `5-tiket/LUBANG-SPESIFIKASI.md` L-1 · L-4 · L-6 · L-8 | penutupan dan pernyataan lingkup |
| `alat/` — **6** perkakas baru | disimpan sebagai berkas, bukan perintah sekali pakai |

### Satu definisi, dua keluaran

`alat/buat-kamus-dan-ddl.py` **mengurai** §10 menjadi `alat/definisi-skema-treaty-masuk.json`;
`alat/tulis-kamus-dan-ddl.py` menulis **kamus dan DDL dari berkas definisi yang sama**. Keduanya
**tidak dapat berbeda**, dan tidak ada satu pun nama kolom yang diketik ulang dengan tangan.

---

## 3. Yang sengaja TIDAK ada di dalam `ddl-usulan/`, dan masing-masing berupa pernyataan keputusan

Ketiganya ditulis sebagai **apa · kenapa · akibat · dilihat di · ditagih** — bukan sebagai `TODO`,
supaya orang berikutnya **berhenti memperbaikinya dengan cara yang justru merusaknya**.

### 3.1 Bentuk fisik **paket uang** — ~~dan ini yang paling berat~~

> **DIGANTI 24 September 2026 oleh P-8 — lihat §7.1.** Sepuluh dari enam belas sudah berbentuk:
> golongan **A** menjadi tabel anak per mata uang, golongan **B** memperoleh kolom
> `KODE_MATA_UANG`. **`INV-08`, `INV-09`, dan `INV-12` dibebaskan.** Yang tersisa **golongan C**,
> enam kolom tanpa denominasi, menunggu `T-6`.
>
> Uraian di bawah **dipertahankan sebagai alasannya**, bukan sebagai keadaan sekarang.

**16 atribut di 9 entitas** bertipe `U1 paket uang`. Masing-masing muncul sebagai **satu kolom
angka**; **kolom mata uangnya tidak ada.**

§0 `SPEC-MODEL-DATA.md` menyatakan peleburannya terang-terangan — *"nilai + mata uang + kurs menjadi
**satu** atribut"* — tetapi peleburan itu **konseptual**. Bentuk fisiknya keputusan gerbang sesi DDL
(§12.6).

**Akibatnya serius, dan disebut supaya tidak terbaca sebagai kerapian:**

1. tabel-tabel itu menyimpan **jumlah uang tanpa denominasi**;
2. **tiga kunci alami tidak dapat dikompilasi** — `INV-08` (`RETENSI_CEDANT`) dan `INV-09` (`EGNPI`)
   berbunyi *"kelompok treaty + **kode mata uang**"*, `INV-12` (`TERMIN`) *"nomor termin + **kode
   mata uang**"* (§10.16a). Ketiganya berstatus **klaim, bukan penegakan**;
3. **ADR-0053** tentang daftar mata uang tidak dapat ditegakkan.

> Menulis `UNIQUE` **tanpa** mata uang bukan jalan tengah — ia **mengubah artinya**, melarang dua
> baris bermata uang berbeda, dan **menolak data yang sah**. Itu persis kekeliruan lingkup yang
> sudah tiga kali terjadi di modul ini.

### 3.2 Constraint tingkat pencatatan — enam paket

`BATAS_MAKSIMUM_KELOMPOK`, `BATAS_MAKSIMUM_NON_KELOMPOK`, `BATAS_PILIHAN`, `NILAI_BATAS`,
`LIMIT_AGREGAT`, `NILAI_EGNPI` — `INV-39` dan `INV-40` tidak dipasang. Ditagih ketika `T-1`…`T-4`
kembali.

### 3.3 Kunci asing untuk dua induk polimorfik

`POTONGAN.ID_INDUK_POTONGAN` (§14.1) dan `PENYEBARAN.ID_INDUK_PENYEBARAN` (§12.3) berinduk **dua**.
`INV-17` melarang rujukan yang sasarannya bergantung nilai kolom lain, jadi bentuk fisiknya —
dua kolom nullable atau tabel jembatan — keputusan sesi DDL.

---

## 4. Temuan baru atas berkas induk — dilaporkan, TIDAK ditambal

| # | Temuan | Bukti |
|---|---|---|
| **F-1** | **`7-2-ACTUALVALUE.md` §1.2–§1.4 terbalik.** Kode tindakan `1` = **lompat**, dikalibrasi 66/66 atas kedua ekspor. Yang dipotret jenis **1 dan 2**, bukan 3 | `TEMUAN-0-5-ACTUALVALUE-TERBALIK.md` |
| **F-2** | ~~**§10.2 mengaku 48 atribut, memuat 44**~~ — **DITUTUP 24 Sep 2026.** Pencarian keempatnya diselesaikan atas sumber mandiri: **nol calon tersisa**; angka judulnya yang salah, kini **45** | §12.2 · `COCOK-SILANG-CACAH-ATRIBUT.md` §4 |
| **F-3** | ~~**§10.21 mengaku 6 atribut, memuat 5**~~ — **DITUTUP 24 Sep 2026.** Nilainya sudah ada; satu baris memuat **dua nama** sehingga pengurainya menolaknya. Yang benar-benar hilang **peran** (ADR-0045 isi minimal). Judulnya **8** | §12.2 |
| **F-4** | ~~`STRUKTUR-DATA.md` belum memuat `PEMULIHAN_LIMIT`~~ — **DITUTUP P-3.** Dan lubangnya ternyata **dua entitas, bukan satu**: `PERISTIWA_KONTRAK` juga tidak ada, begitu pula di `ERD.md` dan perkakas pembangkit skema | butir 4 urutan wewenang |
| **F-5** | Kedua procedure rinci **tidak punya cabang `ELSE`**; jalur perbarui tidak berbuat apa pun dan tetap `StsSave := 1` | `PEMETAAN-PROCEDURE.md` §4.1 |
| **F-6** | **`TREATYINDETAILEDM` tertinggal tujuh kolom** dari `TREATYINDETAIL` — himpunan bagian murni | idem §4.2 |
| **F-7** | `DETAIL_PROPORSIONAL_FAKULTATIF` **tepat 30 bita** — mentok batas Oracle tanpa sisa | keluaran `buat-peta-nama-tabel-treatyin.py` |

| **F-15** | **`PremiumEarnedList` dan `ROLPct` DIHITUNG aturan yang hidup** — dan `PremiumEarnedList` sudah diadili **DITURUNKAN** di dua artefak induk ber-`EVIDENCED` | `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` |
| **F-16** | **§10 tidak memuat entitas peran maupun penugasan bertanggal**, padahal ADR-0044 menempatkan keduanya di gelombang 1. `PERAN_PELAKU` karenanya berdiri tanpa daftar nilai yang sah | §10.21, kotak `F-3` |
| **F-18** | **Lima tabel anak paket uang golongan A tidak pernah terdaftar di `STRUKTUR-DATA.md`** — berkas yang **MENGIKAT** soal daftar entitas. Mereka berdiri di `KAMUS-KOLOM.md` dan `ddl-usulan/` sejak `P-8`. **Instans ketiga** setelah `F-4` dan `L-6`. **Didaftarkan 24 Sep 2026** — menyalin keputusan yang sudah ada, bukan menambah keputusan | `STRUKTUR-DATA.md` §1.3a |
| **F-17** | **§10.19a memuat dua tabel**, yang kedua milik `VERSI_KONTRAK`. Pengurai memberikannya ke `DOKUMEN_ADDENDUM` — **kunci asing `VERSI_KONTRAK` → `DOKUMEN_ADDENDUM` hilang tanpa suara** | §10.2 kotak `F-17`; selisih DDL §12.3 |

> **F-2 dan F-3 SUDAH ditutup 24 September 2026, dan tidak satu pun ditambal dengan karangan** —
> yang satu ditutup dengan menyelesaikan pencariannya lebih dulu, yang lain dengan membaca ADR yang
> melahirkannya. `KAMUS-KOLOM.md` §0 kini berbunyi **"tidak ada selisih"**, dan kalimat *"nol
> selisih bukan lengkap"* tetap berdiri di sana.
>
> **Empat temuan baru — `F-15`, `F-16`, `F-17`, `F-18`.** Tiga di antaranya **dilaporkan, tidak
> ditambal**, masing-masing dengan pemilik dan penagih yang disebut namanya. Yang **dikerjakan**
> hanya `F-18`, dan sebabnya satu kalimat: mendaftarkan kelima tabel anak itu **tidak menambah
> keputusan apa pun** — `P-8` sudah memutuskannya, yang kurang hanya salinannya di berkas yang
> mengikat.

---

## 5. Daftar sisa — dipisah menurut pembacanya

> **Daftar ini bertanggal SEBELUM P-1…P-10 dikerjakan.** Yang berlaku sekarang **§7.4**; yang di
> bawah dipertahankan sebagai jejak, dengan butir yang sudah tertutup dicoret di tempatnya.

### 5.1 Yang MENAHAN

| # | Sisa | Siapa |
|---|---|---|
| ~~S-1~~ | ~~bentuk fisik paket uang~~ — **DIPUTUSKAN P-8** untuk 10 dari 16; sisanya `T-6`. Daftar sisa yang berlaku ada di **§7.4** | — |
| S-2 | angka presisi per kelompok tipe | gerbang sesi DDL |
| S-3 | bentuk fisik dua induk polimorfik | gerbang sesi DDL |
| S-4 | pengurai §10 kehilangan 4 atribut `VERSI_KONTRAK` dan 1 `JEJAK_PERUBAHAN` — **karena §10 sendiri tidak memuatnya** (F-2, F-3) | sesi to-spec |
| ~~S-5~~ | **DITUTUP P-3** | — |
| S-6 | `NILAI_PENYEBARAN` — kedudukannya tertahan **Uji AD**; invariannya belum dinomori | pemilik data, lalu to-spec |
| S-7 | enam paket uang belum bertingkat → `T-1`…`T-5`, **ditambah `T-6`** | **teknik treaty** |
| S-8 | **butir 3 urutan wewenang TETAP DICORET** — diputuskan P-7, dan syarat pengembaliannya kini tertulis di §4.3 | **pemilik proses** |

### 5.2 Yang hanya MENGUKUR KERUSAKAN

| # | Sisa | Yang diukurnya |
|---|---|---|
| M-1 | **Uji AN** (baru) — baris ganda `TREATYINDETAIL` | akibat procedure tanpa `ELSE` |
| M-2 | konsumen hilir `TREATYINDETAILEDM` yang kurang tujuh kolom | siapa yang terdampak |
| M-3 | **Uji AM** — pertanyaannya **berubah** akibat F-1 | |
| M-4 | **340 properti** titik buta gelombang 2 | luas yang belum diadili |
| M-5 | tujuh berkas biner belum pernah dibuka, termasuk `ERD-ORACLE.xlsx` dan `Rekap SQL.zip` | mungkin memuat DDL yang gerbang DDL cari |
| M-6 | `L-9` lingkungan ekspor · `L-10` 21 jenis aturan | **kantor**; tetap tiket tertahan |

### 5.3 Eskalasi

| | |
|---|---|
| **NAIK** | `GET_TOKEN_STORAGE` menerbitkan token dari `STANDARD_HASH('ASMAPP' \|\| timestamp,'MD5')` — **dapat ditebak sepenuhnya, tanpa unsur acak, berlaku satu menit**. **Terbuka sekarang**, bukan utang migrasi. Di luar Treaty In |
| **TURUN** | butir `EGNPI` — sesudah penulisnya terbaca (`= .AmountEPI`), ia pertanyaan faktual yang sama dengan `EPI`. **Usulan: digabung ke `T-4`** |

---

## 6. Usulan suntingan berkas induk yang MENUNGGU persetujuan

| Berkas | Usulan |
|---|---|
| `4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` §1.2–§1.4, kepala | pembacaannya terbalik (F-1) |
| `CONTEXT.md` §2.1 | prasyarat **hidup** pun dibaca bersama **kode tindakannya** |
| ~~`SPEC-MODEL-DATA.md` §10.2, §10.21~~ | ~~selisih jumlah atribut (F-2, F-3)~~ — **DIKERJAKAN 24 Sep 2026** atas perintah putaran penutup; lihat §12.2 |
| `SPEC-MODEL-DATA.md` §10.23a | kolom Tingkat diisi untuk 8 paket; 6 sisanya menyebut nomor pertanyaannya |
| `DAFTAR-ESKALASI-MANAJEMEN.md` | butir token; butir `EGNPI` turun golongan |
| `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` | `PEMULIHAN_LIMIT` ditambahkan (F-4) |

**Tidak satu pun dikerjakan** — kecuali berkas yang langkah 3 dan 4 perintahkan langsung.

---

**To-ticket tidak dimulai.** Ia hanya dimulai atas perintah pemilik proses.

---

## 7. Sepuluh keputusan pemilik proses — dijawab dan dikerjakan, 24 September 2026

| # | Keputusan | Yang berubah di disk |
|---|---|---|
| **P-1** | koreksi `7-2-ACTUALVALUE.md` diterapkan sebagai blok bertanggal, uraian lama dipertahankan sebagai alasan | kepala, §1.2, §1.3, §1.4 — kolom *Dipotret?* dibalik; §1.4 berpindah alamat ke jenis 1 dan 2 |
| **P-2** | selisih jumlah atribut dinyatakan, **angka judul tidak diturunkan** | §10.2 dan §10.21 memperoleh blok selisih; dan sapuannya melahirkan **§3.2a** |
| **P-3** | dua entitas yang tertinggal ditambahkan di seluruh lapisan struktur data | `STRUKTUR-DATA.md`, `ERD.md` §2.3b, `alat/buat-skema-treaty-masuk.py` dijalankan ulang — **33 tabel, 44 relasi** |
| **P-4** | hasil tingkat pencatatan masuk §10.23a, keenam sisanya menyebut nomor pertanyaannya | §10.23a; angka lama 16/22/5 dipertahankan sebagai jejak |
| **P-5** | dua butir eskalasi ditulis | **butir 9** token dapat ditebak · **butir 10** tingkat perkiraan premi · tiga baris lampiran pengukuran |
| **P-6** | aturan baca prasyarat masuk `CONTEXT.md` §2.1 | perluasan bertanggal, dengan tabel kalibrasi kode tindakan |
| **P-7** | **butir 3 urutan wewenang TETAP DICORET** | §4.3 memperoleh peninjauan bertanggal dan **syarat pengembalian yang sebenarnya** |
| **P-8** | paket uang: golongan **A** dan **B** diputuskan, **C** ditahan | lihat §7.1 |
| **P-9** | satu berkas permintaan ke teknik treaty, sekaligus menjadi `daftar wawancara` | `2-to-spec/PERMINTAAN-TEKNIK-TREATY.md` — sebelas butir |
| **P-10** | sesudah ini: **kembali ke C3** grilling Adjustment | — |

### 7.1 P-8 — dan tiga invarian yang dibebaskannya

| Gol. | Jumlah | Putusan |
|---|---:|---|
| **A** | 5 | menjadi **tabel anak** per mata uang: `NILAI_MDP`, `NILAI_MDP_MINIMUM`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM`, `NILAI_CADANGAN_PREMI` |
| **B** | 5 | memperoleh kolom **`KODE_MATA_UANG`** pada barisnya sendiri |
| **C** | 6 | **ditahan** — menunggu `T-6`; pernyataan keputusan berdiri di berkas DDL masing-masing |

> **`INV-08`, `INV-09`, dan `INV-12` dibebaskan.** Ketiganya sempat berstatus **klaim** karena kunci
> alaminya menyebut *"kode mata uang"* dan kolomnya tidak ada. Golongan B menjawabnya sendiri.
> Rinciannya `SPEC-INVARIAN.md` §8.

### 7.2 Angka `2-to-spec/` sesudah kesepuluh keputusan

| | Sebelum P-1…P-10 | Sesudah |
|---|---:|---:|
| entitas di `KAMUS-KOLOM.md` | 29 | **34** |
| kolom | 215 | **240** |
| berkas `ddl-usulan/` | 31 | **36** |
| kunci asing | 28 | **33** |
| `UNIQUE` kunci alami | 20 | **23** |
| pengenal melebihi 30 bita (§16) | 0 | **0** — diperiksa ulang |

### 7.3 Temuan yang lahir SAAT mengerjakan kesepuluhnya

| # | Temuan |
|---|---|
| **F-8** | **`TreatyIn.NusareSharePct` dan `TreatyIn.BrokeragePct`** — nol penulis, dan keduanya menulis ke `TREATY_IN.NUSARESHAREPCT` dan `TREATY_IN.BROKERAGEPCT` pada setiap simpan. **Dua dari dua puluh kolom bisnis selalu kosong.** `BrokeragePct` bahkan **tidak ada di daftar §3.2**, jadi ia di luar semesta 667 yang pernah diadili. Ditulis sebagai §3.2a |
| **F-9** | `ddl-usulan/` **tidak dibersihkan antar-jalan**, sehingga satu tabel sempat berdiri dua kali dengan nomor berbeda. Perkakas diperbaiki: ia kini menghapus keluaran lama **dan melaporkan berapa yang dihapusnya** |
| **F-10** | lima tabel anak golongan A **belum punya nomor invarian** untuk kunci alaminya — `SPEC-INVARIAN.md` §8.3, dengan pemilik dan penagihnya |

### 7.4 Yang TETAP menahan sesudah kesepuluhnya

| # | Sisa | Siapa |
|---|---|---|
| **S-1** | angka **presisi** per kelompok tipe | gerbang sesi DDL |
| **S-2** | bentuk fisik **dua induk polimorfik** (`POTONGAN`, `PENYEBARAN`) | gerbang sesi DDL |
| **S-3** | paket uang **golongan C** — enam kolom tanpa denominasi | `T-6`, teknik treaty |
| **S-4** | empat atribut `VERSI_KONTRAK` dan satu `JEJAK_PERUBAHAN` yang **diumumkan tetapi tidak pernah didaftar** | sesi to-spec |
| **S-5** | kunci alami lima tabel anak golongan A belum bernomor | sesi to-spec |
| **S-6** | `NILAI_PENYEBARAN` — kedudukannya tertahan **Uji AD** | pemilik data |
| **S-7** | sebelas butir `PERMINTAAN-TEKNIK-TREATY.md` | teknik treaty |
| **S-8** | **butir 3 urutan wewenang tetap dicoret** sampai gerbang DDL tertutup | pemilik proses |

> **Gerbang sesi DDL kini tertahan oleh dua hal, bukan tiga.** Bentuk fisik paket uang — yang semula
> paling berat karena tiga invarian menggantung padanya — **sudah diputuskan untuk sepuluh dari enam
> belas**, dan keenam sisanya tidak menahan invarian mana pun.

---

## 8. Penelusuran jalur `JSONDATA` — LANGKAH 1, ditulis 24 September 2026

`2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md` — **satu baris per jalur**, bukan per tabel.

| | Jumlah |
|---|---:|
| jalur beradjudikasi `PETA-TELUSUR-JSON.md` §6 | **667** |
| **dikecualikan** — pohon cermin milik modul Adjustment | **272** |
| **lingkup Treaty In** | **395** |
| DIPETAKAN · cocok tunggal | 101 |
| DIPETAKAN · cocok ganda *(calon, bukan putusan)* | 59 |
| DIPETAKAN · tidak tercocokkan | 53 |
| DITURUNKAN | 166 |
| DIBUANG | 12 |
| DITUNDA | 4 |

**Menutup: 395 + 272 = 667. Tidak ada jalur tanpa nasib.**

> **Klaim *"seluruh jalur JSON punya rumah"* tidak ditulis di mana pun**, dan itu disengaja.
> Yang ditulis: *seluruh jalur di dalam semesta yang diperiksa punya nasib, dan semestanya kurang
> **340 properti** (`L-8`).*

### 8.1 Lima puluh tiga yang tidak tercocokkan — digolongkan bersebab

| Golongan | Jumlah | Artinya |
|---|---:|---|
| GEL-2 | 17 | cabang fakultatif keluar, di luar gelombang 1 |
| TABEL | 13 | jalur menamai **daftarnya** — ia menjadi tabel, bukan kolom |
| NAMA-MASTER | 7 | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) |
| DIBUANG-ULANG | 4 | sudah diputuskan di §12.5, §13, §14.5 |
| GOL-C | 4 | mata uang paket uang golongan C — menunggu `T-6` |
| **BERSUSUN** | **4** | **daftar kelas bisnis / kelompok treaty — belum punya rumah bernama** |
| AGREGAT | 2 | berpasangan dengan besaran agregat yang tidak disimpan |
| **YATIM** | **2** | **tidak ada rumah, dan tidak ada alasan tertulis** |

### 8.2 Temuan baru

| # | Temuan |
|---|---|
| **F-11** | **`TreatyYear` tidak punya kolom di skema baru** — tidak ada `TAHUN_TREATY` maupun padanan. Ia bertanda **DIPETAKAN** di peta telusur, dan merupakan **salah satu dari 20 kolom bisnis `TREATY_IN`**, ditulis setiap penyimpanan. Menebaknya sebagai turunan tanggal mulai akan **benar untuk hampir semua kontrak** — persis *"cacat di balik nilai bawaan"*. Menjadi **`T-7`** |
| **F-12** | **Empat jalur BERSUSUN** — daftar kelas bisnis per kelompok treaty per layer **tidak punya tempat**. Kandidat kekeliruan lingkup yang sudah tiga kali terjadi (§2.2c); lingkupnya harus dibaca dari **aktivitas penambah barisnya** |
| **F-13** | **Lima kolom `CLOB` adalah keputusan perkakas saya, bukan §10** — §10 menulis `T` untuk kelimanya, dan satu entri peta menunjuk kolom yang **tidak ada**. Dicabut; kini **nol `CLOB`/`BLOB`**. Lubang penggantinya: **255 karakter kemungkinan terlalu pendek** untuk pengecualian dan ketentuan khusus → **Uji AP** |
| **F-14** | **Penambahan paket uang golongan A/B semula hanya hidup di penulis**, tidak tersimpan ke berkas definisi — sehingga klaim *"satu definisi"* belum sepenuhnya benar dan perkakas penelusur membaca definisi yang kurang. Dipindahkan ke pembuat definisi; tak-tercocokkan turun **62 → 53** |

---

## 9. Tabel kanonik yang BOLEH dirujuk modul Adjustment

Ditetapkan supaya modul Adjustment tidak menggantung pada bentuk yang masih bergerak.

### 9.1 Boleh dirujuk — bentuknya sudah terkunci

| Tabel | Kenapa aman |
|---|---|
| `KONTRAK` | lapisan beku; kunci alaminya memperingatkan, tidak melarang (ADR-0040) |
| **`VERSI_KONTRAK`** | **inti sambungan** — addendum **adalah** `VERSI_KONTRAK` (GRL-01). `ID_VERSI_KONTRAK_DASAR` menunjuk versi pendahulunya |
| `LAYER`, `DETAIL_PROPORSIONAL`, `BAGIAN` | cabang terkunci §12.3, kunci alaminya INV-05, INV-06, INV-64 |
| `PENYEBARAN`, `RINCIAN_PENYEBARAN` | INV-16, INV-65 |
| `MATA_UANG_KONTRAK`, `RETENSI_CEDANT`, `EGNPI`, `TERMIN` | INV-07, INV-08, INV-09, INV-12 — **keempatnya kini menyebut `KODE_MATA_UANG`** dan dapat dikompilasi |
| `PERIODE_PELAPORAN`, `PERIODE_AKUMULASI`, `SKALA_KOASURANSI`, `BATAS_PER_BAHAYA`, `PORTOFOLIO`, `DOKUMEN_KONTRAK` | INV-10…INV-14, INV-66, INV-67 |
| keenam tabel acuan | INV-68; ADR-0038 |
| `CATATAN_PERSETUJUAN`, `JEJAK_PERUBAHAN`, `PERISTIWA_KONTRAK` | tanpa kunci alami, **dan itu keputusan** (§7.4 `SPEC-INVARIAN.md`) |

### 9.2 BELUM boleh dirujuk — bentuknya masih bergerak

| Tabel / kolom | Kenapa |
|---|---|
| `POTONGAN` — kolom `ID_INDUK_POTONGAN` | induk **polimorfik**; bentuk fisiknya keputusan gerbang DDL |
| `PENYEBARAN` — kolom `ID_INDUK_PENYEBARAN` | idem |
| `NILAI_PENYEBARAN` | kedudukan entitasnya **ditangguhkan** Uji AD (§10.23c butir 3) |
| lima tabel anak `NILAI_*` golongan A | baru lahir P-8; kunci alaminya **belum bernomor** |
| enam kolom paket uang **golongan C** | berdiri **tanpa denominasi** sampai `T-6` dijawab |
| lima kolom teks bebas `VERSI_KONTRAK` | lebarnya belum diukur — Uji AP |

### 9.3 Yang modul Adjustment TAMBAHKAN sendiri, bukan diminta dari sini

`NILAI_SELISIH` dan `BESARAN_DAPAT_DISESUAIKAN` — §11.3. Keduanya **tidak ada** di `ddl-usulan/`,
dan `ERD.md` mencatat **dua relasi lintas sekat** yang bermuara ke `NILAI_SELISIH`: keduanya
**mengikat dua sesi sekaligus dan tidak boleh diubah sepihak**.

---

## 10. Urutan wewenang — butir 3 tetap dicoret, penggantinya dipasang (24 September 2026)

| | |
|---|---|
| **Keputusan** | butir 3 *(berkas DDL mengalahkan spec soal bentuk tabel)* **TETAP DICORET**. Ditambahkan **butir 5**: *soal **nama kolom dan tipe**, `2-to-spec/KAMUS-KOLOM.md` **mengikat*** |
| **Kenapa bukan pemulihan** | butir 3 dibuat untuk mengadili **pertikaian** antara DDL dan spec soal bentuk tabel. Pertikaian itu **tidak dapat terjadi lagi** — keduanya dibangkitkan dari **satu berkas definisi** yang diurai dari §10 |
| **Kenapa penggantinya lebih kuat** | berkas DDL yang **ditulis tangan** bisa menyimpang dari spec; yang **dibangkitkan** tidak bisa. Bila menyimpang, **alatnya** yang salah |
| **Akibatnya pada to-ticket** | penahan *"kriteria selesai tiket tidak dapat menyebut nama kolom"* **jatuh** — `KAMUS-KOLOM.md` memberi **240 kolom bernama** |

**Syarat pengembalian butir 3 tidak berubah** — presisi, bentuk fisik paket uang golongan C, dan
bentuk fisik dua induk polimorfik. Dua di antaranya masih terbuka.

### 10.1 Penahan to-ticket yang benar-benar tersisa — **SATU** *(diperbarui 24 Sep 2026)*

| # | Penahan | Siapa mencabutnya |
|---|---|---|
| 1 | paket **`REV-1` … `REV-6`** belum ditanggapi | pemilik ADR induk |

**Tiga yang dicabut putaran penutup, beserta sebabnya — bukan sekadar dicoret:**

| Bekas penahan | Dicabut oleh | Sebab pencabutannya |
|---|---|---|
| angka **presisi** per kelompok tipe | **`KTV-A`** | ia tidak pernah memblokir **model**, hanya **kepastian**. Ongkos salahnya tidak setangkup — terlalu lebar murah, terlalu sempit memotong data — jadi sisi murahnya diambil, dengan **syarat pembalikan bertenggat**: sesi DDL boleh mempersempit, **hanya sebelum data dimuat** |
| bentuk fisik **dua induk polimorfik** | **`KTV-B`** | dua kolom bernama **tidak melanggar `INV-17`**, dapat dilebur menjadi jembatan kelak, dan **memasang dua kunci asing sungguhan** yang sebelumnya tidak ada sama sekali |
| enam paket uang **golongan C** tanpa denominasi | **`KTV-C`** | kolom nullable murah, dan `T-6` **menyempitkan** — bukan memblokir. Satu di antara keenamnya ternyata golongan **B** (§10.18 menulis *"+ mata uangnya"*), jadi yang tersisa **lima** |

> **Yang sudah TIDAK menahan:** §10 (selesai, cacahnya cocok untuk seluruh 35 entitas), butir 3
> urutan wewenang, presisi, induk polimorfik, dan golongan C.
>
> **Satu penahan yang tersisa bukan soal bentuk data.** `REV-1`…`REV-6` merevisi **alasan** dan
> **daftar keadaan**, dan nol dari enam menentukan letak kolom — diperiksa satu per satu di
> `treaty-in-adjustment/KEPUTUSAN-TANPA-VERIFIKASI-ADJUSTMENT.md` butir `KTV-4`. Yang paling
> mungkin berakibat: penolakan **`REV-3`**, sebab **keadaan adalah kolom**.

**Dua hal yang bukan penahan to-ticket tetapi menunggu satu kalimat pemilik proses:**

| # | Butir | Akibat bila didiamkan |
|---|---|---|
| 1 | **`F-15`** — nasib `PREMIUM_EARNED` dan `PERSEN_ROL` | satu kolom berdiri di `LAYER` dengan bukti yang membantahnya, dan satu entitas yang diperintahkan **tidak** berdiri. Keduanya terlihat di `KAMUS-KOLOM.md` |
| 2 | **`D-7`** modul Adjustment | parkir yang **dasarnya sudah lewat**; `F-1` sudah diterapkan |

## 11. Dua temuan diserahkan dari modul Adjustment — 24 September 2026

Keduanya `SISI` **IRISAN** — temuannya **milik Treaty In**, diadili di modul Adjustment untuk jalur
addendum saja. Keduanya **belum melewati ronde TDA mana pun**; diserahkan ke sesi ini karena §10
**sudah terbuka**, dan mengadilinya di sini lebih murah daripada membuka ronde baru.

### 11.1 `TDA-17` — tiga kolom tanpa rumah, dan sebabnya `L-8`

`TREATYINDETAILEDM` tertinggal tujuh kolom dari `TREATYINDETAIL`. **Empat punya rumah** —
`DEDUCTIBLE` → `LAYER.DEDUCTIBLE`, `MDP` → `LAYER.MDP`, `MDP_PCT` → `LAYER.PERSEN_MINIMUM_DEPOSIT`,
`ADJ_RATE` → `LAYER.PERSEN_PENYESUAIAN`.

**Tiga tidak: `DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT`.**

> **Ketiganya BUKAN sengaja dibuang.** Tidak ada keputusan mana pun yang membuangnya. Mereka muncul
> di `PETA-TELUSUR-JSON.md` **hanya di dalam pohon cermin**, tidak pernah di pohon utama, dan dua ada
> di `datar-titik-buta-pohon.csv` bertanda `TIDAK`. **§10 tidak melewatkannya — §10 tidak pernah
> diperlihatkan kepadanya.**

**Nol aturan membacanya** dari tabel datar: ketujuhnya hanya muncul di `SaveTreatyInDetail.xml`,
sebuah **penulis**. Kalibrasi lulus (`TREATYID`, 21 rujukan).

> ### DIADILI 24 September 2026 — dan pembingkaian di atas **tidak bertahan**
>
> Kalimat *"muncul hanya di pohon cermin"* **keliru**: ketiganya ada di pohon utama
> `datar-treatyin-lama.csv` (`TDA-17-PUTUSAN.md` §1). Dan *"nol pembaca"* benar tetapi menjawab
> pertanyaan lain — ia menyapu **kolom tabel datar**, bukan besarannya.
>
> | Kolom | Putusan | Dasar |
> |---|---|---|
> | `DEDUCTIBLE2` | **DISIMPAN** — `LAYER.DEDUCTIBLE_KEDUA` | sapuan penulisnya **bersih**: tidak ada aturan yang menurunkannya dari atribut lain |
> | `ROL_PCT` | **DIPERTIKAIKAN** — kolomnya berdiri, **ditandai** | `DetailCalculationROL`: `premi ÷ limit × 100` |
> | `PREMIUM_EARNED` | **TIDAK dibuat** | `DetailCalculation` langkah 25, hidup: `EGNPI × tarif penyesuaian`; dan dua artefak induk sudah mengadilinya **DITURUNKAN** ber-`EVIDENCED` |
>
> **Penanda *"PULIH dari `L-8`"* terhitung DUA, bukan tiga**, dan itu keadaan yang benar sampai
> pemilik proses memutuskan.

| | |
|---|---|
| **Siapa menutup** | ~~sesi ini~~ **pemilik proses** — dua jalan tertulis di `F-15` §5 |
| **Yang menagih** | ketiadaan `PREMIUM_EARNED` di `KAMUS-KOLOM.md`, penanda ⚠ pada `PERSEN_ROL` di §10.3, dan berkas `F-15` |

### 11.2 `TDA-18` — **PERUBAHAN**, dan ia mengubah klasifikasi materialitas

| | |
|---|---|
| **Temuan** | **`ReinstatementPct` disalin, tidak dikurangi.** Di `TreatyEDMDifferenceLimits` bentuknya `ValueDifference…ReinstatementPct = .ReinstatementPct` |
| **Bukti** | sapuan **seluruh korpus**, 708 berkas, kedua ekspor, dua bentuk penulis. **Kalibrasi LULUS** — `Limit` ditemukan **6 pengurangan nyata** terhadap `OLDDATA`. Hasil: **10 penugasan, NOL pengurangan**; `ReinstatementValue` **4 penugasan, nol** |
| **`SISI`** | **IRISAN** — `TreatyEDMDifferenceLimits` ada di kedua ekspor |
| **Akibat di sistem lama** | **perubahan persentase reinstatement tidak pernah menghasilkan baris selisih** |

> ### LABEL: **PERUBAHAN** — dan ia masuk daftar PERUBAHAN, bukan hanya daftar temuan
>
> **Versi yang HANYA mengubah persentase reinstatement terbaca NON MATERIAL di sistem lama, dan
> MATERIAL di sistem baru.**
>
> Berubah dari: *"perubahan reinstatement tidak menghasilkan selisih, sehingga tidak pernah
> menaikkan materialitas"*.
>
> **Itu perubahan perilaku yang DISENGAJA** — sistem baru menangkap perubahan yang sistem lama
> lewatkan, sebab `INV-69` melarang versi `TIDAK_MATERIAL` punya baris selisih bertipe porsi, dan
> persentase reinstatement **adalah porsi**.
>
> **Yang mengubah klasifikasi materialitas akan dilihat penyetuju di hari pertama.** Ia harus
> diberitahukan **sebelum** peralihan, bukan sesudah — dan karena itu ia berdiri di daftar
> PERUBAHAN, bukan hanya sebagai temuan yang dicatat.

| | |
|---|---|
| **Arah dampak bila salah** | bila bisnis menghendaki reinstatement **tidak** menaikkan materialitas, `INV-69` dikecualikan untuk besaran itu — **satu baris pengecualian**, dan ia harus ditulis sebagai pengecualian bernama, bukan dilupakan |
| **Siapa menutup** | **pemilik proses**, bersama pemberitahuan pra-peralihan |

---

## 12. Putaran penutup to-spec — 24 September 2026

Empat butir diperintahkan. **Tiga selesai; butir 1 berhenti dan dilaporkan**, sesuai aturan
*"bila isinya ternyata keliru saat diterapkan, berhenti dan laporkan — jangan betulkan diam-diam"*.

| # | Perintah | Keadaan |
|---|---|---|
| 1 | entitas `NILAI_PREMI_DIPEROLEH` berdiri; penanda *"PULIH dari `L-8`"* terhitung **3** | **BERHENTI** — `F-15`. Penandanya tetap **2** |
| 2 | `F-2` dan `F-3` ditutup | **SELESAI**, dan **`F-17`** lahir saat menutupnya |
| — | *(tidak diperintahkan)* pencocokan dua arah `STRUKTUR-DATA.md` ↔ `KAMUS-KOLOM.md` | **`F-18`** — lima entitas hilang dari berkas yang mengikat; **didaftarkan** |
| 3 | `KTV-A`, `KTV-B`, `KTV-C` tertulis beralasan | **SELESAI** — `KEPUTUSAN-TANPA-VERIFIKASI.md` §7 |
| 4 | keadaan to-ticket dinyatakan | **SELESAI** — §10.1 di atas, dan §12.4 di bawah |

### 12.1 Butir 1 — kenapa entitasnya tidak berdiri

**Bentuk menentukan BAGAIMANA sesuatu disimpan; ia tidak menentukan APAKAH sesuatu disimpan.**

| Yang diperintahkan | Yang terbaca dari sumber |
|---|---|
| *"`PremiumEarnedList` adalah `Page List` berisi `Currency` dan `Value` — golongan A menurut `P-8`"* | benar, dan tidak dibantah |
| *"Bentuknya sudah menentukan rumahnya"* | **tidak.** Nasibnya sudah diadili **DITURUNKAN** di dua artefak induk dengan dasar `EVIDENCED`, dan putaran ini menemukan **rumus yang menghitungnya**: `DetailCalculation` langkah 25, hidup, `PremiumEarnedList = EGNPI per mata uang × tarif penyesuaian` |

Temuan yang sama mengenai kolom yang **sudah mendarat**: `PERSEN_ROL` dihitung
`DetailCalculationROL` sebagai `premi ÷ limit × 100`. **Tidak dicabut sendiri** — ditandai di §10.3
dan ditagihkan.

Kalibrasinya lulus dan itu yang membuatnya dapat dipercaya: `Deductible2` disapu dengan cara yang
sama dan kembali **bersih** — **`DEDUCTIBLE_KEDUA` bertahan**.

Dua jalan beserta ongkos salahnya, dan rekomendasi saya, ada di
[`F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md`](F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md). **Satu
kalimat pemilik proses cukup.**

### 12.2 Butir 2 — `F-2`, `F-3`, dan `F-17`

| Lubang | Bagaimana ditutup |
|---|---|
| **`F-2`** §10.2 mengaku 47/48, memuat 44 | **angka judulnya yang salah, bukan tabelnya.** Pencariannya diselesaikan lebih dulu atas sumber mandiri: kelima calon §3.4 diadili, **nol menjadi atribut**. Judulnya kini **45** — 43 + §10.2a + §10.19a — dan ketiga angkanya dapat dijumlahkan pembaca |
| **`F-3`** §10.21 mengaku 6, memuat 5 | **nilainya sudah ada**; yang menyembunyikannya adalah **satu baris memuat dua nama**, yang pengurainya tolak. Barisnya dipecah. Yang benar-benar hilang terbaca dari **ADR-0045**: *"dan di bawah **peran apa**"* — `PERAN_PELAKU` ditambahkan. Judulnya **8** |
| **`F-17`** *(baru)* | §10.19a memuat **dua** tabel; yang kedua milik `VERSI_KONTRAK`. Pengurai memberikannya kepada `DOKUMEN_ADDENDUM`, sehingga **`VERSI_KONTRAK` kehilangan kunci asingnya ke `DOKUMEN_ADDENDUM`** tanpa suara. Ditutup di perkakas |

**Kelima calon `F-2`, beserta putusannya** — rinciannya `COCOK-SILANG-CACAH-ATRIBUT.md` §4:

| Calon | Putusan |
|---|---|
| `TreatyYear` | **turunan** — `@substring(TreatyIn.Commencement,0,4)` |
| `RevisionDate` | **dibuang** — `@CurrentDateTime()`, jejak mesin |
| `CoInScale` | **bukan atribut** — nol penulis, tak terdeklarasi, hanya di layar |
| `Ceding` · `LeadingReinsSource` | **sudah bersebab di §10.1** — sebabnya dicari di bagian yang salah |

**Penutupan `F-2` mencabut parkir `D-1` modul Adjustment.** `D-1` sudah **diterapkan** di §10.2a
dengan nomor urut yang benar (**ke-44**, bukan ke-48), dan tujuh diff Adjustment kini **enam
diterapkan, satu diparkir** — `D-7`, yang menunggu satu kalimat izin dan dasarnya sudah lewat.

Satu lubang baru yang **tidak ditambal**: **`F-16`** — ADR-0044 menempatkan **entitas peran dan
penugasan bertanggal** di gelombang 1, dan §10 tidak memuat satu pun dari keduanya. Akibatnya
`PERAN_PELAKU` berdiri sebagai teks beku **tanpa daftar nilai yang sah**. Pemiliknya pemilik proses,
bersama sesi yang mengerjakan ADR-0044.

### 12.3 Butir 3 — selisih `ddl-usulan/`, dicetak

Dibangkitkan ulang seluruhnya; keluaran lama dihapus lebih dulu dan **jumlahnya dilaporkan
perkakasnya** (37 berkas).

```
tabel sebelum 35, sesudah 35 — tidak ada tabel baru, tidak ada yang hilang
kolom ditambah/dihapus : 16
tipe berubah           : 85
```

| Perubahan | Berapa | Sebab |
|---|---:|---|
| `ID_*` `NUMBER(9)` → `NUMBER(19)` | 34 | `KTV-A` — kunci utama dan kunci asing tidak boleh berbeda lebar |
| `VARCHAR2(255 CHAR)` → `VARCHAR2(1000 CHAR)` | 41 | `KTV-A` — 255 **lebih sempit daripada sistem lama** |
| `VARCHAR2(255 CHAR)` → `VARCHAR2(4000 CHAR)` | 5 | `KTV-A` — teks bebas `VERSI_KONTRAK` |
| `NUMBER(11,8)` → `NUMBER(38,20)` | 4 | `KTV-A` — `P2` diseragamkan; `9989998` tidak muat |
| `ID_INDUK_POTONGAN` / `ID_INDUK_PENYEBARAN` dihapus | 2 | `KTV-B` |
| `ID_BAGIAN` + `ID_DETAIL_PROPORSIONAL` ditambahkan | 4 | `KTV-B` — **dengan kunci asing sungguhan**, dan `CHECK` tepat satu terisi |
| kolom mata uang golongan C bernama | 5 | `KTV-C` |
| `BATAS_PER_BAHAYA.KODE_MATA_UANG` | 1 | `KTV-C` koreksi 1 — `NILAI_BATAS` ternyata golongan **B** |
| `JEJAK_PERUBAHAN` — `PERAN_PELAKU`, `NILAI_SEBELUM`, `NILAI_SESUDAH` | 3 | `F-3` |
| `DOKUMEN_ADDENDUM.ID_DOKUMEN_ADDENDUM` menjadi `NOT NULL` | 1 | **`F-17`** — kunci utamanya sempat memakai keterisian kolom **kembarannya** yang salah tempat |
| `VERSI_KONTRAK.ID_DOKUMEN_ADDENDUM` | 1 | **`F-17`** — kunci asing yang hilang tanpa suara |

**Yang tidak bertambah: entitas.** Tetap **35**, dan itu disengaja — satu-satunya usul penambahan
putaran ini adalah `NILAI_PREMI_DIPEROLEH`, yang berhenti di `F-15`.

**`KAMUS-KOLOM.md` §0 kini berbunyi "tidak ada selisih"** untuk seluruh 35 entitas — pertama kali
sejak berkas itu ada. Kalimat *"nol selisih tetap bukan lengkap"* **tetap berdiri di sana**, sebab
semestanya masih kurang 340 properti.

### 12.4 Butir 4 — cacah kemampuan, dihitung ulang dan **tidak dicocokkan ke angka yang diberikan**

Perintah putaran ini menyebut *"62 kemampuan `P-01`…`P-59`"*. Dihitung mekanis atas
`5-tiket/DAFTAR-PEKERJAAN.md`:

| | |
|---|---:|
| kode yang **punya baris tabel** | **65** — `P-01`…`P-58` dan `P-60`…`P-66` |
| dicoret sebagai SIFAT (U-6) | 3 — `P-24`, `P-25`, `P-26` |
| **aktif** | **62** |

> **Angka 62 benar; rentangnya yang bukan `P-01`…`P-59`.** Ia `P-01`…`P-66`, dan **`P-59` tidak
> punya baris tabel sama sekali** — ia diputuskan di §2.2b *("syarat berbeda tiap pemulihan limit",
> bergolongan BARU)* dan tidak pernah didaftar. Itu lubang tersendiri, **dilaporkan bukan ditambal**:
> menambahkan barisnya berarti mengarang isi kolom Gol, Asal, dan Entitas yang tidak pernah
> ditetapkan siapa pun.

**Dan dua cacah di dalam berkas itu sendiri sudah basi**, keduanya ke arah yang berbeda — judul §2
menulis *"60 didaftar, 57 aktif"*, §4 Hitungan menulis *"58 didaftar, 55 aktif"*. **Tidak disunting
dari sini**: `DAFTAR-PEKERJAAN.md` milik sesi tiket, dan mengubah angkanya lewat sesi yang bukan
pemiliknya adalah kekeliruan yang sudah tercatat di §1a. **Ditagih pada sesi to-ticket, di gerbang
pembukanya.**

#### Tiket yang akan bertanda `TERTAHAN`, dan berapa dari 62

| Penahan | Kemampuan yang kena | Dari 62 |
|---|---|---:|
| **`REV-3`** — ADR-0055, **keadaan adalah kolom** | `P-01`, `P-29`, `P-31`, `P-32`, `P-33`, `P-34`, `P-35`, `P-37`, `P-42`, `P-56`, `P-57`, `P-58` — seluruh baris yang menyebut ADR-0055 sebagai **Asal**-nya | **12** |
| **`T-6`** — lima paket golongan C | **nol baris menyebutnya.** Kolomnya `BATAS_MAKSIMUM_KELOMPOK`, `BATAS_MAKSIMUM_NON_KELOMPOK`, `BATAS_PILIHAN`, `DEDUCTIBLE`, `LIMIT_AGREGAT`; yang paling dekat `P-04` *("mengisi seluruh kepala versi")*, yang memuatnya tanpa menyebutnya | **0 tersurat · 1 tersirat** |
| **`Uji AD`** (`S-6`, `NILAI_PENYEBARAN`) | `P-18` — dan ia **sudah** `TERTAHAN` oleh dua penghalang lain | **1, sudah terhitung** |

> **Batas pengukuran ini, dan ia wajib ikut terbaca:** ia menghitung baris yang **menyebut**
> penahannya di kolom *Asal* atau *Entitas*. Sebuah kemampuan dapat bergantung pada `REV-3` tanpa
> menyebut ADR-0055 — `P-04` adalah contohnya untuk `T-6`. **Angkanya lantai, bukan langit-langit.**

**`TERTAHAN` hari ini seluruhnya tujuh:** `P-18`, `P-20`, `P-41`, `P-49`, `P-60`, `P-61`, `P-62`.
Ketiga yang terakhir milik jalur Adjustment dan bertanda `DIASUMSIKAN-CLEAR`. **`P-62` bertenggat** —
kolomnya dicabut **sebelum data masuk** bila `DB-16b` dibantah.

#### Syarat pengembalian butir 3 urutan wewenang — dua dari tiga dicabut

| Syarat | Keadaan |
|---|---|
| angka **presisi** | **DICABUT** — `KTV-A` |
| bentuk fisik **dua induk polimorfik** | **DICABUT** — `KTV-B` |
| bentuk fisik **paket uang golongan C** | **DICABUT** — `KTV-C`; `T-6` menyempitkan, tidak memblokir |

**Ketiganya dicabut. Butir 3 tetap dicoret** — bukan karena syaratnya belum terpenuhi, melainkan
karena pertikaian yang melahirkannya **tidak dapat terjadi lagi**: DDL dan kamus dibangkitkan dari
satu berkas definisi. Yang berlaku **butir 5**: soal nama kolom dan tipe, `2-to-spec/KAMUS-KOLOM.md`
mengikat.
