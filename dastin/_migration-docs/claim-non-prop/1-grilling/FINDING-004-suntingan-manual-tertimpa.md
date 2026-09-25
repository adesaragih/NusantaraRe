# FINDING-004 — Hasil suntingan manual alokasi tidak terlindungi

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

**Jenis**: laporan kondisi sistem lama (bukan keputusan migrasi)
**Status**: BERSYARAT — berlaku bila hipotesis F2 benar, yaitu `.IsEditClaim` memang tidak punya pembaca di mana pun
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

Dokumen ini terpisah dari keputusan desain untuk sistem baru. Keputusan itu ada di ADR-0008 dan berlaku apa pun hasil temuan ini.

---

## 1. Mekanismenya

`IsEditClaim` adalah `Rule-Obj-Property` milik class `ASM-FW-GISFW-Data-SpreadingRisk`.

| Rule | Aksi |
|---|---|
| `Activity\EditXOLAlokasi.xml` | `.IsEditClaim = 1` — ditandai saat petugas menyunting alokasi XOL |
| `Activity\CountLossAllocation_act.xml` (dua langkah) | `SpreadingRisk(<LAST>).IsEditClaim = 0` — ditimpa kembali menjadi 0 |

Lima kemunculan di seluruh 279 berkas, tidak lebih. **Tidak ada satu pun pembacaan**: tidak di `pyStepsPreCondParamsWhen`, tidak di `pyVisible` section mana pun, tidak di Report Definition, tidak di SQL.

## 2. Akibatnya bila memang tidak ada pembaca

Bendera yang hanya ditulis tidak menjaga apa pun. Rangkaiannya:

1. Petugas menyunting nilai alokasi lewat `EditXOLAlokasi`.
2. Sistem menandai baris itu `IsEditClaim = 1`.
3. Perhitungan alokasi dijalankan lagi.
4. `CountLossAllocation_act` menghitung ulang nilainya **dan** menyetel penandanya kembali ke `0`.
5. Tidak ada peringatan, tidak ada pesan, tidak ada jejak bahwa suntingan pernah ada.

**Berapa besar peluang langkah 3 terjadi:** `CountLossAllocation_act` dapat dipicu dari tiga jalur yang terbaca di XML — `CountClaimTNP_Act` langkah 15, `AddAkseptasiCNP_Act` langkah 18, `InputAkseptasi_PreAct` langkah 5 — dan **juga langsung dari kontrol UI** pada tiga Section (`AdjustmentDetailNP`, `InputAcceptation`, `OutstandingClaim(1)`). Lihat FINDING-003 bagian 3.

Karena pemicunya adalah kontrol layar, petugas yang membuka kembali layar yang sama dan menekan kontrol hitung akan menghapus suntingannya sendiri tanpa menyadarinya.

## 3. Yang membuat ini bukan sekadar bug kecil

Fitur suntingnya **ada dan berfungsi** — layarnya ada, penandanya ditulis. Yang tidak ada adalah pihak yang membaca penanda itu. Bagi penggunanya, sistem tampak menyediakan kemampuan yang sebenarnya tidak dijamin.

Ini berbeda dari fitur yang tidak ada sama sekali: fitur yang tidak ada tidak menyesatkan siapa pun.

## 4. Konsekuensi yang meringankan pekerjaan migrasi

Satu akibat yang justru menyederhanakan, dan perlu dinyatakan supaya tidak dikerjakan dua kali:

> **Bila suntingan selalu tertimpa, maka nilai yang tersimpan di produksi hari ini adalah nilai hasil hitungan, bukan nilai hasil suntingan.**

Artinya **tidak ada "nilai tersunting" yang perlu diselamatkan** saat migrasi data. Tidak perlu kolom khusus, tidak perlu rekonsiliasi, tidak perlu menanyakan ke pengguna nilai mana yang benar. Satu cabang pekerjaan migrasi data terpotong.

Ini berlaku **hanya bila hipotesis F2 benar.** Bila F1 yang benar (pembacanya ada di modul Komite), kesimpulan ini gugur dan nilai tersunting mungkin memang bertahan.

## 5. Batas klaim

1. **Hipotesis F1 belum disingkirkan.** `.IsEditClaim` mungkin dibaca oleh rule di modul Komite. Folder Komite tidak dibuka pada sesi ini. Lihat `_selesai/OPEN-QUESTIONS.md` bagian F.
2. **Ketiadaan pembacaan disimpulkan dari ketiadaan rule, bukan dari adanya rule.** Ekspor ini memuat 279 berkas dari satu folder; pembacaan bisa saja berada di rule yang tidak ikut terekspor.
3. **Belum diuji apakah suntingan pernah benar-benar dilakukan di produksi.** Kalau fitur `EditXOLAlokasi` tidak pernah dipakai, temuan ini benar secara mekanis tetapi nol dampaknya.

## 6. Yang menutup temuan ini

**Satu pemeriksaan tunggal**: sapu modul Komite atas `IsEditClaim`. Satu sesi, satu grep. Hasilnya menentukan F1 atau F2, dan sekaligus menutup atau menegakkan seluruh dokumen ini.

Sampai itu terjadi, **tidak boleh** ada keputusan desain yang mengandaikan salah satunya benar. Keputusan untuk sistem baru sudah diambil terpisah dan tidak menunggu hasil ini — lihat ADR-0008.
