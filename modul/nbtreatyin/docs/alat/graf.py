"""Graf panggilan antar-rule korpus NB Treaty In dan jangkauan dari titik masuk nyata.

HANYA MEMBACA. Rujukan diambil dari dua sumber yang digabung:

  1. POSISI TERSTRUKTUR per jenis rule (korpus.py):
       Flow            implementasi shape, connector bersyarat Action (FlowAction) / When
       FlowAction      section, pre/post activity, pre/post data transform, local action
       Section/Harness section disertakan, aksi medan/button (activity, harness, local
                       action, section refresh), sumber daftar, nama When pada syarat
       Activity        target Call, RequestType RDB-List, Apply-DataTransform,
                       Property-Map-DecisionTable, Param.pyReportName, nama When
       DataTransform   nama When pada langkah
  2. INDEKS RUJUKAN PEGA (pxRuleReferences) dan tag bernama-rule (pyDeferLoadActivity,
     pyEditAction, pyGridRDName, pyContainerVisibleWhen, pyModelName, pyPreDataTransform, ...)
     - salinan section tertanam (pyIncludedRuleXML) DIBUANG lebih dulu.

RESOLUSI SADAR-KELAS. Pega mencari rule bernama X mulai dari kelas konteks lalu naik ke
leluhurnya (pewarisan pola: ASM-FW-GISFW-Data-PolicyTreatyIn -> ASM-FW-GISFW-Data -> ...
-> @baseclass). Bila kelas konteks rujukan diketahui dan rule di korpus berkelas yang
BUKAN leluhur kelas itu, rujukan TIDAK disambung - dicatat sebagai `bedaKelas`.
Contoh nyata: ListSuggest memanggil Protection_Act dengan kelas konteks
Data-PolicyTreatyIn, sedangkan Protection_Act di korpus berkelas ASM-FW-GISFW-Work.
"""
from __future__ import annotations

import collections
import re

TITIK_MASUK = [("Harness", "SFAPortalOpportunities"), ("Flow", "InputRealizationTreatyIn")]

# Sambungan yang DIPUTUS, masing-masing dengan bukti yang dapat diuji ulang.
PUTUS = {
    (("Section", "ListSuggest"), ("Activity", "Protection_Act")):
        "Protection_Act di korpus ini md5 8d6fd933b1e10bbbdab2945ffa681ff6 = varian NB FacIn/RNW Fac In "
        "(sidik-jari-Protection_Act-sebelum.json), kelas ASM-FW-GISFW-Work. ListSuggest memanggilnya "
        "dengan pyActivityClass ASM-FW-GISFW-Data-PolicyTreatyIn - Work bukan leluhur kelas itu, "
        "sehingga Pega memilih varian Data-PolicyTreatyIn (157.304 B, 16 langkah) yang TIDAK ada di "
        "korpus ini. Keputusan work owner 23-09-2026 butir 6: rantai spreading milik Fac In.",
}

# Kelas kerja turunan berarah (directed inheritance) yang tidak tampak dari nama.
_INDUK_BERARAH = {"ASM-FW-GISFW-Work-NB": "ASM-FW-GISFW-Work"}


def leluhur(kelas):
    """Kelas itu sendiri dan seluruh leluhur pola-nya, plus @baseclass."""
    out, k = [], kelas
    while k:
        out.append(k)
        if k in _INDUK_BERARAH:
            k = _INDUK_BERARAH[k]
            continue
        k = k.rsplit("-", 1)[0] if "-" in k else ""
    out.append("@baseclass")
    return out


def cocok_kelas(kelas_rule, kelas_konteks):
    if not kelas_konteks or not kelas_rule:
        return True
    lk = [x.upper() for x in leluhur(kelas_konteks)]
    return kelas_rule.upper() in lk


def _langkah_rata(ls):
    for s in ls:
        yield s
        yield from _langkah_rata(s.get("anak", []))


