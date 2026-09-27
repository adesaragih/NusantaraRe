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

