# -*- coding: utf-8 -*-
"""Bangkitkan ERD-ORACLE.xlsx - ERD kotak-entitas untuk TABEL Oracle saja.

Dua sistem, dua kelompok sheet, tidak dicampur:
  - sistem lama : 32 TABLE Oracle yang benar-benar berdiri (pengetahuan/SCHEMA-ACTUAL.csv)
  - skema baru  : 22 TABLE usulan KLAIMNP (ddl-usulan/*.sql)

VIEW, PROCEDURE, FUNCTION, dan halaman clipboard Pega TIDAK digambar di sini.
Sumber dibaca langsung; berkas turunan tidak pernah menjadi sumber.
"""
import os, re, csv, html, collections

from openpyxl import Workbook
from openpyxl.styles import Font, PatternFill, Alignment, Border, Side
from openpyxl.utils import get_column_letter

AKAR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
XML = r"D:\XML_NURE\Claim Non Prop"
SKEMA = os.path.join(AKAR, "pengetahuan", "SCHEMA-ACTUAL.csv")
USULAN = os.path.join(AKAR, "2-to-spec", "ddl-usulan")
OUT = os.path.join(AKAR, "4-erd-dan-tabel-datar", "ERD-ORACLE.xlsx")

TANGGAL = "19 September 2026"

# ----------------------------------------------------------------- gaya
FONT = "Arial"
H1 = Font(name=FONT, size=14, bold=True)
H2 = Font(name=FONT, size=11, bold=True)
TXT = Font(name=FONT, size=10)
TXT_K = Font(name=FONT, size=10, bold=True)
KECIL = Font(name=FONT, size=9)
KECIL_M = Font(name=FONT, size=9, italic=True, color="666666")
PUTIH = Font(name=FONT, size=10, bold=True, color="FFFFFF")

TIPIS = Side(style="thin", color="B0B0B0")
TEBAL = Side(style="medium", color="404040")
KOTAK_LUAR = Border(left=TEBAL, right=TEBAL, top=TEBAL, bottom=TEBAL)
SEL = Border(left=TIPIS, right=TIPIS, top=TIPIS, bottom=TIPIS)

# warna kepala kotak per kelompok domain
WARNA = {
    "Klaim":      "1F4E79",
    "Treaty":     "2E6C3E",
    "Polis":      "7A4E1D",
    "Pihak":      "6B3A7A",
    "Dokumen":    "1F6F73",
    "Integrasi":  "8A2B2B",
    "Pega":       "555555",
    "Skema baru": "1F4E79",
}
WARNA_ISI = {
    "Klaim":      "DCE6F1",
    "Treaty":     "DFEDE1",
    "Polis":      "F5E7D6",
    "Pihak":      "EDE1F2",
    "Dokumen":    "D9EEEF",
    "Integrasi":  "F6DEDE",
    "Pega":       "E8E8E8",
    "Skema baru": "DCE6F1",
}

KELOMPOK = {
    "OS_AKSEPTASI_KLAIM": "Klaim", "CLAIMREJECTED": "Klaim", "CLAIMXOL2": "Klaim",
    "JSON_KLAIM": "Klaim", "CATASTROPHE": "Klaim",
    "TREATYINDETAIL": "Treaty", "TREATYINDETAILEDM": "Treaty", "TREATYINPRODUCTION": "Treaty",
    "TREATYBUSINESS": "Treaty", "TREATYCONTRACT": "Treaty", "TREATYGROUP": "Treaty",
    "PROPORTIONALARRG": "Treaty", "REINSURANCETYPE": "Treaty", "M_TREATY_IN": "Treaty",
    "M_TREATY_IN_EDM": "Treaty", "M_TREATY_OUT": "Treaty", "M_TREATY_OUT_DETAIL": "Treaty",
    "JSON_POLIS": "Polis",
    "M_CLIENT": "Pihak", "AGENT": "Pihak", "BANKACCOUNT": "Pihak", "LST_BANK_GROUP": "Pihak",
    "MARKETINGOFFICER": "Pihak", "ADJUSTERCONSULTANT": "Pihak", "BUSINESS": "Pihak",
    "RW": "Pihak", "KODE_PRODUKSI": "Pihak", "EMAILKOMITE": "Pihak",
    "T_STORAGE_IMAGE": "Dokumen", "T_FOLDER_IMAGE": "Dokumen",
    "DIRECTTOKASIR_LOG": "Integrasi",
    "PC_ASM_FW_GCNMFW_WORK": "Pega",
}

