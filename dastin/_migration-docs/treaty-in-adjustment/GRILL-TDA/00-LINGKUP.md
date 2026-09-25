> Modul  : Treaty In Adjustment · **ronde TDA** — sapuan penutup Lacak TDA (BUKAN cabang D)
> Dibuat : 2026-09-24
> Peran  : penggrill menyapu; AI grilling mengadili dengan kutipan; putusan milik pemilik proses
> Masukan: `KEPUTUSAN-GRILLING-ADJUSTMENT.md` (Lacak TDA) · `GRILL-A/06-PUTUSAN.md` ·
>          `GRILL-B/06-PUTUSAN.md` · `GRILL-C/06-PUTUSAN.md` · ADR induk 0034, 0037, 0040, 0044,
>          0052, 0055 · `SPEC-INVARIAN.md` (INV-24, INV-25)
> Status : **TERBUKA** — satu residu menunggu keputusan
> Sifat  : grilling. Tidak ada spesifikasi, DDL, struct Go, komponen React, endpoint, atau tiket.

# Ronde TDA — lingkup

## 1. Kenapa ronde ini ada

Ia **tidak direncanakan**. Ia lahir dari titik periksa `GRILL-C/07-AUDIT.md` §4a, yang mencocokkan
tabel **Lacak TDA** terhadap tabel **Ronde** di indeks — dua tabel di dalam berkas yang sama — dan
menemukan **tujuh TDA tanpa nasib**, seluruhnya ditugaskan kepada cabang yang sudah ditutup.

`METODE` §6.3 menuntut setiap TDA membawa nasib: **diperbaiki · dilestarikan · ditunda**.
`METODE` §7.2 menyatakan grilling selesai bila **semua bisa di-spec**.

> TDA tanpa nasib **tidak bisa di-spec.** Ia sampai ke sesi to-spec sebagai cacat yang tercatat,
> ditugaskan kepada cabang yang tidak akan pernah bersidang, dan **tidak ada yang akan tahu bahwa ia
> menunggu**.

Itu sebab ronde ini bukan kerapian.

## 2. Yang disapu

| TDA | Pokok | Ditugaskan ke | Keadaan cabang itu |
|---|---|---|---|
| `TDA-01` | penjaga duplikat mati; tabrakan jadi `UPDATE` yang menimpa | cabang B | DITUTUP |
| `TDA-07` | dua kontrol layar mengosongkan status akseptasi kontrak | cabang D | ditutup tanpa ronde |
| `TDA-09` | peran dari `pyWorkBasketList(2)` / `pyTelephone` / nama tersemat | cabang G | ditutup tanpa ronde |
| `TDA-11` | picker menyatukan kontrak dan addendum tanpa pembeda | cabang H | ditutup tanpa ronde |
| `TDA-12` | offset pengurai nomor revisi meleset satu | cabang B | DITUTUP |
| `TDA-13` | jenis dan materialitas dipilih bebas di radio picker | menunggu `EXP-1` | `EXP-1` DITUTUP |
| `TDA-15` | penggandaan pohon di dalam satu `JSONDATA` | cabang F | ditutup tanpa ronde |

## 3. Aturan sapuan — dan ia yang membuat ronde ini bukan stempel

Sebuah TDA **hanya** boleh dinyatakan tertutup bila ada **kutipan** dari putusan atau ADR yang
menutupnya, **beserta nomor barisnya**. Tiga bentuk keluaran yang sah, dan tidak ada yang keempat:

| Keluaran | Artinya |
|---|---|
| **diperbaiki oleh ‹putusan›** | putusan yang sudah terkunci menghapus sebab cacatnya. Kutipannya wajib |
| **KONFIRMASI ‹sumber›** | sudah dijawab ADR atau spesifikasi induk. Kutipannya wajib |
| **pertanyaan sungguhan** | tidak ada putusan yang menutupnya; ia masuk frontier |

Yang **tidak** sah: menyatakan tertutup karena *"cabangnya sudah ditutup"*. Cabang yang ditutup
tanpa ronde ditutup atas dasar **pemilahan**, bukan atas dasar TDA-nya — dan itu justru yang
membuat ketujuhnya lolos.

## 4. Hasil sapuan, dan ia berbeda dari perkiraan

Perkiraan di `GRILL-C/03-LUBANG.md` §4: **tiga** pertanyaan sungguhan — `TDA-07`, `TDA-11`,
`TDA-13`.

**Hasilnya satu.** `TDA-07` dan `TDA-11` keduanya tertutup oleh putusan yang sudah ada, dan
kutipannya ditemukan di `GRILL-A/06-PUTUSAN.md` dan `GRILL-B/06-PUTUSAN.md`. Perkiraan itu dicatat
apa adanya beserta salahnya, karena ia dibuat sebelum berkas rondenya dibaca — bentuk `TA-08`:
*baca temuan bernomor sendiri sebelum membangun rantai sebab*.

| | |
|---|---|
| **tertutup dengan kutipan** | `TDA-01`, `TDA-07`, `TDA-09`, `TDA-11`, `TDA-12`, `TDA-15` |
| **pertanyaan sungguhan** | **`TDA-13` residu** — lihat `02-SIDANG.md` ST-6 |

## 5. Aturan berhenti

Sama dengan ronde C, dan tidak satu pun terpicu sampai berkas ini ditulis.

> Selain ketiganya: kerjakan, catat, lanjut.
