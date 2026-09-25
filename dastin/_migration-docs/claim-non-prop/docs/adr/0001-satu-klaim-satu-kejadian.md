---
status: accepted
label: DECIDED
---

# Aggregate root adalah Klaim, satu klaim satu kejadian kerugian

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Struktur XML mendukung tafsir ini tetapi tidak memaksakannya: seluruh properti tanggal dan penyebab kerugian di `.ClaimData` bersifat skalar, ada 18 PageList di bawahnya, dan sapuan atas 114 activity menunjukkan **tidak ada satu pun rule yang menulis `.DateOfLoss`** — nilainya hanya masuk lewat input Registrasi. Kami menetapkan satu Klaim = satu kejadian kerugian atas satu polis treaty non-proporsional, dengan Adjustment sebagai transaksi di bawahnya, sehingga Klaim menjadi aggregate root dan ke-18 PageList menjadi tabel anak.

## Consequences

Natural key Klaim belum ditetapkan. Tidak ada rule yang mencegah duplikat `PolicyNo` + `DateOfLoss` + `IDMaster`; `CheckDateDOL_Act` hanya menampilkan riwayat klaim atas polis yang sama sebagai informasi, tidak menolak. Sampai ada konfirmasi bahwa pencegahan duplikat dilakukan secara prosedural, kunci teknis tetap identitas klaim yang digenerasi sistem, dan kombinasi di atas diperlakukan sebagai indeks pendukung, bukan unique constraint.
