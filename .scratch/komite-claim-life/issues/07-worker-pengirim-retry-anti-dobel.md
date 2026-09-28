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

## Implementasi — 28-09-2026 (giliran 10)

### Yang dipakai ulang, bukan disalin

`PekerjaEfek` Claim Life (`antrean.go`): pungut `FOR UPDATE SKIP LOCKED` → jalankan → tuntaskan, satu
transaksi; `Backoff`; `percobaanMaksimum = 8`; jejak menyerah. Sejak temuan `/code-review` giliran ini
pemungutannya **disaring `MODUL`** — pekerja Komite (`NewPekerjaEfekModul(…, "KOMITELIFE")`) hanya
memungut baris Komite, dan pekerja Claim Life tidak menyentuhnya. Aditif pada berkas Claim Life:
`NewPekerjaEfekModul`, dan `ErrKasirBelumDisetujui` di daftar galat permanen `LayakDicobaUlang`.

### Yang dibangun — `services/komite_pengirim.go`

- `PelaksanaKomite.Laksanakan` untuk `arasapas-komite`, `email-komite`, `kasir-komite`:
  - **Anti-dobel Email & Kasir** (AC 22): sebelum mengirim, `EfekSudahSelesai` memeriksa baris **lain**
    kasus yang sama berjenis sama yang sudah `selesai` → bila ada, baris ini tuntas **tanpa** kirim.
    Baris yang sudah `selesai` sendiri tidak pernah dipungut lagi (pemungutan hanya `antre`). Kiriman
    membawa ID baris (`SEQ_LOG_SERVICE_RNM`) sebagai kunci idempoten.
  - **Kunci hilang ≠ jaringan gagal**: `ErrEndpointTidakDitemukan` dan `…BelumDisetujui` permanen
    (gagal-permanen → tiket 08), galat lain dicoba ulang.
  - **Non-produksi = pengirim stub yang mencatat** (km4, ADR-0005): baris tuntas, nol resolusi alamat,
    nol panggilan keluar. Produksi: alamat di-resolve sungguhan (`M_LINK_SERVICE`, ADR-0013) lalu
    berhenti terang sampai panggilan nyata disetujui manusia.
- `PekerjaKomiteOracle(svc)` — pekerja siap pakai.
- Gerbang retro langkah 9 **dibuang** (`[keputusan work owner]`, OQ-064): nol cabang retro di pelaksana.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| dikirim sampai berhasil, atau perlu intervensi | ✅ ulang sampai 8, lalu gagal-permanen + jejak |
| ID idempoten unik | ✅ ID baris outbox |
| cek sudah-sukses sebelum kirim ulang Email/Kasir | ✅ `TestKasirTidakDikirimDuaKali` |
| alamat runtime dari `M_LINK_SERVICE`, nol URL | ✅ |
| kunci hilang tidak diulang, jaringan diulang | ✅ `TestKunciHilangPermanenJaringanBukan` |
| tuntas hanya sesudah keempat efek berhasil | ⚠️ keadaan per efek ada di outbox; ringkasan "tuntas/perlu intervensi" per kasus = tiket 08 |
| semua klaim menjalankan keempat efek | ✅ gerbang retro dibuang · ⚠️ efek 1 (`InsertJsonClaimLife_Act`) JSON — dibuang, jadi **tiga** efek |
| kegagalan kirim tidak membatalkan keputusan | ✅ pekerja berjalan di transaksinya sendiri |
| ⚠️ pekerja berjalan terjadwal | **belum ada penjadwal** di `cmd/` — `SatuPutaran` siap dipanggil; menjalankannya terus-menerus (proses terpisah) adalah keputusan operasi |

### Angka

Go **588 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **357** · tsc bersih.

### Temuan `/code-review` 28-09-2026 (Standards) — diperbaiki

- ⛔ **Anti-dobel menelan email tingkat 2..n**: email diantre setiap keputusan dengan `RUJUKAN = kasus`, jadi email tingkat 1 yang `selesai` membuat sisanya dilewati. Kini `RUJUKAN` email = `kasus#T<tingkat>` (`RujukanEfekKomite`); pembaca efek kasus mencocokkan keduanya; laporan harian memotong akhiran.
- ⛔ **Stub non-produksi tampak seperti kiriman nyata**: stub dulu menuntaskan baris `selesai`, dan anti-dobel menganggap Kasir sudah dibayar bila basis data non-produksi dipromosikan. Kini stub gagal **permanen** dengan `ErrPengirimStubNonProduksi` — barisnya tidak pernah `selesai` tanpa kiriman. ⚠️ Harganya: di DEV setiap efek Komite tampil "perlu intervensi" dengan sebab yang menyebut dirinya stub.
- Dicatat, tidak diubah: pekerja belum punya penjadwal (sama dengan Claim Life); `periksaGiliran` di penyelesai akhir membaca snapshot yang sama (gerbang nyata = UPDATE bersyarat); `PastikanKasusTerbuka` dibaca di luar transaksi (pola Claim Life); setter/`RowsAffected` berulang; dua definisi tingkat berjalan (`MIN(KOMITE_URUT)` lawan `KOMITE_COUNT`) yang dijaga sama oleh penulis bersyarat.
