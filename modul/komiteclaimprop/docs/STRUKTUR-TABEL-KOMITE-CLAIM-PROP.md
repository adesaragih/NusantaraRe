# Struktur Tabel — Komite Claim Prop

Acuan bentuk tabel untuk aplikasi Go. Dibuat 2026-09-18 atas perintah work owner.
Berkas ini menggambarkan **BENTUK**, bukan alasan — alasannya ada di `spec.md` dan
`grilling-ronde-1.md` … `grilling-ronde-8.md`.

⚠️ **Berkas ini adalah acuan TUNGGAL nama kolom. Seluruh nama kolom snake_case**
`[keputusan work owner 2026-09-18]`. Ejaan yang muncul di `spec.md` dan berkas grilling adalah
**ejaan korpus** — itu **bukti asal kolom**, bukan nama kolom.

**Tipe ditulis sebagai kategori logis:** teks · angka desimal · bilangan bulat · DATE.
Seluruh uang, persen, dan kurs adalah **angka desimal**, tidak pernah float (**ADR-0003**).

**Sumber tiap kolom** adalah salah satu dari dua:

- **korpus** — kolom yang dipakai sistem berjalan, disertai nama rule dan nomor langkahnya
- **keputusan** — kolom yang sudah ditetapkan work owner, disertai tanggalnya

Rule korpus yang dirujuk di berkas ini:

| Singkatan | Rule | Kelas |
| --- | --- | --- |
| **AddKomiteChild** | `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT!ADDKOMITETREATYCHILD_ACT` |
| **KomitePost** | `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITEPOSTADJUSTMENT` |
| **KomiteRouter** | `Komite Claim Prop/Activity/KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITEROUTER` |
| **ShowTransfer** | `Komite Claim Prop/Section/ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY!SHOWTRANSFER` |
| **SetKomiteList** | `Komite Claim Prop/Activity/SetKomiteList_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY!SETKOMITELIST_ACT` |
| **Comitee** | kelas data daftar penyetuju `ASM-FW-GCNMFW-Data-Comitee` | — |

⚠️ Seluruh kolom **nullable** kecuali PK dan yang disebut NOT NULL; wajib-isi ditegakkan di Go.
⛔ **Tidak ada `CREATE TABLE` di berkas ini.** Bentuk kolom ditulis dalam kalimat.

---

## T_WORK_CLAIM — baris komite

⚠️ **Tabel ini BUKAN milik modul Komite dan tidak didefinisikan ulang di sini.** Ia tabel
**lintas-lini**, satu baris per work object. Definisi penuhnya ada di berkas struktur modul klaim.

Yang mengikat modul ini hanya **aturan baris komitenya**:

| Aturan | Isi | Sumber |
| --- | --- | --- |
| `ID` | berawalan **`TKMT-`** (teks berformat `TKMT-xxxxxx`) | **work owner 2026-09-18** |
| `COVER_KEY` | **wajib terisi** — `ID` baris klaim induknya (berawalan **`CLMP-`**) | **work owner 2026-09-18** |
| `LINI` | **`PROP`** — sama dengan baris klaim induknya | **work owner 2026-09-18** |

Baris komite **lahir saat penyerahan ke Komite diterima**; sebelum itu ia tidak ada.

⚠️ **`T_WORK_CLAIM` lintas-lini**, jadi definisi kolom `ID`-nya menyebut **delapan awalan** —
sepasang per lini:

| LINI | baris klaim | baris komite |
| --- | --- | --- |
| **FAC** | `CLM-` | `KMT-` |
| **PROP** | `CLMP-` | **`TKMT-`** |
| **NONPROP** | `CLMNP-` | `KMTNP-` |
| **LIFE** | `CLMLF-` | `KMTLF-` |

Berkas ini modul **Prop**, jadi contoh dan ilustrasinya memakai pasangan Prop saja
(`CLMP-` / `TKMT-`).

`T_GENERAL_KOMITE` memakai **`ID` yang sama persis** dengan baris komite ini — **shared primary
key**, tanpa kolom penyambung.

