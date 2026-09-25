# ERD TABEL Oracle — sistem lama dan skema baru, dalam Excel

<!-- STEMPEL ASAL -->
> **Dibangkitkan** `alat/buat-erd-excel.py` pada **19 September 2026**.
>
> **Yang digambar: TABEL Oracle saja.** VIEW, PROCEDURE, FUNCTION, dan halaman clipboard Pega tidak digambar — yang terakhir itulah cacat bentuk sebelumnya, dan bentuk itu sekarang ada di `sampah/`.
>
> **Dua sistem, sheet terpisah, tidak dicampur.** Yang lama **tidak punya satu pun foreign key**; yang baru punya 17 dan DDL-nya **belum pernah dijalankan**.
>
> **TURUNAN.** Sumber dibaca langsung. Bila workbook berbeda dari sumbernya, **alatnya yang salah**.

### → [`ERD-ORACLE.xlsx`](ERD-ORACLE.xlsx)

```
python alat/buat-erd-excel.py
```

---

## 1. Sepuluh sheet

| Sheet | Isi |
|---|---|
| `BACA-DULU` | stempel asal, legenda, dua belas pencacah **berupa rumus**, batas bukti |
| `KELOMPOK` | tujuh kelompok domain dan warna kepala kotaknya |
| **`ERD-LAMA`** | **32 kotak entitas** — TABEL Oracle sistem lama, kolomnya di dalam kotak |
| `KOLOM-LAMA` | 627 kolom, satu baris per kolom, berpenyaring |
| `RELASI-LAMA` | 15 relasi **tersirat**, beserta berkas XML yang membuktikannya |
| **`ERD-BARU`** | **22 kotak entitas** — TABEL usulan `KLAIMNP` |
| `ERD-BARU-DAFTAR` | daftar tabel skema baru, cacah kolom, PK, cacah FK keluar |
| `KOLOM-BARU` | 396 kolom, berpenyaring |
| `RELASI-BARU` | 17 foreign key, tiap satunya bernama constraint |
| `CONSTRAINT-BARU` | 60 `CHECK` + 16 `UNIQUE` |

Tata letak seragam di seluruh sheet tabel: judul baris 1, catatan baris 2, kepala kolom baris 4, **data mulai baris 5**. Keseragaman itu bukan kerapian — rumus pencacah di `BACA-DULU` bersandar padanya.

## 2. Sumber

| Sumber | Dipakai untuk |
|---|---|
| `pengetahuan/SCHEMA-ACTUAL.csv` | 32 TABEL sistem lama beserta 627 kolomnya |
| ekspor XML `D:\XML_NURE\Claim Non Prop` | 50 blok `pyBrowseSQL` → 15 relasi tersirat |
| `ddl-usulan/*.sql` | 22 TABEL skema baru, 396 kolom, 17 FK, 76 constraint |

Ketiganya dibaca **langsung**. Tidak ada berkas turunan yang dipakai sebagai sumber.

## 3. Cacah, beserta perintahnya

Kedua belas angka ini **tidak diketik ke dalam sel** — tiap satunya rumus yang menghitung dari sheet lain di workbook yang sama. Nilai yang seharusnya keluar dicetak alat ini setiap kali ia jalan, jadi isi sel dapat diadu dengan hitungan mandiri.

| Angka | Nilai | Perintah pembanding |
|---|---:|---|
| Tabel sistem lama | **32** | `awk -F';' 'NR>1 && $5=="TABLE"{print $4}' pengetahuan/SCHEMA-ACTUAL.csv \| sort -u \| wc -l` |
| Kolom sistem lama | **627** | `awk -F';' 'NR>1 && $5=="TABLE"' pengetahuan/SCHEMA-ACTUAL.csv \| wc -l` |
| Kolom PK sistem lama | **15** | `awk -F';' 'NR>1 && $5=="TABLE" && $14=="Y"' pengetahuan/SCHEMA-ACTUAL.csv \| wc -l` |
| **Foreign key sistem lama** | **0** | `grep -c "FOREIGN KEY" pengetahuan/ddl/TABLE_*.sql \| awk -F: '{s+=$2} END{print s}'` |
| Relasi tersirat sistem lama | **15** | keluaran alat |
| Tabel skema baru | **22** | `grep -l "CREATE TABLE KLAIMNP" ddl-usulan/*.sql \| wc -l` |
| Kolom skema baru | **396** | keluaran alat |
| Kolom `NOT NULL` skema baru | **178** | keluaran alat |
| Kolom turunan `VIRTUAL` | **5** | keluaran alat |
| Foreign key skema baru | **17** | `grep -c "FOREIGN KEY" ddl-usulan/*.sql \| awk -F: '{s+=$2} END{print s}'` |
| `CHECK` constraint skema baru | **60** | keluaran alat |
| `UNIQUE` constraint skema baru | **16** | keluaran alat |

