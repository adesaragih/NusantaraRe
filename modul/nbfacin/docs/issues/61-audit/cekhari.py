# Uji empiris DaysValidity default NB (GetDateValidity_ACT langkah 2): potong vs bulat atas selisih hari (hanya hitungan).
import glob, re, math, datetime as dt, collections

def waktu(s):
    return dt.datetime.strptime(s, "%Y%m%dT%H%M%S.%f GMT").replace(tzinfo=dt.timezone.utc)

hasil = collections.Counter()
for f in glob.glob(r"D:/migrasi/RNM/DDL/CONTOH/*.xml"):
    t = open(f, encoding="utf-8", errors="replace").read()
    m = re.search(r"<DaysValidity>([^<]*)</DaysValidity>", t)
    b = re.search(r"<StartDateTime>([^<]*)</StartDateTime>", t)
    dv = re.search(r"<DateValidity>([^<]*)</DateValidity>", t)
    if not (m and b and dv and m.group(1).isdigit()):
        continue
    hari = int(m.group(1))
    if hari in (0, 30):
        hasil["hari 0 / 30 (tidak membedakan)"] += 1
        continue
    # 'sekarang' = saat DateValidity dihitung: tanggal DateValidity - hari (WIB); jam tidak diketahui -> cari DateSuggest di hari itu
    wib = dt.timezone(dt.timedelta(hours=7))
    tgl_hitung = waktu(dv.group(1)).astimezone(wib).date() - dt.timedelta(days=hari)
    sug = [waktu(s) for s in re.findall(r"<DateSuggest>([^<]*)</DateSuggest>", t)]
    kand = [s for s in sug if s.astimezone(wib).date() == tgl_hitung]
    if not kand:
        hasil["tanpa DateSuggest di hari hitung"] += 1
        continue
    ok_potong = ok_bulat = False
    for now in kand:
        d = (waktu(b.group(1)) - now).total_seconds() / 86400
        if d <= 0:
            continue
        p = math.floor(d) + 1 if math.floor(d) > 0 else 0
        r = round(d) + 1 if round(d) > 0 else 0
        ok_potong |= min(p, 30) == hari
        ok_bulat |= min(r, 30) == hari
    hasil["potong cocok=%s bulat cocok=%s" % (ok_potong, ok_bulat)] += 1
for k, v in sorted(hasil.items()):
    print(v, k)
