"""Pembaca korpus XML Pega modul NB Treaty In - HANYA MEMBACA.

Dipakai `inventaris.py` untuk menyusun docs/INVENTARIS-XML.md. Korpus READ-ONLY:
tidak ada satu berkas pun yang ditulis, dipindah, atau dihapus.

Jalur korpus dibaca dari env var NBTREATYIN_KORPUS; bawaannya jalur di mesin
pengembang (PROMPT-IMPLEMENT-NB-TREATY-IN.md bab 1).

Lima jebakan (KEADAAN-NB-TREATY-IN.md bab 7) dan cara berkas ini menghindarinya:
  0. Activity memakai PropertiesName/PropertiesValue TANPA awalan py; DataTransform
     memakai pyPropertiesName/pyPropertiesValue. Keduanya dibaca terpisah.
  1. <rowdata REPEATINGINDEX="n"/> yang menutup sendiri: dibaca lewat ElementTree,
     bukan regex, sehingga sel kosong tetap punya tempat.
  2. pyRowNum pada pyOrConditions berbasis NOL: dipetakan ke baris (pyRowNum + 1).
  3. pyStepsPage tidak ada; halaman langkah dibaca dari pyStepsObjectName.
  4. Escape ganda: ElementTree membuka satu lapis, html.unescape membuka lapis kedua.
"""
from __future__ import annotations

import glob
import html
import json
import os
import re
import sys
import xml.etree.ElementTree as ET

KORPUS = os.environ.get(
    "NBTREATYIN_KORPUS", r"D:/NUSARE DEV/NusantaraRe/NB Treaty In (Done)")

JENIS = ["Activity", "DataTransform", "DecisionTable", "Flow", "FlowAction",
         "Harness", "RDBList", "ReportDefinition", "Section", "When"]


def t(e, tag):
    """Teks satu anak, escape lapis kedua dibuka, spasi tepi dipangkas."""
    if e is None:
        return ""
    x = e.findtext(tag)
    return html.unescape(x).strip() if x else ""


def raw(e, tag):
    """Teks satu anak TANPA dipangkas - untuk sel tabel keputusan seperti '"26   "'."""
    if e is None:
        return None
    x = e.find(tag)
    if x is None or x.text is None:
        return None
    return html.unescape(x.text)


def daun(e, buang=()):
    """Seluruh anak langsung yang berupa daun bernilai, sebagai dict."""
    out = {}
    if e is None:
        return out
    for ch in e:
        if len(ch) == 0 and ch.text and ch.text.strip() and not ch.tag.startswith(buang):
            out[ch.tag] = html.unescape(ch.text).strip()
    return out


# --------------------------------------------------------------------------- Activity

_BUANG_PARAM = ("pxCreate", "pxObjClass", "pyTempPlaceHolder", "pyExpressionGadget",
                "pyStepsParamUI", "pyParametersParam")


