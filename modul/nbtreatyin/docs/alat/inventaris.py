"""Penyusun docs/INVENTARIS-XML.md dari korpus XML Pega NB Treaty In - HANYA MEMBACA korpus.

    python korpus.py korpus.json          # baca 278 berkas XML
    python inventaris.py korpus.json      # tulis ../INVENTARIS-XML.md

Status "dibangun di / tidak dibangun" dibaca dari status.json di folder ini (diisi tangan,
satu entri per rule). Rule yang TIDAK terjangkau diberi status otomatis beserta buktinya.

KERAHASIAAN (AC 60, 61): nama orang tidak pernah ditulis. Penyamaran berbasis POLA, sehingga
daftar nama itu sendiri tidak pernah tersimpan di repositori:
  - literal yang dibandingkan dengan identitas operator (pyUserIdentifier, pxInsName,
    pxCreateOperator, pyUserName) -> <ID-operator-N>
  - teks status "NB IS IN <NAMA>'S INBOX" -> <nama-N>
  - literal yang dibandingkan dengan medan bernama-orang (MarketingName, ClientName,
    InsuredName, OperatorName, pxCreateOpName, ...) -> <nama-N>
  - nomor polis berbentuk RNM-<tipe>dd.mm.yyyy.nnnnn -> <nomor-polis-N>
  - alamat surel dan URL -> <alamat-disunting>
"""
from __future__ import annotations

import collections
import json
import os
import re
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from graf import PUTUS, TITIK_MASUK, Graf, leluhur  # noqa: E402

DIR = os.path.dirname(os.path.abspath(__file__))
TUJUAN = os.path.join(DIR, "..", "INVENTARIS-XML.md")

# ------------------------------------------------------------------ penyamaran

_ID_OP = re.compile(
    r"((?:pyUserIdentifier|pxInsName|pxCreateOperator|pyUserName)\s*(?:==|!=|=|,\s*\"[=!]+\"\s*,)\s*)(['\"])([A-Za-z][A-Za-z0-9_.]*)\2")
_NB_INBOX = re.compile(r"NB IS IN ([A-Z][A-Z ]*)'S INBOX")
_MEDAN_NAMA = re.compile(
    r"((?:MarketingName|ClientName|InsuredName|OperatorName|CreateOpName|pxCreateOpName|pxUpdateOpName"
    r"|PICSuggest|UserName|Username)\s*(?:==|!=|=)\s*)(['\"])([^'\"]+)\2")
_NOPOLIS = re.compile(r"\b[A-Z]{2,5}-[A-Z]{1,3}\d{0,2}\.\d{2}\.\d{4}\.\d{4,6}\b")
_SUREL = re.compile(r"[\w.+-]+@[\w-]+\.[\w.]+")
_URL = re.compile(r"https?://[^\s\"'<>]+")


class Penyamar:
    def __init__(self):
        self.op, self.nama, self.nopol = {}, {}, {}

    def __call__(self, s):
        if not isinstance(s, str) or not s:
            return s
        s = _ID_OP.sub(lambda m: f"{m.group(1)}{m.group(2)}{self._op(m.group(3))}{m.group(2)}", s)
        s = _NB_INBOX.sub(lambda m: f"NB IS IN {self._nm(m.group(1))}'S INBOX", s)
        s = _MEDAN_NAMA.sub(lambda m: f"{m.group(1)}{m.group(2)}{self._nm(m.group(3))}{m.group(2)}", s)
        s = _NOPOLIS.sub(lambda m: self._np(m.group(0)), s)
        s = _SUREL.sub("<alamat-disunting>", s)
        s = _URL.sub("<alamat-disunting>", s)
        return s

    def _op(self, v):
        if v not in self.op:
            self.op[v] = f"<ID-operator-{len(self.op) + 1}>"
        return self.op[v]

    def _np(self, v):
        if v not in self.nopol:
            self.nopol[v] = f"<nomor-polis-{len(self.nopol) + 1}>"
        return self.nopol[v]

    def _nm(self, v):
        if v not in self.nama:
            self.nama[v] = f"<nama-{len(self.nama) + 1}>"
        return self.nama[v]


SAMAR = Penyamar()