---

## T_GENERAL_KOMITE

Header kasus komite. Satu baris mewakili **satu kasus komite**, yaitu **satu penyerahan sebuah
baris penyesuaian ke Komite**. `ID`-nya **sama persis** dengan baris komite di `T_WORK_CLAIM`
(shared PK).

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | **PK**, shared dengan `T_WORK_CLAIM` · berawalan `TKMT-` | keputusan work owner 2026-09-18 |
| `ADJUSTMENT_ID` | teks | **tidak** | penunjuk ke **baris penyesuaian induk** di Claim Prop · index **biasa** (`IX_GENERAL_KOMITE_ADJ`, migrasi 681 — RALAT 08-10-2026) | keputusan work owner 2026-09-18 · 2026-10-08 |
| `KOMITE_LOOP` | **bilangan bulat** | ya | berapa penyetuju yang dibutuhkan | korpus — **AddKomiteChild**, saat kasus dibuat |
| `KOMITE_COUNT` | **bilangan bulat** | ya | penyetuju ke berapa yang sedang berjalan | korpus — **AddKomiteChild** langkah awal; dinaikkan **KomiteRouter** langkah 5 dan **KomitePost** langkah 40 |
| `ACCEPT_STATUS` | teks | ya | hasil akseptasi komite terakhir · `1` setuju · `2` tolak | korpus — **KomitePost** |
| `KOMITE_USUL_TUTUP` | teks | **tidak** | **usul tutup klaim** dari penyetuju · penanda `'1'`/`'0'`, bawaan `'0'` (migrasi 680) | keputusan work owner **2026-09-19** · korpus — **ShowTransfer** |
| `KOMITE_USUL_CADANG` | teks | **tidak** | **usul cadangkan klaim** dari penyetuju · penanda `'1'`/`'0'`, bawaan `'0'` (migrasi 680) | keputusan work owner **2026-09-19** · korpus — **ShowTransfer** |
| `KOMITE_SUBJECTIVITY` | teks | **tidak** | isian "Subjectivity ?" tingkat 1 disimpan antar tingkat · `'1'`/`'0'`, bawaan `'0'` (migrasi 682) | keputusan work owner **2026-10-08** (OQ-KCP-01 "a") · korpus — **ShowTransfer** `.IsSubjectivity`, KomitePostAdjustment S16-S34 |
| `KOMITE_SUBJECTIVITY_NOTE` | teks | ya | isian "Subjectivity Note" tingkat 1 (lebar 1000 = `T_CLAIM_ADJUSTMENT.SUBJECTIVITY_NOTE`, migrasi 682) | keputusan work owner **2026-10-08** · korpus — **ShowTransfer** `.SubjectivityNote`, S24 |
| `TRANSFER_TYPE` | teks | **tidak** | jenis penyerahan · `'2'` adjustment / `'3'` Reject Claim / `'4'` Close Without Payment, bawaan `'2'` (migrasi `komiteclaimfacin/642`); kasus Prop tidak menulisnya | keputusan work owner **2026-10-10** (KCF-03) · korpus Komite Claim FacIn — **KomitePostAct** S2 / S4 / S5 |
| `KOMITE_CIRCUM_CAUSE_OF_LOSS` | teks | ya | teks pop-up Chronology kasus komite Fac In TT3 / TT4 (lebar 4000 = baris adjustment `KOMITE_CIRCUM_CAUSE_OF_LOSS`); kasus Prop tidak menulisnya | keputusan work owner **2026-10-10** (OQ-KCFI-03) · korpus Claim Fac In — **SendRejectClaimToKomite2** / **SendCloseClaimToKomite** 7.2 |
| `KOMITE_EXTENT_OF_LOSS` | teks | ya | teks pop-up Extent Of Loss kasus komite Fac In TT3 / TT4 (lebar 4000 = baris adjustment `KOMITE_EXTENT_OF_LOSS`); kasus Prop tidak menulisnya | keputusan work owner **2026-10-10** (OQ-KCFI-03) · korpus Claim Fac In — **SendRejectClaimToKomite2** / **SendCloseClaimToKomite** 7.2 |
| `KOMITE_LEGAL_LIABILITY` | teks | ya | teks pop-up Policy Liability kasus komite Fac In TT3 / TT4 (lebar 4000 = baris adjustment `KOMITE_LEGAL_LIABILITY`); kasus Prop tidak menulisnya | keputusan work owner **2026-10-10** (OQ-KCFI-03) · korpus Claim Fac In — **SendRejectClaimToKomite2** / **SendCloseClaimToKomite** 7.2 |

