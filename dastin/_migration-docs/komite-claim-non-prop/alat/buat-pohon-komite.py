# -*- coding: utf-8 -*-
"""Pembangkit pohon properti modul Komite Claim Non Prop.

Menyapu SELURUH berkas XML di `Komite Claim Non Prop/` dan menyusun:
  - pohon properti per halaman akar (pyWorkPage, TempMainWork, dst.)
  - peta kelas Pega dari deklarasi pyPagesAndClassesPage/Class
  - kontrak tulis-balik ke halaman induk (TempMainWork)

Keluaran: tiga CSV di `4-erd-dan-tabel-datar/`.
Berkas ini memerikan SISTEM LAMA. Nama ditulis apa adanya, termasuk salah ejanya.
"""
import io, os, re, csv, collections

SUMBER = r"D:\XML_NURE\Komite Claim Non Prop"
KELUAR = r"D:\XML_NURE\_migration-docs\komite-claim-non-prop\4-erd-dan-tabel-datar"

if not os.path.isdir(KELUAR):
    os.makedirs(KELUAR)

# ---------------------------------------------------------------- sapuan
berkas = []
for akar, _, nama in os.walk(SUMBER):
    for n in nama:
        if n.lower().endswith(".xml"):
            berkas.append(os.path.join(akar, n))
berkas.sort()

isi = {}
for p in berkas:
    with io.open(p, encoding="utf-8", errors="replace") as f:
        isi[p] = f.read()

def rel(p):
    return os.path.relpath(p, SUMBER).replace("\\", "/")

# ---------------------------------------------------- peta halaman -> kelas
halaman_kelas = collections.defaultdict(collections.Counter)
for p, s in isi.items():
    hal = re.findall(r"<pyPagesAndClassesPage>([^<]*)</pyPagesAndClassesPage>", s)
    kls = re.findall(r"<pyPagesAndClassesClass>([^<]*)</pyPagesAndClassesClass>", s)
    for a, b in zip(hal, kls):
        a, b = a.strip(), b.strip()
        if a and b:
            halaman_kelas[a][b] += 1
    # kelas rule itu sendiri
    for m in re.finditer(r"<pyClassName>([^<]*)</pyClassName>", s):
        halaman_kelas["(kelas rule)"][m.group(1).strip()] += 1

AKAR = sorted(set(h.split(".")[0] for h in halaman_kelas if h != "(kelas rule)"))
for tambahan in ("pyWorkPage", "TempMainWork", "Primary", "Local", "Param", "pxRequestor"):
    if tambahan not in AKAR:
        AKAR.append(tambahan)
AKAR = sorted(set(AKAR), key=len, reverse=True)

# ------------------------------------------------------- sapuan jalur
POLA = re.compile(
    r"\b(" + "|".join(re.escape(a) for a in AKAR) + r")"
    r"((?:\([^()<>]{0,40}\))?(?:\.[A-Za-z][A-Za-z0-9_]*(?:\([^()<>]{0,40}\))?)+)"
)
SEG = re.compile(r"([A-Za-z][A-Za-z0-9_]*)(\([^()<>]{0,40}\))?")

simpul = {}          # jalur kanonik -> dict
def daftar(jalur, akar, nama, kedalaman, berlist, berkas_nama):
    d = simpul.setdefault(jalur, {
        "akar": akar, "nama": nama, "kedalaman": kedalaman,
        "list": False, "ref": 0, "berkas": set(), "anak": set()})
    if berlist:
        d["list"] = True
    d["ref"] += 1
    d["berkas"].add(berkas_nama)
    return d

for p, s in isi.items():
    bn = rel(p)
    for m in POLA.finditer(s):
        akar = m.group(1)
        ekor = m.group(2)
        segmen = []
        for sm in SEG.finditer(ekor):
            segmen.append((sm.group(1), bool(sm.group(2))))
        if not segmen:
            continue
        jalur = akar
        akar_berlist = ekor.startswith("(")
        if akar_berlist:
            daftar(akar, akar, akar, 0, True, bn)
        for i, (nm, berlist) in enumerate(segmen):
            induk = jalur
            jalur = jalur + "." + nm
            daftar(jalur, akar, nm, i + 1, berlist, bn)
            if induk in simpul:
                simpul[induk]["anak"].add(jalur)
            elif induk == akar:
                simpul.setdefault(akar, {
                    "akar": akar, "nama": akar, "kedalaman": 0, "list": False,
                    "ref": 0, "berkas": set(), "anak": set()})["anak"].add(jalur)

# buang simpul semu
BUANG = ("pySteps", "pyActionPrompt", "pxResults.px", "pyTemplate")
for j in list(simpul):
    if any(b in j for b in BUANG):
        del simpul[j]
for j in simpul:
    simpul[j]["anak"] = set(a for a in simpul[j]["anak"] if a in simpul)

# ------------------------------------------------------------ tulis CSV
def kelas_untuk(jalur):
    if jalur in halaman_kelas:
        return halaman_kelas[jalur].most_common(1)[0][0]
    return ""

with io.open(os.path.join(KELUAR, "datar-komite-lama.csv"), "w", encoding="utf-8", newline="") as f:
    w = csv.writer(f)
    w.writerow(["JALUR", "AKAR", "KEDALAMAN", "NAMA", "BENTUK", "KELAS_PEGA",
                "CACAH_ANAK", "REF", "CACAH_BERKAS", "BERKAS_CONTOH"])
    for j in sorted(simpul):
        d = simpul[j]
        w.writerow([j, d["akar"], d["kedalaman"], d["nama"],
                    "Page List" if d["list"] else ("Page" if d["anak"] else "skalar"),
                    kelas_untuk(j), len(d["anak"]), d["ref"], len(d["berkas"]),
                    " | ".join(sorted(d["berkas"])[:3])])

with io.open(os.path.join(KELUAR, "datar-komite-kelas.csv"), "w", encoding="utf-8", newline="") as f:
    w = csv.writer(f)
    w.writerow(["HALAMAN", "KELAS_PEGA", "CACAH_DEKLARASI"])
    for h in sorted(halaman_kelas):
        for k, n in halaman_kelas[h].most_common():
            w.writerow([h, k, n])

# ------------------------------------------- kontrak tulis-balik ke induk
tulis = collections.Counter()
for p, s in isi.items():
    for m in re.finditer(r"\bTempMainWork((?:\.[A-Za-z][A-Za-z0-9_]*(?:\([^()<>]{0,40}\))?)+)", s):
        tulis[("TempMainWork" + re.sub(r"\([^()<>]{0,40}\)", "", m.group(1)), rel(p))] += 1
with io.open(os.path.join(KELUAR, "datar-komite-tulis-balik.csv"), "w", encoding="utf-8", newline="") as f:
    w = csv.writer(f)
    w.writerow(["JALUR_INDUK", "BERKAS", "REF"])
    for (j, b), n in sorted(tulis.items()):
        w.writerow([j, b, n])

# ------------------------------------------------------------- ringkasan
print("berkas XML disapu : %d" % len(berkas))
print("simpul            : %d" % len(simpul))
per_akar = collections.Counter(d["akar"] for d in simpul.values())
for a, n in per_akar.most_common(12):
    print("   %-24s %d" % (a, n))
print("kedalaman maks    : %d" % max(d["kedalaman"] for d in simpul.values()))
print("halaman berkelas  : %d" % len([h for h in halaman_kelas if h != "(kelas rule)"]))
kls = set()
for h, c in halaman_kelas.items():
    if h != "(kelas rule)":
        kls |= set(c)
print("kelas berbeda     : %d" % len(kls))
for k in sorted(kls):
    print("     ", k)