ARTI = {
    "OS_AKSEPTASI_KLAIM": "Nilai akseptasi klaim mendarat di sini; isinya JSON di kolom DATA_JSON",
    "CLAIMREJECTED": "Klaim yang ditolak",
    "CLAIMXOL2": "Kaitan case ke alokasi XOL",
    "JSON_KLAIM": "Muatan klaim utuh sebagai JSON",
    "JSON_POLIS": "Muatan polis sebagai JSON; jembatan nomor polis ke ID Pega",
    "CATASTROPHE": "Daftar kejadian katastrofe",
    "TREATYINDETAIL": "Rincian treaty masuk",
    "TREATYINDETAILEDM": "Rincian treaty masuk, jalur EDM",
    "TREATYINPRODUCTION": "Treaty masuk yang sudah berproduksi",
    "TREATYBUSINESS": "Kaitan treaty ke lini usaha",
    "TREATYCONTRACT": "Kontrak treaty",
    "TREATYGROUP": "Pengelompokan treaty",
    "PROPORTIONALARRG": "Susunan proporsional",
    "REINSURANCETYPE": "Jenis reasuransi",
    "M_TREATY_IN": "Master treaty masuk (JSON)",
    "M_TREATY_IN_EDM": "Master treaty masuk jalur EDM (JSON)",
    "M_TREATY_OUT": "Master treaty keluar (JSON)",
    "M_TREATY_OUT_DETAIL": "Rincian treaty keluar (JSON)",
    "M_CLIENT": "Master klien",
    "AGENT": "Agen / perantara",
    "BANKACCOUNT": "Rekening bank penerima pembayaran",
    "LST_BANK_GROUP": "Kelompok bank",
    "MARKETINGOFFICER": "Petugas pemasaran",
    "ADJUSTERCONSULTANT": "Adjuster dan konsultan",
    "BUSINESS": "Lini usaha",
    "RW": "Wilayah administratif",
    "KODE_PRODUKSI": "Kode produksi per lini",
    "EMAILKOMITE": "Alamat surel anggota komite",
    "T_STORAGE_IMAGE": "Dokumen tersimpan",
    "T_FOLDER_IMAGE": "Folder penyimpanan dokumen",
    "DIRECTTOKASIR_LOG": "Catatan kiriman ke kasir",
    "PC_ASM_FW_GCNMFW_WORK": "Tabel work Pega - TIDAK DIMIGRASI",
}


# ============================================================ 1. sistem lama
def baca_tabel_lama():
    rows = list(csv.DictReader(open(SKEMA, encoding="utf-8-sig"), delimiter=";"))
    tab = collections.OrderedDict()
    for r in rows:
        if r["object_type"] != "TABLE":
            continue
        kunci = r["object_name"]
        tab.setdefault(kunci, {"owner": r["owner"], "kolom": []})
        tipe = r["data_type"]
        if r["data_precision"]:
            tipe += "(%s%s)" % (r["data_precision"],
                                "," + r["data_scale"] if r["data_scale"] else "")
        elif r["data_length"]:
            tipe += "(%s)" % r["data_length"]
        tab[kunci]["kolom"].append({
            "no": r["column_id"], "nama": r["column_name"], "tipe": tipe,
            "null": r["nullable"], "default": r["data_default"],
            "pk": r["is_pk"], "unik": r["is_unique"],
        })
    return tab


def baca_relasi_lama():
    """Relasi sistem lama tidak ada di DDL - nol FOREIGN KEY. Ia hanya terbaca
    dari SQL di ekspor XML: comma-join dan subquery berkorelasi."""
    def deun(s):
        for _ in range(4):
            n = html.unescape(s)
            if n == s:
                break
            s = n
        return s

    blok = []
    for dp, dn, fn in os.walk(XML):
        for f in sorted(fn):
            if not f.lower().endswith(".xml"):
                continue
            t = deun(open(os.path.join(dp, f), encoding="utf-8", errors="replace").read())
            for m in re.findall(r"<pyBrowseSQL>(.*?)</pyBrowseSQL>", t, re.S):
                if m.strip():
                    blok.append((os.path.relpath(os.path.join(dp, f), XML).replace("\\", "/"),
                                 " ".join(m.split())))

    EQ = re.compile(r"\b(\w+)\.(\w+)\s*=\s*(\w+)\.(\w+)\b")
    SUB = re.compile(r"\b(\w+)\s*(?:=|IN)\s*\(\s*SELECT\s+(\w+)\s+FROM\s+(?:(\w+)\.)?(\w+)", re.I)
    out = {}
    for berkas, s in blok:
        alias, dipakai = {}, set()
        for potongan in re.split(r"\bFROM\b", s, flags=re.I)[1:]:
            daftar = re.split(r"\bWHERE\b|\bORDER\b|\bGROUP\b|\bCONNECT\b|\bSELECT\b|\)",
                              potongan, flags=re.I)[0]
            for bagian in daftar.split(","):
                tok = [x.strip('"') for x in bagian.strip().split() if x.strip()]
                if not tok or "(" in bagian:
                    continue
                nama = tok[0].split(".")[-1].upper()
                if not re.fullmatch(r"\w+", nama) or nama in ("DUAL", "SELECT"):
                    continue
                dipakai.add(nama)
                if len(tok) > 1 and re.fullmatch(r"\w+", tok[1]):
                    alias[tok[1].upper()] = nama
                alias.setdefault(nama, nama)
        for me in EQ.finditer(s):
            a1, k1, a2, k2 = (x.upper() for x in me.groups())
            t1, t2 = alias.get(a1), alias.get(a2)
            if t1 and t2 and t1 != t2:
                out.setdefault((t1, k1, t2, k2, "comma-join"), set()).add(berkas)
        for ms in SUB.finditer(s):
            kl, kd, own, td = ms.groups()
            td = td.upper()
            luar = sorted(dipakai - {td})
            if luar:
                out.setdefault((luar[0], kl.upper(), td, kd.upper(), "subquery"),
                               set()).add(berkas)
    return out


