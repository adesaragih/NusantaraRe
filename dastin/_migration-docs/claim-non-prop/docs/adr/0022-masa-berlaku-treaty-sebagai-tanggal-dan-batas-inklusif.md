---
status: accepted
label: DECIDED
---

# Masa berlaku treaty disimpan sebagai tanggal, dan batasnya inklusif

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Masa berlaku treaty disimpan sebagai **tanggal**, bukan tanggal-waktu. Kerugian yang terjadi **tepat pada** tanggal akhir treaty **tertutup** — batasnya inklusif.

## Alasan menyimpan sebagai tanggal — bukan yang saya tulis pertama kali

Alasan versi pertama saya: *"untuk masa berlaku treaty, jam tidak membawa arti"*. **Itu salah untuk reasuransi.** Slip treaty lazim menyebut waktu attachment secara eksplisit — misalnya *"at 00:01 hours Local Standard Time, 1st January"* — dan jam bisa sangat berarti untuk kerugian katastrofa yang terjadi di pergantian periode.

Alasan yang benar lebih sederhana dan lebih kuat:

> **Data sumbernya tidak sanggup menyimpan waktu yang dapat dipercaya.** Jalur Treaty di `V_POLIS` berhenti di jam — tanpa menit dan detik — dan nilainya sendiri hasil `SUBSTR` atas teks.

Perbedaan ini penting: bila kelak ada sumber yang membawa waktu attachment sungguhan, alasan versi pertama akan dipakai untuk menolaknya. Alasan yang benar justru mengundangnya.

## Batas inklusif — dan apa yang sebenarnya dilakukan sistem lama

Seluruh perbandingan di sistem lama memakai operator `>` tegas, tidak pernah `>=`:

```
pyWorkPage.ClaimData.DateOfLoss > pyWorkPage.ClaimData.EndDateTreaty
pyWorkPage.ClaimData.PolicyData.EndDateTime > pyWorkPage.ClaimData.EndDateTreaty
@FormatDateTime(...StartDateTime...) > pyWorkPage.ClaimData.EndDateTreaty
```

Karena penolakan terjadi ketika tanggal kerugian **lebih besar** dari tanggal akhir, kerugian yang jatuh **tepat pada** tanggal akhir tidak ditolak. **Maksud kode lama juga inklusif**, dan itu sejalan dengan praktik reasuransi.

Yang tidak dapat dipastikan adalah apakah maksud itu benar-benar terlaksana, karena perbandingannya dilakukan atas dua format teks yang berbeda (`FINDING-005`). Jadi keputusan ini menegaskan maksud yang sudah ada, bukan mengubahnya.

## Consequences

**Ditetapkan sebelum REQ-020 dijalankan, dan itu bukan formalitas.** Klaim yang tanggal kerugiannya jatuh tepat di batas akan **berpindah sisi** tergantung pilihan ini — masuk atau keluar dari kedua daftar di `ADR-0018`. Bila ditetapkan setelah query berjalan, daftarnya harus dibuat ulang.

Perbandingan di REQ-020 karena itu memakai `TO_DATE(DATEOFLOSS,'YYYYMMDD') > TERMINATION`, bukan `>=`.
