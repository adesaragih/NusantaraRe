# Tiket - Komite Klaim

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Komite Claim Fac In | **13** | 12 | 0 | 0 | 0 | 1 |
| Komite Claim Life | **10** | 10 | 0 | 0 | 0 | 0 |
| Komite Claim Prop | **14** | 13 | 0 | 0 | 0 | 1 |
| Komite Claim Non Prop | **13** | 0 | 0 | 0 | 0 | 13 |
| **Jumlah** | **50** | **35** | **0** | **0** | **0** | **15** |

---

# Komite Claim Fac In

Jumlah tiket: **13**

## Komite Claim Fac In - 01 - Pengeluaran pengaju — jenjang pertama saja

**Status:** ready-for-agent
**Blocked by:** 00
**Menutup:** AC 9 · 10 · 11 · 12 · 13 *(5 AC)* — US 13

#### Hasil & nilai pengguna

Hari ini Seorang penilai dapat **menyetujui klaim yang ia ajukan sendiri**, sebab tidak ada yang memeriksanya.

Sesudah tiket ini, Bila **pemegang jenjang pertama adalah pengaju**, ia **dikeluarkan dari susunan jenjang** kasus itu, dan **jumlah jenjang berkurang satu**.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pengeluaran pengaju | rule penyiapan wewenang penyetujuan mengeluarkan anggota **indeks 1** bila ia akun yang sedang membuka kasus, lalu **menghitung ulang** jumlah jenjang |
| ⛔ Lingkupnya | ⛔ **indeks ditanam `1`** — hanya anggota pertama yang diperiksa |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K9** — ⚠️ **DITIRU APA ADANYA** — ⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan berlaku untuk semua jenjang. ⛔ **Disengaja**

#### Yang harus diuji

- [ ] Pengaju di jenjang **pertama** ⇒ **dikeluarkan**, jumlah jenjang **berkurang satu**
- [ ] Jenjang berikutnya **naik menjadi jenjang pertama**
- [ ] ⛔⛔ **Pengaju di jenjang KEDUA ke bawah TETAP di daftar, dan TETAP dapat menyetujui klaim yang ia ajukan sendiri** — ⭐ **lubang yang dibawa masuk dengan sadar, bukan kelalaian**
- [ ] Susunan jenjang **kosong** sesudah pengeluaran ⇒ kasus **tidak dibentuk**, pengajuan **dikembalikan** berikut alasannya

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan | ⚠️ menahan penentuan siapa pemegang jenjang pertama |

#### Seam & verifikasi

**Seam:** lapisan layanan komite, pada pembentukan kasus.
⭐⭐ **UJI DUA ARAH — dan arah kedua ADALAH inti tiket ini:**
   *(a)* pengaju di jenjang **pertama** ⇒ ⭐ **dikeluarkan**, jumlah jenjang berkurang.
   *(b)* pengaju di jenjang **kedua** ⇒ ⛔ **TETAP DI DAFTAR**, dan **dapat menyetujui**.
⚠️ **Arah (b) bukan uji fungsi — ia MENGUNCI keputusan K9** supaya pengembang berikutnya tidak membacanya sebagai bug lalu 'memperbaikinya' diam-diam.
6. Pengaju satu-satunya anggota ⇒ ⭐ kasus **tidak dibentuk**, pengajuan dikembalikan.

## Komite Claim Fac In - 02 - Tangga berjalan — giliran, maju, selesai, tutup seketika

**Status:** ready-for-agent
**Blocked by:** 00 · 01
**Menutup:** AC 14 · 15 · 16 · 17 · 18 · 19 · 20 · 21 · 22 *(9 AC)* — US 14–18

#### Hasil & nilai pengguna

Hari ini Tangga persetujuan **belum berjalan**: tidak ada yang menentukan siapa menerima giliran, kapan tangga maju, dan kapan ia berhenti.

Sesudah tiket ini, Giliran berpindah ke **pemegang jenjang terendah yang belum memutuskan**. ⭐ Komite **selesai** ketika seluruh jenjang menyetujui, dan **tertutup seketika** ketika satu menolak.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penugasan berputar | router menelusuri susunan jenjang; anggota **pertama yang belum memutuskan** menerima giliran, lalu aktivitas **keluar seketika** |
| ⚠️ Pemutus lama | ⛔ di Pega tangga berhenti karena router **tidak menghasilkan penerima** — ⚠️ **perilaku diam yang tidak dapat diuji** |
| Eskalasi | kasus dapat dinaikkan **satu jenjang** bila pemegang giliran berhalangan; ⛔ turun jenjang **dilarang** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K1** — ⭐ **Komite berakhir bila SELURUH jenjang menyetujui; satu menolak ⇒ tutup seketika.** ⭐ Pemutusnya **keadaan tiap jenjang**, ⛔ **bukan pencacah**

#### Yang harus diuji

- [ ] Giliran diberikan ke **pemegang jenjang terendah yang belum memutuskan**
- [ ] ⭐ Seluruh jenjang menyetujui ⇒ kasus **selesai**
- [ ] ⭐ Satu jenjang menolak ⇒ kasus **tertutup seketika**; ⛔ jenjang berikutnya **tidak** menerima giliran
- [ ] ⛔ Sebuah jenjang **hanya dapat memutuskan sekali**; percobaan kedua **ditolak**
- [ ] Eskalasi **naik satu jenjang** tersedia; ⛔ **turun jenjang ditolak**
- [ ] Eskalasi **terekam**: siapa memindahkan, kapan, dari jenjang mana ke mana
- [ ] Jenjang yang **dilewati** eskalasi tercatat **tanpa keputusan** — ⛔ bukan sebagai menyetujui

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi susunan jenjang | ⚠️ menahan jumlah tingkat yang diuji |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — masukkan susunan jenjang, jalankan urutan keputusan.
⛔ **Uji PERILAKU, bukan pencacah.** ⭐ Pencacah boleh ada sebagai tampilan; ia **bukan kebenaran**.
3. Tiga jenjang, semua setuju ⇒ ⭐ selesai sesudah yang **ketiga**.
4. Tiga jenjang, yang **kedua menolak** ⇒ ⭐ tutup **seketika**; jenjang ketiga **tidak** menerima giliran.
5. Keluarkan pengaju lebih dulu ⇒ ⭐ tangga tetap benar walau jumlah jenjang **berubah di awal**.

## Komite Claim Fac In - 03 - Wewenang komite — ditegakkan di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** 00 · `claim-facin\issues\08` *(wewenang klaim)*
**Menutup:** AC 23 · 24 · 25 · 26 · 27 · 28 · 29 *(7 AC)* — US 19–21

#### Hasil & nilai pengguna

Hari ini ⛔ `[terverifikasi]` Hanya **satu** pemeriksaan pemilik giliran ada di korpus, dan ia menempel pada **tombol Submit di layar**. ⚠️ Di modul saudaranya **tidak ada sama sekali**. ⛔ Siapa pun yang dapat memanggil lapisan layanan dapat menyimpan keputusan untuk jenjang mana pun.

Sesudah tiket ini, ⭐ **Hanya akun beku pada jenjang berjalan** yang dapat menyimpan keputusan, dan penolakannya terjadi **di lapisan layanan** — ⛔ bukan di layar.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pemeriksaan pemilik giliran | satu-satunya di korpus; menempel pada **tombol Submit** |
| ⚠️ Penukaran identitas | ⛔ **dua akun ditukar menjadi akun ketiga SEBELUM pemeriksaan** — ⭐ **tidak dibawa** |
| Modul saudara | ⛔ **nol pemeriksaan** — tidak di aktivitas, tidak di layar |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K10** — ⭐ **ADR-U-0014 BERDIRI** — penegakan pemilik giliran **tetap di lapisan layanan**
- **K3 · K12** — ⭐ Pembandingan memakai **identitas akun**, ⛔ bukan alias

#### Yang harus diuji

- [ ] Hanya **akun beku pada jenjang berjalan** dapat menyimpan keputusan
- [ ] Percobaan oleh akun lain **ditolak dengan galat** — ⛔ bukan diabaikan diam-diam
- [ ] Percobaan yang ditolak **terekam** di jejak audit
- [ ] ⛔ Penukaran akun menjadi akun ketiga **tidak dibawa**
- [ ] ⭐ Layar boleh menyembunyikan tindakan tak berwenang — ⛔ **itu bukan penegakan**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar jabatan dan susunan jenjang BELUM ADA** — korpus tidak memuat satu pun rule otorisasi; 13 medan privilese seluruhnya kosong | ⛔⛔ **MENAHAN PEMBANGUNAN tiket ini.** ⭐ Tiketnya tetap ditulis dan jahitannya disebut; ⛔ ia **tidak dapat selesai** sebelum work owner memberikan daftarnya |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — ⭐ **uji lewat sini, bukan lewat layar**.
2. Panggil sebagai **bukan pemegang giliran** ⇒ ⛔ **ditolak**.
3. Sembunyikan tombolnya di layar, lalu panggil layanan **langsung** ⇒ ⛔ **tetap ditolak**.
4. Periksa jejak audit ⇒ ⭐ percobaan yang ditolak **tercatat**.
5. Panggil sebagai akun yang dulu 'menyamar' ⇒ ⛔ **ditolak** — ⭐ penukaran tidak dibawa.

## Komite Claim Fac In - 04 - Layar komite — 92 medan, dua kotak centang terkunci

**Status:** ready-for-agent
**Blocked by:** 00 · 02
**Menutup:** AC 30 · 31 · 32 · 33 · 34 · 35 *(6 AC)* — US 6–11

#### Hasil & nilai pengguna

Hari ini Anggota komite **belum punya layar** untuk menilai penyesuaian, dan aturan siapa boleh menyunting apa belum ditetapkan.

Sesudah tiket ini, Anggota komite melihat rekapitulasi penyesuaian **dalam mata uang asli dan dalam rupiah**, dapat menuliskan catatan, dan ⭐ **dua kotak centang usulan hanya dapat disunting penyetuju pertama**.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Layar komite | 92 medan; ⭐ **tepat dua** terkunci bagi jenjang selain yang pertama |
| Dua kotak centang | usul **menutup klaim** dan usul **mencadangkan** |
| ⚠️ Jebakan alih | ⛔ pembanding di Pega bertipe **teks**, padahal pencacahnya **bilangan** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K4** — ⭐ **Ditiru apa adanya** — usulan tindak lanjut dibuat **penyetuju pertama**; jenjang berikutnya **menilai**, ⛔ tidak mengganti

#### Yang harus diuji

- [ ] Penyetuju **pertama** dapat menyunting kedua kotak centang usulan
- [ ] Penyetuju **kedua ke atas** melihat keduanya **terkunci**
- [ ] ⭐ **90 medan lain TIDAK terkunci** oleh aturan itu — ⛔ jangan menguncinya karena salah membaca
- [ ] Layar menampilkan rekapitulasi **dalam mata uang asli** dan **total dalam rupiah**
- [ ] Anggota komite dapat menuliskan **catatan** pada keputusannya
- [ ] ⛔ Pembandingan memakai **bilangan dengan bilangan** — ⚠️ Pega memaafkan teks, sistem baru tidak

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **11** | ⚠️ 29 medan bergantung jenis objek | ⚠️ menahan bentuk tampilan objek |

#### Seam & verifikasi

**Seam:** lapisan layanan komite untuk data; layar hanya menampilkan.
1. Buka layar sebagai penyetuju **pertama** ⇒ ⭐ kedua kotak centang **dapat disunting**.
2. Buka sebagai penyetuju **kedua** ⇒ ⭐ keduanya **terkunci**.
3. Periksa 90 medan lain pada penyetuju kedua ⇒ ⛔ **tidak ikut terkunci**.
4. Bandingkan total rupiah dengan hitungan manual ⇒ ⭐ cocok.

## Komite Claim Fac In - 05 - Keputusan tersimpan — akun dirujuk, jabatan disalin

**Status:** ready-for-agent
**Blocked by:** 02 · 03
**Menutup:** AC 36 · 37 · 38 · 39 · 40 · 41 · 42 *(7 AC)* — US 12 · 24 · 25

#### Hasil & nilai pengguna

Hari ini Keputusan tiap jenjang **belum tersimpan dengan jejak yang dapat dipertanggungjawabkan** — ⚠️ dan di sistem lama, jabatan yang tercatat **diambil dari daftar yang ditanam di kode**, dengan **nilai bawaan jabatan direksi** bagi siapa pun yang tak dikenal.

Sesudah tiket ini, Setiap keputusan tersimpan bersama **akun pemutus, salinan kode jabatan saat itu, waktu, dan keputusannya** — dan ⭐ **tidak dapat diubah** sesudah tersimpan.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penyimpan keputusan | dua rule khas modul ini, keduanya dari **salinan produksi**; ⛔ **nol gerbang wewenang** di dalamnya |
| ⚠️ Nilai bawaan jabatan | ⛔ siapa pun di luar daftar tercatat sebagai **jabatan direksi** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K13** — ⭐⭐ **Identitas orang DIRUJUK · jabatan DISALIN · tulisan jabatan DIRUJUK**
- **K11** — ⭐ Jabatan dibaca dari **data pengguna**; ⛔ bila tidak diketahui, **tolak** — ⛔ **tidak ada jabatan bawaan**

#### Yang harus diuji

- [ ] Keputusan tersimpan bersama **akun**, **salinan kode jabatan saat itu**, **waktu**, dan keputusannya
- [ ] ⭐ **Jabatan disimpan sebagai SALINAN kode** — ⛔ naik jabatan **tidak mengubah catatan lama**
- [ ] ⭐ **Akun disimpan sebagai RUJUKAN** — nama tampil yang berubah **memang** ikut berubah
- [ ] ⭐ **Tulisan jabatan dirujuk** dari daftar induk — ⛔ catatan lama **tidak boleh berpindah ke jabatan yang berbeda**
- [ ] ⛔ Jabatan pemutus **tidak diketahui** ⇒ penyimpanan **ditolak**; ⛔ **tidak ada bawaan**
- [ ] ⛔ Keputusan yang tersimpan **tidak dapat diubah**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar induk jabatan | ⚠️ menahan pencatatan **jabatan**; ⭐ tidak menahan pencatatan **akun** |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — penyimpanan keputusan.
1. Simpan satu keputusan ⇒ ⭐ akun, kode jabatan, waktu, dan keputusan tercatat.
2. Coba ubah keputusan yang sudah tersimpan ⇒ ⛔ **ditolak**.
3. Simpan keputusan oleh akun yang **jabatannya tidak diketahui** ⇒ ⛔ **ditolak**, tanpa bawaan.

## Komite Claim Fac In - 06 - Efek akhir — terbit sekali, oleh jenjang terakhir

**Status:** ready-for-agent
**Blocked by:** 05
**Menutup:** AC 43 · 44 · 45 · 46 · 47 · 48 · 49 · 50 · 51 *(9 AC)* — US 28 · 29 · 32 · 33

#### Hasil & nilai pengguna

Hari ini Ringkasan akseptasi, nomor akseptasi, dokumen, dan surel **belum terbit sama sekali** — dan ⚠️ di sistem lama, langkah yang menyusunnya terpisah dari langkah yang menyimpannya, sehingga pernah dikira **penerimaan tidak tersimpan**.

Sesudah tiket ini, ⭐ **Efek akhir terbit HANYA ketika jenjang terakhir menyetujui**, dan **tepat satu kali**: ringkasan akseptasi disusun lalu **disimpan**, nomor terbit, dokumen tercetak, surel terkirim.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penyusun ringkasan | 51 langkah; ⛔ langkah penyimpan **di dalamnya** ber-remark |
| ⭐ Penyimpan sesungguhnya | ⭐ langkah penyimpan **diangkat ke rule pemanggil**, bergerbang *keputusan terima* **dan** *jenjang terakhir* |
| ⚠️ Jalur kegagalan | ⛔ **NOL jalur kegagalan pada 25 langkah penyimpanan** di sistem lama |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K1** — ⭐ Efek akhir menempel pada **jenjang terakhir menyetujui**

#### Yang harus diuji

- [ ] ⭐ Efek akhir terbit **hanya** ketika **jenjang terakhir menyetujui**
- [ ] ⭐ Ringkasan akseptasi **disusun lalu disimpan** — penyimpanan terjadi **sekali**
- [ ] ⛔ **RALAT yang wajib dibaca:** vonis lama *“penolakan tersimpan, penerimaan tidak”* **BATAL** — ⭐ keduanya tersimpan; yang berbeda hanya **di rule mana** langkah penyimpannya duduk
- [ ] Urutannya: ringkasan → nomor → dokumen → surel → kasir
- [ ] Surel terkirim **hanya di lingkungan produksi**
- [ ] ⚠️ Kegagalan salah satu efek **tidak boleh** membatalkan efek yang sudah terbit tanpa jejak — ⭐ sistem baru **wajib** punya jalur kegagalan

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | tidak menahan |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — penyelesaian kasus.
1. Jalankan tangga tiga jenjang sampai selesai ⇒ ⭐ efek akhir terbit **tepat satu kali**.
2. Periksa jenjang pertama dan kedua ⇒ ⛔ **nol efek akhir** di keduanya.
3. Paksa gagal di tengah rangkaian efek ⇒ ⭐ kegagalannya **tercatat**, tidak hilang diam-diam.
4. Jalankan di lingkungan bukan produksi ⇒ ⛔ surel **tidak terkirim**.

## Komite Claim Fac In - 07 - Nomor akseptasi — satu rangkaian, dikunci kelas kasus induk

**Status:** ready-for-agent
**Blocked by:** 06
**Menutup:** AC 52 · 53 · 54 · 55 *(4 AC)* — US 28

#### Hasil & nilai pengguna

Hari ini Nomor akseptasi **belum terbit**, dan di korpus ada **tujuh pembangkit per lini usaha** yang seluruhnya **ber-remark** — sehingga tidak jelas mana yang berlaku.

Sesudah tiket ini, ⭐ Nomor akseptasi terbit dari **satu rangkaian**, dikunci pada **jenis kasus induk** — ⛔ bukan per lini usaha — dan **tidak dibangkitkan ulang** untuk akseptasi yang sudah bernomor.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pembangkit yang hidup | satu prosedur tersimpan, dikunci **kelas kasus induk** |
| ⛔ Tujuh pembangkit per lini | ⛔ **ber-remark seluruhnya**; sisiran menyeluruh **tidak menemukan penggantinya** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K1** — ⭐ Nomor terbit sesudah ringkasan tersimpan, di **jenjang terakhir**

#### Yang harus diuji

- [ ] ⭐ Nomor berasal dari **satu rangkaian**, dikunci **jenis kasus induk**
- [ ] ⛔ Tujuh pembangkit per lini usaha **tidak dialihkan**
- [ ] ⛔ Nomor yang sudah terbit **tidak dibangkitkan ulang**
- [ ] ⚠️ Perilaku prosedur tersimpannya **belum terbaca dari korpus** — ⭐ sistem baru **memerikan sendiri** aturan pembangkitannya, ⛔ tidak mewarisi yang tak terbaca

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **31** | Isi prosedur tersimpan pembangkit nomor — `[data DBA]` | ⚠️ menahan **peniruan persis**; ⭐ tidak menahan pembuatan aturan sendiri |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — penerbitan nomor.
1. Selesaikan dua kasus komite ⇒ ⭐ nomornya **berurut dalam satu rangkaian**.
2. Selesaikan kasus dua lini usaha berbeda ⇒ ⭐ keduanya dari **rangkaian yang sama**.
3. Coba terbitkan ulang untuk akseptasi yang sudah bernomor ⇒ ⛔ **ditolak**.

## Komite Claim Fac In - 08 - Penggolongan lini usaha di sisi komite

**Status:** ready-for-agent
**Blocked by:** 00 · `claim-facin\issues\02` *(klasifikasi lini)*
**Menutup:** AC 56 · 57 · 58 · 59 · 60 *(5 AC)* — US 34–36

#### Hasil & nilai pengguna

Hari ini ⛔ `[terverifikasi]` Penggolongan lini usaha di sisi komite **membaca medan yang tidak pernah disalin** ke objek kerja komite. ⚠️ Akibatnya cabang **MBU** dan **Travel** **tidak pernah terbit** — diam-diam, tanpa galat.

Sesudah tiket ini, Setiap lini usaha **tergolong benar di sisi komite**, termasuk **MBU dan Travel** yang dulu tak pernah terbit.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| 13 penggolong hidup | ⭐ dari 49 rule penggolong, **13 dipakai hidup**; ⛔ **36 tidak dipakai sama sekali** |
| ⚠️ Cacat lama | ⛔ sembilan penggolong menguji medan yang **tidak disalin**; **dua dipakai hidup lima kali** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K2** — ⭐ **Data kutipan disalin UTUH** — ⚠️ `[penyimpangan sadar]`, **cacat yang diperbaiki**. ⛔ Bukan daftar medan bernama — **daftar bernama itulah yang melahirkan cacat ini**

#### Yang harus diuji

- [ ] ⭐ **13 penggolong** dialihkan; ⛔ **36 tidak dibangun**
- [ ] ⭐ Penggolongan membaca **data kutipan lengkap**
- [ ] ⭐ Cabang **MBU** dan **Travel** **terbit** bila datanya memenuhi — ⚠️ inilah cacat lama yang ditutup
- [ ] ⛔ Penggolong berhubungan **ATAU** tetap berperilaku **ATAU**, ⚠️ bukan **DAN**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **29** | Cacah baris lama MBU/Travel terdampak — `[data DBA]` | tidak menahan |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — penyusunan ringkasan akseptasi.
1. Jalankan kasus komite lini **MBU** ⇒ ⭐ cabangnya **terbit**.
2. Jalankan lini **Travel** ⇒ ⭐ cabangnya **terbit**.
⚠️ Keduanya **tidak terbit di sistem lama** — ⭐ uji ini membuktikan cacatnya tertutup.
3. Cari pemanggilan salah satu dari 36 penggolong yang tidak dialihkan ⇒ ⛔ **nihil**.

## Komite Claim Fac In - 09 - Jalur kasir — penjaga ganda-bayar eksplisit

**Status:** ready-for-agent
**Blocked by:** 06 · `claim-facin\issues\11` *(efek keluar klaim)*
**Menutup:** AC 61 · 62 · 63 · 64 · 65 · 66 · 67 · 68 · 69 · 70 *(10 AC)* — US 30 · 31

#### Hasil & nilai pengguna

Hari ini Instruksi pembayaran **belum terkirim ke kasir**, dan ⚠️ di sistem lama **penjaga ganda-bayar bersandar pada gerbang bertanda sama dengan TUNGGAL** yang ⛔ belum dapat dipastikan apakah ia pembandingan atau penugasan. ⚠️ Pemberitahuan galat pun terkirim **setiap kali**, sebab gerbangnya ber-bendera mati.

Sesudah tiket ini, ⭐ Instruksi pembayaran terkirim **tepat satu kali**, penandanya tersimpan **dalam transaksi yang sama**, dan ⭐ **pemberitahuan galat terkirim HANYA ketika pengiriman gagal**.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Jalur kasir | 53 langkah, sarang **6 tingkat**, ⛔ **nol parameter** — rule terdalam di modul |
| Dua panggilan luar | keduanya bergerbang **lingkungan produksi**; ⭐ keduanya punya **lompatan kegagalan** — satu-satunya penanganan kegagalan eksplisit di modul |
| ⚠️ Penjaga ganda-bayar | ⛔ gerbangnya memakai tanda sama dengan **tunggal** |
| ⚠️ Pemberitahuan galat | ⛔ gerbang **ber-bendera mati** ⇒ terkirim **setiap kali**; ⚠️ medannya **salah eja** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K5** — ⭐ **Penjaga ganda-bayar dibuat EKSPLISIT** — penanda tersimpan **dalam transaksi yang sama**; panggilan kedua **ditolak**. ⚠️ `[penyimpangan sadar]`
- **K6-lama** — ⭐ **Pemberitahuan galat HANYA pada kegagalan** — ⛔ perilaku lama **tidak ditiru**. ⚠️ `[penyimpangan sadar]`

#### Yang harus diuji

- [ ] ⭐ Pembayaran terkirim **tepat satu kali** per akseptasi
- [ ] ⭐ Penanda terkirim tersimpan **dalam transaksi yang sama** dengan pengirimannya
- [ ] ⛔ Panggilan kedua atas akseptasi yang sama **ditolak**
- [ ] ⭐ Pemberitahuan galat terkirim **hanya pada kegagalan**
- [ ] ⭐ Medan tanggapan layanan **dinamai dengan benar** — ⚠️ di sistem lama salah eja
- [ ] ⭐ Kegagalan pengiriman **tidak membatalkan** akseptasi yang sudah tersimpan; ia **dicatat** dan **dapat diulang**
- [ ] ⛔ Rancangan **tidak bergantung** pada tafsir tanda sama dengan tunggal

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **26** | ⚠️ **Tanda sama dengan TUNGGAL** pada penjaga ganda-bayar — pembandingan atau penugasan | ⛔ **TIDAK menahan** — ⭐ K5 mengurungnya; jawabannya dipakai untuk **memeriksa DATA LAMA**, yakni apakah sistem lama pernah membayar dua kali |
| **27** | Apakah pemberitahuan galat lama **benar-benar sampai** ke seseorang | ⚠️ tidak menahan — ⭐ bila ternyata tidak, penanganan galat lama **sebenarnya tidak ada** |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — pengiriman ke kasir.
⛔ **Uji sebagai PERILAKU**, bukan sebagai isi kolom penanda.
1. Kirim satu instruksi ⇒ ⭐ terkirim, penanda terisi **dalam transaksi yang sama**.
2. Kirim ulang untuk akseptasi yang sama ⇒ ⛔ **ditolak**.
3. Paksa kegagalan ⇒ ⭐ pemberitahuan galat **terkirim**, akseptasi **tetap tersimpan**.
4. Kirim yang **berhasil** ⇒ ⛔ pemberitahuan galat **TIDAK terkirim** — ⚠️ di sistem lama ia terkirim.
5. Jalankan di lingkungan bukan produksi ⇒ ⛔ panggilan luar **tidak dijalankan**.

## Komite Claim Fac In - 10 - Jalur Fac Retro — jalur yang melompati penyiapan wewenang

