# Cocok enam sumber struktur data — Treaty In dan Treaty In Adjustment

**Tanggal:** 25 September 2026 · **Perkakas:** [`alat/cocok-enam-sumber-struktur.py`](../alat/cocok-enam-sumber-struktur.py)
**Hasil gerbang:** **GAGAL** — tiga temuan, ketiganya bernomor di §3.

> **Berkas ini dibangkitkan dari perkakas, bukan ditulis tangan.** Setiap angka di bawah dapat
> dihasilkan ulang dengan satu perintah. Bila angkanya berbeda dari perkakasnya, **perkakasnya yang
> benar dan berkas ini yang basi**.

---

## 0. Kenapa enam, dan kenapa bukan tiga

Pendahulunya, [`alat/cocok-tiga-sumber-entitas.py`](../alat/_arsip/cocok-tiga-sumber-entitas.py), mengadu
tiga sumber dan mencetak **`lulus = True`**.

Ia lulus karena `NILAI_SELISIH` dan `BESARAN_DAPAT_DISESUAIKAN` **dikeluarkan lebih dulu** ke daftar
`LUAR`, dengan sebab yang ditulis di dalam perkakasnya sendiri:

```
'NILAI_SELISIH':             'modul Adjustment -- SPEC-MODEL-DATA sec 11.3',
'BESARAN_DAPAT_DISESUAIKAN': 'modul Adjustment -- tabel acuan, sec 11.3',
```

`SPEC-MODEL-DATA.md` §11.3 berjudul *"Entitas yang sengaja TIDAK dikerjakan"*, dan sebabnya tertulis
satu kata: **embargo**.

**Embargo itu sudah lewat.** Modul Adjustment sudah digrilling (`GRILL-A`…`GRILL-TDA`), sudah
to-spec (`USULAN-DIFF-KE-INDUK.md` — enam diff **DITERAPKAN**), dan sudah to-ticket (13 tiket;
tiket **`06`** bertanda **PEMBUAT PERTAMA `NILAI_SELISIH`**). Sebabnya hilang, **penagihnya tidak
dicabut** — dan pemeriksa yang mengecualikan justru apa yang seharusnya diperiksanya akan menjawab
**"LULUS" untuk kedua kemungkinan**.

Di gerbang yang diperlebar ini daftar `LUAR` menyusut dari empat menjadi **dua**: `RETRO_KELUAR`
(GEL-2) dan `PENCAPAIAN` (GEL-3). Keduanya tetap di luar karena **gelombangnya**, bukan karena
embargo.

---

## 1. Keenam sumber, dan masing-masing MENGIKAT atas hal yang berbeda

| | Sumber | Mengikat atas | Cacah terbaca |
|---|---|---|---:|
| **S1** | [`STRUKTUR-DATA.md`](STRUKTUR-DATA.md) | **daftar entitas** — urutan wewenang butir 4 | **39** |
| **S2** | [`2-to-spec/KAMUS-KOLOM.md`](../2-to-spec/KAMUS-KOLOM.md) | **nama dan tipe kolom** — butir 5 | **35** |
| **S3** | [`2-to-spec/ddl-usulan/`](../2-to-spec/ddl-usulan/) | apa yang **akan dibangun** | **35** |
| **S4** | [`peta-nama-tabel-treatyin.tsv`](peta-nama-tabel-treatyin.tsv) | nama **tabel datar** lama ↔ baru | **44** baris |
| **S5** | [`datar-treatyin-lama.csv`](datar-treatyin-lama.csv) | **pohon Pega apa adanya** (AS-IS) | **985** simpul |
| **S6** | [`PETA-TELUSUR-JSON.md`](../PETA-TELUSUR-JSON.md) §6 | **nasib** tiap jalur JSON | **667** jalur |
| **S7** | [`KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md`](KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md) | entitas dari **G1…G4** | **2** |

Nasib di S6: **DIPETAKAN 485 · DITURUNKAN 166 · DIBUANG 12 · DITUNDA 4**.
Simpul di S5 yang juga ada di ekspor Adjustment: **984 dari 985**.

### Dua berkas yang bukan saingan, dan salah satunya sempat dikira begitu

[`struktur-treatyin-lama.md`](struktur-treatyin-lama.md) dan [`STRUKTUR-DATA.md`](STRUKTUR-DATA.md)
memerikan **dua benda berbeda**, dan berkas yang pertama mengatakannya sendiri di kepalanya:

> Melengkapi `STRUKTUR-DATA.md` yang memerikan **entitas rancangan baru**. Berkas ini memerikan
> **bentuk lama apa adanya**.

