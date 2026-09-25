# Setiap sapuan menyebut nama tag yang dipakainya — dan diperiksa terhadap daftar nama elemen

**Tanggal:** 24 September 2026
**Sebab:** ini kedua kalinya. `<PropertiesName>` versus `<pyPropertiesName>` menyembunyikan 73 berkas
*data transform* selama sebelas sesi. Pertanyaannya lalu: **sapuan mana lagi yang bersandar pada satu
nama tag?**

> ## ATURAN
>
> **SEBUAH SAPUAN ATAS KORPUS YANG BERISI BANYAK JENIS ATURAN TIDAK BOLEH BERSANDAR PADA SATU NAMA
> TAG.** Setiap sapuan **mencatat nama tag yang dipakainya di dalam dirinya sendiri**, dan nama itu
> **diperiksa terhadap `alat/datar-nama-elemen.csv`** sebelum hasilnya dipercaya.
>
> Sebabnya bukan kerapian: jenis aturan yang berbeda memakai **nama yang berbeda untuk gagasan yang
> sama**, dan sapuan satu-nama-tag **diam** terhadap sisanya — diam yang terbaca sebagai nol.

`alat/sapu-nama-elemen.py` menyusun daftarnya secara mekanis: **1.823 nama elemen** dari **708
berkas**, nol berkas gagal diurai. Ia mencetak, untuk tiap nama, **jenis aturan tempat ia muncul** —
dan kolom itulah yang menjawab pertanyaan di atas.

---

## 1. Hasil pemeriksaan — sembilan perkakas, **dua terdampak**

| Perkakas | Nama tag yang dipakainya | Terdampak? |
|---|---|---|
| `sapu-properti-per-kelas.py` | `pxRuleClassName`, `pyRuleName` (indeks aturan Pega) | **tidak** — indeks, bukan bentuk |
| `buat-pohon-treatyin.py` | idem + `pyPagesAndClasses*` | tidak *(lubangnya L-8, sebab lain)* |
| `panen-kolom-dari-penulis.py` | — membaca `PROCEDURE/*.txt`, bukan XML | tidak |
| `periksa-nama-kueri-dba.py` | — membaca DDL | tidak |
| `sapu-nama-elemen.py` | `pxObjClass` saja, untuk melabeli | tidak |
| `sapu-rujukan-aturan.py` | `pxRuleReferences` | tidak *(batas lain, §5)* |
| `sapu-penulis-properti.py` | keempat bentuk, didaftar di kepalanya | **diperbaiki hari ini** |
| **`sapu-tetapan-di-kode.py`** | `<PropertiesName>` + `<PropertiesValue>` — **regex, Activity saja** | **YA — §2** |
| **`langkah-hidup.py`** | `pyStepsBlockName` | **YA — §4, dan ini yang terberat** |

---

## 2. `sapu-tetapan-di-kode.py` — buta terhadap tetapan di *data transform*

Ia mencari pasangan `<PropertiesName>…</PropertiesName><PropertiesValue>…</PropertiesValue>` dengan
regex. *Data transform* memakai `<pyPropertiesName>` / `<pyPropertiesValue>`. Maka **seluruh 449
tetapan yang pernah dilaporkan berasal dari aktivitas saja.**

Dijalankan ulang atas bentuk *data transform*: **26 pasangan (sasaran, nilai tetap)** yang belum
pernah terlihat. Yang berakibat:

| Sasaran | Nilai tetap | Berkas |
|---|---|---|
| **`TreatyIn.PositionUsername`** | **`"IRVANDY"`, `"AGUNGPUTRAANDALAS"`, `"YOHANESKRISTIAWAN"`, `"NANDINA"`** | **`Akseptasi_DT.xml`** — §3 |
| `TreatyIn.ReportingPeriod` | `"quarter"` | `TreatyInSetPeriod.xml` |
| `TreatyIn.AddendumPremi` | `"1"` | `TreatyInSetEditPre.xml` |
| `TreatyIn.EDMState` | `"1"`, `"3"` | `TreatyCreateEDM.xml` |
| `TreatyIn.Position` / `.StatusAkseptasi` / `.RevisionState` | nilai keadaan | `Akseptasi_DT.xml` |

`ReportingPeriod = "quarter"` layak dicatat tersendiri: `SPEC-MODEL-DATA.md` §10.2 mencatatnya
sebagai atribut bertipe **E** yang **boleh kosong**, seolah ia pilihan pengguna. Ada aturan yang
menyetelnya ke satu nilai tetap. **Apakah nilai lain pernah dipakai adalah pertanyaan data**, dan ia
bentuk *"cacat di balik nilai bawaan"* yang sudah dikenal — bukan temuan yang diputuskan di sini.

---

## 3. Empat nama orang tertulis di dalam aturan persetujuan yang HIDUP

`DataTransform/Akseptasi_DT.xml` adalah mesin keadaan persetujuan. Setiap langkahnya menyetel **dua**
hal: `Position` — perannya — **dan** `PositionUsername`, **nama orang, sebagai tetapan**.

