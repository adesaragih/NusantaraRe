> Modul  : Treaty In Adjustment · Ronde B · 2026-09-26
> Peran  : interogator
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 01 · TEMUAN RONDE B

## 1. Tabel temuan

| # | Temuan | Sisi | Mengubah |
|---|---|---|---|
| **NB-01** | **Pengenal bersusun tidak pernah terbentuk** — bukti B2 saya keliru; addendum yang dipilih dari picker menghasilkan nomor berikutnya yang biasa | 56-KHAS | bukti B2, `UA-1`, `MA-09` |
| **NB-02** | Percabangan rantai **tidak terlihat dari pengenal**; satu-satunya saksinya `OLDID` | 56-KHAS | `UA-1`, `DB-2` |
| **NB-04** | **Audit klaim hidup/mati: ketujuhnya bertahan.** Tidak ada penanda mati ketiga di dalam Activity, dan `pyDisabled` **tidak pernah** muncul di Activity | — | tidak ada keputusan yang berubah |
| **NB-05** | Penanda nonaktif terbaca lengkap: `pyDisabledNew = always` (146) **nonaktif selalu**, `= true` + `pyDisabledWhen` (249) **nonaktif bersyarat**, `= false` (1.302) aktif | 56-KHAS | `IND-5`, NA-21, NA-22 |
| **NB-08** | **Selisih EGNPI pada addendum premi selalu nol** — mesin premi membaca `TreatyIn.EGNPI`, pengguna menyunting `ActualValue.EGNPI` | 56-KHAS | **TDA-16**, syarat C3, GRL-12 |
| **NB-09** | Tiga belas sel berpindah golongan antara dua sensus; sebabnya **syarat tersimpan pada mode yang tidak diaktifkan** | — | §7.7, hitungan 96 -> 83 |
| **NB-06** | **Daftar bocor `DB-10` menyusut**: total angsuran, `LayersEDM`, dan `DetailEGNPI` ternyata **nonaktif permanen**; kalimat "`Layers` menjaga, `LayersEDM` tidak" **terbalik** | 56-KHAS | `DB-10`, GRL-08 butir iii |
| **NB-03** | `ID_VERSI_KONTRAK_DASAR` **sudah diputuskan** modul induk — disimpan eksplisit, tepat satu penunjuk, dengan alasan yang persis menjawab B2 | — | GRL-10, `MA-09` |

## 2. Rincian

### NB-01 · Pengenal bersusun tidak pernah terbentuk

Bukti yang saya ajukan untuk B2 salah pada dua tempat, dan keduanya bertentangan dengan temuan saya
sendiri (`MA-09`).

**Kekeliruan pertama — "nomor tertinggi + 1".** `GetTreatyRevisionID` memang menghitung
`TO_NUMBER(SUBSTR(ID,10,2))+1 as HASIL1`, tetapi **nilainya dibuang**. §4.4 butir 1 dan `MA-03`
sudah menetapkannya: `HASIL1` hanya diuji **kosong atau tidak**, dan nomornya dihitung ulang dari
`TreatyIn.ID` pada langkah 5 `TreatyInRevisi_post`.

**Kekeliruan kedua — "addendum yang dipilih menghasilkan `HASIL1` kosong".** Baris addendum itu
**ada di `m_treaty_in_edm`** dan memenuhi `LIKE`-nya sendiri:

| Yang dipilih | `InputData.CARI1` | Apakah cocok dengan dirinya sendiri? |
|---|---|---|
| `1234567/R03` | `1234567/R03%` | **ya** — barisnya sendiri |

Yang memastikannya ada di sana: `PEGA_M_TREATY_IN_EDM` menulis ke **`M_TREATY_IN_EDM` dan
`TREATY_IN_EDM` sekaligus**, pada cabang sisip maupun cabang perbarui. Maka `HASIL1` **tidak
kosong**, langkah 5 berjalan, dan hasilnya:

