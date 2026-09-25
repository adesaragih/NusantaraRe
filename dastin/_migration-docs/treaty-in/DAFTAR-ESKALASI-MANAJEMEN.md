# Temuan yang perlu diketahui manajemen — Treaty In

**Ditulis untuk:** direksi dan pemilik proses NuRe
**Tanggal:** 23 September 2026

Berkas ini bukan laporan teknis. Ia memuat hal-hal yang ditemukan saat membedah sistem Treaty In
yang berjalan sekarang, dan yang **perlu diketahui apa pun nasib proyek migrasi** — termasuk bila
proyek itu ditunda atau dibatalkan.

Tidak ada istilah teknis di sini. Setiap butir menyebut akibatnya bagi bisnis dan apa yang
dibutuhkan untuk memastikan seberapa besar.

---

## 1. Ada tombol yang menyetujui kontrak tanpa melewati pemberi persetujuan

**Apa yang terjadi.** Di layar penawaran treaty terdapat sebuah tombol yang menyetel sebuah kontrak
langsung menjadi **selesai disetujui**, tanpa kontrak itu melewati satu pun pemberi persetujuan.
Tombol itu bukan bagian dari alur kerja; ia dipasang oleh pengembang sistem pada tahun 2020 dan
diberi nama yang menandai dirinya sendiri sebagai alat pengembang.

**Apa yang sudah dipastikan.** Kami sudah memeriksa siapa yang dapat melihat tombol itu.
Jawabannya: tombol itu **dibatasi pada dua nama orang tertentu** — keduanya dari pihak pengembang,
bukan pengguna bisnis. Jadi hari ini tombol itu **tidak terlihat oleh pengguna biasa**.

**Apa yang belum dijawab.** "Tidak terlihat hari ini" bukan jawaban atas "pernah dipakai atau
tidak". Pertanyaan itu hanya bisa dijawab oleh data, dan caranya sederhana: kontrak yang disetujui
tanpa melewati approver akan terlihat persis begitu — **status akhirnya ada, jejak persetujuannya
tidak**. Kueri untuk menghitungnya sudah disiapkan dan berjalan bersama pemeriksaan lain.

**Akibatnya bagi bisnis.** Bila hasilnya nol, butir ini selesai dan tidak diangkat lagi. Bila
hasilnya tidak nol, artinya ada kontrak yang berlaku hari ini tanpa pernah disetujui siapa pun —
dan itu perlu diketahui terlepas dari nasib proyek migrasi.

**Temuan kedua di layar yang sama, dan ini yang terbuka.** Tombol berikutnya, yang **membuka kembali
kunci** sebuah kontrak yang sudah disetujui sehingga isinya dapat disunting lagi, dibatasi lebih
longgar daripada dua tombol pengembang di sebelahnya: ia terlihat oleh **setiap orang yang terdaftar
pada divisi TI** — bukan hanya dua nama, tetapi juga bukan setiap pengguna. Pemeriksaan pertama kami
hanya membaca satu lapis syarat dan melewatkan lapis di atasnya; koreksi ini hasil pemeriksaan ulang. Yang meredam akibatnya: setelah
kunci terbuka, tombol simpan untuk kontrak yang sudah disetujui **tidak ikut muncul** bagi pengguna
biasa — sejauh yang dapat dipastikan dari layar-layar yang kami periksa. Jadi pintunya terbuka,
tetapi jalan keluarnya tidak terlihat. Kami menyatakan ini apa adanya karena batas pemeriksaannya
memang di situ: kami memeriksa tombol, bukan menjalankan sistemnya.

**Tambahan 24 September 2026, DIRALAT 26 September 2026 — peredam itu tidak berlaku untuk satu
tombol lain di layar yang sama.** Di layar penawaran yang sama terdapat **sebuah tombol** yang
mengerjakan pembukaan kunci **dan penyimpanan sekaligus dalam satu tindakan**: sekali ditekan,
status persetujuan kontrak dikosongkan, nama penekannya dicatat, sebuah komentar ditambahkan, dan
seluruh isi kontrak disimpan kembali. Karena penyimpanannya menyatu dengan tombol itu, **tidak ada
tombol simpan terpisah yang perlu muncul** — sehingga kalimat peredam di atas tidak melindungi jalur
ini.

Dua hal kami nyatakan lebih tepat daripada versi 24 September. **Pertama**, tombol itu **dibatasi
pada peran pengisi kontrak**, bukan terbuka untuk siapa pun. **Kedua**, tombol semacam itu ada
**satu**, bukan dua: tombol keduanya yang sempat kami sebut ternyata justru mensyaratkan kontrak
yang **belum** selesai disetujui, sehingga ia tidak dapat menghapus persetujuan yang sudah ada — dan
ia juga tidak memendekkan jalur persetujuannya. Kami memeriksa kemungkinan itu secara khusus, karena
bila benar, kontrak baru dapat disetujui hanya oleh kepala seksi. Hasilnya: tidak benar.

**Dan satu hal yang memperingan, yang baru kami baca 26 September.** Tombol itu, selain
mengosongkan status persetujuan, juga **mengunci sebagian besar angka uang di muka kontrak** —
batas per bahaya, lapisan, mata uang, bagian, dan brokerase. Jadi bentuknya menyerupai sebuah
kebijakan: *perubahan yang tidak menyentuh angka uang cukup ditinjau lebih pendek*. **Yang menjadi
masalah bukan jalur pendeknya, melainkan bahwa kuncinya bocor**: nama cedant, tanggal mulai, dan
tanggal berakhir — tiga hal yang justru menentukan kontrak itu kontrak yang mana — **tidak ikut
terkunci**, dan tetap dapat diubah. Apakah jalur pendek itu memang kebijakan yang dikehendaki
adalah pertanyaan yang kami ajukan ke bagian teknik treaty, bukan sesuatu yang dapat kami putuskan.

**Cara mengenali kontrak yang pernah dibuka lewat tombol itu.** Nama tombolnya kini terbaca:
**Revision**. Di layar yang sama ada **Copy**, **View**, dan **Edit**; hanya Revision yang bekerja
atas kontrak yang sudah selesai disetujui. Tetapi tindakannya meninggalkan tanda yang
jelas: kontrak yang pernah dibuka dengan cara itu membawa **sebuah komentar bertuliskan "Create
Revision"** beserta nama penekannya dan waktunya. Tanda itu cukup untuk menemukannya kembali.

**Tambahan 24 September 2026 — jalur pendeknya kini TERBACA PENUH, dan panjangnya SATU.**

Paragraf di atas menduga bentuknya *"menyerupai sebuah kebijakan: perubahan yang tidak menyentuh
angka uang cukup ditinjau lebih pendek"*. Dugaan itu **sekarang terbukti, dan ukurannya diketahui.**

Mesin keadaan persetujuannya — `DataTransform/Akseptasi_DT.xml` — memuat **dua cabang utama**, dan
yang kedua dipilih persis oleh penanda yang dipasang tombol **Revision**:

```
RevisionState != 1   : pengajuan -> Kepala Seksi -> Kepala Departemen -> Direktur -> selesai
RevisionState == 1   : pengajuan -> Kepala Seksi -> selesai
```

