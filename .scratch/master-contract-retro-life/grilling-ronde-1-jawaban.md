# Grilling — Master Contract Retro Life — Ronde 1: JAWABAN FRONTIER

Tanggal: 2026-09-15
Konteks: `master-contract-retro-life` — **konteks/menu sendiri** `[keputusan work owner]`
Modul: **Master Contract Retro Life** (66 berkas — modul terkecil di korpus)
Sumber jawaban: sesi grilling dengan work owner (Ronde 1).

> **Konvensi.** `[terverifikasi]` = terbukti korpus (class + nama + path);
> `[keputusan work owner]` = keputusan desain disetujui work owner;
> `[fakta bisnis — work owner]` = fakta bisnis dinyatakan work owner (bukan tebakan);
> `[terbuka]` = OQ, diserahkan ke tim (DBA / Product+UW).

Dokumen ini menutup **Frontier Ronde 1** dari `grilling-ronde-1.md`.
**Seluruh 11 pertanyaan terjawab. Frontier kosong.**

---

## Rekap jawaban 11 frontier

| # | Jenis | Jawaban ringkas |
| --- | --- | --- |
| Q1 | desain | Normalkan — `REINSTYPEID` sekali di kontrak, anak mewarisi; `REINSTYPENAME` cache turunan |
| Q2 | fakta bisnis | `B_*`=batas bawah, `*`=batas atas, `*_SELISIH`=atas−bawah (lebar layer) |
| Q3 | fakta bisnis | (a) share tingkat-2 = % **dari share reinsurer induk** (retrosesi atas retrosesi) |
| Q4 | desain | Boleh simpan meski ≠100% + tanda mencolok + laporan; TIDAK diblokir |
| Q5 | desain | **KASKADE + popup konfirmasi** (hapus induk → anak ikut terhapus; konfirmasi sebut jumlah) — *direvisi dari "tolak"* |
| Q6 | fakta bisnis | (a) tahun treaty **abadi** — sengaja tak bisa dihapus |
| Q7 | fakta bisnis | 5 jenis life ID 10196–10200; `Flag=1`=for life; nama dari `Note` |
| Q8 | desain | Pertahankan aturan tahun, banding tanggal asli (bukan potong teks) |
| Q9 | desain | Panggil 5 procedure Oracle apa adanya |
| Q10 | desain | Bawa "terapkan ke semua" + konfirmasi/pratinjau + audit |
| Q11 | desain | Wajib cek `HASIL1`; gagal tampilkan galat (gagal terang-terangan) |

---

## Rincian per pertanyaan

### Q1 — Denormalisasi `REINSTYPEID` / `TREATYYEAR` `[keputusan work owner]`

**Keputusan: NORMALKAN.** Simpan `REINSTYPEID` **hanya di contract**; reinsurer dan business
mewarisinya lewat `TREATYCONTRACTID`. `TREATYYEAR` disimpan di year saja; business mewarisi.

**Pengecualian sadar:** `REINSTYPENAME` (hasil lookup dari `.Note`) boleh disimpan sebagai nilai
tampilan **di-cache/turunan** — ditandai jelas sebagai turunan, **bukan** sumber kebenaran.

Alasan: korpus `[terverifikasi]` menyimpan `REINSTYPEID`/`REINSTYPENAME` **3×** (contract, reinsurer,
business) dan `TREATYYEAR` **2×** (year, business), tanpa mekanisme penjaga konsistensi. Tiga salinan
dapat menyimpang.

> ⚠️ **Penyimpangan sadar dari Pega.**

---

### Q2 — Enam kolom limit pada kontrak `[fakta bisnis — work owner]`

Kolom: `IDR`, `USD`, `B_IDR`, `B_USD`, `IDR_SELISIH`, `USD_SELISIH`
(`[terverifikasi]` parameter `INSERTTREATYCONTRACT_LIFE`).

| Kolom | Arti |
| --- | --- |
| `B_IDR` | **batas bawah** limit (IDR) — prefix `B_` = **bawah** |
| `B_USD` | **batas bawah** limit (USD) |
| `IDR` | **batas atas** limit (IDR) — tanpa `B_` = **atas** |
| `USD` | **batas atas** limit (USD) |
| `IDR_SELISIH` | **atas − bawah** (IDR) |
| `USD_SELISIH` | **atas − bawah** (USD) |

**Tidak ada "batas atas versi B"** — prefix `B_` khusus batas bawah.

**Makna selisih = lebar layer proteksi berjenjang.** Batas atas satu layer = batas bawah layer
berikutnya. Susunan menaik:

