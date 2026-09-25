# PROMPT TO-SPEC — NB Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
>
> ⛔ **Prasyarat:** `PROMPT-VERIFIKASI-KEADAAN-NB-TREATY-IN.md` sudah dijalankan, dan
> `.scratch\nb-treaty-in\VERIFIKASI-KEADAAN.md` sudah ada dengan putusan **layak dipakai**.
> Bila berkas itu belum ada, **berhenti** dan jalankan ronde verifikasi lebih dulu.
>
> ⚠️ **Skill `to-spec` tidak dapat dipanggil sendiri oleh agen** (CLAUDE.md §8). Ketikkan
> `/to-spec` sebagai manusia, lalu berkas ini menjadi briefnya. Bila skill tidak dipakai, tulis
> spec langsung dalam bentuk rumah yang dijelaskan di Bab 3.

---

## 0. LINGKUP — DIKUNCI

Hanya modul **`D:\XML\RNM_BRD\NB Treaty In`**.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang** — jangan dibaca, dikutip, dibandingkan, atau dijadikan sasaran.

Keluaran satu berkas: **`OUTPUT_HASIL_RNM\.scratch\nb-treaty-in\spec.md`**.

⛔ `grilling-ronde-1..4.md` tersegel — jangan disunting.

---

## 1. SUMBER DAN URUTAN KEWENANGAN

Bila dua sumber bertentangan, yang di atas menang. **Tuliskan pertentangan yang Anda temukan.**

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `.scratch\nb-treaty-in\PERTANYAAN-untuk-*.md` — enam lembar jawaban | keputusan work owner, mengikat |
| 2 | `.scratch\nb-treaty-in\VERIFIKASI-KEADAAN.md` | hasil pemeriksaan; **menang atas nomor 3** di mana berbeda |
| 3 | `.scratch\nb-treaty-in\KEADAAN-NB-TREATY-IN.md` | keadaan terukur |
| 4 | `.scratch\nb-treaty-in\grilling-ronde-1..4.md` | latar; **tidak berlaku** di mana bertentangan dengan 1–3 |
| 5 | korpus XML | selalu boleh dipakai untuk membuktikan ulang |

Berkas grilling memuat angka yang sudah basi — itu disengaja, koreksinya ada di nomor 3 Bab 5.
Jangan menyalin angka dari sana tanpa memeriksa silang.

