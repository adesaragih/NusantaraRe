# 07: Worker pengirim — retry + anti-dobel Email & Kasir

**Status:** ready-for-agent

**Blocked by:** 06 (transactional outbox)

## Hasil & nilai pengguna

Sebagai **Finance**, saya ingin klaim yang disetujui komite **pasti** sampai ke Kasir — dan **tidak
pernah dua kali**. Sebagai **operator sistem**, saya ingin efek yang gagal diantre ulang sampai
berhasil, bukan menguap. *(User story 22–25, 27 di spec)*

## Area codebase

worker (proses terpisah — pengambil entri outbox, pengirim, penjadwal retry),
`internal/services` (interface keempat efek keluar; aturan cek-status-sebelum-kirim-ulang),
`internal/repository` (lookup `M_LINK_SERVICE`; pembaruan keadaan entri outbox).

Tidak menyentuh `frontend/` — permukaan manusia ada di tiket 08.

## Rule Pega sumber

| Efek | Rule | Identitas |
| --- | --- | --- |
| 1 · step 8 | `Komite Claim Life/Activity/InsertJsonClaimLife_Act.xml` | `RULE-OBJ-ACTIVITY` |
| 2 · step 10 | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | `[terverifikasi]` **satu-satunya salinan di korpus** — OQ-035 |
| 3 · step 11 | `Komite Claim Life/Activity/SendEmailKlaimLife.xml` | `RULE-OBJ-ACTIVITY` |
| 4 · step 12 | `Komite Claim Life/Activity/HitServiceToKasirKMTLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `HITSERVICETOKASIRKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY` — **Kasir/pembayaran** |
| resolusi alamat | `Komite Claim Life/Activity/GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` — `Obj-Browse` `M_LINK_SERVICE` pada `(KATEGORI_1, KATEGORI_2)`, ambil `.URL`, lalu `Connect-REST` |
| jejak Kasir | `POOLDATA.DIRECTTOKASIR_LOG` | `[terverifikasi]` jejak di sisi database |

`[terverifikasi]` Kunci kategori Arasapas: `Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`
(`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`).

## ADR terkait

**ADR-0015** (at-least-once; ID idempoten unik; cek status sukses sebelum kirim ulang untuk Email &
Kasir — **menyimpang dari ADR-0008**), **ADR-0013** (alamat di-lookup runtime; **dilarang** URL
sebagai literal, konstanta, maupun env var), **ADR-0005** (flag lingkungan menggerbangi efek keluar).

## Acceptance criteria

- [ ] Keempat efek dikirim sampai **berhasil**, atau ditandai **perlu intervensi**. *(AC 20 spec)*
- [ ] Tiap kiriman membawa **ID idempoten unik**, sehingga penerima dapat menolak duplikat.
      *(AC 21 spec)*
- [ ] Sebelum mengirim ulang **Email** atau **Kasir**, sistem memeriksa status "sudah terkirim
      sukses"; bila sudah, **tidak dikirim lagi**. *(AC 22 spec)*
- [ ] Alamat endpoint di-resolve **runtime** dari `M_LINK_SERVICE` lewat `(KATEGORI_1, KATEGORI_2)`;
      **tidak ada URL** sebagai literal, konstanta, maupun env var. *(AC 24 spec; **ADR-0013**)*
- [ ] **Kunci kategori tidak ditemukan** dibedakan dari **jaringan gagal**: yang pertama tidak
      diulang berkali-kali, yang kedua diulang.
- [ ] Keputusan dilaporkan **tuntas** hanya setelah keempat efek berhasil. *(AC 25 spec)*
- [ ] **Semua** klaim menjalankan keempat efek — tidak ada pengecualian berdasarkan identitas retro.
      *(AC 26 spec; **OQ-064**)*
- [ ] Kegagalan pengiriman **tidak** membatalkan keputusan yang sudah tersimpan. *(AC 18 spec)*

## Catatan — gerbang EXIT retro dibuang

⚠️ `[terverifikasi]` Di korpus, **step 9** adalah langkah gerbang tersendiri berlabel
`EXIT JIKA RETROID "L0000141"` (baris 8483), dengan precondition
`RetroID=="L0000141" || SecurityReinsurerID=="L0000134"` (baris 8525) dan `RetroID=="1000013"`
(baris 8548). Ia berdiri **di antara** efek 1 dan efek 2 — saat memicu, step 10–14 tidak berjalan.

`[keputusan work owner]` **Dibuang** (**OQ-064**). **Konsekuensi yang diterima:** klaim ber-retro
tersebut yang selama ini dikecualikan kini menjalankan **seluruh** efek keluar, **termasuk Kasir**.
Ini **perubahan perilaku yang menyentuh uang**, diterima sadar.

`[dugaan]` Bahwa kode transisi numerik step 9 berarti "Exit Activity" **tidak terbukti dari ekspor**
— nilai yang terbaca (`2`) muncul juga di step 8 yang bukan EXIT. Yang terbukti adalah label
penulisnya dan bentuk langkahnya. Tidak berpengaruh pada implementasi: gerbangnya dibuang.

## Catatan — seam kedua

`[asumsi]` Perilaku retry dan anti-dobel **tidak teramati lewat API HTTP saja**, karena pengiriman
berjalan di proses terpisah. Tiket ini mengasumsikan **seam kedua di batas pengirim outbox**, dengan
keempat layanan luar **di-fake**. Bila seam itu tidak disetujui, cara mengujinya perlu ditetapkan
ulang sebelum tiket dikerjakan.

`[terbuka]` **OQ-002** — kontrak layanan **Kasir** tidak ada di korpus: bentuk permintaan, makna
jawaban, dan apakah ia menghormati ID idempoten **belum diketahui**. Tiket ini menetapkan **jaminan
pengirimannya**, bukan bentuk pesannya.

`[terbuka]` **OQ-035** — `serviceInsertArasapasClaimLife_act` satu salinan dipakai dua konteks.
**Bukan pemblokir isi.**

## Perintah verifikasi

```
go test ./internal/...
make check
```