```
1.1  Position=="" atau ReasTreatyInAdmin        (pengajuan)
     1.1.1  OperatorID.pyTelephone=="SPVTREATY1" -> SecHead, Accept, PositionUsername "IRVANDY"
     1.1.2  OperatorID.pyTelephone=="SPVTREATY2" -> SecHead, Accept, PositionUsername "AGUNGPUTRAANDALAS"
     1.1.3  OperatorID.pyTelephone=="TREATY1"    -> SecHead, Accept, PositionUsername "IRVANDY"
     1.1.4  OperatorID.pyTelephone=="TREATY2"    -> SecHead, Accept, PositionUsername "AGUNGPUTRAANDALAS"
1.2  Position=="ReasTreatyInSecHead"   + Accept -> DeptHead, PositionUsername "YOHANESKRISTIAWAN"
1.3  Position=="ReasTreatyInDeptHead"  + Accept -> Director, PositionUsername "NANDINA"
1.4  Position=="ReasTreatyInGroupLeader"        -> Director                          *** MATI ***
1.5  Position=="ReasTreatyInDirector"  + Accept -> "", StatusAkseptasi "Resolve Complete"
```

| | |
|---|---|
| dibuat | **27 Juni 2019**, oleh `DANIELWISESA`, di sistem `pega` |
| keadaan | **tujuh dari delapan penugasan nama HIDUP**; hanya cabang 1.4 (`GroupLeader`) yang ber-`pyDisabled=true` |
| bentuknya | keluarga yang sama dengan `TreatyIn.ID=="1000951"` — **turunan keenam §2.0**, kode berkondisi pada satu baris data tertentu; di sini barisnya **orang** |

**Dan `OperatorID.pyTelephone` dipakai sebagai kunci penyaluran.** Nilainya `"SPVTREATY1"`,
`"SPVTREATY2"`, `"TREATY1"`, `"TREATY2"` — **kode tim, disimpan di ruas nomor telepon**. Ini
`membaca-export-pega-disiplin` §5 dalam bentuk paling telanjang: **satu nama untuk dua arti**, dan
kali ini arti keduanya menentukan **siapa yang menyetujui kontrak**.

> **Akibatnya pada model: tidak ada.** `PositionUsername` sudah **dibuang** (`SPEC-MODEL-DATA.md`
> §12.5) dan `Position` sudah dilebur (ADR-0046). Yang berubah adalah **dasar** pembuangannya: dari
> penalaran ADR-0044 menjadi **`EVIDENCED`** — berkas, langkah, tanggal, dan penulisnya dapat
> ditunjuk.
>
> **Akibatnya pada eskalasi: ada, dan besar.** Ia naik sebagai butir tersendiri.

---

## 4. `langkah-hidup.py` — **ada DUA penanda mati, dan ia hanya mengenal satu**

Ini yang terberat, karena ia menopang setiap kalimat *"langkah ini mati"* dan *"cabang ini tidak
berjalan"*.

| Penanda | Dikenali `langkah-hidup.py`? | Kemunculan | Di jenis aturan mana |
|---|---|---:|---|
| `pyStepsBlockName == "//"` | **ya** | **235** | `Rule-Obj-Activity` — **dan hanya itu** |
| **`pyDisabled == "true"`** | **TIDAK** | **1.306 baris** | `Rule-HTML-Section` 1.082 · `Rule-HTML-Harness` 198 · **`Rule-Obj-Model` 26** |

> **Satuannya BARIS, bukan tempat — dan itu harus disebut, karena tanpanya angkanya salah dibaca.**
> Ke-26 di `Rule-Obj-Model` adalah **13 baris × dua salinan** `Akseptasi_DT.xml`, satu di tiap modul,
> dan kedua salinan itu **identik secara logika**. Ketiga-belas baris itu adalah **empat cabang**
> (`1.4`, `1.4.1`, `1.4.2`, `1.4.3`) beserta **sembilan penugasan** di dalamnya — seluruhnya cabang
> `GroupLeader`, yang **sudah dibuka**.
>
> Jadi **nol cabang mesin keadaan yang belum dilihat**. Yang belum pernah didaftar adalah langkah
> mati di **aktivitas** yang menjalankan mesin itu: **20 langkah**, dan di antaranya satu aktivitas
> yang seluruhnya mati. Daftarnya `CABANG-MATI-DI-JALUR-PERSETUJUAN.md`.

**Irisan keduanya nol.** Penanda yang dikenali perkakas tidak pernah muncul di luar aktivitas; penanda
yang tidak dikenalinya tidak pernah muncul di dalam aktivitas. Maka untuk **setiap jenis aturan selain
aktivitas, perkakas liveness kami tidak pernah memeriksa apa pun** — dan diamnya terbaca sebagai
*"semuanya hidup"*.

**Ke-26 di `Rule-Obj-Model` seluruhnya di `Akseptasi_DT.xml`** — mesin keadaan persetujuan. Cabang
`GroupLeader` ditemukan mati **dengan tangan**, bukan oleh perkakas. Bila ia tidak kebetulan dibuka,
ia akan terbaca sebagai cabang hidup, dan 7.10 akan diputuskan terbalik.