def sel(s, maks=None):
    """Isi sel tabel markdown: disamarkan, pipa di-escape, baris baru jadi spasi."""
    s = SAMAR(s if s is not None else "")
    s = s.replace("\r", " ").replace("\n", " ").replace("|", "\\|").strip()
    if maks and len(s) > maks:
        s = s[:maks] + " …"
    return s


def kode(s):
    return "`" + sel(s).replace("`", "'") + "`" if s else ""


# ------------------------------------------------------------------ makna kode aksi langkah

AKSI_WHEN = {"1": "lompat ke label", "2": "lanjut", "3": "lewati langkah", "4": "keluar iterasi",
             "5": "lewati syarat berikut", "6": "keluar activity"}


def aksi(kd, label):
    if not kd:
        return ""
    t = AKSI_WHEN.get(kd, f"kode {kd}")
    return f"{t} `{label}`" if label else t


def _rata(ls):
    for s in ls:
        yield s
        yield from _rata(s.get("anak", []))


# ------------------------------------------------------------------ status

def muat_status():
    p = os.path.join(DIR, "status.json")
    if os.path.exists(p):
        with open(p, encoding="utf-8") as fh:
            return json.load(fh)
    return {}


def status_otomatis(g, k, R, masuk):
    """Status untuk rule yang TIDAK terjangkau - dengan bukti."""
    pemanggil = sorted({a for a, _ in masuk.get(k, ())})
    if not pemanggil:
        return "tidak dibangun — **yatim**: nol pemanggil di seluruh korpus (graf.py)"
    hidup = [a for a in pemanggil if a in R]
    if hidup and all((a, k) in PUTUS for a in hidup):
        return ("tidak dibangun — **varian Fac In**: satu-satunya pemanggil hidup diputus "
                "dengan bukti (bab 1.2)")
    return ("tidak dibangun — **tidak terjangkau**: hanya dipanggil oleh "
            + ", ".join(f"`{a[1]}`" for a in pemanggil[:6])
            + (" …" if len(pemanggil) > 6 else "") + " — seluruhnya tidak terjangkau")


# ------------------------------------------------------------------ penulis

class Tulis:
    def __init__(self):
        self.b = []

    def __call__(self, *baris):
        self.b.extend(baris)

    def teks(self):
        return "\n".join(self.b).rstrip() + "\n"


