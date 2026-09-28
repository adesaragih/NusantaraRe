# PROMPT — GILIRAN 10 *(**sesi BARU** di `OUTPUT_HASIL_RNM`, cabang `main` @ `021d0f3`)*: **PL 05a bagian 2 (pembaca, penyimpan, rute; langkah 15 = pl7) → 05b → 06 + 08 + 09 → Komite 01 → 09 → Claim Life §3.1** — laporan pertama paling cepat sesudah **05b**

> ⛔ **Buka sesi Claude Code BARU untuk giliran ini** *(induk §0.5)*: sesi lama tinggal 5% konteks. Semua yang dibutuhkan ada di
> disk — brief ini, GILIRAN-3 *(A/B/C)*, 4–9, tiket, PARITAS, LAPORAN. Aturan berhenti GILIRAN-6 berlaku.
> Mekanisme giliran = lanjutan 8 §1.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 1 commit `021d0f3` — rumus rekap uang murni; berhenti "konteks menipis 5%" | `models/polis_summary.go` + uji; alasan sah | ✅ |
| Go 514 · 0 · 38; JS 339 | dijalankan ulang: **514 PASS · 0 FAIL · 38 SKIP** · vitest **339** · vet, gofmt, tsc bersih | ✅ |
| `.PREMIUM` satu cabang per `Type`; uji merah atas versi Σ empat | `TestPremiumSatuCabangBukanEmpat` | ✅ |
| tiga keanehan warisan dikunci; pembulatan hanya di akhir | `polis_summary.go:42–47, 159`; `TestBrokerageRetroDitambahBukanDikurangi` | ✅ |
| langkah 15 `SubmitPremiumList_Act`: `@replaceAll` tidak pernah cocok | b3413 `@replaceAll(.PremiumListSummary.PL_NUMBER, Local.CurrentMMYY, Local.NextMMYY)`; b1142 `Local.CurrentMMYY = "."+MM+"."+**YYYY**+"."` *(empat angka)* sedangkan `PL_NUMBER` memuat `.MM.YY.` *(dua angka, b3021/b3075, b3126)* → **tidak pernah cocok** | ✅ → **pl7** |
| pohon `main` | ⚠️ `frontend/src/App.tsx` memuat **suntingan belum di-commit** *(15 baris, seluruhnya blok komentar `{/* … */}`)* — bukan milik executor; **jangan** di-commit atau dibuang; bila mengganggu, laporkan | ⚠️ |

**pl7** `[DIPUTUSKAN; veto work owner]`: langkah 15 adalah **no-op** di produksi *(pola `.MM.YYYY.` tidak pernah ada di dalam
`PL_NUMBER` berformat `.MM.YY.`)*; **tidak ditiru**; dicatat bertanggal di tiket 03 dan 05a dengan bukti b1142/b3413/b3126
sebagai residu warisan; bila kelak pemilik Pega menyatakan maksudnya *(menggeser periode nomor saat tutup buku)*, itu keputusan
baru, bukan replikasi.

## 1. URUTAN — semua di `main`, satu commit per tiket

| # | Paket | Isi | Rujukan |
| ---: | --- | --- | --- |
| 1 | **PL 05a bagian 2** | pembaca peserta per polis *(SUM/GROUP BY `CURRENCY` di SQL, atau baca lalu hitung dengan rumus murni yang ada — pilih satu, buktikan setara dengan contoh literal 975/978/897/1789)*; penyimpan `T_PREMIUM_LIST_SUMMARY` hapus-lalu-sisip **satu transaksi** bersama penomoran *(`SubmitPremiumList_Act` langkah 1–14, 16+; langkah 15 = pl7)*; **pl2** `M_LIFE_PREMIUM_SUMMARY`/`M_LIFE_PREMIUM_DETAIL` dalam transaksi yang sama, nol `COMMIT` di teks SQL; rute simpan + layar `ShowLifePremiumSummary` *(`Submit` b27471 → `InsertJsonPolisLife_Act` b26414 + `finishAssignment` b26442 = 05b)* | tiket 05a; brief 3-PREMIUMLIST §1–§2 |
| 2 | **PL 05b** JSON polis | `InsertJsonPolisLife_Act` utuh; keputusan tiket/spec; ralat berbukti bila bertentangan | tiket 05b |
| — | **laporan pertama boleh di sini** | | |
| 3 | **PL 06 + 08 + 09** | outbox Arasapas/email + stub; kontrak hilir dua sisi; verifikasi 050–056 lawan STRUKTUR | brief 3-PREMIUMLIST §2 |
| 4 | **Komite 01 → 09** | Inbox Komite pola grid; `ShowTransfer` layar keputusan; `KomitePostAdjustment` utuh; km1–km5 | brief 3-KOMITE |
| 5 | **Claim Life §3.1** | 41 baris sensus diputuskan; celah dibangun | brief 3-CLAIM-LIFE §3.1 |

## 2. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3, GILIRAN-6 aturan berhenti, GILIRAN-9 §2 *(rumus uang dari brief diverifikasi ulang; `WHEN`
dibaca dari urutan aksi)*. Laporan: baris pertama alasan berhenti; tabel **tiket → commit → rule XML → rute/komponen**; ralat tiket;
OQ; angka uji tiap commit; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 28 September 2026 sesudah verifikasi `021d0f3` (uji dijalankan ulang; `SubmitPremiumList_Act.xml` b1142/b3413
dibaca ulang; `polis_summary.go` dan ujinya dicek; pohon kerja diperiksa).*
