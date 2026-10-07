// Layar satu kasus Claim Prop - FlowAction OutstandingClaim (Assignment2) atau InputAcceptation (Assignment1). Isinya
// pohon tata server; aksi dikirim bersama nilai semua medan terbuka. Local action / harness (PrintFile, GeneratePLA,
// GenerateDLATreaty, PreventRejectClaimProp, CommitteeTreaty) dibuka sebagai Modal; pop-up pemilih sebagai `Popup`.

import { useCallback, useEffect, useMemo, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ApiFailure } from '../../../../inti/frontend/klien'
import {
  aksiKasus,
  ambilAcuan,
  bukaKasus,
  pilihanKasus,
  type AcuanStatis,
  type Halaman,
  type Layar,
  type Pilihan,
} from '../api'
import { CP } from '../labels'
import { ambil, masukan, semuaTata, setel } from '../nilai'
import Popup, { type JenisPopup } from './Popup'
import TataView, { type KonteksTata } from './TataView'

/** Aksi yang hanya membuka pop-up harness. */
const POPUP: Record<string, JenisPopup> = {
  PilihMaster: 'master',
  PilihPolis: 'polis',
  PilihSebab: 'sebab',
  BukaKatastrofe: 'katastrofe',
  RingkasanOS: 'ringkasanOS',
}

/** Aksi submit local action: modal ditutup bila berhasil. */
const SUBMIT_MODAL = new Set(['TryMakePLA', 'PrintDLATreatyIn', 'CloseClaimProp', 'AddKomiteTreatyChild'])

/** Aksi yang parameternya nilai pilihan (bukan jalur medan). */
const PARAM_DARI_NILAI = new Set(['PilihRekening', 'PilihAllocation', 'DisableEditRNMShare'])

const JUDUL_MODAL = (m: string): string => {
  if (m === 'pla') return CP.popPLA
  if (m === 'tutupKlaim') return CP.popTutup
  if (m.startsWith('dla:')) return CP.popDLA
  if (m.startsWith('komite:')) return CP.popKomite
  return m
}

interface RekeningBank {
  ClientName: string
  NameOfBank: string
  BranchOfBank: string
  AccountNo: string
}

