# Penyimpanan relasional NB Treaty In — dokumen `DATA_JSON` diganti tabel

> ⛔ **`T_POLIS_BREAKDOWN_SPREAD` DIBATALKAN 23-09-2026.** `[keputusan work owner]` Tabel itu **tidak ada**. Keempat medannya turunan: dua dari master `POOLDATA.PROPORTIONALARRG`, dua dihitung dari `TotalPremium` dan `TotalClaim` yang sudah tersimpan. Rujukan di bawah dicoret, bunyinya tidak dihapus.

---


*Spec penyimpanan. Disusun 23 September 2026, sesudah rancangan tabel selesai dan dua belas butir
penahan ditutup work owner dalam satu ronde tanya-jawab.*

---

> ## ⛔⛔ RALAT 23 September 2026 — sesudah data guide dan satu dokumen nyata
>
> Spec ini ditulis **sebelum** data guide dan satu dokumen EDM sungguhan diterima. Lima hal
> berubah. Yang berlaku adalah **`rancangan-tabel-datar-treaty-in.md` Bab 4quinque**, bukan yang
> tertulis di bawah.
>
> | # | Yang tertulis di spec ini | Yang berlaku |
> | ---: | --- | --- |
> | 1 | skala kolom uang **minimal 9 desimal** | ⭐ **`NUMBER(38,8)`** *(semula ~~`NUMBER(20,8)`~~ — dinaikkan 23-09-2026 sore; **desimalnya tetap delapan**, yang berubah hanya sisi kiri koma)* |
> | 2 | AC uang "tersimpan tanpa kehilangan satu digit pun" | utuh **sampai 8 desimal** |
> | 3 | migrasi "tidak mengubah satu pun nilai uang" | **membulatkan** yang berdesimal > 8 |
> | 4 | — | ⭐ **tabel baru ~~`T_POLIS_BREAKDOWN_SPREAD`~~ ⛔** ← `BreakDownSpreadList`, 4 medan, ada di tingkat polis dan `OldData` |
> | 5 | rumus selisih `[dugaan]` | ⭐ **`[terverifikasi]` dari data produksi** — `baru − lama.nilai`, terbukti pada 6 medan |
>
> ⭐ Dan satu medan baru yang belum ada di spec ini: **`Remark`** *(panjang 128)*.
> `Show` dan `ViewState` tetap tidak dimigrasi.
>
> ⚠️ **Data guide berkedudukan PELENGKAP, bukan sumber kebenaran** — satu dokumen tunggal memuat
> **95 jalur** yang tidak ada di dalamnya. Ia sah untuk `o:length`; ⛔ tidak sah untuk membuktikan
> ketiadaan sebuah medan.

---


## 1 · Cara membaca berkas ini

### Ringkasan

| Ukuran | Jumlah | |
| --- | ---: | --- |
| User story | **41** | |
| Acceptance criteria | **67** | ⚠️ *RALAT 23 Sep — semula ~~63~~; **tiga** ditambahkan pada pemeriksaan silang, **satu lagi** (AC 20b) pada koreksi presisi uang. Angka ~~66~~ sempat tertulis sebelum AC 20b masuk* |
| Butir `[terbuka]` aktif *(Bab 9)* | **4** | |
| Butir `[penyimpangan sadar]` | **5** | ⚠️ *RALAT — semula ~~4~~; blok tambahan AC 64–66 membawa satu lagi* |
| Tabel yang dispec — proporsional | ~~**8**~~ ⛔ ⭐ **7** *(tabel sebaran tambahan dibatalkan 23-09-2026 sore)* | |
| Tabel yang dispec — non-proporsional | **11** *(~~10~~)* | |

**Sebaran penanda** — cara B, kemunculan label di dalam jendela hitung:

| Penanda | Cara B *(berlaku)* | Cara A *(dikutip, keliru)* |
| --- | ---: | ---: |
| `[terverifikasi]` | **82** | ~~47~~ |
| `[keputusan work owner]` | **18** | ~~12~~ |
| `[penyimpangan sadar]` | **4** | ~~3~~ |
| `[terbuka]` | **3** | ~~4~~ |
| `[data DBA]` | **2** | 2 |
| `[dugaan]` | **0** | ~~2~~ |

⚠️ **Cara A keliru pada lima dari enam penanda.** Ia cacah manual saat menulis, dan meleset jauh
— terutama `[terverifikasi]`, yang saya taksir 47 padahal 82. ⛔ Angka lamanya dikutip, tidak
dihapus (CLAUDE.md §4a).

⭐ Dan satu yang perlu dibedakan: **`[terbuka]` muncul 3 kali sebagai label**, tetapi **butir
terbuka aktif ada 4** — keempatnya terdaftar di Bab 9, dan satu di antaranya tidak membawa label
harfiah di dalam kalimatnya. Keduanya benar; yang dihitung berbeda.

⛔ `[dugaan]` ternyata **nol** — tidak ada satu pun pernyataan berlabel dugaan di spec ini.
Yang semula saya kira dugaan ternyata seluruhnya terbukti dari korpus.

**Cara menghitung — dua cara, jendela hitung disebutkan.**
*Cara A* — cacah manual saat menulis, dicatat per bab.
*Cara B* — pola teks atas berkas jadi: user story = baris yang cocok `^\d+\. Sebagai `;
acceptance criteria = baris yang cocok `^\d+\. \x60\[`; penanda = cacah kemunculan tiap label
**di luar** bab ini.
Jendela: seluruh berkas **kecuali** bab 1 dan bab TELEMETRI, supaya ringkasannya tidak menghitung
dirinya sendiri.

⚠️ Bila kedua cara berselisih, **yang berlaku cara B**, dan angka cara A dikutip bukan dihapus.

### Penanda

| Penanda | Artinya |
| --- | --- |
| `[terverifikasi]` | terbukti dari korpus, disertai jalur berkas dan nama rule |
| `[keputusan work owner]` | ditetapkan work owner, tanggalnya disebut |
| `[dugaan]` | pola kuat, **belum** terbukti |
| `[terbuka]` | belum dijawab · pemiliknya disebut · ⛔ **tidak boleh ditutup oleh pelaksana** |
| `[penyimpangan sadar]` | sengaja **berbeda** dari sistem lama, bunyi lamanya dikutip |
| `[data DBA]` | hanya dapat dijawab dari basis data, bukan dari korpus |

### ⚠️ Peringatan yang berlaku untuk seluruh berkas

⛔ **Setiap cacah medan di sini adalah BATAS BAWAH, bukan total.** Sapuan korpus menelusuri jalur
berjangkar `PolicyTreatyIn.…`. Di dalam kalang, anggota daftar dirujuk **relatif** — `.Suggest`,
`.CedingCoName` — dan tidak terlihat oleh sapuan itu. Tiga kali terbukti: `SuggestList` sempat
terbaca **nol** medan · `CedingCoList` **tidak terbaca sama sekali** sampai ditunjukkan work owner ·
lima property EDM terlewat karena `DataTransform` memakai tag `pyPropertiesName`, bukan
`PropertiesName`.

⇒ Rancangan menyediakan jalan keluarnya: lihat Implementation Decision **ID-27**.

---

## 2 · Problem Statement

Data polis treaty inward hari ini disimpan sebagai **satu dokumen JSON** di kolom
`POOLDATA.json_polis.DATA_JSON`. Dokumen itu ditulis Pega lewat stored procedure, dan seluruh
isinya — uang, tanggal, kode, penanda — **bertipe teks**.

Akibatnya, empat hal yang seharusnya mudah menjadi mahal atau mustahil:

**Tidak bisa ditanyai.** Pertanyaan sesederhana *"berapa total premi treaty proporsional kuartal
ini"* memerlukan pembongkaran dokumen, bukan satu `SUM`. Laporan yang ada mengatasinya dengan
menulis salinan datar ke `TREATYINPRODUCTION` — sehingga ada **dua tempat** menyimpan hal yang sama,
dan keduanya bisa berbeda.

**Tidak bisa dijaga.** Tidak ada kunci, tidak ada keunikan, tidak ada tipe. Dua endorsemen yang
berjalan bersamaan bisa mendapat nomor generasi yang sama, dan tidak ada yang menghentikannya.
`[terverifikasi]` `EDM\RDBList\SelectProdKe` mengambil `PRODKE` dengan `MAX + 1` tanpa penguncian.

**Bentuknya tidak seragam.** `ListInstallment` bersarang pada satu bentuk dokumen dan datar pada
bentuk lain. `[terverifikasi]` 101 rujukan bentuk datar di 11 berkas, 58 rujukan bentuk bersarang
di 6 berkas. Pengurai yang menganggapnya seragam akan patah.

**Sejarahnya disimpan tiga kali.** Dokumen endorsemen memuat keadaan sekarang, salinan beku keadaan
lama (`OldData`, 94 jalur), dan selisih keduanya (`TreatyDifference`, 39 jalur) — ketiganya dengan
nama medan yang sama persis. `[terverifikasi]` Dan salinan beku itu bersarang ke dalam dirinya
sendiri pada endorsemen berlapis, sementara **nol rujukan** di korpus pernah membacanya.

Migrasi ke Go menghapus stored procedure. Tanpa penyimpanan relasional, seluruh beban pembongkaran
dokumen pindah ke Go tanpa satu pun masalah di atas terselesaikan.

---

## 3 · Solution

Ganti kolom dokumen `DATA_JSON` dengan **enam tabel baru** pada bentuk proporsional, sembilan pada
non-proporsional, ditambah satu tabel lama yang sudah datar dan dipakai apa adanya.

Tiga gagasan yang menopangnya:

**Generasi, bukan salinan.** Satu polis dan seluruh endorsemennya duduk di tabel yang sama,
dibedakan `PRODKE`. `OldData` tidak menjadi tabel — ia **baris yang ditunjuk `OLD_POLIS_ID`**.
Pembekuannya dijamin aturan: baris generasi lampau tidak boleh disunting. Selisih tidak disimpan di
NB; ia dihitung, dan tabel proyeksinya milik ronde EDM.

⇒ Sekitar **144 kolom cermin lenyap** dari rancangan, dan cacat `OldData` di dalam `OldData` tidak
pernah lahir.

**`NOURUT`, bukan posisi tersirat.** Setiap tabel anak membawa nomor urut baris. Itulah kunci
pasangan antar generasi, menggantikan indeks kalang Pega yang tidak tersimpan di mana pun.

**Tipe yang dipilih, bukan diwarisi.** Uang dan persen menjadi angka presisi tetap; tanggal menjadi
`DATE`; cacah menjadi bilangan bulat. Tetapi **kode dan penanda tetap teks** — sebab nol di depan
pada kode bisnis dan pembandingan teks pada penanda persetujuan keduanya membawa makna yang hilang
bila dikonversi.

Satu skema melayani NB dan EDM. Yang membedakan hanya **baris mana yang terisi**, bukan tabel mana
yang ada.

---

## 4 · User Stories

### Penyimpanan polis baru

1. Sebagai petugas admin treaty, saya ingin polis yang saya simpan tersimpan utuh beserta seluruh
   daftar di dalamnya dalam satu kali simpan, supaya tidak ada bagian yang tertinggal separuh jalan.
2. Sebagai petugas admin treaty, saya ingin penyimpanan yang gagal di tengah membatalkan
   seluruhnya, supaya tidak ada polis yang tersimpan setengah.
3. Sebagai petugas admin treaty, saya ingin nomor polis dan nomor generasi tidak pernah kembar,
   supaya dua orang yang bekerja bersamaan tidak saling menimpa.
4. Sebagai petugas admin treaty, saya ingin baris angsuran yang saya hapus membuat nomor urut
   sisanya dirapatkan, supaya daftarnya tetap runtut di layar.
