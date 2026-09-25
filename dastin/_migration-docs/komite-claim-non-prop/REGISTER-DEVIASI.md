> Modul  : Komite Claim Non Prop · Tahap spesifikasi · 2026-09-21
> Peran  : auditor
> Ronde  : penutupan spesifikasi, 2026-09-21 — bagian 6 ditambahkan: deviasi 23, 24, 25
> Masukan: `KETETAPAN.md` bagian 9, 11, dan 12 · `GRILL-05/04-PARITAS.md` · `GRILL-05/01-TEMUAN.md` · `GRILL-06/01-PEMBACAAN.md` · penyuntingan penutup ronde 6 · `K6-4` · `SPEC-KOMITE-01.md` (2026-09-21)
> Status : TERBUKA
> Sifat  : HIDUP

# REGISTER DEVIASI DIHARAPKAN

Deviasi adalah perbedaan perilaku yang **disengaja** terhadap sistem lama. Tiap deviasi
membawa uji paritas yang **dirancang gagal**: uji itu ada bukan untuk dilewati, melainkan
untuk membuktikan perbedaannya memang yang dimaksud. Uji yang lulus pada baris di bawah ini
adalah tanda deviasinya tidak terpasang.

Dasarnya D-1: di mana grilling menghasilkan cacat ber-`EVIDENCED`, perbaikan menjadi
default dan paritas yang harus dibela.

---

## 1. Deviasi ronde 1 – 4

Tiga belas baris, tercatat utuh di **`KETETAPAN.md` bagian 9**, dan tidak disalin ke sini
agar tidak ada dua salinan yang perlahan berbeda. Rujuk dengan nomor barisnya:
`KETETAPAN.md` §9 baris 1 … 13.

**Peringatan ratifikasi.** `KETETAPAN.md` bagian 10.1 mencatat bahwa hanya uji baris 8 yang
pernah dirancang eksplisit. Dua belas uji lainnya disusun juru catat dari bunyi deviasinya,
bukan disalin dari sumber, dan berstatus **belum diratifikasi**. Status itu tidak berubah
ronde ini.

---

## 2. Deviasi ronde 5

| # | Deviasi | Diperintahkan oleh | Uji paritas yang dirancang gagal | Skenario |
|---|---|---|---|---|
| 14 | Satu penentu jenjang aktif; hitungan jenjang menjadi turunan dan tidak lagi menjadi syarat penugasan maupun penutupan | `K5-1` | Sirkulasi tiga jenjang yang jenjang keduanya memutus lebih dulu; sistem lama menugaskan lewat satu penentu sementara hitungannya menunjuk jenjang lain, sistem baru tidak dapat menirukan selisih itu | P5-02b · *lihat §5.3* |
| 15 | Pembuatan sirkulasi dan akibatnya adalah satu transaksi | `K5-2`, memperluas `H-5` | Pembentukan yang memunculkan pesan validasi; sistem lama melewati pembuatan lalu tetap menulis nomor komite dari sirkulasi sebelumnya ke klaim dan menyimpannya, sistem baru tidak menulis apa pun | P5-07 |
| 16 | Kegagalan pembentukan wajib terlihat | `K5-3` | Adjustment tanpa `SpreadingRisk`; sistem lama keluar diam-diam dan menyimpan klaim dengan penanda yang tak dibaca siapa pun, sistem baru menampilkan kegagalan dan tidak menyimpan | P5-08 |
| 17 | Pembentukan roster bersifat idempoten | `K5-4` | Klaim bersyarat yang dibentuk dua kali; sistem lama menghasilkan roster ganda dengan jumlah jenjang lebih kecil dari jumlah anggota, sistem baru menghasilkan roster yang sama | P5-06 |
| 18 | Tidak ada kombinasi nilai × bagian tanpa aturan seleksi | `K5-5`, melengkapi `H-1` | `RNMShare > 30` dengan nilai berapa pun; sistem lama memanggil penyaring dengan parameter tak terisi dan melahirkan sirkulasi tanpa jenjang, sistem baru memakai kelas yang ditetapkan | P5-09, P5-11 · *lihat §5.2* |
| 19 | Penanda dua-keadaan adalah satu boolean dengan satu ejaan | `K5-6` | Penolakan lewat jalur tutup/tolak; sistem lama menulis `IsCloseFile` sebagai `"0"` pada anak dan kosong pada induk, sistem baru menulis satu nilai yang sama di kedua tempat | P5-05 |
| 20 | Empat pemetaan orang→jabatan yang ditulis di dalam rule pindah ke data | `D-5` diperluas, sidang F-19 | Anggota komite bernama `CHRISTOPMARHASAK`; sistem lama menimpa jabatannya dengan teks di dalam rule, sistem baru memakai jabatan roster | P5-15 |
| 21 | Modul memeriksa hak sebelum menampilkan dan sebelum memutuskan | `ADR-0006`, sidang F-24 | Pengguna mana pun yang dapat membuka klaimnya; sistem lama tidak memeriksa hak apa pun — nol `pyPrivilegeName` berisi pada 59 rule — sistem baru memeriksa | P5-16 |

