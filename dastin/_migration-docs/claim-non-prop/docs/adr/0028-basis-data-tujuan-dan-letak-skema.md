---
status: accepted
label: DECIDED
---

# Basis data tujuan Oracle, skema baru bersebelahan dengan POOLDATA

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**; dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).
> Keputusan diambil 18 September 2026. Ini keputusan arsitektur pertama sistem baru; seluruh ADR di bawahnya mewarisinya.

> **Status `accepted` sejak 19 September 2026.**
>
> Versi pertama dokumen ini berstatus `proposed` dengan alasan yang ditulis terang: isinya disusun dari pilihan A yang tertulis lengkap sementara pilihan B dibiarkan kosong — **pembacaan atas formulir, bukan kalimat orang**. Syaratnya konfirmasi eksplisit beserta tanggal.
>
> **Konfirmasi itu datang sebagai kalimat manusia pada 19 September 2026**, dengan alasan yang dinyatakan penuturnya — dua, bukan satu:
>
> 1. **Pindah mesin membuat ADR-0016 kehilangan arti dan ADR-0023 mustahil.** Foreign data wrapper pada dasarnya database link dengan nama lain, dan view lintas mesin bukan view melainkan salinan yang menyamar. Keduanya harus ditulis ulang lebih dulu, dan tidak ada DDL yang boleh berdiri di atas ADR yang sudah tidak berlaku.
> 2. **Mengganti mesin berbarengan dengan migrasi menambah sumber perbedaan yang tidak dapat dipisahkan dari perbedaan yang sedang diukur ADR-0005.** Shadow-run hanya bermakna bila selisihnya punya satu sebab; dua perubahan serentak membuat setiap selisih ambigu.
>
> Alasan kedua tidak ada di versi pertama dokumen ini. Ia datang dari penutur, dan dicatat apa adanya.

Basis data tujuan adalah **Oracle**. Skema baru berdiri sebagai **satu skema tersendiri di instance yang sama dengan `POOLDATA`**.

Versi ditulis sebagai **sekurangnya 12.1, angka pastinya terbuka**, dan itu **tidak ikut naik** bersama status ini. ~~**Gerbang DDL terangkat; gerbang penamaan tidak.** Sampai REQ-032 kembali, tidak ada nama constraint, nama index, maupun singkatan yang boleh ditulis — batas panjang pengenal belum diketahui.~~ **Gerbang penamaan ikut dicabut 19 September 2026**, dan bukan karena REQ-032 kembali: batas **30 byte ditetapkan sebagai aturan tetap**, sah di setiap versi Oracle, sehingga jawaban REQ-032 tidak dapat mengubah satu nama pun. REQ-032 turun jadi verifikasi. Aturan penamaan final: `SPEC-MODEL-DATA.md` bagian 16. Lantai itu berdasar: `BLUEPRINT.md` §11 menyimpulkan basis data **sumber** berjalan di 12c ke atas dari tiga hal sekaligus — `CHECK (x IS JSON)`, notasi titik `a.data_json.Field`, dan klausa `CREATE OR REPLACE EDITIONABLE` yang baru ada sejak 12.1. **Itu lantai sumber, bukan versi tujuan**; keduanya tidak boleh dicampur, dan angka pasti tujuan masih ditunggu.

## Considered Options

Basis data lain tidak dipilih. Alasannya bukan selera melainkan biaya pada ADR yang sudah berdiri: pindah mesin membuat **ADR-0016** kehilangan arti — foreign data wrapper pada dasarnya database link dengan nama lain — dan membuat **ADR-0023** mustahil, karena view lintas mesin bukan view melainkan salinan yang menyamar. Keduanya harus ditulis ulang lebih dulu, dan tidak ada DDL yang boleh ditulis di atas ADR yang sudah tidak berlaku.

Instance terpisah juga tidak dipilih, atas alasan yang sama dalam bentuk lebih kecil: ADR-0023 mewajibkan sistem hilir diberi **view**, bukan salinan. Antar instance, view menuntut database link, dan itu dilarang ADR-0016.

## Consequences

### Satu pintu tulis tidak bertahan sebagai kesepakatan — ia ditegakkan lewat grant

Menaruh skema baru bersebelahan dengan `POOLDATA` berarti menaruhnya di tempat yang sistem lain sudah terbiasa menulisinya lewat SQL langsung. Itu bukan kekhawatiran teoretis: DDL tabel work Pega memuat

```sql
GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE, ... TO POOLDATA;
```

