# 11: Kontrak hilir ke Claim Life — satu tabel, dua sudut pandang

> **Ralat 01-10-2026** (gelombang 2 brief, `../RALAT-DEV-01-10-2026.md` — ralat mengalahkan isi di bawah). Teks lama yang tidak berlaku:
> - **E2** — *"Penegakannya milik konteks Claim Life"* → Claim Life sudah menyaring; tiket ini menjadi: Endorsement **menulis** `M_LIFE_PREMIUM_DETAIL` persis `SaveMasterLPDet` (`EDMSTATUS` Old/New/Delete/Batal) dan summary seperti `InsertPLSummary`, nol perubahan kode Claim Life/PremiumList, uji meniru penyaring Claim Life.
> - **E4** — kueri ke `M_LIFE_PREMIUM_DETAIL` → selalu berkunci ber-index yang disebut (`_INDEX4` `PL_NUMBER`, `_INDEX21` `PL_NUMBER_EDM`), nol pemindaian penuh.

**Status:** sebagian 01-10-2026 — rekap warisan + uji penyaring Claim Life (c83bf68); peserta warisan dan uji kontrak HTTP lintas modul menunggu OQ-EDM-016

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 09 (baris — termasuk yang negatif — harus sudah tertulis)

## Hasil & nilai pengguna

Sebagai **Finance**, saya ingin seluruh baris endorsement terlihat di laporan saya agar nettonya
benar; dan sebagai **admin klaim Life**, saya ingin peserta yang sudah dibatalkan atau dihapus
**tidak muncul** saat saya mencari peserta untuk klaim — sehingga saya tidak pernah memproses klaim
atas peserta yang tidak lagi ditanggung. *(User story 41–43 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Penulisan baris ke tabel bersama; **kueri kontrak** bergaya baca-klaim |
| `internal/services` | Kontrak yang dipanggil konteks Claim Life |
| `internal/handlers` | Endpoint pencarian rekam premium endorsement |
| — | Penandaan **kontrak lintas konteks** di tempat bentuk rekam didefinisikan |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SaveMasterLPDet.xml` **dan** `PremiumList Life/RDBList/SaveMasterLPDet.xml` | **penulis bersama** — satu rule identik |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | kedua modul | penulis summary |
| `GetPesertaClaim_sql1` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/GetPesertaClaim_sql1.xml` | **konsumen** |
| `LoadDataPesertaSpesifik_Act` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `LOADDATAPESERTASPESIFIK_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/LoadDataPesertaSpesifik_Act.xml` (67.369 byte) | pemanggil konsumen (baris 545) |

### Kolom penanda — **terbaca dari korpus** `[terverifikasi]`

Daftar kolom `INSERT` + klausa `VALUES` di `SaveMasterLPDet`:

| Kolom | Diisi dari | Nilai | Penanda hidup/mati? |
| --- | --- | --- | --- |
| **`EDMSTATUS`** | `TempValue.EDMStatus` | `Old` / `New` / `Delete` / `Batal` | ✅ **ya** |
| `STATUS` | `CARI48` | `0` untuk `QR`/`QP`, `1` untuk `TP`/`TR` | ❌ penanda **jenis transaksi** |
| `STATUSOLD` | `CARI47` | `1`/`0` | ❌ |

⚠️ `[terverifikasi]` **Jalur new business tidak mengisi `EDMSTATUS` sama sekali** — sensus
`PremiumList Life/Activity/InsertLifePremiumDetail_act.xml` (`ASM-FW-GISFW-WORK-LIFE` /
`INSERTLIFEPREMIUMDETAIL_ACT`): **nol** kemunculan. Baris NB masuk dengan `EDMSTATUS` kosong/NULL.

**Peserta hidup** = `EDMSTATUS` **kosong/NULL**, `Old`, atau `New`. **Mati** = `Delete` atau `Batal`.

## ADR terkait

**ADR-0001** (batas konteks — perubahan bentuk rekam = perubahan kontrak), **ADR-0011** (bentuk
`PremiumListSummary` / `PremiumListDetail` yang dipakai mesin status klaim), **ADR-0003**.

## Acceptance criteria

- [ ] Baris endorsement — **termasuk yang bernilai negatif** — tertulis ke `M_LIFE_PREMIUM_DETAIL`
      dan `M_LIFE_PREMIUM_SUMMARY`, **tabel yang sama** dengan jalur new business. *(AC 18 spec)*
- [ ] Peserta ber-`EDMSTATUS` **`Batal`** atau **`Delete`** **tidak muncul** pada jalur baca klaim.
      *(AC 48 spec; `[keputusan work owner]`)*
- [ ] Peserta **new business** (`EDMSTATUS` kosong/NULL) **tetap muncul**. ⚠️ Penyaring naif
      `EDMSTATUS NOT IN ('Delete','Batal')` membuang seluruh peserta NB di Oracle — test wajib memuat
      kasus ini dan **harus gagal** bila penyaringnya naif. *(AC 48a spec)*
- [ ] `STATUS` dan `STATUSOLD` **tidak** dipakai sebagai penanda hidup/mati. *(AC 48b spec)*
- [ ] Jalur baca **akuntansi/ringkasan premium TIDAK menyaring** — ia melihat **seluruh** baris,
      positif maupun negatif. **Satu tabel, dua sudut pandang**, dan itu disengaja. *(AC 49 spec)*
- [ ] Rekam premium hasil endorsement dapat ditemukan lewat **`PL_NUMBER_EDM`** maupun
      **`PL_NUMBER`**. *(AC 50 spec)*
- [ ] Nilai uang yang ditulis hulu dibaca hilir **identik**, termasuk nilai negatif — tidak ada
      pembulatan di perbatasan. (**ADR-0003**)
- [ ] Bentuk kedua rekam ditandai di kode sebagai **kontrak lintas konteks**; mengubahnya memaksa
      pembaruan sadar di sisi Claim Life. *(AC 51 spec; **ADR-0001**)*
- [ ] Ada **test kontrak** yang menembus dari simpan endorsement sampai pencarian bergaya Claim Life

### Hilir membaca tabel ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Claim Life membaca peserta dari **tabel relasional**, bukan dari CLOB JSON. *(AC 54 spec;
      penyimpangan sadar 8)*