def _langkah(row, jalur):
    metode = t(row, "pyStepsActivityName")
    kata = metode.split(" ", 1)
    jenis = kata[0]
    target = kata[1].strip() if len(kata) > 1 else ""
    pre = []
    for p in row.findall("pyStepsPreCondParams/rowdata"):
        d = daun(p, ("pxCreate", "pxObjClass"))
        if d:
            pre.append({
                "when": d.get("pyStepsPreCondParamsWhen", ""),
                "benar": d.get("pyStepsPreCondParamsWhenTrue", ""),
                "benarLabel": d.get("pyStepsPreCondParamsWhenTruePrms", ""),
                "salah": d.get("pyStepsPreCondParamsWhenFalse", ""),
                "salahLabel": d.get("pyStepsPreCondParamsWhenFalsePrms", ""),
            })
    trans = []
    for p in row.findall("pyStepsTransParams/rowdata"):
        d = daun(p, ("pxCreate", "pxObjClass"))
        trans.append({
            "when": d.get("pyStepsTransParamsWhen", ""),
            "benar": d.get("pyStepsTransParamsWhenTrue", ""),
            "benarLabel": d.get("pyStepsTransParamsWhenTruePrms", ""),
            "salah": d.get("pyStepsTransParamsWhenFalse", ""),
            "salahLabel": d.get("pyStepsTransParamsWhenFalsePrms", ""),
        })
    param = []
    for p in row.findall("pyParamArray/rowdata"):
        d = daun(p, _BUANG_PARAM)
        if "PropertiesName" in d or "PropertiesValue" in d:
            param.append([d.get("PropertiesName", ""), d.get("PropertiesValue", "")])
        elif d:
            param.append(["; ".join(f"{k}={v}" for k, v in d.items()), ""])
    panggil = daun(row.find("pyStepsCallParams"), _BUANG_PARAM)
    ulang = daun(row.find("pyStepsRepeatDef"), ("pxCreate", "pxObjClass"))
    anak = [_langkah(r, f"{jalur}.{i}")
            for i, r in enumerate(row.findall("pySteps/rowdata"), 1)]
    return {
        "no": jalur,
        "metode": metode,
        "jenis": jenis,
        "target": target,
        "halaman": t(row, "pyStepsObjectName"),
        "kelas": t(row, "pyStepsClassName"),
        "label": t(row, "pyStepsBlockName"),
        "ket": t(row, "pyStepsDescription"),
        "whenAktif": t(row, "pyStepsPreCondition"),
        "pre": pre,
        "trans": trans,
        "param": param,
        "panggil": panggil,
        "ulang": ulang,
        "java": t(row, "pyStepsJavaSource"),
        "anak": anak,
    }


def activity(r):
    return {
        "parameter": [t(p, "pyParametersParamName") for p in r.findall("pyParameters/rowdata")
                      if t(p, "pyParametersParamName")],
        "lokal": [t(p, "pyLocalVariablesName") or t(p, "pyParametersParamName")
                  for p in r.findall("pyLocalVariables/rowdata")
                  if t(p, "pyLocalVariablesName") or t(p, "pyParametersParamName")],
        "langkah": [_langkah(row, str(i)) for i, row in enumerate(r.findall("pySteps/rowdata"), 1)],
    }


# --------------------------------------------------------------------------- DataTransform

def _dt_langkah(row, jalur):
    d = daun(row, ("pxCreate", "pxObjClass", "pyExpanded"))
    anak = [_dt_langkah(x, f"{jalur}.{i}")
            for i, x in enumerate(row.findall("pyProperties/rowdata"), 1)]
    # Parameter APPLY_MODEL (pyParameters/rowdata: nama = nilai) - tanpa ini
    # pemanggilan DT bersarang tampak tanpa masukan.
    param = [f"{t(p, 'pyParametersParamName')}={t(p, 'pyParametersParamValue')}"
             for p in row.findall("pyParameters/rowdata") if t(p, "pyParametersParamName")]
    if param:
        d["param"] = ", ".join(param)
    return {
        "no": jalur,
        "aksi": d.get("pyActionName", ""),
        "sasaran": d.get("pyPropertiesName", ""),
        "sumber": d.get("pyPropertiesValue", ""),
        "lain": {k: v for k, v in d.items()
                 if k not in ("pyActionName", "pyPropertiesName", "pyPropertiesValue", "pyPropertyStepId")},
        "anak": anak,
    }


def datatransform(r):
    return {"langkah": [_dt_langkah(row, str(i)) for i, row in enumerate(r.findall("pyProperties/rowdata"), 1)],
            "panggilInduk": t(r, "pyCallSuperClassModel")}


# --------------------------------------------------------------------------- Flow

def flow(r):
    mp = r.find("pyModelProcess")
    bentuk = []
    for s in mp.findall("pyShapes/rowdata"):
        router = s.find("pyRouterProp")
        wb = ""
        if router is not None:
            wb = t(router.find("pyCallParams"), "Workbasket")
        bentuk.append({
            "id": t(s, "pyMOId"), "jenis": t(s, "pyShapeType"), "nama": t(s, "pyMOName"),
            "implementasi": t(s, "pyImplementation"), "kelasKeputusan": t(s, "pyDecisionClass"),
            "workbasket": wb,
        })
    sambung = []
    for c in mp.findall("pyConnectors/rowdata"):
        tugas = [[t(p, "pyPropertiesName") or t(p, "PropertiesName"),
                  t(p, "pyPropertiesValue") or t(p, "PropertiesValue")]
                 for p in c.findall("pyPropertyAssigns/rowdata")]
        sambung.append({
            "dari": t(c, "pyFrom"), "ke": t(c, "pyTo"), "nama": t(c, "pyMOName"),
            "syarat": t(c, "pyConditionType"), "ekspresi": t(c, "pyExpression"),
            "kemungkinan": t(c, "pyLikelihood"),
            "tugas": [x for x in tugas if x[0] or x[1]],
        })
    return {"bentuk": bentuk, "sambung": sambung}