Enam angka skema baru cocok dengan yang dihitung pengurai sebelumnya (22 / 396 / 178 / 5 / 17 / 60). Dua pengurai berbeda membaca satu sumber yang sama dan sepakat — itu yang membuat angkanya boleh dipegang.

## 4. Legenda kotak

| Penanda | Artinya |
|---|---|
| `PK` | kunci primer |
| `UQ` | ikut kunci unik |
| `FK` | foreign key **yang tertulis di DDL** — hanya ada di skema baru |
| `~` | relasi **tersirat** — tidak ditegakkan basis data; buktinya SQL di XML |
| nama kolom **tebal** | `NOT NULL` |
| warna kepala kotak | kelompok domain, daftarnya di sheet `KELOMPOK` |

Kotak disusun ke dalam empat jalur berdampingan dengan tinggi diseimbangkan, kotak terbesar lebih dulu. `PC_ASM_FW_GCNMFW_WORK` (100 kolom), `OS_AKSEPTASI_KLAIM` (67), `TREATYINDETAIL` (64), dan `TREATYINPRODUCTION` (60) karena itu berdiri di puncak masing-masing jalur.

## 5. Temuan yang dibuat terlihat diagram ini

### 5.1 Sistem lama tidak menegakkan satu hubungan pun

**Nol `FOREIGN KEY` di seluruh 32 tabel.** Ini bukan kelalaian pencatatan dan bukan longgar-longgar saja: tidak ada satu pun hubungan antar tabel yang dijaga basis data. Seluruh keterkaitan hidup di dalam kode — di SQL yang ditulis tangan, dan di aktivitas Pega.

Akibatnya untuk migrasi: **tidak ada integritas rujukan yang bisa diwarisi.** Setiap hubungan di skema baru adalah hubungan yang **ditetapkan sekarang**, bukan yang dipindahkan. Data lama boleh jadi memuat baris yatim, dan tidak ada apa pun di sistem lama yang mencegahnya.

### 5.2 Kelima belas relasi tersirat, dan dari mana dibacanya

Tidak satu pun berasal dari kesamaan nama kolom. Semuanya dibaca dari **comma-join** atau **subquery berkorelasi** di SQL yang benar-benar dijalankan sistem lama.

