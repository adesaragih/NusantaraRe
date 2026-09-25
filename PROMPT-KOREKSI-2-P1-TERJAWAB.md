# PROMPT KOREKSI KEDUA — P1 terjawab, penahan habis

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
> ⚠️ **Ronde ini MENYUNTING berkas.** Laporkan setiap suntingan.

---

## 0. LINGKUP — DIKUNCI

`D:\XML\RNM_BRD\` **READ-ONLY**. Menyunting hanya di dalam `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

### Berkas yang TIDAK boleh disunting — dan sebabnya berbeda-beda

| Berkas | Sebab |
| --- | --- |
| `grilling-ronde-*.md` | **tersegel** — isi ronde tidak disunting sesudah selesai |
| `VERIFIKASI-KEADAAN.md`, `VERIFIKASI-P18.md`, `KOREKSI-P18-DIJALANKAN.md` | **laporan ronde** — catatan apa yang benar **saat itu**, bertanggal, bukan dokumen hidup |
| `PROMPT-*.md` | **arsip perintah** — menyunting brief lama menghapus jejak apa yang diperintahkan |

⇒ Pernyataan basi di ketiga kelompok itu **dibiarkan**. Yang diperbarui hanya **dokumen hidup**:
lembar jawaban, `spec.md`, tiket, dan `KEADAAN-NB-TREATY-IN.md`.

---

## 1. SEBABNYA — P1 dijawab work owner 22 September 2026

Naskah **empat** stored procedure diterima, ditambah **tiga contoh** `DATA_JSON` sungguhan.
Rinciannya pada jawaban **P1** dan **P29** di `.scratch\nb-treaty-in\PERTANYAAN-untuk-DBA.md`.

**Penahan proyek kini nol.**

### Empat butir yang berubah karenanya

| Butir | Yang berubah |
| --- | --- |
| **P1** | tertutup |
| **P51** | ⛔ **diperbaiki** — penghapusan ternyata **bersyarat**: tidak berjalan bila `STS_KONVERSI = 1` |
| **P29** | lingkup mengecil — data kontrak **sudah relasional** di `POOLDATA.TREATY_IN` (20 kolom) |
| **P54** | ditegaskan — `TANGGAL_CLOSING` dibaca `WHERE ROWNUM = 1`, **satu baris berlaku global** |
| **P2** | ditegaskan tepat — dua meng-commit sendiri, penomoran tidak, penghapus meng-commit |

---

## 2. DAFTAR KERJA

### 2.1 Cabut pernyataan "P1 menahan" dari dokumen hidup

Terakhir dicacah: **20 pernyataan di 12 berkas**. Enam di antaranya berada di kelompok yang tidak
boleh disunting — biarkan.

Yang **harus** diperbarui:

| Berkas | Perkiraan |
| --- | ---: |
| `.scratch\nb-treaty-in\PERTANYAAN-untuk-DBA.md` | 2 |
| `.scratch\nb-treaty-in\PERTANYAAN-untuk-Product-dan-Underwriting.md` | 1 |
| `.scratch\nb-treaty-in\PERTANYAAN-YANG-MASIH-KOSONG.md` | 1 |
| `.scratch\nb-treaty-in\PERTANYAAN-RONDE-4.md` | 1 |
| `.scratch\nb-treaty-in\issues\12-layar-jenjang-ketiga.md` | 2 |
| `.scratch\nb-treaty-in\issues\13-rantai-perhitungan-uang-dikarantina.md` | 3 |
| `.scratch\nb-treaty-in\issues\00-skema-penyimpanan-dan-migrasi.md` | 1 |

**Cacah ulang sendiri.** Angka di atas dari satu pola pencarian; mungkin ada yang terlewat.

### 2.2 Perbaiki P51 — penghapusan bersyarat

Jawaban P51 kini menyatakan penghapusan berjalan pada kegagalan, tanpa menyebut penjaganya.
Naskah procedure menunjukkan penjaganya ada:

```
hanya menghapus bila  STS_KONVERSI IS NULL  atau  STS_KONVERSI <> '1'
lingkup selalu        WHERE IDPEGA = <satu kasus>
Bisnis = 'T'          JSON_POLIS · TREATYINPRODUCTION · TREATYINPRODUCTION_BACKUP
Bisnis = 'F'          JSON_POLIS · FACINPRODUCTION · FACINPRODUCTION_BACKUP · FACOUTPRODUCTION
```

