import axios from 'axios'

export interface JournalEntry {
  id: number
  user_id: number
  user_nip: string
  user_name: string
  activity_date: string
  activity_time: string
  title: string
  description: string
  evidence_name: string
  evidence_mime: string
  has_evidence: boolean
  status: 'PENDING' | 'APPROVED' | 'REJECTED'
  reviewer_note: string
  reviewed_by: string
  reviewed_at: string
  created_at: string
}

export const journalStatusMeta: Record<string, { label: string; badge: string; icon: string }> = {
  PENDING: { label: 'Menunggu Review', badge: 'bg-amber-50 text-amber-700 border-amber-200', icon: 'hourglass_top' },
  APPROVED: { label: 'Disetujui', badge: 'bg-emerald-50 text-emerald-700 border-emerald-200', icon: 'check_circle' },
  REJECTED: { label: 'Tidak Disetujui', badge: 'bg-rose-50 text-rose-700 border-rose-200', icon: 'cancel' }
}

// Evidence is behind the auth header, so <img src> cannot load it directly; fetch it as a blob instead.
const evidenceCache = new Map<number, string>()

export async function loadEvidenceUrl(entryId: number): Promise<string> {
  const cached = evidenceCache.get(entryId)
  if (cached) return cached
  const res = await axios.get(`/api/journal/${entryId}/evidence`, { responseType: 'blob' })
  const url = URL.createObjectURL(res.data)
  evidenceCache.set(entryId, url)
  return url
}

export async function openEvidence(entryId: number) {
  // Open the tab synchronously so popup blockers allow it, then point it at the blob once loaded
  const win = window.open('', '_blank')
  try {
    const url = await loadEvidenceUrl(entryId)
    if (win) win.location.href = url
    else window.location.href = url
  } catch (e) {
    if (win) win.close()
    throw e
  }
}

export function todayWita(): string {
  // en-CA formats as YYYY-MM-DD
  return new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Makassar' }).format(new Date())
}

export function formatDateId(iso: string): string {
  if (!iso) return '-'
  const d = new Date(`${iso}T00:00:00`)
  return d.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}
