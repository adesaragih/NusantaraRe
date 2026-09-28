# PROMPT — GILIRAN 11 *(**sesi BARU** di `OUTPUT_HASIL_RNM`, cabang `main` @ `251cb3b`)*: **paket 0 status tiket jujur → Claim Life `Save to RNM` (SaveOutStandingLife_Act) → av-2 21 medan polis → bj/bk → sisa temuan → PANDUAN UJI LAYAR tiga modul**

> Tiga modul sudah membangun seluruh tiketnya *(PremiumList 00–09, Komite 00–09, Claim Life 01–15 + §3.1)*. Giliran ini menutup
> celah terbesar yang tersisa dan menyiapkan uji layar oleh manusia. GILIRAN-3 *(A/B/C)*, 4–10 tetap rujukan; aturan berhenti
> GILIRAN-6 dan sesi baru tiap giliran *(induk §0.5)* berlaku. Sesi Treaty Contract Out berjalan **terpisah** di worktree-nya;
> sesi ini tidak menyentuh berkas `tco_*`.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Klaim laporan GILIRAN-10 | Diperiksa ulang | Hasil |
| --- | --- | --- |
| urutan brief selesai; berhenti karena selesai, bukan konteks | 25 commit `0cddee0..main`; pohon bersih kecuali `App.tsx` milik work owner | ✅ |
| Go 613 · 0 · **0 SKIP**; vitest 372 | dijalankan ulang di worktree bersih: **613 PASS · 0 FAIL**; tanpa tag `db` 0 SKIP, **dengan** tag `db` **38 SKIP** *(skema uji DBA belum ada)*; vitest **372**; **66** modul; vet, gofmt, tsc bersih | ✅ *(angka SKIP bergantung tag — laporan berikut menyebut keduanya)* |
| nol migrasi baru | `0cddee0..main` tanpa berkas migrasi baru; DEV `T_MIGRASI` 28 langkah, tetap lengkap | ✅ |
| sensus §3.1: 31 ada/tidak perlu, 4 celah dibangun, 6 OQ | `OQ-untuk-tim.md` memuat OQ-M1…M9 | ✅ |
| ⛔ **status tiket** | PremiumList: 01, 02, 05a, 06, 08 masih **`ready-for-agent`** meski dibangun; Komite: 01–09 **seluruhnya `ready-for-agent`**; kotak AC **0 tercentang** di hampir semua tiket kedua modul; Claim Life 15 tiket masih **`claimed`** | ⛔ paket 0 |
| celah terbesar | rute simpan Outstanding *(`SaveOutStandingLife_Act`, tombol `Save to RNM` `InputOSClaimLife.xml` b21102 → b21126/b21234)* belum dibangun; **av-2**: layar mengikat 21 medan `PolicyDataLife`, panel menampilkan 11 | ⛔ paket 1–2 |

## 1. KEPUTUSAN GILIRAN INI

| Butir | Isi |
| --- | --- |
| **bj** — OQ-M8 `[DIPUTUSKAN; veto work owner]` | Gerbang rute DOL lama **disamakan** dengan gerbang tiga tanggal lainnya menurut XML *(prinsip proyek: XML menang)*. Perubahan authz dicatat bertanggal di tiket 07 dan OQ-M8 ditutup |
| **bk** — OQ-M9 `[DIPUTUSKAN; veto work owner]` | Penanda Claim Received **mengikuti XML**: bila ada rule yang menuliskannya ke tabel *(cari penulisnya: `Obj-Save`, `RDB-Save`, SQL)*, ia disimpan di kolom baru *(migrasi `021`)*; bila hanya `Property-Set` di halaman lalu ditampilkan, ia **dihitung saat baca**, tanpa kolom. Bukti dan pilihan ditulis di tiket 06 |
| OQ-M1…M7, OQ-PL-09/10/11, OQ-K-04a/05/05b | **tetap untuk work owner / pemilik ekspor** — tidak diputuskan executor |

## 2. URUTAN — satu commit per paket

| # | Paket | Isi |
| ---: | --- | --- |
| 0 | **Status tiket jujur** *(dokumen; tiga modul)* | tiap tiket: `Status:` = `selesai` / `sebagian — <apa yang kurang>` / `wontfix` sesuai kode yang ada; setiap AC dicentang `[x]` **hanya** dengan rujukan bukti *(berkas:fungsi atau nama uji)*; AC yang belum terpenuhi tetap `[ ]` dengan alasan satu baris. Commit `docs: status tiket dan AC sesuai kode — PremiumList, Komite, Claim Life` |
| 1 | **Claim Life `Save to RNM`** | baca `Activity/SaveOutStandingLife_Act.xml` **utuh sebagai pohon** *(±13.000 baris terpecah; langkah 1–akhir, tiap `RDB-List`/`Call` ke SQL-nya, tiap `WHEN` dari urutan aksi; langkah 23 total peserta sudah ditiru — pakai ulang)*; bangun `POST /api/klaim-life/{id}/outstanding` bergerbang tahap `Outstanding Claim` + pemegang; satu transaksi; tabel warisan hanya dibaca; tombol `Save to RNM` di layar Outstanding. Tiket 03 diralat bertanggal bila XML membantahnya. Commit `claim-life: Save to RNM — SaveOutStandingLife_Act` |
| 2 | **av-2 — 21 medan `PolicyDataLife`** | daftar 21 ikatan `.PolicyDataLife.*` di section Register/Outstanding/Detail *(path + baris)*; `PolisRingkas` *(pl4)* dilebarkan **aditif** dengan kolom yang ada di `T_PREMIUM_LIST`; medan tanpa sumber = penanda bernama, bukan dikarang. Commit `claim-life: av-2 — 21 medan data polis` |
| 3 | **bj + bk** | §1 |
| 4 | **Sisa temuan `/code-review`** | yang dibiarkan GILIRAN-10 dengan sadar: cadangan tahap tersalin lima kali → satu sumber; tiga tanggal satu tipe *(data clump)*; tahap dibandingkan sebagai tipe, bukan string mentah; label `Batal` dari XML atau penanda. Perilaku tidak berubah; uji tetap hijau |
| 5 | **PANDUAN UJI LAYAR** *(untuk work owner)* | `PANDUAN-UJI-LAYAR-TIGA-MODUL.md` di akar `OUTPUT_HASIL_RNM`: per modul, urutan klik dari Beranda → menu → tiap tombol XML → hasil yang diharapkan *(pesan VERBATIM, tahap tujuan, angka contoh dari fixture `UJI-*`)*; cara menyalakan backend + frontend; daftar yang **belum** dapat diuji *(OQ terbuka, pelaksana stub)*. Tanpa menjalankan tulis ke DEV |

## 3. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3, GILIRAN-6, GILIRAN-9 §2. Laporan: baris pertama alasan berhenti; tabel **paket → commit → rule XML →
rute/komponen**; ralat tiket; OQ dibuka/ditutup; angka uji tiap commit **dengan dan tanpa tag `db`**; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 28 September 2026 sesudah verifikasi `0cddee0..251cb3b` (uji dijalankan ulang dengan dan tanpa tag `db` di worktree
bersih; status dan kotak AC 36 tiket dicacah; `T_MIGRASI` DEV diperiksa).*
