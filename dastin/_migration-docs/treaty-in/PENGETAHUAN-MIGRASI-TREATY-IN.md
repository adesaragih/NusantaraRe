# Pengetahuan migrasi Treaty In — seluruhnya

**Disusun:** 24 September 2026, dari sebelas sesi pembedahan dan perancangan.
**Pasangannya:** `METODE-GRILLING.md` — berkas itu berisi CARA; berkas ini berisi APA.
**Untuk:** siapa pun yang meneruskan pekerjaan ini, termasuk sesi Treaty In Adjustment dan modul
berikutnya.

> Setiap angka di berkas ini adalah angka pada tanggal di atas. Yang berubah cepat ditandai.
> Yang tidak ditandai adalah hal yang sudah stabil selama beberapa sesi.

---

# BAGIAN 1 — PETA PROYEK

## 1.1 Apa yang sedang dikerjakan

Modul **Treaty In** di Pega 8.8 dipindahkan ke **React + Go + Oracle**. Modul **Treaty In
Adjustment** menyusul; sambungannya sudah dirancang, perilakunya belum dibedah.

Tiga alur, berurutan:

| Alur | Menghasilkan | Keadaan |
|---|---|---|
| **grilling** | pemahaman apa-adanya, dengan setiap perilaku diadili | selesai |
| **to-spec** | model data, invarian, ERD, DDL | §10 selesai 27 entitas; ERD dan DDL berjalan |
| **to-ticket** | tiket kerja | gerbang selesai, 58 kemampuan didaftar, tiket belum dibentuk |

## 1.2 Peta berkas

