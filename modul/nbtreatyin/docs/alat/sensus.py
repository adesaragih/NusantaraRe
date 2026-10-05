"""Sensus properti PolicyTreatyIn dan Quotation dari rule yang TERJANGKAU - tiket 00.

HANYA MEMBACA korpus.json (hasil korpus.py). Menulis ../SENSUS-PROPERTI-POLICYTREATYIN.md.

Rujukan yang dihitung:
  - jalur berjangkar: `PolicyTreatyIn.X`, `pyWorkPage.PolicyTreatyIn.X`, `Quotation.X`,
    `PolicyTreatyIn.QuotationData.X`
  - rujukan RELATIF `.X` di rule/langkah yang halamannya `PolicyTreatyIn`
    (kelas ASM-FW-GISFW-Data-PolicyTreatyIn) - titik buta sapuan berjangkar
    yang dicatat DAFTAR-MEDAN (SuggestList, CedingCoList, BreakDownSpreadList)
  - anggota daftar: `.<Daftar>(n).Y`, dan `.Y` di dalam kalang berhalaman daftar itu

Penggolongan skalar lawan daftar: nama yang muncul dengan indeks `(...)`, sebagai
halaman langkah kalang, atau sebagai PageList di sel grid = daftar.
"""
import collections
import json
import os
import re
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from graf import Graf, _langkah_rata  # noqa: E402

DIR = os.path.dirname(os.path.abspath(__file__))
KELAS_POLIS = "ASM-FW-GISFW-Data-PolicyTreatyIn"
NAMA = r"[A-Za-z_][A-Za-z0-9_]*"
SISTEM = re.compile(r"^(px|py|pz)")


class Sensus:
    def __init__(self):
        self.skalar = collections.defaultdict(set)      # nama -> {rule}
        self.daftar = collections.defaultdict(set)      # nama daftar -> {rule}
        self.anggota = collections.defaultdict(lambda: collections.defaultdict(set))
        self.quotation = collections.defaultdict(set)
        self.layar = set()

    def jalur(self, teks, rule, relatif_polis=False, daftar_kalang=None):
        if not teks:
            return
        for m in re.finditer(r"(?:pyWorkPage\.)?PolicyTreatyIn\.(" + NAMA + r")(\([^)]*\))?(?:\.(" + NAMA + r"))?", teks):
            self._catat(m.group(1), bool(m.group(2)), m.group(3), rule, quot=False)
        for m in re.finditer(r"(?<![\w.])(?:pyWorkPage\.)?Quotation\.(" + NAMA + r")", teks):
            if not SISTEM.match(m.group(1)):
                self.quotation[m.group(1)].add(rule)
        for m in re.finditer(r"PolicyTreatyIn\.QuotationData\.(" + NAMA + r")", teks):
            self.quotation[m.group(1)].add(rule)
        if relatif_polis or daftar_kalang:
            for m in re.finditer(r"(?<![\w)\].])\.(" + NAMA + r")(\([^)]*\))?(?:\.(" + NAMA + r"))?", teks):
                nama, idx, sub = m.group(1), m.group(2), m.group(3)
                if SISTEM.match(nama) or nama in ("PolicyTreatyIn", "Quotation", "QuotationData"):
                    continue
                if daftar_kalang:
                    self.anggota[daftar_kalang][nama].add(rule)
                else:
                    self._catat(nama, bool(idx), sub, rule, quot=False)

    def _catat(self, nama, berindeks, sub, rule, quot):
        if SISTEM.match(nama) or nama == "QuotationData":
            return
        if berindeks:
            self.daftar[nama].add(rule)
            if sub and not SISTEM.match(sub):
                self.anggota[nama][sub].add(rule)
        else:
            self.skalar[nama].add(rule)


def _semua(x):
    if isinstance(x, str):
        yield x
    elif isinstance(x, dict):
        for k, v in x.items():
            if k not in ("rujukan",):
                yield from _semua(v)
    elif isinstance(x, list):
        for v in x:
            yield from _semua(v)


