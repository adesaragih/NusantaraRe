# Keputusan pembagian tiket — Treaty In

**Tanggal:** 23 September 2026 · **direvisi 24 September 2026** (butir 4 urutan wewenang; U-6; ongkos kedua dan PEMBUAT PERTAMA;
mode peringatan 4 butir; §4.1–§4.4 jalan (iii); perumusan pengecualian setingkat modul; L-6, L-7)
**Keadaan:** hasil gerbang. Ditulis **sebelum satu tiket pun ada.**

Berkas ini menjawab tiga pertanyaan yang harus dijawab lebih dulu: **sebuah tiket itu sebesar apa**,
**apa yang masuk penyerahan pertama**, dan **tiket yang penghalangnya belum terjawab diapakan**.
Ketiganya dipakai seterusnya tanpa kecuali. Bila salah satunya perlu diubah, ia diubah **di sini**
dan seluruh tiket disesuaikan — bukan dilanggar diam-diam pada satu-dua tiket.

---

## 0. Dua hal yang harus dibaca sebelum apa pun, oleh siapa pun yang tidak ikut sesi sebelumnya

### 0.1 Huruf `G` di proyek ini sudah punya tiga arti berbeda

Ini bukan kerewelan. Salah baca di sini menghasilkan ruang lingkup yang salah.

| Tulisan | Artinya | Di berkas mana |
|---|---|---|
| **GEL-2**, **GEL-3** | **gelombang** penyerahan — pekerjaan yang sengaja ditunda ke penyerahan berikutnya | `SPEC-MODEL-DATA.md` §2.3 menulisnya `G2`/`G3` |
| **G1, G1b, G1c, G2, G3, G4** | gerbang **sambungan ke modul Adjustment**, dijawab sesi ERD | `4-erd-dan-tabel-datar/KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` |
| **G1, G2, G3** *(berkas ini)* | gerbang **pembagian tiket** | berkas ini saja |

> **Mulai sekarang, di seluruh folder `5-tiket/`, gelombang ditulis `GEL-2` dan `GEL-3`.**
> `SPEC-MODEL-DATA.md` §2.3 menulis `RETRO_KELUAR | G2` dan `PENCAPAIAN | G3` — **bacalah itu
> sebagai GEL-2 dan GEL-3**, bukan sebagai gerbang mana pun.

### 0.2 Urutan wewenang bila dua masukan bertentangan

1. **ADR** mengalahkan segalanya. (`docs/adr/`, ADR-0034…0055 ditambah warisan 0003–0033)
2. **`SPEC-INVARIAN.md`** mengalahkan `SPEC-MODEL-DATA.md` soal **aturan**.
3. ~~**Berkas DDL** mengalahkan spec soal **bentuk tabel** yang sudah ditetapkan.~~
   **DICORET — folder `ddl/` TIDAK ADA.** Lihat §4.3.
4. **Soal DAFTAR ENTITAS** — entitas apa saja yang ada —
   **`4-erd-dan-tabel-datar/STRUKTUR-DATA.md` mengikat.**
   `SPEC-MODEL-DATA.md` §2.3 **tidak lengkap dan tidak dipakai untuk menghitung.**
5. **Soal NAMA KOLOM dan TIPE — `2-to-spec/KAMUS-KOLOM.md` mengikat.**
   **Ditambahkan 24 September 2026** sebagai **pengganti butir 3 yang dicoret**, bukan sebagai
   pemulihannya.
6. **Ekspor Pega bukan wewenang lagi.** Ia sudah diadili. Tiket dibangun dari spesifikasi.

> ### Kenapa butir 5 ada, dan kenapa ia BUKAN butir 3 yang dihidupkan kembali
>
> Butir 3 dibuat untuk menyelesaikan **pertikaian** antara berkas DDL dan spesifikasi soal bentuk
> tabel. Pertikaian itu **tidak dapat terjadi lagi**: `2-to-spec/ddl-usulan/` dan
> `2-to-spec/KAMUS-KOLOM.md` **dibangkitkan dari satu berkas definisi** yang diurai dari
> `SPEC-MODEL-DATA.md` §10. Nama tabel, nama kolom, tipe, dan keterisian seluruhnya **turunan** —
> bila salah satunya menyimpang dari §10, **alatnya yang salah**, bukan spesifikasinya.
>
> Satu-satunya tempat DDL dan spec berbeda adalah **tiga pernyataan keputusan**, dan ketiganya
> berbunyi *"sengaja belum diputuskan"* — bukan *"diputuskan berbeda"*. Tidak ada yang perlu
> diadili di antara keduanya.
>
> **Maka yang dibutuhkan bukan pemutus pertikaian, melainkan sumber nama yang dapat disebut
> kriteria selesai tiket.** Itu `KAMUS-KOLOM.md`, dan ia **lebih kuat** daripada butir 3 yang lama:
> berkas DDL yang ditulis tangan **bisa** menyimpang dari spec; yang dibangkitkan **tidak bisa**.
>
> **Butir 3 tetap dicoret**, dan syarat pengembaliannya di §4.3 tidak berubah.

