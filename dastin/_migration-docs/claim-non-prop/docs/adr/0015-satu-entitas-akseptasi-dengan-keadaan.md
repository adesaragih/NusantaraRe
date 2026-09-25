---
status: accepted
label: DECIDED
---

# Akseptasi adalah satu entitas dengan keadaan, bukan dua tabel

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Akseptasi biasa dan akseptasi bersyarat (Subjectivity) adalah **satu entitas** dengan keadaan yang berbeda, bukan dua entitas. Keadaannya: biasa, bersyarat, dan bersyarat-gugur — beserta daftar syarat dan status pemenuhannya.

Sistem lama memisahkannya jadi dua tabel: `OS_AKSEPTASI_KLAIM` dan `OS_AKSEPTASI_SUBJECTIVITY`.

## Consequences

Daur hidup yang berbeda adalah **keadaan**, bukan identitas. Dua hal menjadi entitas berbeda bila identitasnya berbeda; bila satu akseptasi bersyarat yang syaratnya terpenuhi berubah menjadi akseptasi biasa **yang sama**, maka sejak awal ia satu benda yang sedang berada di keadaan tertentu.

Bukti yang menguatkan datang dari sistem lama sendiri: view `CLAIMXOL` harus meng-`UNION ALL` kedua tabel untuk mendapat gambaran utuh. Sesuatu yang harus selalu digabungkan kembali tidak seharusnya dipisah.

Syarat yang gugur menjadi **transisi keadaan yang tercatat**, bukan penghapusan atau perpindahan tabel.

**Temuan yang menyertainya dan belum terjawab**: di folder Claim, `.IsSubjectivity` hanya **dibaca** (`==true`), disalin ke `Local.Subjectivity`, dan **tidak ada satu pun rule yang menggugurkan, membatalkan, atau mengakhiri akseptasi bersyarat yang syaratnya tidak pernah terpenuhi.** Delapan belas kemunculan `Subjectivity` di 279 berkas, tidak satu pun berupa transisi keadaan. Jadi dugaan "menggantung selamanya" konsisten dengan isi folder ini — *disimpulkan dari ketiadaan rule, bukan dari adanya rule*. Pemastiannya `DEFERRED-TO-KOMITE-SESSION`.