**Status:** ready-for-agent
**Blocked by:** 00 · 03
**Menutup:** AC 71 · 72 · 73 · 74 · 75 · 76 *(6 AC)* — US 37 · 38

#### Hasil & nilai pengguna

Hari ini Penyesuaian bertanda **Fac Retro** **melompati penyiapan wewenang penyetujuan** — ⛔ dan itu **tidak terlihat oleh siapa pun yang membaca aturan**, sebab ia terkubur di dalam detail.

Sesudah tiket ini, ⭐ Jalur Fac Retro **dikenali sebagai jalur tersendiri**, disebut terang-terangan, sehingga siapa pun yang membaca aturan **melihat jalur itu ada** — ⛔ bukan menemukannya sesudah sesuatu terjadi.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penanda Fac Retro | jenis treaty tertentu menandai penyesuaian sebagai Fac Retro |
| ⛔ Akibatnya | ⛔ penyiapan wewenang penyetujuan **keluar seketika** bila penanda itu menyala |
| ⚠️ Tiga jalur pemicu | ⛔ `[terverifikasi]` penanda dipicu **setidaknya tiga jalur**, ⚠️ salah satunya **tanpa gerbang sama sekali** |
| ⚠️ Medan jenis treaty | ⛔ **mencampur kode master dengan teks harfiah** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K7** — ⚠️ **DITIRU APA ADANYA** — ⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan menjadikan penanda jenis treaty sebagai **data**. ⛔ **Disengaja**

#### Yang harus diuji

- [ ] ⭐ Penyesuaian bertanda Fac Retro mengikuti **jalur tersendiri** yang **melompati penyiapan wewenang**
- [ ] ⚠️ Nilai penanda jenis treaty tetap **tetapan di dalam kode** ⇒ ⛔ **setiap perubahan daftar treaty menuntut RILIS perangkat lunak**
- [ ] ⛔⛔ **Spec dan tiket ini menyebut jalur itu SECARA TERBUKA** sebagai jalur yang melompati pemeriksaan wewenang — ⭐ supaya ia **terlihat**, bukan terkubur
- [ ] ⚠️ Ketiga jalur pemicu **ditelusuri** sebelum dibangun — ⛔ termasuk yang **tanpa gerbang**
- [ ] ⚠️ Pencampuran kode master dengan teks harfiah pada medan jenis treaty **wajib beres** sebelum medan itu boleh menjadi kolom bernilai kode

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **28** | ⚠️ **Arti nilai penanda jenis treaty**, dan apakah **ketiga jalur pemicu** memang dikehendaki | ⚠️ tidak menahan pembangunan — ⭐ menahan **kepastian bahwa seluruh jalurnya sudah dikenali** |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — penyiapan wewenang penyetujuan.
1. Jalankan penyesuaian **bertanda Fac Retro** ⇒ ⭐ penyiapan wewenang **dilewati**.
2. Jalankan penyesuaian **biasa** ⇒ ⭐ penyiapan wewenang **dijalankan**.
⚠️ Picu penanda lewat **ketiga jalur** ⇒ ⭐ ketiganya menghasilkan perilaku **yang sama**.
3. Periksa dokumentasi jalur ⇒ ⭐ ia **disebut sebagai jalur tersendiri**, bukan catatan kaki.

## Komite Claim Fac In - 11 - Jabatan, tabel login, dan jejak audit

**Status:** ready-for-agent
**Blocked by:** 05 · `claim-facin\issues\12` *(jejak audit klaim)*
**Menutup:** AC 77 · 78 · 79 · 80 · 81 · 82 · 83 · 84 · 85 *(9 AC)* — US 21 · 22 · 23 · 26

#### Hasil & nilai pengguna

Hari ini Jabatan **ditanam di dalam kode** sebagai rantai pencocokan ID pengguna perorangan, dan ⚠️ **orang yang sama tertulis dengan jabatan berbahasa berbeda antar modul**. ⛔ Mengubah susunan jenjang menuntut **rilis**.

Sesudah tiket ini, ⭐ Jabatan menjadi **kolom berisi KODE di tabel login**, susunan jenjang komite ada di **daftar terpisah**, dan ⭐ **perubahan jabatan masuk jejak audit** — sebab jabatan **menentukan siapa boleh menyetujui**.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| ⚠️ Tabel jabatan lama | ⛔ rantai pencocokan **ID pengguna perorangan** dengan jabatan |
| ⚠️ Dua bahasa | ⛔ orang yang sama tertulis berbeda antar modul, dan **kuncinya beda medan** |
| Roster lama | ⛔ menyimpan **nama orang**, bukan jabatan — ⭐ itulah akar penukaran akun |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K13** — ⭐⭐ **Empat syarat:** *(1)* isinya **KODE** menunjuk daftar induk jabatan · *(2)* **urutan jenjang di daftar terpisah** · *(3)* catatan keputusan **menyimpan SALINAN kode jabatan** · *(4)* perubahan kolom jabatan **masuk jejak audit**
- **K12** — ⭐ **Roster berbasis JABATAN**, ⛔ bukan nama orang
- **K11** — ⭐ Jabatan dibaca dari **data pengguna**; ⛔ **nilai bawaan tidak dibawa**

#### Yang harus diuji

- [ ] ⭐ Kolom jabatan berisi **kode**, menunjuk **daftar induk jabatan** — ⛔ bukan teks bebas
- [ ] ⭐ **Susunan jenjang komite** disimpan **terpisah** dari tabel login
- [ ] ⭐ Mengubah susunan jenjang adalah **mengubah baris data** — ⛔ bukan rilis
- [ ] ⭐ **Perubahan kolom jabatan masuk jejak audit**: siapa mengubah, kapan, dari kode apa ke kode apa
- [ ] ⛔ Nama orang **tidak menjadi kunci** di mana pun
- [ ] ⭐ Tulisan jabatan **seragam** — ⚠️ di sistem lama berbeda bahasa antar modul

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | ⛔⛔ **Isi daftar induk jabatan dan susunan jenjang BELUM ADA** | ⛔⛔ **MENAHAN PEMBANGUNAN** |
| **33-lama** | Alasan historis pemetaan akun menjadi nama orang lain | tidak menahan |

#### Seam & verifikasi

**Seam:** lapisan layanan — pembacaan jabatan dan penulisan jejak audit.
⭐⭐ **UJI SALINAN JABATAN — inti tiket ini:**
   *(a)* catat satu keputusan komite;
   *(b)* **naikkan jabatan** orang yang memutuskannya;
   *(c)* ⭐ **catatan lama TIDAK berubah**.
⚠️ **Ini bukan uji fungsi** — ia yang **membedakan SALINAN dari RUJUKAN**. ⛔ Tanpa uji ini, pengembang berikutnya akan menggantinya dengan rujukan karena terlihat lebih rapi.
1. Ubah nama tampil orangnya ⇒ ⭐ catatan lama **ikut nama baru** — ⭐ itu memang dikehendaki.
2. Ubah kolom jabatan seseorang ⇒ ⭐ perubahannya **tercatat di jejak audit**.
3. Ubah susunan jenjang ⇒ ⭐ cukup **mengubah baris**, tanpa rilis.

## Komite Claim Fac In - 12 - Kronologi komite dan penandaan data lama

**Status:** ready-for-agent
**Blocked by:** 05 · 11
**Menutup:** AC 86 · 87 · 88 · 89 · 90 · 91 · 92 · 93 · 94 · 95 *(10 AC)* — US 24 · 25 · 27

#### Hasil & nilai pengguna

Hari ini Catatan kronologi komite **disimpan sebagai kalimat jadi**, sehingga ⚠️ jabatan yang keliru **membeku di dalam teks**. ⚠️ Dan catatan akseptasi lama lini **MBU dan Travel** dibangun **tanpa cabangnya** — tanpa penanda apa pun yang membedakannya.

Sesudah tiket ini, ⭐ Catatan kronologi disusun dari **medan tersimpan** saat ditampilkan, ⛔ bukan disimpan sebagai kalimat jadi. Dan ⭐ **catatan lama MBU/Travel DITANDAI** sehingga laporan **dapat menyaringnya**.

#### Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Catatan kronologi | setiap keputusan komite menambah satu catatan |
| ⛔ RALAT penting | ⛔ sasarannya **catatan kronologi kasus** — ⭐ **jejak internal**, ⚠️ **bukan** dokumen akseptasi yang keluar perusahaan |
| ⚠️ Nilai bawaan | ⛔ siapa pun di luar daftar tercatat sebagai **jabatan direksi** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

#### Keputusan work owner yang mengikat

- **K8** — ⭐ **Data lama MBU dan Travel DIBIARKAN, tetapi DITANDAI** — ⛔ tidak dibangun ulang, sebab angkanya mungkin **sudah dilaporkan keluar**
- **K13** — ⭐ Catatan disusun dari **medan tersimpan** — akun, kode jabatan, keputusan, waktu

#### Yang harus diuji

- [ ] ⭐ Setiap keputusan komite menambah **catatan kronologi**
- [ ] ⭐ Catatan disusun dari **medan tersimpan**, ⛔ **bukan disimpan sebagai kalimat jadi**
- [ ] ⛔⛔ Menghapus kasus **TIDAK menghapus** kronologinya — ⭐ jejak yang ikut terhapus **berhenti menjadi jejak**
- [ ] ⛔ **RALAT:** sasaran catatan itu **jejak internal**, ⚠️ bukan dokumen keluar — ⭐ tetapi akibatnya tetap serius: **jejak audit mencatat jabatan yang keliru**
- [ ] ⭐ Catatan akseptasi lama MBU/Travel **ditandai** dengan kolom penanda ber-nilai dua keadaan
- [ ] ⭐ **Alasan penandanya KOLOM, bukan catatan prosa:** supaya **laporan dapat menyaringnya** — ⛔ bukan hanya pembaca manusia yang kebetulan membaca catatan kaki
- [ ] ⛔ Catatan **baru** tidak menerima penanda itu
- [ ] ⛔ Data lama **tidak dibangun ulang**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **29** | Cacah baris MBU/Travel terdampak — `[data DBA]` | ⚠️ menahan **cacahnya**, ⭐ tidak menahan cara penandaannya |

#### Seam & verifikasi

**Seam:** lapisan layanan komite — penulisan kronologi dan penandaan migrasi.
1. Putuskan satu jenjang ⇒ ⭐ kronologi bertambah satu baris.
2. Hapus kasus komitenya ⇒ ⛔ kronologinya **tetap ada**.
3. Ubah tulisan sebuah jabatan di daftar induk ⇒ ⭐ catatan lama **ikut tulisan baru**, ⛔ tetapi **tidak berpindah ke jabatan yang berbeda**.
4. Periksa catatan akseptasi lama MBU ⇒ ⭐ penandanya **menyala**.
5. Terbitkan catatan **baru** ⇒ ⛔ penandanya **padam**.
6. Jalankan laporan dengan saringan penanda ⇒ ⭐ kedua kelompok **terpisah bersih**.

## Komite Claim Fac In - uru - Urutan tiket — Komite Claim Fac In

⛔ **Ini bukan tiket.** Berkas ini hanya menjelaskan **tiket mana harus selesai sebelum tiket mana**.
Isinya **tidak menambah keputusan apa pun**.

Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5 ·
`claim-facin\RELASI-TABEL-CLAIM-FACIN.md` relasi 12–15.

---

#### ⛔⛔ Nomor di dua folder BERDIRI SENDIRI

⚠️ **Modul ini dan modul Claim Fac In punya folder tiket masing-masing**, dan **keduanya bernomor
mulai dari `00`**.

> ⛔ **Tiket `00` di sini BUKAN tiket `00` di `claim-facin\issues\`.**
> ⭐ **Rujukan lintas-modul WAJIB berawalan foldernya** — contoh: `claim-facin\issues\13`.

⭐ **Sebabnya folder dipisah:** tiket memerikan **perilaku sebuah modul**, dan modul ini punya
**spec sendiri dengan 100 AC sendiri**. ⭐ **Tiket tinggal bersama spec yang ia tutup** — sehingga
`Menutup: AC …` tidak pernah ambigu.

⚠️ **Tabel basis datanya tetap satu** — `[keputusan work owner]` **K4 · K6**, komite lini FAC dan
lini PROP memakai tabel yang sama. ⭐ **Tabel bersama tidak menuntut tiket bersama.**

---

#### Tiga belas tiket

| # | Judul | Blocked by |
| --- | --- | --- |
| ⭐ **00** | Kasus komite lahir dari penyesuaian — jenjang dibekukan | `claim-facin\issues\13` |
| ⚠️ **01** | Pengeluaran pengaju — jenjang pertama saja | 00 |
| ⭐ **02** | Tangga berjalan — giliran, maju, selesai, tutup seketika | 00 · 01 |
| ⛔ **03** | Wewenang komite — ditegakkan di lapisan layanan | 00 · `claim-facin\issues\08` |
| **04** | Layar komite — 92 medan, dua kotak centang terkunci | 00 · 02 |
| ⭐ **05** | Keputusan tersimpan — akun dirujuk, jabatan disalin | 02 · 03 |
| ⭐ **06** | Efek akhir — terbit sekali, oleh jenjang terakhir | 05 |
| **07** | Nomor akseptasi — satu rangkaian, dikunci kelas kasus induk | 06 |
| **08** | Penggolongan lini usaha di sisi komite | 00 · `claim-facin\issues\02` |
| ⚠️ **09** | Jalur kasir — penjaga ganda-bayar eksplisit | 06 · `claim-facin\issues\11` |
| ⚠️ **10** | Jalur Fac Retro — jalur yang melompati penyiapan wewenang | 00 · 03 |
| ⭐ **11** | Jabatan, tabel login, dan jejak audit | 05 · `claim-facin\issues\12` |
| **12** | Kronologi komite dan penandaan data lama | 05 · 11 |

---

#### ⭐ Rantai terdalam

```
claim-facin\issues\00 ──▶ … ──▶ claim-facin\issues\13
                                     └─▶ 00 ──▶ 01 ──▶ 02 ──▶ 05 ──▶ 06 ──▶ 09
```

⭐ **Enam tingkat di dalam modul ini**, ⚠️ **menumpang rantai delapan tingkat** di modul Claim Fac
In. ⛔ Jadi tiket **09** adalah yang **paling jauh dari titik mulai** pada seluruh pekerjaan Fac In.

#### ⭐ Yang dapat berjalan bersamaan

| Sesudah selesai | Dapat mulai bersamaan |
| --- | --- |
| **00** | 01 · **08** *(penggolongan hanya butuh kasus lahir)* · **10** |
| **02** | 04 · *(03 bila wewenang klaim sudah selesai)* |
| **05** | 06 · **11** |
| **06** | 07 · 09 |
| **11** | 12 |

---

#### ⚠️ Lima tiket bergantung pada modul Claim Fac In

| Tiket di sini | Menunggu | Kenapa |
| --- | --- | --- |
| **00** | `claim-facin\issues\13` | ⭐ **kontrak muatan** — kasus komite lahir dari penyerahan |
| **03** | `claim-facin\issues\08` | ⭐ **wewenang dibangun sekali**, dipakai kedua modul |
| **08** | `claim-facin\issues\02` | ⭐ **penggolongan lini usaha dibangun sekali** |
| **09** | `claim-facin\issues\11` | ⭐ **jalur kasir dibangun sekali** |
| **11** | `claim-facin\issues\12` | ⭐ **jejak audit dibangun sekali** |

⭐ **Kelimanya benar** — ⛔ keduanya **bukan** salinan yang berdiri sendiri.

---

#### ⛔ Dua tiket yang TERTAHAN

| Tiket | Penahan | Sifat |
| --- | --- | --- |
| ⛔ **03 · Wewenang** | **butir 6** — ⛔ isi daftar jabatan dan susunan jenjang **belum ada** | ⛔⛔ **MENAHAN PEMBANGUNAN.** ⭐ Tiketnya **tetap ditulis**, jahitannya disebut |
| ⛔ **11 · Jabatan & jejak audit** | **butir 6** yang sama | ⛔⛔ **MENAHAN PEMBANGUNAN** |

⚠️ **Tiket 00 tertahan sebagian** — ⭐ kasus dapat lahir, ⛔ tetapi **jumlah jenjang** menunggu butir
6, dan **butir 39** *(satu jabatan, beberapa pemegang)* menahan bila keadaan itu terjadi.

---

#### ⭐⭐ Liputan Acceptance Criteria

`komite-claim-facin\spec.md` memuat **100 AC**, bernomor **1–100**, ⛔ tanpa nomor hilang.

| | Jumlah |
| --- | ---: |
| ⭐ tertutup **tepat satu kali** | ⭐ **100** |
| tertutup **lebih dari satu kali** | ⛔ **0** |
| ⛔ **tidak tertutup** | ⭐ **0** |

⭐ **Dihitung DUA CARA** — *(a)* dari baris `Menutup:` tiap tiket; *(b)* sisiran pola nomor pada
seluruh berkas tiket. ✅ **Keduanya sepakat.**

⭐ **Liputan penuh, dan itu berbeda dari sisi klaim.** ⚠️ Di `claim-facin\issues\`, **5 dari 114 AC**
tidak tertutup — sebab spec di sana punya bab **"Butir yang belum punya sasaran uji"**.
⭐ `[terverifikasi]` **Spec modul ini TIDAK punya bab serupa** — 15 bab AC-nya seluruhnya perilaku
yang dapat diuji.

---

#### ⭐ Dua uji yang MENGUNCI keputusan, bukan menguji fungsi

⛔ **Keduanya wajib ada, dan sebab keberadaannya wajib tertulis di tiketnya.**

| Uji | Tiket | Kenapa |
| --- | --- | --- |
| ⭐ **Pengeluaran pengaju — DUA ARAH** | **01** | ⛔ Arah kedua membuktikan pengaju di jenjang **kedua TETAP di daftar** — ⭐ **mengunci K9** supaya tidak "terperbaiki" diam-diam |
| ⭐ **Salinan jabatan** | **11** | ⛔ Naik jabatan ⇒ catatan lama **tidak berubah** — ⭐ inilah yang **membedakan salinan dari rujukan** |

---

#### Bukti berkas lain tidak disentuh

| Berkas / folder | Keadaan |
| --- | --- |
| ⛔ `claim-life\` · `komite-claim-life\` · `premiumlist-life\` · `endorsement-life\` | ✅ **NOL disentuh** |
| `komite-claim-facin\spec.md` dan 5 berkas grilling/keputusan | ✅ **NOL disunting** — hanya dibaca |
| `claim-facin\issues\00-14` | ✅ **NOL disunting** |
| `claim-facin\issues\urutan-tiket.md` | ⭐ disunting **hanya pada blok RALAT** |
| `claim-prop\` · `komite-claim-prop\` · `docs\adr\` · `CLAUDE.md` · struktur mana pun | ✅ **NOL disunting** |
| korpus | ✅ md5 tidak berubah |

⛔ Kode **NOL** · DDL **NOL** · `CREATE TABLE` **NOL** · jalur berkas Go **NOL** · nomor baris XML
**NOL**.

# Komite Claim Life

Jumlah tiket: **10**

## Komite Claim Life - 01 - Terima kasus dari Claim Life + inbox komite per posisi

**Status:** ready-for-agent

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)** · **CL-01** (kerangka aplikasi + seam API) · **kontrak muatan penyerahan CL-10**
— bukan penyelesaian CL-10. Lihat §"Kontrak muatan penyerahan".

#### Hasil & nilai pengguna

Sebagai **anggota komite**, kasus yang diserahkan Claim — Life muncul di **inbox saya** — dan hanya
kasus pada tingkat yang saya miliki — lengkap dengan rincian klaim dan baris `AdjustmentList` yang
diputuskan. Saya tahu apa yang harus saya kerjakan tanpa melihat pekerjaan orang lain.
*(User story 1, 2, 7 di spec)*

Ini **tracer bullet** konteks Komite: menembus React → HTTP → handler → service → repository →
Oracle, memakai ulang seam yang sudah ditetapkan **CL-01**.

#### Area codebase

`internal/models` (kasus komite + koleksi `KomiteList` per tingkat), `internal/repository` (pembacaan
kasus + roster), `internal/services` (penyaringan inbox per posisi), `internal/handlers` (endpoint
inbox + detail kasus), `frontend/` (halaman inbox komite + layar detail).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Flow/KomiteLife_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` Assignment "KomiteRouter", `pyImplementation = WorkList`, `pyRouteTo = Custom` — antrean **per-pengguna** |
| `Komite Claim Life/FlowAction/ViewTransferDtl.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `VIEWTRANSFERDTL` / `RULE-OBJ-FLOWACTION`, 31.963 byte | layar detail kasus |
| `Komite Claim Life/Section/ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION`, 1.125.234 byte | bentuk tampilan kasus — **berkas terbesar modul** |

`[terverifikasi]` Kelas kasus: `ASM-FW-GCNMFW-Work-KomiteLife`. Muatan penyerahan membawa `CLMNO`,
`KomiteCount`, `KomiteLoop`, `IndexAdjustment`, `IndexPremiumList`, dan `KomiteList` berisi
`KomiteID` / `IDKomite` / `KomiteAproval` / `KomiteEmail` — ditambah nilai klaim, `CURRENCY`, dan
`STS_REJECT` saat penyerahan (**ADR-U-0001**, keputusan Ronde 2 Q15 Claim Life).

#### ADR terkait

**ADR-U-0001** (kontrak batas — penyerahan masuk), **ADR-U-0014** (inbox hanya menampilkan kasus sesuai
posisi roster), **ADR-U-0003** (uang non-float pada tampilan nilai klaim), **ADR-U-0011** (yang
diputuskan adalah **baris `AdjustmentList`**, bukan klaim).

#### Acceptance criteria

- [ ] Kasus yang diserahkan Claim — Life dapat dibuka lewat API dan menampilkan klaim beserta
      **baris `AdjustmentList`** yang diputuskan.
- [ ] Inbox seorang anggota komite **hanya** memuat kasus pada tingkat yang ia miliki — bukan
      seluruh antrean komite. *(AC 10 spec)*
- [ ] Muatan penyerahan yang diterima memuat nilai klaim dan `CURRENCY`; nilai uang ditampilkan
      tanpa melewati *binary floating point*. *(AC 17 spec)*
- [ ] `KomiteLoop` yang diterima dari muatan dipakai apa adanya sebagai batas tangga — **tidak**
      dihitung ulang di konteks ini.
- [ ] ⚠️ Muatan ber-`KomiteLoop < 1` atau ber-`KomiteList` kosong **ditolak di lapisan layanan**;
      `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` **tidak dibuat**, dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`
      tidak terisi. **Menolak bukan menghitung ulang** — checkbox di atas tetap berlaku.
      *(AC 4 spec; `[keputusan work owner]` — penjaga berlapis; `[terverifikasi]` Pega tidak
      memuat gerbang ini di sisi mana pun)*
- [ ] Status ditampilkan sebagai kata (Outstanding / Aksep / Ditolak), **bukan** nama field
      `STS_REJECT` dan bukan angka. *(AC 29 spec)*
- [ ] Ada satu test ujung-ke-ujung yang menggerakkan sistem lewat HTTP dan memeriksa hasilnya lewat

##### Penyimpanan kasus komite ⚠️ BARU 2026-09-16 — spec §9