```
@substring("1234567/R03", 10, 12)  ->  "3"        (hanya digit satuan - offset yang meleset satu)
@toInt("3") + 1                    ->  4
@substring("1234567/R03", 0, 7)    ->  "1234567"
hasil                              ->  "1234567/R04"
```

**Pengenal bersusun `‹kontrak›/R03/R01` tidak pernah terbentuk.** Ia hanya mungkin bila baris
addendum itu **tidak ada** di `m_treaty_in_edm` — dan prosedur simpan memastikan ia ada.

Ini sekaligus menegaskan deret `R09 -> R010` (§4.4, NA-02): langkah 5 **memang** berjalan atas
revisi yang dipilih, dan justru itulah sebabnya offsetnya berakibat.

### NB-02 · Percabangan hanya terlihat dari `OLDID`

Yang tetap berdiri dari bukti B2: **`OLDID` menunjuk baris yang dipilih**
(`TreatyInRevisi_post` baris 794-795, `TreatyIn.OLDID = TreatyIn.ID`), dan picker meng-`UNION`
kontrak dengan addendum tanpa saringan (TDA-11). Jadi rantai lama **dapat bercabang** — sebuah
addendum dapat dibuat di atas addendum yang bukan versi terakhir.

Tetapi percabangan itu **tidak meninggalkan jejak pada bentuk pengenal**: baik yang dibuat dari
versi terakhir maupun dari versi lama menghasilkan `‹kontrak›/Rnn` yang serupa. Satu-satunya
saksinya `OLDID`.

`UA-1` butir (c) karena itu diganti: bukan "pengenal yang memuat lebih dari satu `/R`", melainkan
**addendum yang `OLDID`-nya menunjuk baris addendum**, dan — yang sebenarnya menandai percabangan —
**addendum yang `OLDID`-nya bukan versi terakhir yang ada pada saat ia dibuat**.

### NB-03 · Rujukan versi dasar sudah diputuskan induk

`4-erd-dan-tabel-datar/STRUKTUR-DATA.md` §1.1:

> **`ID_VERSI_KONTRAK_DASAR` adalah satu-satunya penunjuk yang disimpan, dan ia tidak punya
> pasangan.** Sisi "baru" dari sebuah selisih **adalah versi induk barisnya sendiri** — tidak perlu
> kolom untuk itu.
>
> Ia disimpan **eksplisit**, bukan disimpulkan dari `NOMOR_URUT_VERSI - 1`, karena *picker* sistem
> lama meng-UNION kontrak dengan addendum sehingga dasar sebuah penyesuaian **belum tentu** versi
> tepat sebelumnya.

Alasannya **persis** alasan yang saya bangun ulang dari nol di B2 — dan kesimpulannya **berlawanan**
dengan rekomendasi saya, yang mengusulkan menurunkan alih-alih menyimpan. `ERD.md` §2.6 menegaskan
sisi lainnya: tidak ada relasi dari `NILAI_SELISIH` ke "versi lama"; sisi lama dibaca lewat
`VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR` pada induknya sendiri.

Jadi "bentuk simpan rantai" **bukan pertanyaan terbuka** — ia **KONFIRMASI**.

### NB-04 · Audit klaim hidup/mati — ketujuhnya bertahan

`SAPUAN-DAN-NAMA-TAGNYA.md` §4 menyebut cacat `langkah-hidup.py` induk dalam satu kalimat: **ia
hanya mengenal `pyStepsBlockName == "//"`, dan buta terhadap `pyDisabled == "true"`** — 1.306
kemunculan, seluruhnya di `Rule-HTML-Section` (1.082), `Rule-HTML-Harness` (198), dan
`Rule-Obj-Model` (26). **Irisan kedua penanda nol**: yang dikenali perkakas tidak pernah muncul di
luar aktivitas, yang tidak dikenalinya tidak pernah muncul di dalamnya. Akibatnya untuk setiap jenis
aturan **selain aktivitas**, perkakas itu tidak pernah memeriksa apa pun, dan diamnya terbaca sebagai
"semuanya hidup".

