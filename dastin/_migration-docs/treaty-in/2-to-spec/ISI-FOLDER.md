# Isi folder `2-to-spec/` — Treaty In

**Tanggal:** 24 September 2026 · **Fase:** to-spec, lapisan data. Bukan to-ticket.

> **Bentuknya meniru `../../claim-non-prop/2-to-spec/`.** Penyimpangan dari contoh itu dicatat di §3
> beserta sebabnya — bukan dibiarkan terbaca sebagai kelalaian.

---

## 1. Isinya, dan sifat masing-masing

| Berkas | Sifat | Isi |
|---|---|---|
| **`PENELUSURAN-JSON-KE-KOLOM.md`** | **dibangkitkan** | satu baris per jalur `JSONDATA` — 395 dalam lingkup, 272 dikecualikan, menutup ke 667 |
| **`KAMUS-KOLOM.md`** | **dibangkitkan** | satu baris per kolom skema baru beserta asalnya di sistem lama — **35 entitas, 257 kolom** *(24 Sep 2026, sesudah `KTV-A`/`KTV-B`/`KTV-C`, `F-3`, dan `F-17`)* |
| **`ddl-usulan/`** | **dibangkitkan** | **37 berkas: 35 tabel**, `00_SKEMA_DAN_AKUN.sql`, `Z00_KUNCI_ALAMI.sql`. **Nol dijalankan** |
| `TINGKAT-PENCATATAN-14-PAKET-UANG.md` | ditulis tangan | 8 tingkat terbaca dari ekspor, 6 menjadi pertanyaan |
| `PEMETAAN-PROCEDURE.md` | ditulis tangan | lima stored procedure, tiga golongan, dua cacat baru |
| `PERMINTAAN-TEKNIK-TREATY.md` | ditulis tangan | daftar pertanyaan ke orang — sekaligus `daftar wawancara` yang dirujuk enam tempat |
| `INVENTARIS-BERKAS-INDUK.md` | ditulis tangan | 99 berkas induk, 7 biner ditolak dengan namanya, satu yatim |
| `TEMUAN-0-5-ACTUALVALUE-TERBALIK.md` | ditulis tangan | koreksi berkas induk |
| `COCOK-SILANG-CACAH-ATRIBUT.md` | separuh dibangkitkan | cocok-silang cacah atribut seluruh §10.x; §4 memuat adjudikasi kelima calon yang menutup `F-2` |
| `TDA-17-PUTUSAN.md` | ditulis tangan | putusan tiga kolom tanpa rumah — **sebagian dibantah `F-15`** |
| `F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` | ditulis tangan | temuan pembatal: dua dari tiga kolom `TDA-17` **dihitung aturan yang hidup**. Satu kalimat pemilik proses menutupnya |
| `F-13-MV-TANPA-PEMANTAU-KEBASIAN.md` | ditulis tangan | tiga *materialized view* induk menegakkan tanpa pemantau kebasian |
| `SISA-DAN-SERAH-TERIMA.md` | ditulis tangan | gerbang selesai, sisa, dan sepuluh keputusan pemilik proses |
| `ISI-FOLDER.md` | ditulis tangan | berkas ini |

### Yang dibangkitkan TIDAK disunting dengan tangan

Ketiganya lahir dari **satu berkas definisi**, `alat/definisi-skema-treaty-masuk.json`, yang
**diurai dari `SPEC-MODEL-DATA.md` §10** — bukan diketik ulang.

```
alat/buat-kamus-dan-ddl.py      §10  ->  definisi-skema-treaty-masuk.json
alat/tulis-kamus-dan-ddl.py     definisi  ->  KAMUS-KOLOM.md + ddl-usulan/
alat/telusur-jalur-ke-kolom.py  §6 PETA-TELUSUR + definisi  ->  telusur-jalur.json
alat/tulis-penelusuran.py       telusur-jalur.json  ->  PENELUSURAN-JSON-KE-KOLOM.md
```

**Bila salah satu keluaran menyimpang dari §10, ALATNYA yang salah.** Menyunting hasilnya dengan
tangan membuat ketiganya dapat berbeda, dan perbedaannya tidak akan terlihat.

---

## 2. Tiga larangan `ddl-usulan/`, dan bagaimana masing-masing dipenuhi

| Larangan | Keadaan | Diperiksa dengan |
|---|---|---|
| **Tidak ada `CREATE PROCEDURE`** | **0** — juga nol `CREATE FUNCTION` dan nol trigger pembawa aturan bisnis | sapuan atas baris yang **bukan komentar**; komentar di kepala berkas menyebut kata itu dan sempat memberi angka palsu 36 |
| **Tidak ada kolom blob** | **0** `CLOB`, **0** `BLOB` | §2.1 |
| **Tidak ada constraint untuk paket uang yang tingkatnya belum ditentukan** | **6 pernyataan keputusan**, bukan `TODO` | §2.2 |

### 2.1 Lima kolom `CLOB` yang DICABUT — dan kenapa ia cacat perkakas, bukan keputusan

Perkakas pembangkit sempat memaksa lima kolom menjadi `CLOB`: `KETERANGAN`, `CATATAN`,
`PENGECUALIAN`, `KETENTUAN_KHUSUS`, `CATATAN_BORDEREAUX` — seluruhnya pada `VERSI_KONTRAK`.