- [ ] ⚠️ Menerima penyerahan **membuat `T_GENERAL_KOMITE`** (header) **beserta satu baris
      `T_KOMITE_KOMITELIST` per anggota roster**, masing-masing ber-`KOMITE_APROVAL = 0` dan
      ber-`KOMITE_URUT` sesuai jenjangnya. *(AC 30 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Menerima penyerahan **melahirkan baris baru di `T_WORK_CLAIM`** — `ID` = identitas kasus
      komite (**teks berformat `KMT-xxxxxx`**), `COVER_KEY` = `ID` baris klaim — dan header
      `T_GENERAL_KOMITE` **memakai `ID` yang sama persis** (**shared primary key**).
      **REVISI 2026-09-18:** ⛔ **tidak ada kolom `WORK_CLAIM_ID`** — dibuang; hubungannya dijamin
      oleh ID identik, bukan kolom penyambung. Test yang menemukan kolom `WORK_CLAIM_ID` **gagal**.
      *(§9 spec; REVISI 2026-09-17 — menggantikan `T_WORK_CLAIM.KMT_NO` yang dibuang;
      REVISI 2026-09-18 — shared PK, `[keputusan work owner]`)*
- [ ] ⚠️ **Penunjuk dua arah konsisten dalam SATU transaksi**: `T_GENERAL_KOMITE.ADJUSTMENT_ID` dan
      `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` terisi bersama. Test yang menemukan salah satunya kosong
      sementara yang lain terisi **gagal**. *(AC 32 spec)*
- [ ] `KOMITE_LOOP` diisi **jumlah tingkat** hasil hitung roster aktif; `KOMITE_COUNT` dimulai pada
      tingkat pertama. *(AC 30 spec; tiket 02)*
- [ ] ⚠️ Roster dibaca dari master `EMAILKOMITE`; master itu **tidak ditulis**. *(`[data DBA]`)*
      HTTP, terhadap skema uji Oracle.

#### Kontrak muatan penyerahan — disepakati di muka

Tiket ini bergantung pada **bentuk data** yang diserahkan Claim — Life, **bukan** pada selesainya
CL-10. Bentuk itu ditetapkan di sini agar kedua konteks dapat dikembangkan **paralel**, dengan
muatan penyerahan **di-fake di seam** sampai CL-10 nyata.

| Bagian muatan | Isi |
| --- | --- |
| Penunjuk baris | `IndexAdjustment`, `IndexPremiumList` — baris `AdjustmentList` yang diputuskan |
| Nilai klaim | jumlah klaim yang menjadi dasar pita roster |
| `CURRENCY` | mata uang nilai klaim (**ADR-U-0003**) |
| `KomiteList` | roster tingkat — tiap entri berisi `KomiteID`, `IDKomite`, `KomiteAproval`, `KomiteEmail` |
| `KomiteLoop` | jumlah tingkat tangga — **sudah dihitung** Claim Life |
| Status baris saat serah | `STS_REJECT` pada saat penyerahan |

`[terverifikasi]` Bentuk ini terbaca dari `Claim Life/Activity/CreateKMTLife_Act.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 121.652 byte) —
properti `childPageKomite.*` sebelum `Call pxAddChildWork`. Tiga tambahan (nilai klaim, `CURRENCY`,
`STS_REJECT`) adalah `[keputusan work owner]` dari Ronde 2 Q15 Claim Life (**ADR-U-0001**).

**Perubahan pada bentuk ini adalah perubahan kontrak lintas konteks**, bukan perubahan internal.

#### Catatan — batas kepemilikan dengan Claim Life

`[terverifikasi]` **Perhitungan roster dan `KomiteLoop` adalah kode Claim Life, bukan Komite.**
`Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
`GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) dan `Claim Life/Activity/CreateKMTLife_Act.xml` keduanya
berada di modul `Claim Life`.

Konsekuensi: AC "COUNT roster", "nilai mutlak klaim negatif", dan "gagal bila roster kosong" adalah
AC **CL-10**. Di konteks ini ia **ekspektasi kontrak** — Komite **menerima** `KomiteList` dan
`KomiteLoop` yang sudah dihitung, dan **tidak menghitung ulang**.

`[terbuka]` **OQ-007 / OQ-021** — pemetaan `KomiteID` → identitas akun bergantung pada model RBAC
yang belum ditetapkan. **Asumsi tiket ini: satu `KomiteID` memetakan ke satu identitas akun.**

`[terbuka]` **OQ-035** — `UpdateWorkObject` dan `serviceInsertArasapasClaimLife_act` dipanggil dari
Komite tetapi salinannya ada di modul lain. **Bukan pemblokir isi**; perilakunya sudah terbaca.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Komite Claim Life - 02 - Mesin tangga — rute, naik tingkat, berhenti saat Tolak

**Status:** ready-for-agent

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)**, 01 (terima kasus + inbox per posisi)

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin tangga persetujuan **berjalan sendiri**: kasus berpindah ke
anggota berikutnya saat disetujui, dan berhenti saat ditolak — tanpa siapa pun perlu mengatur
urutannya manual. *(User story 3, 13, 14, 15 di spec)*

Ini **inti konteks Komite**: tangga yang di Pega tidak terlihat di graf, dinyatakan eksplisit.

#### Area codebase

`internal/models` (`KomiteCount`, `KomiteLoop`, entri `KomiteList` per tingkat),
`internal/services` (mesin tangga: pemilihan tingkat, transisi, penghentian),
`internal/handlers` (endpoint simpan keputusan), `frontend/` (layar keputusan).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`, 26.387 byte | `[terverifikasi]` `param.AssignTo = .KomiteID` (baris ~294), bergerbang `.KomiteAproval == 0` (baris ~382) |
| `Komite Claim Life/When/IsKomiteLoop.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `ISKOMITELOOP` / `RULE-OBJ-WHEN` | `[terverifikasi]` `.AcceptStatus = "1"` **DAN** `.KomiteCount <= .KomiteLoop` |
| `Komite Claim Life/Flow/KomiteLife_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` 1 Assignment + 1 Decision + 4 connector — **satu assignment yang di-loop** |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 515.675 byte | `[terverifikasi]` jejak per tingkat (baris ~792, ~908); `KomiteCount + 1` (baris ~9020) |
| `Komite Claim Life/Section/ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION` | `[terverifikasi]` dropdown `AcceptStatus` wajib, baris 32607 — `pyFormat = pxDropdown`, `pyRequired = true` |

#### ADR terkait

**ADR-U-0011** (unit keputusan = baris `AdjustmentList`), **ADR-U-0007** (jejak tiap transisi),
**ADR-U-0001** (tangga adalah isi konteks ini; batasnya ke Claim Life di tiket 05).

#### Acceptance criteria

- [ ] Kasus baru dirutekan ke baris roster **pertama** yang ber-`KomiteAproval == 0`. *(AC 1 spec)*
- [ ] Keputusan **Setuju** pada tingkat bukan-terakhir menaikkan `KomiteCount` satu dan **tidak**
      menyentuh tabel akseptasi. *(AC 5 spec)*
- [ ] Keputusan **Tolak** menghentikan tangga pada tingkat mana pun ia terjadi. *(AC 6 spec)*
- [ ] Tangga berlanjut **hanya** bila keputusan Setuju **dan** `KomiteCount <= KomiteLoop`.
- [ ] Tiap tingkat menghasilkan satu entri berisi **keputusan, komentar, dan waktu**. *(AC 7 spec)*
- [ ] Nilai keputusan adalah **enum tertutup `{1 = Setuju, 2 = Tolak}`**; nilai lain **ditolak
      terang-terangan**, bukan menghentikan tangga diam-diam. *(AC 35 spec; `[keputusan work owner]`)*
- [ ] Tidak ada padanan `TransferType` di kode — lihat catatan. *(AC 26 spec)*

##### Penyimpanan tangga ⚠️ BARU 2026-09-16 — spec §9

- [ ] ⚠️ **`KOMITE_COUNT` dan `KOMITE_LOOP` di-persist di `T_GENERAL_KOMITE`** — bukan hanya hidup di
      halaman kerja. Tingkat berjalan terbaca kembali setelah proses dimulai ulang. *(AC 30 spec;
      penyimpangan sadar 1)*
- [ ] ⚠️ Keputusan di setiap tingkat **menulis baris `T_KOMITE_KOMITELIST`** yang bersesuaian —
      `KOMITE_APROVAL`, `KOMITE_COMMENT`, `DATE_APPROVE`. *(AC 31 spec)*
- [ ] ⚠️ Baris yang ditulis dipilih lewat **`KOMITE_URUT` + `DATA_KOMITE_ID`**, bukan lewat indeks
      posisi. *(AC 34 spec; penyimpangan sadar 2)*

#### Catatan — `TransferType` tidak direplikasi

`[terverifikasi]` `TransferType` **tidak pernah diisi** di korpus: **nol `Property-Set`** terhadapnya
di `Claim Life` maupun `Komite Claim Life`. Ia hanya **dibaca** di dua tempat — precondition kedua
`KomiteRouter` (baris ~442) dan visible-when `.TransferType==2` di `ShowTransfer.xml`
(baris ~29177). Karena tak pernah di-set, `TransferType == '2'` **selalu FALSE** → cabang mati.

`[keputusan work owner]` **Dibuang**, beserta bagian UI `ShowTransfer` yang bergantung padanya.
Routing tingkat digerakkan **hanya** oleh `.KomiteAproval == 0`. (**OQ-033** tertutup.)

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Komite Claim Life - 03 - Penegakan wewenang per `KomiteID` + eskalasi naik satu tingkat

**Status:** ready-for-agent

**Blocked by:** 02 (mesin tangga) — gerbang perlu tindakan nyata untuk dijaga

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin keputusan pada sebuah tingkat **hanya** dapat disimpan oleh
anggota komite yang memang ditunjuk tingkat itu — dan penolakannya terjadi di lapisan layanan,
sehingga tidak dapat dilewati dengan memanggil API langsung.

Sebagai **admin komite**, saya ingin memindahkan kasus **naik satu tingkat** bila anggota tingkat
berjalan berhalangan, supaya kasus tidak macet menunggu satu orang.
*(User story 12, 16, 30 di spec)*

⚠️ **Penyimpangan sadar.** Sistem lama tidak menegakkan apa pun.

#### Area codebase

`internal/services` (pemeriksaan wewenang sebelum menyimpan keputusan; aksi eskalasi),
`internal/handlers` (identitas pemanggil), `internal/repository` (pencatatan eskalasi),
`frontend/` (kontrol eskalasi untuk admin; kontrol keputusan hanya bagi yang berwenang).

#### Rule Pega sumber

| Rule | Identitas | Keadaan sekarang |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` hanya **menempatkan** tugas: `param.AssignTo = .KomiteID`, `pyImplementation = WorkList` |
| `Komite Claim Life/Section/ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION` | `[terverifikasi]` keputusan masuk lewat dropdown wajib — **tanpa** pemeriksaan pemilik |

`[terverifikasi]` **Tidak ada penegakan di sistem lama.** Sapuan 47 berkas modul tidak menemukan satu
pun pemeriksaan bahwa penyimpan keputusan adalah pemilik `KomiteID` tingkat berjalan; dan
`AcceptStatus` **tidak ditulis rule mana pun** (nol `<PropertiesName>…AcceptStatus</PropertiesName>`).
Di Pega, routing `WorkList` hanya menaruh kasus di antrean — **penempatan, bukan penegakan**.

`[terverifikasi]` Komite Claim Life **tidak memuat identitas orang ter-hardcode** (0 berkas) —
berbeda dari Komite Claim FacIn (4) dan Komite Claim Prop (3). Wewenangnya memang berbasis data.

#### ADR terkait

**ADR-U-0014** (penegakan per `KomiteID` + pengecualian eskalasi naik satu tingkat), **ADR-U-0002**
(peran ditegakkan di lapisan layanan), **ADR-U-0007** (eskalasi dan perubahan roster masuk jejak
audit).

#### Acceptance criteria

- [ ] Pengguna yang **bukan** pemilik `KomiteList(KomiteCount).KomiteID` **ditolak** saat menyimpan
      keputusan, meskipun ia dapat membuka kasusnya. *(AC 8 spec)*
- [ ] Penolakan terjadi **di lapisan layanan**, dan tetap terjadi meskipun kontrol UI ditampilkan.
      *(AC 9 spec)*
- [ ] Admin dapat memindahkan kasus **naik satu tingkat**; eskalasi **turun** ditolak.
      *(AC 11 spec)*
- [ ] Eskalasi tercatat: **siapa** memindahkan, **kapan**, dari tingkat mana ke tingkat mana.
      *(AC 12 spec)*
- [ ] Pemutus di tingkat **yang sama** yang bukan pemilik `KomiteID` tetap ditolak — eskalasi bukan
      pintu belakang.
- [ ] Perubahan roster tercatat, karena ia mengubah **siapa yang berwenang**.
- [ ] Tidak ada nama orang ter-hardcode di lapisan mana pun.

#### Catatan

⚠️ `[keputusan work owner]` Eskalasi **memperpendek** tangga: tingkat yang dilewati tidak pernah
memberi keputusan, sehingga entri `KomiteAproval` untuk tingkat itu kosong. `KomiteCount` ikut naik.
Bentuk pencatatannya adalah keputusan implementasi.

`[terbuka]` **OQ-007 / OQ-021** — model RBAC lintas konteks belum ditetapkan. **Asumsi tiket ini:
satu `KomiteID` memetakan ke satu identitas akun yang dapat diautentikasi.** Bila pemetaan ternyata
banyak-ke-banyak, aturan penegakan perlu ditinjau ulang.

`[terverifikasi]` Identity & Access ditandai **ABSENT** dari korpus (`discovery/context-map.md`) —
nol rule otorisasi. Ia dibangun dari nol.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Komite Claim Life - 04a - 04a: Nomor akseptasi — lahir sekali, di keputusan final

**Status:** ready-for-agent

**Blocked by:** 02 (mesin tangga) — tingkat final harus dapat dikenali lebih dulu

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin nomor akseptasi lahir **sekali saja** — pada keputusan Setuju di
tingkat tertinggi — sehingga satu klaim tidak pernah memperoleh dua nomor, dan tingkat-tingkat di
bawahnya tidak membakar sequence. *(User story 17, 19 di spec)*

#### Area codebase

`internal/repository` (pemanggilan kedua SQL penomoran), `internal/services` (gerbang tingkat final;
perakitan nomor per `Type`), `internal/handlers` (nomor tampil pada respons keputusan),
`frontend/` (nomor akseptasi terlihat setelah keputusan final).

#### Rule Pega sumber

Rantai **aktif** di dalam step 4 `KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, versi 2026-09-15).
Nomor sub-step dari `<pyStepPageReference>`:

| Sub-step | Rule | Deskripsi |
| --- | --- | --- |
| 4.1 | `RDB-List` | get tanggal produksi |
| **4.7** | `GetKodeProdLife_SQL` (`ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETKODEPRODLIFE_SQL` / `RULE-CONNECT-SQL`) | **AMBIL KODE PROD** — prefix |
| **4.9** | `GetSequenceNumber_SQL` (`ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL`) | **generate MM.YYYY DAN SEQUENCE** |
| 4.11 | `Property-Set` | cabang `QR,QP` |
| 4.12 | `Property-Set` | cabang `TR,TP` |

`[terverifikasi]` Gerbang tingkat final: `pyWorkPage.AcceptStatus = 1 && pyWorkPage.KomiteCount ==
pyWorkPage.KomiteLoop` (baris **5695**).

`[data DBA]` `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(p_class, p_jenis, p_proddate, OUT p_bulan,
OUT p_seq_number)` — sequence per `(class, jenis, tahun)` di tabel `GENERATE_SEQUENCE_NUMBER`
(PK komposit), memakai `SELECT … FOR UPDATE`; `p_jenis` membedakan retro/non-retro; periode digulir
lewat `POOLDATA.TANGGAL_CLOSING`, dengan aturan cutover `TRUNC(now) <= 02/01/2026` → periode
`12.2025`; keluaran `LPAD(seq, 5, '0')`.

#### ADR terkait

**ADR-U-0006** (penomoran lewat stored procedure — **jangan replikasi logikanya**), **ADR-U-0011**
(keputusan per baris), **ADR-U-0015** (batas transaksi dipegang Go; commit segera setelah nomor
terbentuk agar lock `FOR UPDATE` lekas lepas).

#### Acceptance criteria

- [ ] Nomor akseptasi dibuat **hanya** pada keputusan **Setuju** di tingkat terakhir
      (`KomiteCount == KomiteLoop`). *(AC 13 spec)*
- [ ] Tingkat bukan-terakhir **tidak** memanggil jalur penomoran sama sekali — sequence tidak
      bergerak. *(AC 5 spec)*
- [ ] Keputusan **Tolak** tidak menghasilkan nomor akseptasi.
- [ ] Nomor diperoleh lewat rantai `GetKodeProdLife_SQL` → `GetSequenceNumber_SQL`; aplikasi
      **tidak** memuat logika pembentukan format nomor. *(AC 14 spec; **ADR-U-0006**)*
- [ ] Prefix diperoleh lewat **lookup** ke `POOLDATA.KODE_PRODUKSI`, **tidak** ditanam sebagai
      konstanta.
- [ ] Percabangan per `Type` (`QR`/`QP` versus `TR`/`TP`) menentukan skema nomor yang dipakai.
- [ ] Commit terjadi **segera setelah nomor terbentuk**, sehingga lock `SELECT … FOR UPDATE` tidak
      menahan pemutus lain.
- [ ] Dua keputusan final berurutan menghasilkan dua nomor berbeda.
- [ ] Tidak ada padanan `Generate_NoAccept_KMT_Life` / `Generate_NoAccept_KMT_LifeRetro` di kode —
      `[terverifikasi]` keduanya `RequestType` di step ter-remark (`//`)
      `Komite Claim Life/Activity/KomitePostAdjustment.xml`
      (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`),
      baris 1769 dan 1986 (blok remark di 1725 dan 1943). *(AC 28 spec)*

#### Catatan — dua rule penomoran lama tidak dimigrasikan

`[terverifikasi]` `Generate_NoAccept_KMT_Life` (step **4.4**) dan `Generate_NoAccept_KMT_LifeRetro`
(step **4.5**) **mati**, dibuktikan **dua sinyal bebas**:

| Sinyal | Hasil |
| --- | --- |
| `<pyStepsBlockName>` | berisi `//` → **REMARK** (baris 1725 dan 1943) |
| Asimetri indeks rujukan | `<RequestType>` 1 × , terindeks `<pyRuleName>` **0 ×** |

Pola yang sama sudah ditetapkan **ADR-U-0006** untuk Claim Life. **Jangan dimigrasikan.**

`[terverifikasi]` Step **4.6** `Set Nilai Akseptasi` juga ter-remark (baris 2160).

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Komite Claim Life - 04b - 04b: Rekam akseptasi — satu jalur simpan berparameter status

**Status:** ready-for-agent

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)**, 04a (nomor akseptasi)

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin hasil keputusan komite tersimpan sebagai **rekam akseptasi**
lengkap dengan nomornya — lewat **satu** jalur simpan, apa pun keputusannya — sehingga jalur aksep
dan jalur tolak tidak pernah menyimpang diam-diam satu sama lain.
*(User story 17, 20, 21 di spec)*

⚠️ **Penyimpangan sadar dari paritas struktural.**

#### Area codebase

`internal/models` (bentuk rekam akseptasi), `internal/repository` (penulisan ke tabel akseptasi),
`internal/services` (satu jalur simpan berparameter status), `frontend/` (dokumen akseptasi dapat
diunduh setelah keputusan final).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | blok PL/SQL `INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE (…)` — **55 kolom terbaca langsung** |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` step **4.15** | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | **Insert ke OS** — jalur aksep, gerbang `AcceptStatus = 1 && KomiteCount == KomiteLoop` (baris 5695) |
| idem, step **5.6** | idem | **Insert ke OS** — jalur reject, gerbang `AcceptStatus==2 && KomiteCount == KomiteLoop` (baris 8119) |
| idem, step **4.17 / 4.18** | `Call PrintAkseptasiPDF` / `Call LoadDocumentLife_ACT` | dokumen akseptasi |

⚠️ `[terverifikasi]` **Rule ini dipakai bersama Claim — Life** — hash ternormalisasi `c50bfd9a12`
identik di kedua modul, dan **tidak terdaftar** di register konflik OQ-011. Mengubah bentuknya adalah
**perubahan kontrak lintas konteks**.

`[data DBA]` `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`: `STS_REJECT NUMBER(38)`; delapan kolom uang
bertipe `NUMBER` **tanpa presisi**; `CURRENCY VARCHAR2(100)`; `TYPE VARCHAR2(10)`.

#### ADR terkait

**ADR-U-0003** (uang non-float — `NUMBER` tanpa presisi menuntut desimal presisi arbitrer),
**ADR-U-0011** (unit keputusan = baris), **ADR-U-0001** (tabel akseptasi adalah kontrak bersama),
**ADR-U-0015** (batas transaksi dipegang Go).

#### Acceptance criteria

- [ ] Rekam akseptasi ditulis **hanya** pada tingkat terakhir — baik untuk keputusan Setuju maupun
      Tolak. *(AC 13 spec)*
- [ ] Penyimpanan memakai **satu jalur berparameter status**; **tidak ada dua jalur kembar** di kode.
      *(AC 16 spec)*
- [ ] Tidak ada nilai uang sebagai *binary floating point* di lapisan mana pun maupun di JSON.
      *(AC 17 spec)*
- [ ] Nilai retro yang ditulis adalah nilai **apa adanya dari data policy** — tidak ada logika
      penukaran. *(AC 27 spec; **OQ-065**)*
- [ ] Dokumen akseptasi dapat dihasilkan dan diunduh setelah keputusan Setuju final.
- [ ] Perubahan bentuk rekam akseptasi diperlakukan sebagai **perubahan kontrak lintas konteks**,

##### Penyimpanan keputusan final ⚠️ BARU 2026-09-16 — spec §9

- [ ] ⚠️ Keputusan final **juga menetapkan `T_GENERAL_KOMITE.ACCEPT_STATUS`** (`1` aksep / `2` tolak).
      *(AC 30 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Keputusan final **mengisi `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`** sehingga baris adjustment
      menunjuk kasus komite yang memutuskannya. *(AC 32 spec; User story 42)*
- [ ] ⚠️ Ketiganya — rekam akseptasi, `ACCEPT_STATUS`, dan `KOMITE_ID` — ditulis dalam **satu
      transaksi**; kegagalan pada salah satunya **membatalkan seluruhnya**. *(AC 32 spec)*
      dan ditandai demikian di kode.

#### Catatan — mengapa dua blok disatukan

`[terverifikasi]` Kedua blok tulis didahului rangkaian precondition **yang sama persis**, dan
mengisi himpunan properti **identik** — diff strukturalnya **nol beda**. Yang membedakan hanya nilai
status yang ditulis dan satu gerbang tambahan.

Itu adalah **salin-tempel jalur aksep menjadi jalur tolak**, sepanjang ribuan baris.
`[keputusan work owner]` **Disatukan** — dua blok kembar adalah tempat divergensi diam tumbuh.

#### Catatan — langkah tukar-RetroName tidak direplikasi

`[terverifikasi]` Step **4.14** (jalur aksep, baris 4937–5185) dan step **5.5** (jalur reject,
baris 7663–7911), keduanya berjudul **"Tukar SecurityReinsurer dengan RetroName"**, ber-
`<pyStepsBlockName>//` → **REMARK**.

Ketiga precondition retro berada **di dalam** langkah yang mati itu — termasuk cutover
`ProdDateTime < "20250207T000000.000 GMT"` (baris 5144 dan 7870).

`[keputusan work owner]` Nilai `RetroID`/`RetroName` **sudah di-set di langkah sebelumnya**, mentah
dari data policy. Sistem baru memakainya **apa adanya**; **tidak ada logika penukaran yang
direplikasi**. (**OQ-065** tertutup, **OQ-034** diperkuat.)

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Komite Claim Life - 05 - Jalur balik — tulis `STS_REJECT` ke dua tingkat baris

**Status:** ready-for-agent

**Blocked by:** 04b (rekam akseptasi — satu jalur simpan)

#### Hasil & nilai pengguna

Sebagai **ReasLifeSPV di Claim — Life**, saya ingin hasil keputusan komite **memantul ke baris
`AdjustmentList` saya**, sehingga status di sistem saya selalu mencerminkan keputusan terakhir — dan
bila ditolak, saya dapat mengajukan ulang dengan baris baru. *(User story 18 di spec)*

Ini **kontrak keluar** konteks Komite. Yang **membaca dan menampilkannya** adalah Claim — Life
(**CL-11**); yang **menulisnya** adalah tiket ini.

#### Area codebase

`internal/services` (penerapan hasil keputusan ke dua tingkat baris), `internal/repository`
(penulisan status baris), `internal/models` (bentuk baris `AdjustmentList` dan `PremiumListDetail`).

Tidak menyentuh `frontend/` konteks ini — tampilannya milik Claim — Life.

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 515.675 byte | `[terverifikasi]` menulis `STS_REJECT` ke **dua tingkat baris** — `…PremiumListDetail(idx).STS_REJECT` dan `…AdjustmentList(idx).STS_REJECT` |

`[terverifikasi]` Sensus penulis: **6** `Property-Set` bernilai `1` dan **2** bernilai `2`.
Gerbang baris yang boleh diubah: `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"`.
Gerbang tingkat final: baris **5695** (aksep) dan **8119** (tolak).

`[keputusan work owner]` Pemetaan hasil: **Setuju → `STS_REJECT = 1`**, **Tolak → `STS_REJECT = 2`**.
⚠️ Nama field menyesatkan — nilai `1` berarti **diaksep**.

#### ADR terkait

**ADR-U-0011** (unit keputusan = baris `AdjustmentList`; `PremiumListDetail` dan header klaim adalah
**cerminan** baris terakhir), **ADR-U-0001** (Kontrak 2 — jalur balik; `AcceptStatus` adalah kosakata
konteks ini dan **dipetakan di batas**, tidak disimpan sebagai status kedua di Claim Life),
**ADR-U-0007** (penerapan hasil merekam siapa + kapan).

#### Acceptance criteria

- [ ] Hasil **Setuju** di tingkat terakhir membuat baris yang diserahkan berstatus **Aksep**.
- [ ] Hasil **Tolak** di tingkat terakhir membuat baris berstatus **Ditolak**, dan klaim **tetap**
      dapat menerima baris baru di sisi Claim — Life.
- [ ] `STS_REJECT` tertulis pada **dua tingkat baris** — `PremiumListDetail` dan `AdjustmentList` —
      dengan nilai **yang sama**, dalam satu operasi. *(AC 15 spec)*
- [ ] Perubahan status **hanya** terjadi ketika putaran mencapai **tingkat terakhir**
      (`KomiteCount == KomiteLoop`). *(AC 13 spec)*
- [ ] Hanya baris yang **masih Outstanding**, sudah dipilih, dan **belum bernomor akseptasi** yang
      dapat diubah oleh hasil keputusan.
- [ ] Setiap penerapan hasil merekam **pelaku dan waktu**. *(**ADR-U-0007**)*
- [ ] `AcceptStatus` **tidak** diteruskan ke Claim — Life sebagai status kedua; ia dipetakan ke
      `STS_REJECT` **di batas ini**. *(**ADR-U-0001**)*

#### Catatan — kepemilikan lintas konteks

`[terverifikasi]` **Penulisan ini milik Komite, bukan Claim Life** — pembuktiannya ada di korpus:
`KomitePostAdjustment.xml` berada di modul `Komite Claim Life` dengan class
`ASM-FW-GCNMFW-Work-KomiteLife`.

Konsekuensinya untuk **CL-11**: tiket itu **membaca dan menampilkan** hasil yang ditulis di sini, dan
**tidak menulis `STS_REJECT` sendiri**. Cakupan CL-11 sudah diselaraskan.

#### Perintah verifikasi

```
go test ./internal/...
make check
```

## Komite Claim Life - 06 - Transactional outbox — keputusan + daftar efek dalam satu transaksi

**Status:** ready-for-agent

**Blocked by:** 04b (rekam akseptasi — satu jalur simpan)

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin keputusan komite dan **antrean efek keluarnya** tersimpan dalam
**satu transaksi** — sehingga tidak ada efek yang hilang bila proses mati tepat setelah keputusan
disimpan. *(User story 24, 28 di spec)*

Ini fondasi jaminan "wajib berhasil"; pengirimannya sendiri ada di tiket 07.

#### Area codebase

`internal/models` (entri outbox: efek apa, muatan apa, keadaan apa), `internal/repository`
(penulisan outbox dalam transaksi yang sama dengan keputusan), `internal/services` (penyusunan
daftar empat efek saat keputusan final).

Belum menyentuh `frontend/` — permukaan manusia ada di tiket 08.

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` step **6** | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` `Obj-Save` — **persist**; keempat efek berjalan **sesudahnya** |

`[terverifikasi]` Urutan efek yang harus masuk antrean, dibaca dari `<pyStepPageReference>` versi
2026-09-15:

| Urutan | Step | Efek |
| ---: | ---: | --- |
| 1 | **8** | `Call InsertJsonClaimLife_Act` |
| 2 | **10** | `Call serviceInsertArasapasClaimLife_act` |
| 3 | **11** | `Call SendEmailKlaimLife` |
| 4 | **12** | `Call HitServiceToKasirKMTLife_Act` — Kasir/pembayaran |

`[terverifikasi]` Step **10** dan **12** digerbangi `AcceptStatus = 1 && KomiteCount == KomiteLoop`
(baris 8648, 8887) — hanya pada keputusan Setuju di tingkat final.

#### ADR terkait

**ADR-U-0015** (efek keluar wajib berhasil — transactional outbox, at-least-once; **menyimpang dari
ADR-U-0008**), **ADR-U-0013** (alamat endpoint di-lookup runtime, bukan disimpan di outbox sebagai URL),
**OQ-013** (batas transaksi dipegang Go).

#### Acceptance criteria

- [ ] Keputusan dan **daftar keempat efeknya** tersimpan dalam **satu transaksi database**.
      *(AC 19 spec)*
- [ ] Bila transaksi gagal, **tidak ada** keputusan tersimpan **dan tidak ada** entri outbox — tidak
      ada keadaan separuh.
- [ ] Bila proses mati tepat setelah commit, antrean efek **tetap ada** dan dapat diambil worker.
- [ ] Tiap entri outbox membawa **ID idempoten unik** sejak dibuat. *(AC 21 spec)*
- [ ] Entri outbox menyimpan **kunci kategori** endpoint, **bukan URL** — alamat di-resolve saat
      kirim (**ADR-U-0013**).
- [ ] Efek hanya diantrekan pada keputusan yang benar-benar final; keputusan di tingkat bukan-
      terakhir tidak menghasilkan entri outbox.
- [ ] Keputusan dilaporkan **tersimpan**, belum **tuntas** — kedua keadaan itu dibedakan.
      *(**ADR-U-0015**)*

#### Catatan — mengapa Komite berbeda dari Claim Life

⚠️ **ADR-U-0015 menyimpang dari ADR-U-0008.** Di Claim — Life, efek keluar boleh gagal tanpa memblokir.
Di sini tidak — pemicunya **Kasir**, integrasi **pembayaran** yang `[terverifikasi]` tidak ada di
Claim Life.

Kegagalan diam di Kasir berarti **klaim disetujui tetapi tidak pernah sampai ke pembayaran**, dan
tidak ada yang tahu.

`[keputusan work owner]` Step **14** `SetInformationData` **bukan** efek keluar — hanya temporary
penampung hasil submit. **Tidak** masuk outbox.

#### Perintah verifikasi

```
go test ./internal/...
make check
```

## Komite Claim Life - 07 - Worker pengirim — retry + anti-dobel Email & Kasir

**Status:** ready-for-agent

**Blocked by:** 06 (transactional outbox)

#### Hasil & nilai pengguna

Sebagai **Finance**, saya ingin klaim yang disetujui komite **pasti** sampai ke Kasir — dan **tidak
pernah dua kali**. Sebagai **operator sistem**, saya ingin efek yang gagal diantre ulang sampai
berhasil, bukan menguap. *(User story 22–25, 27 di spec)*

#### Area codebase

worker (proses terpisah — pengambil entri outbox, pengirim, penjadwal retry),
`internal/services` (interface keempat efek keluar; aturan cek-status-sebelum-kirim-ulang),
`internal/repository` (lookup `M_LINK_SERVICE`; pembaruan keadaan entri outbox).

Tidak menyentuh `frontend/` — permukaan manusia ada di tiket 08.

#### Rule Pega sumber

| Efek | Rule | Identitas |
| --- | --- | --- |
| 1 · step 8 | `Komite Claim Life/Activity/InsertJsonClaimLife_Act.xml` | `RULE-OBJ-ACTIVITY` |
| 2 · step 10 | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | `[terverifikasi]` **satu-satunya salinan di korpus** — OQ-035 |
| 3 · step 11 | `Komite Claim Life/Activity/SendEmailKlaimLife.xml` | `RULE-OBJ-ACTIVITY` |
| 4 · step 12 | `Komite Claim Life/Activity/HitServiceToKasirKMTLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `HITSERVICETOKASIRKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY` — **Kasir/pembayaran** |
| resolusi alamat | `Komite Claim Life/Activity/GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` — `Obj-Browse` `M_LINK_SERVICE` pada `(KATEGORI_1, KATEGORI_2)`, ambil `.URL`, lalu `Connect-REST` |
| jejak Kasir | `POOLDATA.DIRECTTOKASIR_LOG` | `[terverifikasi]` jejak di sisi database |

`[terverifikasi]` Kunci kategori Arasapas: `Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`
(`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`).

#### ADR terkait

**ADR-U-0015** (at-least-once; ID idempoten unik; cek status sukses sebelum kirim ulang untuk Email &
Kasir — **menyimpang dari ADR-U-0008**), **ADR-U-0013** (alamat di-lookup runtime; **dilarang** URL
sebagai literal, konstanta, maupun env var), **ADR-U-0005** (flag lingkungan menggerbangi efek keluar).

#### Acceptance criteria

- [ ] Keempat efek dikirim sampai **berhasil**, atau ditandai **perlu intervensi**. *(AC 20 spec)*
- [ ] Tiap kiriman membawa **ID idempoten unik**, sehingga penerima dapat menolak duplikat.
      *(AC 21 spec)*
- [ ] Sebelum mengirim ulang **Email** atau **Kasir**, sistem memeriksa status "sudah terkirim
      sukses"; bila sudah, **tidak dikirim lagi**. *(AC 22 spec)*
- [ ] Alamat endpoint di-resolve **runtime** dari `M_LINK_SERVICE` lewat `(KATEGORI_1, KATEGORI_2)`;
      **tidak ada URL** sebagai literal, konstanta, maupun env var. *(AC 24 spec; **ADR-U-0013**)*
- [ ] **Kunci kategori tidak ditemukan** dibedakan dari **jaringan gagal**: yang pertama tidak
      diulang berkali-kali, yang kedua diulang.
- [ ] Keputusan dilaporkan **tuntas** hanya setelah keempat efek berhasil. *(AC 25 spec)*
- [ ] **Semua** klaim menjalankan keempat efek — tidak ada pengecualian berdasarkan identitas retro.
      *(AC 26 spec; **OQ-064**)*
- [ ] Kegagalan pengiriman **tidak** membatalkan keputusan yang sudah tersimpan. *(AC 18 spec)*

#### Catatan — gerbang EXIT retro dibuang

⚠️ `[terverifikasi]` Di korpus, **step 9** adalah langkah gerbang tersendiri berlabel
`EXIT JIKA RETROID "L0000141"` (baris 8483), dengan precondition
`RetroID=="L0000141" || SecurityReinsurerID=="L0000134"` (baris 8525) dan `RetroID=="1000013"`
(baris 8548). Ia berdiri **di antara** efek 1 dan efek 2 — saat memicu, step 10–14 tidak berjalan.

`[keputusan work owner]` **Dibuang** (**OQ-064**). **Konsekuensi yang diterima:** klaim ber-retro
tersebut yang selama ini dikecualikan kini menjalankan **seluruh** efek keluar, **termasuk Kasir**.
Ini **perubahan perilaku yang menyentuh uang**, diterima sadar.

`[dugaan]` Bahwa kode transisi numerik step 9 berarti "Exit Activity" **tidak terbukti dari ekspor**
— nilai yang terbaca (`2`) muncul juga di step 8 yang bukan EXIT. Yang terbukti adalah label
penulisnya dan bentuk langkahnya. Tidak berpengaruh pada implementasi: gerbangnya dibuang.

#### Catatan — seam kedua

`[asumsi]` Perilaku retry dan anti-dobel **tidak teramati lewat API HTTP saja**, karena pengiriman
berjalan di proses terpisah. Tiket ini mengasumsikan **seam kedua di batas pengirim outbox**, dengan
keempat layanan luar **di-fake**. Bila seam itu tidak disetujui, cara mengujinya perlu ditetapkan
ulang sebelum tiket dikerjakan.

`[terbuka]` **OQ-002** — kontrak layanan **Kasir** tidak ada di korpus: bentuk permintaan, makna
jawaban, dan apakah ia menghormati ID idempoten **belum diketahui**. Tiket ini menetapkan **jaminan
pengirimannya**, bukan bentuk pesannya.

`[terbuka]` **OQ-035** — `serviceInsertArasapasClaimLife_act` satu salinan dipakai dua konteks.
**Bukan pemblokir isi.**

#### Perintah verifikasi

```
go test ./internal/...
make check
```

## Komite Claim Life - 08 - Status "perlu intervensi" di UI Komite + laporan harian

**Status:** ready-for-agent

**Blocked by:** 07 (worker pengirim — retry + anti-dobel)

#### Hasil & nilai pengguna

Sebagai **anggota komite**, saya ingin melihat kasus yang efek keluarnya **gagal dan perlu
intervensi** — dan sebagai **operator**, saya ingin daftar itu sampai ke saya setiap hari tanpa
harus mencarinya. *(User story 26 di spec)*

Tanpa tiket ini, "wajib berhasil" (**ADR-U-0015**) hanya **memindahkan kegagalan diam ke antrean** —
tidak ada yang melihatnya. Itulah sebabnya ia tiket terpisah dari 07: **ini permukaan untuk manusia,
bukan mesin.**

#### Area codebase

`internal/models` (keadaan "perlu intervensi" pada entri outbox), `internal/services` (kriteria
kapan sebuah entri masuk keadaan itu; penyusunan laporan harian), `internal/handlers` (endpoint
daftar kasus bermasalah), `frontend/` (penanda status di inbox dan layar kasus).

#### Rule Pega sumber

**Tidak ada padanan di korpus.** `[terverifikasi]` Sistem lama **tidak memuat** permukaan pemantauan
apa pun untuk kegagalan efek keluar. Yang ada hanya pencatatan di sisi database:

| Objek | Peran |
| --- | --- |
| `POOLDATA.DIRECTTOKASIR_LOG` | `[terverifikasi]` jejak panggilan Kasir |
| `POOLDATA.MONITORING_KLAIM_LOG` | `[terverifikasi]` log layanan — dipakai jalur Claim Life |

Keduanya **log**, bukan antrean kerja: tidak ada yang menampilkannya kepada manusia, dan tidak ada
keadaan "belum selesai" yang dapat ditindaklanjuti.

**Tiket ini karena itu perilaku baru sepenuhnya**, bukan paritas — ia konsekuensi langsung
**ADR-U-0015**.

#### ADR terkait

**ADR-U-0015** (status "perlu intervensi" di UI Komite + laporan harian adalah **syarat pengaman
wajib**, bukan tambahan), **ADR-U-0007** (kegagalan masuk jalur audit, bukan hanya log layanan),
**ADR-U-0014** (siapa yang melihat daftar itu mengikuti wewenang).

#### Acceptance criteria

- [ ] Entri outbox yang **tidak pulih** setelah retry masuk keadaan **"perlu intervensi"** yang
      eksplisit — bukan diam di antrean. *(AC 23 spec)*
- [ ] Kasus ber-status itu **terlihat di UI Komite**, menempel pada kasusnya, bukan di halaman
      terpisah yang harus dicari.
- [ ] Tampilan menyebut **efek mana** yang gagal (InsertJson / Arasapas / Email / Kasir) dan
      **sejak kapan**.
- [ ] Ada **laporan harian** berisi seluruh kasus ber-status itu.
- [ ] Laporan harian tetap terkirim **meskipun kosong** — ketiadaan laporan tidak boleh ambigu
      dengan ketiadaan masalah.
- [ ] Keputusan yang efeknya belum tuntas dilaporkan **tersimpan**, bukan **tuntas**.
      *(AC 25 spec; **ADR-U-0015**)*
- [ ] Kegagalan juga tercatat di **jalur audit**, bukan hanya di log layanan. *(**ADR-U-0007**)*

#### Catatan

⚠️ **Kasus yang paling penting terlihat adalah Kasir.** `[terverifikasi]` `HitServiceToKasirKMTLife_Act`
tidak ada di Claim Life — ia khas Komite, dan ia memicu jalur **pembayaran**. Klaim yang disetujui
tetapi gagal sampai ke Kasir adalah kegagalan yang paling mahal bila tidak terlihat.

`[keputusan work owner]` Siapa yang menindaklanjuti keadaan "perlu intervensi" — anggota komite,
admin, atau operator — **belum ditetapkan**. Tiket ini menampilkannya; **alur penanganannya**
keputusan terpisah. Tidak memblokir: yang penting kegagalan berhenti tersembunyi.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Komite Claim Life - 09 - Riwayat tangga persetujuan di UI

**Status:** ready-for-agent

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)**, 03 (penegakan wewenang + eskalasi) — riwayat harus memuat eskalasi juga

#### Hasil & nilai pengguna

Sebagai **pengguna mana pun**, saya ingin melihat **riwayat lengkap tangga** pada satu kasus — siapa
memutuskan apa, kapan, dengan komentar apa — sehingga saya paham mengapa kasus berada di keadaannya
sekarang tanpa bertanya kepada orang.

Sebagai **auditor**, saya ingin eskalasi ikut terbaca di riwayat yang sama, sehingga pengetatan
wewenang tidak dilubangi diam-diam. *(User story 9, 10, 29, 30, 31 di spec)*

#### Area codebase

`internal/services` (penyusunan riwayat dari entri per tingkat + catatan eskalasi),
`internal/handlers` (endpoint riwayat kasus), `frontend/` (tampilan riwayat pada layar kasus).

#### Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 515.675 byte | `[terverifikasi]` menulis entri per tingkat |

`[terverifikasi]` Tiga `Property-Set` mengisi elemen list yang diindeks pencacah tangga, dengan
`Local.Komite = pyWorkPage.KomiteCount`:

| Properti | Nilai |
| --- | --- |
| `…AdjustmentList(IndexAdjustment).KomiteList(Local.Komite).KomiteAproval` | `pyWorkPage.AcceptStatus` |
| `…KomiteList(Local.Komite).KomiteComment` | `pyWorkPage.Comment` |
| `…KomiteList(Local.Komite).DateApprove` | `@CurrentDateTime()` |

Hal yang sama juga ditulis ke `pyWorkPage.KomiteList(Local.Komite)`.

**Jadi setiap tingkat tangga menghasilkan satu entri berisi keputusan, komentar, dan waktu** —
riwayatnya sudah ada di data; tiket ini membuatnya terbaca.

#### ADR terkait

**ADR-U-0007** (jejak audit setiap transisi — termasuk eskalasi), **ADR-U-0014** (eskalasi adalah
tindakan yang direkam: siapa, kapan, dari tingkat mana ke tingkat mana), **ADR-U-0011** (riwayat
melekat pada **baris `AdjustmentList`**, bukan pada klaim).

#### Acceptance criteria

- [ ] Riwayat menampilkan **tiap tingkat** tangga secara berurutan, dengan keputusan, komentar, dan
      waktu. *(AC 7 spec)*
- [ ] Keputusan ditampilkan sebagai **kata** (Setuju / Tolak), bukan angka dan bukan nama field.
      *(AC 29 spec)*
- [ ] **Eskalasi ikut terbaca** di riwayat yang sama: siapa memindahkan, kapan, dari tingkat mana ke
      tingkat mana. *(AC 12 spec)*
- [ ] Tingkat yang **dilewati** karena eskalasi terlihat sebagai dilewati — bukan hilang tanpa jejak.
- [ ] Riwayat melekat pada **baris `AdjustmentList`** yang diputuskan; satu klaim dengan beberapa
      baris menampilkan riwayat per baris.
- [ ] Riwayat dapat dibaca **tanpa** wewenang memutuskan — melihat bukan memutuskan.

##### Sumber riwayat ⚠️ BARU 2026-09-16 — spec §9

⚠️ **Koreksi premis.** Catatan lama menyebut *"riwayatnya sudah ada di data"* — yang dimaksud adalah
**page runtime Pega**. Di sistem baru page itu **dibuang**; riwayat punya tabelnya sendiri.

- [ ] ⚠️ Riwayat tangga dibaca dari **`T_KOMITE_KOMITELIST` diurut `KOMITE_URUT`** — **bukan** dari page
      runtime maupun JSON. Test yang menemukan pembacaan dari page **gagal**. *(AC 33 spec;
      penyimpangan sadar 1)*
- [ ] Tiap baris riwayat menampilkan **anggota pemutus**, **keputusannya**, **komentarnya**, dan
      **tanggal putus** — langsung dari kolom `KOMITE_ID`, `KOMITE_APROVAL`, `KOMITE_COMMENT`,
      `DATE_APPROVE`. *(AC 31 spec)*
- [ ] Tingkat yang **belum memutus** terbaca sebagai `KOMITE_APROVAL = 0`, dan ditampilkan sebagai
      kata — bukan angka. *(AC 29, 31 spec)*

#### Catatan

`[keputusan work owner]` Eskalasi menaikkan `KomiteCount` tanpa tingkat yang dilewati memberi
keputusan, sehingga entri `KomiteAproval` untuk tingkat itu **kosong**. Bentuk penampilannya —
"dilewati (eskalasi)" atau serupa — adalah keputusan implementasi; yang mengikat adalah **ia tidak
boleh tampak seperti tingkat yang belum diputuskan**.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

# Komite Claim Prop

Jumlah tiket: **14**

## Komite Claim Prop - 01 - Terima penyerahan dari Claim Prop — kasus komite lahir

**Status:** ready-for-agent

**Blocked by:** **00 (skema — PREFACTOR)** · **kontrak muatan penyerahan Claim Prop** — bukan
penyelesaian modul Claim Prop.

#### Hasil & nilai pengguna

Sebagai **petugas klaim**, ketika saya menyerahkan satu baris penyesuaian ke komite, **kasus komite
terbentuk** lengkap dengan daftar penyetuju, pencacah tangga, dan penunjuk ke baris penyesuaian
saya. Saya tahu penyerahan saya diterima.

*(User story 5, 19 di spec)*

Ini **tracer bullet** pertama yang terlihat: menembus penerimaan penyerahan → penyimpanan →
pembacaan kembali.

#### ⚠️ Penyerahan membekukan baris penyesuaian induknya — **PERILAKU BARU**

`[keputusan work owner]` 2026-09-19 — **saat kasus komite lahir, baris penyesuaian induknya
ditandai BEKU di sisi Claim Prop:** tidak dapat diubah dan tidak dapat dihapus, **selamanya**.
**Bekunya tidak mencair meskipun komite menolak** — perbaikan dilakukan dengan **baris penyesuaian
baru**, bukan dengan menyunting baris lama.

⛔ **PENEGAKANNYA MILIK CLAIM PROP, BUKAN TIKET INI.** Daftar penyesuaian adalah milik modul Claim
Prop, jadi penolakan sunting dan penolakan hapus dibangun di sana — pada tiket **08** *(baris
adjustment)* dan tiket **11** *(penyerahan ke Komite)* modul itu. **Tiket ini hanya pemicunya:**
ia mencatat penyerahan, dan pencatatan itulah yang menjadikan baris induknya beku.

⚠️ **Ini perilaku BARU, bukan paritas.** `[terverifikasi]` Di seluruh ekspor modul ini: **nol
pemeriksaan, nol pesan kesalahan, nol penanganan** untuk baris penyesuaian yang berubah atau hilang
sesudah diserahkan. **Aturan ini dibangun, bukan dimigrasikan.**

**Lingkup:** Claim Prop dan Claim Non Prop; `[keputusan work owner]` **Claim Non Prop menyusul —
fokus sekarang Claim Prop saja. Claim — Life tidak termasuk.**

⛔ **Tidak berlaku surut** — kasus komite lama hasil migrasi tidak terlindungi; lihat butir **14**
register `[terbuka]` dan tiket **13**.

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` | `[terverifikasi]` **langkah 13** menetapkan jenis pengajuan dan penunjuk posisi baris · **langkah 14** menyalin isi baris penyesuaian · **langkah 15** menyalin catatan komite · **langkah 22.1** membangun daftar penyetuju, tiap baris ber-keputusan awal `0` |
| idem | `[terverifikasi]` pencacah **jumlah penyetuju** diisi dari **jumlah baris daftar penyetuju**, sekali saat kasus dibuat; pada jalur tutup/tolak yang memakai satu penyetuju tetap nilainya `1` |

⭐ `[terverifikasi]` **Daftar penyetuju ditetapkan saat kasus dibuat dan tidak berubah di tengah
jalan** — nol penulis di modul komite ini.

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Pencacah jumlah penyetuju **disimpan**, diisi ulang di titik yang sama dengan Pega — **bukan** dihitung ulang tiap dibaca |
| 2026-09-18 | **Aturan awalan**: data berawalan titik dan halaman kasus = **komite**; berawalan halaman induk = **klaim**, **dibaca saja** |

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 1 · 8 · 9 · 16 · **76** · **77** · **78**

- [ ] Penyerahan satu baris penyesuaian **melahirkan** baris komite di tabel lintas-lini **dan**
      header kasus komite ber-kunci sama.
- [ ] Penyerahan **membuat satu baris penyetuju per anggota daftar**, masing-masing ber-keputusan
      awal kosong dan ber-urutan sesuai jenjangnya.
- [ ] Pencacah **jumlah penyetuju** terisi dari jumlah baris daftar; pencacah **penyetuju keberapa**
      dimulai di penyetuju pertama.
- [ ] Penunjuk ke baris penyesuaian induk terisi, dan penunjuk balik dari sisi klaim terisi
      **dalam satu transaksi yang sama**. Test yang menemukan salah satunya kosong **gagal**.
- [ ] Muatan penyerahan tanpa daftar penyetuju, atau dengan jumlah penyetuju kurang dari satu,
      **ditolak** — kasus komite **tidak dibuat**.
- [ ] ⚠️ Sesudah kasus komite lahir, baris penyesuaian induknya **tidak dapat disunting** dan
      **tidak dapat dihapus** — termasuk lewat kaskade hapus klaim. Ditegakkan **di sisi Claim
      Prop**; tiket ini menguji bahwa **penandanya benar-benar terpasang** saat penyerahan tercatat.
- [ ] ⚠️ Bekunya **tidak mencair** sesudah komite menolak; test yang **berhasil menyunting** baris
      sesudah penolakan **gagal**.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Berapa kali perulangan berputar** di Pega belum terbaca — menyentuh berapa kali nilai
  penyesuaian ditulis saat kasus dibuat.

#### Seam & verifikasi

Memakai ulang seam Claim Prop. Muatan penyerahan **dipalsukan di seam** sampai sisi Claim Prop
nyata, supaya kedua modul bisa dikerjakan paralel.

## Komite Claim Prop - 02 - Kotak kerja penyetuju dan rute giliran

**Status:** ready-for-agent

**Blocked by:** **01 (kasus komite lahir)**

#### Hasil & nilai pengguna

Sebagai **penyetuju komite**, kasus yang menunggu keputusan saya muncul di **kotak kerja saya** —
dan **hanya saat giliran saya**. Saya tidak melihat pekerjaan orang lain, dan saya tidak diminta
menilai sesuatu yang belum dilihat penyetuju sebelum saya.

*(User story 1, 2, 3 di spec)*

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Flow/KomiteTreaty_Flow.xml` · `ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITETREATY_FLOW` | `[terverifikasi]` satu **assignment** bernama `KomiteRouter`, dirutekan **khusus** ke satu pengguna, berbentuk antrean per-pengguna |
| `Komite Claim Prop/Activity/KomiteRouter.xml` | `[terverifikasi]` 8 langkah, **6 di-remark**. Yang hidup **langkah 6** dan anaknya **6.1**: bergerbang *"keputusan penyetuju masih kosong"*, lalu menetapkan tujuan rute dari **akun operator** penyetuju yang sedang giliran. Pencacah dinaikkan **langkah 5** |

⭐ `[terverifikasi]` **Tangga persetujuan yang lama — antrean berjenjang `komitepnc1..4` — sudah
mati**; keenam langkah itulah yang di-remark. ⛔ **Jangan dipindahkan.**

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Hak akses per activity diabaikan** — daftar hak akses Pega menyebut kelas **tanpa menyebut hak**, jadi ia tidak menegakkan apa pun |

⚠️ Konsekuensinya: **siapa boleh melihat apa ditentukan aturan peran sistem baru**, bukan disalin
dari Pega. Tiket ini hanya menegakkan **giliran**, bukan wewenang.

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 3 · 4 · 10 · 11 · 14 · 67

- [ ] Kotak kerja seorang penyetuju **hanya** memuat kasus yang sedang **gilirannya**.
- [ ] Penyetuju yang belum tiba gilirannya **tidak melihat** kasus itu.
- [ ] Penyetuju yang sudah memutuskan **tidak melihatnya lagi** di kotak kerja.
- [ ] Tujuan rute diambil dari **akun operator** penyetuju, bukan jabatan dan bukan nama antrean.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Pewarisan kelas kerja** tidak ada di ekspor Pega — menyentuh bagaimana properti kasus komite
  diwarisi. Tidak menghambat tiket ini.

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 03 - Layar komite — 93 kolom, tiga wajah menurut jalur

**Status:** ready-for-agent

**Blocked by:** **01 (kasus komite lahir)**

#### Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya membuka kasus dan melihat **data klaim induknya** — tertanggung,
nomor polis, tanggal kejadian, penyebab, lokasi — **beserta nilai penyesuaian yang diajukan**.
Layarnya **menyesuaikan diri dengan jenis pengajuan**, jadi saya tidak dibingungkan kolom yang tidak
relevan.

*(User story 3, 4, 6, 7, 8, 9, 10, 11, 12 di spec)*

#### Bentuk layar

`[terverifikasi]` **93 kolom** — **62 milik kasus komite**, **31 milik kasus klaim induk**. Angka
ini **SAH** `[keputusan work owner]` 2026-09-18.

⭐ **Tiga wajah menurut jenis pengajuan.** Dari 62 kolom komite, **24 hanya tampil pada jalur
penyesuaian**. Tiga judul bagian layar juga bergerbang jenis pengajuan.

⚠️ **Dua kolom yang wajib ikut terbawa meski di luar hitungan 93:** **"dibuat oleh"** dan
**"tanggal dibuat"**. Keduanya **tampil di layar**; keduanya di luar hitungan hanya karena aturan
pencatatan memisahkan properti bawaan Pega.

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/FlowAction/ViewTransferDtl.xml` · `ASM-FW-GCNMFW-WORK-KOMITETREATY!VIEWTRANSFERDTL` | `[terverifikasi]` **pembungkus, nol kolom sendiri**. Kontraknya: penyiap memuat → layar menampilkan → penyimpan menyimpan |
| `Komite Claim Prop/Section/ShowTransfer.xml` · `ASM-FW-GCNMFW-WORK-KOMITETREATY!SHOWTRANSFER` | `[terverifikasi]` 348 sel; **dua famili gerbang** — satu di tingkat sel, satu di tingkat layout. **12 sel di balik gerbang mati** |
| `Komite Claim Prop/Activity/SetKomiteList_Act.xml` | `[terverifikasi]` penyiap layar; **langkah 1** menghentikan activity bila jenis pengajuan bukan penyesuaian |
| `Komite Claim Prop/Section/ViewDetailInterest.xml` | `[terverifikasi]` layar kedua — **3 kolom, ketiganya hanya-baca**, dan **ketiganya sudah ada di layar utama** |

⭐ `[terverifikasi]` **Layar kedua tidak menambah satu kolom pun.**

#### ⭐ Daftar penyetuju DITAMPILKAN — **perilaku BARU**

`[keputusan work owner — atas rekomendasi asisten]` 2026-09-19 — **layar komite menampilkan
daftar penyetuju**: **urutan** · **siapa** · dan untuk yang **sudah memutuskan** — **keputusan,
catatan, dan tanggalnya**. Yang **belum** memutuskan tampil **tanpa isi**.

⚠️ **Di Pega ia tidak tampil sama sekali.** `[terverifikasi]` Ketiga penggerak tangga —
`KomiteLoop`, `KomiteCount`, `KomiteList` — **nol di antaranya tampil di layar**, padahal
`KomiteList` justru berisi daftar penyetuju berikut keputusan, catatan, dan tanggalnya. **Akibatnya
di Pega: setiap penyetuju memutuskan tanpa konteks penyetuju sebelumnya.**

⛔ **NOL kolom basis data baru.** Ketiganya sudah tersimpan di baris penyetuju — kesembilan kolom
yang dikunci tiket **00**. Ini **murni penambahan tampilan**, dan ia **menutup US 3 dan US 4**,
dua user story yang sebelumnya **tidak ditutup satu AC pun**.

⚠️ **Tandanya `atas rekomendasi asisten`, jadi murah dicabut** — pencabutannya menyentuh **AC 79**
dan **AC 80**, satu baris bab 11, satu keputusan bab 10, dan dua butir uji di bawah. ⛔ **Nol kolom
tersentuh.** Alasan lengkapnya di `spec.md` **bab 4**.

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Properti jenis treaty **dibuang** — hanya ada di balik gerbang mati, tidak pernah tampil |
| 2026-09-18 | Nomor klaim dan dua kotak total estimasi di balik gerbang mati **dibuang** |
| 2026-09-18 | **Angka 93 · komite 62 · klaim induk 31 SAH** |
| 2026-09-18 | Data klaim **dibaca** dari Claim Prop, **tidak digarap** modul ini |

#### ⚠️ TITIK YANG SENGAJA DIUBAH — ketelitian tampilan

`[keputusan work owner]` 2026-09-18 — **angka di layar ditampilkan dengan 4 angka di belakang
koma**. Pega tidak punya aturan tampilan yang seragam. **Ini perubahan sadar.**

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 2 · 18 · 19 · 20 · 21 · 22 · 26 · 27 · 28 · 41 · **79** · **80**

- [ ] Membuka kasus komite menampilkan data klaim induk dan nilai penyesuaian yang diajukan.
- [ ] Pada jalur **penyesuaian** ke-24 kolom bergerbang **tampil**; pada jalur **tolak** dan
      **tutup** ke-24 kolom itu **tidak tampil**.
- [ ] Ketiga judul bagian layar berganti sesuai jenis pengajuan.
- [ ] Kolom yang tidak boleh diubah penyetuju tampil sebagai **bacaan saja**.
- [ ] **"Dibuat oleh"** dan **"tanggal dibuat"** tampil di layar.
- [ ] Properti jenis treaty, nomor klaim di balik gerbang mati, dan dua kotak total estimasi
      **tidak dibuat sama sekali**. Test yang menemukannya **gagal**.
- [ ] Angka uang di layar tampil dengan **4 angka di belakang koma**.
- [ ] ⚠️ Layar menampilkan **daftar penyetuju** berikut **urutan** dan **siapa**-nya, untuk
      **seluruh** tingkat — bukan hanya yang sudah lewat. Test yang tidak menemukannya **gagal**
      *(AC 79 spec)*
- [ ] ⚠️ Untuk penyetuju yang **sudah memutuskan**, layar menampilkan **keputusan, catatan, dan
      tanggal**-nya; yang **belum** memutuskan tampil **tanpa isi**. Test yang menemukan kolom
      basis data baru untuk ini **gagal** — datanya sudah ada di baris penyetuju *(AC 80 spec)*

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Dari mana layar kedua dibuka** belum terbaca — layar kedua tidak dirujuk layar utama,
  pembungkusnya, maupun berkas alur.

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 04 - Keputusan penyetuju tersimpan

**Status:** ready-for-agent

**Blocked by:** **02 (kotak kerja dan rute giliran)** · **03 (layar komite)**

#### Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya menyatakan **setuju atau tolak** dan menulis **catatan**, lalu
menekan kirim. Keputusan saya tersimpan pada baris saya, lengkap dengan **tanggal dan identitas
saya** — tanpa saya perlu mengisinya.

*(User story 13, 14, 18 di spec)*

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Section/ShowTransfer.xml` | `[terverifikasi]` **dua isian selalu dapat disunting dan wajib isi**: keputusan setuju/tolak, dan catatan. Keduanya **tanpa gerbang** |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 6** menulis keputusan, catatan, dan tanggal ke baris penyesuaian induk · **langkah 27** menulis catatan · tanggal diisi **waktu saat itu** |
| idem | `[terverifikasi]` **langkah 4** membuka kasus klaim induk **dengan penguncian**, dan kuncinya **dilepas saat penyimpanan** — bukan dipegang selama komite menimbang |

⭐ `[terverifikasi]` **Penguncian hanya menyelimuti proses simpan**, beberapa detik. Penyetuju bisa
membuka layar berjam-jam tanpa mengunci apa pun.

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Satu kolom tanggal persetujuan saja** — properti tanggal kedua di Pega **sengaja tidak dijadikan kolom**; jangan ditambahkan kembali karena "ada di korpus" |
| 2026-09-18 | Pada jalur **tolak** dan **tutup**, penyetuju memang hanya mengisi **dua** hal ini. **Ditiru apa adanya** |

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 23 · 25 · 29 · 53 · 54

- [ ] Keputusan dan catatan tersimpan pada **baris penyetuju yang bersangkutan**, bukan di header.
- [ ] Tanggal dan identitas penyetuju **terisi otomatis**; penyetuju tidak bisa mengubahnya.
- [ ] Keduanya **wajib isi** — pengiriman tanpa salah satunya **ditolak**.
- [ ] Pada jalur **tolak** dan **tutup**, hanya kedua isian ini yang tersedia.
- [ ] Penguncian kasus klaim induk **hanya berlangsung selama proses simpan**, dan **dilepas**
      sesudahnya. Test yang menemukan kunci tertahan sesudah simpan **gagal**.
- [ ] Hanya **satu** kolom tanggal persetujuan yang ada. Test yang menemukan kolom tanggal kedua
      **gagal**.

#### Butir `[terbuka]` yang menyentuh tiket ini

Tidak ada.

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 05 - Persetujuan bersyarat dan dua penanda usul

**Status:** ready-for-agent

**Blocked by:** **04 (keputusan penyetuju tersimpan)**

#### Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya bisa menandai bahwa persetujuan saya **bersyarat**, lalu menulis
**isi syaratnya**. Saya juga bisa **mengusulkan penutupan** atau **pencadangan** klaim. Keempatnya
hanya muncul ketika memang relevan, jadi layar tidak penuh isian yang tidak berlaku.

*(User story 15, 16, 17 di spec)*

#### Rantai gerbangnya

`[terverifikasi]` Keempat isian ini **bersyarat**, dan rantainya **bertingkat**:

| Isian | Muncul bila |
| --- | --- |
| penanda persetujuan bersyarat | penyetuju **menyetujui** **dan** jalur **penyesuaian** |
| isi syarat | penanda bersyarat **dicentang** |
| usul tutup klaim | jalur **penyesuaian** |
| usul cadangkan klaim | jalur **penyesuaian** |

⛔ **RALAT ke catatan lama:** keempatnya pernah disebut *"selalu dapat disunting"*. **Itu keliru** —
yang selalu dapat disunting hanya **dua** isian di tiket 04.

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Section/ShowTransfer.xml` | `[terverifikasi]` keempatnya bergerbang **famili sel**; kedua penanda usul berupa **kotak centang**, dapat disunting, **tidak wajib isi** |
| idem | `[terverifikasi]` **nol penulis di activity mana pun** — keduanya ditulis **langsung dari layar** |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 24** menulis penanda bersyarat dan isi syaratnya ke baris penyesuaian induk |

#### ⭐ Kedua penanda usul punya DUA tempat di Pega

`[terverifikasi]` **KomitePostAdjustment langkah 11** — *"add propose close and propose reserve"*,
**tanpa gerbang**, jadi selalu berjalan — menyalin keduanya ke **kasus klaim induk dengan nama
berbeda**:

```
kasus komite . usul tutup klaim       ->  kasus klaim . penanda tutup berkas
kasus komite . usul cadangkan klaim   ->  kasus klaim . penanda klaim dicadangkan
```

⚠️ **RALAT ke catatan lama:** pernyataan *"tidak ada tulis-menulis lintas modul"* **hanya benar
untuk layar**. Lewat activity, komite **memang menulis ke kasus klaim induk**.

⛔ Apakah sisi klaim ikut disimpan adalah **urusan modul Claim Prop**, bukan tiket ini.

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Dua penanda usul **disimpan di tabel komite**, karena penyetuju yang mengisinya |
| **2026-09-19** | Keduanya disimpan di **header kasus komite** — *"komite menyimpan catatan usulnya sendiri, terpisah dari nilai akhir yang tercatat di kasus klaim"* |

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 24 · 55 · **74**

- [ ] Penanda bersyarat **hanya muncul** bila penyetuju menyetujui **dan** jalurnya penyesuaian.
- [ ] Isi syarat **hanya muncul** setelah penanda bersyarat dicentang.
- [ ] Kedua penanda usul **hanya muncul** pada jalur penyesuaian.
- [ ] Keempatnya **tidak wajib isi** — pengiriman tanpa mengisinya tetap diterima.
- [ ] Kedua penanda usul tersimpan di **header kasus komite**.
- [ ] Nilai yang sama juga **sampai ke kasus klaim induk**, sesuai perilaku Pega.
- [ ] Kedua penanda usul tersimpan **hanya** sebagai `'1'` atau `'0'`; kotak yang **tidak disentuh**
      tersimpan `'0'`, bukan kosong, dan nilai lain **ditolak**.

#### Butir `[terbuka]` yang menyentuh tiket ini

- ✅ **TERJAWAB** `[keputusan work owner]` 2026-09-19 — **kedua penanda usul disimpan sebagai teks
  satu huruf: `'1'` bila diusulkan, `'0'` bila tidak. Kolomnya `CHAR(1)` dan wajib isi — kotak
  centang yang tidak pernah disentuh tersimpan `'0'`, bukan kosong.** ⚠️ Migrasi data lama: nilai
  yang di Pega bermakna benar/salah dikonversi otomatis menjadi `'1'`/`'0'`, dan baris lama yang
  nilainya kosong menjadi `'0'`.
  > ⛔ **Kalimat lamanya dikutip, tidak dihapus:** *"**Nilai apa yang tersimpan** untuk kedua
  > penanda usul — di Pega ia kotak centang, tetapi nilai tersimpannya **tidak terbaca dari
  > korpus**. ⛔ Jangan ditebak."*
  >
  > ⚠️ Butir ini **tidak pernah masuk register 13 butir** di `spec.md`, jadi menjawabnya
  > **tidak mengubah jumlah register**. ⛔ **Nol butir register ditutup.**

- **Tidak ada butir `[terbuka]` lain yang menyentuh tiket ini.**

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 06 - Tangga maju ke penyetuju berikutnya, atau selesai

**Status:** ready-for-agent

**Blocked by:** **04 (keputusan penyetuju tersimpan)**

#### Hasil & nilai pengguna

Sebagai **pengelola proses**, kasus **otomatis berpindah** ke penyetuju berikutnya setelah satu
penyetuju memutuskan, dan **selesai** ketika penyetuju terakhir menyetujui. Tidak ada langkah manual
dan tidak ada tahap menggantung.

*(User story 19, 20 di spec)*

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Flow/KomiteTreaty_Flow.xml` | `[terverifikasi]` **satu jalur balik**: dari gerbang keputusan **kembali ke assignment yang sama**, bersyarat *"masih ada penyetuju berikutnya"*. Bila tidak, kasus **selesai** dengan status akhir *selesai-tuntas* |
| idem | `[terverifikasi]` **empat bentuk saja**: mulai → assignment → gerbang keputusan → selesai. Aksi pengguna pada layar komitelah yang memindahkan kasus dari assignment ke gerbang |
| `Komite Claim Prop/Activity/KomiteRouter.xml` | `[terverifikasi]` **langkah 5** menaikkan pencacah penyetuju |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 40** juga menaikkan pencacah, tetapi **gerbangnya dimatikan** |

⭐ `[terverifikasi]` **Tangga berputar DI DALAM satu tahap.** Setiap penyetuju adalah **kunjungan
baru ke tahap yang sama**, bukan tahap baru. Yang berubah tiap putaran hanya **kepada siapa kasus
dirutekan** dan **nilai pencacah**.

⚠️ `[terverifikasi]` Pembacaan riwayat pada langkah 8 dan 9 terjadi **jauh sebelum** penambahan
pencacah di langkah 40 — jadi tidak ada pembacaan indeks di luar daftar.

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Pencacah jumlah penyetuju **disimpan** dan **tidak dihitung ulang** tiap dibaca |

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 5 · 6 · 7 · 12 · 15 · 17

- [ ] Sesudah satu penyetuju memutuskan, kasus **berpindah** ke penyetuju berikutnya.
- [ ] Penyetuju terakhir yang **menyetujui** **menyelesaikan** kasus.
- [ ] Kasus yang selesai **tidak muncul lagi** di kotak kerja siapa pun.
- [ ] Pencacah penyetuju naik **tepat satu** per putaran; test yang menemukannya melompat **gagal**.
- [ ] Jumlah putaran **tidak melebihi** jumlah penyetuju yang dibutuhkan.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Arti dua jenis perulangan Pega** dan **halaman apa yang diulang** belum terbaca.
- **Berapa kali perulangan berputar** belum terbaca — menyentuh berapa kali langkah di dalamnya
  berjalan.
- **Urutan pemeriksaan bila dua keluarga gerbang sama-sama terisi** tidak terbaca dari struktur.

⚠️ Ketiganya **tidak menghambat** perilaku tangga di atas, yang terbaca dari berkas alur, bukan dari
perulangan langkah.

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 07 - Penolakan menghentikan tangga dan menolak sisa penyetuju

**Status:** ready-for-agent

**Blocked by:** **06 (tangga maju atau selesai)**

#### Hasil & nilai pengguna

Sebagai **pengelola proses**, **satu penolakan menghentikan rangkaian** — penyetuju berikutnya tidak
diminta menilai sesuatu yang sudah ditolak. Sisa penyetuju **ditandai otomatis**, jadi daftarnya
tidak tertinggal setengah terisi.

*(User story 21, 22 di spec)*

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 12** — bila keputusan **tolak**, melompat ke tanda `EXT`, yaitu **langkah 25** |
| idem | `[terverifikasi]` **langkah 25** memaksa pencacah ke nilai akhir sehingga tangga **berhenti** |
| idem | `[terverifikasi]` **langkah 26** menyisir daftar penyetuju; anaknya **26.1** menandai sisa penyetuju **tertolak otomatis**. Langkah 26 adalah **langkah berulang** |
| idem | `[terverifikasi]` **langkah 23 · 25** menulis status penolakan ke baris penyesuaian induk |

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Lima perintah lompat menggantung dibuang** — kelimanya menunjuk tanda yang tidak ada, dan kelimanya bergerbang mati. ⛔ Tidak ditebak ke mana seharusnya melompat, tidak ditandai cacat, **dibuang** |

⚠️ **Lompatan yang HIDUP tetap dipindahkan.** Yang dibuang hanya yang **mati**. Lompatan ke tanda
`EXT` pada langkah 12 **hidup dan tandanya ada** — ia bagian dari tiket ini.

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 5 · 13 · 62 · 68

- [ ] Satu keputusan **tolak** menghentikan rangkaian; penyetuju berikutnya **tidak menerima** kasus.
- [ ] Seluruh sisa penyetuju **ditandai tertolak**; tidak ada baris penyetuju yang tertinggal
      berkeputusan kosong.
- [ ] Status penolakan sampai ke **baris penyesuaian induk**.
- [ ] Kasus **selesai** sesudah penolakan — tidak menggantung.
- [ ] ⛔ Tidak ada perintah lompat menggantung yang ikut dipindahkan. Test yang menemukan lompatan
      ke tanda yang tidak ada **gagal**.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Berapa kali perulangan berputar** — menyentuh berapa kali penandaan sisa penyetuju dikerjakan.

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 08 - Nomor akseptasi — terbit sekali, di tingkat akhir

**Status:** ready-for-agent

**Blocked by:** **06 (tangga maju atau selesai)**

#### Hasil & nilai pengguna

Sebagai **bagian akseptasi**, nomor akseptasi terbit **tepat sekali**, saat penyetuju terakhir
menyetujui **tanpa syarat**. Tidak ada nomor ganda, dan nomor **tidak terbit** ketika persetujuannya
masih bersyarat.

*(User story 23, 24 di spec)*

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` blok **langkah 16** dan anak **16.1**–**16.9**. **Langkah 16** adalah induk, bergerbang **bukan persetujuan bersyarat**, dan merupakan **langkah berulang** |
| idem | `[terverifikasi]` **16.1 · 16.2 · 16.3 · 16.4 di-remark** — tidak berjalan · **16.5** mengambil kode produksi · **16.7** membangkitkan bulan-tahun dan nomor urut · **16.8** merangkai · **16.9** menuliskan nomor dan menandai status akseptasi |
| idem | `[terverifikasi]` **16.9** bergerbang *"penyetuju terakhir dan disetujui"* |

⭐ **Gerbangnya berlapis dua**: di tingkat induk *"bukan persetujuan bersyarat"*, di tingkat anak
*"penyetuju terakhir dan disetujui"*. Itulah yang membuat nomor terbit **sekali saja**.

**Rule pengambil data yang dipanggil** `[terverifikasi]`:

| Langkah | Rule |
| --- | --- |
| 16.1 *(remark)* | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETTANGGALCLOSING_SQL` |
| 16.4 *(remark)* | `ASM-FW-GCNMFW-INT-V_POLIS!GCNM!GENERATENOACCEPTTREATY` |
| **16.5** | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETKODEPRODNONLIFE_SQL` |
| **16.7** | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETSEQUENCENUMBER_SQL` |

⚠️ Dua rule yang langkahnya di-remark **tidak ada di ekspor modul ini** — konsisten: yang dimatikan
tidak ikut diekspor. ⛔ **Jangan dipindahkan.**

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 31 · 32 · 33 · 34 · 35

- [ ] Nomor akseptasi terbit **tepat sekali** — pada penyetuju terakhir yang menyetujui tanpa syarat.
- [ ] Nomor **tidak terbit** ketika persetujuan **bersyarat**, meski penyetuju terakhir menyetujui.
- [ ] Nomor **tidak terbit** pada penyetuju selain yang terakhir.
- [ ] Menjalankan ulang alur pada kasus yang sudah bernomor **tidak menerbitkan nomor kedua**.
- [ ] Status akseptasi pada baris penyesuaian induk terisi bersama nomornya.
- [ ] ⛔ Kedua rule yang langkahnya di-remark **tidak dipanggil**.

#### Butir `[terbuka]` yang menyentuh tiket ini

Tidak ada.

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 09 - Angka uang — total penyesuaian per mata uang

**Status:** ready-for-agent

**Blocked by:** **03 (layar komite)**

#### Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya melihat **total nilai penyesuaian per mata uang** beserta
**totalnya dalam rupiah**, sehingga saya bisa membandingkan antar mata uang sebelum memutuskan.

*(User story 8, 32, 33 di spec)*

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/SetKomiteList_Act.xml` | `[terverifikasi]` **satu-satunya penghitung uang berkelas komite**. **Langkah 1** menghentikan activity bila jenis pengajuan bukan penyesuaian. **Langkah 6** membuang baris bermata-uang-sama dari daftar sementara |
| idem | `[terverifikasi]` **langkah 7.1.1**, di dalam perulangan, menjumlahkan **per mata uang**: total kotor, total bersih, dan keduanya **dikali kurs** untuk rupiah. **Langkah 7.2** menaruh kedua total pertama ke properti kasus |
| idem | `[terverifikasi]` gerbang penjumlahan: baris **berstatus tolak tidak ikut dijumlahkan**, dan penjumlahan hanya untuk mata uang yang sedang diproses |

⭐ `[terverifikasi]` **Nol perhitungan uang di potongan Java** — seluruh rumus terbaca dari penugasan
properti biasa.

#### Kurs

`[keputusan work owner]` 2026-09-18 — **"ikuti aja query di xml CurrencyStandard."**

`[terverifikasi]` Kurs diambil lewat **fungsi tersimpan basis data** dengan **dua masukan saja**:
kode mata uang dan **tanggal server saat query dijalankan**. Tidak ada nomor kasus, tidak ada lini,
tidak ada tanggal kasus.

⭐ **Kurs yang terkunci adalah kurs TANGGAL PENCARIAN DIJALANKAN** — bukan tanggal kasus dibuat,
bukan tanggal komite menyetujui. ⚠️ Isi fungsi tersimpannya **tidak ada di korpus** dan
`[keputusan work owner]` **tidak diminta ke DBA**.

`[terverifikasi]` Penguncian nilainya terjadi **di luar modul ini** — pola *ambil-hanya-bila-kosong*
ada di modul Claim Prop. **Di dalam modul ini nol pola penguncian untuk uang.**

#### ⚠️ TITIK YANG SENGAJA DIUBAH — ketelitian angka

`[keputusan work owner]` 2026-09-18, diralat 2026-09-19. **Ini perubahan sadar, bukan peniruan.**

`[terverifikasi]` Di Pega **tidak ada pembulatan sama sekali**, dan persen dibagi seratus dengan
**tiga cara berbeda** — dua di antaranya di activity yang sama untuk menghitung hal yang sama.

| Di sistem baru | Ketetapan |
| --- | --- |
| Hitungan | sampai **20 angka di belakang koma**, **nol pembulatan di tengah jalan** |
| Penyimpanan | **mengikuti bentuk kolom yang sudah ada — 20 digit, 8 di belakang koma** |
| Tampilan | **4 angka di belakang koma** |
| Tipe | **desimal, bukan bilangan pecahan biner** |

⚠️ **Pembulatan terjadi di batas penyimpanan, diterima sadar.** Angka persen yang di Pega dihitung
sampai 10 desimal **akan menjadi 8** saat disimpan. **Hasil sistem baru AKAN BERBEDA dari Pega di
angka belakang koma, dan itu disengaja.**

⚠️ `[keputusan work owner]` 2026-09-19 — bila nilai melewati batas digit, basis data **menolak
menyimpan**, bukan membulatkan. **Kegagalan itu harus terlihat.**

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 30 · 36 · 37 · 38 · 39 · 40 · 41 · 42 · 43 · 44 · 71

- [ ] Total per mata uang benar, dan baris **berstatus tolak tidak ikut dijumlahkan**.
- [ ] Total rupiah = total mata uang **dikali kurs** yang tersimpan pada baris itu.
- [ ] Baris bermata-uang-sama **digabung** sebelum dijumlahkan.
- [ ] Pada jalur **bukan penyesuaian**, penghitungan **tidak berjalan sama sekali**.
- [ ] Nilai uang melewati aplikasi **tanpa melewati bilangan pecahan biner**.
- [ ] Hitungan berjalan sampai **20 angka di belakang koma**; **tidak ada pembulatan di tengah**.
- [ ] Nilai yang melewati batas digit **ditolak dengan galat yang terlihat**.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Ketelitian pembagian persen di Pega tidak seragam** — mana yang sah belum ditetapkan. Ini
  menentukan bagaimana **selisih terhadap data lama** diterangkan, bukan apa yang dikerjakan
  sistem baru.
- **Seberapa besar selisih** akibat persen 10 desimal menjadi 8 saat disimpan — **belum diukur**,
  dan **menumpuk lewat perulangan** sebelum sampai ke tiket 12.
- **Berapa kali perulangan berputar** — menyentuh berapa kali angka berubah.

#### Seam & verifikasi

Memakai ulang seam Claim Prop.

## Komite Claim Prop - 10 - Efek keluar ke basis data — akseptasi, riwayat, log

**Status:** ready-for-agent

**Blocked by:** **08 (nomor akseptasi)**

#### Hasil & nilai pengguna

Sebagai **bagian akseptasi**, baris akseptasi tercatat di daftar akseptasi klaim sehingga angkanya
bisa ditelusuri. Sebagai **auditor**, setiap keputusan komite tercatat di riwayat akseptasi, dan
penolakan klaim tercatat tersendiri.

*(User story 25, 30, 31 di spec)*

#### Perilaku Pega yang ditiru

`[terverifikasi]` Empat dari delapan efek keluar bersifat penulisan basis data, di
`Komite Claim Prop/Activity/KomitePostAdjustment.xml`:

| Langkah | Rule yang dipanggil | Yang ditulis |
| --- | --- | --- |
| **17** | `SaveAcceptation_Act` | baris akseptasi ke daftar akseptasi klaim |
| **28** | `InsertJsonClaimTreaty_act` | data klaim bentuk JSON |
| **31** | `InsertLogServiceClaim` | baris log pemantauan |
| **33** | `InsertHistoryAkseptasiPega_Sql` | **6 kolom** riwayat akseptasi |

`[terverifikasi]` Penulisan penolakan klaim ada di `KomitePost_Close` **langkah 13** dan
`KomitePost_Reject` **langkah 13**, keduanya **di dalam langkah berulang**.

⚠️ `[terverifikasi]` Kolom **kotak kerja** pada riwayat akseptasi diisi **literal tetap**, selalu,
tanpa cabang. `[terbuka]` di modul lain: baris komite **ikut terhitung** oleh pembaca riwayat di
modul Fac In / Treaty In. ⛔ **Bukan urusan tiket ini, dan jangan mengusulkan membuang kolomnya.**

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | **Sembilan rule yang menyimpan sendiri tetap menyimpan sendiri.** *"Ditiru apa adanya."* |

> ⚠️ **Risiko diterima sadar:** bila langkah sesudahnya gagal, **kesembilan tabel tetap terisi**,
> dan **tidak ada rule pembatal di korpus** — tidak ada penghapusan, pembatalan, maupun penanda
> batal. **Data separuh jadi adalah perilaku yang disengaja ditiru, bukan cacat yang terlewat.**

#### ⚠️ TITIK YANG SENGAJA DIUBAH — urutan terhadap penyimpanan

`[keputusan work owner]` 2026-09-18. **Ini perubahan sadar, bukan peniruan.**

`[terverifikasi]` Di Pega **kedelapan efek keluar berjalan SEBELUM penyimpanan**. Penyimpanan ada di
**langkah 41**, sedangkan pengiriman ke Kasir di **34** dan email di **35**. **Nol yang sesudah.**

⚠️ Akibatnya di Pega: bila penyimpanan gagal, **uang sudah dikirim, email sudah sampai, dokumen
sudah dibuat, dan tiga tabel log sudah terisi** — sementara kasusnya sendiri tidak tersimpan.

> `[keputusan work owner]` — **URUTAN INI TIDAK DITIRU.** Fakta di atas dicatat sebagai **fakta**,
> **tidak mengikat rancangan**. Urutannya **akan disesuaikan**.
>
> ⛔ **Jangan mengunci urutan Pega sebagai syarat.** ⛔ Urutan penggantinya **belum diputuskan**,
> dan **tidak ditetapkan di tiket ini**.

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 45 · 56 · 57 · 58 · 59 · 62 · 63 · 69 · 70 · 73

- [ ] Nomor akseptasi terbit → baris akseptasi, riwayat, dan log **tercatat**.
- [ ] Riwayat akseptasi terisi **6 kolom**; kolom ketujuh yang tidak pernah diisi jalur komite
      **tetap kosong**.
- [ ] Penolakan klaim tercatat **tersendiri**.
- [ ] Keempat penulisan **menyimpan sendiri** dan **tidak dibatalkan** bila langkah sesudahnya gagal.
- [ ] Kolom kotak kerja terisi **literal tetap**, tanpa cabang.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Kolom akun operator** pada riwayat akseptasi **tidak pernah diisi jalur komite**. Adakah penulis
  lain di luar modul ini, atau kolom itu memang selalu kosong. ⛔ Tidak ditebak, **tidak diusulkan
  dibuang**.
- **Urutan terhadap penyimpanan** belum ditetapkan — lihat blok di atas.

#### Seam & verifikasi

`[keputusan work owner]` 2026-09-18 — **efek keluar diuji dengan layanan sungguhan, bukan pengganti
tiruan**, di **lingkungan uji terpisah, bukan produksi**.
⚠️ **Pengujian urutan efek keluar MENUNGGU** urutan barunya ditetapkan.

## Komite Claim Prop - 11 - Efek keluar — dokumen PDF akseptasi

**Status:** ready-for-agent

**Blocked by:** **08 (nomor akseptasi)**

#### Hasil & nilai pengguna

Sebagai **bagian akseptasi**, dokumen akseptasi **dibuat otomatis** sebagai PDF dan tersimpan, jadi
tidak perlu disusun manual.

*(User story 26 di spec)*

#### Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 21** memanggil pembuat dokumen akseptasi |
| `Komite Claim Prop/Activity/HTMLToPDF.xml` | ⚠️ `[terverifikasi]` **rule bawaan Pega**, berkelas dasar — **bukan buatan Nusantara Re**. ⛔ **Tidak perlu dipindahkan**; yang dipindahkan adalah **pemakaiannya** |
| `Komite Claim Prop/Activity/InsertDocument_Act.xml` | `[terverifikasi]` menyimpan baris dokumen lalu memanggil pengunggah berkas |
| `Komite Claim Prop/Activity/InsertGoogleStorage_Act.xml` | `[terverifikasi]` mengunggah berkas ke penyimpanan luar, lalu mencatatnya ke tabel penyimpanan berkas |
| `Komite Claim Prop/Activity/GetLinkService.xml` | `[terverifikasi]` mengambil **alamat layanan** dari tabel alamat, disaring **dua kunci** — kategori dan sub-kategori — sehingga hasilnya satu baris |

#### ⚠️ Sikap terhadap kegagalan — **berhenti dengan galat**

`[terverifikasi]` Rantai PDF **melempar galat dan berhenti** bila gagal — **berbeda** dari sikap
pengiriman ke Kasir, yang melompat dan melanjutkan (tiket 12). **Empat titik penanganannya:**

| Di mana | Bila gagal |
| --- | --- |
| pembuat PDF, langkah 6 | **melempar galat** bila markup kosong |
| pembuat PDF, langkah 9 | **mencatat galat dan berhenti** bila pembuat tidak menghasilkan isi |
| tiga langkah pengubah berkas jadi teks | **melempar galat** bila lampiran gagal |

⭐ **Modul ini punya dua sikap berbeda terhadap kegagalan, dan keduanya disengaja.**

⚠️ `[terverifikasi]` **Pengunggahan berkas TIDAK punya penanganan gagal.** Catatan pengembangnya
sendiri berbunyi *"FIX ERROR HANDLING"*, tetapi **jejaknya tidak terbaca** di keempat keluarga wadah
yang sudah disisir habis. ⛔ Jangan ditebak.

#### ⚠️ TITIK YANG SENGAJA DIUBAH — urutan terhadap penyimpanan

`[keputusan work owner]` 2026-09-18. **Ini perubahan sadar, bukan peniruan.**

`[terverifikasi]` Di Pega **kedelapan efek keluar berjalan SEBELUM penyimpanan**. Penyimpanan ada di
**langkah 41**, sedangkan pengiriman ke Kasir di **34** dan email di **35**. **Nol yang sesudah.**

⚠️ Akibatnya di Pega: bila penyimpanan gagal, **uang sudah dikirim, email sudah sampai, dokumen
sudah dibuat, dan tiga tabel log sudah terisi** — sementara kasusnya sendiri tidak tersimpan.

> `[keputusan work owner]` — **URUTAN INI TIDAK DITIRU.** Fakta di atas dicatat sebagai **fakta**,
> **tidak mengikat rancangan**. Urutannya **akan disesuaikan**.
>
> ⛔ **Jangan mengunci urutan Pega sebagai syarat.** ⛔ Urutan penggantinya **belum diputuskan**,
> dan **tidak ditetapkan di tiket ini**.

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 48 · 49 · 50 · 51 · 52 · 69 · 70 · 73

- [ ] Nomor akseptasi terbit → dokumen PDF **dibuat** dan **tersimpan**, dan baris dokumen tercatat.
- [ ] Kegagalan pembuatan PDF **menghentikan dengan galat yang terlihat**, bukan diam-diam dilewati.
- [ ] Alamat layanan penyimpanan diambil dengan **dua kunci penyaring**, bukan baris pertama tanpa
      saringan.
- [ ] ⚠️ Kegagalan **pengunggahan berkas** — perilakunya mengikuti apa yang ada, dan bila memang
      tidak tertangani, **itu didokumentasikan**, bukan ditambal diam-diam.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Pengunggahan berkas dan email tanpa penanganan gagal** — belum dibawa ke work owner.
- **Catatan pengembang "FIX ERROR HANDLING" yang tidak berjejak.**
- **Urutan terhadap penyimpanan** belum ditetapkan.
- **Baris mana di tabel alamat layanan yang menunjuk lingkungan uji** belum ditetapkan.

#### Seam & verifikasi

Efek keluar diuji dengan **layanan sungguhan di lingkungan uji terpisah**
`[keputusan work owner]` 2026-09-18.

## Komite Claim Prop - 12 - Efek keluar — Kasir, arasapas, dan email

**Status:** ready-for-agent

**Blocked by:** **08 (nomor akseptasi)** · **10 (efek keluar ke basis data)**

#### Hasil & nilai pengguna

Sebagai **bagian keuangan**, data pembayaran **terkirim ke Kasir** setelah akseptasi, sehingga
pembayaran bisa diproses — dan hanya bila **nomor akseptasi sudah tercatat**. Sebagai **pihak
terkait**, saya menerima **email pemberitahuan** hasil akseptasi tanpa membuka sistem.

*(User story 27, 28, 29 di spec)*

#### Perilaku Pega yang ditiru

| Langkah | Rule | Yang keluar |
| --- | --- | --- |
| **29** | `KonversiKlaim_Act` | panggilan layanan **arasapas** |
| **34** | `HitServiceToKasirKMT_Act` | data pembayaran ke **Kasir** |
| **35** | `SendEmailKlaim_KMT` | **email** ke penerima akseptasi |

`[terverifikasi]` Alamat tujuan diambil dari **tabel alamat layanan**, disaring **dua kunci**.
Pasangan kuncinya berbeda per tujuan: satu untuk Kasir, satu untuk klaim, dua untuk penyimpanan
berkas.

#### Syarat sebelum mengirim ke Kasir

`[data work owner]` 2026-09-18 — pengiriman ke Kasir **hanya boleh jalan bila nomor akseptasi sudah
tercatat** di daftar akseptasi klaim. **Itu syarat resmi yang dipertahankan.**

`[terverifikasi]` Rule pemeriksanya mengambil nomor akseptasi, membuang titiknya, mencari di tabel
akseptasi, dan menyalakan penanda hanya bila ketemu. ⛔ SQL-nya di luar modul ini — **tidak
ditelusuri**.

#### ⚠️ Sikap terhadap kegagalan — **melompat dan melanjutkan**

`[keputusan work owner]` 2026-09-18 — **kiriman ke Kasir gagal → email TETAP dikirim, kasus tetap
jalan, nol percobaan ulang. DITIRU APA ADANYA.**

`[terverifikasi]` Mekanismenya: bila pengiriman gagal, alur **melompat ke langkah pengiriman email**
lalu selesai. Pola serupa pada arasapas, yang melompat ke langkah pencatatan log. Di **dalam** rule
pengirim Kasir ada jalur galat tersendiri yang **mengirim email khusus kegagalan**.

> ⚠️ **Risiko diterima sadar:** kegagalan **tidak terlihat** sampai ada orang yang memeriksa
> belakangan. Tidak ada pemberitahuan bahwa transfer gagal; yang sampai ke penerima justru **email
> keberhasilan akseptasi**. **Selisih antara "email terkirim" dan "uang terkirim" hanya ketahuan
> dari pemeriksaan manual.**

⚠️ `[terverifikasi]` **Email sendiri tidak punya penanganan gagal sama sekali.**

#### ⚠️ TITIK YANG SENGAJA DIUBAH — urutan terhadap penyimpanan

`[keputusan work owner]` 2026-09-18. **Ini perubahan sadar, bukan peniruan.**

`[terverifikasi]` Di Pega **kedelapan efek keluar berjalan SEBELUM penyimpanan**. Penyimpanan ada di
**langkah 41**, sedangkan pengiriman ke Kasir di **34** dan email di **35**. **Nol yang sesudah.**

⚠️ Akibatnya di Pega: bila penyimpanan gagal, **uang sudah dikirim, email sudah sampai, dokumen
sudah dibuat, dan tiga tabel log sudah terisi** — sementara kasusnya sendiri tidak tersimpan.

> `[keputusan work owner]` — **URUTAN INI TIDAK DITIRU.** Fakta di atas dicatat sebagai **fakta**,
> **tidak mengikat rancangan**. Urutannya **akan disesuaikan**.
>
> ⛔ **Jangan mengunci urutan Pega sebagai syarat.** ⛔ Urutan penggantinya **belum diputuskan**,
> dan **tidak ditetapkan di tiket ini**.

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 45 · 46 · 47 · 49 · 50 · 52 · 69 · 70 · 73

- [ ] Nomor akseptasi tercatat → data pembayaran **terkirim ke Kasir**.
- [ ] Nomor akseptasi **belum** tercatat → pengiriman ke Kasir **tidak berjalan**.
- [ ] Kegagalan pengiriman ke Kasir **tidak menghentikan kasus**, dan **email tetap terkirim**.
- [ ] Tidak ada percobaan ulang otomatis.
- [ ] Kegagalan pengiriman **tercatat** di jalur galatnya sendiri.
- [ ] Alamat tujuan diambil dengan **pasangan kunci yang benar** untuk masing-masing tujuan.

#### Butir `[terbuka]` yang menyentuh tiket ini

- **Unggah berkas dan email tanpa penanganan gagal** — belum dibawa ke work owner.
- **Urutan terhadap penyimpanan** belum ditetapkan.
- **Baris mana di tabel alamat layanan yang menunjuk lingkungan uji** belum ditetapkan.
- **Seberapa besar selisih** akibat pembulatan di batas penyimpanan — angka yang dikirim ke Kasir
  berasal dari hitungan di tiket 09.

#### Seam & verifikasi

Efek keluar diuji dengan **layanan sungguhan di lingkungan uji terpisah**
`[keputusan work owner]` 2026-09-18. ⚠️ **Pengujian urutan MENUNGGU** urutan barunya ditetapkan.

## Komite Claim Prop - 13 - Migrasi data lama

**Status:** ready-for-agent

**Blocked by:** **00 (skema — PREFACTOR)**

#### Hasil & nilai pengguna

Sebagai **pengguna sistem baru**, kasus komite lama beserta daftar penyetuju dan keputusannya
**terbaca di sistem baru**, sehingga riwayat tidak hilang saat pindah.

#### Yang dipindahkan

| Dari Pega | Ke sistem baru |
| --- | --- |
| kasus komite | baris komite di tabel lintas-lini + header kasus komite |
| daftar penyetuju | satu baris per penyetuju, beserta keputusan, catatan, dan tanggalnya |
| pencacah tangga | kedua pencacah apa adanya |
| penunjuk ke baris penyesuaian induk | penunjuk baris penyesuaian di header kasus |
| **dua penanda usul** | ⭐ **benar/salah → `'1'`/`'0'` · kosong → `'0'`** — lihat aturan konversi di bawah |

##### ⭐ Aturan konversi dua penanda usul

`[keputusan work owner]` 2026-09-19 — kolom tujuannya **`CHAR(1)`, wajib isi, nilai sah hanya
`'1'` dan `'0'`** *(tiket 00)*. Karena itu:

| Nilai lama di Pega | Menjadi |
| --- | --- |
| bermakna **benar** *(dicentang)* | **`'1'`** |
| bermakna **salah** *(tidak dicentang)* | **`'0'`** |
| **kosong / tidak ada nilainya** | **`'0'`** |

⛔ **Konversinya otomatis dan tanpa pengecualian** — migrasi **tidak boleh** meninggalkan kolom
usul kosong, dan **tidak boleh** menolak baris hanya karena nilainya kosong.

⚠️ **Ini satu-satunya tempat migrasi memperbaiki nilai.** Ia **tidak** melonggarkan keputusan
2026-09-19 tentang penunjuk posisional di bawah — penunjuk tetap **dipindahkan apa adanya**.

#### ⚠️ Kunci lama bersifat POSISIONAL

`[terverifikasi]` Di Pega, kasus komite menunjuk baris penyesuaian induknya lewat **nomor urut
baris**, bukan pengenal. Kuncinya dua bagian: kunci kasus klaim induk, ditambah **posisi** di dalam
daftar penyesuaian klaim itu.

⚠️ **Artinya: bila sebuah baris penyesuaian pernah disisipkan atau dihapus di klaim induk sesudah
kasus komitenya dibuat, kasus komite itu menunjuk baris yang KELIRU — dan tidak ada yang
menyadarinya.** `[terverifikasi]` **Tidak ada penanganan untuk hal itu di Pega** — tidak ada
pemeriksaan, tidak ada pesan.

#### Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| **2026-09-19** | **Penunjuk posisional DIPINDAHKAN APA ADANYA, termasuk yang sudah salah alamat di Pega. Tidak diperiksa dulu, tidak diperbaiki.** Konsisten dengan preseden *"tiru apa adanya"* |

> ⚠️ **Risikonya, ditulis tegas:** **penunjuk yang salah di Pega akan TETAP SALAH sesudah migrasi.**
> Sistem baru **mewarisi kekeliruan itu apa adanya, tanpa penanda.** ⛔ Migrasi **tidak boleh**
> memperbaikinya diam-diam, dan **tidak boleh** menolak baris yang mencurigakan.

#### Yang harus diuji

**Diverifikasi oleh:** spec.md AC 61 · 64 · 65 · 66 · 71 · **75**

- [ ] Kasus komite lama terbaca lengkap dengan daftar penyetuju dan keputusannya.
- [ ] Penunjuk posisional diterjemahkan menjadi penunjuk baris penyesuaian **apa adanya** — termasuk
      yang mengarah ke baris yang berbeda dari yang semestinya.
- [ ] ⛔ Migrasi **tidak menolak** dan **tidak memperbaiki** baris yang penunjuknya mencurigakan.
      Test yang menemukan baris ditolak atau diperbaiki **gagal**.
- [ ] Kedua pencacah tangga dipindahkan apa adanya.
- [ ] Kedua penanda usul terisi `'1'` atau `'0'` sesudah migrasi — **nol baris berkolom usul
      kosong**. Test yang menemukan kolom usul kosong sesudah migrasi **gagal**.
- [ ] Nilai uang lama yang **melewati batas digit** menyebabkan **kegagalan yang terlihat**, bukan
      pembulatan diam-diam.

#### Butir `[terbuka]` yang menyentuh tiket ini

- ⭐ **DIPERSEMPIT 2026-09-19 — kini hanya untuk KASUS LAMA HASIL MIGRASI.** Kasus komite lama yang
  penunjuk baris penyesuaiannya **sudah salah alamat atau menunjuk baris yang tidak ada**: apa yang
  dilakukan aplikasi saat kasus itu dibuka — **menolak terang-terangan**, atau **diam seperti
  Pega**. Terdaftar sebagai butir **14** di register `spec.md`.
  > ⛔ **Kalimat lamanya dikutip, tidak dihapus:** *"**Apa yang terjadi bila baris penyesuaian
  > induk tidak ketemu** *(saat berjalan, bukan saat migrasi)* — tidak ada penanganannya di
  > korpus."*
  >
  > **Kenapa menyempit:** `[keputusan work owner]` 2026-09-19 membekukan baris penyesuaian sejak ia
  > diserahkan ke komite *(bab 11)*, sehingga **kasus baru tidak bisa lagi kehilangan baris
  > induknya**. ⚠️ **Kunci beku itu tidak berlaku surut**, jadi pertanyaannya **tidak hilang** —
  > ia hanya tinggal berlaku untuk **data lama hasil migrasi**. ⛔ **Butir ini TIDAK ditutup.**

#### Seam & verifikasi

Migrasi diverifikasi terhadap **skema uji**, bukan produksi.

## Komite Claim Prop - uru - Urutan tiket — Komite Claim Prop

⛔ **Ini bukan tiket.** Berkas ini hanya menjelaskan **tiket mana harus selesai sebelum tiket mana**.
Isinya tidak menambah keputusan apa pun.

Sumber: `spec.md` dan `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`.

---

#### Empat belas tiket

| # | Judul | Blocked by |
| --- | --- | --- |
| **00** | Skema penyimpanan komite — tiga tabel + dua kolom usul — **PREFACTOR** | skema klaim Claim Prop |
| **01** | Terima penyerahan dari Claim Prop — kasus komite lahir | 00 · kontrak muatan Claim Prop |
| **02** | Kotak kerja penyetuju dan rute giliran | 01 |
| **03** | Layar komite — 93 kolom, tiga wajah | 01 |
| **04** | Keputusan penyetuju tersimpan | 02 · 03 |
| **05** | Persetujuan bersyarat dan dua penanda usul | 04 |
| **06** | Tangga maju ke penyetuju berikutnya, atau selesai | 04 |
| **07** | Penolakan menghentikan tangga dan menolak sisa penyetuju | 06 |
| **08** | Nomor akseptasi — terbit sekali, di tingkat akhir | 06 |
| **09** | Angka uang — total penyesuaian per mata uang | 03 |
| **10** | Efek keluar ke basis data — akseptasi, riwayat, log | 08 |
| **11** | Efek keluar — dokumen PDF akseptasi | 08 |
| **12** | Efek keluar — Kasir, arasapas, dan email | 08 · 10 |
| **13** | Migrasi data lama | 00 |

---

#### Urutan pengerjaan

```
                       00  PREFACTOR
                        |
              +---------+---------+
              |                   |
             01                  13  migrasi
              |
        +-----+-----+
        |           |
       02          03
        |           |
        +-----+-----+-----------+
              |                 |
             04                09  angka uang
              |
        +-----+-----+
        |           |
       05          06
                    |
              +-----+-----+
              |           |
             07          08
                          |
                    +-----+-----+
                    |     |     |
                   10    11    12
                          (12 juga menunggu 10)
```

**Frontier awal:** hanya **00**. Sesudah 00 selesai, **01** dan **13** terbuka bersamaan.

**Jalur terpanjang:** `00 → 01 → 02/03 → 04 → 06 → 08 → 10 → 12` — delapan tingkat.

**Yang bisa dikerjakan paralel:**

- **02** dan **03** sesudah 01
- **05**, **06**, dan **09** sesudah prasyaratnya masing-masing
- **07** dan **08** sesudah 06
- **10** dan **11** sesudah 08
- **13** berjalan sendiri sepanjang waktu, hanya menunggu 00

---

#### ⚠️ Tiket yang menunggu modul Claim Prop

`[keputusan work owner]` 2026-09-18 — **seam memakai ulang milik Claim Prop, tidak menambah seam
baru.** Akibatnya menyentuh **seluruh empat belas tiket**: tidak satu pun bisa **diverifikasi
ujung-ke-ujung** sampai seam Claim Prop berdiri.

Di luar itu, **tiga tiket punya ketergantungan isi** pada modul Claim Prop:

| Tiket | Menunggu apa dari Claim Prop |
| --- | --- |
| **00** | **tabel klaim** — baris komite tinggal di tabel lintas-lini yang sama |
| **01** | **kontrak muatan penyerahan** — bentuk data yang diserahkan ke komite. **Bukan** penyelesaian modul Claim Prop; muatannya dapat dipalsukan di seam supaya kedua modul jalan paralel |
| **09** | **rumus angka klaim** — empat dari lima penghitung berkelas kasus klaim, jadi **pemiliknya Claim Prop**. Modul ini hanya menjumlahkan penyesuaian per mata uang |

⚠️ Tiket **12** menyentuh Claim Prop lewat syarat nomor akseptasi, tetapi **rule pemeriksanya sudah
ada di modul ini** — yang di luar hanyalah SQL-nya, dan itu **tidak ditelusuri**
`[keputusan work owner]`.

---

#### Pengujian — berlaku untuk seluruh tiket

`[keputusan work owner]` 2026-09-18:

| Hal | Ketetapan |
| --- | --- |
| **Seam** | **memakai ulang seam Claim Prop**, tidak menambah |
| **Efek keluar** | diuji dengan **layanan sungguhan**, bukan pengganti tiruan, di **lingkungan uji terpisah** |
| **Urutan efek keluar** | pengujiannya **MENUNGGU** urutan barunya ditetapkan |

⚠️ `[terbuka]` **Baris mana di tabel alamat layanan yang menunjuk lingkungan uji** belum ditetapkan.
Menyentuh tiket **11** dan **12**.

---

#### ⚠️ Dua titik yang SENGAJA DIUBAH — muncul di tiket mana saja

| Titik | Muncul di tiket |
| --- | --- |
| **Urutan efek keluar terhadap penyimpanan** | **10** · **11** · **12** |
| **Ketelitian angka** | **00** · **03** *(tampilan)* · **09** *(hitungan)* |

⛔ Keduanya **bukan peniruan**. Jangan diperlakukan sebagai cacat pemindahan bila hasilnya berbeda
dari Pega.

# Komite Claim Non Prop

Jumlah tiket: **13**

## Komite Claim Non Prop - 01 - Roster dan tabel seleksi berdiri sebagai data yang dikelola

---
status: aktif
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemilik proses menyemai dan menyunting roster jenjang dan tabel seleksi
tanpa menyentuh kode; konfigurasi seleksi yang meninggalkan satu pun kombinasi masukan tanpa
aturan **ditolak saat disimpan**, bukan ditemukan oleh klaim pertama yang jatuh ke dalamnya.

`ROSTER_JENJANG` dan `ATURAN_SELEKSI` beserta sequence dan constraint-nya; operasi
`KelolaRoster` dan `KelolaTabelSeleksi`; jalur baca dan tulis keduanya; layar pengelolaan.

**Persyaratan:** `S-050`, `S-051`, `S-053`, `S-054`, `S-055`, `S-056`, `S-057`, `S-036`,
`S-069`.

**Tidak termasuk:** `S-052` — pergantian pemegang tercatat sebagai peristiwa. Ia menunggu
objek `PERISTIWA_ROSTER` ditambahkan ke model data; `PERISTIWA_SIRKULASI` tidak dapat
menampungnya karena pengenal sirkulasinya wajib. **Nama** peran administratif juga tidak
ditetapkan di sini — gerbangnya dibangun, namanya menunggu pemilik proses.

**Jalur gagal:** baris yang menduplikasi kombinasi masukan yang sudah ada → `422`, tabel yang
berlaku tidak berubah · penyimpanan yang meninggalkan kombinasi tak tercakup → `422` dengan
galat yang **menyebut kombinasinya** · delegasi menunjuk orang yang tidak aktif → `422`, baris
roster tidak berubah · dua baris roster aktif berderajat sama dalam satu kelas → ditolak
constraint · pemanggil tanpa peran administratif → `403`.

**Uji:** P5-09, P5-10, P5-11 (kelas dari kombinasi masukan) · P5-16 (hak diperiksa) · BARU
untuk penolakan cakupan tidak lengkap dan untuk kosong-berarti-tidak-membatasi.

**Menggantikan:** `FilterEmailKomiteWithLimit` — penyaring roster; `.LIMIT_TOP` **dibuang**,
ia diambil tetapi tidak pernah menyaring (`F-14`) · `CreateChildKomiteCNP_Act`·12, 13 — dua
tetapan kelas di dalam langkah · `CreateChildKomiteCNP_Act`·26.8.3 — substitusi satu orang ke
orang lain, **dibuang** (`F-19`) · `SetKomiteList_Act` — pemetaan nama orang ke jabatan
(`F-18`) · `EMAILKOMITE` sebagai sumber runtime.

**Blocked by:**

- None (can start immediately)


**Dasar:** DECIDED(`K5-5`, `K6-1`, `K6-3`, `D-5`, `H-1`, `H-3`, `H-7`, keputusan beku no. 5,
ADR-D-CNP-0006, ADR-D-CNP-0016). EVIDENCED: `F-14`, `F-16`, `F-18`, `F-19`.

- [ ] Aturan seleksi disunting sebagai data; perubahan berlaku tanpa rilis.
- [ ] Penyimpanan yang meninggalkan satu pun kombinasi masukan tanpa aturan **ditolak**, dan
      galatnya menyebut kombinasi yang tidak tercakup.
- [ ] Keempat kolom batas dan penanda bersyarat pada `ATURAN_SELEKSI` boleh kosong, dan
      **kosong berarti "tidak membatasi", bukan nol** — dinyatakan eksplisit di komentar
      kolom **dan** di kaki berkas DDL. Tanpa pernyataan itu baris semai akan "dirapikan"
      jadi nol oleh pembaca berikutnya, dan seluruh isinya berubah arti (`ADR-D-CNP-0019`, `K5-6`).
- [ ] Isi awal disemai dari nilai yang **berpenulis** di sistem lama, meninggalkan batas
      nilainya kosong; nol angka ambang di dalam kode, dan angka yang hanya ada pada
      deskripsi langkah **tidak** dipakai (`K6-1`).
- [ ] Baris roster yang dinonaktifkan **tetap tersimpan**; mencabut seseorang tidak menghapus
      jejaknya.
- [ ] Delegasi tetap adalah atribut roster, bukan penukaran nama di dalam aturan.
- [ ] Nol cabang berdasarkan identitas orang di jalur mana pun — diuji dengan pencarian atas
      kode, bukan dengan memanggil operasi.
- [ ] Nama tabel, kolom, constraint, dan sequence di bawah 30 byte, berbentuk
      `<peran>_<tabel>[_n]` yang tidak mengeja kolom; pemeriksa penamaan keluar dengan kode
      gagal bila satu pun melewati batas.

**Ketidakpastian:** **Nama peran administratif** belum ditetapkan — keputusan pemilik proses.
Gerbangnya dibangun sekarang; bila jawabannya lain, yang berubah adalah **nama peran yang
dicocokkan**, bukan adanya gerbang. Menunggui, tidak menahan.

## Komite Claim Non Prop - 02 - Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Analis mengirim satu usulan pembayaran ke komite; sistem memilih kelas
kewenangan dari tabel seleksi, menyusun jenjang dari roster urut derajat menaik, dan
sirkulasi itu terbaca kembali lengkap dengan jenjangnya. Kombinasi masukan yang tidak
tercakup aturan **ditolak sebagai galat konfigurasi**, dan sirkulasi tanpa jenjang **ditolak
saat dibuat** — bukan lahir lalu diam.

`SIRKULASI` dan `JENJANG_SIRKULASI` beserta sequence dan constraint-nya; operasi
`BentukSirkulasi` dan `BacaSirkulasi`; jalur pembentukan dan jalur baca satu sirkulasi serta
seluruh sirkulasi sebuah klaim.

Ini peluru penjejak inti: ia memotong skema, perintah, kueri, dan uji sekaligus. Ia tidak
dipecah karena tiap belahannya menjadi irisan mendatar yang tidak dapat didemokan sendiri.

**Persyaratan:** `S-001`, `S-002`, `S-003`, `S-004`, `S-005`, `S-006`, `S-007`, `S-013`,
`S-014`, `S-016`, `S-044`, `S-047`, `S-064`, `S-065`, `S-067`, `S-070`.

**Tidak termasuk:** jalur menutup dan menolak klaim — tiket `07`. Jalur gagal pembentukan dan
maksud pengajuan — tiket `03`. Memutus — tiket `04`.

**Jalur gagal:** jenis sirkulasi kosong atau di luar tiga nilai → `422` menyebut medan jenis ·
kombinasi masukan tak tercakup aturan seleksi → `422` menyebut ketiga nilainya, nol sirkulasi
dan nol roster tersusun · roster terpilih kosong → `422` menyebut kelas kewenangan yang tidak
punya pemegang aktif · pemanggil tidak berhak mengajukan atas klaim itu → `403`, dapat
dibedakan dari `503` · kunci idempotensi sama dengan muatan berbeda → `409` · sirkulasi atas
usulan itu masih berjalan → `409`.

**Uji:** P5-01, P5-03 (derajat dan giliran pertama) · P5-09, P5-11 (kombinasi tak tercakup) ·
P5-10 (kelas dari nilai dan bagian) · P5-13 (nilai usulan yang diajukan, **wajib memakai klaim
dengan dua usulan beda kelas** — pada klaim berusulan tunggal deviasi ini lolos tanpa
terlihat) · P5-06 (roster idempoten) · P5-16 (hak) · BARU untuk idempotensi dan nomor urut.

**Menggantikan:** `KomiteTreaty_Flow` — mesin keadaan sirkulasi ·
`CreateChildKomiteCNP_Act`·9, 26.13–26.15 — jumlah jenjang dari cacah baris laporan, menjadi
akibat isi roster (`F-13`) · `CreateChildKomiteCNP_Act`·10 — bagian treaty ke pembanding ·
`CreateChildKomiteCNP_Act`·11 — nilai baris terakhir sebagai penentu (`F-17`) ·
`CreateChildKomiteCNP_Act`·14 — cabang yang syaratnya tidak pernah dapat benar, **dibuang**
(`F-16`) · `CreateChildKomiteCNP_Act`·28 — pembatalan hanya pada alokasi layer kosong
(`F-22`) · `KomitePostAdjustment`·16 — nomor urut sirkulasi atas satu usulan (`F-11`) ·
`KomiteCount` dan `KomiteLoop` sebagai kolom, **dibuang** — keduanya turunan (`D-2`, `K5-1`).

**Blocked by:**

- `01` — Roster dan tabel seleksi berdiri sebagai data yang dikelola


**Dasar:** DECIDED(`K5-1`, `K5-4`, `K5-5`, `K6-4`, `H-6`, `D-2`, `E-2`, keputusan beku no. 8,
ADR-D-KCNP-0030, ADR-D-KCNP-0031, ADR-D-KCNP-0032). EVIDENCED: `F-11`, `F-13`, `F-14`, `F-16`, `F-17`, `F-22`,
`F-24`.

- [ ] Jenis sirkulasi datang dari **pemanggil**; nol pembacaan kolom analisis.
- [ ] Kelas kewenangan ditentukan satu pencarian ke tabel seleksi dengan tiga masukan; nol
      cabang kelas di dalam kode.
- [ ] Masukan nilai adalah nilai **usulan yang sedang diajukan**, bukan baris terakhir daftar
      dan bukan jumlah seluruh usulan.
- [ ] Derajat disalin dari roster saat pembentukan, urut menaik, dan **tidak dihitung ulang**
      sesudahnya.
- [ ] Dua jenjang berderajat sama dalam satu sirkulasi **tidak dapat ditulis** — ditolak
      constraint, bukan diperiksa kode. Ini yang mewujudkan invarian `I-2`.
- [ ] Membentuk ulang roster sebuah sirkulasi menghasilkan roster **sama panjang**, bukan
      bertambah.
- [ ] Sirkulasi tanpa jenjang **ditolak saat dibuat**; nol berkas lahir.
- [ ] Nol kolom menyimpan jumlah jenjang — dibuktikan pencarian atas katalog skema yang
      pulang kosong. Ini yang mewujudkan invarian `I-8`.
- [ ] Dua sirkulasi bernomor urut sama atas satu usulan ditolak constraint.
- [ ] Permintaan berulang berkunci sama mengembalikan sirkulasi yang sama, bukan yang kedua.
- [ ] Tali klaim–sirkulasi dibaca **dari sisi sirkulasi**; nol pembacaan yang bersandar pada
      pengenal tersimpan di sisi klaim.
- [ ] Pembaca yang berhak membuka klaim dapat membaca sirkulasinya, tanpa peran tambahan.

**Ketidakpastian:** Tidak ada.

## Komite Claim Non Prop - 03 - Pembentukan yang gagal terlihat, dan klaim tidak membawa jejaknya

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Analis yang pembentukannya gagal **diberi tahu alasannya** dan klaimnya
tidak berubah — tidak menerima pengenal sirkulasi, tidak berpindah keadaan, tidak menerima
akibat apa pun. Maksud pengajuan tercatat pada kolom tersendiri sejak pengajuan dan tetap
tercatat meski sirkulasinya gagal lahir.

Kolom maksud pengajuan pada tabel klaim; batas transaksi pembentukan; jalur galat yang
terlihat di permukaan.

**Persyaratan:** `S-009`, `S-010`, `S-011`, `S-012`, `S-015`.

**Tidak termasuk:** kolom maksud ini menyentuh **tabel milik sisi Klaim**. Penambahannya
perlu persetujuan pemilik tabelnya; itu di luar papan ini.

**Jalur gagal:** gangguan di tengah transaksi → seluruh transaksi batal, `503` yang dapat
dibedakan dari penolakan, klaim tanpa pengenal sirkulasi dan tanpa niat pemberitahuan
tercatat · validasi tidak terpenuhi → `422` menyebut medan yang gagal, nol berkas, nol
roster, nol tulisan ke klaim · maksud tidak dapat ditulis → pembentukan batal seluruhnya,
klaim tidak berubah sama sekali.

**Uji:** P5-07 (pembentukan gagal karena pesan validasi) · P5-08 (alokasi penyebaran kosong) ·
BARU untuk ketiadaan jalur keluar senyap.

**Menggantikan:** `CreateChildKomiteCNP_Act`·29 — penjaga yang hanya melewati satu langkah
sementara langkah sesudahnya menulis pengenal milik sirkulasi lain lalu menyimpan ·
`CreateChildKomiteCloseNP_Act`·6 — penanda maksud ditulis ke induk pada langkah yang membuat
berkas anak, sebelum komite memutus apa pun · `ProteksiSendKomiteCNP_Act` — gerbang yang
menandai galat dan tidak menghentikan (`F-21`) · `IsFlagError` — penanda yang tidak dibaca
satu rule pun, **dibuang**.

**Blocked by:**

- `02` — Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya


**Dasar:** DECIDED(`K5-2`, `K5-3`, `K5-7`, `H-5`, ADR-D-KCNP-0031). EVIDENCED: `F-21`, `F-22`;
GRILL-05 N-07, N-08.

- [ ] Tiap jalur pembentukan yang berakhir tanpa sirkulasi berakhir dengan galat yang
      **terlihat pengaju** dan tercatat; pemanggil tidak pernah menerima jawaban berhasil
      tanpa sirkulasi.
- [ ] Klaim yang pembentukannya gagal **tidak berubah selain maksudnya** — diperiksa dengan
      membandingkan seluruh barisnya sebelum dan sesudah.
- [ ] Nol klaim membawa pengenal sirkulasi tanpa sirkulasi yang bersesuaian. Ini yang
      mewujudkan invarian `I-11`.
- [ ] Maksud dan akibat menempati **kolom yang berbeda**; maksud tidak pernah menjadi akibat.
- [ ] Validasi **menolak** kasus-guna; nol jalur yang memasang pesan lalu melanjutkan.

**Ketidakpastian:** Penambahan kolom maksud pada tabel klaim menunggu persetujuan pemilik
tabel sisi Klaim. Menunggui, tidak menahan: bila ditolak, yang berubah adalah **tempat maksud
disimpan**, bukan bahwa maksud dan akibat terpisah.

## Komite Claim Non Prop - 04 - Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang aktif memutus dengan komentarnya; siapa pun yang bukan
pemegang jenjang aktif **ditolak di server** dengan pencocokan identitas penuh. Penolakan
menutup sirkulasi dan menyisakan jenjang di bawahnya berkeadaan `TIDAK_SAMPAI` — bukan
"menolak" atas nama orang yang tidak pernah bertindak.

Operasi `CatatKeputusan`; jalur pencatatan keputusan; keadaan jenjang dan keadaan sirkulasi.

**Persyaratan:** `S-017`, `S-018`, `S-019`, `S-020`, `S-021`, `S-022`, `S-023`, `S-026`,
`S-027`, `S-028`, `S-029`, `S-030`, `S-068`.

**Tidak termasuk:** akibat pada klaim dan tulis-balik — tiket `05`. Versi usulan dan penanda
usulan — tiket `06`. Nomor akseptasi — tiket `08`.

**Jalur gagal:** pemanggil bukan pemegang jenjang aktif → `403` menyatakan giliran bukan
miliknya, **berbeda bentuk** dari `503` · keputusan kedua atas jenjang yang sudah memutus →
`409`, catatan yang ada tidak berubah · keputusan atas sirkulasi yang sudah selesai → `409` ·
nilai penanda di luar dua keadaan → `422` · kunci idempotensi sama dengan muatan berbeda →
`409` · gangguan di tengah → `503`, nol keputusan tersimpan.

**Uji:** P5-01, P5-02 (urutan giliran, **paritas** — dirancang lulus) · P5-02b (dua penentu
berselisih, **dirancang gagal**) · P5-14, P5-15 (peran dan jabatan dari roster) · P5-05
(penanda dua-keadaan) · BARU untuk `TIDAK_SAMPAI`, ketetapan keputusan, dan idempotensi.

**Menggantikan:** `KomiteRouter`·1–4 — penugasan lewat hitungan jenjang, **dibuang** (`K5-1`) ·
`KomiteRouter`·6, 6.1 — baris pertama yang belum memutuskan, **dipertahankan sebagai
paritas** · `KomiteRouter` — empat cabang mati, **dibuang** · `IsKomiteLoop` — gerbang loop,
**dibuang** · `DateApproval` dan `DateApprove` — dua kolom untuk satu fakta, dilebur jadi satu.

**Blocked by:**

- `02` — Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya


**Dasar:** DECIDED(`K5-1`, `K5-6`, `E-2`, `E-1`, keputusan beku no. 1, 2, 3, 4, ADR-D-KCNP-0030).
EVIDENCED: `F-24` — nol `pyPrivilegeName` berisi pada 59 rule; GRILL-06 P6-3 — pemeriksaan
wewenang memasang pesan dan tidak menghentikan apa pun.

- [ ] Identitas dicocokkan dengan **kesetaraan penuh**; pengguna yang identitasnya mengandung
      pemegang sebagai potongan teks **ditolak**.
- [ ] Jenjang aktif adalah jenjang berderajat terendah yang belum memutuskan, dan **hitungan
      jenjang tidak pernah** menjadi syarat penugasan maupun penutupan.
- [ ] Sirkulasi berjalan punya **tepat satu** jenjang aktif — dihitung per sirkulasi. Ini yang
      mewujudkan invarian `I-7`.
- [ ] Penolakan menutup sirkulasi, dan seluruh jenjang berderajat lebih tinggi berkeadaan
      `TIDAK_SAMPAI` **tanpa komentar dan tanpa waktu milik orang lain**. Invarian `I-5`.
- [ ] Persetujuan jenjang terakhir menutup sirkulasi; penutupan lahir dari **keadaan
      keputusan**, bukan dari perbandingan dua bilangan. Invarian `I-6`.
- [ ] Keputusan tidak dapat disunting sesudah tercatat. Invarian `I-9`.
- [ ] Tabel jenjang memuat **tepat satu** kolom waktu keputusan — pencarian katalog skema
      menemukan tidak ada yang kedua. Invarian `I-10`.
- [ ] Peran dan jabatan disalin pada saat keputusan diambil; mutasi jabatan sesudahnya tidak
      mengubah riwayat.
- [ ] Galat gangguan dan penolakan wewenang **berbeda bentuk** di permukaan.
- [ ] Komentar, penanda bersyarat, dan penanda usulan yang sudah diisi **tidak hilang** ketika
      penyimpanan gagal.

**Ketidakpastian:** Tidak ada.

## Komite Claim Non Prop - 05 - Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Nasib klaim dihitung dari tiga masukan — jenis sirkulasi, bersyarat, dan
hasil — oleh satu fungsi yang ditulis sekali, lalu ditulis ke klaim di dalam transaksi
keputusan yang sama. Pembaca melihat "komite menyetujui usulan penutupan" sebagai **klaim
ditutup**, tidak pernah sebagai "klaim disetujui".

Fungsi akibat berikut ketujuh barisnya; tulis-balik ke sisi Klaim; model baca yang menyajikan
nilai tampil dari fungsi yang sama.

**Persyaratan:** `S-031`, `S-032`, `S-033`, `S-034`, `S-035`, `S-046`, `S-049`.

**Tidak termasuk:** nomor akseptasi — tiket `08`. Efek ke luar — tiket `11`.

**Jalur gagal:** kombinasi yang tidak tercakup ketujuh baris → `422` sebagai galat
konfigurasi, keputusan tidak tersimpan · salah satu tulis-balik gagal → seluruh transaksi
batal, `503`.

**Uji:** BARU — **ketujuh baris** fungsi akibat, ketiga nilai tampil diperiksa **dari model
baca**, bukan dari kode · BARU untuk keputusan bersyarat mendarat pada jenjang yang memutus ·
BARU — uji DDL: nol kolom menyimpan "sedang disirkulasikan".

**Menggantikan:** `KomitePostAdjustment` — akibat jalur A, tersebar ·
`KomitePostAdjustmentCWP` — akibat jalur B · `KomitePostAdjustment`·14.12 — penanda usul tutup
mendarat di induk pada persetujuan terakhir non-bersyarat (`F-5`) ·
`KomitePostAdjustment`·20 — penolakan satu usulan menolak seluruh klaim; **dipertahankan**
sebagai baris ketiga tabel, ditandai `D-1` paritas-dengan-pertanyaan-bisnis, bukan cacat ·
`InsertChronology_DT` — kronologi, **tanpa** pengecualian per peran · `KOMITECOUNT` dan
`FLAGONGOINGCOMMITTE` sebagai kolom, **dibuang**.

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`D-3`, `J-3`, `D-1`, `D-2`, `K-07`, keputusan beku no. 9, ADR-D-CNP-0016,
ADR-D-KCNP-0032). EVIDENCED: `F-5`; GRILL-06 P6-4 — jalur bersyarat menulis ke halaman yang
berbeda, bukan hanya ke baris yang berbeda.

