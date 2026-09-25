# SEAM-ADJUSTMENT — apa yang Treaty In sediakan untuk Treaty In Adjustment

**Tanggal:** 23 September 2026
**Berlaku untuk:** Treaty In
**Sumber mengikat:** ADR-0048 (arahan pemilik proses), ADR-0036, ADR-0039, ADR-0040, ADR-0045
**Ditulis untuk:** perancang modul Treaty In Adjustment dan perancang basis data Treaty In

---

> **CATATAN, 23 September 2026 — embargo dicabut; keputusan berkas ini tetap berlaku.**
> Modul Treaty In Adjustment sudah dibedah; hasilnya di `../treaty-in-adjustment/PENGETAHUAN.md`.
> Satu **fakta** berubah, bukan satu keputusan: sistem lama ternyata **memang menyimpan potret
> nilai lama** — `TreatyIn.OLDDATA`, salinan penuh kontrak yang dibekukan ke dalam `JSONDATA`
> addendum saat addendum dibuat (`Activity/TreatyInSetAddendumToHistory.xml`).
> Kalimat "DATA LAMA TIDAK PERNAH DISALIN" di bawah adalah **aturan untuk sistem baru** dan tetap
> berdiri; ia hanya tidak boleh dibaca sebagai pemerian sistem lama.
> Akibat praktisnya satu: migrasi **tidak perlu merekonstruksi** nilai lama, dan potret yang ada
> dapat dipakai untuk **memeriksa** hasil migrasi.


## Batas berkas ini

Modul **Treaty In Adjustment berada di bawah embargo** dan tidak dibedah. Berkas ini **tidak**
merancang tabel selisih, layar, alur kerja, maupun aturan bisnis Adjustment. Ia hanya menetapkan
**apa yang Treaty In wajib sediakan** agar modul itu bisa dirancang nanti tanpa Treaty In harus
diubah lagi.

Tiga hal. Tidak lebih.

Tidak ada DDL di berkas ini. Tipe dinyatakan **secara konseptual**; bentuk fisiknya ditetapkan di
sesi DDL berikutnya.

---

## Aturan yang ditulis ke dalam seam, dan berlaku atas segalanya di bawah

> ### DATA LAMA TIDAK PERNAH DISALIN — IA DI-SELECT.

Tidak ada penyalinan nilai kontrak ke dalam catatan Adjustment. Tidak ada tabel bayangan, tidak ada
kolom `OLDDATA`, tidak ada potret yang dibuat saat Adjustment dibuat. Adjustment menyimpan
**penunjuk**, dan membaca nilainya lewat bentuk baca di bagian 3.

**Bentuk model membuat aturan ini gratis, bukan mahal.** `KONTRAK` hanya memuat lapisan identitas
yang dibekukan; **seluruh nilai kontrak tinggal di `VERSI_KONTRAK` dan anak-anaknya**. Karena itu
"nilai lama" adalah **versi sebelumnya dari kontrak yang sama** — sebuah SELECT biasa atas baris
yang memang sudah ada. Tidak ada mekanisme penyalinan yang perlu dibangun, karena tidak ada yang
perlu disalin.

Ini penerapan ADR-0023: satu bentuk kanonik, hilir diberi rujukan bukan salinan.

---

## 1. Identitas versi kontrak yang stabil dan bisa ditunjuk dari luar

### Yang disediakan

| Hal | Sifat |
|---|---|
| `ID_VERSI_KONTRAK` | pengenal **buatan sistem**, tanpa makna, unik seumur hidup sistem |
| `ID_KONTRAK` | pengenal kontrak yang memiliki versi itu |
| `NOMOR_URUT_VERSI` | urutan versi di dalam kontraknya, bilangan bulat mulai dari 1 |

**`ID_VERSI_KONTRAK` adalah satu-satunya hal yang boleh disimpan konteks lain sebagai penunjuk.**
Bukan nomor offer, bukan gabungan nomor offer dengan urutan versi, bukan tanggal.

### Empat sifat yang dijaminkan, dan kenapa masing-masing ada

1. **Tidak terbentuk dari teks yang bisa disunting.** Nomor offer bisa diperbaiki orang; identitas
   yang terbentuk darinya ikut berubah dan seluruh penunjuk dari luar putus. (Pola terlarang:
   identitas terbentuk dari penyuntingan teks.)

2. **Tidak terbentuk dari waktu.** Dua versi yang lahir pada milidetik yang sama bukan kemustahilan
   yang boleh diandalkan. (Pola terlarang: kunci dari cap waktu.)