**Dua pemeriksaan sebelum klaim mana pun dipercaya.**

1. **Adakah penanda mati ketiga di dalam Activity?** Seluruh tag ber-"disab", "enabl", "active",
   "skip", "bypass" di folder `Activity/` disapu. Yang muncul: `pyIsStaleRuleCheckEnabled` (155,
   setelan formulir aturan), `pyEnableAddNewField` (8), `SkipVisibility` dan `bypasscondition`
   (parameter aktivitas, bukan penonaktif). **`pyDisabled` nol di seluruh folder `Activity/`** —
   sejalan dengan pernyataan induk bahwa irisannya nol. Jadi untuk aktivitas, `pyStepsBlockName`
   memang penanda yang benar.
2. **Kalibrasi** (`TA-04`): `SaveTreatyIn_EDM_Act` langkah 12 — diketahui mati — terbaca
   `blockName = "//"`; langkah 14 — diketahui hidup — terbaca kosong. Penyapunya membedakan
   keduanya.

**Hasil sapuan ulang, dengan pewarisan dari langkah induk ikut dihitung:**

| Klaim | Cara pemeriksaan ulang | Hasil | Keputusan yang bersandar |
|---|---|---|---|
| `SaveTreatyIn_EDM_Act` 12-13 mati, 14 hidup | blockName + preCondition + pyDisabled, berjenjang | **bertahan** — 12 dan 13 `//`, 14 kosong | GRL-01, GRL-11 label PERUBAHAN |
| `TreatyInEdmCheckDuplicate` mati | idem, `TreatyInEDMSetValue` langkah 7 | **bertahan** — `//` | TDA-01 |
| lima langkah pro rata mati | seluruh `TreatyEDMProRateCalculation` disapu | **bertahan, dan lebih tepat: 10 dari 53 langkah mati** — lima blok `//` (2.1, 2.8, 2.9, 2.10, 8) ditambah lima langkah bersarang yang ikut mati | E3a |
| `TreatyInSetToDirector` mati seluruhnya | idem, termasuk pewarisan | **bertahan** — langkah 1-4 `//`, dan 4.1-4.7 mati **karena induknya** | REV-3 butir i |
| `AddCommentList_Act` langkah 1 mati | idem | **bertahan** — `//`; langkah 2 dan 3 hidup | UA-16 |
| `Akseptasi_DT` cabang 1.4 mati | `pyDisabled = true`, 13 baris | **bertahan** | GRL-07, REV-2 |
| `TreatyEDMDifferenceDeduction` 4 dan 5 mati | blockName | **bertahan** | E3b |

**Tidak satu pun berbalik. Tidak ada label atau keputusan yang berubah.**

Satu kekeliruan penamaan saya sendiri yang muncul saat audit dan tidak menyentuh dokumen mana pun:
`SetTreatyIn_Act` (yang memuat langkah 6, 7, dan 9 pada §4.5 dan NA-21) dan `SaveTreatyIn_Act` (yang
langkah 9-nya sebuah panggilan mati) adalah **dua aktivitas berbeda dengan nama yang hampir sama**.
Keduanya diperiksa terpisah: `SetTreatyIn_Act` 6 dan 7 **hidup**, sesuai NA-21.

### NB-05 · `pyDisabled` di Section tidak berarti "kendali mati"

Ekspor Adjustment memuat **395** `pyDisabled = true` di badan `Section/` — terpusat di
`DetailLimits` (68), `Layers` (41), `TreatyInTabsNonProportional` (41), `TreatyInTabsProportional`
(37). Karena angka itu menyentuh justru seksi yang saya pakai untuk NA-21 dan NA-22, sensus sel
disunting diulang dengan penanda itu ikut dibaca:

| | Jumlah |
|---|---|
| dijaga `ViewState` **dan** punya mode ber-`pyDisabled` | **96** |
| tanpa penjaga, punya mode ber-`pyDisabled` | 120 |
| tanpa penjaga, tanpa mode mati | 762 |

