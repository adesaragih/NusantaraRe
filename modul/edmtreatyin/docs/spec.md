# Spesifikasi — EDM Treaty In

## Migrasi Treaty Inward — Endorsemen · modul *EDM Treaty In*

---

## 1 · Cara membaca berkas ini

### 1.1 Urutan kewenangan

Bila dua sumber bertentangan, yang di atas menang. ⭐ **Pertentangan yang ditemukan ditulis, bukan
didiamkan** — lihat Bab 10.1.

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `PERTANYAAN-RONDE-1.md` — **P50–P60**, terjawab 11/11 | `[keputusan work owner]`, mengikat |
| 2 | `..\nb-treaty-in\PERTANYAAN-untuk-*.md` — **P1–P49** | `[keputusan work owner]`, berlaku lintas modul |
| 3 | `KEADAAN-EDM-TREATY-IN.md` | keadaan terukur, diuji dua cara |
| 4 | `..\nb-treaty-in\KEADAAN-NB-TREATY-IN.md` Bab 4 | dua belas ketetapan; berlaku **kecuali terbukti sebaliknya** |
| 5 | `grilling-ronde-1.md` | ⛔ latar saja — **empat pernyataannya terbukti keliru** |
| 6 | korpus XML | selalu boleh dipakai membuktikan ulang |

⛔ **Nol angka disalin dari sumber 5.** Dua angkanya — **378** langkah dan **1.062** pasangan —
**ditarik** karena tidak dapat direproduksi oleh enam rekonstruksi.

### 1.2 Penanda

| Penanda | Arti |
| --- | --- |
| `[terverifikasi]` | ada **jalur berkas + tipe rule + nama rule**, dan perintah ujinya |
| `[keputusan work owner]` | diputuskan yang berwenang — ⛔ **tidak diubah spec ini** |
| `[penyimpangan sadar]` | ⚠️ **sengaja berbeda dari Pega**, alasannya tertulis |
| `[terbuka]` | belum terjawab — ⛔ **tidak ditutup spec ini** |
| `[data DBA]` | hanya dapat dijawab dari basis data |
| `[dugaan]` | ⛔ **tidak dipakai sebagai dasar keputusan mana pun** |

### 1.3 ⭐ Sensus berkas ini — dihitung DUA CARA

