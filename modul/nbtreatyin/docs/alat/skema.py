"""Bangkitkan migrasi 320-328 dan docs/STRUKTUR-TABEL-NB-TREATY-IN.md dari katalog Go.

Sumber tunggal: backend/models/katalog.go (Kolom properti -> kolom -> golongan). Uji
`backend/models/katalog_test.go` menagih bahwa DDL dan katalog tetap sepakat.

    python skema.py

Aturan penjaga repo yang dipatuhi (inti/backend/penjaga, claimlife):
  - NUMBER hanya (38,8) / (5) / (10); penanda VARCHAR2, bukan NUMBER(1)
  - CREATE di kolom 0, satu kolom per baris, `)` sendiri di barisnya, pemisah `/`
  - {skema} di setiap pernyataan, nol COMMIT, nol kata "total_*" dan nol kata
    persetujuan berbahasa Inggris yang dilarang penjaga skema
  - LF, tanpa BOM; pengenal <= 30 byte
"""
import os
import re

DIR = os.path.dirname(os.path.abspath(__file__))
MODUL = os.path.normpath(os.path.join(DIR, "..", ".."))
KATALOG = os.path.join(MODUL, "backend", "models", "katalog.go")
MIGRASI = os.path.join(MODUL, "backend", "migrations")
STRUKTUR = os.path.join(MODUL, "docs", "STRUKTUR-TABEL-NB-TREATY-IN.md")

GOL = {"kTeks": "teks", "kKode": "kode", "kPenanda": "penanda", "kUang": "uang", "kPersen": "persen",
       "kTgl": "tanggal", "kTglWaktu": "tanggal-waktu", "kCacah": "cacah"}


def baca_katalog():
    src = open(KATALOG, encoding="utf-8").read()
    tabel = {}
    for m in re.finditer(r'var (\w+) = Tabel\{Nama: "(\w+)"(?:, Daftar: ([^,]+))?, Kolom: \[\]Kolom\{(.*?)\n\}\}', src, re.S):
        kol = []
        for k in re.finditer(r'(k\w+)\((?:pt\+)?"([^"]+)", "(\w+)"(?:, (\d+))?\)', m.group(4)):
            fn, prop, nama, n = k.groups()
            kol.append({"properti": prop, "kolom": nama, "gol": GOL[fn],
                        "panjang": int(n) if n else (16 if fn == "kPenanda" else 0)})
        tabel[m.group(2)] = {"var": m.group(1), "daftar": (m.group(3) or "").strip(), "kolom": kol}
    return tabel


def tipe_ddl(k):
    g = k["gol"]
    if g in ("uang", "persen"):
        return "NUMBER(38,8)"
    if g in ("tanggal", "tanggal-waktu"):
        return "DATE"
    if g == "cacah":
        return "NUMBER(10)"
    return f"VARCHAR2({k['panjang']})"


def tipe_struktur(k):
    g = k["gol"]
    if g in ("uang", "persen"):
        return "angka desimal"
    if g in ("tanggal", "tanggal-waktu"):
        return "DATE"
    if g == "cacah":
        return "bilangan bulat"
    return "teks"


def blok(nama, kolom_kunci, kolom, constraint):
    lebar = max(len(c[0]) for c in kolom_kunci + [(k["kolom"], "") for k in kolom]) + 2
    baris = [f"CREATE TABLE {{skema}}.{nama} ("]
    for n, tp in kolom_kunci:
        baris.append(f"  {n.ljust(lebar)}{tp},")
    for k in kolom:
        baris.append(f"  {k['kolom'].ljust(lebar)}{tipe_ddl(k)},")
    for i, c in enumerate(constraint):
        baris.append("  " + c + ("," if i < len(constraint) - 1 else ""))
    baris.append(")")
    return "\n".join(baris)