> **Revisi atas kontrak yang sudah disetujui Direktur disahkan kembali oleh SATU tingkat.**
> Kepala Departemen dan Direktur dilewati, dan **tidak ada ambang nilai di mana pun pada cabang
> itu** — revisi sebesar apa pun melewati jalur yang sama.

**Kenapa ini digabung ke butir ini alih-alih menjadi butir tersendiri.** Ia temuan yang sama dilihat
dari arah lain: butir ini menemukan **tombolnya**, pembacaan ini menemukan **jalurnya**. Satu butir
dengan dua bukti lebih kuat daripada dua butir yang pembacanya harus menyadari sama.

**Sifat pembuktiannya `EVIDENCED`** — terbaca dari aturan yang hidup, bukan disimpulkan. Yang belum
diketahui hanya **berapa kontrak yang pernah melewatinya**, dan itu pengukur kerugian, bukan
pemblokir model; penandanya sudah disebut di paragraf sebelumnya — komentar *"Create Revision"*.

**Pertanyaan bisnisnya, satu kalimat:** *"Apakah perubahan atas kontrak yang sudah disetujui
Direktur memang boleh disahkan kembali oleh Kepala Seksi saja?"* Selama belum dijawab, **ADR-0055
tidak memuat perpindahan untuk jalur revisi** — kekosongan yang dinyatakan, bukan kekosongan.

**Batas pemeriksaan butir tambahan ini, dinyatakan sejujurnya sama seperti paragraf di atas.**
Perilaku kedua tombol ini **terbaca dari aturan yang mengatur layar itu, dan belum kami jalankan**.
Kami tidak menekan tombolnya dan tidak mengamati akibatnya pada sistem yang berjalan. Keduanya ada
di **satu tempat**, yaitu susunan layar penawaran treaty; kemunculannya yang kedua pada berkas
pembungkus layar itu adalah salinan dari susunan yang sama, bukan tombol tambahan.

**Temuan ketiga, 26 September 2026, dan ini pada layar addendum.** Peredam "tombol simpan tidak ikut
muncul" berlaku untuk tombol simpan yang biasa: addendum yang sudah disetujui memang tidak
menampilkannya. Tetapi di layar yang sama ada **tombol simpan kedua**, bertanda pengembang, yang
**tidak menyebut status persetujuan sama sekali**. Ia terlihat oleh setiap orang yang terdaftar pada
divisi TI — golongan yang sama yang melihat tombol pembuka kunci.

Gabungan keduanya berarti: seorang petugas TI dapat membuka kunci sebuah addendum yang **sudah
disetujui**, mengubah isinya, dan menyimpannya di tempat — tanpa versi baru dan tanpa persetujuan
ulang.

**Batas pemeriksaan butir ini, dinyatakan sama seperti paragraf di atas.** Perilaku kedua tombol ini
**terbaca dari aturan yang mengatur layar itu, dan belum kami jalankan**. Kami tidak menekan
tombolnya dan tidak mengamati akibatnya pada sistem yang berjalan.

**Dan ini dapat diukur — kami sempat menulis sebaliknya, lalu memeriksanya.** Ada dua tanda yang
saling menguatkan. **Pertama**, tabel addendum menyimpan tanggal pembaruan yang **disegarkan setiap
kali baris itu disimpan**; addendum yang tanggal pembaruannya lebih baru daripada tanggal
persetujuan terakhirnya adalah baris yang disimpan sesudah disetujui. **Kedua**, tombol simpan itu
**tidak menghitung ulang tabel selisih**, sehingga baris yang disunting sesudah disetujui akan
menyimpan selisih yang tidak lagi cocok dengan isinya sendiri. Keduanya tercatat sebagai **UA-16**
pada lampiran pengukuran.

**Yang belum diukur:** berapa kontrak yang kini berstatus belum disetujui padahal pernah disetujui,
dan berapa yang membawa tanda "Create Revision" — tercatat sebagai **UA-11** pada lampiran
pengukuran.

**Pertanyaan yang tersisa bukan pertanyaan teknis.** Membatasi tombol "selesai disetujui" kepada dua
pengembang tidak menghapus perkaranya; ia mengubah bentuknya menjadi pertanyaan wewenang: **siapa
yang mengizinkan pihak pengembang menyelesaikan persetujuan pada data produksi, dan kepada siapa
pemakaiannya dilaporkan?** Kami tidak menjawabnya — ia milik manajemen, bukan milik pemeriksaan
artefak.

**Yang dibutuhkan:** tidak ada keputusan manajemen yang tertahan oleh butir ini. Ia dicantumkan
karena bobotnya berbeda jenis dari butir-butir lain di bawah — yang lain berbicara tentang kendali
yang **lemah**, butir ini tentang kendali yang dapat **dilewati**.

**Di sistem baru.** Kedua tombol tidak dibawa. Tidak ada jalan menyetel keadaan sebuah kontrak
selain lewat perpindahan yang sah, dan tidak ada jalan menyunting kontrak yang sudah disetujui
selain dengan membuat versi baru yang menempuh persetujuannya sendiri.

---

## 2. Angka yang sudah disetujui tidak dijamin tetap

**Apa yang terjadi.** Ketika sebuah kontrak treaty disetujui, angka-angkanya — pembagian kapasitas,
nilai rupiah, bagian NuRe — tidak dikunci. Bila kontrak itu disentuh lagi kemudian, sebagian
angkanya dihitung ulang dari data acuan yang berlaku **saat itu**, bukan yang berlaku saat
persetujuan diberikan.

**Ada empat sebab yang berdiri sendiri:**

1. Pembagian kapasitas disusun ulang dari tabel acuan setiap kali kontrak dihitung.
2. Tabel acuan itu bisa disunting kapan saja, termasuk setelah kontrak disetujui.
3. Nilai rupiah memakai kurs **tahun berjalan saat perhitungan**, bukan kurs yang berlaku saat
   kontrak disepakati.
4. Bagian NuRe yang dipakai untuk menghitung selalu bagian **hari ini**, juga untuk premi yang
   diterima bertahun-tahun lalu.

**Akibatnya bagi bisnis.** Kontrak yang sudah melewati empat tingkat persetujuan dapat berubah
pembagian kapasitasnya dan nilai rupiahnya tanpa ada yang menyetujui perubahan itu, dan tanpa ada
yang mengetahuinya.

**Dan ada akibat yang lebih tegas:** karena tidak ada tanggal perhitungan yang tersimpan di mana
pun, **nilai yang pernah disetujui tidak bisa direproduksi.** Bukan sulit — tidak bisa. Bila
auditor menanyakan berapa angka yang dilihat direksi saat menyetujui sebuah kontrak, tidak ada
jawaban.

**Kenapa ini mendesak sekarang.** Rencana Treaty In Adjustment mensyaratkan pengambilan nilai lama
untuk menghitung selisih. Selisih hanya bisa dipertanggungjawabkan kalau yang lama tidak berubah.
Selama keempat sebab di atas masih ada, selisih yang sama dihitung dua kali pada waktu berbeda bisa
menghasilkan dua angka.