# --------------------------------------------------------------------------- Section / Harness / FlowAction

def _sel_milik(e):
    """Sel milik rule ini sendiri - salinan section tertanam (pyIncludedRuleXML) dilompati."""
    for ch in e:
        if ch.tag == "pyIncludedRuleXML":
            continue
        if ch.tag == "pyCells":
            for row in ch.findall("rowdata"):
                yield row
                yield from _sel_milik(row)
        else:
            yield from _sel_milik(ch)


def _aksi(e):
    out = []
    if e is None:
        return out
    for b in e.iter("rowdata"):
        if b.find("pyAction") is None:
            continue
        api = b.find("pyActionAPI")
        d = daun(api, ("pxObjClass", "pySelectedMobile", "pyAlwaysRender", "pyActivityClassOrig"))
        prm = [f"{t(p, 'pyName')}={t(p, 'pyValue')}" for p in (api.findall("pyActivityParams/rowdata") if api is not None else [])
               if t(p, "pyName")]
        out.append({"event": t(b, "pyEvent"), "aksi": t(b, "pyAction"),
                    "api": d, "param": prm})
    # buang duplikat (pyBehaviors dan pyActionSets sering memuat aksi yang sama)
    uniq, seen = [], set()
    for a in out:
        k = json.dumps(a, sort_keys=True)
        if k not in seen:
            seen.add(k)
            uniq.append(a)
    return uniq


def _sel(c):
    ud = c.find("pyUserData")
    modes = c.findall("pyModes/rowdata")
    md = {}
    for m in modes:
        md.update(daun(m, ("pxObjClass",)))
    sumber = {}
    for m in modes:
        lds = m.find("pyListDataSource")
        if lds is not None:
            for x in lds.iter():
                if len(x) == 0 and x.text and x.text.strip() and x.tag not in (
                        "pxObjClass", "pyValidPage", "pyIsValidDataPage", "pyNoSelectionText"):
                    sumber[x.tag] = html.unescape(x.text).strip()
    kontrol = t(c, "pyFormat")
    if kontrol in ("pxButton", "pxLink", "pxIcon"):
        # teks tombol ada di mode kontrol (pyLabel); pyLabelFieldValue berisi "Button"
        label = md.get("pyLabel", "") or t(c, "pyLabelPreview") or t(c, "pyLabelFieldValue")
    else:
        label = t(c, "pyLabelFieldValue") or md.get("pyLabel", "") or t(c, "pyLabelPreview")
    return {
        "tipe": t(c, "pyType"), "kontrol": kontrol, "nilai": t(c, "pyValue"),
        "label": label,
        "labelFor": t(c, "pyLabelFor"),
        "wajib": t(c, "pyRequired"), "wajibMode": md.get("pyRequired", ""),
        "wajibWhen": t(c, "pyRequiredCondition") or md.get("pyRequiredWhen", "") or md.get("pyRequiredCondition", ""),
        "kunci": t(c, "pyReadOnly"),
        "kunciWhen": t(ud, "pyReadOnlyCondition") if ud is not None else "",
        "nonaktif": md.get("pyDisabled", ""), "nonaktifWhen": md.get("pyDisabledWhen", ""),
        "tampil": t(ud, "pyVisible") if ud is not None else "",
        "tampilWhen": t(ud, "pyCondition") if ud is not None else "",
        "sumberDaftar": sumber,
        "aksi": _aksi(c.find("pyModes")) + _aksi(c.find("pyActionSets")),
        "sertakan": t(c, "pyInclude") if t(c, "pyType") in ("SUB_SECTION", "SECTION") else "",
    }


