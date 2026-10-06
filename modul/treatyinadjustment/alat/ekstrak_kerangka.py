"""Pembangkit kerangka tab layar Adjustment dari ekspor Pega.

Membaca Section di folder korpus `Treaty In Adjustment`, membuang SELURUH
`<pyIncludedRuleXML>` dengan MENGHITUNG KEDALAMAN SARANG (bukan regex
non-greedy — lihat `docs/LAYAR-ADJUSTMENT.md` §1), lalu menulis
`frontend/ekspor/kerangka.gen.ts`: blok, grid (kolom · lebar · desimal per
sel), medan, teks, tombol, dan include — tiap butir dengan offset bitanya.

Butir yang MATI (penjaganya `1=2`, `1==2`, `3=4`, `Never`, `false`, atau
`… && NEVER`) TIDAK ditulis ke kerangka; ia dicatat di `DIBUANG` beserta
alasannya. ⛔ Yang diperiksa adalah butir YANG DIJAGA — bukan wilayah di
sekitarnya: penjaga yang membungkus tombol tidak mematikan nilai di
sebelahnya.

Pemakaian (dari APP_RNM):
    python modul/treatyinadjustment/alat/ekstrak_kerangka.py "D:/XML_NURE/Treaty In Adjustment"
"""

import json
import os
import re
import sys
import xml.parsers.expat

TAG = re.compile(r"<(/?)pyIncludedRuleXML\b[^>]*?(/?)>")


def buang_pratinjau(t):
    """Buang `<pyIncludedRuleXML>` dengan menghitung kedalaman sarang."""
    out, i, depth, start = [], 0, 0, 0
    for m in TAG.finditer(t):
        tutup, sendiri = m.group(1), m.group(2)
        if sendiri:
            if depth == 0:
                out.append(t[i:m.start()])
                i = m.end()
            continue
        if not tutup:
            if depth == 0:
                out.append(t[i:m.start()])
                start = m.start()
            depth += 1
        else:
            depth -= 1
            if depth == 0:
                i = m.end()
    if depth != 0:
        raise ValueError("pyIncludedRuleXML tidak seimbang")
    out.append(t[i:])
    return "".join(out)


class N:
    __slots__ = ("tag", "off", "kids", "txt", "par")

    def __init__(s, tag, off, par):
        s.tag, s.off, s.kids, s.txt, s.par = tag, off, [], "", par

    def k(s, tag):
        return [c for c in s.kids if c.tag == tag]

    def v(s, tag):
        for c in s.kids:
            if c.tag == tag:
                return c.txt.strip()
        return ""

    def iter(s, tag=None):
        if tag is None or s.tag == tag:
            yield s
        for c in s.kids:
            yield from c.iter(tag)


def muat(path):
    t = buang_pratinjau(open(path, encoding="utf-8", errors="replace").read())
    b = t.encode("utf-8")
    p = xml.parsers.expat.ParserCreate()
    root = N("#", 0, None)
    st = [root]

    def se(n, a):
        x = N(n, p.CurrentByteIndex, st[-1])
        st[-1].kids.append(x)
        st.append(x)

    p.StartElementHandler = se
    p.EndElementHandler = lambda n: st.pop()

    def cd(d):
        st[-1].txt += d

    p.CharacterDataHandler = cd
    p.Parse(b, True)
    return root


MATI = re.compile(
    r"^\s*(1\s*==?\s*2|3\s*==?\s*4|never|false)\s*$|&&\s*(1\s*==?\s*2|never|false)\s*$|^\s*(never|false)\s*&&",
    re.I,
)
IDENTITAS = re.compile(r"OperatorID\.(pxInsName|pyUserName|pyUserIdentifier)")
PENANDA = {"Title", "Mobile reveal", "Mobile dismiss", "Desktop dismiss", "Reveal", "New item", ""}
LABEL_KOSONG = {"Text Input", "Spacer", "Button", "Checkbox", "Label", "Button Template", ""}


def syarat_sendiri(n):
    s = []
    cvw = n.v("pyContainerVisibleWhen")
    if cvw:
        s.append(cvw)
    for u in n.k("pyUserData"):
        if u.v("pyVisible") == "OTHER" and u.v("pyCondition"):
            s.append(u.v("pyCondition"))
    return s


