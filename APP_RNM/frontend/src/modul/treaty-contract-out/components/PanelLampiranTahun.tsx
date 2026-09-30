// Panel lampiran tahun treaty — tiket 12 Treaty Contract Out (FITUR BARU).
//
// Meniru `Section/GridTreatyArrangementAttachment.xml` yang disertakan
// `Section/InputTreatyContract.xml` b13074 di bawah judul `Attachment for`
// b11721: `Add attachment` b578, `Refresh` b1023, `Download All` b2659,
// kolom `File Name` b3032 (tautan → `TreatyOutDownloadOne`) dan `Type` b3170
// (`.pyCategory`), `Delete` b3897.
//
// ⚠️ Penyimpangan sadar 9: jalur lampiran Pega belum rampung. Status
// terkirim / tertunda / gagal, `Ulangi`, dan `Periksa keselarasan` adalah
// tambahan tiket 12 (AC 55, 58, 61) — kosakata kami.
//
// ⛔ Unduhan lewat `unduhBerkasBeridentitas` (fetch berheader identitas),
// BUKAN tautan biasa: rute unduh bergerbang identitas.

import { useCallback, useEffect, useState } from 'react'

import { LAMPIRAN_TCO } from '../labels'
import { Gagal, Kosong, Memuat, Pilih } from '../../../../../inti/frontend/components/ui/dasar'
import {
  ambilKategoriLampiranTCO,
  ambilLampiranTahun,
  hapusLampiranTahun,
  jalurIsiLampiran,
  jalurSemuaLampiran,
  periksaSelarasLampiran,
  ulangiLampiranTahun,
  unggahLampiranTahun,
  type LampiranTahun,
  type StatusLampiranTahun,
  type TemuanSelarasLampiran,
} from '../api'
import { unduhBerkasBeridentitas } from '../../../../../inti/frontend/klien'

/** Label status di layar. */
export function labelStatusLampiran(s: StatusLampiranTahun): string {
  switch (s) {
    case 'terkirim':
      return LAMPIRAN_TCO.statusTerkirim
    case 'gagal':
      return LAMPIRAN_TCO.statusGagal
    default:
      return LAMPIRAN_TCO.statusTertunda
  }
}

/** Hanya berkas yang sudah dipastikan ada di penyimpanan yang dapat diunduh. */
export function bolehUnduh(l: LampiranTahun): boolean {
  return l.status === 'terkirim'
}

/** Yang belum terkirim dapat diulang; backend menolak bila tidak ada yang diperbaiki. */
export function bolehUlangi(l: LampiranTahun): boolean {
  return l.status !== 'terkirim'
}

/** `Download All` hanya berarti bila ada yang terkirim. */
export function adaTerkirim(daftar: LampiranTahun[]): boolean {
  return daftar.some(bolehUnduh)
}