Contoh bentuk rumah: sebelas `spec.md` di `.scratch\*\` — ukurannya 40–143 KB.

---

## 2. KEADAAN YANG DIWARISI

Ringkas, supaya tidak perlu ditemukan ulang. Rinciannya di sumber nomor 3.

| Hal | Angka |
| --- | ---: |
| Berkas dalam modul | 278 |
| `Activity` terjangkau / yatim | 64 / 28 |
| Langkah `Property-Set` terjangkau / yatim | 377 / 561 |
| Pertanyaan terjawab | 47 dari 49 |
| Penahan tersisa | **P1** dan **P18** |

Folder modul adalah **dependency closure** — ia memuat aturan milik modul lain karena pewarisan
kelas. Dua puluh delapan aturan yatim adalah keluarga spreading/premi milik Fac In, dan tidak
dimigrasi.

Sumber data pindah dari kolom dokumen `JSONDATA` ke view relasional
`POOLDATA.TREATYINDETAILJOINEDM` — 39 kolom, dan 33 dari 33 medan yang dipakai laporan tersedia
di sana.

---

## 3. BENTUK SPEC

Bab wajib, berurutan:

```
Cara membaca berkas ini      Problem Statement      Solution
User Stories                 Implementation Decisions
Testing Decisions            Acceptance Criteria
Out of Scope                 Butir [terbuka] — daftar penuh
Further Notes                Lampiran
```

Jangan mengejar ukuran. Modul ini terbesar yang pernah digarap, tetapi lingkupnya sudah menyusut
60 % — kejar kelengkapan.

### Blok ringkasan di "Cara membaca berkas ini"

Wajib memuat cacah: user story · acceptance criteria · sebaran penanda · butir `[terbuka]` aktif ·
butir `[penyimpangan sadar]`.

**Dihitung dua cara yang berbeda**, dan sebutkan jendela hitungnya. Bila Anda menulis angka sebelum
mengukurnya lalu memperbaikinya — **kutip angka lamanya, jangan dihapus**. Itu aturan rumah
(CLAUDE.md §4a).

### Aturan menulis Acceptance Criteria

Bab itu **tidak memutuskan apa pun**. Ia menyatakan ulang keputusan yang sudah ada di bab
sebelumnya, dalam bentuk yang dapat diuji dari luar. Bila sebuah butir terasa seperti keputusan
baru, ia salah tulis.

Bentuknya:

```
N. `[terverifikasi]` <pernyataan>. Test yang menemukan <keadaan sebaliknya> **gagal**. *(Bab X)*
```

Setiap butir membawa **penanda** dan **rujukan bab**. Butir tanpa penanda adalah cacat.

---

## 4. YANG WAJIB ADA DI ACCEPTANCE CRITERIA

### Dua belas ketetapan yang mengikat implementasi

Seluruhnya dari Bab 4 berkas keadaan, masing-masing menjadi satu atau lebih butir yang dapat diuji:

1. `IsApproved`: `0` = ditolak, selain itu = disetujui; aturan hidup di `DecisionTable`, bukan
   `When`; dibandingkan sebagai **teks** *(P24, P6)*
2. Admin menolak → berkas diselesaikan sebagai ditolak. Atasan menolak → kembali ke admin *(P24)*
3. Nilai uang berpresisi penuh; pembulatan hanya di titik penyajian *(P29)*
4. `DEDUCTION1` `DEDUCTION2` `BROKERAGE` `RNM_SHARE` adalah **persentase**; `12.5` = 12,5 persen *(P29)*
5. `COMMENCEMENT` `TERMINATION` dibaca apa adanya; P32 tetap berlaku untuk penulisan baru *(P29)*
6. Setiap query menulis skema `POOLDATA.` eksplisit *(P3)*
7. Seluruh urutan penyimpanan dibungkus **satu transaksi** *(P2)*
8. `OPERATORID` dari identitas akses login; `PIC` dari nama tampilan *(P4, P33)*
9. Pencarian berdasarkan nomor urut antrean diganti pemeriksaan keanggotaan *(P25)*
10. Bila keterangan rule `When` berbeda dari syarat yang dijalankan, **yang dijalankan benar** *(P23)*
11. 38 medan tetap terkunci; 27 medan wajib tetap wajib *(P45, P47)*
12. `TypeTax` dan potongan 2,2 % ditiru apa adanya, termasuk ketidakseragaman presisi *(P46)*

### Ditambah

- **27 medan wajib, berbeda menurut tingkat.** Layar Dept Head tidak mewajibkan enam medan yang
  wajib di layar admin — `ClaimPaymentType` `ClaimType` `IDCurrency` `Quartal` `TypeTax`
  `YearOfQuartal` — tetapi mewajibkan `ResultOnp1` yang tidak wajib di sana.
- **38 medan terkunci permanen**: 36 di layar Dept Head, 2 di layar biasa.
- **Sembilan medan wajib sekaligus terkunci** di layar Dept Head, dan akibatnya: selama **P18**
  kosong, layar itu **tidak dapat disimpan**.
- **Empat cabang putusan alur**: admin terima → naik · admin tolak → selesai sebagai ditolak ·
  atasan terima → naik · atasan tolak → kembali ke admin.
- **36 baris penggolong** `BusinessType_DeT`, 128 kode bisnis, bawaan `"UNKNOWN"`, berhenti di
  baris pertama yang cocok.

---

## 5. YANG WAJIB MASUK OUT OF SCOPE

| Yang dikeluarkan | Sebab | Butir |
| --- | --- | --- |
| 28 aturan yatim, 561 langkah | tidak terjangkau; milik Fac In | keadaan Bab 2 |
| 6 aturan pembongkar JSON | penyimpanan pindah ke view relasional | P29, P15 |
| rantai perhitungan uang | isi 268 langkah belum terkirim | **P18** |
| 20 nomor polis dalam aturan uang | pindah ke modul Fac In | P20 |
| dua aturan uang yang hanya berupa catatan | pindah ke modul Fac In | P21 |
| pencarian berdasarkan nomor urut antrean | diganti pemeriksaan keanggotaan | P25 |
| rule `When\isApproved` | bukan aturan yang hidup | P6 |

Setiap butir menyebut **sebabnya** dan **butir asalnya**. Jangan ada yang dikeluarkan tanpa alasan
tertulis.

### Bab karantina rantai uang

Tulis tersendiri, memuat tiga hal:

1. apa yang **sudah** diketahui bentuknya — medan, sumber data, arah aliran;
2. apa yang **menunggu P18** — rumus di dalam 268 langkah;
3. apa **akibatnya bila P18 tidak pernah dijawab** — termasuk bahwa layar Dept Head tidak dapat
   disimpan.

---

## 6. DISIPLIN

Setiap pernyataan membawa bukti: **jalur berkas + tipe rule + nama rule**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`.

**Jangan menutup pertanyaan terbuka sendiri.** Penutupan milik work owner, DBA, Product &
Underwriting, Aktuaria, Finance, atau IAM.

**Jangan menyalin nilai berupa nama orang.** Nomor polis **tidak** ditulis apa adanya di dalam spec.
Nol rahasia, nol token, nol data nasabah, nol cuplikan data produksi.

Uang tidak pernah `float` (CLAUDE.md §7). Arah ketergantungan `handlers → services → repository`
(§5).

### Empat jebakan yang sudah pernah menjerat

1. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>`.
2. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
3. `pyStepsPage`, `pyCriteriaValue`, `pyResult`, `pyShapeName`, `pyTaskLabel` **tidak ada** di
   ekspor ini.
4. Menyatakan sesuatu nihil tanpa menyebut lingkup penelusuran.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.

---

## 7. BAB WAJIB — TELEMETRI EKSEKUSI

`spec.md` ditutup dengan bab `## TELEMETRI EKSEKUSI`.

Angka token sejati tidak terlihat dari dalam sesi. Ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Catat token keluaran · cache-read · jumlah panggilan alat · durasi · biaya · byte dibaca.
Bila pengukuran dari luar tidak dilakukan, **katakan begitu** — jangan menaksir lalu menyajikannya
sebagai angka terukur.

---

## 8. YANG MENANDAKAN RONDE INI BERHASIL

- `spec.md` lengkap sebelas bab, seluruh AC berpenanda dan berujuk bab.
- Dua belas ketetapan Bab 4 seluruhnya terwakili di Acceptance Criteria.
- Rantai uang dikarantina jelas dalam babnya sendiri, bukan ditebak atau dihilangkan.
- Setiap butir Out of Scope menyebut sebab dan butir asalnya.
- Blok ringkasan dihitung dua cara, dengan jendela hitung disebutkan.
- Pertentangan antar sumber — bila ada — ditulis, bukan didiamkan.
- Bab telemetri terisi angka terukur, bukan taksiran.

---

*Disusun 22 September 2026, sesudah empat ronde grilling, 47 dari 49 pertanyaan terjawab, dan satu
perbaikan berkas korpus yang menyusutkan lingkup modul sebesar 60 %.*