5. Sebagai petugas admin treaty, saya ingin nilai uang yang saya masukkan tersimpan tepat seperti
   yang saya ketik, supaya tidak ada pembulatan diam-diam.
6. Sebagai petugas admin treaty, saya ingin medan yang saya kosongkan tersimpan sebagai kosong,
   bukan sebagai nol, supaya "belum diisi" dan "nol rupiah" tidak tertukar.
7. Sebagai petugas admin treaty, saya ingin kode bisnis berawalan nol tersimpan apa adanya, supaya
   penggolongan jenis bisnis tetap menemukan barisnya.
8. Sebagai petugas admin treaty, saya ingin daftar ceding tersimpan sebagai baris-baris tersendiri,
   supaya bisa ditambah dan dikurangi tanpa menyunting teks gabungan.

### Pembacaan dan penyajian

9. Sebagai petugas admin treaty, saya ingin polis yang saya buka kembali tampil persis seperti saat
   disimpan, supaya saya tidak perlu memeriksa ulang setiap medan.
10. Sebagai petugas admin treaty, saya ingin daftar angsuran tampil dalam urutan yang sama seperti
    saat saya masukkan, supaya nomor angsuran tidak berpindah-pindah.
11. Sebagai Section Head, saya ingin melihat riwayat persetujuan polis beserta siapa dan kapan,
    supaya saya tahu tahapan mana yang sudah dilewati.
12. Sebagai Dept Head, saya ingin keputusan setuju atau tolak pada tiap baris usulan tersimpan
    terpisah dari keputusan atas polisnya, supaya keduanya tidak tertukar.
13. Sebagai analis Finance, saya ingin menjumlahkan premi seluruh polis dalam satu kuartal dengan
    satu pertanyaan, supaya tidak perlu membongkar dokumen satu per satu.
14. Sebagai analis Finance, saya ingin menyaring polis menurut jenis treaty proporsional atau
    non-proporsional, supaya laporan keduanya tidak tercampur.
15. Sebagai analis Finance, saya ingin mengurutkan polis menurut tanggal mulai periode, supaya
    urutannya benar dan bukan urutan teks.

### Bentuk proporsional dan non-proporsional

16. Sebagai pengembang, saya ingin satu skema melayani kedua bentuk treaty, supaya tidak ada dua
    rangkaian tabel yang harus dijaga sejajar.
17. Sebagai pengembang, saya ingin tabel khas non-proporsional bernilai nol baris pada polis
    proporsional, supaya pembacaan tidak perlu bercabang.
18. Sebagai pengembang, saya ingin penulisan menolak baris XOL pada polis proporsional, supaya
    bentuk yang salah tidak masuk tanpa ketahuan.
19. Sebagai pengembang, saya ingin penanda layer tersimpan pada baris layer, bukan pada polis,
    supaya polis dengan beberapa layer tidak kehilangan layer kedua dan seterusnya.
20. Sebagai pengembang, saya ingin jadwal angsuran proporsional berhenti pada satu tingkat, supaya
    tidak ada tabel anak yang selamanya kosong.

### Generasi dan endorsemen

21. Sebagai pengembang, saya ingin setiap generasi polis menjadi baris tersendiri, supaya riwayat
    perubahannya utuh tanpa menyalin dokumen.
22. Sebagai pengembang, saya ingin generasi lampau tidak dapat disunting, supaya keadaan lama
    benar-benar beku.
23. Sebagai pengembang, saya ingin satu generasi hanya boleh punya satu penerus, supaya rantai
    endorsemen tidak bercabang menjadi pohon.
24. Sebagai pengembang, saya ingin generasi sebelumnya ditunjuk langsung lewat kunci tamu, supaya
    tidak perlu mencari lewat nomor polis dan urutan.
25. Sebagai pengembang, saya ingin skema ini sudah menampung endorsemen sejak awal walau NB belum
    memakainya, supaya ronde EDM tidak perlu mengubah bentuk tabel.

### Migrasi

26. Sebagai pemilik data migrasi, saya ingin dokumen lama dipindahkan ke tabel baru tanpa mengubah
    satu pun nilai uang, supaya laporan lama tetap cocok.
27. Sebagai pemilik data migrasi, saya ingin galat presisi yang sudah ada di data lama ikut terbawa
    apa adanya, supaya tidak ada angka yang berubah diam-diam.
28. Sebagai pemilik data migrasi, saya ingin tahu berapa banyak dokumen yang gagal diurai dan
    kenapa, supaya bisa ditindaklanjuti dan bukan ditelan diam-diam.
29. Sebagai pemilik data migrasi, saya ingin pemuatan lewat jalur yang sama dengan penyimpanan
    biasa, supaya data lama tunduk pada penjaga yang sama dengan data baru.
30. Sebagai pemilik data migrasi, saya ingin medan yang tidak dikenal dalam dokumen lama tersimpan,
    bukan dibuang, supaya tidak ada yang hilang sebelum ketahuan pentingnya.

### Pemeliharaan dan pemeriksaan

31. Sebagai DBA, saya ingin setiap query menyebut skema secara eksplisit, supaya tidak bergantung
    pada skema bawaan sambungan.
32. Sebagai DBA, saya ingin tipe dan presisi kolom ditetapkan sekali di satu tempat, supaya tidak
    ada dua tabel yang menyimpan hal sama dengan presisi berbeda.
33. Sebagai pengembang, saya ingin perbandingan nilai uang tidak pernah memakai kesamaan persis,
    supaya galat presisi tidak membuat dua nilai yang sama dinilai berbeda.
34. Sebagai pengembang, saya ingin identitas penyimpan diambil dari identitas login, bukan diketik,
    supaya jejaknya tidak dapat dipalsukan.
35. Sebagai pengembang, saya ingin nama tampilan dan identitas akses disimpan di kolom berbeda,
    supaya keduanya tidak tertukar.
36. Sebagai pembaca laporan lewat SQL, saya ingin tahu kolom mana yang menampung dua satuan berbeda,
    supaya tidak menjumlahkan persen dengan rupiah.
37. Sebagai pengembang berikutnya, saya ingin setiap perbedaan dari sistem lama tercatat beserta
    bunyi aslinya, supaya tidak dikira cacat lalu "diperbaiki" balik.
38. Sebagai pengembang berikutnya, saya ingin tahu medan mana yang belum pernah terukur, supaya
    tidak mengira daftar kolomnya sudah lengkap.

### Keamanan dan kepatuhan

39. Sebagai pemilik data, saya ingin nama orang dan nomor polis tidak tersalin ke berkas rancangan
    maupun test, supaya data nasabah tidak menyebar.
40. Sebagai pemilik data, saya ingin test berjalan di atas data buatan, bukan cuplikan produksi,
    supaya lingkungan uji tidak memuat data sungguhan.
41. Sebagai pemilik data, saya ingin perubahan skema memerlukan persetujuan manusia sebelum
    dijalankan, supaya tidak ada migrasi yang berjalan tanpa sepengetahuan.

---

## 5 · Implementation Decisions

### Modul dan antarmuka

**ID-1** `[keputusan work owner]` Arah ketergantungan tetap `handlers → services → repository`
*(CLAUDE.md §5)*. Penyimpanan hidup seluruhnya di `repository`; aturan dagang di `services`.

**ID-2** Satu antarmuka `repository` menjadi **satu-satunya pintu** ke penyimpanan polis treaty
inward. Bentuk minimalnya dua perintah: menyimpan satu polis beserta seluruh anaknya, dan mengambil
satu polis beserta seluruh anaknya menurut `(NOPOLIS, PRODKE)`.

**ID-3** Pemuat data migrasi **menulis lewat antarmuka yang sama**, bukan jalur tersendiri. ⇒ data
lama tunduk pada penjaga yang sama dengan data baru, dan tidak ada seam kedua yang harus diuji.

**ID-4** Pemecah dokumen *(JSON → objek)* dan penyusunnya *(objek → baris)* adalah rincian di dalam
`repository`, **bukan seam tersendiri**. Keduanya teruji lewat perjalanan pulang-pergi.

### Bentuk tabel

**ID-5** Tujuh tabel pada proporsional, sepuluh pada non-proporsional:

```
T_WORK_POLIS                            akar · LINTAS-LINI
 └ T_GENERAL_POLIS                      1:1 SHARED PK · satu baris per GENERASI
    ├ T_POLIS_QUOTATION                 1:1
    │   └ T_POLIS_CEDING                1:N
    ├ T_POLIS_INSTALMENT                1:N
    │   └ T_POLIS_INSTALMENT_DETAIL     1:N   ← hanya non-proporsional
    ├ T_POLIS_SPREADING                 1:N
    ├ T_POLIS_XOL                       1:N   ← hanya non-proporsional
    │   └ T_POLIS_XOL_LAYER             1:N   ← hanya non-proporsional
    └ POOLDATA.HISTORYAKSEPTASIPRODUCTION  1:N · sudah datar, tidak dibuat ulang
```

**ID-6** `[terverifikasi]` `T_WORK_POLIS` adalah tabel **lintas-lini** yang sudah dipakai
PremiumList. Baris NB dan baris EDM duduk **sejajar** di dalamnya dan **tidak saling menunjuk**;
kaitannya lewat `(NOPOLIS, PRODKE)` dan `OLD_POLIS_ID`, bukan lewat hubungan induk-anak.

**ID-7** `T_GENERAL_POLIS` berbagi kunci utama dengan `T_WORK_POLIS` — **tidak ada kolom kunci tamu
tersendiri**, sejalan pola `T_GENERAL_CLAIM` dan `T_PREMIUM_LIST` yang sudah ada.

### Kunci dan generasi

**ID-8** `[keputusan work owner]` Kunci alami `(NOPOLIS, PRODKE)`. NB selalu `PRODKE = 0`; EDM
`1, 2, 3, …`. Keunikan `(NOPOLIS, PRODKE)` ditegakkan di basis data.

⚠️ Ini memperbaiki bahaya nyata. `[terverifikasi]` `EDM\RDBList\SelectProdKe` mengambil nomor
generasi berikutnya dengan `order by prodke desc` lalu menambah satu, tanpa penguncian — dua
endorsemen serentak membaca angka yang sama. Di sistem baru yang kedua **gagal**, bukan bentrok
diam-diam.

**ID-9** `[keputusan work owner]` `T_GENERAL_POLIS.OLD_POLIS_ID → T_WORK_POLIS.ID`, nullable,
**di NB selalu kosong**. Keunikannya ditegakkan: satu generasi hanya boleh punya **satu** penerus,
sehingga rantai endorsemen tidak dapat bercabang.

`[terverifikasi]` Sistem lama **tidak punya penunjuk ini** — ia mencari tiap kali:
`EDM\RDBList\FetchPolisJsonPolis` dengan `where NOPOLIS = ? order by PRODKE desc fetch first 1 row
only`. Dan pencarian keduanya memakai `substr(nopolis, 1, 24)`, yang **menganggap nomor polis selalu
24 karakter**. Kunci tamu menghapus jebakan itu; panjang teks tidak lagi berpengaruh.

**ID-10** `[keputusan work owner]` Baris generasi lampau **tidak boleh disunting**. Ini yang
menggantikan pembekuan `OldData` *(P58)*, dan ditegakkan di `services`, bukan di tabel.

### `NOURUT`

**ID-11** `[keputusan work owner]` Setiap tabel anak membawa `NOURUT` — nomor urut baris di dalam
dokumen. Itulah kunci pasangan antar generasi.

**ID-12** `[keputusan work owner]` **Di NB baris boleh dihapus, dan `NOURUT` dinomori ulang** rapat
`1..n`. Aman, sebab pada `PRODKE 0` belum ada generasi untuk dipasangkan. Begitu generasi ditutup,
`NOURUT` beku; di EDM tidak ada penghapusan dan nomornya terbawa apa adanya.