> **Kenapa butir 4 ada.** Urutan wewenang lama hanya menyelesaikan pertikaian tentang **aturan**,
> dan tidak pernah menyebut siapa yang menang soal **daftar entitas**. Pembaca yang membuka §2.3 dan
> `STRUKTUR-DATA.md` berdampingan sampai di persimpangan yang sama, dan belum tentu memilih yang
> sama. §2.3 kekurangan **tujuh** entitas — `BATAS_PER_BAHAYA` dan enam tabel acuan (L-6).
>
> **Sapuan sudah dijalankan.** Seluruh berkas `.md` di folder modul disisir untuk pernyataan yang
> menghitung entitas. Hasilnya: **tidak ada satu pun klaim jumlah entitas di seluruh korpus** selain
> angka "20" yang saya tulis sendiri di berkas ini pada 23 September dan sudah dikoreksi menjadi 27.
> `SPEC-MODEL-DATA.md` §11 menyebut "entitas selebihnya" tanpa angka, jadi **tidak ada kesimpulan
> lain yang bersandar pada hitungan §2.3.** Kekurangan tujuh itu tidak pernah merambat.

> **Butir 3 dicoret, bukan diberi catatan kaki.** Sesi DDL belum dijalankan, jadi anak tangga
> ketiga itu tidak ada. Urutan wewenang yang salah satu anak tangganya tidak ada akan dipakai orang
> **seolah ia ada** — karena itu ia dihapus dari urutannya, dan dikembalikan hanya ketika folder
> `ddl/` benar-benar lahir. Sampai saat itu, **tidak ada tiket yang boleh mengaku bersandar
> padanya**, dan bentuk tabel bukan wewenang siapa pun di alur ini.

---

## 1. GERBANG G1 — sebuah tiket itu sebesar apa

### Keputusan

> **SATU TIKET = SATU KEMAMPUAN, DIIRIS TEGAK.**
>
> Satu kemampuan adalah **satu hal yang dapat dilakukan seorang pelaku bernama, dari ujung ke
> ujung**, yang sebelumnya tidak dapat ia lakukan. Satu tiket menembus lapisan yang diperlukannya —
> Oracle, Go, React — dan selesai ketika pelaku itu benar-benar dapat melakukannya.
>
> **Tidak ada tiket lapisan.** Tidak ada tiket "buat tabel", "buat endpoint", "buat halaman".

### Kenapa bukan yang lain

| Pilihan | Kenapa ditolak |
|---|---|
| **per lapisan** (Oracle → Go → React) | melanggar uji G1 secara telak: tiket React menunggu Go, Go menunggu Oracle. Hampir setiap tiket menunggu tiket lain, dan tidak satu pun dapat dinyatakan selesai sendiri. Urutan lapisan memang ada, tetapi ia urutan **di dalam** satu tiket, bukan pembagian antar tiket |
| **per entitas** | dua puluh tujuh entitas — 21 milik modul + 6 tabel acuan, dihitung `DAFTAR-PEKERJAAN.md` §5.3 — dan sebagiannya tidak berarti sendirian — `NILAI_PENYEBARAN` tanpa `PENYEBARAN` bukan apa-apa. Lebih buruk: invarian terpenting justru **melintasi entitas** (INV-47, INV-50, INV-51 mengikat penyebaran ke induknya), sehingga pemisahan per entitas memotong tepat di tempat aturannya paling kuat |
| **per layar** | **tidak ada satu pun spesifikasi layar di seluruh masukan.** `SPEC-INVARIAN.md` nol menyebut layar, `ERD.md` nol. Membagi per layar berarti mengarang layarnya lebih dulu — dan itu menambal lubang, yang dilarang |
| **campuran bertingkat** | ukuran campur tidak dapat diurutkan, tidak dapat ditaksir, dan tidak dapat dibagi ke orang. Sekali satu tiket berukuran "seluruh cabang non-proporsional" berdampingan dengan "tambah satu kolom", perbandingan apa pun di antara keduanya tidak berarti |

### Aturan ukuran, supaya "kemampuan" tidak menjadi karet

Lima aturan — U-1, U-2, U-3, U-4, U-6. **U-5 sengaja tidak ada**: ia diusulkan, diuji, dan
**ditolak** karena melanggar U-1; nomornya dibiarkan kosong supaya penolakannya tidak hilang.
Kelimanya dapat diperiksa orang lain terhadap tiket yang sudah jadi.

**U-1 — Pelakunya disebut namanya, dan ia di luar tiket.**
Pelaku yang sah: **pengisi kontrak**, **pemberi persetujuan** (per tingkat: section head, dept head,
direktur), **pemeriksa jejak**, **pelaksana migrasi**, **konsumen hilir**. Pelaku yang **tidak** sah:
"sistem", "pengembang", "tiket berikutnya". Bila satu-satunya yang menikmati hasil sebuah tiket
adalah pengembang tiket lain, itu bukan tiket — isinya dilebur ke kemampuan pertama yang
membutuhkannya.

> Ini yang menghapus tiket lapisan tanpa aturan tambahan. Tabel, endpoint, dan halaman tidak pernah
> jadi tiket sendiri; ia **ikut** kemampuan pertama yang memerlukannya, dan tiket itu menyatakannya
> di medan **YANG TEGAS BUKAN BAGIAN TIKET INI** milik tetangganya.

**U-2 — Batas atas: satu tiket tidak melintasi lebih dari satu perpindahan keadaan.**
Ada **tiga belas** perpindahan sah (ADR-0055, termasuk perubahan 24 Sep 2026). "Mengajukan kontrak untuk persetujuan" satu tiket;
"menyetujui di tingkat section head" tiket lain. Kemampuan yang melintasi dua perpindahan dipecah
per perpindahan.

**U-3 — Batas bawah: kalimat SATU KALIMAT harus dapat ditulis tanpa kata "supaya nanti".**
"Pengisi kontrak dapat mencatat layer beserta limit dan deductible-nya, dan melihatnya kembali" —
sah. "Tabel layer tersedia supaya nanti dapat diisi" — bukan tiket.