- [ ] ⚠️ Baris ber-`EDMSTATUS` **`Delete`** dan **`Batal`** **tidak pernah** muncul sebagai kandidat
      peserta klaim — keduanya tetap tersimpan, tetapi tersaring di jalur baca klaim. *(AC 66, 67
      spec; kontrak §14)*
- [ ] ⚠️ Baris bernilai **negatif** hasil jurnal balik **tidak pernah** sampai ke layar klaim maupun
      ke perhitungan klaim, sementara jalur akuntansi tetap melihat **seluruh** baris. *(kontrak §14)*
      — **satu seam**, API HTTP, terhadap **skema uji Oracle nyata**.

## Catatan — penegakan penyaring ada di konteks Claim Life

⚠️ `[terverifikasi]` `GetPesertaClaim_sql1` hari ini **tidak memuat penyaring status apa pun**;
sensus modul `Claim Life/` menemukan **nol** kemunculan `EDMSTATUS`. Apakah Pega menyaring di
lapisan lain **tidak terbukti dari korpus** — **jangan menebak mekanismenya**.

**Penegakannya milik konteks Claim Life:** spec Claim Life **§16** + **AC 25–30**, dan tiket
**CL-02**. Tiket ini mengikatnya sebagai **syarat kontrak** yang dibuktikan lewat test kontrak
lintas konteks. Diikat juga di **PL-08**.

## Blocker

**Tidak ada pemblokir.** ⚠️ Satu catatan yang **tidak** memblokir: tipe dan nullability `EDMSTATUS`
belum terbaca (`NULL` atau string kosong) — **OQ-001 (sisa)**, pemilik **DBA**. Sampai dipastikan,
implementasi menangani **keduanya** (`IS NULL` *atau* `= ''`).

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Status 01-10-2026 — selesai (K5 keputusan work owner 01-10-2026, OQ-EDM-016)