Tambahkan sebagai **penajaman**, bukan pembatalan — keputusan work owner tetap berlaku.
⭐ Yang baru diketahui: **data yang sudah berhasil dikonversi tidak dapat tersentuh.**

### 2.3 Dua butir `[terbuka]` baru

| Butir | Isinya |
| --- | --- |
| **Tabel `_BACKUP` ikut dihapus** | `TREATYINPRODUCTION_BACKUP` dan `FACINPRODUCTION_BACKUP` dibersihkan bersama data utamanya. Namanya menyiratkan jaring pengaman, tetapi **tidak ada pemulihan**. Perlu dipastikan work owner |
| **Tanggal ditulis mati di penomoran** | `IF TRUNC(v_now) <= TO_DATE('02/01/2026') THEN v_mm_yyyy := '12.2025'`. Seluruh penomoran sebelum 2 Januari 2026 dipaksa ke periode Desember 2025. Masih dikehendaki? Diperlukan lagi saat pindah ke 2027? |

Catat pada lembar pemilik yang tepat, **jangan ditutup sendiri**.

### 2.4 Tiket `00-skema-penyimpanan-dan-migrasi`

Statusnya `needs-info` karena menunggu P1 dan P29. **Keduanya kini terjawab.**

Perbarui: yang ditunggu sekarang bukan lagi naskah procedure maupun contoh JSON, melainkan
**sensus properti** — lihat Bab 3. Sebutkan juga yang sudah pasti: data kontrak sudah relasional,
tujuh dari delapan kolom `json_polis` sudah datar, hanya `DATA_JSON` yang dipecah.

---

## 3. TEMUAN BARU YANG WAJIB DICATAT — dokumen JSON tidak berbentuk tetap

**Tiga** contoh `DATA_JSON` dari kasus berbeda dibandingkan:

| | Contoh 1 | Contoh 2 | Contoh 3 |
| --- | ---: | ---: | ---: |
| Jenis | proporsional | proporsional *(SOA)* | **non-proporsional XOL** |
| Medan skalar | 37 | 64 | 47 |
| **Ada di ketiganya** | | **23** | |
| Gabungan tiga | | **74** | |

Daftar bersarangnya berbeda-beda:

```
1 : LocationList > OccupationList > AnekaList > CoverageList > (ClauseList, DeductibleList)
2 : ListInstallment (datar) · SpreadingRiskList
3 : ListInstallment > InstallmentList · TreatyXOLList > ValueList · OldData > TreatyXOLList
    BreakDownSpreadList
```

⛔ **`ListInstallment` bersarang di contoh 3 tetapi datar di contoh 2** — nama daftar sama,
kedalaman berbeda. Pengurai yang menganggapnya seragam akan patah.

⭐ **Penentunya terlihat:** `QuotationData.ProportionalType` bernilai `"Proportional"` lawan
`"NonProportional"`, dan `IsNewPolicyNonProp` bernilai `"0"` lawan `"1"`. Bentuk dokumen mengikuti
jenis treaty.

⭐ Contoh 3 juga memunculkan **`OldData`** — halaman data lama yang selama ini dikenal dari EDM —
**di dalam berkas polis baru**, berisi `TreatyXOLList` kosong.

⭐ **Dan korpus mengenal lebih banyak lagi**: penelusuran awal menemukan **97 properti** pada
`PolicyTreatyIn` dan **sembilan** daftar bersarang — termasuk `TreatyXOLList` (dirujuk 228 kali),
`TreatyXOLDifferenceList` (123), `OldData` (258), `TreatyDifference` (97), yang **tidak muncul di
ketiga contoh**.

⚠️ Angka 97 itu **batas bawah** — dari satu pola pencarian saja. Jangan diperlakukan sebagai
sensus. Sensus sebenarnya adalah ronde tersendiri.

**Catat sebagai `[terbuka]` pada P29:** rancangan tabel **tidak dapat** dibangun dari contoh;
ia harus dibangun dari sensus properti korpus, dengan contoh sebagai **penguji**.

### Dan satu uji silang yang memperkuat P46

Contoh 2 ber-`TypeTax` = `"Inclusive"`, `Deduction1` = `32.516`,
`BrokerageFeeSebenarnya` = `31.8160`.

```
32.516 / 1.022 = 31.81604696673…  ->  dibulatkan 4 desimal  =  31.8160   COCOK
```

⇒ Potongan 2,2 % pada **P46** kini **terverifikasi dari data produksi**, bukan hanya dari
pembacaan rumus. Catat pada jawaban P46.