**U-4 — Uji G1 dijalankan pada tiap tiket, bukan diasumsikan.**
Tiap tiket harus dapat **dinyatakan selesai tanpa menunggu tiket lain**, kecuali
ketergantungannya disebut terang di medan **BERGANTUNG PADA**. Bila lebih dari sepertiga tiket
saling menunggu, pembagiannya salah dan berkas ini diubah — bukan tiketnya dipaksakan.

**U-6 — SIFAT bukan KEMAMPUAN. Sifat dibawa setiap tiket sebagai invarian; bentuk fisiknya dibangun
pembuat pertama.**

Sebagian hal yang terasa seperti pekerjaan sebenarnya **sifat yang harus dipenuhi setiap tempat**,
bukan sesuatu yang seorang pelaku "dapat lakukan". Tandanya terbaca dari kolom entitasnya: bila ia
berbunyi *"seluruh entitas ber-…"*, hampir pasti ia sifat.

Ujinya satu pertanyaan: **adakah satu pelaku yang dapat melakukannya, dan dapatkah ia dinyatakan
selesai?** Bila jawabannya tidak, ia sifat.

Di bawah irisan tegak, sifat yang dipaksa menjadi tiket punya dua nasib yang sama buruknya:

| Nasib | Akibatnya |
|---|---|
| menjadi **tiket lapisan terselubung** | U-1 melarangnya — "bentuk paket uang dibangun" tidak punya pelaku di luar tiket |
| **tersebar** ke setiap tiket yang menyentuh uang | dikerjakan lima kali dengan lima tafsir |

**Perlakuannya, dan ini berlaku untuk semua sifat, bukan hanya paket uang:**

| | |
|---|---|
| **sifatnya** | dicoret dari daftar kemampuan; **nomor invariannya** masuk medan **INVARIAN YANG HARUS DIPENUHI** pada **setiap** tiket yang menyentuh hal itu |
| **bentuk fisiknya** | dibangun tiket ber-**PEMBUAT PERTAMA**, bersama tabel dan constraint-nya |

Jadi: **tidak ada tiket paket uang; ada kewajiban paket uang yang dibawa setiap tiket.**

> Pilihan ketiga — menulis ulang sifatnya sebagai kemampuan berpelaku lain, misalnya *"pemeriksa
> jejak dapat membaca setiap nilai uang beserta kurs dan tingkat pencatatannya"* — **ditolak sebagai
> aturan umum**, karena ia menjadikan setiap sifat dapat diselundupkan menjadi kemampuan hanya
> dengan mencari pembaca.

#### Garis yang memisahkan sifat dari kemampuan

Penolakan di atas menimbulkan pertanyaan yang sah: bila mencari pembaca dilarang, kenapa P-27 dan
P-28 boleh ditulis ulang menjadi *"PK melihat X"* sementara P-24 tidak? Dari luar keduanya terlihat
sebagai gerakan yang sama. Garisnya ini, dan ia ditulis supaya dapat **dibantah**, bukan dipercaya:

> ## SEBUAH SIFAT MENJADI KEMAMPUAN BILA PELANGGARANNYA ADALAH SESUATU YANG PELAKUNYA SADARI DAN KELUHKAN — bukan sekadar bila ada orang yang dapat membacanya.

Ujinya bukan *"adakah pembaca?"* melainkan *"bila ini dilanggar, apakah ada orang yang tahu ia
sedang dirugikan?"*

| | Lolos? | Sebab |
|---|---|---|
| **P-28** — kegagalan hitung menghasilkan keterangan, bukan nol | **ya** | melihat angka **nol** di layar padahal yang terjadi adalah perhitungan gagal adalah pengalaman yang berbeda dan keliru. Pengisi kontrak menyadarinya dan mengeluhkannya |
| **P-27** — kurs dibekukan pada versinya | **ya** | angka yang berubah sendiri bulan depan padahal tidak ada yang menyentuhnya adalah pengalaman yang berbeda — dan **justru keluhan itulah yang dulu menemukan cacatnya** |
| **P-24** — setiap nilai uang membawa mata uangnya | **tidak** | tidak ada saat di mana pelakunya menyadari bedanya. Ia sifat **penyimpanan**, dan pelanggarannya baru terasa jauh di hilir sebagai angka yang tidak dapat ditafsirkan |

Dengan garis ini tertulis, ketiga pencoretan dan kedua penyelamatan menjadi keputusan yang dapat
dibantah orang lain — bukan penilaian yang harus dipercaya.

### Dua ongkos pilihan ini, dan yang kedua lebih mengganggu

**Ongkos pertama — dua tiket dapat menyentuh tabel yang sama.**
`URUTAN-DAN-PENGHALANG.md` menyebut, untuk tiap entitas, **tiket mana yang pertama kali
membuatnya**; tiket sesudahnya menyatakan ketergantungan itu di medan **BERGANTUNG PADA**.

**Ongkos kedua — UKURAN YANG TIMPANG, dan U-4 tidak menangkapnya.**

Tiket **pertama** yang menyentuh sebuah entitas membawa seluruh pembuatannya: tabelnya,
constraint-nya, hak aksesnya, jejak perubahannya, pemindahan datanya. Tiket **kedua** yang
menyentuh entitas yang sama hampir tidak membawa apa-apa. Tiket pertama di tiap rantai bisa **tiga
sampai lima kali lebih besar** daripada tetangganya, padahal keduanya sama-sama disebut "satu
kemampuan".