**Sudah diperiksa sebagian, dan hasilnya melegakan.** Pemeriksaan atas nilai rupiah yang tersimpan
menunjukkan **belum ada kerugian yang terjadi**: angka yang tercatat masih sejalan dengan kurs tahun
kontraknya masing-masing. Artinya kontrak lama selama ini **jarang disentuh ulang** — dan justru
itu sebabnya lubang ini belum pernah menimbulkan akibat.

**Tetapi lubangnya tetap terbuka.** Yang menahannya bukan pengaman, melainkan kebiasaan. Satu
addendum atas kontrak lama sudah cukup untuk memicunya, dan tidak ada yang akan memberi tahu bahwa
itu terjadi.

**Yang dibutuhkan untuk tiga sebab lainnya:** beberapa kueri baca-saja ke basis data, sudah
disiapkan.

---

### Tambahan 23 September 2026, DIRALAT 24 September 2026 — sebagian selisih yang pernah dilaporkan tidak dapat dipertanggungjawabkan ulang

Ditemukan saat merancang sambungan ke modul penyesuaian kontrak, dan **diralat** setelah modul itu
dibedah.

**Ralat.** Sebab yang ditulis semula keliru. Sistem **memang** membekukan keadaan kontrak pada saat
penyesuaian dibuat: seluruh isinya disalin ke dalam catatan penyesuaian itu sendiri. Nilai lamanya
karena itu **dapat** dibaca ulang hari ini, dan butir ini **tidak lagi** menyatakan sebaliknya.

**Yang tetap menjadi masalah, dengan sebab yang berbeda.** Perbandingan nilai lama dengan nilai
baru dapat salah dalam **tiga keadaan**, dan hanya dalam ketiganya:

1. sebuah baris **disisipkan atau dihapus di tengah daftar** — perbandingannya memasangkan baris
   menurut urutan, bukan menurut identitasnya, sehingga seluruh baris sesudahnya bergeser;
2. **mata uang sebuah besaran berbeda** antara keadaan lama dan keadaan baru — pengurangannya
   tidak memeriksanya;
3. **angka ringkasan bagian** yang sempat tampil di layar — ia ditimpa pada penyimpanan
   berikutnya dan tidak pernah tersimpan.

**Di luar ketiga keadaan itu, selisihnya benar.**

**Akibatnya bagi bisnis.** Untuk penyesuaian yang pernah mengalami salah satu dari ketiga keadaan
itu, angka selisih yang pernah **dilaporkan** tidak dapat dipertanggungjawabkan ulang — dan bila
dihitung ulang dengan cara yang benar, yang berbeda adalah **angka yang dulu salah**, bukan angka
yang sekarang bergeser.

**Berapa banyak yang terdampak belum diukur.** Dua pengukuran menjawabnya: berapa penyesuaian yang
daftar barisnya berubah panjang, dan berapa yang mata uangnya berubah — tercatat sebagai **UA-5**
dan **UA-6** pada lampiran pengukuran.

**Kenapa ini menguatkan butir 2, bukan mengulanginya.** Butir 2 menyatakan angka yang disetujui
tidak dijamin tetap. Butir ini menunjukkan **akibat yang sudah terjadi karenanya**: bukan hanya
angkanya dapat bergeser, tetapi seluruh perhitungan selisih yang bersandar padanya kehilangan
dasarnya.

**Di sistem baru.** Penyesuaian menunjuk **satu versi kontrak tertentu**, dan versi yang sudah
disetujui tidak dapat berubah. Selisih karena itu dapat dihitung ulang kapan saja, dan
ketidakcocokan antara yang tersimpan dan yang terhitung menjadi **alat deteksi**.

---

### Tambahan 23 September 2026 — perubahan kontrak yang disetujui tidak ikut ke tabel rinciannya

**Apa yang terjadi.** Ketika sebuah perubahan kontrak (addendum) disetujui, sistem menyimpannya ke
tabel addendum. Langkah yang memperbarui tabel rincian kontrak — tabel datar yang dibaca sistem lain
— ada di dalam sistem, lengkap, **dan dimatikan**. Begitu pula langkah yang memperbarui baris
kontrak induknya.

**Apa yang dipertaruhkan, dan ini bukan soal pelaporan.** Angka yang keluar dari sistem ini menjadi
dasar penawaran ke pihak lawan, pembagian ke retrosesioner, dan pembukuan. Bila rinciannya tidak
pernah menerima addendum, maka selama umur sistem ini **angka yang keluar dibuat di atas syarat yang
sudah digantikan** — syarat yang sudah diubah dan disetujui, tetapi tidak pernah sampai ke tabel
yang dibaca orang lain.

Itu bukan lubang pelaporan. Itu uang dan pihak lawan.

**Seberapa besar — belum diketahui, dan diukur satu kueri.** Pemeriksaan kami bersandar pada
pembacaan kode, dan kode itu dipakai bersama dua modul sehingga ada kemungkinan jalur lain yang
belum kami telusuri. **Uji Z** membandingkan langsung: untuk setiap kontrak yang punya addendum,
apakah nilai di kontrak induknya sama dengan nilai di addendum terakhirnya.

| Hasil | Artinya |
|---|---|
| selalu berbeda | tidak satu pun addendum pernah sampai; angkanya jumlah kontrak terdampak |
| sebagian sama | ada jalur lain yang belum ditemukan, dan ia harus dicari |
| semuanya sama | pembacaan kami keliru, dan butir ini gugur |

**Akibat untuk rencana migrasi, dan ini yang paling mendesak.** Keadaan kontrak yang berlaku **tidak
dapat dibaca dari tabel kontrak**. Ia harus disusun ulang dengan menelusuri rantai addendumnya.
Rencana migrasi mana pun yang memindahkan tabel kontrak apa adanya akan memindahkan kontrak
**tanpa addendumnya** — memindahkan syarat yang salah, dengan rapi, ke sistem baru.

Kabar baiknya: rantai itu **dapat ditelusuri**. Pengenal addendum diturunkan dari pengenal
kontraknya, sehingga induknya selalu dapat dikenali.

**Di sistem baru.** Tidak ada baris "kontrak saat ini" yang terpisah dari versinya; yang berlaku
adalah versi terakhir yang disetujui, dan tidak ada langkah penyalin yang dapat dimatikan.

---

### Tambahan 23 September 2026 — tabel yang dibaca sistem lain tidak punya pengisi yang hidup

Ini temuan yang berdiri sendiri, dan ia yang terberat dari ketiga tambahan di butir ini.

**Apa yang terjadi.** `TREATYINOFFER` adalah tabel yang dibaca sistem-sistem lain untuk mengetahui
syarat sebuah kontrak. Penelusuran menyeluruh atas seluruh jalur yang dapat menulisinya menemukan
**tidak satu pun yang dapat dicapai**:

| Jalur | Keadaan |
|---|---|
| dari penyimpanan kontrak biasa | langkahnya **dimatikan** |
| dari penyimpanan addendum | tidak pernah ada |
| dari sebuah tombol di layar penawaran | tombolnya bersyarat **`NEVER`** — tidak terlihat oleh siapa pun, termasuk pengembang yang namanya tertulis di sebelahnya |

