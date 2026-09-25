---
status: accepted
label: DECIDED
---

# Setiap nilai uang disimpan berpasangan, dan ambang kewenangan dibandingkan terhadap nilai IDR

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Klaim non-proporsional berjalan dalam banyak mata uang sekaligus: `ListClaimAmount`, `SpreadingRisk`, dan `ListTotalEstimation` semuanya ber-kunci `Currency`, sementara properti bersufiks `IDR` menunjukkan IDR dipakai sebagai mata uang penyetaraan. Kami menetapkan setiap nilai uang disimpan sebagai satu paket — **nilai asli, mata uang, kurs, tanggal/sumber kurs, dan nilai IDR hasil konversi** — dan setiap perbandingan terhadap ambang kewenangan Komite dilakukan terhadap nilai IDR hasil konversi tersebut.

## Consequences

**Tanggal dan sumber kurs adalah bagian wajib dari paket**, bukan pelengkap. Tanpa keduanya, angka IDR tidak dapat diaudit ulang: nilai yang sama bisa dibenarkan atau disalahkan tergantung kurs kapan yang dipakai, dan rekonsiliasi dengan akuntansi menjadi tidak dapat diselesaikan. Sumber kurs yang terlihat di sistem lama adalah fungsi `getcurrencystandard` dan `TreatyInMaster.CurrencyList.Conversion`; keduanya perlu dicatat identitasnya, bukan hanya hasilnya.

Keputusan tentang ambang ini **tidak menunggu** hasil investigasi kondisi sistem lama yang tercatat di `FINDING-001-threshold-currency.md`. Apa pun temuan di sana, sistem baru membandingkan terhadap IDR.

Presisi dan skala setiap kolom dalam paket ini belum ditetapkan dan menunggu `pengetahuan/SCHEMA-ACTUAL.csv` — lihat ADR-0003 yang masih berstatus draft.