U-4 memeriksa **saling menunggu**. Ini bukan soal menunggu — dua tiket yang tidak saling menunggu
pun tetap tidak sebanding. Akibatnya nyata: taksiran menyesatkan, pembagian ke orang tidak adil,
dan tiket pertama selalu terlambat sementara sisanya selalu lebih cepat dari dugaan — lalu orang
menyimpulkan **taksirannya** buruk, padahal **pembagiannya** yang timpang.

#### Keputusan: U-5 ditolak, penandaan yang dipakai

| Pilihan | Putusan |
|---|---|
| **tambah aturan U-5** yang memecah *pembuatan* entitas dari *pemakaiannya* | **ditolak** — ia melanggar U-1. "Entitas dibuat" **tidak punya pelaku di luar tiket**; satu-satunya yang menikmatinya adalah pengembang tiket berikutnya. Menerima U-5 berarti membatalkan aturan yang justru menghapus tiket lapisan |
| **akui dan tandai** | **dipakai** — jujur, murah, dan tidak mengubah apa pun |

**Bentuk penandaannya — medan wajib tambahan:**

> **PEMBUAT PERTAMA — `<daftar entitas>`**
> Diisi hanya pada tiket yang **pertama kali** membuat entitas itu; kosong pada tiket lain.
> Entitas yang disebut di sini adalah yang tabel, constraint, hak akses, jejak perubahan, dan
> pemindahan datanya **ikut di tiket ini**.

Dan satu aturan pembacaan yang mengikat siapa pun yang menaksir:

> **TAKSIRAN TIKET BER-PEMBUAT PERTAMA TIDAK DIBANDINGKAN DENGAN TAKSIRAN TIKET BIASA.**
> Keduanya bukan hal yang sejenis. Membandingkannya menghasilkan kesimpulan yang salah tentang
> orangnya, bukan tentang pekerjaannya.

Ini ditulis di sini justru supaya **tidak ditemukan lagi sebagai kejutan** saat taksiran pertama
meleset.

---

## 2. GERBANG G2 — apa yang masuk penyerahan pertama, dan apa yang tegas tidak

### Tiga jenis pekerjaan, dan kenapa pemisahannya wajib

| Golongan | Artinya | Diuji bagaimana |
|---|---|---|
| **PELESTARIAN** | sistem lama melakukannya, dan tetap dilakukan | **dapat** diuji terhadap data lama — hasilnya harus cocok |
| **PERUBAHAN** | sistem lama melakukannya, dan sengaja kita ubah | **tidak dapat** diuji terhadap data lama — data lama justru menunjukkan perilaku yang dibuang |
| **BARU** | sistem lama tidak pernah melakukannya | **tidak punya pembanding sama sekali** |

Mencampur ketiganya dalam satu keranjang membuat kegagalan uji tidak terbaca: cocok-tidaknya
terhadap data lama berarti hal yang berlawanan untuk pelestarian dan untuk perubahan.

### Keputusan

> **Penyerahan pertama memuat:**
> **(a) seluruh PELESTARIAN** yang diperlukan untuk mencatat dan menyetujui satu kontrak di kedua
> cabang;
> **(b) PERUBAHAN yang MENENTUKAN BENTUK** — yaitu yang menetapkan skema dan tidak dapat dipasang
> belakangan tanpa memindahkan ulang data yang sudah pindah;
> **(c) BARU hanya dalam MODE PERINGATAN** — menolak belum, memperingatkan sudah.

### (b) — daftar PERUBAHAN yang menentukan bentuk, disebut satu per satu

Bukan "perubahan penting". Ujinya sempit dan dapat diperiksa: **apakah memasangnya belakangan
menuntut migrasi ulang atau audit ulang kontrak yang sudah pindah?** Bila ya, ia masuk sekarang.

| ADR | Isinya | Kenapa tidak bisa belakangan |
|---|---|---|
| ADR-0040 | identitas dipecah `KONTRAK` / `VERSI_KONTRAK`; kunci alami **memperingatkan**, tidak melarang | memecah identitas setelah data pindah berarti memindahkan ulang seluruhnya |
| ADR-0046 + ADR-0055 | satu keadaan siklus hidup, tinggal pada **versi** | keadaan yang salah tempat harus dibaca ulang dari riwayat, dan riwayatnya tidak lengkap |
| ADR-0054 | `WARISAN_TAK_TERPETAKAN` + keadaan asli tetap tersimpan | bila nilai warisan tidak disimpan saat pindah, ia **hilang** — tidak dapat dipulihkan kemudian |
| ADR-0036 | beku saat disetujui, hitung saat dibaca | kontrak yang sudah disetujui tanpa pembekuan sudah terlanjur bergerak; membekukannya kemudian membekukan angka yang salah |
| ADR-0035 | kegagalan tidak pernah disamarkan menjadi nilai | nilai nol yang sebenarnya kegagalan tidak dapat dibedakan lagi setelah tersimpan |
| ADR-0037 | masukan versus turunan; turunan **tidak** disimpan | menghapus kolom turunan belakangan menuntut pembuktian tiap konsumennya sudah berhenti membacanya |
| ADR-0041 | satu fakta, satu penulis | ditegakkan **hak akses per objek**; memasangnya belakangan menuntut mencabut hak yang sudah dipakai |
| ADR-0039 (+0007) | paket uang membawa **tingkat pencatatan** | paket uang tanpa tingkat tidak dapat ditafsirkan ulang — 100 % treaty dan bagian NuRe tidak terbedakan |
| ADR-0045 | jejak perubahan sebagai fakta mesin | jejak yang tidak dicatat sejak awal **tidak ada**, titik |
| ADR-0042 | sejarah pindah apa adanya, aturan sentuh-perbaiki | sama: yang tidak ikut saat pindah tidak dapat disusulkan |
| ADR-0053 | mata uang sebagai **daftar** | skalar yang sudah terisi tidak menyimpan baris kedua yang dulu ada |

