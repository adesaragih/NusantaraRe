"""Urai pohon langkah sebuah Rule-Obj-Activity Pega (ekspor XML).

Alat audit before-image EDM (docs/LAPORAN-IMPLEMENTASI-BEFORE-IMAGE.md §5).
Membaca korpus READ-ONLY; tidak menulis apa pun.

Yang dicetak per langkah: nomor bertingkat, metode, objek loop
(pyStepsObjectName), kelas, deskripsi, pyStepsPreCondition, ekspresi
prakondisi (pyStepsPreCondParamsWhen) beserta transisi WhenTrue/WhenFalse,
lalu pasangan pyParamArray PropertiesName = PropertiesValue.

⛔ PRIVASI (E-Q15): pyStepsRepeatDef dan isi pyStepsPreCondParams selain
ekspresi `...When` memuat nama orang dan alamat email. Keduanya TIDAK dibaca.

Pemakaian (pemanggil `py`, awali PYTHONIOENCODING=utf-8 - CLAUDE.md §4a):
  py langkah.py <berkas.xml>          cetak pohon langkah
  py langkah.py --dt <berkas.xml>     aksi DataTransform
  py langkah.py --lokasi-pola <xml>   cacah pola nomor polis per tag (nilai tak dicetak)
  py langkah.py --uji <folder-korpus>  uji instrumen atas angka yang sudah
                                       diketahui (folder `Endorsment Fac In`)
"""
import os
import re
import sys
import xml.etree.ElementTree as ET


def teks(el, jalur):
    c = el.find(jalur)
    return (c.text or "").strip() if c is not None else ""


def langkah(akar):
    """Hasilkan (nomor, rowdata) seluruh langkah, berurutan, bertingkat."""
    def jalan(rows, awalan):
        for i, row in enumerate(rows, 1):
            nomor = f"{awalan}{i}"
            yield nomor, row
            anak = row.find("pySteps")
            if anak is not None:
                yield from jalan(anak.findall("rowdata"), nomor + ".")
    yield from jalan(akar.findall("pySteps/rowdata"), "")


def penugasan(row):
    """Pasangan (nama, nilai) Property-Set satu langkah."""
    for p in row.findall("pyParamArray/rowdata"):
        n = p.find("PropertiesName")
        if n is not None:
            yield (n.text or "").strip(), teks(p, "PropertiesValue")


# Pola nomor polis produksi, TIDAK peka huruf (jebakan sensus §4a no. 4).
POLA_NOPOLIS = re.compile(r"RNM-[A-Z0-9.\-]+", re.IGNORECASE)


def samar(s):
    """Nomor polis produksi (pola `RNM-…`) tidak pernah dicetak (CLAUDE.md §4.10)."""
    return POLA_NOPOLIS.sub("<NOPOLIS>", s)


def cetak(berkas):
    for nomor, row in langkah(ET.parse(berkas).getroot()):
        ind = "  " * nomor.count(".")
        pre = row.find("pyStepsPreCondParams/rowdata")
        # Setiap baris prakondisi: ekspresi, WhenTrue(->target), WhenFalse(->target).
        baris = []
        for r in row.findall("pyStepsPreCondParams/rowdata"):
            w = teks(r, "pyStepsPreCondParamsWhen")
            if w:
                tt, ft = teks(r, "pyStepsPreCondParamsWhenTruePrms"), teks(r, "pyStepsPreCondParamsWhenFalsePrms")
                baris.append(f"{samar(w)} T={teks(r, 'pyStepsPreCondParamsWhenTrue')}{'->' + tt if tt else ''}"
                             f" F={teks(r, 'pyStepsPreCondParamsWhenFalse')}{'->' + ft if ft else ''}")
        when = baris
        wt = teks(pre, "pyStepsPreCondParamsWhenTrue") if pre is not None else ""
        wf = teks(pre, "pyStepsPreCondParamsWhenFalse") if pre is not None else ""
        label = teks(row, 'pyStepsBlockName')
        label = f" label={label!r}" if label else ""
        print(f"{ind}[{nomor}]{label} {teks(row, 'pyStepsActivityName')} obj={teks(row, 'pyStepsObjectName')!r} "
              f"class={teks(row, 'pyStepsClassName')!r} desc={samar(teks(row, 'pyStepsDescription'))!r} "
              f"pre={teks(row, 'pyStepsPreCondition')!r} when={when}")
        # Langkah tanpa metode (`pyStepsActivityName` kosong, mis. loop EMBEDDED):
        # pasangan pyParamArray-nya mungkin sisa yang TIDAK dieksekusi (temuan sesi
        # nbfacin 01-10-2026) - ditandai, bukan dicetak sebagai SET biasa.
        tanda = "SET" if teks(row, 'pyStepsActivityName') else "SET? (langkah tanpa metode)"
        for n, v in penugasan(row):
            if n:
                print(f"{ind}    {tanda} {samar(n)} = {samar(v)}")