def tulis(nama_berkas, kepala, pernyataan, turun):
    os.makedirs(MIGRASI, exist_ok=True)
    isi = "\n".join("-- " + x if x else "--" for x in kepala) + "\n" + "\n/\n".join(pernyataan) + "\n/\n"
    open(os.path.join(MIGRASI, nama_berkas + ".sql"), "w", encoding="utf-8", newline="\n").write(isi)
    isi_turun = f"-- Jalur mundur {nama_berkas[:3]}.\n" + "\n/\n".join(turun) + "\n/\n"
    open(os.path.join(MIGRASI, nama_berkas + "_down.sql"), "w", encoding="utf-8", newline="\n").write(isi_turun)


def main():
    t = baca_katalog()
    g = t["T_GENERAL_POLIS"]

    # ------------------------------------------------------------ 320
    kunci = [("ID", "VARCHAR2(32) NOT NULL"), ("NOPOLIS", "VARCHAR2(64)"), ("PRODKE", "NUMBER(10) DEFAULT 0 NOT NULL"),
             ("NOENDORS", "VARCHAR2(64)"), ("OLD_POLIS_ID", "VARCHAR2(32)"), ("IDPEGA", "VARCHAR2(128)"),
             ("TGL_INPUT", "DATE"), ("USERNAME", "VARCHAR2(64)"), ("TGL_TUTUP", "DATE")]
    tulis("320_t_general_polis", [
        "320 - T_GENERAL_POLIS: satu baris per GENERASI polis treaty inward (tiket 16).",
        "",
        "Kunci utama BERSAMA T_WORK_POLIS (tabel kasus lintas-lini milik premiumlistlife,",
        "migrasi 050/059): ID adalah ID baris T_WORK_POLIS, tanpa kolom kunci tamu",
        "tersendiri (spec-penyimpanan ID-7, AC 5).",
        "",
        "Kunci alami (NOPOLIS, PRODKE) unik - ditegakkan BASIS DATA (ID-8, AC 1).",
        "NOPOLIS kosong selama realisasi belum bernomor, jadi indeks uniknya hanya",
        "memuat baris bernomor (indeks berfungsi CASE): dua draf tanpa nomor tidak bentrok.",
        "PRODKE bilangan bulat lebar (KEPUTUSAN-RONDE-12 butir 1 dan 8); NB selalu 0.",
        "OLD_POLIS_ID menunjuk generasi sebelumnya - unik, kosong di NB (ID-9, AC 3, 4).",
        "TGL_TUTUP terisi = generasi tertutup, tidak boleh disunting (ID-10, AC 6) -",
        "ditegakkan services.",
        "",
        "Kolom lain DIBANGKITKAN dari backend/models/katalog.go (docs/alat/skema.py).",
        "Uang dan persen NUMBER(38,8) (KEPUTUSAN 23-09-2026 sore), tanggal DATE (P32),",
        "kode dan penanda teks (ID-16, ID-17). Nol COMMIT.",
    ], [
        blok("T_GENERAL_POLIS", kunci, g["kolom"], [
            "CONSTRAINT PK_GENERAL_POLIS PRIMARY KEY (ID)",
            "CONSTRAINT FK_GENERAL_POLIS_WORK FOREIGN KEY (ID) REFERENCES {skema}.T_WORK_POLIS (ID)",
            "CONSTRAINT FK_GENERAL_POLIS_OLD FOREIGN KEY (OLD_POLIS_ID) REFERENCES {skema}.T_WORK_POLIS (ID)",
            "CONSTRAINT UQ_GENERAL_POLIS_OLD UNIQUE (OLD_POLIS_ID)",
        ]),
        "CREATE UNIQUE INDEX {skema}.UQ_GENERAL_POLIS_NOPOLIS ON {skema}.T_GENERAL_POLIS "
        "(CASE WHEN NOPOLIS IS NOT NULL THEN NOPOLIS END, CASE WHEN NOPOLIS IS NOT NULL THEN PRODKE END)",
    ], ["DROP TABLE {skema}.T_GENERAL_POLIS CASCADE CONSTRAINTS"])

    # ------------------------------------------------------------ 321 quotation (1:1)
    q = t["T_POLIS_QUOTATION"]
    tulis("321_t_polis_quotation", [
        "321 - T_POLIS_QUOTATION: halaman Quotation polis, 1:1 dengan T_GENERAL_POLIS",
        "(spec-penyimpanan ID-23; tiket 19 AC 27). Kunci utama = POLIS_ID.",
    ], [blok("T_POLIS_QUOTATION", [("POLIS_ID", "VARCHAR2(32) NOT NULL")], q["kolom"], [
        "CONSTRAINT PK_POLIS_QUOTATION PRIMARY KEY (POLIS_ID)",
        "CONSTRAINT FK_POLIS_QUOTATION_POLIS FOREIGN KEY (POLIS_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID)",
    ])], ["DROP TABLE {skema}.T_POLIS_QUOTATION CASCADE CONSTRAINTS"])

    # ------------------------------------------------------------ anak 1:N
    def anak(no, nama, induk_kol, induk_tabel, pk, uq, fk, kepala):
        k = t[nama]
        tulis(f"{no}_{nama.lower()}", kepala, [blok(nama, [("ID", "VARCHAR2(32) NOT NULL"),
                                                          (induk_kol, "VARCHAR2(32) NOT NULL"),
                                                          ("NOURUT", "NUMBER(5) NOT NULL")], k["kolom"], [
            f"CONSTRAINT {pk} PRIMARY KEY (ID)",
            f"CONSTRAINT {fk} FOREIGN KEY ({induk_kol}) REFERENCES {{skema}}.{induk_tabel} (ID)",
            f"CONSTRAINT {uq} UNIQUE ({induk_kol}, NOURUT)",
        ])], [f"DROP TABLE {{skema}}.{nama} CASCADE CONSTRAINTS"])

    nourut = ["NOURUT = nomor urut baris di dalam induknya, unik per induk (ID-11, AC 8, 10);",
              "kunci pasangan antar generasi, bukan kunci dagang (ID-13, AC 11)."]
    anak("322", "T_POLIS_CEDING", "POLIS_ID", "T_GENERAL_POLIS", "PK_POLIS_CEDING", "UQ_POLIS_CEDING_NOURUT",
         "FK_POLIS_CEDING_POLIS", ["322 - T_POLIS_CEDING <- QuotationData.CedingCoList (ID-24; tiket 19 AC 28-30)."] + nourut)
    anak("323", "T_POLIS_INSTALMENT", "POLIS_ID", "T_GENERAL_POLIS", "PK_POLIS_INSTALMENT", "UQ_POLIS_INSTALMENT_NOURUT",
         "FK_POLIS_INSTALMENT_POLIS", ["323 - T_POLIS_INSTALMENT <- PolicyTreatyIn.ListInstallment (ID-26)."] + nourut)
    anak("324", "T_POLIS_INSTALMENT_DETAIL", "INSTALMENT_ID", "T_POLIS_INSTALMENT", "PK_POLIS_INSTALMENT_DETAIL",
         "UQ_POLIS_INST_DETAIL_NOURUT", "FK_POLIS_INST_DETAIL_INST",
         ["324 - T_POLIS_INSTALMENT_DETAIL <- ListInstallment().InstallmentList - hanya",
          "non-proporsional (rancangan-tabel-datar §3.2; tiket 19 AC 31)."] + nourut)
    anak("325", "T_POLIS_SPREADING", "POLIS_ID", "T_GENERAL_POLIS", "PK_POLIS_SPREADING", "UQ_POLIS_SPREADING_NOURUT",
         "FK_POLIS_SPREADING_POLIS", ["325 - T_POLIS_SPREADING <- PolicyTreatyIn.SpreadingRiskList (ID-28).",
                                      "Total share dan nilai TIDAK disimpan - turunan baris ini."] + nourut)
    anak("326", "T_POLIS_XOL", "POLIS_ID", "T_GENERAL_POLIS", "PK_POLIS_XOL", "UQ_POLIS_XOL_NOURUT",
         "FK_POLIS_XOL_POLIS", ["326 - T_POLIS_XOL <- PolicyTreatyIn.TreatyXOLList - induk XOL tanpa penanda",
                                "layer (ID-29); DEDUCTION di sini UANG (ID-30)."] + nourut)
    anak("327", "T_POLIS_XOL_LAYER", "XOL_ID", "T_POLIS_XOL", "PK_POLIS_XOL_LAYER", "UQ_POLIS_XOL_LAYER_NOURUT",
         "FK_POLIS_XOL_LAYER_XOL", ["327 - T_POLIS_XOL_LAYER <- TreatyXOLList().ValueList - pemegang penanda",
                                    "layer (ID-29; tiket 19 AC 32-35)."] + nourut)
    # ------------------------------------------------------------ STRUKTUR
    w = ["# Struktur Tabel — NB Treaty In", "",
         "Acuan bentuk tabel modul `nbtreatyin`. **Dibangkitkan** `docs/alat/skema.py` dari",
         "`backend/models/katalog.go` — jangan disunting tangan; sunting katalognya.", "",
         "Tipe ditulis sebagai kategori logis: teks · angka desimal · bilangan bulat · DATE.",
         "Uang dan persen **angka desimal** `NUMBER(38,8)`, tidak pernah float (ADR-0003).",
         "Golongan (uang / persen / kode / penanda / tanggal) ada di kolom *Golongan*.", "",
         "Tabel yang **dibaca, tidak dibuat** modul ini dideklarasikan di `MODUL.md`.", ""]

    def sect(nama, kunci_rows, kolom, cat):
        w.extend([f"## {nama}", "", cat, "", "| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |",
                  "| --- | --- | --- | --- | --- | --- |"])
        for r in kunci_rows:
            w.append(f"| `{r[0]}` | {r[1]} | {r[2]} | {r[3]} | {r[4]} | {r[5]} |")
        for k in kolom:
            w.append(f"| `{k['kolom']}` | {tipe_struktur(k)} | ya |  | {k['gol']} | `{k.get('properti', '')}` |")
        w.append("")

    sect("T_GENERAL_POLIS", [
        ("ID", "teks", "tidak", "PK, FK T_WORK_POLIS", "kode", "kasus"),
        ("NOPOLIS", "teks", "ya", "UQ (NOPOLIS, PRODKE)", "kode", "`PolicyTreatyIn.PolicyNo`"),
        ("PRODKE", "bilangan bulat", "tidak", "UQ (NOPOLIS, PRODKE)", "cacah", "generasi; NB = 0"),
        ("NOENDORS", "teks", "ya", "", "kode", "json_polis"),
        ("OLD_POLIS_ID", "teks", "ya", "UQ, FK T_WORK_POLIS", "kode", "generasi sebelumnya"),
        ("IDPEGA", "teks", "ya", "", "kode", "json_polis"),
        ("TGL_INPUT", "DATE", "ya", "", "tanggal-waktu", "json_polis"),
        ("USERNAME", "teks", "ya", "", "kode", "identitas akses login (P4)"),
        ("TGL_TUTUP", "DATE", "ya", "", "tanggal-waktu", "generasi tertutup"),
    ], g["kolom"], "Satu baris per generasi polis; kunci utama bersama `T_WORK_POLIS` (ID-7).")
    sect("T_POLIS_QUOTATION", [("POLIS_ID", "teks", "tidak", "PK, FK T_GENERAL_POLIS", "kode", "induk")],
         q["kolom"], "Halaman `Quotation` / `PolicyTreatyIn.QuotationData`, 1:1 (ID-23).")
    for nama, induk, cat in [
            ("T_POLIS_CEDING", ("POLIS_ID", "T_GENERAL_POLIS"), "← `QuotationData.CedingCoList` (ID-24)."),
            ("T_POLIS_INSTALMENT", ("POLIS_ID", "T_GENERAL_POLIS"), "← `PolicyTreatyIn.ListInstallment` (ID-26)."),
            ("T_POLIS_INSTALMENT_DETAIL", ("INSTALMENT_ID", "T_POLIS_INSTALMENT"), "← `ListInstallment().InstallmentList`, non-proporsional."),
            ("T_POLIS_SPREADING", ("POLIS_ID", "T_GENERAL_POLIS"), "← `PolicyTreatyIn.SpreadingRiskList` (ID-28)."),
            ("T_POLIS_XOL", ("POLIS_ID", "T_GENERAL_POLIS"), "← `PolicyTreatyIn.TreatyXOLList` (ID-29)."),
            ("T_POLIS_XOL_LAYER", ("XOL_ID", "T_POLIS_XOL"), "← `TreatyXOLList().ValueList` (ID-29).")]:
        sect(nama, [("ID", "teks", "tidak", "PK", "kode", "baris"),
                    (induk[0], "teks", "tidak", f"FK {induk[1]}, UQ ({induk[0]}, NOURUT)", "kode", "induk"),
                    ("NOURUT", "bilangan bulat", "tidak", f"UQ ({induk[0]}, NOURUT)", "cacah", "urutan baris (ID-11)")],
             t[nama]["kolom"], cat)
    # ------------------------------------------------------------ tabel warisan (tidak dibuat)
    w.extend(["## HISTORYAKSEPTASIPRODUCTION", "",
              "⛔ **Tabel WARISAN `POOLDATA` — tidak dibuat, tidak diubah strukturnya** (MODUL.md *Tabel warisan*;",
              "keputusan work owner K4 03-10-2026; spec-penyimpanan ID-31). Catatan `PolicyTreatyIn.SuggestList` ditulis",
              "ke sini (`repository/usulan.go`, pengganti `RDBList/InsertViewSuggest_SQL`) dan dibaca balik untuk layar.",
              "Tipe fisik milik tabel lama (belum dicek katalog Oracle — butir terbuka).", "",
              "| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |", "| --- | --- | --- | --- | --- | --- |"])
    for kol, sumber in [("IDPEGA", "`pyWorkPage.pzInsKey`"), ("TYPE_POLIS", "`@replaceAll(pyWorkIDPrefix,\"-\",\"\")` = `NB`"),
                        ("NOURUT", "berikutnya per IDPEGA (XML `.pxListSubscript`)"), ("POSISI", "`\"Policy\"`"),
                        ("PIC", "`.OperatorName` — nama tampilan (P33)"), ("TGL_INP", "`.Date`"),
                        ("DIV", "`OperatorID.pyOrgDivision` — NULL, tanpa sumber (butir terbuka)"),
                        ("TYPE", "`Quotation.BusinessFac` = `T`"), ("PUTARAN", "`\"2\"`"),
                        ("APPROVAL", "`.IsApproved` 1 = Accept, 0 = Reject"), ("KETERANGAN", "`substr(.Suggest, 0, 3990)`"),
                        ("AKSES_LOGIN", "`OperatorID.pyUserIdentifier` — identitas login (P4)"),
                        ("B2B", "`OfferFacIn.IsB2B` — NULL di NB"), ("BUSINESS_CODE", "`Quotation.BusinessCode`"),
                        ("PERCENT_RNM", "`OfferFacIn.PercentShare` — NULL di NB")]:
        w.append(f"| `{kol}` | warisan | ya |  | — | {sumber} |")
    w.append("")
    open(STRUKTUR, "w", encoding="utf-8", newline="\n").write("\n".join(w).rstrip() + "\n")
    print("ok:", ", ".join(sorted(t)))


if __name__ == "__main__":
    main()
