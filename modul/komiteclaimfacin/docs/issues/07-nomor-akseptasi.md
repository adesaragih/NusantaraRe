# 07: Nomor akseptasi — satu rangkaian, dikunci kelas kasus induk

**Status:** ready-for-agent
**Blocked by:** 06
**Menutup:** AC 52 · 53 · 54 · 55 *(4 AC)* — US 28

## Hasil & nilai pengguna

Hari ini Nomor akseptasi **belum terbit**, dan di korpus ada **tujuh pembangkit per lini usaha** yang seluruhnya **ber-remark** — sehingga tidak jelas mana yang berlaku.

Sesudah tiket ini, ⭐ Nomor akseptasi terbit dari **satu rangkaian**, dikunci pada **jenis kasus induk** — ⛔ bukan per lini usaha — dan **tidak dibangkitkan ulang** untuk akseptasi yang sudah bernomor.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pembangkit yang hidup | satu prosedur tersimpan, dikunci **kelas kasus induk** |
| ⛔ Tujuh pembangkit per lini | ⛔ **ber-remark seluruhnya**; sisiran menyeluruh **tidak menemukan penggantinya** |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Pembangkit yang hidup | satu prosedur tersimpan,
> dikunci **kelas kasus induk**"* → pembangkitnya kini terbaca utuh di `KomitePost_Adjustment` S7.2.1.2.4–7:
>
> - S7.2.1.2.4 `GetKodeProdNonLife_SQL` → `ParamSeq.HASIL3` = `KODE_PRODUKSI.KODE` ber-`TYPE = 'NONLIFE'`;
> - S7.2.1.2.5 `ParamSeq.CARI1 := pyWorkCover.pxObjClass` (kunci rangkaian = kelas kasus induk),
>   `ParamSeq.CARI2 := HASIL3 + "A"`;
> - S7.2.1.2.6 `GetSequenceNumber_SQL` → `HASIL1` / `HASIL2`;
> - S7.2.1.2.7 `OutputData.START_DATE := CARI2 + OfferFacIn.QuotationData.BusinessOldId + "." + HASIL1 + "." + HASIL2`,
>   lalu S7.2.1.5 `.AcceptedNo := OutputData.START_DATE`.
>
> Formatnya **`KODE_PRODUKSI(NONLIFE) + "A" + BusinessOldId + "." + HASIL1 + "." + HASIL2`**, **tanpa `.TP`** (Komite
> Prop) atau `.TX` (Komite Non Prop). Panjang 21 / 22 cocok dengan gerbang kasir `CLM` (`HitServiceToKasirKMT_Act`
> L14.2 `@length(.AcceptedNo) = "21" || "22"`). `GetSequenceNumber_SQL` memanggil `PROC_GENERATE_SEQUENCE_NUMBER` +
> `COMMIT`, jadi penomoran berjalan **lewat `inti/backend/penomor` di dalam transaksi Submit** (ADR-0043
> `OUTPUT_HASIL_RNM/docs/bersama/adr/0043-penomoran-di-aplikasi.md`, menggantikan ADR-0006); celah `CARI3` mengikuti
> penanganan claimfacin tahap 1. Satu nomor per KMT (S7.2.1.2 bergerbang `AcceptanceStatus` "" / "0" dan tingkat akhir).

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K1** — ⭐ Nomor terbit sesudah ringkasan tersimpan, di **jenjang terakhir**

## Yang harus diuji

- [ ] ⭐ Nomor berasal dari **satu rangkaian**, dikunci **jenis kasus induk**
- [ ] ⛔ Tujuh pembangkit per lini usaha **tidak dialihkan**
- [ ] ⛔ Nomor yang sudah terbit **tidak dibangkitkan ulang**
- [ ] ⚠️ Perilaku prosedur tersimpannya **belum terbaca dari korpus** — ⭐ sistem baru **memerikan sendiri** aturan pembangkitannya, ⛔ tidak mewarisi yang tak terbaca

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **31** | Isi prosedur tersimpan pembangkit nomor — `[data DBA]` | ⚠️ menahan **peniruan persis**; ⭐ tidak menahan pembuatan aturan sendiri |

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**K1** — ⭐ Nomor terbit sesudah ringkasan
> tersimpan, di **jenjang terakhir**"* dan sel *"⚠️ menahan **peniruan persis**; ⭐ tidak menahan pembuatan aturan
> sendiri"* → di XML nomor terbit **sebelum** ringkasan disimpan (S7.2.1.2 sebelum `SaveAccept_ACT` S7.2.1.17 dan OS
> S8), tetap di jenjang terakhir yang setuju. Butir 31 tidak menahan: penomoran pindah ke aplikasi (ADR-0043), dan
> `inti/backend/penomor` adalah penggantinya. Seam langkah 3 dibaca begini: Submit ditolak bila `AcceptedNo` sudah terisi
> (`ShowTransfer` tombol Submit NA `pyWorkPage.Adjustment.AcceptedNo != ''`).

## Seam & verifikasi

**Seam:** lapisan layanan komite — penerbitan nomor.
1. Selesaikan dua kasus komite ⇒ ⭐ nomornya **berurut dalam satu rangkaian**.
2. Selesaikan kasus dua lini usaha berbeda ⇒ ⭐ keduanya dari **rangkaian yang sama**.
3. Coba terbitkan ulang untuk akseptasi yang sudah bernomor ⇒ ⛔ **ditolak**.
