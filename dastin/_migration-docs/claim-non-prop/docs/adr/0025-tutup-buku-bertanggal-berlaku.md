---
status: accepted
label: DECIDED
---

# Tanggal tutup buku disimpan bertanggal berlaku, bukan satu baris tanpa riwayat

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Aturan penggeseran periode dipertahankan: nomor yang diterbitkan setelah hari tutup buku masuk periode berikutnya. Yang dirancang ulang adalah **penyimpanannya**.

Tanggal tutup buku disimpan dalam tabel **bertanggal berlaku** — setiap perubahan menjadi baris baru dengan periode berlakunya sendiri. Ditambah kolom **lingkup**, dengan nilai bawaan global.

## Consequences

Sistem lama menyimpannya sebagai **satu baris tanpa riwayat**: `SELECT TO_NUMBER(tanggal) INTO v_day_closing FROM POOLDATA.TANGGAL_CLOSING WHERE ROWNUM = 1`.

**Dan itulah sebabnya tambalan lahir.** `PROC_GENERATE_SEQUENCE_NUMBER` memuat:

```sql
IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN
    v_mm_yyyy := '12.2025';  v_tahun := 2025;
```

Perlakuan khusus untuk satu periode tidak punya tempat di data, jadi ia menjadi cabang di dalam kode. Dengan tanggal berlaku, kasus seperti itu cukup menjadi **satu baris data** — dan kelas tambalan ini kehilangan alasan untuk lahir lagi.

Kolom lingkup disediakan sekarang meski belum ada yang memerlukannya; biayanya hampir nol di tahap ini dan mahal ditambahkan setelah tabel terisi.

**Satu masalah yang tidak diselesaikan oleh keputusan ini** dan harus ditangani terpisah: aturan tanggal 25 ada di **dua sumber kebenaran**. `Activity\HitServiceToKasir_Act.xml` menuliskan `@if(TempKasir.CARI19 > 25, ...)` langsung di kode Pega, sementara prosedur membacanya dari tabel. Mengubah tabel tidak mengubah kode. Dibawa ke akuntansi sebagai `ASK-AKUNTANSI.md` pertanyaan 6.3.