> **RALAT 10-10-2026** (izin work owner 10-10-2026, Komite Claim Fac In tahap 2, keputusan KCF-03). `ADJUSTMENT_ID` kini **boleh kosong** (migrasi `komiteclaimfacin/641`, `MODIFY ... NULL`): kasus komite Fac In TT3 Reject Claim dan TT4 Close Without Payment (`T_WORK_CLAIM.LINI = 'FACIN'`, awalan `KMT-`) lahir tanpa baris adjustment; kasus Life / Prop / Non Prop tetap mengisinya. Kolom `TRANSFER_TYPE` ditambahkan migrasi `komiteclaimfacin/642` (`CHAR(1) DEFAULT '2' NOT NULL`, CHECK `'2'`/`'3'`/`'4'` = jenis penyerahan `pyWorkPage.TransferType`); baris lini lain bernilai bawaan `'2'` (penyerahan adjustment). Mundur 641 = `NOT NULL NOVALIDATE` (baris kosong yang ada dibiarkan).

> **RALAT 10-10-2026 (migrasi 643)** (jawaban work owner 10-10-2026, Komite Claim Fac In, OQ-KCFI-03). Kolom `KOMITE_CIRCUM_CAUSE_OF_LOSS` / `KOMITE_EXTENT_OF_LOSS` / `KOMITE_LEGAL_LIABILITY` (`VARCHAR2(4000)`, boleh kosong, ADD saja) ditambahkan migrasi `komiteclaimfacin/643`: teks pop-up Chronology / Extent Of Loss / Policy Liability kasus komite Fac In TT3 Reject Claim / TT4 Close Without Payment (`SendRejectClaimToKomite2` / `SendCloseClaimToKomite` 7.2 -> `childPageKomite.Komite.*`, Section `ShowTransfer` LS39). Kasus lini lain dan TT2 tidak menulisnya (TT2 menyimpan teks yang sama di baris adjustment).

> **RALAT 08-10-2026** (prompt implementasi §7 butir 2). Kalimat lama baris `ADJUSTMENT_ID`: *"index **UNIK**"* — diganti indeks biasa (keputusan work owner 08-10-2026: ID adjustment Prop `SEQ_T_CLAIM` dan Life `SEQ_CLAIMLF_ADJ` sama-sama angka polos; keunikan satu adjustment ↔ satu kasus komite dijaga `KOMITE_ID UNIQUE` di `T_CLAIM_ADJUSTMENT` dan `T_CLAIMLF_ADJUSTMENT`). Kalimat lama dua kolom usul: *"teks *(penanda)* | ya"* — bertentangan dengan keputusan 21 (2026-09-19: wajib isi, hanya `'1'`/`'0'`); kini `CHAR(1) DEFAULT '0' NOT NULL` + `CHECK`. Ketiga tabel komite sudah ada (`claimlife/013`, `komiteclaimlife/030`); modul ini hanya `ADD` dua kolom usul (680) dan mengganti indeks (681).

⚠️ `KOMITE_LOOP` dan `KOMITE_COUNT` adalah **pencacah** — **bilangan bulat biasa**, **tidak** ikut
aturan ketelitian angka (lihat §Ketelitian angka).

`[keputusan work owner]` 2026-09-18 — **`KOMITE_LOOP` adalah kolom yang DISIMPAN**, diisi ulang di
titik yang sama dengan Pega mengisinya. **Bukan `COUNT(*)` setiap kali dibaca.**