Maka **tidak ada yang perlu dipilih di antara keduanya**:

| Berkas | Memerikan | Kedudukan |
|---|---|---|
| `STRUKTUR-DATA.md` | **model baru** | **MENGIKAT** soal daftar entitas (butir 4) |
| `struktur-treatyin-lama.md` | **sistem lama** | AS-IS; **asal**, bukan rancangan |
| `pohon-treatyin.txt` | sistem lama, bentuk **ringkas** | pendamping berkas di atas |
| `datar-treatyin-lama.csv` | sistem lama, bentuk **penuh** | yang dipakai perkakas |

`pohon-treatyin.txt` **tidak dipakai menghitung apa pun** di sini, dan sebabnya disebut: ia
**ringkasan** — daftar skalarnya dipotong dengan `(+91)`. Cacah simpul diambil dari CSV-nya.

---

## 2. Hasil per pemeriksaan

| | Pemeriksaan | Hasil |
|---|---|---|
| **A** | himpunan entitas S1 ↔ S2 ↔ S3 | **GAGAL** — 2 entitas didaftar tetapi tidak berkolom dan tidak ber-DDL |
| **B** | tipe kunci asing = tipe kunci utama yang dirujuk | **LULUS** — 38 kunci asing, 0 menggantung, **0 tipe berbeda** |
| **C** | padanan di peta nama masih ada di DDL | **GAGAL** — 17 padanan menyebut nama yang tidak ada |
| **D** | akar jalur Pega di peta nama ada di pohon | **LULUS** — 0 akar asing |
| **E** | tiap simpul pohon punya nasib di S6 | **358 dari 984 tanpa nasib** — ini `L-8`, diukur ulang |
| **F** | kolom `ID_*` yang bukan PK punya kunci asing | **11 tanpa kunci asing** — keluarga `F-19`/`F-20` |

**Pemeriksaan B yang lulus itu bukan formalitas.** Ia menjawab pertanyaan yang tidak pernah
ditanyakan: apakah ada kolom kunci asing yang tipenya berbeda dari kunci utama yang dirujuknya.
Jawabannya nol dari 38 — dan sekarang jawaban itu **tercatat**, bukan diandaikan.

---

## 3. Tiga temuan, bernomor

### `S-1` — dua entitas modul Adjustment TIDAK AKAN TERBANGUN

| | |
|---|---|
| **Apa** | `NILAI_SELISIH` dan `BESARAN_DAPAT_DISESUAIKAN` ada di S1 yang **mengikat**, dan **nol** di S2, **nol** di S3 |
| **Sebabnya** | `SPEC-MODEL-DATA.md` §11.3 — **embargo yang sudah lewat, penagih yang tidak dicabut** |
| **Arah dampaknya** | siapa pun yang membangun dari `ddl-usulan/` membangun **35 tabel**, lalu tiket `06` — PEMBUAT PERTAMA `NILAI_SELISIH` — **tidak punya tabel untuk dibuat**. Dan seluruh papan Adjustment bersandar padanya |
| **Bukan** | ini **bukan** keputusan yang kurang. G1…G4 sudah memutuskan bentuknya, dan `SEAM-ADJUSTMENT.md` §3 sudah menetapkan `KUNCI_PADANAN`-nya. Yang hilang adalah **pendaratannya**, bukan putusannya |

**Bentuknya sudah ada; atributnya belum.** Kedua entitas ini **belum pernah melewati §10**
`SPEC-MODEL-DATA.md` — penamaan, tipe, dan keterisian tingkat atribut — sebagaimana ke-35 lainnya.
Yang dapat dikumpulkan dari sumber yang ada: **5 kolom** berdasar untuk
`BESARAN_DAPAT_DISESUAIKAN`, **7 kolom** berdasar untuk `NILAI_SELISIH`. Ketujuh-belas-nya membawa
kutipan sumbernya masing-masing di [`ERD-SKEMA-BARU.html`](ERD-SKEMA-BARU.html).

**Tiga hal yang sumbernya DIAM, dan karena itu tidak diisi:**

1. **Paket uangnya belum lengkap.** G3 butir 2 berbunyi baris ber-satuan `UANG` *"wajib membawa
   mata uang **dan paket uangnya**"*. Paket uang menurut ADR-0007 / 0029 / 0039 lebih dari dua
   kolom — `SEAM-ADJUSTMENT.md` §3 menyebut `TINGKAT_PENCATATAN`, `PERSEN_BAGIAN_DIPAKAI`, `KURS`,
   `TANGGAL_KURS`, `SUMBER_KURS`, `NILAI_IDR` **pada bentuk BACA**. Mana di antaranya ikut
   **tersimpan** belum diputuskan.
