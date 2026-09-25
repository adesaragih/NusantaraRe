# KAMUS KOLOM — modul Treaty In (skema `TREATY_MASUK`)

<!-- BERKAS INI DIBANGKITKAN. Jangan disunting dengan tangan. -->
> **Dasar bukti:** `SPEC-MODEL-DATA.md` §10 — **diurai**, bukan diketik ulang.
> Dibangkitkan `alat/buat-kamus-dan-ddl.py` + `alat/tulis-kamus-dan-ddl.py` dari definisi yang
> **sama** dengan `ddl-usulan/`. **Keduanya tidak dapat berbeda.** Bila salah satu menyimpang
> dari §10, **alatnya** yang salah.

> ### BATAS BERKAS INI, dinyatakan di muka
>
> Ia dapat menyatakan sebuah kolom **ADA di §10**. Ia **tidak** dapat menyatakan sebuah entitas
> **LENGKAP** — bila §10 menulis atribut dalam bentuk yang pengurainya tidak kenal, atribut itu
> hilang tanpa suara. Karena itu §0 **membandingkan** jumlah terurai dengan jumlah yang §10 tulis
> di judulnya sendiri, dan melaporkan setiap selisih. **"Nol selisih" bukan "lengkap"** — ia
> hanya berarti judul dan isinya sepakat.

**Cacah berkas ini:** **35** entitas, **257** kolom.

---

## 0. Pemeriksaan diri — terurai versus yang DITULIS §10

**Tidak ada selisih** antara jumlah terurai dan jumlah yang ditulis di judul §10 —
seluruh 35 entitas. Tiga lubang cacah yang pernah berdiri sudah ditutup 24 September 2026:
`F-2` (§10.2), `F-3` (§10.21), dan `F-17` (§10.19a, tabel kedua yang terhisap ke
entitas judulnya).

> **"Nol selisih" tetap BUKAN "lengkap".** Ia berarti judul dan isinya sepakat. Semesta
> yang dipakai menyusun §10 masih kurang **340 properti titik buta** (`L-8`, ditagih
> `M-4`), dan tidak ada pemeriksaan di berkas ini yang dapat melihatnya.

---

## 1. Kelompok tipe — **DIPUTUSKAN `KTV-A`, tanpa verifikasi**

§10 pendahuluan menyatakan tipe konseptualnya **mewarisi kelompok ADR-0003** dari modul Claim
Non Prop, dan bahwa **angka presisinya diserahkan ke gerbang sesi DDL**. Angka itu **sudah
diputuskan 24 September 2026**, butir **`KTV-A`** — bukan lagi warisan yang menunggu.

**Terverifikasi: tidak.** Dasarnya satu kalimat: terlalu lebar di Oracle **murah**, terlalu
sempit **memotong data** dan potongannya baru ketahuan sesudah data masuk — ongkosnya tidak
setangkup, jadi sisi murahnya yang diambil. **Syarat pembalikannya: sesi DDL boleh
mempersempit, dan HANYA sebelum data dimuat.** Alasan tiap angka ada di
`KEPUTUSAN-TANPA-VERIFIKASI.md` §7.

| Kel. | Usulan tipe Oracle | Isi |
|---|---|---|
| `U1` | `NUMBER(38,20)` | nilai uang |
| `P1` | `NUMBER(38,20)` | persentase dan porsi |
| `P2` | `NUMBER(38,20)` | tarif brokerage dan pajak |
| `K1` | `NUMBER(38,20)` | kurs |
| `L1` | `NUMBER(38,20)` | limit layer dan premi deposit |
| `C1` | `NUMBER(9)` | cacah dan nomor urut |
| `T` | `VARCHAR2(1000 CHAR)` | teks |
| `D` | `DATE` | tanggal |
| `E` | `VARCHAR2(40 CHAR)` | himpunan tertutup |
| `R` | `NUMBER(19)` | kunci asing |

**Tidak ada satu pun kolom `CLOB` maupun `BLOB`.** Kalimat yang berdiri di sini sebelumnya
— *"enam kolom memakai `CLOB`"* — **sudah tidak benar sejak keputusan itu dicabut**, dan
ia bertahan karena berkas ini dulu menyimpan salinan tabel tipenya sendiri. Salinan itu
dihapus: tipe kini dihitung **sekali** di `alat/buat-kamus-dan-ddl.py` dan dibawa di dalam
berkas definisi.

**Dua penyimpangan dari tabel di atas, keduanya `KTV-A`:**

| Yang menyimpang | Jadi | Kenapa |
|---|---|---|
| kolom bernama `ID_*` yang §10 tulis `C1` | `NUMBER(19)` | kunci asing `R` sudah `NUMBER(19)`; **dua ujung satu relasi tidak boleh berbeda lebar** |
| tujuh kolom teks bebas | `VARCHAR2(4000 CHAR)` | `KETERANGAN`, `CATATAN`, `PENGECUALIAN`, `KETENTUAN_KHUSUS`, `CATATAN_BORDEREAUX` pada `VERSI_KONTRAK`; `NILAI_SEBELUM` dan `NILAI_SESUDAH` pada `JEJAK_PERUBAHAN`. Lebarnya mengikuti lebar yang sistem lama pakai untuk teks bebasnya sendiri |