**Yang TIDAK masuk butir (b)** meski ia PERUBAHAN: perubahan yang dapat dipasang belakangan tanpa
menyentuh data yang sudah pindah — misalnya penajaman pesan galat, atau penarikan sebuah invarian
dari aplikasi ke basis data. Itu menunggu penyerahan berikutnya.

### (c) — BARU dalam mode peringatan, dan kenapa bukan mode menolak

Kasus utamanya **enam dari delapan syarat pengajuan** (`SPEC-INVARIAN.md` §3, daftar K1): K1-3
sampai K1-8 **mati di sistem lama**. Artinya pengguna **tidak pernah mengalaminya**, dan sebagian
kontrak yang selama ini lolos memang tidak memenuhinya.

Menyalakan keenamnya serentak di hari cut-over menolak pekerjaan orang pada hari yang sudah paling
rapuh. `SPEC-INVARIAN.md` §3 sudah menetapkan langkah CO-6, CO-7, CO-8 untuk ini, dan **setiap tiket
bergolongan BARU wajib mengisi medan CARA MENYALAKANNYA** yang merujuk ketiganya.

> **Ini bukan alasan tidak membangunnya.** Yang ditunda adalah **penolakannya**, bukan
> pemeriksaannya. Pemeriksaan tetap berjalan sejak hari pertama dan hasilnya tercatat; yang belum
> menyala hanyalah kuasa memblokir.

#### Mode peringatan tanpa pembaca peringatan adalah fitur yang dimatikan, ditambah biaya log

Pemeriksaan yang berjalan dan mencatat menuntut **seseorang membaca catatannya**. Bila tidak, mode
peringatan bukan penyalaan bertahap — ia fitur mati berbiaya. Dan yang lebih sering terjadi: **mode
peringatan tidak pernah berakhir**, karena tidak pernah ada yang memutuskan kapan ia berakhir.

Karena itu medan **CARA MENYALAKANNYA** pada setiap tiket bergolongan **BARU** memuat **empat** hal,
bukan satu:

| # | Isi | Bentuk yang diterima |
|---|---|---|
| 1 | rujukan ke **CO-6, CO-7, CO-8** | `SPEC-INVARIAN.md` §3 |
| 2 | **SIAPA** yang membaca hasil pemeriksaan, dan **seberapa sering** | peran bernama + selang waktu. "Tim teknik" bukan jawaban; "pemilik proses Treaty In, mingguan" jawaban |
| 3 | **APA YANG HARUS TERJADI** supaya kuasa memblokir dinyalakan | **angka yang bisa diperiksa**, bukan perasaan. "Tingkat pelanggaran di bawah 2 % selama 4 minggu berturut-turut" bentuknya; "kalau sudah stabil" bukan |
| 4 | **SIAPA yang berwenang** menyalakannya | satu peran, bernama |

> Tanpa butir 2 dan 3, keenam syarat pengajuan itu akan berstatus "peringatan" **selamanya**, dan
> kita sudah membangunnya tanpa pernah memakainya.

### Aturan yang akan berlaku lagi saat GEL-2 diberi ruang lingkup

> ## PENGECUALIAN SETINGKAT MODUL TIDAK MENGURANGI JUMLAH ENTITAS DI DALAM MODUL YANG DIPERTAHANKAN.

Kesalahan yang sudah terjadi sekali, ditulis supaya tidak terjadi dua kali: sebelum
`DAFTAR-PEKERJAAN.md` disusun, ruang lingkup G2 dikira **memotong isi kontrak**, sehingga sisa
pekerjaan §10 dikira mengecil. **Ia memotong modul dan gelombang, bukan isi kontrak.**

Sebabnya sederhana dan akan berlaku lagi:

> **Sebuah kontrak menuntut hampir seluruh bagiannya untuk dapat dicatat dan disetujui.
> Tidak ada kontrak yang setengah dicatat.**

Yang dikeluarkan G2 — `RETRO_KELUAR`, `PENCAPAIAN`, `NILAI_SELISIH` — hanya **tiga**, dan ketiganya
memang sudah di luar daftar entitas inti sejak awal. Dua puluh satu entitas milik modul tetap
terpakai seluruhnya, ditambah enam tabel acuan. Hasilnya **27**, bukan angka yang lebih kecil.

**Saat GEL-2 diberi ruang lingkup nanti, kesalahan yang sama menunggu di tempat yang sama.** Cara
menghindarinya sudah terbukti: **susun daftar kemampuannya lebih dulu, baru turunkan entitasnya
dari situ** — jangan menaksir jumlah entitas dari batas ruang lingkupnya.

### Yang TEGAS TIDAK masuk penyerahan pertama

Ruang lingkup tanpa tepi akan tumbuh sendiri. Sembilan tepi berikut mengikat.