3. **Tidak pernah dipakai ulang.** Versi yang dibatalkan tidak mengembalikan nomornya. Penunjuk yang
   menggantung harus tetap menggantung dan terlihat menggantung — bukan diam-diam menunjuk benda
   lain.

4. **Bertahan melewati migrasi.** Baris warisan mendapat `ID_VERSI_KONTRAK` yang **dicetak satu
   kali** saat pemindahan dan tidak pernah dicetak ulang. Identitas lamanya (`ID` pada
   `M_TREATY_IN`, dan pasangan `ID`/`OLDID` pada `M_TREATY_IN_EDM`) disimpan **berdampingan**
   sebagai atribut tersendiri, bukan sebagai pengganti — ADR-0042.

### Apa yang TIDAK disediakan, dan itu disengaja

Tidak ada pengenal yang menggabungkan kontrak dan versi menjadi satu teks yang bisa dibaca manusia.
Kalau Adjustment butuh label yang bisa dibaca orang untuk ditampilkan, ia **membacanya** lewat
bagian 3, dan tidak menyimpannya.

---

## 2. Jaminan bahwa versi yang sudah disetujui tidak pernah berubah nilainya

### Isi jaminannya

Sebuah versi kontrak yang berkeadaan `DISETUJUI` (ADR-0055) **tidak dapat berubah nilainya, oleh
siapa pun, lewat jalur mana pun**. Perubahan menghasilkan **versi baru**, tidak pernah penimpaan.

### Empat hal yang menegakkannya

1. **Mesin keadaan tidak punya perpindahan keluar dari `DISETUJUI`.** ADR-0055 menetapkan daftar
   perpindahan yang sah dan `DISETUJUI` tidak muncul sebagai asal pada satu pun. Keempat aturan
   sistem lama yang membatalkannya dari samping tidak dibawa.

2. **Penulisan ke versi yang sudah disetujui ditolak, bukan diabaikan.** Penolakan adalah kegagalan
   yang dilaporkan, bukan operasi yang diam-diam tidak berbuat apa-apa — ADR-0035. Yang diam-diam
   gagal akan tampak seperti berhasil di mata penulisnya.

3. **Setiap penulisan meninggalkan jejak.** ADR-0045 menjadikan siapa dan kapan sebagai fakta mesin
   atas setiap perubahan. Larangan yang tidak bisa diperiksa bukan jaminan; jejak inilah yang
   membuat pelanggaran **terdeteksi**, bukan sekadar terlarang.

4. **Nilai yang bisa bergerak dicatat PADA versi, bukan dirujuk dari luarnya.** Ini konsekuensi yang
   paling mudah terlewat, dan tanpanya ketiga hal di atas tidak cukup.

### Konsekuensi (4), diuraikan — karena di sinilah jaminan ini paling mudah bocor

Membekukan baris versi tidak membekukan apa pun yang baris itu **rujuk**. Empat hal di sistem lama
bergerak di luar kontrak dan ikut mengubah angkanya:

| Yang bergerak di sistem lama | Akibatnya bila hanya dirujuk |
|---|---|
| kurs tahunan, dipilih menurut **tahun berjalan** | nilai rupiah versi lama berubah tiap ganti tahun |
| bagian NuRe (`RNMShare`) | premi bagian NuRe ikut bergeser |
| susunan kapasitas pada tabel acuan luar | penyebaran tersusun ulang dengan angka berbeda |
| nama dan atribut master cedant / mata uang | label berubah; nilai tidak, tetapi laporan lama jadi tak cocok |

Karena itu: **kurs yang dipakai, tanggal dan sumber kurs, bagian NuRe yang dipakai, dan susunan
penyebaran yang dipakai dicatat sebagai NILAI pada versi**, bukan sebagai rujukan yang dibaca ulang
saat diminta. Yang dirujuk hanyalah hal yang memang tidak boleh beku — misalnya nama cedant untuk
ditampilkan.

Ini bukan pelanggaran ADR-0037 (turunan tidak disimpan). Kurs yang dipakai pada sebuah transaksi
adalah **fakta yang dibukukan**, bukan turunan: ia tidak bisa dihitung ulang dari masukan lain
tanpa kehilangan kebenarannya. ADR-0036 menyebutnya beku saat disetujui.

### Apa yang membatalkan jaminan ini

Bila ditemukan satu saja jalur penulisan ke versi `DISETUJUI` yang tidak melewati mesin keadaan —
misalnya pekerjaan latar, prosedur basis data, atau perbaikan data manual — jaminan ini **batal
seluruhnya**, dan Adjustment tidak bisa dipertanggungjawabkan. Karena itu jalur semacam itu harus
dilarang secara eksplisit di rancangan basis data, bukan diandaikan tidak ada.