---

## T_KOMITE_KOMITELIST

Satu baris per **penyetuju** dalam satu kasus komite. Inilah tangga persetujuannya.

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | **PK** | keputusan work owner |
| `DATA_KOMITE_ID` | teks | **tidak** | FK → `T_GENERAL_KOMITE.ID` | keputusan work owner |
| `KOMITE_URUT` | **bilangan bulat** | ya | penyetuju ke berapa | keputusan work owner |
| `KOMITE_OPERATORID` | teks | ya | **akun operator** penyetuju | korpus — ejaan Pega `KomiteID` |
| `KOMITE_JABATAN` | teks | ya | **jabatan** penyetuju | korpus — ejaan Pega `IDKomite` |
| `KOMITE_EMAIL` | teks | ya | email penyetuju | korpus — ejaan Pega `KomiteEmail` |
| `KOMITE_APPROVAL` | teks | ya | keputusan · `0` belum · `1` setuju · `2` tolak | korpus — ejaan Pega `KomiteAproval`, **ejaan dibetulkan** |
| `KOMITE_COMMENT` | teks | ya | catatan penyetuju | korpus — ejaan Pega `KomiteComment` |
| `DATE_APPROVE` | DATE | ya | tanggal diputuskan | korpus — ejaan Pega `DateApprove` |

`[keputusan work owner]` 2026-09-18 — **satu kolom tanggal saja.** Di Pega ada properti tanggal
kedua yang diisi **langkah yang sama** dengan nilai identik, dan **tidak pernah ditampilkan**. Ia
**sengaja tidak dijadikan kolom** — jangan menambahkannya kembali karena "ada di korpus".

### Kesembilan kolom dicocokkan dengan properti Pega yang mengisinya

`[terverifikasi]` Sensus penulis se-korpus:

| Kolom | Properti Pega | Penulisan | Pasangannya jelas? |
| --- | --- | --- | --- |
| `KOMITE_OPERATORID` | `KomiteID` | **22** | ✅ jelas |
| `KOMITE_JABATAN` | `IDKomite` | **37** | ✅ jelas |
| `KOMITE_EMAIL` | `KomiteEmail` | **20** | ✅ jelas |
| `KOMITE_APPROVAL` | `KomiteAproval` | **54** | ✅ jelas |
| `KOMITE_COMMENT` | `KomiteComment` | **30** | ✅ jelas |
| `DATE_APPROVE` | `DateApprove` | **30** | ✅ jelas |
| `ID` | — | — | ⚠️ **struktural** — tidak ada padanan Pega |
| `DATA_KOMITE_ID` | — | — | ⚠️ **struktural** — tidak ada padanan Pega |
| `KOMITE_URUT` | — | — | ⚠️ **struktural** — di Pega ia **posisi baris di daftar**, bukan properti tersimpan |

⭐ **Enam dari sembilan punya pasangan jelas. Tiga struktural.**

⚠️ `KOMITE_URUT` perlu perhatian: di Pega urutan penyetuju adalah **posisi dalam daftar**, dan
`KOMITE_COUNT` menunjuk posisi itu. Menjadikannya kolom berarti urutan **menjadi eksplisit**, bukan
tersirat dari susunan baris.

---

## Baris penyesuaian induk — **DIBACA, tidak disalin**

`[keputusan work owner]` 2026-09-18 — sembilan belas properti baris penyesuaian dipilah dua:

| Berapa | Perlakuan |
| --- | --- |
| **17** yang **hanya dibaca** | **DIBACA dari baris penyesuaian milik Claim Prop** — tidak disalin |
| **2** yang **disunting penyetuju** | **DISIMPAN di tabel komite** — lihat §berikutnya |

### Ketujuh belas yang dibaca

`[terverifikasi]` Seluruhnya hanya-baca di layar komite, dan bergerbang **jalur penyesuaian**
(kecuali dua yang tanpa gerbang):

