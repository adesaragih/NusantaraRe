---
status: selesai
---

# 20: Objek, kronologi, dan dokumen klaim


> **SELESAI 19 September 2026 — tanpa perubahan struktur. Ketiga kriteria sudah terpenuhi bentuk yang ada.**
>
> Ditulis begini supaya penutupan ini dapat diperiksa, bukan dipercaya:
>
> | Kriteria | Yang menegakkannya |
> |---|---|
> | URL terisi tanpa kedaluwarsa ditolak, dan sebaliknya | `CK_DOKUMEN_KLAIM_1` — **dua arah, dan di sini dua arah memang benar**: URL tanpa masa berlaku tidak dapat dipakai, masa berlaku tanpa URL tidak menyatakan apa pun. Berbeda dari pasangan IDR–kurs, yang satu arah (tiket `16`) |
> | Dokumen tanpa URL sama sekali diterima | kedua kolom nullable; ADR-0027 — berkasnya di penyimpanan awan, URL diterbitkan saat diperlukan |
> | Kronologi tanpa jenis tindakan ditolak; tanpa keterangan diterima | `JENIS_TINDAKAN NOT NULL`, `KETERANGAN` nullable |
>
> `JENIS_TINDAKAN` **tanpa `CHECK` domain**, dan itu sejalan ADR-0009: yang terstruktur adalah **keputusannya**, dan jenis tindakan baru muncul seiring proses tanpa mengubah skema. `KETERANGAN` **tidak pernah dibaca mesin** — itu pokok ADR-0009, dan sebabnya terbaca di sistem lama: `@contains(.CommentSuggest,"Accepted by Himawan")` membuat satu salah ketik mengubah jalur (E7).
>
> **`ObjectList` tetap lubang tercatat, bukan tabel.** Apa yang membedakannya dari `InterestList` **TIDAK DITEMUKAN DI XML**, dan tabel yang dibuat atas dasar tebakan lebih buruk daripada lubang yang tercatat — `SPEC-MODEL-DATA.md` bagian 10.

*Asal: `T-09` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Objek pertanggungan, kronologi, dan rujukan dokumen tersimpan; URL tanpa masa berlaku ditolak.

`OBJEK_PERTANGGUNGAN`, `KRONOLOGI_KLAIM`, `DOKUMEN_KLAIM` beserta `CK_DOKUMEN_KLAIM_1`.

**Tidak termasuk:** `ObjectList` — **lubang tercatat**, bukan tabel. Apa yang membedakannya dari `InterestList` **TIDAK DITEMUKAN DI XML**.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-0027): berkas di Google Cloud Storage, yang tersimpan hanya metadata dan URL beserta `EXPDATE`. DECIDED(ADR-0009): kronologi terstruktur, keterangan tidak pernah dibaca mesin.

- [x] Dokumen dengan URL terisi tetapi tanpa kedaluwarsa **ditolak**, dan sebaliknya — `CK_DOKUMEN_KLAIM_1`.
- [x] Dokumen tanpa URL sama sekali **diterima** — URL belum pernah diterbitkan.
- [x] Kronologi tanpa jenis tindakan **ditolak**; tanpa keterangan **diterima**.

**Ketidakpastian:** Tidak ada REQ. Lubang `ObjectList` dicatat di `SPEC-MODEL-DATA.md` bagian 10.
