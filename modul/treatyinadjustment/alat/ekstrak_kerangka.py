"""Pembangkit kerangka tab layar Adjustment dari ekspor Pega.

Membaca Section di folder korpus `Treaty In Adjustment`, membuang SELURUH
`<pyIncludedRuleXML>` dengan MENGHITUNG KEDALAMAN SARANG (bukan regex
non-greedy — lihat `docs/LAYAR-ADJUSTMENT.md` §1), lalu menulis
`frontend/ekspor/kerangka.gen.ts`: blok, grid (kolom · lebar · desimal per
sel), medan, teks, tombol, dan include — tiap butir dengan offset bitanya.

RINCIAN BARIS (`pyEditingMode = expandPane`): grid yang membuka panel di
bawah barisnya menyebut FLOW ACTION di `pyEditAction`; flow action itu
menunjuk Section lewat `pySectionReference`. Section itu dibangkitkan ke
`KERANGKA_RINCIAN` dalam MODE RELATIF — ikatannya `.X` (halaman baris),
bukan `TreatyIn.X`. ⛔ Templat `…!pyGridRowDetails` hanya bingkainya; isinya
flow action tadi, dan flow action itu ADA di ekspor.

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
    # ⛔ `pyIsVisibilityOption = ALWAYS` MENIMPA `pyContainerVisibleWhen` —
    # pasangan yang sama dengan `pyVisible = ALWAYS` atas `pyCondition` pada
    # sel. Syarat yang tertinggal (`1=2` pada grid Co-Ins Scale Prop) hanya
    # sisa; membacanya sebagai penjaga pernah membuang grid yang tampil di
    # Pega (gambar 18 Treaty In).
    if cvw and n.v("pyIsVisibilityOption") != "ALWAYS":
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


def kunci_baca(c):
    """Kapan sel ini BACA-SAJA: "selalu", daftar syarat (salah satu benar ->
    baca-saja), atau None (dapat disunting di mode Edit).

    Diukur 7 Oktober 2026 atas Section panel New (sel berikatan, tanpa
    tombol): `Auto` tanpa syarat 103 · `Auto` + `pyDisabledWhen` 41 ·
    `Read-only` + `pyReadOnlyCondition` 23 (+ `pyDisabledWhen` 19) ·
    `Read-only` TANPA syarat 18. Jadi `Read-only` sendiri berarti selalu;
    `Read-only` bersyarat berarti baca-saja HANYA bila syaratnya benar.
    `pyDisabledWhen` (sel `pyDisabled = true`) menambah syarat kunci —
    termasuk kunci per Material Type (`TreatyIn.EDMMaterialType = 2`).
    """
    if c.v("pyFormat") == "pxDisplayText":
        return "selalu"
    roc = ""
    for u in c.k("pyUserData"):
        roc = roc or u.v("pyReadOnlyCondition")
    dis = []
    for m in c.k("pyModes"):
        for r in m.k("rowdata"):
            if r.v("pyDisabled") == "true":
                dis.append(r.v("pyDisabledWhen") or "selalu")
    if (c.v("pyEditOptions") == "Read-only" and not roc) or "selalu" in dis:
        return "selalu"
    s = ([roc] if roc else []) + dis
    return s or None


def sumber_pilihan(c):
    """Sumber daftar sel `pxDropdown`/`pxAutoComplete`, atau None.

    `pyListDataSource/pyListSource`: `reportdefinition` membawa nama RD
    (`pyReportDefinitionPage/pySourceName`) beserta medan NILAI (`pyValue`)
    dan TAMPIL (`pyPrompt`); `associated` = daftar milik rule Property yang
    TIDAK ikut diekspor; `pageList` = halaman sesi (`pyCBPage`).
    """
    if c.v("pyFormat") not in ("pxDropdown", "pxAutoComplete"):
        return None
    sumber, rd, nilai, tampil, halaman = "", "", "", "", ""
    for m in c.iter("pyListDataSource"):
        sumber = sumber or m.v("pyListSource")
        for r in m.k("pyReportDefinitionPage"):
            rd = rd or r.v("pySourceName")
            nilai = nilai or r.v("pyValue")
            tampil = tampil or r.v("pyPrompt")
        for r in m.k("pyCBPage"):
            halaman = halaman or r.v("pySourceName")
    if not sumber:
        return None
    out = {"sumber": sumber}
    if sumber == "reportdefinition" and rd:
        out["rd"] = rd
        if nilai:
            out["nilai"] = nilai.lstrip(".")
        if tampil:
            out["tampil"] = tampil.lstrip(".")
    if sumber == "pageList" and halaman:
        out["halaman"] = halaman
    return out


ANGKA = re.compile(r"-?\d+(\.\d+)?")


def syarat_aksi(br, milik):
    """`pyActionConditions` satu baris perilaku -> SATU teks syarat, atau "".

    Tiap baris satu perbandingan: `Other Property` = `pyOtherOperandLeft`
    + `pyOtherOperandLogic` + `pyOtherOperandRightText`; `Associated
    Property` = properti KONTROL itu sendiri (`milik`) + operand Associated.

    ⚠️ PENGGABUNG BARIS. `pyLogic` tersimpan di SETIAP baris (bawaan `And`)
    dan arah sambungnya tidak terbaca dari korpus (claimprop grilling-
    ronde-2 R7-baru). Dibaca begini, dan HANYA begini:
      - semua baris membandingkan properti YANG SAMA dengan `=` -> `||`.
        `.Note = 'SURPLUS' && .Note = '2ND SURPLUS'` tidak pernah benar,
        jadi pembacaan AND berarti aksinya mati; Treaty In membacanya
        sebagai salah satu (`SURPLUS_MATA_UANG`, `syaratSurplus`);
      - selain itu -> `&&` (bawaan `And`), mis. `EDMMaterialType =
        7897987 && EDMState != 3` — Treaty In membacanya begitu juga.
    Jenis syarat lain MENGGAGALKAN pembangkitan — tidak ditebak."""
    baris = []
    for ac in br.k("pyActionConditions"):
        for r in ac.k("rowdata"):
            jenis = r.v("pyCondition")
            if jenis == "Other Property":
                kiri, op, kanan = r.v("pyOtherOperandLeft"), r.v("pyOtherOperandLogic"), r.v("pyOtherOperandRightText")
            elif jenis == "Associated Property":
                kiri, op, kanan = milik, r.v("pyAssociatedOperandLogic"), r.v("pyAssociatedOperandRightText")
            else:
                raise ValueError(f"pyActionConditions jenis tak dikenal: {jenis!r} @{r.off}")
            if not kiri:
                raise ValueError(f"pyActionConditions tanpa operand kiri @{r.off}")
            baris.append((kiri, op or "=", kanan))
    if not baris:
        return ""

    def nilai(v):
        if v in ('""', ""):
            return "''"
        return v if ANGKA.fullmatch(v) else f"'{v}'"

    teks = [f"{k} {o} {nilai(v)}" for k, o, v in baris]
    sama = len({k for k, _, _ in baris}) == 1 and all(o == "=" for _, o, _ in baris)
    return (" || " if sama else " && ").join(teks)


def aksi_tombol(c, peristiwa=None):
    """Aksi klik sebuah tombol, berurutan: `refresh` + Activity + parameter,
    `setValue` + pasangan, `localAction`, `addRow`, `deleteRow`, ...

    `peristiwa` menyaring baris perilaku menurut `pyEvent` (himpunan;
    `change` untuk medan dan sel); None = semua (tombol). Aksi yang dijaga
    `pyActionConditions` membawa `syarat` — lihat `syarat_aksi`."""
    out = []
    for m in c.k("pyModes"):
        for r in m.k("rowdata"):
            for b in r.k("pyBehaviors"):
                for br in b.k("rowdata"):
                    if peristiwa is not None and br.v("pyEvent") not in peristiwa:
                        continue
                    a = {"aksi": br.v("pyAction")}
                    for api in br.k("pyActionAPI"):
                        for tag in ("pyActivity", "pyLocalAction", "pyFlowAction", "pyDataTransform"):
                            if api.v(tag) and "aktivitas" not in a:
                                a["aktivitas"] = api.v(tag)
                        # DataTransform PRA-refresh (`pyPreDataTransform/pyName`) —
                        # Add Accumulation memakai `TreatyInAddAccumulation` begini.
                        for pdt in api.k("pyPreDataTransform"):
                            if pdt.v("pyName"):
                                a["transformasi"] = pdt.v("pyName")
                                # Parameter DataTransform itu — `AddLimitRetentionCession`
                                # (`type = limit|retention|cession`), `CalculateReinstatement`
                                # (`Subscript = .pxListSubscript`), ...
                                pdp = {}
                                for dp in pdt.k("pyDataTransformParams"):
                                    for pr in dp.k("rowdata"):
                                        if pr.v("pyName"):
                                            pdp[pr.v("pyName")] = pr.v("pyValue")
                                if pdp:
                                    a["paramDT"] = pdp
                        prm = {}
                        for ap in api.k("pyActivityParams"):
                            for pr in ap.k("rowdata"):
                                if pr.v("pyName"):
                                    prm[pr.v("pyName")] = pr.v("pyValue")
                        for nv in api.k("pyNameValuePairs"):
                            for pr in nv.k("rowdata"):
                                if pr.v("pyName"):
                                    prm[pr.v("pyName")] = pr.v("pyValue")
                        if prm:
                            a["param"] = prm
                    sy = syarat_aksi(br, c.v("pyValue"))
                    if sy:
                        a["syarat"] = sy
                    if a["aksi"]:
                        out.append(a)
    return out


def aksi_ubah(c):
    """Aksi `change` sebuah medan/sel — Activity yang Pega jalankan sesudah
    isiannya berubah (mis. `Installment` → `TreatyInSetValueInstallment`).
    ⛔ Hanya `change`: `keyboard`/`enter` mengulang aksi yang sama.
    ⭐ Kotak centang memakai `click` (Share Across The Board →
    `TreatyInXOLAddSpreading`, `TreatyInSetBrokerage`): bagi centang, klik
    ADALAH perubahan nilainya."""
    if c.v("pyFormat") == "pxCheckbox":
        return aksi_tombol(c, ("change", "click"))
    return aksi_tombol(c, ("change",))


def tombol(c, syarat):
    """Satu tombol, LENGKAP — termasuk tombol ikon tanpa `pyLabel`.

    ⛔ Dulu tombol ikon dibuang dan kolom tombol grid dibuang sebagai "jalur
    tulis". Keduanya bagian dari tata letak Pega (Add/Delete grid, Tambah
    Co-Ins) dan sekarang dibawa; yang memutuskan hidup-matinya perender.
    """
    lbl = [x.txt.strip() for x in c.iter("pyLabel") if x.txt.strip()]
    img = [x.txt.strip() for x in c.iter("pyImage") if x.txt.strip()]
    t = {"t": "tombol", "at": c.off, "label": lbl[0] if lbl else "", "syarat": syarat, "aksi": aksi_tombol(c)}
    if img:
        t["ikon"] = img[0].split("/")[-1]
    dis = []
    for m in c.k("pyModes"):
        for r in m.k("rowdata"):
            if r.v("pyDisabled") == "true" and r.v("pyDisabledWhen"):
                dis.append(r.v("pyDisabledWhen"))
    if dis:
        t["nonaktif"] = dis
    return t


def seksi_flowaction(korpus, nama):
    """`FlowAction/<nama>.xml` → Section yang ia tampilkan, atau "" bila tidak ada."""
    f = os.path.join(korpus, "FlowAction", nama + ".xml")
    if not os.path.exists(f):
        return ""
    m = re.search(r"<pySectionReference>([^<]+)</pySectionReference>", open(f, encoding="utf-8", errors="replace").read())
    if m is None:
        return ""
    sek = m.group(1).strip()
    return sek if os.path.exists(os.path.join(korpus, "Section", sek + ".xml")) else ""


def ikat(prop, sisi_lama):
    """`TreatyIn.OLDDATA.X` → sisi; `TreatyIn.X` → akar (Section Old) / sisi (New)."""
    if prop.startswith("TreatyIn.OLDDATA."):
        return "sisi", prop[len("TreatyIn.OLDDATA."):]
    if prop.startswith("TreatyIn."):
        return ("akar" if sisi_lama else "sisi"), prop[len("TreatyIn."):]
    return None, prop


# Awalan halaman sesi yang kerangka bawa (lihat `Pembangkit.ikatan`).
SESI = ("SearchData.", "FlagExcel.")


class Pembangkit:
    def __init__(self, berkas, sisi_lama, korpus="", relatif=False):
        self.berkas, self.sisi_lama, self.dibuang = berkas, sisi_lama, []
        # Mode RELATIF = Section rincian baris: ikatannya `.X` atas halaman
        # baris. Ikatan absolut di dalamnya DIBUANG — konteks baris tidak
        # membawa halaman akar untuk ditulisi.
        self.korpus, self.relatif = korpus, relatif
        # Section rincian yang dirujuk grid di Section ini (antrean bangkit).
        self.rincian = []

    def ikatan(self, n, jenis, prop):
        """(dari, kunci) satu ikatan, atau None bila dibuang."""
        # ⭐ Halaman SESI Pega (`SearchData`, `FlagExcel` — Achievement rincian
        # Limits): bukan bagian dokumen `TreatyIn`, tetapi dibaca dan ditulis
        # Activity (`GetAchievement`, `Reset_DT`). Kuncinya utuh, `dari = sesi`.
        if prop.startswith(SESI):
            return "sesi", prop
        if self.relatif:
            if prop.startswith(".pyTemplate"):
                self.buang(n, jenis, f"kontrol tanpa properti: {prop!r}")
                return None
            if not prop.startswith("."):
                self.buang(n, jenis, f"ikatan absolut di Section rincian: {prop!r}")
                return None
            return "sisi", prop[1:]
        if not prop.startswith("TreatyIn."):
            self.buang(n, jenis, f"ikatan bukan halaman TreatyIn: {prop!r}")
            return None
        return ikat(prop, self.sisi_lama)

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
        if self.relatif:
            ik = self.ikatan(n, "grid", prop)
            if ik is None:
                return None
            dari, larik = ik
        else:
            dari, larik = ikat(prop, self.sisi_lama)
        baris = []
        for t in n.k("pyTable"):
            for r in t.k("pyRows"):
                for row in r.k("rowdata"):
                    baris.append([c for cs in row.k("pyCells") for c in cs.k("rowdata")])
        kepala = baris[0] if baris else []
        badan = baris[1] if len(baris) > 1 else []
        kol = {"kolom": [], "kunci": [], "lebar": [], "desimal": [], "format": [], "syaratSel": [], "atSel": [],
               "baca": [], "tombol": [], "tombolKepala": [], "pilihan": [], "aksiUbah": []}
        for i, c in enumerate(badan):
            nilai, fmt = c.v("pyValue"), c.v("pyFormat")
            sy = syarat_sendiri(c)
            if any(MATI.search(x) for x in sy):
                self.buang(c, "kolom", f"penjaga mati: {sy}")
                continue
            h = kepala[i] if i < len(kepala) else None
            hk = h is not None and (h.v("pyFormat") == "pxButton" or "pyTemplateButton" in h.v("pyValue"))
            hs = syarat_sendiri(h) if hk else []
            kol["tombolKepala"].append(tombol(h, hs) if hk and not any(MATI.search(x) for x in hs) else None)
            if fmt == "pxButton" or "pyTemplateButton" in nilai:
                # Kolom tombol baris (Delete/Hapus) — kuncinya kosong.
                kol["kolom"].append("")
                kol["kunci"].append("")
                kol["baca"].append("selalu")
                kol["tombol"].append(tombol(c, sy))
            else:
                kol["kolom"].append(h.v("pyValue") if h is not None and not hk else "")
                kol["kunci"].append(nilai[1:] if nilai.startswith(".") else nilai)
                kol["baca"].append(kunci_baca(c))
                kol["tombol"].append(None)
            kol["pilihan"].append(sumber_pilihan(c))
            kol["aksiUbah"].append(aksi_ubah(c) or None)
            kol["lebar"].append(int(c.v("pyWidth") or 0))
            kol["desimal"].append(desimal(c))
            kol["format"].append(fmt)
            kol["syaratSel"].append(sy[0] if sy else None)
            kol["atSel"].append(c.off)
        g = {"t": "grid", "at": n.off, "prop": prop, "dari": dari, "larik": larik, "syarat": syarat, **kol}
        tmpl = n.v("pyGridTemplateName")
        if tmpl:
            g["templatBaris"] = tmpl
        # Mode sunting dan flow action-nya ada di `pyGridProps` grid ini.
        gp = (n.k("pyGridProps") or [None])[0]
        aksi = gp.v("pyEditAction") if gp is not None else ""
        if gp is not None and gp.v("pyEditingMode") == "expandPane" and aksi:
            sek = seksi_flowaction(self.korpus, aksi)
            if sek:
                g["rincian"] = sek
                self.rincian.append(sek)
            else:
                # Flow action (atau Section-nya) tidak ada di ekspor — panel
                # rincian TIDAK dikarang; layar menyebut namanya.
                g["rincianHilang"] = aksi
        return g

    def sel(self, c, syarat):
        nilai, fmt, typ = c.v("pyValue"), c.v("pyFormat"), c.v("pyType")
        label = c.v("pyLabelFieldValue")
        if fmt == "pxButton" or "pyTemplateButton" in nilai:
            return tombol(c, syarat)
        if typ == "LABEL":
            if not nilai:
                return None
            return {"t": "teks", "at": c.off, "teks": nilai, "syarat": syarat}
        ik = self.ikatan(c, "medan", nilai)
        if ik is None:
            return None
        dari, kunci = ik
        cap = [x.txt.strip() for x in c.iter("pyCheckboxCaption") if x.txt.strip()]
        m = {
            "t": "medan", "at": c.off, "label": "" if label in LABEL_KOSONG else label,
            "dari": dari, "kunci": kunci, "format": fmt, "desimal": desimal(c), "syarat": syarat,
        }
        if cap:
            m["caption"] = cap[0]
        b = kunci_baca(c)
        if b is not None:
            m["baca"] = b
        sp = sumber_pilihan(c)
        if sp is not None:
            m["pilihan"] = sp
        au = aksi_ubah(c)
        if au:
            m["aksiUbah"] = au
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
                    g = self.grid(c, s)
                    if g is not None:
                        out.append(g)
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
                # ⛔ `pyIncludeHeader = false` menyembunyikan kepala berjudul
                # (`hiddden` tab Share, `hidden, reference` Section Share).
                # BUKAN `pyContainerFormat = NOHEADER`: blok `Exclusions` dan
                # `Summarry of RNM Share` ber-NOHEADER, dan judulnya TAMPIL di
                # gambar Pega 34 dan 39 — keduanya `pyIncludeHeader = true`.
                tampil = judul if (judul not in PENANDA and kepala in ("BAR", "TABBED")
                                   and c.v("pyIncludeHeader") != "false") else ""
                blok = {"t": "blok", "at": c.off, "judul": tampil, "syarat": s, "anak": anak}
                # ⭐ 8 Oktober 2026 — layout group `pyHeaderType = TABBED`:
                # blok-blok seperti ini yang BERURUTAN adalah satu strip tab
                # (mis. sebelas tab `DetailLimits`: Event Limits … Achievement).
                if kepala == "TABBED":
                    blok["tab"] = True
                out.append(blok)
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
    semua = []
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
        # ⭐ 7 Oktober 2026 — cabang ADJUST PREMIUM: `TreatyInNONProportional`
        # @566980 menampilkan Section ini bila `ProportionType='NonProportional'
        # && EDMState=3` (dan `TreatyInTabsNonProportional` bila `EDMState!=3`).
        # Halaman `TreatyIn.ActualValue` hidup di sini. `Actual Retro` ikut
        # dibangkitkan sebagai wadah (isinya Retro, §17) — layar
        # menyembunyikannya bersama tab Retro lain.
        ("TreatyInTabsNonProportionalAdjustPremi", False): [
            "Actual GNPI", "Actual Limits", "Actual Share", "Actual Retro", "Premium Adjustment",
            "Information & Submit"],
    }
    for (nama, lama), daftar in TAB.items():
        r = muat(sec(nama))
        p = Pembangkit(nama, lama, korpus)
        semua.append(p)
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
        # Isi tab cabang Adjust Premium (Actual Retro TIDAK — §17).
        ("TreatyInActualLimits", False), ("TreatyInActualShare", False), ("TreatyInActualSumary", False),
    ]:
        r = muat(sec(nama))
        p = Pembangkit(nama, lama, korpus)
        semua.append(p)
        include[nama] = p.jalan(r)
        dibuang += p.dibuang
    kurs = {}
    for nama, lama, prop in [("TreatyInNONProportionalOldData", True, "TreatyIn.OLDDATA.CurrencyList"),
                             ("TreatyInNONProportional", False, "TreatyIn.CurrencyList")]:
        r = muat(sec(nama))
        p = Pembangkit(nama, lama, korpus)
        semua.append(p)
        for e in r.iter("pyPageListProperty"):
            if e.txt.strip() == prop:
                kurs["lama" if lama else "baru"] = p.grid(e.par, [])
                break
        dibuang += p.dibuang
    # ⭐ Section rincian baris — termasuk rincian DI DALAM rincian (CoBList
    # di dalam Layers), sampai antreannya habis.
    rincian = {}
    antre = [x for p in semua for x in p.rincian]
    while antre:
        nama = antre.pop(0)
        if nama in rincian:
            continue
        p = Pembangkit(nama, False, korpus, relatif=True)
        rincian[nama] = p.jalan(muat(sec(nama)))
        dibuang += p.dibuang
        antre += p.rincian
    return kerangka, include, kurs, dict(sorted(rincian.items())), dibuang


def main():
    korpus = sys.argv[1]
    kerangka, include, kurs, rincian, dibuang = bangkit(korpus)
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
        f.write("// Section rincian baris (`expandPane`) — ikatan RELATIF atas halaman baris.\n")
        f.write(f"export const KERANGKA_RINCIAN: Readonly<Record<string, readonly ButirKerangka[]>> = {j(rincian)}\n\n")
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
    print(f"{len(kerangka)} tab, {len(include)} include, {len(rincian)} rincian, {len(dibuang)} butir dibuang, "
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