def _sertakan(r):
    """Section yang disertakan oleh rule ini (badan pyInclude dan sel SUB_SECTION)."""
    out = []
    for body in r.iter("rowdata"):
        if body.find("pyBodyType") is not None and t(body, "pyInclude"):
            out.append(t(body, "pyInclude"))
    for c in _sel_milik(r):
        if t(c, "pyType") in ("SUB_SECTION", "SECTION") and t(c, "pyInclude"):
            out.append(t(c, "pyInclude"))
    for x in r.iter("pyIncludedRuleXML"):
        n = t(x, "pyStreamName") or t(x, "pyRuleName")
        if n:
            out.append(n)
    return sorted(set(out))


def _badan_tampil(r):
    """Syarat tampil/kunci tingkat layout (badan section) - bukan sel."""
    out = []
    for body in r.iter("rowdata"):
        if body.find("pyBodyType") is None:
            continue
        ud = body.find("pyUserData")
        vis = t(ud, "pyVisible") if ud is not None else ""
        cond = t(ud, "pyCondition") if ud is not None else ""
        ro = t(ud, "pyReadOnlyCondition") if ud is not None else ""
        if cond or ro or (vis and vis != "ALWAYS"):
            out.append({"sertakan": t(body, "pyInclude"), "halaman": t(body, "pyUsingPage"),
                        "tampil": vis, "tampilWhen": cond, "kunciWhen": ro})
    return out


def section(r):
    sel = [_sel(c) for c in _sel_milik(r)]
    return {"sel": [s for s in sel if s["tipe"] in ("FIELD",) or s["aksi"] or s["sertakan"]
                    or s["tampilWhen"] or s["kunciWhen"]],
            "sertakan": _sertakan(r),
            "badan": _badan_tampil(r),
            "halaman": sorted({t(b, "pyUsingPage") for b in r.iter("rowdata")
                               if b.find("pyBodyType") is not None and t(b, "pyUsingPage")})}


def flowaction(r):
    kunci = ["pySectionReference", "pyPreProcessingActivity", "pyPreProcessingTransformRule",
             "pyActionActivity", "pyPostProcessingActivity", "pyActionTransformRule",
             "pyLocalActionActivity", "pyValidateRule", "pyValidation", "pyValidateActivity",
             "pyConfirmHarness", "pyconfirmchoice", "pySubmitLabel", "pyCancelLabel",
             "pyAuditActivity", "pyUsedAs", "pyActionName"]
    d = {k: t(r, k) for k in kunci if t(r, k)}
    for nama in ("pyPreProcessingActivityParams", "pyActionActivityParams", "pyPostProcessingActivityParams"):
        e = r.find(nama)
        if e is not None:
            dd = daun(e, ("pxObjClass", "pyTempPlaceHolder"))
            if dd:
                d[nama] = dd
    d.update(section(r))
    return d


def harness(r):
    d = section(r)
    d["kontainer"] = t(r, "pyContainerType")
    return d


# --------------------------------------------------------------------------- RDBList / ReportDefinition

_TABEL = re.compile(r"\b(?:from|join|into|update|table)\s+([A-Za-z_][\w$#]*(?:\.[A-Za-z_][\w$#]*)?)", re.I)
_PROC = re.compile(r"\b(?:call|exec(?:ute)?|begin)\s+([A-Za-z_][\w$#]*(?:\.[A-Za-z_][\w$#]*){0,2})", re.I)


def rdblist(r):
    sql = {k: t(r, k) for k in ("pyBrowseSQL", "pySaveSQL", "pyDeleteSQL", "pyOpenSQL") if t(r, k)}
    semua = " ".join(sql.values())
    return {"permintaan": t(r, "pyRequestType"), "akses": t(r, "pyRWAccess"), "sql": sql,
            "tabel": sorted({x.upper() for x in _TABEL.findall(semua)}),
            "prosedur": sorted({x.upper() for x in _PROC.findall(semua)
                                if x.upper() not in ("DUAL",)})}


