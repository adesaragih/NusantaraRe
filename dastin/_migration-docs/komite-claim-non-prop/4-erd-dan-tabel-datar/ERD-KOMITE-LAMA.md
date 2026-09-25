# ERD dan tabel datar — Komite Claim Non Prop, sistem lama

<!-- STEMPEL ASAL -->
> **Dibangkitkan** `alat/buat-skema-komite.py` pada **22 September 2026**.
>
> **Yang digambar: pohon halaman kerja Komite, didatarkan menjadi tabel.** Sumbernya `struktur-komite-lama.md` dan `datar-komite-lama.csv` — 688 simpul dari sapuan **59 berkas XML seluruhnya**.
>
> **Ini SISTEM LAMA, bukan rancangan skema baru.** Rancangan yang berlaku hidup di [`SPEC-KOMITE-01.md`](../SPEC-KOMITE-01.md) bagian *Model data* — **tujuh objek** `KLAIMNP`, 56 kolom. Empat belas kotak di sini memerikan bentuk yang **ada**, termasuk cacatnya.
>
> **TURUNAN.** Bila berkas ini berbeda dari `datar-komite-lama.csv`, **alatnya yang salah**.

### → [`ERD-KOMITE-LAMA.html`](ERD-KOMITE-LAMA.html) · [`Diagram-Skema-Tabel-Komite.xlsx`](Diagram-Skema-Tabel-Komite.xlsx)

```
python alat/buat-skema-komite.py
```

---

## 1. Empat keluaran

| Berkas | Isi |
|---|---|
| [`datar-komite-tabel.csv`](datar-komite-tabel.csv) | **tabel datar** — 85 baris, satu baris per kolom, berikut jalur Pega asalnya |
| [`RELASI-KOMITE.csv`](RELASI-KOMITE.csv) | 14 relasi, satu baris per relasi |
| `ERD-KOMITE-LAMA.html` | ERD 14 kotak + tabel datar berpenyaring, terang dan gelap |
| `Diagram-Skema-Tabel-Komite.xlsx` | empat sheet: `BACA-DULU` · `TABEL-DATAR` · `RELASI` · `DAFTAR-TABEL` |

Tata letak sheet seragam dengan folder sisi Klaim: judul baris 1, catatan baris 2, kepala kolom baris 4, **data mulai baris 5**.

## 2. Sumber

| Sumber | Dipakai untuk |
|---|---|
| `4-erd-dan-tabel-datar/datar-komite-lama.csv` | 688 simpul, cacah `REF` dan `CACAH_BERKAS` per kolom |
| `4-erd-dan-tabel-datar/struktur-komite-lama.md` | bentuk pohon, tingkat keyakinan, peta kelas |
| `claim-non-prop/alat/buat-skema-claimnp.py` | konvensi penamaan, bentuk kotak, warna kelompok |
| `Diagram-Skema-Tabel-NusantaraRe.xlsx` | akar `T_WORK_CLAIM` dan awalan `T_KOMITE_` |

Cacah `REF` **tidak diketik**; ia dibaca dari CSV sapuan.

## 3. Penamaan tabel — dan satu tebakan yang kini dapat diperiksa

Awalan mengikuti keputusan work owner 20 September 2026: `T_WORK_CLAIM` sebagai akar lintas-lini, `T_GENERAL_*` untuk badan 1:1 ber-SHARED PK, lalu awalan lini. Modul ini memakai **`T_KOMITE_`**.

Dua nama **sudah ada** di rancangan sisi Klaim dan dipakai ulang apa adanya, bukan dibuat baru:

| Nama | Dari | Keadaan sesudah foldernya dibuka |
|---|---|---|
| `T_GENERAL_KOMITE` | `buat-skema-claimnp.py`, kelompok `komite` | **dikoreksi** — lihat di bawah |
| `T_KOMITE_KOMITELIST` | idem | **dikukuhkan** |

Berkas sisi Klaim menulis peringatannya sendiri: *"⚠ Folder Komite Claim Non Prop **TIDAK dibuka**. Bentuk tabel ini diambil dari sheet Komite Claim Prop yang sudah ada — lintas-lini, **tujuh kolom**."* Foldernya kini dibuka, dan hasilnya dua-duanya perlu dicatat:

- **`T_KOMITE_KOMITELIST` — tebakannya benar.** Ia memang **sembilan kolom**, persis seperti yang ditulis tanpa membuka folder. Lima terbaca jalur penuh, empat lewat rujukan relatif.
- **`T_GENERAL_KOMITE` — tebakannya kurang.** Bukan tujuh kolom melainkan **enam belas**, karena halaman kerja Komite Non Prop membawa dua penyimpan giliran, tiga penanda jalur, dan delapan jejak Pega yang tidak ada di sheet Claim Prop.

Selisih itu **tidak membatalkan** rancangan sisi Klaim: `T_GENERAL_KOMITE` di sana digambar sebagai kotak lintas-lini dengan hanya PK dan FK, dan kolomnya memang tidak diklaim. Yang berubah adalah cacahnya kini berbukti, bukan dipinjam.

Dua belas nama sisanya lahir di sini, seluruhnya dari simpul yang ada di pohon.

## 4. Empat belas tabel

| # | Tabel | Kard. | Induk | Kolom | Kelompok |
|---|---|---|---|---:|---|
| 1 | `T_WORK_CLAIM` | akar | — | 3 | akar |
| 2 | `T_GENERAL_KOMITE` | 1:1 | `T_WORK_CLAIM` SHARED PK | 16 | komite |
| 3 | `T_KOMITE_KOMITELIST` | 1:N | `T_GENERAL_KOMITE` | 10 | komite |
| 4 | `T_KOMITE_LOSS_CONTEXT` | 1:1 | `T_GENERAL_KOMITE` | 5 | komite |
| 5 | `T_KOMITE_ADJ_SNAPSHOT` | 1:1 | `T_GENERAL_KOMITE` | 10 | **salinan** |
| 6 | `T_KOMITE_CLAIM_SNAPSHOT` | 1:1 | `T_GENERAL_KOMITE` | 6 | **salinan** |
| 7 | `T_KOMITE_CLAIM_POLICY` | 1:1 | `T_KOMITE_CLAIM_SNAPSHOT` | 1 | tipis |
| 8 | `T_KOMITE_CLAIM_ADJ` | 1:N | `T_KOMITE_CLAIM_SNAPSHOT` | 1 | tipis |
| 9 | `T_KOMITE_ADJ_CURRENCY` | 1:N | `T_KOMITE_CLAIM_ADJ` | 2 | tipis |
| 10 | `T_KOMITE_CLAIM_OBJECT` | 1:N | `T_KOMITE_CLAIM_SNAPSHOT` | 1 | tipis |
| 11 | `T_KOMITE_OBJECT_ITEM` | 1:N | `T_KOMITE_CLAIM_OBJECT` | 1 | tipis |
| 12 | `T_KOMITE_OBJECT_ITEM_ADJ` | 1:N | `T_KOMITE_OBJECT_ITEM` | 2 | tipis |
| 13 | `T_KOMITE_QUOTATION` | 1:1 | `T_GENERAL_KOMITE` | 1 | tipis |
| 14 | `T_KOMITE_TREATY_MASTER` | 1:1 | `T_GENERAL_KOMITE` | 1 | tipis |

**60 kolom** di luar kunci · **14 relasi**, satu di antaranya menyeberang ke sisi Klaim (`T_GENERAL_KOMITE.ADJUSTMENT_ID → T_CLAIM_ADJUSTMENT.ID`).

### 4.1 Kunci — PK = `ID` di mana-mana, dan apa yang sebenarnya ada di baliknya

Tiap tabel berkunci primer **`ID`**. Dua di antaranya berdiri di atas pengenal yang **benar-benar ada** di sistem lama; dua belas sisanya **surogat yang lahir dari pendataran ini**, dan itu dinyatakan per baris di kolom `KEYAKINAN`, bukan disamarkan.

| Tabel | PK berasal dari | Keyakinan |
|---|---|---|
| `T_WORK_CLAIM` | `pyWorkPage.pyID` | `jalur-penuh` |
| `T_GENERAL_KOMITE` | `pyWorkPage.pzInsKey` — SHARED PK dengan akarnya | `jalur-penuh` |
| dua belas lainnya | `(surogat pendataran)` | `konvensi` |

**Sebabnya satu dan mendasar: Pega tidak memberi kunci pada baris.** Halaman bersarang dan
page list tidak punya kolom pengenal sama sekali. Sebuah baris dikenali **posisinya** —
`pxSubscript`, yang muncul **695 kali di seluruh 59 berkas**, terbanyak kedua setelah
`pzInsKey`. Itulah sebabnya kedua belas PK di atas tidak dapat menunjuk apa pun di sistem
lama: yang mereka gantikan bukan kolom, melainkan urutan.