Ketiganya tertutup. **Tidak ada peristiwa bisnis apa pun yang menerbitkan sebuah kontrak ke tabel
itu hari ini.**

**Akibatnya bagi bisnis, dan ini bukan soal pelaporan.** Apa pun yang ada di dalam tabel itu
sekarang ditulis pada suatu masa lampau yang tidak diketahui, dan tidak bertambah maupun berubah
sejak itu. Setiap sistem yang membacanya membaca **potret beku** — bukan keadaan kontrak hari ini,
dan bukan pula keadaan saat kontrak disetujui.

Bila angka yang keluar dari tabel itu dipakai untuk penawaran ke pihak lawan, pembagian ke
retrosesioner, atau pembukuan, maka angka itu **tidak berasal dari kontrak yang berlaku**.

**Yang dibutuhkan, dan ini pertanyaan bisnis bukan pertanyaan teknis:** **siapa saja yang membaca
`TREATYINOFFER`.** Sistem lain, laporan berkala, atau pihak luar — daftarnya menentukan siapa yang
harus diberi tahu. Sistem ini tidak dapat menjawabnya; hanya orang yang dapat.

**Dua akibat untuk rencana migrasi, dinyatakan sekarang supaya tidak ditemukan ulang:**

| | |
|---|---|
| Ia **bukan sumber migrasi** | memindahkannya berarti memindahkan potret beku, bukan data kontrak |
| Ia **bukan dasar rekonsiliasi** | mencocokkan sistem baru terhadapnya akan menghasilkan selisih yang mengukur kapan tabel itu berhenti diisi, bukan kebenaran data |

Keadaan kontrak yang berlaku **tetap dapat disusun ulang** — dari kontrak beserta rantai
addendumnya, yang penelusurannya sudah dipastikan mungkin. Yang gugur hanya tabel penerbitannya.

---

## 3. Tidak ada penjaga otomatis, dan penjaga manusianya bisa dipersingkat

**Apa yang terjadi.** Pemeriksaan kelengkapan yang pernah dibangun di sistem ini — kewajiban
mengisi jenis treaty, kelas bisnis, susunan layer, dan nilai limit — **sudah dimatikan**. Yang
masih berjalan hanya dua: nama cedant tidak boleh kosong, dan asal bisnis tidak boleh kosong.

Pemeriksaan kontrak ganda juga sudah dimatikan seluruhnya. Hari ini tidak ada yang mencegah dua
kontrak dengan cedant, asal bisnis, dan periode yang sama dibuat berdampingan.

Ditambah: jalur revisi memendekkan persetujuan dari empat tingkat menjadi dua, dan pemilihan jalur
itu dilakukan lewat tombol yang ditekan orang yang sama.

**Akibatnya bagi bisnis.** Satu-satunya penjaga yang tersisa adalah mata para approver — dan jalur
yang mereka lalui bisa dipersingkat atas kehendak pengisinya.

**Yang dibutuhkan:** daftar aturan yang dulu dimatikan masih tersimpan di dalam sistem dan bisa
dijalankan ke data sekarang untuk melihat apa saja yang lolos selama ini.

---

## 4. Perubahan pada kontrak tidak meninggalkan jejak

**Apa yang terjadi.** Sistem hanya mencatat siapa dan kapan pada perpindahan persetujuan.
Penyimpanan biasa — termasuk yang mengubah angka uang — tidak meninggalkan catatan apa pun. Tidak
ada kolom tanggal ubah pada tabel kontrak.

**Akibatnya bagi bisnis.** Bila sebuah angka berbeda dari yang diingat orang, tidak ada cara
mengetahui siapa mengubahnya dan kapan. Ini juga yang membuat butir 1 di atas tidak bisa dilacak.

### Tambahan 24 September 2026 — addendum yang ditolak dihapus, bukan ditandai

> **DIPERTAJAM pada hari yang sama, dan ke arah yang lebih berat.** Kalimat *"sistem lama tidak
> pernah mencatat penolakan addendum"* kurang tepat. **Pencatatannya ADA, ditulis orang, lalu
> dimatikan** — empat langkah pertama `TreatyInDeclineConfirmation_postactEDM` (`AddCommentList`,
> `IsApproved = "Decline"`, `StatusAkseptasi = "Decline"`, simpan) seluruhnya bertanda blok mati —
> **sementara dua langkah penghapusannya (`RDB remove EDM`) dibiarkan hidup.**
>
> Bandingkan dengan jalur kontrak dasar, `TreatyInDeclineConfirmation_postact`: **keempat langkahnya
> hidup**, dan penolakan kontrak memang meninggalkan jejak. **Dua jalur, satu dimatikan.**
>
> **Kenapa bedanya penting bagi pembaca butir ini:** *"tidak pernah ada"* mengundang pertanyaan
> *"apakah memang perlu"*. *"Pernah ada, dimatikan, dan yang menghapus dibiarkan"* mengundang
> pertanyaan **kapan, oleh siapa, dan apakah ada yang tahu** — dan itu pertanyaan yang benar.

**Apa yang terjadi.** Menolak sebuah addendum **menghapus barisnya**. Bukan menandainya, bukan
memindahkannya ke arsip — menghapus. Aktivitas `TreatyInDeclineConfirmation_postactEDM` langkah 6
dan 7, **keduanya berjalan**, menghapus baris di `M_TREATY_IN_EDM` dan `TREATY_IN_EDM`.

Jalur kontrak biasa **tidak** begitu: penolakan kontrak hanya menyetel keadaannya, dan barisnya
bertahan. **Dua jalur, dua perlakuan, tanpa alasan bisnis yang masuk akal untuk membedakannya.**

**Akibatnya bagi bisnis — tiga, dan yang ketiga paling mengena:**

1. **Tidak ada catatan penolakan addendum yang pernah ada.** Bukan catatannya tidak lengkap; tidak
   ada sama sekali. Setiap addendum yang pernah ditolak sepanjang umur sistem sudah lenyap tanpa
   jejak.
2. **Nomor revisi mungkin dipakai ulang.** Pengenal addendum berbentuk `‹kontrak›/Rnn` dan nomornya
   dihitung dari baris yang **ada**. Bila `/R02` dihapus, addendum berikutnya dapat memperoleh
   `/R02` lagi — sehingga dua addendum berbeda pernah memakai pengenal yang sama. Setiap dokumen di
   luar sistem — surat, slip, berkas akuntansi — yang menyebut `‹kontrak›/R02` menjadi **tidak
   tertentu**.
3. **Tidak ada seorang pun yang pernah dapat menjawab berapa kali sebuah kontrak gagal diubah.**
   Itu pertanyaan yang biasanya ditanyakan auditor, dan datanya tidak pernah ada untuk menjawabnya.

**Yang diminta dari manajemen:** tidak ada keputusan — ini pemberitahuan. Pengukurannya **Uji AC**
di berkas permintaan DBA, dan ia hanya menunggu izin kueri baca-saja yang sudah diminta di butir 2.

