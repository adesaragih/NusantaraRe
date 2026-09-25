# Ringkasan untuk Komite Pengarah - Nusantara Re

> Migrasi Pega ke Go + React. Keadaan per **25 September 2026**.
> Dokumen ini memuat keputusan dan angka, bukan cerita. **Nol rekomendasi** pada bagian 5 -
> pilihan dan akibatnya disajikan, penentuannya milik komite.

---

# 1. Keadaan dalam satu halaman

| Ukuran | Angka |
| --- | ---: |
| Modul yang sudah punya spesifikasi **dan** tiket | **20 dari 20** |
| Bounded context | **9** |
| Tiket seluruhnya | **354** |
| Tiket siap dikerjakan | **179** |
| Tiket tertahan | **52** |
| Keputusan arsitektur | **105**, dalam **tiga seri** yang masing-masing mulai dari nomor satu |
| Rangkaian kerja penyusun bahan | **3** |

[~] **Angka tiket di atas hasil pengukuran ulang.** Catatan audit sebelumnya menyebut **365**;
selisih **11** seluruhnya **berkas indeks** yang ikut terhitung. Yang dipakai di sini **berkas
tiket sebenarnya**.

**Tiga hal yang perlu diketahui komite dalam satu kalimat masing-masing:**

1. Seluruh modul sudah bertiket, dan sebagian besar dapat mulai dikerjakan sekarang.
2. Yang tertahan hampir seluruhnya di **Fakultatif**, dan penahannya **bahan dari pihak luar**,
   bukan pekerjaan yang belum dirumuskan.
3. Ada **lima keputusan** yang hanya komite dapat ambil, dan dua di antaranya menyangkut
   **satu tabel fisik dengan dua rancangan kolom yang saling bertentangan**.

---

# 2. Yang siap dibangun sekarang

| Bounded context | Tiket | Siap | Tertahan | Catatan |
| --- | ---: | ---: | ---: | --- |
| Klaim Jiwa | 15 | **15** | 0 | tanpa penahan |
| Penawaran Jiwa | 21 | **18** | 0 | tanpa penahan |
| Data Induk Jiwa | 22 | **22** | 0 | tanpa penahan |
| Klaim Non-Jiwa | 68 | **29** | 0 | tanpa penahan |
| Komite Klaim | 50 | **35** | 0 | tanpa penahan |
| Realisasi Treaty | 40 | **24** | 16 | sebagian menunggu butir terbuka |
| Kontrak Treaty | 66 | **0** | 0 | tanpa penahan |
| Penempatan Treaty | 12 | **12** | 0 | tanpa penahan |
| Fakultatif | 60 | **24** | 36 | penahannya bahan dari luar |
| **Jumlah** | **354** | **179** | **52** | |

[~] **Sebagian tiket tidak membawa baris status yang dapat dibaca mesin** - terhitung
**120** tiket. Ia bukan berarti tertahan; ia berarti statusnya ditulis dengan bentuk lain,
dan perlu dibaca satu per satu sebelum dijadwalkan.

---

# 3. Urutan pembangunan yang diusulkan

Urutan ini mengikuti **gerbang antar tiket** yang sudah tertulis di tiketnya, bukan preferensi.

| Tahap | Yang dikerjakan | Kenapa di sini |
| ---: | --- | --- |
| **1** | kerangka penyimpanan polis dan kunci generasi | seluruh tiket penyimpanan lain menulis ke kerangka ini |
| **2** | nomor urut baris anak, dan tipe kolom | memasangkan baris antar generasi bergantung padanya |
| **3** | pemecah dokumen menjadi baris | menghasilkan data yang dipakai seluruh modul treaty |
| **4** | transaksi tunggal dan skema eksplisit | penjaga keutuhan, dipasang sebelum jalur tulis bercabang |
| **5** | modul yang tidak menunggu siapa pun - Jiwa, Komite, Penempatan Treaty | dapat berjalan sejajar dengan tahap 1-4 |
| **6** | rantai generasi endorsemen dan perhitungan selisih | menuntut tahap 1-3 selesai |
| **7** | pemuat dokumen lama | menuntut lingkup pemindahan diputuskan lebih dulu |
| **8** | Fakultatif | menuntut bahan dari DBA, Product, dan pemilik ekspor |

[!] **Gerbang yang paling menentukan jadwal:** tahap 7 dan 8 **tidak dapat dijadwalkan** sampai
butir pada bagian 4 ditutup. Tahap 1-6 dapat dijadwalkan hari ini.

---

# 4. Butir terbuka yang memblokir pembangunan