- [ ] Fungsi akibat ditulis **sekali**, di satu tempat; ketujuh barisnya diuji.
- [ ] Hasil tidak pernah berarti nasib klaim; keadaan sirkulasi dan akibat pada klaim
      menempati **kolom yang berbeda**.
- [ ] Model baca menyajikan ketiga nilai tampil dari fungsi yang sama; **nol** pembaca
      menafsirkan hasil sendiri.
- [ ] Seluruh tulis-balik terjadi di dalam transaksi keputusan — nol pesan antar-konteks, nol
      tautan basis data.
- [ ] Keputusan bersyarat mendarat pada **jenjang yang memutus**, bukan pada baris terakhir
      daftar dan bukan pada halaman yang berbeda. Invarian `I-3`.
- [ ] Penanda "sedang disirkulasikan" **tidak disimpan**; ia dihitung dari jenjang saat
      dibaca.
- [ ] Pengaju dan maksud tersaji **terpisah** dari akibat.

**Ketidakpastian:** Tidak ada.

## Komite Claim Non Prop - 06 - Versi usulan naik saat penanda berubah, dan penjaga menolak versi yang berselisih

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang berderajat terendah menyunting penanda usulan sambil
memutus, dan suntingan itu terekam sebagai **dua peristiwa dalam satu transaksi** — usulan
diubah, membawa nilai sebelum dan sesudah, dan jenjang memutus. Jenjang berikutnya yakin
bahwa yang ia putuskan sama persis dengan yang dilihat jenjang sebelumnya.

