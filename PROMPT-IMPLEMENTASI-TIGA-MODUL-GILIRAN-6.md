# PROMPT — GILIRAN 6 *(sesi tunggal di `OUTPUT_HASIL_RNM`, cabang `main` @ `06b1289`)*: **paket 0 bg (Beranda + 17 kelompok + PaletMenu + pola grid) → PremiumList 01 bagian 2 layar → 02 → 03 → 04 (+ av) → 05a → 05b → 06 → 08 → 09 → Komite 01 → 09 → Claim Life §3.1** — laporan pertama **paling cepat sesudah PremiumList tiket 03**

> GILIRAN-3 *(A/B/C)*, GILIRAN-4, dan GILIRAN-5 **tetap berlaku**; brief ini menetapkan keadaan awal, urutan, dan **aturan
> berhenti** yang lebih ketat. Mekanisme giliran = lanjutan 8 §1.

## ⛔ ATURAN BERHENTI GILIRAN INI

Empat giliran terakhir masing-masing berhenti sesudah **satu** paket dengan kalimat *"batas tiket"*. Itu bukan alasan.
Mulai giliran ini:

1. Laporan **pertama** boleh dikirim paling cepat sesudah **paket 0 bg + PremiumList tiket 01 bagian 2 + tiket 02 + tiket 03**
   selesai dan ter-commit *(empat commit)*.
2. Berhenti lebih awal hanya dengan kalimat **"konteks menipis: kira-kira N% tersisa"** di baris pertama laporan, disertai
   apa yang sudah ter-commit dan tiket mana yang berikutnya. Tanpa kalimat itu, laporan dini = pelanggaran brief.
3. Sesudah laporan pertama, lanjutkan sampai kelompok berikutnya *(PL 04 + av + 05a/05b; lalu 06 + 08 + 09; lalu Komite
   01–03; 04a–05; 06–09; lalu Claim Life §3.1)* dengan aturan yang sama.
4. Tidak ada pertanyaan "mana yang didahulukan" — urutannya §2.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Klaim laporan `06b1289` | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 20 aturan `ProtectAccept`, pesan VERBATIM; galat dikumpulkan lalu ditampilkan sekali | `ProtectAccept.xml` b430 `Property-Set` *(langkah 2, `ProtectLife.CARI1`)*, b4462 `Page-Set-Messages` *(langkah 12)*, b4563 prasyarat `ProtectLife.CARI1==1`; `ProtectLife.CARI1` disebut 22 ×; pesan `Please choose no offer !`, `COB can't null`, `System Reinsurance can't null` ada di `models/polis_validasi.go`; uji membaca `ProtectAccept.xml` langsung | ✅ |
| tiga penyimpangan dipertahankan | b2705 `.SUM_REASURED==0`, b2867 `.RATE==0` *(tanpa `<0`)*, b4190 `…DISCOUNT_PREMIUM_RETRO<0` — persis; kosong = nol hanya pada tiga pembanding bernama, `Money` tidak berubah | ✅ diterima |
| penjaga `TestTypeKlaimHanyaSatuRumahTersimpan` menagih `Type` polis | `satutype_test.go`; `Type` polis = `T_PREMIUM_LIST` *(051)*, bukan `T_WORK_CLAIM` | ✅ |
| Go 415 · 0 · 38; JS 254 | dijalankan ulang asisten *(lihat pesan verifikasi)* | ✅ |
| urutan: ProtectAccept didahulukan sebelum bg | diterima sekali ini *(aturan murni tanpa layar)*; **bg kini wajib pertama** | ⚠️ |
| `.env`, backend `:8080`, `UNGGAHAN_DIR` | sudah beres *(GILIRAN-5 §0)* | ✅ |

## 1. PAKET 0 — bg *(GILIRAN-5 §1, tidak diulang di sini)*

Beranda; sidebar **17 kelompok** dari folder korpus dengan butir **hanya** pada tiga modul berbukti; `PaletMenu` Ctrl+K;
logo bila logo perusahaan; pola grid + saringan *(`RecordForm`, `BilahSaringRegistry`, `CariSebaris`, `PilihCari`,
`exportXlsx`)* dipakai Inbox Claim Life. Commit `shell: Beranda + 17 kelompok modul + PaletMenu + pola grid dari REFERENSI_UI (bg)`.

## 2. URUTAN — semua di `main`, satu commit per tiket

| # | Bagian | Rujukan |
| ---: | --- | --- |
| 0 | **bg** | GILIRAN-5 §1 |
| 1 | PremiumList **01 bagian 2 — layar** `Input Offer` *(`InputOfferLife` + `ConfirmSection`)* memakai 20 aturan yang sudah ada; kotak masuk `PremiumList` dengan pola grid; tombol `Input Offer`/`Input Premium` → `CreateInputLife` | brief 3-PREMIUMLIST §1–§2 |
| 2 | PremiumList **02** tutup buku · **03** detail + nomor PL | brief 3-PREMIUMLIST §2 |
| — | **laporan pertama boleh di sini** | — |
| 3 | PremiumList **04** unggah CSV + pl4 → segera **av** Claim Life · **05a** · **05b** | brief 3-PREMIUMLIST §2; 3-CLAIM-LIFE §4 |
| 4 | PremiumList **06** · **08** · **09** | brief 3-PREMIUMLIST §2 |
| 5 | Komite **01 → 09** *(Inbox Komite memakai pola grid; `ShowTransfer` = layar keputusan)* | brief 3-KOMITE §1–§2 |
| 6 | Claim Life **§3.1** — 41 baris sensus diputuskan | brief 3-CLAIM-LIFE §3.1 |

## 3. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3. Laporan: tabel **paket → tombol/kelas XML → rute/komponen**, ralat tiket bertanggal, angka uji
tiap commit, bab **TELEMETRI EKSEKUSI** per paket, dan baris pertama menyebut apakah berhenti karena selesai atau karena
**konteks menipis**.

---

*Disusun 28 September 2026 sesudah verifikasi `06b1289` (uji dijalankan ulang; `ProtectAccept.xml` b430/b2705/b2867/b4190/
b4462/b4563 dibaca ulang; tiga pesan VERBATIM dicocokkan ke `models`).*
