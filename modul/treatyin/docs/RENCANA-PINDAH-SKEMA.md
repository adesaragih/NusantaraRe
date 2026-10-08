# Rencana perpindahan `treatyin` ke model acuan `Diagram-Skema-Tabel-NusantaraRe.xlsx`

**Keputusan pemilik proses 5 Oktober 2026:** modul `treatyin` **ADA DI DALAM** lingkup model
acuan. Ketiga daftar di [`PEMETAAN-SKEMA-NUSANTARARE.md`](PEMETAAN-SKEMA-NUSANTARARE.md) karena
itu berarti **"pekerjaan belum dimulai"**, bukan "beda lingkup".

⛔ **Dokumen ini RENCANA. Nol migrasi, nol DDL, nol perubahan jalur baca dijalankan saat ia
ditulis.** ⛔ **Nol tanggal, nol perkiraan lama** — yang diminta urutan dan ketergantungan.

---

## 0 · ⛔ PERTANYAAN PERTAMA — diuji terhadap data, dan jawabannya **kemungkinan 2**

*"Di mana kontrak master tinggal di dalam model berpusat-polis?"*

### 0.1 Kemungkinan 1 — `T_GENERAL_POLIS` memang menampungnya · **DISANGGAH**

**Uji A — irisan medan.** Kunci puncak dokumen dibandingkan langsung:

| | Dokumen | Kunci puncak berbeda |
| --- | --- | ---: |
| acuan | `POOLDATA.JSON_POLIS.DATA_JSON` (= `PolicyTreatyIn`) | 151 dokumen · **312** |
| kita | `POOLDATA.M_TREATY_IN.JSONDATA` (= `TreatyIn`) | 1.854 dokumen · **140** |

**Irisan: 14.** Dan tujuh di antaranya **bukan medan dagang** — `pxObjClass`,
`pxCreateDateTime`, `pxCreateOpName`, `pxCreateOperator`, `pxCreateSystemID`, `pyLabel`,
`pyRuleHarness` (perabot Pega yang ada di **setiap** dokumen Pega).

⇒ **Tujuh medan dagang yang beririsan**, dari 312 lawan 140: `Comment` · `CurrencyList` · `ID` ·
`Installment` · `Position` · `RNMShare` · `TreatyYear`. **298 kunci hanya milik polis; 126 hanya
milik kontrak master.**

**Uji B — kardinalitas.** Lewat kaidah yang sudah terbukti sendiri
(`RDBList/FetchTreatyInProductionUsingNooffer`: `SUBSTR(NOOFFER,1,7)`):

| | |
| --- | ---: |
| kontrak master (`TREATY_IN`) | **1.854** |
| polis berbeda yang terpasang padanya (`TREATYINPRODUCTION`) | **38.314** |
| kontrak dengan **1** polis | 627 |
| kontrak dengan **0** polis | 146 |
| kontrak **terbanyak** — `1000246` | **385 polis** |

⛔ **`T_GENERAL_POLIS` adalah "satu baris per GENERASI" satu polis, berkunci alami
`(NOPOLIS, PRODKE)`.** Satu kontrak master dengan 385 polis tidak dapat tinggal di baris yang
kuncinya nomor polis. **Kemungkinan 1 gugur.**

### 0.2 Kemungkinan 3 — kontrak master tidak punya tempat · **DISANGGAH oleh acuan sendiri**

Acuan menaruh `TREATY_IN` dan `M_TREATY_IN` di bawah judul *"TABEL LAMA YANG SUDAH DATAR — di
luar pohon, tidak dirancang ulang"* (`NB Treaty In Prop` **B85**, **F98**) dan menyebutnya
**"dibaca saja"**. Acuan karena itu **mengharapkan keduanya terus dibaca**, bukan ditinggalkan.

⇒ Menghentikan pembacaan `TreatyIn` akan membuang layar bagi **1.854 kontrak** tanpa satu kalimat
pun di acuan yang memintanya. **Kemungkinan 3 gugur.**