def tulis(data):
    g = Graf(data)
    R = g.terjangkau()
    masuk = g.masuk()
    status = muat_status()
    K = g.per_kunci
    w = Tulis()
    urut = sorted(data, key=lambda d: (d["jenis"], d["nama"].lower()))
    n_jenis = collections.Counter(d["jenis"] for d in data)
    n_hidup = collections.Counter(k[0] for k in R)

    w("# Inventaris XML — NB Treaty In",
      "",
      "> **Berkas bangkitan.** Ditulis `docs/alat/inventaris.py` dari korpus XML Pega "
      "(READ-ONLY). Jangan disunting tangan selain lewat `docs/alat/status.json`; jalankan ulang:",
      ">",
      "> ```",
      "> cd modul/nbtreatyin/docs/alat",
      "> python korpus.py korpus.json && python inventaris.py korpus.json",
      "> ```",
      ">",
      "> Jalur korpus: env var `NBTREATYIN_KORPUS` (bawaan: `D:/NUSARE DEV/NusantaraRe/NB Treaty In (Done)`).",
      "> Nama orang, surel, dan alamat **disamarkan berbasis pola** (AC 60, 61) — lihat kepala `inventaris.py`.",
      "",
      "## Isi",
      "",
      "1. Jangkauan dan bukti lingkup",
      "2. Tabel seluruh berkas (278)",
      "3. Menu dan titik masuk",
      "4. Alur `Flow/InputRealizationTreatyIn`",
      "5. Layar — Harness, FlowAction, Section",
      "6. Activity terjangkau — langkah demi langkah",
      "7. Data Transform terjangkau",
      "8. RDBList — naskah SQL",
      "9. Report Definition",
      "10. Decision Table",
      "11. When — syarat yang DIJALANKAN",
      "12. Rujukan ke rule yang tidak ada di korpus, dan selisih kelas",
      "")

    # ------------------------------------------------------------- 1
    w("## 1 · Jangkauan dan bukti lingkup", "",
      "### 1.1 Cacah", "",
      "| Jenis | Berkas | Terjangkau | Tidak terjangkau |", "| --- | ---: | ---: | ---: |")
    for j in sorted(n_jenis):
        w(f"| `{j}` | {n_jenis[j]} | {n_hidup.get(j, 0)} | {n_jenis[j] - n_hidup.get(j, 0)} |")
    w(f"| **Jumlah** | **{len(data)}** | **{len(R)}** | **{len(data) - len(R)}** |", "",
      "**Titik masuk nyata:** " + " · ".join(f"`{j}/{n}`" for j, n in TITIK_MASUK)
      + " — harness portal (akar `Struktur_MenuNBTreatyIn.xlsx`) dan flow yang dibuat tombol *Create* "
      "portal itu. Jangkauan diikuti lewat rujukan terstruktur dan indeks `pxRuleReferences` Pega "
      "(`graf.py`).", "",
      "### 1.2 Satu sambungan diputus — dengan bukti", "")
    for (a, b), bukti in PUTUS.items():
        w(f"- `{a[0]}/{a[1]}` → `{b[0]}/{b[1]}`: {bukti}")
    w("", "Akibatnya rantai spreading/premi **Fac In** tidak terjangkau. Rule yang hanya hidup lewat "
      "sambungan ini tercatat di bab 2 dengan status *varian Fac In* atau *tidak terjangkau*.", "")
    lama = os.path.join(DIR, "..", "jangkauan.json")
    if os.path.exists(lama):
        with open(lama, encoding="utf-8") as fh:
            j = json.load(fh)
        mine = {k[1] for k in R if k[0] == "Activity"}
        dok = set(j.get("terjangkau", []))
        w("### 1.3 Uji silang dengan `jangkauan.json` (tim migrasi)", "",
          f"- Activity terjangkau menurut inventaris ini: **{len(mine)}**; menurut `jangkauan.json`: **{len(dok)}**.",
          f"- Ada di `jangkauan.json`, tidak di sini: {', '.join(f'`{x}`' for x in sorted(dok - mine)) or 'nihil'}"
          " — `Protection_Act` varian Treaty, yang berkasnya tidak ada di korpus ini (bab 1.2).",
          f"- Ada di sini, tidak di `jangkauan.json`: {', '.join(f'`{x}`' for x in sorted(mine - dok)) or 'nihil'}.",
          "")

    # ------------------------------------------------------------- 2
    w("## 2 · Tabel seluruh berkas", "",
      "Kolom *Status*: **dibangun di** `<berkas>` · **diganti** (keputusan work owner) · "
      "**tidak dibangun** + alasan + bukti. Sumber: `alat/status.json` untuk rule terjangkau; "
      "otomatis untuk yang tidak terjangkau.", "",
      "| # | Jenis | Nama | Kelas | Dipanggil dari | Memanggil | Status |",
      "| ---: | --- | --- | --- | --- | --- | --- |")
    for i, d in enumerate(urut, 1):
        k = (d["jenis"], d["nama"])
        dari = sorted({a for a, _ in masuk.get(k, ())})
        ke = sorted({b for b, _ in g.keluar.get(k, ())})
        st = status.get(f"{d['jenis']}/{d['nama']}")
        if st:
            s = st if isinstance(st, str) else st.get("status", "")
        elif k in R:
            s = "⚠️ **BELUM DIPETAKAN**"
        else:
            s = status_otomatis(g, k, R, masuk)
        hidup = "" if k in R else " ⛔"
        w(f"| {i} | {d['jenis']} | `{d['nama']}`{hidup} | {sel(d['kelas'])} | "
          f"{sel(', '.join(f'{a[1]}' for a in dari), 160)} | {sel(', '.join(f'{b[1]}' for b in ke), 160)} | {sel(s)} |")
    w("", "⛔ = tidak terjangkau dari titik masuk.", "")

    # ------------------------------------------------------------- 3 menu
    w("## 3 · Menu dan titik masuk", "")
    menu = os.path.join(DIR, "menu.json")
    if os.path.exists(menu):
        with open(menu, encoding="utf-8") as fh:
            m = json.load(fh)
        w(f"Dari `Struktur_MenuNBTreatyIn.xlsx` (sheet `Struktur`, {m['baris']} baris, kedalaman {m['kedalaman']}). "
          "Tingkat 1–3 pohon pemanggilan:", "",
          "| No | Jenis rule | Nama | Dipanggil dari |", "| --- | --- | --- | --- |")
        for r in m["atas"]:
            w(f"| {sel(r['no'])} | {sel(r['jenis'])} | `{sel(r['nama'])}` | {sel(r['dari'])} |")
        w("", m.get("catatan", ""), "")
    else:
        w("⚠️ `alat/menu.json` belum dibuat — jalankan `python menu.py`.", "")

    # ------------------------------------------------------------- 4 alur
    fl = K[("Flow", "InputRealizationTreatyIn")]
    w("## 4 · Alur `Flow/InputRealizationTreatyIn`", "",
      f"Kelas `{fl['kelas']}`. **{len(fl['bentuk'])} shape**, **{len(fl['sambung'])} connector**.", "",
      "### 4.1 Shape", "",
      "| Id | Jenis | Nama | Implementasi | Workbasket | Dimasuki connector |",
      "| --- | --- | --- | --- | --- | ---: |")
    masuk_shape = collections.Counter(c["ke"] for c in fl["sambung"])
    for s in fl["bentuk"]:
        w(f"| `{s['id']}` | {sel(s['jenis'].replace('Data-MO-', ''))} | {sel(s['nama'])} | "
          f"{kode(s['implementasi'])} {sel(s['kelasKeputusan'].replace('Data-MO-Gateway-', ''))} | "
          f"{kode(s['workbasket'])} | {masuk_shape.get(s['id'], 0)} |")
    w("", "### 4.2 Connector", "",
      "| Dari | Ke | Nama | Syarat | Ekspresi | Tugas properti |", "| --- | --- | --- | --- | --- | --- |")
    for c in sorted(fl["sambung"], key=lambda c: (c["dari"], c["ke"])):
        tg = "; ".join(f"{a} = {b}" for a, b in c["tugas"])
        w(f"| `{c['dari']}` | `{c['ke']}` | {sel(c['nama'])} | {sel(c['syarat'])} | {kode(c['ekspresi'])} | {sel(tg)} |")
    tak = [s["id"] for s in fl["bentuk"] if masuk_shape.get(s["id"], 0) == 0 and not s["id"].startswith("Start")]
    w("", f"**Shape tanpa connector masuk (tidak pernah dicapai):** {', '.join(f'`{x}`' for x in tak) or 'nihil'}.", "")

    # ------------------------------------------------------------- 5 layar
    w("## 5 · Layar — Harness, FlowAction, Section", "",
      "Setiap medan milik rule itu sendiri (salinan section tertanam `pyIncludedRuleXML` dilompati). "
      "*Wajib*: `pyRequired`/mode; *Kunci*: `pyReadOnly` atau syarat baca-saja; *Tampil*: `pyVisible` + syarat; "
      "*Aksi*: event → aksi → activity/harness/local action yang dipanggil.", "")
    for j in ("Harness", "FlowAction", "Section"):
        for d in [x for x in urut if x["jenis"] == j]:
            k = (d["jenis"], d["nama"])
            w(f"### 5.{j[0]} `{j}/{d['nama']}`" + ("" if k in R else " ⛔ tidak terjangkau"), "",
              f"Kelas `{d['kelas']}`" + (f" · halaman `{', '.join(d.get('halaman', []))}`" if d.get("halaman") else ""))
            if j == "FlowAction":
                for kk in ("pySectionReference", "pyPreProcessingActivity", "pyPreProcessingTransformRule",
                           "pyActionTransformRule", "pyActionActivity", "pyPostProcessingActivity",
                           "pyLocalActionActivity", "pySubmitLabel", "pyconfirmchoice"):
                    if d.get(kk):
                        w(f"- `{kk}` = `{sel(d[kk])}`")
            if d.get("sertakan"):
                w(f"- menyertakan: {', '.join(f'`{x}`' for x in d['sertakan'])}")
            for b in d.get("badan", []):
                w(f"- layout `{sel(b['sertakan'])}` tampil `{sel(b['tampil'])}` {kode(b['tampilWhen'])} kunci {kode(b['kunciWhen'])}")
            sel_ = d.get("sel", [])
            if sel_:
                w("", "| Medan | Kontrol | Label | Wajib | Kunci | Tampil | Aksi |",
                  "| --- | --- | --- | --- | --- | --- | --- |")
                for c in sel_:
                    wajib = c["wajib"] if c["wajib"] == "true" else ""
                    if c["wajibMode"] == "true":
                        wajib = "true"
                    if c["wajibWhen"]:
                        wajib += f" `{sel(c['wajibWhen'])}`"
                    kunci = c["kunci"] if c["kunci"] == "true" else ""
                    if c["kunciWhen"]:
                        kunci += f" `{sel(c['kunciWhen'])}`"
                    if c["nonaktif"] == "true" or c["nonaktifWhen"]:
                        kunci += f" nonaktif `{sel(c['nonaktifWhen'])}`"
                    tampil = "" if c["tampil"] in ("", "ALWAYS") else c["tampil"]
                    if c["tampilWhen"]:
                        if c["tampil"] in ("OTHER", "NOTBLANK", "NOTZERO"):
                            tampil += f" `{sel(c['tampilWhen'])}`"
                        else:
                            # pyVisible=ALWAYS: pyCondition hanya sisa, TIDAK dijalankan
                            tampil += f" selalu (sisa syarat tak berlaku: `{sel(c['tampilWhen'])}`)"
                    ak = []
                    for a in c["aksi"]:
                        api = a["api"]
                        tgt = api.get("pyActivity") or api.get("pyHarnessName") or api.get("pyLocalAction") \
                            or api.get("pySection") or api.get("pyDataTransform") or api.get("pyFunctionName") or ""
                        prm = ("(" + ", ".join(a["param"]) + ")") if a["param"] else ""
                        ak.append(f"{a['event'] or '-'}→{a['aksi']} {tgt}{prm}".strip())
                    ak = list(dict.fromkeys(ak))
                    src = c["sumberDaftar"].get("pySourceName") or ""
                    kontrol = c["kontrol"] or c["tipe"]
                    if src:
                        kontrol += f" ⟵ `{src}`"
                    if c["sertakan"]:
                        kontrol += f" ⟵ `{c['sertakan']}`"
                    w(f"| {kode(c['nilai'])} | {sel(kontrol)} | {sel(c['label'] or c['labelFor'], 60)} | {sel(wajib)} | "
                      f"{sel(kunci)} | {sel(tampil)} | {sel('; '.join(ak))} |")
            w("")

    # ------------------------------------------------------------- 6 activity
    w("## 6 · Activity terjangkau — langkah demi langkah", "",
      "Kode aksi syarat/transisi Pega dicantumkan apa adanya beserta tafsirannya: "
      + " · ".join(f"`{k}` {v}" for k, v in AKSI_WHEN.items())
      + ". Syarat langkah (`when`) hanya berlaku bila kotak *When* langkah itu aktif (`whenAktif=true`).", "")
    for d in [x for x in urut if x["jenis"] == "Activity" and (x["jenis"], x["nama"]) in R]:
        w(f"### 6 · `{d['nama']}`", "",
          f"Kelas `{d['kelas']}` · {d['byte']:,} B".replace(",", ".")
          + (f" · parameter: {', '.join(f'`{p}`' for p in d['parameter'])}" if d["parameter"] else "")
          + (f" · keterangan: {sel(d['ket'])}" if d["ket"] else ""), "")
        _tulis_langkah(w, d["langkah"], 0)
        w("")

    # ------------------------------------------------------------- 7 DT
    w("## 7 · Data Transform", "")
    for d in [x for x in urut if x["jenis"] == "DataTransform"]:
        k = (d["jenis"], d["nama"])
        w(f"### 7 · `{d['nama']}`" + ("" if k in R else " ⛔"), "", f"Kelas `{d['kelas']}`", "",
          "| No | Aksi | Sasaran | Sumber | Lain |", "| --- | --- | --- | --- | --- |")
        for s in _rata(d["langkah"]):
            lain = "; ".join(f"{a}={b}" for a, b in s["lain"].items())
            w(f"| {s['no']} | {sel(s['aksi'])} | {kode(s['sasaran'])} | {kode(s['sumber'])} | {sel(lain, 600)} |")
        w("")

    # ------------------------------------------------------------- 8 RDB
    w("## 8 · RDBList — naskah SQL", "")
    for d in [x for x in urut if x["jenis"] == "RDBList"]:
        k = (d["jenis"], d["nama"])
        pemakai = sorted({a[1] for a, _ in masuk.get(k, ())})
        w(f"### 8 · `{d['nama']}`" + ("" if k in R else " ⛔"), "",
          f"Kelas `{d['kelas']}` · akses `{d['akses']}` · tabel: {', '.join(f'`{x}`' for x in d['tabel']) or '—'}"
          + (f" · prosedur: {', '.join(f'`{x}`' for x in d['prosedur'])}" if d["prosedur"] else "")
          + f" · dipakai: {', '.join(f'`{x}`' for x in pemakai) or '—'}", "")
        for kk, v in d["sql"].items():
            w(f"`{kk}`:", "", "```sql", SAMAR(v), "```", "")

    # ------------------------------------------------------------- 9 RD
    w("## 9 · Report Definition", "", "| Nama | Kelas | Kolom | Filter | Dipakai |", "| --- | --- | --- | --- | --- |")
    for d in [x for x in urut if x["jenis"] == "ReportDefinition"]:
        k = (d["jenis"], d["nama"])
        pemakai = sorted({a[1] for a, _ in masuk.get(k, ())})
        w(f"| `{d['nama']}`{'' if k in R else ' ⛔'} | {sel(d['kelas'])} | {sel(', '.join(d['kolom']), 300)} | "
          f"{sel('; '.join(d['filter']), 300)} | {sel(', '.join(pemakai))} |")
    w("")

    # ------------------------------------------------------------- 10 DT tabel
    w("## 10 · Decision Table", "")
    for d in [x for x in urut if x["jenis"] == "DecisionTable"]:
        k = (d["jenis"], d["nama"])
        w(f"### 10 · `{d['nama']}`" + ("" if k in R else " ⛔"), "",
          f"Kelas `{d['kelas']}` · bawaan `{d['bawaan']}` ({d['aksiBawaan']}) · evaluasi semua baris: `{d['evaluasiSemua'] or 'no'}` "
          "(baris pertama yang cocok menang) · kolom: "
          + ", ".join(f"`{c['properti']}` ({c['tipe']}, `{c['operator']}`)" for c in d["kolom"]), "")
        hdr = " | ".join(f"`{c['properti']}`" for c in d["kolom"])
        w(f"| Baris | {hdr} | Hasil |", "| ---: |" + " --- |" * len(d["kolom"]) + " --- |")
        for i, h in enumerate(d["hasil"], 1):
            cells = []
            for c in d["kolom"]:
                v = c["sel"][i - 1] if i - 1 < len(c["sel"]) else None
                o = c["atau"].get(i)
                s_ = "*(apa saja)*" if v is None else f"`{v}`"
                if o:
                    s_ = ("" if v is None else s_ + " ATAU ") + "salah satu " + ", ".join(f"`{x}`" for x in o)
                cells.append(s_)
            w(f"| {i} | " + " | ".join(cells) + f" | `{h}` |")
        w("", "*Sel bertanda kutip ditiru apa adanya — termasuk spasi di dalam kutip.* "
          "`pyRowNum` daftar-ATAU berbasis nol, dipetakan ke baris = `pyRowNum + 1`.", "")

    # ------------------------------------------------------------- 11 When
    w("## 11 · When — syarat yang DIJALANKAN", "",
      "Yang dijalankan adalah `pyLogic` + `pyCondition` (P23). Kolom *penampil* hanya untuk "
      "perbandingan; bila berbeda, yang dijalankan benar.", "",
      "| Nama | Kelas | Logika dijalankan | Kondisi | Penampil (tidak dijalankan) | Dipakai |",
      "| --- | --- | --- | --- | --- | --- |")
    for d in [x for x in urut if x["jenis"] == "When"]:
        k = (d["jenis"], d["nama"])
        kond = "; ".join(f"{c['label']}: {c['teks']}" for c in d["kondisi"])
        pen = f"{d['logikaPenampil']}: " + "; ".join(d["kondisiPenampil"])
        pemakai = sorted({a[1] for a, _ in masuk.get(k, ())})
        w(f"| `{d['nama']}`{'' if k in R else ' ⛔'} | {sel(d['kelas'])} | `{sel(d['logika'])}` | {sel(kond, 600)} | "
          f"{sel(pen, 200)} | {sel(', '.join(pemakai), 200)} |")
    w("")

    # ------------------------------------------------------------- 12 hilang
    w("## 12 · Rujukan ke rule yang tidak ada di korpus, dan selisih kelas", "",
      "Hanya dari rule yang terjangkau. Rule bawaan Pega (`py*`, `pz*`, `px*`, `Always`, `Never`) "
      "tidak dicantumkan.", "",
      "| Dari | Jenis | Nama dirujuk | Lewat | Kelas konteks |", "| --- | --- | --- | --- | --- |")
    for a in sorted(g.hilang):
        if a not in R:
            continue
        for j, n, lewat, kls in sorted(g.hilang[a]):
            if re.match(r"^(py|pz|px)", n) or n in ("Always", "Never", "NEVER", "null", "undefined"):
                continue
            w(f"| `{a[1]}` | {j} | `{sel(n)}` | {sel(lewat)} | {sel(kls)} |")
    w("", "**Selisih kelas** (nama ada di korpus, kelas berkasnya bukan leluhur kelas konteks rujukan):", "",
      "| Dari | Ke | Lewat | Kelas konteks | Kelas berkas |", "| --- | --- | --- | --- | --- |")
    for a in sorted(g.beda_kelas):
        if a not in R:
            continue
        for b, lewat, kk, kb in sorted(g.beda_kelas[a]):
            w(f"| `{a[1]}` | `{b[0]}/{b[1]}` | {sel(lewat)} | {sel(kk)} | {sel(kb)} |")
    w("")
    return w.teks()