Seluruhnya di `D:\XML_NURE\_migration-docs\treaty-in\`.

| Berkas | Isinya |
|---|---|
| `CONTEXT.md` | aturan kerja — §2.0 sampai §2.9b, §16 penamaan |
| `PENGETAHUAN.md` | analisis apa-adanya sistem lama |
| `INVENTARIS-STRUKTUR-DATA.md` | properti dan kelas Pega |
| `SPEC-MODEL-DATA.md` | model — §10 per entitas, §11 sisanya, §12–§14 pembahasan |
| `SPEC-INVARIAN.md` | 63 invarian + 2 tidak-dapat-dilanggar, §3 kelengkapan perpindahan |
| `docs/adr/` | ADR 0003–0033 warisan, 0034–0055 rekayasa ini |
| `4-erd-dan-tabel-datar/` | STRUKTUR-DATA.md · ERD · SEAM-ADJUSTMENT.md · TITIK-BUTA-POHON.md · AUDIT-PENYEBUT-POHON.md · SAPUAN-DAN-NAMA-TAGNYA.md |
| `5-tiket/` | KEPUTUSAN-PEMBAGIAN-TIKET.md · DAFTAR-PEKERJAAN.md · BENTUK-TIKET.md · LUBANG-SPESIFIKASI.md |
| `PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql` | seluruh uji data A…AN — namanya tidak diubah supaya rujukan lama tidak putus |
| `DAFTAR-ESKALASI-MANAJEMEN.md` | yang harus diputuskan orang berwenang |
| `alat/` | perkakas sapuan, masing-masing menyebut batasnya sendiri |

---

# BAGIAN 2 — DOMAIN REASURANSI, SECUKUPNYA

Tanpa ini sebagian besar temuan tidak terbaca.

## 2.1 Dua cabang

| | **Proporsional** | **Non-proporsional** |
|---|---|---|
| dasar | penanggung ulang menanggung **persentase** setiap risiko | menanggung **lapisan kerugian** di atas ambang |
| bentuk | Quota Share (persen tetap) · Surplus (kelipatan retensi, "lines") | XOL berlapis (layer) |
| besaran khas | persen cession, retensi, EPI/EGNPI | limit, deductible, MDP, reinstatement, ROL |

Satu kontrak dapat memuat keduanya: contoh nyata memuat `Limits` baris 1 QUOTA SHARE dan baris 2
SURPLUS.

## 2.2 Istilah yang dipakai berulang

| Istilah | Artinya |
|---|---|
| **cedant** | perusahaan asuransi yang menyerahkan risiko |
| **retensi** | bagian yang ditahan cedant sendiri |
| **EPI / EGNPI** | perkiraan premi setahun — dasar banyak perhitungan |
| **MDP** | premi minimum dan deposit pada XOL |
| **reinstatement** | pemulihan limit setelah terpakai klaim; dapat berulang, bersyarat sendiri-sendiri |
| **ROL** | rate on line — premi dibagi limit |
| **brokerage / ceding / overriding commission** | potongan atas premi |
| **profit commission** | bagi hasil bila rasio klaim baik |
| **cash loss** | ambang klaim yang boleh ditagih di luar siklus rekening |
| **claim cooperation** | ambang klaim yang menuntut penanggung ulang dilibatkan |
| **retrosesi** | penanggung ulang menyerahkan lagi sebagian risikonya |
| **penyebaran (spreading)** | pembagian bagian NuRe ke susunan retro internal |
| **fakultatif** | per risiko, bukan per portofolio; `FAC-IN` masuk, `FAC-OUT` keluar |

## 2.3 NuRe

Nusantara Re — penanggung ulang. **Treaty In** = kontrak yang MASUK ke NuRe dari cedant.
Karena itu "share" di modul ini berarti bagian NuRe atas kontrak, dan "penyebaran" berarti apa yang
NuRe serahkan lagi ke dalam susunan retronya sendiri.

---

# BAGIAN 3 — SISTEM LAMA

## 3.1 Bentuknya, dan apa yang TIDAK ada

Pega 8.8, ruleset `GISFW`. Yang ada: Harness, Section, Flow Action, Activity, Data Transform,
RDB List, Report Definition, System Settings — **dua belas jenis aturan**.

**Yang NOL kemunculan — dan ini penting:** `Rule-Obj-Flow`, `Rule-Obj-CaseType`,
`Rule-Declare-Expressions`, `Rule-Declare-Trigger`, `Rule-Access-Role-Obj`, dan enam belas lain.

> **Nol kemunculan di ekspor bukan nol di sistem.** Pertanyaan itu belum terjawab (L-10), dan
> sampai terjawab, "tidak ada Flow" berarti "tidak ada Flow DI EKSPOR INI".

Akibatnya: mesin keadaan yang kita susun adalah mesin keadaan **yang terlihat dari sisi aktivitas**.

## 3.2 Tabel Oracle sistem lama

| Tabel | Isinya |
|---|---|
| `POOLDATA.M_TREATY_IN` | **dua kolom saja** — `ID` dan `JSONDATA`. Seluruh kontrak dalam satu dokumen JSON |
| `POOLDATA.TREATY_IN` | proyeksi relasional 20 kolom — **diturunkan**, bukan sumber |
| `M_TREATY_IN_EDM` / `TREATY_IN_EDM` | pasangan yang sama untuk addendum |
| `TREATYINDETAIL` / `TREATYINDETAILEDM` | 62 dan 55 kolom |
| `PROPORTIONALARRG` | master susunan retro — seluruh kolom `VARCHAR2(1000)`, tanpa PK maupun indeks |
| `TREATYINOFFER` | penerbitan ke hilir |
| `TREATYEXCHANGEYEARLY` | kurs per tahun |

Penulisnya satu: prosedur `PEGA_TREATY_IN` menulis `M_TREATY_IN` dan `TREATY_IN` dalam satu
panggilan. **Tidak ada satu pun `INSERT`/`UPDATE` langsung dari aturan.**

**Tidak ada kolom waktu di `M_TREATY_IN`.** Itu sebab tunggal dari temuan terberat di bagian 6.

## 3.3 Bentuk JSON-nya

Kepala kontrak, lalu daftar bersarang: `Limits[] → Detail[] → SpreadingList[]`, `CurrencyList[]`,
`CoInScale[]`, `Portfolio[]`, `ReportingPeriodList[]`, `Installment[] → InstallmentList[]`,
`CommentList[]`.

Dan **dua cermin** berbentuk sama dengan kepala: `ActualValue` (potret nilai sekarang) dan
`ValueDifference` (selisih). Keduanya kosong pada kontrak yang tidak pernah disesuaikan.

---

# BAGIAN 4 — ATURAN KERJA (isi `CONTEXT.md`)

## 4.1 §2.0 — Struktur yang terlihat bukan struktur yang berlaku

Induk dari hampir seluruh temuan. **Delapan turunan**, masing-masing dengan cara menemukannya
sendiri — dan cara menemukannya berbeda-beda, itu yang membuat daftarnya perlu:

| | Bentuk | Cara menemukannya |
|---|---|---|
| **a** | rujukan ≠ panggilan ≠ langkah hidup | tiga tingkat dipisah sebelum menghitung |
| **b** | `pyStepsBlockName = "//"` — langkah mati di aktivitas | periksa blok, bukan daftar langkah |
| **c** | `pyVisible = ALWAYS` mengabaikan `pyVisibleWhen` tersimpan | baca yang menang, bukan yang tersimpan |
| **d** | `pzIndexOwnerKey` menentukan cabang, bukan nama berkas | jangan percaya nama berkas |
| **e** | deskripsi langkah tidak diperbarui | baca isinya |
| **f** | kode berkondisi pada **satu baris data tertentu** | sapu tetapan |
| **g** | nilai keadaan yang **tidak punya penulis** | sapu nilai yang dibaca, cari penulisnya |
| **h** | **tidak ada di ekspor bukan tidak ada di sistem** | daftar jenis, lalu tanyakan |
| **i** | daftar bentuk penulis disusun orang, jadi kurang | susun mekanis dari nama elemen |
| **j** | sapuan satu-nama-tag buta terhadap jenis aturan lain | tiap sapuan menyebut tagnya, diperiksa terhadap daftar |

Dan satu yang menutup celah §2.0-i: **daftar nama elemen menjawab "nama lain untuk gagasan ini";
ia tidak menjawab "gagasan ini berbentuk apa lagi". Yang menjawab yang kedua adalah NOL YANG
MUSTAHIL** — properti yang punya layarnya sendiri tetapi nol penulis.

## 4.2 Aturan proses

| | Aturan |
|---|---|
| **keluaran** | ditulis saat langkahnya selesai, bukan saat sesinya selesai |
| **lubang** | dilaporkan, tidak ditambal — dan membawa **siapa yang menutupnya** dan **apa yang menagihnya** (§2.9b) |
| **cacat** | ditemukan, dicatat, lalu diputuskan — tidak pernah diperbaiki diam-diam |
| **golongan** | empat, tidak ada yang kelima |
| **perubahan** | PELESTARIAN / PERUBAHAN / BARU, dan PERUBAHAN menyebut dari apa |
| **eskalasi** | butir eskalasi bukan pekerjaan teknis |
| **alasan** | setiap hasil membawa **KENAPA BEGINI** |
| **wewenang** | ADR menetapkan apa yang dibangun; ia tidak pernah menetapkan apa yang sistem lama lakukan |
| **lingkup** | ditulis dari penulisnya, bukan dari bentuk yang terlihat |

## 4.3 §16 — penamaan

Bahasa Indonesia · `UPPER_SNAKE_CASE` · kata utuh, bukan singkatan · **≤ 30 bita, dihitung** ·
nama constraint tidak mengeja kolomnya (`CK_x_1`).

Skema: **`TREATY_MASUK`** dan **`TREATY_MASUK_APP`** — bukan `TREATYIN`, karena `POOLDATA.TREATY_IN`
sudah ada di instans yang sama dan keduanya akan disebut orang setiap hari selama masa berdampingan.

---

# BAGIAN 5 — KEPUTUSAN RANCANGAN

## 5.1 ADR yang paling sering dipakai

| ADR | Isinya |
|---|---|
| **0034** | JSON menjadi arsip, tanpa jalur baca. Bukan sumber kanonik |
| **0035** | kegagalan tidak pernah disamarkan menjadi nilai — bukan nol, bukan kurs satu |
| **0036** | angka dasar persetujuan **dibekukan saat disetujui**; angka posisi dihitung saat dibaca. Hasil beku membawa **penunjuk ke masukan yang dipakai** |
| **0037** | masukan versus turunan; turunan tidak disimpan, kecuali fakta terbukukan berpenunjuk |
| **0038** | himpunan yang dapat bertambah = tabel acuan, bukan `CHECK` |
| **0039** | paket uang membawa **tingkat pencatatan** — 100 % treaty atau bagian NuRe |
| **0040** | identitas dipecah `KONTRAK` / `VERSI_KONTRAK`; kunci alami **memperingatkan**, tidak melarang |
| **0041** | satu fakta, satu penulis |
| **0042** | sejarah pindah apa adanya; nomor warisan dipertahankan |
| **0043** | menghitung ulang adalah peristiwa bisnis tersendiri, bukan bagian migrasi |
| **0045** | jejak perubahan sebagai fakta mesin, tidak dapat dilewati |
| **0046 + 0055** | satu keadaan siklus hidup, tinggal pada **versi**; delapan keadaan, tiga belas perpindahan |
| **0049** | materialitas addendum adalah **turunan** |
| **0051** | hilir yang tidak dikenal tetap dapat membaca identitas dalam bentuk lama |
| **0052** | tidak ada jalan pintas persetujuan |
| **0053** | mata uang adalah **daftar**, bukan skalar |
| **0054** | keadaan warisan yang tidak berpadanan mendarat di `WARISAN_TAK_TERPETAKAN`, nilai aslinya tersimpan |

## 5.2 Keputusan besar di luar ADR

**Addendum adalah VERSI, bukan dokumen tersendiri.** Dibuktikan dengan menelusuri 19 properti
`Int-treaty_in_edm` satu per satu — nol fakta tersisa yang tidak punya rumah. Maka tidak ada entitas
`PENYESUAIAN`.

**Sistem lama berdampingan; sistem baru mengganti.** Faktanya: addendum yang disetujui tidak pernah
menimpa baris kontrak. Rancangannya: versi disetujui menggantikan pendahulunya. **Itu PERUBAHAN,
ditandai** — bukan pelestarian.

**Tabel proporsional dan non-proporsional DIPISAH**, menurut ambang: entitas dipisah bila **lebih
dari separuh ruasnya hanya dipakai satu cabang**. `LAYER` 13/15 → dipisah. `VERSI_KONTRAK` 1/47 →
tidak. `POTONGAN` adalah pengujinya: ruasnya identik di kedua pelekatan, jadi ambang yang menyuruh
memisahkannya adalah ambang yang salah.

**`OLDID` punya TIGA arti**, dan di sistem baru menjadi tiga atribut: versi pendahulu
(`ID_VERSI_KONTRAK_DASAR`), kontrak asal salinan (`ID_KONTRAK_DISALIN_DARI`), nomor penawaran
pra-Pega (`NOMOR_PENAWARAN_WARISAN`).

**Keadaan kedelapan `DIBATALKAN`** ditambahkan karena sistem lama tidak punya cara membuang draf
yang tidak merusak apa pun — satu-satunya jalan adalah mengajukannya supaya ditolak.

---

# BAGIAN 6 — TEMUAN

## 6.1 Yang menyentuh uang atau pihak luar

**1. Addendum yang disetujui tidak pernah sampai ke hilir.**
`SaveTreatyIn_EDM_Act` langkah 12 dan 13 mati; hanya `…DetailEdm_Act` yang hidup. `TreatyInSelectAll`
membaca `m_treaty_in` saja.

**2. `TREATYINOFFER` tidak punya satu pun penulis yang terjangkau.**
`SaveTreatyIn_Act` langkah 9 mati. Penulis hidup satu-satunya di balik tombol ber-`pyCondition =
NEVER && OperatorID.pyUserName = 'ALDO SAPUTRA'`. Isinya potret beku dari masa yang tidak diketahui
— **dan bercampur tiga arti** karena `SaveData1.NOOFFER = TreatyIn.OLDID` tanpa syarat.
Maka: **bukan sumber migrasi, bukan dasar rekonsiliasi.** Dan pertanyaan yang belum terjawab: bila
tabel itu mati, **apa yang menyuapi hilir hari ini?**

**3. Angka yang sudah disetujui dapat bergeser sendiri.** Empat sebab: penyebaran dihitung ulang ·
baris master disunting · kurs tahun berjalan · bagian NuRe historis.

**4. Nilai yang pernah disetujui tidak dapat direproduksi.** `M_TREATY_IN` dua kolom, tanpa waktu.
**EVIDENCED**, dari DDL.
> Tetapi lihat 6.4: potret bersarang mungkin menyisakan sebagiannya.

**5. Persetujuan dapat terjadi tanpa penyetuju.** `TreatyInForceResolveComplete` terkunci dua orang;
`TreatyInForceEdit` yang membukanya ber-`pyVisible = ALWAYS`.

**6. Jalur revisi memendekkan persetujuan.** `RevisionState == 1` → satu tingkat. Kepala Departemen
dan Direktur dilewati, tanpa ambang nilai.

**7. Menolak addendum MENGHAPUS barisnya.** Dan kode pencatatannya **ada dan dimatikan** —
`TreatyInDeclineConfirmation_postactEDM` langkah 1–4 mati, dua langkah `RDB remove` hidup. Akibatnya
nomor revisi mungkin dipakai ulang, dan tidak ada yang pernah dapat menjawab berapa kali sebuah
kontrak gagal diubah.

**8. Addendum diajukan tanpa satu pun pemeriksaan.** `TreatyInSubmitEDM` — CheckID, CheckError, dan
keluar-bergalat seluruhnya mati. Untuk jalur addendum, **kedelapan syarat K1 mati**, bukan enam.

**9. Aturan menyebut individu alih-alih peran.** Empat instans: tombol ber-`pyUserName = 'ALDO
SAPUTRA'` · dua nomor kontrak `1000951` dan `1000069` di enam belas cabang · **empat nama orang
sebagai tetapan di penyaluran persetujuan yang hidup sejak 2019** · kode tim disimpan di ruas
`pyTelephone`.
Akibat operasionalnya: penyaluran putus tanpa galat bila orangnya pindah · orang baru tidak dapat
menerima sampai ruas teleponnya disetel, dan langkah itu tidak tertulis di mana pun · layar menyebut
nama yang tetap, sehingga dapat menampilkan penyetuju yang sudah tidak memegang perannya.