def reportdef(r):
    kol = []
    for c in r.iter("rowdata"):
        if c.find("pyFieldName") is not None and t(c, "pyFieldName"):
            kol.append(t(c, "pyFieldName"))
    filt = []
    for f in r.iter("rowdata"):
        # Filter RD Pega: baris Embed-ReportFilter (pyFilterName/Operation/Value
        # + pyLogicLabel). Bentuk lama pyLeft/pyCondition/pyRight tetap dibaca.
        if f.find("pyFilterName") is not None and t(f, "pyFilterName"):
            lab = t(f, "pyLogicLabel")
            teks = f"{t(f, 'pyFilterName')} {t(f, 'pyFilterOperation')} {t(f, 'pyFilterValue')}".strip()
            filt.append(f"{lab}: {teks}" if lab else teks)
        elif f.find("pyLeft") is not None or f.find("pyCondition") is not None and f.find("pyRight") is not None:
            l, op, rr = t(f, "pyLeft"), t(f, "pyCondition") or t(f, "pyOperator"), t(f, "pyRight")
            if l or rr:
                filt.append(f"{l} {op} {rr}".strip())
    logika = ""
    for el in r.iter("pyFilterLogic"):
        if (el.text or "").strip():
            logika = el.text.strip()
            break
    return {"kolom": list(dict.fromkeys(kol)), "filter": list(dict.fromkeys(filt)),
            "logika": logika,
            "kelasLaporan": t(r, "pyClassName")}


# --------------------------------------------------------------------------- When / DecisionTable

def when(r):
    kond = []
    for c in r.findall("pyCondition/rowdata"):
        kond.append({"label": t(c, "pyConditionLabel"),
                     "teks": t(c, "pyConditionValue1String") or t(c, "pyConditionValue1StringLabel"),
                     "ekspresi": t(c, "pyConditionValue1")})
    viewer = r.find("pyConditionViewer")
    return {"logika": t(r, "pyLogic"), "kondisi": [k for k in kond if k["teks"] or k["ekspresi"]],
            "logikaPenampil": t(viewer, "pyLogic"),
            "kondisiPenampil": [t(x, "pyConditionString") for x in (viewer.findall("pyNestedConditions/rowdata") if viewer is not None else [])]}


def decisiontable(r):
    kolom = []
    for col in r.findall("pyColumns/rowdata"):
        sel = [raw(c, ".") if False else (html.unescape(c.text) if c.text is not None else None)
               for c in col.findall("pyCondition/rowdata")]
        ors = {}
        for o in col.findall("pyOrConditions/rowdata"):
            n = t(o, "pyRowNum")
            if n and n != "-1":
                ors[int(n) + 1] = [html.unescape(x.text) if x.text is not None else None
                                   for x in o.findall("pyOrCondition/rowdata")]
        kolom.append({"properti": t(col, "pyProperty"), "tipe": t(col, "pyColumnDataType"),
                      "operator": t(col, "pyDefaultOperator"), "sel": sel, "atau": ors})
    hasil = [html.unescape(c.text) if c.text is not None else None for c in r.findall("pyResults/rowdata")]
    setprop = []
    for pc in r.findall("pyPropertyColumns/rowdata"):
        setprop.append({"properti": t(pc, "pyProperty"),
                        "nilai": [html.unescape(c.text) if c.text is not None else None
                                  for c in pc.findall("pyPropertyValues/rowdata")]})
    return {"kolom": kolom, "hasil": hasil, "setProperti": setprop,
            "bawaan": t(r, "pyDefaultResult"), "aksiBawaan": t(r, "pyDefaultAction"),
            "evaluasiSemua": t(r.find("pyDelegatedRestrictions/rowdata"), "pyEvaluateAllRows")}


# --------------------------------------------------------------------------- rujukan mentah

# pxRuleObjClass Pega -> jenis folder korpus
_JENIS_PEGA = {
    "Rule-Obj-Activity": "Activity", "Rule-Obj-When": "When", "Rule-HTML-Section": "Section",
    "Rule-HTML-Harness": "Harness", "Rule-Obj-FlowAction": "FlowAction", "Rule-Obj-Model": "DataTransform",
    "Rule-Declare-DecisionTable": "DecisionTable", "Rule-Connect-SQL": "RDBList",
    "Rule-Obj-Report-Definition": "ReportDefinition", "Rule-Obj-Flow": "Flow",
}