def _tulis_langkah(w, langkah, tingkat):
    pad = "  " * tingkat
    for s in langkah:
        kepala = f"{pad}- **{s['no']}** `{sel(s['metode']) or '(kosong)'}`"
        if s["halaman"]:
            kepala += f" · halaman `{sel(s['halaman'])}`"
        if s["label"]:
            kepala += f" · label `{sel(s['label'])}`"
        if s["ket"]:
            kepala += f" — {sel(s['ket'], 200)}"
        w(kepala)
        if s["ulang"]:
            w(f"{pad}  - ulang: " + ", ".join(f"{k}=`{sel(v)}`" for k, v in s["ulang"].items()))
        for p in s["pre"]:
            if p["when"] or p["benar"] not in ("", "2") or p["salah"] not in ("", "2"):
                aktif = "" if s["whenAktif"] == "true" else " *(When langkah tidak dicentang)*"
                w(f"{pad}  - syarat {kode(p['when'])}: benar → {aksi(p['benar'], p['benarLabel'])}; "
                  f"salah → {aksi(p['salah'], p['salahLabel'])} [kode {p['benar'] or '-'}/{p['salah'] or '-'}]{aktif}")
        for p in s["trans"]:
            if p["when"] or p["benar"] not in ("", "2") or p["salah"] not in ("", "2"):
                w(f"{pad}  - transisi {kode(p['when'])}: benar → {aksi(p['benar'], p['benarLabel'])}; "
                  f"salah → {aksi(p['salah'], p['salahLabel'])} [kode {p['benar'] or '-'}/{p['salah'] or '-'}]")
        for n, v in s["param"]:
            w(f"{pad}  - {kode(n)} = {kode(v)}" if v else f"{pad}  - {kode(n)}")
        if s["panggil"]:
            w(f"{pad}  - parameter: " + "; ".join(f"{k}=`{sel(v)}`" for k, v in s["panggil"].items()))
        if s["java"]:
            w(f"{pad}  - Java:", "", "```java", SAMAR(s["java"]), "```", "")
        if s["anak"]:
            _tulis_langkah(w, s["anak"], tingkat + 1)


if __name__ == "__main__":
    sumber = sys.argv[1] if len(sys.argv) > 1 else "korpus.json"
    with open(sumber, encoding="utf-8") as fh:
        data = json.load(fh)
    teks = tulis(data)
    with open(TUJUAN, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(teks)
    print(f"{len(teks):,} karakter -> {os.path.normpath(TUJUAN)}")
    print("disamarkan:", len(SAMAR.op), "ID operator,", len(SAMAR.nama), "nama")