`VERSI_USULAN` dan `PERISTIWA_SIRKULASI` beserta constraint-nya; penjaga versi di batas tulis.

**Persyaratan:** `S-024`, `S-025`.

**Tidak termasuk:** `S-052` — peristiwa pergantian pemegang pada roster. Ia menunggu objek
`PERISTIWA_ROSTER`; `PERISTIWA_SIRKULASI` tidak dapat menampungnya.

**Jalur gagal:** penyuntingan datang dari jenjang bukan yang pertama → `422`, penanda usulan
**dan** versi **dan** keputusan sama-sama tidak tersimpan · versi usulan berubah oleh jenjang
selain yang pertama → `409` menyebut versi yang diharapkan dan versi yang ditemukan.

**Uji:** BARU — penyuntingan oleh jenjang bukan pertama ditolak · BARU — dua peristiwa lahir
dalam satu transaksi, dan keduanya hilang bersama ketika transaksi batal.

**Menggantikan:** suntingan yang menimpa nilai sebelumnya **tanpa jejak** — tidak ada padanan
di sistem lama; bentuk yang dijamin 296 kendali layar saja, kini ditegakkan juga sebagai
penjaga murah di batas tulis (`J-1`).

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`E-1`, `J-1`, `J-2`, `E-3`). EVIDENCED: `F-2`, `F-8`;
`3-to-tickets/FAKTA-LAYAR-01.md` — keempat medan berkunci jenjang pertama bernama.