> **Sifat pembuktiannya perlu diketahui sebelum hasilnya dibaca.** Uji AC dapat **mematahkan**
> dugaan "tidak ada yang pernah dihapus", tetapi tidak dapat **mengesahkannya**: baris yang dihapus
> tidak meninggalkan apa pun, sehingga "tidak pernah ditolak" dan "ditolak lalu nomornya dipakai
> ulang" menghasilkan deret nomor yang persis sama.

**Di sistem baru** keduanya tidak terbawa: baris tidak pernah dihapus (ADR-0055 perubahan 24 Sep
2026), nomor versi unik per kontrak (INV-04), dan penolakan addendum **meninggalkan catatan** —
sebuah kemampuan **baru**, bukan pelestarian.

---

## 5. Batas wewenang persetujuan tidak dibaca sistem

**Apa yang terjadi.** Mesin persetujuan menentukan ke siapa sebuah kontrak dikirim **tanpa membaca
satu pun besaran** — tidak membaca limit, tidak membaca premi, tidak membaca bagian NuRe.

**Akibatnya bagi bisnis.** Bila NuRe punya dokumen kebijakan yang menyatakan batas wewenang
berjenjang menurut besaran — misalnya kontrak di atas nilai tertentu harus naik ke direksi — maka
batas itu tidak pernah ditegakkan oleh sistem selama ini.

**Ini bukan temuan migrasi.** Bila dokumen batas wewenang itu ada, ini temuan kepatuhan yang berdiri
sendiri dan tidak hilang bila proyek dibatalkan.

**Yang dibutuhkan:** salinan dokumen batas wewenang yang berlaku. Bukan penjelasan lisan dari
ingatan — dokumennya.

---

## 6. Sebagian perhitungan diketahui salah, dan arah kesalahannya menghilangkan nilai

Beberapa kesalahan perhitungan ditemukan dan sudah dipastikan dari sumbernya. Yang perlu diketahui
manajemen bukan rinciannya, melainkan **arahnya**:

| Yang terjadi | Arah kesalahannya |
|---|---|
| Premi dalam mata uang asing yang kursnya tidak ditemukan | nilainya menjadi **nol** — lenyap dari total |
| Potongan pada perubahan kontrak di tengah periode | potongan tidak ikut menyusut, sehingga **pendapatan bersih tercatat terlalu rendah**, dan pada kasus tertentu bisa menjadi negatif tanpa ada pengembalian premi yang sungguh terjadi |
| Rate on Line pada layer bermata uang lebih dari satu | **terbagi** sebanyak jumlah mata uangnya — dan angkanya **tampak wajar**, sehingga tidak ada yang akan mencurigainya |
| Ketentuan reinstatement yang tarifnya bukan 100% | tarif kontraknya dapat **terhapus** oleh perhitungan |
| Persentase bagian fakultatif yang tampil di layar | diambil dari **salinan** yang dibuat saat penyebaran disusun; bila angkanya diubah sesudah itu, layar tetap menampilkan yang lama — **angka yang dibaca orang bukan angka yang berlaku** |

**Ada kelas ketiga, dan ia yang paling sunyi.** Sebagian kesalahan tidak menggelembungkan dan tidak
menghilangkan — ia **saling meniadakan**. Angka pencapaian pada kontrak proporsional dihitung
terhadap bagian NuRe tingkat kontrak meskipun tiap baris punya bagiannya sendiri; baris yang
bagiannya lebih kecil tercatat terlalu rendah, yang lebih besar tercatat terlalu tinggi.

Akibatnya: **totalnya tampak wajar.** Siapa pun yang memeriksa jumlah tidak akan menemukan apa pun,
sementara setiap barisnya salah.

> Kesalahan yang saling meniadakan di agregat lebih sulit ditemukan, dan karena itu **lebih lama
> bertahan**. Ia tidak pernah memicu pertanyaan dari siapa pun yang membaca total.

Pemeriksaannya sudah disesuaikan: ia membandingkan **per baris**, bukan per total, dan melaporkan
sebaran simpangan ke dua arah. Satu hitungan murah mendahuluinya dan menentukan apakah persoalan ini
terjangkau sama sekali — berapa kontrak yang baris-barisnya benar-benar punya bagian berbeda. Bila
nol, ia risiko yang belum terjadi, bukan kerusakan yang sedang berjalan.

**Satu di antaranya bukan salah hitung melainkan salah tampil, dan itu jenis tersendiri.** Persentase bagian fakultatif disalin ke tempat kedua ketika penyebaran disusun, dan layar membaca salinan itu — sementara yang disunting orang adalah aslinya. Angkanya tidak pernah salah dihitung; ia hanya **sudah tidak berlaku**. Siapa pun yang mengambil keputusan dari layar itu mengambilnya dari angka lama tanpa tanda apa pun bahwa ia lama.

**Kenapa arahnya penting.** Kesalahan yang **menggelembungkan** angka mengundang pertanyaan;
kesalahan yang **menghilangkan** angka tidak. Empat dari temuan di atas menghilangkan, bukan
menggelembungkan — karena itu keempatnya jauh lebih mungkin sudah berlangsung lama tanpa ada yang
menyadarinya.

**Sebagian sudah diperiksa.** Dua dari empat sudah dipastikan **belum menyentuh data**: nilai rupiah
yang tersimpan masih sejalan dengan kurs tahun kontraknya, dan tidak ada kontrak yang menyimpan
angka penanda kegagalan perhitungan.

**Dua sisanya belum diperiksa, dan keduanya justru yang paling sulit terlihat** — karena keduanya
menghasilkan angka yang tampak wajar, bukan angka yang mencolok. Potongan yang tidak ikut menyusut
dan Rate on Line yang terbagi tidak meninggalkan jejak apa pun yang bisa dicurigai orang.

**Yang dibutuhkan:** menjalankan perhitungan yang benar atas data yang ada dan melaporkan
selisihnya dalam rupiah, per tahun. Itu bisa dilakukan tanpa mengubah apa pun, dan menghasilkan
angka, bukan perkiraan.

---

## 7. Ada tabel acuan penting yang tidak punya pemilik

**Apa yang terjadi.** Pembagian kapasitas NuRe diambil dari sebuah tabel yang **tidak pernah
ditulis oleh sistem Treaty In**. Siapa yang memeliharanya, dengan dasar apa, dan dengan persetujuan
siapa — tidak diketahui dari sistem.

Tabel itu juga tidak punya kunci, tidak punya indeks, dan seluruh kolomnya bertipe teks bebas.
Perhitungan kapasitas mempercayainya sepenuhnya: bila angka di dalamnya tidak berjumlah seratus
persen, hasilnya tetap dipakai dan tidak ada yang menangkapnya.

**Akibatnya bagi bisnis.** Angka yang menentukan berapa kapasitas NuRe tersebar ke mana, berasal
dari sumber yang tidak bertuan.

---

## Butir 9 — SISTEM MENERBITKAN KUNCI MASUK YANG DAPAT DITEBAK