| Tidak masuk | Sebab | Dilacak di |
|---|---|---|
| `RETRO_KELUAR` | **GEL-2** — arah keluar, bukan pokok yang sama | `SPEC-MODEL-DATA.md` §2.3 |
| `PENCAPAIAN` | **GEL-3** — dimodelkan sekarang, tidak dibangun sekarang | `SPEC-MODEL-DATA.md` §2.3 |
| perilaku modul **Treaty In Adjustment** | embargo masih berlaku; yang dibangun hanya **sambungannya** | `TEMUAN-ADJUSTMENT-DITUNDA.md` |
| jalur **penerbitan ke luar** | `TREATYINOFFER` tidak punya satu pun penulis yang terjangkau; penggantinya belum diketahui | `DAFTAR-ESKALASI-MANAJEMEN.md` butir 3 |
| **perbaikan data lama** dan hitung ulang | ADR-0043: menghitung ulang adalah **peristiwa bisnis tersendiri**, bukan bagian migrasi | ADR-0043 |
| invarian atas `AccountingMode` | menunggu Uji X-2 | `SPEC-INVARIAN.md` §6 |
| apa pun yang **mewujudkan butir eskalasi yang belum diputuskan** | butir 1, 2, 3, 5, 6, 7 belum ada keputusannya | `DAFTAR-ESKALASI-MANAJEMEN.md` |
| layar di luar yang dituntut kemampuan di penyerahan ini | tidak ada spesifikasi layar; menambahnya berarti mengarang | §4 berkas ini |
| **tombol pengembang** sistem lama | sudah diputuskan tidak dibawa, dan tidak dibahas ulang | keputusan pemilik proses |

---

## 3. GERBANG G3 — tiket yang bergantung pada jawaban yang belum ada

### Keputusan

> **DIPECAH.** Tiket dipotong tepat di garis tempat kepastian berhenti. Bagian yang pasti menjadi
> tiket biasa — lengkap, dapat ditaksir, dapat dimulai. Bagian yang menunggu menjadi **tiket
> tertahan** bernomor sendiri, dengan medan **PENGHALANG**, dan **tanpa taksiran**.

### Ini bukan pilihan selera — dua larangan sudah memaksanya

Larangan yang berlaku di alur ini:

- *"Satu tiket memuat hal yang sudah pasti dan hal yang belum diputuskan"* — **dilarang**.
- *"Taksiran pada tiket yang penghalangnya belum terjawab"* — **dilarang**.

Cocokkan dengan tiga bentuk yang ditawarkan:

| Bentuk | Bertahan? |
|---|---|
| **DITAHAN seluruhnya** | menahan juga bagian yang sudah pasti, padahal bagian itu dapat dikerjakan besok. Membuang pekerjaan yang tersedia tanpa sebab |
| **DITULIS dengan medan penghalang, boleh ditaksir** | **melanggar larangan kedua secara langsung** |
| **DIPECAH** | satu-satunya yang tidak melanggar keduanya |

### Bentuk tiket tertahan

Tiket tertahan tetap **bernomor** dan tetap muncul di `PETA-LIPUTAN.md`. Tiket tanpa nomor tidak
dapat ditunjuk, dan yang tidak dapat ditunjuk akan hilang.

Medan **PENGHALANG** wajib memuat tiga hal, bukan satu:

| Isi | Kenapa |
|---|---|
| **apa** yang ditunggu | supaya jelas kapan ia berhenti tertahan |
| **siapa** yang dapat menjawabnya | penghalang tanpa pemilik tidak pernah terjawab |
| **apa yang berubah** pada tiketnya ketika jawabannya datang | supaya jawaban yang datang langsung dapat dipakai, bukan memicu pembahasan dari nol |

Taksiran: **kosong**, dan kekosongan itu disengaja — bukan lupa diisi.

### Kasus merosot yang diakui, bukan disembunyikan

Ada kemampuan yang **tidak dapat dipecah** karena hal yang belum diputuskan duduk tepat di
tengahnya — misalnya kemampuan yang seluruhnya bergantung pada batas wewenang persetujuan yang
dokumennya belum ada (eskalasi butir 5). Untuk itu **seluruh tiketnya tertahan**, dan pada
tiketnya ditulis **"tidak dapat dipecah, dan sebabnya: …"**.

Itu tetap hasil dari aturan yang sama, bukan pengecualian terhadapnya: pemecahan dicoba, dan
kegagalannya dicatat.

---

## 4. Lubang yang ditemukan saat menjawab ketiga gerbang

Dilaporkan, **tidak ditambal**. Menambalnya dari ekspor akan menjadi keputusan tak bertanda yang
tidak pernah ditinjau siapa pun.

**L-1, L-2, L-4, L-5, L-6, L-7 masuk `LUBANG-SPESIFIKASI.md`. L-3 TIDAK** — ia lubang lingkungan, bukan
lubang spesifikasi, dan tempatnya di daftar yang harus disediakan kantor (§4.4).

