# 08: Perbedaan perhitungan sebaran dan rincian angsuran

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` ⛔ **Tabel sebaran tambahan DIBATALKAN** — keempat medannya turunan, dua dari master `POOLDATA.PROPORTIONALARRG`, dua dihitung dari total yang sudah tersimpan. ⚠️ **Lingkup tiket ini menyusut**: bagian sebaran tambahan keluar, bagian rincian angsuran tetap berlaku.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 3.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** ⛔ `[work owner]` **beda dagang tabel sebaran tambahan** dari sebaran risiko belum dijelaskan~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **19** *(pemecah dokumen menjadi baris)*
**Menutup:** AC **35–38** · AC **56** *(5 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-39 · ID-40

## Hasil & nilai pengguna

Tabelnya sama dengan polis baru; **perhitungannya berbeda**. Tiket ini menyatakan perbedaan itu
dengan tepat, dan membuatnya **diuji terpisah** — bukan diseragamkan.

⚠️ **Perbedaan ini ditiru apa adanya.** Ia perilaku lama yang disalin sadar, bukan cacat yang
diperbaiki.

## Yang dibangun

| Perbedaan | Endorsemen | Polis baru |
| --- | --- | --- |
| presisi pembagian sebaran | **20** | **10** |
| baris pertama sebaran | menerima nilai bawaan **100** bila persentasenya kosong | — |
| rincian angsuran bertingkat | **hidup** pada bentuk non-proporsional | tidak ada |
| asal rincian angsuran | ⭐ dapat berasal dari **salinan master kontrak** | diketik pengguna |

⭐ **Sebab endorsemen memerlukan nilai bawaan dan polis baru tidak:** `[terverifikasi]` aturan yang
**mengisi** medan pembagi itu ada di modul polis baru dan **tidak ada** di modul endorsemen —
sehingga di endorsemen medannya dapat tiba kosong, sementara rumusnya **membaginya**.

~~Ditambah satu tabel sebaran tambahan dengan empat medannya.~~ ⛔ **DICABUT 23-09-2026 sore** — sebaran tambahan **tidak dimigrasi sama sekali**, tabel maupun perhitungannya.

## Batas — yang TIDAK termasuk

⛔ Pemecahan dokumen menjadi baris sebaran — tiket NB **19**.
⛔ Selisih sebaran antar generasi — tiket **07**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji berdampingan wajib:** kasus yang sama disimpan sebagai polis baru
dan sebagai endorsemen — hasilnya **harus berbeda** sesuai ketetapan presisi.
⛔ Uji yang menemukan keduanya sama berarti perbedaannya hilang.

## Acceptance criteria

- [ ] **AC 35** — sebaran endorsemen memakai presisi **20**; polis baru **10**
- [ ] **AC 36** — baris pertama sebaran menerima nilai bawaan **100** ketika persentasenya kosong
- [ ] **AC 37** — rincian angsuran bertingkat **tersimpan** pada endorsemen non-proporsional
- [ ] **AC 38** — rincian dari salinan master kontrak tersimpan **sama** seperti yang diketik
- [x] ~~**AC 56** — tabel sebaran tambahan menyimpan empat medan~~ ⛔ **AC DITARIK DARI LINGKUP.** Cacah AC spec EDM **58 → 57**.

## ⛔ Kenapa tiket ini `blocked`

~~`[work owner]` **Beda dagang tabel sebaran tambahan dari tabel sebaran risiko belum dijelaskan.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Itu tidak ada"* lalu *"itu tidak usah di migrasi perhitungannya"*.

⭐ Terbukti dari korpus: keempat medannya **turunan**. `TreatyType` dan `SharePercentage` diambil dari master `POOLDATA.PROPORTIONALARRG`; `PremiumSpreaded` dan `ClaimSpreaded` dihitung `TotalPremium × SharePercentage/100` dan `TotalClaim × SharePercentage/100`.

⇒ **Lingkup tiket ini menyusut.** Yang tetap berlaku: AC 35 · 36 *(penyebaran risiko)* dan AC 37 · 38 *(rincian angsuran)* — seluruhnya tidak terdampak.

⭐ **AC 35–38 dapat dikerjakan tanpa menunggu** — keempatnya tidak menyentuh tabel itu.

⚠️ **Dan satu ketergantungan yang wajib disebut:** rincian angsuran dapat berasal dari **salinan
master kontrak**, yang berarti ketergantungan pada dokumen kontrak **tetap ada** walau dokumen itu
di luar lingkup penyimpanan.