```
QS  →  2nd QS  →  Surplus  →  2nd Surplus
```

Selisih per layer:
- **2nd QS − QS**
- **Surplus − 2nd QS**
- **2nd Surplus − Surplus**

**Konsekuensi desain (perbaikan sadar):** karena `SELISIH = atas − bawah` **turunan**, di sistem baru
**dihitung, jangan disimpan mentah** — agar tak menyimpang dari kedua batasnya. (Sistem lama menyimpan
kolom terpisah — risiko tidak sinkron.)

> ⚠️ **Penyimpangan sadar dari Pega** (selisih dihitung, bukan disimpan).

**Sisa untuk DBA (OQ-001):** tipe/presisi kolom uang; dan konfirmasi kewajiban `USD` (di Pega
`SaveTreatyLimit_Act` `[terverifikasi]`: `IDR`, `B_IDR`, `B_USD` **wajib**, `USD` **tidak** — apakah
aturan sah (transaksi utama IDR) atau kelalaian validasi).

---

### Q3 — Dua tingkat share `[fakta bisnis — work owner]`

**Jawaban: (a).** Share `TREATYSECURITYREINSURER_LIFE.PCTSHARE` = **porsi dari share reinsurer
induknya** (`TREATYREINSURER_LIFE.PCTSHARE`), **bukan** dari keseluruhan treaty.

**Retrosesi atas retrosesi.** Contoh: Reinsurer A dapat 40% treaty; Security Reinsurer X di bawah A
= 10% → eksposur X terhadap treaty = **10% × 40% = 4%**.

`[terverifikasi]` Bukti hierarki: `INSERTSECURITYREINSURER_LIFE` menerima `p_TREATYREINSURERID`
(menunjuk baris reinsurer induk), COMMIT sendiri di procedure.

**Konsekuensi:** perhitungan eksposur harus mengalikan berjenjang (share anak × share induk).

---

### Q4 — Total share: validasi 100%? `[keputusan work owner]`

**Keputusan: TIDAK diblokir.** Total share ditampilkan mencolok; jika ≠ 100% diberi **tanda
peringatan "belum 100%"**, tetapi penyimpanan **tetap diperbolehkan**.

**Tambahan:** laporan "kontrak dengan total share ≠ 100%" agar celah terlihat.

`[terverifikasi]` Sesuai temuan korpus: `CountingPercentShare_Act` hanya **menjumlahkan** `PCTSHARE`
→ taruh di field `STDRATING` untuk **ditampilkan**; tidak ada perbandingan ke 100 di seluruh modul.
Validasi yang ada hanya batas per baris `0..100` (`SetErrorMessageReinsurer`).

---

### Q5 — Hapus induk yang punya anak `[keputusan work owner]` — **DIREVISI**