Foreign key mengikuti pembedaan yang sama, dan **dua bentuknya tidak boleh diratakan**:

| Bentuk induk | FK di sistem lama | Tabel |
|---|---|---|
| **page list** `[]` | `(posisi baris di daftar induk — pxSubscript)` | 6 tabel |
| **page bersarang** `{}` | `(halaman bersarang; tidak ada kolom penghubung)` | 6 tabel |

Yang pertama punya banyak baris berurutan; yang kedua **satu halaman yang menempel** pada
induknya dan tidak punya baris sama sekali. Menulis keduanya sebagai "FK" yang sama akan
menyembunyikan bahwa separuhnya bukan relasi, melainkan pengelompokan.

**Kolom `SUBSCRIPT` dipasang di keenam tabel `1:N`**, memuat `pxSubscript`. Tanpanya urutan
hilang saat didatarkan — dan urutan itu bukan hiasan: `KomiteRouter` langkah 6 beriterasi
menurut posisi dan berhenti di baris pertama ber-`KomiteAproval==0`. Giliran jenjang **adalah**
urutan baris.

**Satu-satunya tali lintas berkas yang berupa properti nyata** adalah
`pyWorkPage.pxCoverInsKey` — berkas sirkulasi ke klaim induknya, 6 rujukan di 3 berkas. Tali
ke usulan **tidak** berupa properti: ia posisional, lewat `Adjustment.IndexObject` (12 rujukan)
yang dipakai `KomitePostAdjustment` sebagai `Local.IdxParent` (24) dan `Local.IdxAdjustment`
(19). `.KomiteNo` — nomor sirkulasi yang ditulis balik ke klaim — muncul hanya **2 kali di 1
berkas**, dan `ADR-0031` sudah menetapkan ia **tidak dipercaya** sebagai tali.

## 5. Legenda kotak

| Penanda | Artinya |
|---|---|
| `PK` / `FK` | kunci primer / foreign key |
| `SHARED PK` | kunci primer yang sama dengan induknya — tanpa kolom FK terpisah |
| `*` di depan kolom | keyakinan **`rujukan-relatif`** — tidak pernah ditulis dengan jalur penuh |
| `n=` | cacah rujukan di 59 berkas |
| kotak **kuning** | **salinan** — rumahnya di sisi Klaim, ini potret beku |
| kotak **abu putus-putus** | **tipis** — nol atau satu kolom |

## 6. Temuan yang dibuat terlihat diagram ini

### 6.1 Enam dari empat belas tabel tidak boleh ikut dibuat

Delapan kotak bertanda **tipis** dan dua bertanda **salinan** lahir bukan karena modul ini menyimpan sesuatu, melainkan karena halaman kerja Komite **menyalin sebagian klaim ke dalam dirinya**. Mendatarkannya dengan jujur menghasilkan tabel seperti `T_KOMITE_QUOTATION` — satu kolom, `InsuredName` — dan `T_KOMITE_TREATY_MASTER` — satu kolom, `ReportingStart`.

Tabel satu kolom yang menyalin nilai dari tabel lain **bukan entitas**; ia jejak cara Pega memindahkan data antar halaman. Skema baru menetapkan tujuh objek, dan tidak satu pun dari keenam salinan itu ada di dalamnya — itu keputusan, bukan kelalaian.

### 6.2 Tiga tabel tidak punya satu kolom bisnis pun

`T_KOMITE_CLAIM_ADJ`, `T_KOMITE_CLAIM_OBJECT`, dan `T_KOMITE_OBJECT_ITEM` **tidak punya satu pun kolom bisnis yang pernah dirujuk**. Satu-satunya kolom mereka adalah `SUBSCRIPT` — posisinya sendiri. Ketiganya hanya dilewati untuk mencapai satu daun di bawahnya: `CurrencyID`, dua kali, di kedalaman 4 dan 5.

Artinya: modul ini menembus tiga tingkat bersarang **hanya untuk membaca satu pengenal mata uang**. Bila ketiganya dibuat sebagai tabel, hasilnya tiga tabel yang isinya kunci dan urutan — tanpa satu pun fakta.

### 6.3 Satu klaim, tiga salinan, tiga arah

`T_KOMITE_CLAIM_SNAPSHOT` adalah salinan **ketiga** dari klaim yang sama. Dua lainnya tidak digambar di sini karena bukan milik modul: `pyWorkCover.ClaimData` (64 simpul, jalur baca) dan `TempMainWork.ClaimData` (39 simpul, jalur tulis). Rinciannya di [`struktur-komite-lama.md`](struktur-komite-lama.md) bagian 5.