nilai penyesuaian · nilai kotor penyesuaian · nilai kotor · mata uang · kode mata uang ·
jenis penyesuaian · persen RNM · nama bank · cabang bank · nomor rekening · kode swift ·
dibayarkan kepada · dapat dibayar · jenis risiko perorangan · persen risiko perorangan ·
nilai risiko perorangan · nilai penyesuaian yang diusulkan.

### ⭐ Jalur bacanya — dari kasus komite sampai ke baris penyesuaian

`[terverifikasi]` **Tiga langkah**, dan kuncinya **dua bagian**:

```
1.  kasus komite  ->  kunci kasus klaim induk
       bukti: KomitePost langkah 4 membuka kasus klaim dengan kunci induk
              (Obj-Open-By-Handle, dikunci, dilepas saat simpan)

2.  kasus komite  ->  POSISI baris penyesuaian di dalam klaim itu
       bukti: KomitePost langkah 3  ->  penanda posisi disalin ke variabel kerja
              asalnya AddKomiteChild langkah 13, yang mengisinya dari nomor urut
              baris saat kasus komite dibuat

3.  baris penyesuaian  =  daftar penyesuaian milik klaim itu, pada posisi tersebut
       bukti: KomitePost langkah 6 · 7 · 15 · 16.9 · 23 · 24 · 25 · 27 seluruhnya
              menulis balik ke daftar penyesuaian pada posisi itu
```

> ## ⚠️ KUNCINYA POSISIONAL, BUKAN IDENTITAS
>
> `[terverifikasi]` Penghubungnya adalah **nomor urut baris di dalam daftar**, bukan sebuah
> pengenal baris. Artinya: **bila baris penyesuaian disisipkan atau dihapus di klaim induk,
> penunjuk milik kasus komite menunjuk baris yang salah** — tanpa ada yang menyadarinya.
>
> ⛔ **Tidak ada penanganan untuk hal itu di Pega** — tidak ada pemeriksaan, tidak ada pesan.
> `[terbuka]`

⚠️ **Apa yang terjadi bila baris induknya tidak ketemu:** `[terbuka]` — **tidak ada
penanganannya di korpus**. Langkah-langkah di atas menulis ke posisi itu **tanpa memeriksa lebih
dulu** apakah barisnya ada.

⭐ **Untuk Go, kolom `ADJUSTMENT_ID` di `T_GENERAL_KOMITE` menggantikan pasangan kunci itu** —
itulah gunanya ia ada, dan itulah sebabnya ia **NOT NULL** dan ber-index **UNIK**.

---

## Dua kolom baru — usul tutup dan usul cadangkan

`[keputusan work owner]` 2026-09-18 — kedua penanda ini **DISIMPAN di tabel komite**, karena
**penyetuju yang mengisinya**.

`[terverifikasi]` Sifatnya di layar: **kotak centang**, **dapat disunting**, **tidak wajib isi**,
bergerbang famili A **jalur penyesuaian** — bukti **ShowTransfer**. **Nol penulis di activity mana
pun**: keduanya ditulis **langsung dari layar**, tanpa perantara.

### ⚠️ Di Pega keduanya punya DUA tempat

`[terverifikasi]` Sesudah penyetuju mengisinya, **KomitePost langkah 11** — *"add propose close and
propose reserve"*, **tanpa gerbang** — menyalin keduanya ke **kasus klaim induk** dengan nama lain:

```
kasus klaim . penanda tutup berkas       <-  kasus komite . usul tutup klaim
kasus klaim . penanda klaim dicadangkan  <-  kasus komite . usul cadangkan klaim
```

⛔ **Dicatat sebagai temuan.** Apakah sisi klaim ikut disimpan adalah **urusan modul Claim Prop**,
bukan berkas ini.

