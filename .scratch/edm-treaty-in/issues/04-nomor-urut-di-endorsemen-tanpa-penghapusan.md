# 04: Nomor urut di endorsemen — tanpa penghapusan

**Status:** ready-for-agent
**Blocked by:** **03**
**Bergantung pada tiket NB:** **17** *(nomor urut baris anak)*
**Menutup:** AC **10–13** *(4 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-16 · ID-17 · ID-18 · ID-34

## Hasil & nilai pengguna

Di endorsemen, baris rincian **tidak dapat dihapus**, dan nomor urut lama **terbawa apa adanya**.

⚠️ **Ini berbeda dari polis baru**, tempat baris boleh dihapus dan nomornya dinomori ulang rapat —
aman di sana karena belum ada generasi untuk dipasangkan. ⭐ Begitu generasi pertama ditutup, nomor
urut **beku selamanya**.

## Yang dibangun

Aturan penomoran khas endorsemen:

| Aturan | Sebab |
| --- | --- |
| baris rincian **tidak dapat dihapus** | pemasangan antar generasi lewat nomor urut menuntutnya |
| nomor urut lama **terbawa apa adanya** | pasangan lama tidak boleh bergeser |
| baris baru mendapat **maksimum + 1** | penyisipan di tengah menggeser seluruh pasangan sesudahnya |

⭐ **Akibat yang menguntungkan:** karena tidak ada penghapusan, pasangan yang bergeser berarti
**anomali sungguhan** — bukan derau yang wajar. Itu membuat penandanya berguna.

## Batas — yang TIDAK termasuk

⛔ Penandaan pasangan bergeser saat migrasi — tiket **09**.
⛔ Pembatalan — tiket **05**; ia **bukan** penghapusan baris.

## Cara mengujinya

Lewat seam `repository`. ⭐ Uji utama: permintaan menghapus baris rincian pada endorsemen harus
**ditolak**; dan baris baru yang ditambahkan harus mendapat nomor sesudah yang terbesar, bukan
mengisi lubang.

## Acceptance criteria

- [ ] **AC 10** — pada endorsemen, baris rincian **tidak dapat dihapus**
- [ ] **AC 11** — nomor urut lama **terbawa apa adanya** ke generasi berikutnya
- [ ] **AC 12** — baris baru mendapat nomor urut **maksimum + 1**
- [ ] **AC 13** — pasangan yang bergeser diperlakukan sebagai **anomali**, bukan keadaan wajar