**Seluruh 96 sel terjaga memilikinya — seratus persen.** Sebuah penanda yang benar berarti "kendali
mati" tidak mungkin melekat pada **setiap** kendali yang dijaga kondisi. Sebabnya terbaca dari
letaknya: `pyDisabled` duduk di dalam `pyModes/rowdata` ber-`pxObjClass = Embed-Control-Mode`,
yaitu **per-mode**, bukan per-sel. Satu sel punya beberapa mode, dan mode yang tidak dipakai
ditandai mati.

**Maka NA-01, NA-17b, NA-21, dan NA-22 bertahan** — dan diperiksa satu per satu: kedua sel radio di
`PickerTreatyInMasterRevisi` dan keempat sel lapisan beku di kedua layar **tidak memuat satu pun**
mode ber-`pyDisabled`.

Ini **mempertajam** peringatan induk sendiri, yang menyebut `pyDisabled` di Section *"hampir pasti
menandai kendali layar yang dimatikan"*. Di ekspor ini ia menandai **mode**, bukan kendali — dan
perbedaan itu menentukan apakah 395 kemunculan dibaca sebagai 395 kendali mati atau sebagai
struktur biasa. Diusulkan balik sebagai **`IND-5`**.

### NB-06 · Penanda nonaktif ditafsirkan, dan daftar bocor menyusut

Rinciannya di §7.7 sebagai blok koreksi. Yang menentukan: **kalibrasi** atas satu sel yang diketahui
hanya-baca (`Installments..AmountTotal`, kolom total hasil hitung) dan satu yang diketahui dapat
disunting (`InputTreatyInAdjustment..Ceding`). Keduanya terpisah bersih oleh `pyDisabledNew`:
`always` versus `false`.

Sebarannya **terbaca, bukan disimpulkan**: 146 mode `always` + 249 mode `true`-berkondisi = **395**,
tepat sama dengan jumlah `pyDisabled = true` di badan `Section/`; dan kedua kelompok tidak
bertumpang tindih — tidak ada `true` tanpa syarat, tidak ada `always` dengan syarat.

**Tiga butir dicabut dari daftar bocor `DB-10`,** dan dugaan pemilik proses tentang kolom total
benar: `Installments` `.AmountTotal` dan `.PctTotal`, empat belas sel `LayersEDM` termasuk
`.Deductible` dan `.AgregateLimit`, dan `DetailEGNPI` `.AmountIDR` dan `.Proportion` — seluruhnya
**nonaktif permanen**.

**Dan kalimat paling tajam saya terbalik.** *"`Layers` menjaga `.Deductible`; kembarannya
`LayersEDM` tidak"* — **`LayersEDM` justru lebih ketat**: empat belas selnya nonaktif permanen,
sedangkan `Layers` hanya menjaganya bersyarat. Kalimat itu dicabut.

**Yang tetap bocor, dan tetap memuat angka uang:** `TreatyInTabsNonProportional` 36 sel (termasuk
`.Deductible`, `.AggregateLimit`, `.Amount`), `TreatyRetroList` 26 (termasuk `.Deductible`,
`.Limit`, `.NetPremi`), `Share` dan `ShareRetro` 5 masing-masing (`.Value`, `.Pct`, `.Currency`),
`DetailLimits` 4 (termasuk `.Value` dan `.Currency`), dan tujuh seksi lain dengan 1-9 sel.

**Kesimpulan `DB-10` tidak berubah; daftarnya berubah.** Deductible dan batas agregat memang masih
dapat diubah sesudah `Revision` — lewat `TreatyInTabsNonProportional` dan `TreatyRetroList`, bukan
lewat `LayersEDM`. Total angsuran **tidak**.

### NB-07 · Tiga syarat E1 diperiksa