| # | Lubang | Akibatnya pada tiket |
|---|---|---|
| **L-1** | **Folder `ddl/` tidak ada.** Sesi DDL ditangguhkan; wewenang urutan ke-3 menunjuk berkas yang belum lahir | nama constraint, angka presisi INV-45, bentuk fisik `POTONGAN` berinduk dua, dan pilihan menurunkan INV-35 jadi constraint — seluruhnya belum ada. Tiket yang menyentuhnya **dipecah** menurut G3 |
| **L-2** | **`SPEC-MODEL-DATA.md` §10 baru meliputi 4 entitas** — `KONTRAK`, `VERSI_KONTRAK`, `LAYER`, `DETAIL_PROPORSIONAL`. §11 menyatakannya sendiri: sisanya **belum dikerjakan** | entitas selebihnya tidak punya nama atribut, tipe, maupun keterisian. Kriteria selesai yang menyebut field tidak dapat ditulis untuk entitas itu. **Ini lubang terbesar di seluruh alur.** Jumlah pastinya baru diketahui setelah langkah 2 — lihat `DAFTAR-PEKERJAAN.md` §5.3: **23 entitas**, bukan 16 |
| **L-7** | **ADR-0055 tidak punya cara membuang draf.** `DRAFT` punya tepat satu perpindahan keluar: `AJUKAN`. Digabung INV-25 dan INV-22, pengisi kontrak yang salah pencet **terkunci** | satu-satunya jalan keluar adalah mengajukannya supaya ditolak. **Keputusan milik ADR-0055**, dan harus diambil **sebelum** sesi to-spec melanjutkan §10 — lihat `LUBANG-SPESIFIKASI.md` L-7 |
| **L-6** | **Daftar entitas di `SPEC-MODEL-DATA.md` §2.3 tidak lengkap.** `BATAS_PER_BAHAYA` diperkenalkan §10.0b dan ditegakkan INV-14, tetapi **tidak ada di §2.3**. Enam tabel acuan (`MATA_UANG`, `JENIS_POTONGAN`, `JENIS_REASURANSI`, `BAHAYA`, `KELOMPOK_TREATY`, `KELAS_BISNIS`) juga tidak didaftar di sana | siapa pun yang menghitung ruang lingkup dari §2.3 akan **kekurangan tujuh entitas**. Daftar yang mengikat adalah `4-erd-dan-tabel-datar/STRUKTUR-DATA.md`, bukan §2.3 |
| **L-3** | **Empat invarian berstatus BELUM DIBUKTIKAN** — INV-47, INV-50, INV-51, dan bentuk lintas baris INV-31, ditegakkan lewat *materialized view* `REFRESH ON COMMIT`. Uji negatifnya **belum pernah dijalankan**; tidak ada instans Oracle yang terjangkau | bila salah satu gagal, keempatnya **turun ke aplikasi** dan tiketnya berubah bentuk. Tiket yang menegakkannya ditulis dengan penghalang ini disebut |
| **L-4** | **Tidak ada spesifikasi layar di mana pun.** `SPEC-INVARIAN.md` nol menyebut layar; `ERD.md` nol | bentuk layar tidak dapat dijadikan kriteria selesai. Kriteria dinyatakan sebagai **perilaku yang dapat gagal**, bukan sebagai susunan tampilan |
| **L-5** | **Huruf `G` bermakna tiga hal** (§0.1) | sudah ditangani di sini dengan menetapkan `GEL-2`/`GEL-3`. Dicatat supaya pembaca `SPEC-MODEL-DATA.md` §2.3 tidak salah baca |

### 4.1 L-1, L-2, L-3 adalah satu kalimat yang sama

Ketiganya bukan tiga lubang yang kebetulan berdampingan:

| | Lubang | Kalimatnya |
|---|---|---|
| L-1 | folder `ddl/` tidak ada | sesi **DDL** belum dijalankan |
| L-2 | §10 baru meliputi 4 dari 20 | sesi **to-spec** belum selesai |
| L-3 | empat invarian belum terbukti | **uji negatifnya** belum dijalankan |

> **LANGKAH SEBELUMNYA BELUM SELESAI.** Menjawab L-2 sendirian membuat percakapan yang sama
> terulang persis di L-1, lalu di L-3. Karena itu ketiganya diputuskan sekaligus, di sini.

### 4.2 L-2 — jalan (iii), dan kenapa bukan (i) maupun (ii)

Dua jalan yang mula-mula terlihat, keduanya salah:

- **(i) selesaikan §10 untuk 16 entitas sisanya** — **mengerjakan lebih banyak daripada yang
  diperlukan.** G2 sudah mengeluarkan `RETRO_KELUAR`, `PENCAPAIAN`, perilaku Adjustment, dan jalur
  penerbitan hilir dari penyerahan pertama. Sebagian dari keenam belas entitas itu hampir pasti
  tidak terpakai sekarang, dan §10 untuk entitas yang tidak dibangun dikerjakan **dua kali**:
  sekarang dengan setengah pengetahuan, dan lagi nanti saat gilirannya tiba.
- **(ii) tulis keenam belasnya sebagai tiket tertahan** — menghasilkan belasan tiket tertahan
  sekaligus. Itu tanda ada yang salah, dan yang salah adalah **urutannya**, bukan bentuk tiketnya.

**JALAN (iii) — yang dikerjakan:**

| | Langkah | Kenapa bisa |
|---|---|---|
| **a** | **jalankan langkah 2 sekarang** — `DAFTAR-PEKERJAAN.md` | ia **tidak memerlukan §10**. Sebuah kemampuan berbunyi "pelaku X dapat melakukan Y", dan kalimat itu tidak butuh satu pun nama atribut. Yang dibutuhkannya sudah lengkap: daftar entitas, `SPEC-INVARIAN.md`, dua belas perpindahan, dan batas G2. Ia juga belum ditaksir dan belum jadi tiket, jadi tidak ada larangan yang tersentuh |
| **b** | dari daftar itu, turunkan **entitas mana saja yang benar-benar tersentuh** penyerahan pertama | keluarannya satu daftar pendek |
| **c** | **baru §10 dikerjakan** — hanya untuk entitas di daftar itu | pekerjaan **sesi to-spec**, bukan sesi ini. Sisanya tetap kosong dan **tercatat kosong** |
| **d** | sesudah itu tiket ditulis | ruang lingkupnya sudah diketahui, bukan ditebak |

Jalan ini mengerjakan §10 **sekali**, untuk yang terpakai saja.

### 4.3 L-1 — butir 3 urutan wewenang DICORET, bukan diberi catatan kaki