**Ditambahkan 24 September 2026. Ini bukan utang migrasi — ia terbuka sekarang, di sistem yang
sedang berjalan.**

**Apa yang terjadi.** Untuk menyimpan dan mengambil berkas lampiran, sistem menerbitkan sebuah kunci
masuk sementara. Kunci itu dibentuk dari **satu kata tetap yang sama setiap kali, ditambah waktu
saat itu** — tanpa satu pun unsur acak. Siapa pun yang mengetahui cara pembentukannya dapat membuat
kunci yang **sah** untuk detik mana pun, tanpa pernah meminta izin kepada sistem.

Masa berlakunya satu menit. Itu memperkecil jendelanya, **tidak** menutup celahnya: yang dapat
menebak kunci untuk satu detik dapat menebaknya untuk setiap detik.

**Akibatnya bagi bisnis.** Akses ke penyimpanan berkas lampiran tidak benar-benar dijaga oleh kunci
itu. Berkas lampiran kontrak memuat slip, korespondensi, dan lampiran bernilai komersial.

**Ini bukan milik Treaty In.** Penyimpanan itu dipakai lebih dari satu modul; yang menemukannya
kebetulan sesi ini. **Yang diminta: penetapan siapa pemiliknya, lalu keputusannya.**

---

## Butir 10 — PERKIRAAN PREMI: BELUM ADA YANG MEMASTIKAN ANGKANYA UNTUK SIAPA

**Ditambahkan 24 September 2026.** Butir ini **sudah ditetapkan sebagai eskalasi** oleh sesi to-spec
pada 23 September dan **tidak pernah sampai ke daftar ini** — kekeliruan pencatatan, bukan
keputusan.

**Apa yang terjadi.** Perkiraan pendapatan premi yang disepakati saat penawaran dibawa apa adanya
menjadi dasar perhitungan premi kontrak. **Belum ada yang memastikan apakah angka itu berarti
seluruh treaty, atau sudah bagian NuRe saja.**

Selama ini rumus pencapaian mengandaikan yang pertama. Andaian itu tidak pernah ditulis di mana pun
dan tidak pernah dikonfirmasi ke bagian teknik.

**Akibatnya bagi bisnis.** Bila andaian itu terbalik, yang harus diperbaiki **bukan hanya cara
menghitung ke depan** — **angka pencapaian yang sudah dilaporkan harus disajikan ulang.** Besarnya
tidak diketahui sampai pertanyaannya dijawab.

**Yang diminta.** Satu jawaban dari bagian teknik treaty; pertanyaannya satu kalimat dan sudah
disiapkan sebagai butir `T-4` di `2-to-spec/PERMINTAAN-TEKNIK-TREATY.md`. Butir ini dan `T-4`
menanyakan hal yang sama kepada orang yang berbeda: `T-4` meminta **jawabannya**, butir ini memberi
tahu manajemen **apa yang bergantung padanya**. Jawaban lewat salah satu jalur menutup keduanya.

---

## Ringkasan yang dibutuhkan dari manajemen

| Butir | Yang diminta |
|---|---|
| 1 | tidak ada — pemeriksaannya sudah berjalan, hasilnya dilaporkan menyusul |
| 2, 6 | izin menjalankan kueri baca-saja ke basis data produksi |
| 3 | pertemuan singkat dengan bagian teknik untuk menilai daftar aturan yang dulu dimatikan |
| 5 | salinan dokumen batas wewenang persetujuan yang berlaku |
| 7 | penetapan siapa pemilik tabel acuan pembagian kapasitas |
| 3 | **daftar siapa saja yang membaca `TREATYINOFFER`** — sistem lain, laporan, atau pihak luar |
| **1** | **jawaban bisnis:** apakah perubahan atas kontrak yang sudah disetujui Direktur memang boleh disahkan kembali oleh Kepala Seksi saja |
| **9** | **penetapan siapa pemilik penyimpanan berkas lampiran**, lalu keputusan atas cara penerbitan kunci masuknya |
| **10** | **jawaban bagian teknik treaty:** perkiraan pendapatan premi dicatat untuk seluruh treaty atau untuk bagian NuRe |
| **8** | penetapan tata kelola atas penulisan aturan langsung di lingkungan produksi |
| **9** | **jawaban dua kalimat:** apakah keempat nama yang tertulis di penyaluran persetujuan masih memegang peran itu, dan apakah penyaluran memang dimaksudkan menyebut orang alih-alih jabatan |

---

## Lampiran — pengukuran yang belum dijalankan

**Diperbarui 23 September 2026.** Pemilik proses memutuskan seluruh butir yang tertahan diputuskan
tanpa menunggu pemeriksaan data, dan **pembangunan sistem baru berjalan penuh**. Keputusan itu
tidak mengubah butir 2 sampai 7 di atas.

Yang berubah adalah kedudukan pemeriksaan datanya. Pemeriksaan di bawah **tidak lagi menghambat
pembangunan**. Ia hanya dibutuhkan untuk menjawab satu pertanyaan, dan itu pertanyaan manajemen,
bukan pertanyaan teknis:

> **Berapa banyak kontrak yang terlanjur salah, dan berapa rupiah nilainya?**

| Yang diukur | Untuk butir |
|---|---|
| Berapa nilai rupiah menyimpang akibat kurs yang dipilih menurut tahun perhitungan | 2, 6 |
| Berapa kontrak yang pembagian kapasitasnya berubah setelah disetujui | 2 |
| Berapa kontrak yang bagian NuRe-nya berubah, lalu dipakai menghitung premi lama | 2 |
| Berapa kontrak yang tersimpan tanpa kelengkapan yang dulu diwajibkan | 3 |
| Berapa kelompok kontrak yang berbagi cedant, asal bisnis, dan periode yang sama | 3 |
| Berapa layer yang Rate on Line-nya terbagi jumlah mata uang | 6 |
| Berapa kontrak yang potongannya tidak ikut menyusut pada perubahan di tengah periode | 6 |
| Berapa kontrak yang pendapatan bersihnya tercatat negatif tanpa pengembalian premi yang nyata | 6 |
| Berapa ketentuan reinstatement yang tarif kontraknya terhapus perhitungan | 6 |
| Berapa susunan kapasitas yang persentasenya tidak berjumlah seratus | 7 |
| Berapa kontrak yang baris rinciannya bertambah setiap kali kontraknya disimpan ulang, alih-alih diperbarui | **baru 24 Sep** |
| Berapa kontrak yang persentase bagian NuRe dan brokerage-nya kosong di tabel yang dibaca sistem lain | **baru 24 Sep** |
| Berapa addendum yang nilai aktualnya tidak pernah terisi sama sekali | **baru 24 Sep** |
| Berapa panjang teks terpanjang pada pengecualian, ketentuan khusus, keterangan, catatan, dan catatan bordereaux — untuk menentukan lebar kolomnya tanpa memotong | **baru 24 Sep** |
| Berapa kontrak yang berstatus selesai disetujui tanpa jejak pemberi persetujuan | 1 |
| Berapa baris yang angka pencapaiannya menyimpang, dipilah menurut arah simpangannya | 6 |
| Berapa kontrak yang tabel rinciannya tidak memuat addendum yang sudah disetujui | 3 |