2. **Dua sumber bertentangan soal berapa penunjuk versi yang disimpan** — lihat `S-3`.
3. **Tipe Oracle-nya belum melewati `KTV-A`.** Kelompok tipe di `KAMUS-KOLOM.md` §1 diputuskan
   untuk 35 tabel; kedua tabel ini tidak ikut.

### `S-2` — peta nama tabel datar BASI pada 17 baris

Dari 44 baris `peta-nama-tabel-treatyin.tsv`, padanannya terbagi empat, dan **pembagiannya memakai
daftar bernama, bukan tebakan dari bentuk namanya**:

| Golongan | Cacah | Keterangan |
|---|---:|---|
| masih ada di `ddl-usulan/` | **14** | sehat |
| sengaja **di luar gelombang** | **8** | fakultatif (6) · `BAGIAN_RETRO` GEL-2 · `PENCAPAIAN` GEL-3 |
| **penunjang**, bukan entitas model | **4** | `MIGRASI_*` (3) · `ARSIP_MUATAN_KELUAR` |
| akan mendarat bersama `S-1` | **1** | `NILAI_SELISIH` |
| **BELUM TERPADANKAN** | **17** | ← inilah yang basi |

Ketujuh belas itu: `KELAS_BISNIS_LAYER` · `NILAI_LAYER` · `BESARAN_LAYER` · `KELOMPOK_LAYER` ·
`KELAS_BISNIS_KELOMPOK` · `PENYEBARAN_XOL` · `NILAI_TERSEBAR` · `NILAI_BAGIAN` · `RETENSI` ·
`RINGKASAN_LIMIT` · `ANGSURAN` · `RINCIAN_ANGSURAN` · `AKUMULASI` · `BATAS_BAHAYA` ·
`REKAP_KONTRAK` · `NILAI_SEBELUM_PRO_RATE` · `BATAS_WEWENANG`.

**Sebagiannya jelas berganti nama** — `RETENSI` → `RETENSI_CEDANT`, `AKUMULASI` →
`PERIODE_AKUMULASI`, `ANGSURAN` → `TERMIN`, `BATAS_BAHAYA` → `BATAS_PER_BAHAYA`. **Sebagiannya
mungkin tidak menjadi tabel sama sekali** — `RINGKASAN_LIMIT` dan `REKAP_KONTRAK` berbunyi seperti
agregat, dan agregat bergolongan DITURUNKAN, tidak disimpan.

> **Keduanya tidak dipadankan di sini, dan itu disengaja.** Menebak padanan adalah persis cara
> sebuah peta nama menjadi salah tanpa ada yang tahu. Yang dilakukan berkas ini hanya
> **menghitungnya dan menyebut namanya**.

**Akibat yang terukur:** asal Pega hanya dapat ditempelkan ke **14** dari 37 entitas lewat peta
nama. Sesudah sumber kedua dipakai — `SPEC-MODEL-DATA.md` §2.3 dan §10.22, yang mutakhir — cakupan
naik menjadi **30 dari 37**. Tujuh sisanya memang tidak punya asal di pohon Pega, dan itu benar:
`PERISTIWA_KONTRAK` (§10.21a, BARU), kelima tabel anak paket uang (`F-18`, lahir dari `P-8`), dan
`BESARAN_DAPAT_DISESUAIKAN` (tabel acuan BARU).

### `S-3` — dua sumber bertentangan: BERAPA penunjuk versi yang disimpan

| Sumber | Bunyinya |
|---|---|
| `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` **G2** | *"Dua penunjuk disimpan, bukan satu"* — `ID_VERSI_KONTRAK_LAMA` dan `ID_VERSI_KONTRAK_BARU` |
| `STRUKTUR-DATA.md` **§1.1** | `ID_VERSI_KONTRAK_DASAR` *"adalah satu-satunya penunjuk yang disimpan, **dan ia tidak punya pasangan**"* |

Keduanya menjawab pertanyaan yang sama dengan jawaban yang berbeda, dan keduanya artefak to-spec
yang berlaku.

**Yang dipakai di gambar: SATU**, yaitu induk `VERSI_KONTRAK` — atas dasar butir 4 urutan wewenang,
yang menjadikan `STRUKTUR-DATA.md` mengikat soal daftar entitas **dan induknya**.

> **Itu pilihan wewenang, bukan adjudikasi.** Perkakas memilih sumber mana yang menang; ia tidak
> memeriksa mana yang benar. Putusannya milik pemilik proses.