def uji(folder):
    """Instrumen diuji atas butir yang jawabannya SUDAH diketahui dari bahan
    spec (docs/07-edm/08-bahan-spec-before-image.md §2-§3), yang dihitung
    dengan cara lain (cacah baris). Gagal = instrumennya yang rusak."""
    akt = os.path.join(folder, "Activity")
    old = [(nm, v) for _, r in langkah(ET.parse(os.path.join(akt, "SetOldData.xml")).getroot())
           for nm, v in penugasan(r) if nm.endswith("Old")]
    nilai = dict()
    for nomor, r in langkah(ET.parse(os.path.join(akt, "SetValueToEDMWork.xml")).getroot()):
        nilai[nomor] = list(penugasan(r))
    harapan = [
        ("SetOldData: penugasan *Old (bahan §3)", len(old), 52),
        ("SetOldData: RateOld self-referential (bahan §3.3)",
         sum(1 for nm, v in old if nm.endswith("RateOld") and v.startswith('@If(.RateOld!=""')), 6),
        ("SetOldData: PremiumOld self-referential (bahan §3.3)",
         sum(1 for nm, v in old if nm.endswith(".PremiumOld") and v.startswith('@If(.PremiumOld!=""')), 1),
        ("SetOldData: cadangan string \"0\" (bahan §3.2)", sum(1 for _, v in old if v.endswith(',"0")')), 2),
        ("SetValueToEDMWork 14.3: penugasan (bahan §2)", len(nilai.get("14.3", [])), 54),
        ("SetValueToEDMWork 14.3: bersumber OldData (bahan §2)",
         sum(1 for _, v in nilai.get("14.3", []) if ".OldData." in v), 51),
    ]
    gagal = 0
    for label, dapat, mau in harapan:
        tanda = "OK " if dapat == mau else "GAGAL"
        gagal += dapat != mau
        print(f"{tanda} {label}: {dapat} (mau {mau})")
    sys.exit(1 if gagal else 0)


def cetak_dt(berkas):
    """Rule-Obj-Model (DataTransform): baris aksi berurutan, bertingkat.
    `pyDisabled=true` = dilewati (K-048, keputusan work owner)."""
    def jalan(el, ind):
        for row in el.findall("rowdata"):
            aksi = teks(row, "pyActionName")
            if aksi:
                mati = " [pyDisabled]" if teks(row, "pyDisabled") == "true" else ""
                print(f"{'  ' * ind}{aksi}{mati}: {samar(teks(row, 'pyPropertiesName'))} = {samar(teks(row, 'pyPropertiesValue'))}")
            for anak in row:
                if anak.tag != "pyExpressionGadget" and anak.findall("rowdata"):
                    jalan(anak, ind + 1)
    for el in ET.parse(berkas).getroot():
        if el.findall("rowdata") and el.tag not in ("pyPagesAndClasses", "pyRuleVersionsList"):
            jalan(el, 0)


def lokasi_pola(berkas):
    """Cacah kemunculan pola nomor polis per NAMA TAG tempatnya tinggal - nilai
    tidak pernah dicetak. Membedakan logika yang dieksekusi
    (`pyStepsPreCondParamsWhen`, `PropertiesValue`) dari cache editor
    (`pyExpressionGadget/...`)."""
    import collections
    c = collections.Counter()

    def jalan(el, jalur):
        for ch in el:
            p = jalur + "/" + ch.tag
            if ch.text:
                n = len(POLA_NOPOLIS.findall(ch.text))
                if n:
                    c[p.replace("/rowdata", "")] += n
            jalan(ch, p)
    jalan(ET.parse(berkas).getroot(), "")
    for jalur, n in c.most_common():
        print(n, jalur)
    print("total", sum(c.values()))


if __name__ == "__main__":
    if len(sys.argv) == 3 and sys.argv[1] == "--lokasi-pola":
        lokasi_pola(sys.argv[2])
    elif len(sys.argv) == 3 and sys.argv[1] == "--dt":
        cetak_dt(sys.argv[2])
    elif len(sys.argv) == 3 and sys.argv[1] == "--uji":
        uji(sys.argv[2])
    elif len(sys.argv) == 2:
        cetak(sys.argv[1])
    else:
        sys.exit(__doc__)