### 0.3 ⭐ Kemungkinan 2 — kontrak master butuh pohonnya sendiri · **YANG TERSISA**

Disimpulkan **dengan menyisihkan**, dan bentuk permintaannya pun sudah terbaca di acuan:

> `T_WORK_POLIS` — *"akar · **LINTAS-LINI** · 1 baris per work object"*, *"diambil dari
> `pyWorkPage.pzInsKey`"*, *"tabel yang sama dengan akar sheet 'PremiumList NB + EDM' — baris NB
> dan baris EDM duduk **SEJAJAR**, tidak saling menunjuk"* (`NB NonProp` **B5–B7**)

⭐ `T_WORK_POLIS` sudah **lintas-lini** dan sudah menampung dua lini yang duduk sejajar. Kontrak
master Treaty In adalah **work object Pega** juga (`ASM-FW-GISFW-Int-TREATY_IN`), jadi bentuk yang
paling kecil perubahannya: **satu baris `T_WORK_POLIS` lagi, dengan pohonnya sendiri di
sampingnya** — sejajar dengan `T_GENERAL_POLIS`, bukan di dalamnya.

⛔ **Bentuk itu TIDAK dikarang di sini.** Ia **permintaan ke pemilik acuan** — butir 7 di §5.
Rencana di bawah disusun supaya berlaku apa pun nama yang kelak diberikan: tahap 1–3 tidak
menyentuh bentuk sasaran sama sekali.

---

## 1 · Urutan yang menjaga data

Dua kenyataan yang menentukan urutannya:

| | |
| --- | --- |
| **Tabel tab kini 0 baris** (9 tabel, dikosongkan migrasi `436`) | ⇒ **paling murah dipindahkan lebih dulu** — nol data berisiko |
| **`TREATY_IN` 1.854 baris · `M_TREATY_IN` 1.854 baris** | ⛔ **tidak boleh disentuh kapan pun**, di tahap mana pun |

⭐ Dan satu kenyataan yang membuat seluruh rencana ini lebih murah daripada kelihatannya:
**keempat tab Limits · Share · Event Limits · RNM Share sudah tidak punya tabel** — ronde
sebelumnya mencabut `M_TREATY_IN2` sebagai sumber, dan keempatnya membaca dokumen langsung.
Yang tersisa dipindahkan hanya sembilan tabel tab.

---

## 2 · Tahapan

⛔ **Tiap tahap berdiri sendiri dan meninggalkan pohon HIJAU.** Tahap yang hanya benar bila tahap
berikutnya ikut jalan bukan tahap, dan tidak dimasukkan.

### Tahap 1 · Permintaan bentuk sasaran — **nol kode**

| | |
| --- | --- |
| **Tabel disentuh** | nol |
| **Jalur baca berubah** | nol |
| **Uji ikut** | nol |
| **Hijau sendirian?** | ⭐ ya — tidak ada yang berubah |
| **Pembalikan** | tidak perlu |

Isinya: butir 7 §5 dikirim, dan jawabannya ditunggu. ⛔ **Tahap 2 dan seterusnya tidak dapat
dimulai tanpa jawabannya** — tanpa nama tabel sasaran, perpindahan hanya dapat menebak.

⚠️ **Ini satu-satunya tahap yang memblokir.** Sisa rencana ditulis supaya terbaca sekarang, bukan
supaya dijalankan sekarang.

### Tahap 2 · Sembilan tabel tab → nama sasaran

| | |
| --- | --- |
| **Tabel disentuh** | `T_TREATY_REPORTING_PERIOD` · `T_TREATY_PORTFOLIO` · `T_TREATY_ACCUMULATION` · `T_TREATY_EGNPI` · `T_TREATY_RETENTION` · `T_TREATY_INSTALLMENT` · `T_TREATY_INSTALLMENT_ITEM` · `T_VIEW_COMMENT` · `M_TREATYIN_COINSCALE` |
| **Data berisiko** | ⭐ **NOL baris** — kesembilannya kosong sejak `436` |
| **Jalur baca berubah** | `repository/pendaratan_baca.go` · `pendaratan_peta.go` · `pendaratan_muat.go` |
| **Uji ikut** | `migrasi_pendaratan_test.go` (peta `namaDDLPendaratan`, `badanCreateTable`, penghitung `DROP COLUMN`, daftar kunci tak-mendarat) · `pendaratan_*_db_test.go` |
| **Hijau sendirian?** | ⭐ ya — tabel kosong dinamai ulang, pemuatnya mengikuti, layar tidak berubah |
| **Pembalikan** | migrasi `down` yang mengembalikan nama; nol data hilang sebab nol baris |

