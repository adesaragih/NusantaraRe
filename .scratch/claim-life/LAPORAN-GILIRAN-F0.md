# Laporan giliran F0.3 → A3 Register

**Aturan giliran** *(brief lanjutan 8 §1)*: nol teks ke manusia sesudah tiap commit paket; catatan
kemajuan ditulis di berkas ini; satu pesan ke manusia hanya sesudah A3 kelompok Register ter-commit,
atau sesudah berhenti sah *(gerbang yang hanya manusia dapat buka, disebut namanya)*.

Titik tetap: **`9186078`** *(docs: brief lanjutan 8)*.

---

## Verifikasi F0.2 yang diminta work owner

Dijalankan sebelum giliran ini dimulai; hasilnya sejalan dengan tabel §0 brief 8.

| Klaim | Diperiksa | Hasil |
| --- | --- | --- |
| 12 berkas, 5 baru | `git show --numstat 5b49954` | ✅ 12 diubah, 5 baru, +1.082/−353 |
| axios dilepas | `package.json` + impor di `src/` | ✅ nol; sisa hanya kata di komentar dan nama uji |
| build 41 modul, 158,83 kB | `npm run build` | ✅ |
| 76 uji JS; Go 261 · 0 · 34 | vitest; `go test -tags=db -v` | ✅ |

---

## F0.3 — Shell dan menu

**Isi:** Shell referensi (sidebar terlipat + laci ponsel, topbar, menu profil, palet Ctrl+K,
`PagarGalat`), `lib/lipatMenu.ts` + `components/KelompokMenu.tsx` + `hooks/useHalaman.ts` disalin,
`MENU` + `MODUL_LAIN_TERLARANG` di `labels.ts`, `App.tsx` F0.2 dibuang.

**Menu — dua butir, keduanya berbukti:**

| Butir | Bukti |
| --- | --- |
| kelompok `Claim Life` | `Flow/Register_Flow.xml:270` `<pyWorkTypeName>ClaimLife</pyWorkTypeName>` |
| `Inbox Claim Life` | `[tidak ada di korpus]` — kosakata `InboxPremiumList` / `Struktur_InboxClaimLife`; nol harness portal di ekspor |
| `Register` | VERBATIM `Flow/Register_Flow.xml:155` `<pyLabel>Register</pyLabel>` |

**Penjaga yang dibuktikan menggigit:** butir menu ketiga (`Detail & Tutup`) -> dua uji jatuh; menu
modul lain (`NB Treaty In`) -> uji jatuh.

**Ralat penjaga saya sendiri:** larangan nama modul memeriksa SELURUH teks berkas dan menuduh
`ShellProps` karena memuat "Prop". Disempitkan ke LABEL MENU - penjaga yang menuduh nama tipe akan
membuat orang mengganti nama tipenya, bukan memperbaiki menunya.

**Telemetri:** uji JS 76 -> 97, modul build 41 -> 44, Go tidak tersentuh. XML dibaca ulang: 1
(`Register_Flow.xml`, baris 155 · 270 · 268 · 313 · 343 · 358). Taksiran token +-40k (taksiran).

---

## at/au — kolom TAHAP dan TGL_CREATE  ·  SHA `cdd7b98`

Migrasi 016; `serahTerimaSah` dari peta PERAN menjadi peta TAHAP; `TahapDanPeran` membaca keduanya
dalam satu kueri; penjaga optimis `UPDATE` memakai `NVL(TAHAP, …)`; pendaftaran menulis `TAHAP`,
`TGL_CREATE`, `CREATE_OP`, dan memakai SATU jam. Ralat tiket 08 dan 14; STRUKTUR diperbarui.

XML dibaca ulang: `Register_Flow.xml` (358 · 343 · 268 · 313), `InputOSClaimLife.xml` (21349 ·
21404 · 21433 · 21839 · 21863), `InboxPremiumList.xml` (736). Go 261 → 263 PASS.

---

## F0.4 backend — rute kotak masuk  ·  SHA `6b1618c`

`GET /api/klaim-life?tahap=&halaman=&ukuran=`; `repository/inbox.go` + `services/inbox.go` +
`handlers/inbox.go`. Worklist (pribadi) vs workbasket (bersama) dari `Register_Flow.xml`
1511 · 1508 · 1300 · 1351 · 993 · 1072 · 1119 · 1198 — kedelapan baris diperiksa satu per satu.

⛔ **RALAT ANGKA DI PESAN COMMIT `6b1618c`**: pesannya menulis "Go 263 → 274 PASS"; angka
sebenarnya **272**. Dicatat di sini alih-alih menulis ulang sejarah. Cacah sesudah paket ini tetap
272 PASS · 0 FAIL · 34 SKIP.

XML dibaca ulang: `InboxPremiumList.xml` (592 · 721 · 733 · 736 · 751 · 765 · 942),
`Register_Flow.xml` (delapan baris perutean).

---

## F0.4 frontend — halaman Inbox

`pages/InboxClaimLife.tsx` + `ambilKotakMasuk` di `api.ts` + gaya. Empat tab = empat assignment,
judul VERBATIM `pyTaskName`, terjemahan Indonesia DI SAMPING. Tab yang perannya tidak dipegang
TIDAK dirender. Label kolom VERBATIM `InboxPremiumList.xml`, dan ujinya MEMBUKA korpus pada baris
yang tiap label sebut.

Uji JS 97 → 114; build 45 modul.