---

## 3. Bentuk baca nilai kontrak pada versi tertentu

### Namanya dan sifatnya

**`NILAI_VERSI_KONTRAK`** — bentuk **baca**, diturunkan dari entitas di belakangnya, **bukan
kanonik**, tidak ditulis oleh siapa pun, dan dapat dibangun ulang kapan saja dari entitas itu.

Butirannya: **satu baris per (versi kontrak × besaran uang)**.

### Isinya

| Kelompok | Isi | Kenapa ada |
|---|---|---|
| **Versi** | `ID_VERSI_KONTRAK`, `ID_KONTRAK`, `NOMOR_URUT_VERSI` | penunjuk yang dipegang Adjustment |
| | `KEADAAN_SIKLUS_HIDUP`, `TANGGAL_DISETUJUI` | agar pembaca bisa **menolak** memakai versi yang belum disetujui, alih-alih memakainya tanpa sadar |
| **Letak besaran** | `JENIS_INDUK`, `KUNCI_PADANAN` | menjawab "besaran ini milik baris yang mana" |
| **Besaran** | `NAMA_BESARAN` | besaran apa: limit, deposit premium, EPI, dan seterusnya |
| **Paket uang** (ADR-0007, 0029, 0039) | `NILAI`, `KODE_MATA_UANG` | nilai dan mata uangnya |
| | `TINGKAT_PENCATATAN` | 100% treaty atau bagian NuRe — **sifat besarannya**, bukan konvensi |
| | `PERSEN_BAGIAN_DIPAKAI` | bagian yang dipakai, bila besaran itu besaran tingkat bagian |
| | `KURS`, `TANGGAL_KURS`, `SUMBER_KURS` | kurs yang dipakai saat itu, dicatat bukan dicari ulang |
| | `NILAI_IDR` | nilai rupiah yang dibukukan |

Tipe konseptual: `NILAI`, `KURS`, `PERSEN_BAGIAN_DIPAKAI`, `NILAI_IDR` adalah **bilangan eksak**
(tidak pernah teks, tidak pernah pecahan biner). `KODE_MATA_UANG`, `TINGKAT_PENCATATAN`,
`JENIS_INDUK`, `NAMA_BESARAN`, `KEADAAN_SIKLUS_HIDUP` adalah **himpunan nilai tertutup**.

### `KUNCI_PADANAN` — bagian yang paling mudah dilupakan, dan tanpanya seam ini tidak berguna

Selisih menuntut **pemadanan baris**: baris mana pada versi lama yang sepadan dengan baris mana pada
versi baru. `ID_VERSI_KONTRAK` tidak bisa dipakai untuk itu — ia justru berbeda antar versi. Begitu
pula pengenal buatan sistem milik baris anak.

Maka setiap baris membawa **identitas bisnisnya yang stabil lintas versi**: nomor layer bagi layer,
pengenal pihak bagi bagian dan penyebaran, dan seterusnya. Itulah `KUNCI_PADANAN`.

> ### Ditambahkan 24 September 2026 — untuk sebagian entitas, `KUNCI_PADANAN` harus MAJEMUK
>
> Kalimat di atas mengandaikan setiap kunci alami unik **di dalam versinya**. **Sembilan memang
> begitu; tiga tidak** — `SPEC-INVARIAN.md` §4.1 memisahkannya menjadi dua bentuk:
>
> | Bentuk | Lingkup kunci alaminya | `KUNCI_PADANAN`-nya |
> |---|---|---|
> | **A** — sembilan entitas | **versi** | kunci alaminya sendiri, **cukup** |
> | **B** — INV-06 `DETAIL_PROPORSIONAL`, INV-15 `POTONGAN`, INV-16 `PENYEBARAN` | **induk langsung** | **kunci padanan INDUKNYA lebih dulu, baru kuncinya sendiri** |
>
> Contohnya: kelompok treaty **tidak** unik di dalam satu versi — ia unik di dalam satu `LAYER`, dan
> kontrak berlayer banyak memang wajar memuat kelompok treaty yang sama di dua layer berbeda.
> Memadankan `DETAIL_PROPORSIONAL` antar versi hanya dengan kelompok treaty karena itu **memadankan
> baris yang salah**, dan selisihnya akan dihitung terhadap pasangan yang keliru.
>
> **Bentuk yang benar untuk B:** `‹KUNCI_PADANAN induk›` + `‹kunci alami sendiri›` — untuk
> `DETAIL_PROPORSIONAL` berarti *(nomor layer + bagian layer)* + *kelompok treaty*.
>
> **`POTONGAN` menuntut satu ruas lagi**, karena §14.1 menetapkan ia **berinduk dua**: kunci
> padanannya harus menyebut **induk mana** yang dimaksud, selain kunci padanan induk itu dan jenis
> potongannya. Tanpa ruas itu, potongan pada kedua induk saling tertukar saat dipadankan.