> Alasan tiap angka di tabel kelompok tipe ada di `KEPUTUSAN-TANPA-VERIFIKASI.md` §7,
> butir **`KTV-A`** — beserta **syarat pembalikannya**: sesi DDL boleh mempersempit,
> dan hanya **sebelum data dimuat**.

---

## `KONTRAK`

*§10.1 · 8 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | pengenal buatan sistem; bukan dari teks, bukan dari cap waktu |
| `NOMOR_KONTRAK_WARISAN` | `VARCHAR2(1000 CHAR)` | Y | ID | hanya terisi pada baris hasil migrasi (ADR-0042) |
| `ID_KONTRAK_DISALIN_DARI` | `NUMBER(19)` | Y | OLDID` pecahan 1 | rujukan ke `KONTRAK` lain |
| `ID_CEDANT` | `NUMBER(19)` | N | CedingID | nama cedant **tidak disalin** |
| `ID_ASAL_BISNIS` | `NUMBER(19)` | N | LeadingReinsSourceID | *source of business* |
| `SIFAT_PROPORSI` | `VARCHAR2(40 CHAR)` | N | ProportionType | `PROPORSIONAL` / `NON_PROPORSIONAL` |
| `TANGGAL_MULAI` | `DATE` | N | Commencement | batas **inklusif** (ADR-0022) |
| `TANGGAL_BERAKHIR` | `DATE` | N | Termination | batas **inklusif** |

## `VERSI_KONTRAK`

*§10.2 · 48 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `NOMOR_URUT_VERSI` | `NUMBER(9)` | N | **baru** — tidak ada di sistem lama |  |
| `KEADAAN_SIKLUS_HIDUP` | `VARCHAR2(40 CHAR)` | N | Position` + `StatusAkseptasi` + 3 bendera | ⚠ **DELAPAN nilai** — `DIBATALKAN` ditambahkan ADR-0055 perubahan 24 Sep 2026 |
| `KEADAAN_WARISAN_ASLI` | `VARCHAR2(1000 CHAR)` | Y | **baru** — tidak ada di sistem lama |  |
| `JENIS_ADDENDUM` | `VARCHAR2(40 CHAR)` | Y | EDMState |  |
| `TANGGAL_BERLAKU_ADDENDUM` | `DATE` | Y | EDMEffective |  |
| `NAMA_KONTRAK` | `VARCHAR2(1000 CHAR)` | N | TreatyContractName |  |
| `LINGKUP_WILAYAH` | `VARCHAR2(1000 CHAR)` | Y | TeritorialScope |  |
| `KELAS_BISNIS_KONTRAK` | `NUMBER(19)` | Y | ClassofBusiness |  |
| `KODE_MATA_UANG_KONTRAK` | `NUMBER(19)` | N | Currency |  |
| `ID_REASURADUR_PEMIMPIN` | `NUMBER(19)` | Y | LeadingReinsID |  |
| `ID_KETUA_TREATY` | `NUMBER(19)` | Y | TreatyLeader |  |
| `PERSEN_BAGIAN_NURE` | `NUMBER(38,20)` | N | RNMShare` + `RNMShareP |  |
| `BAGIAN_NURE_SERAGAM` | `VARCHAR2(40 CHAR)` | N | RNMShareAcrossTheBoard |  |
| `PERSEN_BAGIAN_NURE_DIPOTONG` | `NUMBER(38,20)` | Y | RnmShareDeducted |  |
| `PERSEN_BROKERAGE` | `NUMBER(38,20)` | Y | BrokeragePercent` + `…P |  |
| `PERSEN_BAGIAN_FAKULTATIF` | `NUMBER(38,20)` | Y | FacultativeShare` (`FacShare` dibuang, salinan — §14.3) |  |
| `PERSEN_BROKERAGE_FAKULTATIF` | `NUMBER(38,20)` | Y | FacultativeShareBrokerage` (`FacShareBrokerage` dibuang, salinan) |  |
| `PENGECUALIAN` | `VARCHAR2(4000 CHAR)` | Y | Exclusions` + `ExclusionsP |  |
| `KETENTUAN_KHUSUS` | `VARCHAR2(4000 CHAR)` | Y | SpecialConditions` + `SpecialConditionsP |  |
| `KETERANGAN` | `VARCHAR2(4000 CHAR)` | Y | Information |  |
| `CATATAN` | `VARCHAR2(4000 CHAR)` | Y | Comment |  |
| `CATATAN_BORDEREAUX` | `VARCHAR2(4000 CHAR)` | Y | BordereauxNote |  |
| `MEMAKAI_BORDEREAUX` | `VARCHAR2(40 CHAR)` | N | Bordeaux |  |
| `CARA_PEMBUKUAN` | `VARCHAR2(40 CHAR)` | N | AccountingMode |  |
| `CARA_PEMBUKUAN_XOL` | `VARCHAR2(40 CHAR)` | Y | AccountingModeNonProp |  |
| `PERIODE_PELAPORAN_KONTRAK` | `VARCHAR2(40 CHAR)` | Y | ReportingPeriod |  |
| `SELANG_PELAPORAN` | `NUMBER(9)` | Y | ReportingInterval |  |
| `TANGGAL_MULAI_PELAPORAN` | `DATE` | Y | ReportingStart |  |
| `TANGGAL_AKHIR_PELAPORAN` | `DATE` | Y | ReportingEnd |  |
| `HARI_BATAS_PENYERAHAN` | `NUMBER(9)` | Y | ReportingSubmission |  |
| `HARI_BATAS_KONFIRMASI` | `NUMBER(9)` | Y | ReportingConfirmation |  |
| `HARI_BATAS_PELUNASAN` | `NUMBER(9)` | Y | ReportingSettlement |  |
| `HARI_PENGINGAT` | `NUMBER(9)` | Y | ReminderDays |  |
| `PERIODE_AKUMULASI_KONTRAK` | `VARCHAR2(40 CHAR)` | Y | AccumulationPeriod |  |
| `JUMLAH_TERMIN` | `NUMBER(9)` | Y | InstallmentNo |  |
| `MEMAKAI_PRORATA` | `VARCHAR2(40 CHAR)` | N | IsProRate |  |
| `BATAS_MAKSIMUM_KELOMPOK` | `NUMBER(38,20)` | Y | MaxCoGroup |  |
| `MATA_UANG_BATAS_KELOMPOK` | `VARCHAR2(1000 CHAR)` | Y | **tidak ada di sistem lama** | **denominasi `BATAS_MAKSIMUM_KELOMPOK`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `BATAS_MAKSIMUM_NON_KELOMPOK` | `NUMBER(38,20)` | Y | MaxCoNonGroup |  |
| `MATA_UANG_BATAS_NON_KELOMPOK` | `VARCHAR2(1000 CHAR)` | Y | **tidak ada di sistem lama** | **denominasi `BATAS_MAKSIMUM_NON_KELOMPOK`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `BATAS_PILIHAN` | `NUMBER(38,20)` | Y | OptionLimit |  |
| `MATA_UANG_BATAS_PILIHAN` | `VARCHAR2(1000 CHAR)` | Y | **tidak ada di sistem lama** | **denominasi `BATAS_PILIHAN`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `RETRO_BERGANDA` | `VARCHAR2(40 CHAR)` | N | IsMultipleRetro |  |
| `ID_DAFTAR_RETRO` | `NUMBER(19)` | Y | RetroList |  |
| `SIFAT_MATERIAL_ADDENDUM` | `VARCHAR2(40 CHAR)` | Y | EDMMaterialType | dua nilai **`MATERIAL`** / **`TIDAK_MATERIAL`**. **Masukan, bukan turunan** (`GRL-20`; `GRL-12` BATAL). Beku sejak `AJUKAN`, berjejak selama `DRAFT` (`KTV-1` Adjustment) |
| `ID_DOKUMEN_ADDENDUM` | `NUMBER(19)` | Y | **baru** — tidak ada di sistem lama | satu dokumen memayungi banyak versi, **lintas kontrak** (`DB-3`, `DB-4` dibantah). **Persetujuan tetap per versi** |