Kalimat lama *"sebagian 01-10-2026 — rekap warisan + uji penyaring Claim Life (c83bf68); peserta warisan dan uji kontrak HTTP lintas
modul menunggu OQ-EDM-016"* tidak berlaku lagi.

| AC | Bukti |
| --- | --- |
| Baris endorsement — termasuk negatif — ke `M_LIFE_PREMIUM_DETAIL` **dan** `M_LIFE_PREMIUM_SUMMARY` | `TulisPesertaWarisan` + `TulisRekapWarisan`; db: baris `Delete` `GROSS_PREMIUM_REFUND` −7,25 |
| `Batal`/`Delete` tidak muncul di jalur baca klaim; NB (NULL) tetap muncul | `TestPutuskanMenulisPesertaWarisan` (penyaring Claim Life DITIRU, tidak diimpor) + db `TestPutuskanTerhadapOracle` (`COUNT` hidup 2 dari 3) + `models.TestStatusPesertaSejalanPenyaringClaimLife` (kesetaraan dengan teks `penyaringHidup` Claim Life) *(diperjelas code review 01-10-2026: yang dibuktikan adalah baris yang DITULIS modul ini — uji `db` menghitung baris ber-`IDPEGA` kasus itu saja)* |
| `STATUS`/`STATUSOLD` bukan penanda hidup | `models.TestStatusJenisBukanPenandaHidup` |
| Jalur akuntansi tidak menyaring | rekap 12 menjumlah seluruh peserta kasus (tidak berubah) |
| Ditemukan lewat `PL_NUMBER_EDM` maupun `PL_NUMBER` | kedua kolom ditulis (`:1`, `:2`) |
| Test kontrak simpan endorsement → pencarian bergaya Claim Life | seam layanan (tiruan) dan seam repository terhadap Oracle (db, SKIP tanpa `ORACLE_DSN`); seam HTTP tiruan: `TestRutePutuskan` memeriksa `pesertaWarisan` di jawaban Confirm. ⚠️ Seam HTTP terhadap Oracle **tidak** dibangun: modul ini belum punya uji `db` di lapisan handler |

⚠️ **Temuan, ditiru persis:** `EM_PERCENT` ← `{TempInputDetail.CARI49}` dan `RISK` ← `{TempInputDetail.CARI50}` (`SaveMasterLPDet`
b249/b250) tidak pernah ditetapkan `InsertJsonPolisLife_Act` 11.1 (CARI12 ditetapkan `.EM_PERCENT` lalu ditimpa
`.GROSS_PREMIUM_REFUND`) — baris endorsement Pega di tabel warisan ber-`EM_PERCENT`/`RISK` NULL, dan begitu pula yang ditulis
modul ini. Jalur new business PremiumList mengisi keduanya.

⚠️ **Risiko hilir untuk pemilik Claim Life (code review 01-10-2026, tidak diperbaiki di modul ini):** seperti Pega,
Confirm hanya MENYISIP baris versi baru (`SaveMasterLPDet` satu-satunya rule korpus Endorsement yang menyentuh
`M_LIFE_PREMIUM_DETAIL`; nol `UPDATE`/`DELETE` atas baris versi lama). Baris versi sebelumnya — new business ber-`EDMSTATUS`
NULL, atau `Old`/`New` endorsement lama — tetap HIDUP bagi penyaring Claim Life. Akibatnya pencarian peserta klaim berkunci
`PL_NUMBER` (+ `CERTIFICATE_NO`, `AmbilUntukKlaim` `FETCH FIRST 1 ROWS ONLY` tanpa `ORDER BY`) dapat memilih baris pra-endorsement,
dan sertifikat yang di-`Delete` endorsement tetap terlihat lewat baris NB-nya. Ini sudah berlaku untuk endorsement buatan Pega
di DEV; pemilihan versi di jalur baca klaim milik konteks Claim Life (kontrak §14) — diserahkan, bukan ditebak di sini.

