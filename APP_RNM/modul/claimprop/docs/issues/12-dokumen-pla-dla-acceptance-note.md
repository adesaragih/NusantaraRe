# 12: Dokumen — PLA, DLA, dan Acceptance Note

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 08 (baris adjustment)
**Menutup:** AC 84 · 85 · 86 · 87 · 88 *(5 AC)* — US 45–49

## Hasil & nilai pengguna

Claim Admin mencetak PLA pada tahap outstanding, DLA per reinsurer per baris adjustment, dan
Acceptance Note otomatis saat akseptasi disimpan. Semua dokumen tersimpan dan dapat dibuka kembali
lewat tautan, sehingga tidak perlu mencetak ulang — dan nomor dokumen **tidak terbakar** saat
pembuatan gagal (dijamin tiket 00).

## Area codebase

Pembuatan PDF · penyimpanan berkas dan metadata · tautan berbatas waktu · gerbang anti-cetak-ganda.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/TryMakePLA_Act.xml` | 8 · 11 · 14 · 26 | PLA; nomor diambil lalu hanya ke clipboard, step lanjut dapat melempar exception, **nol `Obj-Save`** |
| `Activity/PrintDLATreatyIn.xml` | 7 | DLA; menarik nomor lewat generator DLA |
| `Activity/PrintDLATreatyIn.xml` | 15.11.1 | cabang yang **hidup** — diikuti apa adanya, termasuk pola halaman sementara berindeks tetap |
| `Activity/PrintDLATreatyIn.xml` | 15.11.2 · 15.11.3 | ⚠️ cabang salvage dan adjuster fee **di-remark** → **tidak dimigrasikan** |
| `Activity/PrintFileAcceptance.xml` | — | Acceptance Note tercetak saat akseptasi disimpan |
| `RDBList/GetTokenStorage_SQL.xml` · `RDBList/Insert_T_Storage_SQL.xml` | — | token → unggah → simpan metadata |
| `Activity/CloseClaimProp.xml` | — | ⚠️ **tidak** memeriksa nomor PLA maupun DLA — pencetakan bukan gerbang |

## ADR terkait

**ADR-0010** (penyimpanan berkas tetap di Google Storage) · **ADR-0006** (penomoran lewat stored
procedure).

## Acceptance criteria

- [ ] `[keputusan work owner]` **DLA = Definite Loss Advise** — satu PDF **per reinsurer per baris adjustment**, dengan saudara retro *(AC 84 spec)*
- [ ] `[terverifikasi]` Ketiga dokumen **disimpan** ke penyimpanan berkas dengan metadata di database; URL berbatas waktu dengan jalur penyegaran *(AC 85 spec)*
- [ ] `[terverifikasi]` **Pencetakan bukan gerbang** — penutupan klaim tidak memeriksa nomor dokumen. Ditiru apa adanya *(AC 86 spec)*
- [ ] `[terverifikasi]` Gerbang **anti-cetak-ganda** dipertahankan *(AC 87 spec)*
- [ ] ⚠️ **RISIKO DITERIMA SADAR — dipertahankan, bukan diperbaiki.** Cabang **salvage** (step 15.11.2) dan **adjuster fee** (step 15.11.3) pada DLA di-remark → **tidak dimigrasikan**, tanpa pengecualian; **step 15.11.1 diikuti apa adanya**, termasuk pola halaman sementara berindeks tetap. **Konsekuensinya diterima:** DLA untuk salvage dan adjuster fee menampilkan **nilai klaim**, sama seperti produksi hari ini. **Alasan menyimpang:** aturan induk — langkah ber-remark tidak ditulis sama sekali *(AC 88 spec)*

## Perintah verifikasi

```
jalankan test "cetak PLA -> berkas tersimpan, metadata tercatat, tautan dapat dibuka"
jalankan test "cetak DLA -> satu PDF per reinsurer per baris adjustment"
jalankan test "simpan akseptasi -> Acceptance Note tercetak otomatis"
jalankan test "pembuatan PDF gagal -> nomor dokumen TIDAK terbakar"
jalankan test "cetak dua kali -> gerbang anti-cetak-ganda menahan"
jalankan test "tutup klaim tanpa cetak PLA/DLA -> LOLOS (paritas)"
jalankan test "DLA salvage menampilkan nilai klaim (risiko diterima sadar)"
```