> ## ⛔ RALAT — *"tidak ada tulis-menulis lintas modul"* hanya benar untuk LAYAR
>
> Catatan lama menyebut bahwa isi baris penyesuaian sampai ke kasus komite lewat **penyalinan
> saat kasus dibuat**, dan **"tidak ada tulis-menulis lintas modul"**. `[terverifikasi]` Kalimat
> itu **hanya benar untuk layar**.
>
> **Lewat activity, komite MEMANG menulis ke kasus klaim induk.** Buktinya
> **KomitePostAdjustment langkah 11** — bukan `KomitePost`, yang hanya penyalur tiga langkah —
> menyalin kedua penanda usul ke kasus klaim **dengan nama berbeda**, dan langkah itu **tanpa
> gerbang**:
>
> ```
> kasus komite . usul tutup klaim       ->  kasus klaim . penanda tutup berkas
> kasus komite . usul cadangkan klaim   ->  kasus klaim . penanda klaim dicadangkan
> ```
>
> ⭐ **Jadi nilai yang sama hidup di DUA tempat dengan dua nama.** Itu fakta Pega, bukan
> rancangan.

### ✅ Tabel yang memuatnya: **`T_GENERAL_KOMITE`** — **DITUTUP**

`[terverifikasi]` **Bukti B1: halaman penyesuaian pada kasus komite adalah SATU, bukan daftar.**
Dasarnya tiga:

| Bukti | Isinya |
| --- | --- |
| **AddKomiteChild** langkah 14 | mengisi **tanpa indeks** — `Adjustment.«kolom» := «nilai»` |
| **ShowTransfer** | membacanya **tanpa indeks**, 26 kali |
| deklarasi halaman | dideklarasikan sebagai **satu halaman** berkelas `Data-Adjustment`, bukan daftar |

⚠️ Dua rujukan **berindeks** yang tampak mirip ternyata **daftar lain yang lebih dalam** milik
kasus klaim — bukan halaman penyesuaian milik komite.

`[keputusan work owner]` **2026-09-19** — **pilihan A: kedua kolom disimpan di
`T_GENERAL_KOMITE`.** **Butir terbukanya TUTUP.**

**Alasan work owner, ditulis apa adanya:** *komite menyimpan catatan usulnya sendiri, terpisah
dari nilai akhir yang tercatat di kasus klaim.*

Cocok dengan bukti B1: satu kasus komite = satu baris penyesuaian, jadi header kasus adalah
tempat yang tepat.

**Usulan penamaan** — *ini penamaan, bukan rancangan*, dan mengikuti pola nama yang sudah dipakai
di berkas ini (`KOMITE_` + kata, snake_case, huruf besar):

| Untuk | Usulan nama | Tipe |
| --- | --- | --- |
| usul tutup klaim | `KOMITE_USUL_TUTUP` | teks *(penanda)* |
| usul cadangkan klaim | `KOMITE_USUL_CADANG` | teks *(penanda)* |

⚠️ **Nilainya `[terbuka]`** — di Pega ia kotak centang, tetapi **nilai yang tersimpan tidak terbaca
dari korpus** (tidak ada nilai tercentang/tidak-tercentang yang dinyatakan di berkas layar).

---

## Ketelitian angka

`[keputusan work owner]` 2026-09-18:

| Golongan kolom | Bentuknya |
| --- | --- |
| **Angka uang, persen, dan kurs** | **mengikuti bentuk yang sudah ada di produksi — 20 digit seluruhnya, 8 di antaranya di belakang koma** |
| **Pencacah** — `KOMITE_LOOP` · `KOMITE_COUNT` · `KOMITE_URUT` | **bilangan bulat biasa**, tidak ikut aturan di atas |
| **Tipe di Go** | **tipe desimal, BUKAN bilangan pecahan biner** `[keputusan work owner — atas rekomendasi asisten]` |
| **Hitungan di Go** | sampai **20 angka di belakang koma**, **nol pembulatan di tengah jalan** |
| **Tampilan layar** | **4 angka di belakang koma** |

⚠️ **Pembulatan terjadi di batas penyimpanan, dan itu diterima sadar.** Angka persen yang di Pega
dihitung sampai 10 desimal **akan menjadi 8** saat disimpan.

