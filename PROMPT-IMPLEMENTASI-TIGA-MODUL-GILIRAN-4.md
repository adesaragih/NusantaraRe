# PROMPT — GILIRAN 4 *(sesi tunggal di `OUTPUT_HASIL_RNM`, cabang `main` @ `5482a00`)*: **PremiumList tiket 01 bagian 2 → 09 (+ av) → Komite tiket 01 → 09 → Claim Life §3.1 triase 41 baris**, tanpa berhenti di antara tiket

> Ketiga brief GILIRAN-3 *(`…-3-PREMIUMLIST.md`, `…-3-KOMITE.md`, `…-3-CLAIM-LIFE.md`)* **tetap berlaku utuh**; brief ini hanya
> menetapkan keadaan awal dan urutannya. Mekanisme giliran = lanjutan 8 §1. Berhenti sah hanya pada **batas tiket** dengan pohon
> bersih, LAPORAN terbarui, SHA disebut — dan hanya bila konteks benar-benar menipis, bukan sesudah satu paket.

---

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Hal | Keadaan |
| --- | --- |
| `main` | `5482a00` = `632b53b` + merge Komite `96ba363` *(ralat ShowTransfer)* + merge PremiumList *(pl6 ralat 056, penjaga `TestNolKataCadanganOracleSebagaiKolom`, tiket 01 bagian 1)*; kedua worktree modul **sudah** di `5482a00`; tiga pohon bersih *(suntingan 016 sudah dikembalikan work owner)* |
| Uji pohon terpadu | **Go 402 PASS · 0 FAIL · 38 SKIP** · vitest **254** · **51** modul · vet, gofmt, tsc bersih |
| DEV `T_MIGRASI` | **28 langkah lengkap**: `001`–`020`, `030`, `050`–`056` *(056 terpasang 28-09 10.10; `T_VIEW_SUGGEST` berkolom `INITIAL_SUGGEST`)*. Tidak ada langkah yang menunggu |
| `.env` work owner | `UNGGAHAN_DIR` **belum** ada → rute unggah dokumen Claim Life akan menolak sampai diisi *(di luar sesi; executor tidak menyunting `.env`)* |
| Log GILIRAN-3 yang sudah diverifikasi | sensus 135 rule `eb9bcef` ✅; ralat ShowTransfer `6409915` ✅ *(rantai b763 → b91 → section; `KomiteRouter` b442 `[terbuka — pemilik ekspor]` diterima)*; paket 0 PremiumList `2f743b4` ✅ |

## 1. URUTAN GILIRAN INI — semua di `main`, satu commit per tiket, nol pesan di antara tiket

| # | Bagian | Rujukan | Catatan |
| ---: | --- | --- | --- |
| 1 | **PremiumList tiket 01 bagian 2** — layar `Input Offer` *(`InputOfferLife` + `ConfirmSection`)* dan kotak masuk `PremiumList` *(dua tombol b16472/b16964 → `CreateInputLife`)* | brief 3-PREMIUMLIST §1–§2 | menu = kelompok **PremiumList Life** → satu butir `PremiumList` di Shell yang ada; label VERBATIM |
| 2 | **PremiumList tiket 02 → 09** *(02, 03, 04, 05a, 05b, 06, 08, 09)* | brief 3-PREMIUMLIST §2 | migrasi baru `057`+ hanya dari keputusan tercatat; penjaga kata cadangan sudah menjaga nama kolom |
| 3 | **Claim Life av** — segera sesudah tiket 04 *(pl4 `repository.PolisRingkas`)* ada | brief 3-CLAIM-LIFE §4 | `PanelDataPolis` sepuluh medan; ambang **ba**; penanda "menunggu modul PremiumList Life" dicabut |
| 4 | **Komite tiket 01 → 09** *(01, 02, 03, 04a, 04b, 05, 06, 07, 08, 09)* | brief 3-KOMITE §1–§2 | `ShowTransfer` = layar keputusan; `TransferType` tetap `[terbuka]`; dokumen komite memakai jalur dokumen Claim Life yang ada |
| 5 | **Claim Life §3.1** — 41 baris sensus diputuskan: **tidak perlu** + sebab XML, atau **celah → dibangun** | brief 3-CLAIM-LIFE §3.1 | nol baris tanpa keputusan |

Bila konteks menipis sebelum nomor 5: berhenti pada batas tiket, laporkan, dan giliran berikut mulai dari nomor berikutnya —
bukan mengulang yang sudah ada.

## 2. ATURAN YANG TIDAK BERUBAH

XML sebagai pohon *(path + baris)* di tiket sebelum kode; ralat bertanggal bila XML membantah tiket; menu hanya yang berbukti
korpus; uji murni + handler + JS per tiket; `PARITAS-LAYAR-DAN-AKSI.md` milik tiap modul; `LAPORAN-GILIRAN-F0.md` +1 bab per
tiket; `-migrate` **tidak** dijalankan executor; endpoint luar **tidak** dipanggil di DEV *(pelaksana stub)*; nol kebocoran
*(nama orang, nomor polis, kredensial, alamat layanan)*; pembagian berkas induk §2.

## 3. LAPORAN · TELEMETRI

Satu pesan di akhir *(atau batas tiket)*: tabel per modul **tiket → commit → menu/tombol XML → rute/kontrol**; ralat tiket;
OQ dibuka/ditutup; angka uji tiap commit; bab **TELEMETRI EKSEKUSI** per tiket *(lanjutan 8 §6–§7)*.

---

*Disusun 28 September 2026 sesudah verifikasi `5482a00` (uji pohon terpadu di worktree bersih), katalog DEV (`T_MIGRASI` 28
langkah, kolom `T_VIEW_SUGGEST`), dan `.env` work owner (kunci saja, tanpa nilai).*
