---
status: aktif
---

# 45: Daftar keadaan dan perpindahan sah berdiri, dan perpindahan di luar daftar ditolak

*Asal: `DAFTAR-PEKERJAAN.md` `P-01` (pecahan keadaan) dan `P-37` · `ADR-0055` · `ADR-0046` · `INV-20`…`INV-23`.*

**What to build:** Kolom `KEADAAN_SIKLUS_HIDUP` hanya menerima **delapan** nilai sah, dan hanya **tiga belas**
perpindahan yang diterima. Setiap perpindahan lain ditolak — termasuk **setiap** perpindahan keluar
dari `DISETUJUI`, `DITOLAK`, dan `DIBATALKAN`.

Ini **mesinnya**. Tiket `14` sudah membuat kolomnya; tiket ini memberi kolom itu **arti**.

**Persyaratan:** `ADR-0055` (delapan keadaan, tiga belas perpindahan) · `ADR-0046` (satu keadaan, himpunan tertutup) · `INV-20`…`INV-23` · `ADR-0038` (himpunan tertutup **boleh** `CHECK`, sebab ia tidak bertambah tanpa mengubah arti)

**Tidak termasuk:** **Siapa boleh memicu perpindahan mana** — itu wewenang, bukan daftar (`ADR-0044`), dan ia
milik tiket `48`…`52`.
**Catatan persetujuan** yang menyertai tiap perpindahan — tiket `54`.
**Pembekuan terminal** — tiket `46`; daftar putih menolak *perpindahan*, tidak menjaga *nilai*.

**Jalur gagal:** Menyetel keadaan ke nilai di luar delapan -> **ditolak** · `DISETUJUI` -> `DRAFT` ->
**ditolak**, dan pesannya menyebut bahwa keadaan itu terminal · Melompati `DIAJUKAN` langsung ke
`DISETUJUI` -> ditolak.

**Uji:** **Negatif:** kedelapan keadaan diuji terhadap setiap perpindahan yang **tidak** ada di daftar.
**POSITIF — dan ia yang membuat tiket ini tidak boleh dipecah:** **ketiga belas** perpindahan sah
diuji satu per satu dan **seluruhnya diterima**.

> Daftar putih yang kurang **lolos setiap uji negatif** — uji negatif hanya memeriksa bahwa yang
> terlarang ditolak. Daftar yang baru memuat lima perpindahan menolak delapan yang sah, dan
> **satu-satunya yang menangkapnya adalah uji positif yang lengkap.**

**Menggantikan:** `TDA-08` — *rantai memendek lewat `RevisionState`, dan jejaknya dihapus.* Dan lebih
luas: sistem lama **tidak punya daftar perpindahan sama sekali**. Keadaan disetel lewat
`StatusAkseptasi` dari enam layar berbeda, termasuk satu yang menyetelnya ke nilai `"test"` —
nilai yang tidak punya padanan di antara keadaan sah mana pun (`ADR-0054`).

**Blocked by:** `14`

**Dasar:**
```
EVIDENCED(TreatyInSetValue@ekspor-2026-09 - StatusAkseptasi="test" di enam layar)
        EVIDENCED(Akseptasi_DT@ekspor-2026-09 - cabang 1.4 pyDisabled, jalur kelompok tak terjangkau)
        DECIDED(ADR-0055, ADR-0046, ADR-0038)
```

- [ ] kolom keadaan menerima **tepat delapan** nilai, dan menolak yang lain
- [ ] **tiga belas** perpindahan diterima; **uji positif ketiga belasnya lulus satu per satu**
- [ ] setiap perpindahan keluar dari `DISETUJUI`, `DITOLAK`, `DIBATALKAN` ditolak
- [ ] pesan penolakan menyebut **keadaan asal dan tujuan**, bukan hanya "tidak sah"