## `LAYER`

*§10.3 · 24 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_LAYER` | `NUMBER(19)` | N | Limits[].ID |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `NOMOR_LAYER` | `NUMBER(9)` | N | Layer | **kunci alami** di dalam versi |
| `BAGIAN_LAYER` | `NUMBER(9)` | Y | LayerPart | bagian dari kunci alami |
| `JENIS_LAYER` | `VARCHAR2(40 CHAR)` | N | LayerType |  |
| `JENIS_BAGIAN_LAYER` | `VARCHAR2(40 CHAR)` | Y | LayerPartType |  |
| `CAKUPAN` | `VARCHAR2(1000 CHAR)` | Y | Cover |  |
| `LIMIT` | `NUMBER(38,20)` | N | Limit` + `Currency | tingkat **100% treaty** |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` / `CurrencyID` pada baris yang sama | **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `DEDUCTIBLE` | `NUMBER(38,20)` | N | Deductible | tingkat **100% treaty** |
| `MATA_UANG_DEDUCTIBLE` | `VARCHAR2(1000 CHAR)` | Y | **tidak ada di sistem lama** | **denominasi `DEDUCTIBLE`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `PERSEN_PENYESUAIAN` | `NUMBER(38,20)` | Y | AdjRate |  |
| `PERSEN_MINIMUM_DEPOSIT` | `NUMBER(38,20)` | Y | MDPPct |  |
| `PORSI_PEMULIHAN_LIMIT` | `NUMBER(38,20)` | Y | ReinstatementPct | porsi limit yang dipulihkan |
| `TARIF_PREMI_PEMULIHAN` | `NUMBER(38,20)` | Y | ReinstatementValue | tarif premi untuk pemulihan itu — **pembacaannya punya saingan**, §10.3b |
| `LIMIT_AGREGAT` | `NUMBER(38,20)` | Y | AgregateLimit` *(salah eja ada di sumbernya) | **BARU** — batas total sepanjang periode, terpisah dari limit per kejadian. Tingkat **BELUM DITENTUKAN** |
| `MATA_UANG_LIMIT_AGREGAT` | `VARCHAR2(1000 CHAR)` | Y | **tidak ada di sistem lama** | **denominasi `LIMIT_AGREGAT`** — `KTV-C`. Sistem lama tidak merekam mata uang untuk besaran ini sama sekali; kolomnya disediakan sekarang karena menambahkannya sesudah data masuk mahal. **Dicabut sebelum data dimuat bila `T-6` menyatakan ia selalu mata uang kontrak** |
| `MDP` | `NUMBER(38,20)` | Y | MDPList[] | deposit premium; tingkat **BELUM DITENTUKAN** |
| `MDP_MINIMUM` | `NUMBER(38,20)` | Y | MDPMinList[] | **BARU** — tanpa ia, *"deposit premium sama dengan minimum premium"* tidak dapat dinyatakan |
| `PERSEN_MDP_MINIMUM` | `NUMBER(38,20)` | Y | MDPMinPct | **BARU** |
| `MDP_DIGABUNG` | `VARCHAR2(40 CHAR)` | N | IsCombineMDP | **BARU** — MDP dihitung per layer atau gabungan |
| `TANPA_HITUNG_PREMI_PEMULIHAN` | `VARCHAR2(40 CHAR)` | N | NoRIPCalculation | **BARU** — mematikan perhitungan premi pemulihan pada layer ini |
| `DEDUCTIBLE_KEDUA` | `NUMBER(38,20)` | Y | Limits[].Deductible2 | **PULIH dari `L-8`** — `TDA-17`. Tersimpan sebagai **teks** di sistem lama (`@toDecimal` di kedua sisi rumus selisih); migrasi wajib mengubah tipe, dan baris yang gagal konversi **tidak disamarkan jadi nol** (ADR-0035) |
| `PERSEN_ROL` | `NUMBER(38,20)` | Y | Limits[].ROLPct | **PULIH dari `L-8`** — `TDA-17`. Rate on line. ⚠ **DIBANTAH `F-15`** — pertanyaan *"selalu dihitung atau pernah disepakati"* **terjawab dari sumber**: `DetailCalculationROL` menghitungnya `premi ÷ limit × 100`. Kolomnya **tidak dicabut diam-diam**; keputusannya pemilik proses, `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` |

## `DETAIL_PROPORSIONAL`

*§10.4 · 12 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_DETAIL_PROPORSIONAL` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_LAYER` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_KELOMPOK_TREATY` | `NUMBER(19)` | N | TreatyGroupID | **kunci alami** di dalam **`LAYER`** — INV-06 |
| `JENIS_TREATY` | `VARCHAR2(40 CHAR)` | N | TreatyType | `QUOTA_SHARE` / `SURPLUS` |
| `PERSEN_QUOTA_SHARE` | `NUMBER(38,20)` | Y | QSPct | terisi hanya bila `JENIS_TREATY = QUOTA_SHARE` |
| `JUMLAH_LINES_SURPLUS` | `NUMBER(9)` | Y | Surplus | terisi hanya bila `JENIS_TREATY = SURPLUS` |
| `PERSEN_KOMISI_KOTOR` | `NUMBER(38,20)` | Y | RIOGR | kepanjangan **belum diketahui** |
| `PERSEN_KOMISI_BERSIH` | `NUMBER(38,20)` | Y | RIONR | kepanjangan **belum diketahui** |
| `PERSEN_CADANGAN_PREMI` | `NUMBER(38,20)` | Y | PremiumReservePct | dari inventaris kelas Pega |
| `CADANGAN_PREMI` | `NUMBER(38,20)` | Y | ReserveList[] | **BARU** — pasangan nilai dari persentase di atas; persentase tanpa nilainya adalah setengah fakta. Tingkat **BELUM DITENTUKAN** |
| `PERSEN_KAPASITAS_SURPLUS` | `NUMBER(38,20)` | Y | IOOPct | **BARU** |
| `ID_SUSUNAN_RETRO` | `NUMBER(19)` | Y | SpreadingTypeID | **BARU — penunjuk susunan baku yang menyemai penyebaran.** Syarat berdirinya `RINCIAN_PENYEBARAN` sebagai fakta terbukukan (INV-58); lihat §10.4a |