**Arah dampak bila salah:** bila G2 yang benar, `NILAI_SELISIH` kurang satu kolom, dan selisih tidak
dapat dibaca ulang tanpa menelusuri induknya lebih dulu — **satu kolom ditambahkan**, tanpa
kehilangan data. Bila `STRUKTUR-DATA` yang benar dan G2 diikuti, satu fakta tersimpan **dua kali**
— `NILAI_SELISIH.ID_VERSI_KONTRAK_BARU` dan induknya sendiri — dan dua salinan satu fakta adalah
tempat lahirnya ketidakcocokan.

---

## 4. Yang bukan temuan baru, tetapi terukur ulang di sini

### `L-8` diukur mandiri: **358**, bukan **340**

`PETA-TELUSUR-JSON.md` §1 menyatakan *"Di luar keempat kategori: **nol**. Kriteria terima
terpenuhi"* atas semesta **667 jalur**. Pohon Pega yang sebenarnya berisi **985 simpul**.

Diadu satu per satu: **358 simpul tidak punya nasib** — **189 skalar, 169 daftar**.

Berkas itu sudah menyatakan semestanya kurang **340** (`L-8`, ditagih `M-4`). Pengukuran mandiri ini
memberi **358**. Selisih keduanya **18**.

> **Selisih 18 itu tidak diklaim sebagai temuan.** Perkakas menyamakan jalur dengan membuang `[]`
> dan awalan `TreatyIn.`, sehingga jalur yang berbeda **hanya pada tanda itu** terbaca sama, dan
> batas itu dicetak di dalam keluarannya sendiri. Yang dapat dikatakan: **angka 340 berordo benar,
> dan pengukuran kedua tidak membantahnya.**

### `F-19` / `F-20` — 11 kolom `ID_*` tanpa kunci asing

`JENIS_REASURANSI.ID_INDUK` · `KONTRAK.ID_KONTRAK_DISALIN_DARI` · `KONTRAK.ID_CEDANT` ·
`KONTRAK.ID_ASAL_BISNIS` · `VERSI_KONTRAK.ID_REASURADUR_PEMIMPIN` · `VERSI_KONTRAK.ID_KETUA_TREATY` ·
`VERSI_KONTRAK.ID_DAFTAR_RETRO` · `DOKUMEN_KONTRAK.ID_DOKUMEN` · `BAGIAN.ID_SUSUNAN_RETRO` ·
`DETAIL_PROPORSIONAL.ID_SUSUNAN_RETRO` · `PENYEBARAN.ID_JENIS_REASURANSI_INDUK`

**Sebagiannya memang tidak boleh berkunci asing** — empat menunjuk ke luar skema `TREATY_MASUK`
(`CEDANT`, `ASAL_BISNIS`, `REASURADUR`, `DOKUMEN`), dan itu keputusan yang tertulis di
`STRUKTUR-DATA.md` §4. Sisanya **menunjuk ke dalam skema sendiri** dan tidak punya alasan tertulis.

Perkakas ini **tidak memisahkan keduanya**, dan itu batasnya: ia melihat ketiadaan kunci asing, ia
tidak melihat alasannya.

---

## 5. Cara menjalankan ulang

```
cd D:\XML_NURE\_migration-docs\treaty-in
set PYTHONIOENCODING=utf-8
python alat\bangkitkan-erd-skema-baru.py      # ddl-usulan/  -> erd-skema-baru.json  (35)
python alat\cocok-enam-sumber-struktur.py     # enam sumber  -> himpunan-struktur.json
python alat\lengkapi-erd-adjustment.py        # + Adjustment -> erd-skema-terpadu.json (37)
python alat\bangkitkan-erd-html.py            #              -> ERD-SKEMA-BARU.html
python alat\bangkitkan-struktur-data-xlsx.py  #              -> STRUKTUR-DATA-SKEMA-BARU.xlsx
```

Kelimanya dapat dijalankan berkali-kali. `cocok-enam-sumber-struktur.py` **keluar dengan kode 1**
selama gerbangnya gagal — ia dapat dipasang sebagai penjaga.

---

## 6. Batas berkas ini, dan ia bagian dari isinya

1. **Ia mencocokkan NAMA, bukan ARTI.** Entitas yang namanya cocok di enam tempat dengan **arti**
   berbeda tetap lulus.
2. **"Cocok" bukan "lengkap".** Lihat `L-8` di §4.
3. **Ia tidak membangun apa pun.** Menggambar sebuah tabel bukan membangunnya — `NILAI_SELISIH` dan
   `BESARAN_DAPAT_DISESUAIKAN` kini **ada di gambar** dan tetap **nol berkas** di `ddl-usulan/`.
4. **Satu pertentangan dipilih, bukan diadili** — `S-3`.