**ID-13** Kunci dagang — `TREATY_NAME + CURRENCY_ID`, `INSTALLMENT_NO`,
`LAYER + LAYER_PART + ID_CURRENCY` — **tidak dipakai memasangkan**, hanya memeriksa. ⇒ syarat
"kunci harus unik dalam satu polis" **gugur**.

### Tipe kolom

**ID-14** `[keputusan work owner]` Konversi terjadi **sekali saat masuk**, bukan tiap kali dibaca:

| Golongan | Tipe | Dasar |
| --- | --- | --- |
| uang · persen | ⭐ **`NUMBER(38,8)`** *(semula ~~`NUMBER(20,8)`~~ — dinaikkan 23-09-2026 sore; **desimalnya tetap delapan**, yang berubah hanya sisi kiri koma)* | P29 disempurnakan 23-09-2026 · ADR-0003 |
| tanggal | `DATE` | dua format masuk |
| cacah | bilangan bulat | `NOURUT` `PRODKE` `INSTALLMENT_NO` |
| ⭐ kode | **teks** | nol di depan wajib utuh |
| ⭐ penanda | **teks** | `""` berbeda dari `"0"` |

**ID-15** `[keputusan work owner]` P29 — galat presisi pada data lama **diikuti apa adanya**, tidak
dibulatkan.

⚠️ **DISEMPURNAKAN 23-09-2026.** Bunyi lama dikutip: *"Itu sebabnya skalanya minimal sembilan
desimal."* Yang berlaku sekarang⭐ **`NUMBER(38,8)`** *(semula ~~`NUMBER(20,8)`~~ — dinaikkan 23-09-2026 sore; **desimalnya tetap delapan**, yang berubah hanya sisi kiri koma)*
lebih dari delapan **dibulatkan saat dimuat**; `ShareValue` yang bersimpan 24 desimal dan
`PremiumSpreaded` 20 desimal keduanya terpotong.

⭐ Ekor galat `2,76 × 10⁻⁷` **tetap terlihat** — ia jatuh di desimal ketujuh. ⛔ `[terbuka]`
Dua belas digit di depan koma **belum diuji** terhadap nilai terbesar yang pernah tersimpan.

`[terverifikasi]` Galat yang dimaksud nyata di data produksi: premi angsuran
`148.157.378,220000069` dikali empat menghasilkan `592.629.512,880000276`, sementara nilai
sebenarnya `592.629.512,88` — selisih `2,76 × 10⁻⁷`. Galat itu **lahir di rantai perhitungan, bukan
di penyimpanan**: besarnya 57 kali galat `float64` untuk bilangan sebesar itu, arahnya berlawanan,
dan pembagian desimalnya tepat tanpa sisa.

⚠️ Go memakai desimal, jadi perhitungan baru **tidak menghasilkan galat yang sama**. Satu kolom akan
memuat dua macam nilai — baris lama bergalat, baris baru bersih.

**ID-16** `[terverifikasi]` Kode **tetap teks**, sebab nol di depan membawa makna:
`GroupPanel = "006"` dan `BusinessOldId = "01"` adalah masukan penggolong `BusinessType_DeT` —
36 baris, 128 kode, bawaan `"UNKNOWN"`, berhenti di baris pertama yang cocok. Dijadikan bilangan,
`"006"` menjadi `6` dan pencarian 36 barisnya tidak menemukan apa pun.

**ID-17** `[terverifikasi]` Penanda **tetap teks**, sebab `IsApproved` dibandingkan sebagai teks —
`= 0` berarti ditolak, selain itu disetujui *(P24, P6)* — dan `""` adalah keadaan sah yang
**berbeda** dari `"0"`.

**ID-18** `[keputusan work owner]` Teks kosong `""` yang masuk kolom angka atau tanggal menjadi
**`NULL`**, bukan `0`. Nol adalah nilai uang yang bermakna, dan tanggal tidak punya nol.

**ID-19** `[terverifikasi]` Dua format tanggal masuk dalam satu dokumen: `YYYYMMDD` dan cap waktu
Pega bersufiks ` GMT`. Pengurai wajib menerima keduanya.

**ID-20** ⛔ Uang tidak pernah `float` *(ADR-0003)*. Perbandingan nilai uang **tidak boleh
sama-persis** — bertoleransi, atau dibandingkan dalam bentuk terbulatkan.

> ⛔ **RALAT** 2026-10-04 (P11, tiket 18) — bunyi lama ID-20: *"Perbandingan nilai uang **tidak boleh
> sama-persis** — bertoleransi, atau dibandingkan dalam bentuk terbulatkan."* → bunyi baru: **Pembandingan
> DUA nilai uang satu sama lain tidak boleh sama-persis** — bertoleransi, atau dalam bentuk terbulatkan.
> Pembandingan nilai uang **lawan tetapan nol** (tanda: `>= 0`, `< 0`, `!= 0`, `> 0`, kosong / `"0"`)
> mengikuti XML **apa adanya, eksak** — toleransi di sana mengubah perilaku rule. Bukti (sisir 176 rule
> terjangkau, tiket 18 bab P11): nol pembandingan dua nilai uang yang hidup; yang ada hanya uang lawan nol
> (`Activity/SetDueTo_act.xml` langkah 1–2 `.BalanceDueTo>=0` / `<0`; `Activity/CountNetPremi_act.xml` langkah 6
> `.Deduction1!="" || .Deduction1!=0`; `Activity/CountOGPONP_Act.xml` langkah 1–2 `.PremiOgp == "" ||.PremiOgp ==  "0"`
> dan langkah 8 `.Claim!=0&&.Claim!=""`; `Activity/InputPolicyTreatyInDetail_preACT.xml` langkah 17 `@if(.Limit>0,"IDR","")`;
> `Section/DetailPolicyTreatyIn.xml` `.Claim != '' && .Claim != 0`; `Section/DetailDeptHeadTreatyIn_UW.xml` `.BalanceDueTo < 0` /
> `>= 0`). Satu-satunya pembandingan bertoleransi atas rasio dua nilai uang, `((.ResultOgp2/.PremiOgp)-.OveriddingCommOgp)<=0.01`
> dan tiga saudaranya, ada di langkah berlabel `//` (dinonaktifkan): `CountOverridingCommOgp_Act` 1/2/4,
> `CountOverridingCommOnp_Act` 1/2, `CountRiCommOgp_act` 4, `CountRiCommOnp_act` 1; `Local.NETPREMI>200000000.00` hanya di
> `CekLimitTreatyAcc_Act` (tidak dibangun, K2). Uji: `backend/models/pembandingan_uang_test.go`
> `TestPortTidakMembandingkanDuaNilaiUang` (penjaga AST: setiap `.Cmp` lawan tetapan rule, nol teks medan uang
> dibandingkan persis dengan teks medan lain) dan `TestTandaUangLawanNolEksakSepertiXML`.

### Isi tabel

**ID-21** `T_GENERAL_POLIS` memuat 79 medan skalar tingkat atas `PolicyTreatyIn` — ⭐ ditambah **`REMARK`** *(panjang 128)*, medan yang sensus korpus lewatkan dan baru terlihat dari data guide — ditambah tujuh
kolom yang sudah datar di `POOLDATA.json_polis`: `IDPEGA` `NOPOLIS` `NOENDORS` `PRODKE` `TGL_INPUT`
`TGL_PROD` `USERNAME`.

