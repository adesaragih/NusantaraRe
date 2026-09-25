> Modul  : Treaty In Adjustment · ronde C
> Dibuat : 2026-09-24
> Sifat  : buku besar ronde — apa yang dijalankan, apa yang ditolak, apa yang salah, dan apa yang
>          belum diperiksa. **Termasuk yang hasilnya bersih.**

# Audit ronde C

## 1. Buku besar — apa yang dijalankan dan apa yang dikembalikannya

Seluruh pembacaan ekspor lewat perkakas. **Tidak satu berkas XML pun dibaca langsung** — kedua
ekspor berjumlah 185 MB.

| Perkakas | Sasaran | Hasil | Ditolak / titik buta |
|---|---|---|---|
| `tulis.py sapu` | `.IsProRate` | 2 penulisan, keduanya BIND `pxCheckbox` **dapat disunting**, dua ekspor | 0 salinan terbungkus · 0 gagal urai · **32 Java, 191 SQL/REST tidak terurai** |
| `tulis.py sapu` | `.EDMEffective` | 5 penulisan; sesudah `Param.*` dipisah → **1 penulis** `TreatyIn.EDMEffective` | idem |
| `tulis.py sapu` | `.ProRatePercent` | 6 DT + 4 BIND; seluruh BIND **hanya-baca** | idem |
| `pre.py` | `TreatyEDMProRateCalculation` | blok mati: 2.1, 2.8, 2.9, 2.10, 8 (+ anak) | — |
| `pre.py` | `TreatyEDMDifferenceDeduction` | langkah 4 dan 5 **MATI** | — |
| `dt.py` | `TreatyCalculateProratePct` | 4 langkah + 7 anak; gerbang `IsProRate == true`; bawaan 0 | — |
| `cat56.py` | katalog khas | **56 aturan**; hanya `TreatyInSetEditPre` di antara yang disentuh | — |
| pencarian teks | 4 nama aturan, dua ekspor | pemanggil `TreatyCalculateProratePct`: **nol Activity** | pencarian teks **kebal nama tag** — itu kekuatannya di sini |

> **Titik buta yang tidak boleh dilupakan:** 32 langkah Java dan 191 langkah SQL/REST tidak terurai
> oleh penyapu mana pun. Setiap angka "nol penulis" di atas berarti **nol di antara lima bentuk yang
> terurai**, bukan nol di sistem.

## 2. Kesalahan yang dinyatakan terang

### `PC-1` — saya hampir melaporkan lima penulis `EDMEffective`

Sapuan mengembalikan **5 penulisan** untuk `.EDMEffective`. Empat di antaranya bukan penulis
`TreatyIn.EDMEffective`: dua menyasar **`Param.EDMEffective`** (penyapu mencocokkan **akhiran**
jalur), dua lagi kontrol layar **hanya-baca**.

Bila angka 5 dilaporkan apa adanya, kesimpulan NC-01 **terbalik** — "banyak penulis" alih-alih
"satu penulis yang memakukannya pada `Commencement`". Ditangkap dengan membaca kolom sasaran satu
per satu, bukan oleh perkakas.

**Diusulkan sebagai catatan pemakaian `tulis.py`:** cetak `Param.*` dan `TreatyIn.*` sebagai dua
kelompok terpisah, karena keduanya menjawab pertanyaan yang berbeda.

### `PC-2` — kolom sisi hampir ditentukan dari nama berkas

Saya sempat menggolongkan seluruh `TreatyEDM*` sebagai 56-KHAS karena awalannya. `cat56.py`
membantahnya: **tidak satu pun** ada di katalog 56. Aturan *"kolom sisi diperiksa per berkas lewat
`pzInsKey`, bukan dari nama berkas"* menangkapnya — tetapi hanya karena aturannya dijalankan, bukan
karena saya curiga.