> ### Peninjauan 24 September 2026 — folder DDL sudah lahir, dan butir 3 TETAP DICORET
>
> `2-to-spec/ddl-usulan/` kini berdiri: 29 berkas tabel, `00_SKEMA_DAN_AKUN.sql`,
> `Z00_KUNCI_ALAMI.sql`. Dibaca harfiah, syarat pengembalian di bawah — *"dikembalikan hanya ketika
> folder `ddl/` benar-benar lahir"* — terpenuhi. **Ia tetap tidak dikembalikan**, dan tiga hal yang
> terbaca dari isinya sendiri yang menentukan:
>
> 1. namanya **`ddl-usulan/`**, dan kepala setiap berkas berbunyi **"USULAN. BELUM PERNAH
>    DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN."** Ia keluaran sesi **to-spec**, bukan
>    keluaran sesi DDL yang kalimat di bawah maksud;
> 2. isinya **belum ditinjau siapa pun**;
> 3. **dan ini yang memutuskan** — isinya **sengaja bertentangan dengan spesifikasi di tiga
>    tempat**, masing-masing bertuliskan pernyataan keputusan. Yang terberat: **paket uang belum
>    punya kolom mata uang**. Bila butir 3 dipulihkan, maka *"berkas DDL mengalahkan spec soal
>    bentuk tabel"* berarti **"jumlah uang tanpa denominasi" menang atas §10** — dan tiket dapat
>    ditulis di atasnya.
>
> **Anak tangga yang ADA tetapi berisi usulan yang belum ditinjau lebih berbahaya daripada anak
> tangga yang tidak ada** — ia tidak menuntut siapa pun berpura-pura; ia benar-benar di sana dan
> tampak berwibawa.
>
> **Syarat pengembaliannya diperjelas, supaya pertanyaan yang sama tidak kembali dengan jawaban yang
> berbeda.** Yang mengembalikan butir 3 **bukan lahirnya sebuah folder**, melainkan **tertutupnya
> gerbang sesi DDL** — yaitu ketika ketiganya sudah diputuskan:
>
> | Yang harus diputuskan | Keadaan 24 Sep 2026 |
> |---|---|
> | angka presisi per kelompok tipe | **belum** — mewarisi modul tetangga sebagai usulan |
> | bentuk fisik **paket uang** | **sebagian** — golongan A dan B diputuskan 24 Sep; golongan C menunggu `T-6` |
> | bentuk fisik **induk polimorfik** (`POTONGAN`, `PENYEBARAN`) | **belum** |
>
> **Dan satu hal yang TIDAK menunggu ketiganya.** Kebutuhan yang sebenarnya — *sumber nama kolom
> yang dapat disebut kriteria selesai tiket* — dipenuhi **butir 5** urutan wewenang, ditambahkan
> pada hari yang sama. `2-to-spec/KAMUS-KOLOM.md` memberi **240 kolom bernama** dan **tidak dapat
> menyimpang dari §10**, karena ia dibangkitkan darinya.
>
> **Maka to-ticket tidak lagi tertahan oleh butir 3.** Yang masih menahan adalah keempat penahan
> lain di `DAFTAR-PEKERJAAN.md` §1.1 — bukan ketiadaan nama kolom.



Sesi DDL dijalankan **setelah** §10 selesai untuk entitas yang terpakai (jalan (iii) langkah c).

Selama folder `ddl/` belum ada, **butir 3 urutan wewenang tidak berlaku dan dicoret** — lihat §0.2.
Bukan diberi catatan kaki: urutan wewenang yang salah satu anak tangganya tidak ada akan dipakai
orang **seolah ia ada**.

### 4.4 L-3 — bukan lubang spesifikasi, melainkan lubang LINGKUNGAN

Uji negatif *materialized view* tidak dapat dijalankan karena **tidak ada instans Oracle yang
terjangkau**. Itu bukan sesuatu yang dapat ditutup sesi mana pun — bukan sesi to-spec, bukan sesi
DDL, bukan sesi ini.

Karena itu L-3 **tidak masuk `LUBANG-SPESIFIKASI.md`**. Ia masuk daftar **yang harus disediakan
kantor**, bersama perkara lingkungan lain. Menaruhnya di daftar lubang spesifikasi membuatnya
tampak seakan ada sesi yang bisa mengerjakannya, dan tidak ada.

---

## 5. Ringkasan ketiga keputusan, satu tabel

| Gerbang | Keputusan | Ujinya |
|---|---|---|
| **G1 — ukuran** | satu tiket = **satu kemampuan, diiris tegak**, pelakunya disebut namanya; tidak ada tiket lapisan | tiap tiket dapat dinyatakan selesai sendiri; ketergantungan disebut terang; tidak lebih dari sepertiga tiket saling menunggu |
| **G2 — penyerahan pertama** | seluruh **PELESTARIAN** + **PERUBAHAN yang menentukan bentuk** (11 ADR, didaftar) + **BARU dalam mode peringatan** | sembilan tepi yang tegas tidak masuk, didaftar |
| **G3 — yang tertahan** | **DIPECAH**; bagian pasti jadi tiket biasa, bagian menunggu jadi tiket tertahan bernomor, tanpa taksiran, dengan penghalang bersebut apa/siapa/apa-yang-berubah | kasus yang tidak dapat dipecah ditulis sebabnya, bukan disembunyikan |

---

## 6. Apa yang terjadi berikutnya

Menurut urutan kerja yang berlaku, **langkah 2** adalah `DAFTAR-PEKERJAAN.md` — daftar hal yang
harus ada di sistem baru, dengan asal-usulnya, **belum dibentuk jadi tiket dan belum ditaksir**.

Sebelum itu dimulai, satu hal menunggu jawaban: **jalan (i) atau (ii) untuk lubang L-2.**
