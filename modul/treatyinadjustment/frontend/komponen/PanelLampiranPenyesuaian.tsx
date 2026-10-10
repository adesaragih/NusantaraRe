// Panel Attachment layar Adjustment — HIDUP, sama persis dengan Treaty In
// (8 Oktober 2026).
//
// ⛔ DISALIN dari `modul/treatyin/frontend/components/PanelLampiran.tsx`,
// tidak diimpor (`inti/frontend/lapisan.guard.test.ts`). Harness Adjustment
// memakai Section yang SAMA: `InputTreatyInAdjustment.xml` `pyInclude
// WorkAttachments` → `WorkAttachments.xml` (Category · Count · Upload file ·
// View File · Download All · Refresh · spanduk biru) dan
// `ShowAttachmentTreaty.xml` (File Name · Type · View Office Online · Delete ·
// Change Category/Save). Syarat tampilnya disalin dari ekspor modul INI:
//
//   Upload file / Delete  `TreatyIn.ViewState !='1' || TreatyIn.RevisionState='1'`
//   Change Category       `TreatyIn.StatusAkseptasi != 'Resolve Complete' &&
//                         TreatyIn.StatusAkseptasi != 'Decline'` (wadah), lalu
//                         `StatusDoc.CARI30 != 1` / `= 1` (Change Category ⇄ Save)
//   Download All          `pyVisible ALWAYS` menimpa `pyCondition never`
//   Refresh               `!pyIsMobile`
//   View Office Online    `.pyFileMimeType` xls/xlsx/doc/docx/ppt/pptx
//
// Perbedaannya dengan Treaty In hanya sumber data: panel MEMBACA SENDIRI
// lewat rute baca modul ini (`ambilLampiran`, menamai kode 00007 menurut
// `jenis`), dan tombol tulisnya memanggil rute lampiran Treaty In lewat HTTP.

