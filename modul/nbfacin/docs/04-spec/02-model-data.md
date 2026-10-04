# Spesifikasi Model Data — Siklus New Business

> Sumber: `02-layar\01-skema-field-nb.md` (432 Section NB, 6.314 field terikat) dan
> `01-activity\01-inventaris-activity-nb.md`, dibangun dari `D:\migrasi\RNM\NB FacIn\`.
> Label mengikuti `CLAUDE.md` §3. Rancangan ditandai **[usulan]**.

---

## 1. Kendala yang mendahului semua rancangan

### 1.1 ⛔ Tidak ada satu pun rule `Property` — seluruh tipe data **belum terverifikasi**

`[terverifikasi]`

```powershell
Test-Path "D:\migrasi\RNM\NB FacIn\Property"                                    # -> False
(Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern '<pyPropertyName>' -List).Count  # -> 3
```

Korpus memuat **1.805 ekspresi properti unik** (1.096 nama daun) tanpa satu pun deklarasi tipe.
Satu-satunya petunjuk tipe adalah **kontrol UI** yang mengikatnya — dan petunjuk itu **saling
bertentangan di 150 ekspresi**.

Bukti terkuat, diverifikasi langsung:

```powershell
$c = [IO.File]::ReadAllText("D:\migrasi\RNM\NB FacIn\Section\CoverageSpreadingList_IsUW.xml")
[regex]::Matches($c, '<pyValue>\.ASMDateOfBirth</pyValue>').Count   # -> 2
# pyFormat di sekitar pengikatan pertama : pxDateTime
# pyFormat di sekitar pengikatan kedua   : pxInteger
```

Properti yang sama, **berkas yang sama**, dua tipe kontrol yang tidak mungkin keduanya benar.

**Konsekuensi [usulan]:** setiap field di model Go membawa komentar tipe dengan status
`belum terverifikasi` sampai DDL Oracle diperoleh. Untuk 150 ekspresi bertentangan, tipe **tidak
boleh dipilih** oleh migrasi — itu pertanyaan bisnis/DBA.

### 1.2 ⛔ Tidak ada aturan validasi selain wajib-isi

`[terverifikasi]` Di seluruh lapisan Section NB: **nol** aturan panjang, format, rentang, atau pola.
Satu-satunya validasi yang terbaca adalah wajib-isi (477 sel / 191 properti / 73 Section), dan tipe
rule `Edit Validate` / `Validate` **tidak ada di ekspor NB**.

**Konsekuensi:** sistem baru **tidak dapat menurunkan validasinya dari korpus**. Menambahkan validasi
yang masuk akal terdengar benar, tetapi akan menolak data yang hari ini diterima sistem lama — dan
itu perubahan perilaku, bukan migrasi (`CLAUDE.md` §1).

### 1.3 ⛔ 604 dropdown tanpa satu pun daftar nilainya

`[terverifikasi]` 604 sel `pxDropdown`/`pxRadioButtons`, dan korpus tidak memuat enumerasinya.
**Tidak ada layar isian yang dapat dibangun lengkap** sebelum daftar pilihan diperoleh.

### 1.4 ⚠️ Mata uang sering tidak menyertai nilai uang

`[terverifikasi]` Hanya 93 pengikatan properti mata uang di seluruh NB, sementara **112 Section
menampilkan nilai uang tanpa field mata uang mana pun**.

Ini konfirmasi independen atas `CLAUDE.md` §4.1: **korpus memang kehilangan informasi mata uang**.
Karena itu `Money.Currency` wajib di target — bukan karena kerapian, melainkan karena sumbernya
memang tidak membawanya dan harus ditentukan di batas input.

### 1.5 ⚠️ Presisi desimal berbeda untuk properti yang sama

`[terverifikasi]` 20 dari 439 properti desimal punya `pyDecimalPlaces` berbeda antar layar:
`.Premium` 0/2/4 · `.SharePercentage` −1/2/4 · `.RateLife` 10/4 · `.LimitofLiability` 4/−999.

⛔ Arti sentinel `−1` dan `−999` **belum terverifikasi**. Dan presisi mana yang **tersimpan** (bukan
sekadar ditampilkan) tidak dapat dijawab dari lapisan Section.

---

## 2. Bentuk agregat — satu inti + lima bagian opsional

`[terverifikasi]` Dari 1.096 nama daun unik: **963 khas satu lini bisnis, hanya 133 bersama.** Rasio
itu menolak rancangan satu struct tunggal untuk semua lini.

`[usulan]` Yang didukung bukti adalah **satu inti + lima bagian opsional**:

```go
// Nama field JSON MENGIKUTI nama properti Pega apa adanya — agregat diserialkan ke
// JSON_POLIS.DATA_JSON dengan nama identik; mengubahnya memutus pembacaan data historis.
type OfferFacIn struct {
    Quotation QuotationData `json:"QuotationData"`
    Policy    PolicyData    `json:"PolicyData"`

    // Lima bagian opsional. Yang aktif ditentukan lini bisnis (§3).
    Property *PropertyPart `json:"PropertyList,omitempty"`
    Cargo    *CargoPart    `json:"CargoList,omitempty"`
    Vehicle  *VehiclePart  `json:"VehicleList,omitempty"`
    Person   *PersonPart   `json:"PersonList,omitempty"`
    Aneka    *AnekaPart    `json:"AnekaList,omitempty"`
}
```

⚠️ Nama field JSON di atas **[dugaan]** — harus dikunci terhadap sampel `DATA_JSON` nyata sebelum
dipakai. Korpus tidak memuat contoh dokumennya.

> ### 📌 Temuan 17 September 2026 — sampelnya kemungkinan sudah ada
>
> `[terverifikasi]` Lima berkas kasus di `D:\migrasi\RNM\DDL\` (`P-5 NB-181231 (FIRE)`,
> `P-5 NB-184233 ( MARINE CARGO)`, `P-5 RNW-10579 (FIRE)`, `P-5 EDM-13445 (FIRE)`,
> `P-5 NB-176005 ( AS. KREDIT )`) berformat **JSON**, berawalan `{"CedingRetention":"0",…`.
>
> `[dugaan]` Kelimanya **kemungkinan sampel `DATA_JSON` nyata** — persis yang dibutuhkan untuk
> menaikkan status nama field agregat di atas dari `[dugaan]` menjadi `[terverifikasi]`.
>
> ⛔ **Belum dikerjakan, dan sengaja.** Kelimanya memuat **data pelanggan** (nama tertanggung, alamat)
> dan **nama operator** (`pxCreateOpName` di kelimanya). Pembacaannya untuk mengunci skema hanya boleh
> dilakukan **setelah de-identifikasi** — lihat tiket de-identifikasi pada rencana `05-tickets\`.
> Dicatat di sini agar temuannya tidak hilang, bukan sebagai pekerjaan yang tertunda diam-diam.

---

## 3. Pemilih lini bisnis

`[terverifikasi]` Lini bisnis ditentukan **satu properti**: `Quotation.BusinessType`, diverifikasi
dari `pyUnmodifiedPath` rule `When`.

⛔ **Tetapi dua rule membacanya dari jalur berbeda** di dalam agregat — `IsMarineCargo` dan
`IsCustomBonds` memakai 2–3 jalur. **[pertanyaan terbuka] MEMBLOKIR:** apakah ketiga jalur selalu
sinkron? Bila tidak, dua kasus yang sama diklasifikasikan berbeda tergantung rule mana yang jalan.

Sampai dijawab, `BusinessType` **tidak boleh** dinormalisasi menjadi satu field di model — ketiga
jalur dipertahankan dan perbedaannya dicatat saat runtime.

⚠️ `[terverifikasi]` Sebagian nilai `BusinessType` muncul **dengan spasi di depan/belakang**.
Normalisasi di batas input, dan catat sebagai kandidat perbaikan.

---

## 4. Daftar bersarang

`[terverifikasi]` 629 RepeatGrid di 281 Section, mengikat **181 PageList unik** → slice di Go,
komponen tabel di React.

⚠️ **Kedalaman bersarang tidak terukur langsung.** Penyarangan grid di dalam satu berkas = **0**;
ekspor tidak memuat isi Section yang disisipkan. Kedalaman hanya terukur lewat graf penyisipan:
**121 pasangan, 33 induk**. Jadi angka kedalaman apa pun bersifat `[dugaan]`.

---

## 5. Read-only: 76,6 % field

`[terverifikasi]` 4.834 dari 6.314 field terikat bersifat read-only (`pyEditOptions=Read-only`).

Ini memberi tahu sesuatu tentang bentuk aplikasi: sebagian besar layar adalah **tampilan**, bukan
isian. Rancangan React yang memperlakukan semuanya sebagai form akan salah ukuran.

⚠️ Terpisah dari itu: 183 field **selalu** read-only dan 106 field **tidak pernah tampil**
(kondisi yang selalu benar/salah). Apakah itu disengaja adalah **keputusan bisnis**, bukan keputusan
migrasi.

⚠️ **Catatan metodologis:** ejaan huruf tidak seragam menyembunyikan angka. Regex peka huruf memberi
367/94 kondisi mati; dengan `(?i)` hasilnya **407/105** — `never` huruf kecil muncul 40 kali,
`Always` 13 kali. Audit apa pun atas kondisi mati **wajib** case-insensitive.

---

## 6. Komponen bersama UW / non-UW

`[terverifikasi]` 81 dari 87 Section `*_IsUW` punya kembaran non-UW dengan skema field praktis
identik.

`[usulan]` Satu komponen React dengan prop `mode`, bukan dua pohon komponen. Perbedaan yang tersisa
(6 pasangan) diperiksa satu per satu sebelum digabung.

---

## 7. Yang **tidak dapat** dispesifikasikan sekarang

| Bagian | Penghalang |
| --- | --- |
| Tipe setiap field | Tidak ada rule `Property`; 150 ekspresi berkontrol bertentangan; korpus nol DDL |
| Isi dropdown | 604 sel tanpa enumerasi di korpus |
| Validasi | Tidak ada aturan selain wajib-isi; tipe rule validasi tidak ikut terekspor |
| Presisi tersimpan | Lapisan Section hanya menunjukkan presisi **tampilan** |
| Skema `DATA_JSON` | Hanya 11 jalur terbaca dari SQL; dokumen sesungguhnya tidak terdokumentasi |
| Mata uang per nilai | 112 Section menampilkan uang tanpa field mata uang |

**Yang dapat dimulai:** `pkg/money`, kerangka agregat inti + lima bagian opsional (§2), dan
komponen React read-only (76,6 % field).

---

*Lampiran A–C di `02-layar\01-skema-field-nb.md` memuat skema per Section, katalog 1.805 ekspresi
properti, dan seluruh RepeatGrid — dipakai sebagai rujukan saat menulis struct, bukan disalin.*