## `BAGIAN`

*§10.5 · 8 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_BAGIAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_LAYER` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | **relasi**, bukan salinan — §12.2 |
| `CAKUPAN` | `VARCHAR2(1000 CHAR)` | Y | Cover | tetap atribut bagian (§12.2) |
| `PERSEN_BAGIAN_NURE` | `NUMBER(38,20)` | Y | Share[].RNMShare | terisi bila `BAGIAN_NURE_SERAGAM` mati; berbeda dari atribut senama pada versi (§10.2) yang berlaku seragam |
| `PREMI_BRUTO` | `NUMBER(38,20)` | N | GrossPremiumList[] | tingkat **BELUM DITENTUKAN** — §10.23 |
| `PREMI_BRUTO_MINIMUM` | `NUMBER(38,20)` | Y | GrossPremiumMinList[] | **BARU** — tidak pernah terlihat pohon |
| `ID_SUSUNAN_RETRO` | `NUMBER(19)` | Y | SpreadingTypeIDXOL | penunjuk susunan baku yang menyemai penyebaran — **syarat INV-58**, §10.4a |
| `PERSEN_BAGIAN_DIPAKAI` | `NUMBER(38,20)` | Y | **baru** — tidak ada di sistem lama | wajib bila paket uangnya bertingkat `BAGIAN_NURE` (INV-40) |

## `PENYEBARAN`

*§10.6 · 6 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_PENYEBARAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_BAGIAN` | `NUMBER(19)` | Y | **baru** — tidak ada di sistem lama | `KTV-B` — salah satu dari **dua** pelekatan. **Tepat satu** dari kedua kolom terisi, dijaga `CHECK`. Menggantikan `ID_INDUK_PENYEBARAN`, yang tidak dapat punya kunci asing karena sasarannya bergantung nilai kolom lain (INV-17) |
| `ID_DETAIL_PROPORSIONAL` | `NUMBER(19)` | Y | **baru** — tidak ada di sistem lama | `KTV-B` — salah satu dari **dua** pelekatan. **Tepat satu** dari kedua kolom terisi, dijaga `CHECK`. Menggantikan `ID_INDUK_PENYEBARAN`, yang tidak dapat punya kunci asing karena sasarannya bergantung nilai kolom lain (INV-17) |
| `ID_JENIS_REASURANSI` | `NUMBER(19)` | N | ReinsTypeID | **kunci alami** |
| `ID_JENIS_REASURANSI_INDUK` | `NUMBER(19)` | Y | ParentReinsTypeID | jenis reasuransi **bersusun**; rujukan ke baris lain di tabel acuan yang sama |
| `PERSEN_PENYEBARAN` | `NUMBER(38,20)` | N | Pct | INV-50 menjumlahkannya |