# ============================================================ 2. skema baru
def baca_tabel_baru():
    tab = collections.OrderedDict()
    fk, ck, uq, pk_tbl = [], [], [], {}
    CT = re.compile(r"CREATE\s+TABLE\s+KLAIMNP\.(\w+)\s*\((.*?)\n\)\s*;", re.S | re.I)
    for f in sorted(os.listdir(USULAN)):
        if not f.lower().endswith(".sql"):
            continue
        teks = open(os.path.join(USULAN, f), encoding="utf-8", errors="replace").read()
        bersih = "\n".join(b for b in teks.splitlines() if not b.strip().startswith("--"))
        for m in CT.finditer(bersih):
            nama, badan = m.group(1).upper(), m.group(2)
            kolom = []
            for baris in badan.split("\n"):
                b = baris.strip().rstrip(",").strip()
                if not b or b.upper().startswith(("CONSTRAINT", "PRIMARY KEY", "UNIQUE",
                                                  "FOREIGN KEY", "CHECK")):
                    continue
                mm = re.match(r"(\w+)\s+(.+)$", b)
                if not mm:
                    continue
                sisa = mm.group(2)
                kolom.append({
                    "no": len(kolom) + 1, "nama": mm.group(1).upper(),
                    "tipe": re.split(r"\s+(?:DEFAULT|NOT|NULL|GENERATED|CONSTRAINT)\b",
                                     sisa, flags=re.I)[0].strip(),
                    "null": "N" if re.search(r"NOT\s+NULL", sisa, re.I) else "Y",
                    "default": (re.search(r"DEFAULT\s+([^\s,]+)", sisa, re.I).group(1)
                                if re.search(r"DEFAULT\s+", sisa, re.I) else ""),
                    "turunan": "Y" if re.search(r"GENERATED\s+ALWAYS", sisa, re.I) else "",
                })
            tab[nama] = {"owner": "KLAIMNP", "kolom": kolom, "berkas": f}
        for m in re.finditer(r"ALTER\s+TABLE\s+KLAIMNP\.(\w+)\s+ADD\s+CONSTRAINT\s+(\w+)\s+"
                             r"FOREIGN\s+KEY\s*\(([^)]+)\)\s*REFERENCES\s+KLAIMNP\.(\w+)\s*\(([^)]+)\)",
                             bersih, re.I):
            fk.append({"tabel": m.group(1).upper(), "constraint": m.group(2).upper(),
                       "kolom": m.group(3).strip().upper(), "ref_tabel": m.group(4).upper(),
                       "ref_kolom": m.group(5).strip().upper(), "berkas": f})
        for m in re.finditer(r"ALTER\s+TABLE\s+KLAIMNP\.(\w+)\s+ADD\s+CONSTRAINT\s+(\w+)\s+"
                             r"CHECK\s*\((.+?)\);", bersih, re.I | re.S):
            ck.append({"tabel": m.group(1).upper(), "constraint": m.group(2).upper(),
                       "isi": " ".join(m.group(3).split()), "berkas": f})
        for m in re.finditer(r"ALTER\s+TABLE\s+KLAIMNP\.(\w+)\s+ADD\s+CONSTRAINT\s+(\w+)\s+"
                             r"UNIQUE\s*\(([^)]+)\)", bersih, re.I):
            uq.append({"tabel": m.group(1).upper(), "constraint": m.group(2).upper(),
                       "kolom": m.group(3).strip().upper(), "berkas": f})
        for m in re.finditer(r"ALTER\s+TABLE\s+KLAIMNP\.(\w+)\s+ADD\s+CONSTRAINT\s+(\w+)\s+"
                             r"PRIMARY\s+KEY\s*\(([^)]+)\)", bersih, re.I):
            pk_tbl[m.group(1).upper()] = m.group(3).strip().upper()
    # PK yang ditulis inline di dalam CREATE TABLE
    for f in sorted(os.listdir(USULAN)):
        if not f.lower().endswith(".sql"):
            continue
        teks = open(os.path.join(USULAN, f), encoding="utf-8", errors="replace").read()
        bersih = "\n".join(b for b in teks.splitlines() if not b.strip().startswith("--"))
        for m in CT.finditer(bersih):
            nama, badan = m.group(1).upper(), m.group(2)
            mp = re.search(r"CONSTRAINT\s+\w+\s+PRIMARY\s+KEY\s*\(([^)]+)\)", badan, re.I)
            if mp:
                pk_tbl[nama] = mp.group(1).strip().upper()
    return tab, fk, ck, uq, pk_tbl


