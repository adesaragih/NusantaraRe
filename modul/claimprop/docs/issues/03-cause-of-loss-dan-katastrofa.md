# 03: Cause of Loss dua tingkat dan penandaan katastrofa

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis)
**Menutup:** AC 96 · 97 · 98 · 118 · 119 · 120 · 121 · 122 *(8 AC)* — US 7–9

## Hasil & nilai pengguna

Claim Admin mengisi penyebab kerugian dari master dua tingkat yang bertaut, sehingga penyebab
tercatat seragam dan dapat dilaporkan. Klaim tidak dapat disimpan tanpa penyebab. Klaim yang berasal
dari satu peristiwa besar dapat ditandai katastrofa dan dikelompokkan di bawah satu event, sehingga
akumulasi kerugian per peristiwa terlihat.

## Area codebase

Master Cause of Loss (induk, anak, relasi lini bisnis) · entitas event katastrofa · validasi simpan.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `ReportDefinition/BrowseVMCauseOfLoss_RD.xml` | — | tingkat **induk** — id induk, deskripsi, id lama |
| `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` | — | tingkat **anak**, disaring oleh id induk yang sedang dipilih; ⚠️ memuat satu rujukan kelas salah ketik |
| `Activity/CNMInsertCauseOfLoss_act.xml` | — | perawatan master induk; ⚠️ menerima muatan JSON lewat parameter bernama seolah kolom teks |
| `Activity/CNMInsertDetailCauseOfLoss_act.xml` | — | perawatan master anak; ⚠️ pola parameter yang sama |
| `RDBList/GetLBUID_SQL.xml` | — | **tingkat ketiga** — relasi penyebab kerugian ke lini bisnis; ⚠️ kedua tabel **tanpa prefiks schema**, digabung dengan koma gaya lama |
| `Activity/ProteksiData_act.xml` | 14 | gerbang wajib — simpan ditolak saat penyebab kerugian kosong; flag `true`, arah normal |
| `Activity/SaveCatasrtope_Act.xml` | 3 | master event katastrofa disimpan lewat `Obj-Save`; ⚠️ step SQL-nya (step 2) **di-remark** → tidak ditulis |
| `Activity/SetDefNonCatastrope_Act.xml` | 1 · 2 | penanda katastrofa ya/tidak dan jenis non-katastrofa |
| `Activity/InputCatastrope.xml` · `Activity/SetCatastrope_act.xml` · `Activity/SetEditCatastrope.xml` | — | pembuatan dan penyuntingan event |
| `ReportDefinition/GetCatastrope_RD.xml` | — | pemilihan event yang sudah ada |

## ADR terkait

Tidak ada ADR yang mengikat langsung; keputusan skema mengikuti **tiket 00**.

## Acceptance criteria

- [ ] `[terverifikasi]` Cause of Loss dipilih dari **master dua tingkat yang bertaut** — tingkat anak disaring oleh id induk yang sedang dipilih *(AC 118 spec)*
- [ ] `[terverifikasi]` Cause of Loss dua tingkat dihubungkan kunci induk, ditambah relasi ke lini bisnis; **wajib** sebelum simpan *(AC 97 spec)*
- [ ] `[terverifikasi]` Tingkat ketiga — relasi Cause of Loss ke lini bisnis — tersedia. ⚠️ Kedua tabelnya diprefiks schema eksplisit dengan `JOIN` yang tertulis; di Pega keduanya telanjang dan digabung koma gaya lama *(AC 119 spec)*
- [ ] ⚠️ Parameter procedure pembaruan Cause of Loss **dinamai sesuai isinya**. **Alasan menyimpang:** kedua procedure Pega menerima muatan JSON lewat parameter bernama seolah kolom teks, sehingga tanda tangannya menyesatkan pembaca *(AC 98 spec)*
- [ ] ⚠️ `[terverifikasi]` Simpan **ditolak** bila Cause of Loss belum diisi. **Jebakan:** rule yang sama menyiapkan **dua pesan kembar** berbunyi sama dan memakai **penanda kesalahan kedua** berambang `>1` — pola yang sama dengan penanda yang dipecah di tiket 00 (AC 105); jangan menyatukan keduanya tanpa membaca kodenya *(AC 121 spec)*
- [ ] `[terverifikasi]` Katastrofa punya **entitas event sendiri** dan satu event **mengelompokkan banyak klaim** *(AC 96 spec)*
- [ ] ⚠️ Penanda katastrofa **ya/tidak** dicatat terpisah dari event-nya, dengan **satu ejaan dan satu enum**. **Alasan menyimpang:** di Pega nama propertinya bahasa Indonesia sedangkan nilainya bahasa Inggris, dan modul mengeja konsep ini dalam **empat bentuk** — `Catastrope` (82 kemunculan) · `Catastrophe` (16) · `Catasrtope` (8) · `Catastrofe` (3 berkas) *(AC 122 spec)*
- [ ] ⚠️ Rujukan kelas salah ketik di `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` **tidak dibawa**; penamaan mengikuti satu bentuk *(AC 120 spec)*

## Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` memuat rujukan kelas yang kurang satu
huruf, berdampingan dengan ejaan benar di berkas yang sama. Hampir pasti salah ketik lama; perbaikannya
sepele. Pemilik: **pemilik export Pega**.

⚠️ Aturan induk berlaku: `Activity/SaveCatasrtope_Act.xml` step **2** di-remark → **step itu tidak
ditulis sama sekali**; penyimpanan mengikuti jalur step 3.

## Perintah verifikasi

```
jalankan test "pilih penyebab induk -> daftar anak tersaring oleh induk itu"
jalankan test "simpan tanpa Cause of Loss -> DITOLAK"
jalankan test "dua klaim menunjuk satu event katastrofa -> keduanya terkelompok"
jalankan test "penanda katastrofa tersimpan sebagai enum tertutup"
cari ejaan katastrofa di kode                             -> tepat satu bentuk
cari rujukan kelas salah ketik                            -> nihil
```