- [ ] Versi usulan naik **hanya** ketika penanda usulan berubah.
- [ ] Tiap keputusan mencatat versi yang diputusnya.
- [ ] Penyuntingan terekam dua peristiwa dalam **satu** transaksi; keduanya batal bersama.
- [ ] Jenjang pertama yang **menolak** tetap meninggalkan suntingannya tercatat sebagai
      peristiwa, tanpa akibat hilir.
- [ ] Nilai sebelum dan sesudah menyimpan **penanda**, bukan uang — nol kolom uang di kedua
      objek.

**Ketidakpastian:** Tidak ada.

## Komite Claim Non Prop - 07 - Sirkulasi menutup dan menolak klaim, pemegangnya diselesaikan dari peran

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Analis mengirim usulan menutup atau menolak klaim; sirkulasinya tetap satu
jenjang, dan pemegangnya diselesaikan **dari roster berdasarkan peran** — bukan dari nama
yang ditulis di dalam aturan. Bila tidak ada baris roster aktif yang memegang peran itu,
sirkulasi **ditolak saat dibuat**.

Jenis sirkulasi menutup dan menolak klaim pada `BentukSirkulasi`; baris fungsi akibat untuk
keduanya.

**Persyaratan:** `S-008`.

**Tidak termasuk:** —

