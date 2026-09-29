# PROMPT — GILIRAN 15 *(folder `OUTPUT_HASIL_RNM`, cabang `main` @ `ce58a33` atau lebih baru)*: **menerapkan enam jawaban work owner — br, N7, N10, N11, PL-15, kosong-sebagai-nol di 7.8**

> Hanya konteks **Claim Life, PremiumList Life, Komite Claim Life**. GILIRAN-3 *(A/B/C)*, 4–14 tetap rujukan; aturan berhenti GILIRAN-6;
> **setiap pembacaan activity mencetak `pyStepsBlockName`**. Berkas `tco_*` dan worktree `.worktrees/desain-template` tidak disentuh;
> `App.tsx` berisi suntingan work owner — jangan di-commit atau dibuang. Bila sesi Treaty lanjutan 3 sedang berjalan di folder yang
> sama, **tunggu ia selesai** — kedua sesi menyentuh `migrasi_test.go`.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 29-09-2026

| Hal | Keadaan |
| --- | --- |
| GILIRAN-14 | 5 commit `055f552` → `df36af5`; uji `main`: **911 PASS · 0 FAIL**, dengan tag `db` **56 SKIP**, vitest **614** *(termasuk Treaty)* |
| Ralat brief 14 diterima | prasyarat langkah 7.8 = **b3919** *(b3631 milik 7.7)* |
| br | grid adjustment `ClaimLifeDetailGCNM.xml` b17126: **4** kolom data `Read-only`; `CLAIM_GROSS` dibaca 4 activity, **ditulis nol** |
| DEV | `T_MIGRASI` 29 langkah *(057 terpasang)*; `SEQ_WORK_POLIS` `last_number` **4**; `T_WORK_POLIS` berisi **2** kasus uji work owner `NBLF-2`, `NBLF-3`; nomor lama tertinggi yang terlihat **22373** *(`JSON_POLIS`/`POLICYJSONLIFE`, 33 baris ber-`NBLF-`; agregat)*; tabel penghitung Pega `PC_DATA_UNIQUEID` **tidak terlihat** dari akun `POOLDATA` |

## 1. JAWABAN WORK OWNER *(29-09-2026: "ikuti rekomendasi")*

| Butir | Keputusan | Kerjakan |
| --- | --- | --- |
| **br** | ikut XML: **tidak ada** sunting sel adjustment | tutup br dan OQ-N8 dengan blok bertanggal; `TestGridAdjustmentNolSelDapatDisunting` tetap; tidak ada rute `PUT` |
| **N7** | Delete baris adjustment **tidak berlaku** *(ADR-U-0031 melarang hapus fisik; tabel tanpa kolom penanda)* | tombol `Delete` b19120 tidak dirender *(bukan tombol mati)*; alasannya di PARITAS dan tiket 03; OQ-N7 ditutup |
| **N10** | **tiru** pembulatan empat angka langkah 7.7 | baca langkah 7.7 utuh *(b3600–b3671; rumus `@divide(@toDecimal(@replaceAll(x,",",".")),1,4)` per medan)*; terapkan pada nilai peserta yang disalin, dengan `bulat` bersama yang sudah ada; uji contoh literal; OQ-N10 ditutup |
| **N11** | tanya pemilik ekspor Pega; sementara **kosong** | `CLAIM_GROSS` tidak diisi aplikasi; baris OQ-N11 dipindah ke daftar **pemilik ekspor** dengan bukti *(4 pembaca, 0 penulis)* |
| **PL-15** | awal `SEQ_WORK_POLIS` **dimajukan** di atas nomor lama | migrasi **`058_seq_work_polis_mulai_ulang.sql`** *(+ down)*: `DROP` lalu `CREATE SEQUENCE SEQ_WORK_POLIS START WITH 22374` *(tertinggi terlihat + 1)* dengan komentar bukti; label `[sementara — DBA memastikan pyLastReservedID awalan NBLF- di PC_DATA_UNIQUEID sebelum data nyata]`; penjaga kata cadangan dan penghitung migrasi disesuaikan; `-migrate` oleh work owner. Dua kasus uji `NBLF-2`/`NBLF-3` di DEV **tidak** disentuh executor *(lihat §3)* |
| **Kosong = nol** | ikut Pega **hanya di langkah 7.8** | delapan medan baris adjustment yang lahir saat Submit: sumber kosong → `0` *(`@toDecimal("")`)*, dikunci uji; ADR-U-0027 tetap berlaku di tempat lain; penyimpangan terhadap ADR dicatat bertanggal di tiket 02 dan di komentar satu fungsi konversi |

## 2. URUTAN — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | br + N7 + N11 *(dokumen, tombol Delete tidak dirender)* | `claim-life: br/N7/N11 — sel adjustment baca-saja, Delete tidak berlaku` |
| 2 | N10 + kosong-sebagai-nol di 7.8 | `claim-life: N10 — pembulatan 7.7; 7.8 kosong dibaca nol seperti Pega` |
| 3 | PL-15 migrasi 058 | `premiumlist-life: PL-15 — SEQ_WORK_POLIS mulai 22374` |
| 4 | status tiket, panduan uji, OQ | `docs: status tiket dan OQ — giliran 15` |

## 3. UNTUK WORK OWNER

- Dua kasus uji `NBLF-2` dan `NBLF-3` sudah ada di `T_WORK_POLIS` DEV *(dari mencoba tombol Input Offer/Premium)*. Nomor itu mungkin sama
  dengan nomor kasus lama di Pega. Bila tidak dipakai lagi, minta asisten menghapusnya sesudah 058 terpasang.
- `-migrate` untuk 058 dijalankan **sesudah** sesi Treaty lanjutan 3 menyelesaikan paket 1 *(berkas 300–307 dibuang)*.

## 4. LAPORAN

Baris pertama alasan berhenti; tabel **butir → commit → rule XML → kode**; ralat tiket; OQ ditutup; angka uji tiap commit dengan dan tanpa
tag `db`; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 29 September 2026 dari jawaban work owner ("ikuti rekomendasi"), verifikasi `df36af5`, pembacaan grid b17126 dan penulis
`CLAIM_GROSS`, dan katalog DEV (`JSON_POLIS`/`POLICYJSONLIFE` nomor `NBLF-` tertinggi 22373, `SEQ_WORK_POLIS` last_number 4, `T_WORK_POLIS`
2 baris — agregat dan pengenal saja).*