# ============================================================ 3. gambar kotak
def gambar_erd(ws, tabel, kelompok_of, arti_of, relasi_kolom, pk_of, judul, catatan):
    """Susun kotak entitas dalam beberapa jalur vertikal, tinggi diseimbangkan."""
    LEBAR = 4            # kolom isi per kotak
    JALUR = 4            # banyak jalur berdampingan
    GAP_KOL = 1
    GAP_BARIS = 2

    ws["A1"] = judul
    ws["A1"].font = H1
    ws["A2"] = catatan
    ws["A2"].font = KECIL_M
    baris_awal = 4

    # lebar kolom
    for j in range(JALUR):
        c0 = 1 + j * (LEBAR + GAP_KOL)
        ws.column_dimensions[get_column_letter(c0)].width = 5      # penanda kunci
        ws.column_dimensions[get_column_letter(c0 + 1)].width = 30  # nama kolom
        ws.column_dimensions[get_column_letter(c0 + 2)].width = 18  # tipe
        ws.column_dimensions[get_column_letter(c0 + 3)].width = 30  # relasi / catatan
        if j < JALUR - 1:
            ws.column_dimensions[get_column_letter(c0 + LEBAR)].width = 3

    tinggi = [baris_awal] * JALUR
    peta = {}
    for nama in sorted(tabel, key=lambda x: (-len(tabel[x]["kolom"]), x)):
        j = tinggi.index(min(tinggi))
        r = tinggi[j]
        c0 = 1 + j * (LEBAR + GAP_KOL)
        info = tabel[nama]
        grup = kelompok_of.get(nama, "Skema baru")
        warna_kepala = WARNA.get(grup, "1F4E79")
        warna_isi = WARNA_ISI.get(grup, "DCE6F1")

        # kepala: nama tabel
        ws.merge_cells(start_row=r, start_column=c0, end_row=r, end_column=c0 + LEBAR - 1)
        sel = ws.cell(row=r, column=c0)
        sel.value = "%s.%s" % (info["owner"], nama)
        sel.font = PUTIH
        sel.fill = PatternFill("solid", fgColor=warna_kepala)
        sel.alignment = Alignment(horizontal="left", vertical="center", indent=1)

        # baris arti + cacah kolom
        ws.merge_cells(start_row=r + 1, start_column=c0, end_row=r + 1, end_column=c0 + LEBAR - 1)
        sk = ws.cell(row=r + 1, column=c0)
        sk.value = "%d kolom  -  %s" % (len(info["kolom"]), arti_of.get(nama, ""))
        sk.font = KECIL_M
        sk.fill = PatternFill("solid", fgColor=warna_isi)
        sk.alignment = Alignment(horizontal="left", vertical="center", indent=1, wrap_text=False)

        for i, k in enumerate(info["kolom"]):
            rr = r + 2 + i
            penanda = []
            if k.get("pk") == "Y" or pk_of.get(nama) == k["nama"]:
                penanda.append("PK")
            if k.get("unik") == "Y":
                penanda.append("UQ")
            rel = relasi_kolom.get((nama, k["nama"]), "")
            if rel:
                penanda.append("FK" if k.get("fk") else "~")
            ws.cell(row=rr, column=c0, value=" ".join(penanda)).font = TXT_K
            cn = ws.cell(row=rr, column=c0 + 1, value=k["nama"])
            cn.font = TXT_K if penanda else TXT
            ws.cell(row=rr, column=c0 + 2, value=k["tipe"]).font = KECIL
            ws.cell(row=rr, column=c0 + 3, value=rel).font = KECIL_M
            if k.get("null") == "N":
                cn.font = Font(name=FONT, size=10, bold=True)
            for cc in range(c0, c0 + LEBAR):
                ws.cell(row=rr, column=cc).border = SEL

        tinggi_kotak = 2 + len(info["kolom"])
        # garis luar kotak
        for cc in range(c0, c0 + LEBAR):
            atas = ws.cell(row=r, column=cc)
            atas.border = Border(top=TEBAL, bottom=TIPIS,
                                 left=TEBAL if cc == c0 else TIPIS,
                                 right=TEBAL if cc == c0 + LEBAR - 1 else TIPIS)
            bawah = ws.cell(row=r + tinggi_kotak - 1, column=cc)
            b = bawah.border
            bawah.border = Border(left=TEBAL if cc == c0 else b.left,
                                  right=TEBAL if cc == c0 + LEBAR - 1 else b.right,
                                  top=b.top, bottom=TEBAL)
        for rr in range(r, r + tinggi_kotak):
            kiri = ws.cell(row=rr, column=c0)
            b = kiri.border
            kiri.border = Border(left=TEBAL, right=b.right, top=b.top, bottom=b.bottom)
            kanan = ws.cell(row=rr, column=c0 + LEBAR - 1)
            b = kanan.border
            kanan.border = Border(left=b.left, right=TEBAL, top=b.top, bottom=b.bottom)

        peta[nama] = (r, c0)
        tinggi[j] = r + tinggi_kotak + GAP_BARIS
    ws.freeze_panes = "A4"
    return peta


def kepala_tabel(ws, judul, kolom, lebar, catatan=""):
    """Tata letak SERAGAM di semua sheet: judul baris 1, catatan baris 2,
    kepala kolom baris 4, data mulai baris 5. Diseragamkan supaya rumus cacah
    di BACA-DULU boleh memakai satu pola jangkauan untuk semua sheet."""
    ws["A1"] = judul
    ws["A1"].font = H1
    if catatan:
        ws["A2"] = catatan
        ws["A2"].font = KECIL_M
    r = 4
    for i, (nm, w) in enumerate(zip(kolom, lebar), start=1):
        c = ws.cell(row=r, column=i, value=nm)
        c.font = PUTIH
        c.fill = PatternFill("solid", fgColor="404040")
        c.alignment = Alignment(horizontal="left", vertical="center", indent=1)
        c.border = SEL
        ws.column_dimensions[get_column_letter(i)].width = w
    ws.freeze_panes = ws.cell(row=r + 1, column=1).coordinate
    return r + 1


