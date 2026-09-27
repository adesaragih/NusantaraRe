# PROMPT — lanjutan 11 *(sesi Claim Life, cabang `main`)*: **Detail & Tutup (2) → Dokumen → Medis → Akseptasi → Komite-handoff → A4** dalam satu giliran; dua sesi modul lain berjalan di worktree-nya

> Baca `PROMPT-INDUK-TIGA-MODUL.md` *(pembagian dan kontrak)*, lalu brief lanjutan 10 §2–§7 *(isi tiap
> kelompok tetap berlaku dan tidak diulang di sini)*. Mekanisme giliran = lanjutan 8 §1.
> ⛔ *"Belum dibaca sebagai pohon, jadi saya berhenti"* **bukan** berhenti sah: membaca pohon adalah
> **bagian pekerjaan giliran ini**, bukan syarat sebelum giliran dimulai. Berhenti sah hanya gerbang
> yang manusia sendiri harus buka *(lanjutan 2 §1–2)*. Konteks menipis → catat ke
> `LAPORAN-GILIRAN-F0.md`, lanjutkan *(lanjutan 3 §1 butir 4)*.

---

## 0. KEADAAN SESUDAH `85065b9` — DIVERIFIKASI ULANG

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 8 commit dari `29a7ebf` | 9 termasuk dua commit asisten *(ax kode + docs)*; 7 milik executor | ✅ |
| Go 285 · 0 · 34; JS 169; 46 modul; vet bersih; nol migrasi baru | dijalankan ulang, angka sama; migrasi tetap `001`–`016` | ✅ |
| `DeletePesertaClaimLife` **tidak menghapus** | `Activity` = `Property-Set` ×2 + `Obj-Save`, nol `Page-Remove`; `deleteRow` hanya di `InputOSClaimLife.xml` 17874/18017; syarat `CLAIM_NO == ''` 18082 | ✅ ralat brief 10 §1 diterima |
| ralat tiket 03, PARITAS 81, OQ-E/F/G | tiket 03 baris 535; PARITAS `⛔ RALAT`; OQ 148/162/173 | ✅ |
| `AUTH_STUB` tidak terdokumentasi → `.env.example`, PANDUAN, uji | commit `3ce9fc7`, 3 berkas; `.env` work owner kini `AUTH_STUB=true` | ✅ |
| nol kebocoran | seluruh berkas yang berubah | ✅ |

Yang belum *(laporan executor, benar)*: §2 brief 10 tujuh dari delapan butir; §3–§7 nol.

### 0.1 Dua sebab pita merah *"Permintaan ditolak backend"* di layar work owner *(diperiksa 27-09-2026 16.10)*

| # | Sebab | Bukti | Siapa membereskan |
| ---: | --- | --- | --- |
| 1 | backend berjalan dengan **`stubPelaku=false`**: proses `api` dimulai 16.06 di jendela yang memuat `.env` **sebelum** `AUTH_STUB=true` ditambahkan *(16.00)* — `muat-env.ps1` memang mengekspor semua kunci, tetapi hanya saat dijalankan | `GET /api/klaim-life?tahap=1` **dengan** header stub → `401 {"galat":"permintaan tanpa identitas pelaku ditolak"}`; `healthz` → database terjangkau | **work owner**: di jendela backend `Ctrl+C` → `. .\muat-env.ps1` *(harus mencetak `AUTH_STUB=true`)* → `go run .\cmd\api` |
| 2 | **kunci envelope galat tidak sama**: backend menulis `{"galat": …}` *(`handlers/klaimlife.go:47`, `dol.go:62`, sejak tiket 01)*, frontend membaca `o.error` *(`services/api.ts:181`)* sehingga pesan apa pun jatuh ke teks bawaan *(`api.ts:633`)*. Komentar `api.ts:592` yang menyebut `{"error": …}` **keliru** | tidak ada tiket/ADR yang menetapkan kuncinya; yang lebih tua dan diuji adalah `galat` | **executor, paket 0 giliran ini**: `api.ts` membaca `galat`; komentar diralat; **uji kontrak dua sisi** *(handler menulis `galat`; klien membaca `galat`, dan `error` **tidak** diterima diam-diam)*; commit `fix: envelope galat — klien membaca "galat" seperti yang backend tulis; uji kontrak` |

## 1. DUA KEPUTUSAN ATAS "BELUM PUNYA PEMANGGIL" `[DIPUTUSKAN — pemanggil XML-nya ada; veto work owner]`

