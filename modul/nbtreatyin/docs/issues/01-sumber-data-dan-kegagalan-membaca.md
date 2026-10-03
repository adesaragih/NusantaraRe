# 01: Sumber data realisasi treaty — dibaca dari view relasional, gagal baca menghentikan proses

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** —
**Menutup:** AC 15 · 16 · 17 · 36 · 37 · 38 · 57 · 58 · 89 *(9 AC)* — US 21 · 23 · 24 · 37

## Hasil & nilai pengguna

Hari ini data kontrak dibaca dengan **membongkar satu dokumen teks** menjadi properti saat
halaman dibuka. `[terverifikasi]` Enam langkah Java identik melakukannya, dan bila pembongkarannya
**gagal**, galatnya **hanya ditulis ke log** — aktivitas **tetap lanjut** dengan halaman kosong
atau separuh terisi. ⛔ Pengguna tidak diberi tahu apa pun.

Sesudah tiket ini, data kontrak dibaca dari **sumber relasional yang bentuknya dapat dinyatakan**,
dan ⭐ bila pembacaan gagal, pengguna **diberi tahu** dan prosesnya **berhenti** — tidak ada lagi
berkas yang tersimpan dari pembacaan yang gagal.

## Area codebase

- Lapisan repository: pembacaan data realisasi treaty
- Lapisan repository: pembacaan penempatan keluar *(hanya baca)*
- Lapisan service: penanganan kegagalan pembacaan

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembongkaran dokumen — **tidak dimigrasi** | 6 langkah Java di `Activity\FetchMasterTreatyIn`, `SetTreatyIn_Act`, `InputPolicyTreatyInDetail_*`, `InputPolicyTreatyOutDetail_*` |
| Medan yang dipakai laporan | `ReportDefinition\BrowseTreatyInDetail.xml` — **33 medan** |
| Penempatan keluar, hanya baca | `RDBList\BrowseTreatyOut.xml` · `RDBList\BrowseTreatyOutDetail.xml` — keduanya `SELECT` |

## ADR terkait

- **ADR-0009** — migrasi penuh, tidak ada koeksistensi dua penulis

## Acceptance criteria

- [x] **AC 15** — data dibaca dari sumber relasional, bukan dari dokumen
- [x] **AC 16** — sistem baru **tidak menulis** dokumen
- [ ] 🟡 **AC 17** — ke-**33** medan yang dipakai laporan tersedia
- [ ] 🟡 **AC 89** — nol medan yang dipakai tetapi tidak tersedia
- [x] **AC 36** — kegagalan pembacaan **menghentikan** proses
- [x] **AC 37** — kegagalan pembacaan **menampilkan galat kepada pengguna**
- [x] **AC 38** — nol kasus tersimpan dari pembacaan yang gagal
- [ ] ⛔ **AC 57** — penempatan keluar **dapat dibaca** dari konteks ini
- [x] **AC 58** — penempatan keluar **tidak pernah ditulis** dari sini

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **P29** | contoh isi dokumen — hanya untuk **migrasi data lama**, bukan untuk pembacaan baru | ⭐ **tidak menahan** |
| **19** | satu aturan **bernama** menulis penempatan keluar tetapi tidak ada perintah tulisnya | tidak menahan |

## Perintah verifikasi

1. Buka sebuah realisasi treaty — ⭐ seluruh **33** medan laporan terisi.
2. Putus sumber data di tengah pembacaan — ⭐ pengguna **melihat galat**, ⭐ dan **nol** berkas
   tersimpan.
3. Coba menulis penempatan keluar dari konteks ini — ⭐ **ditolak**.

## Catatan

⚠️ `[penyimpangan sadar]` Menghentikan proses saat gagal baca **berbeda** dari perilaku Pega.
Alasannya tertulis di `spec.md` §5.9: **kegagalan yang terlihat lebih murah daripada yang
tersembunyi**.

## ⛔ RALAT implementasi 2026-10-03

1. **Nilai master `TreatyIn.*` tidak ada di view.** Bunyi lama (Hasil): *"data dibaca dari sumber
   relasional"*. Benar untuk 33 kolom RD — tetapi `TreatyIn.RNMShareP`, `RNMShare`,
   `BrokeragePercentP`, `CurrencyList`, `INSTALLMENT`, `Limits/Share` milik JSON master (P29) dan
   **tidak punya kolom padanan**; `RNM_SHARE` view belum boleh dipakai dalam perhitungan
   (PERTANYAAN-untuk-DBA). Langkah rantai uang yang membaginya **dilewati selama nilainya kosong**
   (`models.MasterTersedia`) — penyimpangan sadar, dicatat.
2. **`InputPolicyTreatyInDetail_preACT` dibangun sebagian** — langkah 3-8, 11, 14, 15 (seluruhnya
   membaca view dan tabel acuan); langkah 9-10, 13, 16-18 (JSON master) tidak.
   `pxResults(1).CURRENCYID` bukan kolom RD maupun view → ID mata uang selalu dicari menurut nama
   (langkah 4).
3. **AC 57 tidak dapat dipenuhi seperti tertulis.** Seluruh bagian treaty keluar
   (`InputPolicyTreatyOutDetail_*`, `BusinessAndSOBListRetro`, `DetailPolicyTreatyOutNonProportional`)
   membaca JSON `M_TREATY_OUT` (P29). Data treaty keluar tidak dapat ditampilkan tanpa sumber
   relasional baru — `[terbuka]`.