---

## 3. Yang **bukan** deviasi, meski sempat dicatat begitu

| Hal | Mengapa bukan deviasi |
|---|---|
| Urutan giliran jenjang | Bacaan pertama ronde 5 menyimpulkan sistem lama memberi giliran kepada baris terakhir yang belum memutuskan, sehingga mempertahankan urutan naik akan menjadi deviasi. Bacaan itu keliru (`GRILL-05/06-PUTUSAN.md` M5-01): sistem lama sudah memberi giliran kepada baris **pertama**. Urutan naik adalah **paritas**, dan P5-01 serta P5-02 adalah uji paritas biasa yang dirancang **lulus** |
| Pengalihan sesaat tidak diekspos (`H-4`) | Pilihan rancangan, bukan perbedaan perilaku terhadap fitur yang ada — sistem lama pun tidak mengeksposnya. Disarankan pindah ke sini dari `REGISTER-PAGAR.md` PG-02 hanya bila kelak diekspos |

---

## 4. Aturan yang berlaku bagi register ini

1. Tiap deviasi menyebut ketetapan yang memerintahkannya. Deviasi tanpa ketetapan adalah
   selera, dan selera tidak masuk register.
2. Tiap deviasi membawa satu uji yang dirancang gagal, dan uji itu menyebut **keadaan
   masukan** yang membuatnya gagal — bukan sekadar menyebut perbedaannya.
3. Uji yang disusun oleh juru catat, bukan disalin dari sumber keputusannya, ditandai
   **belum diratifikasi** sampai pemilik keputusannya membacanya.
4. Deviasi yang ternyata paritas dipindahkan ke bagian 3 dengan alasannya, bukan dihapus.

---

## 5. Pemutakhiran ronde 6 — 2026-09-20

### 5.1 Deviasi 11 dilebur ke deviasi 18

Baris 11 pada `KETETAPAN.md` §9 — "cabang mati `Flagkomite` langkah 14 ditutup" — **tidak
lagi dicatat sebagai deviasi tersendiri**. `K5-5` menelan perintahnya dengan bingkai yang
lebih baik: setiap kombinasi nilai × bagian harus punya aturan, dan kombinasi yang tidak
tercakup adalah **galat konfigurasi**, bukan roster kosong. Baris lama tidak dihapus dari
`KETETAPAN.md` §9; ia diberi status di sini.

| # | Status baru | Oleh |
|---|---|---|
| 11 | **DILEBUR ke deviasi 18** — jangan dicatat dua kali | `K5-5`, `K6-1` |

Uji yang menyertai baris 11 ("`RNMShare` di antara 15 dan 30…") ikut gugur, dan itu tepat:
`K6-1` melarang angka 15 dipakai, sehingga uji yang memakainya tidak dapat dijalankan.

### 5.2 Deviasi 18 diperluas

| # | Deviasi | Diperintahkan oleh | Uji paritas yang dirancang gagal | Skenario |
|---|---|---|---|---|
| 18 | Tidak ada kombinasi **nilai × bagian × bersyarat** tanpa aturan seleksi | `K5-5`, `K6-1`, `K6-3` | `RNMShare > 30` dengan nilai berapa pun, dan nilai > 50 juta dengan `RNMShare ≤ 30`; sistem lama memanggil penyaring dengan parameter tak terisi dan melahirkan sirkulasi tanpa jenjang, sistem baru memakai kelas yang ditetapkan. Kombinasi yang tidak tercakup pada sistem baru adalah galat konfigurasi, bukan roster kosong | P5-09, P5-11 |

Masukan ketiga — bersyarat — masuk ke tabel yang sama dengan isi awal yang **mereproduksi
perilaku sekarang** (`K6-3`), sehingga tidak melahirkan deviasi baru.

### 5.3 Skenario baru pada deviasi 14

| # | Skenario tambahan | Uji paritas yang dirancang gagal | Ratifikasi |
|---|---|---|---|
| 14 | Sirkulasi **bersyarat** berjenjang lebih dari satu | **Keadaan masukan:** usulan ber-`IsSubjectivity=true` yang rosternya lebih dari satu jenjang. **Sistem lama:** langkah 6 `KomitePostAdjustment` dilewati, sehingga `KomiteAproval` pada roster sirkulasi tak pernah terisi; `KomiteRouter` langkah 6.1 menyeleksi `==0` dan menahan giliran pada jenjang pertama tanpa batas — sementara sub-langkah 26.6 memberi klaim bersyarat `LIMIT_BOTTOM = 0`, yaitu roster terluas. **Sistem baru:** giliran berpindah. | **belum diratifikasi** |

Bukan temuan baru: `K5-1` dan `K5-6` sudah menyelesaikannya; yang kurang hanya ujinya.
Asal: `GRILL-06/01-PEMBACAAN.md` P6-4 dan `GRILL-05/01-TEMUAN.md` N-10.

### 5.4 Deviasi baru