Tanpa ini, Adjustment hanya bisa membandingkan total dan tidak bisa menunjukkan **apa** yang
berubah. Dengan ini, selisih bisa dihitung per baris, dan baris yang **hilang** atau **muncul**
antar versi ikut terlihat — bukan hanya baris yang nilainya bergeser.

### Aturan bacanya

1. **Dibaca menurut versi, bukan menurut waktu.** Diberi `ID_VERSI_KONTRAK`, bentuk ini
   mengembalikan **persis apa yang tercatat pada versi itu**. Ia tidak pernah membaca data acuan
   yang berlaku hari ini. Inilah yang menyambungkannya ke jaminan bagian 2.

2. **Besaran yang tidak ada tidak menjadi nol.** Bila sebuah besaran tidak tercatat pada versi itu,
   barisnya **tidak ada**. Ia tidak muncul dengan nilai nol, dan tidak muncul dengan angka
   pengganti. ADR-0035.

3. **Nilai rupiah yang tidak bisa dihitung tidak menjadi nol.** Bila kurs tidak tercatat,
   `NILAI_IDR` **tidak terisi** dan barisnya menyatakan dirinya tidak lengkap. Di sistem lama,
   kurs yang tidak ditemukan membuat nilainya menjadi nol dan lenyap dari total — kesalahan yang
   menghilangkan, bukan menggelembungkan, dan karena itu tidak pernah menimbulkan pertanyaan.

4. **Baca-saja.** Tidak ada jalur tulis lewat bentuk ini, sekarang maupun nanti.

### Satu keputusan saya, ditandai

`JENIS_INDUK` + `KUNCI_PADANAN` adalah **rujukan bergaya polimorfik**: satu pasang kolom yang
menunjuk ke beberapa jenis induk berbeda. Bentuk itu **dilarang di model tersimpan** (pola
terlarang: kolom yang menampung dua arti menurut cabang), dan saya memakainya di sini secara sadar
dengan tiga syarat:

- ia hanya ada di **bentuk baca**, yang tidak kanonik dan tidak ditulis siapa pun;
- kanoniknya tetap entitas di belakangnya, yang tetap bertipe tegas dan berkunci asing sungguhan;
- ia dapat dibangun ulang, sehingga kesalahan di dalamnya tidak pernah menjadi kehilangan data.

Alternatifnya adalah satu bentuk baca terpisah per jenis induk. Saya tidak memilihnya karena seam
ini menyeberangi batas modul yang **sedang diembargo**: menambah bentuk baca baru nanti berarti
mengubah Treaty In setelah Adjustment dirancang, dan itulah yang seam ini ada untuk dicegah.
Membandingkan selisih pada dasarnya operasi yang seragam, jadi satu bentuk seragam melayaninya.

**Apa yang membatalkan keputusan ini:** bila ternyata Adjustment memperlakukan tiap jenis induk
dengan aturan yang berbeda-beda — bukan sekadar membandingkan angkanya — maka keseragaman itu
palsu, dan bentuk baca dipecah per jenis induk.

---

## Apa yang Adjustment simpan, dan apa yang tidak

Dinyatakan di sini hanya sejauh ia mengikat Treaty In. Selebihnya bukan urusan sesi ini.

| Adjustment **menyimpan** | Adjustment **tidak menyimpan** |
|---|---|
| penunjuk ke versi lama dan versi baru (`ID_VERSI_KONTRAK` keduanya) | nilai kontrak dalam bentuk apa pun |
| baris selisih, di tabel tersendiri | salinan baris kontrak |
| kapan selisih itu dihitung | apa pun yang bisa dibaca lewat bagian 3 |

Karena nilai lama selalu dapat di-SELECT dan nilai baru diketahui, **selisih yang tersimpan dapat
direkonsiliasi ulang kapan saja**. Ketidakcocokan antara selisih tersimpan dan selisih terhitung
bukan masalah — **ia alat deteksi**, dan kemungkinan itu harus dipertahankan, bukan dihilangkan.

---

## Yang tetap tidak dikerjakan

- Tabel selisih, layar, dan alur kerja Adjustment: **tidak dirancang**.
- Folder `D:\XML_NURE\Treaty In Adjustment`: **tetap diembargo dan tidak dibuka**.
- DDL: **tidak ada di berkas ini**, sesuai batas sesi.