## 6.2 Cacat di balik nilai bawaan — lima instans

Pola yang paling sering menyembunyikan kesalahan: rumus yang salah menghasilkan angka yang benar
karena nilainya kebetulan sama.

| | Yang tersembunyi |
|---|---|
| **D = M** | clamp premi tidak pernah teruji |
| **OR + RI** | penjumlahan yang kebetulan benar |
| **ROL 100/100** | dua persentase **tertukar** — `ReinstatementAmount` memakai limit, `AdditionalAmount` memakai MDP |
| **`SetReinstatementPct` menulis "100"/"100" tetap** | sistem lama tidak dapat menyatakan "pemulihan pertama 100 %, kedua 50 %" |
| **`EDMEffective = Commencement` disemai** | `ProRatePercent` selalu 100 sampai ada yang menggeser tanggalnya |

Dan satu instans yang berbeda bentuknya: **INV-50 tidak pernah gagal** karena persentase penyebaran
datang dari master yang memang berjumlah 100. Invariannya benar, penegakannya benar, dan ia tidak
pernah menguji apa pun.

## 6.3 Cacat perkakas kita sendiri

Tiga, dan ketiganya menyembunyikan hal yang nyata:

| | Cacatnya | Yang tersembunyi |
|---|---|---|
| **L-8** | penjaga `primary_ok` **membuang** properti dari kelas anak alih-alih menggantungkannya | **414 properti**, 74 di dalam lingkup |
| **sapuan penulis** | mencari `<PropertiesName>` saja; data transform memakai `<pyPropertiesName>` | **73 berkas**, 39 sasaran |
| **`langkah-hidup.py`** | mengenal `pyStepsBlockName` saja; Section/Harness/Data Transform memakai `pyDisabled` | 1.306 kemunculan, ~20 langkah mati yang berakibat |