⭐ **Tahap inilah yang menutup `TestTCONolNamaTabelBaruDiKode`** — penjaga lintas-aplikasi yang
hari ini merah karena menganggap awalan `T_TREATY_*` milik tabel warisan modul lain. Begitu nama
sasaran dipakai, penjaga itu hijau tanpa dilonggarkan.

⚠️ **Yang HILANG di tahap ini:** keempat penjaga di `migrasi_pendaratan_test.go` yang baru saja
diajari nama `T_TREATY_*` harus diajari **lagi**. Nol di antaranya boleh dilonggarkan — yang
berubah daftar namanya, bukan ketegasannya.

### Tahap 3 · `T_VIEW_COMMENT` dipisahkan dari `T_VIEW_SUGGEST`

| | |
| --- | --- |
| **Tabel disentuh** | `T_VIEW_COMMENT` (kita, 0 baris) |
| **Jalur baca berubah** | `repository/warisan_riwayat.go` di modul Adjustment ikut |
| **Uji ikut** | `warisan_lampiran_db_test.go` (panel History) |
| **Hijau sendirian?** | ⭐ ya |
| **Pembalikan** | nama dikembalikan |

⛔ **Sebabnya ada di §4.3 pemetaan:** migrasi `436` menamai tabel komentar kontrak master
`T_VIEW_COMMENT`, sementara acuan punya `T_VIEW_SUGGEST` (Daftar Relasi **D18**, **D31**) yang
isinya `PolicyTreatyIn.SuggestList` — **tabel yang berbeda dengan pola nama yang sama**. Dibiarkan,
yang berikutnya akan mengira keduanya satu dan menggabungkannya.

### Tahap 4 · Kontrak master memperoleh akarnya

| | |
| --- | --- |
| **Tabel disentuh** | tabel baru hasil jawaban Tahap 1; **`TREATY_IN` dan `M_TREATY_IN` NOL** |
| **Jalur baca berubah** | nol — jalur lama tetap membaca dokumen |
| **Uji ikut** | uji bentuk tabel baru saja |
| **Hijau sendirian?** | ⭐ ya — tabel baru berdiri kosong di samping jalur yang berjalan |
| **Pembalikan** | `DROP` tabel baru; nol baris, nol pembaca |

⚠️ Tahap ini sengaja **tidak memindahkan data**. Tabel berdiri dulu, diisi kemudian — itu yang
membuatnya dapat dibalikkan tanpa kehilangan apa pun.

### Tahap 5 · Pemuatan kontrak master, jalur baca **ganda**

| | |
| --- | --- |
| **Tabel disentuh** | tabel baru (tulis); `M_TREATY_IN` **baca saja** |
| **Jalur baca berubah** | ditambah, bukan diganti — layar masih membaca dokumen |
| **Uji ikut** | uji silang: baris baru **sama dengan** yang dokumen berikan |
| **Hijau sendirian?** | ⭐ ya — dua sumber hidup berdampingan, dan ujilah yang mengadu |
| **Pembalikan** | kosongkan tabel baru; pembacanya belum ada |

⭐ **Di sinilah kebenarannya dibuktikan**, bukan di tahap pindah layar: selama uji silang merah,
tahap 6 tidak boleh dimulai.

### Tahap 6 · Layar berpindah ke sumber baru