**Keputusan FINAL: (b) KASKADE + KONFIRMASI SEBELUM HAPUS (popup opsi 1).** Menghapus induk
**ikut menghapus semua anaknya** (kontrak → reinsurer + security reinsurer + business; reinsurer →
security reinsurer). **Popup konfirmasi muncul LEBIH DULU** (sebelum apa pun terhapus), menyebut
**apa & berapa** yang akan ikut terhapus (mis. "Menghapus kontrak ini akan ikut menghapus 3 reinsurer
dan 5 business. Lanjut? [Ya] [Batal]"). **Ya → kaskade jalan; Batal → tidak ada yang terhapus.**
`[keputusan work owner — popup opsi 1: konfirmasi sebelum hapus]`

> **Revisi dari keputusan awal.** Sebelumnya dicatat (c) tolak-bila-punya-anak. Work owner mengubah
> menjadi **(b) kaskade + popup konfirmasi**. Alasan work owner: hapus induk memang dimaksudkan
> menghapus seluruh sub-pohonnya.

`[terverifikasi]` Korpus **tanpa kaskade** — 4 `DELETE … WHERE ID=…` datar; menghapus induk
meninggalkan anak yatim. Sistem baru **memperbaiki** dengan kaskade sadar + konfirmasi.

`[keputusan work owner]` DB kini punya **FK antar tabel** → mode `ON DELETE` harus **`CASCADE`**
agar sejalan (DB ikut menghapus anak). Bila FK terpasang RESTRICT, harus diubah ke CASCADE.

> ⚠️ **Penyimpangan sadar dari Pega** (Pega hapus datar tanpa kaskade & tanpa konfirmasi; sistem
> baru kaskade + konfirmasi).

---

### Q6 — Tahun treaty tidak dapat dihapus `[fakta bisnis — work owner]`

**Jawaban: (a) DISENGAJA.** Tahun treaty **abadi** setelah dibuat — memang tidak boleh dihapus.
Dijadikan **aturan eksplisit** di spec (bukan kelalaian).

`[terverifikasi]` Korpus: 5 penulis, hanya 4 penghapus — `TREATYYEAR_LIFE` tanpa jalur hapus.

---

### Q7 — Enumerasi jenis reasuransi `[fakta bisnis — work owner]` (menutup OQ-057 untuk life)

Tabel master: class Pega `ASM-FW-GISFW-Int-REINSURANCETYPE`.
`[terverifikasi]` RD `BrowseReinsuranceTypeLimit_RD` menyaring `.Flag = 1`, menampilkan `.Note`,
urut `.ID` DESC.

**Q7a — daftar nilai (life):**

| ID | Note (nama tampil) |
| --- | --- |
| 10196 | QS |
| 10197 | 2ND QS |
| 10198 | SURPLUS |
| 10199 | 2ND SURPLUS |
| 10200 | OR |

→ **5 jenis** life, rentang ID **10196–10200**. ("OR" dicatat apa adanya — tidak ditafsirkan.)

**Q7b — `Flag = 1` = "for life"** (penanda lini life; **bukan** aktif/nonaktif, **bukan** untuk
limit). Cocok dengan memo korpus `[terverifikasi]` `pyUsage: "Parameter Flag, 1 for life"`.

**Q7c — beda kolom:**
- `.Note` = **nama tampil** (dipakai UI) ✅ `[terverifikasi]` + work owner.
- `.Flag` = **penyaring**, `=1` for life ✅.
- `.Code`, `.Type` = **tidak dipakai** di modul ini → **OQ kecil DBA** (isi tabel), tidak memblokir.

**Sisa untuk DBA (OQ kecil):** nama tabel fisik Oracle (korpus hanya menyebut *class*, bukan nama
tabel); isi kolom `.Code` dan `.Type`.

---

### Q8 — Gerbang konsistensi tahun `[keputusan work owner]`

**Keputusan: PERTAHANKAN aturan, PERBAIKI cara.** Aturan tetap: "tahun pada tanggal-mulai kontrak
harus = tahun treaty induk". Tetapi bandingkan **tahun dari nilai tanggal (tipe date)**, bukan
memotong teks.

`[terverifikasi]` Pega memakai `@substring(TREATYSTARTDATE,6,10) <> TREATYYEAR_LIFE` (potong teks
posisi 6–10) di `SaveBusinessLife_Act`, `SaveSecurityLife_Act`, `SaveSecurityReinsurerLife_Act` —
rapuh terhadap perubahan format.

> ⚠️ **Penyimpangan sadar dari Pega** (banding tanggal asli, bukan substring).

---

### Q9 — Lima procedure penulis `[keputusan work owner]`

**Keputusan: PANGGIL APA ADANYA.** Go memanggil 5 stored procedure Oracle yang ada; **tidak** ditulis
ulang. Konsisten dengan 3 konteks Life sebelumnya & ADR-0006.

Prosedur `[terverifikasi]`:
- `INSERTTREATYYEAR_LIFE`
- `INSERTTREATYCONTRACT_LIFE`
- `INSERTREINSURER_LIFE`
- `INSERTSECURITYREINSURER_LIFE`
- `INSERTBUSINESS_LIFE`

> ✅ **OQ-002 DITUTUP** `[data DBA]` (`procedure-bodies-from-dba.md`). Body 5 procedure = **upsert
> dikunci `ID` + audit**, ID dibuat DB via sequence (`'1'+lpad(seq,6)`), `TGLUPDATE`=`SYSDATE`,
> `COMMIT` internal, **nol logika bisnis**. `HASIL1`/`o_message` = pesan galat teks (berisi HTML),
> kosong/NULL = sukses. Seluruh aturan bisnis wajib di lapisan Go.

---

### Q10 — Fitur "terapkan ke semua" `[keputusan work owner]`

**Keputusan: BAWA + PENGAMAN.** Fitur `SaveBusinessToAllLife_Act` (`[terverifikasi]` menyimpan ke
semua baris `.REINSTYPEID == Param.REINSTYPEID`) dibawa, tetapi dengan:
- **konfirmasi + pratinjau** jumlah baris yang akan terpengaruh sebelum eksekusi;
- hasil **tercatat di jejak audit**.

> ⚠️ **Penyimpangan sadar dari Pega** (Pega tanpa konfirmasi/audit).

---

### Q11 — Keluaran galat procedure (`HASIL1`) `[keputusan work owner]`

**Keputusan: TEGAKKAN.** `HASIL1` **wajib diperiksa** setelah tiap pemanggilan procedure; kegagalan
**ditampilkan** ke user. Penyimpanan yang ditolak DB **tidak boleh** tampak berhasil.

`[terverifikasi]` Korpus: 5 procedure mengembalikan `HASIL1 out`; **tidak satu pun** activity `Save*`
membacanya (hanya `DeleteRowBusiness`). Galat jatuh diam-diam.

> ⚠️ **Penyimpangan sadar dari Pega** ("gagal terang-terangan", sekeluarga dengan tutup buku
> PremiumList Life). `[data DBA]` `HASIL1`/`o_message` = **pesan galat teks berisi HTML**;
> kosong/NULL = sukses, non-empty = gagal. Di Go: strip HTML, sajikan pesan bersih.

---

## Penyimpangan sadar (rekap) — semua akan jadi AC bertanda ⚠️ di spec

| # | Penyimpangan | Dari |
| --- | --- | --- |
| 1 | Normalisasi `REINSTYPEID`/`TREATYYEAR` (sekali di induk, anak mewarisi) | Q1 |
| 2 | `*_SELISIH` dihitung, bukan disimpan mentah | Q2 |
| 3 | Kaskade hapus induk→anak + **popup konfirmasi Ya/Batal sebelum hapus** (bukan tolak) | Q5 |
| 4 | Gerbang tahun via banding tanggal asli (bukan substring) | Q8 |
| 5 | "Terapkan ke semua" + konfirmasi/pratinjau + audit | Q10 |
| 6 | `HASIL1` wajib diperiksa, gagal terang-terangan (strip HTML) | Q11 |
| 7 | `ID` PRIMARY KEY di kelima tabel — ✅ selesai di DB (work owner) | DDL |
| 8 | FK penjaga antar tabel — ✅ selesai di DB (work owner); `ON DELETE CASCADE` selaras Q5 | DDL |

---

## Status OQ setelah Ronde 1

| OQ | Status | Pemilik | Isi |
| --- | --- | --- | --- |
| **OQ-002** | ✅ **DITUTUP** `[data DBA]` | — | Body 5 procedure diterima: upsert dikunci `ID`+audit, ID via sequence, `COMMIT` internal, `o_message`=pesan galat (HTML), nol logika bisnis. `procedure-bodies-from-dba.md` |
| **OQ-001** | ✅ **DITUTUP** `[data DBA]` | — | DDL 5 tabel: uang/share=`NUMBER`, `RIRATE`=teks, semua nullable, nol constraint, `ID` bukan PK, seq START 44. `ddl-tables-from-dba.md` |
| **OQ-057** | **DITUTUP untuk life** | — | Nilai enumerasi jenis reasuransi life = ID 10196–10200 (Q7a); `Flag=1`=for life (Q7b) |
| OQ kecil | `[terbuka]` | DBA | Nama tabel fisik `REINSURANCETYPE`; isi kolom `.Code` & `.Type` |
| OQ-054 | `[terbuka]` | — | Berkas non-Pega di `Claude outputs/` — tidak dibaca, dicatat |

---

## Kesiapan `/to-spec` — SEMUA PERTANYAAN SELESAI

- **Frontier kosong** — 11 pertanyaan terjawab.
- **7 keputusan desain** (Q1, Q4, Q5, Q8, Q9, Q10, Q11); Q5 final = **kaskade + popup konfirmasi
  Ya/Batal sebelum hapus**.
- **4 fakta bisnis** dinyatakan work owner (Q2, Q3, Q6, Q7).
- **8 penyimpangan sadar** siap jadi AC bertanda ⚠️ (2 terakhir sudah dieksekusi di DB).
- **Pemblokir: NIHIL.** ✅ OQ-002 (body procedure) & OQ-001 (DDL + PK + FK) **ditutup** dari DBA.
- **Data DBA tambahan:** lihat `procedure-bodies-from-dba.md` (perilaku simpan) dan
  `ddl-tables-from-dba.md` (tipe kolom, nullability, PK, FK).
- **OQ tak-memblokir tersisa:** isi kolom `.Code`/`.Type` + nama tabel fisik `REINSURANCETYPE`
  (tidak dipakai sistem) — dibiarkan terbuka atas keputusan work owner.

**Semua tiket dapat `ready` (termasuk migrasi).** Siap dipicu: `/mattpocock-skills:to-spec`
(spec ditulis ulang bersih dari file ini + data DBA).
