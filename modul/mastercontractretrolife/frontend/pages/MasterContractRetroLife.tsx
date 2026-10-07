// Halaman awal modul - tiket 01. `Section/GridRetrocessionLife.xml` (judul b1017, menyertakan
// `InputRetrocessionLife` b2475) → `InputRetrocessionLife.xml` (+ form `InputDtlRetrocessionLife`).
//
// ⭐ Pintu masuk: keempat harness korpus dibuka DARI DALAM (`showHarness` popup) - tidak satu pun pintu
// masuk; `GridRetrocessionLife` tidak dirujuk rule mana pun di korpus modul (PARITAS §1).
//
// Urutan XML: form `Input New Data` (wadah b774 `DATASHOW3 = 1`) → wadah grid b8631 (`IsFire`
// ber-`ALWAYS` → selalu tampil, RALAT R6): tombol b8888 (label sel `End Period` - TIDAK ditampilkan
// sejak 02-10-2026, keputusan work owner; teks `Add`, tooltip
// `Add New Data` → `NewInputTreatyYear_Life_Act`), `pyGridPaginator` b9192, grid b9375 RD
// `BrowseTreatyYear_Life_RD` (urut `.ID ASC` - server; SEMUA baris, tanpa paginasi sejak 04-10-2026) dengan tombol baris `Edit`
// b11645 (`SetTreatyYearLife_Act`) dan `ReinsType` b11988 (`showHarness` `InboxRetroLimitReinsurers`).
//
// ⛔ Identitas tidak pernah diketik: `ID` hanya dibaca (tampil bila tidak kosong, `NOTBLANK` b1798);
// baris baru dikirim TANPA id (POST), ubah lewat PUT /{id}. ⛔ Tahun treaty tidak dapat dihapus -
// XML tidak punya tombol hapus di layar ini.

import { useCallback, useEffect, useState } from 'react'