| Butir | Keputusan | Bukti pemanggil |
| --- | --- | --- |
| **az** `lib/pilihSemua.ts` | **dipertahankan**; dipasang pada tombol `Select All` layar **Outstanding** di kelompok berikutnya *(bukan Register)* | `InputOSClaimLife.xml` 24489 `pyButtonLabel Select All`; `SelectAllClaimLife_act` menyalakan `.IsAccept` seluruh baris |
| **ba** `models.PenandaBatasHari` | **dipertahankan**; ambangnya dari **view** `PRODUCTINWARD_LIFE` `[data DBA: view 40 kolom, kolom `MAXEXPIREDCLAIM`, `MAXDATARECEIVE` — dibaca `GetProductName`]* lewat repository baru berkunci `ID` = `PolicyDataLife.ProductNameID`; **kuncinya** menunggu av *(PremiumList Life)*, jalannya dibangun sekarang dengan uji, dan layar Detail menampilkan "menunggu modul PremiumList Life" sampai kuncinya ada | `RDBList/GetProductName.xml`: `SELECT * FROM POOLDATA.PRODUCTINWARD_LIFE WHERE ID = {pyWorkPage.PolicyDataLife.ProductNameID}` |

## 2. URUTAN GILIRAN INI — paket 0 lalu lima kelompok, tanpa pesan di antaranya

| # | Kelompok | Isi *(rincian di brief 10)* | Commit |
| ---: | --- | --- | --- |
| 1 | **Detail & Tutup (2)** | `ClaimLifeDetailGCNM` sebagai pohon *(field, `<pyCondition>`, tombol → activity; `Save Adjustment` 22590 → rute akseptasi; `Diagnose_Harness` 5080/5224 → `BelumTersedia` bernama)*; `ShowEditClaimLife` + `EditDateClaimLife_Section` → rute tanggal kejadian; `SetSTS_Reject`; `SetIndexAdjustmentList`; `CloseClaim` + `CloseClaim_Section` 1028 + `ProtectCloseClaim_act`; `DocumentLife` tampil daftar; tombol `Send Back to Register`/`Send to Medical Check` sudah ada *(aw)* | `claim-life: A3 — Detail & Tutup (2)` |
| 2 | **Dokumen** | brief 10 §3 *(ar1, an, `GetMimeType` → tabel data, outbox `T_LOG_SERVICE_RNM`, pengirim stub di DEV, unggahan berbatas ke folder aplikasi)* | `claim-life: A3 — Dokumen` |
| 3 | **Medis** | brief 10 §4 *(`DISEASE_LIFE` 97.586 baris → pencarian berbatas; `.DiagnoseList`; keputusan **al** dari bukti; `SendtoAdmin` dari Medical berjalan sesuai XML)* | `claim-life: A3 — Medis` |
| 4 | **Akseptasi** | brief 10 §5 | `claim-life: A3 — Akseptasi` |
| 5 | **Komite-handoff** | brief 10 §6 — **tanpa** menyentuh berkas `komite_*` sesi Komite | `claim-life: A3 — Komite` |
| 6 | **A4** *(bila giliran masih hidup)* | brief 10 §7 | `claim-life: A4 — migrasi data` |

Tiap kelompok: bab **Pembacaan ulang XML** dengan path + baris di tiket yang terkait, ralat tiket bila
XML membantah, `PARITAS-LAYAR-DAN-AKSI.md` diperbarui, uji murni + handler + JS, `LAPORAN-GILIRAN-F0.md`
bertambah satu bab. Nol migrasi baru kecuali dari keputusan yang tercatat; nomor `017`–`029`.

## 3. YANG SUDAH DIBERESKAN DI LUAR SESI INI *(jangan diulang)*

- Worktree `modul/premiumlist-life` dan `modul/komite-claim-life` **sudah dibuat asisten** 27-09-2026
  dari `85065b9` beserta `.env` per port dan `npm install` *(induk §1)*. Sesi ini **tidak** membuat
  atau menyentuhnya.
- `AUTH_STUB=true` sudah ada di `.env` work owner; backend berjalan dengan database terjangkau.
- Keputusan **ax** *(`T_LOG_SERVICE_RNM`)* sudah di-commit `e3d537a`.

## 4. LANGKAH 0 · LAPORAN · TELEMETRI

Langkah 0: brief ini dan induk **sudah di-commit asisten** sebagai `38cc700`; kedua worktree modul
sudah di-fast-forward ke SHA yang sama. Yang tersisa untuk executor: `git status --porcelain` kosong →
uji hijau *(285 · 34 SKIP · 169 JS · 46 modul)* → langsung paket 0. Laporan akhir dan telemetri
persis lanjutan 8 §6–§7; satu pesan, sesudah kelompok 5 *(atau A4)*.

---

*Disusun 27 September 2026 sesudah verifikasi `85065b9` (uji, vet, build dijalankan ulang; 9 commit
dibaca statistiknya), pembacaan ulang `DeletePesertaClaimLife.xml`, `InputOSClaimLife.xml`
(`deleteRow`, `CLAIM_NO`), `GetProductName.xml`, tiket 03, PARITAS, `OQ-untuk-tim.md`, `config.go`
(`AUTH_STUB`), dan pembuatan dua worktree modul.*
