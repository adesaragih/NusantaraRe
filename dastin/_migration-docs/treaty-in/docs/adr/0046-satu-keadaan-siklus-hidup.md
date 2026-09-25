# ADR-0046 — Satu keadaan siklus hidup kontrak

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Lima properti tumpang tindih menjawab pertanyaan yang sama — apakah kontrak ini boleh diubah, dan
di mana ia berada dalam alurnya:

| Properti | Rujukan di ekspor |
|---|---|
| `Position` | 5.634 |
| `ViewState` | 661 |
| `StatusAkseptasi` | 165 |
| `IsEditData` | 125 |
| `RevisionState` | 21 |

`ViewState` disetel **imperatif** di lima tempat, termasuk satu aturan yang namanya sendiri sudah
mengakui bendanya: `TreatyInForceEdit`. Padahal isinya sepenuhnya bisa diturunkan dari posisi
persetujuan dan mode revisi kontrak itu. Ia turunan yang disimpan dan disetel oleh klik.

Perbandingannya ditulis dengan **delapan ejaan berbeda** untuk maksud yang sama, sehingga satu
salah ketik membuka atau mengunci sebuah field tanpa ada yang tahu.

## Keputusan

Kelima bendera dilebur menjadi **satu keadaan siklus hidup** kontrak: satu nilai bertipe tegas
dengan himpunan nilai tertutup.

**"Boleh diubah" adalah turunan** dari keadaan itu, dihitung saat ditanya, dan tidak pernah
disimpan sebagai bendera tersendiri.

`ViewState` **tidak dinamai ulang — ia tidak ada di model baru.** Menamainya ulang hanya akan
melestarikan bendanya dan memperbaiki labelnya.

## Konsekuensi

- Delapan ejaan perbandingan tidak direkonsiliasi dan tidak dicari mana yang benar — tidak ada yang
  diangkut.
- Salah ketik tidak mungkin lolos, karena himpunan nilainya tertutup.
- Ini instans ADR-0037: turunan tidak disimpan dan tidak disetel oleh klik.