**Kenapa ini tetap perlu meski pembangunan jalan terus.** Sistem baru memperbaiki mekanismenya ke
depan. Ia **tidak** memperbaiki angka yang sudah tercatat, dan memang tidak boleh —
mengubah angka yang pernah disetujui adalah persis persoalan di butir 1.

Jadi bila ada kerugian yang terlanjur terjadi, ia akan tetap ada di pembukuan sampai seseorang
memutuskan apa yang dilakukan terhadapnya. **Keputusan itu milik manajemen, dan tidak bisa diambil
tanpa angkanya.**

**Ongkosnya rendah:** seluruhnya kueri baca-saja, tidak mengubah apa pun, dan sudah tersusun di
`PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql`. Yang dibutuhkan hanya izin menjalankannya.

---

## Butir 8 — ATURAN DITULIS LANGSUNG DI LINGKUNGAN PRODUKSI

**Ditambahkan 24 September 2026.** Ditemukan lewat sapuan `pxCreateSystemID` yang semula dirancang
untuk pertanyaan lain.

### Bentuknya: KEBIASAAN dengan satu hari padat — bukan satu insiden

Angka mentahnya 32; sesudah dihilangkan salinan antarfolder, **17 aturan berbeda**.

| | |
|---|---|
| Pembuat | **tiga orang** |
| Rentang | **2 Juni 2025 – 3 September 2026** — lima belas bulan |
| Hari terpadat | **2 Juni 2025: delapan aturan, tiga orang, satu hari** |
| Tanggal lain | 5 Jun · 9 Jun · 7 Jul · 16 Sep 2025 · 3 Sep 2026 |

**Pengelompokannya memutuskan kalimat eskalasinya.** Bila seluruhnya terkumpul pada satu-dua orang
dalam beberapa hari, ia **insiden** — ada yang mendesak, dan pertanyaannya *"apa yang terjadi waktu
itu"*. Yang terbaca **bukan** itu: delapan aturan pada satu hari **memang** berbentuk insiden atau
saat rilis, tetapi **lima tanggal berikutnya tersebar lima belas bulan oleh dua orang yang sama**.

> **Maka ini soal TATA KELOLA, bukan satu kejadian.** Menulis aturan langsung di produksi adalah
> cara kerja yang berjalan, bukan pengecualian yang pernah terjadi.

### Dan ya — sebagian berada di jalur HIDUP

Pertanyaan ini berbeda dari *"apakah temuan kita bersandar padanya"*, dan jawabannya **ya**:

| Aturan | Langkah | Mati | Perujuk | Terjangkau dari |
|---|---:|---:|---:|---|
| `TreatyInSetBrokerage` | 54 | 2 | 9 | `TreatyInNONProportional.xml` — **layar** |
| `TreatyInNonAddItem` | 41 | 3 | 9 | `TreatyInNONProportional.xml` — **layar** |
| `AddSpreadingXOL` | 4 | 0 | 1 | `Share.xml` — **layar** |

> **Kode yang paling sedikit ditinjau adalah kode yang berjalan.** Ketiganya hidup, terjangkau dari
> layar, dan ditulis di lingkungan tempat tidak ada tinjauan yang memaksa.

### Penajaman kedua — apa yang SEBENARNYA disunting di produksi

Enam belas berkas terakhir disimpan di produksi, dan **bentuknya memberi keterangan yang mengubah
berat butir ini**:

| Golongan | Berkas | Risikonya |
|---|---|---|
| **pelaporan dan pencarian** | `Browse*_RD` (4), `CountDocTreaty_SQL`, `GetMasterTreatyCategory_SQL`, `GetNilaiTotal`, `TreatyInSummaryLimit`, `TreatyInSummaryLimitShare`, `BrowseTreatyOutDetailEDM`, `Layers`, `ShowAttachmentTreaty`, `InputTreatyInAdjustment`, `DetailCalculationROL` | **terbatas** — laporan yang salah **menyesatkan orang**; ia tidak mengubah angka yang tersimpan |
| **LOGIKA TRANSAKSI** | **`TreatyInNonAddItem`** · **`TreatyInSetBrokerage`** · `AddSpreadingXOL` | **nyata** — ketiganya menulis nilai kontrak, hidup, dan terjangkau dari layar |

> **Kebiasaan menyunting di produksi ada, dan ia TERUTAMA menyentuh pelaporan.** Tiga pengecualian
> menyentuh logika transaksi, dan ketiganya disebut dengan nama di atas.

**Kalimat ini sengaja tidak digeneralisasi.** Butir yang berbunyi *"tujuh belas aturan disunting di
produksi"* tanpa membedakan keduanya **akan dibaca lebih berat daripada keadaannya**, dan pembaca
yang memeriksanya sendiri akan menganggap seluruh butir ini dilebih-lebihkan. Yang menyebut bedanya
**akan dipercaya pada butir berikutnya**.

### KOREKSI atas laporan sebelumnya

Saya menulis *"tidak satu pun temuan kita bersandar pada aturan yang dibuat di produksi"*. **Itu
salah.** `TreatyInNonAddItem` dan `TreatyInSetBrokerage` adalah **dua dari bukti 5a** — keduanya
dipakai untuk membuktikan cabang mana yang menulis `Limits[]` dan `Share[]`.

**Arahnya justru menguntungkan, dan itu perlu disebut supaya tidak dibaca terbalik:** aturan yang
**dibuat di produksi** menurut definisinya **ada di produksi**. Ia satu-satunya golongan yang
**L-9 tidak dapat ancam**. Yang tetap terpapar L-9 adalah sisi proporsional bukti 5a —
`TreatyInMappingDataconvertProp` dan `TreatyInPropAdd`, keduanya dibuat di `pega`/dev.

### Yang diminta

| Kepada | Apa |
|---|---|
| tata kelola TI | apakah penulisan aturan langsung di produksi memang dibolehkan, dan bila ya, dengan tinjauan apa |
| ketiga pembuat | apakah ketujuh belas aturan itu sudah dirapikan kembali ke jalur pengembangan, atau masih hanya ada di produksi |
| pemilik proses | ketiganya hidup dan terjangkau layar — apakah perilakunya sudah pernah ditinjau siapa pun |

**Penagihnya:** permintaan **ekspor kedua** (`LUBANG-SPESIFIKASI.md` §8). Bila ekspor produksi
diperoleh dan ketujuh belas aturan ini **tidak ada** di ekspor pengembangan, itu bukti langsung
bahwa keduanya berbeda — dan cakupan L-9 terukur seketika.

### Penajaman 24 September 2026 — `pxUpdateSystemID` menutup separuh yang kurang

`pxCreateSystemID` membuktikan sebuah aturan **pernah ada** di produksi saat dibuat. Ia **tidak**
membuktikan bahwa salinan yang kita pegang sama dengan yang berjalan di sana sekarang — aturan dapat
disunting di tempat lain sesudahnya. Kolom kedua menutupnya.

**Dari ketujuh belas aturan yang dibuat di produksi:**

