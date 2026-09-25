---
status: accepted
label: DECIDED
---

# Layer dan Retensi Cedant adalah dua entitas terpisah

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Susunan lapisan XOL dan Retensi Cedant disimpan di **dua tabel berbeda**. Sistem lama menyimpan keduanya di satu PageList `SpreadingRisk`, dengan Retensi Cedant sebagai baris ber-`TreatyName="UR"` yang di-`APPEND` oleh `CountLossAllocation_act`.

## Consequences

Alasan yang menentukan bukan soal penyaring, melainkan **perilaku baris itu sendiri**: `ClaimPercentage` selalu `0` dan `ClaimSpreaded` selalu `0`. Porsi reasuradur untuk baris itu tidak pernah ada. Sesuatu yang selalu nol di dalam daftar alokasi bukan anggota daftar itu — ia konteks bagi daftar tersebut.

Argumen bahwa "kalau ia lapisan nol, tidak akan ada rule yang perlu menyaringnya" **tidak dipakai** sebagai dasar, karena penyaring bisa ada karena sebab lain.

Setelah dipisah, pertanyaan "Retensi Cedant ikut dihitung atau tidak" menjadi eksplisit pada setiap query — bukan bergantung pada apakah penulisnya ingat menuliskan penyaring. Di sistem lama, sembilan rule menyentuh `SpreadingRisk` antara 16 dan 66 kali tanpa pernah menyebut `"UR"`.

Ini mengunci model data inti dan sangat mahal diubah kemudian: setiap query alokasi, setiap agregasi, dan setiap laporan ikut berubah bentuknya.

## Tambahan 18 September 2026 — satu kolom yang selalu nol ikut masuk ke tabel ini

Baris Retensi Cedant menerima `TotalClaim = Local.TotalUR` (`Activity\CountLossAllocation_act.xml` baris 7015–7120).

Dan `Local.TotalUR` pada cabang `Local.Currency == .Currency` dihitung sebagai `(.Deductible * 0 * Local.ProrateClaim/100)` — **dikalikan nol**, sehingga selalu nol. Penetapan `Local.UR` pada cabang yang sama, di berkas yang sama, tidak memakai `* 0`.

Artinya: untuk seluruh klaim bermata uang sama, tabel Retensi Cedant yang baru akan memuat kolom `TotalClaim` yang **selalu bernilai nol** — bila perilaku lama dipertahankan.

Apakah perilaku itu dipertahankan atau diperbaiki **belum diputuskan**, dan bukan keputusan teknis: ia dibawa ke akuntansi sebagai `ASK-AKUNTANSI.md` pertanyaan **2b**. Keputusan ADR ini tentang **bentuk tabel** tetap berlaku apa pun jawabannya.