| Kiri | Kolom | Kanan | Kolom | Bentuk | Keadaan |
|---|---|---|---|---|---|
| `JSON_KLAIM` | `IDPEGA` | `CLAIMREJECTED` | `INSKEY` | comma-join | kedua ujung TABEL |
| `M_CLIENT` | `ID` | `AGENT` | `CLIENTID` | subquery | kedua ujung TABEL |
| `PROPORTIONALARRG` | `REINSTYPEID` | `TREATYBUSINESS` | `REINSTYPEID` | subquery | kedua ujung TABEL |
| `PROPORTIONALARRG` | `REINSTYPEID` | `TREATYCONTRACT` | `REINSTYPEID` | subquery | kedua ujung TABEL |
| `PROPORTIONALARRG` | `TREATYGROUPID` | `TREATYBUSINESS` | `TREATYGROUPID` | subquery | kedua ujung TABEL |
| `PROPORTIONALARRG` | `TREATYYEARID` | `TREATYCONTRACT` | `IDTREATYYEAR` | subquery | kedua ujung TABEL |
| `REINSURANCETYPE` | `ID` | `PROPORTIONALARRG` | `REINSTYPEID` | comma-join | kedua ujung TABEL |
| `TREATYINDETAIL` | `TREATYGROUPID` | `TREATYBUSINESS` | `TREATYGROUPID` | comma-join | kedua ujung TABEL |
| `TREATYINDETAILEDM` | `TREATYGROUPID` | `TREATYBUSINESS` | `TREATYGROUPID` | comma-join | kedua ujung TABEL |
| `TREATYINPRODUCTION` | `MARKETINGOFFICERCODE` | `MARKETINGOFFICER` | `ID` | comma-join | kedua ujung TABEL |
| `CITY` | `ID` | `RW` | `CITYID` | comma-join | `CITY` sebuah VIEW |
| `CLAIMXOL` | `CASEID` | `CLAIMXOL2` | `CASEID` | comma-join | `CLAIMXOL` sebuah VIEW |
| `CLAIMXOL` | `CASEID` | `OS_AKSEPTASI_KLAIM` | `CASEID` | subquery | `CLAIMXOL` sebuah VIEW |
| `V_D_CAUSE_OF_LOSS_BUSINESS` | `BISNISID` | `BUSINESS` | `ID` | comma-join | ujung kiri sebuah VIEW |
| `TREATYBUSINESS` | `TREATYGROUPID` | `TREATYYEAR` | `TREATYGROUPID` | comma-join | **`TREATYYEAR` tidak punya DDL** |

Sepuluh relasi berujung TABEL di kedua sisinya; empat menyentuh VIEW; satu menyentuh `TREATYYEAR`, tabel yang **disebut SQL produksi tetapi belum ada DDL-nya** — ia salah satu dari 25 objek di `pengetahuan/PULL-LIST.csv` yang masih menunggu tarikan.

### 5.3 Nol `JOIN` berkata kunci di seluruh 50 blok SQL

Sistem lama tidak menulis `INNER JOIN` / `LEFT JOIN` satu kali pun. Yang dipakai comma-join gaya lama dan subquery. Sapuan yang mencari kata `JOIN` akan melaporkan nol relasi dan menyimpulkan tabel-tabelnya berdiri sendiri — kesimpulan yang salah. Ini tercatat di sini supaya sapuan berikutnya tidak mengulanginya.

## 6. Yang TIDAK digambar

- **VIEW** — 10 di `pengetahuan/ddl/`, 9 di skema baru. Ia bukan tabel. Empat di antaranya muncul di `RELASI-LAMA` karena SQL memang menjoin lewatnya, dan di sana ia ditandai pada kolom `KEADAAN`.
- **PROCEDURE (6) dan FUNCTION (1)** — logika, bukan bentuk data. Dua di antaranya memuat logika bisnis yang tidak terbaca dari mana pun selain badan prosedurnya.
- **Halaman clipboard Pega** — ia bukan tabel Oracle. Inilah yang salah pada bentuk sebelumnya.
- **Index, sequence, hak akses** — bukan bentuk, bukan hubungan.
- **25 objek tanpa DDL** di `pengetahuan/PULL-LIST.csv`. Yang disebut SQL tetap muncul di `RELASI-LAMA` dengan keadaannya tertulis; ia tidak diberi kotak karena kolomnya tidak diketahui.
- **Modul Komite Claim Non Prop** — foldernya tidak dibuka, mengikuti batas `BLUEPRINT.md`.

## 7. Batas

- **Rumus tanpa nilai tersimpan.** LibreOffice tidak terpasang di mesin ini, jadi workbook tidak dapat dihitung ulang sebelum dikirim. Sebagai gantinya `fullCalcOnLoad` dinyalakan — **Excel menghitung ulang saat berkas dibuka**. Pembaca yang hanya membaca nilai tersimpan (pandas, penampil ringan) akan melihat kedua belas sel pencacah kosong sampai berkasnya pernah dibuka Excel sekali. Angkanya sendiri ada di bagian 3 berkas ini.
- **Nol DDL dijalankan. Nol berkas sumber berubah.** Alat ini hanya membaca.
- Skema baru **usulan** — ia memerikan apa yang akan ditegakkan, bukan apa yang berdiri.
