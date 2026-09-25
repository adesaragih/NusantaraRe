# R07: Konversi produksi kasus renewal

**What to build:** Kasus renewal yang disetujui dikonversi ke produksi **lewat jalur yang sama persis
dengan New Business** — bukan jalur kedua yang bisa menyimpang.

`[terverifikasi]` Dasarnya kuat: `Activity\SaveJsonPolicyFacIn_Act` — activity simpan produksi —
**identik byte-per-byte** antara NB dan RNW (371.819 byte, SHA-256 sama). Renewal **tidak memerlukan
penghasil nomor polis tersendiri**.

`[terverifikasi]` `Activity\serviceInsertArasapasRNW_act` — kelas `ASM-FW-GISFW-Work`, **10 langkah**:

```
1 Call serviceInsertArasapas_act    6 Connect-REST
2 Property-Set                      7 Property-Set
3 RDB-List                          8 Call InsertLogServiceProd
4 Property-Set                      9 Page-Remove
5 Call …Int-M_LINK_SERVICE.GetLinkService   10 RDB-List
```

Polanya: **endpoint diambil dari tabel layanan → panggilan REST → catat log layanan**. Ini menguatkan
`CLAUDE.md` §4.4 langsung dari korpus — endpoint memang datang dari tabel, bukan literal.

## ⛔ Blocker eksternal: jalur produksi NB belum ditiketkan

Tiket ini **ditulis penuh sekarang, tetapi belum dapat dieksekusi.** Jalur simpan produksi berada di
**Out of Scope butir 3 spec NB**, menunggu **struktur tabel flat** (keputusan work owner) dan
**`ALL_SOURCE`** untuk prosedur yang dilewati seluruh tulisan produksi.

⚠️ Celah **K-004** juga tetap terbuka: rule konversi asli kelas `Work` (`serviceInsertArasapas_act`
tanpa sufiks) **tidak ada di korpus mana pun**. Yang tersedia hanya dua saudara yang memperlihatkan
polanya. Arah rancangan mengikuti pola NB (K-035); rule aslinya belum pernah terbaca.

**Blocked by:** R01 · ⛔ **EKSTERNAL — jalur produksi NB** (Out of Scope butir 3 spec NB: tabel flat +
`ALL_SOURCE`; belum ada tiketnya)

**Status:** ready-for-agent

- [ ] Konversi renewal memakai **jalur simpan yang sama** dengan NB — satu implementasi, bukan dua
- [ ] Endpoint diambil dari tabel layanan saat runtime; **tidak ada endpoint literal di kode**
- [ ] Panggilan servis dicatat ke log layanan sesuai perilaku lama
- [ ] Urutan kesepuluh langkah dipertahankan
- [ ] Komentar mencatat bahwa **rule kelas `Work` aslinya tidak ada di korpus** (K-004), dan bahwa rancangan ini mengikuti pola dari dua saudara yang terbaca
- [ ] ⛔ Tiket **tidak dinyatakan selesai** sebelum jalur produksi NB tersedia — blocker eksternalnya dicatat, bukan disiasati
