// Layar tahun treaty — tiket 03 Treaty Contract Out.
//
// Meniru `Harness/InboxTreatyContract.xml` → `Section/GridTreatyContract.xml`
// (judul b1057) → `Section/InputTreatyContract.xml`: grid `BrowseTreatyYear_RD`
// (b17355; sort `.ID DESC` b672) dengan tombol baris `Edit` b19939 /
// `ReinsType` b20778 / `List Description` b22196, tombol `Add` b16387, dan
// form lengkap "Input New Data" dari `Section/InputDtlTreatyContact.xml`
// (disertakan b1206; `Save` b10332 → `SaveTreatyYear_Act`, `Cancel` b10622).
//
// ⛔ Identitas tidak pernah diketik (AC 5): medan `ID` hanya dibaca; baris
// baru dikirim TANPA id (POST), pembaruan lewat PUT /{id}.
//
// ⛔ Gerbang ada di SERVER: periode terbalik (422), anti-dobel (409), wajib
// isi (422). Layar menampilkan kalimatnya apa adanya (`Gagal`).
//
// ⛔ Tombol `Copy` b20459 dan form From/To b2374/b3818 TIDAK dibawa — fitur
// salin tahun treaty dibuang (AC 72). `ReinsType` membuka `PanelKontrakTahun`
// (tiket 04), `List Description` membuka `PanelKlausulTahun` (tiket 08).
//
// [keputusan work owner 30-09-2026] Panel yang dibuka dari baris tahun hanya
// SATU sekaligus, dan selama ia terbuka tabel utama DISEMBUNYIKAN — fokus pada
// grid yang diklik; menutup panel mengembalikan tabel. Catatan pengembang di
// layar (label bersilang OQ-TCO-05, "simpan dulu") dibuang.
//
// ⚠️ OQ-TCO-05: label tahun bersilang antara grid dan form di korpus; keduanya
// dibawa apa adanya.
//
// [keputusan work owner 11-10-2026 "SEMUA HEADER BUAT FILTERNYA"] Baris saringan
// di bawah keenam kepala kolom, disaring SERVER (daftar berhalaman): tahun dan
// tanggal diketik (Contains), Treaty Group dan Reinsurance Type dipilih.

import { useCallback, useEffect, useRef, useState } from 'react'

