> Modul  : Treaty In Adjustment · ronde TDA
> Dibuat : 2026-09-24
> Sifat  : buku besar ronde — apa yang dijalankan, apa yang salah, dan apa yang belum diperiksa.
>          **Termasuk yang hasilnya bersih, dan termasuk kekeliruan saya sendiri.**

# Audit ronde TDA

## 1. Buku besar

| Yang dijalankan | Atas apa | Hasil |
|---|---|---|
| pembacaan berkas ronde | `GRILL-A/06-PUTUSAN`, `GRILL-B/06-PUTUSAN`, `GRILL-C/06-PUTUSAN` | tujuh TDA diadili; **enam tertutup dengan kutipan beserta nomor barisnya** |
| pembacaan ADR induk | 0034, 0037, 0044, 0045, 0055 | tiga penutupan bersandar padanya |
| `scratchpad/props.py` | 13 berkas aturan, dua ekspor | nama properti asli; **0 berkas tidak ditemukan** |
| `perbarui-indeks-tda.py` | indeks | **16 penggantian, 0 gagal** |

**Tidak ada sapuan ekspor baru untuk mencari perilaku.** Ronde ini menyapu catatan, bukan sistem
lama — dinyatakan begitu supaya tidak ada yang mengira ekspor diperiksa ulang.

## 2. Kesalahan yang dinyatakan terang

### `PT-1` — perkiraan tiga pertanyaan, hasilnya satu

`GRILL-C/03-LUBANG.md` §4 memperkirakan **tiga** pertanyaan sungguhan: `TDA-07`, `TDA-11`,
`TDA-13`. Hasilnya **satu**.

`TDA-07` dan `TDA-11` sudah tertutup oleh `GRL-08` dan `GRL-10`/`B-4` — dan **keduanya menyebut
nomor TDA-nya di dalam kolom "dari apa" putusannya sendiri**. Perkiraan itu dibuat **sebelum**
berkas rondenya dibuka.

Bentuk `TA-08`: *baca temuan bernomor sendiri sebelum membangun rantai sebab.* Ongkosnya di sini
nol — perkiraan yang terlalu tinggi hanya membuat ronde terlihat lebih berat. **Ongkos arah
sebaliknya tidak nol**, dan itu yang perlu diingat: perkiraan yang terlalu rendah akan membuat
sebuah pertanyaan tidak pernah diajukan.

### `PT-2` — nama ronde bertabrakan dengan nama cabang

Folder semula `GRILL-D/`, mengikuti konvensi *"ronde = cabang"*. Indeks sudah memuat **cabang D**
yang ditutup tanpa ronde. Tertangkap saat indeks hendak diperbarui, **sebelum** ada rujukan dari
luar; diubah menjadi `GRILL-TDA/`. Rinciannya `NT-02`.

### `PT-3` — "lima belas TDA" di berkas saya sendiri, beberapa menit sesudah menuliskan aturannya

Draf pertama `05-GERBANG.md` §2 berbunyi *"Lima belas TDA, lima belas nasib"*. Angka itu benar
sampai ronde C; **`TDA-16` lahir 27 September** dan penyebutnya tidak ikut diperbarui.

> Saya menulis `NT-01` — *"sebuah keadaan diperbarui di satu tabel dan tidak di tabel sebelahnya"* —
> lalu melakukannya sendiri di berkas berikutnya. Dicatat apa adanya, sebab ia bukti terbaik bahwa
> `TA-09` menuntut **langkah**, bukan kewaspadaan.

Dikoreksi di tempatnya, dengan koreksinya ikut tertulis di sana.

## 3. Titik periksa — tiga arah, dan hasil bersih ikut dicatat

### 3a. Adakah TDA yang masih tanpa nasib? → **BERSIH**

Enam belas TDA, enam belas nasib. Tiga di antaranya **"diperbaiki dengan sisa ditunda"**, dan
sisanya dicatat sebagai butir bernomor di `03-LUBANG.md` §2 — bukan dibiarkan melekat pada TDA yang
sudah dinyatakan tertutup.

### 3b. Adakah cabang yang masih menunggu tanpa ronde? → **BERSIH sesudah dikoreksi**

Enam baris Status cabang yang basi (`NT-01`) diperbaiki. Tinggal **cabang K** yang terbuka, dan
penahannya bukan pekerjaan.

### 3c. Adakah penutupan tanpa kutipan? → **BERSIH**

Aturan `00-LINGKUP.md` §3 melarang menyatakan sebuah TDA tertutup hanya karena cabangnya ditutup.
Ketujuh penutupan membawa kutipan: lima dari berkas ronde, dua dari ADR induk, dan satu (`TDA-13`)
dari putusan ronde ini sendiri.

## 4. Yang BELUM diperiksa, dan dinyatakan begitu

| Hal | Kenapa | Siapa dapat menutupnya |
|---|---|---|
| apakah nilai `TreatyIn.EDMState` di produksi hanya memuat lima kombinasi sah | **pertanyaan data**, bukan pertanyaan aturan | `UA-2` |
| berapa banyak baris warisan yang jenisnya pernah diubah sesudah pengajuan | **tidak dapat dijawab sama sekali** — sistem lama tidak menyimpan jejaknya | tidak ada; dinyatakan sebagai **tidak dapat diketahui** |
| `CABANG-K-PEMETAAN-TO-SPEC.md` belum memuat `GRL-12`…`GRL-18` | lahir sesudah berkas itu ditulis | **sesi to-spec**, sebagai pekerjaan pertamanya |

Baris kedua layak dibaca dua kali: ia bukan uji data yang belum dijalankan, melainkan pertanyaan
yang **tidak punya sumber**. Menuliskannya sebagai `UA` akan menyesatkan.

## 5. Aturan berhenti — tidak satu pun terpicu

1. **membatalkan GRL terkunci** — tidak. `GRL-18` **melengkapi** `GRL-13` pada sumbu yang GRL-13
   memang tidak sentuh.
2. **wewenang pemilik proses di tengah blok** — tidak. `ST-7` diputuskan pemilik proses pada 24
   September.
3. **fakta yang tidak terbaca dan tidak ada perkakasnya** — satu ditemukan (§4 baris kedua), dan ia
   tidak menghentikan apa pun: `GRL-18` mengatur sistem baru, dan tidak bergantung pada berapa
   sering sistem lama disalahgunakan.
