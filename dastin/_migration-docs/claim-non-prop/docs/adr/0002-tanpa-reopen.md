---
status: accepted
label: DECIDED
---

# Klaim yang sudah ditutup tidak dapat dibuka kembali

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Di seluruh 279 berkas folder ini hanya ada satu status akhir, `Resolved-Completed`, dan tidak ditemukan satu pun rule yang mengembalikan klaim dari status akhir ke status terbuka. Kami menetapkan tidak ada mekanisme reopen: koreksi atas klaim yang sudah ditutup dilakukan melalui Adjustment baru, bukan dengan membuka ulang klaim.

## Consequences

Kesimpulan ini **disimpulkan dari ketiadaan rule, bukan dari adanya rule**. Ketiadaan jalur reopen di XML tidak membuktikan ketiadaan reopen di produksi — pembukaan ulang secara manual lewat database tidak meninggalkan jejak di rule. Verifikasinya tercatat sebagai EXTERNAL di `_selesai/OPEN-QUESTIONS.md`. Bila ternyata reopen manual memang dilakukan, keputusan ini harus ditinjau ulang sebelum model status dikunci, karena konsekuensinya menyentuh keterlacakan audit.