> ⭐ **Jendela: seluruh berkas ini MINUS blok sensus 1.3 sendiri** — **909 baris** ⚠️ *(RALAT — semula ditulis **885**, diukur sebelum Bab 11.5 dan blok telemetri ditambahkan; angka lama tidak dihapus)*.
>
> | Yang dicacah | Cara A | Cara B | |
> | --- | ---: | ---: | :---: |
> | User story | **42** | **42** | ✅ |
> | Acceptance Criteria | **59** | **59** | ✅ |
> | ⛔ AC **tanpa penanda** | **0** | **0** | ✅ |
> | ⛔ AC **tanpa rujukan bab** | **0** | **0** | ✅ |
> | butir `[terbuka]` aktif | **13** | **13** | ✅ |
> | butir `[penyimpangan sadar]` | **5** | **5** | ✅ |
>
> ⭐ **Cara A** — cacah pola penomoran di dalam babnya *(`^N. **Sebagai**` Bab 4; `` ^N. `[ `` Bab 7;
> baris tabel Bab 9 dan 10.3)*.
> ⭐ **Cara B** — periksa **keberurutan nomornya**: user story **1–42 tanpa lompat**,
> AC **1–59 tanpa lompat**; penyimpangan bernomor **1–5**.
>
> **Sebaran penanda, kemunculan di jendela:** `[terverifikasi]` **46** · `[keputusan work owner]` **38** · `[data DBA]` **21** · `[penyimpangan sadar]` **11** · `[terbuka]` **5** · `[dugaan]` **1**
>
> ⚠️ **Kemunculan penanda ≠ jumlah butir** — satu butir dapat membawa dua penanda.

### 1.4 ⛔ Yang TIDAK dilakukan spec ini

- ⛔ **Tidak menutup satu pun butir `[terbuka]`.** Penutupan milik work owner, DBA,
  Product & Underwriting, Finance, atau IAM.
- ⛔ **Tidak membuat ADR baru.** Lima belas ADR yang ada tetap berlaku.
- ⛔ **Tidak merancang tabel penyimpanan.** Itu menunggu sensus properti — tiket `00` NB.
- ⛔ **Tidak memilih** di antara tiga pilihan galat uang *(Bab 8.9)*.

---

## 2 · Problem Statement

Hari ini, mengubah kontrak treaty inward yang sudah berjalan dikerjakan di dalam sistem Pega yang
akan dihentikan. Endorsemen — perubahan atas polis yang sudah terbit — adalah **jenis kasus
tersendiri** di sana, dengan kelas kerjanya sendiri, tangganya sendiri, dan perhitungan selisihnya
sendiri.

Yang membuat pemindahannya sulit bukan besarnya, melainkan **empat sifat yang tidak terlihat dari
luar**:

1. ⛔ **Bentuk data tidak seragam.** Dokumen yang menyimpan polis berubah bentuk menurut jenis
   treaty. Tiga contoh produksi memberi **37**, **64**, dan **47** medan skalar tingkat atas; hanya
   **23** ada di ketiganya, sedangkan gabungannya **74**. `[data DBA]`
2. ⛔ **Satu nama daftar, dua kedalaman.** `ListInstallment` **bersarang** pada satu bentuk dan
   **datar** pada bentuk lain — keduanya nyata di korpus. Pengurai yang menganggapnya seragam akan
   patah. `[terverifikasi]`
3. ⛔⛔ **Angka uang di produksi sudah membawa galat.** Ekor sebesar **2,76 × 10⁻⁷** terbukti
   tersimpan permanen. EDM **menghitung selisih**, dan selisih dua angka bergalat membawa **jumlah**
   galatnya. `[data DBA]` `[terverifikasi]`
4. ⛔ **Endorsemen berlapis menumpuknya.** Satu polis dapat diendorse berkali-kali, dan tiap
   endorsemen menghitung selisih terhadap keadaan **tepat sebelumnya**. `[keputusan work owner]`

Dan satu sifat yang menguntungkan, yang wajib dinyatakan supaya tidak disalahpahami:

⭐⭐ **Lingkup EDM tidak menyusut.** Folder `NB Treaty In` adalah *dependency closure* — ia memuat
aturan milik modul lain, sehingga 28 aturan dan 561 langkah dibuang sebagai tak terjangkau. Folder
`EDM Treaty In` **bukan**: **66 dari 66** `Activity` terjangkau, **380 dari 380** langkah
`Property-Set` terjangkau. `[terverifikasi]`

⇒ ⛔ **Jangan ada yang mengharapkan penyusutan 60 % seperti NB. Ke-380 langkah dikerjakan
seluruhnya.**

---

## 3 · Solution

Endorsemen treaty inward dibangun ulang di Go + React + Oracle sebagai **jenis kasus tersendiri**,
sejajar dengan polis baru, bukan sebagai tahap di dalamnya.

| Yang dibangun | Bentuknya di sistem baru |
| --- | --- |
| **Jenis kasus** | kasus endorsemen berdiri sendiri, mengikuti kelas kerja lama `…Work-EndorsementTreaty` |
| **Tangga persetujuan** | tiga jenjang — Admin → Sec Head → Dept Head — sejalan **P13** |
| **Pembekuan data lama** | nilai polis induk **disalin dan beku** saat berkas dibuat |
| **Perhitungan selisih** | selisih terhadap keadaan **tepat sebelumnya**, bukan polis asli |
| **Penomoran** | `NoPolis + "/E" + ProdKe` dua digit |
| **Pembatalan** | **jenis endorsemen**, dipilih di awal — bukan alur terpisah |
| **Konversi ke produksi** | satu panggilan keluar, dengan penanganan gagal yang menghapus bersyarat |
| **Penyimpanan** | tabel relasional; ⛔ bentuk tabelnya menunggu sensus properti |

⭐ **Yang membuat ini mungkin sekarang, dan tidak mungkin dua minggu lalu:** isi 380 langkah
penetapan nilai **terbaca penuh** dari ekspor — **1.309 pasangan nama=nilai**, **379 dari 380**
langkah membawa isinya. Rumusnya tidak perlu ditebak. `[terverifikasi]`

---

## 4 · User Stories

### Membuat endorsemen

1. Sebagai **admin reasuransi**, saya ingin membuat berkas endorsemen dari sebuah polis yang sudah
   terbit, supaya perubahan kontrak tercatat sebagai dokumen tersendiri, bukan menimpa polis asal.
2. Sebagai **admin reasuransi**, saya ingin nilai polis induk tersalin otomatis ke berkas endorsemen
   saat berkas dibuat, supaya saya tidak mengetik ulang puluhan medan.
3. Sebagai **admin reasuransi**, saya ingin nilai yang tersalin itu **beku**, supaya perubahan pada
   polis induk sesudahnya tidak diam-diam mengubah selisih yang sudah saya hitung.
4. Sebagai **admin reasuransi**, saya ingin memilih **jenis endorsemen** di awal — termasuk
   pembatalan — supaya saya tidak perlu mencari alur yang berbeda untuk membatalkan polis.
5. Sebagai **admin reasuransi**, saya ingin nomor endorsemen terbentuk otomatis dari nomor polis
   induk, supaya nomornya dapat dilacak kembali ke polis asalnya tanpa tabel bantu.
6. Sebagai **admin reasuransi**, saya ingin endorsemen kedua atas polis yang sama mendapat nomor
   urut berikutnya, supaya urutan perubahan terbaca dari nomornya saja.
7. Sebagai **underwriter**, saya ingin melihat nilai lama dan nilai baru berdampingan, supaya saya
   dapat menilai perubahannya tanpa membuka dua layar.
8. Sebagai **underwriter**, saya ingin selisihnya dihitung sistem, supaya angka yang saya setujui
   bukan hasil hitungan tangan.
9. Sebagai **underwriter**, saya ingin selisih dihitung terhadap **keadaan tepat sebelumnya**,
   supaya endorsemen ketiga tidak mengulang selisih yang sudah diakui pada endorsemen kedua.

### Tangga persetujuan

10. Sebagai **Sec Head**, saya ingin berkas endorsemen masuk ke antrean saya sesudah admin selesai,
    supaya saya tahu apa yang menunggu keputusan saya.
11. Sebagai **Dept Head**, saya ingin dapat menyetujui atau menolak, supaya kewenangan jenjang
    ketiga terpakai sebagaimana mestinya.
12. Sebagai **Dept Head**, saya ingin dapat mengisi medan keputusan saya sendiri, supaya saya
    bukan sekadar pengamat.
13. Sebagai **pemegang jenjang mana pun**, saya ingin penolakan saya mengembalikan berkas ke admin,
    bukan menutupnya, supaya perbaikan masih mungkin.
14. Sebagai **admin reasuransi**, saya ingin penolakan saya sendiri menutup berkas sebagai ditolak,
    supaya berkas yang salah sejak awal tidak menyumbat antrean atasan.
15. Sebagai **auditor internal**, saya ingin setiap perpindahan jenjang tercatat beserta perannya,
    supaya jejak persetujuan dapat ditelusuri tanpa bertanya kepada orangnya.
16. Sebagai **pemilik sistem**, saya ingin wewenang ditentukan **peran**, bukan nama orang, supaya
    kepindahan jabatan tidak mematahkan alur secara diam-diam.

### Perhitungan uang

17. Sebagai **underwriter**, saya ingin premi, komisi, dan potongan dihitung dengan rumus yang sama
    seperti sistem lama, supaya angka yang keluar dapat dibandingkan dengan riwayat.
18. Sebagai **Finance**, saya ingin potongan pajak brokerage dihitung `Deduction / 1,022` ketika
    jenis pajaknya `Inclusive`, supaya perlakuan pajaknya tidak berubah diam-diam.
19. Sebagai **Finance**, saya ingin nilai uang **tidak pernah** disimpan sebagai bilangan
    mengambang, supaya sistem baru tidak menambah galat baru.
20. Sebagai **Finance**, saya ingin persentase dibedakan dari uang, supaya `12.5` tidak pernah
    terbaca sebagai dua belas setengah rupiah.
21. Sebagai **Finance**, saya ingin mengetahui bahwa angka historis membawa galat kecil sebelum
    memutuskan cara memperlakukannya, supaya keputusannya sadar, bukan kejutan.
22. Sebagai **underwriter**, saya ingin penyebaran risiko dihitung ulang ketika endorsemen mengubah
    nilai, supaya bagian tiap penanggung tetap benar.
23. Sebagai **underwriter**, saya ingin kasus satu baris penyebaran dengan persentase kosong tetap
    menghasilkan angka yang benar, supaya bentuk endorsemen yang paling umum tidak jadi kasus gagal.

### Bentuk data

24. Sebagai **pengembang sistem baru**, saya ingin pembacaan dokumen lama menangani **kedua** bentuk
    daftar angsuran, supaya polis non-proporsional tidak gagal dibaca.
25. Sebagai **pengembang sistem baru**, saya ingin jenis treaty terbaca dari medan penentunya,
    supaya bentuk dokumen dapat diketahui sebelum diurai, bukan sesudah gagal.
26. Sebagai **pengembang sistem baru**, saya ingin medan yang tidak ada pada suatu bentuk dianggap
    **kosong**, bukan galat, supaya 51 dari 74 medan yang tidak universal tidak mematahkan pembacaan.
27. Sebagai **pemilik data**, saya ingin data dibaca dari **view relasional** yang sudah ada, bukan
    dari dokumen teks, supaya bentuknya dapat dinyatakan.

### Konversi ke produksi

28. Sebagai **admin reasuransi**, saya ingin endorsemen yang disetujui terkirim ke sistem produksi,
    supaya pencatatannya tidak perlu diulang di sana.
29. Sebagai **admin reasuransi**, saya ingin tahu **sebelum** data terhapus bahwa kegagalan konversi
    akan menghapus data produksi yang sudah tersisip, supaya saya tidak kehilangan pekerjaan tanpa
    peringatan.
30. Sebagai **pemilik sistem**, saya ingin penghapusan itu **tidak menyentuh** data yang sudah
    berhasil dikonversi, supaya kegagalan berulang tidak merusak konversi yang sah.
31. Sebagai **pemilik sistem**, saya ingin penghapusan hanya lewat **satu jalur**, supaya tidak ada
    tempat kedua di dalam kode yang dapat menghapus data produksi.
32. Sebagai **admin reasuransi**, saya ingin melihat pesan galat yang dapat dibaca ketika konversi
    gagal, supaya saya tahu harus mengulang atau menghubungi IT.
33. Sebagai **pengelola lingkungan**, saya ingin alamat sistem tujuan berupa **variabel
    lingkungan**, supaya lingkungan uji tidak pernah mengirim ke produksi.
34. Sebagai **pengelola lingkungan**, saya ingin lingkungan uji memakai basis data terpisah, supaya
    percobaan tidak menyentuh data sungguhan.

### Penyimpanan dan riwayat

35. Sebagai **pemilik sistem**, saya ingin seluruh urutan penyimpanan dibungkus **satu transaksi**,
    supaya kegagalan di tengah tidak meninggalkan separuh data permanen.
36. Sebagai **auditor internal**, saya ingin riwayat akseptasi tercatat dengan identitas pengguna
    yang sah, supaya kolom pelaku tidak kosong seperti di sistem lama.
37. Sebagai **pemilik data**, saya ingin setiap query menulis skema basis data **eksplisit**, supaya
    tidak ada tabel yang terbaca dari skema yang salah.
38. Sebagai **Finance**, saya ingin tanggal tutup buku berlaku **global** dan perubahannya terasa
    seragam, supaya periode penomoran tidak berbeda antar berkas.
39. Sebagai **pengembang sistem baru**, saya ingin nomor urut dihasilkan di dalam satu transaksi
    dengan penyimpanannya, supaya nomor tidak terbakar ketika penyimpanan gagal.

### Pemantauan dan pencapaian

40. Sebagai **pemilik sistem**, saya ingin pencatatan pemantauan bersifat **idempoten**, supaya
    percobaan ulang tidak menggandakan barisnya.
41. Sebagai **Finance**, saya ingin pencatatan pencapaian tidak tergandakan oleh selisih pecahan
    yang tak berarti, supaya angka kuartalan tidak terhitung dua kali.
42. Sebagai **pengembang sistem baru**, saya ingin pesan galat dari lapisan penyimpanan **tidak
    membawa markah tampilan**, supaya penyajiannya tetap milik lapisan tampilan.

---

## 5 · Implementation Decisions

### 5.1 Bentuk kasus dan alur

⭐ **Endorsemen adalah jenis kasus tersendiri.** `[terverifikasi]`
`Flow\InputAddendumTreatyIn` berkelas `ASM-FW-GISFW-Work-EndorsementTreaty` — kelas kerja yang
**berbeda** dari kasus polis baru. Alurnya memuat **8** kotak keputusan, **4** penugasan,
**2** utilitas, dan **22** sambungan.

⭐ **Tiga antrean, sejalan P13:** `ReasTreatyInAdmin` · `ReasTreatyInSecHead` ·
`ReasTreatyInDeptHead`. ⛔ Tidak ada jenjang keempat; tidak ada penyimpangan dari tangga NB.

**Modul yang dibangun:** layanan kasus endorsemen · layanan perhitungan selisih · layanan
penyebaran risiko · layanan penomoran · penghubung konversi produksi · repository penyimpanan.
Arah ketergantungan tetap `handlers → services → repository` *(CLAUDE.md §5)*.

### 5.2 Pembekuan data lama — `[keputusan work owner]` P58

Nilai polis induk **disalin** ke berkas endorsemen saat berkas dibuat, lalu **tidak pernah dibaca
ulang**. `[terverifikasi]` Jejaknya di `Activity\CreateEDMT` — 20 langkah, lima yang menentukan:
langkah 7 `RDB-List` membaca polis induk; langkah 8 memetakan kolom; langkah 10 `Page-Copy`
menyalin **seluruh halaman** ke dalam `OldData`; langkah 11 menyalin data kuotasi; langkah 13
`Page-Remove` **membuang halaman sumbernya**.

⇒ Sesudah langkah 13 tidak ada lagi sambungan ke polis induk. Perubahan pada polis induk sesudahnya
**tidak berpengaruh**.

⚠️ **Konsekuensi yang dibawa ke sistem baru:** penyalinan dilakukan **satu kali, di titik
pembuatan**, dan snapshot itu bagian dari berkas endorsemen — bukan rujukan yang dibaca saat
ditampilkan.

### 5.3 Perhitungan selisih — `[keputusan work owner]` P57

Selisih dihitung terhadap **keadaan tepat sebelumnya**, bukan terhadap polis asli.
`[terverifikasi]` Pemilihnya medan `OldData.EDMNo`:

| Keadaan | Layar |
| --- | --- |
| `.OldData.EDMNo = ''` — data lama **belum pernah** diendorse | `Section\DetailPolicyTreatyInPropOldData` |
| `.OldData.EDMNo != ''` — data lama **sudah** berupa endorsemen | `Section\DetailPolicyTreatyInPropOldData2` |

`Activity\EDMTCalculateTreatyDifference` memakai syarat yang sama untuk memilih cara menghitung.

⭐ **Akhiran `2` berarti dua KEADAAN, bukan dua tingkat tampilan.**

⚠️ **Pengujinya medan, bukan keberadaan simpul.** `OldData` **ada bahkan pada polis baru** —
`Activity\TreatyRealizationCheckXOLList` di NB menuliskannya saat realisasi. Uji berbasis keberadaan
akan selalu benar dan tidak berguna. `[terverifikasi]` `[data DBA]`

### 5.4 Penomoran endorsemen — `[keputusan work owner]` P55

```
EDMNo = OldData.PolicyNo + "/E" + @If(ProdKe < 10, "0" + ProdKe, ProdKe)
ProdKe = @toInt(<nomor urut terakhir dari basis data>) + 1
```

Ditetapkan `Activity\SetEDMTNoPolis` langkah 3, dipanggil `Activity\CreateEDMT`. `[terverifikasi]`
Bentuknya: nomor polis induk, `/E`, lalu **dua digit**. Batas teknis `ProdKe` **99**.

⛔⛔ **Pertentangan yang ditemukan — lihat Bab 10.1.** Korpus memuat **jalur penomoran kedua** yang
menghasilkan bentuk berbeda. **P55 menang** sebagai keputusan; jalur kedua **tidak dimigrasi** dan
masuk Out of Scope, ⛔ **dengan catatan bahwa pertentangannya belum dijelaskan work owner**.

### 5.5 Pembatalan — `[keputusan work owner]` P56

Pembatalan adalah **jenis endorsemen**, dipilih di awal saat berkas dibuat. `[terverifikasi]`
`Activity\CreateEDMT` langkah 16 memanggil `SetEDMTCancel` dengan catatan *"when Param.edmtype = 4,
run this activity"*. ⛔ **Bukan alur terpisah**, bukan tombol tersendiri.

### 5.6 Penyebaran risiko — `[keputusan work owner]` P60

`Activity\CountSpreading_Act` EDM punya **8** langkah bermetode lawan **7** di NB.
`[terverifikasi]` Langkah tambahannya nomor **4**:

```
.SpreadingRiskList(1).SplitRNMSharePct = 100      ⟵ "JIKA SPREADINGLIST 1 DAN SplitRNMSharePct NULL || 0"
```

Dan rumus di bawahnya **berbeda antar modul**:

| | Bila `SharePercentage` kosong | Presisi |
| --- | --- | ---: |
| **NB** | `100 / @LengthOfPageList(SpreadingRiskList)` | **10** |
| **EDM** | `@divide(.SplitRNMSharePct, Primary.RNMShare, 20)` | **20** |

⭐ **Sebab EDM memerlukan nilai bawaan dan NB tidak:** aturan yang **mengisi** `SplitRNMSharePct` —
`Activity\TreatyNonPropSetSpreading` — **ada di NB dan tidak ada di EDM**. `[terverifikasi]` Di
polis baru daftar dibangun dari nol dan medannya diisi di muka; di endorsemen daftar itu datang dari
polis yang diendorse dan medannya dapat tiba kosong — sementara rumus EDM **membaginya**.

⚠️ Ketidakseragaman presisi **10 lawan 20** ditiru apa adanya, sejalan ketetapan NB 12 dan **P46**.
⛔ Ia **bukan** keputusan baru; ia perilaku lama yang disalin sadar.

### 5.7 Konversi ke produksi — `[keputusan work owner]` P50, P51, P52, P53

**Satu panggilan keluar.** `[terverifikasi]` EDM **tidak punya FacOut**; cabang pengiriman kedua
**tidak pernah berjalan** dan **tidak dibangun**.

**Pemeriksaan keberhasilan** memeriksa **kedua** penanda, dalam tiga cabang:
`FacIn = 1` **dan** *(`FacOut = ""` **atau** `FacOut = 1`)*. `[terverifikasi]`
⚠️ Ronde 1 menyatakan hanya FacIn diperiksa — **keliru**, ia membaca keterangan manusia, bukan
bentuk yang dijalankan. Ketetapan NB 10 berlaku: **yang dijalankan benar**.

**Penanganan gagal:** data produksi yang sudah tersisip **dihapus** — ⭐ **bersyarat**, hanya bila
`STS_KONVERSI` belum bernilai `1`. Lingkupnya **satu kasus** *(`IDPEGA`)*, pada tiga tabel untuk
treaty. `[data DBA]`

⛔ **Penghapusan hanya lewat satu jalur** — satu fungsi repository, tidak dipanggil dari tempat lain.
⭐ **Pengguna diberi tahu lebih dulu** bahwa kegagalan akan menghapus.

⚠️ `[penyimpangan sadar]` **1** — logika penghapusan **ditulis ulang di Go**, bukan memanggil
stored procedure; mandat proyek melarang pemanggilan procedure. Yang dipertahankan **aturannya**,
termasuk penjaga `STS_KONVERSI`.

**Langkah 10** `Page-Set-Messages` **dihidupkan** — di sistem lama ia mati, sehingga pengguna tidak
pernah melihat pesan galat konversi. **Langkah 18** `Call ASMForceCaseClose` **tetap mati**.
`[keputusan work owner]` P52.

**Autentikasi:** ⛔ **tidak ada dari sisi kita** — jalur internal; autentikasi milik penerima.
`[keputusan work owner]` P53. ⭐ **Alamat tujuan variabel lingkungan, tidak pernah harfiah.**

### 5.8 Sumber data

| Hal | Ketetapan |
| --- | --- |
| pembacaan kontrak | **view relasional** yang sudah ada — 39 kolom, kecukupan sudah diuji *(P29)* |
| `RDBList` | **36 dari 36** membawa naskah SQL — nol yang harus ditebak `[terverifikasi]` |
| skema | ⭐ setiap query menulis skema **eksplisit** *(ketetapan NB 6)* |
| tanggal tutup buku | dibaca `WHERE ROWNUM = 1` — ⭐ **satu baris berlaku global** `[data DBA]` P54 |
| lingkungan | uji dan produksi memakai **basis data berbeda** `[keputusan work owner]` P59 |

⚠️ **Sistem lama tidak menaati ketetapan 6.** `[terverifikasi]` Dari 31 pernyataan `FROM`, hanya
**16** memakai skema `POOLDATA.`; `JSON_POLIS` muncul **dua ejaan**; dan `DATAPEGA` adalah **skema
kedua** yang disentuh modul ini. ⭐ **Ketetapannya tetap berlaku untuk sistem baru** — yang berubah
hanya catatan bahwa sistem lama tidak menaatinya.

### 5.9 Enam stored procedure — naskahnya kini lengkap

| # | Procedure | Tabel yang disentuh | `COMMIT` di dalam |
| ---: | --- | --- | :---: |
| 1 | `PEGA_TREATY_IN` | dua tabel | ya |
| 2 | `PEGA_JSON_POLIS_TREATYIN` | satu tabel, 8 kolom | ya |
| 3 | `PROC_GENERATE_SEQUENCE_NUMBER` | deret nomor | ⛔ **tidak** |
| 4 | `PEGA_DELETE_ERROR_KONVERSI` | tiga/empat tabel | ya |
| 5 | `INSERTJSONPOLISMONITORING` | satu tabel pemantauan | ya |
| 6 | `InsertUpdateAchievment` | satu tabel pencapaian, 18 kolom | ⛔ **tidak** |

⚠️ `[penyimpangan sadar]` **2** — **seluruh urutan penyimpanan dibungkus satu transaksi**, lebih
ketat dari sistem lama yang meng-`COMMIT` di empat tempat berbeda. *(Ketetapan NB 7, P2.)*

**Empat hal dari kedua naskah yang baru diterima:**

1. ⛔ **`INSERTJSONPOLISMONITORING` bukan "insert-or-update".** `[data DBA]` Bila baris ber-`IDPEGA`
   sama sudah ada, ia **berhenti tanpa berbuat apa-apa** — cabang bertanda *"Update"* isinya hanya
   `RETURN`. ⭐ **Idempoten, ya; memperbarui, tidak.** Kolom nomor produksi **selalu diisi `0`**,
   bukan dari parameter.
   ⇒ Sistem baru **mempertahankan sifat idempotennya** dan **menamainya jujur** — ia penyisip
   sekali-jalan, bukan pembaru.
2. ⛔ **Pesan galatnya membawa markah HTML** dan berubah bunyi menurut penanda lingkungan.
   `[data DBA]` ⚠️ `[penyimpangan sadar]` **3** — di sistem baru **pesan galat tidak membawa markah**;
   penyajian milik lapisan tampilan, dan penanda lingkungan tidak merembes ke lapisan penyimpanan.
   Menyentuh **ADR-0005**; ⛔ **tidak ada ADR baru.**
3. ⭐ **`InsertUpdateAchievment` hanya menyisip** — tidak ada jalur pembaruan sama sekali.
   ⭐⭐ **Parameter uangnya bertipe `NUMBER`**, satu-satunya dari enam yang menerima uang sebagai
   bilangan; lima lainnya menerimanya sebagai teks atau di dalam dokumen. `[data DBA]`
   ⇒ **Di sinilah presisi ditetapkan saat melintasi batas.** Sistem baru menyeberangkan uang sebagai
   desimal berpresisi tetap, **tidak pernah `float`** *(ADR-0003, CLAUDE.md §7)*.
4. ⛔⛔ **Penjaga anti-duplikatnya dipatahkan galat uang.** `[data DBA]` `[terverifikasi]`
   Pemeriksaannya membandingkan **seluruh 17 kolom, termasuk enam kolom uang**. Selisih
   **2,76 × 10⁻⁷** — yang **terbukti ada di data produksi** — membuat baris yang seharusnya sama
   dinilai **berbeda**, sehingga baris kedua tersisip. ⛔ **Penjaganya gagal justru pada kasus yang
   paling perlu dijaga.**
   ⇒ ⚠️ `[penyimpangan sadar]` **4** — sistem baru **tidak meniru** pemeriksaan duplikat berbasis
   nilai uang; kuncinya **medan pengenal**. ⛔ **Bentuk kunci itu menunggu keputusan work owner** —
   Bab 9.

### 5.10 Bentuk dokumen tidak seragam

**Penentunya** `QuotationData.ProportionalType` *(`Proportional` / `NonProportional`)* dan
`IsNewPolicyNonProp` *(`0` / `1`)*. `[terverifikasi]` Korpus menyebutnya **170** dan **87** kali.

⛔⛔ **Yang paling berbahaya — satu nama daftar, dua kedalaman:**

| Bentuk yang dirujuk aturan | Rujukan | Berkas |
| --- | ---: | ---: |
| `ListInstallment(n).InstallmentList` — **bersarang** | **58** | 6 |
| `ListInstallment(n).<medan skalar>` — **datar** | **101** | 11 |

⇒ **Pembaca dokumen memilih bentuk dari medan penentu lebih dulu**, lalu mengurai sesuai bentuk itu.
⛔ **Pengurai tunggal yang menganggap bentuknya seragam akan patah.**

⚠️ **Medan yang tidak ada pada suatu bentuk diperlakukan kosong, bukan galat.** Dari **74** medan
gabungan tiga contoh produksi, hanya **23** ada di ketiganya. `[data DBA]`

### 5.11 Peran, bukan nama orang

⚠️ `[penyimpangan sadar]` **5** — wewenang ditentukan **peran**, bukan pemeriksaan
*"apakah pengguna ini orang tertentu"*. Sistem lama menanam nama orang ke dalam teks status berkas.
`[terverifikasi]` ⛔ **Nilai itu tidak disalin ke artefak mana pun** — hanya nama medan dan
jumlahnya yang dicatat. *(Ketetapan NB 8, P4, P33; pemetaan nama → peran milik `[IAM]`.)*

---

## 6 · Testing Decisions

### 6.1 Apa yang membuat test baik di sini

Test menguji **perilaku dari luar**, bukan susunan di dalamnya. Test yang pecah ketika sebuah fungsi
dipecah dua, padahal keluarannya sama, adalah test yang salah tulis.

⭐ **Seam yang dipakai sedapat mungkin yang sudah ada, dan setinggi mungkin.** Empat seam, dan tidak
lebih:

| # | Seam | Yang diuji dari sana |
| ---: | --- | --- |
| **1** | **handler kasus endorsemen** *(HTTP)* | pembuatan berkas, tangga, penomoran, pembatalan — seluruh alur dari luar |
| **2** | **pembaca dokumen polis lama** | kedua bentuk dokumen, medan hilang, daftar bersarang lawan datar |
| **3** | **layanan perhitungan** | selisih, penyebaran, pajak, presisi |
| **4** | **penghubung konversi produksi** *(ganda)* | pengiriman, kegagalan, penghapusan bersyarat |

⛔ **Repository tidak diuji langsung.** Ia diuji lewat seam 1 dengan basis data sungguhan.

### 6.2 Yang wajib diuji terpisah

⛔⛔ **Dua bentuk `ListInstallment` diuji sebagai dua kasus terpisah.** Test yang hanya mencakup
satu bentuk **dinyatakan tidak memadai** dan tidak boleh dihitung sebagai cakupan.

⭐ **Tiga bentuk dokumen** diuji dari tiga contoh nyata — proporsional, proporsional SOA, dan
non-proporsional XOL. ⛔ **Contohnya tidak disimpan di dalam repositori** — ketiganya memuat nama
orang. Yang disimpan **cetakan bentuknya**: nama medan, kedalaman, dan jumlah.

⭐ **Endorsemen berlapis diuji sampai lapis ketiga**, supaya selisih terhadap keadaan tepat
sebelumnya benar-benar terbukti — bukan hanya lapis pertama yang kebetulan sama dengan polis asli.

### 6.3 Prior art

Bentuk test mengikuti `..\nb-treaty-in\spec.md` Bab 6 — seam handler untuk alur, seam perhitungan
untuk uang, dan uji paritas terhadap data lama. ⚠️ **Ambang paritas uang belum dapat ditetapkan** —
lihat Bab 8.9.

---

## 7 · Acceptance Criteria

⛔⛔ **Bab ini tidak memutuskan apa pun.** Ia menyatakan ulang keputusan yang sudah ada, dalam
bentuk yang dapat diuji dari luar. Bila sebuah butir terasa seperti keputusan baru, ia salah tulis.

### Jenis kasus dan alur

1. `[terverifikasi]` Endorsemen berjalan sebagai **jenis kasus tersendiri**, bukan tahap di dalam
   kasus polis baru. Test yang menemukan berkas endorsemen berbagi daur hidup dengan polis baru
   **gagal**. *(Bab 5.1)*
2. `[terverifikasi]` Tangga persetujuan punya **tiga jenjang** — Admin, Sec Head, Dept Head. Test
   yang menemukan jenjang keempat **gagal**. *(Bab 5.1)*
3. `[keputusan work owner]` Penolakan oleh **atasan** mengembalikan berkas ke admin. Test yang
   menemukan berkas tertutup **gagal**. *(Bab 5.1)*
4. `[keputusan work owner]` Penolakan oleh **admin** menutup berkas sebagai ditolak. Test yang
   menemukan berkas naik ke atasan **gagal**. *(Bab 5.1)*
5. `[keputusan work owner]` `IsApproved` dibandingkan sebagai **teks**: `0` berarti ditolak, selain
   itu disetujui. Test yang membandingkannya sebagai bilangan atau boolean **gagal**. *(Bab 5.1)*
6. `[penyimpangan sadar]` Wewenang ditentukan **peran**, bukan nama orang. Test yang menemukan
   pemeriksaan identitas perorangan di jalur wewenang **gagal**. *(Bab 5.11)*

### Pembekuan data lama — P58

7. `[keputusan work owner]` Nilai polis induk **disalin** ke berkas endorsemen pada saat berkas
   dibuat. Test yang menemukan pembacaan ulang dari polis induk sesudah itu **gagal**. *(Bab 5.2)*
8. `[terverifikasi]` Perubahan pada polis induk sesudah berkas endorsemen dibuat **tidak mengubah**
   nilai lama di berkas itu. Test yang menemukan nilainya ikut berubah **gagal**. *(Bab 5.2)*
9. `[terverifikasi]` Snapshot data lama adalah **bagian dari berkas endorsemen**, bukan rujukan yang
   dibaca saat ditampilkan. Test yang menemukan pembacaan tertunda **gagal**. *(Bab 5.2)*

### Selisih dan endorsemen berlapis — P57

10. `[keputusan work owner]` Selisih dihitung terhadap **keadaan tepat sebelumnya**, bukan terhadap
    polis asli. Test pada endorsemen ketiga yang menemukan selisih dihitung terhadap polis asli
    **gagal**. *(Bab 5.3)*
11. `[terverifikasi]` Pemilih cara hitung memakai **isi medan `EDMNo` pada data lama** *(kosong lawan
    terisi)*. Test yang memilih berdasarkan **keberadaan** simpul data lama **gagal** — simpul itu
    ada bahkan pada polis baru. *(Bab 5.3)*
12. `[keputusan work owner]` Satu polis dapat diendorse **berkali-kali**. Test yang menolak
    endorsemen kedua **gagal**. *(Bab 5.3)*

### Penomoran — P55

13. `[keputusan work owner]` Nomor endorsemen berbentuk **nomor polis induk + `/E` + dua digit**.
    Test yang menghasilkan bentuk lain **gagal**. *(Bab 5.4)*
14. `[terverifikasi]` Nomor urut endorsemen **bertambah satu** dari nomor urut terakhir polis yang
    sama. Test yang menemukan nomor terulang atau melompat **gagal**. *(Bab 5.4)*
15. `[terverifikasi]` Nomor urut di bawah sepuluh **dibubuhi nol di depan**. Test yang menghasilkan
    satu digit **gagal**. *(Bab 5.4)*
16. `[keputusan work owner]` Adendum dan premi tambahan memakai **jalur penomoran yang sama**. Test
    yang menemukan format nomor kedua **gagal**. *(Bab 5.4, Bab 8.3)*
17. `[terverifikasi]` Nomor urut dihasilkan **di dalam transaksi yang sama** dengan penyimpanannya.
    Test yang menemukan nomor terpakai walau penyimpanan gagal **gagal**. *(Bab 5.9)*

### Pembatalan — P56

18. `[keputusan work owner]` Pembatalan dipilih sebagai **jenis endorsemen** di awal. Test yang
    menemukan alur pembatalan terpisah **gagal**. *(Bab 5.5)*
19. `[terverifikasi]` Pembatalan menolkan nilai lewat jalur yang sama dengan endorsemen lain. Test
    yang menemukan jalur penolan tersendiri **gagal**. *(Bab 5.5)*

### Penyebaran risiko — P60

20. `[terverifikasi]` Bila daftar penyebaran berisi **satu baris** dan persentase pembagiannya
    **kosong atau nol**, nilai bawaan **100** ditetapkan lebih dulu. Test yang menemukan pembagian
    dengan nol **gagal**. *(Bab 5.6)*
21. `[terverifikasi]` Persentase bagian dihitung `SplitRNMSharePct / RNMShare` dengan presisi
    **20**. Test yang memakai presisi 10 **gagal**. *(Bab 5.6)*
22. `[keputusan work owner]` Ketidakseragaman presisi antara modul **ditiru apa adanya**. Test yang
    menyeragamkannya **gagal**. *(Bab 5.6)*
23. `[terverifikasi]` Penyebaran dihitung ulang ketika endorsemen mengubah nilai premi. Test yang
    menemukan penyebaran lama bertahan **gagal**. *(Bab 5.6)*

### Uang dan pajak

24. `[keputusan work owner]` Nilai uang **tidak pernah** dibaca, disimpan, atau dihitung sebagai
    bilangan mengambang. Test yang menemukan tipe mengambang di jalur uang **gagal**.
    *(Bab 5.9, ADR-0003)*
25. `[keputusan work owner]` Potongan brokerage dibagi **1,022** ketika jenis pajak bernilai
    `Inclusive` **persis**. Test yang membaginya pada ejaan lain **gagal**. *(Bab 5.6)*
26. `[terverifikasi]` PPH **2 %** dan PPN **2,2 %** dihitung dari potongan sesudah pembagian itu.
    Test yang menghitungnya dari potongan mentah **gagal**. *(Bab 5.6)*
27. `[keputusan work owner]` Empat medan potongan dan bagian adalah **persentase**, bukan uang;
    `12.5` berarti 12,5 persen. Test yang memperlakukannya sebagai nilai uang **gagal**. *(Bab 5.8)*
28. `[terverifikasi]` Nilai dibaca dan disimpan dengan **presisi penuh**; pembulatan hanya di titik
    penyajian. Test yang menemukan pembulatan di lapisan repository **gagal**. *(Bab 5.8)*
29. `[data DBA]` Uang yang melintas ke pencatatan pencapaian diseberangkan sebagai **desimal
    berpresisi tetap**. Test yang menemukan penyeberangan sebagai mengambang **gagal**. *(Bab 5.9)*

### Bentuk dokumen

30. `[terverifikasi]` Jenis treaty dibaca dari **medan penentunya** sebelum dokumen diurai. Test
    yang mengurai lebih dulu lalu menebak bentuknya **gagal**. *(Bab 5.10)*
31. `[terverifikasi]` Daftar angsuran dibaca benar pada **bentuk bersarang**. Test yang gagal
    membacanya **gagal**. *(Bab 5.10)*
32. `[terverifikasi]` Daftar angsuran dibaca benar pada **bentuk datar**. Test yang gagal membacanya
    **gagal**. *(Bab 5.10)*
33. `[terverifikasi]` ⛔ Cakupan test yang **hanya mencakup satu** dari kedua bentuk itu dinyatakan
    **tidak memadai**. Test suite yang lulus tanpa menguji keduanya **gagal**. *(Bab 6.2)*
34. `[data DBA]` Medan yang tidak ada pada suatu bentuk dokumen diperlakukan **kosong**, bukan
    galat. Test yang melempar galat atas medan yang hilang **gagal**. *(Bab 5.10)*
35. `[data DBA]` Pembacaan berhasil atas **ketiga** bentuk contoh produksi. Test yang hanya mencakup
    satu bentuk **gagal**. *(Bab 6.2)*

### Konversi ke produksi — P50 sampai P53

36. `[keputusan work owner]` Jalur konversi **FacOut tidak dibangun**. Test yang menemukan kode
    pengiriman kedua **gagal**. *(Bab 5.7, Bab 8.1)*
37. `[terverifikasi]` Pemeriksaan keberhasilan menguji **kedua** penanda, dengan cabang yang
    memperlakukan penanda kedua **kosong** sebagai sah. Test yang hanya menguji penanda pertama
    **gagal**. *(Bab 5.7)*
38. `[keputusan work owner]` Kegagalan konversi **menghapus** data produksi yang sudah tersisip.
    Test yang menemukan data tertinggal **gagal**. *(Bab 5.7)*
39. `[data DBA]` Penghapusan itu **tidak berjalan** bila kasusnya sudah bertanda konversi berhasil.
    Test yang menemukan data sah terhapus **gagal**. *(Bab 5.7)*
40. `[keputusan work owner]` Penghapusan hanya lewat **satu fungsi repository**. Test yang menemukan
    perintah hapus data produksi di tempat kedua **gagal**. *(Bab 5.7)*
41. `[keputusan work owner]` Pengguna **diberi tahu lebih dulu** bahwa kegagalan akan menghapus.
    Test yang menemukan penghapusan berjalan tanpa pemberitahuan **gagal**. *(Bab 5.7)*
42. `[keputusan work owner]` Pesan galat konversi **sampai ke pengguna**. Test yang menemukan
    kegagalan berlalu diam-diam **gagal**. *(Bab 5.7)*
43. `[keputusan work owner]` Penutupan paksa berkas **tidak dibangun**. Test yang menemukannya
    **gagal**. *(Bab 5.7, Bab 8.1)*
44. `[keputusan work owner]` Sistem baru **tidak menambahkan autentikasi** pada jalur ke produksi.
    Test yang menemukan kredensial dikirim **gagal**. *(Bab 5.7)*
45. `[penyimpangan sadar]` Alamat sistem tujuan dibaca dari **variabel lingkungan**. Test yang
    menemukan alamat harfiah di dalam kode **gagal**. *(Bab 5.7)*
46. `[keputusan work owner]` Lingkungan uji memakai **basis data berbeda** dari produksi. Test yang
    menemukan keduanya menunjuk basis data yang sama **gagal**. *(Bab 5.8)*

### Penyimpanan

47. `[keputusan work owner]` Seluruh urutan penyimpanan dibungkus **satu transaksi**. Test yang
    menemukan data separuh tersimpan sesudah kegagalan **gagal**. *(Bab 5.9)*
48. `[keputusan work owner]` Setiap query menulis **skema basis data eksplisit**. Test yang
    menemukan rujukan tabel tanpa skema **gagal**. *(Bab 5.8)*
49. `[data DBA]` Tanggal tutup buku dibaca sebagai **satu baris yang berlaku global**. Test yang
    menemukan pembacaan per periode atau per lini bisnis **gagal**. *(Bab 5.8)*
50. `[keputusan work owner]` Kolom pelaku pada riwayat akseptasi diisi dari **identitas akses
    login**. Test yang menemukannya kosong **gagal**. *(Bab 5.11)*
51. `[keputusan work owner]` Pencarian berdasarkan **nomor urut antrean** tidak dibangun; diganti
    pemeriksaan keanggotaan. Test yang menemukannya **gagal**. *(Bab 8.1)*

### Pemantauan dan pencapaian

52. `[data DBA]` Pencatatan pemantauan bersifat **idempoten**: pemanggilan kedua atas kasus yang
    sama **tidak menambah baris**. Test yang menemukan baris kedua **gagal**. *(Bab 5.9)*
53. `[data DBA]` Pencatatan pemantauan **tidak memperbarui** baris yang sudah ada. Test yang
    menemukan pembaruan **gagal** — perilakunya berhenti, bukan memperbarui. *(Bab 5.9)*
54. `[penyimpangan sadar]` Pemeriksaan duplikat pencapaian **tidak memakai nilai uang** sebagai
    pembanding. Test yang menemukan perbandingan kolom uang di jalur anti-duplikat **gagal**.
    *(Bab 5.9)*
55. `[penyimpangan sadar]` Pesan galat dari lapisan penyimpanan **tidak membawa markah tampilan**.
    Test yang menemukan markah di dalam pesan **gagal**. *(Bab 5.9)*
56. `[penyimpangan sadar]` Penanda lingkungan **tidak mengubah bunyi pesan galat** di lapisan
    penyimpanan. Test yang menemukannya **gagal**. *(Bab 5.9)*

### Lingkup

57. `[terverifikasi]` **Ke-380 langkah penetapan nilai dibangun** — tidak ada yang dibuang sebagai
    tak terjangkau. Test yang menemukan langkah terjangkau dilewati **gagal**. *(Bab 2)*
58. `[terverifikasi]` Rumus perhitungan berasal dari **isi langkah yang terbaca di ekspor**, bukan
    dari tebakan. Test yang menemukan rumus tanpa asal-usul di korpus **gagal**. *(Bab 3)*
59. `[terverifikasi]` Penelusuran langkah menelusuri **sampai enam tingkat sarang**. Test yang
    berhenti di tingkat pertama dan kehilangan 215 dari 380 langkah **gagal**. *(Bab 11.2)*

---

## 8 · Out of Scope

⚠️ **Tidak ada aturan yatim di EDM.** Yang dikeluarkan di sini adalah hal yang **mati**, **tidak
pernah berjalan**, atau **sengaja tidak ditiru** — ⛔ **bukan** aturan yang tak terjangkau.

### 8.1 Daftar penuh — setiap butir menyebut sebab dan asalnya

| # | Yang dikeluarkan | Sebab | Asal |
| ---: | --- | --- | --- |
| 1 | jalur konversi **FacOut** | EDM tidak punya FacOut; ⛔ **tidak pernah berjalan** | **P50** |
| 2 | langkah 18 `Call ASMForceCaseClose` | ⛔ tetap mati, keputusan work owner | **P52** |
| 3 | jalur penomoran kedua *(`RNM-E…`)* | ⛔ **satu jalur penomoran**, bukan dua — ⚠️ **lihat pertentangan Bab 10.1** | **P55** |
| 4 | `OldData` **di dalam** `OldData` | dibuat penyalinan seluruh halaman; ⛔ **nol rujukan** di korpus | keadaan Bab 9 |
| 5 | pencarian berdasarkan **nomor urut antrean** | diganti pemeriksaan keanggotaan | **P25** |
| 6 | pemeriksaan duplikat berbasis **nilai uang** | ⛔ dipatahkan galat presisi yang terbukti ada | Bab 5.9 |
| 7 | **markah HTML** di dalam pesan galat basis data | penyajian milik lapisan tampilan | Bab 5.9 |
| 8 | **penanda lingkungan** di dalam lapisan penyimpanan | lingkungan bukan urusan procedure | Bab 5.9 |
| 9 | blok yang dikomentari di procedure pemantauan | ⛔ tidak pernah berjalan | Bab 5.9 |
| 10 | `Struktur_InputAddendumTreatyIn.xlsx` | ⭐ lembar kerja, **bukan aturan Pega** | keadaan Bab 2.1 |
| 11 | dokumen JSON sebagai bentuk penyimpanan | diganti tabel relasional | **P29** |
| 12 | pemanggilan **stored procedure** dari Go | mandat proyek; logikanya ditulis ulang | **P51** |

### 8.2 ⛔⛔ Karantina galat uang — ditulis tersendiri

⭐ **Ditulis tersendiri karena ia satu-satunya bagian yang menunggu keputusan, bukan bahan.**

#### a. Apa yang SUDAH diketahui

```
premium angsuran   148157378,220000069     x 4  =  592629512,880000276
NetPremium         592629512,880000276          <- sama persis
nilai bulat dua desimal                         =  592629512,88
selisih                                         =  2,76 x 10^-7
```

`[data DBA]` `[terverifikasi]` Perkaliannya dihitung ulang dengan aritmetika desimal presisi 40 dan
**cocok persis**.

⭐ **Galatnya lahir di rantai perhitungan, bukan di penyimpanan:**

| | Nilai |
| --- | --- |
| galat `float64` untuk besaran ini | meleset **ke bawah**, sekitar 4,8 × 10⁻⁹ |
| galat yang benar-benar tersimpan | meleset **ke atas**, **2,76 × 10⁻⁷** |
| nisbahnya | sekitar **57 kali lebih besar**, dan **arah berlawanan** |
| pembagian empat dalam desimal | **tepat tanpa sisa** |

⇒ ⛔ Pembagian bersih tidak akan menghasilkan ekor ini. ⛔ **Sebab pastinya tidak dapat dinyatakan
dari satu kasus, dan tidak dikarang di sini.**

#### b. Kenapa EDM lebih terdampak daripada NB

⭐⭐ **EDM menghitung selisih.** Selisih dua angka yang masing-masing membawa galat membawa
**jumlah** galatnya. Dan endorsemen **berlapis** *(P57)* menumpuknya lagi pada tiap putaran.

#### c. Apa yang MENUNGGU keputusan work owner

⛔ **Tiga pilihan pada P29 — spec ini tidak memilih.**

| | Pilihan | Akibat | Ambang uji paritas |
| --- | --- | --- | --- |
| **a** | ikuti apa adanya | paritas sempurna; ⛔ cacatnya **diwariskan** dan tiap penjumlahan menambahnya | paritas **persis** |
| **b** | bulatkan saat migrasi | bersih ke depan; ⛔ **angka historis berubah**, laporan lama tidak cocok | paritas **tidak berlaku** untuk data lama |
| **c** | simpan apa adanya, bulatkan saat dihitung | sejarah utuh, galat **berhenti berlipat**; ⚠️ hasil hitung berbeda dari sistem lama | paritas **dengan toleransi** |

#### d. ⛔ Akibatnya bila TIDAK PERNAH diputuskan

| # | Akibat |
| ---: | --- |
| ⛔ **1** | **Ambang uji paritas tidak dapat ditetapkan.** Sistem baru berhitung benar, sehingga tidak menghasilkan ekor itu — ⛔ uji yang menuntut kesamaan persis **gagal justru karena sistem barunya benar**. |
| ⛔ **2** | **Penjaga anti-duplikat pencapaian tetap patah.** Selama kuncinya belum ditetapkan, baris yang seharusnya sama tetap dapat tersisip dua kali. |
| ⚠️ **3** | **Migrasi data lama tertunda arah.** Tanpa keputusan, tidak jelas apa yang dianggap benar ketika angka lama dipindahkan. |
| ⚠️ **4** | Selisih endorsemen berlapis **menumpuk galat tanpa batas yang dinyatakan**. |

⭐ **ADR-0003 tetap ditegakkan apa pun pilihannya** — sistem baru **tidak memakai `float`**,
sehingga **tidak menambah galat baru**. ⛔ **Tidak ada ADR baru yang diperlukan.**

---

## 9 · Butir `[terbuka]` — daftar penuh

⛔⛔ **Spec ini TIDAK menutup satu pun butir di bawah.** Penutupan milik pemilik perannya.

### 9.1 Yang MENAHAN pembangunan

⭐⭐ **NIHIL.** `[terverifikasi]` Pertanyaan EDM terjawab **11 dari 11**; pertanyaan NB **49 dari
49**. ⛔ Tidak ada butir yang menahan dimulainya pekerjaan.

### 9.2 Yang menentukan bentuk, tetapi tidak menahan

| # | Butir | Pemilik | Yang bergantung padanya |
| ---: | --- | --- | --- |
| 1 | ⛔⛔ **Galat angka uang tersimpan** — tiga pilihan pada P29 | `[work owner]` · `[Finance]` | ambang uji paritas · migrasi data lama |
| 2 | ⛔ **Kunci anti-duplikat pencapaian** — medan pengenal mana yang dipakai | `[work owner]` | AC 54 |
| 3 | ⛔⛔ **Dua jalur penomoran di sistem lama** — mana yang benar untuk adendum | `[work owner]` · `[Product+Underwriting]` | AC 16 · Bab 8.1 butir 3 |
| 4 | ⚠️ **Sensus properti `PolicyTreatyIn`** belum dikerjakan — 94 nama unik dan 9 simpul bersarang baru **batas bawah** | tim migrasi | bentuk tabel penyimpanan |
| 5 | ⚠️ **Selisih mana yang dipakai laporan** — yang terakhir, atau jumlah seluruhnya | `[Finance]` | pengakuan premi |
| 6 | ⚠️ **Batas berapa kali satu polis boleh diendorse** — batas teknis 99 | `[Product+Underwriting]` | validasi |
| 7 | ⚠️ **Tiga pemeriksaan pencegah yang dimatikan** di pembuatan berkas | `[work owner]` | endorsemen serentak |
| 8 | ⚠️ **`OldData` di dalam `OldData`** — dibuat, tidak pernah dibaca | `[work owner]` | Bab 8.1 butir 4 |
| 9 | ⚠️ **Skema tidak konsisten di SQL lama** — `JSON_POLIS` dua ejaan, `DATAPEGA` skema kedua | `[data DBA]` | AC 48 |
| 10 | ⚠️ **Maksud dagang dasar `SplitRNMSharePct / RNMShare`** | `[Product+Underwriting]` | AC 21 |
| 11 | ⚠️ **Apakah enam stored procedure sudah seluruhnya** — korpus mencatat lebih dari sepuluh nama | `[data DBA]` | Bab 5.9 |
| 12 | ⚠️ **Ketidakseragaman presisi 4 lawan 8 desimal** pada rumus pajak yang sama | `[Finance]` | AC 22 |
| 13 | ⚠️ **Pemetaan nama → peran** | `[IAM]` · `[work owner]` | AC 6 |

---

## 10 · Further Notes

### 10.1 ⛔⛔ Pertentangan antar sumber — ditemukan dan ditulis

#### Pertentangan 1 — **satu jalur penomoran, atau dua?**

| Sumber | Bunyinya |
| --- | --- |
| **P55** *(peringkat 1)* | penomoran **satu jalur**; adendum dan premi tambahan ikut jalur yang sama |
| **korpus** *(peringkat 6)* | ⛔ **dua jalur, keduanya terjangkau** |

`[terverifikasi]` Terukur:

| Jalur | Dipanggil dari | Bentuk yang dihasilkan |
| --- | --- | --- |
| `Activity\SetEDMTNoPolis` | `Activity\CreateEDMT` | nomor polis + `/E` + dua digit |
| `Activity\GeneratePolicyNoTreatyAddendum_Act` → `RDBList\GenerateNoEDMTreaty` | ⛔ `Section\DetailPolicyTreatyInAddendum` | `RNM-E` + kode bisnis + `.` + periode + `.` + **lima digit** |

⭐ **P55 menang** sesuai urutan kewenangan, sehingga jalur kedua **tidak dimigrasi**.
⛔ **Tetapi pertentangannya belum dijelaskan**, dan akibat salah pilih berat: nomor polis **tidak
dapat ditarik kembali sesudah terbit**. Tercatat sebagai butir `[terbuka]` **3**.

⚠️ **Dan satu ralat atas brief ronde ini:** brief menyebut `Activity\GenerateNoEDMTreaty` sebagai
*"aturan mati yang dijaga syarat `EDMNo` kosong"*. `[terverifikasi]` Terukur: ia bertipe
**`RDBList`**, bukan `Activity`; ⛔ **tidak ditemukan syarat penjaga apa pun**; dan pemanggilnya
dirujuk dari sebuah layar. ⇒ ⛔ **Ia bukan aturan mati — ia jalur kedua yang hidup.**

#### Pertentangan 2 — **ketetapan skema eksplisit lawan perilaku sistem lama**

Ketetapan NB **6** mewajibkan skema eksplisit; sistem lama hanya menaatinya pada **16 dari 31**
pernyataan. ⭐ **Tidak benar-benar bertentangan** — ketetapan itu mengikat sistem **baru**, bukan
memerikan yang lama. Dicatat supaya tidak terbaca sebagai paritas yang gagal.

### 10.2 Empat pernyataan ronde 1 yang sudah diralat

⛔ `grilling-ronde-1.md` **tersegel dan tidak disunting.** Ralatnya hidup di
`KEADAAN-EDM-TREATY-IN.md` Bab 6.

| Bunyi ronde 1 | Yang benar |
| --- | --- |
| ⛔ *"isi langkah penetapan nilai tidak terekspor"* | ⭐ **ada** — `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**; 379 dari 380 langkah, 1.309 pasangan |
| ⛔ *"pemeriksaan keberhasilan hanya memeriksa FacIn"* | ⭐ memeriksa **keduanya**, tiga cabang |
| ⛔ *"`BusinessType_DeT` tinggal 9 baris dari 34"* | ⭐ **36 baris utuh**; terverifikasi ulang dari **data produksi** |
| ⚠️ *"beda hanya metadata ekspor"* untuk 69 aturan | ⭐ **bertahan untuk ke-69-nya**; yang beda hanya **4** dari 73 |

⛔ **Dua angka ditarik dari peredaran:** **378** langkah dan **1.062** pasangan — enam rekonstruksi
dicoba, tidak satu pun menghasilkannya.

### 10.3 Lima penyimpangan sadar

| # | Penyimpangan | Alasan |
| ---: | --- | --- |
| ⚠️ **1** | Logika penghapusan data produksi **ditulis ulang di Go**, bukan memanggil procedure | mandat proyek melarang pemanggilan stored procedure; **aturannya** dipertahankan utuh |
| ⚠️ **2** | **Satu transaksi** menggantikan commit yang tersebar | sistem lama meng-`COMMIT` di dalam empat procedure; ini lebih ketat, disengaja |
| ⚠️ **3** | Pesan galat **tanpa markah tampilan** dan tanpa penanda lingkungan | penyajian milik lapisan tampilan; menyentuh **ADR-0005** |
| ⚠️ **4** | Anti-duplikat **tidak memakai nilai uang** | penjaga lama dipatahkan galat presisi yang terbukti ada |
| ⚠️ **5** | Wewenang lewat **peran**, bukan nama orang | aturan lama berhenti bekerja diam-diam ketika orangnya pindah |

### 10.4 Yang tidak boleh masuk artefak mana pun

⛔ **Nol nama orang.** Sistem lama menanam nama operator ke dalam teks status berkas; nilai itu
**tidak disalin**. ⛔ **Nomor polis tidak ditulis apa adanya.** ⛔ **Contoh JSON tidak disimpan** —
ketiganya memuat nama pemasar, operator, dan tertanggung. ⛔ **Nol rahasia, nol token, nol data
nasabah.** ⭐ Alamat dan hos internal adalah **variabel lingkungan**, tidak pernah harfiah.

---

## 11 · Lampiran

### 11.1 Keadaan terukur yang diwarisi

| Hal | Angka | Cara A | Cara B |
| --- | ---: | --- | --- |
| berkas `.xml` | **163** | `glob **/*.xml` | `os.walk` → 164 berkas *(satu `.xlsx`)* |
| ukuran modul | **17.639.259 B** | jumlah `getsize` | md5 `24afd5726e083c9df6a692da94728139` |
| `Activity` terjangkau / yatim | ⭐ **66 / 0** | urai `Call` dari titik masuk | nama disebut >1 kali |
| langkah `Property-Set` | **380** | `ElementTree`, `pySteps` puncak | pola teks |
| langkah membawa isinya | **379** *(99,7 %)* | — | — |
| pasangan nama=nilai terisi | **1.309** | pasangan dalam satu `rowdata` | pola teks bersebelahan |
| `RDBList` punya SQL | **36 / 36** | — | — |

### 11.2 Kedalaman sarang langkah

tingkat 1 = **165** · 2 = 83 · 3 = 65 · 4 = 29 · 5 = 20 · 6 = 18 — total **380**.

⚠️ **Penelusur yang berhenti di tingkat pertama kehilangan 215 dari 380 langkah.**

### 11.3 ⛔ Tujuh jebakan yang sudah menjerat proyek ini

1. ⛔⛔ **Menebak nama tag.** ⭐ `Activity` memakai `PropertiesName`/`PropertiesValue` **tanpa**
   awalan `py`; `Flow` dan `DataTransform` memakai **dengan** awalan. **Kebalikan satu sama lain.**
   `pyPropRef`, `pyPropertiesValue`, `pyStepsPage`, `pyCriteriaValue`, `pyResult`, `pyShapeName`
   **tidak ada**. ⇒ **Periksa daftar tag yang benar-benar ada lebih dulu.**
2. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>`.
3. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
4. Memasangkan dua daftar menurut **urutan**, bukan menurut blok `rowdata` yang **sama**.
5. Membaca `pyConditionString` — keterangan manusia — alih-alih bentuk yang dijalankan *(P23)*.
6. Menyatakan sesuatu nihil **tanpa menyebut lingkup penelusuran**.
7. ⭐ `<pagedata>` membungkus seluruh berkas, sehingga pola `<(\w+)>(.*?)</\1>` **menelan isinya**
   dan melaporkan nol untuk tag yang jelas berisi.

⚠️ **Buang blok `pyExpressionGadget` lebih dulu** — peninggalan layar penyunting.
⚠️ **Jangan `html.unescape` sebelum mencocokkan pola** — entitas berubah menjadi tanda kurung sudut
dan memecah polanya.
⚠️ **Jalur korpus memuat spasi** — pakai Python, bukan loop shell.

### 11.4 Cara menguji ulang angka spec ini

Perintah lengkapnya ada di `KEADAAN-EDM-TREATY-IN.md` Bab 8, siap jalan, dengan **dua cara** untuk
tiap angka. ⚠️ **Bila md5 korpus berubah, seluruh angka wajib diukur ulang** — korpus proyek ini
pernah berubah di tengah jalan.

---

### 11.5 ⭐ Peta cakupan P50–P60 ke Acceptance Criteria

⚠️ **Dapat diperiksa langsung:** tiap butir keputusan EDM punya sedikitnya satu AC.

| Butir | Ketetapan | AC |
| --- | --- | --- |
| **P50** | EDM tidak punya FacOut; pengiriman kedua tidak dibangun | **36** · 37 |
| **P51** | penghapusan bersyarat, satu jalur, pengguna diberi tahu | **38** · **39** · **40** · **41** |
| **P52** | langkah pesan dihidupkan; penutupan paksa tetap mati | **42** · **43** |
| **P53** | tanpa autentikasi dari sisi kita; alamat variabel lingkungan | **44** · **45** |
| **P54** | tanggal tutup buku satu baris berlaku global | **49** |
| **P55** | penomoran satu jalur, dua digit | **13** · **14** · **15** · **16** · 17 |
| **P56** | pembatalan jenis endorsemen, dipilih di awal | **18** · **19** |
| **P57** | endorsemen berlapis; selisih terhadap keadaan tepat sebelumnya | **10** · **11** · **12** |
| **P58** | data lama disalin dan beku | **7** · **8** · **9** |
| **P59** | lingkungan uji dan produksi berbasis data berbeda | **46** |
| **P60** | `CountSpreading_Act` 8 langkah, presisi 20 | **20** · **21** · **22** · 23 |

⭐ **Sebelas dari sebelas terwakili.** Ditambah **dua belas ketetapan NB** yang diperiksa
keberlakuannya: AC 5 · 6 · 24 · 25 · 26 · 27 · 28 · 47 · 48 · 50 · 51 — dan
⚠️ **pelanggaran skema SQL oleh sistem lama dinyatakan**, tidak didiamkan *(Bab 5.8, Bab 10.1)*.

---

## TELEMETRI EKSEKUSI

### ⛔ Pengukuran dari luar TIDAK dilakukan

Brief menyarankan `claude --print --output-format json "<prompt>" > hasil.json`.
⛔ **Itu tidak dijalankan**, dan sengaja: perintah itu akan memulai **sesi baru dan terpisah**, yang
ongkosnya bukan ongkos ronde ini.

⭐ **Angka di bawah adalah hasil pengurangan terhadap baseline** yang diambil dari catatan sesi pada
awal ronde. ⚠️ **Angka token sejati tidak terlihat dari dalam sesi**; ini pendekatan terbaik yang
tersedia, dan **bukan** pengukuran dari luar.

| Ukuran | Nilai | Sifat |
| --- | ---: | --- |
| panggilan model | **21** | terukur, selisih baseline |
| token keluar | **114.633** | terukur, selisih baseline |
| token cache ditulis | **122.850** | terukur, selisih baseline |
| token cache dibaca | **12.393.534** | terukur, selisih baseline |
| panggilan alat | ⛔ **tidak diukur terpisah** | baseline-nya tidak diambil di awal ronde; ⚠️ **tidak ditaksir lalu disajikan sebagai terukur** |
| berkas korpus dibaca | **163** | terukur |
| byte korpus dibaca | **17.639.259** | terukur |
| durasi · biaya | ⛔ **tidak diukur** | tidak tersedia dari dalam sesi |

### Disiplin ronde ini

| ⛔ Larangan | Dipatuhi? | Bukti |
| --- | :---: | --- |
| korpus read-only | ✅ | **0** berkas korpus disunting |
| `grilling-ronde-1.md` tidak disentuh | ✅ | **0** tulis |
| jangan menutup butir `[terbuka]` | ✅ | **0** ditutup; **13** tercatat di Bab 9 |
| jangan membuat ADR baru | ✅ | **0** — ADR-0003 dan ADR-0005 dirujuk, tidak ditambah |
| jangan memilih di antara tiga pilihan galat uang | ✅ | ketiganya ditulis dengan akibatnya; **0** dipilih |
| nol nama orang · nomor polis · contoh JSON | ✅ | **0** · **0** · **0** |
| nol kode, nol DDL | ✅ | cuplikan adalah **kutipan korpus**, bukan kode sistem |

---

*Disusun 22 September 2026, sesudah ronde keadaan EDM mengoreksi empat pernyataan ronde 1, menarik
dua angka yang tidak dapat direproduksi, dan menutup penahan terakhir.*