def kumpulkan(data, R):
    s = Sensus()
    for d in data:
        if (d["jenis"], d["nama"]) not in R:
            continue
        rule = f'{d["jenis"]}/{d["nama"]}'
        polis_kelas = d["kelas"] == KELAS_POLIS
        if d["jenis"] == "Activity":
            def jalan(ls, halaman_induk):
                for st in ls:
                    hal = st["halaman"] or halaman_induk
                    m = re.search(r"PolicyTreatyIn\.(" + NAMA + r")$", hal or "")
                    kalang = None
                    if st["ulang"] and m:
                        kalang = m.group(1)
                        s.daftar[kalang].add(rule)
                    relatif = (polis_kelas and not hal) or bool(re.search(r"PolicyTreatyIn$", hal or ""))
                    teks = [st["ket"]] + [p["when"] for p in st["pre"]] + [a for pr in st["param"] for a in pr] \
                        + list(st["panggil"].values())
                    for tks in teks:
                        s.jalur(tks, rule, relatif_polis=relatif and not kalang, daftar_kalang=kalang)
                    jalan(st["anak"], hal if not st["ulang"] else hal)
            jalan(d["langkah"], "")
        elif d["jenis"] in ("Section", "Harness", "FlowAction"):
            for c in d.get("sel", []):
                for tks in (c["nilai"], c["tampilWhen"], c["kunciWhen"], c["wajibWhen"], c["nonaktifWhen"]):
                    s.jalur(tks, rule, relatif_polis=polis_kelas)
                if polis_kelas and c["nilai"].startswith(".") and c["kontrol"]:
                    s.layar.add(c["nilai"][1:].split(".")[0])
        else:
            for tks in _semua({k: v for k, v in d.items() if k not in ("nama", "berkas", "kelas", "jenis")}):
                s.jalur(tks, rule, relatif_polis=polis_kelas)
    # nama yang ternyata daftar tidak dihitung skalar
    for n in list(s.skalar):
        if n in s.daftar:
            del s.skalar[n]
    return s


def tulis(s):
    w = ["# Sensus properti `PolicyTreatyIn` dan `Quotation` — dari XML yang terjangkau", "",
         "> Bangkitan `docs/alat/sensus.py` atas korpus.json (rule TERJANGKAU saja, `graf.py`). "
         "Tiket 00: *\"Cocokkan ulang ke XML (termasuk DataTransform, When, dan DecisionTable)\"*.", "",
         f"- Skalar `PolicyTreatyIn`: **{len(s.skalar)}**",
         f"- Daftar (PageList) `PolicyTreatyIn`: **{len(s.daftar)}**",
         f"- Properti `Quotation` / `QuotationData`: **{len(s.quotation)}**",
         f"- Tampil di layar (sel berhalaman PolicyTreatyIn): **{len(s.layar)}**", "",
         "## Skalar `PolicyTreatyIn`", "", "| Properti | Layar | Dirujuk (rule) |", "| --- | :---: | --- |"]
    for n in sorted(s.skalar, key=str.lower):
        r = sorted(s.skalar[n])
        w.append(f"| `{n}` | {'⭐' if n in s.layar else ''} | {len(r)}: {', '.join(x.split('/')[1] for x in r[:6])}{' …' if len(r) > 6 else ''} |")
    w += ["", "## Daftar `PolicyTreatyIn` dan anggotanya", ""]
    for n in sorted(s.daftar, key=str.lower):
        ang = s.anggota.get(n, {})
        w.append(f"- `{n}` — {len(s.daftar[n])} rule; anggota: "
                 + (", ".join(f"`{a}`" for a in sorted(ang, key=str.lower)) or "—"))
    w += ["", "## `Quotation` / `QuotationData`", "", ", ".join(f"`{n}`" for n in sorted(s.quotation, key=str.lower)), ""]
    return "\n".join(w) + "\n"


if __name__ == "__main__":
    sumber = sys.argv[1] if len(sys.argv) > 1 else "korpus.json"
    with open(sumber, encoding="utf-8") as fh:
        data = json.load(fh)
    g = Graf(data)
    s = kumpulkan(data, g.terjangkau())
    teks = tulis(s)
    with open(os.path.join(DIR, "..", "SENSUS-PROPERTI-POLICYTREATYIN.md"), "w", encoding="utf-8", newline="\n") as fh:
        fh.write(teks)
    print(len(s.skalar), "skalar,", len(s.daftar), "daftar,", len(s.quotation), "quotation")