import { useEffect, useRef, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { pesanGalat } from '../../../../inti/frontend/klien'
import {
  ambilLampiran,
  ambilTautanLampiran,
  hapusLampiran,
  ubahKategoriLampiran,
  unduhLampiran,
  unduhSemuaLampiran,
  unggahLampiran,
  type BarisKategoriLampiran,
  type BarisLampiranWarisan,
  type PanelLampiranTreatyIn,
} from '../api'
import { KOLOM_LAMPIRAN } from '../labels'
import { KOLOM_LIHAT_BERKAS, LAMPIRAN_KELOLA as LAMPIRAN } from '../labelsLampiran'
import { gabungBerkas, unggahBerurutan } from './unggahBerkas'

export interface PropsPanelLampiranPenyesuaian {
  /**
   * Pengenal yang lampirannya DITAMPILKAN — penyesuaian tersimpan, atau
   * kontrak asal bagi draf (salinan lampirannya menunggu Save).
   */
  idKontrak: string
  /** `TreatyIn.ProportionType` — nama Non-Prop untuk kode `00007`. */
  jenis?: string
  /**
   * Draf yang BELUM tersimpan: lampiran yang tampil milik kontrak asal, jadi
   * nol tombol tulis (Upload/Delete/Change Category) — persis Treaty In untuk
   * kontrak yang belum ber-ID.
   */
  draf?: boolean
  /** `TreatyIn.ViewState !='1' || TreatyIn.RevisionState='1'`. */
  bisaUnggah?: boolean
  /** `TreatyIn.StatusAkseptasi` — syarat wadah Change Category. */
  statusAkseptasi?: string
}

/**
 * Panel Attachment layar Adjustment — membaca lampiran pengenal yang tampil
 * (dibaca ULANG bila pengenalnya berganti, mis. draf yang baru tersimpan),
 * lalu menyerahkannya ke `PanelLampiranIsi`.
 */
export default function PanelLampiranPenyesuaian(props: PropsPanelLampiranPenyesuaian) {
  const { idKontrak, jenis = '' } = props
  const [kategori, setKategori] = useState<readonly BarisKategoriLampiran[]>([])
  const [berkas, setBerkas] = useState<readonly BarisLampiranWarisan[]>([])
  const [memuat, setMemuat] = useState(idKontrak !== '')
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    setGalat(null)
    if (idKontrak === '') {
      setKategori([])
      setBerkas([])
      setMemuat(false)
      return
    }
    let dibuang = false
    setMemuat(true)
    ambilLampiran(idKontrak, jenis)
      .then((l) => {
        if (dibuang) return
        setKategori(l.kategori)
        setBerkas(l.berkas)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
      .finally(() => {
        if (!dibuang) setMemuat(false)
      })
    return () => {
      dibuang = true
    }
  }, [idKontrak, jenis])

  return (
    <>
      {galat !== null && <Gagal galat={galat} />}
      <PanelLampiranIsi {...props} kategori={kategori} berkas={berkas} memuat={memuat} />
    </>
  )
}

/**
 * Panel Attachment — Category · Count · Upload file · View File.
 *
 * ⭐ Bentuk modal unggah = Treaty In (bentuk Master Product Name Life): kotak
 * seret-lepas, berkas terpilih bisa dibuang satu-satu, diunggah SATU PER
 * SATU dengan progres `2/5 nama`, berkas gagal tinggal beserta alasannya.
 */
export function PanelLampiranIsi({
  kategori: kategoriAwal,
  berkas: berkasAwal,
  memuat = false,
  idKontrak,
  jenis = '',
  draf = false,
  bisaUnggah = false,
  statusAkseptasi = '',
}: PropsPanelLampiranPenyesuaian & {
  kategori: readonly BarisKategoriLampiran[]
  berkas: readonly BarisLampiranWarisan[]
  /** Pembacaan pertama belum selesai. */
  memuat?: boolean
}) {
  // Pengenal untuk tombol TULIS — kosong bagi draf (lampiran menempel pada
  // ID tersimpan, dan yang tampil milik kontrak asal).
  const idTulis = draf ? '' : idKontrak
  // Isi panel — dari pembacaan pertama, lalu dari Refresh/Attach sendiri.
  const [kategori, setKategori] = useState<readonly BarisKategoriLampiran[]>(kategoriAwal)
  const [berkas, setBerkas] = useState<readonly BarisLampiranWarisan[]>(berkasAwal)
  useEffect(() => {
    setKategori(kategoriAwal)
    setBerkas(berkasAwal)
  }, [kategoriAwal, berkasAwal])
  // Modal `ASM Attach Content` — kategori yang dipilih (`SetkategoriDoc`).
  const [unggahKe, setUnggahKe] = useState<BarisKategoriLampiran | null>(null)
  const [dipilih, setDipilih] = useState<File[]>([])
  const [proses, setProses] = useState<string | null>(null)
  const [gagalUnggah, setGagalUnggah] = useState<{ berkas: File; galat: unknown }[]>([])
  const [seret, setSeret] = useState(false)
  const [galatUnggah, setGalatUnggah] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  // Galat Refresh / unduh / Delete / Change Category.
  const [galatPanel, setGalatPanel] = useState('')
  const galat = (e: unknown) => {
    setGalatPanel(e instanceof Error ? e.message : String(e))
  }
  /** Jawaban rute tulis — panel yang dibaca ulang server. */
  const terapkanPanel = (p: PanelLampiranTreatyIn) => {
    setKategori(p.kategoriLampiran)
    setBerkas(p.lampiran)
  }
  /**
   * `Refresh` (`GetMasterTreatyCategory_Act`) — panel saja, lewat rute baca
   * modul ini; isian form yang belum di-Save tidak tersentuh.
   */
  const segarkan = () => {
    if (idKontrak === '') return
    setSibuk(true)
    setGalatPanel('')
    ambilLampiran(idKontrak, jenis)
      .then((l) => {
        setKategori(l.kategori)
        setBerkas(l.berkas)
      })
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  // ⭐ Modal `View File` — `ShowAttachmentTreaty`.
  // Change Category: `StatusDoc.CARI30` 1 = sel kategori menjadi dropdown.
  const [gantiKategori, setGantiKategori] = useState(false)
  const [kategoriBaru, setKategoriBaru] = useState<Record<string, string>>({})
  const bolehGantiKategori = statusAkseptasi !== 'Resolve Complete' && statusAkseptasi !== 'Decline'
  /** Tautan nama berkas (`DownloadAttachmentTreaty`) — `fetch` beridentitas. */
  const unduh = (b: BarisLampiranWarisan) => {
    setSibuk(true)
    setGalatPanel('')
    unduhLampiran(idKontrak, b.id, b.namaBerkas)
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  /**
   * `View Office Online` — penampil kantor di POPUP berbingkai: popup dibuka
   * lebih dulu, lalu bingkainya diisi URL penampil dari backend. Jawaban untuk
   * popup yang sudah ditutup dibuang (`giliran`).
   */
  const [penampil, setPenampil] = useState<{ nama: string; siap: boolean } | null>(null)
  const bingkai = useRef<HTMLIFrameElement>(null)
  const giliran = useRef(0)
  const tutupPenampil = () => {
    giliran.current++
    setPenampil(null)
  }
  const bukaOffice = (b: BarisLampiranWarisan) => {
    const ke = ++giliran.current
    setPenampil({ nama: b.namaBerkas, siap: false })
    ambilTautanLampiran(idKontrak, b.id, true)
      .then((j) => {
        if (ke !== giliran.current) return
        if (bingkai.current !== null) bingkai.current.src = j.url
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      })
      .catch((e: unknown) => {
        if (ke === giliran.current) tutupPenampil()
        galat(e)
      })
  }
  const hapus = (idLampiran: string) => {
    if (idTulis === '') return
    setSibuk(true)
    setGalatPanel('')
    hapusLampiran(idTulis, idLampiran)
      .then(terapkanPanel)
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  const simpanKategori = () => {
    if (idTulis === '') return
    const perubahan = Object.entries(kategoriBaru).map(([id, kode]) => ({ id, kategori: kode }))
    setSibuk(true)
    setGalatPanel('')
    ubahKategoriLampiran(idTulis, perubahan)
      .then((p) => {
        terapkanPanel(p)
        setGantiKategori(false)
        setKategoriBaru({})
      })
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  const bukaUnggah = (k: BarisKategoriLampiran) => {
    setUnggahKe(k)
    setDipilih([])
    setGagalUnggah([])
    setGalatUnggah(null)
    setProses(null)
  }
  /**
   * Attach (`TreatySaveAttachment`) — SATU berkas per permintaan, berurutan.
   * Tiap jawaban membawa panel yang dibaca ulang, jadi `Count` naik selagi
   * berkas berikutnya dikirim.
   */
  const lampirkan = () => {
    if (unggahKe === null || sibuk || idTulis === '') return
    const kode = unggahKe.kode
    const antre = dipilih
    setSibuk(true)
    setGalatUnggah(null)
    setGagalUnggah([])
    void (async () => {
      try {
        if (antre.length === 0) {
          // Kalimat backend apa adanya: "Tidak ada berkas yang dipilih."
          await unggahLampiran(idTulis, kode, [])
          return
        }
        const gagal = await unggahBerurutan(
          antre,
          async (f) => {
            const h = await unggahLampiran(idTulis, kode, [f])
            terapkanPanel(h)
            // Berkas DITOLAK dijawab 200 dengan `berhasil: false` — ia tetap
            // kegagalan berkas ini.
            const b = h.berkas.at(0)
            if (b !== undefined && !b.berhasil) throw new Error(b.pesan)
          },
          (f, i) => {
            setProses(`${String(i + 1)}/${String(antre.length)} ${f.name}`)
          },
        )
        setGagalUnggah(gagal)
        setDipilih(gagal.map((g) => g.berkas))
        if (gagal.length === 0) setUnggahKe(null)
      } catch (e) {
        setGalatUnggah(e)
      } finally {
        setProses(null)
        setSibuk(false)
      }
    })()
  }
  const [berkasDilihat, setBerkasDilihat] = useState<string | null>(null)
  const kategoriDilihat = kategori.find((k) => k.kode === berkasDilihat) ?? null
  const sibukPanel = sibuk || memuat
  return (
    <Panel judul={LAMPIRAN.judul}>
      {/* Kelas pengait kerapatan panel ini — terpisah dari isian tab. */}
      <div className="tria__lampiran">
        {/* Spanduk biru — ATURAN nama berkas, bukan hiasan. */}
        <span className="tria__spanduk" role="note">
          {LAMPIRAN.spanduk}
        </span>
        {/* `Download All` (`DownloadAll_Act`) lalu `Refresh` — di KANAN atas grid. */}
        <div className="tria__aksi tria__lampiran-aksi">
          <button
            type="button"
            className="btn btn--sm"
            disabled={sibukPanel || idKontrak === ''}
            onClick={() => {
              setGalatPanel('')
              unduhSemuaLampiran(idKontrak).catch((e: unknown) => {
                setGalatPanel(pesanGalat(e) ?? (e instanceof Error ? e.message : LAMPIRAN.gagal))
              })
            }}
          >
            {LAMPIRAN.unduhSemua}
          </button>
          <button type="button" className="btn btn--sm" disabled={sibukPanel || idKontrak === ''} onClick={segarkan}>
            {LAMPIRAN.segarkan}
          </button>
        </div>
        {galatPanel !== '' && berkasDilihat === null && (
          <div className="alert alert--error" role="alert">
            {galatPanel}
          </div>
        )}
        {bisaUnggah && idTulis === '' && <p className="tria__redup">{LAMPIRAN.simpanDulu}</p>}

        <div className="table-wrap">
          <table className="tria__tabel tria__lampiran-grid">
            {/* Lebar kolom dari gambar Pega: Category ±64%, sisanya ±12%. */}
            <colgroup>
              <col style={{ width: '64%' }} />
              <col style={{ width: '12%' }} />
              <col style={{ width: '12%' }} />
              <col style={{ width: '12%' }} />
            </colgroup>
            <thead>
              <tr>
                {KOLOM_LAMPIRAN.map((k, i) => (
                  <th key={k} scope="col" className={i > 0 ? 'tria__lampiran-tengah' : undefined}>
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {memuat && kategori.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_LAMPIRAN.length}>
                    <Memuat />
                  </td>
                </tr>
              )}
              {!memuat && kategori.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_LAMPIRAN.length}>
                    <Kosong pesan={LAMPIRAN.tanpaLampiran} petunjuk={LAMPIRAN.petunjukLampiran} />
                  </td>
                </tr>
              )}
              {kategori.map((k) => (
                <tr key={k.kode !== '' ? k.kode : k.nama}>
                  {/* NAMA katalog; kode hanya bila tanpa nama. */}
                  <td>{k.nama !== '' ? k.nama : k.kode}</td>
                  {/* ⛔ Cacah TIDAK diformat — ia butir, bukan uang. */}
                  <td>{String(k.cacah)}</td>
                  <td>
                    {/* `Upload file` — `TreatyIn.ViewState !='1' ||
                        TreatyIn.RevisionState='1'`; butuh ID tersimpan. */}
                    {bisaUnggah && (
                      <button
                        type="button"
                        className="tria__lampiran-ikon tria__lampiran-ikon--unggah"
                        aria-label={`${LAMPIRAN.unggah} ${k.nama}`}
                        title={LAMPIRAN.unggah}
                        disabled={sibukPanel || idTulis === '' || k.kode === ''}
                        onClick={() => {
                          bukaUnggah(k)
                        }}
                      >
                        <svg viewBox="0 0 24 24" aria-hidden="true" focusable="false">
                          <rect x="1" y="1" width="22" height="22" rx="3" fill="currentColor" />
                          <path d="M12 5.5 7.5 10h3v4.5h3V10h3L12 5.5Zm-5 11v2h10v-2H7Z" fill="#fff" />
                        </svg>
                      </button>
                    )}
                  </td>
                  <td>
                    {/* `View File` — selalu tampil; membuka modal `ShowAttachmentTreaty`. */}
                    <button
                      type="button"
                      className="tria__lampiran-ikon tria__lampiran-ikon--lihat"
                      aria-label={`${LAMPIRAN.lihatBerkas} ${k.nama}`}
                      title={LAMPIRAN.lihatBerkas}
                      onClick={() => {
                        setGalatPanel('')
                        setGantiKategori(false)
                        setKategoriBaru({})
                        setBerkasDilihat(k.kode)
                      }}
                    >
                      <svg viewBox="0 0 24 24" aria-hidden="true" focusable="false">
                        <rect x="1" y="1" width="22" height="22" rx="3" fill="currentColor" />
                        <path d="M12 15 7.5 10.5h3V6h3v4.5h3L12 15Zm-5 1.5v2h10v-2H7Z" fill="#fff" />
                      </svg>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {unggahKe !== null && (
          <Modal
            judul={LAMPIRAN.judulUnggah}
            onTutup={() => {
              setUnggahKe(null)
            }}
            labelBatal={LAMPIRAN.batal}
            aksi={
              <button type="button" className="btn btn--primary" disabled={sibuk} onClick={lampirkan}>
                {LAMPIRAN.lampirkan}
              </button>
            }
          >
            <p>
              <strong>{unggahKe.nama}</strong>
            </p>
            {galatUnggah !== null && <Gagal galat={galatUnggah} />}
            <label
              className={'tria__unggah' + (seret ? ' tria__unggah--seret' : '')}
              onDragOver={(e) => {
                e.preventDefault()
                setSeret(true)
              }}
              onDragLeave={() => {
                setSeret(false)
              }}
              onDrop={(e) => {
                e.preventDefault()
                setSeret(false)
                const jatuh = Array.from(e.dataTransfer.files)
                setDipilih((b) => gabungBerkas(b, jatuh))
              }}
            >
              <span>{LAMPIRAN.seretBerkas}</span>
              <input
                type="file"
                multiple
                aria-label={LAMPIRAN.pilihBerkas}
                disabled={sibuk}
                onChange={(e) => {
                  const baru = Array.from(e.target.files ?? [])
                  setDipilih((b) => gabungBerkas(b, baru))
                  // Dikosongkan: memilih berkas yang sama lagi tetap memicu onChange.
                  e.target.value = ''
                }}
              />
            </label>
            {dipilih.length > 0 && (
              <ul className="tria__unggah-daftar">
                {dipilih.map((f) => (
                  <li key={f.name}>
                    <span>{f.name}</span>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      disabled={sibuk}
                      onClick={() => {
                        setDipilih((b) => b.filter((x) => x !== f))
                      }}
                    >
                      {LAMPIRAN.buangPilihan}
                    </button>
                  </li>
                ))}
              </ul>
            )}
            {proses !== null && (
              <p className="muted">
                {LAMPIRAN.mengunggah} {proses}
              </p>
            )}
            {gagalUnggah.length > 0 && (
              <ul className="tria__unggah-gagal">
                {gagalUnggah.map((g) => (
                  <li key={g.berkas.name}>
                    {g.berkas.name}: {pesanGalat(g.galat) ?? (g.galat instanceof Error ? g.galat.message : LAMPIRAN.gagal)}
                  </li>
                ))}
              </ul>
            )}
          </Modal>
        )}

        {berkasDilihat !== null && (
          <ModalLihatBerkas
            kategori={kategori}
            berkas={berkasKategori(berkas, kategoriDilihat)}
            idKontrak={idKontrak}
            bolehTulis={idTulis !== ''}
            bisaUnggah={bisaUnggah}
            bolehGantiKategori={bolehGantiKategori}
            gantiKategori={gantiKategori}
            kategoriBaru={kategoriBaru}
            sibuk={sibuk}
            galatPanel={galatPanel}
            onTutup={() => {
              setBerkasDilihat(null)
            }}
            onGantiKategori={() => {
              setGantiKategori(true)
            }}
            onSimpanKategori={simpanKategori}
            onPilihKategori={(id, kode) => {
              setKategoriBaru((x) => ({ ...x, [id]: kode }))
            }}
            onUnduh={unduh}
            onOffice={bukaOffice}
            onHapus={hapus}
          />
        )}

        {/* `View Office Online` — penampil kantor di bingkai popup. */}
        {penampil !== null && (
          <Modal judul={penampil.nama} onTutup={tutupPenampil} labelBatal={LAMPIRAN.tutup} penuh>
            {!penampil.siap && <p className="tria__redup">{LAMPIRAN.memuatPenampil}</p>}
            <iframe ref={bingkai} title={penampil.nama} className="tria__penampil" />
          </Modal>
        )}
      </div>
    </Panel>
  )
}

/**
 * Modal `View File` — `ShowAttachmentTreaty`: DUA kolom (`File Name` ·
 * `Type`). Nama berkas = tautan unduh; sel Type memuat jenis, `View Office
 * Online`, dropdown Change Category, dan `Delete`.
 *
 * Dipisah dari panel supaya isinya dapat dirender statis oleh uji; keadaan
 * modalnya tetap milik panel.
 */
export function ModalLihatBerkas({
  kategori,
  berkas,
  idKontrak,
  bolehTulis,
  bisaUnggah,
  bolehGantiKategori,
  gantiKategori,
  kategoriBaru,
  sibuk,
  galatPanel,
  onTutup,
  onGantiKategori,
  onSimpanKategori,
  onPilihKategori,
  onUnduh,
  onOffice,
  onHapus,
}: {
  /** Seluruh kategori — pilihan dropdown Change Category. */
  kategori: readonly BarisKategoriLampiran[]
  /** Berkas milik kategori yang dilihat. */
  berkas: readonly BarisLampiranWarisan[]
  idKontrak: string
  /** Pengenal TERSIMPAN ada — Delete/Change Category boleh. */
  bolehTulis: boolean
  /** `TreatyIn.ViewState !='1' || TreatyIn.RevisionState='1'`. */
  bisaUnggah: boolean
  /** `StatusAkseptasi` bukan Resolve Complete / Decline. */
  bolehGantiKategori: boolean
  /** `StatusDoc.CARI30 = 1`. */
  gantiKategori: boolean
  kategoriBaru: Readonly<Record<string, string>>
  sibuk: boolean
  galatPanel: string
  onTutup: () => void
  onGantiKategori: () => void
  onSimpanKategori: () => void
  onPilihKategori: (idLampiran: string, kode: string) => void
  onUnduh: (b: BarisLampiranWarisan) => void
  onOffice: (b: BarisLampiranWarisan) => void
  onHapus: (idLampiran: string) => void
}) {
  return (
    <Modal judul={LAMPIRAN.judulLihatBerkas} onTutup={onTutup} labelBatal={LAMPIRAN.tutup} lebar>
      {/* `Change Category` / `Save` — wadah bersyarat `StatusAkseptasi`
          bukan Resolve Complete / Decline; tombolnya bergantian menurut
          `StatusDoc.CARI30`. */}
      {bolehGantiKategori && bolehTulis && (
        <div className="tria__aksi">
          {gantiKategori ? (
            <button type="button" className="btn btn--primary btn--sm" disabled={sibuk} onClick={onSimpanKategori}>
              {LAMPIRAN.simpanKategori}
            </button>
          ) : (
            <button type="button" className="btn btn--sm" disabled={sibuk} onClick={onGantiKategori}>
              {LAMPIRAN.gantiKategori}
            </button>
          )}
        </div>
      )}
      {galatPanel !== '' && (
        <div className="alert alert--error" role="alert">
          {galatPanel}
        </div>
      )}
      <div className="table-wrap">
        <table className="tria__tabel tria__lihat-berkas">
          <colgroup>
            <col className="tria__lihat-nama" />
            <col className="tria__lihat-tipe" />
          </colgroup>
          <thead>
            <tr>
              {KOLOM_LIHAT_BERKAS.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {berkas.length === 0 && (
              <tr>
                <td colSpan={KOLOM_LIHAT_BERKAS.length}>{LAMPIRAN.tanpaBaris}</td>
              </tr>
            )}
            {berkas.map((b) => (
              <tr key={b.id}>
                {/* Nama berkas = tautan `DownloadAttachmentTreaty`. */}
                <td>
                  <button
                    type="button"
                    className="tria__tautan"
                    disabled={idKontrak === '' || sibuk}
                    onClick={() => {
                      onUnduh(b)
                    }}
                  >
                    {b.namaBerkas}
                  </button>
                </td>
                <td>
                  <div className="tria__lihat-jenis">
                    <span>{b.jenisMime}</span>
                    <div className="tria__lihat-aksi">
                      {/* `View Office Online` — xls/xlsx/doc/docx/ppt/pptx. */}
                      {JENIS_OFFICE.includes(b.jenisMime.toLowerCase()) && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={idKontrak === ''}
                          onClick={() => {
                            onOffice(b)
                          }}
                        >
                          {LAMPIRAN.viewOffice}
                        </button>
                      )}
                      {/* Mode Change Category: kategori menjadi dropdown. */}
                      {gantiKategori && (
                        <select
                          className="field__input"
                          aria-label={LAMPIRAN.gantiKategori}
                          value={kategoriBaru[b.id] ?? b.kodeKategori}
                          onChange={(e) => {
                            onPilihKategori(b.id, e.target.value)
                          }}
                        >
                          {kategori
                            .filter((k) => k.kode !== '')
                            .map((k) => (
                              <option key={k.kode} value={k.kode}>
                                {k.nama}
                              </option>
                            ))}
                        </select>
                      )}
                      {/* `Delete` — `ViewState !='1' || RevisionState='1'`. */}
                      {bisaUnggah && bolehTulis && (
                        <button
                          type="button"
                          className="btn btn--sm tl-hapus"
                          disabled={sibuk}
                          onClick={() => {
                            onHapus(b.id)
                          }}
                        >
                          {LAMPIRAN.hapus}
                        </button>
                      )}
                    </div>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Modal>
  )
}

/** Syarat `View Office Online` — `ShowAttachmentTreaty`. */
const JENIS_OFFICE: readonly string[] = ['xls', 'xlsx', 'doc', 'docx', 'ppt', 'pptx']

/**
 * Berkas milik SATU kategori — disaring menurut `kode` bila ada, NAMA bila
 * tidak (sama dengan Treaty In).
 */
export function berkasKategori(
  berkas: readonly BarisLampiranWarisan[],
  kategori: BarisKategoriLampiran | null,
): readonly BarisLampiranWarisan[] {
  if (kategori === null) return []
  if (kategori.kode !== '') return berkas.filter((b) => b.kodeKategori === kategori.kode)
  return berkas.filter((b) => b.namaKategori === kategori.nama)
}