> ⛔ **RALAT putaran 2 — 03-10-2026 (kolom ikut diagram, bab 0 butir 12)** atas ID-10, ID-21, ID-23,
> ID-24. Rinciannya per kolom: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`; ditagih
> `TestTabelDanKolomMengikutiDiagramGrilling`.
>
> - **ID-21** bunyi lama: *"79 medan skalar tingkat atas `PolicyTreatyIn` — ditambah `REMARK` — ditambah
>   tujuh kolom … `json_polis`"*. Bunyi baru: 79 = **69 kolom** katalog (termasuk `TGL_PROD` dan
>   `REMARK`) + `NOPOLIS` + 9 tanpa kolom (`LAYER*` 4 dicoret diagram F26, `Total*` 4 turunan,
>   `isApprovedtoDeptHead` P36). Ditambah tiga kolom halaman kerja `POSITION_NOTE` `NB_STATUS`
>   `TREATY_IN_ID` (RALAT: dibaca connector Flow / tampil portal / dibaca `BrowseTreatyIn`).
>   `IS_OJK_NOPOLIS` dan `BROKERAGE_FEE` **dibuang** (hanya ditulis).
> - **ID-23** bunyi lama: *"`T_POLIS_QUOTATION` — 1:1, 10 medan"*. Bunyi baru: 10 medan diagram **+ 6
>   RALAT** yang XML buktikan dibaca/tampil: `BusinessName` (`InputPolicyTreatyInPre_Act` langkah 2;
>   portal), `BusinessFac` (`SaveViewSuggest` CARI7; `GetListOpportunity` filter E), `InsuredID` dan
>   `InsuredName` (`InputPolicyTreatyInDetail_preACT` langkah 3, 14.1, 14.3; portal), `NoOfferSlip` dan
>   `IsSurveyReport` (tampil di Section NB, AC 64). Sebelas kolom putaran 1 lainnya dibuang.
> - **ID-24** bunyi lama: *"`.CedingCo` sebagai id"*. Bunyi baru: kolomnya `CEDING_CO_ID` (diagram R43) dan
>   tabel menunjuk induknya lewat `QUOTATION_ID → T_POLIS_QUOTATION` (diagram O39), bukan `POLIS_ID`.
> - **ID-10**: tanpa kolom penanda tutup. Putaran 1 menambah `TGL_TUTUP` (di luar diagram) — **dibuang**;
>   generasi tertutup = ada baris penerus yang `OLD_POLIS_ID`-nya menunjuknya.
> - Angsuran (ID-26): `T_POLIS_INSTALMENT` = rancangan §4.3 tanpa `PAYMENT_DATE` (diagram R61) + RALAT
>   `PPN` `PPH` `PAYMENT_TOTAL_AFTER_PPN` `PAYMENT_TOTAL_AFTER_TAX` (dibaca rumus preACT 18.3.4.1);
>   `T_POLIS_INSTALMENT_DETAIL` = rancangan §4.3 (11), tanpa `PPN`/`PPH` anak (hanya ditulis).

**ID-22** `[terverifikasi]` ⛔ `LAYER` `LAYER_TYPE` `LAYER_PART` `LAYER_PART_TYPE` **tidak
disimpan** di `T_GENERAL_POLIS`. Di sistem lama keempatnya **pantulan**, dibaca balik dari kolom
tabel: `InputPolicyTreatyInDetail_preACT` menetapkan
`POLIS.LayerType = pyReportContentPage.pxResults(1).LAYERTYPE`.

⚠️ Dan `pxResults(1)` hanya mengambil **baris pertama** ⇒ nilai di tingkat polis itu **layer pertama
saja**, bukan ringkasan seluruh layer. Laporan yang memakainya sebagai "layer polis ini" mengabaikan
layer kedua dan seterusnya.

**ID-23** `T_POLIS_QUOTATION` — 1:1, 10 medan: `ProportionalType` `MOID` `BusinessCode`
`BusinessOldId` `GroupPanel` `SourceOfBusiness` `Type` `EdmType` `OldPolicyNo` `MarketingName`.

**ID-24** `T_POLIS_CEDING` ← `QuotationData.CedingCoList()`, dua medan: `.CedingCo` sebagai id dan
`.CedingCoName` sebagai nama.

**ID-25** `[keputusan work owner]` `CEDING_CO_NAME` dan `CEDING_CO` di `T_GENERAL_POLIS`
**disalin apa adanya** dari dokumen, **tidak pernah dirangkai ulang** dari `T_POLIS_CEDING`.

`[terverifikasi]` Keduanya daftar yang digabung, berakhiran `"; "`. Perangkaiannya terjadi di
`SetCedingCo_Act` dengan pola `@If(kosong, nilai baris, sambung)`, dan penghapusannya di
`DeleteCeding_Act` dengan `@replaceAll` pada teks gabungan — sehingga menghapus `"PT A"` dari
`"PT ABC; PT A; "` juga merusak `"PT ABC"`. Nilai tersimpan karena itu **tidak dapat direproduksi**
dengan merangkai ulang.

⚠️ `[penyimpangan sadar]` Konsekuensi yang diterima work owner: bentuk gabungan dan `T_POLIS_CEDING`
**bisa tidak sinkron**, dan tidak ada penjaganya — sama seperti sistem lama. Penghapusan berbasis
ganti-teks **tidak ditiru**; di sistem baru cukup hapus barisnya.

**ID-26** `[keputusan work owner]` Pada bentuk **proporsional**, `T_POLIS_INSTALMENT` **satu
tingkat**: dokumen berhenti di `ListInstallment`, tanpa sarang `InstallmentList`.

`[terverifikasi]` Croscheck korpus mendukungnya. Tujuh berkas menyentuh
`.ListInstallment().InstallmentList()`, dan ketujuhnya terjelaskan: lima bernama `_NonProp`/`_NP`;
`FillPaymentInstallmentEDMT` yang keterangannya sendiri berbunyi *"Generate Installment based on
XOLDifferenceList"*; dan `SetInstallmentValue`, yang lewat `FillMasterInstallment` **menyalin dari
master kontrak** `TreatyIn.Installment(n).InstallmentList` — dari sanalah sarang itu berasal.

**ID-27** ⭐ **Tabel penampung medan tak dikenal.** Sensus medan berasal dari rujukan di dalam
aturan, dan aturan hanya menyentuh medan yang dihitung atau ditampilkan. Medan yang **hanya lewat**
tidak terlihat. Tiga contoh dokumen nyata memberi 37, 64 dan 47 medan skalar tingkat atas dengan
gabungan **74**, sementara sapuan korpus memberi **79** — ⛔ kedua angka **tidak dapat dipertemukan**
tanpa `JSON_DATAGUIDE`.

⇒ Pemuat menyimpan medan yang tidak dikenal alih-alih membuangnya, dan penampung itu **wajib kosong
sebelum rancangan dinyatakan selesai**.

> ⛔ **RALAT putaran 2 — 03-10-2026** `[keputusan work owner]` **K17** (PROMPT-NB-TREATY-IN-PUTARAN-2.md
> bab 2). Bunyi lama: *"⭐ **Tabel penampung medan tak dikenal.**"* — dibangun putaran 1 sebagai
> `T_POLIS_MEDAN_LAIN` (migrasi 329). Bunyi baru: tabel itu **tidak ada di diagram grilling** dan
> dihapus. Penampungnya **berkas laporan CSV per jalankan pemuat** — `POLIS_ID`, `JALUR`, `NILAI` — di
> folder keluaran yang ditentukan operator (`MODUL.md` *Pemuat dokumen lama*, tiket 22). Medan tetap
> tersimpan beserta nilainya, jumlahnya dicetak, dan **wajib 0** sebelum pekerjaan dinyatakan selesai.

> ⛔ **RALAT putaran 3 — 04-10-2026** `[keputusan work owner]` **F3** (PROMPT-NB-TREATY-IN-PUTARAN-3.md
> bab 2). Bunyi lama (RALAT K17 di atas), dikutip: *"Penampungnya **berkas laporan CSV per jalankan
> pemuat** — `POLIS_ID`, `JALUR`, `NILAI` — … Medan tetap tersimpan beserta nilainya, jumlahnya dicetak,
> dan **wajib 0** sebelum pekerjaan dinyatakan selesai."* Bunyi baru: **tidak ada penampung.** Setiap
> medan dokumen lama tanpa kolom **diputuskan per medan** dari XML: (a) dibaca rule NB terjangkau →
> berkolom di salah satu dari 8 tabel (bab 0 butir 12; satu-satunya: `PolicyTreatyIn.EDMType` →
> `T_GENERAL_POLIS.EDM_TYPE`, dibaca syarat `InputPolicyTreatyInPre_Act` langkah 10), atau (b) tidak
> dibaca → **dibuang** dengan alasan + bukti tertulis (`backend/models/medan_abaikan_lama.json` bagian
> `pola`). `SuggestList` dokumen lama **disalin** ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (ID-31) dengan
> penjaga dobel menurut `IDPEGA`. Berkas CSV pemuat tetap ditulis sebagai **arsip audit pemuatan**
> (`nbtreatyin-arsip-medan-<stempel>.csv`: `POLIS_ID`, `JALUR`, `NILAI`, `KEPUTUSAN`) untuk setiap medan
> yang tidak masuk kolom — bukan penampung untuk dimigrasi kelak. Yang wajib **0** sebelum pekerjaan
> dinyatakan selesai: medan **belum diputuskan** (`KEPUTUSAN` = `BELUM DIPUTUSKAN`).

**ID-27b** ⛔ **DICABUT 23-09-2026.** `[keputusan work owner]` Tabel sebaran tambahan **tidak ada** — keempat medannya turunan. Bunyi lamanya dikutip di blok kepala berkas ini.

⚠️ Sensus korpus **melewatkannya** — tidak satu pun aturan merujuk anggotanya. Ini persis jebakan
"medan yang hanya lewat" yang diperingatkan Bab 1, dan pembuktian bahwa peringatan itu bukan
kehati-hatian berlebihan.

⚠️ Mirip `T_POLIS_SPREADING` tetapi **tidak sama**: spreading punya `CURRENCY` dan `CURRENCY_ID`,
breakdown tidak. ⛔ `[terbuka]` Apa bedanya secara dagang — pemiliknya `[work owner]`.

**ID-28** `T_POLIS_SPREADING` — `TREATY_NAME` `TREATY_TYPE` `CURRENCY` `CURRENCY_ID`
`SHARE_PERCENTAGE` `SPLIT_RNM_SHARE_PCT` `CLAIM_PERCENTAGE` *(persen)* · `PREMIUM_SPREADED`
`CLAIM_SPREADED` *(uang)*.

`[terverifikasi]` **P60** — NB membagi `100 / jumlah baris` dengan presisi **10**; EDM menetapkan
baris pertama `SPLIT_RNM_SHARE_PCT = 100` lalu membagi dengan presisi **20**. Ditiru apa adanya, dan
perbedaannya wajib punya test tersendiri.

**ID-29** `[terverifikasi]` `T_POLIS_XOL` induk **tanpa penanda layer sama sekali** — 13 medan,
seluruhnya uang. `T_POLIS_XOL_LAYER` yang memegang `LAYER` `LAYER_TYPE` `LAYER_PART`
`LAYER_PART_TYPE`, 12 rujukan masing-masing, di induk **nol**.

Pemisahannya bukan selera: uji korpus menemukan **9 berkas memakai kedua tingkat** dan **0 memakai
induk saja**.

**ID-30** `[terverifikasi]` `DEDUCTION` pada XOL adalah **uang**, bukan persen. Buktinya ia
**dijumlahkan** antar layer — `local.Deduction = local.Deduction + .Deduction` — dan persentase tidak
dijumlahkan antar layer. Polanya identik dengan `GrossPremi` `NetPremi` `DueToValue`.

**ID-31** `POOLDATA.HISTORYAKSEPTASIPRODUCTION` **tidak dibuat ulang** — sudah datar, 15 kolom.
`[terverifikasi]` Isinya berasal dari `pyWorkPage.PolicyTreatyIn.SuggestList` lewat
`Activity\SaveViewSuggest → RDBList\InsertViewSuggest_SQL`, satu-satunya pemanggilnya.

Pemetaan yang mengikat: `.Suggest → KETERANGAN` dipotong `substr(…, 0, 3990)` ·
`.IsApproved → APPROVAL` dengan `"1"` = Accept dan `"0"` = Reject · `AKSES_LOGIN` dari
`OperatorID.pyUserIdentifier` *(P4)* · `PIC` dari nama tampilan *(P33)*.

⭐ `IsApproved` di sini **per baris usulan**, berbeda dari `IsApproved` tingkat polis.
Menggabungkan keduanya adalah cacat.

> ⭐ **Penerapan putaran 2 — 03-10-2026** `[keputusan work owner]` **K4**: ID-31 berlaku apa adanya.
> `SuggestList` ditulis ke tabel lama ini (`repository/usulan.go`, pemetaan `models/usulan.go`) dan
> dibaca balik untuk layar. `[penyimpangan sadar]` terhadap XML, dasar grilling ID-31/AC 39 dan K4:
> (1) syarat `Quotation.BusinessFac == "F"` `SaveViewSuggest` langkah 2 tidak ditiru; (2) baris ditulis
> pada submit yang menambahkannya di **ketiga** jenjang, transaksi yang sama — XML hanya memanggilnya dari
> `InputPolicyTreatyInPost_Act` langkah 4 (pasca-submit admin); (3) `NOURUT` = berikutnya per `IDPEGA` di
> bawah kunci kasus (XML: `.pxListSubscript`); (4) `TGL_INP` jam 24 (XML memformat `hh` lalu
> `To_date(…,'HH24…')` — catatan sore tersimpan pagi); (5) `DIV` (`OperatorID.pyOrgDivision`) NULL — tanpa
> sumber di `inti.Pelaku` (butir terbuka). Tabel dideklarasikan *Tabel warisan* di `MODUL.md`.

### Penulisan dan transaksi

**ID-32** `[keputusan work owner]` Seluruh urutan penyimpanan dibungkus **satu transaksi** *(P2)*.
Kegagalan pada tabel mana pun membatalkan seluruhnya.

**ID-33** `[keputusan work owner]` Setiap query menulis skema `POOLDATA.` **eksplisit** *(P3)*.

⚠️ `[terverifikasi]` Sistem lama tidak konsisten — pada EDM hanya 16 dari 31 pernyataan `FROM`
memakai skema, `JSON_POLIS` muncul dua ejaan, dan `DATAPEGA` adalah skema kedua. Ketetapan ini
**berlaku untuk penulisan baru**, bukan pemerian sistem lama.

**ID-34** `[keputusan work owner]` `OPERATORID` diambil dari identitas akses login, `PIC` dari nama
tampilan — dua kolom berbeda, tidak boleh tertukar *(P4, P33)*.

**ID-35** `services` **menolak** baris XOL pada polis proporsional, dan menolak generasi yang tidak
memuat setiap `NOURUT` generasi sebelumnya. Aturan bentuk hidup di `services`, bukan di tabel.

---

## 6 · Acceptance Criteria

Bab ini **tidak memutuskan apa pun**. Ia menyatakan ulang keputusan pada Bab 5 dalam bentuk yang
dapat diuji dari luar.

### Kunci dan generasi

1. `[terverifikasi]` Menyimpan dua polis dengan `NOPOLIS` sama dan `PRODKE` sama **ditolak**. Test
   yang menemukan keduanya tersimpan **gagal**. *(Bab 5, ID-8)*
2. `[terverifikasi]` Polis NB tersimpan dengan `PRODKE = 0`. Test yang menemukan nilai lain
   **gagal**. *(ID-8)*
3. `[terverifikasi]` `OLD_POLIS_ID` pada polis NB bernilai `NULL`. Test yang menemukan nilai terisi
   **gagal**. *(ID-9)*
4. `[terverifikasi]` Dua baris berbeda tidak boleh memiliki `OLD_POLIS_ID` yang sama. Test yang
   menemukan dua penerus untuk satu generasi **gagal**. *(ID-9)*
5. `[terverifikasi]` `T_GENERAL_POLIS` dan `T_WORK_POLIS` berbagi kunci utama; tidak ada kolom kunci
   tamu tersendiri. Test yang menemukan kolom penyambung terpisah **gagal**. *(ID-7)*
6. `[terverifikasi]` Menyunting baris generasi yang sudah ditutup **ditolak**. Test yang menemukan
   perubahan tersimpan **gagal**. *(ID-10)*
7. `[terverifikasi]` Baris NB dan baris EDM pada `T_WORK_POLIS` tidak saling menunjuk. Test yang
   menemukan kunci tamu di antara keduanya **gagal**. *(ID-6)*

### `NOURUT`

8. `[terverifikasi]` Setiap tabel anak memiliki kolom `NOURUT`. Test yang menemukan tabel anak tanpa
   kolom itu **gagal**. *(ID-11)*
9. `[terverifikasi]` Pada satu polis NB dengan tiga baris angsuran, menghapus baris kedua
   menghasilkan `NOURUT` `1` dan `2`. Test yang menemukan `1` dan `3` **gagal**. *(ID-12)*
10. `[terverifikasi]` `NOURUT` unik dalam satu induk. Test yang menemukan dua baris ber-`NOURUT` sama
    **gagal**. *(ID-11)*
11. `[terverifikasi]` Pemasangan antar generasi memakai `NOURUT`, bukan kunci dagang. Test yang
    menemukan pemasangan berdasarkan `TREATY_NAME` **gagal**. *(ID-13)*

### Tipe kolom

12. `[terverifikasi]` `GroupPanel` bernilai `"006"` tersimpan dan terbaca kembali sebagai `"006"`.
    Test yang menemukan `"6"` atau `6` **gagal**. *(ID-16)*
13. `[terverifikasi]` `BusinessOldId` bernilai `"01"` tersimpan utuh. Test yang menemukan `"1"`
    **gagal**. *(ID-16)*
14. `[terverifikasi]` Penggolong `BusinessType_DeT` dengan masukan `GroupPanel = "006"` dan
    `BusinessOldId = "01"` menemukan barisnya. Test yang menemukan bawaan `"UNKNOWN"` **gagal**.
    *(ID-16)*
15. `[terverifikasi]` `IsApproved` bernilai `""` tersimpan sebagai `""`, bukan `NULL` dan bukan
    `"0"`. Test yang menemukan ketiganya menyatu **gagal**. *(ID-17)*
    > ⛔ **RALAT** 2026-10-04 (P11, tiket 18) — bunyi lama dikutip: *"`IsApproved` bernilai `""` tersimpan sebagai
    > `""`, bukan `NULL` dan bukan `"0"`. Test yang menemukan ketiganya menyatu **gagal**."* → bunyi baru:
    > **`IsApproved` bernilai `""` tersimpan sebagai `NULL` dan terbaca kembali sebagai `""`; `"0"` tersimpan
    > `'0'` dan terbaca kembali `"0"`. Test yang menemukan `""` tersimpan atau terbaca sebagai `"0"`, atau `"0"`
    > tersimpan NULL / terbaca `""`, gagal.** Sebab: Oracle menyamakan teks kosong `''` dengan `NULL` pada `VARCHAR2`
    > (*Oracle Database SQL Language Reference*, bab "Nulls": nilai karakter berpanjang nol diperlakukan sebagai null)
    > — "tersimpan sebagai `""`, bukan `NULL`" mustahil secara fisik. Yang dijaga ID-17 tetap utuh: `""` dan `"0"`
    > tidak pernah menyatu, di kolom (`NULL` lawan `'0'`) maupun di halaman (`""` lawan `"0"`), sebab XML
    > membandingkannya sebagai teks: `DecisionTable/isApproved.xml` kolom `pyWorkPage.PolicyTreatyIn.IsApproved`
    > bertipe `text`, satu baris `= 0` → `No`, selain itu bawaan (disetujui) — `""` disetujui, `"0"` ditolak (AC 16).
    > Uji: `backend/repository/kolom_test.go` `TestNilaiTulisKosongJadiNULL` (penanda `""` → NULL),
    > `TestNilaiBaca` (NULL → `""`); `backend/repository/penyimpanan_db_test.go` `TestIsApprovedKosongDanNolTetapBerbeda`
    > (kolom dibaca langsung, bertag `db`, belum dijalankan — K11); `backend/models/tangga_test.go`
    > `TestIsApprovedSelainNolDisetujuiTermasukKosong`.
16. `[terverifikasi]` `IsApproved` bernilai `"0"` menghasilkan keputusan **ditolak**; nilai lain
    menghasilkan **disetujui**. Test yang menemukan sebaliknya **gagal**. *(ID-17)*
17. `[terverifikasi]` Teks kosong `""` pada medan uang tersimpan sebagai `NULL`. Test yang menemukan
    `0` **gagal**. *(ID-18)*
18. `[terverifikasi]` Teks kosong `""` pada medan tanggal tersimpan sebagai `NULL`. Test yang
    menemukan tanggal bawaan **gagal**. *(ID-18)*
19. `[keputusan work owner]` Nilai uang `592629512.880000276` tersimpan sebagai
    `592629512.88000028` — dibulatkan pada desimal kedelapan, **bukan** dipotong ke
    `592629512.88`. Test yang menemukan `592629512.88` **gagal**. *(ID-15)*
    ⚠️ Bunyi lama dikutip: ~~*"tersimpan dan terbaca kembali tanpa kehilangan satu digit pun"*~~.
20. `[keputusan work owner]` Kolom uang bertipe `NUMBER(38,8)` *(semula ~~`NUMBER(20,8)`~~)* — delapan angka di belakang koma,
    dua belas di depan. Test yang menemukan skala lain **gagal**. *(ID-14)*
    ⚠️ Bunyi lama dikutip: ~~*"menerima sekurangnya sembilan angka di belakang koma"*~~.
20b. `[terverifikasi]` Nilai berdesimal lebih dari delapan **dibulatkan, bukan ditolak**. Test yang
    menemukan kegagalan pemuatan pada `ShareValue` berdesimal 24 **gagal**. *(ID-15)*
21. `[terverifikasi]` Tanggal berformat `YYYYMMDD` terurai benar. Test yang menemukan kegagalan urai
    **gagal**. *(ID-19)*
22. `[terverifikasi]` Cap waktu Pega bersufiks ` GMT` terurai benar. Test yang menemukan kegagalan
    urai **gagal**. *(ID-19)*
23. `[terverifikasi]` Mengurutkan polis menurut tanggal mulai menghasilkan urutan kronologis, bukan
    urutan teks. Test yang menemukan urutan leksikal **gagal**. *(ID-14)*
24. `[terverifikasi]` Tidak ada kolom uang bertipe `float` di seluruh skema. Test yang menemukan satu
    pun **gagal**. *(ID-20)*
25. `[terverifikasi]` Pembandingan dua nilai uang memakai toleransi atau bentuk terbulatkan, bukan
    kesamaan persis. Test yang menemukan pembandingan persis **gagal**. *(ID-20)*
    > ⛔ **RALAT** 2026-10-04 (P11, tiket 18) — bunyi lama dikutip: *"Pembandingan dua nilai uang memakai toleransi
    > atau bentuk terbulatkan, bukan kesamaan persis. Test yang menemukan pembandingan persis **gagal**."* → bunyi
    > baru: **Port rumus tidak membandingkan dua nilai uang satu sama lain — 176 rule terjangkau pun tidak; bila
    > kelak ada, ia bertoleransi atau terbulatkan. Pembandingan uang lawan nol mengikuti XML eksak. Test yang
    > menemukan pembandingan persis antara dua nilai uang di port, atau toleransi pada pembandingan lawan nol,
    > gagal.** Bukti XML dan uji: RALAT ID-20 (bab 5) dan tiket 18 bab P11.

### Isi tabel

26. `[terverifikasi]` `T_GENERAL_POLIS` tidak memiliki kolom `LAYER` `LAYER_TYPE` `LAYER_PART`
    `LAYER_PART_TYPE`. Test yang menemukan satu pun **gagal**. *(ID-22)*
27. `[terverifikasi]` `T_POLIS_QUOTATION` menyimpan sepuluh medan `QuotationData`, termasuk
    `OldPolicyNo` dan `GroupPanel`. Test yang menemukan kurang **gagal**. *(ID-23)*
28. `[terverifikasi]` `T_POLIS_CEDING` menyimpan satu baris per ceding, dengan id dan nama terpisah.
    Test yang menemukan keduanya dalam satu kolom **gagal**. *(ID-24)*
29. `[keputusan work owner]` `CEDING_CO_NAME` pada `T_GENERAL_POLIS` tersimpan persis seperti di
    dokumen, termasuk ekor `"; "`. Test yang menemukan hasil rangkaian ulang **gagal**. *(ID-25)*
30. `[keputusan work owner]` Menghapus satu ceding menghapus **barisnya**, bukan mengganti teks
    gabungan. Test yang menemukan `@replaceAll` atau padanannya **gagal**. *(ID-25)*
31. `[terverifikasi]` Polis proporsional tersimpan **tanpa** baris `T_POLIS_INSTALMENT_DETAIL`. Test
    yang menemukan baris di sana **gagal**. *(ID-26)*
32. `[terverifikasi]` Polis proporsional tersimpan **tanpa** baris `T_POLIS_XOL` dan
    `T_POLIS_XOL_LAYER`. Test yang menemukan baris di sana **gagal**. *(ID-5)*
33. `[terverifikasi]` Menyimpan baris XOL pada polis ber-`ProportionalType = 'Proportional'`
    **ditolak**. Test yang menemukan baris tersimpan **gagal**. *(ID-35)*
34. `[terverifikasi]` Polis non-proporsional menyimpan `LAYER*` pada `T_POLIS_XOL_LAYER`, bukan pada
    `T_POLIS_XOL`. Test yang menemukan `LAYER` terisi di induk **gagal**. *(ID-29)*
35. `[terverifikasi]` `DEDUCTION` pada `T_POLIS_XOL` dan `T_POLIS_XOL_LAYER` bertipe uang, dengan
    presisi sama dengan kolom uang lainnya. Test yang menemukan tipe persen **gagal**. *(ID-30)*
36. `[terverifikasi]` `T_POLIS_SPREADING` pada NB membagi `100 / jumlah baris` dengan presisi 10.
    Test yang menemukan presisi 20 **gagal**. *(ID-28)*
37. `[terverifikasi]` `SHARE_PERCENTAGE` `CLAIM_PERCENTAGE` `SPLIT_RNM_SHARE_PCT` disimpan sebagai
    persen; nilai `12.5` berarti 12,5 persen. Test yang memperlakukannya sebagai pecahan **gagal**.
    *(ID-28)*
38. `[terverifikasi]` `DEDUCTION1` `DEDUCTION2` `TOTAL_SHARE_PERCENTAGE_PREMIUM`
    `TOTAL_SHARE_PERCENTAGE_CLAIM` disimpan sebagai persen. Test yang memperlakukannya sebagai uang
    **gagal**. *(Bab 5, ketetapan lama P29)*
    > ⛔ **RALAT** 2026-10-03 (K3, paket P2) — bunyi lama: *"`DEDUCTION1` `DEDUCTION2`
    > `TOTAL_SHARE_PERCENTAGE_PREMIUM` `TOTAL_SHARE_PERCENTAGE_CLAIM` disimpan sebagai persen. Test
    > yang memperlakukannya sebagai uang **gagal**"* → bunyi baru: **`DEDUCTION1` `DEDUCTION2`
    > `T_GENERAL_POLIS` bergolongan uang** (`NUMBER(38,8)`, sama tipe fisiknya), mengikuti pemakaian
    > XML (keputusan WO K3); `TOTAL_SHARE_PERCENTAGE_*` tetap persen. Test yang memperlakukan
    > `DEDUCTION1/2` sebagai persen **gagal**. Bukti: sel `.Deduction1` `.Deduction2` `pxCurrency` di `Section/DetailPolicyTreatyIn.xml` dan `Section/DetailDeptHeadTreatyIn_UW.xml`; `Activity/CountNetPremi_act` langkah 4 mengurangkan keduanya dari premi; `Activity/SetPPNPPH` langkah 4 membagi `.Deduction1` dengan 1,022. Kode: `models/katalog.go` (`kUang`).
39. `[terverifikasi]` `HISTORYAKSEPTASIPRODUCTION` tidak dibuat ulang; polis menulis ke tabel yang
    sudah ada. Test yang menemukan tabel riwayat baru **gagal**. *(ID-31)*
40. `[terverifikasi]` `KETERANGAN` dipotong pada 3990 karakter. Test yang menemukan nilai lebih
    panjang tersimpan **gagal**. *(ID-31)*
41. `[terverifikasi]` `APPROVAL` bernilai `Accept` untuk `.IsApproved = "1"` dan `Reject` untuk
    `"0"`. Test yang menemukan pemetaan lain **gagal**. *(ID-31)*
42. `[terverifikasi]` `IsApproved` per baris usulan tersimpan terpisah dari `IsApproved` tingkat
    polis. Test yang menemukan keduanya dibaca dari kolom sama **gagal**. *(ID-31)*
43. `[terverifikasi]` `AKSES_LOGIN` terisi dari identitas login, bukan dari masukan pengguna. Test
    yang menemukan nilai dapat diketik **gagal**. *(ID-34)*
44. `[terverifikasi]` `PIC` terisi dari nama tampilan, bukan dari identitas akses. Test yang
    menemukan keduanya bernilai sama **gagal**. *(ID-34)*

### Transaksi dan penulisan

45. `[terverifikasi]` Kegagalan menulis tabel anak mana pun membatalkan seluruh penyimpanan. Test
    yang menemukan baris induk tersisa **gagal**. *(ID-32)*
46. `[terverifikasi]` Seluruh penyimpanan satu polis terjadi dalam satu transaksi. Test yang
    menemukan lebih dari satu komit **gagal**. *(ID-32)*
47. `[terverifikasi]` Setiap pernyataan SQL menyebut skema `POOLDATA.` eksplisit. Test yang menemukan
    nama tabel tanpa skema **gagal**. *(ID-33)*
48. `[terverifikasi]` Penyimpanan tidak memanggil satu pun stored procedure. Test yang menemukan
    pemanggilan procedure **gagal**. *(Problem Statement)*

### Pulang-pergi dan bentuk

49. `[terverifikasi]` Polis proporsional yang disimpan lalu dibaca kembali menghasilkan nilai yang
    sama pada seluruh medan. Test yang menemukan satu medan berbeda **gagal**. *(ID-2)*
50. `[terverifikasi]` Polis non-proporsional yang disimpan lalu dibaca kembali menghasilkan nilai
    yang sama, termasuk seluruh layer. Test yang menemukan layer hilang **gagal**. *(ID-2)*
51. `[terverifikasi]` Urutan baris anak saat dibaca sama dengan urutan `NOURUT`. Test yang menemukan
    urutan acak **gagal**. *(ID-11)*
52. `[terverifikasi]` Dokumen berbentuk `ListInstallment` datar terurai benar. Test yang menemukan
    kegagalan **gagal**. *(ID-26)*
53. `[terverifikasi]` Dokumen berbentuk `ListInstallment` bersarang terurai benar pada bentuk
    non-proporsional. Test yang menemukan kegagalan **gagal**. *(ID-26)*
54. `[terverifikasi]` Test yang hanya mencakup satu bentuk dokumen dinyatakan **tidak memadai**.
    *(ID-26)*

### Migrasi

55. `[keputusan work owner]` Dokumen lama dimuat **tanpa pembulatan ke presisi mata uang**;
    pembulatan hanya terjadi pada desimal kesembilan ke atas, mengikuti `NUMBER(38,8)` *(semula ~~`NUMBER(20,8)`~~)*. Test yang
    menemukan nilai dibulatkan ke dua desimal **gagal**. *(ID-15)*
    ⚠️ Bunyi lama dikutip: ~~*"dimuat tanpa mengubah satu pun nilai uang"*~~ — tidak lagi benar
    secara harfiah sejak presisi ditetapkan `NUMBER(38,8)` *(semula ~~`NUMBER(20,8)`~~)*.
56. `[terverifikasi]` Pemuat migrasi menulis lewat antarmuka `repository` yang sama dengan
    penyimpanan biasa. Test yang menemukan jalur tulis terpisah **gagal**. *(ID-3)*
57. `[terverifikasi]` Medan dokumen yang tidak dikenal tersimpan di penampung, bukan dibuang. Test
    yang menemukan medan hilang **gagal**. *(ID-27)*
58. `[terverifikasi]` Dokumen yang gagal diurai dilaporkan beserta sebabnya, bukan dilewati diam-diam.
    Test yang menemukan kegagalan senyap **gagal**. *(ID-3)*
59. `[terbuka]` Penampung medan tak dikenal **wajib kosong** sebelum rancangan dinyatakan selesai.
    *(ID-27)*
    ⛔ **RALAT AC 57 dan 59 — putaran 2, 03-10-2026** (K17). Bunyi lama AC 57: *"Medan dokumen yang tidak
    dikenal tersimpan di penampung"*; AC 59: *"Penampung medan tak dikenal wajib kosong"*. Bunyi baru:
    "penampung" = **berkas CSV** `POLIS_ID,JALUR,NILAI` per jalankan pemuat, bukan tabel. AC 57: test
    yang menemukan medan tak dikenal tidak tertulis di berkas itu beserta nilainya **gagal**. AC 59:
    berkas itu wajib **nol baris data** sebelum pekerjaan dinyatakan selesai; pemuat mencetak jumlahnya
    dan keluar dengan kode bukan nol selama jumlahnya > 0.
    ⛔ **RALAT AC 57 dan 59 — putaran 3, 04-10-2026** (`[keputusan work owner]` **F3**, RALAT ID-27). Bunyi
    lama (RALAT K17 di atas), dikutip: *"AC 57: test yang menemukan medan tak dikenal tidak tertulis di
    berkas itu beserta nilainya **gagal**. AC 59: berkas itu wajib **nol baris data** sebelum pekerjaan
    dinyatakan selesai"*. Bunyi baru: berkas CSV = **arsip audit pemuatan** `POLIS_ID,JALUR,NILAI,KEPUTUSAN`,
    bukan penampung. **AC 57**: test yang menemukan medan dokumen yang tidak masuk kolom tetapi tidak
    tertulis di arsip beserta nilai dan keputusannya **gagal** (`TestLaporanArsipMedanTanpaKolomBerkasCSV`).
    **AC 59**: **nol medan yang belum diputuskan** sebelum pekerjaan dinyatakan selesai — setiap medan
    berkolom, disalin (`SuggestList` → `HISTORYAKSEPTASIPRODUCTION`), atau dibuang dengan alasan + bukti
    XML (`backend/models/medan_abaikan_lama.json`); pemuat keluar dengan kode bukan nol hanya bila ada medan
    `BELUM DIPUTUSKAN` atau dokumen gagal (`TestRingkasanSelesaiHanyaBilaNolGalatDanNolBelumDiputuskan`).
    Seluruh jalur daun panduan bentuk dokumen (378 entri) sudah diputuskan (`TestPanduanBentukDokumenNolMedanBelumDiputuskan`);
    panduan itu basi (diagram R45), jadi AC 59 atas data nyata baru dapat dipastikan pada uji-kering
    pemuat (urutan F7 langkah 2, `MODUL.md` bab *Migrasi*) — tetap 🟡, penahannya kini F7/K11, bukan F3.

### Keamanan

60. `[terverifikasi]` Tidak ada nama orang tersalin ke berkas rancangan, spec, maupun test.
    *(Bab 9)*
61. `[terverifikasi]` Tidak ada nomor polis ditulis apa adanya di dalam berkas proyek. *(Bab 9)*
62. `[terverifikasi]` Test berjalan di atas data buatan, bukan cuplikan produksi. *(Bab 7)*
63. `[terverifikasi]` Perubahan skema memerlukan persetujuan manusia sebelum dijalankan. *(Bab 9)*

### ⭐ Tiga yang ditambahkan pada pemeriksaan silang 23 September 2026

> ⚠️ **Ditambahkan sesudah pemeriksaan silang terhadap ketiga puluh satu ketetapan brief.**
> `[penyimpangan sadar]` Tiga ketetapan **terwakili di Implementation Decisions tetapi belum punya
> AC**: ketetapan **16** *(ID-21)*, **21** *(ID-25)*, dan **31** *(ID-1)*. ⛔ **Tidak ada keputusan
> baru di sini** — ketiganya dinyatakan ulang dalam bentuk yang dapat diuji.

64. `[terverifikasi]` `T_GENERAL_POLIS` menampung **79 medan skalar tingkat atas** `PolicyTreatyIn`
    ditambah **tujuh kolom datar** dari `POOLDATA.json_polis`. Test yang menemukan salah satu medan
    itu tidak punya kolom, atau menemukan kolom kedelapan dari `json_polis`, **gagal**.
    ⚠️ Angka 79 adalah **batas bawah** — lihat Bab 1 peringatan dan butir `[terbuka]` 1.
    *(ID-21)*
65. `[keputusan work owner]` Bentuk gabungan `CEDING_CO_NAME` / `CEDING_CO` di `T_GENERAL_POLIS`
    **boleh tidak sinkron** dengan baris `T_POLIS_CEDING`, dan **tidak ada penjaga** yang
    memaksanya sinkron. Test yang **menolak** penyimpanan karena keduanya berbeda **gagal** —
    ketidaksinkronan itu konsekuensi yang diterima, bukan cacat. *(ID-25)*
66. `[keputusan work owner]` Arah ketergantungan `handlers → services → repository` tidak pernah
    dibalik: `repository` tidak memanggil `services`, dan `services` tidak memanggil `handlers`.
    Test yang menemukan pemanggilan ke arah sebaliknya **gagal**. *(ID-1, CLAUDE.md §5)*

---


### ⛔ RALAT dan pertentangan yang ditemukan saat implementasi — 2026-10-03

| AC | Bunyi lama (dikutip) | Temuan | Yang dibangun |
| ---: | --- | --- | --- |
| 15 | *"`IsApproved` bernilai `""` tersimpan sebagai `""`, bukan `NULL`"* | Oracle menyimpan `''` sebagai NULL | dibaca kembali `""`; setara di halaman, tidak di SQL |
| 15 ⛔ RALAT (P11, 04-10-2026) | baris di atas: *"dibaca kembali `""`; setara di halaman, tidak di SQL"* — temuan tanpa bunyi AC baru | Oracle `''` ≡ NULL | bunyi baru AC 15 ditulis di bab AC: `""` → NULL → `""`, `"0"` → `'0'` → `"0"`; keduanya tidak menyatu (`DecisionTable/isApproved.xml` kolom `text`); `repository/penyimpanan_db_test.go` `TestIsApprovedKosongDanNolTetapBerbeda` (K11) |
| 25 ⛔ RALAT (P11, 04-10-2026) | *"Pembandingan dua nilai uang memakai toleransi atau bentuk terbulatkan, bukan kesamaan persis"* | 176 rule terjangkau: nol pembandingan dua nilai uang yang hidup (rasio bertoleransi `<=0.01` hanya di langkah `//`); uang lawan nol eksak | bunyi baru AC 25 + RALAT ID-20; `models/pembandingan_uang_test.go` |
| 38 | *"`DEDUCTION1` `DEDUCTION2` `TOTAL_SHARE_PERCENTAGE_PREMIUM` `TOTAL_SHARE_PERCENTAGE_CLAIM` disimpan sebagai persen"* | `TOTAL_*` turunan baris spreading; penjaga repo melarang nama `TOTAL_` di migrasi | DEDUCTION1/2 persen; `TOTAL_*` dihitung saat dibaca (`HitungTotalSpreading`) |
| 38 ⛔ RALAT (P9, 04-10-2026) | baris di atas: *"DEDUCTION1/2 persen"* | `[keputusan work owner]` **K3** (rumus XML apa adanya) — bertentangan dengan kode sejak paket P2 | `DEDUCTION1/2` **uang** (`models/katalog.go` `kUang`, `NUMBER(38,8)`; RALAT AC 38 di bab AC): sel `.Deduction1/2` `pxCurrency` (`Section\DetailPolicyTreatyIn.xml`, `Section\DetailDeptHeadTreatyIn_UW.xml`), `CountNetPremi_act` langkah 4, `SetPPNPPH` langkah 4; `TOTAL_*` tetap turunan |
| 39 | *"`HISTORYAKSEPTASIPRODUCTION` tidak dibuat ulang"* | `SaveViewSuggest` hanya menulis bila `BusinessFac == "F"`; treaty "T" | ⚠️ `T_POLIS_SUGGEST` (baris dokumen SuggestList, bukan tabel riwayat) — keputusan agen, mohon konfirmasi WO |
| 39 ⛔ RALAT atas RALAT (putaran 2, 03-10-2026) | baris di atas: *"⚠️ `T_POLIS_SUGGEST` (baris dokumen SuggestList, bukan tabel riwayat) — keputusan agen, mohon konfirmasi WO"* | `[keputusan work owner]` **K4**: konfirmasi **ditolak** — tabel di luar diagram grilling (bab 0 butir 11) | ✅ AC 39 apa adanya: migrasi 328 dihapus; `SuggestList` ditulis ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION` yang sudah ada (15 kolom `InsertViewSuggest_SQL`) dan dibaca balik; syarat `BusinessFac == "F"` = `[penyimpangan sadar]` (lihat ID-31) |
| 21-22, 52-59 | format dan pemuat dokumen lama | pemuat (tiket 22) belum dibangun | `T_POLIS_MEDAN_LAIN` (329) sudah ada tanpa penulis |

⛔ **RALAT putaran 2 — 03-10-2026 — baris *21-22, 52-59* di atas.** Bunyi lama: *"pemuat (tiket 22) belum
dibangun"* / *"`T_POLIS_MEDAN_LAIN` (329) sudah ada tanpa penulis"*. Bunyi baru: pemuat dibangun
(`backend/alat/pemuatlama`, `services/pemuat.go`, `models/dokumenlama.go`, `models/laporanlama.go`,
`repository/lama.go`); penampung = berkas CSV (K17, RALAT ID-27 dan AC 57/59). AC 21, 22, 52, 53, 54, 56,
57, 58 diuji; AC 55 dan 59 🟡 (uji kolom bertag `db` belum dijalankan; jumlah atas data nyata belum
diketahui). `[penyimpangan sadar]` atas **ID-4**: pemecah dokumen lama tinggal di `models` sebagai fungsi
murni (seam 3 `spec.md` §6.2) supaya uji-kering pemuat tidak menyentuh tabel baru; penulisannya tetap
lewat antarmuka `repository` yang sama (ID-3). Cap waktu ` GMT` dibaca sebagai jam dinding
Asia/Jakarta (rule `GeneratePolicyNoTreaty_Act` langkah 5.3 membaca hari dalam Asia/Jakarta); tanggal
ambigu tidak ditebak (K15).


## 7 · Testing Decisions

### Apa yang membuat sebuah test baik di sini

Test menguji **perilaku dari luar**, bukan rincian di dalamnya. Ia tidak tahu nama fungsi pemecah
dokumen, tidak tahu urutan pemanggilan, tidak tahu bentuk struktur di dalam Go. Yang diuji: apa yang
masuk, apa yang tersimpan, apa yang keluar.

Test yang patah ketika rincian di dalam diubah **tanpa mengubah perilaku** adalah test yang buruk,
dan harus ditulis ulang, bukan diperbaiki.

### Seam — satu, di `repository`

Seluruh pengujian penyimpanan lewat antarmuka `repository`. Pemecah dokumen dan penyusunnya diuji
**lewat** seam itu, bukan langsung.

Alasannya: makin sedikit seam, makin sedikit yang harus dijaga sejajar. Dan pemecah dokumen tidak
punya nilai di luar penyimpanan — mengujinya tersendiri berarti mengunci bentuk yang seharusnya
bebas berubah.

### ⚠️ Pulang-pergi saja tidak cukup

Kalau tulis dan baca sama-sama salah secara **simetris**, uji pulang-pergi tetap hijau.

Contoh nyata dari korpus ini: `GroupPanel = "006"` disimpan sebagai bilangan `6`, lalu dibaca balik
menjadi `"6"`. Pulang-pergi lolos; penggolong `BusinessType_DeT` tetap gagal menemukan barisnya.

⇒ **Sebagian test memeriksa nilai kolom langsung**, bukan hasil bacanya. Masih seam yang sama —
hanya cara menegaskan yang berbeda.

| Cara menegaskan | Dipakai untuk |
| --- | --- |
| pulang-pergi | keutuhan 79 medan inti dan seluruh tabel anak |
| **nilai kolom langsung** | nol di depan · `""` lawan `"0"` lawan `NULL` · skala desimal · tipe kolom |
| pemeriksaan skema | ketiadaan kolom `LAYER*` di inti · ketiadaan tipe `float` · keunikan kunci |
| pemeriksaan perilaku | penolakan baris XOL pada polis prop · pembatalan transaksi |

### Modul yang diuji

| Modul | Yang dijaga |
| --- | --- |
| `repository` polis treaty inward | seluruh AC 1–54 |
| pemuat migrasi | AC 55–59, lewat seam yang sama |
| penggolong `BusinessType_DeT` | AC 14 — satu-satunya aturan dagang yang diuji di sini, sebab ia bergantung langsung pada tipe kolom |

### Data uji

Data **buatan**, tidak diambil dari produksi. Nomor polis dan nama pihak dibuat-buat.

⭐ Tetapi nilai yang **bermasalah** ditiru bentuknya: `592629512.880000276` dipakai apa adanya
sebagai nilai uji, sebab justru itu yang harus bertahan. Nilai seperti itu **bukan data pribadi** —
ia bentuk, bukan isi.

Sekurangnya tiga berkas contoh dokumen: proporsional biasa · proporsional dengan daftar angsuran
datar · non-proporsional dengan beberapa layer dan beberapa mata uang.

### Prior art

⛔ **Tidak dapat disebut.** Repo implementasi `D:\XML\nusantara-re\` **terlarang** menurut aturan
proyek, dan proyeknya belum memiliki kode Go. Test yang ada karena itu tidak dapat dijadikan
rujukan.

Yang dapat dirujuk hanya **preseden dokumen**: `.scratch\claim-life\spec-penyimpanan-relasional.md`
memakai bentuk spec yang sama untuk persoalan yang sama pada lini Life.

⇒ Begitu test pertama ditulis, berkas ini **diperbarui** untuk menyebutnya sebagai prior art bagi
ronde berikutnya.

---

## 8 · Out of Scope

| Yang dikeluarkan | Sebab |
| --- | --- |
| Pohon `LocationList → OccupationList → AnekaList → CoverageList → ClauseList/DeductibleList` | ⭐ milik **Fac In**. `[terverifikasi]` nol rujukan di NB Treaty In · `[keputusan work owner]` contoh dokumen yang memuatnya dipastikan milik Fac In |
| `M_TREATY_IN.JSONDATA` dan keluarganya | dokumen **kontrak**, bukan polis. `[terverifikasi]` dibaca 11 SQL di 8 modul, termasuk pemeriksaan batas limit oleh Claim Prop dan Claim Non Prop |
| Tiga tabel proyeksi selisih | milik ronde **EDM**. Di NB nilainya nol baris, sebab `OLD_POLIS_ID` selalu kosong |
| `TREATYINPRODUCTION` · `HISTORYAKSEPTASIPEGA` · `JSON_POLIS_MONITORING` · `TREATY_IN` · `ACHIEVEMENT` | sudah datar — **dipertahankan apa adanya**, tidak dirancang ulang |
| Penjaga duplikat `ACHIEVEMENT` berbasis nilai uang | `[penyimpangan sadar]` — dipatahkan galat presisi; kuncinya medan pengenal |
| Penghapusan ceding lewat `@replaceAll` pada teks gabungan | `[penyimpangan sadar]` — tidak ditiru |
| Perhitungan premi, komisi, pajak | bukan persoalan penyimpanan |
| Alur persetujuan dan antrean | sudah dispec di `spec.md` |
| DDL dan presisi fisik | dicocokkan DBA di dalam tiket |

### ⚠️ Peringatan yang wajib dibaca pemakai `TREATYINPRODUCTION`

`[terverifikasi]` Tabel itu dipertahankan apa adanya, dan membawa dua jebakan:

**Kolom `DEDUCTION1` menampung dua satuan.** Jalur proporsional mengisinya dari `.Deduction1` — sebuah
**persen**. Jalur XOL mengisinya dari `.Deduction` — sebuah **nilai uang**. Tidak ada penanda yang
membedakannya selain `PROPORTIONALTYPE` di baris yang sama.
> ⛔ **RALAT** 2026-10-03 (K3, paket P2) — bunyi lama: *"Jalur proporsional mengisinya dari
> `.Deduction1` — sebuah **persen**"* → bunyi baru: `.Deduction1` jalur proporsional adalah **jumlah
> uang** (XML: `pxCurrency`, dikurangkan dari premi di `CountNetPremi_act` langkah 4) — kedua jalur
> mengisi `DEDUCTION1` dengan uang (peringatan "persen dengan rupiah" di sini gugur untuk kolom ini).

⇒ Pembaca lewat SQL **wajib menyaring `PROPORTIONALTYPE` lebih dulu**, kalau tidak ia menjumlahkan
persen dengan rupiah.

**`LAYER*` pada polis proporsional bernilai `"0"`, bukan kosong** — jalur proporsional menulis `"0"`
ke sana.

Dan satu hal lagi: `[terverifikasi]` ada **empat titik sisip** di `InsetTreatyInProdAddendum_Act`,
dengan jalur XOL menyisip **per mata uang** ⇒ satu polis dapat menghasilkan beberapa baris.

---

## 9 · Butir `[terbuka]` — daftar penuh

⛔ **Spec ini `needs-info`.** Empat butir belum terjawab, dan **tidak boleh ditutup oleh pelaksana.**

| # | Butir | Pemilik | Yang tertahan olehnya |
| ---: | --- | --- | --- |
| 1 | ⭐ **`JSON_DATAGUIDE(DATA_JSON)` dari DBA** — satu-satunya yang dapat menutup selisih **74 lawan 79** medan antara contoh dokumen dan sapuan korpus | `[data DBA]` | kelengkapan daftar kolom · ID-27 · AC 59 |
| 2 | **Presisi fisik kolom uang** — skala pastinya | `[data DBA]` | DDL |
| 3 | ⚠️ **Treaty Out masuk lingkup atau tidak.** `[terverifikasi]` Folder `NB Treaty In` memuat `InsertToTreatyOutXOLList` · `InputPolicyTreatyOutDetail_preACT` · `DetailPolicyTreatyOutNonProportional` · `BrowseTreatyOut` yang membaca `M_treaty_out`. Sapuan berjangkar `PolicyTreatyIn` **tidak menampungnya** | `[work owner]` | lingkup seluruh spec — bila ikut, spec **diperluas**, bukan ditambal |
| 4 | **Migrasi dokumen lama** — dipindahkan seluruhnya, atau sebagian dibaca lewat jalur lama | `[work owner]` | rencana pemuatan · AC 55–59 |

Butir **3** yang paling menentukan: bila Treaty Out ternyata ikut, halaman `PolicyTreatyOut`, tabel
`M_TREATY_OUT` dan `FACOUTPRODUCTION` masuk lingkup, dan sensus medan harus diulang dengan jangkar
kedua.

---

## 10 · Further Notes

### Yang berubah dari sistem lama — tiga `[penyimpangan sadar]`

**1 · Penghapusan ceding.** Sistem lama memakai `@replaceAll` pada teks gabungan; sistem baru
menghapus barisnya.

> Bunyi lama, dikutip: `DeleteCeding_Act` menetapkan
> `Quotation.CedingCoName = @replaceAll(Quotation.CedingCoName, Param.cedingco, "")`.

**2 · Penjaga duplikat `ACHIEVEMENT`.** Sistem lama membandingkan seluruh 17 kolom termasuk enam
kolom uang; sistem baru memakai medan pengenal.

> Sebabnya: selisih `2,76 × 10⁻⁷` membuat baris yang seharusnya sama dinilai berbeda, lalu tersisip
> dua kali. Penjaganya gagal justru pada kasus yang paling perlu dijaga.

**3 · Keunikan `(NOPOLIS, PRODKE)`.** Sistem lama mengambil nomor generasi dengan `MAX + 1` tanpa
penguncian; sistem baru menolak yang kedua.

### Tiga tabel proyeksi yang menunggu ronde EDM

`T_POLIS_DIFFERENCE` dan dua anaknya tidak dispec di sini, tetapi **bentuk tabel ini sudah
menampungnya** — `OLD_POLIS_ID` sudah ada, dan `NOURUT` sudah menjadi kunci pasangannya. Ronde EDM
tidak perlu mengubah satu pun tabel yang dispec di berkas ini.

⭐ Keputusan yang sudah diambil dan berlaku di sana: rumus selisih **dihitung di Go**, bukan di
Oracle; nilainya **tidak pernah dihitung ulang** untuk data lama; dan kolom `SUMBER` membedakan
baris hasil migrasi dari baris hasil hitungan baru.

### Go tetap harus mampu membaca JSON

Walau dokumen **polis** ditinggalkan, dokumen **kontrak** `M_TREATY_IN.JSONDATA` tetap ada dan tetap
dibaca. ⇒ Kemampuan mengurai JSON tidak hilang dari sistem baru; yang hilang hanya kewajiban
**menulisnya** untuk data polis.

⚠️ Ini meralat kalimat yang sempat beredar di proyek ini — *"tidak ada lagi dokumen JSON"* — yang
ternyata terlalu jauh.

### Usulan yang belum diputuskan, tidak menahan spec ini

`[terbuka]` `PROC_GENERATE_SEQUENCE_NUMBER` memuat penyesuaian akhir tahun bertanggal mati
*(`02/01/2026`)*, yang menurut work owner dipakai bila tanggal closing bergeser. Di sistem baru tidak
ada procedure untuk disunting; kalau penyesuaian itu ditulis di Go, tiap akhir tahun butuh rilis
kode.

⇒ **Usul:** nilainya menjadi **data**, satu baris setelan, bukan baris program. Menunggu persetujuan.

### Tidak ada ADR baru

Spec ini tidak menuntut ADR baru. ADR-0003 *(uang bukan `float`)* dan ADR-0005 *(penanda lingkungan)*
keduanya tetap berlaku apa adanya.

### Empat jebakan yang sudah menjerat proyek ini

Diwariskan supaya tidak terulang pada ronde berikutnya:

1. **Menebak nama tag.** `Activity` memakai `PropertiesName`/`PropertiesValue` **tanpa** awalan `py`;
   `DataTransform` dan `Flow` memakai **dengan** awalan. Menyapu satu saja kehilangan yang lain.
2. **Sapuan berjangkar tidak melihat rujukan relatif.** Terbukti tiga kali.
3. **Mengambil nilai mayoritas korpus, bukan nilai di aturan yang bersangkutan.** Sekali ini
   menghasilkan laporan palsu bahwa dua kolom tertukar.
4. **`pyStepsPreCondParams` adalah percabangan lompat**, bukan gerbang hidup/mati.

---

## ⭐ KEPUTUSAN 23-09-2026 — lingkup pemindahan dokumen lama

`[keputusan work owner]` **Seluruh** dokumen polis dipindahkan — setiap polis, **setiap generasinya**. Tidak ada penyaringan menurut status, tahun buku, atau lini usaha.

Bunyi lama: ~~*"dipindahkan seluruhnya atau sebagian — belum diputuskan"*~~.

⭐ Alasannya: selisih endorsemen dihitung terhadap generasi sebelumnya. Bila hanya generasi terakhir dipindah, **pembandingnya hilang**.

⛔ `[terbuka]` Besaran datanya belum diketahui — empat angka diminta ke DBA *(cacah baris, tahun paling awal, ukuran total, nomor generasi tertinggi)*.

---

## ⭐ KEPUTUSAN 23-09-2026 sore — presisi uang dan daftar medan

### Presisi — **`NUMBER(38,8)`**

`[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."*

Bunyi lama: ~~*`NUMBER(20,8)` — dua belas digit di depan koma*~~ ⛔ **tidak cukup.**

⭐ **30 digit di depan koma, delapan di belakang** — batas tertinggi Oracle.

| Bukti korpus | Digit kiri |
| --- | ---: |
| `99999999999999999.99` sentinel "tanpa batas" | **17** |
| `181500000000.00` ambang kapasitas treaty | **12** |
| `150000000000.00` ambang kapasitas treaty | **12** |

⚠️ Ambang dagang nyata **tepat di batas** dua belas digit — nol ruang lega, padahal penjumlahan lintas mata uang dapat melewatinya.

⭐ `NUMBER` di Oracle berpanjang **berubah-ubah**: hanya digit bermakna yang tersimpan. Pelebaran ini **tidak memakan ruang tambahan**, sehingga tidak ada alasan menahannya di angka yang pas-pasan.

⛔ Butir `[data DBA]` presisi **ditutup** — oleh bukti, bukan oleh jawaban DBA.

### Daftar medan — dari Activity dan Section

`[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."*

⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, **232** di antaranya tampil di layar. Hasilnya: **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**.

Data guide turun derajat menjadi **penambal**. Ia sah hanya untuk **panjang maksimum per medan**, dan untuk 24 nama yang tidak tersapu aturan — `OldData` · `ProdKe` · `EDMNo` · `OldPolicyNo` · `BranchCode` dan lainnya, yang ditulis sistem, bukan diketik pengguna.

⛔ Butir `[data DBA]` panduan bentuk dokumen **ditutup**.

⚠️ **Cacah 394 tetap batas bawah, bukan total.** Penampung medan tak dikenal tetap wajib, dan wajib **kosong** sebelum pekerjaan dinyatakan selesai.

---

## TELEMETRI EKSEKUSI

⛔ **Pengukuran dari luar TIDAK dilakukan** untuk ronde ini. Spec ditulis di dalam sesi yang sama
dengan seluruh ronde rancangan, sehingga tidak ada pemanggilan `claude --print --output-format json`
yang dapat memberi angka terukur.

Yang dapat dihitung dari dalam sesi, dan **hanya ini yang bukan taksiran**:

| Ukuran | Nilai | Cara |
| --- | ---: | --- |
| berkas korpus disapu | **9.430** | `os.walk` seluruh `D:\XML\RNM_BRD` kecuali `OUTPUT_HASIL_RNM` |
| jalur properti unik | **322** | pola berjangkar `PolicyTreatyIn` di dua modul |
| SQL `pyBrowseSQL` diperiksa | 36 EDM + 41 NB | folder `RDBList` |
| stored procedure dibaca | **6** | naskah dari work owner |
| skrip sapuan dijalankan | **92** | catatan sesi |

⚠️ **Taksiran, bukan ukuran** — dan dinyatakan begitu supaya tidak dikira terukur:

| Ukuran | Taksiran |
| --- | ---: |
| token masuk | ± 985–1.080 ribu |
| token keluar | ± 184–205 ribu |

Cara mengukur yang sebenarnya, untuk ronde berikutnya:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Keluarannya memuat `total_cost_usd`, `usage`, `duration_ms`, dan `num_turns`.

---

*Spec penyimpanan relasional NB Treaty In · 23 September 2026 · `needs-info` karena empat butir
terbuka pada Bab 9.*

---

### ⭐ Pemeriksaan silang 23 September 2026 — sesi kedua

⚠️ **Berkas ini ditulis satu sesi, lalu diperiksa sesi lain terhadap kesembilan kriteria
keberhasilan brief.** Yang ditemukan dan dikerjakan:

| Kriteria brief | Hasil |
| --- | --- |
| sembilan bab lengkap | ✅ **sepuluh** bab bernomor + `TELEMETRI EKSEKUSI` |
| seluruh AC berpenanda | ✅ **0** tanpa penanda |
| seluruh AC berujuk | ✅ **0** tanpa rujukan — ⭐ memakai skema `ID-1..ID-35` sendiri, **lebih kuat** daripada nomor bab yang diminta brief |
| ⛔ **31 ketetapan Bab 4 terwakili di AC** | ⛔ **tiga belum** — ketetapan **16**, **21**, **31**. ⭐ **Diperbaiki:** AC **64**, **65**, **66** ditambahkan |
| NB lawan EDM dinyatakan terang-terangan | ✅ `NOURUT` dinomori ulang di NB · satu tingkat angsuran pada proporsional · tabel yang sama untuk keduanya |
| peringatan `DEDUCTION1` dua satuan dan `LAYER* = "0"` | ✅ keduanya tertulis |
| empat butir `[terbuka]` dibawa utuh | ✅ **4** terdaftar, **0** ditutup |
| blok ringkasan dua cara, jendela disebutkan | ✅ — dan ⭐ **cara A yang keliru dikutip, tidak dihapus** |
| nol `CREATE TABLE` · nol ADR baru · nol nama orang | ✅ **0** · **0** · **0** |

⚠️ **Dua kekeliruan alat pemeriksa itu sendiri, dicatat:**
**(1)** Ia mula-mula melaporkan **57 dari 63 AC tanpa rujukan bab** — keliru; AC memakai bentuk
`*(ID-n)*`, dan polanya hanya mencari `*(Bab n)*`.
**(2)** Ia mula-mula melaporkan **delapan** ketetapan tanpa AC; sesudah diperiksa satu per satu,
yang benar-benar kurang **tiga**. Lima sisanya tertangkap dengan kata yang berbeda — mis. ketetapan
9 tertulis *"sembilan angka di belakang koma"*, bukan *"skala 9 desimal"*.

**Ongkos pemeriksaan silang ini** — hasil pengurangan terhadap baseline, ⛔ **bukan** pengukuran
dari luar:

| Ukuran | Nilai |
| --- | ---: |
| panggilan model | **28** |
| token keluar | **47.183** |
| token cache ditulis | **124.721** |
| token cache dibaca | **18.209.300** |
| AC ditambahkan | **3** |
| bab lain disunting | **0** |
