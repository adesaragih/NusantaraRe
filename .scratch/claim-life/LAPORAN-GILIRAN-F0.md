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

---

## F0.5 — panduan  ·  SHA `16d567f`

`PANDUAN-MENJALANKAN.txt` bab 4 diperbarui: layar masuk stub, Shell + menu dua butir beserta
buktinya, Inbox empat tab, pembedaan antrian PRIBADI vs BERSAMA. Cacah uji JS 5 -> 114.

Satu kalimat yang sengaja ditulis terang: tab Admin akan KOSONG sampai pemakai mendaftarkan klaim
sendiri (worklist = `CREATE_OP` miliknya). Tanpa itu orang melaporkannya sebagai bug.

---

## A3 kelompok 1 — REGISTER

**Backend:** `repository/rujukan.go` (tiga pembaca berbatas), `services/rujukan.go` (gerbang,
himpunan tertutup), `handlers/rujukan.go` + rute `GET /api/rujukan/{jenis}?cari=`.

**Frontend:** `components/PanelDataPolis.tsx` (himpunan medan PERSIS section), label `REGISTER`
dengan 15 nomor baris, `cariRujukan` di `api.ts`, `RegisterKlaim.tsx` memakai label VERBATIM.

**Dua temuan yang mengubah kode:**

1. Penjaga arsitektur `TestHandlersTidakMengimporRepository` menangkap `handlers/rujukan.go`
   mengimpor repository. Diperbaiki dengan ALIAS di services, bukan struct kembar.
2. Urutan pemeriksaan salah: jenis rujukan yang salah ketik dijawab "database belum dikonfigurasi".
   Permintaan diperiksa lebih dulu, baru infrastrukturnya.

**XML dibaca ulang (6 berkas):** `InputRegisterClaimLife.xml` (15 baris label + `pyValue`/`pyFormat`
tiap medan), `BrowseCedingCoLife_RD.xml`, `BrowseBusinessLife_RD.xml`, `BrowseMarketingOfficer_RD.xml`,
`InputOSClaimLife.xml`, `Register_Flow.xml`.

**Telemetri:** Go 272 -> 274 PASS · 0 FAIL · 34 SKIP; JS 114 -> 136; build 46 modul.
Taksiran token giliran ini **+-1,1 juta** (taksiran; angka sejati tidak terlihat dari dalam sesi).

---

## F0.6 — form login ditiadakan  ·  SHA `daa9b2c`

Halaman `Masuk` + ujinya dibuang; identitas dari env Vite (`VITE_AUTH_STUB`, `VITE_STUB_PELAKU`,
`VITE_STUB_PERAN`). Bawaan peran KETIGA-tiganya. Tombol `Keluar` dibuang. JS 136 → 129.

## Ralat av  ·  SHA `1255809`

⛔ **RALAT PEMBACAAN SAYA SENDIRI.** Empat medan (`Type`, `Marketing Officer`, `Ceding`,
`Class of Business`) `pyReadOnly` true / `Read-only` di Pega, terikat `.PolicyDataLife.*`. Ronde
sebelumnya saya membangunnya sebagai dropdown + rute + tiga pembaca tabel — **kode mati di tiga
lapis, nol `.tsx` memakainya**. Sebabnya: blok kontrol berdiri **sebelum** labelnya di DOM, dan
jendela pembacaan saya menghadap **ke depan**. Itu grep dengan langkah tambahan, bukan pohon.

Dibuang: `repository/rujukan.go`, `services/rujukan.go` (+uji), `handlers/rujukan.go`, rutenya,
`api.ts cariRujukan`. Tiket 02 + PARITAS §6 + `OQ-untuk-tim.md` (4 pertanyaan).

⛔ **RALAT ANGKA DI PESAN COMMIT `1255809`**: pesannya menulis "Go tetap 274"; angka sebenarnya
**272** — membuang `rujukan_test.go` menghapus 2 uji. Dicatat di sini, bukan dengan menulis ulang
sejarah.

## A3 Outstanding — butir aw

Rute `POST /api/klaim-life/{id}/tahap/{tujuan}` (tujuan KATA), `handlers/tahap.go`,
`pages/OutstandingClaimLife.tsx`, `OUTSTANDING` + `TOMBOL_OS` di `labels.ts` dengan 16 nomor baris.

**Pertanyaan pohon terjawab:** `SendtoAdmin_Act` b338 berprasyarat
`pyPosition=="ReasLifeMedicalAdvisor"` (WhenFalse=3 = lewati) padahal layar Outstanding dipegang
Admin → kedua tombol **tidak menulis apa pun** di Pega. Cacat rule warisan; ditiru **maksudnya**,
cacatnya dilaporkan OQ-C. Ralat tiket 08.

**Telemetri:** Go 272 → 274 PASS · 0 FAIL · 34 SKIP; JS 129 → 150; build 46 modul. XML dibaca ulang
(4 berkas): `InputRegisterClaimLife.xml` (blok kontrol read-only), `InputOSClaimLife.xml` (16 baris),
`Register_Flow.xml` (8 penyambung), `setDetailClaim_act.xml`. Taksiran token **±1,4 juta** giliran
ini (taksiran).