def main():
    lama = baca_tabel_lama()
    relasi = baca_relasi_lama()
    baru, fk, ck, uq, pk_baru = baca_tabel_baru()

    print("sistem lama : %d tabel, %d kolom" % (len(lama), sum(len(v["kolom"]) for v in lama.values())))
    print("skema baru  : %d tabel, %d kolom, %d FK, %d CHECK, %d UNIQUE"
          % (len(baru), sum(len(v["kolom"]) for v in baru.values()), len(fk), len(ck), len(uq)))
    print("relasi lama dari SQL XML: %d" % len(relasi))

    wb = Workbook()

    # ---------------------------------------------------------- BACA-DULU
    ws = wb.active
    ws.title = "BACA-DULU"
    ws.column_dimensions["A"].width = 34
    ws.column_dimensions["B"].width = 16
    ws.column_dimensions["C"].width = 86
    r = 1
    ws.cell(row=r, column=1, value="ERD TABEL ORACLE - Claim Non Prop").font = H1
    r += 2
    for judul, isi in [
        ("STEMPEL ASAL", "Dibangkitkan alat/buat-erd-excel.py pada %s." % TANGGAL),
        ("", "Sumber sistem lama: pengetahuan/SCHEMA-ACTUAL.csv (698 baris; 598 dari "
             "DDL_Script_ClaimNonProp.xls, 100 dari DDL tabel work Pega)."),
        ("", "Sumber relasi lama: SQL di ekspor XML D:\\XML_NURE\\Claim Non Prop (50 blok pyBrowseSQL)."),
        ("", "Sumber skema baru: ddl-usulan/*.sql - USULAN, belum pernah dijalankan."),
        ("", "Berkas turunan tidak pernah menjadi sumber. Bila workbook ini berbeda dari "
             "sumbernya, ALATNYA yang salah."),
    ]:
        ws.cell(row=r, column=1, value=judul).font = H2
        ws.cell(row=r, column=3, value=isi).font = TXT
        r += 1
    r += 1

    ws.cell(row=r, column=1, value="YANG DIGAMBAR").font = H2
    ws.cell(row=r, column=3, value="TABEL Oracle saja.").font = TXT_K
    r += 1
    for t in ["VIEW (10 di pengetahuan/ddl/, 9 di skema baru) TIDAK digambar - bukan tabel.",
              "PROCEDURE (6) dan FUNCTION (1) TIDAK digambar - bukan bentuk data.",
              "Halaman clipboard Pega TIDAK digambar - ia bukan tabel Oracle.",
              "Index, sequence, dan hak akses TIDAK digambar - bukan bentuk, bukan hubungan."]:
        ws.cell(row=r, column=3, value=t).font = TXT
        r += 1
    r += 1

    ws.cell(row=r, column=1, value="LEGENDA KOTAK").font = H2
    r += 1
    for a, b in [("PK", "kunci primer"), ("UQ", "ikut kunci unik"),
                 ("FK", "foreign key yang tertulis di DDL (hanya skema baru)"),
                 ("~", "relasi TERSIRAT - tidak ditegakkan basis data; buktinya SQL di XML"),
                 ("nama kolom tebal", "NOT NULL"),
                 ("warna kepala kotak", "kelompok domain; daftarnya di sheet KELOMPOK")]:
        ws.cell(row=r, column=2, value=a).font = TXT_K
        ws.cell(row=r, column=3, value=b).font = TXT
        r += 1
    r += 1

    ws.cell(row=r, column=1, value="CACAH").font = H2
    ws.cell(row=r, column=2, value="nilai").font = H2
    ws.cell(row=r, column=3, value="rumusnya menghitung sendiri dari sheet lain").font = KECIL_M
    r += 1
    n_kl = 4 + sum(len(v["kolom"]) for v in lama.values())   # baris terakhir KOLOM-LAMA
    n_kb = 4 + sum(len(v["kolom"]) for v in baru.values())   # baris terakhir KOLOM-BARU
    n_rl = 4 + len(relasi)
    n_rb = 4 + len(fk)
    n_cb = 4 + len(ck) + len(uq)
    n_db = 4 + len(baru)
    cacah = [
        ("Tabel sistem lama",
         "=SUMPRODUCT(('KOLOM-LAMA'!$B$5:$B$%d<>\"\")/COUNTIF('KOLOM-LAMA'!$B$5:$B$%d,"
         "'KOLOM-LAMA'!$B$5:$B$%d&\"\"))" % (n_kl, n_kl, n_kl), "tabel unik di KOLOM-LAMA"),
        ("Kolom sistem lama", "=COUNTA('KOLOM-LAMA'!$D$5:$D$%d)" % n_kl, ""),
        ("Kolom PK sistem lama", "=COUNTIF('KOLOM-LAMA'!$H$5:$H$%d,\"Y\")" % n_kl, ""),
        ("Foreign key sistem lama",
         "=COUNTIF('RELASI-LAMA'!$F$5:$F$%d,\"DITEGAKKAN\")" % n_rl,
         "nol - dan itu hasil hitung, bukan kalimat"),
        ("Relasi tersirat sistem lama", "=COUNTA('RELASI-LAMA'!$A$5:$A$%d)" % n_rl, ""),
        ("Tabel skema baru", "=COUNTA('ERD-BARU-DAFTAR'!$A$5:$A$%d)" % n_db, ""),
        ("Kolom skema baru", "=COUNTA('KOLOM-BARU'!$D$5:$D$%d)" % n_kb, ""),
        ("Kolom NOT NULL skema baru", "=COUNTIF('KOLOM-BARU'!$F$5:$F$%d,\"N\")" % n_kb, ""),
        ("Kolom turunan skema baru", "=COUNTIF('KOLOM-BARU'!$H$5:$H$%d,\"Y\")" % n_kb, ""),
        ("Foreign key skema baru", "=COUNTA('RELASI-BARU'!$A$5:$A$%d)" % n_rb, ""),
        ("CHECK constraint skema baru", "=COUNTIF('CONSTRAINT-BARU'!$C$5:$C$%d,\"CHECK\")" % n_cb, ""),
        ("UNIQUE constraint skema baru", "=COUNTIF('CONSTRAINT-BARU'!$C$5:$C$%d,\"UNIQUE\")" % n_cb, ""),
    ]
    for nm, rumus, ket in cacah:
        ws.cell(row=r, column=1, value=nm).font = TXT
        c = ws.cell(row=r, column=2, value=rumus)
        c.font = TXT_K
        c.alignment = Alignment(horizontal="left")
        if ket:
            ws.cell(row=r, column=3, value=ket).font = KECIL_M
        r += 1
    r += 1
    ws.cell(row=r, column=1, value="BATAS BUKTI").font = H2
    for t in ["Sistem lama tidak punya SATU PUN foreign key. Seluruh 15 relasinya TERSIRAT - "
              "terbaca dari comma-join dan subquery di SQL, bukan ditegakkan basis data.",
              "Relasi hanya disapu dari folder Claim Non Prop; folder Komite Claim Non Prop "
              "tidak dibuka, mengikuti batas BLUEPRINT.md.",
              "25 dari 74 objek di pengetahuan/PULL-LIST.csv belum punya DDL. Tabel yang "
              "disebut SQL tetapi tak ada DDL-nya terdaftar di sheet RELASI-LAMA kolom KEADAAN.",
              "Skema baru BELUM PERNAH DIJALANKAN. Ia memerikan apa yang AKAN ditegakkan."]:
        ws.cell(row=r, column=3, value=t).font = TXT
        r += 1

    # ---------------------------------------------------------- KELOMPOK
    wsk = wb.create_sheet("KELOMPOK")
    rr = kepala_tabel(wsk, "Kelompok domain dan warnanya",
                      ["KELOMPOK", "WARNA KEPALA", "TABEL"], [16, 16, 110])
    for grup in ["Klaim", "Treaty", "Polis", "Pihak", "Dokumen", "Integrasi", "Pega"]:
        anggota = sorted(t for t in lama if KELOMPOK.get(t) == grup)
        wsk.cell(row=rr, column=1, value=grup).font = TXT_K
        c = wsk.cell(row=rr, column=2, value="#" + WARNA[grup])
        c.fill = PatternFill("solid", fgColor=WARNA[grup])
        c.font = Font(name=FONT, size=9, color="FFFFFF")
        wsk.cell(row=rr, column=3, value=", ".join(anggota)).font = TXT
        rr += 1

    # ---------------------------------------------------------- relasi lama
    rel_kolom = {}
    rel_rows = []
    nama_lama = set(lama)
    for (t1, k1, t2, k2, jenis), berkas in sorted(relasi.items()):
        keadaan = []
        for t in (t1, t2):
            if t not in nama_lama:
                keadaan.append("%s bukan TABLE ber-DDL" % t)
        rel_rows.append({
            "kiri": t1, "kiri_kol": k1, "kanan": t2, "kanan_kol": k2,
            "jenis": jenis, "tegak": "TERSIRAT",
            "keadaan": "; ".join(keadaan) if keadaan else "kedua ujung TABLE ber-DDL",
            "berkas": " | ".join(sorted(berkas)[:3]),
        })
        if t1 in nama_lama:
            rel_kolom[(t1, k1)] = "~ %s.%s" % (t2, k2)
        if t2 in nama_lama:
            rel_kolom[(t2, k2)] = "~ %s.%s" % (t1, k1)

    # ---------------------------------------------------------- ERD-LAMA
    wsl = wb.create_sheet("ERD-LAMA")
    gambar_erd(wsl, lama, KELOMPOK, ARTI, rel_kolom, {},
               "ERD sistem lama - 32 TABEL Oracle",
               "Nol foreign key. Tanda ~ adalah relasi TERSIRAT dari SQL di XML, bukan yang "
               "ditegakkan basis data. Kolom tebal = NOT NULL.")

    # ---------------------------------------------------------- KOLOM-LAMA
    wkl = wb.create_sheet("KOLOM-LAMA")
    rr = kepala_tabel(wkl, "Kolom seluruh TABEL Oracle sistem lama",
                      ["OWNER", "TABEL", "NO", "KOLOM", "TIPE", "NULL_BOLEH", "DEFAULT",
                       "PK", "UNIK", "KELOMPOK", "RELASI_TERSIRAT"],
                      [12, 30, 6, 32, 20, 12, 22, 6, 7, 12, 34],
                      "Satu baris per kolom. Sumber: pengetahuan/SCHEMA-ACTUAL.csv.")
    for nama, info in sorted(lama.items()):
        for k in info["kolom"]:
            for i, v in enumerate([info["owner"], nama, k["no"], k["nama"], k["tipe"],
                                   k["null"], k["default"], k["pk"], k["unik"],
                                   KELOMPOK.get(nama, ""),
                                   rel_kolom.get((nama, k["nama"]), "")], start=1):
                c = wkl.cell(row=rr, column=i, value=v)
                c.font = TXT
                c.border = SEL
            rr += 1
    wkl.auto_filter.ref = "A4:K%d" % (rr - 1)

    # ---------------------------------------------------------- RELASI-LAMA
    wrl = wb.create_sheet("RELASI-LAMA")
    rr = kepala_tabel(wrl, "Relasi sistem lama - SELURUHNYA TERSIRAT",
                      ["TABEL_KIRI", "KOLOM_KIRI", "TABEL_KANAN", "KOLOM_KANAN",
                       "BENTUK_DI_SQL", "DITEGAKKAN", "KEADAAN", "BERKAS_XML"],
                      [28, 22, 28, 22, 14, 12, 34, 52],
                      "Nol FOREIGN KEY di seluruh 32 tabel. Relasi di bawah ini dibaca dari "
                      "comma-join dan subquery berkorelasi di 50 blok SQL ekspor XML.")
    for x in rel_rows:
        for i, v in enumerate([x["kiri"], x["kiri_kol"], x["kanan"], x["kanan_kol"],
                               x["jenis"], x["tegak"], x["keadaan"], x["berkas"]], start=1):
            c = wrl.cell(row=rr, column=i, value=v)
            c.font = TXT
            c.border = SEL
        rr += 1
    wrl.auto_filter.ref = "A4:H%d" % (rr - 1)

    # ---------------------------------------------------------- skema baru
    fk_kolom = {}
    for x in fk:
        fk_kolom[(x["tabel"], x["kolom"])] = "-> %s.%s" % (x["ref_tabel"], x["ref_kolom"])
    for nama in baru:
        for k in baru[nama]["kolom"]:
            if (nama, k["nama"]) in fk_kolom:
                k["fk"] = True

    wsb = wb.create_sheet("ERD-BARU")
    gambar_erd(wsb, baru, {}, {n: "skema usulan KLAIMNP" for n in baru}, fk_kolom, pk_baru,
               "ERD skema baru - 22 TABEL usulan KLAIMNP",
               "USULAN, belum pernah dijalankan. Tanda -> adalah FOREIGN KEY yang tertulis "
               "di DDL. Kolom tebal = NOT NULL.")

    wsd = wb.create_sheet("ERD-BARU-DAFTAR")
    rr = kepala_tabel(wsd, "Daftar tabel skema baru",
                      ["TABEL", "BERKAS_DDL", "CACAH_KOLOM", "PK", "CACAH_FK_KELUAR"],
                      [32, 34, 14, 28, 16])
    for nama, info in sorted(baru.items()):
        n_fk = sum(1 for x in fk if x["tabel"] == nama)
        for i, v in enumerate([nama, info["berkas"], len(info["kolom"]),
                               pk_baru.get(nama, ""), n_fk], start=1):
            c = wsd.cell(row=rr, column=i, value=v)
            c.font = TXT
            c.border = SEL
        rr += 1

    wkb = wb.create_sheet("KOLOM-BARU")
    rr = kepala_tabel(wkb, "Kolom seluruh TABEL skema baru KLAIMNP",
                      ["OWNER", "TABEL", "NO", "KOLOM", "TIPE", "NULL_BOLEH", "DEFAULT",
                       "TURUNAN", "PK", "FK_KE"],
                      [12, 30, 6, 34, 26, 12, 22, 10, 6, 34],
                      "Satu baris per kolom. Sumber: ddl-usulan/*.sql.")
    for nama, info in sorted(baru.items()):
        for k in info["kolom"]:
            for i, v in enumerate(["KLAIMNP", nama, k["no"], k["nama"], k["tipe"], k["null"],
                                   k["default"], k.get("turunan", ""),
                                   "Y" if pk_baru.get(nama) == k["nama"] else "",
                                   fk_kolom.get((nama, k["nama"]), "")], start=1):
                c = wkb.cell(row=rr, column=i, value=v)
                c.font = TXT
                c.border = SEL
            rr += 1
    wkb.auto_filter.ref = "A4:J%d" % (rr - 1)

    wrb = wb.create_sheet("RELASI-BARU")
    rr = kepala_tabel(wrb, "Foreign key skema baru - DITEGAKKAN BASIS DATA",
                      ["TABEL", "KOLOM", "REF_TABEL", "REF_KOLOM", "CONSTRAINT", "BERKAS_DDL"],
                      [32, 30, 32, 30, 32, 34],
                      "Seluruhnya tertulis di DDL dan dapat ditunjuk nama constraint-nya. "
                      "Tidak ada garis yang disimpulkan dari kesamaan nama.")
    for x in sorted(fk, key=lambda y: (y["tabel"], y["constraint"])):
        for i, v in enumerate([x["tabel"], x["kolom"], x["ref_tabel"], x["ref_kolom"],
                               x["constraint"], x["berkas"]], start=1):
            c = wrb.cell(row=rr, column=i, value=v)
            c.font = TXT
            c.border = SEL
        rr += 1
    wrb.auto_filter.ref = "A4:F%d" % (rr - 1)

    wcb = wb.create_sheet("CONSTRAINT-BARU")
    rr = kepala_tabel(wcb, "Constraint skema baru selain foreign key",
                      ["TABEL", "CONSTRAINT", "JENIS", "ISI", "BERKAS_DDL"],
                      [30, 34, 12, 96, 32],
                      "CHECK menegakkan aturan di dalam satu baris - ia tidak punya garis di ERD.")
    for x in sorted(ck, key=lambda y: (y["tabel"], y["constraint"])):
        for i, v in enumerate([x["tabel"], x["constraint"], "CHECK", x["isi"], x["berkas"]],
                              start=1):
            c = wcb.cell(row=rr, column=i, value=v)
            c.font = TXT
            c.border = SEL
        rr += 1
    for x in sorted(uq, key=lambda y: (y["tabel"], y["constraint"])):
        for i, v in enumerate([x["tabel"], x["constraint"], "UNIQUE", x["kolom"], x["berkas"]],
                              start=1):
            c = wcb.cell(row=rr, column=i, value=v)
            c.font = TXT
            c.border = SEL
        rr += 1
    wcb.auto_filter.ref = "A4:E%d" % (rr - 1)

    # LibreOffice tidak tersedia di mesin ini, jadi openpyxl menulis rumus TANPA nilai
    # tersimpan. Tanpa penanda ini, pembaca yang mengandalkan nilai tersimpan akan
    # melihat sel kosong. Dengan penanda ini Excel menghitung ulang saat berkas dibuka.
    # PETA-NAMA-T - jembatan antara nama DDL (Indonesia) dan nama T_ yang
    # berlaku di Diagram-Skema-Tabel-NusantaraRe.xlsx sesudah 20-09-2026.
    # Nama DDL di ddl-usulan/ TIDAK diubah; sheet ini yang menerjemahkan.
    import io as _io
    PETA_TSV = os.path.join(AKAR, "4-erd-dan-tabel-datar", "peta-nama-tabel.tsv")
    if os.path.exists(PETA_TSV):
        _teks = _io.open(PETA_TSV, encoding="utf-8").read().splitlines()[1:]
        _br = [x.split(chr(9)) for x in _teks if x.strip()]
        wpn = wb.create_sheet("PETA-NAMA-T")
        rr = kepala_tabel(wpn, "Nama DDL (Indonesia) <-> nama T_ yang berlaku",
                          ["NAMA_T", "PADANAN DI ddl-usulan/", "KELOMPOK",
                           "INDUK", "KUNCI TAMU", "ASAL DI PEGA"],
                          [32, 34, 12, 30, 22, 74],
                          "ddl-usulan/ memakai nama Indonesia; berkas rujukan "
                          "memakai T_. Keduanya SAH - ini penerjemahnya. "
                          "Skema Oracle tetap KLAIMNP.")
        for x in sorted(_br, key=lambda y: y[0]):
            for i, v in enumerate(x[:6], start=1):
                c = wpn.cell(row=rr, column=i, value=v)
                c.font = TXT
                c.border = SEL
            rr += 1
        wpn.auto_filter.ref = "A4:F%d" % (rr - 1)

    wb.calculation.fullCalcOnLoad = True

    wb.save(OUT)
    print("->", OUT)
    print("sheet:", wb.sheetnames)

    # Verifikasi mandiri: hitung di Python apa yang seharusnya dikeluarkan tiap rumus.
    harap = {
        "Tabel sistem lama": len(lama),
        "Kolom sistem lama": sum(len(v["kolom"]) for v in lama.values()),
        "Kolom PK sistem lama": sum(1 for v in lama.values() for k in v["kolom"] if k["pk"] == "Y"),
        "Foreign key sistem lama": 0,
        "Relasi tersirat sistem lama": len(relasi),
        "Tabel skema baru": len(baru),
        "Kolom skema baru": sum(len(v["kolom"]) for v in baru.values()),
        "Kolom NOT NULL skema baru": sum(1 for v in baru.values() for k in v["kolom"]
                                         if k["null"] == "N"),
        "Kolom turunan skema baru": sum(1 for v in baru.values() for k in v["kolom"]
                                        if k.get("turunan") == "Y"),
        "Foreign key skema baru": len(fk),
        "CHECK constraint skema baru": len(ck),
        "UNIQUE constraint skema baru": len(uq),
    }
    print("\nnilai yang seharusnya keluar dari tiap rumus di BACA-DULU:")
    for k, v in harap.items():
        print("   %-32s %s" % (k, v))


if __name__ == "__main__":
    main()