def desimal(c):
    for m in c.k("pyModes"):
        for d in m.iter("pyDecimalPlaces"):
            if d.txt.strip():
                return int(d.txt.strip())
    return None


def ikat(prop, sisi_lama):
    """`TreatyIn.OLDDATA.X` → sisi; `TreatyIn.X` → akar (Section Old) / sisi (New)."""
    if prop.startswith("TreatyIn.OLDDATA."):
        return "sisi", prop[len("TreatyIn.OLDDATA."):]
    if prop.startswith("TreatyIn."):
        return ("akar" if sisi_lama else "sisi"), prop[len("TreatyIn."):]
    return None, prop


class Pembangkit:
    def __init__(self, berkas, sisi_lama):
        self.berkas, self.sisi_lama, self.dibuang = berkas, sisi_lama, []

    def buang(self, n, jenis, alasan):
        self.dibuang.append({"berkas": self.berkas, "at": n.off, "jenis": jenis, "alasan": alasan})

    def periksa(self, n, jenis):
        """Syarat butir ini; None bila butirnya mati / bersyarat identitas."""
        s = syarat_sendiri(n)
        for x in s:
            if MATI.search(x):
                self.buang(n, jenis, f"penjaga mati: {x}")
                return None
            if IDENTITAS.search(x):
                self.buang(n, jenis, "syarat identitas operator (blok dev)")
                return None
        return s

    def grid(self, n, syarat):
        prop = n.v("pyPageListProperty")
        dari, larik = ikat(prop, self.sisi_lama)
        baris = []
        for t in n.k("pyTable"):
            for r in t.k("pyRows"):
                for row in r.k("rowdata"):
                    baris.append([c for cs in row.k("pyCells") for c in cs.k("rowdata")])
        kepala = baris[0] if baris else []
        badan = baris[1] if len(baris) > 1 else []
        kol = {"kolom": [], "kunci": [], "lebar": [], "desimal": [], "format": [], "syaratSel": [], "atSel": []}
        for i, c in enumerate(badan):
            nilai, fmt = c.v("pyValue"), c.v("pyFormat")
            sy = syarat_sendiri(c)
            if fmt == "pxButton" or "pyTemplateButton" in nilai:
                self.buang(c, "kolom", "kolom tombol (jalur tulis)")
                continue
            if any(MATI.search(x) for x in sy):
                self.buang(c, "kolom", f"penjaga mati: {sy}")
                continue
            h = kepala[i] if i < len(kepala) else None
            kol["kolom"].append(h.v("pyValue") if h is not None else "")
            kol["kunci"].append(nilai[1:] if nilai.startswith(".") else nilai)
            kol["lebar"].append(int(c.v("pyWidth") or 0))
            kol["desimal"].append(desimal(c))
            kol["format"].append(fmt)
            kol["syaratSel"].append(sy[0] if sy else None)
            kol["atSel"].append(c.off)
        g = {"t": "grid", "at": n.off, "prop": prop, "dari": dari, "larik": larik, "syarat": syarat, **kol}
        tmpl = n.v("pyGridTemplateName")
        if tmpl:
            g["templatBaris"] = tmpl
        return g

    def sel(self, c, syarat):
        nilai, fmt, typ = c.v("pyValue"), c.v("pyFormat"), c.v("pyType")
        label = c.v("pyLabelFieldValue")
        if fmt == "pxButton" or "pyTemplateButton" in nilai:
            lbl = [x.txt.strip() for x in c.iter("pyLabel") if x.txt.strip()]
            if not lbl:
                self.buang(c, "tombol", "tombol tanpa label (ikon)")
                return None
            return {"t": "tombol", "at": c.off, "label": lbl[0], "syarat": syarat}
        if typ == "LABEL":
            if not nilai:
                return None
            return {"t": "teks", "at": c.off, "teks": nilai, "syarat": syarat}
        if not nilai.startswith("TreatyIn."):
            self.buang(c, "medan", f"ikatan bukan halaman TreatyIn: {nilai!r}")
            return None
        dari, kunci = ikat(nilai, self.sisi_lama)
        cap = [x.txt.strip() for x in c.iter("pyCheckboxCaption") if x.txt.strip()]
        m = {
            "t": "medan", "at": c.off, "label": "" if label in LABEL_KOSONG else label,
            "dari": dari, "kunci": kunci, "format": fmt, "desimal": desimal(c), "syarat": syarat,
        }
        if cap:
            m["caption"] = cap[0]
        return m

    def jalan(self, n):
        out = []
        for c in n.kids:
            if c.tag != "rowdata":
                out += self.jalan(c)
                continue
            if c.v("pyPageListProperty"):
                s = self.periksa(c, "grid")
                if s is not None:
                    out.append(self.grid(c, s))
                continue
            if c.v("pyInclude"):
                s = self.periksa(c, "include")
                if s is not None:
                    out.append({"t": "include", "at": c.off, "nama": c.v("pyInclude"), "syarat": s})
                continue
            if c.par is not None and c.par.tag == "pyCells" and c.v("pyType") in ("FIELD", "LABEL"):
                s = self.periksa(c, "sel")
                if s is not None:
                    b = self.sel(c, s)
                    if b:
                        out.append(b)
                continue
            judul, kepala = c.v("pyTitle"), c.v("pyHeaderType")
            s = syarat_sendiri(c)
            if s or (judul and judul not in PENANDA):
                s = self.periksa(c, "blok")
                if s is None:
                    continue
                anak = self.jalan(c)
                tampil = judul if (judul not in PENANDA and kepala in ("BAR", "TABBED")) else ""
                out.append({"t": "blok", "at": c.off, "judul": tampil, "syarat": s, "anak": anak})
                continue
            out += self.jalan(c)
        return out