| Pemilik | Yang ditunggu | Jumlah | Menahan |
| --- | --- | ---: | --- |
| **DBA** | isi stored procedure jalur tulis produksi | **35** | seluruh jalur tulis Fakultatif |
| **DBA** | definisi tabel fisik | **nol tersedia** | penyelarasan presisi di seluruh modul |
| **DBA** | presisi kolom uang - lebar digit di depan koma belum diuji terhadap nilai terbesar | 1 butir | tiket tipe kolom di dua rangkaian |
| **DBA** | panduan bentuk dokumen lama, yang ada terbukti basi | 1 butir | daftar kolom lengkap |
| **Product** | daftar nilai untuk daftar-pilihan layar | **604** | layar Fakultatif |
| **Product & UW** | arti kode lini bisnis, batas endorse, tarif pajak menurut waktu | beberapa | validasi dan perhitungan |
| **pemilik ekspor** | aturan yang perlu diekspor ulang | **77** | perilaku yang belum terbaca |
| **pemilik ekspor** | baris tabel keputusan underwriting | tidak terekspor | keputusan underwriting otomatis |
| **pemilik pekerjaan** | lingkup pemindahan dokumen lama | 1 butir | pemuat migrasi di dua rangkaian |
| **pemilik pekerjaan** | rancangan kolom tabel akar bersama | 1 butir | penyimpanan lintas lini |
| **Finance** | perlakuan atas galat angka yang sudah tersimpan | 1 butir | ambang uji paritas |
| **IAM** | pemetaan nama orang menjadi peran | **12 tempat** | jalur wewenang |

**Kuesioner Fakultatif** memuat **27 pertanyaan** yang belum terjawab, dialamatkan ke
Underwriting, Product, IT, dan DBA.

---

# 5. Keputusan yang diminta komite

[AWAS] **Tiap butir disajikan dengan pilihan dan akibatnya. Tidak ada rekomendasi.**

## 5a. Satu tabel fisik, dua rancangan kolom

Tabel akar dan tabel generasi dirancang **dua kali**, oleh dua rangkaian kerja:

| Sisi | Rancangan | Dasarnya |
| --- | --- | --- |
| Treaty | **79 medan** tingkat atas ditambah kolom keterangan | lembar rancangan tabel proyek |
| Fakultatif | **71 kolom**, diturunkan dari **113** contoh dokumen | penelusuran dokumen produksi |

Sisi Fakultatif menyatakan rekonsiliasi antar kedua rancangan **gugur**, dan menerima irisan
kolom serta tabrakan ruang kunci sebagai **risiko sadar**.

**Yang diminta:** siapa **pemilik rancangan tabel** itu.

| Pilihan | Akibat |
| --- | --- |
| satu pihak ditetapkan sebagai pemilik, sisi lain menyesuaikan | satu sisi mengerjakan ulang sebagian rancangan dan tiketnya |
| kedua rancangan disatukan lewat rekonsiliasi | menunda kedua sisi sampai rekonsiliasi selesai |
| dua tabel fisik terpisah, satu per lini | menggugurkan keputusan arsitektur tabel akar lintas-lini |

## 5b. Kolom pembeda lini

Satu keputusan Fakultatif menyatakan kolom pembeda lini **di luar proyek** dan mencabutnya dari
rancangan. Satu keputusan arsitektur menyatakan tabel kerja adalah **akar lintas-lini** - yang
mengandaikan ada yang membedakan barisnya.

**Akibat bila dibiarkan:** satu tabel akan memuat baris Fakultatif dan baris Treaty **tanpa
kolom yang membedakannya**.

| Pilihan | Akibat |
| --- | --- |
| kolom pembeda dikembalikan | rancangan Fakultatif berubah; tiketnya disesuaikan |
| pencabutan dipertahankan | pembedaan harus diturunkan dari kolom lain; pembaca SQL wajib tahu caranya |
| tabel dipisah per lini | menggugurkan keputusan arsitektur tabel akar lintas-lini |

## 5c. Penamaan tabel untuk properti yang sama

| Properti sistem lama | Rangkaian Treaty | Rangkaian Fakultatif |
| --- | --- | --- |
| daftar ceding | `T_POLIS_CEDING` | `T_CEDINGCOLIST` |
| daftar angsuran | `T_POLIS_INSTALMENT` | `T_LISTINSTALLMENT` |
| daftar sebaran | `T_POLIS_SPREADING` | `T_SPREADINGLIST` |
| data kuotasi | `T_POLIS_QUOTATION` | `T_QUOTATIONDATA` |

**Akibat bila dibiarkan:** dua tabel untuk satu hal, dan laporan yang harus tahu keduanya.

| Pilihan | Akibat |
| --- | --- |
| satu gaya penamaan ditetapkan | sisi yang lain mengubah rancangan dan tiketnya |
| dibiarkan, dengan pemetaan tertulis | pemeliharaan dua nama selamanya; setiap laporan lintas lini butuh pemetaan |

## 5d. Tiga seri keputusan arsitektur

Keputusan lahir di tiga rangkaian, ketiganya memulai penomoran dari nomor satu. Dokumen
*Keputusan Arsitektur* memberi **awalan seri** dan memuat **konkordansi 105 baris** serta
**pasangan sebidang** - keputusan dari seri berbeda yang membahas hal yang sama.

