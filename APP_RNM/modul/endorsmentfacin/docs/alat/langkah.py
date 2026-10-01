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
  py langkah.py --uji <folder-korpus>  uji instrumen atas angka yang sudah
                                       diketahui (folder `Endorsment Fac In`)
"""
import os
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


def cetak(berkas):
    for nomor, row in langkah(ET.parse(berkas).getroot()):
        ind = "  " * nomor.count(".")
        pre = row.find("pyStepsPreCondParams/rowdata")
        when = [(w.text or "").strip() for w in row.findall("pyStepsPreCondParams/rowdata/pyStepsPreCondParamsWhen") if w.text]
        wt = teks(pre, "pyStepsPreCondParamsWhenTrue") if pre is not None else ""
        wf = teks(pre, "pyStepsPreCondParamsWhenFalse") if pre is not None else ""
        print(f"{ind}[{nomor}] {teks(row, 'pyStepsActivityName')} obj={teks(row, 'pyStepsObjectName')!r} "
              f"class={teks(row, 'pyStepsClassName')!r} desc={teks(row, 'pyStepsDescription')!r} "
              f"pre={teks(row, 'pyStepsPreCondition')!r} when={when} T={wt} F={wf}")
        for n, v in penugasan(row):
            print(f"{ind}    SET {n} = {v}")


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


if __name__ == "__main__":
    if len(sys.argv) == 3 and sys.argv[1] == "--uji":
        uji(sys.argv[2])
    elif len(sys.argv) == 2:
        cetak(sys.argv[1])
    else:
        sys.exit(__doc__)