## `RINCIAN_PENYEBARAN`

*§10.7 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_RINCIAN_PENYEBARAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_PENYEBARAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_JENIS_REASURANSI` | `NUMBER(19)` | N | ReinsID` (prop) / `ReinsTypeID` (XOL) | **kunci alami**. Ruas proporsional **bernama pihak tetapi berisi jenis** — §4.4a |
| `PERSEN_RINCIAN` | `NUMBER(38,20)` | N | SharePct | INV-50 menjumlahkannya per induk |

## `NILAI_PENYEBARAN`

*§10.8 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_NILAI_PENYEBARAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_RINCIAN_PENYEBARAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `NILAI` | `NUMBER(38,20)` | N | Amount` + `Currency |  |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` / `CurrencyID` pada baris yang sama | **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |

## `POTONGAN`

*§10.9 · 6 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_POTONGAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | §14.2 |
| `ID_BAGIAN` | `NUMBER(19)` | Y | **baru** — tidak ada di sistem lama | `KTV-B` — salah satu dari **dua** pelekatan. **Tepat satu** dari kedua kolom terisi, dijaga `CHECK`. Menggantikan `ID_INDUK_POTONGAN`, yang tidak dapat punya kunci asing karena sasarannya bergantung nilai kolom lain (INV-17) |
| `ID_DETAIL_PROPORSIONAL` | `NUMBER(19)` | Y | **baru** — tidak ada di sistem lama | `KTV-B` — salah satu dari **dua** pelekatan. **Tepat satu** dari kedua kolom terisi, dijaga `CHECK`. Menggantikan `ID_INDUK_POTONGAN`, yang tidak dapat punya kunci asing karena sasarannya bergantung nilai kolom lain (INV-17) |
| `ID_JENIS_POTONGAN` | `NUMBER(19)` | N | Comment | §14.2 |
| `DASAR_PERHITUNGAN` | `VARCHAR2(40 CHAR)` | N | **baru** — tidak ada di sistem lama | §14.2 |
| `PERSEN_POTONGAN` | `NUMBER(38,20)` | N | DeductionPct | §14.2 |

## `MATA_UANG_KONTRAK`

*§10.10 · 6 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_MATA_UANG_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `KODE_MATA_UANG` | `NUMBER(19)` | N | CurrencyID | **kunci alami**; `Currency` (nama) tidak disimpan — INV-59 |
| `KURS` | `NUMBER(38,20)` | N | Conversion | INV-38: selalu lebih besar dari nol |
| `TANGGAL_MULAI_BERLAKU` | `DATE` | N | PeriodStart | batas inklusif, ADR-0022 |
| `TANGGAL_AKHIR_BERLAKU` | `DATE` | N | PeriodEnd | batas inklusif |

## `RETENSI_CEDANT`

*§10.11 · 7 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_RETENSI_CEDANT` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_KELOMPOK_TREATY` | `NUMBER(19)` | N | TreatyGroupID | bagian **kunci alami**; `TreatyGroup` (nama) tidak disimpan |
| `NILAI_RETENSI` | `NUMBER(38,20)` | N | Amount` + `Currency`/`CurrencyID | bagian **kunci alami** lewat kode mata uangnya. Tingkat **BELUM DITENTUKAN** — §10.23 |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` / `CurrencyID` pada baris yang sama | **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `CATATAN` | `VARCHAR2(1000 CHAR)` | Y | Note | **BARU** — ada di kelas, tidak ada di §3.5 |
| `PERSEN_BAGIAN_DIPAKAI` | `NUMBER(38,20)` | Y | **baru** — tidak ada di sistem lama | bila paket uangnya bertingkat `BAGIAN_NURE` (INV-40) |