Dan bentuk kelima penulis ditemukan lewat **nol yang mustahil**: `<pyStepsCallParams><CopyInto>` —
sasaran di parameter langkah, tak tersaring oleh kata `Target`/`Property`/`Value`.

Pelajarannya: **sapuan pemanggil justru KEBAL**, karena ia pencarian teks dan tidak mengenal tag
sama sekali. Perkakas yang lebih kasar ternyata lebih lengkap, dan sebabnya bentuk, bukan kebajikan.

## 6.4 Cacat yang mungkin menyelamatkan sesuatu

`TreatyInSetEditPre` langkah 3 membuang `OLDDATA` dan `ValueDifference` tetapi **tidak** membuang
`ActualValue` — sehingga potret lama ikut tersalin dan **potretnya bersarang**, makin dalam tiap
addendum premi.

Itu cacat. Tetapi tiap lapis sarang adalah **potret kontrak pada satu addendum sebelumnya** — dan
itu satu-satunya sisa dari 6.1 butir 4.

> **Jangan bersihkan saat migrasi sebelum dibaca sekali.**

---

# BAGIAN 7 — MODEL BARU

## 7.1 Entitas

**27 entitas tersentuh penyerahan pertama** — 21 milik Treaty In, 6 tabel acuan. §10 selesai untuk
seluruhnya.