def tab(root, judul):
    for x in root.iter("rowdata"):
        if x.v("pyTitle") == judul:
            return x
    raise KeyError(judul)


def bangkit(korpus):
    sec = lambda n: os.path.join(korpus, "Section", n + ".xml")
    kerangka, include, dibuang = {}, {}, []
    TAB = {
        ("TreatyInTabsNonProportional", False): [
            "Maximum Retention", "Event Limits", "EGNPI", "Limits", "Share", "Retro", "Installment",
            "Value Difference", "Exclusions", "Special Conditions", "Information & Submit"],
        ("TreatyInTabsNonProportionalOldData", True): [
            "Maximum Retention", "EGNPI", "Limits", "Share", "Retro", "Installment", "Exclusions", "Special Conditions"],
        ("TreatyInTabsProportional", False): [
            "Reporting Period", "Portfolio", "Limits", "Share", "Retro", "Co-Ins Scale", "Accumulation",
            "Exclusions", "Special Conditions", "Information & Submit", "Achievement In IDR"],
        ("TreatyInTabsProportionalOldData", True): [
            "Reporting Period", "Portfolio", "Limits", "Share", "Accumulation", "Exclusions", "Special Conditions"],
    }
    for (nama, lama), daftar in TAB.items():
        r = muat(sec(nama))
        p = Pembangkit(nama, lama)
        for j in daftar:
            t = tab(r, j)
            kerangka[f"{nama}#{j}"] = {"at": t.off, "syarat": syarat_sendiri(t), "isi": p.jalan(t)}
        dibuang += p.dibuang
    # ⛔ Isi Retro (TreatyInFacultativeShareCalculation*, TreatyInFacultativeRetro)
    # dan WorkAttachments TIDAK dibangkitkan — keputusan §17 / panel sendiri.
    for nama, lama in [
        ("TreatyInTabsNonProportionalOldDataShare", True), ("TreatyInShareProp", False),
        ("TreatyInfoSubmit", False), ("TreatyInTabsAchievement", False),
        ("TreatyInTabsNonProportionalValueDifference", False),
        ("TreatyInTabsNPValueDifferenceProRate", False), ("TreatyInTabsNPValueDifference_NoProRate", False),
    ]:
        r = muat(sec(nama))
        p = Pembangkit(nama, lama)
        include[nama] = p.jalan(r)
        dibuang += p.dibuang
    kurs = {}
    for nama, lama, prop in [("TreatyInNONProportionalOldData", True, "TreatyIn.OLDDATA.CurrencyList"),
                             ("TreatyInNONProportional", False, "TreatyIn.CurrencyList")]:
        r = muat(sec(nama))
        p = Pembangkit(nama, lama)
        for e in r.iter("pyPageListProperty"):
            if e.txt.strip() == prop:
                kurs["lama" if lama else "baru"] = p.grid(e.par, [])
                break
        dibuang += p.dibuang
    return kerangka, include, kurs, dibuang