## `EGNPI`

*§10.12 · 9 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_EGNPI` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_KELOMPOK_TREATY` | `NUMBER(19)` | N | TreatyGroupID | bagian **kunci alami** |
| `ID_KELAS_BISNIS` | `NUMBER(19)` | Y | ClassOfBusiness | **BARU** — tidak ada di §3.5 |
| `NILAI_EGNPI` | `NUMBER(38,20)` | N | Amount` + `Currency`/`CurrencyID | tingkat **BELUM DITENTUKAN**, dan ini **satu dari lima §7 yang naik ke eskalasi** |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` / `CurrencyID` pada baris yang sama | **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `TANGGAL_BERLAKU` | `DATE` | Y | AsDate |  |
| `PROPORSI` | `NUMBER(38,20)` | Y | Proportion | **BARU** — belum terpakai di rumus mana pun; dibawa karena membuangnya menghilangkan fakta yang tidak dapat dipulihkan |
| `CATATAN` | `VARCHAR2(1000 CHAR)` | Y | Note | **BARU** |

## `PORTOFOLIO`

*§10.13 · 5 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_PORTOFOLIO` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ARAH_PORTOFOLIO` | `VARCHAR2(40 CHAR)` | N | Type | masuk atau keluar; bagian **kunci alami** |
| `JENIS_PORTOFOLIO` | `VARCHAR2(40 CHAR)` | N | TypePortfolio | bagian **kunci alami** |
| `KETERANGAN` | `VARCHAR2(1000 CHAR)` | Y | Description |  |

## `PERIODE_PELAPORAN`

*§10.14 · 7 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_PERIODE_PELAPORAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `PERIODE` | `VARCHAR2(1000 CHAR)` | N | Period | **kunci alami** |
| `TANGGAL_AWAL` | `DATE` | N | InitialDate |  |
| `BATAS_PENYERAHAN` | `DATE` | N | SubmissionDue |  |
| `BATAS_KONFIRMASI` | `DATE` | N | ConfirmationDue |  |
| `BATAS_PELUNASAN` | `DATE` | N | SettlementDue |  |

## `PERIODE_AKUMULASI`

*§10.15 · 6 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_PERIODE_AKUMULASI` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `PERIODE` | `VARCHAR2(1000 CHAR)` | N | Period | **kunci alami** |
| `TANGGAL_LAPOR` | `DATE` | N | ReportDate |  |
| `HARI_BATAS_PENYERAHAN` | `NUMBER(9)` | Y | SubDays | dari inventaris kelas, §12.4 |
| `BATAS_PENYERAHAN` | `DATE` | Y | SubDueDate | dari inventaris kelas, §12.4 |

## `TERMIN`

*§10.16 · 9 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_TERMIN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `NOMOR_TERMIN` | `NUMBER(9)` | N | InstallmentList[].Installment | bagian **kunci alami**, bersama kode mata uang — §10.16a |
| `PERSEN_TERMIN` | `NUMBER(38,20)` | N | InstallmentPct |  |
| `NILAI_TERMIN` | `NUMBER(38,20)` | Y | Amount` + `Currency`/`CurrencyID | tingkat **BELUM DITENTUKAN** — §10.23 |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` / `CurrencyID` pada baris yang sama | **denominasi paket uang di sebelahnya** — bagian kunci alami (INV-08 / INV-09 / INV-12). Ditambahkan 24 Sep 2026, P-8 golongan B |
| `TANGGAL_JATUH_TEMPO` | `DATE` | N | DueDate |  |
| `TANGGAL_BAYAR` | `DATE` | Y | PaymentDate | kosong selama belum dibayar |
| `WPC` | `VARCHAR2(1000 CHAR)` | Y | WPC | **singkatan yang kepanjangannya tidak ada di korpus** — lihat di bawah |

## `SKALA_KOASURANSI`

*§10.17 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_SKALA_KOASURANSI` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `PERSEN_LIMIT` | `NUMBER(38,20)` | N | PctLimit |  |
| `PERSEN_BAGIAN` | `NUMBER(38,20)` | N | CoInShare |  |

## `BATAS_PER_BAHAYA`