Inti: `KONTRAK` (8 atribut, lapisan beku) → `VERSI_KONTRAK` (47) → anak-anaknya:
`MATA_UANG_KONTRAK` · `RETENSI_CEDANT` · `EGNPI` · `PORTOFOLIO` · `PERIODE_PELAPORAN` ·
`PERIODE_AKUMULASI` · `TERMIN` · `SKALA_KOASURANSI` · `BATAS_PER_BAHAYA` · `DOKUMEN_KONTRAK` ·
`CATATAN_PERSETUJUAN` · `PERISTIWA_KONTRAK` · `JEJAK_PERUBAHAN` · `LAYER` → (`DETAIL_PROPORSIONAL` |
`BAGIAN` · `PEMULIHAN_LIMIT`) · `POTONGAN` · `PENYEBARAN` → `RINCIAN_PENYEBARAN` →
`NILAI_PENYEBARAN`.

Tabel acuan: `MATA_UANG` · `JENIS_POTONGAN` · `JENIS_REASURANSI` · `BAHAYA` · `KELOMPOK_TREATY` ·
`KELAS_BISNIS`.

Di luar gelombang ini: `RETRO_KELUAR` (GEL-2) · `PENCAPAIAN` (GEL-3) · `NILAI_SELISIH` dan
`BESARAN_DAPAT_DISESUAIKAN` (Adjustment).