class Graf:
    def __init__(self, data):
        self.data = data
        self.per_kunci = {(d["jenis"], d["nama"]): d for d in data}
        self.per_nama = collections.defaultdict(list)
        for d in data:
            self.per_nama[d["nama"].lower()].append((d["jenis"], d["nama"]))
        self.when = {d["nama"] for d in data if d["jenis"] == "When"}
        self.keluar = collections.defaultdict(set)     # kunci -> {(kunci, alasan)}
        self.hilang = collections.defaultdict(set)     # dirujuk, nama tak ada di korpus
        self.beda_kelas = collections.defaultdict(set)  # nama ada, kelas bukan leluhur konteks
        for d in data:
            getattr(self, "_" + d["jenis"].lower())(d)
            for r in d.get("rujukan", []):
                self._ke(d, (r["jenis"],), r["nama"], f"{r['lewat']}", r["kelas"])

    # ------------------------------------------------------------- bantu
    def _ke(self, d, jenis_boleh, nama, alasan, kelas=""):
        if not nama:
            return
        nama = nama.strip().strip('"')
        if "." in nama and not nama.startswith("Rule-Obj-") and "Activity" in jenis_boleh:
            kelas, nama = nama.rsplit(".", 1)
        asal = (d["jenis"], d["nama"])
        kandidat = [k for k in self.per_nama.get(nama.lower(), []) if k[0] in jenis_boleh]
        if not kandidat:
            self.hilang[asal].add(("/".join(jenis_boleh), nama, alasan, kelas))
            return
        # Disambung BERBASIS NAMA - sama seperti closure ekspor Pega. Selisih kelas hanya
        # dicatat (beda_kelas); pemutusan dilakukan lewat PUTUS, dan hanya dengan bukti.
        for k in kandidat:
            if k != asal:
                self.keluar[asal].add((k, alasan))
            if not cocok_kelas(self.per_kunci[k]["kelas"], kelas):
                self.beda_kelas[asal].add((k, alasan, kelas, self.per_kunci[k]["kelas"]))

    def _when_dalam(self, d, teks, alasan):
        for w in self.when:
            if re.search(r"(?<![\w.])" + re.escape(w) + r"(?![\w(])", teks or ""):
                self.keluar[(d["jenis"], d["nama"])].add((("When", w), alasan))

    # ------------------------------------------------------------- per jenis
    def _flow(self, d):
        for s in d["bentuk"]:
            if s["implementasi"] and s["kelasKeputusan"] == "Rule-Declare-DecisionTable":
                self._ke(d, ("DecisionTable",), s["implementasi"], f"shape {s['id']}", d["kelas"])
            elif "Utility" in s["jenis"]:
                self._ke(d, ("Activity",), s["implementasi"], f"shape {s['id']}", d["kelas"])
        for c in d["sambung"]:
            if c["syarat"] == "Action":
                self._ke(d, ("FlowAction",), c["ekspresi"], f"connector {c['dari']}->{c['ke']}")
            elif c["syarat"] == "When":
                self._ke(d, ("When",), c["ekspresi"], f"connector {c['dari']}->{c['ke']}")

    def _flowaction(self, d):
        self._ke(d, ("Section",), d.get("pySectionReference"), "section")
        for k in ("pyPreProcessingActivity", "pyActionActivity", "pyPostProcessingActivity", "pyLocalActionActivity"):
            self._ke(d, ("Activity",), d.get(k), k, d["kelas"])
        for k in ("pyPreProcessingTransformRule", "pyActionTransformRule"):
            self._ke(d, ("DataTransform",), d.get(k), k)
        self._section(d)

    def _harness(self, d):
        self._section(d)

    def _section(self, d):
        for s in d.get("sertakan", []):
            if s != d["nama"]:
                self._ke(d, ("Section",), s, "sertakan")
        for b in d.get("badan", []):
            for k in ("tampilWhen", "kunciWhen"):
                self._when_dalam(d, b[k], "syarat layout")
        for c in d.get("sel", []):
            for k in ("tampilWhen", "kunciWhen", "wajibWhen", "nonaktifWhen"):
                self._when_dalam(d, c[k], f"syarat medan {c['nilai']}")
            for a in c["aksi"]:
                api = a["api"]
                kls = api.get("pyActivityClass", "")
                medan = f"aksi {a['aksi']} medan {c['nilai']}"
                self._ke(d, ("Activity",), api.get("pyActivity"), medan, kls)
                self._ke(d, ("Harness", "Section"), api.get("pyHarnessName"), medan)
                self._ke(d, ("FlowAction",), api.get("pyLocalAction"), medan)
                self._ke(d, ("DataTransform",), api.get("pyDataTransform"), medan)
                self._ke(d, ("Section",), api.get("pySection"), medan)
            src = c["sumberDaftar"]
            for k in ("pySourceName", "pyReportDefinition", "pyReportDefName", "pyRDName", "pyReportName"):
                self._ke(d, ("ReportDefinition", "RDBList", "Activity"), src.get(k), f"sumber daftar medan {c['nilai']}")

    def _activity(self, d):
        for s in _langkah_rata(d["langkah"]):
            j = s["jenis"].lower()
            kls = s["kelas"] if s["kelas"] and not s["kelas"].startswith("Rule-") else ""
            if j == "call" and s["target"]:
                tgt = s["target"]
                if tgt in ("pxShowReport", "pxRetrieveReportData", "Rule-Obj-Report-Definition.pxRetrieveReportData"):
                    rn = s["panggil"].get("pyReportName")
                    if rn:
                        self._ke(d, ("ReportDefinition",), rn, f"langkah {s['no']} {tgt}")
                else:
                    self._ke(d, ("Activity",), tgt, f"langkah {s['no']} Call", kls)
            elif j == "rdb-list":
                self._ke(d, ("RDBList",), s["panggil"].get("RequestType"), f"langkah {s['no']} RDB-List")
            elif j == "apply-datatransform":
                self._ke(d, ("DataTransform",), s["panggil"].get("DataTransform"), f"langkah {s['no']} Apply-DataTransform")
            elif j == "property-map-decisiontable":
                for v in s["panggil"].values():
                    self._ke(d, ("DecisionTable",), v, f"langkah {s['no']} Property-Map-DecisionTable")
            for n, v in s["param"]:
                if n.lower() == "param.pyreportname":
                    self._ke(d, ("ReportDefinition",), v, f"langkah {s['no']} Param.pyReportName")
            for p in s["pre"]:
                self._when_dalam(d, p["when"], f"syarat langkah {s['no']}")
            for p in s["trans"]:
                self._when_dalam(d, p["when"], f"transisi langkah {s['no']}")

    def _datatransform(self, d):
        for s in _langkah_rata(d["langkah"]):
            for teks in (s["sasaran"], s["sumber"], *s["lain"].values()):
                self._when_dalam(d, teks, f"langkah {s['no']}")

    def _when(self, d):
        for k in d["kondisi"]:
            self._when_dalam(d, k["teks"], "kondisi")

    def _decisiontable(self, d):
        pass

    def _rdblist(self, d):
        pass

    def _reportdefinition(self, d):
        pass

    # ------------------------------------------------------------- jangkauan
    def masuk(self):
        m = collections.defaultdict(set)
        for a, ks in self.keluar.items():
            for b, alasan in ks:
                m[b].add((a, alasan))
        return m

    def terjangkau(self, titik=TITIK_MASUK, buang=()):
        seen, antri = set(), [k for k in titik if k in self.per_kunci]
        while antri:
            k = antri.pop()
            if k in seen or k in buang:
                continue
            seen.add(k)
            antri.extend(b for b, _ in self.keluar.get(k, ()) if (k, b) not in PUTUS)
        return seen