*§10.18 · 6 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_BATAS_PER_BAHAYA` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_BAHAYA` | `NUMBER(19)` | N | nama kolom lama | **kunci alami**; tabel acuan, ADR-0038 |
| `NILAI_BATAS` | `NUMBER(38,20)` | N | Earthquake` / `FloodJab` / `FloodNation` / `RSMDLimit` + mata uangnya | tingkat **BELUM DITENTUKAN** — §10.23 |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | Y | `Currency` / `CurrencyID` pada baris yang sama | **denominasi paket uang di sebelahnya**. Golongan **B**, dipindahkan dari golongan C oleh `KTV-C` — §10.18 menulis asalnya *"+ mata uangnya"*. **Boleh kosong**: ia bukan bagian kunci alami (INV-14 memakai `ID_BAHAYA`) |
| `PERSEN_BAGIAN_DIPAKAI` | `NUMBER(38,20)` | Y | **baru** — tidak ada di sistem lama | bila bertingkat `BAGIAN_NURE` (INV-40) |

## `DOKUMEN_KONTRAK`

*§10.19 · 5 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_DOKUMEN_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_DOKUMEN` | `NUMBER(19)` | N | lampiran | **kunci alami**; rujukan ke luar skema |
| `JENIS_DOKUMEN` | `VARCHAR2(40 CHAR)` | Y | lampiran |  |
| `TANGGAL_LAMPIR` | `DATE` | N | lampiran |  |

## `DOKUMEN_ADDENDUM`

*§10.19a · 3 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_DOKUMEN_ADDENDUM` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `NOMOR_DOKUMEN` | `VARCHAR2(1000 CHAR)` | N | tidak ada di sistem lama | **kunci alami**, unik **global** — `INV-71`. Nomor yang beredar di luar sistem: ditulis di kertas dan disebut orang |
| `TANGGAL_BERLAKU` | `DATE` | Y | **baru** — tidak ada di sistem lama | `KTV-2` — bila kosong, berlaku mengikuti versi yang dipayunginya. **Dicabut bila `DB-16b` dibantah** |

## `CATATAN_PERSETUJUAN`

*§10.20 · 6 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_CATATAN_PERSETUJUAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | menggantung pada **versi**, bukan kontrak |
| `WAKTU_KEPUTUSAN` | `DATE` | N | Date |  |
| `NAMA_PEMUTUS` | `VARCHAR2(1000 CHAR)` | N | OperatorName | **dibekukan sebagai teks dengan sengaja** — ia fakta historis siapa memutuskan apa dan kapan (§12.5, ADR-0045) |
| `DISETUJUI` | `VARCHAR2(40 CHAR)` | N | IsApproved |  |
| `ALASAN` | `VARCHAR2(1000 CHAR)` | Y | Suggest |  |

## `JEJAK_PERUBAHAN`

*§10.21 · 8 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_JEJAK_PERUBAHAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `WAKTU_PERUBAHAN` | `DATE` | N | **baru** — tidak ada di sistem lama | **`DATE` beresolusi detik.** Dua perubahan dalam detik yang sama tidak dapat dipisahkan olehnya — dan itu tidak merusak apa pun di sini, sebab entitas ini **sengaja tanpa kunci alami** |
| `PELAKU` | `VARCHAR2(1000 CHAR)` | N | **baru** — tidak ada di sistem lama | dibekukan, sekeluarga dengan `NAMA_PEMUTUS` |
| `PERAN_PELAKU` | `VARCHAR2(1000 CHAR)` | N | **baru** — tidak ada di sistem lama | **peran yang berlaku SAAT ITU** — ADR-0045 isi minimal butir kelima. **POTRET, bukan rujukan** (`KTV-D`): teks, final, **tidak dinormalisasi ulang** ketika entitas peran kelak lahir. Dari mana nilainya diambil masih `F-16` |
| `RUAS_YANG_BERUBAH` | `VARCHAR2(1000 CHAR)` | N | **baru** — tidak ada di sistem lama |  |
| `NILAI_SEBELUM` | `VARCHAR2(4000 CHAR)` | Y | **baru** — tidak ada di sistem lama | teks, karena ruasnya beragam tipe |
| `NILAI_SESUDAH` | `VARCHAR2(4000 CHAR)` | Y | **baru** — tidak ada di sistem lama | idem |

## `PERISTIWA_KONTRAK`

*§10.21a · 5 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_PERISTIWA_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_VERSI_KONTRAK` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | menggantung pada **versi** |
| `WAKTU` | `DATE` | N | CommentList[].Date |  |
| `NAMA_PELAKU` | `VARCHAR2(1000 CHAR)` | N | CommentList[].OperatorName | dibekukan sebagai teks, sekeluarga `NAMA_PEMUTUS` (§12.5) |
| `JENIS_PERISTIWA` | `VARCHAR2(40 CHAR)` | N | CommentList[].Suggest` **diurai | `SUNTINGAN_INTERNAL` · `ADDENDUM_EKSTERNAL` · `ADDENDUM_PREMI` · `DISALIN` · `REVISI_DIBUAT` |

## `PEMULIHAN_LIMIT`

*§10.3b · 6 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_PEMULIHAN_LIMIT` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_LAYER` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `NOMOR_URUT_PEMULIHAN` | `NUMBER(9)` | N | Reinstatement_List[].ReinstatementValue |  |
| `PERSEN_PEMULIHAN` | `NUMBER(38,20)` | N | Reinstatement_List[].ReinstatementPct |  |
| `PERSEN_TAMBAHAN` | `NUMBER(38,20)` | Y | Reinstatement_List[].AdditionalPct |  |
| `CATATAN` | `VARCHAR2(1000 CHAR)` | Y | Reinstatement_List[].ReinstatementNote |  |

## `MATA_UANG`

*§10.22 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_MATA_UANG` | `NUMBER(19)` | N | tabel acuan (ADR-0038) |  |
| `KODE` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) | **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) |  |
| `AKTIF` | `VARCHAR2(40 CHAR)` | N | tabel acuan (ADR-0038) | nilai lama **tidak pernah dihapus**; ia dimatikan |