export default function LayarKasus({ id, pelaku, onKembali }: { id: string; pelaku: string; onKembali: () => void }) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [h, setH] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [info, setInfo] = useState<string | null>(null)
  const [modal, setModal] = useState<string | null>(null)
  const [popup, setPopup] = useState<JenisPopup | null>(null)
  const [acuan, setAcuan] = useState<AcuanStatis | null>(null)
  const [opsiKasus, setOpsiKasus] = useState<Record<string, Pilihan[]>>({})
  const [saranPeta, setSaranPeta] = useState<Record<string, Pilihan[]>>({})

  const terima = useCallback((l: Layar) => {
    setLayar(l)
    setH(l.halaman)
    setGalat(null)
  }, [])

  useEffect(() => {
    bukaKasus(id).then(terima, (g: unknown) => setGalat(g))
    ambilAcuan().then(setAcuan, (g: unknown) => setGalat(g))
  }, [id, terima])

  const idMaster = h?.nilai['ClaimData.IDMaster'] ?? ''
  useEffect(() => {
    if (idMaster === '') return
    for (const s of ['limits', 'spreading', 'shareRNM']) {
      pilihanKasus<Pilihan[]>(id, s).then(
        (d) => setOpsiKasus((o) => ({ ...o, [s]: d ?? [] })),
        () => undefined,
      )
    }
  }, [id, idMaster])

  const kirim = useCallback(
    (aksi: string, indeks = 0, ubahan: Record<string, string> = {}, param = '') => {
      if (!layar || !h) return
      let h2 = h
      for (const [j, v] of Object.entries(ubahan)) h2 = setel(h2, j, v)
      setH(h2)
      const prm = param !== '' ? param : PARAM_DARI_NILAI.has(aksi) ? (Object.values(ubahan)[0] ?? '') : ''
      const semua = semuaTata(layar.tata, [layar.adjustment ?? {}, layar.modal ?? {}])
      setSibuk(true)
      setInfo(null)
      aksiKasus(id, { aksi, indeks, param: prm, tahap: layar.kasus.tahap, masukan: masukan(h2, semua) }).then(
        (l) => {
          terima(l)
          setSibuk(false)
          if (l.info) setInfo(l.info)
          if (SUBMIT_MODAL.has(aksi)) setModal(null)
          if (aksi === 'BukaKomite') {
            const ok = l.halaman.nilai['Protect.CARI1'] === '1' && l.halaman.nilai['Protect.CARI2'] === '1'
            setModal(ok ? `komite:${indeks}` : null)
          }
          setPopup(null)
        },
        (g: unknown) => {
          setSibuk(false)
          setGalat(g)
        },
      )
    },
    [id, layar, h, terima],
  )

  const aksi = useCallback(
    (nama: string, indeks = 0, ubahan: Record<string, string> = {}) => {
      if (POPUP[nama]) {
        setPopup(POPUP[nama] ?? null)
        return
      }
      if (nama === 'BukaPLA') return setModal('pla')
      if (nama === 'BukaTutupKlaim') return setModal('tutupKlaim')
      if (nama === 'BukaDLA') return setModal(`dla:${indeks}`)
      if (nama === 'TutupModal') return setModal(null)
      kirim(nama, indeks, ubahan)
    },
    [kirim],
  )

  const saran = useCallback(
    (sumber: string, indeks: number, cari: string) => {
      if (!['adjuster', 'provinsi', 'rekening'].includes(sumber)) return
      pilihanKasus<unknown[]>(id, sumber, indeks, cari).then(
        (d) => {
          const op: Pilihan[] =
            sumber === 'rekening'
              ? ((d ?? []) as RekeningBank[]).map((r) => ({
                  nilai: `${r.AccountNo}|${r.NameOfBank}`,
                  label: `${r.NameOfBank} - ${r.BranchOfBank} - ${r.AccountNo} (${r.ClientName})`,
                }))
              : ((d ?? []) as Pilihan[])
          setSaranPeta((s) => ({ ...s, [`${sumber}|${indeks}`]: op }))
        },
        () => undefined,
      )
    },
    [id],
  )

  const opsi = useCallback(
    (sumber: string, indeks: number): Pilihan[] => {
      if (!h) return []
      if (sumber === 'mataUang') return acuan?.mataUang ?? []
      if (sumber === 'jenisReas') return acuan?.jenisReas ?? []
      if (sumber === 'jenisReas4') return acuan?.jenisReas4 ?? []
      if (sumber.startsWith('kode:')) return (acuan?.kode[sumber.slice(5)] ?? []).map((k) => ({ nilai: k, label: k }))
      if (sumber === 'mataUangAdj') {
        return (h.daftar[`ClaimData.AdjustmentList(${indeks}).CurencyAdjustment`] ?? []).map((b) => ({
          nilai: b.CurrencyID ?? '',
          label: b.Currency ?? '',
        }))
      }
      if (sumber === 'allocation') {
        return (h.daftar[`ClaimData.AdjustmentList(${indeks}).LossAllocation`] ?? []).map((b) => ({
          nilai: b.TreatyName ?? '',
          label: `${b.TreatyName ?? ''} (${b.SharePercentage ?? ''})`,
        }))
      }
      if (opsiKasus[sumber]) return opsiKasus[sumber] ?? []
      return saranPeta[`${sumber}|${indeks}`] ?? []
    },
    [h, acuan, opsiKasus, saranPeta],
  )

  const k: KonteksTata | null = useMemo(() => {
    if (!layar || !h) return null
    const dasar: KonteksTata = {
      h,
      ubah: (j, v) => setH((x) => (x ? setel(x, j, v) : x)),
      aksi,
      opsi,
      saran,
      pesanMedan: layar.pesanMedan ?? {},
      sibuk,
    }
    return {
      ...dasar,
      rincian: (_j: string, n: number) => (
        <div className="claimprop__rinci">
          <div className="claimprop__label">{CP.popAdjustment}</div>
          <TataView tata={layar.adjustment?.[String(n)] ?? []} k={dasar} />
        </div>
      ),
    }
  }, [layar, h, aksi, opsi, saran, sibuk])

  if (!layar || !h || !k) {
    return (
      <section className="inbox claimprop__akar">
        {galat ? <Gagal galat={galat} /> : <Memuat pesan={CP.memuat} />}
      </section>
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? galat.detail.message : null
  return (
    <section className="inbox claimprop__akar">
      <header className="inbox__kepala">
        <button type="button" className="btn btn--sm" onClick={onKembali}>
          {CP.kembali}
        </button>
        <h2 className="inbox__judul">
          {layar.kasus.id} - {layar.label || layar.kasus.statusWork}
        </h2>
      </header>
      {!layar.bolehKerja && <div className="alert">{CP.hanyaLihat}</div>}
      {galat !== null &&
        (pesanGalat ? <div className="alert alert--error">{pesanGalat}</div> : <Gagal galat={galat} />)}
      {(layar.pesan ?? []).length > 0 && (
        <div className="alert alert--error">
          <ul>
            {(layar.pesan ?? []).map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        </div>
      )}
      {info && <div className="alert">{info}</div>}
      <div className="claimprop__layar">
        <TataView tata={layar.tata} k={k} />
      </div>
      {modal && (
        <Modal judul={JUDUL_MODAL(modal)} onTutup={() => setModal(null)} lebar>
          <TataView tata={layar.modal?.[modal] ?? []} k={{ ...k, rincian: undefined }} />
        </Modal>
      )}
      {popup && (
        <Popup
          id={id}
          jenis={popup}
          pelaku={pelaku}
          sts={{ sts: ambil(h, 'ClaimData.StsKatastrofe'), non: ambil(h, 'ClaimData.NonKatastrofeType') }}
          onPilih={(a, p) => kirim(a, 0, {}, p)}
          onTutup={() => setPopup(null)}
        />
      )}
    </section>
  )
}