**Batas yang harus disebut sebelum angka 1.306 dipakai:** `pyDisabled` di `Section` dan `Harness`
hampir pasti menandai **kendali layar yang dimatikan**, bukan kode mati — fakta yang berbeda,
walaupun keduanya menjawab *"dapatkah ini dijalankan"*. Hanya ke-26 di `Rule-Obj-Model` yang terbukti
**cabang logika mati**. Ketiganya tetap harus dikenali perkakas; yang tidak boleh adalah
**menjumlahkannya sebagai satu besaran**.

---

## 5. Sapuan pemanggil — **TIDAK terdampak, dan sebabnya bukan kehati-hatian**

Kekhawatiran yang wajar: bila sapuan pemanggil juga memakai satu nama tag, maka *"tidak ada jalur
terjangkau yang menulis `TREATYINOFFER`"* dan *"aktivitas ini tidak pernah dipanggil"* ikut runtuh.

**Ia tidak memakai nama tag sama sekali.** Sapuan pemanggil dijalankan sebagai **pencarian teks atas
seluruh berkas**, dan pencarian teks **tidak mengenal tag** — karena itu ia kebal terhadap kegagalan
ini **secara bentuk**, bukan karena ada yang memikirkannya.

Dua pemeriksaan dijalankan, dan keduanya diukur.

**Pemeriksaan pertama — indeks rujukan bawaan Pega terhadap pencarian teks:**

| | |
|---|---:|
| aktivitas di korpus yang **indeks rujukan Pega** sebut **nol perujuk** | **63** |
| di antaranya yang **ditemukan pencarian teks** di berkas lain | **63 — seluruhnya** |

**Pemeriksaan kedua — sapuan sebelas nama tag rujukan terhadap pencarian teks:**

| | |
|---|---:|
| aktivitas di korpus | **157** |
| perujuknya ditemukan **oleh keduanya** | **153** |
| **ditemukan teks, NOL oleh sapuan sebelas tag** | **4** |

Keempatnya: `RefreshAchievement`, `TreatyInEDMSetValue`, `TreatyLoadMasterJoinEdm`,
`TreatyLoadMasterJoinEdmXOL` — dirujuk dari `PickerTreatyInMaster`, `PickerTreatyInMasterRevisi`,
`DetailLimitsOldData`, lewat nama elemen yang **tidak masuk daftar sebelas** yang saya susun.

> **Empat dari 157 adalah angka kecil, dan justru itu pelajarannya.** Sapuan berbasis tag yang
> disusun dengan hati-hati **tetap** kurang, karena daftar tagnya disusun oleh orang. Pencarian teks
> tidak dapat kurang dengan cara itu.
>
> **Koreksi atas percobaan pertama saya:** percobaan pertama melaporkan **157 dari 157** terlewat.
> Itu **salah** — ia mencocokkan `pyStepsActivityName` dengan nama aktivitas secara penuh, padahal
> nilainya berbunyi `"Call AddCommentList_Act"`: **metode dan nama yang dipanggil menyatu dalam satu
> nilai**. Angka yang benar 4, dan angka yang salah tidak sempat dikutip.

> **Indeks `pxRuleReferences` tidak dapat dipakai untuk menyimpulkan "tidak dipanggil".** Ia mencatat
> properti dan kelas yang dirujuk, bukan setiap pemanggilan aktivitas. Perkakas
> `alat/sapu-rujukan-aturan.py` menuliskan batas itu di dalam dirinya sendiri, bersama **86 berkas
> yang tidak memuat indeks itu sama sekali** (84 RDBList, 2 SystemSettings) — dan **diamnya untuk
> ke-86 itu bukan nol rujukan**.

Jadi: yang dipakai selama ini benar, yang baru dibangun hari ini **lebih lemah** untuk pertanyaan itu,
dan **keduanya sekarang menyatakan batasnya sendiri.**

---

## 6. Yang ditagih dari aturan ini

| Perkakas | Yang harus dikerjakan | Keadaan |
|---|---|---|
| `sapu-penulis-properti.py` | menyebut keempat bentuk di kepalanya | **selesai** |
| `sapu-nama-elemen.py`, `sapu-rujukan-aturan.py` | menyebut batas kemampuannya | **selesai** |
| `sapu-tetapan-di-kode.py` | **diperluas ke `pyPropertiesName`**, lalu 449 dihitung ulang | **tertunda** — hasil bentuk keduanya sudah ada di §2; penggabungan angkanya belum |
| `langkah-hidup.py` | **mengenali `pyDisabled`**, dan melaporkan kedua penanda **terpisah** | **tertunda** |

Keduanya tertunda bukan karena sulit, melainkan karena **angka gabungannya akan dikutip**, dan
penggabungan dua besaran yang artinya berbeda lebih mahal daripada dua angka terpisah yang jujur.
