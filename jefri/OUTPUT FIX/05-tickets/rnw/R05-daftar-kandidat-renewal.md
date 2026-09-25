# R05: Daftar kandidat renewal dan portalnya

**What to build:** Underwriter membuka daftar polis yang mendekati jatuh tempo dan memilih mana yang
akan diperpanjang — tanpa perlu tahu nomor polisnya lebih dulu. Daftar yang sama muncul di portal.

`[terverifikasi]` `ReportDefinition\RenewalList_RD` — **11 kolom**:

```
.pyID · .pzInsKey · .pxCreateDateTime · .pxCreateOperator
.Quotation.OldPolicyNo · .Quotation.InsuredName · .Quotation.MarketingName
.OfferFacIn.PolicyData.StartDateTime · .OfferFacIn.PolicyData.EndDateTime
.NBStatus · .NBStatusNew
```

Berkunci **`OldPolicyNo`** dan tanggal berakhirnya polis — itulah yang membuatnya daftar *kandidat*,
bukan sekadar daftar kasus.

`[terverifikasi]` `Section\SFAPortal_Renewal` mengikat 9 properti dengan bentuk yang sama.

### ⚠️ `Work-Renewal` adalah irisan pelaporan, bukan kelas kerja

`[terverifikasi]` `RenewalList_RD` berkelas `ASM-FW-GISFW-Work-Renewal`, **tetapi** `pyWorkClass` pada
flow renewal adalah `ASM-FW-GISFW-Work` — **sama dengan NB** — dan `Work-Renewal` tidak dipakai rule
lain mana pun.

⛔ **Jangan membuat kelas kerja atau antrean tersendiri untuk renewal.** Kasus renewal hidup di kelas
kerja yang sama dengan NB; `Work-Renewal` hanya cara melaporkannya.

**Blocked by:** R01

**Status:** ready-for-agent

- [ ] Daftar kandidat menampilkan kesebelas kolom, berkunci nomor polis lama dan tanggal berakhir
- [ ] Portal menampilkan daftar yang sama
- [ ] Memilih satu baris mengarah ke alur masuk renewal (R01) dengan nomor polis terisi
- [ ] **Tidak ada kelas kerja atau antrean baru** untuk renewal — kasus tetap di kelas kerja NB
- [ ] Penyaringan status dan kelompok tim direproduksi sesuai perilaku lama
- [ ] ⚠️ Kolom nama tertanggung dan nama marketing **ditampilkan sesuai perilaku lama**, tetapi **nilainya tidak pernah masuk fixture, test, atau artefak mana pun**