## `JENIS_POTONGAN`

*§10.22 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_JENIS_POTONGAN` | `NUMBER(19)` | N | tabel acuan (ADR-0038) |  |
| `KODE` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) | **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) |  |
| `AKTIF` | `VARCHAR2(40 CHAR)` | N | tabel acuan (ADR-0038) | nilai lama **tidak pernah dihapus**; ia dimatikan |

## `JENIS_REASURANSI`

*§10.22 · 5 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_JENIS_REASURANSI` | `NUMBER(19)` | N | tabel acuan (ADR-0038) |  |
| `KODE` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) | **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) |  |
| `AKTIF` | `VARCHAR2(40 CHAR)` | N | tabel acuan (ADR-0038) | nilai lama **tidak pernah dihapus**; ia dimatikan |
| `ID_INDUK` | `NUMBER(19)` | Y | tabel acuan (ADR-0038) | **hanya `JENIS_REASURANSI`** — ia bersusun, §10.6 |

## `BAHAYA`

*§10.22 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_BAHAYA` | `NUMBER(19)` | N | tabel acuan (ADR-0038) |  |
| `KODE` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) | **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) |  |
| `AKTIF` | `VARCHAR2(40 CHAR)` | N | tabel acuan (ADR-0038) | nilai lama **tidak pernah dihapus**; ia dimatikan |

## `KELOMPOK_TREATY`

*§10.22 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_KELOMPOK_TREATY` | `NUMBER(19)` | N | tabel acuan (ADR-0038) |  |
| `KODE` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) | **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) |  |
| `AKTIF` | `VARCHAR2(40 CHAR)` | N | tabel acuan (ADR-0038) | nilai lama **tidak pernah dihapus**; ia dimatikan |

## `KELAS_BISNIS`

*§10.22 · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_KELAS_BISNIS` | `NUMBER(19)` | N | tabel acuan (ADR-0038) |  |
| `KODE` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) | **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | `VARCHAR2(1000 CHAR)` | N | tabel acuan (ADR-0038) |  |
| `AKTIF` | `VARCHAR2(40 CHAR)` | N | tabel acuan (ADR-0038) | nilai lama **tidak pernah dihapus**; ia dimatikan |

## `NILAI_PREMI_BRUTO`

*§10.x · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_NILAI_PREMI_BRUTO` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_BAGIAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` pada baris daftarnya | **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | `NUMBER(38,20)` | N | `Value` pada baris daftarnya |  |

## `NILAI_PREMI_BRUTO_MINIMUM`

*§10.x · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_NILAI_PREMI_BRUTO_MINIMUM` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_BAGIAN` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` pada baris daftarnya | **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | `NUMBER(38,20)` | N | `Value` pada baris daftarnya |  |

## `NILAI_CADANGAN_PREMI`

*§10.x · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_NILAI_CADANGAN_PREMI` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_DETAIL_PROPORSIONAL` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` pada baris daftarnya | **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | `NUMBER(38,20)` | N | `Value` pada baris daftarnya |  |

## `NILAI_MDP`

*§10.x · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_NILAI_MDP` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_LAYER` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` pada baris daftarnya | **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | `NUMBER(38,20)` | N | `Value` pada baris daftarnya |  |

## `NILAI_MDP_MINIMUM`

*§10.x · 4 kolom*

| Kolom | Tipe | Null | Asal di sistem lama | Catatan |
|---|---|---|---|---|
| `ID_NILAI_MDP_MINIMUM` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama |  |
| `ID_LAYER` | `NUMBER(19)` | N | **baru** — tidak ada di sistem lama | induknya — paket uang ini asalnya **daftar per mata uang** di sistem lama |
| `KODE_MATA_UANG` | `VARCHAR2(1000 CHAR)` | N | `Currency` pada baris daftarnya | **kunci alami** di dalam induknya — belum bernomor |
| `NILAI` | `NUMBER(38,20)` | N | `Value` pada baris daftarnya |  |