**1. Status `SEAM-ADJUSTMENT.md` §3.** Berkas itu menyatakan dirinya *"Sumber mengikat: ADR-0048
(arahan pemilik proses), ADR-0036, ADR-0039, ADR-0040, ADR-0045"* dan menetapkan *"apa yang Treaty
In **wajib** sediakan"*. Bannernya 23 September berbunyi *"embargo dicabut; **keputusan berkas ini
tetap berlaku**"*. Jadi ia **keputusan yang sudah diterima**, bukan usulan — dan E1 boleh turun ke
KONFIRMASI.

**2. Premis.** `KUNCI_PADANAN` **tidak** bersandar pada premis yang dibantah, dan berkas itu sudah
memisahkannya sendiri:

> Kalimat "DATA LAMA TIDAK PERNAH DISALIN" di bawah adalah **aturan untuk sistem baru** dan tetap
> berdiri; ia hanya tidak boleh dibaca sebagai pemerian sistem lama.

Kebutuhan memadankan baris antar versi berdiri **terlepas dari** apakah sistem lama menyalin atau
mengambil lewat `SELECT` — ia lahir dari kenyataan bahwa pengenal baris anak berbeda antar versi,
bukan dari G2. Satu catatan letak: `KUNCI_PADANAN` tinggal di `NILAI_VERSI_KONTRAK`, yang berkas itu
sebut **bentuk baca, diturunkan, bukan kanonik** — jadi ia bagian sisi baca seam, bukan atribut
kanonik entitas anak.

**3. Keunikan kunci adalah pertanyaan data** → **`UA-18`**: untuk setiap entitas anak, berapa daftar
di `JSONDATA` yang memuat dua baris dengan kunci padanan yang sama. `SEAM` §3 menetapkan bentuknya;
ia tidak menjamin data warisan mematuhinya.

### NB-08 · Selisih EGNPI pada addendum premi selalu nol

Rinciannya di §5.4 sebagai temuan bertanggal, dan di **TDA-16**. Bukti pemutusnya satu tag:
`TreatyEDMDifferencePremium` langkah 1 ber-`pyStepsObjectName = TreatyIn.EGNPI`, sedangkan
`TreatyEDMDifferenceShare` langkah 1 ber-`= TreatyIn.ActualValue`. Di berkas premi, kata
`ActualValue` muncul **satu kali saja**, sebagai deklarasi `pyPagesAndClasses` — tidak pernah
sebagai objek langkah maupun nilai.

**Ia menjadi syarat yang wajib dijawab C3**, persis seperti diminta pemilik proses: di bawah GRL-12,
versi `PENYESUAIAN_PREMI` yang hanya mengubah premi aktual akan terbaca **non-material**.

### NB-09 · Tiga belas sel yang berpindah golongan — sebabnya

Sensus lama menghitung 96 sel "dijaga `ViewState`"; sensus baru menghitung **83**. Totalnya sama
(978), jadi **tiga belas** berpindah. Sebabnya satu kalimat: **sensus lama membaca syarat
`ViewState` di mana pun di dalam subpohon sel, termasuk pada mode yang `pyDisabledNew`-nya bukan
`true`** — yaitu syarat yang tersimpan tetapi **tidak diaktifkan**, pola yang sama dengan
`pyVisible = ALWAYS` yang menyimpan kondisi tanpa mengevaluasinya.

Ke mana ketiga belasnya pergi: **sembilan** ke "selalu nonaktif" — `DetailLimitsOldData` (3),
`LayersEDM` `.Currency`, `LayersOldData` `.Currency`, `MaxRetentionOldData` (2),
`TreatyInTabsProportionalOldData` (2) — dan **empat** ke "nonaktif bersyarat oleh syarat lain",
yaitu `MaxRetention` `.Amount`/`.Currency` dan `TreatyInTabsProportional` `MaxCoGroup`/`MaxCoNonGroup`,
yang syarat aktifnya `TreatyIn.IsEditData` dan `EDMMaterialType`, bukan `ViewState`.

Tidak satu pun dari ketiga belas berada di daftar bocor `DB-10`, sehingga **tidak ada kesimpulan
yang berubah**.
