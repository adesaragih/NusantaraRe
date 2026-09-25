---
status: selesai
---

# 27: Domain kolom keadaan

> **SELESAI 2026-09-18.** RANCANGAN selesai — CHECK sudah ada di DDL. Pemasangannya ikut tiket tabelnya masing-masing; tidak ada pekerjaan tersendiri di sini.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-12` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Kolom keadaan yang domainnya kita tetapkan sendiri terikat `CHECK`; yang domainnya warisan dan belum lengkap **sengaja tidak terikat**.

`CK` `KEADAAN_BARIS` di **11 dari 11** tabel yang punya kolom itu; `CK` `KEADAAN_AKSEPTASI`.

**Tidak termasuk:** `STATUS_KLAIM` dan `KEPUTUSAN_KOMITE` — **tanpa `CHECK`, dan itu keputusan**. Domain tertutup yang isinya belum lengkap akan menolak nilai yang sah pada hari pertama modul Komite tersambung.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`22`~~ *(selesai)* — Koreksi bernilai tercatat menggantikan tambalan di dalam kode


**Dasar:** DECIDED(ADR-0019 lapis 2). EVIDENCED: sapuan S2 empat lapisan — `CNPStatusCase` hanya **2 nilai ditulis**, **nol rule menguji**, dan `"CLAIM ACCEPTED"`/`"CLAIM REJECTED"` **nihil di 279 berkas**.

- [ ] Nilai `KEADAAN_BARIS` di luar `LENGKAP` · `MENUNGGU_KURS` · `GAGAL_URAI` **ditolak**, di kesebelas tabel.
- [ ] `STATUS_KLAIM` bernilai apa pun **diterima**, dan alasannya tercatat.

**Ketidakpastian:** Kelengkapan enum `CNPStatusCase` hanya dapat dibuktikan dengan membuka folder `Komite Claim Non Prop` — **belum diizinkan**.
