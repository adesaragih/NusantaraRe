# PROMPT — GILIRAN 18 *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang `main`, **sesudah** `PROMPT-REFACTOR-NAMA-MODUL.md` selesai)*: **PL-09 dari badan prosedur + OQ-N13 ikut XML** — dua butir terakhir yang dapat diputuskan

> Hanya konteks **Claim Life, PremiumList Life, Komite Claim Life**, pada struktur bentuk B *(`APP_RNM/inti`, `APP_RNM/modul/<nama>`)* dengan nama
> folder dari tabel nama. Aturan berhenti GILIRAN-6; **setiap pembacaan activity mencetak `pyStepsBlockName`**; commit dengan jalur eksplisit.

## 0. VERIFIKASI GILIRAN 17 *(asisten, 30-09-2026)*

Tujuh commit `5d1d646` → `9fc918f` ada; `DAFTAR-SERAH-TERIMA-TIGA-MODUL.md` ada. Uji pohon terkini *(sesudah refactor bentuk B)*: Go **957 PASS · 0
FAIL**, dengan tag `db` **59 SKIP**; vitest **678**; build **98** modul. DEV: `T_MIGRASI` **32** langkah — **021 dan 022 sudah terpasang**.

## 1. DUA BUTIR

| Butir | Bukti | Kerjakan |
| --- | --- | --- |
| **PL-09** `M_LIFE_PREMIUM_SUMMARY` | executor tidak dapat membaca `ALL_SOURCE`; **asisten membacanya**: penulis tabel itu = prosedur **`PEGA_M_LIFE_PREMIUM_SUMMARY`** *(135 baris; tabel 38 kolom di katalog DEV)*. Badannya, tanpa baris komentar, ada di **`.scratch/premiumlist-life/dba-procedure-PEGA_M_LIFE_PREMIUM_SUMMARY.md`** *(commit `526fc93`)* — parameter `P_COB`, `P_PL_NUMBER`, `P_PL_NUMBER_EDM`, `P_CURRENCY`, `P_PREMIUM`, `P_COMMISSION`, … `P_CLAIM`, `P_BALANCE`, dst. | tiru prosedur di Go *(keputusan o)*: tabel, kolom, bentuk teks nilai, cabang insert/update, di **transaksi simpan summary yang sama** *(pl2)*; nilai dari rumus rekap 05a yang sudah ada *(`.PREMIUM`/`.COMMISSION`/`.BALANCE` per cabang `Type`)*; kolom uang kosong = `0` *(PL-10)*; nol `COMMIT` di teks SQL; uji murni + `db` *(SKIP)*; PL-09 ditutup, dipindah **keluar** dari daftar serah terima DBA |
| **OQ-N13** status cermin di Save to RNM | `RDBList/InsertJsonKlaimLife_sql.xml` **b176** menulis `'0'` — Save to RNM di Pega **menyetel** status cermin `OS_AKSEPTASI_KLAIM_LIFE` | **ikut XML** `[keputusan asisten dari bukti; veto work owner]`: cermin ditulis `'0'` saat Save to RNM, sehingga klaim ganda antarklaim baru tertangkap tanpa menunggu Komite; larangan "tabel warisan baca-saja di Save to RNM" dari brief lama **dicabut untuk kolom ini saja**, dicatat bertanggal di tiket 03; uji `db` |

## 2. URUTAN

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | PL-09 | `premiumlist-life: PL-09 — M_LIFE_PREMIUM_SUMMARY seperti PEGA_M_LIFE_PREMIUM_SUMMARY` |
| 2 | OQ-N13 | `claim-life: N13 — Save to RNM menyetel status cermin '0' seperti InsertJsonKlaimLife_sql b176` |
| 3 | daftar serah terima + status tiket + OQ ditutup | `docs: giliran 18` |

## 3. LAPORAN

Satu pesan: tabel **butir → commit → bukti → kode**; angka uji dengan dan tanpa tag `db`; bab **TELEMETRI EKSEKUSI**. Tidak ada migrasi baru.

---

*Disusun 30 September 2026 dari `ALL_SOURCE` `PEGA_M_LIFE_PREMIUM_SUMMARY` (DEV, baca-saja), katalog `M_LIFE_PREMIUM_SUMMARY`,
`InsertJsonKlaimLife_sql.xml` b176, dan `T_MIGRASI` DEV.*