| | |
| --- | --- |
| **Tabel disentuh** | nol DDL |
| **Jalur baca berubah** | `services/warisan_kontrak.go` membaca tabel baru |
| **Uji ikut** | seluruh uji layar dan `db` modul ini |
| **Hijau sendirian?** | ⭐ ya — bila tahap 5 hijau, tahap ini hanya menukar sumber |
| **Pembalikan** | tukar kembali; tabel lama dan dokumen tidak disentuh |

⚠️ **Yang HILANG di tahap ini, dan ini kerugian terbesar seluruh rencana:** seluruh uji yang
mengadu hasil layar dengan **dokumen** berubah artinya. `TestLubang510Tertutup`,
`TestTigaTingkatPohonLimitsTerisi`, `TestMedanPuncakBercabang`, dan kedua belas uji
`warisan_layer_dokumen*` mengukur penguraian `Limits[]` — begitu layar membaca tabel, uji-uji itu
mengukur jalur yang **tidak lagi dipakai layar**.

⛔ Keduanya harus dinyatakan saat itu, bukan dihapus diam-diam: entah diubah artinya menjadi
penjaga **pemuat** (tahap 5), atau dicabut beserta sebabnya. Pola yang sama dengan pencabutan
`M_TREATY_IN2` — tujuh uji, masing-masing dengan nasibnya tertulis.

### Tahap 7 · `OLD_POLIS_ID` dan `NOURUT` diadopsi

| | |
| --- | --- |
| **Tabel disentuh** | `NILAI_SELISIH` · `NILAI_SEBELUM_PRO_RATE` (modul Adjustment) |
| **Jalur baca berubah** | modul Adjustment |
| **Hijau sendirian?** | ⚠️ **hanya bila tahap 4–6 selesai** — ia bergantung pada akar kontrak master |
| **Pembalikan** | kunci alami lama dipulihkan |

⚠️ **Tahap ini BERGANTUNG**, dan karena itu ia yang terakhir. Ia disebut supaya terbaca bahwa
`OLD_POLIS_ID` (pengganti `OldData`) dan `NOURUT` (`selisih[n] = baru[n] − lama[n]`) **belum**
diadopsi — keduanya tercatat di §5 pemetaan, tidak dipasang.

---

## 3 · Ketiga bentrok yang terbawa

### 3.1 Presisi — **dapat berdiri bersama, dengan satu ketegangan yang dinyatakan**

| | |
| --- | --- |
| Acuan `F20` + `rancangan…` §4q.1 | *"uang · persen → angka presisi tetap, skala **minimal 9 desimal** (P29: galat lama diikuti apa adanya, tidak dibulatkan)"* — mengatur **kolom, migrasi, pembandingan** = **penyimpanan** |
| KEPUTUSAN §24 | desimal **per kolom** dari gambar; nol di ekor dipertahankan — mengatur **tampilan** |

⭐ Keduanya **tidak bertabrakan secara mekanis**: kolom `NUMBER(p,9)` yang menyimpan
`2484250.000000001` tetap dapat ditampilkan `2.484.250,00`.

⛔ **Ketegangan yang tetap dinyatakan:** tampilan 2 desimal **menyembunyikan ekor galat yang
penyimpanan sengaja pertahankan**. Siapa yang memeriksa selisih di layar tidak akan melihatnya.
P29 tidak berbicara tentang tampilan — ia tidak melarangnya, dan tidak pula mengizinkannya.

⇒ **Tahap 4–5 harus memakai skala ≥ 9 pada kolom barunya**, dan §24 tetap berlaku di layar. Bila
kelak terbukti pemeriksa memerlukan ekornya, yang berubah §24, bukan skalanya.

### 3.2 `T_TREATY_*` — dicabut di **Tahap 2**

Acuan nol kali memuat awalan itu. ⛔ **Tidak dicabut ronde ini**; ia Tahap 2, dan tahap itu
sekaligus menutup `TestTCONolNamaTabelBaruDiKode`.

### 3.3 `T_VIEW_COMMENT` ≠ `T_VIEW_SUGGEST` — dipisahkan di **Tahap 3**

Dua tabel berbeda, pola nama sama. Isinya berbeda sumber: komentar kontrak master lawan
`PolicyTreatyIn.SuggestList`.