def main():
    korpus = sys.argv[1]
    kerangka, include, kurs, dibuang = bangkit(korpus)
    sini = os.path.dirname(os.path.abspath(__file__))
    tujuan = os.path.join(sini, "..", "frontend", "ekspor", "kerangka.gen.ts")
    j = lambda o: json.dumps(o, ensure_ascii=False, indent=1)
    with open(tujuan, "w", encoding="utf-8", newline="\n") as f:
        f.write("// ⛔ BERKAS BANGKITAN — JANGAN DISUNTING TANGAN.\n")
        f.write("// Dibangkitkan `alat/ekstrak_kerangka.py` dari Section ekspor `Treaty In Adjustment`,\n")
        f.write("// sesudah `<pyIncludedRuleXML>` dibuang dengan menghitung kedalaman sarang.\n")
        f.write("// Offset `at` = posisi bita di berkas Section SESUDAH pembuangan itu.\n\n")
        f.write("import type { ButirKerangka, GridKerangka, Kerangka, Terbuang } from './jenis'\n\n")
        f.write(f"export const KERANGKA_TAB: Readonly<Record<string, Kerangka>> = {j(kerangka)}\n\n")
        f.write(f"export const KERANGKA_INCLUDE: Readonly<Record<string, readonly ButirKerangka[]>> = {j(include)}\n\n")
        f.write(f"export const GRID_KURS: Readonly<Record<'lama' | 'baru', GridKerangka>> = {j(kurs)}\n\n")
        f.write(f"export const DIBUANG: readonly Terbuang[] = {j(dibuang)}\n")
    medan, larik = kunci_terbaca(kerangka, include, kurs)
    go = os.path.join(sini, "..", "backend", "repository", "kunci_kerangka_gen.go")
    daftar = lambda xs: "".join(f"\t{json.dumps(x)},\n" for x in xs)
    with open(go, "w", encoding="utf-8", newline="\n") as f:
        f.write("// Code generated by alat/ekstrak_kerangka.py. DO NOT EDIT.\n\n")
        f.write("package repository\n\n")
        f.write("// MedanKerangka - kunci skalar yang kerangka tab bangkitan baca, termasuk\n")
        f.write("// kunci yang disebut SYARAT tampilnya. Dibaca bersama `MedanDibaca`.\n")
        f.write("var MedanKerangka = []string{\n" + daftar(medan) + "}\n\n")
        f.write("// LarikKerangka - larik yang grid kerangka tab bangkitan baca.\n")
        f.write("var LarikKerangka = []string{\n" + daftar(larik) + "}\n")
    print(f"{len(kerangka)} tab, {len(include)} include, {len(dibuang)} butir dibuang, "
          f"{len(medan)} medan, {len(larik)} larik -> {os.path.normpath(tujuan)}")


SYARAT_KUNCI = re.compile(r"TreatyIn\.(\w+)")


def kunci_terbaca(kerangka, include, kurs):
    medan, larik = set(), set()

    def jalan(bs):
        for b in bs:
            for s in b.get("syarat", []):
                medan.update(SYARAT_KUNCI.findall(s))
            if b["t"] == "blok":
                jalan(b["anak"])
            elif b["t"] == "medan":
                medan.add(b["kunci"])
            elif b["t"] == "grid":
                larik.add(b["larik"])
                for s in b["syaratSel"]:
                    if s:
                        medan.update(SYARAT_KUNCI.findall(s))

    for k in kerangka.values():
        for s in k["syarat"]:
            medan.update(SYARAT_KUNCI.findall(s))
        jalan(k["isi"])
    for bs in include.values():
        jalan(bs)
    jalan(list(kurs.values()))
    return sorted(medan), sorted(larik)


if __name__ == "__main__":
    main()
