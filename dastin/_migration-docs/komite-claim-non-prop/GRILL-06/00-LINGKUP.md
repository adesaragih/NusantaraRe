> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : pembaca
> Masukan: HitServiceToKasirKMT_Act.xml · InsertXOLKlaimCNP.xml · KomitePostAdjustment.xml (ekspor Pega 2026-09-09) · PUTUSAN-01.md (2026-09-20) · KETETAPAN.md (2026-09-20) · GRILL-05/06-PUTUSAN.md M5-01, M5-02 (2026-09-20)
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

## 1. Bentuk ronde ini

Ronde 6 **bukan grilling**. Tidak ada temuan yang dicari dan tidak ada pertanyaan yang
diajukan. Empat pekerjaan, seluruhnya dari bahan yang sudah dipegang:

1. Membaca ulang empat temuan lama yang berdiri di atas sandi transisi, dengan alat yang
   sudah diperbaiki.
2. Membakukan istilah modul dalam `CONTEXT.md`.
3. Menulis tiga ADR dari ketetapan yang sudah ada.
4. Memutakhirkan lima berkas hidup.

## 2. Sebab pekerjaan 1 ada

`M5-01` menemukan bahwa sepasang `WhenTrue/WhenFalse` yang terbaca di jendela teks ternyata
milik langkah lain, dan bahwa alat baca menerjemahkan sandi `6` dengan kata yang dapat
dibaca dua arah. Bacaan itu nyaris membalik urutan wewenang seluruh modul dan hanya
tertahan oleh uji kenyataan.

Empat temuan lama berdiri di atas pembacaan yang sejenis dan belum pernah diperiksa dengan
penambat langkah: `G-03`, `G-08`, `K-01`, `K-07`.

## 3. Bukti yang dipakai

| Bukti | Asal | Tanggal | Untuk |
|---|---|---|---|
| `Komite Claim Non Prop/Activity/HitServiceToKasirKMT_Act.xml` | ekspor Pega | berkas 2026-09-09 | P6-1 (`G-03`) |
| `Komite Claim Non Prop/Activity/InsertXOLKlaimCNP.xml` | ekspor Pega | rule diubah 2024-11-26 | P6-2 (`G-08`) |
| `Komite Claim Non Prop/Activity/KomitePostAdjustment.xml` | ekspor Pega | rule diubah 2025-09-08 | P6-3 (`K-01`), P6-4 (`K-07`) |
| `PUTUSAN-01.md` baris 22, 31, 57, 104, 130 | berkas | 2026-09-20 | memastikan status lama tiap temuan, bukan mengandalkan ingatan atau ringkasan |
| `alat/dump_act.py` | scratchpad | diperbaiki 2026-09-20 sebelum pembacaan | seluruh P6 |

## 4. Bukti yang sengaja TIDAK dipakai

| Tidak dipakai | Alasan |
|---|---|
| Isi muatan kiriman Kasir sebagai dasar kesimpulan | pagar `PG-04`; strukturnya boleh dibaca, artinya tidak boleh disimpulkan |
| Bentuk simpan baris akseptasi dan layer | pagar `PG-05` |
| Isi surat dan dokumen | pagar `PG-06` |
| Ringkasan status temuan dari ronde sebelumnya | status lama diambil dari `PUTUSAN-01.md` langsung — ronde ini menemukan satu status yang selama ini dikutip keliru (lihat `01-PEMBACAAN.md` §5) |

## 5. Aliran yang beku

`A-4` (`PG-04`) · `A-5` (`PG-05`) · isi `A-6` (`PG-06`) · `A-3` terbatas (`PG-03`).
Rujukan lengkap di `REGISTER-PAGAR.md`. **Alasan pun tunduk pada pagar, bukan hanya
kesimpulan.**

## 6. Aturan berhenti

1. **Pembacaan berhenti pada struktur.** Yang dilaporkan: langkah mana pemilik sebuah
   pasangan transisi, dan sandi apa yang tertulis. Yang tidak dilaporkan: apa artinya bagi
   uang yang terkirim, bila itu berada di balik pagar.
2. **Uji kenyataan dijalankan sebelum menyimpulkan, bukan sesudah.** Bila bacaan baru
   menyiratkan sistem tidak mungkin berjalan seperti yang terbukti berjalan, yang dicurigai
   pertama adalah bacaannya.
3. **Temuan baru tidak dikerjakan.** Bila sesuatu muncul, ia dicatat di `07-AUDIT.md`
   bagian akhir sebagai di luar cakupan.
4. **Dugaan yang tertulis di instruksi tidak diistimewakan.** Instruksi ronde ini
   mengandaikan `G-03` akan runtuh dan `PG-04` akan menyusut. Andaian itu diuji seperti
   andaian lain, dan hasilnya dilaporkan apa adanya.
