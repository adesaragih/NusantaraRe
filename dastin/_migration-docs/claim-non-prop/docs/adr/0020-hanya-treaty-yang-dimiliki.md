---
status: accepted
label: DECIDED
---

# Fakultatif dan Treaty dua entitas berbeda; sistem ini hanya memiliki Treaty

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Polis fakultatif dan polis treaty adalah **dua entitas berbeda**. Sistem ini **hanya memiliki Treaty**; Fakultatif dibaca sebagai rujukan, tidak dimodelkan, tidak dimiliki. Bila kelak ada modul yang memerlukannya, modul itu yang memilikinya.

## Consequences

**Alasan pemisahannya adalah domain, bukan bentuk penyimpanan.** Fakultatif dinegosiasikan per risiko, satu per satu; treaty adalah kontrak payung atas satu portofolio. Beda cara lahirnya, beda dokumennya, beda kewajiban para pihaknya.

Argumen bentuk data — bahwa `V_POLIS` bercabang pada `QuotationData.BusinessFac` dan membaca jalur JSON yang berbeda (`$.PolicyData.StartDateTime` versus `$.StartDate`) — adalah **akibat, bukan sebab**. Ia pendukung yang baik dan bukan bukti utama. Menyimpulkan identitas dari bentuk penyimpanan akan bertentangan dengan `ADR-0015`, yang justru menyatukan dua tabel karena identitasnya sama.

**Yang membatasi kepemilikan**: lingkup modul ini adalah klaim non-proporsional, dan itu berjalan di jalur Treaty. Memodelkan Fakultatif berarti membangun entitas yang tidak pernah ditulis dan tidak pernah dihitung — biaya tanpa penerima manfaat.

**Pemeriksaan yang diminta, dan hasilnya tidak sebagaimana diharapkan.** Pertanyaannya: adakah klaim non-prop yang pernah berjalan di cabang `'F'`? Sapuan 279 berkas menemukan **`BusinessFac` nol kemunculan** — modul Claim **tidak pernah membacanya sama sekali**.

Artinya modul Claim tidak dapat membedakan kedua cabang; ia menerima apa pun yang dikeluarkan `V_POLIS`. Jadi pertanyaan itu **tidak dapat dijawab dari XML** dan menjadi pertanyaan data: apakah ada baris `json_polis` ber-`BusinessFac='F'` yang pernah dirujuk sebuah klaim non-prop. → **REQ-030**.

Konsekuensi yang perlu dicatat sekarang: bila ternyata ada, klaim itu menerima `STARTDATETIME`/`ENDDATETIME` berpresisi penuh, sementara jalur Treaty menerima bentuk yang terpotong di jam (`FINDING-005` bagian 6.1) — dua bentuk teks berbeda panjang masuk ke perbandingan yang sama.

Batas kepemilikan ini ditulis di `CONTEXT.md`.