⭐ **`LAYER` bukan bentrok** — ditutup, dan tidak dibuka lagi di rencana ini.

---

## 4 · Apa yang hilang, diringkas

| Tahap | Yang artinya berubah |
| --- | --- |
| 2 | empat penjaga `migrasi_pendaratan_test.go` diajari ulang nama sasaran — **nol dilonggarkan** |
| 3 | uji panel History di modul Adjustment |
| 6 | ⚠️ **dua belas uji `warisan_layer_dokumen*` + `TestLubang510Tertutup`** mengukur jalur yang tidak lagi dipakai layar — tiap satunya dinyatakan nasibnya, pola pencabutan `M_TREATY_IN2` |
| 7 | kunci alami `NILAI_*` modul Adjustment |

---

## 5 · Permintaan — kini **tujuh** butir, satu daftar

⛔ Jangan dipecah.

| # | Yang diminta | Untuk | Kepada |
| --: | --- | --- | --- |
| 1 | Rule `Rule-Obj-Property` **`AccountingMode`** | label `underwriting`/`accounting` | tim Pega |
| 2 | Rule **`AccountingModeNonProp`** | label `loss`/**`risk`** (21 kontrak) | tim Pega |
| 3 | Rule **`Bordeaux`** | label `reporting`/**`nonreporting`** (687 kontrak) | tim Pega |
| 4 | **`pyDecimalPlaces`** tiap kontrol angka | 32 kolom berdesimal `null` | tim Pega |
| 5 | Padanan nama **`TREATY_IN_POLIS*` ↔ `T_*`** | pemetaan tingkat kolom | pemilik acuan |
| 6 | ⚠️ **Ketidakkonsistenan `T_POLIS_BREAKDOWN_SPREAD`** — hidup di `Daftar Relasi D77`, dibatalkan di `NB NonProp B4`/`J63` | agar tidak dihitung sebagai tabel hidup | pemilik acuan |
| 7 | ⭐ **Bentuk akar kontrak master** — §0.3 | memblokir Tahap 2 dst. | pemilik acuan |
| 8 | ⭐ **Cacat dua penjaga lintas-aplikasi** — di bawah | agar `./inti/...` hijau | pemilik penjaga |

### 5.1 Butir 8 · `RENAME` sequence tidak dapat ditulis ber-skema — dibuktikan

`TestSetiapPernyataanSahDanBerskema` dan `TestPernyataanMulaiDenganPerintah` menuntut sesuatu
yang **tidak dapat ditulis di Oracle**. Dibuktikan terhadap instans yang dipakai, **Oracle
12.2.0.1.0**, dengan nama yang dijamin tidak ada sehingga nol objek nyata tersentuh:

| Pernyataan | Jawaban Oracle | Artinya |
| --- | --- | --- |
| `ALTER SEQUENCE POOLDATA.x RENAME TO y` | **ORA-02286** *no options specified for ALTER SEQUENCE* | `ALTER SEQUENCE` nol klausa `RENAME` |
| `ALTER SEQUENCE x RENAME TO y` | **ORA-02286** | idem tanpa skema |
| `RENAME POOLDATA.x TO y` | **ORA-01765** *specifying owner's name of the table is not allowed* | ⛔ awalan skema **dilarang keras** |
| `RENAME x TO POOLDATA.y` | **ORA-01765** | idem di sisi sasaran |
| `RENAME x TO y` | **ORA-04043** *object does not exist* | ⭐ sintaksnya **SAH** |
| `ALTER TABLE POOLDATA.x RENAME TO y` | **ORA-00942** *table or view does not exist* | ⭐ **pembanding**: tabel punya bentuk ber-skema |

⭐ Baris terakhir yang mengunci bacaannya: **tabel punya bentuk ber-skema, sequence tidak.**

**Perbaikan yang diminta:** terima `RENAME` sebagai perintah SQL, dan kecualikan `RENAME` dari
syarat ber-skema. ⛔ **Tidak dikerjakan dari sini** — keduanya penjaga lintas-aplikasi.