import '../mcrl.css'
import { ambilTahun, simpanTahun, type TahunMasuk, type TahunTreaty } from '../api'
import PanelKontrak from '../components/PanelKontrak'
import { TAHUN_MCRL, UMUM_MCRL } from '../labels'
import { keTanggalKabel, operatorKini, sel, selTanggal, selWaktu, waktuKini } from '../tampilan'
import { Field, FieldTanggal, Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'

/** Isian form - `InputTreatyYear.*`. */
export interface FormTahun {
  id: string
  /** `TRANSACTION YEAR`. */
  treatyYear: string
  underwritingYear: string
  startDate: string
  endDate: string
  /** `Modified Date` / `Inputor` (ro). */
  tglUpdate: string
  userId: string
}

/**
 * `NewInputTreatyYear_Life_Act` b393: kosongkan `ID, TREATYYEAR, UNDERWRITINGYEAR, STARTDATE, ENDDATE`,
 * `TGLUPDATE ← @CurrentDateTime()`; `Inputor` = `OperatorID.pyUserName` (`DATASHOW2 = 0`, b3783).
 */
export function formTahunBaru(operator: string, kini: string): FormTahun {
  return { id: '', treatyYear: '', underwritingYear: '', startDate: '', endDate: '', tglUpdate: kini, userId: operator }
}

/**
 * `SetTreatyYearLife_Act` b556: salin `ID, TREATYYEAR, UNDERWRITINGYEAR, STARTDATE, ENDDATE`,
 * `USERID ← OperatorID`, `TGLUPDATE ← @CurrentDateTime()` (`Inputor` b3967 = `USERID`).
 */
export function formTahunDari(t: TahunTreaty, operator: string, kini: string): FormTahun {
  return {
    id: t.id, treatyYear: t.treatyYear, underwritingYear: t.underwritingYear,
    startDate: t.startDate, endDate: t.endDate, tglUpdate: kini, userId: operator,
  }
}

/** Badan simpan; tanggal diseragamkan ke `YYYY-MM-DD`, teks bukan tanggal dikirim apa adanya. */
/**
 * START DATE sesudah END DATE (tanggal mundur) - Save dikunci dan pesannya
 * tampil di END DATE (keputusan work owner 04-10-2026). Sama hari diterima.
 * Server menolak hal yang sama (`PesanTanggalMundur`).
 */
export function tanggalMundur(f: FormTahun): boolean {
  const mulai = keTanggalKabel(f.startDate)
  const akhir = keTanggalKabel(f.endDate)
  const iso = /^\d{4}-\d{2}-\d{2}$/
  return iso.test(mulai) && iso.test(akhir) && mulai > akhir
}

export function keTahunMasuk(f: FormTahun): TahunMasuk {
  return {
    id: f.id,
    treatyYear: f.treatyYear.trim(),
    underwritingYear: f.underwritingYear.trim(),
    startDate: keTanggalKabel(f.startDate),
    endDate: keTanggalKabel(f.endDate),
  }
}

export default function MasterContractRetroLife() {
  const [daftar, setDaftar] = useState<TahunTreaty[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [form, setForm] = useState<FormTahun | null>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  // Tombol `ReinsType` b11988: popup `InboxRetroLimitReinsurers` untuk tahun itu.
  const [kontrakDari, setKontrakDari] = useState<TahunTreaty | null>(null)

  const muat = useCallback(async () => {
    try {
      const d = await ambilTahun()
      setDaftar(d.daftar)
      setGalat(null)
    } catch (e) {
      // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong.
      setGalat(e)
    }
  }, [])

  useEffect(() => {
    void muat()
  }, [muat])

  function buka(f: FormTahun): void {
    setGalatSimpan(null)
    setForm(f)
  }

  function ubah(k: 'treatyYear' | 'underwritingYear' | 'startDate' | 'endDate') {
    return (v: string) => {
      setForm((f) => (f === null ? f : { ...f, [k]: v }))
    }
  }

  async function simpan(): Promise<void> {
    if (form === null || menyimpan) return
    setMenyimpan(true)
    setGalatSimpan(null)
    try {
      await simpanTahun(keTahunMasuk(form))
      // `SaveTreatyYearLife_Act` langkah 7 b1290: `DATASHOW3 = ""` - form tertutup, grid dimuat ulang.
      setForm(null)
      await muat()
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setMenyimpan(false)
    }
  }

  // Popup di atas halaman awal: tabel tidak dirender, keadaannya (halaman, form) tetap.
  if (kontrakDari !== null) {
    return (
      <section className="inbox mcrl">
        <PanelKontrak
          key={kontrakDari.id}
          tahun={kontrakDari}
          onTutup={() => {
            setKontrakDari(null)
          }}
        />
      </section>
    )
  }

  const semua = daftar ?? []

  return (
    <section className="inbox mcrl">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TAHUN_MCRL.judul}</h2>
      </header>

      {form !== null && (
        <section className="panel">
          <h3 className="panel__title">{TAHUN_MCRL.inputNewData}</h3>
          {galatSimpan !== null && <Gagal galat={galatSimpan} />}
          <div className="form-grid">
            {form.id !== '' && <Field label={TAHUN_MCRL.formId} value={form.id} onChange={() => undefined} readOnly />}
            <Field label={TAHUN_MCRL.formUnderwritingYear} value={form.underwritingYear} onChange={ubah('underwritingYear')} required />
            <Field label={TAHUN_MCRL.formTransactionYear} value={form.treatyYear} onChange={ubah('treatyYear')} required />
            <FieldTanggal label={TAHUN_MCRL.formStartDate} value={form.startDate} onChange={ubah('startDate')} required />
            <FieldTanggal
              label={TAHUN_MCRL.formEndDate}
              value={form.endDate}
              onChange={ubah('endDate')}
              required
              error={tanggalMundur(form) ? TAHUN_MCRL.tanggalMundur : undefined}
            />
            <Field label={TAHUN_MCRL.formModifiedDate} value={selWaktu(form.tglUpdate)} onChange={() => undefined} readOnly />
            <Field label={TAHUN_MCRL.formInputor} value={form.userId} onChange={() => undefined} readOnly />
          </div>
          <div className="aksi-baris">
            <button
              type="button"
              className="btn btn--primary"
              disabled={menyimpan || tanggalMundur(form)}
              onClick={() => void simpan()}
            >
              {TAHUN_MCRL.save}
            </button>{' '}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setForm(null)
                setGalatSimpan(null)
              }}
            >
              {TAHUN_MCRL.cancel}
            </button>
          </div>
        </section>
      )}

      <div className="aksi-baris mcrl-aksi-grid">
        {/* Label sel `End Period` (b8927, `TAHUN_MCRL.labelSelAdd`) SENGAJA tidak ditampilkan -
            keputusan work owner 02-10-2026 ("tulisan end period di hapus"); RALAT 02-10-2026. */}
        <button
          type="button"
          className="btn btn--primary"
          title={TAHUN_MCRL.tooltipAdd}
          onClick={() => buka(formTahunBaru(operatorKini(), waktuKini(new Date())))}
        >
          {TAHUN_MCRL.add}
        </button>
      </div>
      {daftar === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {daftar !== null && semua.length === 0 && <Kosong pesan={UMUM_MCRL.kosong} />}
      {semua.length > 0 && (
        <div className="mcrl-tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{TAHUN_MCRL.kolomId}</th>
                <th>{TAHUN_MCRL.kolomUnderwritingYear}</th>
                <th>{TAHUN_MCRL.kolomTransactionYear}</th>
                <th>{TAHUN_MCRL.kolomStartDate}</th>
                <th>{TAHUN_MCRL.kolomEndDate}</th>
                <th className="table__actions" />
              </tr>
            </thead>
            <tbody>
              {semua.map((t) => (
                <tr key={t.id} className="inbox__baris">
                  <td>{sel(t.id)}</td>
                  <td>{sel(t.underwritingYear)}</td>
                  <td>{sel(t.treatyYear)}</td>
                  <td>{selTanggal(t.startDate)}</td>
                  <td>{selTanggal(t.endDate)}</td>
                  <td className="table__actions">
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      title={TAHUN_MCRL.tooltipEdit}
                      onClick={() => buka(formTahunDari(t, operatorKini(), waktuKini(new Date())))}
                    >
                      {TAHUN_MCRL.edit}
                    </button>{' '}
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setKontrakDari(t)
                      }}
                    >
                      {TAHUN_MCRL.reinsType}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}