export default function PanelLampiranTahun({ tahunID }: { tahunID: string }) {
  const [daftar, setDaftar] = useState<LampiranTahun[] | null>(null)
  const [kategori, setKategori] = useState<string[]>([])
  const [pilihan, setPilihan] = useState('')
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [peringatan, setPeringatan] = useState<string | null>(null)
  const [temuan, setTemuan] = useState<TemuanSelarasLampiran[] | null>(null)

  const muat = useCallback(async () => {
    try {
      setDaftar(await ambilLampiranTahun(tahunID))
    } catch (e) {
      setGalat(e)
    }
  }, [tahunID])

  useEffect(() => {
    void muat()
    ambilKategoriLampiranTCO()
      .then(setKategori)
      .catch((e: unknown) => {
        setGalat(e)
      })
  }, [muat])

  async function jalankan(kerja: () => Promise<string | undefined>): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    setPeringatan(null)
    try {
      const p = await kerja()
      if (p !== undefined && p !== '') setPeringatan(p)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  const baris = daftar ?? []

  return (
    <section className="panel">
      <h3 className="panel__title">
        {LAMPIRAN_TCO.attachmentFor} {tahunID}
      </h3>
      <p className="polis__catatan">{LAMPIRAN_TCO.forTreatyContractOut}</p>
      {galat !== null && <Gagal galat={galat} />}
      {peringatan !== null && <p role="alert">{peringatan}</p>}

      <div className="form-grid">
        <Pilih
          label={LAMPIRAN_TCO.kolomType}
          value={pilihan}
          onChange={setPilihan}
          opsi={kategori.map((k) => ({ value: k, label: k }))}
          kosong={LAMPIRAN_TCO.pilihKategori}
          required
        />
        <label className="field">
          <span className="field__label">{LAMPIRAN_TCO.addAttachment}</span>
          <input
            type="file"
            disabled={sibuk || pilihan === ''}
            onChange={(e) => {
              const f = e.target.files?.[0]
              // Kotak berkas dikosongkan: memilih berkas yang sama dua kali
              // tidak memicu `change` bila nilainya masih tertinggal.
              e.target.value = ''
              if (f === undefined) return
              void jalankan(async () => (await unggahLampiranTahun(tahunID, f, pilihan)).peringatan)
            }}
          />
        </label>
      </div>
      <div className="aksi-baris">
        <button type="button" className="btn btn--ghost" disabled={sibuk} onClick={() => void jalankan(async () => undefined)}>
          {LAMPIRAN_TCO.refresh}
        </button>{' '}
        <button
          type="button"
          className="btn btn--ghost"
          disabled={sibuk || !adaTerkirim(baris)}
          onClick={() =>
            void jalankan(async () => {
              await unduhBerkasBeridentitas(jalurSemuaLampiran(tahunID), `treaty-year-attachments-${tahunID}.zip`)
              return undefined
            })
          }
        >
          {LAMPIRAN_TCO.downloadAll}
        </button>{' '}
        <button
          type="button"
          className="btn btn--ghost"
          disabled={sibuk}
          onClick={() =>
            void jalankan(async () => {
              setTemuan(await periksaSelarasLampiran(tahunID))
              return undefined
            })
          }
        >
          {LAMPIRAN_TCO.periksaSelaras}
        </button>
      </div>

      {temuan !== null && temuan.length === 0 && <p role="status">{LAMPIRAN_TCO.selarasBersih}</p>}
      {temuan !== null && temuan.length > 0 && (
        <ul role="alert">
          {temuan.map((t) => (
            <li key={t.lampiranId}>
              {t.lampiranId} {t.fileName}: {t.masalah} ({t.perbaikan === 'ulangi' ? LAMPIRAN_TCO.perbaikanUlangi : LAMPIRAN_TCO.perbaikanHapus})
            </li>
          ))}
        </ul>
      )}

      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && baris.length === 0 && <Kosong pesan={LAMPIRAN_TCO.kosong} />}
      {baris.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th scope="col">{LAMPIRAN_TCO.kolomFileName}</th>
              <th scope="col">{LAMPIRAN_TCO.kolomType}</th>
              <th scope="col">{LAMPIRAN_TCO.kolomStatus}</th>
              <th scope="col" aria-label={LAMPIRAN_TCO.kolomAksi} />
            </tr>
          </thead>
          <tbody>
            {baris.map((l) => (
              <tr key={l.id}>
                <td>
                  {bolehUnduh(l) ? (
                    <button
                      type="button"
                      className="btn btn--link"
                      disabled={sibuk}
                      onClick={() =>
                        void jalankan(async () => {
                          await unduhBerkasBeridentitas(jalurIsiLampiran(tahunID, l.id), l.fileName)
                          return undefined
                        })
                      }
                    >
                      {l.fileName}
                    </button>
                  ) : (
                    l.fileName
                  )}
                </td>
                <td>{l.category}</td>
                <td title={l.galat !== '' ? l.galat : undefined}>
                  {labelStatusLampiran(l.status)}
                  {l.galat !== '' && <small> — {l.galat}</small>}
                </td>
                <td>
                  {bolehUlangi(l) && (
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      disabled={sibuk}
                      onClick={() => void jalankan(async () => (await ulangiLampiranTahun(tahunID, l.id)).peringatan)}
                    >
                      {LAMPIRAN_TCO.ulangi}
                    </button>
                  )}{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    disabled={sibuk}
                    onClick={() => void jalankan(async () => (await hapusLampiranTahun(tahunID, l.id)).peringatan)}
                  >
                    {LAMPIRAN_TCO.delete}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
