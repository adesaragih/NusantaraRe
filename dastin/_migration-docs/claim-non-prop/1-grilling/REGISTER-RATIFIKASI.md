# REGISTER RATIFIKASI — putusan teknis yang menunggu ratifikasi akuntansi

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` versi 2026-09-18 10:36, 48 objek, terurai ke `pengetahuan/ddl/`; dan `MEMORI_PEMAHAMAN.MD`.
> Disusun 18 September 2026.

**Delapan butir.** Seluruhnya berlabel **DECIDED-TEKNIS**: diturunkan dari XML dan DDL yang sudah ada, bukan dari angka baru, dan **bukan** dari jawaban akuntansi.

Aturan register ini:

- Kedelapan butir **berlaku penuh sekarang** dan tidak ditanyakan ulang.
- **Ratifikasi tidak menahan pekerjaan.** Akuntansi meratifikasi kemudian; sampai itu terjadi, pekerjaan berjalan di atas putusan ini.
- Bila ratifikasi mengubah salah satunya, yang berubah dicatat di baris butir itu beserta tanggalnya — butirnya tidak dihapus.

`ASK-AKUNTANSI.md` **ditutup** oleh register ini. Ia tetap ada sebagai catatan pertanyaan yang pernah diajukan beserta bukti berangkanya.

| # | Butir | Status ratifikasi |
|---|---|---|
| AK-1 | Presisi dan kelompok domain | belum |
| AK-2a | Nilai lama dimigrasi apa adanya | belum |
| AK-2b | Cacat perhitungan lama | belum |
| AK-3 | Toleransi shadow-run | belum |
| AK-4 | Tidak ada ambang penghentian migrasi | belum |
| AK-5 | Tarif | belum |
| AK-6.1 | Ambang kewenangan Komite | belum |
| AK-6.3 | Cut-off tutup buku | belum |

---

## AK-1 — Presisi

Skala perhitungan dan penyimpanan nilai antara = **20**, mengikuti `@divide(...,20)` pada **126 ekspresi** di seluruh rule inti klaim non-proporsional. Dasarnya `BLUEPRINT.md` §6.3: **skala mengikuti modul, bukan jenis nilai**.

| Kelompok | Isi | Tipe |
|---|---|---|
| **U1** | nilai uang antara | `NUMBER(38,20)` |
| **P1** | persentase dan porsi | `NUMBER(38,20)` |
| **K1** | kurs | `NUMBER(38,20)` |
| **L1** | limit layer dan premi deposit | `NUMBER(38,20)` |
| **P2** | tarif pajak dan brokerage | `NUMBER(11,8)` — skala 8, mengikuti `SetPPNPPH` |
| **C1** | cacah dan nomor urut | `NUMBER(9)` |

**Tidak ada kelompok U2.** Pembulatan dua desimal terjadi **hanya di tepi**: lapisan view kompatibilitas dan muatan keluar. Bila muatan keluar perlu diarsipkan, arsipnya **tabel tersendiri**, bukan kolom di tabel bisnis.

Alasannya satu ekspresi: `TempKasir.CARI9 = .TotalClaim - .PremiumSpreaded` (`HitServiceToKasir_Act`) — `.TotalClaim` adalah nilai antara **sekaligus** bahan instruksi bayar. Satu kolom tidak dapat berskala 20 dan 2 sekaligus.

**Akibat**: `ADR-0003` naik dari `proposed` ke **`accepted`**.

## AK-2a — Nilai lama dimigrasi apa adanya

Nilai lama dimigrasi **tanpa dibulatkan**. Pembulatan dua desimal berlaku hanya pada nilai yang **baru diterbitkan** sistem baru.

Preseden: ADR-0018 — klaim terdampak penjaga tanggal yang salah dimigrasi apa adanya, dan didaftar.

## AK-2b — Cacat perhitungan lama

Tidak diwarisi ke depan; angka lama tidak diubah ke belakang; selisihnya **didaftar**.

| Cacat | Perlakuan sistem baru | Bukti |
|---|---|---|
| `TotalUR` dikalikan `0` pada cabang mata uang sama | dihitung seperti Retensi Cedant biasa. **Tidak ada kolom yang selalu nol** | ADR-0010 tambahan 18 Sep; `(.Deductible * 0 * Local.ProrateClaim/100)` |
| Dua rumus premi pemulihan | acuan **`CountReinstatement_Act`**; `AdjClaimCNP_Act` baris **3109** tidak diwarisi | FINDING-007 |
| `TotalValueAdjust` memakai `=` bukan `+=` di dalam loop | sistem baru **menjumlahkan seluruh Adjustment** | FINDING-001 bagian 3 |

Ketiganya masuk **daftar pengecualian bernama** pada shadow-run — **bukan kegagalan cutover**.

## AK-3 — Toleransi shadow-run

| Golongan nilai | Toleransi |
|---|---|
| **Nilai yang masuk jurnal** — tujuh besaran `ASK-AKUNTANSI.md` §0.2: `TotalClaim`, `PremiumSpreaded`, `AdjustmentValue`, `AdjusterFeeValue`, `SalvageValue`, `IndividualRiskRNM`, `IndividualRiskPercentage` | **sama persis sampai 2 desimal, tanpa toleransi** |
| **Nilai antara** | 0,01% relatif **dan** Rp 1.000 mutlak; yang dilanggar lebih dulu yang berlaku. Setiap selisih di atas ambang dijelaskan **satu per satu** |
| **Nilai hasil konversi mata uang** | dibandingkan pada **mata uang asli**, bukan IDR |

Alasan yang ketiga: kurs lama tidak dapat direproduksi (ADR-0005 tambahan). Diperkuat FINDING-006 — `GETCURRENCYSTANDARD` mengembalikan `1` bila kurs tidak ditemukan, dan `1` tidak dapat dibedakan dari kurs yang sah.

## AK-4 — Tidak ada angka ambang penghentian migrasi

ADR-0014 sudah mewajibkan **setiap** baris yang tidak dapat diurai diselesaikan satu per satu, tanpa kecuali.

Migrasi berhenti bila ada **satu** baris yang tidak dapat diurai **dan** tidak dapat diselesaikan. **Tidak ada baris yang dibuang karena "cuma sedikit".**

## AK-5 — Tarif

Diisi ke `TARIF_BERLAKU` sebagai **satu baris**, berlaku sejak sebelum data tertua, dengan kolom lingkup bernilai **global**. Preseden: ADR-0025.

| Tarif | Nilai | Skala |
|---|---|---|
| brokerage | 2,5% | 8 |
| PPh | 2% | 8 |
| PPN | 2,2% | 8 |

Faktor **102,2 tidak disimpan sebagai angka**; ia diturunkan sebagai `(100 + tarif PPN)`. Di sistem lama ia tertanam sebagai `@divide(102.2,100,8)` — angka yang diam-diam memuat tarif PPN di dalamnya, sehingga perubahan tarif harus diingat di dua tempat.

**Tidak ada bukti tarif pernah berubah**, jadi tidak ada riwayat yang dibuat-buat. Satu baris, bukan rangkaian baris berantai.

## AK-6.1 — Ambang kewenangan Komite

`LimitPersenMax = LimitPersenMaxDivHead = 30,00` — apa adanya, `CreateChildKomiteCNP_Act` langkah 10.

Cabang `CekLimitPersen > 30 AND <= 30` adalah **cabang mati** dan **TIDAK DIMIGRASI** (FINDING-002 bagian 10). Jadi **dua jenjang, bukan tiga**.

Ambang nilai **30.000.000** dan **50.000.000**, disimpan sebagai **data**, dan dibandingkan terhadap **nilai IDR hasil konversi** — DECIDED(ADR-0007).

## AK-6.3 — Cut-off tutup buku

**Satu sumber kebenaran**: tabel `TUTUP_BUKU` bertanggal berlaku — DECIDED(ADR-0025).

Angka **25** yang tertanam di `HitServiceToKasir_Act` **tidak diwarisi**. Baris pertama diisi `25` karena itu satu-satunya angka yang terbaca.

Nilai mati di `PROC_GENERATE_SEQUENCE_NUMBER` — `TO_DATE('02/01/2026')` menghasilkan `'12.2025'` — menjadi **satu baris data**, bukan cabang kode.

---

## Batas klaim atas register ini

Kedelapan butir diturunkan dari XML dan DDL. Yang **tidak** dapat disimpulkan darinya, dan karena itu tetap terbuka:

| Hal | Sebabnya |
|---|---|
| Apakah tarif 2,5% / 2% / 2,2% **masih berlaku hari ini** | tarif tertanam di kode **tanpa tanggal berlaku**. Yang dibuat adalah tempatnya; nilainya menunggu ratifikasi — A8 |
| Apakah ambang 30 memang dimaksudkan, atau salah satunya seharusnya 15 | kode menetapkan keduanya `30.00`. Bahwa satu jenjang tidak pernah berjalan adalah akibat yang terbaca; apakah itu disengaja **TIDAK DITEMUKAN** |
| Presisi yang sesungguhnya diterima akuntansi selama ini | nilai uang klaim **IN-BLOB** dan keluar sebagai `varchar2`; basis data tidak pernah menyimpan presisinya. Profil data tidak akan mengubah kenyataan itu |
| Apakah `25` masih tanggal tutup buku yang berlaku | satu-satunya angka yang terbaca, dari dua tempat yang tidak saling tahu |