### ⛔⛔ Dan satu temuan yang menyentuh keputusan uang — WAJIB dicatat

Contoh 3 memperlihatkan **galat bilangan mengambang yang sudah tersimpan permanen di data
produksi**:

```
Premium angsuran   148157378.220000069     x 4  =  592629512.880000276
NetPremium         592629512.880000276          <- sama persis
selisih dari nilai bulat dua desimal            =  2,76 x 10^-7
```

Ekor `0002760` **bukan presisi** — itu galat pembulatan yang berlipat: premi dibagi empat angsuran,
tiap angsuran membawa galatnya, lalu dijumlahkan kembali.

⚠️ **Ini menyentuh ketetapan P29 nomor 3** *(presisi penuh, pembulatan hanya di titik penyajian)*.
Diterapkan apa adanya, sistem baru **melestarikan galat ini selamanya**, dan tiap penjumlahan
menambahnya.

**Catat sebagai `[terbuka]` pada P29, dengan tiga pilihan dan akibat masing-masing:**

| | Pilihan | Akibat |
| --- | --- | --- |
| **a** | ikuti apa adanya | paritas sempurna dengan sistem lama, cacatnya diwariskan |
| **b** | bulatkan saat migrasi | bersih ke depan, **angka historis berubah**, laporan lama tidak cocok |
| **c** | simpan apa adanya, bulatkan saat dihitung | sejarah utuh, galat berhenti berlipat |

⛔ **Jangan memilih sendiri.** Ini keputusan work owner, dan presisi uang yang wajar per mata uang
adalah pertanyaan **Finance**.

⭐ ADR-0003 tetap ditegakkan apa pun pilihannya: sistem baru tidak memakai `float`, sehingga tidak
menambah galat baru.

⛔ **Jangan menyimpan contoh JSON apa adanya.** Ketiganya memuat nama orang pada medan pemasar,
operator, dan tertanggung. Catat nama medan dan jumlahnya saja.

---

## 4. YANG TIDAK BOLEH DIKERJAKAN

- ⛔ Jangan menyunting berkas pada tabel Bab 0.
- ⛔ Jangan menutup butir `[terbuka]` yang tidak disebut di sini.
- ⛔ Jangan merancang tabel. Itu ronde berikutnya, sesudah sensus properti.
- ⛔ Jangan membuat ADR baru. Bila perlu, katakan di laporan.
- ⛔ Nol kode, nol DDL, nol nama orang.

---

## 5. SESUDAH MENYUNTING — WAJIB DIPERIKSA ULANG

| Uji | Yang diharapkan |
| --- | --- |
| cacah pertanyaan NB, dua cara | **49 dari 49** — nol terbuka |
| cacah pertanyaan EDM | 11 dari 11 |
| sinkronisasi tiga salinan lembar pertanyaan | nol nomor berbeda isi |
| tiket `blocked` | hanya tiket 13, **karena P30** — bukan P1 maupun P18 |
| "P1 menahan" di dokumen hidup | **nihil** |
| butir `[terbuka]` aktif di `spec.md` | dicacah ulang dua cara; angka lama **dikutip**, bukan dihapus |

---

## 6. KELUARAN

`.scratch\nb-treaty-in\KOREKSI-2-DIJALANKAN.md`, lima bab:
suntingan · P51 dan P46 ditajamkan · dua butir terbuka baru · hasil enam pemeriksaan ·
**TELEMETRI EKSEKUSI**.

Telemetri diukur dari luar bila bisa:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Bila tidak, katakan begitu. Bila hasil pengurangan baseline, katakan itu juga.

---

## 7. YANG MENANDAKAN RONDE INI BERHASIL

- Nol pernyataan "P1 menahan" tersisa di dokumen hidup; yang di arsip dan laporan **dibiarkan**.
- P51 **ditajamkan**, bukan dibatalkan — keputusan work owner tetap berlaku.
- P46 membawa uji silang dari data produksi.
- Dua butir terbuka baru tercatat pada lembar pemilik yang tepat, **tidak ditutup sendiri**.
- Tiket `00` menyatakan dengan tepat apa yang masih ditunggu: **sensus properti**, bukan P1.
- Keenam pemeriksaan dijalankan, hasilnya ditulis apa adanya.
- Nol contoh JSON tersimpan, nol nama orang tersalin.

---

*Disusun 22 September 2026, sesudah naskah empat stored procedure dan tiga contoh `DATA_JSON`
diterima — dan penahan terakhir proyek ini habis.*