## 7.2 Dua pemisahan yang penting

**`CATATAN_PERSETUJUAN` dan `PERISTIWA_KONTRAK` dipisah.** `CommentList` sistem lama memuat dua
jenis baris: persetujuan (`IsApproved` terisi) dan peristiwa (`IsApproved` kosong — "Copied from
ID …", "Had Created Internal Edit"). Peristiwa tidak punya penyetuju, tingkat, maupun perpindahan.
Aturan migrasinya mekanis: `IsApproved` terisi → persetujuan; kosong → peristiwa.

**Setiap nilai uang adalah paket**: nilai + mata uang + (kurs, tanggal kurs, sumber kurs bila IDR) +
tingkat pencatatan. Ini **sifat yang dibawa setiap tiket**, bukan kemampuan tersendiri.

## 7.3 Invarian

63 bernomor + 2 tidak-dapat-dilanggar.

| Tempat | Jumlah |
|---|---:|
| CONSTRAINT terbukti | 38 |
| CONSTRAINT **belum terbukti** | 4 |
| INDEKS UNIK | 4 → **2** setelah pemisahan tabel |
| TRIGGER | 6 |
| APLIKASI | 11 |

**Empat yang belum terbukti** — INV-47, INV-50, INV-51, dan bentuk lintas baris INV-31 — ditegakkan
lewat *materialized view* ber-`REFRESH ON COMMIT` dengan `CHECK` padanya. **Belum pernah dijalankan
di Oracle.** Sampai uji negatifnya lulus, ia **klaim, bukan fakta** — dan bila gagal, keempatnya
turun ke aplikasi dan bagian APLIKASI menjadi 23,8 %.

> **MV yang gagal me-refresh berhenti menegakkan tanpa satu galat pun.** Itu §2.0 yang lahir dari
> rancangan kita sendiri, dan ia menuntut uji negatif serta pemantau kebasian.

**TERMINAL BERARTI BEKU, BUKAN HANYA BERHENTI BERPINDAH** — INV-24 menjangkau `DISETUJUI`,
`DITOLAK`, dan `DIBATALKAN`, beserta seluruh entitas anaknya, dengan trigger
`BEFORE INSERT OR UPDATE OR DELETE` pada anak — karena UPDATE/DELETE saja masih mengizinkan versi
disetujui **bertambah** baris.

## 7.4 Siklus hidup

Delapan keadaan, tiga belas perpindahan. Tiga tingkat persetujuan: Section Head → Dept Head →
Direktur. Pengaju tidak boleh menyetujui; tiap tingkat harus berbeda orang dari tingkat sebelumnya.

`DRAFT → DIBATALKAN` baru, oleh pembuatnya sendiri, barisnya tidak dihapus, nomornya tidak dipakai
ulang (INV-04 sudah menutupnya tanpa mekanisme tambahan).

**K1** — delapan syarat pengajuan, **enam mati di jalur kontrak dan kedelapannya mati di jalur
addendum**. Menyalakannya adalah kemampuan BARU, bukan pelestarian, dan menuntut rencana: ukur dulu
terhadap produksi, nyalakan berperingatan, tetapkan **angka** yang memicu blokir.

---

# BAGIAN 8 — YANG TERBUKA

## 8.1 Lubang spesifikasi

| | Isinya | Siapa menutup |
|---|---|---|
| L-1 | folder `ddl/` belum ada | sesi DDL |
| L-2 | §10 — **DITUTUP** | — |
| L-4 | tidak ada spesifikasi layar di mana pun | mungkin tidak perlu ditutup |
| L-5 | huruf `G` bermakna tiga hal — **DITUTUP** dengan `GEL-2`/`GEL-3` | — |
| L-6 | daftar entitas §2.3 tidak lengkap; `STRUKTUR-DATA.md` mengikat | to-spec |
| L-7 | tidak ada cara membuang draf — **DITUTUP** | — |
| L-8 | titik buta pohon, 414 properti | to-spec |
| L-9 | asal ekspor tidak diketahui — `LinkService` menunjuk QA | kantor |
| L-10 | 21 jenis aturan nol kemunculan | kantor |

## 8.2 Yang hanya kantor dapat sediakan

1. **Satu instans Oracle** — tanpa ini empat invarian tetap klaim.
2. **Izin kueri baca-saja ke produksi** — uji A sampai AN.
3. **Ekspor kedua, ~15 berkas**, untuk dibandingkan byte per byte; plus **daftar jenis aturan** yang
   ada untuk aplikasi ini. Dua persen dari ekspor, dan ia menutup L-9 dan L-10.
4. **Dokumen batas wewenang persetujuan** — P-41 tertahan dan tidak dapat dipecah.
5. **Pemilik tabel acuan pembagian kapasitas** — P-49.

## 8.3 Yang hanya orang dapat jawab

- Ekspor ini diambil dari lingkungan mana, dan oleh siapa?
- Angka kontrak yang keluar ke akuntansi, ke retro, dan ke pihak lawan — **datangnya dari mana, dan
  siapa yang mengerjakannya?** (Jangan sebut nama tabel; tidak ada yang tahu nama tabel.)
- Kalau ada yang salah membuat draf atau addendum yang tidak jadi, selama ini mereka melakukan apa?
- Siapa yang berwenang memutuskan keadaan sah sebuah kontrak warisan?
- Kontrak `1000951` dan `1000069` itu apa, dan kenapa diperlakukan khusus?
- Empat singkatan yang korpusnya sudah disapu habis: **RSMD, WPC, RIOGR, RIONR**.
- Satu slip non-proporsional — menjawab lima hal sekaligus.

---

# BAGIAN 9 — BEKAL UNTUK SESI ADJUSTMENT

**Ekspor "Treaty In Adjustment" berisi 379 berkas; 323 di antaranya BYTE-IDENTIK dengan Treaty In.**
Permukaan yang sebenarnya **56 berkas**: 14 seksi `*OldData`, tujuh aturan `*EDM*`, dua picker, dan
kelas `Int-treaty_in_edm` dengan 19 properti.

**Mulai dari 56, bukan dari 379.** Dan untuk setiap temuan, sebutkan berkasnya ada di sisi mana:
irisan 323 → temuan **Treaty In**; 56 khas → diparkir sampai sesinya.

Sudah diputuskan dan tidak dibuka ulang: addendum adalah versi · pengenal `‹kontrak›/Rnn` · tabel
selisih berbentuk **sempit**, `NILAI_SELISIH` menggantung pada `VERSI_KONTRAK`.

Sudah diketahui dan menunggu: `ROWNUM = 1` mendahului `ORDER BY` pada kueri nomor revisi · penomoran
hanya `R01..R99` · offset tetap mengandaikan pengenal kontrak tujuh karakter · delapan butir di
`TEMUAN-ADJUSTMENT-DITUNDA.md`.

Aturan bisnis yang mengikat: **OLDDATA tidak disimpan ulang, harus di-SELECT dari data lama** ·
**nilai selisih = nilai sekarang − nilai lama** · tabel tersendiri untuk baris nilai selisih.

---

# BAGIAN 10 — KALIMAT YANG TERBUKTI BERGUNA

> Ini kebutuhan bisnis, atau akibat dari cara sistem lama dibangun?

> Struktur yang terlihat bukan struktur yang berlaku.

> Diamnya satu aturan bukan diamnya sistem.

> Rujukan bukan panggilan. Panggilan bukan langkah hidup.

> Tidak ada di ekspor bukan tidak ada di sistem.

> Bukan bukti bahwa ia akan gagal — bukti bahwa ia gagal.

> Kriteria yang tidak bisa gagal bukan kriteria.

> Keadaan seperti apa yang harus ada agar pelanggarannya dapat dinyatakan?

> Daftar dengan alasan per baris dapat diperiksa orang lain. Persentase tidak.

> Selisih antara "ada lubang, kami tahu" dan "74 di dalam lingkup, ini daftarnya" adalah selisih
> antara catatan dan bukti.

> Kolom yang mencatat di mana sesuatu terjadi pada aturan tidak pernah mencatat dari mana salinan
> ini diambil.

> Lubang yang dicatat tetapi tidak ditindaklanjuti berperilaku persis seperti lubang yang tidak
> diketahui.

> Penghalang yang dipindah ke tempat yang salah akan ditunggu di tempat yang salah.

> Terminal berarti beku, bukan hanya berhenti berpindah.

> Keluaran ditulis saat langkahnya selesai, bukan saat sesinya selesai.