| | |
|---|---:|
| dibuat **dan** terakhir disimpan di `pegaprdnusare` | **13** |
| terakhir disimpan di `pegadevnusare2` | **4** — `GetEstimasiClaim_act`, `GetHistoryClaim_act`, `GetMasterTreatyCategory_Act`, `SetkategoriDoc` |

> **Kedua aturan yang menyangga bukti 5a — `TreatyInNonAddItem` dan `TreatyInSetBrokerage` —
> dibuat DAN terakhir disimpan di produksi.** Pembalikan itu berdiri **penuh** untuk keduanya:
> keduanya keluar dari daftar terdampak L-9.

**Sebaran atas seluruh 708, dan ia memberi bacaan KETIGA yang belum kami sebut:**

| Dibuat → terakhir disimpan | Berkas |
|---|---:|
| `pega` → `pega` | 442 |
| `pegadevnusare2` → `pegadevnusare2` | 206 |
| `pegaprdnusare` → `pegaprdnusare` | 24 |
| `pega` → `pegadevnusare2` | 17 |
| `pegaprdnusare` → `pegadevnusare2` | 8 |
| `pega` → `pegaprdnusare` | 5 |

**29 berkas terakhir disimpan di produksi, 231 di pengembangan.** Bila ekspor ini murni dari
pengembangan, ke-29 itu hanya dapat hadir bila perubahan produksi pernah ditarik kembali ke
pengembangan.

> **Bacaan ketiga: ekspornya CAMPURAN, dan "ekspor ini dari lingkungan mana" mungkin tidak punya
> satu jawaban.** Yang punya jawaban adalah pertanyaan per berkas — dan `pxUpdateSystemID` sudah
> menjawabnya untuk ketujuh belas yang paling berakibat.

Ini **petunjuk terkuat yang dapat diperoleh dari ekspor**, dan ia tetap **bukan bukti**.

---


---

### Penajaman ketiga 24 September 2026 — "terjangkau dari layar" DIVERIFIKASI, dan ketiganya bertahan

Kalimat *"ketiganya hidup serta terjangkau dari layar"* menyangga penyebutan tiga aturan logika
transaksi di atas, dan sampai hari ini ia berdiri di atas pemeriksaan yang **tidak pernah membaca
penanda pematian kendali layar** (`pyDisabled`). Kekurangan itu ditutup.

| Aturan | kendali hidup | mati bersyarat | mati tanpa syarat |
|---|---:|---:|---:|
| `AddSpreadingXOL` | 8 | 0 | **0** |
| `TreatyInNonAddItem` | 132 | 42 | **0** |
| `TreatyInSetBrokerage` | 35 | 28 | **7** — seluruhnya di layar cermin/nilai lama, dan memang seharusnya mati |

> **Tidak ada kalimat yang dicabut.** Butir ini tetap menyebut **tiga** aturan logika transaksi.
> Pemeriksaannya: `4-erd-dan-tabel-datar/CABANG-MATI-DI-JALUR-PERSETUJUAN.md` §4.

---

## Butir 9 — ATURAN YANG MENYEBUT INDIVIDU ALIH-ALIH PERAN

**Empat instans, satu cara kerja.** Didaftar bersama dengan sengaja: empat butir terpisah akan
dibaca sebagai empat kekhilafan kecil, sementara satu butir berdaftar empat dibaca sebagai apa
adanya — **cara kerja yang menempatkan orang, dan baris data tertentu, di dalam kode.**

| | Instans | Bentuknya |
|---|---|---|
| **a** | tombol yang hanya muncul bagi `OperatorID.pyUserName = 'ALDO SAPUTRA'` | satu **nama orang** sebagai syarat tampil |
| **b** | dua nomor kontrak, `1000951` dan `1000069`, di **enam belas** cabang yang hidup | dua **baris data** sebagai syarat jalan |
| **c** | **empat nama orang** di penyaluran persetujuan (`Akseptasi_DT.xml`) | empat **nama orang** sebagai nilai yang ditulis |
| **d** | kode tim disimpan di ruas **nomor telepon** (`OperatorID.pyTelephone`) | satu **ruas berarti dua hal**, dan arti keduanya menentukan penyaluran |

### Instans (c) dan (d) bekerja berpasangan — dan itu yang membuatnya berakibat

```
1.1  pyTelephone == "SPVTREATY1" / "TREATY1"  ->  penerima berikutnya: "IRVANDY"
1.1  pyTelephone == "SPVTREATY2" / "TREATY2"  ->  penerima berikutnya: "AGUNGPUTRAANDALAS"
1.2  Kepala Seksi menyetujui                  ->  penerima berikutnya: "YOHANESKRISTIAWAN"
1.3  Kepala Departemen menyetujui             ->  penerima berikutnya: "NANDINA"
```

Dibuat **27 Juni 2019**, di sistem produksi, dan **hidup** — tujuh dari delapan penugasan nama
berjalan.

### Tiga akibat yang dirasakan orang, dan ketiganya belum pernah dinyatakan

| | Akibat |
|---|---|
| **a** | **Penyaluran putus diam-diam ketika orangnya pergi.** Bila salah satu dari keempat nama itu berpindah tugas atau berhenti, aturannya tetap menulis nama mereka. **Tidak ada galat**; yang terjadi adalah kontrak disalurkan kepada orang yang tidak lagi memegang peran itu. Memperbaikinya menuntut **mengubah aturan**, bukan mengubah data |
| **b** | **Orang baru tidak dapat menerima penyaluran** sampai `pyTelephone`-nya disetel ke salah satu dari empat kode tim. Itu **langkah penyiapan pengguna yang tidak tertulis di mana pun**, dan yang mengerjakannya harus tahu bahwa ruas telepon bukan ruas telepon |
| **c** | **Layar menyebut nama yang tetap.** Pengaju melihat *"berikutnya: NANDINA"* apa pun keadaan sebenarnya — dan bila (a) terjadi, **layar berbohong tanpa ada yang menyadarinya** |

Dan satu akibat tambahan dari (d) sendiri: siapa pun yang menyunting ruas telepon seorang pengguna —
pekerjaan administratif yang tampak tidak berbahaya — **mengubah ke siapa kontraknya disalurkan.**

### Yang ditanyakan ke manajemen, dua kalimat

> **Apakah keempat nama itu masih memegang peran tersebut hari ini?**
> **Dan apakah penyaluran persetujuan memang dimaksudkan menyebut orang, bukan jabatan?**

**Jawaban apa pun berguna.** *"Masih"* berarti sistem lama kebetulan masih benar, dan yang tersisa
adalah risiko ke depan. *"Tidak lagi"* berarti ada kontrak yang **sekarang** disalurkan kepada orang
yang salah, dan jumlahnya dapat dihitung.

**Akibatnya pada sistem baru sudah diputuskan dan tidak menunggu jawaban ini:** penyaluran dibaca
dari **peran bertanggal** (ADR-0044), dan nama orang tidak disimpan di kontrak
(`SPEC-MODEL-DATA.md` §12.5). Butir ini naik untuk **sistem yang berjalan sekarang**, bukan untuk
rancangan.