# tag yang nilainya NAMA RULE, di luar yang dibaca per jenis
_TAG_RUJUKAN = {
    "pyDeferLoadActivity": "Activity", "pyEditAction": "FlowAction", "pyGridRDName": "ReportDefinition",
    "pyContainerVisibleWhen": "When", "pyModelName": "DataTransform", "pyLocalAction": "FlowAction",
    "pyHarnessName": "Harness", "pyRefreshWhen": "When", "pyVisibleWhen": "When",
}


def _tanpa_tanaman(r):
    """Salinan root tanpa subtree pyIncludedRuleXML (salinan section lain)."""
    for p in r.iter():
        for ch in list(p):
            if ch.tag == "pyIncludedRuleXML":
                p.remove(ch)
    return r


def rujukan(r):
    out = []
    for rr in r.findall("pxRuleReferences/rowdata"):
        j = _JENIS_PEGA.get(t(rr, "pxRuleObjClass"))
        if not j:
            continue
        nama = t(rr, "pyRuleName")
        if j == "RDBList":            # "Kelas ASM NamaPermintaan"
            nama = nama.split(" ")[-1]
        out.append({"jenis": j, "kelas": t(rr, "pxRuleClassName"), "nama": nama, "lewat": "pxRuleReferences"})
    for x in r.iter():
        if x.tag in _TAG_RUJUKAN and x.text and x.text.strip():
            v = html.unescape(x.text).strip().lstrip("!").strip()
            if re.fullmatch(r"[A-Za-z_][\w\-]*", v):
                out.append({"jenis": _TAG_RUJUKAN[x.tag], "kelas": "", "nama": v, "lewat": x.tag})
        if x.tag == "pyPreDataTransform":
            n = t(x, "pyName")
            if n:
                out.append({"jenis": "DataTransform", "kelas": "", "nama": n, "lewat": "pyPreDataTransform"})
        if x.tag == "pyPreActivity" and x.text and x.text.strip():
            out.append({"jenis": "Activity", "kelas": "", "nama": x.text.strip(), "lewat": "pyPreActivity"})
    uniq, seen = [], set()
    for o in out:
        k = (o["jenis"], o["kelas"], o["nama"], o["lewat"])
        if k not in seen:
            seen.add(k)
            uniq.append(o)
    return uniq


# --------------------------------------------------------------------------- pemuat

_PEMBACA = {"Activity": activity, "DataTransform": datatransform, "DecisionTable": decisiontable,
            "Flow": flow, "FlowAction": flowaction, "Harness": harness, "RDBList": rdblist,
            "ReportDefinition": reportdef, "Section": section, "When": when}


def nama_rule(jenis, r, berkas):
    # Nama berkas = nama rule. pyLabel tidak dipakai: pada Harness/ReportDefinition ia
    # berisi keterangan ("Displays a list of Opportunities"), bukan nama.
    return os.path.splitext(os.path.basename(berkas))[0]


def muat(korpus=KORPUS):
    out = []
    for jenis in JENIS:
        for f in sorted(glob.glob(os.path.join(korpus, jenis, "*.xml"))):
            r = ET.parse(f).getroot()
            d = {"jenis": jenis, "berkas": f"{jenis}/{os.path.basename(f)}",
                 "byte": os.path.getsize(f),
                 "nama": nama_rule(jenis, r, f), "kelas": t(r, "pyClassName"),
                 "ket": t(r, "pyDescription"), "versi": t(r, "pyRuleSetVersion")}
            d.update(_PEMBACA[jenis](r))
            d["rujukan"] = rujukan(_tanpa_tanaman(r))
            out.append(d)
    return out


if __name__ == "__main__":
    data = muat()
    tujuan = sys.argv[1] if len(sys.argv) > 1 else "korpus.json"
    with open(tujuan, "w", encoding="utf-8") as fh:
        json.dump(data, fh, ensure_ascii=False, indent=1)
    print(len(data), "rule ->", tujuan)