`POOLDATA` memegang `INSERT`, `UPDATE`, `DELETE`, bahkan `ALTER` atas tabel milik Pega. Apakah hak itu benar-benar dipakai belum diketahui (REQ-021, BLOCKER) — tetapi hak yang ada akan dipakai cepat atau lambat, dan **ADR-0017 yang hanya berupa kesepakatan akan runtuh di tempat seperti ini.**

Maka penegakannya masuk DDL, bukan catatan:

| Pihak | Hak |
|---|---|
| Akun pemilik skema | seluruhnya — hanya dipakai saat pemasangan dan migrasi |
| Akun aplikasi | `INSERT`, `UPDATE`, `DELETE`, `SELECT` atas tabel kanonik |
| Seluruh akun lain, termasuk `POOLDATA` | `SELECT` **hanya atas view** ADR-0023 — tidak ada satu pun hak atas tabel kanonik |

`GRANT` dan `REVOKE`-nya ditulis di DDL sebagai objek tersendiri, bukan dijalankan tangan. Hak yang tidak tertulis di DDL adalah hak yang tidak dapat diaudit.

### Versi 12.1 sebagai dasar, dengan penanda di tempat yang lebih ringkas di atasnya

DDL ditulis supaya berjalan di 12.1. Setiap tempat yang lebih ringkas di versi lebih baru ditandai di spec, tidak dipakai diam-diam. Yang menggigit paling awal adalah **batas panjang pengenal**: 30 byte sampai 12.1, 128 byte sejak 12.2. Nama constraint dan index menabraknya lebih dulu daripada nama tabel.

Selama versi pasti belum diketahui, spec wajib memuat **tabel singkatan tertutup** — daftar tetap, ditulis sekali, tidak boleh ditambah saat menulis DDL. Singkatan ad-hoc adalah cara `UR` dan `MDP` lahir, dan itu persis yang glosarium tutup.

### Tidak ada kolom JSON untuk data yang dimiliki

Seluruh alasan migrasi ini adalah memberi tiap nilai kolomnya sendiri. Sistem lama menyimpan nilai uang sebagai **teks di dalam JSON**, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2` (`BLUEPRINT.md` §13.2, §13.5).

JSON hanya boleh muncul di **tabel pendaratan migrasi**. Itu menjadikan tabel pendaratan satu-satunya tempat bentuk lama boleh masuk utuh, dan membuat aturan penguraian ADR-0014 bekerja di satu batas yang jelas — antara pendaratan dan kanonik — bukan tersebar di banyak tempat.

### Yang diwarisi ADR lain

**ADR-0016** tetap berlaku apa adanya: data pegawai masuk lewat API atau salinan tersinkron. Instance yang sama tidak mengubahnya — `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` yang ditemukan di `V_MST_USER_TEKNIS` adalah link ke mesin **lain**, dan itu yang dilarang.

**ADR-0023** menjadi dapat dilaksanakan: view lintas skema di satu instance adalah view sungguhan, tanpa salinan dan tanpa link.

**ADR-0017** berpindah dari kesepakatan menjadi konfigurasi yang bisa diperiksa, lewat tabel grant di atas.

## Akibat yang lahir bersama kenaikan status ini

### D20 naik dari catatan menjadi pekerjaan

Selama dokumen ini `proposed`, temuan **D20** — `POOLDATA` sudah memakai database link ke instance lain (`hrdasm.v_hrd_mst@asmd.sinarmas.co.id` di `V_MST_USER_TEKNIS`) — hanya catatan tentang sistem lama. Sekarang skema baru benar-benar berdiri di instance itu, dan ADR-0016 melarang apa yang sudah berjalan di sebelahnya.

~~**K1 harus diputuskan sebelum tiket `35` dikerjakan**, dan kedua cabangnya dibawa ke pemilik keputusan — **tidak dipilih sendiri**.~~

**Gerbang itu gugur 19 September 2026.** K1 ditutup sebagai **K1a**: ADR-0016 tetap berlaku tanpa revisi, karena modul ini **tidak mengambil data HRD sama sekali** — sehingga larangan database link tidak pernah bertabrakan dengan kebutuhan kita. Kekhawatiran yang dicatat di sini, bahwa K1b *"melemahkan alasan yang menopang ADR-0023 dan ADR-0028"*, **tidak terjadi**: K1b tidak dipilih, dan kedua ADR itu tidak ditinjau ulang. Rinciannya di badan ADR-0016.

### Penegakan ADR-0017 lewat grant adalah pekerjaan, bukan catatan

Satu pintu tulis tidak bertahan sebagai kesepakatan. Penegakannya lewat `GRANT` dan `REVOKE` masuk lingkup tiket `01` **sebagai bagian pekerjaannya**, bukan catatan operasional yang menyusul. Alasannya satu kalimat: **hak yang tidak tertulis tidak dapat diaudit.**