| # | Deviasi | Diperintahkan oleh | Uji paritas yang dirancang gagal | Skenario |
|---|---|---|---|---|
| 22 | Maksud dan akibat adalah dua penanda berbeda; "diajukan untuk ditutup" tidak pernah menjadi "ditutup" | `K5-7`, `ADR-0031` | **Keadaan masukan:** klaim diajukan untuk ditutup lewat `TypeComentAnalysis != "5"`, dan pembentukan sirkulasi gagal — daftar penyebaran risiko kosong, atau halaman anak membawa pesan validasi. **Sistem lama:** klaim tersimpan membawa `IsCloseFile` terisi meski tidak ada sirkulasi yang pernah lahir, dan `CNPStatusCase` sudah berbunyi "COMITEE ACCEPTANCE (DEPT. HEAD)". **Sistem baru:** maksud tercatat, akibat tidak, dan keadaan klaim tidak berubah | P5-07, P5-08 |

### 5.5 Status ratifikasi

Status **belum diratifikasi** pada dua belas uji lama **tidak berubah** ronde ini. Uji pada
§5.3 dan §5.4 lahir ronde 5 dan 6 dan juga **belum diratifikasi**.

---

## 6. Pemutakhiran tahap spesifikasi — 2026-09-21

### 6.1 Tiga deviasi turunan `K6-4`

`K6-4` ditetapkan pada penutup ronde 6, **sesudah** pemutakhiran terakhir register ini. Tiga
perubahan perilaku yang diperintahkannya karena itu tercatat di `SPEC-KOMITE-01.md` sebagai
`BERUBAH` tanpa punya baris di sini. Deviasi bertanda `BERUBAH` tanpa baris register adalah
kebalikan dari deviasi yatim — sama-sama putus. Ketiganya diberi nomor sekarang.

Deret uji **`PS-xx`** dibuka di sini: `P` paritas, `S` tahap spesifikasi. Ia **bukan** ronde
tujuh — grilling tetap tertutup, dan `KETETAPAN.md` bagian 7 mengatur kapan ronde dibuka.
Deret tersendiri dipakai agar tidak mengulang tabrakan `G-01…G-20` melawan `G-1…G-4`.

| # | Deviasi | Diperintahkan oleh | Uji paritas yang dirancang gagal | Skenario |
|---|---|---|---|---|
| 23 | Jenis sirkulasi **dinyatakan pemanggil**, tidak diturunkan dari nilai kolom analisis | `K6-4` | **Keadaan masukan:** klaim yang kolom analisisnya bernilai selain `"5"`, sementara pengaju memilih tindakan **menolak** klaim — dua sumber yang berselisih. **Sistem lama:** `CreateChildKomiteCloseNP_Act` langkah 12 membedakan jenis dari nilai kolom analisis, sehingga jenis yang lahir mengikuti kolom dan bukan tindakan, tanpa ada yang tahu keduanya berselisih. **Sistem baru:** jenis datang dari pemanggil dan kolom analisis tidak dibaca, sehingga selisih itu tidak dapat ditirukan | PS-01 |
| 24 | Pemegang jalur tutup dan tolak **diselesaikan dari roster berdasarkan peran**, bukan dari nama di dalam aturan | `K6-4` | **Keadaan masukan:** sirkulasi menutup klaim dibentuk ketika orang yang namanya ditulis di dalam aturan sudah tidak menduduki kedudukan itu. **Sistem lama:** langkah 8 dan 9 menulis nama orang, surel, inisial, dan **dua sebutan jabatan berbeda untuk satu kedudukan** (`F-20`); sirkulasi lahir membawa pemegang yang sudah tidak menjabat. **Sistem baru:** pemegang diselesaikan dari roster berdasarkan peran, sehingga yang lahir membawa pemegang yang berlaku | PS-02 |
| 25 | **Penolakan saat dibuat** bila tidak ada baris roster aktif yang memegang peran itu | `K6-4` | **Keadaan masukan:** sirkulasi menutup klaim dibentuk ketika nol baris roster aktif memegang peran itu. **Sistem lama:** pemegang adalah tetapan di dalam aturan, sehingga keadaan "tidak ada pemegang" tidak pernah dapat terjadi pada jalur ini; pada jalur pembayaran, keadaan setara — penyaring roster pulang kosong — tetap melahirkan berkas yang lalu diam (`F-22`). **Sistem baru:** pembentukan ditolak dengan galat yang menyebut peran yang kosong | PS-03 |

Ketiganya **belum diratifikasi** — disusun juru catat dari bunyi `K6-4`, bukan disalin dari
sumber keputusannya.

### 6.2 Akibatnya pada cacah

Dua puluh lima baris deviasi. Baris 11 tetap **gugur** karena dilebur ke 18. Dua puluh tiga
dari dua puluh lima berstatus **belum diratifikasi**; hanya baris 8 yang pernah dirancang
eksplisit, dan status itu tidak berubah.

### 6.3 Apa yang tidak dikerjakan di sini

Bagian 3 — "yang bukan deviasi" — tidak berubah. `H-4` tetap bukan deviasi. Tidak ada baris
lama yang disunting bunyinya, dan tidak ada baris yang dihapus.