import { TAHUN_TCO } from '../labels'
import PanelKlausulTahun from '../components/PanelKlausulTahun'
import PanelKontrakTahun from '../components/PanelKontrakTahun'
import { PROPORSI_TCO, labelProporsi } from '../proporsi'
import { Field, FieldTanggal, Gagal, Halaman, Kosong, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { formatDate } from '../../../../inti/frontend/lib/format'
import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import {
  ambilAkhirBawaanTahun,
  ambilGrupTreaty,
  ambilTahunTreaty,
  simpanTahunTreaty,
  type GrupTreaty,
  type HalamanTahunTreaty,
  type SaringTahunTreaty,
  type TahunTreaty,
  type TahunTreatyMasuk,
} from '../api'

/** Isian form — nama medan mengikuti `InputTreatyYear.*`. */
export interface FormTahun {
  id: string
  treatyYear: string
  underwritingYear: string
  treatyGroupId: string
  treatyGroupName: string
  proportion: string
  /** Bentuk apa pun yang `keInputTanggal` kenal; dikirim sebagai YYYY-MM-DD. */
  startDate: string
  endDate: string
  /** Hanya dibaca — `Modified Date` b9097 dan `Username` b9282. */
  tglUpdate: string
  userId: string
}

/** Form baris BARU — `NewInputTreatyYear_Act` mengosongkan ID/TreatyYear/TreatyGroupID (b469–b512). */
export function formKosong(): FormTahun {
  return {
    id: '', treatyYear: '', underwritingYear: '', treatyGroupId: '', treatyGroupName: '',
    proportion: '', startDate: '', endDate: '', tglUpdate: '', userId: '',
  }
}

/** Form dari baris terpilih — `SetTreatyYear_Act` menyalin seluruh medan (b1110–b1353). */
export function formDari(t: TahunTreaty): FormTahun {
  return {
    id: t.id, treatyYear: t.treatyYear, underwritingYear: t.underwritingYear,
    treatyGroupId: t.treatyGroupId, treatyGroupName: t.treatyGroupName, proportion: t.proportion,
    startDate: t.startDate, endDate: t.endDate, tglUpdate: t.tglUpdate, userId: t.userId,
  }
}

/**
 * Badan permintaan dari form. Tanggal diseragamkan ke YYYY-MM-DD; teks yang
 * bukan tanggal dikirim APA ADANYA supaya server yang menolaknya dengan
 * pesan yang tepat — layar tidak menebak.
 */
export function keMasuk(f: FormTahun): TahunTreatyMasuk {
  const tanggal = (v: string): string => {
    const t = v.trim()
    if (t === '') return ''
    const iso = keInputTanggal(t)
    return iso === '' ? t : iso
  }
  return {
    id: f.id,
    treatyYear: f.treatyYear.trim(),
    underwritingYear: f.underwritingYear.trim(),
    treatyGroupId: f.treatyGroupId,
    treatyGroupName: f.treatyGroupName,
    proportion: f.proportion,
    startDate: tanggal(f.startDate),
    endDate: tanggal(f.endDate),
  }
}

/** Tahun (YYYY) sebuah tanggal isian; kosong bila belum tanggal sah. */
export function tahunDari(tanggal: string): string {
  const iso = keInputTanggal(tanggal.trim())
  return iso === '' ? '' : iso.slice(0, 4)
}

/**
 * Start Date diisi. Form BARU [keputusan work owner 30-09-2026]: Underwriting
 * Year dan Transaction Year mengambil tahun Start Date — keduanya tetap dapat
 * diubah sesudahnya. Form yang sudah ber-ID tidak diisi ulang.
 */
export function isiDariMulai(f: FormTahun, mulai: string): FormTahun {
  const tahun = f.id === '' ? tahunDari(mulai) : ''
  return tahun === '' ? { ...f, startDate: mulai } : { ...f, startDate: mulai, treatyYear: tahun, underwritingYear: tahun }
}

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function selTahun(v: string): string {
  return v.trim() === '' ? '—' : v
}

/** Nama grup dari daftar master untuk ID terpilih; kosong bila tidak dikenal. */
export function namaGrup(daftar: readonly GrupTreaty[], id: string): string {
  return daftar.find((g) => g.id === id)?.treatyGroupName ?? ''
}

const UKURAN = 20

/** Saringan kosong - semua baris. */
export const SARING_KOSONG: SaringTahunTreaty = {
  underwritingYear: '', treatyYear: '', startDate: '', endDate: '', treatyGroupId: '', proportion: '',
}

/** Ada isian saringan terisi. */
export function adaSaring(s: SaringTahunTreaty): boolean {
  return Object.values(s).some((v) => v.trim() !== '')
}

/** Jeda sesudah ketikan terakhir sebelum daftar dimuat ulang. */
const JEDA_SARING_MS = 300

/** Isian saringan ketik, urut kolom grid (tahun: angka; tanggal: bentuk tampil DD-MM-YYYY). */
const SARING_KETIK = [
  { k: 'underwritingYear', kolom: TAHUN_TCO.kolomUnderwritingYear, isi: TAHUN_TCO.saring, angka: true },
  { k: 'treatyYear', kolom: TAHUN_TCO.kolomTransactionYear, isi: TAHUN_TCO.saring, angka: true },
  { k: 'startDate', kolom: TAHUN_TCO.kolomStartDate, isi: TAHUN_TCO.saringTanggal, angka: false },
  { k: 'endDate', kolom: TAHUN_TCO.kolomEndDate, isi: TAHUN_TCO.saringTanggal, angka: false },
] as const

/** Panel yang dibuka dari satu baris tahun — `ReinsType` atau `List Description`. */
export type RinciTahun = { jenis: 'kontrak' | 'klausul'; tahun: TahunTreaty }

export default function InboxTreatyContract() {
  const [hal, setHal] = useState<HalamanTahunTreaty | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)
  const [form, setForm] = useState<FormTahun | null>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [info, setInfo] = useState<string | null>(null)
  const [grup, setGrup] = useState<GrupTreaty[]>([])
  const [galatGrup, setGalatGrup] = useState<unknown>(null)
  // Tiket 04: tombol `ReinsType` b20778 membuka editor kontrak tahun itu
  // (Pega: `BrowseReinsTypeYear` + `showHarness` popup `InboxTreatyContractReinsType`).
  // Tiket 08: tombol `List Description` b22196 membuka layar klausul tahun itu
  // (Pega: `showHarness` `InboxTreatyContractDescription`). SATU keadaan untuk
  // keduanya: membuka yang satu menutup yang lain.
  const [rinci, setRinci] = useState<RinciTahun | null>(null)
  // Posisi gulir tabel saat panel dibuka - dikembalikan saat panel ditutup.
  const gulirTabel = useRef(0)
  // Urutan permintaan End Date bawaan - jawaban basi dibuang.
  const urutanAkhir = useRef(0)
  // Saringan per kolom: `saring` = isian layar, `saringKirim` = yang dimuat (sesudah jeda ketik).
  const [saring, setSaring] = useState<SaringTahunTreaty>(SARING_KOSONG)
  const [saringKirim, setSaringKirim] = useState<SaringTahunTreaty>(SARING_KOSONG)
  // Urutan permintaan daftar - jawaban basi (saringan yang sudah diganti) dibuang.
  const urutanMuat = useRef(0)

  const muat = useCallback(async (h: number, s: SaringTahunTreaty) => {
    const ke = ++urutanMuat.current
    setSibuk(true)
    setGalat(null)
    try {
      const d = await ambilTahunTreaty(h, UKURAN, s)
      if (ke === urutanMuat.current) setHal(d)
    } catch (e) {
      // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong.
      if (ke !== urutanMuat.current) return
      setGalat(e)
      setHal(null)
    } finally {
      if (ke === urutanMuat.current) setSibuk(false)
    }
  }, [])

  useEffect(() => {
    void muat(halaman, saringKirim)
  }, [halaman, saringKirim, muat])

  // Saringan berganti: dimuat sesudah jeda ketik, mulai dari halaman pertama.
  useEffect(() => {
    const t = setTimeout(() => {
      setSaringKirim(saring)
      setHalaman(1)
    }, JEDA_SARING_MS)
    return () => clearTimeout(t)
  }, [saring])

  // Master grup treaty untuk pilihan saringan Treaty Group (form memuatnya ulang saat dibuka).
  useEffect(() => {
    let hidup = true
    ambilGrupTreaty()
      .then((d) => {
        if (hidup) setGrup(d.daftar)
      })
      .catch((e: unknown) => {
        if (hidup) setGalatGrup(e)
      })
    return () => {
      hidup = false
    }
  }, [])

  const ubahSaring = (k: keyof SaringTahunTreaty, v: string) => setSaring((s) => ({ ...s, [k]: v }))

  // Master grup treaty dimuat saat form dibuka — kosong dinyatakan (503).
  useEffect(() => {
    if (form === null) return
    let hidup = true
    void (async () => {
      try {
        const d = await ambilGrupTreaty()
        if (hidup) setGrup(d.daftar)
      } catch (e) {
        if (hidup) setGalatGrup(e)
      }
    })()
    return () => {
      hidup = false
    }
  }, [form === null])

  async function simpan(): Promise<void> {
    if (form === null || menyimpan) return
    setMenyimpan(true)
    setGalatSimpan(null)
    try {
      await simpanTahunTreaty(keMasuk(form))
      setForm(null)
      setInfo(TAHUN_TCO.tersimpan)
      await muat(halaman, saringKirim)
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setMenyimpan(false)
    }
  }

  const ubah = (k: keyof FormTahun) => (v: string) => {
    setForm((f) => (f === null ? f : { ...f, [k]: v }))
  }

  // Form BARU: Start Date mengisi tahun (di sini) dan End Date (server, aturan
  // yang sama dengan kontrak); semuanya tetap dapat diubah. Hanya jawaban
  // TERAKHIR yang dipakai, dan hanya bila Start Date belum berubah lagi.
  function ubahMulai(v: string): void {
    setForm((f) => (f === null ? f : isiDariMulai(f, v)))
    const iso = keInputTanggal(v.trim())
    if (form === null || form.id !== '' || iso === '') return
    const ke = ++urutanAkhir.current
    // End Date yang diketik pemakai SELAMA permintaan berjalan tidak ditimpa.
    const akhirSaatMinta = form.endDate
    ambilAkhirBawaanTahun(iso)
      .then((akhir) => {
        if (ke !== urutanAkhir.current) return
        setForm((f) =>
          f === null || f.id !== '' || f.startDate !== v || f.endDate !== akhirSaatMinta ? f : { ...f, endDate: akhir },
        )
      })
      .catch((e: unknown) => {
        if (ke !== urutanAkhir.current) return
        setGalatSimpan(e)
      })
  }

  function bukaRinci(r: RinciTahun): void {
    if (typeof window !== 'undefined') gulirTabel.current = window.scrollY
    setRinci(r)
    // Tabel panjang tersembunyi: mulai dari atas panel, bukan dari tengah layar.
    if (typeof window !== 'undefined') window.scrollTo({ top: 0 })
  }

  // Fokus pada panel yang diklik: tabel utama, form, dan penomoran halaman
  // TIDAK dirender; `Tutup` panel mengembalikannya (state tabel tetap).
  if (rinci !== null) {
    const tutup = () => {
      setRinci(null)
      // Kembali ke baris yang tadi diklik, sesudah tabel dirender ulang.
      const y = gulirTabel.current
      if (typeof window !== 'undefined') window.requestAnimationFrame(() => window.scrollTo({ top: y }))
    }
    return (
      <section className="inbox tco">
        {rinci.jenis === 'kontrak' ? (
          <PanelKontrakTahun key={rinci.tahun.id} tahun={rinci.tahun} onTutup={tutup} />
        ) : (
          <PanelKlausulTahun key={rinci.tahun.id} tahun={rinci.tahun} onTutup={tutup} />
        )}
      </section>
    )
  }

  return (
    <section className="inbox tco">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TAHUN_TCO.judul}</h2>
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setGalatSimpan(null)
            setInfo(null)
            setForm(formKosong())
          }}
        >
          {TAHUN_TCO.add}
        </button>
      </header>
      {info !== null && <p role="status">{info}</p>}

      {form !== null && (
        <section className="panel">
          <h3 className="panel__title">{TAHUN_TCO.inputNewData}</h3>
          {galatSimpan !== null && <Gagal galat={galatSimpan} />}
          {galatGrup !== null && <Gagal galat={galatGrup} />}
          <div className="form-grid">
            <Field label={TAHUN_TCO.formId} value={form.id} onChange={() => undefined} readOnly />
            <Pilih
              label={TAHUN_TCO.formTreatyGroup}
              value={form.treatyGroupId}
              onChange={(v) => {
                setForm((f) =>
                  f === null ? f : { ...f, treatyGroupId: v, treatyGroupName: namaGrup(grup, v) },
                )
              }}
              opsi={grup.map((g) => ({ value: g.id, label: g.treatyGroupName }))}
              required
            />
            {/* OQ-TCO-04 (keputusan work owner 30-09-2026): dua nilai, bukan master jenis reasuransi. */}
            <Pilih
              label={TAHUN_TCO.formReinsuranceType}
              value={form.proportion}
              onChange={ubah('proportion')}
              opsi={PROPORSI_TCO.map((p) => ({ value: p.value, label: p.label }))}
            />
            <FieldTanggal label={TAHUN_TCO.formStartDate} value={form.startDate} onChange={ubahMulai} />
            <FieldTanggal label={TAHUN_TCO.formEndDate} value={form.endDate} onChange={ubah('endDate')} />
            <Field label={TAHUN_TCO.formUnderwritingYear} value={form.treatyYear} onChange={ubah('treatyYear')} required />
            <Field label={TAHUN_TCO.formTransactionYear} value={form.underwritingYear} onChange={ubah('underwritingYear')} />
            <Field label={TAHUN_TCO.formModifiedDate} value={form.tglUpdate} onChange={() => undefined} readOnly />
            <Field label={TAHUN_TCO.formUsername} value={form.userId} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button type="button" className="btn btn--primary" disabled={menyimpan} onClick={() => void simpan()}>
              {TAHUN_TCO.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalatSimpan(null)
              }}
            >
              {TAHUN_TCO.cancel}
            </button>
          </div>
        </section>
      )}

      {sibuk && hal === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {hal !== null && hal.baris.length === 0 && !adaSaring(saring) && !adaSaring(saringKirim) && (
        <Kosong pesan={TAHUN_TCO.kosong} />
      )}
      {hal !== null && (hal.baris.length > 0 || adaSaring(saring) || adaSaring(saringKirim)) && (
        <div className="tco-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{TAHUN_TCO.kolomUnderwritingYear}</th>
                <th>{TAHUN_TCO.kolomTransactionYear}</th>
                <th>{TAHUN_TCO.kolomStartDate}</th>
                <th>{TAHUN_TCO.kolomEndDate}</th>
                <th>{TAHUN_TCO.kolomTreatyGroup}</th>
                <th>{TAHUN_TCO.kolomReinsuranceType}</th>
                <th className="table__actions" />
              </tr>
              <tr className="tco-saring">
                {SARING_KETIK.map((c) => (
                  <th key={c.k}>
                    <input
                      className="field__input"
                      inputMode={c.angka ? 'numeric' : undefined}
                      value={saring[c.k]}
                      placeholder={c.isi}
                      aria-label={`${TAHUN_TCO.saring} ${c.kolom}`}
                      onChange={(e) => ubahSaring(c.k, e.target.value)}
                    />
                  </th>
                ))}
                <th>
                  <select
                    className="field__input"
                    value={saring.treatyGroupId}
                    aria-label={`${TAHUN_TCO.saring} ${TAHUN_TCO.kolomTreatyGroup}`}
                    onChange={(e) => ubahSaring('treatyGroupId', e.target.value)}
                  >
                    <option value="">{TAHUN_TCO.saringSemua}</option>
                    {grup.map((g) => (
                      <option key={g.id} value={g.id}>
                        {g.treatyGroupName}
                      </option>
                    ))}
                  </select>
                </th>
                <th>
                  <select
                    className="field__input"
                    value={saring.proportion}
                    aria-label={`${TAHUN_TCO.saring} ${TAHUN_TCO.kolomReinsuranceType}`}
                    onChange={(e) => ubahSaring('proportion', e.target.value)}
                  >
                    <option value="">{TAHUN_TCO.saringSemua}</option>
                    {PROPORSI_TCO.map((p) => (
                      <option key={p.value} value={p.value}>
                        {p.label}
                      </option>
                    ))}
                  </select>
                </th>
                <th className="table__actions">
                  {adaSaring(saring) && (
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => setSaring(SARING_KOSONG)}>
                      {TAHUN_TCO.saringHapus}
                    </button>
                  )}
                </th>
              </tr>
            </thead>
            <tbody>
              {hal.baris.length === 0 && (
                <tr>
                  <td colSpan={7} className="tco-saring-kosong">
                    {TAHUN_TCO.saringKosong}
                  </td>
                </tr>
              )}
              {hal.baris.map((b) => (
                <tr key={b.id} className="inbox__baris">
                  <td>{selTahun(b.underwritingYear)}</td>
                  <td>{selTahun(b.treatyYear)}</td>
                  <td>{selTahun(formatDate(b.startDate))}</td>
                  <td>{selTahun(formatDate(b.endDate))}</td>
                  <td>{selTahun(b.treatyGroupName)}</td>
                  <td>{selTahun(labelProporsi(b.proportion))}</td>
                  <td className="table__actions">
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setGalatSimpan(null)
                        setInfo(null)
                        setForm(formDari(b))
                      }}
                    >
                      {TAHUN_TCO.edit}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        bukaRinci({ jenis: 'kontrak', tahun: b })
                      }}
                    >
                      {TAHUN_TCO.reinsType}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        bukaRinci({ jenis: 'klausul', tahun: b })
                      }}
                    >
                      {TAHUN_TCO.listDescription}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {hal !== null && (
        <Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={setHalaman} />
      )}
    </section>
  )
}