**Jalur gagal:** nol baris roster aktif memegang peran itu → `422` menyebut **peran yang
kosong**, nol sirkulasi lahir, dan maksud tidak ditulis ke klaim.

**Uji:** P5-04 (jalur tutup klaim) · P5-05 (jalur tolak klaim) · **PS-01** (jenis dinyatakan
pemanggil) · **PS-02** (pemegang dari peran) · **PS-03** (penolakan saat peran kosong).
Ketiga uji `PS-xx` **dirancang gagal** terhadap sistem lama — deviasi 23, 24, 25 pada
`REGISTER-DEVIASI.md` bagian 6, seluruhnya **belum diratifikasi**.

**Menggantikan:** `CreateChildKomiteCloseNP_Act`·8, 9 — satu jenjang dengan pemegang dari nama
orang, surel, inisial, dan **dua sebutan jabatan berbeda untuk satu kedudukan** (`F-20`) ·
`CreateChildKomiteCloseNP_Act`·12 — jenis dibedakan dari nilai kolom analisis.

**Blocked by:**

- `05` — Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama


**Dasar:** DECIDED(`K6-4`, keputusan beku no. 5). EVIDENCED: `F-20`.
`K6-4` menggantikan `H-2`, yang tetap tercatat `BUNYI HILANG`.

- [ ] Jenis sirkulasi datang dari pemanggil; nol pembacaan kolom analisis untuk menentukannya.
- [ ] Pemegang diselesaikan dari roster **berdasarkan peran**; nol nama orang di dalam aturan.
- [ ] Peran yang kosong **menolak pembentukan**, dan galatnya menyebut peran itu.
- [ ] Ketiga uji `PS-xx` **gagal** terhadap sistem lama. Uji `PS-xx` yang **lulus** paritas
      adalah tanda deviasinya tidak terpasang, dan itu kegagalan tiket ini.
- [ ] Deret uji tetap dieja `PS-xx`; ia **tidak** dinormalkan menjadi `P-xx`, karena
      penamaannya sengaja tidak terbaca sebagai ronde tujuh.

**Ketidakpastian:** Tidak ada.

## Komite Claim Non Prop - 08 - Nomor akseptasi terbit sekali, di dalam transaksi keputusan

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Persetujuan jenjang terakhir atas usulan yang **tidak bersyarat**
menerbitkan nomor akseptasi, di dalam transaksi keputusan yang sama. Persetujuan **bersyarat**
menutup sirkulasi tanpa menerbitkan nomor. Gagal menerbitkan membatalkan seluruh keputusan —
tidak pernah ada keputusan tersimpan tanpa nomor ketika syaratnya terpenuhi.

Pembungkus pencatat di atas sumber urutan Oracle yang ada; constraint keunikan pada kolom
nomor akseptasi yang sudah ada di sisi Klaim.

**Persyaratan:** `S-038`, `S-039`, `S-040`, `S-041`, `S-042`, `S-043`, `S-066`.