**Akibatnya besar:** lima dari tujuh temuan ronde ini berpindah sisi, dari "temuan Adjustment"
menjadi **"temuan Treaty In yang berjalan hari ini"**.

## 3. Daftar untuk dibantah — siap kirim

Dua butir baru, `DB-18` dan `DB-19`, bunyinya di [`03-LUBANG.md`](03-LUBANG.md) §2. Keduanya
bergabung dengan `DB-1`…`DB-17` di `GRILL-A/07-AUDIT.md` §3 dan **dikirim sebagai satu paket**.

Tiga kueri baru, `UA-19`…`UA-21`, bunyinya di `03-LUBANG.md` §1. Dua di antaranya menyatakan batas
pembuktiannya di dalam bunyinya sendiri.

## 4. Titik periksa — dua arah, dan hasil bersih ikut dicatat

### 4a. Arah pertama: adakah TDA tanpa nasib? → **TIDAK BERSIH, tujuh ditemukan**

Dicocokkan: tabel **Lacak TDA** terhadap tabel **Ronde** di `../KEPUTUSAN-GRILLING-ADJUSTMENT.md`.
Hasilnya `03-LUBANG.md` §4 — tujuh TDA menunjuk cabang yang sudah ditutup, ditambah empat yang
tertutup ronde ini dan perlu ditandai.

### 4b. Arah kedua: adakah GRL yang tidak punya baris di Lacak TDA maupun di indeks? → **BERSIH**

Keempat putusan ronde ini punya baris di indeks sesudah pembaruan, dan masing-masing menyebut
pembatalnya. **Hasil bersih dicatat supaya arah ini tidak dijalankan lagi tanpa sebab.**

### 4c. Arah ketiga: adakah klaim ronde ini yang bersandar pada berkas penggrill? → **BERSIH**

`PENGETAHUAN-PENGGRILL-ADJUSTMENT.md` menyatakan dirinya **bukan sumber kebenaran**. Ia dipakai
hanya sebagai peta — untuk mengetahui butir mana yang terbuka. Setiap fakta di `01-TEMUAN.md`
disapu ulang atas ekspor atau dikutip dari `PENGETAHUAN.md`.

Dua klaimnya justru **disidangkan dan salah satunya salah kaprah** — `SC-1` dan `SC-2` di
[`02-SIDANG.md`](02-SIDANG.md).

## 5. Yang BELUM diperiksa, dan dinyatakan begitu

| Hal | Kenapa tidak diperiksa | Siapa dapat menutupnya |
|---|---|---|
| nilai `ProRatePercent` yang tersisa di halaman saat pengajuan | **tidak terbaca dari ekspor** — ia keadaan runtime, bukan aturan | `UA-20` |
| apakah 32 langkah Java menulis properti bisnis | teksnya bebas; sapuan teks sebelumnya mengembalikan nol untuk properti bisnis, dan itu **bukan bukti** | sapuan tersendiri, bila pernah dibutuhkan |
| apakah 191 langkah SQL/REST menulis ke `TreatyIn.*` | menulis ke basis data, bukan ke clipboard — dicirikan di sesi sebelumnya, tidak diulang | — |
| enam jalur berbeda versi antara kedua ekspor | di luar lingkup ronde ini; tidak satu pun menjadi bukti TDA (catatan 24 Sep) | ronde D bila menyentuh TDA |

## 6. Aturan berhenti — tidak satu pun terpicu

Ketiga keadaan di `00-LINGKUP.md` §5 diperiksa di setiap putusan:

1. **membatalkan GRL terkunci** — tidak. `GRL-14` **menyusutkan** kalimat GRL-13 (*"penyesuaian
   premi selalu material"* → *"…yang mengubah premi aktual selalu material"*), dan penyusutan itu
   **dinyatakan terang** di `06-PUTUSAN.md` beserta `DB-19`. Menyusut bukan membatalkan.
2. **wewenang pemilik proses di tengah blok** — tidak. Keempat butir diputuskan pemilik proses pada
   24 September dengan *"ikuti rekomendasi"*.
3. **fakta yang tidak terbaca dan tidak ada perkakasnya** — satu ditemukan (nilai `ProRatePercent`
   saat pengajuan), dan ia **tidak menghentikan apa pun**: ia menjadi `UA-20`, dan `GRL-15` tetap
   diputuskan karena putusannya tidak bergantung padanya.

---

## MA-11 — bekerja dari ringkasan yang basi, tanpa memeriksa keadaan sekarang

**Tanggal:** 24 September 2026, sesudah sesi to-spec Treaty In.

### Apa yang terjadi

Sesi terputus dan dilanjutkan dari ringkasan. Ringkasan itu menyatakan cabang C, E, dan I **masih
terbuka** dengan sisa anggaran tiga butir — keadaan yang benar **pada saat ringkasan dibuat**, dan
sudah **tidak benar** ketika saya melanjutkannya.

Saya **mengajukan ulang C3 yang sudah terkunci**, meminta jawaban pemilik proses atasnya, lalu
**menulis `GRL-14` kedua** ke `GRILL-B/06-PUTUSAN.md` — folder cabang yang bahkan bukan tempatnya,
karena C3 milik ronde C.

Yang menemukannya: saya membuka `KEPUTUSAN-GRILLING-ADJUSTMENT.md` untuk memperbarui indeksnya, dan
indeks itu **sudah** memuat `GRL-14` sampai `GRL-17` beserta folder `GRILL-C/` dan `GRILL-TDA/`.

### Kenapa ia lolos

**Indeks itu sendiri sudah memperingatkannya**, dan saya tidak membacanya lebih dulu:

> **Berkas ini indeks, bukan sumber kebenaran.** … **Bila sesi terputus, baca indeks ini lebih
> dulu**, lalu berkas ronde yang statusnya masih TERBUKA.

Peringatan itu ditulis untuk keadaan ini persis. Yang gagal bukan ketiadaan penjaga — penjaganya ada,
di tempat yang benar, dan saya melewatinya.

### Perbaikannya

| Yang dilakukan | |
|---|---|
| duplikat `GRL-14` **dicabut** dari `GRILL-B/06-PUTUSAN.md` | GRILL-B kembali memuat GRL-09 … GRL-13 |
| bukti baru sesi to-spec **dipindahkan** ke `GRL-14` yang asli | penguatan bertanggal di `GRILL-C/06-PUTUSAN.md` — kalibrasi kode tindakan 66/66, dan koreksi `7-2-ACTUALVALUE.md` |
| putusannya **tidak berubah** | GRL-14 asli sudah memutuskan hal yang sama: `ActualValue` dibuang, premi aktual menjadi nilai versi. Yang saya tambahkan hanya **sebab mekanisnya** |

### Aturan yang diusulkan — `TA-09`

> **Sesi yang dilanjutkan dari ringkasan membaca INDEKS lebih dulu, dan memeriksa daftar folder
> ronde, sebelum mengajukan pertanyaan apa pun.**
>
> Ringkasan adalah potret pada satu saat; indeks adalah keadaan sekarang. Bila keduanya berbeda,
> **indeks yang menang** — dan perbedaannya sendiri adalah temuan, karena ia menunjukkan pekerjaan
> yang sudah selesai tetapi tidak terbawa.
>
> Uji termurahnya satu perintah: **daftar folder `GRILL-*/`**, dan bandingkan dengan yang ringkasan
> sebut. Folder yang ada tetapi tidak disebut ringkasan adalah pekerjaan yang hilang dari ingatan.

**Ongkos yang sudah terjadi:** satu pertanyaan diajukan kepada pemilik proses yang tidak perlu
dijawab, dan satu putusan ganda yang harus dicabut. **Tidak ada keputusan yang berubah, dan tidak
ada bukti yang hilang.**