**Yang diminta:** apakah ketiganya dibiarkan berseri, atau dilebur menjadi satu penomoran; dan
**siapa yang memutus** pasangan sebidang.

| Pilihan | Akibat |
| --- | --- |
| dibiarkan berseri | rujukan lama tetap sahih; pembaca wajib menyebut awalan seri |
| dilebur menjadi satu penomoran | seluruh rujukan di spec dan tiket harus diperbarui |
| dilebur hanya untuk keputusan yang sebidang | lebih sedikit perubahan, tetapi dua sistem penomoran hidup bersama |

## 5e. Jadwal Fakultatif

Dari **60** tiket Fakultatif, **36** tertahan - sebagian besar oleh
tiket lain di rantainya, dan pangkal rantainya menunggu bahan pada bagian 4.

**Yang diminta:** kapan butir bagian 4 ditutup, dan oleh siapa.

| Pilihan | Akibat |
| --- | --- |
| Fakultatif dijadwalkan setelah bahannya lengkap | jadwal Fakultatif belum dapat ditetapkan hari ini |
| Fakultatif dimulai dari tiket yang siap | sebagian pekerjaan berjalan; sisanya menunggu |
| bahan diadakan dengan cara lain - cetakan layar, akses baca, dokumentasi | menambah pekerjaan pengumpulan, memperpendek penantian |

---

# 6. Risiko teratas

**1. Dua rancangan kolom untuk satu tabel fisik.** Dua rangkaian merancang tabel yang sama
dengan **79** dan **71** kolom dari dasar berbeda, dan satu sisi menyatakan rekonsiliasinya
gugur. Bila tidak diputuskan, tabrakan baru muncul saat integrasi - ketika kedua sisi sudah
menulis kode.

**2. Galat angka uang yang sudah tersimpan.** Selisih sebesar sepersepuluh juta terbukti ada di
data produksi dan lahir di rantai perhitungan. Sistem baru yang berhitung benar akan
menghasilkan angka **berbeda** - sehingga uji paritas yang menuntut kesamaan persis akan gagal
justru karena sistem barunya benar.

**3. Bahan Fakultatif belum lengkap.** **35** stored procedure tanpa isi, **nol** definisi
tabel, **604** daftar-pilihan tanpa nilai, dan **77** aturan perlu diekspor ulang. Sebagian
besar tiket Fakultatif tertahan karenanya.

**4. Bentuk dokumen lama tidak seragam.** Tiga contoh nyata hanya berbagi **23** dari **74**
medan gabungan, dan satu daftar muncul dalam **dua kedalaman berbeda**. Pengurai yang
menganggapnya seragam akan patah pada sebagian kasus, bukan seluruhnya - sehingga kegagalannya
diam.

**5. Penamaan tabel bercabang.** Empat properti sistem lama mendapat dua nama tabel. Bila
dibiarkan, setiap laporan lintas lini harus mengetahui pemetaannya.

**6. Tiga seri keputusan bernomor sama.** Rujukan `nomor 0003` dapat berarti tiga keputusan
berbeda. Risiko salah rujuk meningkat seiring bertambahnya orang yang membaca dokumen.

**7. Cacah medan adalah batas bawah.** Sensus berasal dari rujukan di dalam aturan; medan yang
tidak pernah dirujuk tidak tertangkap. Rancangan tabel yang dianggap lengkap mungkin belum.

---

# 7. Langkah berikutnya - 30 hari

| Pekan | Yang dikerjakan | Siapa |
| ---: | --- | --- |
| **1** | komite memutus **5a** dan **5b** - pemilik rancangan tabel dan kolom pembeda lini | komite pengarah |
| **1** | permintaan resmi ke DBA: isi stored procedure dan definisi tabel | pemilik pekerjaan |
| **1-2** | pembangunan tahap 1-4 dimulai: kerangka penyimpanan, nomor urut, tipe kolom, transaksi | tim pengembang |
| **2** | komite memutus **5c** penamaan tabel, sebelum kode penyimpanan bercabang | komite pengarah |
| **2-3** | modul tanpa penahan dikerjakan sejajar - Jiwa, Komite, Penempatan Treaty | tim pengembang |
| **3** | Finance memutus perlakuan galat angka lama; ambang uji paritas ditetapkan | Finance |
| **3-4** | Product menyerahkan daftar nilai daftar-pilihan; kuesioner Fakultatif dijawab | Product . UW . IT |
| **4** | jadwal Fakultatif ditetapkan berdasarkan bahan yang sudah masuk | komite pengarah |

[!] **Yang tidak dapat dijadwalkan sebelum pekan 1 selesai:** seluruh pekerjaan penyimpanan
lintas lini, dan seluruh Fakultatif.

---

*Disusun 25 September 2026. Nol rekomendasi pada bagian 5; nol butir terbuka ditutup.*