**Tidak termasuk:** **nilai ambang periode buku** — **PAGAR-01**. Mekanika penomoran ditulis
penuh; nilainya adalah parameter bernama tanpa isi, dan pembukanya adalah isi badan
`PROC_GENERATE_SEQUENCE_NUMBER` pada basis data berjalan (`INVENTARIS-BUKTI.md` §2.5 baris 4).
Isi baris akseptasi dan bentuk simpannya — **PAGAR-02**, tiket `11`.

**Jalur gagal:** sumber urutan tidak menjawab → `503`, seluruh transaksi batal, nol keputusan
tersimpan · penerbitan kedua atas usulan yang sama → ditolak constraint, pemanggil menerima
hasil penerbitan pertama · permintaan menyunting nomor → `409` · parameter periode buku tidak
terisi → `503`; **nol nomor terbit dengan periode terkaan**.

**Uji:** BARU — syarat terbit, idempotensi dengan mengulang perintah ber-kunci sama,
pembatalan transaksi saat penerbitan gagal · BARU — uji DDL untuk keunikan nomor.

**Menggantikan:** `GetSequenceNumber_SQL` dan `PROC_GENERATE_SEQUENCE_NUMBER` — prosedur lama
**dipakai apa adanya** pada fase 1, dibungkus pencatat (`D-4`) · `GetKodeProdNonLife_SQL` —
kode produksi, **berpagar** PAGAR-01 · idempotensi lewat baca-lalu-tulis — diganti constraint.

**Blocked by:**

- `05` — Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama


**Dasar:** DECIDED(`D-4`, `J-3`, `E-2`, keputusan beku no. 7). Dasar faktualnya bertanda
**TAFSIR (menunggu C-01)** — berkas DDL yang dipegang adalah hasil reverse-engineer, bukan
ekspor langsung.

- [ ] Nomor terbit **bila dan hanya bila** jenjang terakhir menyetujui dan usulan tidak
      bersyarat.
- [ ] Persetujuan bersyarat menutup sirkulasi, **tidak** menerbitkan nomor, dan nilai
      tampilnya berbunyi disetujui bersyarat.
- [ ] Idempotensi ditegakkan **constraint**, bukan baca-lalu-tulis.
- [ ] Percobaan ulang ber-kunci sama mengembalikan **nomor yang sama**, bukan nomor kedua.
- [ ] Nomor tidak pernah disunting sesudah terbit.
- [ ] Nol keputusan final tidak bersyarat tersimpan tanpa nomor.
- [ ] Periode buku ditentukan di **satu lapis**, Oracle; sisi aplikasi tidak menggeser ulang
      dan tidak menambal. Ambangnya parameter bernama, bukan angka di dalam kode.

**Ketidakpastian:** **C-01** — isi badan prosedur penerbit pada basis data berjalan belum
pernah dibandingkan. Menunggui, tidak menahan: mekanika sudah ditetapkan `D-4`; bila isinya
berbeda, yang berubah adalah **pembungkusnya**, bukan bahwa nomor terbit di dalam transaksi.

## Komite Claim Non Prop - 09 - Layar persetujuan komite dengan tiga masukan dan kunci yang ditegakkan dua kali

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang membuka layar persetujuan, membaca seluruh nilai sebagai
proyeksi baca-saja, dan mengisi **tiga hal saja** — keputusan, komentar, dan penanda usulan.
Empat penanda berkunci jenjang pertama; medan induk yang baca-saja **tidak memblokir**
keputusan.

Layar persetujuan komite; penegakan kunci sisi klien.

**Persyaratan:** `S-058`, `S-059`, `S-060`, `S-061`, `S-062`, `S-063`.

**Tidak termasuk:** kebenaran tiap aturan kunci — ia diuji **sekali** di seam perintah, bukan
dua kali. Seam layar sempit dengan sengaja: ia memeriksa bahwa penegakan pertama **ada**.

**Jalur gagal:** keputusan tanpa komentar → `422`, keputusan tidak tersimpan · penyuntingan
penanda dari jenjang bukan pertama → `422` **di server**, meski klien mengizinkannya · medan
di luar ketiga masukan terkirim → diabaikan, dan ditolak `422` bila berbeda dari nilai
tersimpan.

**Uji:** BARU — tiap kunci diperiksa **ada** di klien, lalu diperiksa **ditegakkan** di
server dengan permintaan yang melewati klien. Kunci yang hanya ada di klien **bukan kunci**,
dan uji itu yang membuktikannya.

**Menggantikan:** `ShowTransfer` — layar persetujuan · `ViewTransferDtl` — pembuka layar ·
`ResetSubjectivityNote` · medan rekening, 16 kemunculan, **seluruhnya baca-saja** (`F-1`,
`F-6`) · 16 properti nilai, **seluruhnya baca-saja** (`F-9`) · medan bersufiks dua,
**dibuang** — nol pembaca di Activity modul ini (`F-3`) · blok bersyarat yang syaratnya tidak
pernah dapat benar dan blok bersyarat "tidak pernah", **dibuang** · dua `pyDisabledWhen` yang
menunjuk kendali **yang sudah baca-saja**, **dibuang** — lihat `FAKTA-LAYAR-01.md`.

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`E-4`, `J-2`, `J-1`, ADR-D-CNP-0006). EVIDENCED: `F-1`, `F-3`, `F-6`, `F-7`,
`F-8`, `F-9`; `3-to-tickets/FAKTA-LAYAR-01.md` — enam `pyDisabledWhen`, keenamnya terbaca
dari `pyValue` selnya sendiri, bukan dari kedekatan posisi.

- [ ] Satu-satunya masukan komite adalah keputusan, komentar, dan penanda usulan.
- [ ] Seluruh panel nilai baca-saja; komite **tidak dapat** mengubah satu pun angka uang dan
      **tidak pernah** menyunting rekening penerima.
- [ ] Keempat penanda berkunci jenjang pertama adalah `.IsSubjectivity`, `.SubjectivityNote`,
      `.Adjustment.IsProposeClose`, `.Adjustment.IsPropReserved` — bernama, bukan dicacah.
- [ ] Label mata uang baca-saja bagi **seluruh** jenjang.
- [ ] Medan induk baca-saja **tidak memblokir** keputusan; komentar wajib.
- [ ] Tiap blok yang dibuang **tercatat beserta syarat lamanya**, bukan dihilangkan diam-diam.
- [ ] Tiap kunci ditegakkan di klien **dan** di server; nol kunci yang hanya ada di klien.

**Ketidakpastian:** **Bentuk kendali catatan syarat.** `.SubjectivityNote` di sistem lama
adalah `pxDropdown`, bukan kotak teks, sementara `SPEC-KOMITE-01.md` menuliskannya "teks
panjang" dan kolomnya `VARCHAR2(2000 CHAR)`. **Isi daftar pilihannya tidak terbaca**, dan itu
lubang bukti yang sah: `Rule-Obj-FieldValue` memang tidak ikut diekspor —
`INVENTARIS-BUKTI.md` §1, baris "XML rule Komite", kolom "Apa yang TIDAK diliput", entri
**Field Value**. Baris itu sudah ada; tidak ada baris baru ditambahkan.

## Komite Claim Non Prop - 10 - Giliran saya hari ini, dan umur yang tidak pernah disimpan

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang membuka daftar sirkulasi yang menunggu keputusannya, dan
melihat sudah berapa lama tiap sirkulasi berjalan serta berapa lama jenjang aktif menunggu —
tanpa ada yang perlu menuliskannya, dan tanpa satu pun kolom umur tersimpan.

Jalur baca daftar menunggu pemanggil; umur sebagai atribut turunan.

**Persyaratan:** `S-045`, `S-048`.

**Tidak termasuk:** masa tenggat dan eskalasi. Tidak ada keadaan kedaluwarsa dan tidak ada
pekerja latar (`E-5`); umur dikumpulkan sebagai data agar aturan tenggat kelak diputuskan
dengan bukti, bukan dengan angka karangan.

**Jalur gagal:** pemanggil bukan pemegang jenjang mana pun → **daftar kosong, bukan galat** ·
pemanggil tidak berhak membaca klaimnya → `403`, nol sebagian isi bocor lewat pesan galat.

**Uji:** P5-03 — empat nama keranjang sistem lama tidak terpakai; daftar dihitung dari jenjang
aktif · BARU — uji DDL: nol kolom umur tersimpan.

**Menggantikan:** keranjang kerja Pega — daftar dihitung dari jenjang aktif, bukan dari
keranjang yang ditulis terpisah · `KomiteTreaty_Flow` — nol SLA, timer, eskalasi, reassign,
dan delegasi di seluruh modul (`F-4`); **tidak ada yang menggantikannya**, umur menjadi
turunan.

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`E-5`, `K5-1`, `K6-2`, `J-4`). EVIDENCED: `F-4`.

- [ ] Daftar dihitung dari **jenjang aktif**, bukan dari keranjang tersimpan.
- [ ] Umur sirkulasi dan umur tunggu jenjang aktif dihitung **saat dibaca**.
- [ ] Nol kolom umur tersimpan — dibuktikan pencarian atas katalog skema.
- [ ] Nol keadaan kedaluwarsa dan nol pekerja latar.

**Ketidakpastian:** Tidak ada.

## Komite Claim Non Prop - 11 - Lima port hilir berdiri kosong dengan niat tercatat lebih dulu

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Tidak satu pun jalur berakhir di efek hilir tanpa port yang **dinamai**.
Niat memanggil tiap port dicatat **sebelum** panggilan, dengan kunci idempotensi per efek,
dan kegagalan port tidak membatalkan keputusan yang sudah tercatat — ia meninggalkan niat yang
belum terpenuhi dan **terbaca**.

Lima port: pencatatan akseptasi, pembalikan akseptasi, pengiriman kasir, penerbitan surat,
pengunggahan dokumen. Kelimanya **kosong**.

**Persyaratan:** `S-037`.

**Tidak termasuk:** **isi kelima port**, seluruhnya berpagar. PAGAR-02 dan PAGAR-03 —
`PEGA_JSON_OS_AKSEP_KLAIM`, `PEGA_JSON_OS_AKSEP_SUBJECTIVITY`, `HISTORYAKSEPTASIPEGA`
(`INVENTARIS-BUKTI.md` §2.1) dan nol baris `OS_AKSEPTASI_KLAIM.DATA_JSON` (§2.5 baris 2).
PAGAR-04 — `XOL2_AKSEP_KLAIM` (§2.1) dan `SetProtectionEstimation` (§2.3). PAGAR-05 dan
PAGAR-06 — nol baris `POOLDATA.DIRECTTOKASIR_LOG` (§2.5 baris 1). PAGAR-07 —
`PostEmailKomiteCNP` (§2.3). PAGAR-08 — `SendEmailWithAttachments` (§2.3).

**Jalur gagal:** port dipanggil tanpa niat tercatat lebih dulu → uji gagal · port gagal →
keputusan yang sudah tercatat **tidak** dibatalkan; niat tetap ada dan terbaca sebagai belum
terpenuhi.

**Uji:** BARU — port diuji sebagai **port kosong**: uji memeriksa bahwa port dipanggil dengan
niat yang tercatat lebih dulu dan kunci idempotensi, **bukan** memeriksa apa yang dikirimnya.

**Menggantikan:** `KomitePostAdjustment`·19.1 → `PostEmailKomiteCNP` · `InsertOSKlaimCNP`,
`InsertOSSubjectivityCNP` · `InsertXOLKlaimCNP` · `SaveRejectOSKomiteCNP` ·
`HitServiceToKasirKMT_Act`·10.3, 10.7, 10.9 · `GenerateAccCNP_act` → `InsertDocument_Act` ·
`SendEmailKlaim_KMT`, `SendEmailKlaimRejectClose_KMT` · `InsertHistoryAkseptasiPega_Sql`.
Seluruhnya **DIPAGARI**, bukan digantikan — yang dibangun hanyalah pintunya.

**Blocked by:**

- `05` — Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama
- `08` — Nomor akseptasi terbit sekali, di dalam transaksi keputusan


**Dasar:** DECIDED(keputusan beku no. 6, `ADR-D-KCNP-0031`, `ADR-D-KCNP-0032`). Kelima port berpagar
PAGAR-02 sampai PAGAR-08 pada `SPEC-KOMITE-01.md` bagian *Out of Scope*.

- [ ] Kelima port **dinamai** dan **kosong**; nol jalur berakhir di efek hilir tanpa melewati
      salah satunya.
- [ ] Niat dicatat **sebelum** panggilan, dengan kunci idempotensi per efek.
- [ ] Efek ke luar berada **di luar** batas transaksi dan dijalankan sesudahnya.
- [ ] Kegagalan port tidak membatalkan keputusan; niat yang belum terpenuhi terbaca.
- [ ] Nol persyaratan, nol muatan, dan nol penjaga ditulis untuk isi port mana pun —
      **termasuk pada alasannya**, bukan hanya pada kesimpulannya.

**Ketidakpastian:** **`KonversiKlaim_Act`** dan **`InsertJsonClaimTreatyNonProp_act`** adalah
efek ke luar yang dipanggil dari jalur keputusan sistem lama dan **tidak terpeta ke satu pun
dari kelima port**. Berkas keempat rule-nya **ada di repo** (`INVENTARIS-BUKTI.md` §2.4),
sehingga keduanya **lubang rancangan, bukan lubang bukti**, dan tidak boleh dipagari.
Keduanya berkeadaan `BELUM DIPUTUSKAN`; keputusannya milik pemilik proses.

## Komite Claim Non Prop - 12 - Pergantian pemegang jenjang tercatat sebagai peristiwa

---
status: tertahan
---



*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemilik proses mengganti pemegang sebuah baris roster, atau mencatat
delegasi tetap, dan pergantian itu tersimpan sebagai peristiwa dengan **pemegang sebelum dan
sesudah, pelaku, alasan, dan waktunya**. Auditor dapat menelusuri siapa memegang jenjang apa
pada tanggal berapa tanpa bertanya kepada siapa pun — dan tanpa bergantung pada baris roster
yang sudah tertimpa.

`PERISTIWA_ROSTER` beserta sequence dan constraint-nya; penulisan peristiwa di dalam
transaksi `KelolaRoster` yang sama dengan perubahan barisnya.

**Persyaratan:** `S-052`.

**Tidak termasuk:** **pengalihan sesaat.** `H-4` menetapkan ia tetap tidak diekspos: modelnya
menampung lewat `E-3`, permukaannya tidak ada. Kolom objek ini **tidak kurang** — yang tidak
ada adalah jalur yang menulisnya. Jangan menambahkan kolom yang dikira terlupakan.

**Jalur gagal:** peristiwa tidak dapat ditulis → **penyimpanan roster batal seluruhnya**;
baris roster tidak berubah, dan tidak ada peristiwa separuh yang tertinggal · baris yang
kedua kolom pemegangnya kosong sekaligus → ditolak constraint, karena baris semacam itu tidak
mencatat pergantian apa pun · pemanggil tanpa peran administratif → `403`, nol peristiwa
lahir.

**Uji:** BARU — pergantian pemegang menghasilkan tepat satu peristiwa yang membawa nilai
sebelum **dan** sesudah · BARU — kegagalan penulisan peristiwa membatalkan perubahan
rosternya, diperiksa dengan membandingkan baris roster sebelum dan sesudah · BARU — uji DDL
untuk `CK_PERISTIWA_ROSTER_2` dan untuk ketiadaan kolom pengalihan sesaat.

**Menggantikan:** `CreateChildKomiteCNP_Act`·26.8.3 — substitusi satu orang ke orang lain yang
ditulis **di dalam kode** (`F-19`). Cabangnya dibuang oleh keputusan beku no. 5; kebutuhan
yang ditambalnya digantikan dua hal — delegasi tetap sebagai atribut roster (`S-051`, tiket
`01`) dan pergantiannya tercatat sebagai peristiwa di sini (`S-052`). Di sistem lama nilai
`KomiteID` sebelum ditimpa **hilang begitu ditimpa**; itu yang objek ini kembalikan.

**Blocked by:**

- `01` — Roster dan tabel seleksi berdiri sebagai data yang dikelola


**Dasar:** DECIDED(`E-3`, `H-3`, `H-4`, keputusan beku no. 4, `ADR-D-CNP-0019`, `K5-6`).
EVIDENCED: `F-19` — satu substitusi orang ditulis di dalam rule,
`CreateChildKomiteCNP_Act`·26.8.3.

- [ ] Tiap pergantian pemegang menghasilkan **tepat satu** peristiwa dengan pelaku, alasan,
      dan waktunya.
- [ ] Peristiwa membawa **pemegang sebelum dan sesudah**; nilai sebelum tidak hilang ketika
      baris roster ditimpa.
- [ ] Peristiwa dan perubahan baris roster berada dalam **satu transaksi**; keduanya batal
      bersama.
- [ ] Baris yang kedua kolom pemegangnya kosong sekaligus **tidak dapat ditulis** — ditolak
      constraint, bukan diperiksa kode.
- [ ] Kosong pada kolom pemegang berarti **tidak ada pemegang**, bukan nol; nol kolom di objek
      ini yang kosong dan nolnya dapat tertukar (`K5-6`, `ADR-D-CNP-0019`).
- [ ] Peristiwa roster **tidak ikut terbaca** saat sirkulasi dibaca — ia hidup di seam
      `KelolaRoster`, bukan `BacaSirkulasi`.
- [ ] Nol kolom untuk pengalihan sesaat, dan nol jalur yang menulisnya (`H-4`).
- [ ] Cap waktu berzona `Asia/Jakarta` (aturan kerja `G`); nama di bawah 30 byte, constraint
      berbentuk `<peran>_<tabel>[_n]` yang tidak mengeja kolom.

**Ketidakpastian:** Tidak ada.

## Komite Claim Non Prop - REA - Papan tiket — Komite Claim Non Prop

**Tiket: aktif 1 · tertahan 11 · menunggu instance 0 · selesai 0 · mati 0**

> Pencacah ini menghitung **tiket**. Persyaratan `S-xxx` punya pencacahnya sendiri di `SPEC-KOMITE-01.md` bagian *Cakupan dan cacah*; keduanya diberi label bendanya supaya tidak pernah tertukar.

#### Yang sebenarnya dapat dimulai

**1 dari 12 tiket di papan.** Sebelas sisanya tertahan **tiket lain**, bukan sesuatu di luar papan — rantai dependensinya nyata dan seluruhnya ada di sini. Itu keadaan yang berbeda dari papan lapisan data sisi Klaim, tempat seluruh sisa tertahan keputusan di luar papan.

Tidak ada penahan dari luar papan. Tiga tiket **menunggui** jawaban dari luar tanpa tertahan olehnya; daftarnya di bawah.

Diterbitkan 21 September 2026 dari `SPEC-KOMITE-01.md` bagian *Persyaratan* — 70 persyaratan `S-xxx`, **seluruhnya** diketiketkan sejak `PERISTIWA_ROSTER` ditambahkan sebagai objek ketujuh.
Tracker belum terkonfigurasi (`/setup-matt-pocock-skills` belum dijalankan), jadi papannya berkas lokal.

#### Aturan papan

- **Status adalah field di dalam berkas tiket** (`status: aktif | tertahan | selesai | selesai-sebagian`), bukan lokasi foldernya. Pembangkit membaca field, dan **tidak pernah menghapus berkas tiket**.
- **Penahan dari luar papan ditulis dengan nama aslinya** — `C-01`, `PAGAR-01` — tidak diterjemahkan jadi nomor tiket. Yang menahan dari luar harus terlihat berasal dari luar.
- **Penahan yang sudah selesai dicoret, tidak dihapus**, supaya rantainya tetap terbaca: `~~05~~ *(selesai)*`.
- **Nomor yang lompat bukan kekeliruan.** Penomoran tidak disusun ulang, karena menyusunnya ulang memutus setiap rujukan yang sudah ada.
- **Uji hidup di dalam tiketnya**, bukan sebagai tiket tersendiri. Ini berbeda dari papan lapisan data sisi Klaim, dan sebabnya disebut di bagian terakhir.

#### Papan

| # | Tiket | Ditahan oleh |
|---|---|---|
| `01` | Roster dan tabel seleksi berdiri sebagai data yang dikelola | — |
| `02` | Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya | `01` |
| `03` | Pembentukan yang gagal terlihat, dan klaim tidak membawa jejaknya | `02` |
| `04` | Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak | `02` |
| `05` | Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama | `04` |
| `06` | Versi usulan naik saat penanda berubah, dan penjaga menolak versi yang berselisih | `04` |
| `07` | Sirkulasi menutup dan menolak klaim, pemegangnya diselesaikan dari peran | `05` |
| `08` | Nomor akseptasi terbit sekali, di dalam transaksi keputusan | `05` |
| `09` | Layar persetujuan komite dengan tiga masukan dan kunci yang ditegakkan dua kali | `04` |
| `10` | Giliran saya hari ini, dan umur yang tidak pernah disimpan | `04` |
| `11` | Lima port hilir berdiri kosong dengan niat tercatat lebih dulu | `05` · `08` |
| `12` | Pergantian pemegang jenjang tercatat sebagai peristiwa | `01` |

#### Menunggui, tidak menahan

Tiket ini **boleh jalan**. Jawaban yang ditunggu mengubah satu hal yang tiketnya sudah menyebut — sebuah nama peran, sebuah tempat simpan, sebuah pembungkus — bukan rancangannya. Dipisahkan dari penahan sungguhan supaya papan tidak berteriak serigala.

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|
| `01` | Nama peran administratif — keputusan pemilik proses | **Nama peran yang dicocokkan**, bukan adanya gerbang |
| `03` | Persetujuan pemilik tabel sisi Klaim atas kolom maksud pengajuan | **Tempat maksud disimpan**, bukan bahwa maksud dan akibat terpisah |
| `08` | `C-01` — isi badan prosedur penerbit pada basis data berjalan | **Pembungkusnya**, bukan bahwa nomor terbit di dalam transaksi |

#### Dapat dimulai hari pertama

| # | Tiket |
|---|---|
| `01` | Roster dan tabel seleksi berdiri sebagai data yang dikelola |

Satu tiket, dan itu disengaja: tanpa roster dan tabel seleksi sebagai data, tidak ada sirkulasi yang dapat dibentuk sama sekali. Tiket `02` adalah peluru penjejak intinya dan terbuka begitu `01` hijau, bersama tiket `12`.

---

#### Tidak diketiketkan

Dua daftar, dan pembedaannya penting. Yang pertama memuat **persyaratan** — baris `S-xxx` di
spesifikasi yang belum punya tiket. Yang kedua memuat **keputusan yang belum diambil**, yang
bukan persyaratan dan tidak akan pernah ditemukan dengan mencari `S-xxx` di spesifikasi.
Dicampur, pembaca berikutnya akan mengira ada persyaratan hilang dan mencarinya di tempat
yang tidak memuatnya.

##### Persyaratan ditahan

**Kosong.** Ketujuh puluh persyaratan `S-xxx` seluruhnya punya tiket.

`S-052` pernah berada di sini, karena ia menuntut pergantian pemegang tercatat sebagai
peristiwa sementara satu-satunya tabel peristiwa mewajibkan pengenal sirkulasi — dan
pergantian roster tidak punya sirkulasi. Objek ketujuh `PERISTIWA_ROSTER` ditambahkan
21 September 2026, `S-052` lolos `T-1`, dan ia kini menjadi tiket `12`. Judul daftar ini
tidak dihapus: kosong adalah keadaan yang perlu terbaca, bukan bagian yang perlu hilang.

##### Keputusan belum diambil

Bukan persyaratan, dan **bukan lubang bukti**. Ketiganya menunggu orang, bukan berkas.

| Yang ditunggu | Keadaannya |
|---|---|
| **Nama peran administratif** pengelola roster dan tabel seleksi | Keputusan pemilik proses. `H-4` menyebut "peran administratif tersendiri di luar keanggotaan komite" sebagai **default**, bukan sebagai ketetapan. Gerbangnya sudah dispesifikasikan dan dibangun di tiket `01`; yang menunggu hanya namanya, dan itu **menunggui, tidak menahan** |
| **Bentuk migrasi sirkulasi berjalan** | Satu cacah yang belum diambil: berapa klaim di produksi membawa pengenal sirkulasi yang bukan miliknya — `INVENTARIS-BUKTI.md` §2.5 baris 6. Yang **sudah** diputuskan ada di tiket `02` sebagai `S-070` |
| **Port hilir bagi `KonversiKlaim_Act` dan `InsertJsonClaimTreatyNonProp_act`** | Keputusan pemilik proses. Berkas keempat rule-nya **ada di repo** (`INVENTARIS-BUKTI.md` §2.4), sehingga keduanya **lubang rancangan, bukan lubang bukti**, dan berkeadaan `BELUM DIPUTUSKAN` — bukan `DIPAGARI` |

#### Delapan pagar, nol tiket

Aliran A-4, A-5, dan **isi** A-6 tidak punya tiket sama sekali, dan itu benar. Yang dibangun hanyalah **pintunya** — tiket `11`, lima port kosong. PAGAR-01 menahan satu nilai di dalam tiket `08`, bukan tiketnya.

Pagar berlaku bagi **alasan**, bukan hanya bagi kesimpulan: tiket mana pun yang argumennya bersandar pada isi aliran beku tidak sah, meski kesimpulannya berada di aliran terbuka.

#### Dua penyimpangan dari papan lapisan data sisi Klaim

Konvensi berkasnya ditiru apa adanya. Dua hal berbeda, dan sebabnya dicatat supaya tidak terbaca sebagai kelalaian.

**Uji tidak menjadi tiket tersendiri.** Papan sisi Klaim memberi delapan uji tiketnya masing-masing, bertanda `menunggu-instance`, karena batch-nya lapisan data saja — ujinya memang tidak dapat hijau sampai DDL berjalan di instance nyata. Batch ini memotong seluruh lapisan sekaligus, sehingga ujinya hijau di dalam tiketnya sendiri. Memisahkannya justru melanggar aturan irisan tegak.

**Tiap tiket membawa bagian `Persyaratan`, `Jalur gagal`, `Uji`, dan `Menggantikan`.** Keempatnya tidak punya tempat di konvensi sisi Klaim, dan ditambahkan sebagai bagian baru — bukan dibuang. Ketiganya yang pertama datang dari lapisan `S-xxx`; yang terakhir datang dari bagian *Ketertelusuran*, dan tanpanya shadow-run kehilangan dasar pembandingnya.
