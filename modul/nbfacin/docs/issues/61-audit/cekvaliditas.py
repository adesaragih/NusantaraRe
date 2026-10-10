# Uji empiris tafsir @DateTime.addCalendar atas contoh ekspor DDL\CONTOH (hanya hitungan cocok, tanpa nilai).
import glob, re, datetime as dt, calendar, collections

WIB = dt.timezone(dt.timedelta(hours=7))


def waktu(s):
    return dt.datetime.strptime(s, "%Y%m%dT%H%M%S.%f GMT").replace(tzinfo=dt.timezone.utc)


def tgl(s):
    return waktu(s).astimezone(WIB).date()


def tambah_bulan(d, n):
    y, m = divmod(d.month - 1 + n, 12)
    y += d.year
    m += 1
    return dt.date(y, m, min(d.day, calendar.monthrange(y, m)[1]))


hasil = collections.Counter()
for f in glob.glob(r"D:/migrasi/RNM/DDL/CONTOH/*.xml"):
    t = open(f, encoding="utf-8", errors="replace").read()
    def satu(tag):
        m = re.search(r"<%s>([^<]*)</%s>" % (tag, tag), t)
        return m.group(1) if m else ""
    dv, days, wb = satu("DateValidity"), satu("DaysValidity"), satu("WaitingBindDate")
    begin, touw, tba = satu("StartDateTime"), satu("DateToUW"), satu("IsTBA")
    sug = re.findall(r"<DateSuggest>([^<]*)</DateSuggest>", t)
    if dv and days.isdigit() and sug:
        tgl_sug = sorted(tgl(s) for s in sug)
        cocok = any((x + dt.timedelta(days=int(days))) == tgl(dv) for x in tgl_sug)
        hasil["DateValidity = tgl DateSuggest + DaysValidity hari: " + ("cocok" if cocok else "tidak")] += 1
    if wb and begin and sug:
        last = sorted(waktu(s) for s in sug)[-1]
        kandidat = {"Begin+1bln": tambah_bulan(tgl(begin), 1), "LastAksep+1bln": tambah_bulan(last.astimezone(WIB).date(), 1),
                    "LastAksep+3bln": tambah_bulan(last.astimezone(WIB).date(), 3)}
        pas = [k for k, v in kandidat.items() if v == tgl(wb)]
        # rumus L701 dengan tafsir arg2 = bulan
        if touw and begin:
            if waktu(touw) <= waktu(begin):
                harap = "LastAksep+3bln" if tba == "true" else "Begin+1bln"
            else:
                harap = "LastAksep+1bln"
            hasil["WaitingBind sesuai rumus L701 (arg2=bulan, LastAksep=DateSuggest terakhir): " + ("cocok" if harap in pas else "tidak")] += 1
        hasil["WaitingBind cocok salah satu kandidat: " + (",".join(pas) or "tidak ada")] += 1
for k, v in sorted(hasil.items()):
    print(v, k)