**§10 menulis tipe `T` untuk kelimanya.** `CLOB` adalah **keputusan perkakas, bukan keputusan
spesifikasi**, dan tidak pernah ditinjau siapa pun. Satu entri dalam peta itu bahkan menunjuk kolom
yang **tidak ada** — `CATATAN_PERSETUJUAN.ISI_CATATAN` — dan itu tidak pernah gagal, karena peta
hanya dibaca bila kolomnya ada. **Penolakan yang diam, di dalam perkakas saya sendiri.**

Pencabutannya mengembalikan kelimanya ke `VARCHAR2(255 CHAR)`, sesuai kelompok tipe `T`.

> **Lubang yang lahir dari pencabutan ini — dilaporkan, tidak ditambal.**
>
> **255 karakter kemungkinan terlalu pendek** untuk `PENGECUALIAN` dan `KETENTUAN_KHUSUS`: keduanya
> menampung klausul pengecualian dan ketentuan khusus kontrak, yang di slip treaty biasa berupa
> beberapa kalimat, kadang beberapa paragraf.
>
> | | |
> |---|---|
> | **Apa** | panjang kolom kelima teks bebas itu **belum diukur** dari data |
> | **Kenapa tidak ditebak** | menaikkannya ke `CLOB` adalah keputusan yang sama yang baru saja dicabut karena tidak pernah ditinjau. Menaikkan ke `VARCHAR2(4000)` juga tebakan |
> | **Akibat bila terlalu pendek** | **pemotongan senyap** pada migrasi — teks klausul hilang sebagiannya, dan tidak ada yang menangkapnya karena barisnya tetap masuk |
> | **Siapa menutup** | **pemilik data** — satu kueri: panjang maksimum kelima ruas itu di `M_TREATY_IN.JSONDATA` |
> | **Yang menagih** | gerbang sesi DDL, bersama angka presisi |
>
> Diusulkan sebagai **Uji AP** di lampiran pengukuran `DAFTAR-ESKALASI-MANAJEMEN.md`.

### 2.2 Enam constraint yang sengaja tidak dipasang

`BATAS_MAKSIMUM_KELOMPOK`, `BATAS_MAKSIMUM_NON_KELOMPOK`, `BATAS_PILIHAN`, `DEDUCTIBLE`,
`LIMIT_AGREGAT`, `NILAI_BATAS` — paket uang **golongan C**, yang di sistem lama **tidak pernah
dicatat bersama mata uangnya**.

Masing-masing ditulis sebagai **pernyataan keputusan** berbentuk **apa · kenapa · akibat · dilihat
di · ditagih**, bukan `TODO`. Sebabnya dinyatakan di berkasnya sendiri: `TODO` akan dipasang orang
berikutnya **tanpa tahu kenapa ia tidak dipasang, dan tanpa dapat mengujinya**.

---

## 3. Penyimpangan dari bentuk `claim-non-prop/2-to-spec/`, beserta sebabnya

| # | Penyimpangan | Sebab |
|---|---|---|
| 1 | **`SPEC-MODEL-DATA.md` tetap di akar modul**, tidak dipindahkan ke sini | ia dirujuk **ratusan kali** dari akar, dan **tidak satu pun rujukan itu dapat diuji**. Memindahkannya memutus seluruhnya sekaligus; membiarkannya dapat dibatalkan kapan saja. Diputuskan pemilik proses 24 Sep 2026 |
| 2 | ada **`PENELUSURAN-JSON-KE-KOLOM.md`**, yang tidak ada di contoh | modul ini menggantikan **`JSONDATA` sebagai tempat kebenaran**. Mengganti JSON menuntut pembuktian tidak ada yang hilang, dan pembuktian itu **per jalur**, bukan per tabel |
| 3 | tidak ada berkas `V*.sql` (view) maupun `Z01`/`Z02` | belum ada konsumen hilir yang **teridentifikasi** (ADR-0051), dan kebijakan pengawasan tulis milik gerbang sesi DDL |
| 4 | penomoran berkas DDL **dibangkitkan menurut kedalaman induk**, bukan menurut nomor tiket | tidak ada tiket di alur ini; urutan jalannya yang mengikat |

---

## 4. Yang TIDAK ada di sini, dan itu disengaja

| Yang tidak ada | Sebab |
|---|---|
| tabel `NILAI_SELISIH` dan `BESARAN_DAPAT_DISESUAIKAN` | milik modul Adjustment — `SPEC-MODEL-DATA.md` §11.3 |
| tabel `RETRO_KELUAR` dan `PENCAPAIAN` | GEL-2 / GEL-3, di luar gelombang ini |
| angka presisi yang diputuskan | gerbang sesi DDL; yang dipakai **mewarisi** kelompok ADR-0003 modul tetangga sebagai usulan |
| bentuk fisik induk polimorfik (`POTONGAN`, `PENYEBARAN`) | gerbang sesi DDL — `INV-17` melarang rujukan yang sasarannya bergantung nilai kolom lain |
| spesifikasi layar | `L-4` — milik batch lapisan aplikasi |