> **RALAT 08-10-2026** (prompt §7 butir 2). Kalimat lama *"mengikuti bentuk yang sudah ada di produksi — 20 digit
> seluruhnya, 8 di antaranya di belakang koma"* — kolom uang `T_CLAIM_*` yang ada di DEV `NUMBER(38,10)` (keputusan
> work owner 07-10-2026). Bagan relasi di bawah yang menulis `ADJUSTMENT_ID` *"index UNIK"* kini indeks biasa
> (migrasi 681).

⛔ **Di tabel-tabel berkas ini tidak ada satu pun kolom angka uang** — ketiga tabel di atas hanya
memuat teks, tanggal, dan pencacah. Aturan ini berlaku bagi kolom uang di **tabel penyesuaian milik
Claim Prop**, yang dibaca modul ini.

---

## Pohon relasi — sisi komite

```
T_WORK_CLAIM  (lintas-lini, satu baris per work object)
   |
   |  baris klaim   ID = CLMP-xxxxxx   LINI = PROP   COVER_KEY kosong
   |  baris komite  ID = TKMT-xxxxxx   LINI = PROP   COVER_KEY = CLMP-xxxxxx
   |
   +-- T_GENERAL_KOMITE        ID = TKMT-xxxxxx   (shared PK, tanpa kolom penyambung)
   |      |
   |      |  ADJUSTMENT_ID -> baris penyesuaian milik Claim Prop
   |      |                   NOT NULL, index UNIK, 1:1
   |      |
   |      +-- T_KOMITE_KOMITELIST     DATA_KOMITE_ID -> T_GENERAL_KOMITE.ID
   |             satu baris per penyetuju, 1:banyak
   |
   +-- (17 properti penyesuaian DIBACA dari Claim Prop, tidak disalin)
```

---

## Tiga keputusan **2026-09-19** — dan butir yang ditutupnya

### ✅ Tabel dua kolom usul — **DITUTUP**

`[keputusan work owner]` 2026-09-19 — **`T_GENERAL_KOMITE`**. Lihat § di atas.

### ✅ Batas 12 digit di depan koma — **DITUTUP**

`[keputusan work owner]` 2026-09-19 — **ikuti bentuk kolom yang sudah ada apa adanya. Tidak
ditanyakan ke DBA.**

> ⚠️ **Risiko diterima sadar:** bila sebuah nilai melewati batas itu, basis data **menolak
> menyimpan** — bukan membulatkan. Penyimpanan **gagal**, dan kegagalan itu harus terlihat.
> Batas ini **tidak dilebarkan**.

### ✅ Migrasi penunjuk posisional — **DITUTUP**

`[keputusan work owner]` 2026-09-19 — **penunjuk posisional ke baris penyesuaian DIPINDAHKAN APA
ADANYA, termasuk yang sudah salah alamat di Pega. Tidak diperiksa dulu, tidak diperbaiki.**
Konsisten dengan preseden **"tiru apa adanya"**.

> ⚠️ **Risikonya, ditulis tegas:** **penunjuk yang salah di Pega akan TETAP SALAH sesudah
> migrasi.** Bila sebuah baris penyesuaian pernah disisipkan atau dihapus di klaim induk sesudah
> kasus komitenya dibuat, kasus komite itu menunjuk baris yang keliru — dan sistem baru akan
> **mewarisi kekeliruan itu apa adanya**, tanpa penanda.

---

## Butir `[terbuka]` berkas ini — **dua**

⛔ **Tidak satu pun ditutup selain ketiga di atas.**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | **Nilai apa yang tersimpan** untuk kedua penanda usul — tidak terbaca dari korpus | korpus / work owner |
| 2 | **Apa yang terjadi bila baris penyesuaian induk tidak ketemu** *(saat berjalan, bukan saat migrasi)* — tidak ada penanganannya di korpus | work owner |

⚠️ Butir 2 **tidak menghambat berkas ini** — `ADJUSTMENT_ID` sudah menggantikan kunci posisional
untuk data baru. Ia menggigit pada **data hasil migrasi**, yang menurut keputusan di atas
membawa penunjuk lama apa adanya.