Memetakan ketiganya jadi tiga tabel akan melipatgandakan satu entitas; menggabungkannya tanpa mencatat **kapan** masing-masing dipakai akan menghapus fakta bahwa salinannya beku.

### 6.4 Dua penyimpan untuk satu fakta, dua kali

`T_GENERAL_KOMITE` membawa `KOMITE_COUNT` **dan** `KOMITE_LOOP` — dua bilangan yang menentukan hal yang sama dan tidak pernah saling memeriksa. `T_KOMITE_KOMITELIST` membawa `DATE_APPROVAL` **dan** `DATE_APPROVE` — dua kolom tanggal untuk satu keputusan.

Keduanya digambar apa adanya karena begitulah bentuk yang ada. Skema baru menghapus keduanya: hitungan jenjang menjadi turunan, dan waktu keputusan menjadi satu kolom.

### 6.5 `Adjustment` di kedalaman 5 bukan usulan

`T_KOMITE_OBJECT_ITEM_ADJ` berasal dari `ObjectList[].ObjectItemList[].Adjustment[]`. Namanya `Adjustment`, kelasnya bukan `Data-Adjustment`, dan isinya satu `CurrencyID` — ia **rincian mata uang per item objek**, bukan usulan yang diedarkan. Nama yang sama, arti berbeda, kedalaman berbeda.

Ini kembaran dari jebakan yang sudah dicatat sisi Klaim untuk `SpreadingAdjustment`.

### 6.6 Urutan jenjang adalah posisi baris, bukan sebuah nilai

`T_KOMITE_KOMITELIST` **tidak punya kolom derajat**. `DEGREE` hidup di `EMAILKOMITE`, di luar folder ini — nol kemunculan di 59 berkas. Yang ada hanyalah **posisi baris**, dan itu cukup bagi sistem lama: `KomiteRouter` langkah 6 beriterasi menurut posisi dan berhenti di baris pertama yang belum memutus.

Jadi derajat menentukan urutan **saat daftar dibentuk** (di folder Claim), lalu menghilang; sesudah itu urutan hidup sebagai posisi saja. Kolom `SUBSCRIPT` di §4.1 ada untuk menahannya — didatarkan tanpa itu, giliran jenjang hilang tanpa jejak.

## 7. Yang TIDAK digambar

- **Halaman induk** `pyWorkCover` (64 simpul) dan `TempMainWork` (39 simpul). Keduanya klaim induk, bukan milik modul ini; kontrak tulis-baliknya ada di `struktur-komite-lama.md` bagian 6 dan `datar-komite-tulis-balik.csv`.
- **Wadah muatan keluar** `InputParamOs` (17 field) dan `TempKasir` (25 field bernama `CARI…`). Keduanya bukan entitas, melainkan bentuk sesaat saat data keluar dari Pega. Isi `TempKasir` **berpagar** — `PG-04`.
- **Tipe data.** Tidak ada, dan tidak akan ada sampai `Rule-Obj-Property` diekspor.
- **Tabel Oracle sistem lama.** `EMAILKOMITE`, `OS_AKSEPTASI_KLAIM`, `DIRECTTOKASIR_LOG`, dan `PC_ASM_FW_GCNMFW_WORK` digambar di `claim-non-prop/4-erd-dan-tabel-datar/ERD-ORACLE.xlsx` sheet `ERD-LAMA`. Berkas ini menggambar **halaman Pega**, bukan tabel Oracle — pembedaan yang sama yang dipakai folder sisi Klaim.
- **Skema baru.** Tujuh objek `KLAIMNP` beserta 56 kolomnya ada di `SPEC-KOMITE-01.md` bagian *Model data*.

## 8. Batas

- **Nol foreign key diwarisi.** Sistem lama tidak menegakkan satu hubungan pun — sisi Klaim sudah mencatatnya untuk 32 tabel Oracle, dan hal yang sama berlaku di sini: keempat belas relasi di atas adalah **bentuk bersarang halaman Pega**, bukan constraint yang berdiri.
- **Empat kolom bertingkat `rujukan-relatif`.** Keempatnya di `T_KOMITE_KOMITELIST`. Induknya pasti karena repeat-nya bersumber pada `KomiteList`, tetapi jalurnya tidak pernah ditulis penuh.
- **Nol DDL dijalankan. Nol berkas sumber berubah.** Alat ini hanya membaca.
