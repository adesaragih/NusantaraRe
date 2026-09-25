# Inventaris berkas induk `treaty-in` — langkah 0.5 sesi to-spec

**Tanggal:** 24 September 2026 · **Perkakas:** `alat/inventaris-berkas.py`

> ### BATAS KEMAMPUAN perkakas ini — dibaca sebelum hasilnya
>
> Perkakas ini **tidak dapat menyatakan sebuah berkas SUDAH DIBACA.** Ia hanya dapat menyatakan
> sebuah berkas **DISEBUT** oleh berkas lain, atau **TIDAK PERNAH DISEBUT**.
>
> Penyebutan bukan pembacaan — sebuah berkas dapat disebut lima puluh kali tanpa pernah dibuka.
> Yang kuat adalah arah sebaliknya: **tidak pernah disebut** berarti tidak ada satu pun langkah
> yang merujuknya, dan itu calon terkuat untuk "belum pernah dibaca langkah mana pun".

## 1. Cacah

| | |
|---|---:|
| Berkas di `_migration-docs/treaty-in` | **99** |
| Korpus yang dipindai (treaty-in + treaty-in-adjustment + claim-non-prop) | **339** berkas teks |
| **DITOLAK — biner, tidak dipindai** | **7** |

**Yang ditolak, disebut namanya supaya penolakannya tidak diam:**
`Diagram-Skema-Tabel-TreatyMasuk.xlsx`, `Diagram-Skema-Tabel-ClaimNonProp.xlsx`,
`Diagram-Skema-Tabel-ClaimNonProp.BARU.xlsx`, `Diagram-Skema-Tabel-Gabungan.xlsx`,
`Diagram-Skema-Tabel-NusantaraRe.xlsx`, `ERD-ORACLE.xlsx`, `Rekap SQL.zip`.

Keenam `.xlsx` dan satu `.zip` **tidak pernah dibuka oleh langkah mana pun di proyek ini.**
Salah satunya bernama `ERD-ORACLE.xlsx` dan satu lagi `Rekap SQL.zip` — keduanya berjanji memuat
hal yang sesi DDL justru cari. **Dilaporkan, tidak dibuka di langkah ini**; masuk daftar sisa.

| Jenis | Jumlah |
|---|---:|
| `.md` | 59 |
| `.txt` | 13 |
| `.py` | 12 |
| `.csv` | 10 |
| `.html` | 2 |
| `.xlsx` · `.tsv` · `.sql` | 1 masing-masing |

## 2. Sapuan memakai TIGA bentuk penyebutan, bukan satu

Sapuan pertama memakai **satu** bentuk — nama berkas utuh — dan melaporkan **30 berkas yatim**,
di antaranya **seluruh 22 ADR**. Itu sampah, bukan temuan: ADR dirujuk sebagai `ADR-0055`, tidak
pernah sebagai `0055-daftar-keadaan-dan-perpindahan.md`. Aturan `CONTEXT.md` §2.0-j berlaku dan
bentuk kedua serta ketiga ditambahkan.

| Bentuk | Contoh | Untuk |
|---|---|---|
| **B1** nama berkas utuh | `SPEC-INVARIAN.md` | semua |
| **B2** penanda ADR | `ADR-0055` | `docs/adr/*.md` |
| **B3** nama objek tanpa ekstensi | `TREATYINDETAIL`, `PEGA_TREATY_IN` | `Table/`, `PROCEDURE/` |

Dengan ketiganya, yatim turun dari **30** menjadi **1**.

## 3. Yang TIDAK PERNAH DISEBUT — satu berkas

| Berkas | Ukuran | Apa ia |
|---|---:|---|
| `PROMPT-TO-TICKET-DAN-STRUKTUR-DATA.md` | 26.934 B | **prompt kerja**, bukan sumber fakta. Ia menyatakan dirinya begitu di kepalanya: *"Ini bukan spesifikasi, bukan tiket, dan bukan temuan."* |

**Bukan titik buta.** Ia berkas perintah, dan berkas perintah memang tidak dirujuk artefak.

## 4. Yang nyaris tidak pernah disebut — dan satu di antaranya membawa temuan

| Berkas | Sebutan | Catatan |
|---|---:|---|
| `4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` | **1** | **memuat pembacaan yang TERBALIK terhadap sumbernya** — lihat `TEMUAN-0-5-ACTUALVALUE-TERBALIK.md` |
| `4-erd-dan-tabel-datar/7-1-OLDID.md` | 2 | memperkenalkan **bentuk penulis kelima** (`CopyFrom`/`CopyInto`), yang tidak ada di `BENTUK-PENULIS-PROPERTI.md` |
| `alat/panen-kolom-dari-penulis.py` | 2 | perkakas yang langkah 2 sesi ini justru butuhkan |
| `alat/periksa-nama-kueri-dba.py` | 1 | |
| `ERD-TREATY-MASUK.md` | 2 | |
| `4-erd-dan-tabel-datar/PRA-PEMISAHAN-ANGKA-DAN-TABRAKAN.md` | 2 | memuat hasil sapuan `pxCreateSystemID` (L-9) |
| `docs/adr/0047-data-uji.md` | 1 | ADR yang paling sedikit dirujuk |
| `docs/adr/0050`, `0051` | 2 | |
| `4-erd-dan-tabel-datar/AUDIT-PENYEBUT-POHON.md` | 3 | semesta `PETA-TELUSUR-JSON` yang kurang |

## 5. `Table/` — tujuh DDL asli, dan apa yang TIDAK ada di sana

Ketujuh berkas `Table/*.txt` memuat `CREATE TABLE "POOLDATA"...` **asli**:

| Berkas | Tabel yang sebenarnya ada di dalamnya |
|---|---|
| `ACHIEVEMENT.txt` | `POOLDATA.ACHIEVEMENT` |
| `M_ATTACHMENTTREATY_2.txt` | `POOLDATA.M_ATTACHMENTTREATY_2` |
| **`PEGA_M_TREATY_IN_DETAIL.txt`** | **`POOLDATA.M_TREATY_IN`** — **nama berkas tidak sama dengan nama tabelnya** |
| `PROPORTIONALARRG.txt` | `POOLDATA.PROPORTIONALARRG` |
| `TREATYINDETAIL.txt` | `POOLDATA.TREATYINDETAIL` |
| `TREATYINPRODUCTION.txt` | `POOLDATA.TREATYINPRODUCTION` |
| `T_STORAGE_IMAGE.txt` | `POOLDATA.T_STORAGE_IMAGE` |

> **`TREATY_IN` dan `TREATY_IN_EDM` tidak ada di antaranya** — dan itu tepat dua tabel yang
> langkah 2 harus panen kolomnya dari badan procedure. Pernyataan itu **diperiksa, bukan
> diwarisi**: ketujuh berkas dibuka dan nama tabelnya dibaca satu per satu.

## 6. Sisa yang dilaporkan, tidak ditambal

| Sisa | Siapa menutup | Yang menagih |
|---|---|---|
| tujuh berkas biner belum pernah dibuka, termasuk `ERD-ORACLE.xlsx` dan `Rekap SQL.zip` | sesi berikutnya yang menyentuh DDL | gerbang sesi DDL: nama constraint dan presisi belum punya sumber, dan dua berkas ini menjanjikannya |
| `7-1-OLDID.md` memperkenalkan bentuk penulis **kelima** yang tidak masuk `BENTUK-PENULIS-PROPERTI.md` | sesi yang berikutnya memakai `sapu-penulis-properti.py` | **langkah 1 sesi ini** — sapuan lima bentuk yang bentuknya sendiri baru empat akan melaporkan nol yang salah |
